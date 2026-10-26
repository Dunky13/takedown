package releaseflow

import (
	"fmt"
	"regexp"
	"strings"
)

var titlePattern = regexp.MustCompile(`^(build|chore|ci|docs|feat|fix|perf|refactor|revert|style|test)(\([^()\s]+\))?(!)?: ([^\r\n]+)$`)
var breakingPattern = regexp.MustCompile(`(?m)^BREAKING(?: CHANGE|-CHANGE):[ \t]+\S`)
var footerPattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9-]*(?::[ \t]+|[ \t]+#)\S`)

type Commit struct {
	SHA, Subject, Body string
}

func CheckTitle(title string) error {
	matches := titlePattern.FindStringSubmatch(title)
	if matches == nil || strings.TrimSpace(matches[4]) == "" || strings.TrimSpace(matches[4]) != matches[4] {
		return fmt.Errorf("expected type(scope)!: description using build, chore, ci, docs, feat, fix, perf, refactor, revert, style, or test; got %q", title)
	}
	return nil
}

func (r Repository) CheckCommits(base, head string) error {
	baseSHA, err := r.resolveCommit(base)
	if err != nil {
		return err
	}
	headSHA, err := r.resolveCommit(head)
	if err != nil {
		return err
	}
	commits, err := r.commits(baseSHA, headSHA)
	if err != nil {
		return err
	}
	var failures []string
	for _, commit := range commits {
		if err := CheckTitle(commit.Subject); err != nil {
			failures = append(failures, commit.SHA+": "+err.Error())
		}
	}
	if len(failures) > 0 {
		return fmt.Errorf("nonconventional commits:\n%s", strings.Join(failures, "\n"))
	}
	return nil
}

func releaseNotes(plan Plan, commits []Commit) string {
	var out strings.Builder
	fmt.Fprintf(&out, "# %s\n\nSource: `%s`\n", plan.Tag, plan.SHA)
	if plan.PreviousTag != "" {
		fmt.Fprintf(&out, "\nChanges since `%s`.\n", plan.PreviousTag)
	}
	groups := map[string][]string{}
	for _, commit := range commits {
		group := "Other changes"
		matches := titlePattern.FindStringSubmatch(commit.Subject)
		if matches != nil {
			switch matches[1] {
			case "feat":
				group = "Features"
			case "fix":
				group = "Bug fixes"
			case "perf":
				group = "Performance"
			case "docs":
				group = "Documentation"
			case "build", "ci":
				group = "Build and CI"
			}
		}
		line := "- " + escapeMarkdown(commit.Subject) + " (`" + commit.SHA[:12] + "`)"
		groups[group] = append(groups[group], line)
		if (matches != nil && matches[3] == "!") || breakingPattern.MatchString(commit.Body) {
			breaking := line
			for _, bodyLine := range breakingDescription(commit.Body) {
				breaking += "\n  " + escapeMarkdown(bodyLine)
			}
			groups["Breaking changes"] = append(groups["Breaking changes"], breaking)
		}
	}
	for _, name := range []string{"Breaking changes", "Features", "Bug fixes", "Performance", "Documentation", "Build and CI", "Other changes"} {
		if len(groups[name]) > 0 {
			fmt.Fprintf(&out, "\n## %s\n\n%s\n", name, strings.Join(groups[name], "\n"))
		}
	}
	if len(commits) == 0 {
		out.WriteString("\nNo source changes in this range.\n")
	}
	return out.String()
}

// Breaking descriptions can span paragraphs. Stop at the next conventional
// footer so DCO sign-offs and issue references do not become migration guidance.
func breakingDescription(body string) []string {
	var description []string
	collecting := false
	for _, line := range strings.Split(body, "\n") {
		if breakingPattern.MatchString(line) {
			collecting = true
		} else if footerPattern.MatchString(line) {
			collecting = false
		}
		if collecting {
			description = append(description, line)
		}
	}
	for len(description) > 0 && strings.TrimSpace(description[len(description)-1]) == "" {
		description = description[:len(description)-1]
	}
	return description
}

func escapeMarkdown(value string) string {
	return strings.NewReplacer("\\", "\\\\", "`", "\\`", "*", "\\*", "_", "\\_", "[", "\\[", "]", "\\]", "<", "&lt;", ">", "&gt;", "&", "&amp;").Replace(value)
}
