package auth

import (
	"context"
	"crypto/subtle"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Dunky13/takedown/internal/events"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
)

const (
	loginTimeout     = 5 * time.Minute
	callbackPath     = "/oauth2/callback"
	maxCallbackQuery = 16 << 10
)

// OpenURL uses the normal system browser; it does not inspect or control it.
// Caller-supplied URLs and opener diagnostics are never included in errors.
func OpenURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil || !secureURL(u) {
		return errors.New("system browser requires an HTTPS URL without embedded credentials or a nonstandard port")
	}
	if browser.OpenURL(raw) != nil {
		return errors.New("cannot open the system browser; configure a default browser and retry Login")
	}
	return nil
}

// Login performs one bounded, cancellable desktop OAuth flow. All sign-in,
// passkey/MFA and consent interactions remain in the user's ordinary browser.
func (m *Manager) Login(ctx context.Context) error {
	ctx, release, err := m.acquire(ctx)
	if err != nil {
		return err
	}
	defer release()
	if err := m.ensureConfig(); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, loginTimeout)
	defer cancel()
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return errors.New("cannot bind the Google OAuth loopback callback to 127.0.0.1")
	}
	state, verifier := oauth2.GenerateVerifier(), oauth2.GenerateVerifier()
	callback := &callbackHandler{host: listener.Addr().String(), state: state, results: make(chan callbackResult, 1)}
	server := &http.Server{
		Handler: callback, ReadHeaderTimeout: 5 * time.Second, WriteTimeout: 5 * time.Second,
		IdleTimeout: 5 * time.Second, MaxHeaderBytes: 8 << 10,
		ErrorLog:    log.New(io.Discard, "", 0),
		BaseContext: func(net.Listener) context.Context { return ctx },
	}
	served := make(chan error, 1)
	go func() { served <- server.Serve(listener) }()
	defer func() { _ = server.Close(); <-served }()
	config := *m.config
	config.RedirectURL = "http://" + callback.host + callbackPath
	options := []oauth2.AuthCodeOption{
		oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier),
		oauth2.SetAuthURLParam("prompt", "consent"),
		oauth2.SetAuthURLParam("include_granted_scopes", "false"),
	}
	if m.opts.Email != "" {
		options = append(options, oauth2.SetAuthURLParam("login_hint", m.opts.Email))
	}
	authorizationURL := config.AuthCodeURL(state, options...)
	if len(authorizationURL) > 32<<10 {
		return errors.New("Google OAuth authorization request is too large; select fewer resource scopes")
	}
	m.emit(events.Status, "Complete Google sign-in and consent in your system browser; cancel to stop waiting")
	opened := make(chan error, 1)
	// OpenBrowser has no context parameter. Do not let an OS launcher that
	// remains attached to the browser block cancellation of the OAuth flow.
	go func() { opened <- m.opts.OpenBrowser(authorizationURL) }()
	var result callbackResult
selectLoop:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-served:
			// Keep the deferred teardown's single receive valid.
			served <- http.ErrServerClosed
			return errors.New("Google OAuth callback server stopped before authorization completed")
		case openErr := <-opened:
			if openErr != nil {
				return errors.New("cannot open the system browser; configure a default browser and retry Login")
			}
			opened = nil
		case result = <-callback.results:
			break selectLoop
		}
	}
	// Stop accepting callbacks before exchanging the single accepted code.
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second)
	_ = server.Shutdown(shutdownCtx)
	shutdownCancel()
	if result.denied {
		return m.needAuth("Google authorization was denied; run Login when ready to grant the required scopes")
	}
	exchangeCtx, exchangeCancel := context.WithTimeout(ctx, tokenTimeout)
	defer exchangeCancel()
	exchangeCtx = context.WithValue(exchangeCtx, oauth2.HTTPClient, m.tokenHTTP)
	token, err := config.Exchange(exchangeCtx, result.code, oauth2.VerifierOption(verifier))
	if err != nil {
		return m.oauthError(exchangeCtx, err)
	}
	if err := exchangeCtx.Err(); err != nil {
		return err
	}
	if err := m.validateToken(token); err != nil {
		return err
	}
	if err := m.acceptToken(token); err != nil {
		return err
	}
	m.emit(events.Status, "Google OAuth authorization stored securely in the OS keyring")
	return nil
}

type callbackResult struct {
	code   string
	denied bool
}

type callbackHandler struct {
	host, state string
	used        atomic.Bool
	results     chan callbackResult
}

func (h *callbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if r.Method != http.MethodGet || r.Host != h.host || r.URL.IsAbs() || r.URL.Host != "" || r.URL.Path != callbackPath || r.URL.RawPath != "" || len(r.URL.RawQuery) > maxCallbackQuery {
		http.Error(w, "Invalid OAuth callback.", http.StatusBadRequest)
		return
	}
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || !validCallbackQuery(query) || subtle.ConstantTimeCompare([]byte(query.Get("state")), []byte(h.state)) != 1 {
		http.Error(w, "Invalid OAuth callback.", http.StatusBadRequest)
		return
	}
	code, failure := query.Get("code"), query.Get("error")
	if query.Has("code") == query.Has("error") || (code == "") == (failure == "") || (code != "" && (strings.ContainsAny(code, " \t\r\n\x00") || query.Has("error_description") || query.Has("error_uri"))) || (query.Has("iss") && query.Get("iss") != "https://accounts.google.com") {
		http.Error(w, "Invalid OAuth callback.", http.StatusBadRequest)
		return
	}
	if !h.used.CompareAndSwap(false, true) {
		http.Error(w, "OAuth callback already received.", http.StatusGone)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = io.WriteString(w, "Authorization response received. You can return to the application and close this tab.")
	h.results <- callbackResult{code: code, denied: failure != ""}
}

func validCallbackQuery(query url.Values) bool {
	for key, values := range query {
		if len(values) != 1 {
			return false
		}
		switch key {
		case "state", "code", "error", "error_description", "error_uri", "scope", "authuser", "prompt", "hd", "iss":
		default:
			return false
		}
	}
	return query.Get("state") != ""
}
