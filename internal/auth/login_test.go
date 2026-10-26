package auth

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func sendCallback(t *testing.T, raw string, values url.Values) {
	t.Helper()
	u, err := url.Parse(raw)
	if err != nil {
		t.Error(err)
		return
	}
	u.RawQuery = values.Encode()
	client := &http.Client{Timeout: time.Second}
	r, err := client.Get(u.String())
	if err != nil {
		t.Error(err)
		return
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		t.Errorf("callback status: %d", r.StatusCode)
	}
}

func TestLoginDesktopPKCEAndKeyring(t *testing.T) {
	opts := testOptions(t)
	store := &fakeSecrets{}
	var authQuery url.Values
	var tokenCalls atomic.Int32
	var callbackURL string
	opts.OpenBrowser = func(raw string) error {
		u, err := url.Parse(raw)
		if err != nil {
			return err
		}
		authQuery = u.Query()
		if u.Scheme != "https" || u.Host != "accounts.google.com" || u.Path != "/o/oauth2/v2/auth" {
			t.Error("untrusted authorization URL")
		}
		for key, want := range map[string]string{"response_type": "code", "access_type": "offline", "code_challenge_method": "S256", "scope": testScope, "login_hint": opts.Email, "prompt": "consent", "include_granted_scopes": "false"} {
			if authQuery.Get(key) != want {
				t.Errorf("%s: %q", key, authQuery.Get(key))
			}
		}
		if authQuery.Has("client_secret") || authQuery.Has("code_verifier") || authQuery.Get("state") == "" {
			t.Error("unsafe or missing authorization parameters")
		}
		callbackURL = authQuery.Get("redirect_uri")
		redirect, _ := url.Parse(callbackURL)
		if redirect.Scheme != "http" || redirect.Hostname() != "127.0.0.1" || redirect.Port() == "" || redirect.Path != callbackPath {
			t.Error("callback is not ephemeral IPv4 loopback")
		}
		sendCallback(t, callbackURL, url.Values{"state": {authQuery.Get("state")}, "code": {"test-code"}, "scope": {testScope}, "authuser": {"0"}, "prompt": {"consent"}})
		return nil
	}
	m := testManager(t, opts, store, func(r *http.Request) (*http.Response, error) {
		tokenCalls.Add(1)
		if r.Method != "POST" || r.URL.String() != googleTokenURL {
			t.Error("unexpected exchange destination")
		}
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		if r.Form.Get("grant_type") != "authorization_code" || r.Form.Get("code") != "test-code" || r.Form.Get("redirect_uri") != callbackURL || r.Form.Get("client_secret") != "test-client-secret" {
			t.Error("incorrect library exchange form")
		}
		verifier := r.Form.Get("code_verifier")
		if len(verifier) < 43 || oauth2.S256ChallengeFromVerifier(verifier) != authQuery.Get("code_challenge") {
			t.Error("PKCE verifier/challenge mismatch")
		}
		if r.Header.Get("Authorization") != "" {
			t.Error("library probed auth style")
		}
		return response(r, 200, `{"access_token":"login-access","refresh_token":"login-refresh","expires_in":3600,"token_type":"Bearer","scope":"`+testScope+`"}`), nil
	})
	if err := m.Login(context.Background()); err != nil {
		t.Fatal(err)
	}
	if tokenCalls.Load() != 1 {
		t.Fatal("exchange not single shot")
	}
	stored, err := m.vault.load()
	if err != nil || stored.AccessToken != "login-access" || stored.RefreshToken != "login-refresh" {
		t.Fatalf("not stored: %v", err)
	}
	u, _ := url.Parse(callbackURL)
	connection, err := net.DialTimeout("tcp", u.Host, 200*time.Millisecond)
	if err == nil {
		_ = connection.Close()
		t.Fatal("callback listener survived Login")
	}
}

func TestLoginCancellationClosesCallbackAndAllowsQueuedOperation(t *testing.T) {
	opts := testOptions(t)
	opened := make(chan string, 1)
	opts.OpenBrowser = func(raw string) error { opened <- raw; return nil }
	m := testManager(t, opts, &fakeSecrets{}, nil)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- m.Login(ctx) }()
	raw := <-opened
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("Login did not cancel")
	}
	auth, _ := url.Parse(raw)
	callback, _ := url.Parse(auth.Query().Get("redirect_uri"))
	conn, err := net.DialTimeout("tcp", callback.Host, 200*time.Millisecond)
	if err == nil {
		_ = conn.Close()
		t.Fatal("canceled callback still accepts connections")
	}
	if err := m.Refresh(context.Background()); !errors.Is(err, ErrAuthorizationRequired) {
		t.Fatal(err)
	}
}

