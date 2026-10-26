//go:build !headless

package gui

import (
	"errors"
	"strconv"
	"strings"

	"github.com/Dunky13/takedown/internal/config"
)

func (s *screen) preferences() (config.Settings, error) {
	cfg := s.a.Settings
	buffer, err := strconv.ParseUint(strings.TrimSpace(s.buffer.Text), 10, 64)
	if err != nil || buffer == 0 {
		return cfg, errors.New("enter a positive safety buffer in bytes")
	}
	workers, err := strconv.Atoi(s.workers.Selected)
	if err != nil {
		return cfg, errors.New("choose the number of concurrent archive parts")
	}
	interval := 0
	if text := strings.TrimSpace(s.pollInterval.Text); text != "" {
		interval, err = strconv.Atoi(text)
		if err != nil {
			return cfg, errors.New("enter a whole number of seconds for the watch interval")
		}
	}
	cfg.Destination = strings.TrimSpace(s.destination.Text)
	cfg.Email = strings.TrimSpace(s.email.Text)
	cfg.Source = s.selectedSource()
	cfg.OAuthClientFile = strings.TrimSpace(s.oauthClient.Text)
	cfg.DriveFolderID = strings.TrimSpace(s.driveFolder.Text)
	cfg.PortabilityResources = nil
	if cfg.Source == "portability" && strings.TrimSpace(s.resources.Text) != "" {
		for _, item := range strings.Split(s.resources.Text, ",") {
			resource := strings.TrimSpace(item)
			if resource == "" {
				return cfg, errors.New("enter resource group names separated by commas, without empty groups")
			}
			cfg.PortabilityResources = append(cfg.PortabilityResources, resource)
		}
	}
	cfg.PollIntervalSeconds = interval
	cfg.SafetyBuffer = buffer
	cfg.Concurrency = workers
	cfg.KnownHostsPath = strings.TrimSpace(s.knownHosts.Text)
	cfg.PrivateKeyPath = strings.TrimSpace(s.sshKey.Text)
	cfg.AllowUnknownSpace = s.unknown.Checked
	return cfg, cfg.Validate()
}

func (s *screen) save() error {
	if s.busy {
		return errors.New("stop the current job before changing preferences")
	}
	cfg, err := s.preferences()
	if err != nil {
		return err
	}
	changedSource := cfg.SourceIdentity() != s.a.Settings.SourceIdentity()
	if changedSource {
		// Clear URLs even when rebinding fails; no old discovery can authorize
		// downloads for the newly entered source settings.
		s.invalidateExports()
	}
	if err := s.a.Configure(cfg); err != nil {
		return err
	}
	return config.Save(s.a.DataDir, cfg)
}

// Adopt the confirmed install's runtime paths before the worker takes its
// source snapshot. App settings and the fields that ordinary saves read must
// agree, including when a startup backend is injected or installation fails.
func (s *screen) adoptStartupRuntimePaths() error {
	cfg, err := s.preferences()
	if err != nil {
		return err
	}
	normalized, err := cfg.AbsoluteRuntimePaths()
	if err != nil {
		return err
	}
	if cfg.SourceIdentity() != normalized.SourceIdentity() {
		return errors.New("runtime path normalization changed the source selection; review settings before installing")
	}
	if err := s.a.Configure(normalized); err != nil {
		return err
	}
	s.destination.SetText(normalized.Destination)
	// This is the same OAuth file, not a source edit. SourceIdentity canonicalizes
	// its path; keep the account and revision while restoring genuine edit hooks.
	func() {
		changed := s.oauthClient.OnChanged
		s.oauthClient.OnChanged = nil
		defer func() { s.oauthClient.OnChanged = changed }()
		s.oauthClient.SetText(normalized.OAuthClientFile)
	}()
	s.sshKey.SetText(normalized.PrivateKeyPath)
	s.knownHosts.SetText(normalized.KnownHostsPath)
	return nil
}
