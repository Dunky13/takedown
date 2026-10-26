package storage

import (
	"context"
	"errors"
	"os"
	"sync"
)

// fileSet owns every returned handle, including handles the caller forgot to
// close. Cancellation closes the transport first for remote files.
type fileSet struct {
	mu     sync.Mutex
	closed bool
	files  map[*contextFile]struct{}
}

func (s *fileSet) add(ctx context.Context, file archiveFile, abort func(), writing bool) (*contextFile, error) {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		_ = file.Close()
		return nil, os.ErrClosed
	}
	if s.files == nil {
		s.files = make(map[*contextFile]struct{})
	}
	f := &contextFile{ctx: ctx, file: file, owner: s, writing: writing, syncOnClose: writing && abort == nil}
	s.files[f] = struct{}{}
	s.mu.Unlock()
	f.mu.Lock()
	f.stop = context.AfterFunc(ctx, func() {
		if abort != nil {
			abort()
		}
		// Local writers stay open for their owner's final sync/checkpoint even
		// after cancellation. Remote cancellation must interrupt blocked I/O;
		// a failed remote sync then prevents publishing a newer checkpoint.
		if !writing || abort != nil {
			_ = f.Close()
		}
	})
	if f.closed {
		f.stop()
	}
	f.mu.Unlock()
	if err := ctx.Err(); err != nil {
		_ = f.Close()
		return nil, err
	}
	return f, nil
}

func (s *fileSet) close() error {
	s.mu.Lock()
	s.closed = true
	files := make([]*contextFile, 0, len(s.files))
	for f := range s.files {
		files = append(files, f)
	}
	s.mu.Unlock()
	var err error
	for _, f := range files {
		err = errors.Join(err, f.Close())
	}
	return err
}

type contextFile struct {
	ctx         context.Context
	file        archiveFile
	owner       *fileSet
	mu          sync.Mutex
	stop        func() bool
	closed      bool
	writing     bool
	syncOnClose bool
	synced      bool
	once        sync.Once
	closeErr    error
}

func (f *contextFile) check() error {
	if err := f.ctx.Err(); err != nil {
		return err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closed {
		return os.ErrClosed
	}
	return nil
}
func (f *contextFile) result(err error) error {
	if e := f.ctx.Err(); e != nil {
		return e
	}
	return err
}
func (f *contextFile) Read(p []byte) (int, error) {
	if err := f.check(); err != nil {
		return 0, err
	}
	n, err := f.file.Read(p)
	return n, f.result(err)
}
func (f *contextFile) ReadAt(p []byte, offset int64) (int, error) {
	if err := f.check(); err != nil {
		return 0, err
	}
	n, err := f.file.ReadAt(p, offset)
	return n, f.result(err)
}
func (f *contextFile) Write(p []byte) (int, error) {
	if err := f.check(); err != nil {
		return 0, err
	}
	n, err := f.file.Write(p)
	if n > 0 {
		f.mu.Lock()
		f.synced = false
		f.mu.Unlock()
	}
	return n, f.result(err)
}

// Sync is a cleanup operation: cancellation stops new I/O, but must still let a
// local writer make already-written bytes durable before publishing its offset.
func (f *contextFile) Sync() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closed {
		if f.writing && (f.synced || f.closeErr != nil) {
			return f.closeErr
		}
		return os.ErrClosed
	}
	return f.syncLocked()
}
func (f *contextFile) syncLocked() error {
	if f.synced {
		return nil
	}
	if syncer, ok := f.file.(interface{ Sync() error }); ok {
		err := syncer.Sync()
		if err == nil {
			f.synced = true
		}
		return err
	}
	return errors.New("destination writer cannot sync archive data")
}
func (f *contextFile) Seek(offset int64, whence int) (int64, error) {
	if err := f.check(); err != nil {
		return 0, err
	}
	n, err := f.file.Seek(offset, whence)
	return n, f.result(err)
}
func (f *contextFile) Close() error {
	f.once.Do(func() {
		f.mu.Lock()
		if f.stop != nil {
			f.stop()
		}
		if f.syncOnClose {
			f.closeErr = f.syncLocked()
		}
		f.closeErr = errors.Join(f.closeErr, f.file.Close())
		f.closed = true
		f.mu.Unlock()
		f.owner.mu.Lock()
		delete(f.owner.files, f)
		f.owner.mu.Unlock()
	})
	return f.closeErr
}
