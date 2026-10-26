package downloader

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Dunky13/takedown/internal/events"
	"github.com/Dunky13/takedown/internal/model"
)

// safeError keeps the cause available to errors.Is/As without rendering signed
// request URLs from net/http's url.Error in logs, events, or UI messages.
type safeError struct {
	message string
	cause   error
}

func (e *safeError) Error() string { return e.message }
func (e *safeError) Unwrap() error { return e.cause }

type transientStatus struct {
	code  int
	after time.Duration
}

func (e transientStatus) Error() string {
	return fmt.Sprintf("temporary archive HTTP status %d", e.code)
}

func archiveClient(input *http.Client) *http.Client {
	if input == nil {
		input = http.DefaultClient
	}
	c := *input
	previous := c.CheckRedirect
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if isAccountURL(req.URL) {
			return ErrReauthentication
		}
		if previous != nil {
			return previous(req, via)
		}
		if len(via) >= 10 {
			return errors.New("too many archive redirects")
		}
		return nil
	}
	return &c
}

func isAccountURL(u *url.URL) bool {
	if u == nil {
		return false
	}
	return strings.EqualFold(u.Hostname(), "accounts.google.com")
}

func validateLink(link string) error {
	u, err := url.Parse(link)
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") || u.Host == "" || u.User != nil {
		return errors.New("invalid archive download link")
	}
	if isAccountURL(u) {
		return ErrReauthentication
	}
	return nil
}

func (e *Engine) refresh(ctx context.Context, part model.Part, rejected string) (string, error) {
	e.linkMu.Lock()
	defer e.linkMu.Unlock()
	if link := e.links[part.Key()]; link != "" && link != rejected {
		return link, nil
	}
	e.publishPart(events.ReAuthRequired, part, "Refreshing archive authorization", 0)
	if e.opts.Links == nil {
		return "", ErrReauthentication
	}
	var link string
	var err error
	if source, ok := e.opts.Links.(interface {
		RefreshRejectedLink(context.Context, string, int, string) (string, error)
	}); ok {
		// Sources may refresh an entire multipart job at once. Supplying the
		// rejected signature lets a sibling reuse its renewed link without polling
		// again, while a rejected current signature still observes the cooldown.
		link, err = source.RefreshRejectedLink(ctx, part.ExportID, part.Index, rejected)
	} else {
		link, err = e.opts.Links.RefreshLink(ctx, part.ExportID, part.Index)
	}
	if err != nil {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		return "", &safeError{"archive authorization refresh failed", err}
	}
	if err = validateLink(link); err != nil {
		return "", err
	}
	e.links[part.Key()] = link
	return link, nil
}

func strongETag(tag string) bool {
	if len(tag) < 2 || tag[0] != '"' || tag[len(tag)-1] != '"' {
		return false
	}
	for i := 1; i < len(tag)-1; i++ {
		if tag[i] < 0x21 || tag[i] == '"' || tag[i] == 0x7f {
			return false
		}
	}
	return true
}

func responseRange(value string) (start, end, total int64, err error) {
	if !strings.HasPrefix(value, "bytes ") {
		return 0, 0, 0, ErrProtocol
	}
	x := strings.Split(strings.TrimPrefix(value, "bytes "), "/")
	if len(x) != 2 {
		return 0, 0, 0, ErrProtocol
	}
	y := strings.Split(x[0], "-")
	if len(y) != 2 {
		return 0, 0, 0, ErrProtocol
	}
	parse := func(s string) (int64, error) {
		if s == "" || strings.Trim(s, "0123456789") != "" {
			return 0, ErrProtocol
		}
		return strconv.ParseInt(s, 10, 64)
	}
	if start, err = parse(y[0]); err != nil {
		return
	}
	if end, err = parse(y[1]); err != nil {
		return
	}
	total, err = parse(x[1])
	if err != nil || start < 0 || end < start || total <= end {
		err = ErrProtocol
	}
	return
}

