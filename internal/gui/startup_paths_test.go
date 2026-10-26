//go:build !headless

package gui

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/Dunky13/takedown/internal/config"
)

func relativeStartupScreen(t *testing.T) (*screen, *screenAccount, config.Settings) {
	t.Helper()
	root := t.TempDir()
	t.Chdir(root)
	s, account := newTestScreen(t, func(cfg *config.Settings) {
		cfg.Destination = "archives"
		cfg.OAuthClientFile = "client.json"
		cfg.PrivateKeyPath = "id_ed25519"
		cfg.KnownHostsPath = "known_hosts"
	})
	want := s.a.Settings
	want.Destination = filepath.Join(root, "archives")
	want.OAuthClientFile = filepath.Join(root, "client.json")
	want.PrivateKeyPath = filepath.Join(root, "id_ed25519")
	want.KnownHostsPath = filepath.Join(root, "known_hosts")
	return s, account, want
}

func assertStartupPaths(t *testing.T, got, want config.Settings) {
	t.Helper()
	for name, paths := range map[string][2]string{
		"destination":  {got.Destination, want.Destination},
		"OAuth client": {got.OAuthClientFile, want.OAuthClientFile},
		"SSH key":      {got.PrivateKeyPath, want.PrivateKeyPath},
		"known_hosts":  {got.KnownHostsPath, want.KnownHostsPath},
	} {
		if paths[0] != paths[1] {
			t.Errorf("%s path = %q, want %q", name, paths[0], paths[1])
		}
	}
}

func TestConfirmedStartupAdoptsAllRuntimePathsBeforeWorkerAndLaterSave(t *testing.T) {
	s, account, want := relativeStartupScreen(t)
	bindTestDiscovery(s)
	revision := s.sourceRevision
	identity := s.a.Settings.SourceIdentity()
	calls := 0
	s.startupUI.configure = func(_ context.Context, action string, _, dry bool) (string, error) {
		calls++
		if action != "install" || dry {
			t.Error("confirmed installation used different action options")
		}
		assertStartupPaths(t, s.a.Settings, want)
		return "Fake watcher installed. No OS or credential operations performed.", nil
	}
	test.Tap(s.startupUI.buttons["install"])
	if calls != 0 || filepath.IsAbs(s.destination.Text) || filepath.IsAbs(s.oauthClient.Text) {
		t.Fatal("unconfirmed install changed runtime paths or called the backend")
	}
	tapDialogButton(t, s, "Confirm startup change")
	awaitJob(t, s)
	if calls != 1 || !strings.Contains(s.status.Text, "Startup action completed") || s.w.Canvas().Overlays().Top() != nil {
		t.Fatal("same-path normalization caused a false error or reconfirmation")
	}
	assertStartupPaths(t, s.a.Settings, want)
	fields, err := s.preferences()
	if err != nil {
		t.Fatal(err)
	}
	assertStartupPaths(t, fields, want)
	if s.a.Account != account || s.sourceRevision != revision || s.a.Settings.SourceIdentity() != identity || !s.discovery.matches(s) || !s.discoveryBound {
		t.Fatal("same-target normalization changed the account, revision, or discovery authority")
	}
	// Exercise the ordinary Save preferences worker, which previously rewrote
	// normalized config from stale relative entries after installation.
	s.run("Saving preferences…", func(context.Context, sourceSnapshot) error { return nil }, nil, nil)
	awaitJob(t, s)
	saved, err := config.Load(s.a.DataDir)
	if err != nil {
		t.Fatal(err)
	}
	assertStartupPaths(t, saved, want)
	if s.a.Account != account || s.sourceRevision != revision {
		t.Fatal("the later save rebound the source")
	}
	s.oauthClient.SetText(filepath.Join(filepath.Dir(want.OAuthClientFile), "different-client.json"))
	if s.sourceRevision != revision+1 || s.discoveryBound {
		t.Fatal("the OAuth change handler was not restored after normalization")
	}
}

func TestStartupPreviewAndDismissalKeepRelativeRuntimePaths(t *testing.T) {
	s, account, _ := relativeStartupScreen(t)
	relative := s.a.Settings
	revision := s.sourceRevision
	calls := 0
	s.startupUI.configure = func(_ context.Context, action string, _, dry bool) (string, error) {
		calls++
		if action != "install" || !dry {
			t.Error("preview attempted a real startup change")
		}
		assertStartupPaths(t, s.a.Settings, relative)
		return "Fake definition preview", nil
	}
	test.Tap(s.startupUI.buttons["preview"])
	awaitJob(t, s)
	test.Tap(s.startupUI.buttons["install"])
	tapDialogButton(t, s, "Cancel")
	fields, err := s.preferences()
	if err != nil {
		t.Fatal(err)
	}
	assertStartupPaths(t, fields, relative)
	assertStartupPaths(t, s.a.Settings, relative)
	if calls != 1 || s.a.Account != account || s.sourceRevision != revision {
		t.Fatal("preview or dismissed install adopted paths or rebound the source")
	}
}

func TestStartupRelativePathEditsInvalidateOpenInstallConfirmation(t *testing.T) {
	for _, field := range []string{"destination", "client", "key", "known_hosts"} {
		t.Run(field, func(t *testing.T) {
			s, account, _ := relativeStartupScreen(t)
			calls := 0
			s.startupUI.configure = func(context.Context, string, bool, bool) (string, error) { calls++; return "", nil }
			test.Tap(s.startupUI.buttons["install"])
			switch field {
			case "destination":
				s.destination.SetText("other-archives")
			case "client":
				s.oauthClient.SetText("other-client.json")
			case "key":
				s.sshKey.SetText("other-key")
			case "known_hosts":
				s.knownHosts.SetText("other-known-hosts")
			}
			tapDialogButton(t, s, "Confirm startup change")
			if calls != 0 || s.busy || s.a.Account != account || !strings.Contains(s.status.Text, "Settings changed") {
				t.Fatal("a stale path confirmation reached the startup backend")
			}
		})
	}
}

func TestFailedStartupRetainsAdoptedPathsForLaterSave(t *testing.T) {
	s, account, want := relativeStartupScreen(t)
	revision := s.sourceRevision
	s.startupUI.configure = func(context.Context, string, bool, bool) (string, error) {
		return "", errors.New("fake prerequisite failure; no watcher installed")
	}
	test.Tap(s.startupUI.buttons["install"])
	tapDialogButton(t, s, "Confirm startup change")
	awaitJob(t, s)
	if !strings.Contains(s.banner.Text, "fake prerequisite failure") {
		t.Fatal("the actual install error was replaced by a false source-change error")
	}
	tapDialogButton(t, s, "Close")
	if err := s.save(); err != nil {
		t.Fatal(err)
	}
	saved, err := config.Load(s.a.DataDir)
	if err != nil {
		t.Fatal(err)
	}
	assertStartupPaths(t, saved, want)
	if s.a.Account != account || s.sourceRevision != revision {
		t.Fatal("failed install rebound the source or changed its revision")
	}
}
