package source

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

func validSignedURL(raw string) (*url.URL, error) {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.User != nil || u.Fragment != "" || u.Port() != "" || u.Opaque != "" || u.Host != strings.ToLower(u.Host) {
		return nil, errors.New("Google returned an unsafe archive URL")
	}
	host := u.Hostname()
	if host != "storage.googleapis.com" && !strings.HasSuffix(host, ".storage.googleapis.com") {
		return nil, errors.New("Google returned an archive URL outside documented Cloud Storage HTTPS hosts")
	}
	if u.Path == "" || u.Path == "/" {
		return nil, errors.New("Google returned an archive URL without an object path")
	}
	return u, nil
}

func signedClient() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DisableCompression = true
	transport.ResponseHeaderTimeout = 30 * time.Second
	return &http.Client{Transport: transport, Timeout: 30 * time.Second, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 5 {
			return errors.New("archive redirect limit reached")
		}
		if _, err := validSignedURL(req.URL.String()); err != nil {
			return err
		}
		req.Header.Del("Authorization")
		req.Header.Del("Cookie")
		return nil
	}}
}

type archiveMetadata struct {
	Object     string `json:"object"` // hash of stable object URL, without signatures
	Filename   string `json:"filename"`
	Size       int64  `json:"size"`
	MD5        string `json:"md5,omitempty"`
	Generation string `json:"generation,omitempty"`
}

func signedChecksum(h http.Header) (string, error) {
	encoded := h.Get("Content-MD5")
	for _, value := range h.Values("X-Goog-Hash") {
		for _, item := range strings.Split(value, ",") {
			item = strings.TrimSpace(item)
			if strings.HasPrefix(item, "md5=") {
				md5 := strings.TrimPrefix(item, "md5=")
				if encoded != "" && encoded != md5 {
					return "", errors.New("archive checksum metadata conflicts")
				}
				encoded = md5
			}
		}
	}
	if encoded == "" {
		return "", nil
	}
	digest, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil || len(digest) != 16 {
		return "", errors.New("archive checksum metadata is invalid")
	}
	return hex.EncodeToString(digest), nil
}

func objectIdentity(u *url.URL) string {
	// Generation is immutable object identity; signing timestamps and signatures
	// are not. Never persist the signed URL itself.
	stable := u.Scheme + "://" + u.Host + u.EscapedPath() + "?generation=" + u.Query().Get("generation")
	sum := sha256.Sum256([]byte(stable))
	return hex.EncodeToString(sum[:])
}

func filenameFrom(h http.Header, u *url.URL) (string, error) {
	if disposition := h.Get("Content-Disposition"); disposition != "" {
		_, params, err := mime.ParseMediaType(disposition)
		if err != nil {
			return "", errors.New("archive Content-Disposition is malformed")
		}
		if name := params["filename"]; name != "" {
			return name, nil
		}
	}
	name := path.Base(u.Path)
	if !archiveName(name) {
		return "", errors.New("archive response has no safe ZIP/TGZ filename; refusing to invent one")
	}
	return name, nil
}

func safeContent(h http.Header) error {
	if encoding := h.Get("Content-Encoding"); encoding != "" && !strings.EqualFold(encoding, "identity") {
		return errors.New("archive metadata response is content encoded")
	}
	if value := h.Get("Content-Type"); value != "" {
		kind, _, err := mime.ParseMediaType(value)
		if err != nil {
			return errors.New("archive Content-Type is malformed")
		}
		switch strings.ToLower(kind) {
		case "application/zip", "application/x-zip-compressed", "application/gzip", "application/x-gzip", "application/octet-stream", "application/x-tar":
		default:
			return errors.New("archive response is not an archive content type")
		}
	}
	return nil
}

// Signed Cloud Storage URLs may bind the HTTP method to GET. Use that
// method directly; a HEAD request can fail even while its GET signature is valid.
func (m *Manager) probeSigned(ctx context.Context, raw string) (archiveMetadata, error) {
	u, err := validSignedURL(raw)
	if err != nil {
		return archiveMetadata{}, err
	}
	if err = ctx.Err(); err != nil {
		return archiveMetadata{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, raw, nil)
	if err != nil {
		return archiveMetadata{}, errors.New("cannot construct archive metadata request")
	}
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("Range", "bytes=0-0")
	resp, err := m.public.Do(req)
	if err != nil {
		if resp != nil && resp.Body != nil {
			_ = resp.Body.Close()
		}
		return archiveMetadata{}, apiError(ctx, "read unauthenticated archive metadata", err)
	}
	// Close immediately, without reading or draining a server's full archive
	// response if Range was ignored. The requested metadata body is at most one byte.
	_ = resp.Body.Close()
	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		return archiveMetadata{}, errors.New("signed archive link was denied or expired; refresh the same API job (links last six hours, data lasts fourteen days)")
	}
	if resp.StatusCode != http.StatusPartialContent {
		if resp.StatusCode != http.StatusOK {
			return archiveMetadata{}, signedStatus(resp.StatusCode, resp.Header)
		}
		return archiveMetadata{}, errors.New("archive server did not honor the one-byte metadata Range request")
	}
	if err = safeContent(resp.Header); err != nil {
		return archiveMetadata{}, err
	}
	value := resp.Header.Get("Content-Range")
	if !strings.HasPrefix(value, "bytes 0-0/") {
		return archiveMetadata{}, errors.New("archive metadata Range response has an invalid Content-Range")
	}
	total := strings.TrimPrefix(value, "bytes 0-0/")
	if total == "" || strings.Trim(total, "0123456789") != "" {
		return archiveMetadata{}, errors.New("archive metadata has no exact byte size")
	}
	size, err := strconv.ParseInt(total, 10, 64)
	if err != nil || size <= 0 || resp.ContentLength != 1 {
		return archiveMetadata{}, errors.New("archive metadata Range response has an invalid exact size")
	}
	name, err := filenameFrom(resp.Header, u)
	if err != nil {
		return archiveMetadata{}, err
	}
	checksum, err := signedChecksum(resp.Header)
	if err != nil {
		return archiveMetadata{}, err
	}
	generation := resp.Header.Get("X-Goog-Generation")
	if generation != "" && (len(generation) > 40 || strings.Trim(generation, "0123456789") != "") {
		return archiveMetadata{}, errors.New("archive object generation metadata is invalid")
	}
	if err = validateArchive(name, size, checksum); err != nil {
		return archiveMetadata{}, err
	}
	return archiveMetadata{Object: objectIdentity(u), Filename: name, Size: size, MD5: checksum, Generation: generation}, nil
}

func signedStatus(code int, h http.Header) error {
	message := fmt.Sprintf("archive metadata request failed (HTTP %d); retry later or refresh the API job links", code)
	if delay := retryAfter(h.Get("Retry-After"), time.Now()); delay > 0 {
		message += fmt.Sprintf("; retry after %s", delay.Round(time.Second))
	}
	return &sourceError{message: message, temporary: code == 429 || code >= 500 && code <= 599, delay: retryAfter(h.Get("Retry-After"), time.Now())}
}
