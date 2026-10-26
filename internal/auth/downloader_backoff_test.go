package auth

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Dunky13/takedown/internal/downloader"
)

func TestOAuthRetryAfterReachesDownloader(t *testing.T) {
	m := testManager(t, testOptions(t), &fakeSecrets{}, func(r *http.Request) (*http.Response, error) {
		res := response(r, 503, `{"error":"temporarily_unavailable"}`)
		res.Header.Set("Retry-After", "3600")
		return res, nil
	})
	token := validToken()
	token.Expiry = time.Now().Add(-time.Hour)
	seedToken(t, m, token)
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://www.googleapis.com/drive/v3/files/example?alt=media", nil)
	_, err := m.HTTPClient().Do(req)
	if !IsTemporary(err) || RetryDelay(err) != time.Hour {
		t.Fatalf("auth setup failed: %v", err)
	}
	if downloader.RetryDelay(err) != time.Hour {
		t.Fatalf("OAuth says retry after %v; downloader says %v (temporary=%v)", RetryDelay(err), downloader.RetryDelay(err), downloader.IsTemporary(err))
	}
}
