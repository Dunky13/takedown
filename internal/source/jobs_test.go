package source

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Dunky13/takedown/internal/model"
	downloadstate "github.com/Dunky13/takedown/internal/state"

	"google.golang.org/api/dataportability/v1"
	"google.golang.org/api/drive/v3"
)

func TestTerminalJobDoesNotBlockReadySibling(t *testing.T) {
	for _, terminal := range []string{"FAILED", "CANCELLED"} {
		for _, persisted := range []bool{true, false} {
			t.Run(terminal+"/persisted="+map[bool]string{true: "yes", false: "no"}[persisted], func(t *testing.T) {
				now := time.Date(2026, 10, 25, 12, 0, 0, 0, time.UTC)
				m := testManager(t, Options{Provider: "portability", Resources: []string{"myactivity.search"}}, func(r *http.Request) (*http.Response, error) {
					if strings.Contains(r.URL.Path, "broken") {
						return responseJSON(r, 200, &dataportability.PortabilityArchiveState{State: terminal}), nil
					}
					return responseJSON(r, 200, &dataportability.PortabilityArchiveState{State: "COMPLETE", Urls: []string{"https://storage.googleapis.com/bucket/ready.zip?signature=secret"}}), nil
				})
				m.now = func() time.Time { return now }
				installSigned(t, m, func(r *http.Request) (*http.Response, error) { return signedResponse(r, "ready.zip", 100), nil })
				broken := testJob(now.Add(-time.Hour))
				broken.ID = "broken"
				if persisted {
					broken.State = terminal
				}
				ready := testJob(now.Add(-time.Hour))
				ready.ID = "ready"
				seedJobs(t, m, broken, ready)
				exports, err := m.List(context.Background())
				if err != nil || len(exports) != 1 || exports[0].ID != "portability-ready" {
					t.Fatalf("ready sibling blocked: exports=%v err=%v", exports, err)
				}
			})
		}
	}
}

func TestDriveDuplicateNamesAcrossExports(t *testing.T) {
	m := testManager(t, Options{DriveFolderID: "selected"}, func(r *http.Request) (*http.Response, error) {
		if r.URL.Path == "/drive/v3/files/selected" {
			return responseJSON(r, 200, folderFile("selected", "")), nil
		}
		switch r.URL.Query().Get("q") {
		case "'selected' in parents and trashed = false":
			return responseJSON(r, 200, &drive.FileList{Files: []*drive.File{folderFile("first", "selected"), folderFile("second", "selected")}}), nil
		case "'first' in parents and trashed = false":
			return responseJSON(r, 200, &drive.FileList{Files: []*drive.File{archiveFile("a", "takeout.zip", "first")}}), nil
		case "'second' in parents and trashed = false":
			return responseJSON(r, 200, &drive.FileList{Files: []*drive.File{archiveFile("b", "takeout.zip", "second")}}), nil
		}
		t.Fatalf("unexpected Drive request: %s", r.URL.Path)
		return nil, nil
	})
	exports, err := m.List(context.Background())
	if err != nil || len(exports) != 2 || exports[0].ID == exports[1].ID {
		t.Fatalf("distinct Drive exports rejected: %v %v", exports, err)
	}
}

