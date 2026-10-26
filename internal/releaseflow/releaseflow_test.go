package releaseflow

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSemVerAndInstallerLimits(t *testing.T) {
	for _, valid := range []string{"0.0.0", "1.2.3", "65535.65535.65535", "1.2.3-rc.1", "1.2.3-0.alpha+build.012", "1.2.3-alpha-1"} {
		v, err := ParseVersion(valid)
		if err != nil || v.String() != valid {
			t.Errorf("ParseVersion(%q) = %s, %v", valid, v, err)
		}
	}
	for _, invalid := range []string{"", "v1.2.3", "01.2.3", "1.2", "1.2.3-01", "1.2.3-rc..1", "1.2.3+", "65536.0.0", "1.2.3\nmalicious=value", "1.2.3-rc;touch"} {
		if _, err := ParseVersion(invalid); err == nil {
			t.Errorf("accepted invalid version %q", invalid)
		}
	}
	ordered := []string{"1.0.0-alpha", "1.0.0-alpha.1", "1.0.0-alpha.beta", "1.0.0-beta", "1.0.0-beta.2", "1.0.0-beta.11", "1.0.0-rc.1", "1.0.0", "1.0.1", "1.1.0", "2.0.0"}
	for i := 1; i < len(ordered); i++ {
		a, _ := ParseVersion(ordered[i-1])
		b, _ := ParseVersion(ordered[i])
		if a.Compare(b) >= 0 || b.Compare(a) <= 0 {
			t.Errorf("wrong precedence %s / %s", a, b)
		}
	}
	a, _ := ParseVersion("1.0.0+one")
	b, _ := ParseVersion("1.0.0+two")
	if a.Compare(b) != 0 {
		t.Error("build metadata changed version precedence")
	}
}

func TestConventionalTitles(t *testing.T) {
	for _, value := range []string{"feat: add release flow", "fix(api)!: change response", "build(deps): update Go", "ci: publish nightly", "revert: undo change"} {
		if err := CheckTitle(value); err != nil {
			t.Errorf("rejected %q: %v", value, err)
		}
	}
	for _, value := range []string{"update flow", "Feat: title", "feature: title", "fix:", "fix: ", "fix:  padded", "fix: trailing ", "fix(): empty scope", "fix( ): blank scope", "fix(a(b)): nested", "fix: first\nsecond", "Merge pull request #1"} {
		if err := CheckTitle(value); err == nil {
			t.Errorf("accepted %q", value)
		}
	}
}

func TestReleasePlanBumpsAndChangelog(t *testing.T) {
	r := fixture(t)
	old := commit(t, r, "initial historical import", "")
	git(t, r, "tag", "v1.9.0", old)
	commit(t, r, "feat(api)!: new [API] <script>", "BREAKING CHANGE: changed response shape")
	commit(t, r, "fix: repair download", "")
	commit(t, r, "docs: update guide", "")
	for kind, version := range map[string]string{"patch": "1.9.1", "minor": "1.10.0", "major": "2.0.0"} {
		t.Run(kind, func(t *testing.T) {
			plan, notes, err := r.Plan(Options{Kind: kind})
			if err != nil {
				t.Fatal(err)
			}
			if plan.Version != version || plan.PackageVersion != version || plan.Tag != "v"+version || plan.PreviousTag != "v1.9.0" || plan.Prerelease || plan.Skip {
				t.Fatalf("wrong plan: %+v", plan)
			}
			for _, expected := range []string{"## Breaking changes", "## Features", "## Bug fixes", "## Documentation", "changed response shape", `\[API\] &lt;script&gt;`} {
				if !strings.Contains(notes, expected) {
					t.Errorf("notes omit %q:\n%s", expected, notes)
				}
			}
			if strings.Contains(notes, "historical import") {
				t.Error("included commit outside release range")
			}
		})
	}
	plan, _, err := r.Plan(Options{Kind: "custom", Custom: "2.0.0-rc.1+ci.10"})
	if err != nil || plan.Version != "2.0.0-rc.1+ci.10" || plan.PackageVersion != "2.0.0" || !plan.Prerelease {
		t.Fatalf("prerelease: %+v, %v", plan, err)
	}
}

