package storage

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type syncWitness struct {
	*os.File
	calls   int
	failure error
}

func (f *syncWitness) Sync() error { f.calls++; return f.failure }

func TestWriterForwardsSyncAndPropagatesFailure(t *testing.T) {
	file, err := os.Create(filepath.Join(t.TempDir(), "archive.part"))
	if err != nil {
		t.Fatal(err)
	}
	want := errors.New("sync failed")
	witness := &syncWitness{File: file, failure: want}
	var handles fileSet
	w, err := handles.add(context.Background(), witness, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = w.Write([]byte("archive bytes")); err != nil {
		t.Fatal(err)
	}
	if err = w.Sync(); !errors.Is(err, want) {
		t.Fatalf("sync failure lost: %v", err)
	}
	if err = w.Close(); !errors.Is(err, want) {
		t.Fatalf("close lost failed flush receipt: %v", err)
	}
	if witness.calls < 1 {
		t.Fatal("underlying sync never called")
	}
}
func TestCanceledLocalWriterCanSyncBeforeCheckpoint(t *testing.T) {
	b, _ := localForTest(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w, err := b.OpenWriteResume(ctx, "archive.part", 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = w.Write([]byte("retained prefix")); err != nil {
		t.Fatal(err)
	}
	cancel()
	syncer, ok := w.(interface{ Sync() error })
	if !ok {
		t.Fatal("writer hides sync")
	}
	if err = syncer.Sync(); err != nil {
		t.Fatalf("cancellation prevented cleanup sync: %v", err)
	}
	if err = w.Close(); err != nil {
		t.Fatal(err)
	}
	if err = syncer.Sync(); err != nil {
		t.Fatalf("closed writer lost successful sync receipt: %v", err)
	}
}
