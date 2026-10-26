//go:build !linux && !darwin && !windows

package storage

import (
	"errors"
	"os"
)

var errUnsupportedLocal = errors.New("confined local storage is supported on Linux, macOS, and Windows")

func openLocalRoot(string) (*os.File, error)                      { return nil, errUnsupportedLocal }
func openLocalFile(*os.File, string, int, bool) (*os.File, error) { return nil, errUnsupportedLocal }
func renameLocalFile(*os.File, string, string) error              { return errUnsupportedLocal }
func localFreeSpace(*os.File, string) (uint64, error) {
	return 0, &SpaceUnknownError{errUnsupportedLocal}
}
