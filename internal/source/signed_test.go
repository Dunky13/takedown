package source

import (
	"context"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type trackedBody struct{ read, closed bool }

func (b *trackedBody) Read([]byte) (int, error) {
	b.read = true
	return 0, errors.New("must not read or drain archive probe body")
}
func (b *trackedBody) Close() error { b.closed = true; return nil }

func TestSignedGETMethodBoundSignature(t *testing.T) {
	m := testManager(t, Options{}, func(r *http.Request) (*http.Response, error) {
		t.Fatal("authenticated client used for signed metadata")
		return nil, nil
	})
	methods := []string{}
	installSigned(t, m, func(r *http.Request) (*http.Response, error) {
		methods = append(methods, r.Method)
		// Models a URL signed for GET. Its HEAD signature is invalid.
		if r.Method == http.MethodHead {
			return &http.Response{StatusCode: 403, Header: http.Header{}, Body: io.NopCloser(strings.NewReader("SignatureDoesNotMatch")), Request: r}, nil
		}
		resp := signedResponse(r, "export.zip", 987654321)
		resp.Header.Set("X-Goog-Hash", "crc32c=AAAAAA==,md5="+base64.StdEncoding.EncodeToString(make([]byte, 16)))
		return resp, nil
	})
	meta, err := m.probeSigned(context.Background(), "https://storage.googleapis.com/bucket/object?X-Goog-Signature=secret")
	if err != nil || meta.Size != 987654321 || meta.Filename != "export.zip" || meta.MD5 != strings.Repeat("0", 32) {
		t.Fatal(meta, err)
	}
	if len(methods) != 1 || methods[0] != http.MethodGet {
		t.Fatal("method-bound signature was probed with HEAD", methods)
	}
}

func TestSignedMetadataFailsClosedAndClosesBody(t *testing.T) {
	for _, tc := range []struct {
		name   string
		change func(*http.Response)
	}{
		{"full response", func(r *http.Response) { r.StatusCode = 200 }},
		{"html", func(r *http.Response) { r.Header.Set("Content-Type", "text/html") }},
		{"encoded", func(r *http.Response) { r.Header.Set("Content-Encoding", "gzip") }},
		{"unknown size", func(r *http.Response) { r.Header.Set("Content-Range", "bytes 0-0/*") }},
		{"zero size", func(r *http.Response) { r.Header.Set("Content-Range", "bytes 0-0/0") }},
		{"wrong range", func(r *http.Response) { r.Header.Set("Content-Range", "bytes 1-1/100") }},
		{"overflow", func(r *http.Response) { r.Header.Set("Content-Range", "bytes 0-0/9999999999999999999999") }},
		{"body size mismatch", func(r *http.Response) { r.ContentLength = 100 }},
		{"unknown body size", func(r *http.Response) { r.ContentLength = -1 }},
		{"unsafe filename", func(r *http.Response) { r.Header.Set("Content-Disposition", `attachment; filename="../export.zip"`) }},
		{"invented filename", func(r *http.Response) { r.Header.Del("Content-Disposition") }},
		{"invalid checksum", func(r *http.Response) { r.Header.Set("Content-MD5", "bad") }},
		{"invalid generation", func(r *http.Response) { r.Header.Set("X-Goog-Generation", "bad") }},
		{"forbidden", func(r *http.Response) { r.StatusCode = 403 }},
		{"quota transient", func(r *http.Response) { r.StatusCode = 429; r.Header.Set("Retry-After", "3600") }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := &trackedBody{}
			m := testManager(t, Options{}, func(r *http.Request) (*http.Response, error) { t.Fatal("OAuth metadata request"); return nil, nil })
			installSigned(t, m, func(r *http.Request) (*http.Response, error) {
				resp := signedResponse(r, "export.zip", 100)
				tc.change(resp)
				resp.Body = body
				return resp, nil
			})
			_, err := m.probeSigned(context.Background(), "https://storage.googleapis.com/bucket/object?secret=token")
			if err == nil || body.read || !body.closed {
				t.Fatal("bad metadata accepted or body consumed", err, body)
			}
			if strings.Contains(err.Error(), "token") {
				t.Fatal("signed URL leaked")
			}
			if tc.name == "quota transient" {
				if !IsTemporary(err) || RetryDelay(err) != time.Hour {
					t.Fatal("metadata quota lost retry delay", err)
				}
			} else if IsTemporary(err) {
				t.Fatal("malformed metadata or permission failure marked temporary", err)
			}
		})
	}
}

func TestSignedHostPolicyAndFilenameFromURL(t *testing.T) {
	for _, raw := range []string{"http://storage.googleapis.com/b/a.zip", "https://storage.googleapis.com.evil/b/a.zip", "https://user@storage.googleapis.com/b/a.zip", "https://storage.googleapis.com:443/b/a.zip", "https://evil.invalid/a.zip", "https://storage.googleapis.com/b/a.zip#secret"} {
		if _, err := validSignedURL(raw); err == nil {
			t.Errorf("accepted %s", raw)
		}
	}
	m := testManager(t, Options{}, func(r *http.Request) (*http.Response, error) { t.Fatal("OAuth metadata request"); return nil, nil })
	installSigned(t, m, func(r *http.Request) (*http.Response, error) {
		resp := signedResponse(r, "unused.zip", 100)
		resp.Header.Del("Content-Disposition")
		return resp, nil
	})
	meta, err := m.probeSigned(context.Background(), "https://bucket.storage.googleapis.com/nested/export.tar.gz?signature=secret")
	if err != nil || meta.Filename != "export.tar.gz" {
		t.Fatal(meta, err)
	}
	client := signedClient()
	defer client.CloseIdleConnections()
	request, _ := http.NewRequest(http.MethodGet, "https://evil.invalid/", nil)
	if err = client.CheckRedirect(request, []*http.Request{{}}); err == nil {
		t.Fatal("external redirect accepted")
	}
	request, _ = http.NewRequest(http.MethodGet, "https://storage.googleapis.com/b/a.zip", nil)
	request.Header.Set("Authorization", "Bearer secret")
	request.Header.Set("Cookie", "secret=value")
	if err = client.CheckRedirect(request, []*http.Request{{}}); err != nil || request.Header.Get("Authorization") != "" || request.Header.Get("Cookie") != "" {
		t.Fatal("redirect retained credentials", err)
	}
}
