package state

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/Dunky13/takedown/internal/model"
)

func record() Record {
	return Record{Part: model.Part{ExportID: "export", Index: 1, Filename: "takeout.zip", Size: 100, URL: "https://example.test/archive?secret=SIGNED"}, Destination: "test-destination", Offset: 20}
}

func TestVersionTwoMigratesAndCompletionDigestPersists(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	r := record()
	r.Part.URL = ""
	r.Offset, r.Complete = r.Part.Size, true
	legacy, err := json.Marshal(document{Version: 2, Records: []Record{r}})
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, legacy, 0600); err != nil {
		t.Fatal(err)
	}
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	r, ok := s.Get(r.Part.Key())
	if !ok || r.SHA256 != "" || !r.Complete {
		t.Fatalf("legacy record changed: %+v", r)
	}
	r.SHA256 = strings.Repeat("ab", 32)
	if err = s.Put(r); err != nil {
		t.Fatal(err)
	}
	if err = s.Close(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var doc document
	if err = json.Unmarshal(data, &doc); err != nil || doc.Version != 4 {
		t.Fatalf("schema did not migrate: version=%d, %v", doc.Version, err)
	}
	s, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	got, _ := s.Get(r.Part.Key())
	if got.SHA256 != r.SHA256 {
		t.Fatal("digest lost after reopen")
	}
	for _, invalid := range []Record{
		{Part: r.Part, Destination: r.Destination, Offset: r.Part.Size, SHA256: "invalid"},
		{Part: r.Part, Destination: r.Destination, Offset: 1, SHA256: r.SHA256},
	} {
		if err = s.Put(invalid); err == nil {
			t.Fatal("invalid digest binding accepted")
		}
	}
}

func TestStateRoundTripLockAndSecrets(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if _, err = Open(path); !errors.Is(err, ErrLocked) {
		t.Fatalf("second open: %v", err)
	}
	if err = s.Put(record()); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(data, []byte("SIGNED")) || bytes.Contains(data, []byte("https://")) {
		t.Fatalf("link leaked: %s", data)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	// Unix permission bits do not describe Windows ACLs.
	if runtime.GOOS != "windows" && info.Mode().Perm() != 0600 {
		t.Fatalf("state permissions: %v", info.Mode())
	}
	if err = s.Close(); err != nil {
		t.Fatal(err)
	}
	s, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	r, ok := s.Get("export/1")
	if !ok || r.Offset != 20 || r.Part.URL != "" {
		t.Fatalf("checkpoint: %+v", r)
	}
	copy := s.Snapshot()
	copy[0].Offset = 99
	r, _ = s.Get("export/1")
	if r.Offset != 20 {
		t.Fatal("snapshot aliases store")
	}
}

func TestStateProcessLock(t *testing.T) {
	if path := os.Getenv("TAKEDOWN_STATE_LOCK_TEST"); path != "" {
		_, err := Open(path)
		if !errors.Is(err, ErrLocked) {
			os.Exit(2)
		}
		return
	}
	path := filepath.Join(t.TempDir(), "state.json")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	cmd := exec.Command(os.Args[0], "-test.run=^TestStateProcessLock$")
	cmd.Env = append(os.Environ(), "TAKEDOWN_STATE_LOCK_TEST="+path)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("subprocess: %v: %s", err, out)
	}
}

func TestMutationFailurePreservesPriorState(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	r := record()
	if err = s.Put(r); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	s.persist = func([]byte) error { return errors.New("injected disk failure") }
	r.Offset = 90
	if err = s.Put(r); err == nil {
		t.Fatal("expected write failure")
	}
	got, _ := s.Get(r.Part.Key())
	if got.Offset != 20 {
		t.Fatalf("mutation survived failure: %+v", got)
	}
	after, _ := os.ReadFile(path)
	if !bytes.Equal(before, after) {
		t.Fatal("previous file changed")
	}
	if err = s.Delete(r.Part.Key()); err == nil {
		t.Fatal("expected delete failure")
	}
	if _, ok := s.Get(r.Part.Key()); !ok {
		t.Fatal("failed delete changed state")
	}
}

func TestMalformedSchemaFailsClosed(t *testing.T) {
	validPart := `"part":{"export_id":"e","index":1,"filename":"a.zip","size":3},"destination":"d","offset":0,"complete":false`
	for _, data := range []string{
		``, `null`, `{}`, `{"version":1,"records":[]}`, `{"version":2,"records":null}`,
		`{"version":2,"records":[],"unexpected":true}`,
		`{"version":2,"version":2,"records":[]}`,
		`{"version":2,"records":[]} {}`,
		`{"version":2,"records":[{"part":{"export_id":"e","index":1,"filename":"a.zip","size":3},"destination":"d"}]}`,
		`{"version":2,"records":[{` + validPart + `,"offset":1}]}`,
		`{"version":2,"records":[{` + validPart + `,"download_url":"secret"}]}`,
		`{"version":2,"records":[{` + validPart + `},{` + validPart + `}]}`,
	} {
		t.Run(data, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "state.json")
			if err := os.WriteFile(path, []byte(data), 0600); err != nil {
				t.Fatal(err)
			}
			if s, err := Open(path); err == nil {
				s.Close()
				t.Fatal("malformed state accepted")
			}
			after, _ := os.ReadFile(path)
			if string(after) != data {
				t.Fatal("malformed state was replaced")
			}
		})
	}
}

func TestConcurrentCheckpoints(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "state.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	var wg sync.WaitGroup
	for i := 1; i <= 12; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			r := record()
			r.Part.Index = index
			if err := s.Put(r); err != nil {
				t.Error(err)
			}
			_ = s.Snapshot()
		}(i)
	}
	wg.Wait()
	if got := len(s.Snapshot()); got != 12 {
		t.Fatalf("records: %d", got)
	}
}
