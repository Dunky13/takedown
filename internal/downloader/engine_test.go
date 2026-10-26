package downloader

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Dunky13/takedown/internal/events"
	"github.com/Dunky13/takedown/internal/model"
	"github.com/Dunky13/takedown/internal/space"
	"github.com/Dunky13/takedown/internal/state"
	"github.com/Dunky13/takedown/internal/storage"
)

func zipBytes(t *testing.T, size int) []byte {
	t.Helper()
	var out bytes.Buffer
	w := zip.NewWriter(&out)
	f, err := w.CreateHeader(&zip.FileHeader{Name: "Takeout/data.bin", Method: zip.Store})
	if err != nil {
		t.Fatal(err)
	}
	data := make([]byte, size)
	_, _ = rand.New(rand.NewSource(42)).Read(data)
	if _, err = f.Write(data); err != nil {
		t.Fatal(err)
	}
	if err = w.Close(); err != nil {
		t.Fatal(err)
	}
	return out.Bytes()
}

func gzipBytes(t *testing.T) []byte {
	t.Helper()
	var out bytes.Buffer
	w := gzip.NewWriter(&out)
	if _, err := w.Write([]byte(strings.Repeat("Takeout archive payload\n", 2000))); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return out.Bytes()
}

type fixture struct {
	engine    *Engine
	backend   *storage.Local
	store     *state.Store
	bus       *events.Bus
	dir       string
	statePath string
	part      model.Part
}

func newFixture(t *testing.T, data []byte, filename, link string, configure func(*Options)) *fixture {
	t.Helper()
	dir := t.TempDir()
	b, err := storage.NewLocal(filepath.Join(dir, "destination"))
	if err != nil {
		t.Fatal(err)
	}
	s, err := state.Open(filepath.Join(dir, "state.json"))
	if err != nil {
		t.Fatal(err)
	}
	bus := events.New()
	opts := Options{Backend: b, Destination: "fixture-destination", State: s, Bus: bus, Concurrency: 2, CheckInterval: 32 * 1024, MaxRetries: 1}
	if configure != nil {
		configure(&opts)
	}
	e, err := New(opts)
	if err != nil {
		t.Fatal(err)
	}
	e.spacePoll = 10 * time.Millisecond
	f := &fixture{e, b, s, bus, filepath.Join(dir, "destination"), filepath.Join(dir, "state.json"), model.Part{ExportID: "export", Index: 1, Filename: filename, Size: int64(len(data)), URL: link}}
	t.Cleanup(func() { _ = s.Close(); _ = b.Close(); bus.Close() })
	return f
}

func seedPartial(t *testing.T, f *fixture, data []byte, count int, tag string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(f.dir, f.part.Filename+".part"), data[:count], 0600); err != nil {
		t.Fatal(err)
	}
	if err := f.store.Put(state.Record{Part: f.part, Destination: "fixture-destination", Offset: int64(count), ETag: tag}); err != nil {
		t.Fatal(err)
	}
}

func serveArchive(w http.ResponseWriter, r *http.Request, data []byte, tag string) {
	w.Header().Set("Content-Type", "application/octet-stream")
	if tag != "" {
		w.Header().Set("ETag", tag)
	}
	start := int64(0)
	if value := r.Header.Get("Range"); value != "" {
		_, _ = fmt.Sscanf(value, "bytes=%d-", &start)
		if start >= int64(len(data)) {
			w.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", len(data)))
			w.WriteHeader(416)
			return
		}
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, len(data)-1, len(data)))
		w.Header().Set("Content-Length", strconv.Itoa(len(data)-int(start)))
		w.WriteHeader(206)
	} else {
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	}
	_, _ = w.Write(data[start:])
}

