package auth

import (
	"errors"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// temporaryError deliberately has no Unwrap: network and OAuth errors can
// contain secret URLs, form values, and provider response bodies.
type temporaryError struct {
	timeout    bool
	retryAfter time.Duration
}

func (e *temporaryError) Error() string             { return "Google service temporarily unavailable; retry later" }
func (e *temporaryError) Timeout() bool             { return e.timeout }
func (e *temporaryError) Temporary() bool           { return true }
func (e *temporaryError) RetryAfter() time.Duration { return e.retryAfter }

// IsTemporary identifies sanitized OAuth/network outages, including errors
// wrapped by net/http. It never classifies authorization/configuration as an
// outage. Callers still own their cancellation and finite retry policy.
func IsTemporary(err error) bool {
	if err == nil || errors.Is(err, ErrAuthorizationRequired) {
		return false
	}
	var temporary interface{ Temporary() bool }
	return errors.As(err, &temporary) && temporary.Temporary()
}

// RetryDelay returns a provider Retry-After delay when supplied, or zero.
func RetryDelay(err error) time.Duration {
	var retry interface{ RetryAfter() time.Duration }
	if errors.As(err, &retry) {
		return retry.RetryAfter()
	}
	return 0
}

func retryAfter(header string) time.Duration {
	header = strings.TrimSpace(header)
	if seconds, err := strconv.ParseUint(header, 10, 64); err == nil {
		if seconds > uint64(time.Duration(1<<63-1)/time.Second) {
			return time.Duration(1<<63 - 1)
		}
		return time.Duration(seconds) * time.Second
	} else if errors.Is(err, strconv.ErrRange) && !strings.HasPrefix(header, "-") {
		return time.Duration(1<<63 - 1)
	}
	if when, err := http.ParseTime(header); err == nil {
		delay := time.Until(when)
		if delay > 0 {
			return delay
		}
	}
	return 0
}

func temporaryNetwork(err error) error {
	var network net.Error
	var operation *net.OpError
	if errors.As(err, &network) && (network.Timeout() || network.Temporary()) {
		return &temporaryError{timeout: network.Timeout()}
	}
	if errors.As(err, &operation) || errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return &temporaryError{}
	}
	return nil
}
