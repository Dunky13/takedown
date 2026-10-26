package autostart

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func startupLifecycleFixture(t *testing.T, script string) (Options, Definition, string) {
	t.Helper()
	if runtime.GOOS != "darwin" {
		t.Skip("launchd fixture is Darwin-specific")
	}
	root := t.TempDir()
	exe := filepath.Join(root, "takedown")
	if err := os.WriteFile(exe, []byte("#!/bin/sh\nexit 0\n"), 0700); err != nil {
		t.Fatal(err)
	}
	o := Options{Home: root, Executable: exe, DataDir: filepath.Join(root, "data"), UID: "501"}
	d, err := Build(o)
	if err != nil {
		t.Fatal(err)
	}
	fake := filepath.Join(root, "bin")
	if err = os.Mkdir(fake, 0700); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(fake, "launchctl"), []byte("#!/bin/sh\n"+script), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", fake)
	marker := filepath.Join(root, "registered")
	t.Setenv("REVIEW_REGISTERED", marker)
	return o, d, marker
}

func TestStartupRepeatedInstall(t *testing.T) {
	o, _, _ := startupLifecycleFixture(t, `case "$1" in
 bootstrap) if test -e "$REVIEW_REGISTERED"; then echo 'Bootstrap failed: 5: Input/output error' >&2; exit 5; fi; : > "$REVIEW_REGISTERED";;
 print) test -e "$REVIEW_REGISTERED" || exit 113;;
 *) exit 99;;
esac
`)
	if _, err := Install(context.Background(), o); err != nil {
		t.Fatal(err)
	}
	if _, err := Install(context.Background(), o); err != nil {
		t.Fatalf("identical already-loaded installation must be idempotent: %v", err)
	}
}

func TestStartupRemoveUnloadedDefinition(t *testing.T) {
	for _, symlink := range []bool{false, true} {
		t.Run(map[bool]string{false: "regular", true: "stale-symlink"}[symlink], func(t *testing.T) {
			o, d, _ := startupLifecycleFixture(t, `case "$1" in
 bootout|print) echo 'Could not find service in domain for user gui: 501' >&2; exit 113;;
 *) exit 99;;
esac
`)
			if err := os.MkdirAll(filepath.Dir(d.Path), 0700); err != nil {
				t.Fatal(err)
			}
			if symlink {
				if err := os.Symlink(filepath.Join(t.TempDir(), "missing-target.plist"), d.Path); err != nil {
					t.Fatal(err)
				}
			} else if err := os.WriteFile(d.Path, []byte(d.Content), 0600); err != nil {
				t.Fatal(err)
			}
			err := Remove(context.Background(), o)
			_, statErr := os.Lstat(d.Path)
			if err != nil || !os.IsNotExist(statErr) {
				t.Fatalf("unloaded definition must be removable; err=%v remaining=%v", err, statErr)
			}
		})
	}
}

func TestStartupInstallRejectsDefinitionSymlink(t *testing.T) {
	o, d, marker := startupLifecycleFixture(t, `: > "$REVIEW_REGISTERED"
exit 0
`)
	if err := os.MkdirAll(filepath.Dir(d.Path), 0700); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(t.TempDir(), "external.plist")
	if err := os.WriteFile(target, []byte(d.Content), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, d.Path); err != nil {
		t.Fatal(err)
	}
	_, err := Install(context.Background(), o)
	_, called := os.Stat(marker)
	if err == nil || !os.IsNotExist(called) {
		t.Fatalf("must reject definition symlink before OS registration: err=%v command_called=%v", err, called == nil)
	}
	got, readErr := os.ReadFile(target)
	if readErr != nil || !strings.Contains(string(got), "<plist") {
		t.Fatal("fixture target unexpectedly changed", readErr)
	}
}

func TestStartupRemoveMissingDefinitionStillStopsLoadedService(t *testing.T) {
	o, _, marker := startupLifecycleFixture(t, `case "$1" in
 print) exit 0;;
 bootout) printf stopped > "$REVIEW_REGISTERED";;
 *) exit 99;;
esac
`)
	if err := Remove(context.Background(), o); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(marker)
	if err != nil || string(got) != "stopped" {
		t.Fatalf("missing definition left its service loaded: %q, %v", got, err)
	}
}

func TestStartupLifecycleFailureRetainsDefinition(t *testing.T) {
	for _, failAt := range []string{"print", "bootout"} {
		t.Run(failAt, func(t *testing.T) {
			o, d, marker := startupLifecycleFixture(t, `if test "$1" = "$FAIL_AT"; then echo 'Permission denied' >&2; exit 5; fi
case "$1" in
 print) exit 0;;
 *) printf unexpected > "$REVIEW_REGISTERED"; exit 99;;
esac
`)
			t.Setenv("FAIL_AT", failAt)
			if err := os.MkdirAll(filepath.Dir(d.Path), 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(d.Path, []byte(d.Content), 0600); err != nil {
				t.Fatal(err)
			}
			if err := Remove(context.Background(), o); err == nil {
				t.Fatal("unclassified lifecycle error was hidden")
			}
			if _, err := os.Stat(d.Path); err != nil {
				t.Fatalf("failed removal lost its definition: %v", err)
			}
			if _, err := os.Stat(marker); !os.IsNotExist(err) {
				t.Fatal("unexpected OS command was executed")
			}
		})
	}
}

func TestStartupRemoveSymlinkKeepsTargetAndAllowsChangedReinstall(t *testing.T) {
	o, d, marker := startupLifecycleFixture(t, `case "$1" in
 print) if test -f "$REVIEW_REGISTERED"; then read -r state < "$REVIEW_REGISTERED"; test "$state" = running && exit 0; fi; exit 113;;
 bootstrap) printf 'running\n' > "$REVIEW_REGISTERED";;
 bootout) printf 'stopped\n' > "$REVIEW_REGISTERED";;
 *) exit 99;;
esac
`)
	if _, err := Install(context.Background(), o); err != nil {
		t.Fatal(err)
	}
	if err := Remove(context.Background(), o); err != nil {
		t.Fatal(err)
	}
	o.PartByPart = true
	if _, err := Install(context.Background(), o); err != nil {
		t.Fatalf("changed definition could not replace a removed watch: %v", err)
	}
	if err := os.Remove(d.Path); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(t.TempDir(), "unrelated.plist")
	if err := os.WriteFile(target, []byte("KEEP"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, d.Path); err != nil {
		t.Fatal(err)
	}
	if err := Remove(context.Background(), o); err != nil {
		t.Fatal(err)
	}
	if data, err := os.ReadFile(target); err != nil || string(data) != "KEEP" {
		t.Fatal("removal followed a definition symlink", err)
	}
	if data, err := os.ReadFile(marker); err != nil || string(data) != "stopped\n" {
		t.Fatal("symlink removal left the loaded service running", err)
	}
}