func assertCompleted(t *testing.T, f *fixture, data []byte) {
	t.Helper()
	got, err := os.ReadFile(filepath.Join(f.dir, f.part.Filename))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, data) {
		t.Fatal("completed archive differs from original")
	}
	if _, err = os.Stat(filepath.Join(f.dir, f.part.Filename+".part")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("partial remains: %v", err)
	}
	r, ok := f.store.Get(f.part.Key())
	if !ok || !r.Complete || r.Offset != int64(len(data)) {
		t.Fatalf("completion checkpoint: %+v", r)
	}
}

func TestRealArchivesAndExistingCompletion(t *testing.T) {
	for _, name := range []string{"takeout.zip", "takeout.tgz", "takeout.tar.gz"} {
		t.Run(name, func(t *testing.T) {
			data := gzipBytes(t)
			if strings.HasSuffix(name, ".zip") {
				data = zipBytes(t, 260_000)
			}
			var requests atomic.Int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests.Add(1)
				if r.Header.Get("Accept-Encoding") != "identity" {
					t.Error("identity encoding not requested")
				}
				serveArchive(w, r, data, `"archive-v1"`)
			}))
			defer srv.Close()
			f := newFixture(t, data, name, srv.URL, nil)
			if err := f.engine.Download(context.Background(), []model.Part{f.part}); err != nil {
				t.Fatal(err)
			}
			assertCompleted(t, f, data)
			if err := f.engine.Download(context.Background(), []model.Part{f.part}); err != nil {
				t.Fatal(err)
			}
			if requests.Load() != 1 {
				t.Fatal("verified completion was downloaded again")
			}
		})
	}
}

func TestStrongETagResume(t *testing.T) {
	data := zipBytes(t, 300_000)
	const prefix = 48_000
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Range") != fmt.Sprintf("bytes=%d-", prefix) || r.Header.Get("If-Range") != `"v1"` {
			t.Errorf("unsafe resume headers: %v", r.Header)
		}
		serveArchive(w, r, data, `"v1"`)
	}))
	defer srv.Close()
	f := newFixture(t, data, "archive.zip", srv.URL, nil)
	seedPartial(t, f, data, prefix, `"v1"`)
	if err := f.engine.Download(context.Background(), []model.Part{f.part}); err != nil {
		t.Fatal(err)
	}
	assertCompleted(t, f, data)
}

func TestValidatorlessResumeComparesSameResponsePrefix(t *testing.T) {
	data := zipBytes(t, 200_000)
	for _, changed := range []bool{false, true} {
		t.Run(fmt.Sprint(changed), func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Range") != "" {
					t.Error("validatorless resume must prove prefix in one response")
				}
				serveArchive(w, r, data, "")
			}))
			defer srv.Close()
			f := newFixture(t, data, "archive.zip", srv.URL, nil)
			old := bytes.Clone(data)
			if changed {
				old[100] ^= 0xff
			}
			seedPartial(t, f, old, 20_000, "")
			err := f.engine.Download(context.Background(), []model.Part{f.part})
			if changed {
				if !errors.Is(err, ErrIdentityChanged) {
					t.Fatalf("changed prefix: %v", err)
				}
				got, _ := os.ReadFile(filepath.Join(f.dir, f.part.Filename+".part"))
				if !bytes.Equal(got, old[:20_000]) {
					t.Fatal("mismatched prefix was changed")
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				assertCompleted(t, f, data)
			}
		})
	}
}

type linkFunc func(context.Context, string, int) (string, error)

func (f linkFunc) RefreshLink(ctx context.Context, id string, index int) (string, error) {
	return f(ctx, id, index)
}

