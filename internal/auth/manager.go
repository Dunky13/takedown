// Package auth implements Google's documented desktop OAuth flow. Only Login
// opens the user's system browser; background requests never start a login.
package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Dunky13/takedown/internal/events"
	"golang.org/x/oauth2"
)

// ErrAuthorizationRequired means explicit authorization or OAuth configuration
// is needed. Use errors.Is, including on HTTP request errors.
var ErrAuthorizationRequired = errors.New("Google authorization required")

var errClosed = errors.New("authentication manager is closed")

const tokenTimeout = 30 * time.Second

type Options struct {
	DataDir, ClientFile, Email string
	Scopes                     []string
	Bus                        *events.Bus
	OpenBrowser                func(string) error
}

type Manager struct {
	opts      Options
	store     secretStore
	base      http.RoundTripper
	client    *http.Client
	tokenHTTP *http.Client
	life      context.Context
	cancel    context.CancelFunc
	gate      chan struct{}
	closeOnce sync.Once
	revision  atomic.Uint64
	config    *oauth2.Config // all fields below are protected by gate
	vault     *vault
	token     *oauth2.Token
	loaded    bool
	pending   bool // persist a rotation before using the new access token
	required  bool
}

// New performs no filesystem, keyring, network, or browser operations. Empty
// ClientFile and Scopes are permitted for settings/setup; operational use fails
// with ErrAuthorizationRequired until configured.
func New(opts Options) (*Manager, error) { return newManager(opts, osSecrets{}, nil) }

