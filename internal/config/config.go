// Package config stores non-secret preferences. Credentials belong in the OS keyring.
package config

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

type Settings struct {
	Destination          string   `json:"destination"`
	SafetyBuffer         uint64   `json:"safety_buffer"`
	Concurrency          int      `json:"concurrency"`
	Email                string   `json:"email,omitempty"`
	Source               string   `json:"source,omitempty"`
	OAuthClientFile      string   `json:"oauth_client_file,omitempty"`
	DriveFolderID        string   `json:"drive_folder_id,omitempty"`
	PortabilityResources []string `json:"portability_resources,omitempty"`
	PollIntervalSeconds  int      `json:"poll_interval_seconds,omitempty"`
	PrivateKeyPath       string   `json:"private_key_path,omitempty"`
	KnownHostsPath       string   `json:"known_hosts_path,omitempty"`
	AllowUnknownSpace    bool     `json:"allow_unknown_space"`
}

func DefaultDataDir() (string, error) {
	root, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "takedown"), nil
}
func Defaults() (Settings, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return Settings{}, err
	}
	return Settings{Destination: filepath.Join(home, "Downloads", "takedown"), SafetyBuffer: 2_000_000_000, Concurrency: 2, Source: "drive", PollIntervalSeconds: 300, KnownHostsPath: filepath.Join(home, ".ssh", "known_hosts")}, nil
}

func (s Settings) EffectiveSource() string {
	if s.Source == "" {
		return "drive"
	}
	return s.Source
}

func (s Settings) PollInterval() time.Duration {
	if s.PollIntervalSeconds == 0 {
		return 5 * time.Minute
	}
	return time.Duration(s.PollIntervalSeconds) * time.Second
}

// SourceIdentity binds discovered URLs and background work to the full source
// configuration. An email is an OAuth login hint, never a verified account ID.
func (s Settings) SourceIdentity() string {
	resources := slices.Clone(s.PortabilityResources)
	slices.Sort(resources)
	client := s.OAuthClientFile
	invalidClient := false
	if client != "" {
		absolute, err := filepath.Abs(client)
		if err != nil {
			// Validate rejects this configuration. Retain its spelling and an
			// explicit invalid marker so it cannot match a valid source binding.
			invalidClient = true
		} else {
			client = absolute
		}
	}
	data, _ := json.Marshal(struct {
		Provider, Client, Hint, Folder string
		Resources                      []string
		InvalidClient                  bool
	}{s.EffectiveSource(), client, strings.ToLower(strings.TrimSpace(s.Email)), s.DriveFolderID, resources, invalidClient})
	return fmt.Sprintf("%x", sha256.Sum256(data))
}

func (s Settings) Validate() error {
	if s.OAuthClientFile != "" {
		if _, err := filepath.Abs(s.OAuthClientFile); err != nil {
			return errors.New("cannot resolve the OAuth client path; restore a valid working directory")
		}
	}
	if s.EffectiveSource() != "drive" && s.EffectiveSource() != "portability" {
		return errors.New("source must be drive or portability")
	}
	minimum := 60
	if s.EffectiveSource() == "portability" {
		minimum = 300
	}
	if s.PollIntervalSeconds != 0 && (s.PollIntervalSeconds < minimum || s.PollIntervalSeconds > 86400) {
		return fmt.Errorf("poll interval must be between %d and 86400 seconds", minimum)
	}
	for _, value := range []string{s.Email, s.OAuthClientFile, s.DriveFolderID} {
		if strings.ContainsAny(value, "\x00\r\n") {
			return errors.New("source settings contain invalid control characters")
		}
	}
	if s.Destination == "" {
		return errors.New("destination is required")
	}
	if strings.ContainsRune(s.Destination, 0) {
		return errors.New("destination contains a null byte")
	}
	if strings.Contains(s.Destination, "://") {
		u, err := url.Parse(s.Destination)
		if err != nil {
			return errors.New("invalid destination URL")
		}
		if u.Scheme != "smb" && u.Scheme != "sftp" {
			return errors.New("destination URL must use smb:// or sftp://")
		}
		if u.User != nil {
			if _, ok := u.User.Password(); ok {
				return errors.New("destination passwords must use the private password field or environment")
			}
		}
		if u.RawQuery != "" || u.ForceQuery || u.Fragment != "" {
			return errors.New("destination cannot contain query parameters or a fragment")
		}
	}
	if s.Concurrency < 1 || s.Concurrency > 16 {
		return errors.New("concurrency must be between 1 and 16")
	}
	if s.SafetyBuffer < 1 {
		return errors.New("safety buffer must be positive")
	}
	return nil
}
func Load(dir string) (Settings, error) {
	s, err := Defaults()
	if err != nil {
		return s, err
	}
	f, err := os.Open(filepath.Join(dir, "config.json"))
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return s, err
	}
	defer f.Close()
	d := json.NewDecoder(io.LimitReader(f, 1<<20))
	d.DisallowUnknownFields()
	// Older preferences may contain an obsolete Chromium executable path.
	// Read and discard it during migration; it is never executed or saved.
	wire := struct {
		*Settings
		LegacyBrowserPath *string `json:"browser_path,omitempty"`
	}{Settings: &s}
	if err = d.Decode(&wire); err != nil {
		return s, fmt.Errorf("read config: %w", err)
	}
	var extra json.RawMessage
	if err = d.Decode(&extra); err != io.EOF {
		return s, errors.New("configuration must contain one JSON object")
	}
	return s, s.Validate()
}
func Save(dir string, s Settings) error {
	if err := s.Validate(); err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	f, err := os.CreateTemp(dir, ".config-*")
	if err != nil {
		return err
	}
	name := f.Name()
	defer os.Remove(name)
	err = errors.Join(f.Chmod(0600), writeConfig(f, append(data, '\n')), f.Sync(), f.Close())
	if err != nil {
		return err
	}
	return os.Rename(name, filepath.Join(dir, "config.json"))
}
func writeConfig(f *os.File, b []byte) error { _, err := f.Write(b); return err }
