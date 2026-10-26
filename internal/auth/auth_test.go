package auth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Dunky13/takedown/internal/events"
	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

const testScope = "https://www.googleapis.com/auth/dataportability.myactivity.search"

type fakeSecrets struct {
	mu                        sync.Mutex
	values                    map[string]string
	gets, sets, deletes       int
	getErr, setErr, deleteErr error
}

func (s *fakeSecrets) Get(service, account string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.gets++
	if s.getErr != nil {
		return "", s.getErr
	}
	value, ok := s.values[service+":"+account]
	if !ok {
		return "", keyring.ErrNotFound
	}
	return value, nil
}
func (s *fakeSecrets) Set(service, account, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sets++
	if s.setErr != nil {
		return s.setErr
	}
	if s.values == nil {
		s.values = make(map[string]string)
	}
	s.values[service+":"+account] = value
	return nil
}
func (s *fakeSecrets) Delete(service, account string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.deletes++
	if s.deleteErr != nil {
		return s.deleteErr
	}
	delete(s.values, service+":"+account)
	return nil
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func response(r *http.Request, status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}
}

func clientDocument() map[string]any {
	return map[string]any{"installed": map[string]any{
		"client_id": "123-test.apps.googleusercontent.com", "client_secret": "test-client-secret",
		"auth_uri": googleAuthURL, "token_uri": googleTokenURL, "redirect_uris": []string{"http://localhost"},
	}}
}

func writeClient(t *testing.T, document map[string]any) string {
	t.Helper()
	data, err := json.Marshal(document)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "client.json")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func testOptions(t *testing.T) Options {
	t.Helper()
	return Options{DataDir: t.TempDir(), ClientFile: writeClient(t, clientDocument()), Email: "hint@example.test", Scopes: []string{testScope}, OpenBrowser: func(string) error { t.Error("unexpected browser launch"); return errors.New("unexpected browser") }}
}

func testManager(t *testing.T, opts Options, store *fakeSecrets, base roundTripFunc) *Manager {
	t.Helper()
	if base == nil {
		base = func(*http.Request) (*http.Response, error) {
			t.Error("unexpected HTTP request")
			return nil, errors.New("unexpected request")
		}
	}
	m, err := newManager(opts, store, base)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = m.Close() })
	return m
}

func seedToken(t *testing.T, m *Manager, token *oauth2.Token) {
	t.Helper()
	if err := m.ensureConfig(); err != nil {
		t.Fatal(err)
	}
	if err := m.vault.save(token); err != nil {
		t.Fatal(err)
	}
}

func validToken() *oauth2.Token {
	return &oauth2.Token{AccessToken: "test-access-token", RefreshToken: "test-refresh-token", TokenType: "Bearer", Expiry: time.Now().Add(time.Hour)}
}

func TestNewIsLazyAndAllowsUnconfiguredScopes(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "uncreated")
	store := &fakeSecrets{getErr: errors.New("must not read"), setErr: errors.New("must not write")}
	m := testManager(t, Options{DataDir: dir}, store, nil)
	if store.gets+store.sets+store.deletes != 0 {
		t.Fatal("New accessed keyring")
	}
	if _, err := os.Stat(dir); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("New created data directory")
	}
	if m.client.Jar != nil {
		t.Fatal("cookie jar enabled")
	}
	for _, operation := range []func(context.Context) error{m.Login, m.Refresh} {
		if err := operation(context.Background()); !errors.Is(err, ErrAuthorizationRequired) || !strings.Contains(err.Error(), "scope") {
			t.Fatalf("missing scopes: %v", err)
		}
	}
	_, err := m.client.Get("https://www.googleapis.com/drive/v3/files")
	if !errors.Is(err, ErrAuthorizationRequired) {
		t.Fatalf("HTTP missing scope: %v", err)
	}
	if store.gets+store.sets+store.deletes != 0 {
		t.Fatal("unconfigured operation accessed keyring")
	}
}

func TestMissingClientIsActionableAndPublishes(t *testing.T) {
	opts := testOptions(t)
	opts.ClientFile = ""
	opts.Bus = events.New()
	eventsCh, stop := opts.Bus.Subscribe(8)
	defer stop()
	m := testManager(t, opts, &fakeSecrets{}, nil)
	for _, op := range []func(context.Context) error{m.Login, m.Refresh} {
		if err := op(context.Background()); !errors.Is(err, ErrAuthorizationRequired) || !strings.Contains(err.Error(), "Desktop app") {
			t.Fatalf("missing client: %v", err)
		}
	}
	if event := <-eventsCh; event.Kind != events.ReAuthRequired {
		t.Fatalf("event: %+v", event)
	}
	_ = m.Close()
	opts.Bus.Publish(events.Event{Kind: events.Status, Message: "parent still owns bus"})
}

