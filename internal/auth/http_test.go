package auth

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Dunky13/takedown/internal/events"
)

func TestHTTPHostAndCredentialBoundaries(t *testing.T) {
	for _, host := range []string{"www.googleapis.com", "drive.googleapis.com", "dataportability.googleapis.com", "storage.googleapis.com", "bucket.storage.googleapis.com", "files.googleusercontent.com", "drive.usercontent.google.com", "files.usercontent.google.com"} {
		t.Run(host, func(t *testing.T) {
			m := testManager(t, testOptions(t), &fakeSecrets{}, func(r *http.Request) (*http.Response, error) {
				want := ""
				if apiHost(host) {
					want = "Bearer test-access-token"
				}
				if r.Header.Get("Authorization") != want {
					t.Errorf("bearer to %s: %q", host, r.Header.Get("Authorization"))
				}
				for _, key := range []string{"Cookie", "Cookie2", "Proxy-Authorization", "Referer"} {
					if r.Header.Get(key) != "" {
						t.Errorf("%s leaked", key)
					}
				}
				return response(r, 200, `{}`), nil
			})
			seedToken(t, m, validToken())
			r, _ := http.NewRequest("GET", "https://"+host+"/path?signed=secret", nil)
			for _, key := range []string{"Authorization", "Cookie", "Cookie2", "Proxy-Authorization", "Referer"} {
				r.Header.Set(key, "must-not-leak")
				r.Header[strings.ToLower(key)] = []string{"must-not-leak"}
			}
			res, err := m.client.Do(r)
			if err != nil {
				t.Fatal(err)
			}
			_ = res.Body.Close()
			if r.Header.Get("Authorization") != "must-not-leak" {
				t.Fatal("transport mutated caller headers")
			}
		})
	}
}

func TestHTTPRejectsUntrustedInitialAndRedirectRequests(t *testing.T) {
	for _, raw := range []string{
		"http://www.googleapis.com/", "https://www.googleapis.com:444/",
		"https://name:secret@www.googleapis.com/", "https://accounts.google.com/", "https://takeout.google.com/", "https://mail.google.com/",
		"https://googleusercontent.com/", "https://googleapis.com/", "https://www.googleapis.com.evil.example/", "https://evilgoogleusercontent.com/",
		"https://storage.googleapis.com.evil.example/", "https://evilusercontent.google.com/", "https://www.googleapis.com/#fragment",
	} {
		t.Run(raw, func(t *testing.T) {
			m := testManager(t, testOptions(t), &fakeSecrets{}, nil)
			r, err := http.NewRequest("GET", raw, nil)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := m.client.Transport.RoundTrip(r); err == nil {
				t.Fatal("untrusted initial request accepted")
			}
			if err := checkRedirect(r, []*http.Request{r}); err == nil {
				t.Fatal("untrusted redirect accepted")
			}
		})
	}
	m := testManager(t, testOptions(t), &fakeSecrets{}, nil)
	r, _ := http.NewRequest("GET", "https://www.googleapis.com/", nil)
	r.Host = "evil.example"
	if _, err := m.client.Transport.RoundTrip(r); err == nil {
		t.Fatal("Host override accepted")
	}
}

func TestAPIRedirectToSignedDownloadNeverCarriesBearerOrCookies(t *testing.T) {
	var calls int
	m := testManager(t, testOptions(t), &fakeSecrets{}, func(r *http.Request) (*http.Response, error) {
		calls++
		if calls == 1 {
			if r.Header.Get("Authorization") != "Bearer test-access-token" {
				t.Error("API bearer missing")
			}
			res := response(r, 302, "")
			res.Header.Set("Location", "https://bucket.storage.googleapis.com/archive?signed=secret")
			res.Header.Set("Set-Cookie", "SID=must-not-leak; Domain=.googleapis.com")
			return res, nil
		}
		for _, key := range []string{"Authorization", "Cookie", "Referer"} {
			if r.Header.Get(key) != "" {
				t.Errorf("redirect leaked %s", key)
			}
		}
		return response(r, 200, "archive"), nil
	})
	seedToken(t, m, validToken())
	r, err := m.client.Get("https://www.googleapis.com/drive/v3/files/file?alt=media")
	if err != nil {
		t.Fatal(err)
	}
	_ = r.Body.Close()
	if calls != 2 {
		t.Fatal("redirect not followed")
	}
}

func TestSignedDownloadsDoNotReadKeyring(t *testing.T) {
	store := &fakeSecrets{getErr: errors.New("keyring unavailable")}
	m := testManager(t, testOptions(t), store, func(r *http.Request) (*http.Response, error) {
		if r.Header.Get("Authorization") != "" {
			t.Error("download has bearer")
		}
		return response(r, 200, "archive"), nil
	})
	r, err := m.client.Get("https://files.googleusercontent.com/archive?signed=secret")
	if err != nil {
		t.Fatal(err)
	}
	_ = r.Body.Close()
	if store.gets+store.sets != 0 {
		t.Fatal("signed download touched keyring")
	}
}