func (e *Engine) transfer(ctx context.Context, s *session, link string) error {
	if err := validateLink(link); err != nil {
		return err
	}
	part := s.record.Part
	offset := s.record.Offset
	// A strong ETag makes an offset request safe. Without one, one complete
	// response must prove the retained prefix before it can provide new bytes.
	ranged := offset > 0 && offset < part.Size && strongETag(s.record.ETag)
	start := int64(0)
	if ranged {
		start = offset
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		return &safeError{"cannot construct archive HTTP request", err}
	}
	req.Header.Set("Accept-Encoding", "identity")
	if ranged {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", offset))
		req.Header.Set("If-Range", s.record.ETag)
	}
	resp, err := e.client.Do(req)
	if err != nil {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
		if errors.Is(err, ErrReauthentication) {
			return ErrReauthentication
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return &safeError{"archive HTTP request failed", err}
	}
	defer resp.Body.Close()
	if (resp.Request != nil && isAccountURL(resp.Request.URL)) || resp.StatusCode == http.StatusUnauthorized {
		return ErrReauthentication
	}
	if resp.StatusCode == http.StatusForbidden {
		if strings.Contains(strings.ToLower(resp.Header.Get("Content-Type")), "json") {
			var detail struct {
				Error struct {
					Errors []struct {
						Reason string `json:"reason"`
					} `json:"errors"`
				} `json:"error"`
			}
			if json.NewDecoder(io.LimitReader(resp.Body, 64*1024)).Decode(&detail) == nil {
				for _, problem := range detail.Error.Errors {
					switch problem.Reason {
					case "rateLimitExceeded", "userRateLimitExceeded":
						return transientStatus{code: resp.StatusCode, after: retryDelay(resp.Header.Get("Retry-After"), time.Now())}
					case "downloadQuotaExceeded", "dailyLimitExceeded":
						return errors.New("Google download quota is exhausted; wait for recovery before retrying")
					case "insufficientFilePermissions", "cannotDownloadFile", "fileNotDownloadable":
						return errors.New("Google does not permit this file to be downloaded")
					}
				}
			}
		}
		return ErrReauthentication
	}
	if resp.StatusCode == http.StatusRequestedRangeNotSatisfiable {
		value := resp.Header.Get("Content-Range")
		if !strings.HasPrefix(value, "bytes */") {
			return fmt.Errorf("%w: malformed 416 range", ErrProtocol)
		}
		length := strings.TrimPrefix(value, "bytes */")
		total, err := strconv.ParseInt(length, 10, 64)
		if err != nil || strings.Trim(length, "0123456789") != "" || total != part.Size {
			return fmt.Errorf("%w: 416 total differs from archive size", ErrProtocol)
		}
		info, err := e.opts.Backend.Stat(ctx, part.StoragePath()+".part")
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if err == nil && info.Mode().IsRegular() && info.Size() == total && offset == total && !s.proveComplete {
			return nil // The caller still performs full archive integrity checks.
		}
		return fmt.Errorf("%w: 416 does not match a complete partial archive", ErrProtocol)
	}
	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
		return transientStatus{code: resp.StatusCode, after: retryDelay(resp.Header.Get("Retry-After"), time.Now())}
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("%w: HTTP status %d", ErrProtocol, resp.StatusCode)
	}
	if encoding := resp.Header.Get("Content-Encoding"); encoding != "" && !strings.EqualFold(encoding, "identity") {
		return fmt.Errorf("%w: response is content encoded", ErrProtocol)
	}
	if strings.Contains(strings.ToLower(resp.Header.Get("Content-Type")), "html") {
		return ErrReauthentication
	}
	r := bufio.NewReaderSize(resp.Body, 128*1024)
	peek, _ := r.Peek(512)
	lead := bytes.ToLower(bytes.TrimSpace(bytes.TrimPrefix(peek, []byte{0xef, 0xbb, 0xbf})))
	if bytes.HasPrefix(lead, []byte("<!doctype html")) || bytes.HasPrefix(lead, []byte("<html")) {
		return ErrReauthentication
	}
	if ranged && resp.StatusCode == http.StatusOK {
		return fmt.Errorf("%w: server ignored range; partial archive preserved", ErrProtocol)
	}
	if resp.StatusCode == http.StatusPartialContent {
		a, b, total, err := responseRange(resp.Header.Get("Content-Range"))
		if err != nil || a != start || b != part.Size-1 || total != part.Size {
			return fmt.Errorf("%w: Content-Range differs from requested offset or total", ErrProtocol)
		}
	} else if resp.Header.Get("Content-Range") != "" {
		return fmt.Errorf("%w: unexpected Content-Range on full response", ErrProtocol)
	}
	tag, modified := resp.Header.Get("ETag"), resp.Header.Get("Last-Modified")
	if strings.ContainsAny(tag+modified, "\r\n\x00") {
		return fmt.Errorf("%w: invalid archive validator", ErrProtocol)
	}
	if ranged && tag != s.record.ETag {
		return ErrIdentityChanged
	}
	if resp.ContentLength >= 0 && resp.ContentLength != part.Size-start {
		return fmt.Errorf("%w: Content-Length differs from exact archive size", ErrProtocol)
	}
	if offset > 0 && !ranged {
		if err = e.comparePrefix(ctx, s, r, offset); err != nil {
			return err
		}
	}
	// Record response identity only after any existing prefix has been proved.
	e.ioMu.Lock()
	s.record.ETag, s.record.LastModified = tag, modified
	err = e.saveCheckpoint(s)
	e.ioMu.Unlock()
	if err != nil {
		return err
	}
	buf := make([]byte, 128*1024)
	remaining := part.Size - offset
	emptyReads := 0
	for remaining > 0 {
		if err = e.waitReady(ctx); err != nil {
			return err
		}
		n, readErr := r.Read(buf[:min(int64(len(buf)), remaining)])
		if n > 0 {
			emptyReads = 0
			if err = e.write(ctx, s, buf[:n]); err != nil {
				return err
			}
			remaining -= int64(n)
		} else if readErr == nil {
			emptyReads++
			if emptyReads >= 100 {
				return io.ErrNoProgress
			}
		}
		if readErr != nil {
			if errors.Is(readErr, io.EOF) && remaining == 0 {
				return nil
			}
			if errors.Is(readErr, io.EOF) {
				return io.ErrUnexpectedEOF
			}
			return &safeError{"archive response stream failed", readErr}
		}
	}
	// Never write excess bytes, even for chunked responses with no length header.
	if _, err = r.ReadByte(); err == nil {
		return fmt.Errorf("%w: response exceeds exact archive size", ErrProtocol)
	}
	if !errors.Is(err, io.EOF) {
		return &safeError{"archive response did not finish cleanly", err}
	}
	return nil
}

