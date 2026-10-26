package downloader

import (
	"archive/zip"
	"compress/gzip"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/Dunky13/takedown/internal/model"
)

type contextReader struct {
	ctx   context.Context
	r     io.Reader
	ready func(context.Context) error
}

func (r contextReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	if r.ready != nil {
		if err := r.ready(r.ctx); err != nil {
			return 0, err
		}
	}
	return r.r.Read(p)
}

func (e *Engine) verifyArchive(ctx context.Context, name string, part model.Part) error {
	_, err := e.verifyArchiveDigest(ctx, name, part, "")
	return err
}

func (e *Engine) verifyArchiveDigest(ctx context.Context, name string, part model.Part, expectedSHA256 string) (digest string, retErr error) {
	f, err := e.read.OpenRead(ctx, name)
	if err != nil {
		return "", err
	}
	defer func() { retErr = errors.Join(retErr, f.Close()) }()
	size, err := f.Seek(0, io.SeekEnd)
	if err != nil {
		return "", err
	}
	if size != part.Size {
		return "", fmt.Errorf("%w: archive length differs from its exact size", ErrIntegrity)
	}
	if _, err = f.Seek(0, io.SeekStart); err != nil {
		return "", err
	}
	sha := sha256.New()
	md := md5.New() // Drive supplies MD5 for blob integrity, not authentication.
	var sink io.Writer = sha
	if part.MD5 != "" {
		sink = io.MultiWriter(sha, md)
	}
	if _, err = io.Copy(sink, contextReader{ctx, f, e.waitReady}); err != nil {
		return "", err
	}
	digest = hex.EncodeToString(sha.Sum(nil))
	if expectedSHA256 != "" && !strings.EqualFold(digest, expectedSHA256) {
		return "", fmt.Errorf("%w: archive differs from its durable SHA-256 completion digest", ErrIdentityChanged)
	}
	if part.MD5 != "" && !strings.EqualFold(hex.EncodeToString(md.Sum(nil)), part.MD5) {
		return "", fmt.Errorf("%w: archive differs from its upstream checksum", ErrIntegrity)
	}
	if strings.HasSuffix(strings.ToLower(part.Filename), ".zip") {
		// archive/zip permits trailing junk; require an EOCD whose comment ends
		// exactly at EOF, including remote ReaderAt destinations.
		tail := make([]byte, min(size, int64(65535+22)))
		if _, err = f.ReadAt(tail, size-int64(len(tail))); err != nil {
			return "", fmt.Errorf("%w: ZIP end record is unreadable", ErrIntegrity)
		}
		found := false
		for i := len(tail) - 22; i >= 0; i-- {
			if binary.LittleEndian.Uint32(tail[i:i+4]) == 0x06054b50 && i+22+int(binary.LittleEndian.Uint16(tail[i+20:i+22])) == len(tail) {
				if binary.LittleEndian.Uint16(tail[i+4:i+6]) != 0 || binary.LittleEndian.Uint16(tail[i+6:i+8]) != 0 {
					return "", fmt.Errorf("%w: multi-disk ZIP is unsupported", ErrIntegrity)
				}
				found = true
				break
			}
		}
		if !found {
			return "", fmt.Errorf("%w: ZIP end record is missing or has trailing data", ErrIntegrity)
		}
		zr, err := zip.NewReader(f, size)
		if err != nil {
			return "", fmt.Errorf("%w: invalid ZIP directory", ErrIntegrity)
		}
		for _, entry := range zr.File {
			if err = e.waitReady(ctx); err != nil {
				return "", err
			}
			r, err := entry.Open()
			if err != nil {
				return "", fmt.Errorf("%w: unreadable ZIP entry", ErrIntegrity)
			}
			_, copyErr := io.Copy(io.Discard, contextReader{ctx, r, e.waitReady})
			err = errors.Join(copyErr, r.Close())
			if err != nil {
				if ctx.Err() != nil {
					return "", ctx.Err()
				}
				return "", fmt.Errorf("%w: ZIP entry data or CRC is invalid", ErrIntegrity)
			}
		}
		return digest, nil
	}
	if _, err = f.Seek(0, io.SeekStart); err != nil {
		return "", err
	}
	gr, err := gzip.NewReader(contextReader{ctx, f, e.waitReady})
	if err != nil {
		return "", fmt.Errorf("%w: invalid gzip header", ErrIntegrity)
	}
	// Reading through EOF validates CRC32, ISIZE and every concatenated member.
	_, copyErr := io.Copy(io.Discard, gr)
	err = errors.Join(copyErr, gr.Close())
	if ctx.Err() != nil {
		return "", ctx.Err()
	}
	if err != nil {
		return "", fmt.Errorf("%w: gzip data or trailer is invalid", ErrIntegrity)
	}
	return digest, nil
}
