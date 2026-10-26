package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResetWithoutConfirmationCannotStartOAuthOrMutation(t *testing.T) {
	dir := t.TempDir()
	cmd := NewCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"jobs", "reset-authorization", "--data-dir", dir, "--source", "portability", "--resources", "myactivity.search"})
	err := cmd.ExecuteContext(context.Background())
	if err == nil || !strings.Contains(err.Error(), "--confirm") {
		t.Fatalf("missing confirmation guard: %v", err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != 0 {
		t.Fatal("unconfirmed reset wrote profile state")
	}
}

func TestServiceDryRunCreatesNoConfigurationOrActivation(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "profile")
	cmd := NewCommand()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"service", "install", "--dry-run", "--data-dir", dir, "--part-by-part"})
	if err := cmd.ExecuteContext(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "watch") || !strings.Contains(out.String(), "--part-by-part") {
		t.Fatal("watch definition not rendered")
	}
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatal("dry run changed the profile")
	}
}
