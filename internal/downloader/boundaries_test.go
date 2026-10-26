package downloader

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Dunky13/takedown/internal/model"
	"github.com/Dunky13/takedown/internal/space"
	"github.com/Dunky13/takedown/internal/storage"
)

type intervalBackend struct {
	*controlledBackend
	mu           sync.Mutex
	observations []int64
}

func (b *intervalBackend) FreeSpace(ctx context.Context) (uint64, error) {
	b.mu.Lock()
	b.observations = append(b.observations, b.bytes.Load())
	b.mu.Unlock()
	return b.controlledBackend.FreeSpace(ctx)
}

func TestSpaceQueriesAtPartAdmissionAndByteIntervals(t *testing.T) {
	data := zipBytes(t, 1_000_000)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { serveArchive(w, r, data, `"v1"`) }))
	defer srv.Close()
	var b *intervalBackend
	const interval = 300_000
	f := newFixture(t, data, "archive.zip", srv.URL, func(o *Options) {
		b = &intervalBackend{controlledBackend: &controlledBackend{Local: o.Backend.(*storage.Local)}}
		b.free.Store(10_000_000)
		o.Backend = b
		o.Guard = space.New(b, 0, nil)
		o.CheckInterval = interval
	})
	if err := f.engine.Download(context.Background(), []model.Part{f.part}); err != nil {
		t.Fatal(err)
	}
	assertCompleted(t, f, data)
	// Split writes at the check boundary even when a 128 KiB HTTP read spans it.
	want := []int64{0, interval, 2 * interval, 3 * interval}
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.observations) != len(want) {
		t.Fatalf("free-space RPCs: %v; want %v", b.observations, want)
	}
	for i, v := range want {
		if b.observations[i] != v {
			t.Fatalf("free-space RPC offsets: %v; want %v", b.observations, want)
		}
	}
}

func TestDefaultSpaceIntervalAvoidsPerChunkQueries(t *testing.T) {
	data := zipBytes(t, 500_000)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { serveArchive(w, r, data, `"v1"`) }))
	defer srv.Close()
	var b *controlledBackend
	f := newFixture(t, data, "archive.zip", srv.URL, func(o *Options) {
		b = &controlledBackend{Local: o.Backend.(*storage.Local)}
		b.free.Store(10_000_000)
		o.Backend = b
		o.Guard = space.New(b, 0, nil)
		o.CheckInterval = 0
		o.Concurrency = 0
	})
	if f.engine.opts.Concurrency != 2 {
		t.Fatalf("default concurrency: %d", f.engine.opts.Concurrency)
	}
	if err := f.engine.Download(context.Background(), []model.Part{f.part}); err != nil {
		t.Fatal(err)
	}
	if b.spaceCalls.Load() != 1 {
		t.Fatalf("small archive performed %d free-space queries", b.spaceCalls.Load())
	}
}

func TestCompleteOrphansRequireSourceIdentityAcrossRestarts(t *testing.T) {
	data := zipBytes(t, 40_000)
	for _, filename := range []string{"archive.zip.part", "archive.zip"} {
		t.Run(filename, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { serveArchive(w, r, data, `"v1"`) }))
			defer srv.Close()
			f := newFixture(t, data, "archive.zip", srv.URL, nil)
			orphan := bytes.Clone(data)
			orphan[100] ^= 1
			if err := os.WriteFile(filepath.Join(f.dir, filename), orphan, 0600); err != nil {
				t.Fatal(err)
			}
			for i := 0; i < 2; i++ {
				err := f.engine.Download(context.Background(), []model.Part{f.part})
				if !errors.Is(err, ErrIdentityChanged) {
					t.Fatalf("orphan adopted on attempt %d: %v", i, err)
				}
				if r, ok := f.store.Get(f.part.Key()); ok && (r.Offset == f.part.Size || r.Complete) {
					t.Fatalf("unproved orphan became trusted: %+v", r)
				}
			}
			if got, _ := os.ReadFile(filepath.Join(f.dir, filename)); !bytes.Equal(got, orphan) {
				t.Fatal("orphan modified")
			}
		})
	}
}

func TestCompleteMatchingOrphansReconciled(t *testing.T) {
	data := zipBytes(t, 40_000)
	for _, filename := range []string{"archive.zip.part", "archive.zip"} {
		t.Run(filename, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { serveArchive(w, r, data, `"v1"`) }))
			defer srv.Close()
			f := newFixture(t, data, "archive.zip", srv.URL, nil)
			if err := os.WriteFile(filepath.Join(f.dir, filename), data, 0600); err != nil {
				t.Fatal(err)
			}
			if err := f.engine.Download(context.Background(), []model.Part{f.part}); err != nil {
				t.Fatal(err)
			}
			assertCompleted(t, f, data)
		})
	}
}

