package app

import (
	"archive/zip"
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Dunky13/takedown/internal/config"
	"github.com/Dunky13/takedown/internal/events"
	"github.com/Dunky13/takedown/internal/model"
	"github.com/Dunky13/takedown/internal/state"
)

func layoutArchive(t *testing.T, payload string) []byte {
	t.Helper()
	var b bytes.Buffer
	z := zip.NewWriter(&b)
	w, err := z.Create("data.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = w.Write([]byte(payload)); err != nil {
		t.Fatal(err)
	}
	if err = z.Close(); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

func TestDistinctExportsWithSameFilenameDownloadAndRestart(t *testing.T) {
	aBytes, bBytes := layoutArchive(t, "first export"), layoutArchive(t, "second export")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := aBytes
		if r.URL.Path == "/second" {
			data = bBytes
		}
		w.Header().Set("Content-Type", "application/zip")
		http.ServeContent(w, r, "takeout.zip", time.Time{}, bytes.NewReader(data))
	}))
	defer server.Close()
	dir := t.TempDir()
	a := &App{DataDir: filepath.Join(dir, "state"), Settings: config.Settings{Destination: filepath.Join(dir, "archives"), SafetyBuffer: 1, Concurrency: 2}, Bus: events.New(), Account: testAccount{client: server.Client()}}
	defer a.Close()
	parts := []model.Part{
		{ExportID: "first", Index: 1, Filename: "takeout.zip", Size: int64(len(aBytes)), URL: server.URL + "/first"},
		{ExportID: "second", Index: 1, Filename: "takeout.zip", Size: int64(len(bBytes)), URL: server.URL + "/second"},
	}
	for range 2 {
		if err := a.Download(context.Background(), parts, DownloadOptions{}); err != nil {
			t.Fatal(err)
		}
	}
	records, err := state.ReadSnapshot(filepath.Join(a.DataDir, "takedown-state.json"))
	if err != nil || len(records) != 2 {
		t.Fatalf("records: %v, %v", records, err)
	}
	for _, r := range records {
		if !r.Complete || r.SHA256 == "" || r.Part.RelativePath == "" {
			t.Fatalf("unbound completion: %+v", r)
		}
		got, err := os.ReadFile(filepath.Join(a.Settings.Destination, filepath.FromSlash(r.Part.StoragePath())))
		if err != nil {
			t.Fatal(err)
		}
		want := aBytes
		if r.Part.ExportID == "second" {
			want = bBytes
		}
		if !bytes.Equal(got, want) {
			t.Fatal("export contents collided")
		}
	}
}

func TestLegacyFlatPartialRemainsResumable(t *testing.T) {
	data := layoutArchive(t, "legacy partial export")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/zip")
		http.ServeContent(w, r, "archive.zip", time.Time{}, bytes.NewReader(data))
	}))
	defer server.Close()
	dir := t.TempDir()
	dest := filepath.Join(dir, "archives")
	if err := os.MkdirAll(dest, 0700); err != nil {
		t.Fatal(err)
	}
	p := model.Part{ExportID: "legacy", Index: 1, Filename: "archive.zip", Size: int64(len(data)), URL: server.URL}
	if err := os.WriteFile(filepath.Join(dest, p.Filename+".part"), data[:30], 0600); err != nil {
		t.Fatal(err)
	}
	s, err := state.Open(filepath.Join(dir, "state", "takedown-state.json"))
	if err != nil {
		t.Fatal(err)
	}
	id, err := DestinationIdentity(dest)
	if err != nil {
		t.Fatal(err)
	}
	if err = s.Put(state.Record{Part: p, Destination: id, Offset: 30}); err != nil {
		t.Fatal(err)
	}
	if err = s.Close(); err != nil {
		t.Fatal(err)
	}
	a := &App{DataDir: filepath.Join(dir, "state"), Settings: config.Settings{Destination: dest, SafetyBuffer: 1, Concurrency: 2}, Bus: events.New(), Account: testAccount{client: server.Client()}}
	defer a.Close()
	if err = a.Download(context.Background(), []model.Part{p}, DownloadOptions{}); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(dest, p.Filename))
	if err != nil || !bytes.Equal(got, data) {
		t.Fatalf("legacy partial moved or lost: %v", err)
	}
}
