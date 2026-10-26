package space

import (
	"context"
	"errors"
	"io"
	"math"
	"sync"
	"testing"

	"github.com/Dunky13/takedown/internal/events"
	"github.com/Dunky13/takedown/internal/storage"
)

type fakeBackend struct {
	mu           sync.Mutex
	free         uint64
	err          error
	queries      int
	allowUnknown bool
}

func (b *fakeBackend) set(free uint64, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.free, b.err = free, err
}
func (b *fakeBackend) FreeSpace(ctx context.Context) (uint64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.queries++
	return b.free, b.err
}
func (*fakeBackend) Stat(context.Context, string) (storage.FileInfo, error) {
	return nil, errors.New("unused")
}
func (*fakeBackend) OpenWriteResume(context.Context, string, int64) (io.WriteCloser, error) {
	return nil, errors.New("unused")
}
func (*fakeBackend) Rename(context.Context, string, string) error { return errors.New("unused") }
func (*fakeBackend) Close() error                                 { return nil }
func (b *fakeBackend) UnknownSpaceAllowed() bool                  { return b.allowUnknown }

func TestReservationsIncludeAllWorkersAndProgress(t *testing.T) {
	b := &fakeBackend{free: 100}
	g := New(b, 10, nil)
	ctx := context.Background()
	if err := g.Reserve(ctx, "a", 60); err != nil {
		t.Fatal(err)
	}
	if err := g.Reserve(ctx, "b", 30); err != nil {
		t.Fatal(err)
	}
	if err := g.Reserve(ctx, "c", 1); !errors.Is(err, ErrInsufficientSpace) {
		t.Fatalf("overcommit: %v", err)
	}
	b.set(80, nil) // Worker a wrote 20 bytes.
	if err := g.Check(ctx, "a", 40); err != nil {
		t.Fatalf("progress was double counted: %v", err)
	}
	if err := g.Check(ctx, "b", 30); err != nil {
		t.Fatal(err)
	}
	b.set(30, nil) // External consumption forces a pause.
	if err := g.Check(ctx, "a", 20); !errors.Is(err, ErrInsufficientSpace) {
		t.Fatalf("external usage ignored: %v", err)
	}
	b.set(80, nil)
	if err := g.Reserve(ctx, "c", 20); err != nil {
		t.Fatalf("failed Check did not retain updated progress: %v", err)
	}
	g.Release("b")
	g.Release("b")
	if err := g.Reserve(ctx, "d", 30); err != nil {
		t.Fatalf("released capacity not available: %v", err)
	}
}

func TestConcurrentReservationsCannotOvercommit(t *testing.T) {
	g := New(&fakeBackend{free: 100}, 10, nil)
	start := make(chan struct{})
	results := make(chan error, 2)
	var wg sync.WaitGroup
	for _, key := range []string{"a", "b"} {
		wg.Add(1)
		go func(key string) { defer wg.Done(); <-start; results <- g.Reserve(context.Background(), key, 60) }(key)
	}
	close(start)
	wg.Wait()
	close(results)
	admitted, denied := 0, 0
	for err := range results {
		if err == nil {
			admitted++
		} else if errors.Is(err, ErrInsufficientSpace) {
			denied++
		} else {
			t.Fatal(err)
		}
	}
	if admitted != 1 || denied != 1 {
		t.Fatalf("admitted=%d denied=%d", admitted, denied)
	}
}

func TestPreflightCountsReservationsAndExactBuffer(t *testing.T) {
	g := New(&fakeBackend{free: 100}, 10, nil)
	ctx := context.Background()
	if err := g.Preflight(ctx, 90); err != nil {
		t.Fatal(err)
	}
	if err := g.Preflight(ctx, 91); !errors.Is(err, ErrInsufficientSpace) {
		t.Fatalf("buffer ignored: %v", err)
	}
	if err := g.Reserve(ctx, "active", 50); err != nil {
		t.Fatal(err)
	}
	if err := g.Preflight(ctx, 40); err != nil {
		t.Fatal(err)
	}
	if err := g.Preflight(ctx, 41); !errors.Is(err, ErrInsufficientSpace) {
		t.Fatalf("preflight ignored active reservations: %v", err)
	}
}