func TestCloseCancelsLoginAndIsIdempotent(t *testing.T) {
	opts := testOptions(t)
	opened := make(chan struct{})
	opts.OpenBrowser = func(string) error { close(opened); return nil }
	m := testManager(t, opts, &fakeSecrets{}, nil)
	done := make(chan error, 1)
	go func() { done <- m.Login(context.Background()) }()
	<-opened
	if err := m.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("Close did not cancel Login")
	}
	if err := m.Close(); err != nil {
		t.Fatal(err)
	}
	if err := m.Login(context.Background()); !errors.Is(err, errClosed) {
		t.Fatalf("closed Login: %v", err)
	}
}

func TestCancellationDoesNotWaitForBrowserLauncher(t *testing.T) {
	opts := testOptions(t)
	opened, release := make(chan struct{}), make(chan struct{})
	defer close(release)
	opts.OpenBrowser = func(string) error { close(opened); <-release; return nil }
	m := testManager(t, opts, &fakeSecrets{}, nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- m.Login(ctx) }()
	<-opened
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("launcher blocked cancellation")
	}
}

func TestLoginDenialAndBrowserErrorsAreSanitized(t *testing.T) {
	for _, denied := range []bool{false, true} {
		opts := testOptions(t)
		opts.OpenBrowser = func(raw string) error {
			if !denied {
				return errors.New("https://accounts.google.com/?code=secret-code&client_secret=hidden")
			}
			u, _ := url.Parse(raw)
			sendCallback(t, u.Query().Get("redirect_uri"), url.Values{"state": {u.Query().Get("state")}, "error": {"access_denied"}, "error_description": {"secret-code hidden"}, "error_uri": {"https://evil.example/?token=hidden"}})
			return nil
		}
		store := &fakeSecrets{}
		m := testManager(t, opts, store, nil)
		err := m.Login(context.Background())
		if err == nil || strings.Contains(err.Error(), "secret-code") || strings.Contains(err.Error(), "hidden") || strings.Contains(err.Error(), "accounts.google.com") {
			t.Fatalf("unsafe error: %v", err)
		}
		if denied && !errors.Is(err, ErrAuthorizationRequired) {
			t.Fatal(err)
		}
		if store.sets != 0 {
			t.Fatal("denied grant persisted")
		}
	}
}

func TestCallbackStrictValidationAndOneShot(t *testing.T) {
	valid := callbackPath + "?state=expected-state&code=accepted-code"
	for _, tc := range []struct{ name, method, target, host string }{
		{"wrong state", "GET", callbackPath + "?state=wrong&code=accepted-code", "127.0.0.1:9999"},
		{"missing state", "GET", callbackPath + "?code=accepted-code", "127.0.0.1:9999"},
		{"duplicate state", "GET", valid + "&state=expected-state", "127.0.0.1:9999"},
		{"duplicate code", "GET", valid + "&code=other", "127.0.0.1:9999"},
		{"missing code", "GET", callbackPath + "?state=expected-state", "127.0.0.1:9999"},
		{"code and denial", "GET", valid + "&error=access_denied", "127.0.0.1:9999"},
		{"empty code and denial", "GET", callbackPath + "?state=expected-state&code=&error=access_denied", "127.0.0.1:9999"},
		{"description with success", "GET", valid + "&error_description=secret", "127.0.0.1:9999"},
		{"malformed query", "GET", valid + "&unknown=%zz", "127.0.0.1:9999"},
		{"semicolon query", "GET", valid + ";extra=x", "127.0.0.1:9999"},
		{"unknown query", "GET", valid + "&access_token=secret", "127.0.0.1:9999"},
		{"foreign issuer", "GET", valid + "&iss=https%3A%2F%2Fevil.example", "127.0.0.1:9999"},
		{"encoded path", "GET", "/oauth2/%63allback?state=expected-state&code=accepted-code", "127.0.0.1:9999"},
		{"wrong path", "GET", "/wrong?state=expected-state&code=accepted-code", "127.0.0.1:9999"},
		{"wrong host", "GET", valid, "localhost:9999"},
		{"host prefix", "GET", valid, "127.0.0.1.evil.example:9999"},
		{"post", "POST", valid, "127.0.0.1:9999"},
		{"absolute target", "GET", "http://127.0.0.1:9999" + valid, "127.0.0.1:9999"},
		{"large query", "GET", valid + "&scope=" + strings.Repeat("a", maxCallbackQuery), "127.0.0.1:9999"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := &callbackHandler{host: "127.0.0.1:9999", state: "expected-state", results: make(chan callbackResult, 1)}
			r := httptest.NewRequest(tc.method, tc.target, nil)
			r.Host = tc.host
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code != 400 || h.used.Load() {
				t.Fatalf("invalid callback accepted: status %d", w.Code)
			}
			if strings.Contains(w.Body.String(), "accepted-code") || strings.Contains(w.Body.String(), "secret") {
				t.Fatal("callback reflected secret")
			}
			if w.Header().Get("Cache-Control") != "no-store" || w.Header().Get("Referrer-Policy") != "no-referrer" {
				t.Fatal("missing privacy headers")
			}
			good := httptest.NewRequest("GET", valid, nil)
			good.Host = h.host
			first := httptest.NewRecorder()
			h.ServeHTTP(first, good)
			if first.Code != 200 {
				t.Fatal("invalid callback consumed flow")
			}
			if result := <-h.results; result.code != "accepted-code" {
				t.Fatal("incorrect callback result")
			}
			again := httptest.NewRecorder()
			h.ServeHTTP(again, good)
			if again.Code != 410 {
				t.Fatal("callback reused")
			}
		})
	}
}

