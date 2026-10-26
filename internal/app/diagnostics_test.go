package app

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/Dunky13/takedown/internal/config"
	"github.com/Dunky13/takedown/internal/events"
)

func TestLocalDiagnosticsAvoidImplicitSourceRequestsAndStorageCreation(t *testing.T) {
	dir := t.TempDir()
	account := &watchAccount{}
	a := &App{DataDir: filepath.Join(dir, "config"), Settings: config.Settings{Destination: filepath.Join(dir, "archives"), SafetyBuffer: 1, Concurrency: 2}, Bus: events.New(), Account: account}
	defer a.Close()
	r, err := a.Diagnose(context.Background(), DiagnosticOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if r.Healthy() {
		t.Fatal("missing OAuth configuration was reported healthy")
	}
	if account.created != 0 {
		t.Fatal("diagnostics created an export")
	}
	if _, err = os.Stat(a.Settings.Destination); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("local-only diagnostics created the destination")
	}
	if err = a.beginJob(); err != nil {
		t.Fatal("diagnostic retained active-job lock")
	}
	a.endJob()
}
