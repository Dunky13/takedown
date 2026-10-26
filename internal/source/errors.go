package source

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/Dunky13/takedown/internal/auth"
	"google.golang.org/api/googleapi"
)

type sourceError struct {
	message       string
	temporary     bool
	delay         time.Duration
	authorization bool
}

// A confirmed failure of one job's identity/metadata can be reported while
// ready sibling jobs proceed. Authorization, quota, transport and journal
// failures never receive this marker.
type jobReadError struct{ error }

func jobFailure(message string) error { return &jobReadError{errors.New(message)} }
func isJobFailure(err error) bool     { var e *jobReadError; return errors.As(err, &e) }

// IsJobFailure reports a confirmed job-local status/identity/metadata hold.
// A watcher may keep running at its normal bounded cadence after this error.
// Interactive List still returns no partial result for a bad job when no
// unrelated ready export exists. Auth, quota, provider backoff, transport and
// local journal/manifest failures never receive this classification.
func IsJobFailure(err error) bool { return isJobFailure(err) }

func (e *sourceError) Error() string   { return e.message }
func (e *sourceError) Temporary() bool { return e.temporary }
func (e *sourceError) RetryAfter() time.Duration {
	if e.temporary {
		return e.delay
	}
	return 0
}

// Preserve a safe sentinel, never the raw provider or HTTP error chain.
func (e *sourceError) Is(target error) bool {
	return e.authorization && target == auth.ErrAuthorizationRequired
}

// IsTemporary classifies only sanitized read failures produced by this package.
// Creation failures deliberately never carry a retryable classification.
func IsTemporary(err error) bool { var e *sourceError; return errors.As(err, &e) && e.temporary }
func RetryDelay(err error) time.Duration {
	var e *sourceError
	if errors.As(err, &e) && e.temporary {
		return e.delay
	}
	return 0
}

func apiTemporary(g *googleapi.Error) bool {
	reasons := []string{}
	for _, e := range g.Errors {
		reasons = append(reasons, e.Reason)
	}
	for _, detail := range g.Details {
		if d, ok := detail.(map[string]interface{}); ok {
			if r, ok := d["reason"].(string); ok {
				reasons = append(reasons, r)
			}
		}
	}
	rate := false
	for _, reason := range reasons {
		switch reason {
		case "RESOURCE_EXHAUSTED_ONE_TIME", "RESOURCE_EXHAUSTED_TIME_BASED", "DAILY_LIMIT_EXCEEDED", "dailyLimitExceeded", "dailyLimitExceededUnreg", "downloadQuotaExceeded", "quotaExceeded":
			return false
		case "rateLimitExceeded", "userRateLimitExceeded", "RATE_LIMIT_EXCEEDED":
			rate = true
		}
	}
	return g.Code == 429 || g.Code >= 500 && g.Code <= 599 || g.Code == 403 && rate
}

// Raw Google/network errors may contain signed URLs, server response bodies, or
// credentials. Neither render nor wrap them into user-visible error chains.
func apiError(ctx context.Context, operation string, err error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if errors.Is(err, auth.ErrAuthorizationRequired) {
		return &sourceError{message: operation + ": Google authorization required; configure a valid Google OAuth Desktop app client JSON file, then run Login to authorize the selected scopes", authorization: true}
	}
	var g *googleapi.Error
	if errors.As(err, &g) {
		message := "Google API request failed; check connectivity and API configuration"
		switch g.Code {
		case 400:
			message = "Google rejected the request; check selected resource groups and API configuration"
		case 401:
			message = "Google authentication expired or was revoked; sign in again"
		case 403:
			message = "Google denied access; check consent scopes, API enablement, file permissions, and download restrictions"
		case 404:
			message = "Google resource is missing or inaccessible; check the configured folder or recorded job"
		case 409, 412:
			message = "Google export precondition failed; allow 24 hours between time-based exports and check grant validity"
		case 429:
			message = "Google quota or one-time grant is exhausted; wait for quota recovery or explicitly obtain fresh consent; no export was retried"
		default:
			if g.Code >= 500 {
				message = "Google service is temporarily unavailable; retry the read later; export creation is never retried automatically"
			}
		}
		if delay := retryAfter(g.Header.Get("Retry-After"), time.Now()); delay > 0 {
			message += fmt.Sprintf("; retry after %s", delay.Round(time.Second))
		}
		return &sourceError{message: fmt.Sprintf("%s: %s (HTTP %d)", operation, message, g.Code), temporary: apiTemporary(g), delay: retryAfter(g.Header.Get("Retry-After"), time.Now()), authorization: g.Code == http.StatusUnauthorized}
	}
	var network net.Error
	temporary := auth.IsTemporary(err) || errors.As(err, &network) && (network.Timeout() || network.Temporary())
	return &sourceError{message: fmt.Sprintf("%s: Google request failed; check connectivity and retry the read if the connection was interrupted", operation), temporary: temporary, delay: auth.RetryDelay(err)}
}

func retryAfter(value string, now time.Time) time.Duration {
	if secs, err := strconv.ParseInt(value, 10, 32); err == nil && secs > 0 {
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(value); err == nil && t.After(now) {
		return t.Sub(now)
	}
	return 0
}

// Read requests have three attempts, with bounded backoff. A long Retry-After
// is reported to the caller instead of retrying earlier than Google requested.
func readAPI[T any](ctx context.Context, operation string, call func(context.Context) (T, error)) (T, error) {
	var zero T
	for attempt := 0; attempt < 3; attempt++ {
		if err := ctx.Err(); err != nil {
			return zero, err
		}
		callCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		v, err := call(callCtx)
		cancel()
		if err == nil {
			return v, nil
		}
		safe := apiError(ctx, operation, err)
		if !IsTemporary(safe) || attempt == 2 {
			return zero, safe
		}
		delay := time.Second * time.Duration(1<<attempt)
		if r := RetryDelay(safe); r > delay {
			delay = r
		}
		if delay > 30*time.Second {
			return zero, safe
		}
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return zero, ctx.Err()
		case <-timer.C:
		}
	}
	return zero, errors.New("Google read retry limit reached")
}
