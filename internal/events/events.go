// Package events provides bounded, nonblocking subscriptions for UI telemetry.
package events

import (
	"sync"
	"time"
)

type Kind string

const (
	Progress       Kind = "progress"
	PartCompleted  Kind = "part_completed"
	SpaceWarning   Kind = "space_warning"
	ReAuthRequired Kind = "reauth_required"
	Status         Kind = "status"
	Error          Kind = "error"
)

type Event struct {
	Kind           Kind
	Time           time.Time
	PartKey        string
	Filename       string
	Message        string
	Bytes          int64
	Total          int64
	BytesPerSecond float64
	FreeBytes      uint64
}

type Bus struct {
	mu          sync.Mutex
	subscribers map[uint64]chan Event
	next        uint64
	closed      bool
}

func New() *Bus { return &Bus{subscribers: make(map[uint64]chan Event)} }

// Subscribe returns a bounded telemetry stream and an idempotent cancellation.
// Slow consumers lose the oldest event. Durable status belongs in state.Store.
func (b *Bus) Subscribe(capacity int) (<-chan Event, func()) {
	if capacity < 1 {
		capacity = 64
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan Event, capacity)
	if b.closed {
		close(ch)
		return ch, func() {}
	}
	b.next++
	id := b.next
	b.subscribers[id] = ch
	var once sync.Once
	return ch, func() {
		once.Do(func() {
			b.mu.Lock()
			defer b.mu.Unlock()
			if c, ok := b.subscribers[id]; ok {
				delete(b.subscribers, id)
				close(c)
			}
		})
	}
}
func (b *Bus) Publish(e Event) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if e.Time.IsZero() {
		e.Time = time.Now()
	}
	for _, ch := range b.subscribers {
		select {
		case ch <- e:
		default:
			select {
			case <-ch:
			default:
			}
			select {
			case ch <- e:
			default:
			}
		}
	}
}
func (b *Bus) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.closed = true
	for id, ch := range b.subscribers {
		delete(b.subscribers, id)
		close(ch)
	}
}
