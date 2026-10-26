package downloader

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/Dunky13/takedown/internal/model"
)

func TestDownloadQuotaDoesNotTriggerReauthentication(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error":{"errors":[{"reason":"downloadQuotaExceeded"}]}}`))
	}))
	defer server.Close()
	f := newFixture(t, zipBytes(t, 50), "takeout.zip", server.URL, nil)
	err := f.engine.Download(context.Background(), []model.Part{f.part})
	if err == nil || !strings.Contains(err.Error(), "quota") || errors.Is(err, ErrReauthentication) || IsTemporary(err) {
		t.Fatalf("quota handling: %v", err)
	}
}

type providerBackoff struct{}

func (providerBackoff) Error() string             { return "temporary authorization outage" }
func (providerBackoff) Temporary() bool           { return true }
func (providerBackoff) Timeout() bool             { return false }
func (providerBackoff) RetryAfter() time.Duration { return time.Hour }

func TestProviderRetryAfterSurvivesHTTPErrorWrappers(t *testing.T) {
	cause := &safeError{"archive HTTP request failed", &url.Error{Op: "Get", URL: "https://example.test?secret=SIGNED", Err: providerBackoff{}}}
	if !IsTemporary(cause) || RetryDelay(cause) != time.Hour {
		t.Fatalf("provider retry metadata was lost: %v", cause)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := retryWait(ctx, 0, cause)
	if !IsTemporary(err) || RetryDelay(err) != time.Hour || errors.Is(err, context.Canceled) {
		t.Fatalf("long backoff should return immediately to watch: %v", err)
	}
	if strings.Contains(err.Error(), "SIGNED") {
		t.Fatal("provider URL leaked")
	}
}

func TestLongRetryAfterPreservesBackoffWithoutSleepingOrRefreshing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "120")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()
	f := newFixture(t, zipBytes(t, 50), "takeout.zip", server.URL, nil)
	start := time.Now()
	err := f.engine.Download(context.Background(), []model.Part{f.part})
	if !IsTemporary(err) || RetryDelay(err) != 2*time.Minute {
		t.Fatalf("lost retry classification: %v", err)
	}
	if time.Since(start) > 5*time.Second {
		t.Fatal("a short batch slept beyond its retry budget")
	}
}
