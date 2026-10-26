package storage

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// Local also supports Windows UNC paths through the operating system. The root
// is pinned by a directory handle; archive operations do not follow symlinks.
type Local struct {
	root         string
	dir          *os.File
	mu           sync.RWMutex
	closed       bool
	files        fileSet
	once         sync.Once
	closeErr     error
	allowUnknown bool
}

func NewLocal(root string) (*Local, error) {
	if root == "" {
		return nil, errors.New("destination directory is required")
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(abs, 0700); err != nil {
		return nil, err
	}
	// An explicitly selected root may itself be a symlink. Canonicalize it once,
	// then pin it; child paths always reject symlinks and reparse points.
	abs, err = filepath.EvalSymlinks(abs)
	if err != nil {
		return nil, err
	}
	dir, err := openLocalRoot(abs)
	if err != nil {
		return nil, err
	}
	return &Local{root: abs, dir: dir}, nil
}

func (b *Local) UnknownSpaceAllowed() bool { return b.allowUnknown }

func (b *Local) Stat(ctx context.Context, name string) (FileInfo, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := validateName(name); err != nil {
		return nil, err
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.closed {
		return nil, os.ErrClosed
	}
	f, err := openLocalFile(b.dir, name, os.O_RDONLY, false)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.Stat()
}

func (b *Local) OpenWriteResume(ctx context.Context, name string, offset int64) (io.WriteCloser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := validateName(name); err != nil {
		return nil, err
	}
	if offset < 0 {
		return nil, ErrInvalidOffset
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.closed {
		return nil, os.ErrClosed
	}
	flags := os.O_WRONLY
	if offset == 0 {
		flags |= os.O_CREATE
	}
	f, err := openLocalFile(b.dir, name, flags, offset == 0)
	if err != nil {
		return nil, err
	}
	if err := prepareResume(f, offset); err != nil {
		_ = f.Close()
		return nil, err
	}
	return b.files.add(ctx, f, nil, true)
}

func (b *Local) OpenRead(ctx context.Context, name string) (ReadSeekCloser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := validateName(name); err != nil {
		return nil, err
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.closed {
		return nil, os.ErrClosed
	}
	f, err := openLocalFile(b.dir, name, os.O_RDONLY, false)
	if err != nil {
		return nil, err
	}
	return b.files.add(ctx, f, nil, false)
}

func (b *Local) Rename(ctx context.Context, from, to string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := validateName(from); err != nil {
		return err
	}
	if err := validateName(to); err != nil {
		return err
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.closed {
		return os.ErrClosed
	}
	return renameLocalFile(b.dir, from, to)
}

func (b *Local) FreeSpace(ctx context.Context) (uint64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.closed {
		return 0, os.ErrClosed
	}
	n, err := localFreeSpace(b.dir, b.root)
	if err != nil && !errors.Is(err, ErrSpaceUnknown) {
		err = &SpaceUnknownError{err}
	}
	return n, err
}

func (b *Local) Close() error {
	b.once.Do(func() {
		b.mu.Lock()
		b.closed = true
		b.mu.Unlock()
		b.closeErr = errors.Join(b.files.close(), b.dir.Close())
	})
	return b.closeErr
}

var _ StorageBackend = (*Local)(nil)
var _ ReadBackend = (*Local)(nil)