func newManager(opts Options, store secretStore, base http.RoundTripper) (*Manager, error) {
	if strings.TrimSpace(opts.DataDir) == "" {
		return nil, errors.New("authentication data directory is required")
	}
	dir, err := filepath.Abs(opts.DataDir)
	if err != nil {
		return nil, errors.New("invalid authentication data directory")
	}
	opts.DataDir = dir
	opts.Email = strings.ToLower(strings.TrimSpace(opts.Email))
	if strings.ContainsAny(opts.Email, "\r\n\x00") {
		return nil, errors.New("invalid Google login hint")
	}
	opts.Scopes, err = normalizeScopes(opts.Scopes)
	if err != nil {
		return nil, err
	}
	if opts.OpenBrowser == nil {
		opts.OpenBrowser = OpenURL
	}
	if base == nil {
		if transport, ok := http.DefaultTransport.(*http.Transport); ok {
			base = transport.Clone()
		} else {
			base = http.DefaultTransport
		}
	}
	life, cancel := context.WithCancel(context.Background())
	m := &Manager{opts: opts, store: store, base: base, life: life, cancel: cancel, gate: make(chan struct{}, 1)}
	m.gate <- struct{}{}
	m.tokenHTTP = &http.Client{Transport: &grantTransport{base: base, scopes: opts.Scopes}, Timeout: tokenTimeout, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	m.client = &http.Client{Transport: &apiTransport{manager: m}, CheckRedirect: checkRedirect}
	return m, nil
}

func (m *Manager) operation(ctx context.Context) (context.Context, func(), error) {
	if ctx == nil {
		return nil, nil, errors.New("authentication context is required")
	}
	if m.life.Err() != nil {
		return nil, nil, errClosed
	}
	op, cancel := context.WithCancel(ctx)
	stop := context.AfterFunc(m.life, cancel)
	return op, func() { stop(); cancel() }, nil
}

func (m *Manager) acquire(ctx context.Context) (context.Context, func(), error) {
	op, cleanup, err := m.operation(ctx)
	if err != nil {
		return nil, nil, err
	}
	select {
	case <-op.Done():
		cleanup()
		return nil, nil, op.Err()
	case <-m.gate:
		if op.Err() != nil || m.life.Err() != nil {
			m.gate <- struct{}{}
			cleanup()
			if m.life.Err() != nil {
				return nil, nil, errClosed
			}
			return nil, nil, op.Err()
		}
		return op, func() { cleanup(); m.gate <- struct{}{} }, nil
	}
}

func (m *Manager) emit(kind events.Kind, message string) {
	if m.opts.Bus != nil {
		m.opts.Bus.Publish(events.Event{Kind: kind, Message: message})
	}
}

func (m *Manager) needAuth(message string) error {
	m.emit(events.ReAuthRequired, message)
	return fmt.Errorf("%w: %s", ErrAuthorizationRequired, message)
}

func (m *Manager) ensureConfig() error {
	if len(m.opts.Scopes) == 0 {
		return m.needAuth("select at least one required Google API scope before signing in")
	}
	return m.ensureClient()
}

func (m *Manager) ensureClient() error {
	if m.config != nil {
		return nil
	}
	if strings.TrimSpace(m.opts.ClientFile) == "" {
		return m.needAuth("configure a Google OAuth Desktop app client JSON file, then run Login")
	}
	config, err := loadClient(m.opts.ClientFile, m.opts.Scopes)
	if err != nil {
		return m.needAuth(err.Error())
	}
	m.config = config
	m.vault = newVault(m.opts, config.ClientID, m.store)
	return nil
}

// ProfileID identifies a local profile slot, not the Google account that gave
// consent. It excludes scopes so resource histories can share the same slot.
// It reads the validated desktop client lazily but never accesses the keyring.
func (m *Manager) ProfileID(ctx context.Context) (string, error) {
	_, release, err := m.acquire(ctx)
	if err != nil {
		return "", err
	}
	defer release()
	if err := m.ensureClient(); err != nil {
		return "", err
	}
	binding, _ := json.Marshal([]string{m.opts.DataDir, m.config.ClientID, m.opts.Email})
	return fingerprint(string(binding)), nil
}

// Refresh forces a library token refresh. Concurrent callers already waiting
// on the same refresh reuse its result. HTTP expiry uses the same locked path.
func (m *Manager) Refresh(ctx context.Context) error {
	revision := m.revision.Load()
	ctx, release, err := m.acquire(ctx)
	if err != nil {
		return err
	}
	defer release()
	_, err = m.tokenLocked(ctx, m.revision.Load() == revision)
	return err
}

func (m *Manager) tokenLocked(ctx context.Context, force bool) (*oauth2.Token, error) {
	if err := m.ensureConfig(); err != nil {
		return nil, err
	}
	if m.required {
		return nil, m.needAuth("run Login to authorize the required Google API scopes again")
	}
	if !m.loaded {
		token, err := m.vault.load()
		if err != nil {
			if errors.Is(err, ErrAuthorizationRequired) {
				return nil, m.needAuth("no usable OAuth token is stored; run Login")
			}
			return nil, err
		}
		m.token, m.loaded = token, true
	}
	if m.pending {
		if err := m.vault.save(m.token); err != nil {
			return nil, err
		}
		m.pending = false
		m.revision.Add(1)
		force = false
	}
	if m.token == nil {
		return nil, m.needAuth("no OAuth token is stored; run Login")
	}
	if !force && m.token.Valid() {
		return m.token, nil
	}
	if m.token.RefreshToken == "" {
		return nil, m.needAuth("the OAuth grant has no refresh token; run Login for a new grant")
	}
	refreshCtx, cancel := context.WithTimeout(ctx, tokenTimeout)
	defer cancel()
	refreshCtx = context.WithValue(refreshCtx, oauth2.HTTPClient, m.tokenHTTP)
	seed := *m.token
	// The library owns refresh/rotation. Expire a copy for explicit Refresh.
	seed.Expiry = time.Unix(1, 0)
	token, err := m.config.TokenSource(refreshCtx, &seed).Token()
	if err != nil {
		return nil, m.oauthError(refreshCtx, err)
	}
	if err := refreshCtx.Err(); err != nil {
		return nil, err
	}
	if err := m.validateToken(token); err != nil {
		m.required = true
		return nil, err
	}
	if err := m.acceptToken(token); err != nil {
		return nil, err
	}
	return m.token, nil
}

func (m *Manager) validateToken(token *oauth2.Token) error {
	if token == nil || !strings.EqualFold(token.Type(), "Bearer") || token.Expiry.IsZero() || !token.Valid() {
		return m.needAuth("Google returned an unusable or expired OAuth token; run Login again")
	}
	if err := validateGrant(token, m.opts.Scopes); err != nil {
		return m.needAuth("Google did not grant all required API scopes; run Login and grant every selected scope")
	}
	return nil
}

func (m *Manager) acceptToken(token *oauth2.Token) error {
	m.token, m.loaded, m.pending, m.required = token, true, true, false
	if err := m.vault.save(token); err != nil {
		return err
	}
	m.pending = false
	m.revision.Add(1)
	return nil
}

func (m *Manager) oauthError(ctx context.Context, err error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	var response *oauth2.RetrieveError
	if errors.As(err, &response) {
		switch response.ErrorCode {
		case "temporarily_unavailable", "server_error":
			var delay time.Duration
			var timeout bool
			if response.Response != nil {
				delay = retryAfter(response.Response.Header.Get("Retry-After"))
				timeout = response.Response.StatusCode == http.StatusRequestTimeout
			}
			return &temporaryError{timeout: timeout, retryAfter: delay}
		case "invalid_grant", "invalid_token", "access_denied", "insufficient_scope":
			m.required = true
			return m.needAuth("Google rejected or revoked the OAuth grant; run Login again")
		case "invalid_client", "unauthorized_client", "invalid_scope", "invalid_request", "redirect_uri_mismatch", "unsupported_grant_type":
			m.required = true
			return m.needAuth("Google rejected the OAuth client; configure a valid Desktop app client JSON file")
		}
		if response.Response != nil {
			status := response.Response.StatusCode
			if status == http.StatusTooManyRequests || status == http.StatusRequestTimeout || status >= 500 && status <= 599 {
				return &temporaryError{timeout: status == http.StatusRequestTimeout, retryAfter: retryAfter(response.Response.Header.Get("Retry-After"))}
			}
			if status == 400 || status == 401 || status == 403 {
				m.required = true
				return m.needAuth("Google rejected the OAuth request; verify the Desktop app client and scopes, then run Login")
			}
		}
	}
	if errors.Is(err, ErrAuthorizationRequired) {
		m.required = true
		return m.needAuth("Google did not grant all required API scopes; run Login and grant every selected scope")
	}
	if temporary := temporaryNetwork(err); temporary != nil {
		return temporary
	}
	// Provider/transport errors may contain secret bodies, URLs, or form fields.
	return errors.New("Google OAuth token request failed; check the connection and client configuration, then retry")
}

// HTTPClient is concurrency safe, has no cookie jar, and sends bearer tokens
// only to explicit API hosts. Treat its fields as read-only.
func (m *Manager) HTTPClient() *http.Client { return m.client }

// Forget deletes the configured OAuth slot and current hint's legacy secrets.
// It neither revokes Google's remote grant nor enumerates other profiles.
func (m *Manager) Forget(ctx context.Context) error {
	ctx, release, err := m.acquire(ctx)
	if err != nil {
		return err
	}
	defer release()
	m.token, m.loaded, m.pending, m.required = nil, true, false, true
	m.revision.Add(1)
	var errs []error
	if m.config == nil && m.opts.ClientFile != "" && len(m.opts.Scopes) > 0 {
		config, loadErr := loadClient(m.opts.ClientFile, m.opts.Scopes)
		if loadErr != nil {
			errs = append(errs, errors.New("cannot identify the OAuth token slot: restore the configured Desktop app client JSON file and retry Forget"))
		} else {
			m.config = config
			m.vault = newVault(m.opts, config.ClientID, m.store)
		}
	}
	if m.vault != nil {
		errs = append(errs, m.vault.forget())
	}
	if ctx.Err() != nil {
		return errors.Join(append(errs, ctx.Err())...)
	}
	errs = append(errs, forgetLegacy(m.opts, m.store))
	if err := errors.Join(errs...); err != nil {
		return err
	}
	m.emit(events.Status, "Stored authorization for the current profile was forgotten")
	return nil
}

// Close cancels active operations and clears cached credentials. It is
// idempotent and leaves the parent-owned event bus open.
func (m *Manager) Close() error {
	m.closeOnce.Do(func() {
		m.cancel()
		<-m.gate
		defer func() { m.gate <- struct{}{} }()
		m.token = nil
		if m.config != nil {
			m.config.ClientSecret = ""
		}
		if closer, ok := m.base.(interface{ CloseIdleConnections() }); ok {
			closer.CloseIdleConnections()
		}
	})
	return nil
}
