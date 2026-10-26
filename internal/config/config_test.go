package config

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestPreferencesRoundTripAndRejectMalformed(t *testing.T) {
	dir := t.TempDir()
	s, err := Defaults()
	if err != nil {
		t.Fatal(err)
	}
	s.Destination = "sftp://operator@archive.example/exports"
	s.AllowUnknownSpace = true
	if err = Save(dir, s); err != nil {
		t.Fatal(err)
	}
	got, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, s) {
		t.Fatalf("preferences changed: %+v", got)
	}
	if err = os.WriteFile(filepath.Join(dir, "config.json"), []byte(`{"concurrency":0}`), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = Load(dir); err == nil {
		t.Fatal("invalid concurrency accepted")
	}
	if err = os.WriteFile(filepath.Join(dir, "config.json"), []byte(`{"password":"secret"}`), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = Load(dir); err == nil {
		t.Fatal("unknown secret field accepted")
	}
}

func TestSecretInDestinationIsNeverSaved(t *testing.T) {
	dir := t.TempDir()
	s, err := Defaults()
	if err != nil {
		t.Fatal(err)
	}
	s.Destination = "sftp://operator:private-password@nas.example/exports"
	if err = Save(dir, s); err == nil {
		t.Fatal("password-bearing destination saved")
	}
	if _, err = os.Stat(filepath.Join(dir, "config.json")); !os.IsNotExist(err) {
		t.Fatal("invalid preferences created a file")
	}
}

func TestSourceIdentityTracksPermissionsAndIgnoresResourceOrder(t *testing.T) {
	s, err := Defaults()
	if err != nil {
		t.Fatal(err)
	}
	s.Source = "portability"
	s.PortabilityResources = []string{"myactivity.search", "youtube.public_videos"}
	identity := s.SourceIdentity()
	s.PortabilityResources = []string{"youtube.public_videos", "myactivity.search"}
	if s.SourceIdentity() != identity {
		t.Fatal("resource order changed the source identity")
	}
	s.PortabilityResources = []string{"myactivity.search"}
	if s.SourceIdentity() == identity {
		t.Fatal("changed permissions kept stale source identity")
	}
	s.Source = "drive"
	identity = s.SourceIdentity()
	s.DriveFolderID = "another-folder"
	if s.SourceIdentity() == identity {
		t.Fatal("changed Drive folder kept stale source identity")
	}
}

func TestOldBrowserPreferenceIsMigratedWithoutExecutionOrResave(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte(`{"browser_path":"obsolete-chromium"}`), 0600); err != nil {
		t.Fatal(err)
	}
	s, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err = Save(dir, s); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "config.json"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "browser_path") {
		t.Fatal("obsolete browser automation preference retained")
	}
}

func TestProviderPollingLimits(t *testing.T) {
	s, err := Defaults()
	if err != nil {
		t.Fatal(err)
	}
	s.PollIntervalSeconds = 60
	if err = s.Validate(); err != nil {
		t.Fatal(err)
	}
	s.Source = "portability"
	if err = s.Validate(); err == nil {
		t.Fatal("Data Portability over-polling allowed")
	}
	s.PollIntervalSeconds = 300
	if err = s.Validate(); err != nil {
		t.Fatal(err)
	}
}
