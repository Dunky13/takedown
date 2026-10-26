package downloader

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/Dunky13/takedown/internal/model"
	"github.com/Dunky13/takedown/internal/storage"
)

type failingSyncWriter struct {
	io.WriteCloser
	failure error
}

func (w failingSyncWriter) Sync() error { return w.failure }

type failingSyncBackend struct {
	*storage.Local
	failure error
}

func (b failingSyncBackend) OpenWriteResume(ctx context.Context, name string, offset int64) (io.WriteCloser, error) {
	w, err := b.Local.OpenWriteResume(ctx, name, offset)
	if err != nil {
		return nil, err
	}
	return failingSyncWriter{w, b.failure}, nil
}

func TestSyncFailureNeverPublishesCheckpointOrFinalArchive(t *testing.T) {
	data := zipBytes(t, 1000)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { serveArchive(w, r, data, `"v1"`) }))
	defer server.Close()
	f := newFixture(t, data, "archive.zip", server.URL, nil)
	failure := errors.New("injected archive sync failure")
	e, err := New(Options{Backend: failingSyncBackend{f.backend, failure}, Destination: "fixture-destination", State: f.store, Client: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	if err = e.Download(context.Background(), []model.Part{f.part}); !errors.Is(err, failure) {
		t.Fatalf("sync failure not propagated: %v", err)
	}
	r, _ := f.store.Get(f.part.Key())
	if r.Complete || r.Offset != 0 {
		t.Fatalf("state advanced before successful sync: %+v", r)
	}
	if _, err = os.Stat(filepath.Join(f.dir, f.part.Filename)); !os.IsNotExist(err) {
		t.Fatal("unsynced archive was published")
	}
}
