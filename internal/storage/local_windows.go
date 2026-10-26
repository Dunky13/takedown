//go:build windows

package storage

import (
	"errors"
	"os"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

func openLocalRoot(root string) (*os.File, error) {
	p, err := windows.UTF16PtrFromString(root)
	if err != nil {
		return nil, err
	}
	h, err := windows.CreateFile(p, windows.FILE_GENERIC_READ, windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE, nil, windows.OPEN_EXISTING, windows.FILE_FLAG_BACKUP_SEMANTICS|windows.FILE_FLAG_OPEN_REPARSE_POINT, 0)
	if err != nil {
		return nil, err
	}
	if err := rejectReparse(h); err != nil {
		_ = windows.CloseHandle(h)
		return nil, err
	}
	return os.NewFile(uintptr(h), root), nil
}

func rejectReparse(handle windows.Handle) error {
	var info windows.ByHandleFileInformation
	if err := windows.GetFileInformationByHandle(handle, &info); err != nil {
		return err
	}
	if info.FileAttributes&windows.FILE_ATTRIBUTE_REPARSE_POINT != 0 {
		return ErrUnsafePath
	}
	return nil
}

// NtCreateFile supplies true handle-relative opens on Windows, including UNC.
func ntOpen(parent windows.Handle, name string, access, disposition, options uint32) (*os.File, error) {
	objectName, err := windows.NewNTUnicodeString(name)
	if err != nil {
		return nil, err
	}
	oa := &windows.OBJECT_ATTRIBUTES{RootDirectory: parent, ObjectName: objectName, Attributes: windows.OBJ_CASE_INSENSITIVE | windows.OBJ_DONT_REPARSE}
	oa.Length = uint32(unsafe.Sizeof(*oa))
	var iosb windows.IO_STATUS_BLOCK
	var handle windows.Handle
	err = windows.NtCreateFile(&handle, access|windows.SYNCHRONIZE, oa, &iosb, nil, windows.FILE_ATTRIBUTE_NORMAL,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE, disposition,
		options|windows.FILE_OPEN_REPARSE_POINT|windows.FILE_SYNCHRONOUS_IO_NONALERT, 0, 0)
	if err != nil {
		var status windows.NTStatus
		if errors.As(err, &status) {
			err = status.Errno()
		}
		return nil, &os.PathError{Op: "open", Path: name, Err: err}
	}
	if err := rejectReparse(handle); err != nil {
		_ = windows.CloseHandle(handle)
		return nil, err
	}
	return os.NewFile(uintptr(handle), name), nil
}

func localParent(root *os.File, name string, create bool) (*os.File, string, error) {
	parts := strings.Split(name, "/")
	var handle windows.Handle
	if err := windows.DuplicateHandle(windows.CurrentProcess(), windows.Handle(root.Fd()), windows.CurrentProcess(), &handle, 0, false, windows.DUPLICATE_SAME_ACCESS); err != nil {
		return nil, "", err
	}
	parent := os.NewFile(uintptr(handle), "destination directory")
	disposition := uint32(windows.FILE_OPEN)
	if create {
		disposition = windows.FILE_OPEN_IF
	}
	for _, component := range parts[:len(parts)-1] {
		access := uint32(windows.FILE_GENERIC_READ)
		if create {
			access |= windows.FILE_GENERIC_WRITE
		}
		next, err := ntOpen(windows.Handle(parent.Fd()), component, access, disposition, windows.FILE_DIRECTORY_FILE)
		_ = parent.Close()
		if err != nil {
			return nil, "", err
		}
		parent = next
	}
	return parent, parts[len(parts)-1], nil
}

func openLocalFile(root *os.File, name string, flags int, createParents bool) (*os.File, error) {
	parent, base, err := localParent(root, name, createParents)
	if err != nil {
		return nil, err
	}
	defer parent.Close()
	access := uint32(windows.FILE_GENERIC_READ)
	if flags&os.O_WRONLY != 0 {
		access = windows.FILE_GENERIC_WRITE | windows.FILE_READ_ATTRIBUTES
	}
	disposition := uint32(windows.FILE_OPEN)
	if flags&os.O_CREATE != 0 {
		disposition = windows.FILE_OPEN_IF
	}
	f, err := ntOpen(windows.Handle(parent.Fd()), base, access, disposition, windows.FILE_NON_DIRECTORY_FILE)
	if err != nil {
		return nil, err
	}
	var info windows.ByHandleFileInformation
	err = windows.GetFileInformationByHandle(windows.Handle(f.Fd()), &info)
	if err != nil || (flags&os.O_WRONLY != 0 && info.NumberOfLinks != 1) {
		_ = f.Close()
		if err != nil {
			return nil, err
		}
		return nil, ErrUnsafePath
	}
	return f, nil
}

type renameInformation struct {
	ReplaceIfExists uint32
	RootDirectory   windows.Handle
	FileNameLength  uint32
	FileName        [1]uint16
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
	file, err := ntOpen(windows.Handle(oldParent.Fd()), oldBase, windows.DELETE|windows.FILE_READ_ATTRIBUTES, windows.FILE_OPEN, windows.FILE_NON_DIRECTORY_FILE)
	if err != nil {
		return err
	}
	defer file.Close()
	if target, err := ntOpen(windows.Handle(newParent.Fd()), newBase, windows.FILE_READ_ATTRIBUTES, windows.FILE_OPEN, windows.FILE_NON_DIRECTORY_FILE); err == nil {
		_ = target.Close()
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	utf16, err := windows.UTF16FromString(newBase)
	if err != nil {
		return err
	}
	length := len(utf16) - 1
	var layout renameInformation
	offset := int(unsafe.Offsetof(layout.FileName))
	buffer := make([]byte, offset+length*2)
	info := (*renameInformation)(unsafe.Pointer(&buffer[0]))
	info.ReplaceIfExists = 1
	info.RootDirectory = windows.Handle(newParent.Fd())
	info.FileNameLength = uint32(length * 2)
	copy(unsafe.Slice((*uint16)(unsafe.Pointer(&buffer[offset])), length), utf16[:length])
	var iosb windows.IO_STATUS_BLOCK
	err = windows.NtSetInformationFile(windows.Handle(file.Fd()), &iosb, &buffer[0], uint32(len(buffer)), windows.FileRenameInformation)
	var status windows.NTStatus
	if errors.As(err, &status) {
		return status.Errno()
	}
	return err
}

func localFreeSpace(dir *os.File, _ string) (uint64, error) {
	buffer := make([]uint16, 512)
	for {
		// Zero requests a normalized DOS volume path (VOLUME_NAME_DOS).
		n, err := windows.GetFinalPathNameByHandle(windows.Handle(dir.Fd()), &buffer[0], uint32(len(buffer)), 0)
		if err != nil {
			return 0, err
		}
		if int(n) < len(buffer) {
			buffer = buffer[:n+1]
			break
		}
		buffer = make([]uint16, n+1)
	}
	var available uint64
	err := windows.GetDiskFreeSpaceEx(&buffer[0], &available, nil, nil)
	return available, err
}
