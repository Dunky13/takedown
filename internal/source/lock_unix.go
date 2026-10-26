//go:build unix

package source

import (
	"errors"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

func lockState(name string) (func(), error) {
	fd, err := unix.Open(name, unix.O_CREAT|unix.O_RDWR|unix.O_CLOEXEC|unix.O_NOFOLLOW, 0600)
	if err != nil {
		return nil, errors.New("cannot open private archive job lock")
	}
	f := os.NewFile(uintptr(fd), name)
	if info, err := f.Stat(); err != nil || !info.Mode().IsRegular() {
		f.Close()
		return nil, errors.New("archive job lock must be a regular file")
	}
	if err = unix.Flock(fd, unix.LOCK_EX|unix.LOCK_NB); err != nil {
		f.Close()
		return nil, errors.New("archive job state is in use by another process; retry when that operation finishes")
	}
	return func() { _ = unix.Flock(fd, unix.LOCK_UN); _ = f.Close() }, nil
}
func replaceState(old, new string) error {
	if err := os.Rename(old, new); err != nil {
		return errors.New("cannot atomically replace archive job state")
	}
	return syncStateDir(filepath.Dir(new))
}

func syncStateDir(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return errors.New("cannot open archive job directory for durability")
	}
	defer f.Close()
	if err = f.Sync(); err != nil {
		return errors.New("cannot sync archive job directory; creation must remain on hold until state durability is verified")
	}
	return nil
}
