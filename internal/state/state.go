// Package state persists download checkpoints without download credentials.
package state

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/Dunky13/takedown/internal/model"
)

const schemaVersion = 4
const maxStateBytes = 32 << 20

var (
	ErrLocked = errors.New("state is locked by another process")
	ErrClosed = errors.New("state store is closed")
)

type Record struct {
	Part         model.Part `json:"part"`
	Destination  string     `json:"destination"`
	Offset       int64      `json:"offset"`
	Complete     bool       `json:"complete"`
	ETag         string     `json:"etag,omitempty"`
	LastModified string     `json:"last_modified,omitempty"`
	SHA256       string     `json:"sha256,omitempty"`
}

type document struct {
	Version int      `json:"version"`
	Records []Record `json:"records"`
}

type Store struct {
	mu      sync.Mutex
	path    string
	lock    *os.File
	records map[string]Record
	closed  bool
	// Kept on the store so failures can be exercised without changing global I/O.
	persist func([]byte) error
}

// Open holds a nonblocking process lock until Close. A malformed or unsupported
// document is an error; it is never silently reset.
func Open(path string) (*Store, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("state path is required")
	}
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve state path: %w", err)
	}
	if err = os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, fmt.Errorf("create state directory: %w", err)
	}
	lock, err := lockProcess(path + ".lock")
	if err != nil {
		return nil, err
	}
	s := &Store{path: path, lock: lock, records: make(map[string]Record)}
	s.persist = func(data []byte) error { return atomicWrite(path, data) }
	if err = s.load(); err != nil {
		_ = unlockProcess(lock)
		_ = lock.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) load() error {
	info, err := os.Lstat(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("inspect state: %w", err)
	}
	if !info.Mode().IsRegular() || info.Size() > maxStateBytes {
		return errors.New("state must be a regular file within the size limit")
	}
	f, err := os.Open(s.path)
	if err != nil {
		return fmt.Errorf("open state: %w", err)
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, maxStateBytes+1))
	if err != nil || len(data) > maxStateBytes {
		return errors.New("cannot read state within the size limit")
	}
	if err = rejectDuplicateFields(data); err != nil {
		return fmt.Errorf("invalid state JSON: %w", err)
	}
	// Pointer fields distinguish missing/null fields from legitimate zero values.
	var disk struct {
		Version *int `json:"version"`
		Records *[]struct {
			Part         *model.Part `json:"part"`
			Destination  *string     `json:"destination"`
			Offset       *int64      `json:"offset"`
			Complete     *bool       `json:"complete"`
			ETag         string      `json:"etag,omitempty"`
			LastModified string      `json:"last_modified,omitempty"`
			SHA256       string      `json:"sha256,omitempty"`
		} `json:"records"`
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err = dec.Decode(&disk); err != nil {
		return errors.New("invalid state schema")
	}
	if disk.Version == nil || (*disk.Version != schemaVersion && *disk.Version != 3 && *disk.Version != 2) || disk.Records == nil {
		return errors.New("unsupported or incomplete state schema")
	}
	for _, d := range *disk.Records {
		if d.Part == nil || d.Destination == nil || d.Offset == nil || d.Complete == nil {
			return errors.New("incomplete state record")
		}
		r := Record{Part: *d.Part, Destination: *d.Destination, Offset: *d.Offset, Complete: *d.Complete, ETag: d.ETag, LastModified: d.LastModified, SHA256: d.SHA256}
		if err = validate(r); err != nil {
			return fmt.Errorf("invalid state record: %w", err)
		}
		if _, exists := s.records[r.Part.Key()]; exists {
			return errors.New("duplicate part in state")
		}
		s.records[r.Part.Key()] = r
	}
	return nil
}

func validate(r Record) error {
	if r.SHA256 != "" {
		if digest, err := hex.DecodeString(r.SHA256); err != nil || len(digest) != 32 || r.Offset != r.Part.Size {
			return errors.New("completion digest must be a valid SHA-256 bound to the full archive length")
		}
	}
	if err := r.Part.Validate(); err != nil {
		return err
	}
	if strings.TrimSpace(r.Destination) == "" || strings.ContainsRune(r.Destination, 0) {
		return errors.New("destination is required")
	}
	if u, err := url.Parse(r.Destination); err == nil && u.User != nil {
		if _, hasPassword := u.User.Password(); hasPassword {
			return errors.New("destination must not contain a password")
		}
	}
	if r.Offset < 0 || r.Offset > r.Part.Size || (r.Complete && r.Offset != r.Part.Size) {
		return errors.New("invalid checkpoint offset")
	}
	if strings.ContainsAny(r.ETag+r.LastModified, "\r\n\x00") {
		return errors.New("invalid HTTP validator")
	}
	return nil
}

