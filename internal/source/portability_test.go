package source

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"google.golang.org/api/dataportability/v1"
)

func testJob(at time.Time) portabilityJob {
	return portabilityJob{ID: "job-123", Resources: []string{"myactivity.search"}, CreatedAt: at, AccessType: "ACCESS_TYPE_TIME_BASED", State: "IN_PROGRESS"}
}
func seedJobs(t *testing.T, m *Manager, jobs ...portabilityJob) string {
	t.Helper()
	state, name, unlock, err := m.openState(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer unlock()
	state.Jobs = jobs
	if err = saveState(name, state); err != nil {
		t.Fatal(err)
	}
	return name
}

func statePath(t *testing.T, m *Manager) string {
	t.Helper()
	id, err := m.opts.Session.ProfileID(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(m.opts.DataDir, "source", strings.ToLower(id), "portability.json")
}
func signedResponse(r *http.Request, name string, size int64) *http.Response {
	return &http.Response{StatusCode: http.StatusPartialContent, Header: http.Header{"Content-Type": []string{"application/zip"}, "Content-Range": []string{fmt.Sprintf("bytes 0-0/%d", size)}, "Content-Disposition": []string{`attachment; filename="` + name + `"`}, "X-Goog-Generation": []string{"1234"}}, ContentLength: 1, Body: io.NopCloser(strings.NewReader("")), Request: r}
}
func installSigned(t *testing.T, m *Manager, handler roundTripFunc) {
	t.Helper()
	m.public.CloseIdleConnections()
	m.public.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.Header.Get("Authorization") != "" || r.Header.Get("Cookie") != "" {
			t.Fatal("OAuth/cookies on signed URL request")
		}
		if r.Header.Get("Accept-Encoding") != "identity" {
			t.Fatal("encoded metadata")
		}
		if r.Header.Get("Range") != "bytes=0-0" {
			t.Fatal("metadata must be a single-byte GET")
		}
		return handler(r)
	})
}