func TestRedirectLimitAndAllowedPort(t *testing.T) {
	var calls int
	m := testManager(t, testOptions(t), &fakeSecrets{}, func(r *http.Request) (*http.Response, error) {
		calls++
		res := response(r, 302, "")
		res.Header.Set("Location", "https://storage.googleapis.com:443/archive")
		return res, nil
	})
	_, err := m.client.Get("https://storage.googleapis.com:443/archive")
	if err == nil || calls != 10 {
		t.Fatalf("redirects=%d error=%v", calls, err)
	}
}

func TestMutationRedirectIsNeverReplayed(t *testing.T) {
	var calls int
	m := testManager(t, testOptions(t), &fakeSecrets{}, func(r *http.Request) (*http.Response, error) {
		calls++
		res := response(r, 307, "")
		res.Header.Set("Location", "https://www.googleapis.com/redirected")
		return res, nil
	})
	seedToken(t, m, validToken())
	r, _ := http.NewRequest("POST", "https://www.googleapis.com/create", strings.NewReader(`{}`))
	if _, err := m.client.Do(r); err == nil || calls != 1 {
		t.Fatal("mutation redirect replayed")
	}
}

func TestAPIUnauthorizedReturnsAuthorizationRequiredAndPublishes(t *testing.T) {
	opts := testOptions(t)
	opts.Bus = events.New()
	ch, stop := opts.Bus.Subscribe(4)
	defer stop()
	m := testManager(t, opts, &fakeSecrets{}, func(r *http.Request) (*http.Response, error) {
		return response(r, 401, `{"secret":"must-not-be-reported"}`), nil
	})
	seedToken(t, m, validToken())
	_, err := m.client.Get("https://www.googleapis.com/drive/v3/files")
	if !errors.Is(err, ErrAuthorizationRequired) || strings.Contains(err.Error(), "must-not-be-reported") {
		t.Fatalf("wrong/unsafe error: %v", err)
	}
	select {
	case e := <-ch:
		if e.Kind != events.ReAuthRequired {
			t.Fatal(e)
		}
	default:
		t.Fatal("missing ReAuthRequired")
	}
}

func TestCloseCancelsDownloadBodyLifetime(t *testing.T) {
	var requestContext context.Context
	m := testManager(t, testOptions(t), &fakeSecrets{}, func(r *http.Request) (*http.Response, error) {
		requestContext = r.Context()
		return response(r, 200, "body"), nil
	})
	r, err := m.client.Get("https://storage.googleapis.com/archive")
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-requestContext.Done():
		t.Fatal("body lifetime canceled early")
	default:
	}
	_ = m.Close()
	select {
	case <-requestContext.Done():
	case <-time.After(time.Second):
		t.Fatal("Close did not cancel body context")
	}
	_ = r.Body.Close()
}

func TestBodyCloseOrEOFCleansOperation(t *testing.T) {
	for _, read := range []bool{false, true} {
		var ctx context.Context
		m := testManager(t, testOptions(t), &fakeSecrets{}, func(r *http.Request) (*http.Response, error) { ctx = r.Context(); return response(r, 200, "body"), nil })
		r, err := m.client.Get("https://storage.googleapis.com/archive")
		if err != nil {
			t.Fatal(err)
		}
		if read {
			_, _ = io.ReadAll(r.Body)
		} else {
			_ = r.Body.Close()
		}
		select {
		case <-ctx.Done():
		default:
			t.Fatal("body cleanup missing")
		}
		_ = r.Body.Close()
	}
}

func TestHTTPPreservesRequestContextCancellation(t *testing.T) {
	m := testManager(t, testOptions(t), &fakeSecrets{}, func(r *http.Request) (*http.Response, error) { <-r.Context().Done(); return nil, r.Context().Err() })
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r, _ := http.NewRequestWithContext(ctx, "GET", "https://storage.googleapis.com/archive", nil)
	_, err := m.client.Do(r)
	if !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}

func TestSecureURLRejectsOpaqueAndCredentials(t *testing.T) {
	for _, u := range []*url.URL{nil, {Scheme: "https", Opaque: "//www.googleapis.com"}, {Scheme: "https", Host: "www.googleapis.com:"}, {Scheme: "https", Host: "www.googleapis.com", User: url.User("secret")}} {
		if secureURL(u) {
			t.Fatal("insecure URL accepted")
		}
	}
}

