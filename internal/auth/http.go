package auth

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

func secureURL(u *url.URL) bool {
	return u != nil && u.Scheme == "https" && u.Hostname() != "" && u.User == nil && u.Opaque == "" && u.Fragment == "" && (u.Port() == "" || u.Port() == "443") && !strings.HasSuffix(u.Host, ":")
}

func apiHost(host string) bool {
	switch host {
	case "www.googleapis.com", "drive.googleapis.com", "dataportability.googleapis.com":
		return true
	}
	return false
}

func downloadHost(host string) bool {
	return host == "storage.googleapis.com" || strings.HasSuffix(host, ".storage.googleapis.com") || strings.HasSuffix(host, ".googleusercontent.com") || host == "drive.usercontent.google.com" || strings.HasSuffix(host, ".usercontent.google.com")
}

func allowedRequest(r *http.Request) bool {
	return secureURL(r.URL) && (r.Host == "" || r.Host == r.URL.Host) && (apiHost(r.URL.Hostname()) || downloadHost(r.URL.Hostname()))
}

func stripCredentials(header http.Header) {
	for name := range header {
		switch strings.ToLower(name) {
		case "authorization", "proxy-authorization", "cookie", "cookie2", "referer":
			delete(header, name)
		}
	}
}

func checkRedirect(r *http.Request, via []*http.Request) error {
	stripCredentials(r.Header)
	if len(via) > 0 && via[0].Method != http.MethodGet && via[0].Method != http.MethodHead {
		return errors.New("Google API/download redirect refused for a mutation; verify the endpoint before starting a new request")
	}
	if len(via) >= 10 || !allowedRequest(r) {
		return errors.New("Google API/download redirect refused: HTTPS, an allowed host, and port 443 are required (maximum 10 redirects)")
	}
	return nil
}

type apiTransport struct{ manager *Manager }

func (t *apiTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	ownedBody := true
	defer func() {
		if ownedBody && request.Body != nil {
			_ = request.Body.Close()
		}
	}()
	if !allowedRequest(request) {
		return nil, errors.New("Google API/download request refused: HTTPS, an allowed host, and port 443 are required")
	}
	m := t.manager
	ctx, cleanup, err := m.operation(request.Context())
	if err != nil {
		return nil, err
	}
	request = request.Clone(ctx)
	stripCredentials(request.Header)
	access, revision, err := m.authorizeRequest(ctx, request, false, "", 0)
	if err != nil {
		cleanup()
		return nil, err
	}
	for attempt := 0; attempt < 2; attempt++ {
		ownedBody = false // the underlying RoundTripper now owns request.Body
		response, err := m.base.RoundTrip(request)
		if err != nil {
			contextErr := ctx.Err()
			cleanup()
			if contextErr != nil {
				return nil, contextErr
			}
			if temporary := temporaryNetwork(err); temporary != nil {
				return nil, temporary
			}
			return nil, errors.New("Google API/download request failed; check the connection and retry")
		}
		api := apiHost(request.URL.Hostname())
		unauthorized := api && response.StatusCode == http.StatusUnauthorized
		insufficient := api && response.StatusCode == http.StatusForbidden && strings.Contains(response.Header.Get("WWW-Authenticate"), "insufficient_scope")
		if unauthorized || insufficient {
			if response.Body != nil {
				_ = response.Body.Close()
			}
			read := request.Method == http.MethodGet || request.Method == http.MethodHead
			replayable := request.Body == nil || request.Body == http.NoBody || request.GetBody != nil
			if unauthorized && read && replayable && attempt == 0 {
				access, revision, err = m.authorizeRequest(ctx, request, true, access, revision)
				if err != nil {
					cleanup()
					return nil, err
				}
				if request.Body != nil && request.Body != http.NoBody {
					request.Body, err = request.GetBody()
					ownedBody = true
					if err != nil {
						cleanup()
						return nil, errors.New("cannot replay the authorized Google API read")
					}
				}
				continue
			}
			if insufficient || attempt == 1 {
				m.rejectToken(ctx, access, revision)
			}
			cleanup()
			return nil, m.needAuth("Google rejected the API authorization; run Login and grant the required scopes")
		}
		if response.Body == nil {
			cleanup()
		} else {
			response.Body = &operationBody{ReadCloser: response.Body, ctx: ctx, cleanup: cleanup}
		}
		return response, nil
	}
	cleanup()
	return nil, errors.New("Google API authorization read replay limit reached")
}

func (m *Manager) authorizeRequest(ctx context.Context, request *http.Request, refresh bool, rejectedAccess string, rejectedRevision uint64) (string, uint64, error) {
	_, release, err := m.acquire(ctx)
	if err != nil {
		return "", 0, err
	}
	defer release()
	if !apiHost(request.URL.Hostname()) {
		return "", 0, m.ensureConfig()
	}
	force := refresh && m.token != nil && m.token.AccessToken == rejectedAccess && m.revision.Load() == rejectedRevision
	token, err := m.tokenLocked(ctx, force)
	if err != nil {
		return "", 0, err
	}
	token.SetAuthHeader(request)
	return token.AccessToken, m.revision.Load(), nil
}

func (m *Manager) rejectToken(ctx context.Context, access string, revision uint64) {
	_, release, err := m.acquire(ctx)
	if err != nil {
		return
	}
	defer release()
	if m.token != nil && m.token.AccessToken == access && m.revision.Load() == revision {
		m.required = true
	}
}

type operationBody struct {
	io.ReadCloser
	once    sync.Once
	ctx     context.Context
	cleanup func()
}

func (b *operationBody) Read(p []byte) (int, error) {
	n, err := b.ReadCloser.Read(p)
	if err != nil {
		if err != io.EOF {
			if b.ctx.Err() != nil {
				err = b.ctx.Err()
			} else if temporary := temporaryNetwork(err); temporary != nil {
				err = temporary
			} else {
				err = errors.New("Google API/download response was interrupted")
			}
		}
		b.once.Do(b.cleanup)
	}
	return n, err
}

func (b *operationBody) Close() error {
	err := b.ReadCloser.Close()
	b.once.Do(b.cleanup)
	if err != nil {
		return errors.New("cannot close the Google API/download response")
	}
	return nil
}

// Compile-time guard: the transport intentionally does not use a cookie jar or
// an oauth2.Transport that would attach tokens indiscriminately to download URLs.
var _ http.RoundTripper = (*apiTransport)(nil)
