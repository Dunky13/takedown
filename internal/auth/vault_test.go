package auth

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVaultMissingCorruptAndUnavailableNeverFallsBackToDisk(t *testing.T) {
	for _, encoded := range []string{"", `{"version":2}`, `{"version":1,"scopes":[],"token":{}}`, strings.Repeat("x", maxTokenBytes+1)} {
		opts := testOptions(t)
		store := &fakeSecrets{}
		v := newVault(opts, "client", store)
		if encoded != "" {
			_ = store.Set(vaultService, v.account, encoded)
		}
		if _, err := v.load(); !errors.Is(err, ErrAuthorizationRequired) {
			t.Fatalf("invalid token: %v", err)
		}
		entries, err := os.ReadDir(opts.DataDir)
		if err != nil || len(entries) != 0 {
			t.Fatal("token file created")
		}
	}
	store := &fakeSecrets{getErr: errors.New("raw keyring diagnostic secret"), setErr: errors.New("raw keyring diagnostic secret")}
	v := newVault(testOptions(t), "client", store)
	for _, err := range []error{func() error { _, err := v.load(); return err }(), v.save(validToken())} {
		if err == nil || strings.Contains(err.Error(), "diagnostic") {
			t.Fatalf("keyring failure: %v", err)
		}
	}
}

func TestForgetDeletesOnlyConfiguredAndLegacyCurrentProfile(t *testing.T) {
	opts := testOptions(t)
	store := &fakeSecrets{}
	m := testManager(t, opts, store, nil)
	seedToken(t, m, validToken())
	legacyAccount := "credentials-" + fingerprint(m.opts.DataDir) + "-" + fingerprint(m.opts.Email)
	otherLegacy := "credentials-" + fingerprint(m.opts.DataDir) + "-" + fingerprint("other@example.test")
	master := "aes256-" + fingerprint(m.opts.DataDir)
	for _, account := range []string{legacyAccount, otherLegacy, master} {
		_ = store.Set(legacyVaultService, account, "old encrypted secret")
	}
	for _, email := range []string{m.opts.Email, "other@example.test"} {
		for _, prefix := range []string{"session-", "profile-"} {
			if err := os.WriteFile(filepath.Join(opts.DataDir, prefix+fingerprint(email)+".enc"), []byte("encrypted"), 0600); err != nil {
				t.Fatal(err)
			}
		}
	}
	if err := m.Forget(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Get(vaultService, m.vault.account); err == nil {
		t.Fatal("OAuth token retained")
	}
	if _, err := store.Get(legacyVaultService, legacyAccount); err == nil {
		t.Fatal("legacy password/TOTP retained")
	}
	for _, account := range []string{otherLegacy, master} {
		if _, err := store.Get(legacyVaultService, account); err != nil {
			t.Fatal("unrelated/shared legacy key removed")
		}
	}
	for _, prefix := range []string{"session-", "profile-"} {
		if _, err := os.Stat(filepath.Join(opts.DataDir, prefix+fingerprint(m.opts.Email)+".enc")); !errors.Is(err, os.ErrNotExist) {
			t.Fatal("current legacy snapshot retained")
		}
		if _, err := os.Stat(filepath.Join(opts.DataDir, prefix+fingerprint("other@example.test")+".enc")); err != nil {
			t.Fatal("other profile removed")
		}
	}
	if err := m.Refresh(context.Background()); !errors.Is(err, ErrAuthorizationRequired) {
		t.Fatal("forgotten token reused")
	}
	if err := m.Forget(context.Background()); err != nil {
		t.Fatalf("Forget not idempotent: %v", err)
	}
}

func TestForgetWorksBeforeClientAndScopesConfigured(t *testing.T) {
	opts := Options{DataDir: t.TempDir(), Email: "OLD@example.test"}
	store := &fakeSecrets{}
	m := testManager(t, opts, store, nil)
	account := "credentials-" + fingerprint(m.opts.DataDir) + "-" + fingerprint(m.opts.Email)
	_ = store.Set(legacyVaultService, account, "legacy encrypted password")
	if err := m.Forget(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Get(legacyVaultService, account); err == nil {
		t.Fatal("legacy password retained")
	}
}

func TestLegacyCleanupSymlinkRemovesEntryNotTargetAndRefusesDirectory(t *testing.T) {
	opts := Options{DataDir: t.TempDir(), Email: "hint@example.test"}
	target := filepath.Join(t.TempDir(), "keep")
	if err := os.WriteFile(target, []byte("keep"), 0600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(opts.DataDir, "session-"+fingerprint(opts.Email)+".enc")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	directory := filepath.Join(opts.DataDir, "profile-"+fingerprint(opts.Email)+".enc")
	if err := os.Mkdir(directory, 0700); err != nil {
		t.Fatal(err)
	}
	if err := forgetLegacy(opts, &fakeSecrets{}); err == nil {
		t.Fatal("directory accepted as snapshot")
	}
	if _, err := os.Lstat(link); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("symlink not removed")
	}
	if _, err := os.Stat(target); err != nil {
		t.Fatal("symlink target removed")
	}
	if info, err := os.Stat(directory); err != nil || !info.IsDir() {
		t.Fatal("directory removed")
	}
}

func TestForgetReportsKeyringFailureButStillCleansKnownFiles(t *testing.T) {
	opts := testOptions(t)
	store := &fakeSecrets{deleteErr: errors.New("raw diagnostic secret")}
	m := testManager(t, opts, store, nil)
	path := filepath.Join(opts.DataDir, "session-"+fingerprint(m.opts.Email)+".enc")
	_ = os.WriteFile(path, []byte("encrypted"), 0600)
	err := m.Forget(context.Background())
	if err == nil || strings.Contains(err.Error(), "diagnostic") {
		t.Fatalf("unsafe/missing error: %v", err)
	}
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("known file not cleaned after keyring failure")
	}
}

func TestProfileIDUsesClientHintDirectoryButNotScopesOrKeyring(t *testing.T) {
	opts := testOptions(t)
	opts.Email = "  HINT@EXAMPLE.TEST  "
	opts.Scopes = nil
	store := &fakeSecrets{getErr: errors.New("keyring unavailable")}
	m := testManager(t, opts, store, nil)
	profile, err := m.ProfileID(context.Background())
	if err != nil || len(profile) != 64 {
		t.Fatalf("profile: %q error %v", profile, err)
	}
	other := opts
	other.Email = "hint@example.test"
	other.Scopes = []string{testScope}
	m2 := testManager(t, other, store, nil)
	same, err := m2.ProfileID(context.Background())
	if err != nil || same != profile {
		t.Fatal("scopes/hint normalization changed profile slot")
	}
	if store.gets+store.sets+store.deletes != 0 {
		t.Fatal("ProfileID accessed keyring")
	}
	if err := m.Login(context.Background()); !errors.Is(err, ErrAuthorizationRequired) {
		t.Fatal("ProfileID bypassed missing scopes")
	}
	for _, modify := range []func(*Options){func(o *Options) { o.Email = "other@example.test" }, func(o *Options) { o.DataDir = t.TempDir() }, func(o *Options) {
		doc := clientDocument()
		doc["installed"].(map[string]any)["client_id"] = "456-other.apps.googleusercontent.com"
		o.ClientFile = writeClient(t, doc)
	}} {
		changed := opts
		modify(&changed)
		manager := testManager(t, changed, store, nil)
		id, err := manager.ProfileID(context.Background())
		if err != nil || id == profile {
			t.Fatal("profile namespace collision")
		}
	}
}
