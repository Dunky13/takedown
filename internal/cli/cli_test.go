package cli

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Dunky13/takedown/internal/model"
	"github.com/Dunky13/takedown/internal/state"
)

func TestSelectPartsFailsForMissingRequestedExport(t *testing.T) {
	if _, err := SelectParts([]model.Export{{ID: "ready"}}, []string{"missing"}); err == nil {
		t.Fatal("unknown export silently ignored")
	}
}
func TestStatusWorksWhileWriterHoldsStateLock(t *testing.T) {
	dir := t.TempDir()
	s, err := state.Open(filepath.Join(dir, "takedown-state.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if err = s.Put(state.Record{Part: model.Part{ExportID: "export", Index: 1, Filename: "archive.zip", Size: 100}, Destination: "./archives", Offset: 50}); err != nil {
		t.Fatal(err)
	}
	cmd := NewCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"status", "--data-dir", dir, "--json"})
	if err = cmd.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `"offset":50`) {
		t.Fatalf("unexpected status: %s", out.String())
	}
}
func TestHelpDoesNotAuthenticate(t *testing.T) {
	cmd := NewCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--help"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"auth", "download", "status", "gui", "export", "watch", "oauth-client", "drive-folder"} {
		if !strings.Contains(out.String(), name) {
			t.Fatalf("missing command %s", name)
		}
	}
	for _, removed := range []string{"store-credentials", "forget-credentials", "--browser"} {
		if strings.Contains(out.String(), removed) {
			t.Fatalf("obsolete account automation exposed: %s", removed)
		}
	}
}

func TestMissingOAuthConfigurationFailsWithoutStartingBrowser(t *testing.T) {
	cmd := NewCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"auth", "--data-dir", t.TempDir()})
	err := cmd.ExecuteContext(context.Background())
	if err == nil || !strings.Contains(err.Error(), "Desktop app client") {
		t.Fatalf("missing configuration: %v", err)
	}
	if !strings.Contains(out.String(), "read-only access") || strings.Contains(out.String(), "Google password:") {
		t.Fatal("OAuth disclosure missing or credential prompt retained")
	}
}

func TestInvalidPollingIntervalFailsBeforeAuthorization(t *testing.T) {
	cmd := NewCommand()
	cmd.SetArgs([]string{"watch", "--data-dir", t.TempDir(), "--poll-interval", "1s"})
	if err := cmd.ExecuteContext(context.Background()); err == nil || !strings.Contains(err.Error(), "poll interval") {
		t.Fatalf("poll limit bypassed: %v", err)
	}
}

func TestPortabilityEncryptionDisclosurePrecedesAuthorization(t *testing.T) {
	cmd := NewCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"auth", "--data-dir", t.TempDir(), "--source", "portability", "--resources", "myactivity.search"})
	if err := cmd.ExecuteContext(context.Background()); err == nil || !strings.Contains(err.Error(), "Desktop app client") {
		t.Fatalf("expected configuration failure before browser access: %v", err)
	}
	for _, requirement := range []string{"encryption in transit and at rest", "does not encrypt destination files", "SFTP encrypts transit only", "SMB encryption depends"} {
		if !strings.Contains(out.String(), requirement) {
			t.Fatalf("missing disclosure %q", requirement)
		}
	}
}
