// Package model holds the shared archive metadata used by both interfaces.
package model

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

type Part struct {
	ExportID     string `json:"export_id"`
	Index        int    `json:"index"`
	Filename     string `json:"filename"`
	Size         int64  `json:"size"`
	MD5          string `json:"md5,omitempty"`
	RelativePath string `json:"relative_path,omitempty"`
	URL          string `json:"-"`
}

func (p Part) Key() string { return fmt.Sprintf("%s/%d", p.ExportID, p.Index) }

// NamespacePath separates exports and revisions without changing the original
// archive filename. The full digest avoids placing opaque provider IDs on disk.
func NamespacePath(exportID, filename string) string {
	digest := sha256.Sum256([]byte(exportID))
	return fmt.Sprintf("exports/%x/%s", digest, filename)
}

func (p Part) StoragePath() string {
	if p.RelativePath != "" {
		return p.RelativePath
	}
	return p.Filename
}

func (p Part) Validate() error {
	if p.RelativePath != "" && p.RelativePath != NamespacePath(p.ExportID, p.Filename) {
		return fmt.Errorf("archive path must match its immutable export namespace")
	}
	if p.ExportID == "" || p.Index < 1 || p.Size <= 0 {
		return fmt.Errorf("part requires export ID, positive index, and exact positive byte size")
	}
	if p.MD5 != "" {
		if digest, err := hex.DecodeString(p.MD5); err != nil || len(digest) != 16 {
			return fmt.Errorf("invalid archive MD5 checksum")
		}
	}
	if p.Filename == "" || p.Filename == "." || p.Filename == ".." || filepath.Base(p.Filename) != p.Filename || strings.ContainsAny(p.Filename, "\\/:\x00") {
		return fmt.Errorf("unsafe archive filename")
	}
	name := strings.ToLower(p.Filename)
	if !strings.HasSuffix(name, ".zip") && !strings.HasSuffix(name, ".tgz") && !strings.HasSuffix(name, ".tar.gz") {
		return fmt.Errorf("unsupported archive format: %s", p.Filename)
	}
	return nil
}

type Export struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Parts     []Part    `json:"parts"`
}