func TestClientValidation(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func(map[string]any)
	}{
		{"web", func(d map[string]any) { d["web"] = d["installed"]; delete(d, "installed") }},
		{"both", func(d map[string]any) { d["web"] = d["installed"] }},
		{"service account", func(d map[string]any) { d["type"] = "service_account" }},
		{"authorization rewrite", func(d map[string]any) {
			d["installed"].(map[string]any)["auth_uri"] = "https://evil.example/auth?secret=hidden"
		}},
		{"token rewrite", func(d map[string]any) { d["installed"].(map[string]any)["token_uri"] = "https://evil.example/token" }},
		{"token query", func(d map[string]any) {
			d["installed"].(map[string]any)["token_uri"] = googleTokenURL + "?secret=hidden"
		}},
		{"token port", func(d map[string]any) {
			d["installed"].(map[string]any)["token_uri"] = "https://oauth2.googleapis.com:444/token"
		}},
		{"missing id", func(d map[string]any) { delete(d["installed"].(map[string]any), "client_id") }},
		{"missing secret", func(d map[string]any) { delete(d["installed"].(map[string]any), "client_secret") }},
		{"missing redirect", func(d map[string]any) { delete(d["installed"].(map[string]any), "redirect_uris") }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			doc := clientDocument()
			tc.mutate(doc)
			_, err := loadClient(writeClient(t, doc), []string{testScope})
			if err == nil {
				t.Fatal("invalid client accepted")
			}
			if strings.Contains(err.Error(), "hidden") || strings.Contains(err.Error(), "evil.example") || strings.Contains(err.Error(), "test-client-secret") {
				t.Fatalf("unsafe error: %v", err)
			}
		})
	}
	for _, authURL := range []string{googleAuthURL, "https://accounts.google.com/o/oauth2/auth"} {
		doc := clientDocument()
		doc["installed"].(map[string]any)["auth_uri"] = authURL
		config, err := loadClient(writeClient(t, doc), []string{testScope})
		if err != nil || config.Endpoint.AuthStyle != oauth2.AuthStyleInParams {
			t.Fatalf("official client: %v", err)
		}
	}
	if _, err := loadClient(filepath.Join(t.TempDir(), "missing"), nil); err == nil {
		t.Fatal("missing file accepted")
	}
	path := filepath.Join(t.TempDir(), "invalid.json")
	_ = os.WriteFile(path, []byte(`{"installed": null}`), 0600)
	if _, err := loadClient(path, nil); err == nil {
		t.Fatal("null client accepted")
	}
}

func TestScopesAndNamespace(t *testing.T) {
	for _, scopes := range [][]string{{testScope, "openid"}, {testScope, "https://www.googleapis.com/auth/drive.readonly"}, {"scope with spaces"}, {""}} {
		if _, err := New(Options{DataDir: t.TempDir(), Scopes: scopes}); err == nil {
			t.Fatalf("accepted scopes %q", scopes)
		}
	}
	opts := testOptions(t)
	opts.Scopes = []string{testScope + ".z", testScope, testScope}
	m := testManager(t, opts, &fakeSecrets{}, nil)
	if len(m.opts.Scopes) != 2 || m.opts.Scopes[0] != testScope {
		t.Fatal("scopes not sorted/deduplicated")
	}
	base := newVault(m.opts, "client-a", &fakeSecrets{}).account
	for _, modify := range []func(*Options){func(o *Options) { o.Email = "other@example.test" }, func(o *Options) { o.DataDir += "/other" }, func(o *Options) { o.Scopes = []string{testScope} }} {
		other := m.opts
		modify(&other)
		if newVault(other, "client-a", &fakeSecrets{}).account == base {
			t.Fatal("token namespace collision")
		}
	}
	if newVault(m.opts, "client-b", &fakeSecrets{}).account == base {
		t.Fatal("client namespace collision")
	}
}

func TestScopeGrantValidation(t *testing.T) {
	for _, tc := range []struct {
		name  string
		extra map[string]any
		want  bool
	}{
		{"omitted", nil, true}, {"complete", map[string]any{"scope": testScope}, true},
		{"partial", map[string]any{"scope": testScope + ".other"}, false}, {"empty", map[string]any{"scope": ""}, false},
		{"wrong type", map[string]any{"scope": 42}, false}, {"mixed", map[string]any{"scope": testScope + " openid"}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			token := validToken().WithExtra(tc.extra)
			if (validateGrant(token, []string{testScope}) == nil) != tc.want {
				t.Fatal("scope grant acceptance mismatch")
			}
		})
	}
}