func TestPortabilityExplicitCreatePollRefreshAndPersistence(t *testing.T) {
	now := time.Date(2026, 10, 25, 10, 0, 0, 0, time.UTC)
	posts, polls, probes := 0, 0, 0
	signature := "first-secret-signature"
	m := testManager(t, Options{Provider: "portability", Resources: []string{"myactivity.search"}}, func(r *http.Request) (*http.Response, error) {
		switch r.URL.Path {
		case "/v1/accessType:check":
			return responseJSON(r, 200, &dataportability.CheckAccessTypeResponse{TimeBasedResources: []string{"myactivity.search"}}), nil
		case "/v1/portabilityArchive:initiate":
			posts++
			if r.Method != http.MethodPost {
				t.Fatal("creation method")
			}
			var request dataportability.InitiatePortabilityArchiveRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(request.Resources, []string{"myactivity.search"}) {
				t.Fatal("unselected resources", request.Resources)
			}
			return responseJSON(r, 200, &dataportability.InitiatePortabilityArchiveResponse{ArchiveJobId: "job-123", AccessType: "ACCESS_TYPE_TIME_BASED"}), nil
		case "/v1/archiveJobs/job-123/portabilityArchiveState":
			polls++
			if r.Method != http.MethodGet {
				t.Fatal("poll method")
			}
			return responseJSON(r, 200, &dataportability.PortabilityArchiveState{State: "COMPLETE", Urls: []string{"https://storage.googleapis.com/bucket/first.zip?X-Goog-Signature=" + signature, "https://storage.googleapis.com/bucket/second.tgz?X-Goog-Signature=" + signature}}), nil
		}
		t.Fatal("unexpected API request", r.URL.Path)
		return nil, nil
	})
	m.now = func() time.Time { return now }
	installSigned(t, m, func(r *http.Request) (*http.Response, error) {
		probes++
		if r.Method != http.MethodGet {
			t.Fatal("method-bound signed URL must use GET")
		}
		if strings.HasSuffix(r.URL.Path, "first.zip") {
			return signedResponse(r, "takeout-001.zip", 100), nil
		}
		return signedResponse(r, "takeout-002.tgz", 200), nil
	})
	if _, err := m.List(context.Background()); err != nil {
		t.Fatal(err)
	}
	if posts != 0 || polls != 0 {
		t.Fatal("list created a job or polled unrecorded jobs")
	}
	if _, err := m.CreateExport(context.Background()); err != nil {
		t.Fatal(err)
	}
	if posts != 1 {
		t.Fatal("creation count", posts)
	}
	if _, err := m.CreateExport(context.Background()); err == nil || posts != 1 {
		t.Fatal("24h duplicate request")
	}
	if _, err := m.List(context.Background()); err != nil {
		t.Fatal(err)
	}
	if polls != 0 {
		t.Fatal("polled before minimum cadence")
	}
	now = now.Add(5 * time.Minute)
	exports, err := m.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(exports) != 1 || len(exports[0].Parts) != 2 || polls != 1 || probes != 2 {
		t.Fatalf("exports=%v polls=%d probes=%d", exports, polls, probes)
	}
	for i, p := range exports[0].Parts {
		if p.Index != i+1 || p.ExportID != "portability-job-123" {
			t.Fatal("part ordering", p)
		}
	}
	stateFile := statePath(t, m)
	data, err := os.ReadFile(stateFile)
	if err != nil {
		t.Fatal(err)
	}
	for _, secret := range []string{signature, "https://", "X-Goog-", "token"} {
		if strings.Contains(string(data), secret) {
			t.Fatalf("secret in state: %s", secret)
		}
	}
	if !strings.Contains(string(data), "job-123") || !strings.Contains(string(data), "myactivity.search") {
		t.Fatal("missing durable receipt")
	}
	if _, err = m.List(context.Background()); err != nil {
		t.Fatal(err)
	}
	if polls != 1 || probes != 2 {
		t.Fatal("complete cache refetched within cadence")
	}
	now = now.Add(5 * time.Minute)
	signature = "second-secret-signature"
	link, err := m.RefreshLink(context.Background(), exports[0].ID, 2)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(link, signature) || polls != 2 || probes != 4 || posts != 1 {
		t.Fatal("refresh created or mismatched a job")
	}
	// New process/manager has no signed URLs. Metadata and cooldown survive.
	m2 := testManager(t, Options{Provider: "portability", Resources: m.opts.Resources, DataDir: m.opts.DataDir}, roundTripFunc(func(r *http.Request) (*http.Response, error) {
		polls++
		return responseJSON(r, 200, &dataportability.PortabilityArchiveState{State: "COMPLETE", Urls: []string{"https://storage.googleapis.com/bucket/first.zip?new-secret", "https://storage.googleapis.com/bucket/second.tgz?new-secret"}}), nil
	}))
	m2.now = func() time.Time { return now }
	installSigned(t, m2, func(r *http.Request) (*http.Response, error) {
		if strings.HasSuffix(r.URL.Path, "first.zip") {
			return signedResponse(r, "takeout-001.zip", 100), nil
		}
		return signedResponse(r, "takeout-002.tgz", 200), nil
	})
	if result, err := m2.List(context.Background()); err != nil || len(result) != 0 || polls != 2 {
		t.Fatal("restart bypassed persisted cooldown", result, err)
	}
	now = now.Add(5 * time.Minute)
	if _, err = m2.RefreshLink(context.Background(), exports[0].ID, 1); err != nil {
		t.Fatal("durable metadata could not be refreshed", err)
	}
}