func TestUnknownSpaceFailsClosedAndWarnsOnceWhenEnabled(t *testing.T) {
	b := &fakeBackend{err: &storage.SpaceUnknownError{Cause: errors.New("statvfs unsupported")}}
	bus := events.New()
	defer bus.Close()
	stream, cancel := bus.Subscribe(8)
	defer cancel()
	g := New(b, 10, bus)
	ctx := context.Background()
	if err := g.Preflight(ctx, 1); !errors.Is(err, storage.ErrSpaceUnknown) {
		t.Fatalf("unknown space allowed by default: %v", err)
	}
	<-stream // Default failure notification.
	g.AllowUnknownSpace(true)
	if err := g.Reserve(ctx, "a", 60); err != nil {
		t.Fatal(err)
	}
	if err := g.Reserve(ctx, "b", 60); err != nil {
		t.Fatal(err)
	}
	if err := g.Check(ctx, "a", 30); err != nil {
		t.Fatal(err)
	}
	select {
	case event := <-stream:
		if event.Kind != events.SpaceWarning || event.PartKey != "a" || event.FreeBytes != 0 {
			t.Fatalf("wrong warning: %+v", event)
		}
	default:
		t.Fatal("missing opt-in warning")
	}
	select {
	case e := <-stream:
		t.Fatalf("repeated warning: %+v", e)
	default:
	}
	g.AllowUnknownSpace(false)
	if err := g.Check(ctx, "a", 20); !errors.Is(err, storage.ErrSpaceUnknown) {
		t.Fatal("policy could not be disabled")
	}
	broken := errors.New("transport disconnected")
	b.set(0, broken)
	g.AllowUnknownSpace(true)
	if err := g.Check(ctx, "a", 20); !errors.Is(err, broken) {
		t.Fatalf("non-space failure ignored: %v", err)
	}
}

func TestExplicitBackendPolicyEnablesUnknownSpace(t *testing.T) {
	g := New(&fakeBackend{err: storage.ErrSpaceUnknown, allowUnknown: true}, DefaultSafetyBuffer, nil)
	if err := g.Reserve(context.Background(), "a", 100); err != nil {
		t.Fatal(err)
	}
}

func TestSpaceWarningContainsTypedFailureAndReservation(t *testing.T) {
	bus := events.New()
	defer bus.Close()
	stream, cancel := bus.Subscribe(1)
	defer cancel()
	g := New(&fakeBackend{free: 30}, 10, bus)
	err := g.Reserve(context.Background(), "export/1", 21)
	var detail *InsufficientSpaceError
	if !errors.As(err, &detail) || detail.FreeBytes != 30 || detail.RequiredBytes != 31 {
		t.Fatalf("missing detail: %v", err)
	}
	event := <-stream
	if event.Kind != events.SpaceWarning || event.FreeBytes != 30 || event.PartKey != "export/1" || event.Message == "" {
		t.Fatalf("wrong event: %+v", event)
	}
	if _, ok := g.remaining["export/1"]; ok {
		t.Fatal("failed reservation retained")
	}
}

func TestOverflowNeverWrapsIntoAdmission(t *testing.T) {
	ctx := context.Background()
	for _, operation := range []string{"preflight", "reserve"} {
		g := New(&fakeBackend{free: math.MaxUint64}, math.MaxUint64, nil)
		var err error
		if operation == "preflight" {
			err = g.Preflight(ctx, 1)
		} else {
			err = g.Reserve(ctx, "a", 1)
		}
		if !errors.Is(err, ErrInsufficientSpace) {
			t.Fatalf("%s overflow allowed: %v", operation, err)
		}
	}
	g := New(&fakeBackend{free: math.MaxUint64}, 0, nil)
	if err := g.Reserve(ctx, "a", math.MaxInt64); err != nil {
		t.Fatal(err)
	}
	if err := g.Reserve(ctx, "b", math.MaxInt64); err != nil {
		t.Fatal(err)
	}
	if err := g.Reserve(ctx, "c", 2); !errors.Is(err, ErrInsufficientSpace) {
		t.Fatalf("sum overflow allowed: %v", err)
	}
	if err := g.Reserve(ctx, "c", 1); err != nil {
		t.Fatalf("exact uint64 limit rejected: %v", err)
	}
	if err := g.Check(ctx, "c", 2); !errors.Is(err, ErrInsufficientSpace) {
		t.Fatalf("Check overflow allowed: %v", err)
	}
	g.Release("c")
	if err := g.Check(ctx, "a", 0); err != nil {
		t.Fatal(err)
	}
}

func TestInvalidAndCanceledRequestsDoNotReserve(t *testing.T) {
	b := &fakeBackend{free: 100}
	g := New(b, 0, nil)
	ctx := context.Background()
	if err := g.Reserve(ctx, "", 1); err == nil {
		t.Fatal("empty key accepted")
	}
	if err := g.Reserve(ctx, "a", -1); err == nil {
		t.Fatal("negative reservation accepted")
	}
	if err := g.Preflight(ctx, -1); err == nil {
		t.Fatal("negative preflight accepted")
	}
	if err := g.Check(ctx, "missing", 1); err == nil {
		t.Fatal("unreserved worker accepted")
	}
	canceled, cancel := context.WithCancel(ctx)
	cancel()
	if err := g.Reserve(canceled, "a", 1); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancel: %v", err)
	}
	if b.queries != 0 || len(g.remaining) != 0 {
		t.Fatal("invalid calls changed state or queried destination")
	}
	if err := g.Reserve(ctx, "a", 1); err != nil {
		t.Fatal(err)
	}
	if err := g.Reserve(ctx, "a", 1); err == nil {
		t.Fatal("duplicate worker accepted")
	}
}