func TestMissingTokenNeverOpensBrowser(t *testing.T) {
	store := &fakeSecrets{}
	m := testManager(t, testOptions(t), store, nil)
	_, err := m.client.Get("https://dataportability.googleapis.com/v1/portabilityArchive:initiate")
	if !errors.Is(err, ErrAuthorizationRequired) {
		t.Fatalf("missing token: %v", err)
	}
	if err := m.Refresh(context.Background()); !errors.Is(err, ErrAuthorizationRequired) {
		t.Fatalf("refresh missing token: %v", err)
	}
}

func TestConcurrentHTTPRefreshRotatesOnce(t *testing.T) {
	store := &fakeSecrets{}
	var refreshes, apiCalls atomic.Int32
	m := testManager(t, testOptions(t), store, func(r *http.Request) (*http.Response, error) {
		if r.URL.String() == googleTokenURL {
			refreshes.Add(1)
			if err := r.ParseForm(); err != nil {
				t.Error(err)
			}
			if r.Form.Get("refresh_token") != "test-refresh-token" || r.Form.Get("grant_type") != "refresh_token" || r.Form.Get("client_secret") != "test-client-secret" {
				t.Error("incorrect library refresh request")
			}
			return response(r, 200, `{"access_token":"rotated-access","refresh_token":"rotated-refresh","token_type":"Bearer","expires_in":3600,"scope":"`+testScope+`"}`), nil
		}
		apiCalls.Add(1)
		if r.Header.Get("Authorization") != "Bearer rotated-access" {
			t.Error("incorrect bearer token")
		}
		return response(r, 200, `{}`), nil
	})
	token := validToken()
	token.Expiry = time.Now().Add(-time.Hour)
	seedToken(t, m, token)
	var wg sync.WaitGroup
	for range 24 {
		wg.Go(func() {
			r, err := m.client.Get("https://www.googleapis.com/drive/v3/files")
			if err != nil {
				t.Error(err)
				return
			}
			_ = r.Body.Close()
		})
	}
	wg.Wait()
	if refreshes.Load() != 1 || apiCalls.Load() != 24 {
		t.Fatalf("refreshes=%d API=%d", refreshes.Load(), apiCalls.Load())
	}
	stored, err := m.vault.load()
	if err != nil || stored.RefreshToken != "rotated-refresh" {
		t.Fatalf("rotation not persisted: %v", err)
	}
}

func TestRefreshRetainsOmittedRefreshToken(t *testing.T) {
	m := testManager(t, testOptions(t), &fakeSecrets{}, func(r *http.Request) (*http.Response, error) {
		return response(r, 200, `{"access_token":"new","token_type":"Bearer","expires_in":3600}`), nil
	})
	seedToken(t, m, validToken())
	if err := m.Refresh(context.Background()); err != nil {
		t.Fatal(err)
	}
	stored, err := m.vault.load()
	if err != nil || stored.RefreshToken != "test-refresh-token" {
		t.Fatalf("refresh token dropped: %v", err)
	}
}

func TestRefreshScopeAndExpiryFailuresRequireLogin(t *testing.T) {
	for _, body := range []string{
		`{"access_token":"a","expires_in":3600,"scope":""}`,
		`{"access_token":"a","expires_in":3600,"scope":"ungranted"}`,
		`{"access_token":"a","expires_in":3600,"scope":null}`,
		`{"access_token":"a","expires_in":-1}`,
		`{"access_token":"a"}`,
		`{"access_token":"a","token_type":"MAC","expires_in":3600}`,
	} {
		var calls int
		m := testManager(t, testOptions(t), &fakeSecrets{}, func(r *http.Request) (*http.Response, error) { calls++; return response(r, 200, body), nil })
		seedToken(t, m, validToken())
		if err := m.Refresh(context.Background()); !errors.Is(err, ErrAuthorizationRequired) {
			t.Fatalf("expected authorization required: %v", err)
		}
		if err := m.Refresh(context.Background()); !errors.Is(err, ErrAuthorizationRequired) || calls != 1 {
			t.Fatal("rejected grant retried without explicit Login")
		}
	}
}

func TestOAuthErrorSanitizationAndNoAutomaticBrowser(t *testing.T) {
	for _, contentType := range []string{"application/json", "application/x-www-form-urlencoded"} {
		m := testManager(t, testOptions(t), &fakeSecrets{}, func(r *http.Request) (*http.Response, error) {
			body := `{"error":"invalid_grant","error_description":"secret-access code-secret https://evil.example/?token=secret"}`
			if contentType != "application/json" {
				body = "error=invalid_grant&error_description=secret-access+code-secret+https%3A%2F%2Fevil.example%2F%3Ftoken%3Dsecret"
			}
			res := response(r, 400, body)
			res.Header.Set("Content-Type", contentType)
			return res, nil
		})
		seedToken(t, m, validToken())
		err := m.Refresh(context.Background())
		if !errors.Is(err, ErrAuthorizationRequired) || strings.Contains(err.Error(), "secret") || strings.Contains(err.Error(), "evil.example") {
			t.Fatalf("unsafe or wrong error: %v", err)
		}
	}
}

