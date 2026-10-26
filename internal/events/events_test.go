package events

import (
	"sync"
	"testing"
)

func TestSlowConsumerKeepsLatest(t *testing.T) {
	b := New()
	defer b.Close()
	ch, stop := b.Subscribe(1)
	defer stop()
	b.Publish(Event{Kind: Progress, Bytes: 1})
	b.Publish(Event{Kind: Progress, Bytes: 2})
	if e := <-ch; e.Bytes != 2 || e.Time.IsZero() {
		t.Fatalf("unexpected event: %+v", e)
	}
}
func TestConcurrentCloseAndPublish(t *testing.T) {
	b := New()
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, stop := b.Subscribe(2)
			defer stop()
			for j := 0; j < 100; j++ {
				b.Publish(Event{Kind: Status})
			}
		}()
	}
	b.Close()
	wg.Wait()
	ch, stop := b.Subscribe(1)
	defer stop()
	if _, ok := <-ch; ok {
		t.Fatal("closed bus returned open stream")
	}
}
