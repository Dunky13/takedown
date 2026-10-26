//go:build !windows

package state

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

func lockProcess(path string) (*os.File, error) {
	fd, err := unix.Open(path, unix.O_CREAT|unix.O_RDWR|unix.O_CLOEXEC|unix.O_NOFOLLOW|unix.O_NONBLOCK, 0600)
	if err != nil {
		return nil, fmt.Errorf("open state lock: %w", err)
	}
	f := os.NewFile(uintptr(fd), path)
	info, err := f.Stat()
	if err != nil || !info.Mode().IsRegular() {
		_ = f.Close()
		return nil, errors.New("state lock must be a regular file")
	}
	if err = unix.Flock(fd, unix.LOCK_EX|unix.LOCK_NB); err != nil {
		_ = f.Close()
		if errors.Is(err, unix.EWOULDBLOCK) || errors.Is(err, unix.EAGAIN) {
			return nil, ErrLocked
		}
		return nil, fmt.Errorf("lock state: %w", err)
	}
	return f, nil
}

func unlockProcess(f *os.File) error {
	return unix.Flock(int(f.Fd()), unix.LOCK_UN)
}

func replaceFile(old, next string) error { return os.Rename(old, next) }

func syncDir(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	err = f.Sync()
	if errors.Is(err, unix.EINVAL) || errors.Is(err, unix.ENOTSUP) {
		return nil // Filesystems without directory fsync still get atomic rename.
	}
	return err
}