func TestRotationPersistenceFailureBlocksBearerAndRetriesSave(t *testing.T) {
	store := &fakeSecrets{}
	var refreshes, apiCalls int
	m := testManager(t, testOptions(t), store, func(r *http.Request) (*http.Response, error) {
		if r.URL.String() == googleTokenURL {
			refreshes++
			return response(r, 200, `{"access_token":"new-access","refresh_token":"new-refresh","expires_in":3600}`), nil
		}
		apiCalls++
		return response(r, 200, `{}`), nil
	})
	token := validToken()
	token.Expiry = time.Now().Add(-time.Hour)
	seedToken(t, m, token)
	store.setErr = errors.New("sensitive keyring error")
	if _, err := m.client.Get("https://www.googleapis.com/drive/v3/files"); err == nil || strings.Contains(err.Error(), "sensitive") {
		t.Fatalf("save failure: %v", err)
	}
	if apiCalls != 0 {
		t.Fatal("bearer sent before secure persistence")
	}
	store.setErr = nil
	r, err := m.client.Get("https://www.googleapis.com/drive/v3/files")
	if err != nil {
		t.Fatal(err)
	}
	_ = r.Body.Close()
	if refreshes != 1 || apiCalls != 1 {
		t.Fatal("did not retry persistence before refreshing again")
	}
	stored, err := m.vault.load()
	if err != nil || stored.RefreshToken != "new-refresh" {
		t.Fatal("lost rotated refresh token")
	}
}

func TestCancellationInterruptsTokenExchangeAndGateWait(t *testing.T) {
	entered := make(chan struct{})
	m := testManager(t, testOptions(t), &fakeSecrets{}, func(r *http.Request) (*http.Response, error) {
		close(entered)
		<-r.Context().Done()
		return nil, r.Context().Err()
	})
	seedToken(t, m, validToken())
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- m.Refresh(ctx) }()
	<-entered
	queued, stop := context.WithCancel(context.Background())
	stop()
	if err := m.Refresh(queued); !errors.Is(err, context.Canceled) {
		t.Fatalf("queued cancellation: %v", err)
	}
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("refresh did not cancel")
	}
}

func TestExpiredGrantWithoutRefreshRequiresLogin(t *testing.T) {
	m := testManager(t, testOptions(t), &fakeSecrets{}, nil)
	token := validToken()
	token.RefreshToken = ""
	token.Expiry = time.Now().Add(-time.Hour)
	seedToken(t, m, token)
	if _, err := m.client.Get("https://www.googleapis.com/drive/v3/files"); !errors.Is(err, ErrAuthorizationRequired) {
		t.Fatal(err)
	}
}

func TestRefreshCancellationAfterResponseDoesNotPersistRotation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	store := &fakeSecrets{}
	m := testManager(t, testOptions(t), store, func(r *http.Request) (*http.Response, error) {
		cancel()
		return response(r, 200, `{"access_token":"new","refresh_token":"new-refresh","expires_in":3600}`), nil
	})
	seedToken(t, m, validToken())
	if err := m.Refresh(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled response accepted: %v", err)
	}
	if store.sets != 1 {
		t.Fatal("rotation persisted after cancellation")
	}
}

func TestRequestBodyErrorDoesNotExposeTransportDetails(t *testing.T) {
	m := testManager(t, testOptions(t), &fakeSecrets{}, func(*http.Request) (*http.Response, error) {
		return nil, errors.New("private token and URL https://evil.example/?token=secret")
	})
	seedToken(t, m, validToken())
	r, _ := http.NewRequest("GET", "https://www.googleapis.com/drive/v3/files", nil)
	_, err := m.client.Transport.RoundTrip(r)
	if err == nil || strings.Contains(err.Error(), "private") || strings.Contains(err.Error(), "secret") {
		t.Fatalf("unsafe transport error: %v", err)
	}
}

func TestOpenURLRejectsUnsafeURLWithoutLaunching(t *testing.T) {
	for _, raw := range []string{"file:///tmp/no", "http://accounts.google.com/", "https://user:secret@accounts.google.com/", "https://accounts.google.com:444/", "https://accounts.google.com/#fragment"} {
		if err := OpenURL(raw); err == nil || strings.Contains(err.Error(), raw) {
			t.Fatalf("unsafe URL accepted or exposed: %v", err)
		}
	}
}
