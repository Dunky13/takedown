// Package app is the shared job orchestrator for the CLI and the native GUI.
package app

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Dunky13/takedown/internal/auth"
	"github.com/Dunky13/takedown/internal/config"
	"github.com/Dunky13/takedown/internal/downloader"
	"github.com/Dunky13/takedown/internal/events"
	"github.com/Dunky13/takedown/internal/model"
	"github.com/Dunky13/takedown/internal/source"
	"github.com/Dunky13/takedown/internal/space"
	"github.com/Dunky13/takedown/internal/state"
	"github.com/Dunky13/takedown/internal/storage"
)

type Account interface {
	Login(context.Context) error
	List(context.Context) ([]model.Export, error)
	RefreshLink(context.Context, string, int) (string, error)
	HTTPClient() *http.Client
	Forget(context.Context) error
	CreateExport(context.Context) (string, error)
	Close() error
}
type App struct {
	DataDir  string
	Settings config.Settings
	Bus      *events.Bus
	Account  Account
	mu       sync.Mutex
	active   bool
	engine   *downloader.Engine
	paused   bool
}
type DownloadOptions struct {
	Destination     string
	PartByPart      bool
	StoragePassword string
}

func New(dataDir string, s config.Settings) (*App, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}
	b := events.New()
	a, err := newAccount(dataDir, s, b)
	if err != nil {
		b.Close()
		return nil, err
	}
	return &App{DataDir: dataDir, Settings: s, Bus: b, Account: a}, nil
}

func newAccount(dataDir string, s config.Settings, bus *events.Bus) (Account, error) {
	scopes, err := source.Scopes(s.EffectiveSource(), s.PortabilityResources)
	if err != nil {
		return nil, err
	}
	session, err := auth.New(auth.Options{DataDir: dataDir, ClientFile: s.OAuthClientFile, Email: s.Email, Scopes: scopes, Bus: bus})
	if err != nil {
		return nil, err
	}
	account, err := source.New(source.Options{DataDir: dataDir, Provider: s.EffectiveSource(), DriveFolderID: s.DriveFolderID, Resources: s.PortabilityResources, Bus: bus, Session: session})
	if err != nil {
		_ = session.Close()
		return nil, err
	}
	return account, nil
}
func (a *App) Close() error { err := a.Account.Close(); a.Bus.Close(); return err }

// Configure is called before starting work; it refreshes the account selection
// when a GUI user changes the account hint, OAuth client, or export source.
func (a *App) Configure(s config.Settings) error {
	if err := s.Validate(); err != nil {
		return err
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.active {
		return errors.New("stop the current job before changing preferences")
	}
	if s.SourceIdentity() != a.Settings.SourceIdentity() {
		next, err := newAccount(a.DataDir, s, a.Bus)
		if err != nil {
			return err
		}
		if err = a.Account.Close(); err != nil {
			_ = next.Close()
			return err
		}
		a.Account = next
	}
	a.Settings = s
	return nil
}
func (a *App) Pause() {
	a.mu.Lock()
	a.paused = true
	engine := a.engine
	a.mu.Unlock()
	if engine != nil {
		engine.Pause()
	}
}
func (a *App) Resume() {
	a.mu.Lock()
	a.paused = false
	engine := a.engine
	a.mu.Unlock()
	if engine != nil {
		engine.Resume()
	}
}

func (a *App) Download(ctx context.Context, parts []model.Part, opts DownloadOptions) (retErr error) {
	if err := a.beginJob(); err != nil {
		return err
	}
	defer a.endJob()
	return a.download(ctx, parts, opts)
}

func (a *App) beginJob() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.active {
		return errors.New("a job is already active")
	}
	a.active = true
	a.paused = false
	return nil
}

func (a *App) endJob() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.active, a.paused, a.engine = false, false, nil
}