func TestChunkedExcessStaysUntrustedOnNextRun(t *testing.T) {
	data := zipBytes(t, 30_000)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.(http.Flusher).Flush()
		_, _ = w.Write(data)
		_, _ = w.Write([]byte("excess"))
	}))
	defer srv.Close()
	f := newFixture(t, data, "archive.zip", srv.URL, nil)
	for i := 0; i < 2; i++ {
		if err := f.engine.Download(context.Background(), []model.Part{f.part}); !errors.Is(err, ErrProtocol) {
			t.Fatalf("oversized run %d: %v", i, err)
		}
		if _, err := os.Stat(filepath.Join(f.dir, f.part.Filename)); !errors.Is(err, os.ErrNotExist) {
			t.Fatal("oversized response published on restart")
		}
		r, _ := f.store.Get(f.part.Key())
		if r.Offset == f.part.Size || r.Complete {
			t.Fatalf("unclean EOF acquired trusted checkpoint: %+v", r)
		}
	}
}

func TestRange416OnlyReconcilesTrustedExactSize(t *testing.T) {
	data := zipBytes(t, 40_000)
	for _, mode := range []string{"complete", "short", "wrong-total", "malformed", "unproved"} {
		t.Run(mode, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				total := len(data)
				if mode == "wrong-total" {
					total++
				}
				value := fmt.Sprintf("bytes */%d", total)
				if mode == "malformed" {
					value = "bytes */*"
				}
				w.Header().Set("Content-Range", value)
				w.WriteHeader(416)
			}))
			defer srv.Close()
			f := newFixture(t, data, "archive.zip", srv.URL, nil)
			count := len(data)
			if mode == "short" {
				count--
			}
			seedPartial(t, f, data, count, `"v1"`)
			r, _ := f.store.Get(f.part.Key())
			s := &session{record: r, name: f.part.Filename + ".part", proveComplete: mode == "unproved"}
			err := f.engine.transfer(context.Background(), s, srv.URL)
			if mode == "complete" {
				if err != nil {
					t.Fatal(err)
				}
				if err = f.engine.verifyArchive(context.Background(), s.name, f.part); err != nil {
					t.Fatal(err)
				}
			} else if !errors.Is(err, ErrProtocol) {
				t.Fatalf("invalid 416 accepted: %v", err)
			}
		})
	}
}

type failingTransport struct{}

func (failingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return nil, fmt.Errorf("transport failed for %s", r.URL)
}

func TestSignedURLNotInErrorsEventsOrState(t *testing.T) {
	data := zipBytes(t, 100)
	f := newFixture(t, data, "archive.zip", "https://example.invalid/download?secret=DO-NOT-LOG", func(o *Options) { o.Client = &http.Client{Transport: failingTransport{}} })
	stream, stop := f.bus.Subscribe(20)
	defer stop()
	err := f.engine.Download(context.Background(), []model.Part{f.part})
	if err == nil {
		t.Fatal("expected request failure")
	}
	if strings.Contains(err.Error(), "DO-NOT-LOG") {
		t.Fatal("signed URL leaked in error")
	}
	for {
		select {
		case ev := <-stream:
			if strings.Contains(ev.Message, "DO-NOT-LOG") {
				t.Fatal("signed URL leaked in event")
			}
		default:
			goto drained
		}
	}
drained:
	disk, _ := os.ReadFile(f.statePath)
	if bytes.Contains(disk, []byte("DO-NOT-LOG")) {
		t.Fatal("signed URL leaked in state")
	}
}

func TestPartByPartSerializesRequests(t *testing.T) {
	data := zipBytes(t, 100_000)
	var mu sync.Mutex
	active, maximum := 0, 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		active++
		maximum = max(maximum, active)
		mu.Unlock()
		time.Sleep(5 * time.Millisecond)
		serveArchive(w, r, data, `"v1"`)
		mu.Lock()
		active--
		mu.Unlock()
	}))
	defer srv.Close()
	f := newFixture(t, data, "one.zip", srv.URL, func(o *Options) { o.PartByPart = true; o.Concurrency = 4 })
	p2 := f.part
	p2.Index = 2
	p2.Filename = "two.zip"
	if err := f.engine.Download(context.Background(), []model.Part{f.part, p2}); err != nil {
		t.Fatal(err)
	}
	mu.Lock()
	defer mu.Unlock()
	if maximum != 1 {
		t.Fatalf("part-by-part concurrent requests: %d", maximum)
	}
}

func TestAuthRefreshFailurePreservesPartial(t *testing.T) {
	data := zipBytes(t, 40_000)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(403) }))
	defer srv.Close()
	f := newFixture(t, data, "archive.zip", srv.URL, func(o *Options) {
		o.Links = linkFunc(func(context.Context, string, int) (string, error) {
			return "", errors.New("signed URL credentials=PRIVATE")
		})
	})
	seedPartial(t, f, data, 10_000, `"v1"`)
	err := f.engine.Download(context.Background(), []model.Part{f.part})
	if err == nil || strings.Contains(err.Error(), "PRIVATE") {
		t.Fatalf("refresh error: %v", err)
	}
	got, _ := os.ReadFile(filepath.Join(f.dir, f.part.Filename+".part"))
	if !bytes.Equal(got, data[:10_000]) {
		t.Fatal("refresh failure changed partial")
	}
	r, _ := f.store.Get(f.part.Key())
	if r.Offset != 10_000 || r.Complete {
		t.Fatalf("refresh failure checkpoint: %+v", r)
	}
}