func TestPortabilityUncertainCreationNeverDuplicatesAcrossRestart(t *testing.T) {
	for _, failure := range []string{"network", "server", "malformed", "cancel"} {
		t.Run(failure, func(t *testing.T) {
			posts := 0
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			var m *Manager
			m = testManager(t, Options{Provider: "portability", Resources: []string{"myactivity.search"}}, func(r *http.Request) (*http.Response, error) {
				if r.URL.Path == "/v1/accessType:check" {
					return responseJSON(r, 200, &dataportability.CheckAccessTypeResponse{TimeBasedResources: []string{"myactivity.search"}}), nil
				}
				posts++
				data, err := os.ReadFile(statePath(t, m))
				if err != nil || !strings.Contains(string(data), `"intent"`) {
					t.Fatal("remote creation began without durable intent", err)
				}
				switch failure {
				case "network":
					return nil, temporaryNetwork{}
				case "server":
					return responseJSON(r, 503, map[string]any{"error": map[string]any{"code": 503, "message": "secret"}}), nil
				case "malformed":
					return responseJSON(r, 200, &dataportability.InitiatePortabilityArchiveResponse{ArchiveJobId: "job-123"}), nil
				default:
					cancel()
					return nil, context.Canceled
				}
			})
			_, err := m.CreateExport(ctx)
			if err == nil || IsTemporary(err) || posts != 1 {
				t.Fatal("creation not held", err, posts)
			}
			m2 := testManager(t, Options{Provider: "portability", Resources: m.opts.Resources, DataDir: m.opts.DataDir}, func(r *http.Request) (*http.Response, error) {
				posts++
				t.Fatal("duplicate mutation after restart")
				return nil, nil
			})
			if _, err = m2.CreateExport(context.Background()); err == nil || !strings.Contains(err.Error(), "unknown outcome") || posts != 1 {
				t.Fatal("no persisted uncertain hold", err, posts)
			}
		})
	}
}

