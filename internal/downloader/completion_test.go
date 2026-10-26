package downloader

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"

	"github.com/Dunky13/takedown/internal/model"
	"github.com/Dunky13/takedown/internal/state"
)

func sameSizeArchives(t *testing.T) ([]byte, []byte) {
	t.Helper()
	makeArchive := func(fill byte) []byte {
		var b bytes.Buffer
		z := zip.NewWriter(&b)
		w, err := z.CreateHeader(&zip.FileHeader{Name: "data", Method: zip.Store})
		if err != nil {
			t.Fatal(err)
		}
		if _, err = w.Write(bytes.Repeat([]byte{fill}, 1000)); err != nil {
			t.Fatal(err)
		}
		if err = z.Close(); err != nil {
			t.Fatal(err)
		}
		return b.Bytes()
	}
	a, b := makeArchive('a'), makeArchive('b')
	if len(a) != len(b) {
		t.Fatal("fixtures must have the same length")
	}
	return a, b
}

func TestCompletionDigestSurvivesRestartAndRejectsReplacement(t *testing.T) {
	for _, partial := range []bool{false, true} {
		name := "final"
		if partial {
			name = "partial"
		}
		t.Run(name, func(t *testing.T) {
			data, replacement := sameSizeArchives(t)
			var requests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests.Add(1)
				serveArchive(w, r, data, `"v1"`)
			}))
			defer server.Close()
			f := newFixture(t, data, "archive.zip", server.URL, nil)
			if err := f.engine.Download(context.Background(), []model.Part{f.part}); err != nil {
				t.Fatal(err)
			}
			sum := sha256.Sum256(data)
			r, _ := f.store.Get(f.part.Key())
			if r.SHA256 != hex.EncodeToString(sum[:]) {
				t.Fatalf("completion digest missing: %+v", r)
			}
			if err := f.store.Close(); err != nil {
				t.Fatal(err)
			}
			store, err := state.Open(f.statePath)
			if err != nil {
				t.Fatal(err)
			}
			defer store.Close()
			opts := f.engine.opts
			opts.State = store
			e, err := New(opts)
			if err != nil {
				t.Fatal(err)
			}
			if err = e.Download(context.Background(), []model.Part{f.part}); err != nil {
				t.Fatal(err)
			}
			if requests.Load() != 1 {
				t.Fatal("unchanged completion unexpectedly fetched again")
			}
			path := filepath.Join(f.dir, f.part.Filename)
			if partial {
				if err = os.Rename(path, path+".part"); err != nil {
					t.Fatal(err)
				}
				path += ".part"
			}
			if err = os.WriteFile(path, replacement, 0600); err != nil {
				t.Fatal(err)
			}
			if err = e.Download(context.Background(), []model.Part{f.part}); !errors.Is(err, ErrIdentityChanged) {
				t.Fatalf("changed valid same-size archive accepted: %v", err)
			}
			if requests.Load() != 1 {
				t.Fatal("bound completion should be checked against its durable digest")
			}
		})
	}
}

func TestLegacyFullCheckpointRequiresSourceProof(t *testing.T) {
	for _, partial := range []bool{false, true} {
		for _, changed := range []bool{false, true} {
			data, replacement := sameSizeArchives(t)
			var requests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requests.Add(1)
				serveArchive(w, r, data, `"v1"`)
			}))
			f := newFixture(t, data, "archive.zip", server.URL, nil)
			name := f.part.Filename
			if partial {
				name += ".part"
			}
			local := data
			if changed {
				local = replacement
			}
			if err := os.WriteFile(filepath.Join(f.dir, name), local, 0600); err != nil {
				t.Fatal(err)
			}
			if err := f.store.Put(state.Record{Part: f.part, Destination: f.engine.opts.Destination, Offset: f.part.Size, Complete: !partial, ETag: `"v1"`}); err != nil {
				t.Fatal(err)
			}
			err := f.engine.Download(context.Background(), []model.Part{f.part})
			server.Close()
			if changed {
				if !errors.Is(err, ErrIdentityChanged) {
					t.Fatalf("unbound legacy completion accepted: partial=%v: %v", partial, err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
			if requests.Load() != 1 {
				t.Fatalf("legacy completion must prove its bytes: requests=%d", requests.Load())
			}
		}
	}
}
