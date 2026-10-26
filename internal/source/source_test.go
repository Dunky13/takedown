package source

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"google.golang.org/api/dataportability/v1"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type fakeSession struct {
	client                        *http.Client
	profile                       string
	login, refresh, forget, close int
}

func (s *fakeSession) Login(context.Context) error   { s.login++; return nil }
func (s *fakeSession) Refresh(context.Context) error { s.refresh++; return nil }
func (s *fakeSession) Forget(context.Context) error  { s.forget++; return nil }
func (s *fakeSession) HTTPClient() *http.Client      { return s.client }
func (s *fakeSession) ProfileID(context.Context) (string, error) {
	if s.profile != "" {
		return s.profile, nil
	}
	return strings.Repeat("a", 64), nil
}
func (s *fakeSession) Close() error { s.close++; return nil }

func responseJSON(r *http.Request, status int, value any) *http.Response {
	b, _ := json.Marshal(value)
	return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(string(b))), ContentLength: int64(len(b)), Request: r}
}

func testManager(t *testing.T, opts Options, handler roundTripFunc) *Manager {
	t.Helper()
	if opts.DataDir == "" {
		opts.DataDir = t.TempDir()
	}
	if opts.Session == nil {
		opts.Session = &fakeSession{client: &http.Client{Transport: handler}}
	}
	m, err := New(opts)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := m.Close(); err != nil {
			t.Error(err)
		}
	})
	return m
}

func TestScopesAndUnconfiguredGUI(t *testing.T) {
	for _, provider := range []string{"", "drive"} {
		scopes, err := Scopes(provider, nil)
		if err != nil || len(scopes) != 1 || scopes[0] != drive.DriveReadonlyScope {
			t.Fatalf("drive scopes %v %v", scopes, err)
		}
	}
	scopes, err := Scopes("portability", []string{"youtube.private_videos", "myactivity.search"})
	if err != nil || len(scopes) != 2 || scopes[0] != dataportability.DataportabilityMyactivitySearchScope || scopes[1] != dataportability.DataportabilityYoutubePrivateVideosScope {
		t.Fatalf("selected scopes %v %v", scopes, err)
	}
	for _, resources := range [][]string{{"openid"}, {"email"}, {drive.DriveReadonlyScope}, {"https://www.googleapis.com/auth/dataportability.myactivity.search"}, {"myactivity.search", "myactivity.search"}, {"not.real"}, {"myactivity.search?x"}} {
		if _, err := Scopes("portability", resources); err == nil {
			t.Errorf("accepted %v", resources)
		}
	}
	if _, err := Scopes("drive", []string{"myactivity.search"}); err == nil {
		t.Fatal("mixed providers")
	}
	if _, err := Scopes("browser", nil); err == nil {
		t.Fatal("unknown provider")
	}
	for _, provider := range []string{"drive", "portability"} {
		m, err := New(Options{Provider: provider})
		if err != nil {
			t.Fatal(err)
		}
		if _, err = m.List(context.Background()); err == nil {
			t.Fatal("unconfigured list succeeded")
		}
		if err = m.Login(context.Background()); err == nil {
			t.Fatal("unconfigured login succeeded")
		}
		if provider == "portability" {
			if _, err = m.CreateExport(context.Background()); err == nil {
				t.Fatal("empty resources created export")
			}
		}
		_ = m.Close()
	}
}

func TestSessionForwardingAndTakeout(t *testing.T) {
	session := &fakeSession{client: &http.Client{}}
	var opened string
	m, err := New(Options{Session: session, OpenURL: func(s string) error { opened = s; return nil }})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if m.HTTPClient() != session.client {
		t.Fatal("session client not forwarded")
	}
	if err = m.Login(ctx); err != nil {
		t.Fatal(err)
	}
	if err = m.Refresh(ctx); err != nil {
		t.Fatal(err)
	}
	if err = m.Forget(ctx); err != nil {
		t.Fatal(err)
	}
	message, err := m.CreateExport(ctx)
	if err != nil || !strings.Contains(message, "Create export") {
		t.Fatalf("create %q %v", message, err)
	}
	u, err := url.Parse(opened)
	if err != nil || u.Host != "takeout.google.com" || u.Path != "/" || u.RawQuery != "" {
		t.Fatalf("takeout URL %q", opened)
	}
	_ = m.Close()
	_ = m.Close()
	if session.login != 1 || session.refresh != 1 || session.forget != 1 || session.close != 1 {
		t.Fatalf("forwarding %#v", session)
	}
	if _, err = m.List(ctx); err == nil {
		t.Fatal("closed manager operated")
	}
}

