package auth

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"syscall"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func TestOAuthOutagesKeepSanitizedClassificationAndRetryDelay(t *testing.T) {
	for _, status := range []int{408, 429, 500, 502, 503, 504} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			m := testManager(t, testOptions(t), &fakeSecrets{}, func(r *http.Request) (*http.Response, error) {
				res := response(r, status, `{"error":"temporarily_unavailable","error_description":"secret-token https://evil.example/?code=hidden"}`)
				res.Header.Set("Retry-After", "17")
				return res, nil
			})
			seedToken(t, m, validToken())
			err := m.Refresh(context.Background())
			if !IsTemporary(err) || RetryDelay(err) != 17*time.Second || errors.Is(err, ErrAuthorizationRequired) {
				t.Fatalf("outage classification: %v delay %v", err, RetryDelay(err))
			}
			var network net.Error
			var provider *oauth2.RetrieveError
			if !errors.As(err, &network) || network.Timeout() != (status == 408) || errors.As(err, &provider) || errors.Unwrap(err) != nil {
				t.Fatal("lost sanitized net.Error classification or exposed original")
			}
			if strings.Contains(err.Error(), "hidden") || strings.Contains(err.Error(), "secret-token") {
				t.Fatal("secret in error")
			}
		})
	}
}

func TestOAuthNetworkOutageIsTemporaryWithoutRawUnwrap(t *testing.T) {
	for _, original := range []error{
		&url.Error{Op: "Post", URL: googleTokenURL + "?client_secret=hidden", Err: &net.DNSError{Err: "private", Name: "secret-name", IsTimeout: true}},
		&net.OpError{Op: "dial", Net: "tcp", Err: syscall.ECONNREFUSED},
	} {
		m := testManager(t, testOptions(t), &fakeSecrets{}, func(*http.Request) (*http.Response, error) { return nil, original })
		seedToken(t, m, validToken())
		err := m.Refresh(context.Background())
		if !IsTemporary(err) || strings.Contains(err.Error(), "private") || strings.Contains(err.Error(), "hidden") {
			t.Fatalf("network classification: %v", err)
		}
		var raw *net.OpError
		var rawURL *url.Error
		if errors.As(err, &raw) || errors.As(err, &rawURL) {
			t.Fatal("raw error exposed")
		}
	}
}

func TestInvalidGrantRemainsTerminalEvenWithTemporaryHTTPStatus(t *testing.T) {
	m := testManager(t, testOptions(t), &fakeSecrets{}, func(r *http.Request) (*http.Response, error) {
		return response(r, 503, `{"error":"invalid_grant"}`), nil
	})
	seedToken(t, m, validToken())
	err := m.Refresh(context.Background())
	if !errors.Is(err, ErrAuthorizationRequired) || IsTemporary(err) {
		t.Fatal(err)
	}
}

func TestRetryDelayThroughHTTPWrapperAndDate(t *testing.T) {
	err := &url.Error{Op: "Get", URL: "https://www.googleapis.com/", Err: &temporaryError{retryAfter: 9 * time.Second}}
	if !IsTemporary(err) || RetryDelay(err) != 9*time.Second {
		t.Fatal("HTTP wrapper lost classification")
	}
	delay := retryAfter(time.Now().Add(time.Minute).UTC().Format(http.TimeFormat))
	if delay < 58*time.Second || delay > time.Minute {
		t.Fatalf("HTTP date delay: %v", delay)
	}
	for _, invalid := range []string{"-1", "bogus", "0"} {
		if retryAfter(invalid) != 0 {
			t.Fatal("invalid delay accepted")
		}
	}
	if retryAfter("172800") != 48*time.Hour || retryAfter("999999999999999999999999") != time.Duration(1<<63-1) {
		t.Fatal("long Retry-After was discarded")
	}
}
