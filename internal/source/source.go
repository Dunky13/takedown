// Package source discovers archives using Google's documented APIs. It never
// automates a browser or creates an export as a side effect of discovery.
package source

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Dunky13/takedown/internal/auth"
	"github.com/Dunky13/takedown/internal/events"
	"github.com/Dunky13/takedown/internal/model"
	"google.golang.org/api/dataportability/v1"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type Session interface {
	Login(context.Context) error
	Refresh(context.Context) error
	Forget(context.Context) error
	HTTPClient() *http.Client
	ProfileID(context.Context) (string, error)
	Close() error
}

type Options struct {
	DataDir, Provider, DriveFolderID string
	Resources                        []string
	Bus                              *events.Bus
	Session                          Session
	OpenURL                          func(string) error
}

type Manager struct {
	opts        Options
	drive       *drive.Service
	portability *dataportability.Service
	// Signed archive requests use a separate client with no OAuth transport or jar.
	public    *http.Client
	now       func() time.Time
	gate      chan struct{}
	life      context.Context
	cancel    context.CancelFunc
	closeOnce sync.Once
	closeErr  error
	links     map[string]cachedLinks
	persist   func(string, *portabilityState) error
}

// New performs no authentication, network requests, or disk writes. Missing
// setup is reported by the operation that needs it, so a GUI can show settings.
func New(opts Options) (*Manager, error) {
	opts.Provider = strings.ToLower(strings.TrimSpace(opts.Provider))
	if opts.Provider == "" {
		opts.Provider = "drive"
	}
	if _, err := Scopes(opts.Provider, opts.Resources); err != nil {
		return nil, err
	}
	resources, err := normalizeResources(opts.Resources)
	if err != nil {
		return nil, err
	}
	opts.Resources = resources
	if opts.DriveFolderID != "" {
		opts.DriveFolderID, err = folderID(opts.DriveFolderID)
		if err != nil {
			return nil, err
		}
	}
	if opts.OpenURL == nil {
		opts.OpenURL = auth.OpenURL
	}
	life, cancel := context.WithCancel(context.Background())
	m := &Manager{opts: opts, life: life, cancel: cancel, now: time.Now, gate: make(chan struct{}, 1), links: map[string]cachedLinks{}, public: signedClient(), persist: saveState}
	m.gate <- struct{}{}
	if opts.Session != nil && opts.Session.HTTPClient() != nil {
		// SDK debug response logging can include signed URLs. Supply an explicit
		// discard logger even when GOOGLE_SDK_GO_LOGGING_LEVEL enables logging.
		sdkOptions := []option.ClientOption{option.WithHTTPClient(opts.Session.HTTPClient()), option.WithLogger(slog.New(slog.NewTextHandler(io.Discard, nil)))}
		m.drive, err = drive.NewService(life, sdkOptions...)
		if err == nil {
			m.portability, err = dataportability.NewService(life, sdkOptions...)
		}
		if err != nil {
			cancel()
			return nil, errors.New("cannot configure Google API clients")
		}
	}
	return m, nil
}

func (m *Manager) acquire(ctx context.Context) (context.Context, func(), error) {
	if ctx == nil {
		return nil, nil, errors.New("source context is required")
	}
	// Bound individual HTTP calls, not the entire finite traversal: large
	// exports can legitimately need many metadata requests. The caller's
	// deadline and Close cancel the same context throughout the operation.
	ctx, cancel := context.WithCancel(ctx)
	stop := context.AfterFunc(m.life, cancel)
	select {
	case <-ctx.Done():
		stop()
		cancel()
		return nil, nil, ctx.Err()
	case <-m.gate:
		if err := ctx.Err(); err != nil {
			m.gate <- struct{}{}
			stop()
			cancel()
			return nil, nil, err
		}
		if m.life.Err() != nil {
			m.gate <- struct{}{}
			stop()
			cancel()
			return nil, nil, errors.New("archive source is closed")
		}
		return ctx, func() { m.gate <- struct{}{}; stop(); cancel() }, nil
	}
}