func TestExpiredLinkRefreshAfterInterruptedStream(t *testing.T) {
	data := zipBytes(t, 400_000)
	var expired, fresh, refreshes atomic.Int32
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/expired" {
			expired.Add(1)
			if r.Header.Get("Range") != "" {
				w.WriteHeader(403)
				return
			}
			w.Header().Set("ETag", `"v1"`)
			w.Header().Set("Content-Length", strconv.Itoa(len(data)))
			_, _ = w.Write(data[:70_000])
			return // A short body becomes io.ErrUnexpectedEOF at the client.
		}
		fresh.Add(1)
		if r.Header.Get("Range") != "bytes=70000-" {
			t.Errorf("refreshed range: %q", r.Header.Get("Range"))
		}
		serveArchive(w, r, data, `"v1"`)
	}))
	defer srv.Close()
	f := newFixture(t, data, "archive.zip", srv.URL+"/expired?signature=DO-NOT-PERSIST", func(o *Options) {
		o.MaxRetries = 2
		o.Links = linkFunc(func(_ context.Context, id string, index int) (string, error) {
			refreshes.Add(1)
			if id != "export" || index != 1 {
				t.Error("wrong refresh identity")
			}
			return srv.URL + "/fresh?signature=ANOTHER-SECRET", nil
		})
	})
	if err := f.engine.Download(context.Background(), []model.Part{f.part}); err != nil {
		t.Fatal(err)
	}
	assertCompleted(t, f, data)
	if expired.Load() != 2 || fresh.Load() != 1 || refreshes.Load() != 1 {
		t.Fatalf("requests expired=%d fresh=%d refreshes=%d", expired.Load(), fresh.Load(), refreshes.Load())
	}
	disk, _ := os.ReadFile(f.statePath)
	if bytes.Contains(disk, []byte("signature")) {
		t.Fatal("signed link persisted")
	}
}

func TestAuthenticationResponsesRefreshAndKeepPartial(t *testing.T) {
	data := zipBytes(t, 80_000)
	for _, mode := range []string{"401", "403", "html-header", "html-sniff", "html-on-range-200", "redirect"} {
		t.Run(mode, func(t *testing.T) {
			var refreshed atomic.Int32
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/fresh" {
					serveArchive(w, r, data, `"v1"`)
					return
				}
				switch mode {
				case "401":
					w.WriteHeader(401)
				case "403":
					w.WriteHeader(403)
				case "html-header":
					w.Header().Set("Content-Type", "text/html")
					_, _ = io.WriteString(w, "<html>Sign in</html>")
				case "html-sniff":
					w.Header().Set("Content-Type", "application/octet-stream")
					w.Header().Set("ETag", `"v1"`)
					w.Header().Set("Content-Range", fmt.Sprintf("bytes 1000-%d/%d", len(data)-1, len(data)))
					w.WriteHeader(206)
					_, _ = io.WriteString(w, "<!doctype html><html>Sign in</html>")
				case "html-on-range-200":
					w.Header().Set("Content-Type", "application/octet-stream")
					_, _ = io.WriteString(w, "<html>Sign in</html>")
				case "redirect":
					http.Redirect(w, r, "https://accounts.google.com/ServiceLogin?secret=hidden", 302)
				}
			}))
			defer srv.Close()
			f := newFixture(t, data, "archive.zip", srv.URL+"/expired", func(o *Options) {
				o.Links = linkFunc(func(context.Context, string, int) (string, error) { refreshed.Add(1); return srv.URL + "/fresh", nil })
			})
			seedPartial(t, f, data, 1000, `"v1"`)
			if err := f.engine.Download(context.Background(), []model.Part{f.part}); err != nil {
				t.Fatal(err)
			}
			if refreshed.Load() != 1 {
				t.Fatal("authentication was not refreshed once")
			}
			assertCompleted(t, f, data)
		})
	}
}