func (a *App) download(ctx context.Context, parts []model.Part, opts DownloadOptions) (retErr error) {
	defer func() { a.mu.Lock(); a.engine = nil; a.mu.Unlock() }()
	if len(parts) == 0 {
		return errors.New("no ready archive parts selected")
	}
	dest := opts.Destination
	if dest == "" {
		dest = a.Settings.Destination
	}
	identity, err := DestinationIdentity(dest)
	if err != nil {
		return err
	}
	opts.StoragePassword, err = a.storagePassword(dest, opts.StoragePassword)
	if err != nil {
		return err
	}
	backend, err := storage.Open(ctx, dest, storage.Options{Password: opts.StoragePassword, PrivateKeyPath: a.Settings.PrivateKeyPath, KnownHostsPath: a.Settings.KnownHostsPath, AllowUnknownSpace: a.Settings.AllowUnknownSpace})
	if err != nil {
		return err
	}
	defer func() { retErr = errors.Join(retErr, backend.Close()) }()
	s, err := state.Open(filepath.Join(a.DataDir, "takedown-state.json"))
	if err != nil {
		return err
	}
	defer func() { retErr = errors.Join(retErr, s.Close()) }()
	parts = layoutParts(parts, identity, s.Snapshot())
	seen := make(map[string]bool)
	for _, p := range parts {
		if err = p.Validate(); err != nil {
			return err
		}
		name := strings.ToLower(p.StoragePath())
		if seen[name] {
			return errors.New("archive output paths collide; preserve the legacy files and select a separate destination")
		}
		seen[name] = true
	}
	g := space.New(backend, a.Settings.SafetyBuffer, a.Bus)
	g.AllowUnknownSpace(a.Settings.AllowUnknownSpace)
	var remaining int64
	for _, p := range parts {
		need := p.Size
		if info, e := backend.Stat(ctx, p.StoragePath()); e == nil {
			if info.Size() != p.Size {
				return fmt.Errorf("existing archive %s has an unexpected size", p.Filename)
			}
			need = 0
		} else if !isNotExist(e) {
			return e
		}
		if need > 0 {
			if info, e := backend.Stat(ctx, p.StoragePath()+".part"); e == nil {
				if info.Size() > p.Size {
					return fmt.Errorf("partial archive %s is larger than expected", p.Filename)
				}
				need -= info.Size()
			} else if !isNotExist(e) {
				return e
			}
		}
		if need > math.MaxInt64-remaining {
			return errors.New("batch byte count overflow")
		}
		remaining += need
	}
	if !opts.PartByPart {
		if err = g.Preflight(ctx, remaining); err != nil {
			return err
		}
	}
	e, err := downloader.New(downloader.Options{Backend: backend, Destination: identity, State: s, Guard: g, Bus: a.Bus, Client: a.Account.HTTPClient(), Links: a.Account, Concurrency: a.Settings.Concurrency, PartByPart: opts.PartByPart, CheckInterval: 50_000_000, MaxRetries: 3})
	if err != nil {
		return err
	}
	a.mu.Lock()
	a.engine = e
	if a.paused {
		e.Pause()
	}
	a.mu.Unlock()
	a.Bus.Publish(events.Event{Kind: events.Status, Message: "Download started"})
	if err = e.Download(ctx, parts); err != nil {
		return err
	}
	a.Bus.Publish(events.Event{Kind: events.Status, Message: "All selected archives verified and downloaded"})
	return nil
}
func DestinationIdentity(dest string) (string, error) {
	if strings.HasPrefix(dest, "smb://") || strings.HasPrefix(dest, "sftp://") {
		u, err := url.Parse(dest)
		if err != nil {
			return "", errors.New("invalid destination URL")
		}
		if u.User != nil {
			if _, ok := u.User.Password(); ok {
				return "", errors.New("supply network passwords using the password prompt or environment, never the destination URL")
			}
		}
		return u.String(), nil
	}
	return filepath.Abs(dest)
}
