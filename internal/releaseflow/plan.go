package releaseflow

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Repository struct{ Dir string }

type Options struct {
	Kind, Custom, LastNightly string
	Date                      time.Time
}

type Plan struct {
	Version        string `json:"version"`
	PackageVersion string `json:"package_version"`
	Tag            string `json:"tag"`
	SHA            string `json:"sha"`
	PreviousTag    string `json:"previous_tag"`
	Skip           bool   `json:"skip"`
	Prerelease     bool   `json:"prerelease"`
}

var nightlyPattern = regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+-nightly\.[0-9]{14}\.[0-9a-f]{12}$`)

func (r Repository) Plan(options Options) (Plan, string, error) {
	if options.Kind != "custom" && options.Custom != "" {
		return Plan{}, "", fmt.Errorf("custom version is only accepted for kind custom")
	}
	if options.Kind != "nightly" && options.LastNightly != "" {
		return Plan{}, "", fmt.Errorf("last-nightly is only accepted for kind nightly")
	}
	sha, err := r.resolveCommit("HEAD")
	if err != nil {
		return Plan{}, "", err
	}
	base, previous, err := r.latestStable(sha)
	if err != nil {
		return Plan{}, "", err
	}
	if options.Kind == "nightly" && options.LastNightly != "" {
		if !nightlyPattern.MatchString(options.LastNightly) {
			return Plan{}, "", fmt.Errorf("invalid last-nightly tag %q", options.LastNightly)
		}
		if _, err := ParseVersion(strings.TrimPrefix(options.LastNightly, "v")); err != nil {
			return Plan{}, "", err
		}
		lastSHA, err := r.resolveCommit("refs/tags/" + options.LastNightly)
		if err != nil {
			return Plan{}, "", err
		}
		parts := strings.Split(options.LastNightly, ".")
		if _, err := time.Parse("20060102150405", parts[len(parts)-2]); err != nil || parts[len(parts)-1] != lastSHA[:12] {
			return Plan{}, "", fmt.Errorf("last-nightly tag timestamp or source hash is invalid: %q", options.LastNightly)
		}
		if err := r.ancestor(lastSHA, sha); err != nil {
			return Plan{}, "", fmt.Errorf("last nightly is not an ancestor of HEAD: %w", err)
		}
		if lastSHA == sha {
			return Plan{SHA: sha, PreviousTag: options.LastNightly, Skip: true, Prerelease: true}, "", nil
		}
		previous = options.LastNightly
	}
	var next Version
	if options.Kind == "custom" {
		next, err = ParseVersion(options.Custom)
		if err == nil && (next.Prerelease == "nightly" || strings.HasPrefix(next.Prerelease, "nightly.")) {
			err = fmt.Errorf("the nightly prerelease namespace is reserved for scheduled builds")
		}
		if err == nil && next.Compare(base) <= 0 {
			err = fmt.Errorf("custom version must be greater than latest stable %s", base.String())
		}
	} else {
		next, err = bump(base, options.Kind)
	}
	if err != nil {
		return Plan{}, "", err
	}
	if options.Kind == "nightly" {
		if options.Date.IsZero() {
			return Plan{}, "", fmt.Errorf("nightly requires a release timestamp")
		}
		next.Prerelease = "nightly." + options.Date.UTC().Format("20060102150405") + "." + sha[:12]
	}
	plan := Plan{Version: next.String(), PackageVersion: next.Core(), Tag: "v" + next.String(), SHA: sha, PreviousTag: previous, Prerelease: next.Prerelease != ""}
	if _, err := r.git("show-ref", "--verify", "--quiet", "refs/tags/"+plan.Tag); err == nil {
		return Plan{}, "", fmt.Errorf("tag %s already exists", plan.Tag)
	} else if exit, ok := err.(*exec.ExitError); !ok || exit.ExitCode() != 1 {
		return Plan{}, "", fmt.Errorf("check existing tag: %w", err)
	}
	baseSHA := ""
	if previous != "" {
		baseSHA, err = r.resolveCommit("refs/tags/" + previous)
		if err != nil {
			return Plan{}, "", err
		}
	}
	commits, err := r.commits(baseSHA, sha)
	if err != nil {
		return Plan{}, "", err
	}
	return plan, releaseNotes(plan, commits), nil
}

func (r Repository) latestStable(head string) (Version, string, error) {
	data, err := r.git("tag", "--merged", head, "--list", "v*")
	if err != nil {
		return Version{}, "", fmt.Errorf("list tags: %w", err)
	}
	var best Version
	tag := ""
	for _, candidate := range strings.Fields(string(data)) {
		v, err := ParseVersion(strings.TrimPrefix(candidate, "v"))
		if err != nil {
			if versionPattern.MatchString(strings.TrimPrefix(candidate, "v")) {
				return Version{}, "", fmt.Errorf("unsupported existing version tag %s: %w", candidate, err)
			}
			continue
		}
		if v.Prerelease != "" {
			continue
		}
		if tag == "" || v.Compare(best) > 0 {
			best, tag = v, candidate
		}
	}
	return best, tag, nil
}

func (r Repository) resolveCommit(ref string) (string, error) {
	if ref == "" || strings.HasPrefix(ref, "-") {
		return "", fmt.Errorf("invalid commit reference %q", ref)
	}
	data, err := r.git("rev-parse", "--verify", "--end-of-options", ref+"^{commit}")
	if err != nil {
		return "", fmt.Errorf("resolve commit %q: %w", ref, err)
	}
	sha := strings.TrimSpace(string(data))
	if len(sha) != 40 && len(sha) != 64 {
		return "", fmt.Errorf("unexpected Git object ID %q", sha)
	}
	return sha, nil
}

func (r Repository) ancestor(base, head string) error {
	_, err := r.git("merge-base", "--is-ancestor", base, head)
	return err
}

func (r Repository) commits(base, head string) ([]Commit, error) {
	rangeArg := head
	if base != "" {
		rangeArg = base + ".." + head
	}
	data, err := r.git("log", "--no-merges", "--reverse", "--format=%H%x00%s%x00%b%x00", rangeArg, "--")
	if err != nil {
		return nil, fmt.Errorf("read commits: %w", err)
	}
	var commits []Commit
	fields := strings.Split(string(data), "\x00")
	for i := 0; i+2 < len(fields); i += 3 {
		commits = append(commits, Commit{strings.TrimSpace(fields[i]), fields[i+1], fields[i+2]})
	}
	return commits, nil
}

func (r Repository) git(args ...string) ([]byte, error) {
	command := exec.Command("git", args...)
	command.Dir = r.Dir
	return command.Output()
}

// Write stores a complete planning bundle; it refuses to overwrite an existing one.
func Write(dir string, plan Plan, notes string) error {
	if err := os.Mkdir(dir, 0o755); err != nil {
		return fmt.Errorf("create plan directory: %w", err)
	}
	data, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		return err
	}
	outputs := fmt.Sprintf("version=%s\npackage_version=%s\ntag=%s\nsha=%s\nprevious_tag=%s\nskip=%t\nprerelease=%t\n", plan.Version, plan.PackageVersion, plan.Tag, plan.SHA, plan.PreviousTag, plan.Skip, plan.Prerelease)
	for name, contents := range map[string]string{"plan.json": string(data) + "\n", "notes.md": notes, "github-output": outputs} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(contents), 0o644); err != nil {
			return err
		}
	}
	return nil
}
