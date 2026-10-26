package app

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/Dunky13/takedown/internal/config"
	"github.com/Dunky13/takedown/internal/events"
	"github.com/Dunky13/takedown/internal/source"
)

type jobTransport func(*http.Request) (*http.Response, error)

func (f jobTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type jobSession struct{ testAccount }

func (jobSession) Refresh(context.Context) error             { return nil }
func (jobSession) ProfileID(context.Context) (string, error) { return strings.Repeat("a", 64), nil }

func TestWatchKeepsPollingWhenOnlyCurrentJobIsHeld(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var polls atomic.Int32
		client := &http.Client{Transport: jobTransport(func(r *http.Request) (*http.Response, error) {
			body, code := "", 200
			switch r.URL.Path {
			case "/v1/accessType:check":
				body = `{"timeBasedResources":["myactivity.search"]}`
			case "/v1/portabilityArchive:initiate":
				body = `{"archiveJobId":"held-job","accessType":"ACCESS_TYPE_TIME_BASED"}`
			case "/v1/archiveJobs/held-job/portabilityArchiveState":
				polls.Add(1)
				code = 404
				body = `{"error":{"code":404,"message":"missing job"}}`
			default:
				t.Fatalf("unexpected API request: %s", r.URL.Path)
			}
			return &http.Response{StatusCode: code, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
		})}
		dir := t.TempDir()
		manager, err := source.New(source.Options{DataDir: dir, Provider: "portability", Resources: []string{"myactivity.search"}, Session: jobSession{testAccount{client: client}}})
		if err != nil {
			t.Fatal(err)
		}
		if _, err = manager.CreateExport(context.Background()); err != nil {
			t.Fatal(err)
		}
		time.Sleep(5 * time.Minute)
		a := &App{DataDir: dir, Settings: config.Settings{Source: "portability", PortabilityResources: []string{"myactivity.search"}, Destination: t.TempDir(), SafetyBuffer: 1, Concurrency: 2}, Bus: events.New(), Account: manager}
		defer a.Close()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		done := make(chan error, 1)
		go func() { done <- a.Watch(ctx, WatchOptions{}) }()
		synctest.Wait()
		select {
		case err := <-done:
			t.Fatalf("held job terminated watch: %v", err)
		default:
		}
		if polls.Load() != 1 {
			t.Fatalf("first cadence: %d", polls.Load())
		}
		time.Sleep(5 * time.Minute)
		synctest.Wait()
		if polls.Load() != 2 {
			t.Fatalf("watch did not preserve future discovery: %d", polls.Load())
		}
		cancel()
		synctest.Wait()
		<-done
	})
}