func TestConcurrentCallbacksAcceptExactlyOne(t *testing.T) {
	h := &callbackHandler{host: "127.0.0.1:9999", state: "state", results: make(chan callbackResult, 1)}
	var accepted atomic.Int32
	var wg sync.WaitGroup
	for range 20 {
		wg.Go(func() {
			r := httptest.NewRequest("GET", callbackPath+"?state=state&code=code", nil)
			r.Host = h.host
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code == 200 {
				accepted.Add(1)
			}
		})
	}
	wg.Wait()
	if accepted.Load() != 1 || len(h.results) != 1 {
		t.Fatal("callback was not one shot")
	}
}

func TestLoginRejectsIncompleteScopesWithoutSaving(t *testing.T) {
	opts := testOptions(t)
	opts.OpenBrowser = func(raw string) error {
		u, _ := url.Parse(raw)
		sendCallback(t, u.Query().Get("redirect_uri"), url.Values{"state": {u.Query().Get("state")}, "code": {"code"}})
		return nil
	}
	store := &fakeSecrets{}
	m := testManager(t, opts, store, func(r *http.Request) (*http.Response, error) {
		return response(r, 200, `{"access_token":"partial-access","expires_in":3600,"scope":"different"}`), nil
	})
	if err := m.Login(context.Background()); !errors.Is(err, ErrAuthorizationRequired) {
		t.Fatal(err)
	}
	if store.sets != 0 {
		t.Fatal("incomplete grant saved")
	}
}

func TestTokenEndpointRedirectIsNotFollowed(t *testing.T) {
	var calls int
	m := testManager(t, testOptions(t), &fakeSecrets{}, func(r *http.Request) (*http.Response, error) {
		calls++
		res := response(r, 307, `redirect`)
		res.Header.Set("Location", "https://evil.example/?token=hidden")
		return res, nil
	})
	seedToken(t, m, validToken())
	err := m.Refresh(context.Background())
	if err == nil || calls != 1 || strings.Contains(err.Error(), "hidden") || strings.Contains(err.Error(), "evil.example") {
		t.Fatalf("unsafe redirect handling: calls=%d error=%v", calls, err)
	}
}

func TestCallbackResponseHasNoCredentialData(t *testing.T) {
	h := &callbackHandler{host: "127.0.0.1:8888", state: "private-state", results: make(chan callbackResult, 1)}
	r := httptest.NewRequest("GET", callbackPath+"?state=private-state&code=private-code", nil)
	r.Host = h.host
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	body, _ := io.ReadAll(w.Result().Body)
	if strings.Contains(string(body), "private-") {
		t.Fatal("private data reflected in browser")
	}
}
