package source

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Dunky13/takedown/internal/auth"
	"google.golang.org/api/dataportability/v1"
	"google.golang.org/api/drive/v3"
)

// Optional downloader contract: callers that know the failed URL can distinguish
// an old generation from the current cache even for concurrent same-part calls.
type rejectedLinkSource interface {
	RefreshRejectedLink(context.Context, string, int, string) (string, error)
}

func refreshRejected(m *Manager, ctx context.Context, id string, index int, rejected string) (string, error) {
	if richer, ok := interface{}(m).(rejectedLinkSource); ok {
		return richer.RefreshRejectedLink(ctx, id, index, rejected)
	}
	return m.RefreshLink(ctx, id, index)
}

func multipartFixture(t *testing.T) (*Manager, *time.Time, *int, string, []string) {
	t.Helper()
	now := time.Date(2026, 10, 25, 10, 0, 0, 0, time.UTC)
	polls := 0
	m := testManager(t, Options{Provider: "portability", Resources: []string{"myactivity.search"}}, func(r *http.Request) (*http.Response, error) {
		polls++
		if r.Method != http.MethodGet || r.URL.Path != "/v1/archiveJobs/job-123/portabilityArchiveState" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		return responseJSON(r, 200, &dataportability.PortabilityArchiveState{State: "COMPLETE", Urls: []string{fmt.Sprintf("https://storage.googleapis.com/bucket/first.zip?signature=generation-%d", polls), fmt.Sprintf("https://storage.googleapis.com/bucket/second.zip?signature=generation-%d", polls)}}), nil
	})
	m.now = func() time.Time { return now }
	installSigned(t, m, func(r *http.Request) (*http.Response, error) {
		if strings.HasSuffix(r.URL.Path, "first.zip") {
			return signedResponse(r, "first.zip", 100), nil
		}
		return signedResponse(r, "second.zip", 200), nil
	})
	seedJobs(t, m, testJob(now.Add(-time.Hour)))
	exports, err := m.List(context.Background())
	if err != nil || len(exports) != 1 {
		t.Fatalf("fixture setup: %v", err)
	}
	return m, &now, &polls, exports[0].ID, []string{exports[0].Parts[0].URL, exports[0].Parts[1].URL}
}

func TestMultipartRefreshSharesJobGeneration(t *testing.T) {
	for _, concurrent := range []bool{false, true} {
		t.Run(fmt.Sprintf("concurrent=%t", concurrent), func(t *testing.T) {
			m, now, polls, id, old := multipartFixture(t)
			*now = now.Add(6 * time.Hour)
			type result struct {
				index int
				link  string
				err   error
			}
			results := make(chan result, 3)
			refresh := func(index int) {
				link, err := refreshRejected(m, context.Background(), id, index, old[index-1])
				results <- result{index, link, err}
			}
			if concurrent {
				var ready sync.WaitGroup
				ready.Add(3)
				start := make(chan struct{})
				for _, index := range []int{1, 2, 1} {
					go func(index int) { ready.Done(); <-start; refresh(index) }(index)
				}
				ready.Wait()
				close(start)
			} else {
				refresh(1)
				refresh(2)
				refresh(1)
			}
			for range 3 {
				var r result
				select {
				case r = <-results:
				case <-time.After(5 * time.Second):
					t.Fatal("concurrent refresh did not finish")
				}
				if r.err != nil || r.link == old[r.index-1] || !strings.Contains(r.link, "generation-2") {
					t.Errorf("part %d did not share renewed job links: %q %v", r.index, r.link, r.err)
				}
			}
			if *polls != 2 {
				t.Fatalf("job was unnecessarily polled %d times", *polls)
			}
		})
	}
}

func TestRejectedCurrentGenerationWaitsWithoutReturningBadURL(t *testing.T) {
	m, now, polls, id, old := multipartFixture(t)
	// A newly listed URL is already the current generation. It must not be
	// returned unchanged just because a cached job entry exists.
	if link, err := refreshRejected(m, context.Background(), id, 1, old[0]); link != "" || !IsTemporary(err) || RetryDelay(err) != 5*time.Minute {
		t.Fatalf("rejected current URL: %q %v delay=%v", link, err, RetryDelay(err))
	}
	if *polls != 1 {
		t.Fatal("current-generation rejection bypassed cadence")
	}
	*now = now.Add(5 * time.Minute)
	fresh, err := refreshRejected(m, context.Background(), id, 1, old[0])
	if err != nil || fresh == old[0] {
		t.Fatal("new generation was not fetched", err)
	}
	*now = now.Add(time.Minute)
	if link, err := refreshRejected(m, context.Background(), id, 1, fresh); link != "" || !IsTemporary(err) || RetryDelay(err) != 4*time.Minute {
		t.Fatalf("fresh rejection did not preserve remaining cooldown: %q %v %v", link, err, RetryDelay(err))
	}
	// An old caller must not get a current URL that another caller just rejected.
	if link, err := refreshRejected(m, context.Background(), id, 1, old[0]); link != "" || !IsTemporary(err) {
		t.Fatal("known rejected current signature reused", link, err)
	}
	if *polls != 2 {
		t.Fatal("cooldown or old generation triggered an extra API poll", *polls)
	}
	// An unaffected sibling still receives its newly signed URL from the same job.
	if link, err := refreshRejected(m, context.Background(), id, 2, old[1]); err != nil || !strings.Contains(link, "generation-2") {
		t.Fatal("unaffected sibling failed", link, err)
	}
}

