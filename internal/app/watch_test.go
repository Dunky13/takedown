package app

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Dunky13/takedown/internal/config"
	"github.com/Dunky13/takedown/internal/events"
	"github.com/Dunky13/takedown/internal/model"
)

type watchAccount struct {
	testAccount
	exports []model.Export
	created int
}

func (w *watchAccount) List(context.Context) ([]model.Export, error) { return w.exports, nil }
func (w *watchAccount) CreateExport(context.Context) (string, error) {
	w.created++
	return "unexpected", nil
}

func TestWatchDownloadsNewArchivesOnceAndRecoversDeletedFiles(t *testing.T) {
	var body bytes.Buffer
	z := zip.NewWriter(&body)
	f, err := z.Create("Takeout/data.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = f.Write([]byte("export data")); err != nil {
		t.Fatal(err)
	}
	if err = z.Close(); err != nil {
		t.Fatal(err)
	}
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		w.Header().Set("Content-Type", "application/zip")
		http.ServeContent(w, r, "archive.zip", time.Time{}, bytes.NewReader(body.Bytes()))
	}))
	defer server.Close()
	dir := t.TempDir()
	account := &watchAccount{testAccount: testAccount{client: server.Client(), link: server.URL}}
	a := &App{DataDir: filepath.Join(dir, "state"), Settings: config.Settings{Destination: filepath.Join(dir, "archives"), SafetyBuffer: 1, Concurrency: 2}, Account: account, Bus: events.New()}
	defer a.Close()
	known := map[string]watchedArchive{}
	ctx := context.Background()
	if err = a.watchPass(ctx, WatchOptions{}, known); err != nil {
		t.Fatal(err)
	}
	if account.created != 0 || requests.Load() != 0 {
		t.Fatal("watch created an export or downloaded without a ready archive")
	}
	p := model.Part{ExportID: "drive-file-version", Index: 1, Filename: "archive.zip", Size: int64(body.Len()), URL: server.URL}
	account.exports = []model.Export{{ID: p.ExportID, Parts: []model.Part{p}}}
	if err = a.watchPass(ctx, WatchOptions{}, known); err != nil {
		t.Fatal(err)
	}
	if requests.Load() != 1 {
		t.Fatalf("first scan requests = %d", requests.Load())
	}
	if err = a.watchPass(ctx, WatchOptions{}, known); err != nil {
		t.Fatal(err)
	}
	if requests.Load() != 1 {
		t.Fatal("watch downloaded a completed archive again")
	}
	if err = os.Remove(filepath.Join(a.Settings.Destination, filepath.FromSlash(model.NamespacePath(p.ExportID, p.Filename)))); err != nil {
		t.Fatal(err)
	}
	if err = a.watchPass(ctx, WatchOptions{}, known); err != nil {
		t.Fatal(err)
	}
	if requests.Load() != 2 {
		t.Fatal("watch did not recover a removed archive")
	}
}

func TestWatchCancelsWhileWaitingAndKeepsPreferencesLocked(t *testing.T) {
	a := &App{DataDir: t.TempDir(), Settings: config.Settings{Destination: t.TempDir(), SafetyBuffer: 1, Concurrency: 2}, Account: &watchAccount{}, Bus: events.New()}
	defer a.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eventsCh, stop := a.Bus.Subscribe(8)
	defer stop()
	result := make(chan error, 1)
	go func() { result <- a.Watch(ctx, WatchOptions{}) }()
	select {
	case <-eventsCh:
	case <-time.After(5 * time.Second):
		t.Fatal("watch did not enter its wait")
	}
	if err := a.Configure(a.Settings); err == nil {
		t.Fatal("preferences changed while watch was running")
	}
	cancel()
	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("watch cancellation stalled")
	}
}
