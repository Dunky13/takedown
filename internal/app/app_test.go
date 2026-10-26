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

type testAccount struct {
	client *http.Client
	link   string
}

func (a testAccount) Login(context.Context) error                              { return nil }
func (a testAccount) List(context.Context) ([]model.Export, error)             { return nil, nil }
func (a testAccount) RefreshLink(context.Context, string, int) (string, error) { return a.link, nil }
func (a testAccount) HTTPClient() *http.Client                                 { return a.client }
func (a testAccount) Forget(context.Context) error                             { return nil }
func (a testAccount) CreateExport(context.Context) (string, error)             { return "", nil }
func (a testAccount) Close() error                                             { return nil }

func TestDownloadProducesVerifiedArchiveAndDurableStatus(t *testing.T) {
	var data bytes.Buffer
	z := zip.NewWriter(&data)
	f, err := z.Create("Takeout/data.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = f.Write([]byte("export contents")); err != nil {
		t.Fatal(err)
	}
	if err = z.Close(); err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/zip")
		http.ServeContent(w, r, "takeout.zip", time.Time{}, bytes.NewReader(data.Bytes()))
	}))
	defer server.Close()
	dir := t.TempDir()
	dest := filepath.Join(dir, "archives")
	a := &App{DataDir: filepath.Join(dir, "state"), Settings: config.Settings{Destination: dest, SafetyBuffer: 1, Concurrency: 2}, Bus: events.New(), Account: testAccount{client: server.Client(), link: server.URL}}
	defer a.Close()
	p := model.Part{ExportID: "export", Index: 1, Filename: "takeout.zip", Size: int64(data.Len()), URL: server.URL}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = a.Download(ctx, []model.Part{p}, DownloadOptions{}); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(dest, filepath.FromSlash(model.NamespacePath(p.ExportID, p.Filename))))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, data.Bytes()) {
		t.Fatal("archive content changed")
	}
	if _, err = os.Stat(filepath.Join(dest, filepath.FromSlash(model.NamespacePath(p.ExportID, p.Filename))+".part")); !os.IsNotExist(err) {
		t.Fatalf("partial file remains: %v", err)
	}
	records, err := state.ReadSnapshot(filepath.Join(a.DataDir, "takedown-state.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 || !records[0].Complete || records[0].Offset != p.Size {
		t.Fatalf("bad durable status: %+v", records)
	}
	if err = a.Download(ctx, []model.Part{p}, DownloadOptions{}); err != nil {
		t.Fatalf("completed job could not restart: %v", err)
	}
}

func TestDestinationRejectsPassword(t *testing.T) {
	if _, err := DestinationIdentity("sftp://user:secret@nas.example/backup"); err == nil {
		t.Fatal("URL password accepted")
	}
}