func TestPlanUsesHighestReachableStableVersion(t *testing.T) {
	r := fixture(t)
	first := commit(t, r, "feat: initial", "")
	git(t, r, "tag", "v2.9.0", first)
	git(t, r, "tag", "v2.10.0", first)
	git(t, r, "tag", "v90.0.0-rc.1", first)
	git(t, r, "checkout", "-b", "other")
	unrelated := commit(t, r, "feat: other branch", "")
	git(t, r, "tag", "v99.0.0", unrelated)
	git(t, r, "checkout", "main")
	commit(t, r, "fix: main change", "")
	plan, _, err := r.Plan(Options{Kind: "patch"})
	if err != nil || plan.Version != "2.10.1" {
		t.Fatalf("selected wrong base: %+v %v", plan, err)
	}
}

func TestMultilineBreakingMigrationNotes(t *testing.T) {
	r := fixture(t)
	commit(t, r, "feat(api)!: rename response fields", "BREAKING CHANGE: rename fields\n\nMigrate clients from old_name to new_name.\nKeep both during rollout.\n\nSigned-off-by: Example <example@example.invalid>\nRefs #123")
	_, notes, err := r.Plan(Options{Kind: "major"})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"BREAKING CHANGE: rename fields", `Migrate clients from old\_name to new\_name.`, "Keep both during rollout."} {
		if !strings.Contains(notes, want) {
			t.Fatalf("missing migration detail %q:\n%s", want, notes)
		}
	}
	if strings.Contains(notes, "Signed-off-by:") || strings.Contains(notes, "Refs #123") {
		t.Fatalf("unrelated footers became migration guidance:\n%s", notes)
	}
}

func TestCustomAndCollisionFailures(t *testing.T) {
	r := fixture(t)
	commit(t, r, "feat: initial", "")
	git(t, r, "tag", "v1.2.3")
	for _, options := range []Options{
		{Kind: "wrong"},
		{Kind: "custom"},
		{Kind: "custom", Custom: "1.2.3"},
		{Kind: "custom", Custom: "1.2.3-rc.1"},
		{Kind: "custom", Custom: "1.0.0"},
		{Kind: "custom", Custom: "1.3.0-nightly.20261026010000.123456789abc"},
		{Kind: "patch", Custom: "2.0.0"},
		{Kind: "patch", LastNightly: "v1.0.0"},
		{Kind: "nightly", LastNightly: "--malicious"},
		{Kind: "nightly", LastNightly: "v1.2.4-nightly.20261026010000.123456789abc"},
	} {
		if _, _, err := r.Plan(options); err == nil {
			t.Errorf("accepted invalid options %+v", options)
		}
	}
	git(t, r, "tag", "v1.2.4-rc.1")
	if _, _, err := r.Plan(Options{Kind: "custom", Custom: "1.2.4-rc.1"}); err == nil {
		t.Error("accepted tag collision")
	}
	git(t, r, "tag", "v65535.65535.65535")
	for _, kind := range []string{"patch", "minor", "major", "nightly"} {
		if _, _, err := r.Plan(Options{Kind: kind}); err == nil {
			t.Errorf("accepted overflowing %s bump", kind)
		}
	}
}

func TestStableBuildMetadataIsStillVersionBase(t *testing.T) {
	r := fixture(t)
	commit(t, r, "feat: initial", "")
	git(t, r, "tag", "v1.2.3+build.1")
	plan, _, err := r.Plan(Options{Kind: "patch"})
	if err != nil || plan.Version != "1.2.4" || plan.PreviousTag != "v1.2.3+build.1" {
		t.Fatalf("ignored stable build metadata: %+v %v", plan, err)
	}
	git(t, r, "tag", "v65536.0.0")
	if _, _, err := r.Plan(Options{Kind: "patch"}); err == nil {
		t.Error("silently ignored unsupported existing stable version")
	}
}