func snapshot(records map[string]Record) []Record {
	out := make([]Record, 0, len(records))
	for _, r := range records {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Part.Key() < out[j].Part.Key() })
	return out
}

func (s *Store) Snapshot() []Record {
	s.mu.Lock()
	defer s.mu.Unlock()
	return snapshot(s.records)
}

func (s *Store) Get(key string) (Record, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.records[key]
	return r, ok
}

// Put atomically replaces a checkpoint. Link fields are erased in both the
// persisted document and the in-memory snapshot. Failed writes keep the prior map.
func (s *Store) Put(r Record) error {
	r.Part.URL = ""
	if err := validate(r); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return ErrClosed
	}
	next := make(map[string]Record, len(s.records)+1)
	for key, existing := range s.records {
		next[key] = existing
	}
	next[r.Part.Key()] = r
	return s.commit(next)
}

func (s *Store) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return ErrClosed
	}
	next := make(map[string]Record, len(s.records))
	for k, r := range s.records {
		if k != key {
			next[k] = r
		}
	}
	return s.commit(next)
}

func (s *Store) commit(next map[string]Record) error {
	data, err := json.MarshalIndent(document{Version: schemaVersion, Records: snapshot(next)}, "", "  ")
	if err != nil {
		return err
	}
	if len(data)+1 > maxStateBytes {
		return errors.New("state exceeds size limit")
	}
	if err = s.persist(append(data, '\n')); err != nil {
		return fmt.Errorf("persist state: %w", err)
	}
	s.records = next
	return nil
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil
	}
	s.closed = true
	return errors.Join(unlockProcess(s.lock), s.lock.Close())
}

func atomicWrite(path string, data []byte) error {
	dir := filepath.Dir(path)
	f, err := os.CreateTemp(dir, ".state-*")
	if err != nil {
		return err
	}
	name := f.Name()
	defer os.Remove(name)
	if _, err = f.Write(data); err != nil {
		_ = f.Close()
		return err
	}
	if err = f.Sync(); err != nil {
		_ = f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	// Fail before replacement when directory syncing is unavailable for this
	// filesystem. Known unsupported platform operations are handled by syncDir.
	if err = syncDir(dir); err != nil {
		return err
	}
	old, readErr := os.ReadFile(path)
	if readErr != nil && !errors.Is(readErr, os.ErrNotExist) {
		return readErr
	}
	if err = replaceFile(name, path); err != nil {
		return err
	}
	if err = syncDir(dir); err != nil {
		// A failure after rename must not make a rejected mutation visible. Keep
		// the old document available for an atomic rollback.
		rollbackErr := restore(path, old, errors.Is(readErr, os.ErrNotExist))
		return errors.Join(err, rollbackErr)
	}
	return nil
}

func restore(path string, data []byte, absent bool) error {
	if absent {
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("rollback state: %w", err)
		}
		return syncDir(filepath.Dir(path))
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".state-rollback-*")
	if err != nil {
		return fmt.Errorf("rollback state: %w", err)
	}
	defer os.Remove(f.Name())
	_, writeErr := f.Write(data)
	err = errors.Join(writeErr, f.Sync(), f.Close())
	if err == nil {
		err = replaceFile(f.Name(), path)
	}
	if err == nil {
		err = syncDir(filepath.Dir(path))
	}
	if err != nil {
		return fmt.Errorf("rollback state: %w", err)
	}
	return nil
}

func rejectDuplicateFields(data []byte) error {
	d := json.NewDecoder(bytes.NewReader(data))
	var walk func(int) error
	walk = func(depth int) error {
		if depth > 32 {
			return errors.New("excessive JSON nesting")
		}
		t, err := d.Token()
		if err != nil {
			return err
		}
		delim, ok := t.(json.Delim)
		if !ok {
			return nil
		}
		keys := make(map[string]bool)
		for d.More() {
			if delim == '{' {
				key, err := d.Token()
				if err != nil {
					return err
				}
				k, ok := key.(string)
				if !ok || keys[k] {
					return errors.New("duplicate or invalid object field")
				}
				keys[k] = true
			}
			if err = walk(depth + 1); err != nil {
				return err
			}
		}
		_, err = d.Token()
		return err
	}
	if err := walk(0); err != nil {
		return err
	}
	if _, err := d.Token(); err != io.EOF {
		return errors.New("trailing JSON data")
	}
	return nil
}