func TestFolderURLValidation(t *testing.T) {
	for _, s := range []string{"folder_123", "https://drive.google.com/drive/folders/folder_123?usp=sharing", "https://drive.google.com/drive/u/0/folders/folder_123"} {
		id, err := folderID(s)
		if err != nil || id != "folder_123" {
			t.Errorf("folder %q: %q %v", s, id, err)
		}
	}
	for _, s := range []string{"root", "", "https://drive.google.com/file/d/file_123/view", "http://drive.google.com/drive/folders/x", "https://drive.google.com.evil/drive/folders/x", "https://user@drive.google.com/drive/folders/x", "x' in parents", "https://drive.google.com:443/drive/folders/x", "https://drive.google.com/drive/folders/root"} {
		if _, err := folderID(s); err == nil {
			t.Errorf("accepted %q", s)
		}
	}
}

type temporaryNetwork struct{}

func (temporaryNetwork) Error() string {
	return "https://storage.googleapis.com/private.zip?secret=do-not-echo"
}
func (temporaryNetwork) Timeout() bool   { return true }
func (temporaryNetwork) Temporary() bool { return true }

type oauthTemporaryDelay struct{ temporaryNetwork }

func (oauthTemporaryDelay) RetryAfter() time.Duration { return 2 * time.Hour }

func TestSanitizedTemporaryErrors(t *testing.T) {
	for _, tc := range []struct {
		code      int
		reason    string
		temporary bool
	}{
		{429, "", true}, {503, "", true}, {403, "rateLimitExceeded", true}, {403, "userRateLimitExceeded", true},
		{403, "downloadQuotaExceeded", false}, {403, "dailyLimitExceeded", false}, {403, "insufficientPermissions", false},
		{429, "RESOURCE_EXHAUSTED_ONE_TIME", false}, {401, "", false}, {400, "", false}, {404, "", false},
	} {
		g := &googleapi.Error{Code: tc.code, Message: "secret-token", Header: http.Header{"Retry-After": []string{"7200"}}, Errors: []googleapi.ErrorItem{{Reason: tc.reason}}}
		err := apiError(context.Background(), "read", g)
		if IsTemporary(err) != tc.temporary {
			t.Errorf("%d/%s temporary=%v", tc.code, tc.reason, IsTemporary(err))
		}
		if strings.Contains(err.Error(), "secret-token") || errors.Unwrap(err) != nil {
			t.Fatal("unsafe error chain", err)
		}
		if tc.temporary && RetryDelay(err) != 2*time.Hour {
			t.Fatalf("lost long Retry-After: %v", RetryDelay(err))
		}
		if !tc.temporary && RetryDelay(err) != 0 {
			t.Fatal("terminal retry delay")
		}
	}
	err := apiError(context.Background(), "read", &url.Error{Op: "Get", URL: "https://secret.invalid?token=sensitive", Err: temporaryNetwork{}})
	if !IsTemporary(err) || strings.Contains(err.Error(), "sensitive") || errors.Unwrap(err) != nil {
		t.Fatal(err)
	}
	if IsTemporary(errors.New("temporary")) {
		t.Fatal("arbitrary error classified")
	}
	oauth := apiError(context.Background(), "read", &url.Error{Op: "Get", URL: "https://secret.invalid/", Err: oauthTemporaryDelay{}})
	if !IsTemporary(oauth) || RetryDelay(oauth) != 2*time.Hour {
		t.Fatal("OAuth transport lost requested delay", oauth, RetryDelay(oauth))
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if !errors.Is(apiError(ctx, "read", temporaryNetwork{}), context.Canceled) {
		t.Fatal("cancellation lost")
	}
	calls := 0
	_, err = readAPI(context.Background(), "read", func(context.Context) (int, error) {
		calls++
		return 0, &googleapi.Error{Code: 503, Header: http.Header{"Retry-After": []string{"3600"}}}
	})
	if calls != 1 || !IsTemporary(err) || RetryDelay(err) != time.Hour {
		t.Fatal("long retry-after was internally retried or lost", calls, err)
	}
}
