package source

import (
	"errors"
	"os"

	"golang.org/x/sys/windows"
)

func lockState(name string) (func(), error) {
	if info, err := os.Lstat(name); err == nil && !info.Mode().IsRegular() {
		return nil, errors.New("archive job lock must be a regular file")
	}
	f, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, errors.New("cannot open private archive job lock")
	}
	var overlapped windows.Overlapped
	if err = windows.LockFileEx(windows.Handle(f.Fd()), windows.LOCKFILE_EXCLUSIVE_LOCK|windows.LOCKFILE_FAIL_IMMEDIATELY, 0, 1, 0, &overlapped); err != nil {
		f.Close()
		return nil, errors.New("archive job state is in use by another process; retry when that operation finishes")
	}
	return func() { _ = windows.UnlockFileEx(windows.Handle(f.Fd()), 0, 1, 0, &overlapped); _ = f.Close() }, nil
}

func replaceState(old, new string) error {
	from, err := windows.UTF16PtrFromString(old)
	if err != nil {
		return errors.New("invalid archive job state path")
	}
	to, err := windows.UTF16PtrFromString(new)
	if err != nil {
		return errors.New("invalid archive job state path")
	}
	if err = windows.MoveFileEx(from, to, windows.MOVEFILE_REPLACE_EXISTING|windows.MOVEFILE_WRITE_THROUGH); err != nil {
		return errors.New("cannot durably replace archive job state")
	}
	return nil
}
