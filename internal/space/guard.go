// Package space prevents concurrent downloads from overcommitting a destination.
package space

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"

	"github.com/Dunky13/takedown/internal/events"
	"github.com/Dunky13/takedown/internal/storage"
)

const DefaultSafetyBuffer uint64 = 2_000_000_000

var ErrInsufficientSpace = errors.New("insufficient destination space")

type InsufficientSpaceError struct {
	FreeBytes     uint64
	RequiredBytes uint64
}

func (e *InsufficientSpaceError) Error() string {
	return fmt.Sprintf("%v: %d bytes available, %d bytes required including safety buffer", ErrInsufficientSpace, e.FreeBytes, e.RequiredBytes)
}
func (e *InsufficientSpaceError) Unwrap() error { return ErrInsufficientSpace }

type Guard struct {
	mu            sync.Mutex
	backend       storage.StorageBackend
	buffer        uint64
	bus           *events.Bus
	remaining     map[string]uint64
	allowUnknown  bool
	warnedUnknown bool
}

func New(backend storage.StorageBackend, buffer uint64, bus *events.Bus) *Guard {
	g := &Guard{backend: backend, buffer: buffer, bus: bus, remaining: make(map[string]uint64)}
	if policy, ok := backend.(interface{ UnknownSpaceAllowed() bool }); ok {
		g.allowUnknown = policy.UnknownSpaceAllowed()
	}
	return g
}

// AllowUnknownSpace explicitly permits downloads when the provider cannot report
// free space. The first such observation emits a warning; no capacity is invented.
func (g *Guard) AllowUnknownSpace(allow bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.allowUnknown = allow
}

// Preflight checks an additional unreserved batch alongside active reservations.
// It does not reserve bytes. Reserve each worker before starting its writes.
func (g *Guard) Preflight(ctx context.Context, remaining int64) error {
	if remaining < 0 {
		return errors.New("remaining byte count cannot be negative")
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	total, err := g.total(uint64(remaining))
	if err != nil {
		return g.warning("", 0, err)
	}
	return g.checkSpace(ctx, "", total)
}

// Reserve atomically admits a new worker only when all active remaining bytes
// plus the buffer fit. Duplicate keys are rejected to avoid untracked workers.
func (g *Guard) Reserve(ctx context.Context, key string, remaining int64) error {
	if key == "" || remaining < 0 {
		return errors.New("reservation requires a key and nonnegative remaining bytes")
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, exists := g.remaining[key]; exists {
		return errors.New("space reservation already exists")
	}
	total, err := g.total(uint64(remaining))
	if err != nil {
		return g.warning(key, 0, err)
	}
	if err := g.checkSpace(ctx, key, total); err != nil {
		return err
	}
	g.remaining[key] = uint64(remaining)
	return nil
}

// Check records actual progress even when space is now insufficient. Other
// workers always see the current total, including every active reservation.
// The caller decides whether to pause/retry; this method never waits for space.
func (g *Guard) Check(ctx context.Context, key string, remaining int64) error {
	if key == "" || remaining < 0 {
		return errors.New("reservation requires a key and nonnegative remaining bytes")
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, exists := g.remaining[key]; !exists {
		return errors.New("space reservation does not exist")
	}
	g.remaining[key] = uint64(remaining)
	total, err := g.total(0)
	if err != nil {
		return g.warning(key, 0, err)
	}
	return g.checkSpace(ctx, key, total)
}

func (g *Guard) Release(key string) { g.mu.Lock(); defer g.mu.Unlock(); delete(g.remaining, key) }

func (g *Guard) total(extra uint64) (uint64, error) {
	total := extra
	for _, value := range g.remaining {
		if math.MaxUint64-total < value {
			return 0, fmt.Errorf("%w: remaining byte count overflows uint64", ErrInsufficientSpace)
		}
		total += value
	}
	if math.MaxUint64-total < g.buffer {
		return 0, fmt.Errorf("%w: safety buffer overflows uint64", ErrInsufficientSpace)
	}
	return total + g.buffer, nil
}

func (g *Guard) checkSpace(ctx context.Context, key string, required uint64) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	free, err := g.backend.FreeSpace(ctx)
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if errors.Is(err, storage.ErrSpaceUnknown) && g.allowUnknown {
			if !g.warnedUnknown {
				g.warnedUnknown = true
				g.publish(key, 0, "Destination free space is unavailable; continuing under the explicitly enabled warning policy")
			}
			return nil
		}
		return g.warning(key, 0, err)
	}
	if free < required {
		return g.warning(key, free, &InsufficientSpaceError{FreeBytes: free, RequiredBytes: required})
	}
	return nil
}

func (g *Guard) warning(key string, free uint64, err error) error {
	g.publish(key, free, err.Error())
	return err
}
func (g *Guard) publish(key string, free uint64, message string) {
	if g.bus != nil {
		g.bus.Publish(events.Event{Kind: events.SpaceWarning, PartKey: key, FreeBytes: free, Message: message})
	}
}
