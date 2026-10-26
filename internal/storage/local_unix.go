//go:build linux || darwin

package storage

import (
	"errors"
	"os"
	"strings"

	"golang.org/x/sys/unix"
)

func openLocalRoot(root string) (*os.File, error) {
	fd, err := unix.Open(root, unix.O_RDONLY|unix.O_DIRECTORY|unix.O_CLOEXEC|unix.O_NOFOLLOW, 0)
	if err != nil {
		return nil, err
	}
	return os.NewFile(uintptr(fd), root), nil
}

func localParent(root *os.File, name string, create bool) (*os.File, string, error) {
	parts := strings.Split(name, "/")
	fd, err := unix.Openat(int(root.Fd()), ".", unix.O_RDONLY|unix.O_DIRECTORY|unix.O_CLOEXEC|unix.O_NOFOLLOW, 0)
	if err != nil {
		return nil, "", err
	}
	for _, component := range parts[:len(parts)-1] {
		next, e := unix.Openat(fd, component, unix.O_RDONLY|unix.O_DIRECTORY|unix.O_CLOEXEC|unix.O_NOFOLLOW, 0)
		if errors.Is(e, unix.ENOENT) && create {
			if e = unix.Mkdirat(fd, component, 0700); e == nil || errors.Is(e, unix.EEXIST) {
				next, e = unix.Openat(fd, component, unix.O_RDONLY|unix.O_DIRECTORY|unix.O_CLOEXEC|unix.O_NOFOLLOW, 0)
			}
		}
		_ = unix.Close(fd)
		if e != nil {
			return nil, "", &os.PathError{Op: "open directory", Path: name, Err: e}
		}
		fd = next
	}
	return os.NewFile(uintptr(fd), "destination directory"), parts[len(parts)-1], nil
}

func openLocalFile(root *os.File, name string, flags int, createParents bool) (*os.File, error) {
	parent, base, err := localParent(root, name, createParents)
	if err != nil {
		return nil, err
	}
	defer parent.Close()
	fd, err := unix.Openat(int(parent.Fd()), base, flags|unix.O_CLOEXEC|unix.O_NOFOLLOW|unix.O_NONBLOCK, 0600)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: name, Err: err}
	}
	f := os.NewFile(uintptr(fd), name)
	var info unix.Stat_t
	if err := unix.Fstat(fd, &info); err != nil {
		_ = f.Close()
		return nil, err
	}
	if info.Mode&unix.S_IFMT != unix.S_IFREG || (flags&os.O_WRONLY != 0 && info.Nlink != 1) {
		_ = f.Close()
		return nil, ErrUnsafePath
	}
	return f, nil
}

func renameLocalFile(root *os.File, from, to string) error {
	oldParent, oldBase, err := localParent(root, from, false)
	if err != nil {
		return err
	}
	defer oldParent.Close()
	newParent, newBase, err := localParent(root, to, false)
	if err != nil {
		return err
	}
	defer newParent.Close()
	for _, target := range []struct {
		parent   *os.File
		name     string
		optional bool
	}{{oldParent, oldBase, false}, {newParent, newBase, true}} {
		var info unix.Stat_t
		err := unix.Fstatat(int(target.parent.Fd()), target.name, &info, unix.AT_SYMLINK_NOFOLLOW)
		if target.optional && errors.Is(err, unix.ENOENT) {
			continue
		}
		if err != nil {
			return err
		}
		if info.Mode&unix.S_IFMT != unix.S_IFREG {
			return ErrUnsafePath
		}
	}
	return unix.Renameat(int(oldParent.Fd()), oldBase, int(newParent.Fd()), newBase)
}

func localFreeSpace(dir *os.File, _ string) (uint64, error) {
	var info unix.Statfs_t
	if err := unix.Fstatfs(int(dir.Fd()), &info); err != nil {
		return 0, err
	}
	if info.Bsize <= 0 {
		return 0, &SpaceUnknownError{errors.New("invalid filesystem block size")}
	}
	return spaceProduct(info.Bavail, uint64(info.Bsize))
}
