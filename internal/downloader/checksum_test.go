package downloader

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/Dunky13/takedown/internal/model"
)

func TestDriveChecksumBindsExistingArchiveWithoutAnotherDownload(t *testing.T) {
	data := zipBytes(t, 200)
	f := newFixture(t, data, "takeout.zip", "", nil)
	digest := md5.Sum(data)
	f.part.MD5 = hex.EncodeToString(digest[:])
	if err := os.WriteFile(filepath.Join(f.dir, f.part.Filename), data, 0600); err != nil {
		t.Fatal(err)
	}
	// No link source or URL is configured. A server-provided checksum can prove
	// an orphan's identity without fetching an entire existing archive again.
	if err := f.engine.Download(context.Background(), []model.Part{f.part}); err != nil {
		t.Fatal(err)
	}
	r, ok := f.store.Get(f.part.Key())
	if !ok || !r.Complete || r.Part.MD5 != f.part.MD5 {
		t.Fatal("verified checksum not preserved in completion")
	}
}

func TestValidArchiveWithWrongDriveChecksumIsNotPublished(t *testing.T) {
	data := zipBytes(t, 200)
	f := newFixture(t, data, "takeout.zip", "", nil)
	f.part.MD5 = "00000000000000000000000000000000"
	if err := os.WriteFile(filepath.Join(f.dir, f.part.Filename+".part"), data, 0600); err != nil {
		t.Fatal(err)
	}
	if err := f.engine.Download(context.Background(), []model.Part{f.part}); !errors.Is(err, ErrIntegrity) {
		t.Fatalf("wrong checksum: %v", err)
	}
	if _, err := os.Stat(filepath.Join(f.dir, f.part.Filename)); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("wrong archive was published")
	}
}