func TestRangeProtocolFailuresPreservePartial(t *testing.T) {
	data := zipBytes(t, 80_000)
	for _, mode := range []string{"ignored", "start", "total", "end", "wildcard", "etag", "missing-etag", "encoding", "oversized"} {
		t.Run(mode, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if mode == "ignored" {
					serveArchive(w, &http.Request{Header: make(http.Header)}, data, `"v1"`)
					return
				}
				start, end, total := 1000, len(data)-1, len(data)
				switch mode {
				case "start":
					start++
				case "total":
					total++
				case "end":
					end--
				}
				w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, total))
				if mode == "wildcard" {
					w.Header().Set("Content-Range", "bytes 1000-79999/*")
				}
				w.Header().Set("ETag", `"v1"`)
				if mode == "etag" {
					w.Header().Set("ETag", `"v2"`)
				}
				if mode == "missing-etag" {
					w.Header().Del("ETag")
				}
				if mode == "encoding" {
					w.Header().Set("Content-Encoding", "gzip")
				}
				body := data[1000:]
				if mode == "oversized" {
					body = append(bytes.Clone(body), 1)
				}
				w.Header().Set("Content-Length", strconv.Itoa(len(body)))
				w.WriteHeader(206)
				_, _ = w.Write(body)
			}))
			defer srv.Close()
			f := newFixture(t, data, "archive.zip", srv.URL, nil)
			seedPartial(t, f, data, 1000, `"v1"`)
			err := f.engine.Download(context.Background(), []model.Part{f.part})
			if !errors.Is(err, ErrProtocol) && !errors.Is(err, ErrIdentityChanged) {
				t.Fatalf("protocol response accepted: %v", err)
			}
			got, _ := os.ReadFile(filepath.Join(f.dir, f.part.Filename+".part"))
			if !bytes.Equal(got, data[:1000]) {
				t.Fatal("protocol failure changed existing partial")
			}
			if _, err = os.Stat(filepath.Join(f.dir, f.part.Filename)); !errors.Is(err, os.ErrNotExist) {
				t.Fatal("invalid response published")
			}
		})
	}
}

func TestCorruptArchivesAndChunkedExcessNeverPublished(t *testing.T) {
	for _, mode := range []string{"zip-crc", "zip-eocd", "zip-trailing", "gzip-checksum", "gzip-trailer", "chunked-excess"} {
		t.Run(mode, func(t *testing.T) {
			data := zipBytes(t, 100_000)
			name := "archive.zip"
			switch mode {
			case "zip-crc":
				data[100] ^= 0xff
			case "zip-eocd":
				data[len(data)-22] ^= 0xff
			case "zip-trailing":
				data = append(data, 1)
			case "gzip-checksum":
				data = gzipBytes(t)
				name = "archive.tgz"
				data[len(data)-8] ^= 0xff
			case "gzip-trailer":
				data = gzipBytes(t)
				name = "archive.tar.gz"
				data = data[:len(data)-2]
			}
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if mode == "chunked-excess" {
					w.WriteHeader(200)
					w.(http.Flusher).Flush()
					_, _ = w.Write(data)
					_, _ = w.Write([]byte("unexpected"))
					return
				}
				serveArchive(w, r, data, `"v1"`)
			}))
			defer srv.Close()
			f := newFixture(t, data, name, srv.URL, nil)
			err := f.engine.Download(context.Background(), []model.Part{f.part})
			if !errors.Is(err, ErrIntegrity) && !errors.Is(err, ErrProtocol) {
				t.Fatalf("invalid archive result: %v", err)
			}
			if _, err = os.Stat(filepath.Join(f.dir, name)); !errors.Is(err, os.ErrNotExist) {
				t.Fatal("corrupt archive published")
			}
			if info, err := os.Stat(filepath.Join(f.dir, name+".part")); err != nil || info.Size() != int64(len(data)) {
				t.Fatalf("partial not preserved at exact length: %v %v", info, err)
			}
			r, _ := f.store.Get(f.part.Key())
			if r.Complete {
				t.Fatal("corrupt archive marked complete")
			}
		})
	}
}

