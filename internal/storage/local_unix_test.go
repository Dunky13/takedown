//go:build linux || darwin

package storage

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/sys/unix"
)

func TestLocalRejectsFIFOAndWritableHardlinks(t *testing.T) {
	b, root := localForTest(t)
	if err := unix.Mkfifo(filepath.Join(root, "pipe"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := b.OpenRead(context.Background(), "pipe"); !errors.Is(err, ErrUnsafePath) {
		t.Fatalf("FIFO read: %v", err)
	}
	if _, err := b.OpenWriteResume(context.Background(), "pipe", 0); err == nil {
		t.Fatal("FIFO write accepted")
	}
	outside := filepath.Join(t.TempDir(), "outside.zip")
	if err := os.WriteFile(outside, []byte("safe"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Link(outside, filepath.Join(root, "hard.zip")); err != nil {
		t.Fatal(err)
	}
	if _, err := b.OpenWriteResume(context.Background(), "hard.zip", 0); !errors.Is(err, ErrUnsafePath) {
		t.Fatalf("hardlink write: %v", err)
	}
	if data, err := os.ReadFile(outside); err != nil || string(data) != "safe" {
		t.Fatalf("hardlink changed outside file: %q %v", data, err)
	}
}

func TestLocalRootIsPinnedAfterRename(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "destination")
	b, err := NewLocal(root)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()
	moved := filepath.Join(parent, "moved")
	if err := os.Rename(root, moved); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(root, 0700); err != nil {
		t.Fatal(err)
	}
	writeArchive(t, b, "archive.zip", 0, "pinned")
	if data, err := os.ReadFile(filepath.Join(moved, "archive.zip")); err != nil || string(data) != "pinned" {
		t.Fatalf("unpinned write: %q %v", data, err)
	}
	if _, err := os.Stat(filepath.Join(root, "archive.zip")); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("write followed replaced root path")
	}
}