func TestPortabilityRejectedCreationIsTerminalAndExplicitRetryAllowed(t *testing.T) {
	posts := 0
	m := testManager(t, Options{Provider: "portability", Resources: []string{"myactivity.search"}}, func(r *http.Request) (*http.Response, error) {
		if r.URL.Path == "/v1/accessType:check" {
			return responseJSON(r, 200, &dataportability.CheckAccessTypeResponse{TimeBasedResources: []string{"myactivity.search"}}), nil
		}
		posts++
		resp := responseJSON(r, 429, map[string]any{"error": map[string]any{"code": 429, "message": "secret"}})
		resp.Header.Set("Retry-After", "3600")
		return resp, nil
	})
	_, err := m.CreateExport(context.Background())
	if err == nil || IsTemporary(err) || posts != 1 || strings.Contains(err.Error(), "secret") {
		t.Fatal(err, posts)
	}
	state, _, unlock, err := m.openState(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	unlock()
	if state.Intent != nil || len(state.Jobs) != 0 {
		t.Fatal("definitive rejection retained unknown creation")
	}
	_, err = m.CreateExport(context.Background())
	if err == nil || posts != 2 {
		t.Fatal("explicit retry prevented", err, posts)
	}
}

func TestPortabilityRefreshRejectsChangesAndLifetime(t *testing.T) {
	for _, change := range []string{"order", "count", "size", "name", "generation", "object", "expired", "scope"} {
		t.Run(change, func(t *testing.T) {
			now := time.Date(2026, 10, 25, 10, 0, 0, 0, time.UTC)
			urls := []string{"https://storage.googleapis.com/bucket/a.zip?sig=secret", "https://storage.googleapis.com/bucket/b.zip?sig=secret"}
			job := testJob(now.Add(-time.Hour))
			job.State = "COMPLETE"
			for _, raw := range urls {
				u, _ := validSignedURL(raw)
				job.Parts = append(job.Parts, archiveMetadata{Object: objectIdentity(u), Filename: strings.TrimPrefix(u.Path, "/bucket/"), Size: 100, Generation: "1234"})
			}
			switch change {
			case "order":
				urls[0], urls[1] = urls[1], urls[0]
			case "count":
				urls = urls[:1]
			case "object":
				urls[0] = "https://storage.googleapis.com/bucket/new.zip?sig=secret"
			case "expired":
				job.CreatedAt = now.Add(-portabilityDataLife)
			case "scope":
				job.Resources = []string{"youtube.private_videos"}
			}
			polls := 0
			m := testManager(t, Options{Provider: "portability", Resources: []string{"myactivity.search"}}, func(r *http.Request) (*http.Response, error) {
				polls++
				return responseJSON(r, 200, &dataportability.PortabilityArchiveState{State: "COMPLETE", Urls: urls}), nil
			})
			m.now = func() time.Time { return now }
			seedJobs(t, m, job)
			installSigned(t, m, func(r *http.Request) (*http.Response, error) {
				name := strings.TrimPrefix(r.URL.Path, "/bucket/")
				size := int64(100)
				if change == "size" {
					size++
				}
				if change == "name" {
					name = "changed.zip"
				}
				resp := signedResponse(r, name, size)
				if change == "generation" {
					resp.Header.Set("X-Goog-Generation", "1235")
				}
				return resp, nil
			})
			if link, err := m.RefreshLink(context.Background(), "portability-job-123", 1); err == nil || link != "" {
				t.Fatal("changed metadata got link", link, err)
			}
			if (change == "expired" || change == "scope") && polls != 0 {
				t.Fatal("expired or out-of-scope API request")
			}
		})
	}
}

func TestPortabilityCadenceOnFailuresLockAndCorruptState(t *testing.T) {
	now := time.Date(2026, 10, 25, 10, 0, 0, 0, time.UTC)
	polls := 0
	m := testManager(t, Options{Provider: "portability", Resources: []string{"myactivity.search"}}, func(r *http.Request) (*http.Response, error) {
		polls++
		resp := responseJSON(r, 503, map[string]any{"error": map[string]any{"code": 503, "message": "secret"}})
		resp.Header.Set("Retry-After", "7200")
		return resp, nil
	})
	m.now = func() time.Time { return now }
	file := seedJobs(t, m, testJob(now.Add(-time.Hour)))
	_, err := m.List(context.Background())
	if !IsTemporary(err) || RetryDelay(err) != 2*time.Hour || polls != 1 {
		t.Fatal(err, polls)
	}
	if _, err = m.List(context.Background()); err != nil || polls != 1 {
		t.Fatal("failed poll bypassed cooldown", err, polls)
	}
	now = now.Add(6 * time.Minute)
	m2 := testManager(t, Options{DataDir: m.opts.DataDir, Provider: "portability", Resources: m.opts.Resources}, func(r *http.Request) (*http.Response, error) {
		polls++
		t.Fatal("restart polled before the persisted two-hour Retry-After")
		return nil, nil
	})
	m2.now = func() time.Time { return now }
	if _, err = m2.List(context.Background()); err != nil || polls != 1 {
		t.Fatal("restart lost provider-requested backoff", err, polls)
	}
	_, _, unlock, err := m.openState(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if _, err = m.List(context.Background()); err == nil {
		t.Fatal("concurrent state reader did not lock")
	}
	unlock()
	if err = os.WriteFile(file, []byte(`{"version":1,"jobs":[]} trailing`), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err = m.CreateExport(context.Background()); err == nil {
		t.Fatal("corrupt state silently reset")
	}
}

func TestPortabilitySelectedResourcesAndOneTimeGrant(t *testing.T) {
	now := time.Date(2026, 10, 25, 10, 0, 0, 0, time.UTC)
	m := testManager(t, Options{Provider: "portability", Resources: []string{"myactivity.search"}}, func(r *http.Request) (*http.Response, error) { t.Fatal("unexpected network"); return nil, nil })
	m.now = func() time.Time { return now }
	job := testJob(now.Add(-25 * time.Hour))
	job.AccessType = "ACCESS_TYPE_ONE_TIME"
	seedJobs(t, m, job)
	if _, err := m.CreateExport(context.Background()); err == nil || !strings.Contains(err.Error(), "one-time") {
		t.Fatal(err)
	}
	job.Resources = []string{"youtube.private_videos"}
	seedJobs(t, m, job)
	if exports, err := m.List(context.Background()); err != nil || len(exports) != 0 {
		t.Fatal("out of scope job was polled", exports, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := m.CreateExport(ctx); !errors.Is(err, context.Canceled) {
		t.Fatal("cancellation ignored", err)
	}
}

func TestPortabilityProfileIsolationAndSameProfileResourceHistory(t *testing.T) {
	dir := t.TempDir()
	calls := 0
	handler := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		t.Fatal("unexpected remote operation")
		return nil, nil
	})
	a := testManager(t, Options{DataDir: dir, Provider: "portability", Resources: []string{"myactivity.search"}}, handler)
	state, name, unlock, err := a.openState(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	state.Intent = &creationIntent{Resources: a.opts.Resources, CreatedAt: time.Now().UTC()}
	if err = saveState(name, state); err != nil {
		t.Fatal(err)
	}
	unlock()
	// Same client/profile, different resource scopes must still see the hold.
	same := testManager(t, Options{DataDir: dir, Provider: "portability", Resources: []string{"youtube.private_videos"}}, handler)
	if _, err = same.CreateExport(context.Background()); err == nil || !strings.Contains(err.Error(), "unknown outcome") {
		t.Fatal("resource changes bypassed creation intent", err)
	}
	// A different account hint/client profile in the same DataDir is isolated.
	b := testManager(t, Options{DataDir: dir, Provider: "portability", Resources: []string{"myactivity.search"}, Session: &fakeSession{profile: strings.Repeat("b", 64), client: &http.Client{Transport: handler}}}, handler)
	other, otherName, unlock, err := b.openState(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	unlock()
	if other.Intent != nil || len(other.Jobs) != 0 || otherName == name {
		t.Fatal("different profile inherited another profile's job state")
	}
	if exports, err := b.List(context.Background()); err != nil || len(exports) != 0 || calls != 0 {
		t.Fatal("other profile was polled", exports, err)
	}
	// One-time history is shared when this profile returns to the same resources.
	job := testJob(time.Now().Add(-48 * time.Hour))
	job.AccessType = "ACCESS_TYPE_ONE_TIME"
	seedJobs(t, a, job)
	if _, err = a.CreateExport(context.Background()); err == nil {
		t.Fatal("persisted intent unexpectedly cleared")
	}
	// Invalid IDs must fail before any path is created outside the profile root.
	bad := testManager(t, Options{DataDir: dir, Provider: "portability", Resources: []string{"myactivity.search"}, Session: &fakeSession{profile: "../escape", client: &http.Client{Transport: handler}}}, handler)
	if _, err = bad.List(context.Background()); err == nil {
		t.Fatal("unvalidated profile ID used as path")
	}
	if _, err = os.Stat(filepath.Join(dir, "escape")); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("invalid profile wrote outside source state")
	}
}

func TestPortabilityReceiptFailuresHoldBeforeAndAfterPOST(t *testing.T) {
	for _, failure := range []string{"intent", "receipt"} {
		t.Run(failure, func(t *testing.T) {
			posts := 0
			m := testManager(t, Options{Provider: "portability", Resources: []string{"myactivity.search"}}, func(r *http.Request) (*http.Response, error) {
				if r.URL.Path == "/v1/accessType:check" {
					return responseJSON(r, 200, &dataportability.CheckAccessTypeResponse{TimeBasedResources: []string{"myactivity.search"}}), nil
				}
				posts++
				return responseJSON(r, 200, &dataportability.InitiatePortabilityArchiveResponse{ArchiveJobId: "job-123", AccessType: "ACCESS_TYPE_TIME_BASED"}), nil
			})
			writes := 0
			m.persist = func(name string, state *portabilityState) error {
				writes++
				if failure == "intent" || writes == 2 {
					return errors.New("injected durable receipt failure")
				}
				return saveState(name, state)
			}
			if _, err := m.CreateExport(context.Background()); err == nil || IsTemporary(err) {
				t.Fatal("durable receipt failure was acknowledged or retryable", err)
			}
			if failure == "intent" {
				if posts != 0 {
					t.Fatal("POST before durable creation intent")
				}
				return
			}
			if posts != 1 {
				t.Fatal("wrong creation count", posts)
			}
			m2 := testManager(t, Options{DataDir: m.opts.DataDir, Provider: "portability", Resources: m.opts.Resources}, func(r *http.Request) (*http.Response, error) {
				posts++
				t.Fatal("duplicate POST after failed receipt")
				return nil, nil
			})
			if _, err := m2.CreateExport(context.Background()); err == nil || posts != 1 {
				t.Fatal("failed receipt lost creation hold", err)
			}
		})
	}
}
