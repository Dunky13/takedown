package state

import "os"

// ReadSnapshot reads one atomically published document without acquiring the
// writer lock, so status remains available while another process downloads.
func ReadSnapshot(path string) ([]Record, error) {
	if _, err := os.Lstat(path); err != nil {
		return nil, err
	}
	s := &Store{path: path, records: make(map[string]Record)}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s.Snapshot(), nil
}