func TestCancelAndReopenResume(t *testing.T) {
	data := zipBytes(t, 500_000)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Range") != "" {
			serveArchive(w, r, data, `"v1"`)
			return
		}
		w.Header().Set("ETag", `"v1"`)
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		_, _ = w.Write(data[:64_000])
		w.(http.Flusher).Flush()
		<-r.Context().Done()
	}))
	defer srv.Close()
	f := newFixture(t, data, "archive.zip", srv.URL, nil)
	stream, stop := f.bus.Subscribe(128)
	defer stop()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result := make(chan error, 1)
	go func() { result <- f.engine.Download(ctx, []model.Part{f.part}) }()
	for ev := range stream {
		if ev.Kind == events.Progress {
			cancel()
			break
		}
	}
	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("cancel: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("cancellation did not join workers")
	}
	info, err := os.Stat(filepath.Join(f.dir, f.part.Filename+".part"))
	if err != nil || info.Size() <= 0 {
		t.Fatalf("partial after cancel: %v", err)
	}
	r, _ := f.store.Get(f.part.Key())
	if r.Offset != info.Size() || r.Complete {
		t.Fatalf("cancel checkpoint: %+v size=%d", r, info.Size())
	}
	if err = f.store.Close(); err != nil {
		t.Fatal(err)
	}
	s, err := state.Open(f.statePath)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	f.store = s
	e, err := New(Options{Backend: f.backend, Destination: "fixture-destination", State: s, Bus: f.bus, MaxRetries: 1})
	if err != nil {
		t.Fatal(err)
	}
	if err = e.Download(context.Background(), []model.Part{f.part}); err != nil {
		t.Fatal(err)
	}
	assertCompleted(t, f, data)
}

type controlledBackend struct {
	*storage.Local
	free       atomic.Uint64
	spaceCalls atomic.Int64
	bytes      atomic.Int64
	open       atomic.Int32
}

func (b *controlledBackend) FreeSpace(context.Context) (uint64, error) {
	b.spaceCalls.Add(1)
	return b.free.Load(), nil
}
func (b *controlledBackend) OpenWriteResume(ctx context.Context, name string, offset int64) (io.WriteCloser, error) {
	w, err := b.Local.OpenWriteResume(ctx, name, offset)
	if err != nil {
		return nil, err
	}
	b.open.Add(1)
	return &countWriter{w: w, b: b}, nil
}

type countWriter struct {
	w    io.WriteCloser
	b    *controlledBackend
	once sync.Once
	err  error
}

func (w *countWriter) Write(p []byte) (int, error) {
	n, err := w.w.Write(p)
	w.b.bytes.Add(int64(n))
	return n, err
}
func (w *countWriter) Close() error {
	w.once.Do(func() { w.err = w.w.Close(); w.b.open.Add(-1) })
	return w.err
}

func TestLowSpacePausesAllWritersThenResumes(t *testing.T) {
	data := zipBytes(t, 1_000_000)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", `"v1"`)
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		for start := 0; start < len(data); start += 16_000 {
			if _, err := w.Write(data[start:min(start+16_000, len(data))]); err != nil {
				return
			}
			w.(http.Flusher).Flush()
			select {
			case <-r.Context().Done():
				return
			case <-time.After(time.Millisecond):
			}
		}
	}))
	defer srv.Close()
	var backend *controlledBackend
	f := newFixture(t, data, "one.zip", srv.URL, func(o *Options) {
		backend = &controlledBackend{Local: o.Backend.(*storage.Local)}
		backend.free.Store(10_000_000)
		o.Backend = backend
		o.Guard = space.New(backend, 100, nil)
	})
	part2 := f.part
	part2.Index = 2
	part2.Filename = "two.zip"
	stream, stop := f.bus.Subscribe(2048)
	defer stop()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result := make(chan error, 1)
	go func() { result <- f.engine.Download(ctx, []model.Part{f.part, part2}) }()
	dropped, paused := false, false
	for !paused {
		select {
		case ev := <-stream:
			if ev.Kind == events.Progress && !dropped {
				backend.free.Store(0)
				dropped = true
			}
			if ev.Kind == events.SpaceWarning && strings.Contains(ev.Message, "all workers paused") {
				paused = true
			}
		case err := <-result:
			t.Fatalf("finished before low-space pause: %v", err)
		case <-ctx.Done():
			t.Fatal("no low-space pause")
		}
	}
	if backend.open.Load() != 0 {
		t.Fatal("writers still open after pause checkpoint")
	}
	before := backend.bytes.Load()
	for _, r := range f.store.Snapshot() {
		info, err := os.Stat(filepath.Join(f.dir, r.Part.Filename+".part"))
		if err == nil && r.Offset != info.Size() {
			t.Fatalf("unflushed checkpoint: %+v", r)
		}
	}
	time.Sleep(30 * time.Millisecond)
	if backend.bytes.Load() != before {
		t.Fatal("bytes written during low-space pause")
	}
	backend.free.Store(10_000_000)
	f.engine.Resume()
	select {
	case err := <-result:
		if err != nil {
			t.Fatal(err)
		}
	case <-ctx.Done():
		t.Fatal("space recovery did not resume")
	}
	assertCompleted(t, f, data)
	got, err := os.ReadFile(filepath.Join(f.dir, part2.Filename))
	if err != nil || !bytes.Equal(got, data) {
		t.Fatalf("second worker completion: %v", err)
	}
	if backend.bytes.Load() != int64(2*len(data)) {
		t.Fatal("resume discarded or duplicated partial bytes")
	}
}

