package storage

import (
	"context"
	"errors"
	"io"
	"net"
	"os"
	"path"
	"strings"
	"sync"
)

type remoteOps struct {
	lstat  func(context.Context, string) (os.FileInfo, error)
	mkdir  func(context.Context, string) error
	open   func(context.Context, string, int) (archiveFile, error)
	rename func(context.Context, string, string) error
	free   func(context.Context, string) (uint64, error)
	close  func() error
}

// Standard SFTP and the SMB library do not provide handle-relative no-follow
// opens. Every component is checked with Lstat. Remote destinations must therefore
// be controlled by the user and not concurrently mutated by an untrusted party.
type remoteBackend struct {
	root         string
	base         string
	ops          remoteOps
	conn         net.Conn
	mu           sync.RWMutex
	closed       bool
	files        fileSet
	once         sync.Once
	closeErr     error
	allowUnknown bool
}

func (b *remoteBackend) UnknownSpaceAllowed() bool { return b.allowUnknown }

// Canceling any SFTP operation invalidates its shared connection. This interrupts
// blocked library calls without leaking goroutines; callers reopen after cancel.
func remoteCall[T any](ctx context.Context, conn net.Conn, call func() (T, error)) (T, error) {
	var zero T
	if err := ctx.Err(); err != nil {
		return zero, err
	}
	stop := context.AfterFunc(ctx, func() { _ = conn.Close() })
	defer stop()
	value, err := call()
	if ctx.Err() != nil {
		return zero, ctx.Err()
	}
	return value, err
}

func remoteDo(ctx context.Context, conn net.Conn, call func() error) error {
	_, err := remoteCall(ctx, conn, func() (struct{}, error) { return struct{}{}, call() })
	return err
}

func (b *remoteBackend) ensureDirectories(ctx context.Context, directory string, create bool) error {
	current := b.base
	if info, err := b.ops.lstat(ctx, current); err != nil {
		return err
	} else if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return ErrUnsafePath
	}
	rel := directory
	if b.base == "/" {
		rel = strings.TrimPrefix(directory, "/")
	} else if directory == "." {
		rel = ""
	}
	for _, component := range strings.Split(strings.Trim(rel, "/"), "/") {
		if component == "" || component == "." {
			continue
		}
		current = path.Join(current, component)
		info, err := b.ops.lstat(ctx, current)
		if errors.Is(err, os.ErrNotExist) && create {
			if err = b.ops.mkdir(ctx, current); err != nil && !errors.Is(err, os.ErrExist) {
				return err
			}
			info, err = b.ops.lstat(ctx, current)
		}
		if err != nil {
			return err
		}
		if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return ErrUnsafePath
		}
	}
	return nil
}

func (b *remoteBackend) resolve(ctx context.Context, name string, createParents bool, mayNotExist bool) (string, error) {
	if err := validateName(name); err != nil {
		return "", err
	}
	full := path.Join(b.root, name)
	if err := b.ensureDirectories(ctx, path.Dir(full), createParents); err != nil {
		return "", err
	}
	info, err := b.ops.lstat(ctx, full)
	if errors.Is(err, os.ErrNotExist) && mayNotExist {
		return full, nil
	}
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() {
		return "", ErrUnsafePath
	}
	return full, nil
}

func (b *remoteBackend) Stat(ctx context.Context, name string) (FileInfo, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if b.closed {
		return nil, os.ErrClosed
	}
	full, err := b.resolve(ctx, name, false, false)
	if err != nil {
		return nil, err
	}
	return b.ops.lstat(ctx, full)
}

func (b *remoteBackend) OpenWriteResume(ctx context.Context, name string, offset int64) (io.WriteCloser, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if b.closed {
		return nil, os.ErrClosed
	}
	if offset < 0 {
		return nil, ErrInvalidOffset
	}
	full, err := b.resolve(ctx, name, offset == 0, offset == 0)
	if err != nil {
		return nil, err
	}
	flags := os.O_WRONLY
	if offset == 0 {
		flags |= os.O_CREATE
	}
	f, err := b.ops.open(ctx, full, flags)
	if err != nil {
		return nil, err
	}
	if err := remoteDo(ctx, b.conn, func() error { return prepareResume(f, offset) }); err != nil {
		_ = f.Close()
		return nil, err
	}
	return b.files.add(ctx, f, func() { _ = b.conn.Close() }, true)
}

func (b *remoteBackend) OpenRead(ctx context.Context, name string) (ReadSeekCloser, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if b.closed {
		return nil, os.ErrClosed
	}
	full, err := b.resolve(ctx, name, false, false)
	if err != nil {
		return nil, err
	}
	f, err := b.ops.open(ctx, full, os.O_RDONLY)
	if err != nil {
		return nil, err
	}
	return b.files.add(ctx, f, func() { _ = b.conn.Close() }, false)
}

func (b *remoteBackend) Rename(ctx context.Context, from, to string) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if err := ctx.Err(); err != nil {
		return err
	}
	if b.closed {
		return os.ErrClosed
	}
	oldPath, err := b.resolve(ctx, from, false, false)
	if err != nil {
		return err
	}
	newPath, err := b.resolve(ctx, to, false, true)
	if err != nil {
		return err
	}
	return b.ops.rename(ctx, oldPath, newPath)
}

func (b *remoteBackend) FreeSpace(ctx context.Context) (uint64, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	if b.closed {
		return 0, os.ErrClosed
	}
	if err := b.ensureDirectories(ctx, b.root, false); err != nil {
		return 0, err
	}
	n, err := b.ops.free(ctx, b.root)
	if err != nil && ctx.Err() != nil {
		return 0, ctx.Err()
	}
	if err != nil && !errors.Is(err, ErrSpaceUnknown) {
		return 0, &SpaceUnknownError{err}
	}
	return n, err
}

func (b *remoteBackend) Close() error {
	b.once.Do(func() {
		// Interrupt in-flight calls before waiting for their read locks. Graceful
		// protocol cleanup is attempted while bounded by the transport deadline.
		_ = b.conn.SetDeadline(timeNowPlusCloseTimeout())
		b.mu.Lock()
		b.closed = true
		b.mu.Unlock()
		b.closeErr = errors.Join(b.files.close(), b.ops.close(), closeResult(b.conn.Close()))
	})
	return b.closeErr
}

func closeResult(err error) error {
	if errors.Is(err, net.ErrClosed) || errors.Is(err, os.ErrClosed) || errors.Is(err, io.EOF) {
		return nil
	}
	return err
}

var _ StorageBackend = (*remoteBackend)(nil)
var _ ReadBackend = (*remoteBackend)(nil)