func TestNightlySkipAndExactIncrementalRange(t *testing.T) {
	r := fixture(t)
	first := commit(t, r, "historical import", "")
	date := time.Date(2026, 10, 26, 1, 2, 3, 0, time.UTC)
	plan, notes, err := r.Plan(Options{Kind: "nightly", Date: date})
	if err != nil {
		t.Fatal(err)
	}
	if plan.Version != "0.0.1-nightly.20261026010203."+first[:12] || !plan.Prerelease || !strings.Contains(notes, "historical import") {
		t.Fatalf("wrong first nightly: %+v\n%s", plan, notes)
	}
	git(t, r, "tag", plan.Tag)
	skipped, _, err := r.Plan(Options{Kind: "nightly", LastNightly: plan.Tag, Date: date.Add(time.Hour)})
	if err != nil || !skipped.Skip || skipped.SHA != first {
		t.Fatalf("did not skip unchanged main: %+v %v", skipped, err)
	}
	commit(t, r, "fix: new change", "")
	next, nextNotes, err := r.Plan(Options{Kind: "nightly", LastNightly: plan.Tag, Date: date.Add(time.Hour)})
	if err != nil || next.Skip || next.PreviousTag != plan.Tag || strings.Contains(nextNotes, "historical import") || !strings.Contains(nextNotes, "new change") {
		t.Fatalf("wrong incremental nightly: %+v %v\n%s", next, err, nextNotes)
	}
	// An unmarked build/tag does not count as a successfully published nightly.
	unmarked, _, err := r.Plan(Options{Kind: "nightly", Date: date.Add(2 * time.Hour)})
	if err != nil || unmarked.Skip {
		t.Fatalf("tag alone incorrectly skipped nightly: %+v %v", unmarked, err)
	}
	git(t, r, "checkout", "-b", "other")
	other := commit(t, r, "feat: other branch", "")
	otherTag := "v0.0.1-nightly.20261026050000." + other[:12]
	git(t, r, "tag", otherTag)
	git(t, r, "checkout", "main")
	if _, _, err := r.Plan(Options{Kind: "nightly", LastNightly: otherTag, Date: date}); err == nil {
		t.Error("accepted nonancestor nightly")
	}
}

func TestCommitRangeEnforcementAndMergeWrapper(t *testing.T) {
	r := fixture(t)
	base := commit(t, r, "historical import", "")
	commit(t, r, "feat: good commit", "")
	git(t, r, "checkout", "-b", "topic")
	commit(t, r, "fix: topic change", "")
	git(t, r, "checkout", "main")
	git(t, r, "merge", "--no-ff", "topic", "-m", "Merge pull request #2")
	if err := r.CheckCommits(base, "HEAD"); err != nil {
		t.Fatal(err)
	}
	bad := commit(t, r, "unconventional new change", "")
	if err := r.CheckCommits(base, "HEAD"); err == nil || !strings.Contains(err.Error(), bad) {
		t.Fatalf("did not report actual bad commit: %v", err)
	}
	if err := r.CheckCommits("--all", "HEAD"); err == nil {
		t.Error("accepted option injection")
	}
}

func TestWritePlanBundle(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "plan")
	plan := Plan{Version: "1.2.3-rc.1", PackageVersion: "1.2.3", Tag: "v1.2.3-rc.1", SHA: strings.Repeat("a", 40), Prerelease: true}
	if err := Write(dir, plan, "# notes\n"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "plan.json"))
	if err != nil {
		t.Fatal(err)
	}
	var decoded Plan
	if err := json.Unmarshal(data, &decoded); err != nil || decoded != plan {
		t.Fatalf("wrong serialized plan: %+v %v", decoded, err)
	}
	outputs, err := os.ReadFile(filepath.Join(dir, "github-output"))
	if err != nil || !strings.Contains(string(outputs), "prerelease=true\n") || !strings.Contains(string(outputs), "package_version=1.2.3\n") {
		t.Fatalf("wrong workflow outputs: %s %v", outputs, err)
	}
	if err := Write(dir, plan, "overwrite"); err == nil {
		t.Error("overwrote existing plan")
	}
}

func fixture(t *testing.T) Repository {
	t.Helper()
	r := Repository{Dir: t.TempDir()}
	git(t, r, "init", "-b", "main")
	git(t, r, "config", "user.name", "Release Tests")
	git(t, r, "config", "user.email", "tests@example.invalid")
	git(t, r, "config", "commit.gpgsign", "false")
	return r
}

func commit(t *testing.T, r Repository, subject, body string) string {
	t.Helper()
	git(t, r, "commit", "--allow-empty", "-m", subject, "-m", body)
	return strings.TrimSpace(git(t, r, "rev-parse", "HEAD"))
}

func git(t *testing.T, r Repository, args ...string) string {
	t.Helper()
	command := exec.Command("git", args...)
	command.Dir = r.Dir
	data, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, data)
	}
	return string(data)
}
