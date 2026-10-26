package storage

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func localForTest(t *testing.T) (*Local, string) {
	t.Helper()
	root := t.TempDir()
	b, err := NewLocal(root)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := b.Close(); err != nil {
			t.Error(err)
		}
	})
	return b, root
}

func writeArchive(t *testing.T, backend StorageBackend, name string, offset int64, content string) {
	t.Helper()
	w, err := backend.OpenWriteResume(context.Background(), name, offset)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.WriteString(w, content); err != nil {
		_ = w.Close()
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
}

func readArchive(t *testing.T, backend ReadBackend, name string) string {
	t.Helper()
	r, err := backend.OpenRead(context.Background(), name)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestLocalResumeReadAndRename(t *testing.T) {
	b, root := localForTest(t)
	ctx := context.Background()
	writeArchive(t, b, "export/archive.zip.part", 0, "abcdef-stale-tail")
	writeArchive(t, b, "export/archive.zip.part", 3, "XYZ")
	if got := readArchive(t, b, "export/archive.zip.part"); got != "abcXYZ" {
		t.Fatalf("resume produced %q", got)
	}
	if err := b.Rename(ctx, "export/archive.zip.part", "export/archive.zip"); err != nil {
		t.Fatal(err)
	}
	if _, err := b.Stat(ctx, "export/archive.zip.part"); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("source still exists: %v", err)
	}
	info, err := b.Stat(ctx, "export/archive.zip")
	if err != nil || info.Size() != 6 {
		t.Fatalf("stat: %v, %v", info, err)
	}
	r, err := b.OpenRead(ctx, "export/archive.zip")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	buf := make([]byte, 3)
	if n, err := r.ReadAt(buf, 3); n != 3 || err != nil || string(buf) != "XYZ" {
		t.Fatalf("ReadAt: %d %q %v", n, buf, err)
	}
	if _, err := r.Seek(1, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	if _, err := io.ReadFull(r, buf); err != nil || string(buf) != "bcX" {
		t.Fatalf("Seek/Read: %q %v", buf, err)
	}
	writeArchive(t, b, "export/archive.zip", 0, "z")
	if got := readArchive(t, b, "export/archive.zip"); got != "z" {
		t.Fatalf("restart retained stale data: %q", got)
	}
	if info, err := os.Stat(filepath.Join(root, "export", "archive.zip")); err != nil || info.Size() != 1 {
		t.Fatalf("disk result: %v %v", info, err)
	}
	if _, err := b.FreeSpace(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestLocalRejectsInvalidOffsetWithoutChangingFile(t *testing.T) {
	b, root := localForTest(t)
	writeArchive(t, b, "archive.part", 0, "prefix")
	for _, offset := range []int64{-1, 7} {
		if w, err := b.OpenWriteResume(context.Background(), "archive.part", offset); !errors.Is(err, ErrInvalidOffset) {
			if w != nil {
				w.Close()
			}
			t.Fatalf("offset %d: %v", offset, err)
		}
	}
	if got := readArchive(t, b, "archive.part"); got != "prefix" {
		t.Fatalf("failed resume changed data: %q", got)
	}
	if w, err := b.OpenWriteResume(context.Background(), "missing.part", 1); !errors.Is(err, os.ErrNotExist) {
		if w != nil {
			w.Close()
		}
		t.Fatalf("missing resume: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "missing.part")); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("invalid resume created a file")
	}
}

func TestLocalRejectsUnsafeNames(t *testing.T) {
	b, _ := localForTest(t)
	ctx := context.Background()
	writeArchive(t, b, "source.part", 0, "safe")
	names := []string{"", ".", "..", "../outside", "/absolute", "a/../../x", "a//x", "a/./x", "a/../x", `a\x`, `C:\x`, "a:stream", "a\x00.zip", "a\n.zip", "a. ", "a.", "CON", "con.zip", "COM1.zip", "NUL .zip", "COM1 .zip", "LPT³.txt", "a/aux.txt", "*.zip", string([]byte{0xff}) + ".zip"}
	for _, name := range names {
		t.Run(strings.ReplaceAll(name, "/", "_"), func(t *testing.T) {
			if _, err := b.Stat(ctx, name); !errors.Is(err, ErrUnsafePath) {
				t.Errorf("Stat: %v", err)
			}
			if w, err := b.OpenWriteResume(ctx, name, 0); !errors.Is(err, ErrUnsafePath) {
				if w != nil {
					w.Close()
				}
				t.Errorf("OpenWriteResume: %v", err)
			}
			if r, err := b.OpenRead(ctx, name); !errors.Is(err, ErrUnsafePath) {
				if r != nil {
					r.Close()
				}
				t.Errorf("OpenRead: %v", err)
			}
			if err := b.Rename(ctx, "source.part", name); !errors.Is(err, ErrUnsafePath) {
				t.Errorf("Rename destination: %v", err)
			}
			if err := b.Rename(ctx, name, "target.zip"); !errors.Is(err, ErrUnsafePath) {
				t.Errorf("Rename source: %v", err)
			}
		})
	}
}

func TestLocalRejectsSymlinks(t *testing.T) {
	b, root := localForTest(t)
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "victim.zip"), []byte("untouched"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "escape")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if err := os.Symlink(filepath.Join(outside, "victim.zip"), filepath.Join(root, "link.zip")); err != nil {
		t.Fatal(err)
	}
	writeArchive(t, b, "source.part", 0, "inside")
	ctx := context.Background()
	for _, name := range []string{"escape/victim.zip", "escape/new.zip", "link.zip"} {
		if _, err := b.Stat(ctx, name); err == nil {
			t.Errorf("Stat followed %s", name)
		}
		if w, err := b.OpenWriteResume(ctx, name, 0); err == nil {
			w.Close()
			t.Errorf("write followed %s", name)
		}
		if r, err := b.OpenRead(ctx, name); err == nil {
			r.Close()
			t.Errorf("read followed %s", name)
		}
		if err := b.Rename(ctx, "source.part", name); err == nil {
			t.Errorf("rename destination followed %s", name)
		}
		if err := b.Rename(ctx, name, "final.zip"); err == nil {
			t.Errorf("rename source followed %s", name)
		}
	}
	data, err := os.ReadFile(filepath.Join(outside, "victim.zip"))
	if err != nil || string(data) != "untouched" {
		t.Fatalf("outside file changed: %q %v", data, err)
	}
	if _, err := os.Stat(filepath.Join(outside, "new.zip")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("outside file created: %v", err)
	}
}

func TestLocalCancellationAndCloseOwnHandles(t *testing.T) {
	b, _ := localForTest(t)
	ctx, cancel := context.WithCancel(context.Background())
	w, err := b.OpenWriteResume(ctx, "archive.part", 0)
	if err != nil {
		t.Fatal(err)
	}
	cancel()
	if _, err := w.Write([]byte("x")); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled write: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	writeArchive(t, b, "read.zip", 0, "data")
	r, err := b.OpenRead(context.Background(), "read.zip")
	if err != nil {
		t.Fatal(err)
	}
	if err := b.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Read(make([]byte, 1)); !errors.Is(err, os.ErrClosed) {
		t.Fatalf("backend left reader open: %v", err)
	}
	if _, err := b.Stat(context.Background(), "read.zip"); !errors.Is(err, os.ErrClosed) {
		t.Fatalf("closed Stat: %v", err)
	}
	if _, err := b.FreeSpace(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled free space: %v", err)
	}
}

func TestDestinationValidationDoesNotEchoSecrets(t *testing.T) {
	for _, destination := range []string{"sftp://user:do-not-print@host/a", "smb://user@host/share?password=do-not-print", "sftp://user:do-not-print@host:bad/a", "sftp://user@host/a#secret", "sftp://host/a", "smb://user@host/../share", "sftp://user@host/%2e%2e/outside", "sftp://user@host/a%2fb/../x"} {
		_, err := Open(context.Background(), destination, Options{})
		if err == nil {
			t.Fatalf("accepted %q", destination)
		}
		if strings.Contains(err.Error(), "do-not-print") {
			t.Fatalf("credential leaked: %v", err)
		}
	}
}
