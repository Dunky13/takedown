package config

import (
	"path/filepath"
	"testing"
)

func TestAbsoluteRuntimePathsPreserveSourceAndDestinationIdentity(t *testing.T) {
	t.Chdir(t.TempDir())
	original := Settings{Source: "drive", OAuthClientFile: "client.json", Destination: "archives", PrivateKeyPath: "keys/archive", KnownHostsPath: "keys/known_hosts"}
	resolved, err := original.AbsoluteRuntimePaths()
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{resolved.OAuthClientFile, resolved.Destination, resolved.PrivateKeyPath, resolved.KnownHostsPath} {
		if !filepath.IsAbs(path) {
			t.Fatalf("runtime path is relative: %q", path)
		}
	}
	if original.SourceIdentity() != resolved.SourceIdentity() {
		t.Fatal("same OAuth file spelling caused a false source change")
	}
	changed := resolved
	changed.OAuthClientFile = filepath.Join(t.TempDir(), "client.json")
	if resolved.SourceIdentity() == changed.SourceIdentity() {
		t.Fatal("different OAuth client paths shared a source binding")
	}
	for _, destination := range []string{"smb://operator@nas.example/share/archive", "sftp://operator@nas.example/archive"} {
		original.Destination = destination
		original.OAuthClientFile = ""
		original.PrivateKeyPath = ""
		original.KnownHostsPath = ""
		resolved, err = original.AbsoluteRuntimePaths()
		if err != nil || resolved.Destination != destination || resolved.OAuthClientFile != "" || resolved.PrivateKeyPath != "" || resolved.KnownHostsPath != "" {
			t.Fatalf("URL or empty optional path changed: %+v, %v", resolved, err)
		}
	}
}
