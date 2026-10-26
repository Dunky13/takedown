// Package storage provides confined, resumable archive destinations.
package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net/url"
	"os"
	"runtime"
	"strings"
	"unicode"
	"unicode/utf8"
)

type FileInfo = os.FileInfo

type StorageBackend interface {
	Stat(context.Context, string) (FileInfo, error)
	OpenWriteResume(context.Context, string, int64) (io.WriteCloser, error)
	Rename(context.Context, string, string) error
	FreeSpace(context.Context) (uint64, error)
	Close() error
}

type ReadBackend interface {
	OpenRead(context.Context, string) (ReadSeekCloser, error)
}

type ReadSeekCloser interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
}

type Options struct {
	Password, Domain, PrivateKeyPath, KnownHostsPath string
	AllowUnknownSpace                                bool
}

var (
	ErrSpaceUnknown  = errors.New("destination free space is unavailable")
	ErrUnsafePath    = errors.New("unsafe destination path")
	ErrInvalidOffset = errors.New("invalid resume offset")
)

// SpaceUnknownError preserves the provider failure and supports errors.Is.
type SpaceUnknownError struct{ Cause error }

func (e *SpaceUnknownError) Error() string        { return fmt.Sprintf("%v: %v", ErrSpaceUnknown, e.Cause) }
func (e *SpaceUnknownError) Unwrap() error        { return e.Cause }
func (e *SpaceUnknownError) Is(target error) bool { return target == ErrSpaceUnknown }

// Open connects to a destination. Credentials are accepted only through Options;
// URL passwords, queries, and fragments are rejected without echoing their values.
func Open(ctx context.Context, destination string, opts Options) (StorageBackend, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if destination == "" {
		return nil, errors.New("destination is required")
	}
	if !strings.Contains(destination, "://") {
		if runtime.GOOS != "windows" && (strings.HasPrefix(destination, `\\`) || (len(destination) > 1 && destination[1] == ':')) {
			return nil, errors.New("Windows drive and UNC destinations require Windows; use smb:// on other systems")
		}
		b, err := NewLocal(destination)
		if err != nil {
			return nil, err
		}
		b.allowUnknown = opts.AllowUnknownSpace
		return b, nil
	}
	u, err := url.Parse(destination)
	if err != nil {
		return nil, errors.New("invalid destination URL")
	}
	if u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || u.Opaque != "" {
		return nil, errors.New("destination URL cannot contain query, fragment, or opaque data")
	}
	if u.User != nil {
		if _, ok := u.User.Password(); ok {
			return nil, errors.New("supply destination password through options, never in the URL")
		}
	}
	if u.Hostname() == "" {
		return nil, errors.New("destination URL requires a host")
	}
	if u.User == nil || u.User.Username() == "" {
		return nil, errors.New("remote destination requires a username")
	}
	switch u.Scheme {
	case "smb":
		return openSMB(ctx, u, opts)
	case "sftp":
		return openSFTP(ctx, u, opts)
	default:
		return nil, errors.New("unsupported destination scheme (use local path, smb://, or sftp://)")
	}
}

// Names are portable relative paths. Reject alternate data streams, devices,
// ambiguous Windows spellings, control characters, and traversal on every OS.
func validateName(name string) error {
	if name == "" || !utf8.ValidString(name) || strings.HasPrefix(name, "/") || strings.ContainsAny(name, `\:<>"|?*`) {
		return ErrUnsafePath
	}
	for _, r := range name {
		if unicode.IsControl(r) {
			return ErrUnsafePath
		}
	}
	for _, component := range strings.Split(name, "/") {
		if component == "" || component == "." || component == ".." || strings.HasSuffix(component, ".") || strings.HasSuffix(component, " ") {
			return ErrUnsafePath
		}
		stem := strings.ToUpper(strings.TrimRight(strings.SplitN(component, ".", 2)[0], " "))
		if stem == "CON" || stem == "PRN" || stem == "AUX" || stem == "NUL" || stem == "CONIN$" || stem == "CONOUT$" {
			return ErrUnsafePath
		}
		if (strings.HasPrefix(stem, "COM") || strings.HasPrefix(stem, "LPT")) && len([]rune(stem)) == 4 && strings.ContainsRune("123456789¹²³", []rune(stem)[3]) {
			return ErrUnsafePath
		}
	}
	return nil
}

func spaceProduct(values ...uint64) (uint64, error) {
	result := uint64(1)
	for _, v := range values {
		if v != 0 && result > math.MaxUint64/v {
			return 0, &SpaceUnknownError{errors.New("free-space arithmetic overflow")}
		}
		result *= v
	}
	return result, nil
}

type archiveFile interface {
	ReadSeekCloser
	io.Writer
	Stat() (os.FileInfo, error)
	Truncate(int64) error
}

// prepareResume preserves exactly the trusted prefix, discards stale tails, and
// never creates a sparse gap or uses append mode.
func prepareResume(f archiveFile, offset int64) error {
	if offset < 0 {
		return ErrInvalidOffset
	}
	info, err := f.Stat()
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return ErrUnsafePath
	}
	if offset > info.Size() {
		return fmt.Errorf("%w: offset exceeds existing size", ErrInvalidOffset)
	}
	if err := f.Truncate(offset); err != nil {
		return err
	}
	_, err = f.Seek(offset, io.SeekStart)
	return err
}
