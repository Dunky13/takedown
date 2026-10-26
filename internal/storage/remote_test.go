package storage

import (
	"context"
	"errors"
	"math"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pkg/sftp"
)

type smbSpaceInfo struct{ blocks, sector, sectors uint64 }

func (i smbSpaceInfo) BlockSize() uint64           { return i.sector }
func (i smbSpaceInfo) FragmentSize() uint64        { return i.sectors }
func (i smbSpaceInfo) TotalBlockCount() uint64     { return math.MaxInt64 }
func (i smbSpaceInfo) FreeBlockCount() uint64      { return math.MaxInt64 }
func (i smbSpaceInfo) AvailableBlockCount() uint64 { return i.blocks }

func TestRemoteSpaceUsesCallerAvailableBytesAndChecksOverflow(t *testing.T) {
	n, err := availableSMBSpace(smbSpaceInfo{blocks: 3, sector: 512, sectors: 8})
	if err != nil || n != 12288 {
		t.Fatalf("SMB space: %d %v", n, err)
	}
	n, err = availableSFTPSpace(&sftp.StatVFS{Bsize: 512, Frsize: 4096, Bavail: 3, Bfree: 9999})
	if err != nil || n != 12288 {
		t.Fatalf("SFTP available space: %d %v", n, err)
	}
	for _, info := range []smbSpaceInfo{{blocks: 1, sector: 0, sectors: 8}, {blocks: math.MaxInt64, sector: 4096, sectors: 8}, {blocks: math.MaxUint64, sector: 1, sectors: 1}} {
		if _, err := availableSMBSpace(info); !errors.Is(err, ErrSpaceUnknown) {
			t.Fatalf("invalid SMB space accepted: %+v %v", info, err)
		}
	}
	for _, info := range []*sftp.StatVFS{{Frsize: 0, Bavail: 1}, {Frsize: 4096, Bavail: math.MaxUint64}} {
		if _, err := availableSFTPSpace(info); !errors.Is(err, ErrSpaceUnknown) {
			t.Fatalf("invalid SFTP space accepted: %+v %v", info, err)
		}
	}
	if n, err := availableSMBSpace(smbSpaceInfo{blocks: 0, sector: 512, sectors: 8}); n != 0 || err != nil {
		t.Fatalf("full disk was not reported as zero: %d %v", n, err)
	}
}

func TestRemoteURLRootsAndSMBPaths(t *testing.T) {
	for _, tc := range []struct{ input, root, share string }{
		{"smb://user@host/share", ".", "share"},
		{"smb://user@host/share/path/to/export", "path/to/export", "share"},
		{"smb://user@[::1]:1445/share/.hidden/", ".hidden", "share"},
		{"sftp://user@host/", "/", ""},
		{"sftp://user@host/path%20with%20spaces/export", "/path with spaces/export", ""},
	} {
		u, err := url.Parse(tc.input)
		if err != nil {
			t.Fatal(err)
		}
		root, share, err := remotePath(u, u.Scheme == "smb")
		if err != nil || root != tc.root || share != tc.share {
			t.Fatalf("%s: root=%q share=%q err=%v", tc.input, root, share, err)
		}
	}
	if smbPath(".") != "" || smbPath("one/two.zip") != `one\two.zip` {
		t.Fatal("incorrect share-relative SMB path")
	}
}

func TestSharedRemoteOperationsConfineHiddenRootsAndResume(t *testing.T) {
	disk := t.TempDir()
	conn, peer := net.Pipe()
	t.Cleanup(func() { conn.Close(); peer.Close() })
	b := &remoteBackend{root: ".hidden/exports", base: ".", conn: conn}
	onDisk := func(name string) string { return filepath.Join(disk, filepath.FromSlash(name)) }
	b.ops = remoteOps{
		lstat: func(_ context.Context, name string) (os.FileInfo, error) { return os.Lstat(onDisk(name)) },
		mkdir: func(_ context.Context, name string) error { return os.Mkdir(onDisk(name), 0700) },
		open: func(_ context.Context, name string, flags int) (archiveFile, error) {
			if flags&(os.O_APPEND|os.O_TRUNC) != 0 {
				t.Error("unsafe remote open flags")
			}
			return os.OpenFile(onDisk(name), flags, 0600)
		},
		rename: func(_ context.Context, from, to string) error { return os.Rename(onDisk(from), onDisk(to)) },
		free:   func(context.Context, string) (uint64, error) { return 0, errors.New("provider lacks space query") },
		close:  func() error { return nil },
	}
	defer b.Close()
	if err := b.ensureDirectories(context.Background(), b.root, true); err != nil {
		t.Fatal(err)
	}
	writeArchive(t, b, "batch/archive.part", 0, "abcdef-tail")
	writeArchive(t, b, "batch/archive.part", 3, "done")
	if got := readArchive(t, b, "batch/archive.part"); got != "abcdone" {
		t.Fatalf("resume: %q", got)
	}
	if err := b.Rename(context.Background(), "batch/archive.part", "batch/archive.zip"); err != nil {
		t.Fatal(err)
	}
	if data, err := os.ReadFile(filepath.Join(disk, ".hidden", "exports", "batch", "archive.zip")); err != nil || string(data) != "abcdone" {
		t.Fatalf("wrong hidden root: %q %v", data, err)
	}
	if _, err := b.FreeSpace(context.Background()); !errors.Is(err, ErrSpaceUnknown) {
		t.Fatalf("free-space error not typed: %v", err)
	}
	if _, err := b.Stat(context.Background(), "../outside"); !errors.Is(err, ErrUnsafePath) {
		t.Fatalf("traversal: %v", err)
	}
	if err := os.Symlink(disk, filepath.Join(disk, ".hidden", "exports", "link")); err == nil {
		if _, err := b.OpenWriteResume(context.Background(), "link/outside.zip", 0); !errors.Is(err, ErrUnsafePath) {
			t.Fatalf("directory symlink accepted: %v", err)
		}
	}
}

func TestRemoteCallCancellationClosesBlockedTransport(t *testing.T) {
	conn, peer := net.Pipe()
	defer conn.Close()
	defer peer.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	started := make(chan struct{})
	result := make(chan error, 1)
	go func() {
		_, err := remoteCall(ctx, conn, func() (int, error) { close(started); return conn.Read(make([]byte, 1)) })
		result <- err
	}()
	<-started
	cancel()
	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("blocked call cancellation: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("remote cancellation left a blocked operation")
	}
}