func TestLegacyMultipartRefreshSharesNewGeneration(t *testing.T) {
	m, now, polls, id, _ := multipartFixture(t)
	*now = now.Add(6 * time.Hour)
	for _, index := range []int{1, 2} {
		if _, err := m.RefreshLink(context.Background(), id, index); err != nil {
			t.Fatalf("legacy sibling refresh: %v (polls=%d)", err, *polls)
		}
	}
	if *polls != 2 {
		t.Fatal("legacy sibling polled the job again", *polls)
	}
}

func TestRefreshNeverReturnsUnchangedRejectedSignatureAfterPoll(t *testing.T) {
	for _, slowProbe := range []bool{false, true} {
		t.Run(fmt.Sprintf("slow-probe=%t", slowProbe), func(t *testing.T) {
			m, now, polls, id, old := multipartFixture(t)
			*now = now.Add(5 * time.Minute)
			m.HTTPClient().Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
				*polls++
				return responseJSON(r, 200, &dataportability.PortabilityArchiveState{State: "COMPLETE", Urls: old}), nil
			})
			if slowProbe {
				installSigned(t, m, func(r *http.Request) (*http.Response, error) {
					if strings.HasSuffix(r.URL.Path, "first.zip") {
						return signedResponse(r, "first.zip", 100), nil
					}
					*now = now.Add(6 * time.Minute)
					return signedResponse(r, "second.zip", 200), nil
				})
			}
			link, err := refreshRejected(m, context.Background(), id, 1, old[0])
			if link != "" || !IsTemporary(err) {
				t.Fatalf("rejected signature returned after API poll: %q %v", link, err)
			}
			if !slowProbe && RetryDelay(err) != 5*time.Minute {
				t.Fatal("unchanged signature lost cooldown", RetryDelay(err))
			}
			if *polls != 2 {
				t.Fatal("unchanged signature caused extra API polls", *polls)
			}
		})
	}
}

func TestDiscoveryPreservesAuthorizationClassification(t *testing.T) {
	for _, cause := range []string{"no usable OAuth token is stored; run Login", "Google rejected or revoked the OAuth grant; run Login again", "configure a Google OAuth Desktop app client JSON file, then run Login"} {
		t.Run(cause, func(t *testing.T) {
			m := testManager(t, Options{DriveFolderID: "fixture-folder"}, func(r *http.Request) (*http.Response, error) {
				return nil, fmt.Errorf("%w: %s", auth.ErrAuthorizationRequired, cause)
			})
			_, err := m.List(context.Background())
			if !errors.Is(err, auth.ErrAuthorizationRequired) || !strings.Contains(err.Error(), "Login") || !strings.Contains(err.Error(), "Desktop app client") {
				t.Fatalf("authorization became connectivity failure: %v", err)
			}
			if IsTemporary(err) || RetryDelay(err) != 0 || errors.Unwrap(err) != nil {
				t.Fatal("authorization is retryable or wraps an unsafe cause", err)
			}
		})
	}
	// The actual auth manager diagnoses missing client setup before keyring or
	// network access, matching the retained application-level overlay witness.
	for _, provider := range []string{"drive", "portability"} {
		t.Run("real-auth/"+provider, func(t *testing.T) {
			resources := []string(nil)
			scope := drive.DriveReadonlyScope
			if provider == "portability" {
				resources = []string{"myactivity.search"}
				scope = dataportability.DataportabilityMyactivitySearchScope
			}
			dir := t.TempDir()
			a, err := auth.New(auth.Options{DataDir: dir, Scopes: []string{scope}})
			if err != nil {
				t.Fatal(err)
			}
			m, err := New(Options{DataDir: dir, Provider: provider, Resources: resources, DriveFolderID: "fixture-folder", Session: a})
			if err != nil {
				t.Fatal(err)
			}
			defer m.Close()
			_, err = m.List(context.Background())
			if !errors.Is(err, auth.ErrAuthorizationRequired) || !strings.Contains(err.Error(), "Desktop app client") {
				t.Fatal("real auth setup instructions lost", err)
			}
		})
	}
}

func TestAuthorizationClassificationDoesNotExposeURLOrProviderText(t *testing.T) {
	raw := &url.Error{Op: "Get", URL: "https://www.googleapis.com/drive/v3/files?access_token=secret-query", Err: fmt.Errorf("%w: secret-response-body", auth.ErrAuthorizationRequired)}
	err := apiError(context.Background(), "read Drive metadata", raw)
	if !errors.Is(err, auth.ErrAuthorizationRequired) || IsTemporary(err) || errors.Unwrap(err) != nil {
		t.Fatal("authorization classification lost or cause retained", err)
	}
	for _, secret := range []string{"secret-query", "secret-response-body", "https://"} {
		if strings.Contains(err.Error(), secret) {
			t.Fatalf("unsafe auth text: %v", err)
		}
	}
}

func TestCooldownErrorExposesProviderRetryContract(t *testing.T) {
	now := time.Date(2026, 10, 25, 10, 0, 0, 0, time.UTC)
	j := testJob(now.Add(-time.Hour))
	j.LastPoll = now.Add(-time.Minute)
	j.RetryNotBefore = now.Add(2 * time.Hour)
	err := portabilityCooldown(j, now)
	var provider interface {
		Temporary() bool
		RetryAfter() time.Duration
	}
	if !errors.As(err, &provider) || !provider.Temporary() || provider.RetryAfter() != 2*time.Hour || RetryDelay(err) != 2*time.Hour {
		t.Fatal("provider cooldown contract lost", err)
	}
}
