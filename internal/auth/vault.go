package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"slices"

	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

const (
	vaultService       = "github.com/Dunky13/takedown/oauth2"
	legacyVaultService = "github.com/Dunky13/takedown"
	maxTokenBytes      = 64 << 10
)

type secretStore interface {
	Get(service, account string) (string, error)
	Set(service, account, value string) error
	Delete(service, account string) error
}

type osSecrets struct{}

func (osSecrets) Get(s, a string) (string, error) { return keyring.Get(s, a) }
func (osSecrets) Set(s, a, v string) error        { return keyring.Set(s, a, v) }
func (osSecrets) Delete(s, a string) error        { return keyring.Delete(s, a) }

func fingerprint(value string) string {
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}

type vault struct {
	account string
	scopes  []string
	store   secretStore
}

type tokenRecord struct {
	Version int          `json:"version"`
	Scopes  []string     `json:"scopes"`
	Token   oauth2.Token `json:"token"`
}

func newVault(opts Options, clientID string, store secretStore) *vault {
	// JSON preserves field boundaries. New normalized the scope ordering.
	binding, _ := json.Marshal([]any{opts.DataDir, opts.Email, clientID, opts.Scopes})
	return &vault{account: "token-v1-" + fingerprint(string(binding)), scopes: opts.Scopes, store: store}
}

func (v *vault) load() (*oauth2.Token, error) {
	encoded, err := v.store.Get(vaultService, v.account)
	if errors.Is(err, keyring.ErrNotFound) {
		return nil, ErrAuthorizationRequired
	}
	if err != nil {
		return nil, errors.New("OS keyring unavailable: unlock/configure it before reading OAuth tokens")
	}
	var record tokenRecord
	if len(encoded) > maxTokenBytes || json.Unmarshal([]byte(encoded), &record) != nil || record.Version != 1 || !slices.Equal(record.Scopes, v.scopes) {
		return nil, ErrAuthorizationRequired
	}
	token := &record.Token
	if token.AccessToken == "" || token.Expiry.IsZero() || token.Type() != "Bearer" {
		return nil, ErrAuthorizationRequired
	}
	return token, nil
}

func (v *vault) save(token *oauth2.Token) error {
	encoded, err := json.Marshal(tokenRecord{Version: 1, Scopes: v.scopes, Token: *token})
	if err != nil || len(encoded) > maxTokenBytes {
		return errors.New("cannot encode OAuth token within the secure storage limit")
	}
	defer clear(encoded)
	if v.store.Set(vaultService, v.account, string(encoded)) != nil {
		return errors.New("OS keyring unavailable: cannot persist OAuth tokens; unlock/configure it and retry")
	}
	return nil
}

func (v *vault) forget() error {
	err := v.store.Delete(vaultService, v.account)
	if err != nil && !errors.Is(err, keyring.ErrNotFound) {
		return errors.New("OS keyring unavailable: cannot forget OAuth tokens; unlock/configure it and retry")
	}
	return nil
}

func forgetLegacy(opts Options, store secretStore) error {
	var errs []error
	if opts.Email != "" {
		account := "credentials-" + fingerprint(opts.DataDir) + "-" + fingerprint(opts.Email)
		if err := store.Delete(legacyVaultService, account); err != nil && !errors.Is(err, keyring.ErrNotFound) {
			errs = append(errs, errors.New("OS keyring unavailable: cannot forget the current profile's legacy password/TOTP slot"))
		}
	}
	root, err := os.OpenRoot(opts.DataDir)
	if errors.Is(err, os.ErrNotExist) {
		return errors.Join(errs...)
	}
	if err != nil {
		return errors.Join(append(errs, errors.New("cannot open the authentication data directory for legacy cleanup"))...)
	}
	defer root.Close()
	for _, prefix := range []string{"session-", "profile-"} {
		name := prefix + fingerprint(opts.Email) + ".enc"
		info, err := root.Lstat(name)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil || (!info.Mode().IsRegular() && info.Mode()&os.ModeSymlink == 0) {
			errs = append(errs, errors.New("cannot safely remove a legacy authentication snapshot; expected a regular file or symlink"))
			continue
		}
		if root.Remove(name) != nil {
			errs = append(errs, errors.New("cannot remove the current profile's legacy authentication snapshot"))
		}
	}
	// The legacy AES key is shared by all accounts in this directory; retain
	// it so that this scoped cleanup cannot destroy another account's secrets.
	return errors.Join(errs...)
}