func TestManualPauseFlushesAndCancelWakes(t *testing.T) {
	data := zipBytes(t, 300_000)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", `"v1"`)
		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		_, _ = w.Write(data[:64_000])
		w.(http.Flusher).Flush()
		<-r.Context().Done()
	}))
	defer srv.Close()
	f := newFixture(t, data, "archive.zip", srv.URL, nil)
	stream, stop := f.bus.Subscribe(128)
	defer stop()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result := make(chan error, 1)
	go func() { result <- f.engine.Download(ctx, []model.Part{f.part}) }()
	for ev := range stream {
		if ev.Kind == events.Progress {
			break
		}
	}
	f.engine.Pause()
	r, _ := f.store.Get(f.part.Key())
	info, err := os.Stat(filepath.Join(f.dir, f.part.Filename+".part"))
	if err != nil || r.Offset != info.Size() {
		t.Fatalf("pause checkpoint: %+v %v", r, err)
	}
	cancel()
	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("cancel paused run: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("paused cancel blocked")
	}
}

func TestBoundedRetriesAndContextBackoff(t *testing.T) {
	data := zipBytes(t, 100)
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1); w.WriteHeader(503) }))
	defer srv.Close()
	f := newFixture(t, data, "archive.zip", srv.URL, func(o *Options) { o.MaxRetries = 2 })
	if err := f.engine.Download(context.Background(), []model.Part{f.part}); err == nil {
		t.Fatal("transient failures accepted")
	}
	if calls.Load() != 3 {
		t.Fatalf("retry bound: %d requests", calls.Load())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	start := time.Now()
	err := f.engine.Download(ctx, []model.Part{f.part})
	if !errors.Is(err, context.DeadlineExceeded) || time.Since(start) > time.Second {
		t.Fatalf("context backoff: %v", err)
	}
}

func TestDestinationChangeDoesNotReuseValidator(t *testing.T) {
	data := zipBytes(t, 30_000)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Range") != "" || r.Header.Get("If-Range") != "" {
			t.Error("old destination identity reused")
		}
		serveArchive(w, r, data, `"new"`)
	}))
	defer srv.Close()
	f := newFixture(t, data, "archive.zip", srv.URL, nil)
	if err := f.store.Put(state.Record{Part: f.part, Destination: "old-destination", Offset: 20_000, ETag: `"old"`}); err != nil {
		t.Fatal(err)
	}
	if err := f.engine.Download(context.Background(), []model.Part{f.part}); err != nil {
		t.Fatal(err)
	}
	assertCompleted(t, f, data)
}