func editJournal(t *testing.T, m *Manager, edit func(*portabilityState)) {
	t.Helper()
	s, name, unlock, err := m.openState(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	defer unlock()
	edit(s)
	if err = saveState(name, s); err != nil {
		t.Fatal(err)
	}
}
func readJournal(t *testing.T, m *Manager) *portabilityState {
	t.Helper()
	data, err := os.ReadFile(statePath(t, m))
	if err != nil {
		t.Fatal(err)
	}
	var s portabilityState
	if err = json.Unmarshal(data, &s); err != nil {
		t.Fatal(err)
	}
	if err = validateState(&s); err != nil {
		t.Fatal(err)
	}
	return &s
}

func TestJobsLocalAndRefreshBoundedCadence(t *testing.T) {
	now := time.Date(2026, 10, 25, 12, 0, 0, 0, time.UTC)
	calls := 0
	m := testManager(t, Options{Provider: "portability", Resources: []string{"myactivity.search"}}, func(r *http.Request) (*http.Response, error) {
		calls++
		if r.Method != http.MethodGet {
			t.Fatal("refresh mutated Google")
		}
		return responseJSON(r, 200, &dataportability.PortabilityArchiveState{State: "IN_PROGRESS"}), nil
	})
	m.now = func() time.Time { return now }
	jobs := make([]portabilityJob, 101)
	for i := range jobs {
		jobs[i] = testJob(now.Add(-time.Hour))
		jobs[i].ID = fmt.Sprintf("job-%d", i)
	}
	seedJobs(t, m, jobs...)
	snap, err := m.Jobs(context.Background())
	if err != nil || len(snap.Jobs) != 101 || calls != 0 {
		t.Fatal("local status used network", err, calls)
	}
	snap.Jobs[0].Resources[0] = "chrome.history"
	if readJournal(t, m).Jobs[0].Resources[0] != "myactivity.search" {
		t.Fatal("mutable snapshot changed receipt")
	}
	if _, err = m.RefreshJobs(context.Background()); err != nil || calls != 100 {
		t.Fatal("refresh not bounded", err, calls)
	}
	if _, err = m.RefreshJobs(context.Background()); err != nil || calls != 101 {
		t.Fatal("refresh did not advance to remaining due job", err, calls)
	}
	m2 := testManager(t, m.opts, func(r *http.Request) (*http.Response, error) { t.Fatal("unexpected transport"); return nil, nil })
	m2.now = m.now
	if _, err = m2.RefreshJobs(context.Background()); err != nil || calls != 101 {
		t.Fatal("restart bypassed cadence", err, calls)
	}
	now = now.Add(5 * time.Minute)
	if _, err = m.RefreshJobs(context.Background()); err != nil || calls != 201 {
		t.Fatal("due refresh failed", err, calls)
	}
}

func TestJobScopedMetadataFailurePreservesReadyAndSafetyErrors(t *testing.T) {
	for _, failure := range []string{"missing", "ambiguous", "unsafe", "auth", "quota"} {
		t.Run(failure, func(t *testing.T) {
			now := time.Now().UTC()
			calls := 0
			m := testManager(t, Options{Provider: "portability", Resources: []string{"myactivity.search"}}, func(r *http.Request) (*http.Response, error) {
				calls++
				if strings.Contains(r.URL.Path, "broken") {
					switch failure {
					case "missing":
						return responseJSON(r, 404, map[string]any{"error": map[string]any{"code": 404, "message": "secret"}}), nil
					case "auth":
						return responseJSON(r, 403, map[string]any{"error": map[string]any{"code": 403, "message": "secret"}}), nil
					case "quota":
						return responseJSON(r, 429, map[string]any{"error": map[string]any{"code": 429, "message": "secret"}}), nil
					case "unsafe":
						return responseJSON(r, 200, &dataportability.PortabilityArchiveState{State: "COMPLETE", Urls: []string{"https://evil.example/archive.zip?token=secret"}}), nil
					case "ambiguous":
						return responseJSON(r, 200, &dataportability.PortabilityArchiveState{State: "COMPLETE", Urls: []string{"https://storage.googleapis.com/bucket/same.zip?a=1", "https://storage.googleapis.com/bucket/same.zip?a=2"}}), nil
					}
				}
				return responseJSON(r, 200, &dataportability.PortabilityArchiveState{State: "COMPLETE", Urls: []string{"https://storage.googleapis.com/bucket/ready.zip?signature=secret"}}), nil
			})
			m.now = func() time.Time { return now }
			installSigned(t, m, func(r *http.Request) (*http.Response, error) { return signedResponse(r, "archive.zip", 100), nil })
			broken := testJob(now.Add(-time.Hour))
			broken.ID = "broken"
			ready := broken
			ready.ID = "ready"
			seedJobs(t, m, broken, ready)
			exports, err := m.List(context.Background())
			if failure == "auth" || failure == "quota" {
				if err == nil || exports != nil || calls != 1 {
					t.Fatal("ignored safety error", err, calls)
				}
				return
			}
			if err != nil || len(exports) != 1 || exports[0].ID != "portability-ready" {
				t.Fatal("scoped failure blocked ready job", err, exports)
			}
			snap, err := m.Jobs(context.Background())
			if err != nil || !strings.Contains(snap.Jobs[0].Message, "held") {
				t.Fatal("scoped failure invisible", err, snap)
			}
		})
	}
}

func TestCheckAccessSingleAttemptAndHolds(t *testing.T) {
	calls := 0
	status := 200
	m := testManager(t, Options{Provider: "portability", Resources: []string{"myactivity.search"}}, func(r *http.Request) (*http.Response, error) {
		calls++
		if r.Method != http.MethodPost || r.URL.Path != "/v1/accessType:check" {
			t.Fatal("wrong grant endpoint")
		}
		if status != 200 {
			return responseJSON(r, status, map[string]any{"error": map[string]any{"code": status, "message": "https://secret.example/?token=secret"}}), nil
		}
		return responseJSON(r, 200, &dataportability.CheckAccessTypeResponse{OneTimeResources: []string{"myactivity.search"}}), nil
	})
	j := testJob(time.Now().Add(-48 * time.Hour))
	j.AccessType = "ACCESS_TYPE_ONE_TIME"
	j.State = "COMPLETE"
	seedJobs(t, m, j)
	a, err := m.CheckAccess(context.Background())
	if err != nil || len(a.OneTimeResources) != 1 || calls != 1 {
		t.Fatal(a, err, calls)
	}
	if _, err = m.CreateExport(context.Background()); err == nil || calls != 1 {
		t.Fatal("grant check cleared one-time history", err, calls)
	}
	status = 503
	if _, err = m.CheckAccess(context.Background()); err == nil || strings.Contains(err.Error(), "secret.example") || calls != 2 {
		t.Fatal("POST replay or leaked response", err, calls)
	}
}

func TestRetryJobLineageAndUnknownOutcome(t *testing.T) {
	for _, unknown := range []bool{false, true} {
		t.Run(fmt.Sprint(unknown), func(t *testing.T) {
			now := time.Now().UTC()
			posts := 0
			var m *Manager
			m = testManager(t, Options{Provider: "portability", Resources: []string{"myactivity.search"}}, func(r *http.Request) (*http.Response, error) {
				posts++
				if r.Method != http.MethodPost || !strings.HasSuffix(r.URL.Path, ":retry") {
					t.Fatal("wrong retry endpoint")
				}
				s := readJournal(t, m)
				if s.Intent == nil || s.Intent.Kind != "retry" || s.Jobs[0].RetryCount != posts {
					t.Fatal("retry lacked durable intent/budget", s)
				}
				if unknown {
					return nil, temporaryNetwork{}
				}
				return responseJSON(r, 200, &dataportability.RetryPortabilityArchiveResponse{ArchiveJobId: fmt.Sprintf("retry-%d", posts)}), nil
			})
			m.now = func() time.Time { return now }
			j := testJob(now.Add(-time.Hour))
			j.State = "FAILED"
			seedJobs(t, m, j)
			id, err := m.RetryJob(context.Background(), j.ID)
			if unknown {
				if err == nil || IsTemporary(err) || posts != 1 {
					t.Fatal("unknown retry not held", err, posts)
				}
				m2 := testManager(t, m.opts, nil)
				if _, err = m2.RetryJob(context.Background(), j.ID); err == nil || posts != 1 {
					t.Fatal("unknown mutation replay", err, posts)
				}
				if s := readJournal(t, m); s.Intent == nil || s.Intent.Kind != "retry" {
					t.Fatal("lost retry intent")
				}
				return
			}
			if err != nil || id != "retry-1" {
				t.Fatal(id, err)
			}
			if _, err = m.RetryJob(context.Background(), j.ID); err == nil || posts != 1 {
				t.Fatal("predecessor replayed", err, posts)
			}
			for retry := 2; retry <= 3; retry++ {
				editJournal(t, m, func(s *portabilityState) { s.Jobs[len(s.Jobs)-1].State = "FAILED" })
				id, err = m.RetryJob(context.Background(), id)
				if err != nil {
					t.Fatal(err)
				}
			}
			editJournal(t, m, func(s *portabilityState) { s.Jobs[len(s.Jobs)-1].State = "FAILED" })
			if _, err = m.RetryJob(context.Background(), id); err == nil || posts != 3 {
				t.Fatal("retry cap bypass", err, posts)
			}
			if err = m.ArchiveJob(context.Background(), j.ID); err != nil {
				t.Fatal(err)
			}
			if _, err = m.RetryJob(context.Background(), id); err == nil || posts != 3 {
				t.Fatal("archival erased retry budget", err, posts)
			}
		})
	}
}

func TestCancelRestrictionsAndExplicitReconciliation(t *testing.T) {
	now := time.Now().UTC()
	posts, polls := 0, 0
	remote := "IN_PROGRESS"
	var m *Manager
	m = testManager(t, Options{Provider: "portability", Resources: []string{"myactivity.search"}}, func(r *http.Request) (*http.Response, error) {
		if r.Method == http.MethodPost {
			posts++
			s := readJournal(t, m)
			if s.Intent == nil || s.Intent.Kind != "cancel" {
				t.Fatal("cancel without intent")
			}
			return nil, temporaryNetwork{}
		}
		polls++
		return responseJSON(r, 200, &dataportability.PortabilityArchiveState{State: remote, Name: "archiveJobs/job-123/portabilityArchiveState"}), nil
	})
	m.now = func() time.Time { return now }
	for _, j := range []portabilityJob{func() portabilityJob {
		j := testJob(now.Add(-time.Hour))
		j.AccessType = "ACCESS_TYPE_ONE_TIME"
		return j
	}(), func() portabilityJob { j := testJob(now.Add(-time.Hour)); j.State = "COMPLETE"; return j }()} {
		seedJobs(t, m, j)
		if err := m.CancelJob(context.Background(), j.ID); err == nil || posts != 0 {
			t.Fatal("invalid cancellation", err, posts)
		}
	}
	seedJobs(t, m, testJob(now.Add(-time.Hour)))
	if err := m.CancelJob(context.Background(), "job-123"); err == nil || IsTemporary(err) || posts != 1 {
		t.Fatal(err, posts)
	}
	if err := m.ReconcileIntent(context.Background(), "other"); err == nil || polls != 0 {
		t.Fatal("wrong cancellation identity")
	}
	if err := m.ReconcileIntent(context.Background(), "job-123"); err == nil || polls != 1 {
		t.Fatal("in-progress outcome cleared", err, polls)
	}
	if err := m.CancelJob(context.Background(), "job-123"); err == nil || posts != 1 {
		t.Fatal("unknown cancel replay")
	}
	now = now.Add(5 * time.Minute)
	remote = "CANCELLED"
	// Automated refresh records terminal status but deliberately retains intent.
	if _, err := m.RefreshJobs(context.Background()); err != nil {
		t.Fatal(err)
	}
	if s := readJournal(t, m); s.Intent == nil {
		t.Fatal("watch resolved mutation")
	}
	if err := m.ReconcileIntent(context.Background(), "job-123"); err != nil {
		t.Fatal(err)
	}
	if s := readJournal(t, m); s.Intent != nil || s.Jobs[0].State != "CANCELLED" {
		t.Fatal("cancel not reconciled")
	}
}

func TestKnownJobAdoptionUsesDurableIntentAndConservativeGrant(t *testing.T) {
	now := time.Now().UTC()
	created := now.Add(-time.Hour)
	calls := 0
	valid := false
	m := testManager(t, Options{Provider: "portability", Resources: []string{"myactivity.search"}}, func(r *http.Request) (*http.Response, error) {
		calls++
		if r.URL.Path == "/v1/accessType:check" {
			return responseJSON(r, 200, &dataportability.CheckAccessTypeResponse{TimeBasedResources: []string{"myactivity.search"}}), nil
		}
		name := "wrong"
		if valid {
			name = "archiveJobs/evidence-id/portabilityArchiveState"
		}
		return responseJSON(r, 200, &dataportability.PortabilityArchiveState{State: "IN_PROGRESS", Name: name}), nil
	})
	m.now = func() time.Time { return now }
	editJournal(t, m, func(s *portabilityState) {
		s.Intent = &creationIntent{Resources: []string{"myactivity.search"}, CreatedAt: created}
	})
	if err := m.ReconcileIntent(context.Background(), "invented/path"); err == nil || calls != 0 {
		t.Fatal("invalid ID accepted")
	}
	if err := m.ReconcileIntent(context.Background(), "evidence-id"); err == nil || calls != 2 {
		t.Fatal("wrong identity accepted", err, calls)
	}
	if s := readJournal(t, m); s.Intent == nil || len(s.Jobs) != 0 {
		t.Fatal("invalid response cleared hold")
	}
	valid = true
	if err := m.ReconcileIntent(context.Background(), "evidence-id"); err == nil || calls != 2 {
		t.Fatal("adoption cadence bypass")
	}
	now = now.Add(5 * time.Minute)
	if err := m.ReconcileIntent(context.Background(), "evidence-id"); err != nil {
		t.Fatal(err)
	}
	s := readJournal(t, m)
	if s.Intent != nil || len(s.Jobs) != 1 || s.Jobs[0].ID != "evidence-id" || !s.Jobs[0].CreatedAt.Equal(created) || s.Jobs[0].Resources[0] != "myactivity.search" || s.Jobs[0].AccessType != "ACCESS_TYPE_ONE_TIME" {
		t.Fatal("adoption invented authority", s)
	}
}

func TestExpiredIntentRequiresPolicyAndRetainsTombstone(t *testing.T) {
	for _, grant := range []string{"time", "one", "missing"} {
		t.Run(grant, func(t *testing.T) {
			now := time.Now().UTC()
			calls := 0
			m := testManager(t, Options{Provider: "portability", Resources: []string{"myactivity.search"}}, func(r *http.Request) (*http.Response, error) {
				calls++
				if r.URL.Path != "/v1/accessType:check" {
					t.Fatal("expiry mutated Google")
				}
				a := &dataportability.CheckAccessTypeResponse{}
				if grant == "time" {
					a.TimeBasedResources = []string{"myactivity.search"}
				}
				if grant == "one" {
					a.OneTimeResources = []string{"myactivity.search"}
				}
				return responseJSON(r, 200, a), nil
			})
			m.now = func() time.Time { return now }
			editJournal(t, m, func(s *portabilityState) {
				s.Intent = &creationIntent{Resources: []string{"myactivity.search"}, CreatedAt: now.Add(-13 * 24 * time.Hour), AccessType: "ACCESS_TYPE_TIME_BASED"}
			})
			if err := m.ResolveExpiredIntent(context.Background()); err == nil || calls != 0 {
				t.Fatal("early expiry")
			}
			now = now.Add(24*time.Hour + 5*time.Minute)
			err := m.ResolveExpiredIntent(context.Background())
			s := readJournal(t, m)
			if grant != "time" {
				if err == nil || s.Intent == nil || len(s.ExpiredIntents) != 0 {
					t.Fatal("unsafe expiry", err)
				}
				return
			}
			if err != nil || s.Intent != nil || len(s.ExpiredIntents) != 1 {
				t.Fatal("no executable timed expiry", err, s)
			}
			if err = m.ResetAuthorization(WithResetConfirmation(context.Background())); err == nil {
				t.Fatal("unknown tombstone bypassed reset")
			}
		})
	}
}

func completedJob(now time.Time) portabilityJob {
	j := testJob(now.Add(-25 * time.Hour))
	j.State = "COMPLETE"
	j.AccessType = "ACCESS_TYPE_ONE_TIME"
	j.Parts = []archiveMetadata{{Object: strings.Repeat("a", 64), Filename: "archive.zip", Size: 100, MD5: "0123456789abcdef0123456789abcdef"}}
	return j
}
func completeReceipt(t *testing.T, m *Manager, j portabilityJob, alter func(*downloadstate.Record)) {
	t.Helper()
	p := jobExport(j, []string{""}).Parts[0]
	p.RelativePath = model.NamespacePath(p.ExportID, p.Filename)
	r := downloadstate.Record{Part: p, Destination: "local-archive", Offset: p.Size, Complete: true, SHA256: strings.Repeat("a", 64)}
	if alter != nil {
		alter(&r)
	}
	s, err := downloadstate.Open(filepath.Join(m.opts.DataDir, "takedown-state.json"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if err = s.Put(r); err != nil {
		t.Fatal(err)
	}
}

func TestResetRequiresExactDurableCompletionsAndConfirmation(t *testing.T) {
	for _, kind := range []string{"no-confirm", "pending", "failed", "cancelled", "expired-unknown", "intent", "missing-record", "missing-digest", "wrong-filename", "wrong-size", "wrong-md5", "wrong-index", "wrong-id", "partial", "ok", "archived-ok", "expired-complete-ok"} {
		t.Run(kind, func(t *testing.T) {
			now := time.Now().UTC()
			posts := 0
			var m *Manager
			m = testManager(t, Options{Provider: "portability", Resources: []string{"myactivity.search"}}, func(r *http.Request) (*http.Response, error) {
				posts++
				if r.Method != http.MethodPost || r.URL.Path != "/v1/authorization:reset" {
					t.Fatal("wrong reset endpoint")
				}
				s := readJournal(t, m)
				if s.Intent == nil || s.Intent.Kind != "reset" || s.Generation != 0 {
					t.Fatal("reset lacked pre-mutation receipt")
				}
				return responseJSON(r, 200, map[string]any{}), nil
			})
			m.now = func() time.Time { return now }
			j := completedJob(now)
			switch kind {
			case "pending":
				j.State = "IN_PROGRESS"
			case "failed":
				j.State = "FAILED"
			case "cancelled":
				j.State = "CANCELLED"
			case "expired-unknown":
				j.State = "EXPIRED"
			case "archived-ok":
				j.Archived = true
			case "expired-complete-ok":
				j.CreatedAt = now.Add(-15 * 24 * time.Hour)
			}
			seedJobs(t, m, j)
			if kind == "intent" {
				editJournal(t, m, func(s *portabilityState) { s.Intent = &creationIntent{Resources: j.Resources, CreatedAt: now} })
			}
			if kind != "missing-record" {
				completeReceipt(t, m, j, func(r *downloadstate.Record) {
					switch kind {
					case "missing-digest":
						r.SHA256 = ""
					case "wrong-filename":
						r.Part.Filename = "other.zip"
						r.Part.RelativePath = model.NamespacePath(r.Part.ExportID, r.Part.Filename)
					case "wrong-size":
						r.Part.Size++
						r.Offset++
					case "wrong-md5":
						r.Part.MD5 = strings.Repeat("b", 32)
					case "wrong-index":
						r.Part.Index++
					case "wrong-id":
						r.Part.ExportID = "other-export"
						r.Part.RelativePath = model.NamespacePath(r.Part.ExportID, r.Part.Filename)
					case "partial":
						r.Offset--
						r.Complete = false
						r.SHA256 = ""
					}
				})
			}
			ctx := context.Background()
			if kind != "no-confirm" {
				ctx = WithResetConfirmation(ctx)
			}
			err := m.ResetAuthorization(ctx)
			if strings.HasSuffix(kind, "ok") {
				if err != nil || posts != 1 {
					t.Fatal("safe reset rejected", err, posts)
				}
				s := readJournal(t, m)
				if s.Generation != 1 || !s.NeedsConsent || !s.Jobs[0].Archived || s.Intent != nil {
					t.Fatal("reset receipt missing", s)
				}
				return
			}
			if err == nil || posts != 0 {
				t.Fatal("unsafe reset accepted", err, posts)
			}
		})
	}
}

func TestMutationPersistenceFailuresPreventReplay(t *testing.T) {
	for _, kind := range []string{"retry", "cancel", "reset"} {
		for _, failAt := range []int{1, 2} {
			t.Run(fmt.Sprintf("%s/write-%d", kind, failAt), func(t *testing.T) {
				now := time.Now().UTC()
				posts := 0
				m := testManager(t, Options{Provider: "portability", Resources: []string{"myactivity.search"}}, func(r *http.Request) (*http.Response, error) {
					posts++
					if kind == "retry" {
						return responseJSON(r, 200, &dataportability.RetryPortabilityArchiveResponse{ArchiveJobId: "successor"}), nil
					}
					return responseJSON(r, 200, map[string]any{}), nil
				})
				m.now = func() time.Time { return now }
				j := testJob(now.Add(-time.Hour))
				if kind == "retry" {
					j.State = "FAILED"
				}
				if kind == "reset" {
					j = completedJob(now)
				}
				seedJobs(t, m, j)
				if kind == "reset" {
					completeReceipt(t, m, j, nil)
				}
				writes := 0
				m.persist = func(name string, s *portabilityState) error {
					writes++
					if writes == failAt {
						return errors.New("injected durability failure")
					}
					return saveState(name, s)
				}
				act := func() error {
					switch kind {
					case "retry":
						_, err := m.RetryJob(context.Background(), j.ID)
						return err
					case "cancel":
						return m.CancelJob(context.Background(), j.ID)
					default:
						return m.ResetAuthorization(WithResetConfirmation(context.Background()))
					}
				}
				if err := act(); err == nil || posts != failAt-1 {
					t.Fatal("wrong mutation boundary", err, posts)
				}
				if failAt == 1 {
					return
				}
				m.persist = saveState
				if err := act(); err == nil || posts != 1 {
					t.Fatal("lost response receipt permitted mutation replay", err, posts)
				}
				s := readJournal(t, m)
				if s.Intent == nil || intentKind(*s.Intent) != kind {
					t.Fatal("durable hold lost", s)
				}
			})
		}
	}
}

func TestResetUnknownAndFreshConsentCannotClearHolds(t *testing.T) {
	for _, unknown := range []bool{false, true} {
		t.Run(fmt.Sprint(unknown), func(t *testing.T) {
			now := time.Now().UTC()
			resets, creates, checks := 0, 0, 0
			fresh := false
			m := testManager(t, Options{Provider: "portability", Resources: []string{"myactivity.search"}}, func(r *http.Request) (*http.Response, error) {
				switch r.URL.Path {
				case "/v1/authorization:reset":
					resets++
					if unknown {
						return nil, temporaryNetwork{}
					}
					return responseJSON(r, 200, map[string]any{}), nil
				case "/v1/accessType:check":
					checks++
					if !fresh {
						return responseJSON(r, 401, map[string]any{"error": map[string]any{"code": 401}}), nil
					}
					return responseJSON(r, 200, &dataportability.CheckAccessTypeResponse{OneTimeResources: []string{"myactivity.search"}}), nil
				case "/v1/portabilityArchive:initiate":
					creates++
					return responseJSON(r, 200, &dataportability.InitiatePortabilityArchiveResponse{ArchiveJobId: "fresh-job", AccessType: "ACCESS_TYPE_ONE_TIME"}), nil
				}
				t.Fatal("unexpected endpoint")
				return nil, nil
			})
			m.now = func() time.Time { return now }
			j := completedJob(now)
			seedJobs(t, m, j)
			completeReceipt(t, m, j, nil)
			err := m.ResetAuthorization(WithResetConfirmation(context.Background()))
			if unknown {
				if err == nil || IsTemporary(err) {
					t.Fatal("unknown reset exposed retry")
				}
				if err = m.Login(context.Background()); err != nil {
					t.Fatal(err)
				}
				if _, err = m.CreateExport(context.Background()); err == nil || creates != 0 || checks != 0 {
					t.Fatal("login cleared reset hold", err)
				}
				if err = m.ReconcileIntent(context.Background(), j.ID); err == nil {
					t.Fatal("job ID proved unrelated reset")
				}
				now = now.Add(15 * 24 * time.Hour)
				if err = m.ResolveExpiredIntent(context.Background()); err == nil {
					t.Fatal("expiry cleared reset hold")
				}
				if err = m.ResetAuthorization(WithResetConfirmation(context.Background())); err == nil || resets != 1 {
					t.Fatal("reset replayed")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if _, err = m.CreateExport(context.Background()); err == nil || creates != 0 || checks != 1 {
				t.Fatal("old consent survived reset", err, creates, checks)
			}
			fresh = true
			if _, err = m.CreateExport(context.Background()); err != nil || creates != 1 || checks != 2 {
				t.Fatal("new grant workflow unavailable", err, creates, checks)
			}
			if _, err = m.CreateExport(context.Background()); err == nil || creates != 1 {
				t.Fatal("fresh grant exported twice")
			}
		})
	}
}

func TestArchivalReclaimsActiveLimitWithoutRemovingGuards(t *testing.T) {
	now := time.Now().UTC()
	posts := 0
	m := testManager(t, Options{Provider: "portability", Resources: []string{"myactivity.search"}}, func(r *http.Request) (*http.Response, error) {
		if r.URL.Path == "/v1/accessType:check" {
			return responseJSON(r, 200, &dataportability.CheckAccessTypeResponse{TimeBasedResources: []string{"myactivity.search"}}), nil
		}
		posts++
		return responseJSON(r, 200, &dataportability.InitiatePortabilityArchiveResponse{ArchiveJobId: "new-job", AccessType: "ACCESS_TYPE_TIME_BASED"}), nil
	})
	m.now = func() time.Time { return now }
	jobs := make([]portabilityJob, maxJobs)
	for i := range jobs {
		jobs[i] = testJob(now.Add(-48 * time.Hour))
		jobs[i].ID = fmt.Sprintf("historic-%d", i)
		jobs[i].State = "FAILED"
		jobs[i].Resources = []string{"chrome.history"}
	}
	seedJobs(t, m, jobs...)
	if _, err := m.CreateExport(context.Background()); err == nil || posts != 0 {
		t.Fatal("active limit bypass")
	}
	if err := m.ArchiveJob(context.Background(), jobs[0].ID); err != nil {
		t.Fatal(err)
	}
	if _, err := m.CreateExport(context.Background()); err != nil || posts != 1 {
		t.Fatal("archival did not free active slot", err, posts)
	}
	if s := readJournal(t, m); len(s.Jobs) != maxJobs+1 || activeJobs(s) != maxJobs {
		t.Fatal("history removed")
	}
	for _, oneTime := range []bool{false, true} {
		j := testJob(now.Add(-time.Hour))
		j.State = "COMPLETE"
		if oneTime {
			j.CreatedAt = now.Add(-48 * time.Hour)
			j.AccessType = "ACCESS_TYPE_ONE_TIME"
		}
		seedJobs(t, m, j)
		if err := m.ArchiveJob(context.Background(), j.ID); err != nil {
			t.Fatal(err)
		}
		if _, err := m.CreateExport(context.Background()); err == nil || posts != 1 {
			t.Fatal("archival erased cadence or one-time guard", err, posts)
		}
	}
}

func TestHistoricalUnknownOrOneTimeIntentCannotExpireFromCurrentGrant(t *testing.T) {
	for _, historical := range []string{"", "ACCESS_TYPE_ONE_TIME", "ACCESS_TYPE_TIME_BASED"} {
		t.Run(historical, func(t *testing.T) {
			now := time.Now().UTC()
			checks := 0
			m := testManager(t, Options{Provider: "portability", Resources: []string{"myactivity.search"}}, func(r *http.Request) (*http.Response, error) {
				checks++
				return responseJSON(r, 200, &dataportability.CheckAccessTypeResponse{TimeBasedResources: []string{"myactivity.search"}}), nil
			})
			m.now = func() time.Time { return now }
			editJournal(t, m, func(s *portabilityState) {
				s.Intent = &creationIntent{Kind: "create", Resources: []string{"myactivity.search"}, CreatedAt: now.Add(-portabilityDataLife), AccessType: historical}
			})
			if err := m.ResolveExpiredIntent(context.Background()); err == nil || checks != 0 {
				t.Fatal("request budget margin omitted", err, checks)
			}
			now = now.Add(5 * time.Minute)
			err := m.ResolveExpiredIntent(context.Background())
			if historical != "ACCESS_TYPE_TIME_BASED" {
				if err == nil || checks != 0 || readJournal(t, m).Intent == nil {
					t.Fatal("current grant guessed historical grant", err, checks)
				}
				return
			}
			if err != nil || checks != 1 || len(readJournal(t, m).ExpiredIntents) != 1 {
				t.Fatal("historically proven expiry unavailable", err, checks)
			}
		})
	}
}

func TestCreationPersistsVerifiedGrantBeforeSingleMutation(t *testing.T) {
	for _, grant := range []string{"ACCESS_TYPE_ONE_TIME", "ACCESS_TYPE_TIME_BASED"} {
		t.Run(grant, func(t *testing.T) {
			now := time.Now().UTC()
			checks, posts := 0, 0
			var m *Manager
			m = testManager(t, Options{Provider: "portability", Resources: []string{"myactivity.search"}}, func(r *http.Request) (*http.Response, error) {
				if r.URL.Path == "/v1/accessType:check" {
					checks++
					a := &dataportability.CheckAccessTypeResponse{}
					if grant == "ACCESS_TYPE_ONE_TIME" {
						a.OneTimeResources = []string{"myactivity.search"}
					} else {
						a.TimeBasedResources = []string{"myactivity.search"}
					}
					return responseJSON(r, 200, a), nil
				}
				posts++
				s := readJournal(t, m)
				if checks != 1 || s.Intent == nil || s.Intent.AccessType != grant || s.Intent.Resources[0] != "myactivity.search" {
					t.Fatal("mutation preceded historical grant receipt", s, checks)
				}
				return nil, temporaryNetwork{}
			})
			m.now = func() time.Time { return now }
			if _, err := m.CreateExport(context.Background()); err == nil || IsTemporary(err) || posts != 1 || checks != 1 {
				t.Fatal("unknown creation not held", err, posts, checks)
			}
			if _, err := m.CreateExport(context.Background()); err == nil || posts != 1 || checks != 1 {
				t.Fatal("mutation or preflight replayed past unknown hold", err, posts, checks)
			}
		})
	}
}

func TestCreationGrantPreflightCannotCreateWithMissingOrInvalidScopes(t *testing.T) {
	for _, kind := range []string{"missing", "conflict", "unknown", "unavailable"} {
		t.Run(kind, func(t *testing.T) {
			checks, posts := 0, 0
			m := testManager(t, Options{Provider: "portability", Resources: []string{"myactivity.search"}}, func(r *http.Request) (*http.Response, error) {
				if r.URL.Path != "/v1/accessType:check" {
					posts++
					t.Fatal("creation without verified selected grant")
				}
				checks++
				a := &dataportability.CheckAccessTypeResponse{}
				switch kind {
				case "missing":
					a.TimeBasedResources = []string{"chrome.history"}
				case "conflict":
					a.TimeBasedResources = []string{"myactivity.search"}
					a.OneTimeResources = []string{"myactivity.search"}
				case "unknown":
					a.TimeBasedResources = []string{"invented"}
				case "unavailable":
					return responseJSON(r, 503, map[string]any{"error": map[string]any{"code": 503, "message": "secret signed URL"}}), nil
				}
				return responseJSON(r, 200, a), nil
			})
			if _, err := m.CreateExport(context.Background()); err == nil || strings.Contains(err.Error(), "secret signed") || checks != 1 || posts != 0 {
				t.Fatal("unsafe preflight", err, checks, posts)
			}
			if _, err := os.Stat(statePath(t, m)); !errors.Is(err, os.ErrNotExist) {
				t.Fatal("failed preflight created mutation intent", err)
			}
		})
	}
}

func TestIsJobFailureExcludesSafetyAndAllowsLaterReadyJob(t *testing.T) {
	now := time.Now().UTC()
	phase := 0
	calls := 0
	m := testManager(t, Options{Provider: "portability", Resources: []string{"myactivity.search"}}, func(r *http.Request) (*http.Response, error) {
		calls++
		if strings.Contains(r.URL.Path, "broken") {
			return responseJSON(r, 404, map[string]any{"error": map[string]any{"code": 404}}), nil
		}
		if phase == 0 {
			return responseJSON(r, 200, &dataportability.PortabilityArchiveState{State: "IN_PROGRESS"}), nil
		}
		return responseJSON(r, 200, &dataportability.PortabilityArchiveState{State: "COMPLETE", Urls: []string{"https://storage.googleapis.com/bucket/ready.zip?signature=secret"}}), nil
	})
	m.now = func() time.Time { return now }
	installSigned(t, m, func(r *http.Request) (*http.Response, error) { return signedResponse(r, "ready.zip", 100), nil })
	broken := testJob(now.Add(-time.Hour))
	broken.ID = "broken"
	later := broken
	later.ID = "later"
	seedJobs(t, m, broken, later)
	exports, err := m.List(context.Background())
	if !IsJobFailure(err) || len(exports) != 0 || calls != 2 {
		t.Fatal("interactive failure lost classification", exports, err, calls)
	}
	phase = 1
	now = now.Add(5 * time.Minute)
	exports, err = m.List(context.Background())
	if err != nil || len(exports) != 1 || exports[0].ID != "portability-later" {
		t.Fatal("later ready sibling lost", exports, err)
	}
	for _, safety := range []error{&sourceError{message: "auth", authorization: true}, &sourceError{message: "backoff", temporary: true, delay: time.Hour}, errors.New("manifest invalid"), errors.New("journal corrupt"), context.Canceled} {
		if IsJobFailure(safety) {
			t.Fatal("safety failure classified as a job", safety)
		}
	}
}

func TestMissingJobWithProviderBackoffIsNotJobScoped(t *testing.T) {
	now := time.Now().UTC()
	calls := 0
	m := testManager(t, Options{Provider: "portability", Resources: []string{"myactivity.search"}}, func(r *http.Request) (*http.Response, error) {
		calls++
		resp := responseJSON(r, 404, map[string]any{"error": map[string]any{"code": 404}})
		resp.Header.Set("Retry-After", "3600")
		return resp, nil
	})
	m.now = func() time.Time { return now }
	seedJobs(t, m, testJob(now.Add(-time.Hour)))
	if _, err := m.List(context.Background()); err == nil || IsJobFailure(err) || calls != 1 {
		t.Fatal("provider backoff ignored", err, calls)
	}
	s := readJournal(t, m)
	if !s.Jobs[0].RetryNotBefore.Equal(now.Add(time.Hour)) {
		t.Fatal("backoff not durable")
	}
	now = now.Add(5 * time.Minute)
	if _, err := m.RefreshJobs(context.Background()); err != nil || calls != 1 {
		t.Fatal("refresh bypassed provider backoff", err, calls)
	}
}

func TestMetadataProbeTimeDoesNotExtendSignedURLLifetime(t *testing.T) {
	now := time.Now().UTC()
	m := testManager(t, Options{Provider: "portability", Resources: []string{"myactivity.search"}}, func(r *http.Request) (*http.Response, error) {
		return responseJSON(r, 200, &dataportability.PortabilityArchiveState{State: "COMPLETE", Urls: []string{"https://storage.googleapis.com/bucket/ready.zip?signature=secret"}}), nil
	})
	m.now = func() time.Time { return now }
	seedJobs(t, m, testJob(now.Add(-time.Hour)))
	installSigned(t, m, func(r *http.Request) (*http.Response, error) {
		now = now.Add(6 * time.Hour)
		return signedResponse(r, "ready.zip", 100), nil
	})
	exports, err := m.List(context.Background())
	if err != nil || len(exports) != 0 {
		t.Fatal("slow metadata made old signature fresh", err, exports)
	}
}

func TestSlowIntentPersistenceCannotOutliveExpiryMarginBeforePOST(t *testing.T) {
	for _, kind := range []string{"create", "retry"} {
		t.Run(kind, func(t *testing.T) {
			now := time.Now().UTC()
			posts, checks := 0, 0
			m := testManager(t, Options{Provider: "portability", Resources: []string{"myactivity.search"}}, func(r *http.Request) (*http.Response, error) {
				if r.URL.Path == "/v1/accessType:check" {
					checks++
					return responseJSON(r, 200, &dataportability.CheckAccessTypeResponse{TimeBasedResources: []string{"myactivity.search"}}), nil
				}
				posts++
				t.Fatal("mutation escaped bounded expiry margin")
				return nil, nil
			})
			m.now = func() time.Time { return now }
			j := testJob(now.Add(-time.Hour))
			j.State = "FAILED"
			if kind == "retry" {
				seedJobs(t, m, j)
			}
			m.persist = func(name string, s *portabilityState) error {
				if err := saveState(name, s); err != nil {
					return err
				}
				now = now.Add(5 * time.Minute)
				return nil
			}
			var err error
			if kind == "create" {
				_, err = m.CreateExport(context.Background())
			} else {
				_, err = m.RetryJob(context.Background(), j.ID)
			}
			if err == nil || posts != 0 || IsTemporary(err) {
				t.Fatal("slow durable write did not hold", err, posts)
			}
			if s := readJournal(t, m); s.Intent == nil || intentKind(*s.Intent) != kind {
				t.Fatal("lost conservative intent")
			}
			if kind == "create" && checks != 1 {
				t.Fatal("creation bypassed access check")
			}
		})
	}
}