func (e *Engine) comparePrefix(ctx context.Context, s *session, remote io.Reader, length int64) (retErr error) {
	f, err := e.read.OpenRead(ctx, s.name)
	if err != nil {
		return err
	}
	defer func() { retErr = errors.Join(retErr, f.Close()) }()
	localBuf, remoteBuf := make([]byte, 128*1024), make([]byte, 128*1024)
	for length > 0 {
		if err = e.waitReady(ctx); err != nil {
			return err
		}
		n := min(length, int64(len(localBuf)))
		if _, err = io.ReadFull(f, localBuf[:n]); err != nil {
			return fmt.Errorf("%w: partial prefix is unreadable", ErrIdentityChanged)
		}
		if _, err = io.ReadFull(remote, remoteBuf[:n]); err != nil {
			if errors.Is(err, io.EOF) {
				return io.ErrUnexpectedEOF
			}
			return &safeError{"archive prefix response failed", err}
		}
		if !bytes.Equal(localBuf[:n], remoteBuf[:n]) {
			return ErrIdentityChanged
		}
		length -= n
	}
	return nil
}

func retryable(err error) bool {
	var status transientStatus
	if errors.As(err, &status) {
		return true
	}
	if errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.ErrNoProgress) {
		return true
	}
	var network net.Error
	return errors.As(err, &network) && (network.Timeout() || network.Temporary())
}

// IsTemporary identifies failures a long-running watch can retry after backoff.
// Authorization, permission, protocol and integrity failures are terminal.
func IsTemporary(err error) bool { return retryable(err) }

func RetryDelay(err error) time.Duration {
	var provider interface{ RetryAfter() time.Duration }
	if errors.As(err, &provider) {
		return provider.RetryAfter()
	}
	var status transientStatus
	if errors.As(err, &status) {
		return status.after
	}
	return 0
}

func retryDelay(value string, now time.Time) time.Duration {
	if seconds, err := strconv.ParseInt(value, 10, 32); err == nil && seconds > 0 {
		return time.Duration(seconds) * time.Second
	}
	if date, err := http.ParseTime(value); err == nil && date.After(now) {
		return date.Sub(now)
	}
	return 0
}

func retryWait(ctx context.Context, attempt int, cause error) error {
	delay := min(time.Second*time.Duration(1<<min(attempt, 6)), time.Minute)
	delay = max(delay, RetryDelay(cause))
	if delay > time.Minute {
		return &safeError{"Google requested a longer retry delay; automatic retries stopped, retry later", cause}
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