func (m *Manager) configured() error {
	if m.opts.Provider == "portability" && len(m.opts.Resources) == 0 {
		return errors.New("select explicit Data Portability resource groups before signing in or exporting")
	}
	if m.opts.Session == nil || m.HTTPClient() == nil {
		return errors.New("configure Google OAuth and sign in before reading archives")
	}
	return nil
}

func (m *Manager) Login(ctx context.Context) error {
	if err := m.configured(); err != nil {
		return err
	}
	return m.opts.Session.Login(ctx)
}
func (m *Manager) Refresh(ctx context.Context) error {
	if err := m.configured(); err != nil {
		return err
	}
	return m.opts.Session.Refresh(ctx)
}
func (m *Manager) Forget(ctx context.Context) error {
	if m.opts.Session == nil {
		return errors.New("Google OAuth session is not configured")
	}
	return m.opts.Session.Forget(ctx)
}
func (m *Manager) HTTPClient() *http.Client {
	if m.opts.Session == nil {
		return nil
	}
	return m.opts.Session.HTTPClient()
}
func (m *Manager) Close() error {
	m.closeOnce.Do(func() {
		m.cancel()
		<-m.gate
		defer func() { m.gate <- struct{}{} }()
		m.public.CloseIdleConnections()
		if m.opts.Session != nil {
			m.closeErr = m.opts.Session.Close()
		}
	})
	return m.closeErr
}

func (m *Manager) emit(message string) {
	if m.opts.Bus != nil {
		m.opts.Bus.Publish(events.Event{Kind: events.Status, Message: message})
	}
}

func (m *Manager) List(ctx context.Context) ([]model.Export, error) {
	ctx, release, err := m.acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer release()
	if err = m.configured(); err != nil {
		return nil, err
	}
	if m.opts.Provider == "drive" {
		return m.listDrive(ctx)
	}
	return m.listPortability(ctx)
}

func (m *Manager) RefreshLink(ctx context.Context, exportID string, index int) (string, error) {
	return m.refreshLink(ctx, exportID, index, "")
}

// RefreshRejectedLink lets a downloader identify the failed signature. A fresh
// job cache can then satisfy concurrent sibling or same-part requests without
// polling again or returning the same rejected URL.
func (m *Manager) RefreshRejectedLink(ctx context.Context, exportID string, index int, rejectedURL string) (string, error) {
	return m.refreshLink(ctx, exportID, index, rejectedURL)
}

func (m *Manager) refreshLink(ctx context.Context, exportID string, index int, rejectedURL string) (string, error) {
	ctx, release, err := m.acquire(ctx)
	if err != nil {
		return "", err
	}
	defer release()
	if err = m.configured(); err != nil {
		return "", err
	}
	if index < 1 {
		return "", errors.New("archive part index must be positive")
	}
	if m.opts.Provider == "drive" {
		return m.refreshDrive(ctx, exportID, index)
	}
	return m.refreshPortability(ctx, exportID, index, rejectedURL)
}

// TakeoutURL preserves the user's full product choice. Google's documented
// custom route names product IDs; it does not establish that a product-free
// /custom route supports destination/schedule parameters. The user selects Add
// to Drive and the recurring schedule on the canonical Takeout page normally.
func TakeoutURL() string {
	return "https://takeout.google.com/"
}

func (m *Manager) CreateExport(ctx context.Context) (string, error) {
	ctx, release, err := m.acquire(ctx)
	if err != nil {
		return "", err
	}
	defer release()
	if m.opts.Provider == "drive" {
		if err = ctx.Err(); err != nil {
			return "", err
		}
		if err = m.opts.OpenURL(TakeoutURL()); err != nil {
			return "", errors.New("cannot open Google Takeout; open the Takeout settings URL in your browser")
		}
		return "Review products in Google Takeout, select Add to Drive and exports every 2 months, then confirm Create export. Configure the resulting Drive folder ID for downloads.", nil
	}
	if err = m.configured(); err != nil {
		return "", err
	}
	return m.createPortability(ctx)
}