func TestRead401RefreshesFreshLookingTokenAndReplaysOnce(t *testing.T) {
	for _, method := range []string{http.MethodGet, http.MethodHead} {
		var refreshes, reads int
		m := testManager(t, testOptions(t), &fakeSecrets{}, func(r *http.Request) (*http.Response, error) {
			if r.URL.String() == googleTokenURL {
				refreshes++
				return response(r, 200, `{"access_token":"fresh","refresh_token":"rotated","expires_in":3600}`), nil
			}
			reads++
			if r.Header.Get("Authorization") == "Bearer test-access-token" {
				return response(r, 401, `{}`), nil
			}
			if r.Header.Get("Authorization") != "Bearer fresh" {
				t.Error("wrong replay token")
			}
			return response(r, 200, `{}`), nil
		})
		seedToken(t, m, validToken())
		r, _ := http.NewRequest(method, "https://www.googleapis.com/drive/v3/files", nil)
		res, err := m.client.Do(r)
		if err != nil {
			t.Fatal(err)
		}
		_ = res.Body.Close()
		if refreshes != 1 || reads != 2 {
			t.Fatalf("refreshes=%d reads=%d", refreshes, reads)
		}
	}
}

func TestConcurrent401RefreshesOnce(t *testing.T) {
	const workers = 16
	var rejected, refreshes, replayed atomic.Int32
	allInitial := make(chan struct{})
	m := testManager(t, testOptions(t), &fakeSecrets{}, func(r *http.Request) (*http.Response, error) {
		if r.URL.String() == googleTokenURL {
			refreshes.Add(1)
			return response(r, 200, `{"access_token":"fresh","expires_in":3600}`), nil
		}
		if r.Header.Get("Authorization") == "Bearer test-access-token" {
			if rejected.Add(1) == workers {
				close(allInitial)
			}
			select {
			case <-allInitial:
			case <-r.Context().Done():
				return nil, r.Context().Err()
			}
			return response(r, 401, `{}`), nil
		}
		replayed.Add(1)
		return response(r, 200, `{}`), nil
	})
	seedToken(t, m, validToken())
	var wg sync.WaitGroup
	for range workers {
		wg.Go(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			r, _ := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/drive/v3/files", nil)
			res, err := m.client.Do(r)
			if err != nil {
				t.Error(err)
				return
			}
			_ = res.Body.Close()
		})
	}
	wg.Wait()
	if refreshes.Load() != 1 || replayed.Load() != workers {
		t.Fatalf("refreshes=%d replays=%d", refreshes.Load(), replayed.Load())
	}
}

func TestPermanent401IsBoundedAndMutationIsNeverReplayed(t *testing.T) {
	for _, method := range []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete} {
		var refreshes, requests int
		m := testManager(t, testOptions(t), &fakeSecrets{}, func(r *http.Request) (*http.Response, error) {
			if r.URL.String() == googleTokenURL {
				refreshes++
				return response(r, 200, `{"access_token":"fresh","expires_in":3600}`), nil
			}
			requests++
			return response(r, 401, `{}`), nil
		})
		seedToken(t, m, validToken())
		r, _ := http.NewRequest(method, "https://dataportability.googleapis.com/v1/portabilityArchive:initiate", nil)
		_, err := m.client.Do(r)
		if !errors.Is(err, ErrAuthorizationRequired) {
			t.Fatal(err)
		}
		wantRequests, wantRefreshes := 1, 0
		if method == "GET" {
			wantRequests, wantRefreshes = 2, 1
		}
		if requests != wantRequests || refreshes != wantRefreshes {
			t.Fatalf("%s requests=%d refreshes=%d", method, requests, refreshes)
		}
		if method == "GET" {
			_, err = m.client.Get("https://www.googleapis.com/drive/v3/files")
			if !errors.Is(err, ErrAuthorizationRequired) || requests != 2 || refreshes != 1 {
				t.Fatal("terminal grant retried")
			}
		}
	}
}

func TestInsufficientScope403DoesNotRefreshOrReplay(t *testing.T) {
	var calls int
	m := testManager(t, testOptions(t), &fakeSecrets{}, func(r *http.Request) (*http.Response, error) {
		calls++
		res := response(r, 403, `{}`)
		res.Header.Set("WWW-Authenticate", `Bearer error="insufficient_scope"`)
		return res, nil
	})
	seedToken(t, m, validToken())
	_, err := m.client.Get("https://www.googleapis.com/drive/v3/files")
	if !errors.Is(err, ErrAuthorizationRequired) || calls != 1 {
		t.Fatal("insufficient grant refreshed/replayed")
	}
}

type failingBody struct{}

func (failingBody) Read([]byte) (int, error) {
	return 0, errors.New("private response diagnostic token=hidden")
}
func (failingBody) Close() error { return errors.New("private close diagnostic token=hidden") }

func TestResponseBodyErrorsAreSanitized(t *testing.T) {
	m := testManager(t, testOptions(t), &fakeSecrets{}, func(r *http.Request) (*http.Response, error) {
		res := response(r, 200, "")
		res.Body = failingBody{}
		return res, nil
	})
	res, err := m.client.Get("https://storage.googleapis.com/archive")
	if err != nil {
		t.Fatal(err)
	}
	_, readErr := io.ReadAll(res.Body)
	closeErr := res.Body.Close()
	for _, err := range []error{readErr, closeErr} {
		if err == nil || strings.Contains(err.Error(), "hidden") || strings.Contains(err.Error(), "private") {
			t.Fatalf("unsafe response body error: %v", err)
		}
	}
}
