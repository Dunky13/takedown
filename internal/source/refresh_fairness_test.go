package source

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"google.golang.org/api/dataportability/v1"
)

func TestRefreshJobsFairAcrossCadenceAndRestart(t *testing.T) {
	for _, restart := range []bool{false, true} {
		t.Run(fmt.Sprintf("restart=%t", restart), func(t *testing.T) {
			now := time.Date(2026, 10, 25, 12, 0, 0, 0, time.UTC)
			seen := map[string]int{}
			var calls []string
			handler := func(r *http.Request) (*http.Response, error) {
				if r.Method != http.MethodGet {
					t.Fatal("refresh sent a mutation")
				}
				id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/v1/archiveJobs/"), "/portabilityArchiveState")
				seen[id]++
				calls = append(calls, id)
				return responseJSON(r, 200, &dataportability.PortabilityArchiveState{State: "IN_PROGRESS"}), nil
			}
			opts := Options{Provider: "portability", Resources: []string{"myactivity.search"}, DataDir: t.TempDir()}
			m := testManager(t, opts, handler)
			m.now = func() time.Time { return now }
			jobs := make([]portabilityJob, 201)
			for i := range jobs {
				jobs[i] = testJob(now.Add(-time.Hour))
				// Journal insertion order must not define scheduling priority.
				jobs[i].ID = fmt.Sprintf("job-%03d", len(jobs)-1-i)
			}
			seedJobs(t, m, jobs...)
			for round := 0; round < 3; round++ {
				if restart && round > 0 {
					m = testManager(t, opts, handler)
					m.now = func() time.Time { return now }
				}
				before := len(calls)
				if _, err := m.RefreshJobs(context.Background()); err != nil {
					t.Fatal(err)
				}
				if len(calls)-before != 100 {
					t.Fatalf("round %d: reads=%d, want exactly the bounded 100", round, len(calls)-before)
				}
				if want := fmt.Sprintf("job-%03d", 100*round); calls[before] != want {
					t.Fatalf("round %d: first=%s, want oldest unvisited %s", round, calls[before], want)
				}
				if got := readJournal(t, m).Jobs[0].ID; got != "job-200" {
					t.Fatal("scheduling reordered the durable job history", got)
				}
				now = now.Add(5 * time.Minute)
			}
			if len(seen) != len(jobs) {
				t.Fatalf("cadence-spaced refreshes starved jobs: visited=%d of %d", len(seen), len(jobs))
			}
		})
	}
}

func TestRefreshJobsOldestEligiblePollAndDeterministicTies(t *testing.T) {
	now := time.Date(2026, 10, 25, 12, 0, 0, 0, time.UTC)
	var calls []string
	handler := func(r *http.Request) (*http.Response, error) {
		id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/v1/archiveJobs/"), "/portabilityArchiveState")
		calls = append(calls, id)
		return responseJSON(r, 200, &dataportability.PortabilityArchiveState{State: "IN_PROGRESS"}), nil
	}
	opts := Options{Provider: "portability", Resources: []string{"myactivity.search"}, DataDir: t.TempDir()}
	m := testManager(t, opts, handler)
	m.now = func() time.Time { return now }
	jobs := make([]portabilityJob, 100)
	for i := range jobs {
		jobs[i] = testJob(now.Add(-time.Hour))
		jobs[i].ID = fmt.Sprintf("same-%03d", 99-i)
		jobs[i].LastPoll = now.Add(-10 * time.Minute)
	}
	older, newer, backoff, cooldown := testJob(now.Add(-time.Hour)), testJob(now.Add(-time.Hour)), testJob(now.Add(-time.Hour)), testJob(now.Add(-time.Hour))
	older.ID = "z-never-polled"
	newer.ID, newer.LastPoll = "a-newer", now.Add(-5*time.Minute)
	backoff.ID, backoff.RetryNotBefore = "backoff", now.Add(time.Hour)
	cooldown.ID, cooldown.LastPoll = "cooldown", now.Add(-4*time.Minute)
	jobs = append([]portabilityJob{newer, backoff, cooldown}, append(jobs, older)...)
	seedJobs(t, m, jobs...)
	if _, err := m.RefreshJobs(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(calls) != 100 || calls[0] != older.ID {
		t.Fatalf("oldest eligible request did not lead bounded batch: %v", calls)
	}
	for i := 1; i < len(calls); i++ {
		if want := fmt.Sprintf("same-%03d", i-1); calls[i] != want {
			t.Fatalf("tied LastPoll order=%s, want %s", calls[i], want)
		}
	}
	// Persisted LastPoll, rather than an in-memory cursor, orders the next
	// cadence-spaced call after a restart. Backoff still excludes its job.
	now = now.Add(5 * time.Minute)
	m = testManager(t, opts, handler)
	m.now = func() time.Time { return now }
	calls = nil
	if _, err := m.RefreshJobs(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(calls) != 100 || calls[0] != "same-099" || calls[1] != newer.ID || calls[2] != cooldown.ID {
		t.Fatalf("restart did not prioritize oldest due jobs: %v", calls)
	}
	for _, id := range calls {
		if id == backoff.ID {
			t.Fatal("fair scheduling bypassed provider backoff")
		}
	}
	for _, j := range readJournal(t, m).Jobs {
		if j.ID == backoff.ID && (!j.LastPoll.IsZero() || !j.RetryNotBefore.Equal(backoff.RetryNotBefore)) {
			t.Fatal("fair scheduling changed an ineligible job's durable cadence")
		}
	}
}
