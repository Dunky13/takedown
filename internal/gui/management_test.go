//go:build !headless

package gui

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/Dunky13/takedown/internal/app"
	"github.com/Dunky13/takedown/internal/config"
	"github.com/Dunky13/takedown/internal/source"
)

type managedAccount struct {
	*screenAccount
	snapshot source.JobSnapshot
	calls    []string
	jobs     func(context.Context) (source.JobSnapshot, error)
	action   func(context.Context, string, string) error
	browse   func(context.Context, string) ([]source.Folder, error)
	access   func(context.Context) (source.AccessInfo, error)
}

func (a *managedAccount) Jobs(ctx context.Context) (source.JobSnapshot, error) {
	if a.jobs != nil {
		return a.jobs(ctx)
	}
	return a.snapshot, nil
}
func (a *managedAccount) RefreshJobs(ctx context.Context) (source.JobSnapshot, error) {
	a.calls = append(a.calls, "refresh")
	return a.Jobs(ctx)
}
func (a *managedAccount) apply(ctx context.Context, kind, id string) error {
	a.calls = append(a.calls, kind+":"+id)
	if a.action != nil {
		return a.action(ctx, kind, id)
	}
	return nil
}
func (a *managedAccount) RetryJob(ctx context.Context, id string) (string, error) {
	return "retry-receipt", a.apply(ctx, "retry", id)
}
func (a *managedAccount) CancelJob(ctx context.Context, id string) error {
	return a.apply(ctx, "cancel", id)
}
func (a *managedAccount) ArchiveJob(ctx context.Context, id string) error {
	return a.apply(ctx, "archive", id)
}
func (a *managedAccount) ReconcileIntent(ctx context.Context, id string) error {
	return a.apply(ctx, "reconcile", id)
}
func (a *managedAccount) ResolveExpiredIntent(ctx context.Context) error {
	return a.apply(ctx, "expire", "")
}
func (a *managedAccount) ResetAuthorization(ctx context.Context) error {
	return a.apply(ctx, "reset", "")
}
func (a *managedAccount) BrowseFolders(ctx context.Context, id string) ([]source.Folder, error) {
	if a.browse != nil {
		return a.browse(ctx, id)
	}
	return nil, nil
}
func (a *managedAccount) CheckAccess(ctx context.Context) (source.AccessInfo, error) {
	if a.access != nil {
		return a.access(ctx)
	}
	return source.AccessInfo{}, nil
}

func newManagedScreen(t *testing.T) (*screen, *managedAccount) {
	t.Helper()
	s, base := newTestScreen(t, func(cfg *config.Settings) {
		cfg.Source = "portability"
		cfg.PortabilityResources = []string{"myactivity.search"}
	})
	a := &managedAccount{screenAccount: base}
	s.a.Account = a
	s.setBusy(false)
	return s, a
}

func jobFixture(state string) source.Job {
	return source.Job{ID: "opaque-google-job-id", State: state, AccessType: "ACCESS_TYPE_TIME_BASED", Resources: []string{"myactivity.search"}, CreatedAt: time.Date(2026, 10, 25, 14, 0, 0, 0, time.UTC), ExpiresAt: time.Date(2026, 11, 1, 14, 0, 0, 0, time.UTC)}
}

func TestJobsShowIndependentStatesAndArchiveWithoutDeletingFiles(t *testing.T) {
	s, a := newManagedScreen(t)
	for i, state := range []string{"IN_PROGRESS", "FAILED", "EXPIRED", "COMPLETE", "CANCELLED"} {
		j := jobFixture(state)
		j.ID = fmt.Sprintf("opaque-%d", i)
		if state == "FAILED" {
			j.Message = "Google could not prepare this job. Other jobs can still complete."
		}
		a.snapshot.Jobs = append(a.snapshot.Jobs, j)
	}
	archived := jobFixture("COMPLETE")
	archived.ID = "opaque-archived"
	archived.Archived = true
	a.snapshot.Jobs = append(a.snapshot.Jobs, archived)
	var created atomic.Int32
	a.create = func(context.Context) (string, error) { created.Add(1); return "", nil }
	s.tabs.SelectIndex(4)
	test.Tap(s.jobsUI.buttons["load"])
	awaitJob(t, s)
	v := s.jobsUI
	if len(v.visible) != 5 || len(a.calls) != 0 || created.Load() != 0 {
		t.Fatal("local history did not keep job states independent or caused a remote mutation")
	}
	for i, j := range v.visible {
		label := archiveLabel("")
		v.list.UpdateItem(i, label)
		if strings.Contains(label.Text, j.ID) {
			t.Fatal("friendly job row exposed an opaque ID")
		}
	}
	v.list.Select(1)
	if !strings.Contains(v.details.Text, "Other jobs can still complete") || v.buttons["retry"].Disabled() || !v.buttons["cancel"].Disabled() {
		t.Fatal("failed job details/actions are wrong")
	}
	test.Tap(v.archived)
	if len(v.visible) != 6 {
		t.Fatal("archived history is inaccessible")
	}
	v.list.Select(3)
	test.Tap(v.buttons["archive"])
	if !dialogHasText(s.w.Canvas().Overlays().Top(), "Local files are kept") || !dialogHasText(s.w.Canvas().Overlays().Top(), "one-time access guards") || len(a.calls) != 0 {
		t.Fatal("archive mutated before explaining preserved files and guards")
	}
	tapDialogButton(t, s, "Cancel")
	if len(a.calls) != 0 {
		t.Fatal("dismissed archive mutated history")
	}
	a.action = func(_ context.Context, kind, id string) error {
		if kind == "archive" && id == "opaque-3" {
			a.snapshot.Jobs[3].Archived = true
		}
		return nil
	}
	test.Tap(v.buttons["archive"])
	tapDialogButton(t, s, "Confirm action")
	awaitJob(t, s)
	if len(a.calls) != 1 || a.calls[0] != "archive:opaque-3" || !v.snapshot.Jobs[3].Archived {
		t.Fatal("explicit archive did not persist in reloaded history")
	}
	for _, viewport := range []struct {
		size   fyne.Size
		suffix string
	}{{fyne.NewSize(960, 760), "960"}, {fyne.NewSize(800, 640), "800"}} {
		s.w.Resize(viewport.size)
		s.content.Refresh()
		assertWindowFits(t, s, viewport.size)
		captureScreen(t, s, "gui-jobs-"+viewport.suffix+".png")
	}
}

func TestJobEligibilityAndExplicitRefreshNeverCreatesExports(t *testing.T) {
	s, a := newManagedScreen(t)
	for _, tc := range []struct {
		state, access          string
		retries                int
		cooldown               bool
		retry, cancel, archive bool
	}{
		{"FAILED", "ACCESS_TYPE_TIME_BASED", 0, false, true, false, true},
		{"FAILED", "ACCESS_TYPE_TIME_BASED", 3, false, false, false, true},
		{"FAILED", "ACCESS_TYPE_TIME_BASED", 0, true, false, false, true},
		{"IN_PROGRESS", "ACCESS_TYPE_TIME_BASED", 0, false, false, true, false},
		{"IN_PROGRESS", "ACCESS_TYPE_ONE_TIME", 0, false, false, false, false},
		{"EXPIRED", "ACCESS_TYPE_TIME_BASED", 0, false, false, false, true},
		{"COMPLETE", "ACCESS_TYPE_ONE_TIME", 0, false, false, false, true},
	} {
		j := jobFixture(tc.state)
		j.AccessType = tc.access
		j.RetryCount = tc.retries
		if tc.cooldown {
			j.RetryNotBefore = time.Now().Add(time.Hour)
		}
		a.snapshot = source.JobSnapshot{Jobs: []source.Job{j}}
		s.loadJobs(false)
		awaitJob(t, s)
		s.jobsUI.list.Select(0)
		for name, want := range map[string]bool{"retry": tc.retry, "cancel": tc.cancel, "archive": tc.archive} {
			if s.jobsUI.buttons[name].Disabled() == want {
				t.Errorf("%+v: %s enabled=%v", tc, name, !s.jobsUI.buttons[name].Disabled())
			}
		}
	}
	test.Tap(s.jobsUI.buttons["refresh"])
	awaitJob(t, s)
	if len(a.calls) != 1 || a.calls[0] != "refresh" {
		t.Fatalf("refresh performed another action: %v", a.calls)
	}
}

func TestJobMutationsConfirmAndRejectStaleSource(t *testing.T) {
	for _, kind := range []string{"retry", "cancel", "archive", "reset"} {
		t.Run(kind, func(t *testing.T) {
			s, a := newManagedScreen(t)
			state := map[string]string{"retry": "FAILED", "cancel": "IN_PROGRESS", "archive": "COMPLETE", "reset": "COMPLETE"}[kind]
			a.snapshot.Jobs = []source.Job{jobFixture(state)}
			s.loadJobs(false)
			awaitJob(t, s)
			s.jobsUI.list.Select(0)
			test.Tap(s.jobsUI.buttons[kind])
			if len(a.calls) != 0 {
				t.Fatal("mutation occurred before confirmation")
			}
			if kind == "reset" {
				for _, text := range []string{"revokes the Google grant", "remove old remote archives", "cannot be undone", "downloaded and verified"} {
					if !dialogHasText(s.w.Canvas().Overlays().Top(), text) {
						t.Fatalf("reset warning missing %q", text)
					}
				}
			}
			tapDialogButton(t, s, "Cancel")
			if len(a.calls) != 0 {
				t.Fatal("canceled confirmation mutated")
			}
			test.Tap(s.jobsUI.buttons[kind])
			tapDialogButton(t, s, "Confirm action")
			awaitJob(t, s)
			if len(a.calls) != 1 || !strings.HasPrefix(a.calls[0], kind+":") {
				t.Fatalf("confirmed action=%v", a.calls)
			}
			if kind == "reset" {
				s.loadJobs(false)
				awaitJob(t, s)
				s.jobsUI.list.Select(0)
			}
			test.Tap(s.jobsUI.buttons[kind])
			s.email.SetText("other@example.com")
			tapDialogButton(t, s, "Confirm action")
			if len(a.calls) != 1 || s.busy || s.jobsUI.loaded || len(s.jobsUI.visible) != 0 {
				t.Fatal("stale confirmation acted on another source or retained its history")
			}
		})
	}
}

func TestUnknownIntentStaysHeldAndReconciliationRequiresKnownID(t *testing.T) {
	s, a := newManagedScreen(t)
	a.snapshot = source.JobSnapshot{Jobs: []source.Job{jobFixture("FAILED")}, Intent: &source.PendingAction{Kind: "create", Resources: []string{"myactivity.search"}, Message: "The response was lost. Google may have accepted the request."}}
	s.loadJobs(false)
	awaitJob(t, s)
	s.jobsUI.list.Select(0)
	v := s.jobsUI
	if !strings.Contains(v.pending.Text, "unknown result stays held") || !v.buttons["reconcile"].Disabled() || !v.buttons["retry"].Disabled() || !v.buttons["reset"].Disabled() {
		t.Fatal("ambiguity was not held")
	}
	v.knownID.SetText("known-from-receipt")
	test.Tap(v.buttons["reconcile"])
	for _, required := range []string{"EXACT pending request, Google account, resource groups, and request time", "cannot prove its association with this pending request", "Do not guess a job ID", "Known receipt ID: known-from-receipt"} {
		if !dialogHasText(s.w.Canvas().Overlays().Top(), required) {
			t.Fatalf("reconciliation confirmation missing %q", required)
		}
	}
	for _, size := range []fyne.Size{fyne.NewSize(960, 760), fyne.NewSize(800, 640)} {
		s.w.Resize(size)
		for _, name := range []string{"Cancel", "Confirm action"} {
			button := findDialogButton(s.w.Canvas().Overlays().Top(), name)
			if button == nil {
				t.Fatalf("reconciliation confirmation missing %q", name)
			}
			assertButtonFits(t, s, button, size)
		}
		captureScreen(t, s, fmt.Sprintf("gui-reconcile-confirmation-%d.png", int(size.Width)))
	}
	tapDialogButton(t, s, "Confirm action")
	awaitJob(t, s)
	if len(a.calls) != 1 || a.calls[0] != "reconcile:known-from-receipt" || v.snapshot.Intent == nil {
		t.Fatal("reconciliation guessed an ID or erased an unresolved result")
	}
	a.action = func(context.Context, string, string) error {
		return errors.New("request is still within its recovery window; keep this hold")
	}
	test.Tap(v.buttons["expire"])
	tapDialogButton(t, s, "Confirm action")
	awaitJob(t, s)
	if v.snapshot.Intent == nil || !strings.Contains(s.banner.Text, "keep this hold") {
		t.Fatal("failure cleared held ambiguity")
	}
	tapDialogButton(t, s, "Copy error details")
	if !strings.Contains(fyne.CurrentApp().Clipboard().Content(), "keep this hold") {
		t.Fatal("error details were not copyable")
	}
	tapDialogButton(t, s, "Close")
	s.tabs.SelectIndex(4)
	for _, size := range []fyne.Size{fyne.NewSize(960, 760), fyne.NewSize(800, 640)} {
		s.w.Resize(size)
		s.content.Refresh()
		assertWindowFits(t, s, size)
		for _, name := range []string{"load", "refresh", "reconcile", "expire", "retry", "cancel", "archive", "reset"} {
			assertButtonFits(t, s, v.buttons[name], size)
		}
		captureScreen(t, s, fmt.Sprintf("gui-jobs-held-%d.png", int(size.Width)))
	}
}

func TestJobLoadCancellationAndSourceChangeDiscardLateHistory(t *testing.T) {
	for _, change := range []bool{false, true} {
		s, a := newManagedScreen(t)
		entered, release := make(chan struct{}), make(chan struct{})
		a.jobs = func(context.Context) (source.JobSnapshot, error) {
			close(entered)
			<-release
			return source.JobSnapshot{Jobs: []source.Job{jobFixture("COMPLETE")}}, nil
		}
		s.loadJobs(false)
		receive(t, entered)
		if change {
			s.source.SetSelected(driveSourceLabel)
		} else {
			test.Tap(s.cancelButton)
		}
		close(release)
		awaitJob(t, s)
		if s.jobsUI.loaded || len(s.jobsUI.visible) != 0 || s.busy {
			t.Fatal("late job snapshot survived canceled/source-changed operation")
		}
	}
}

func TestPartialJobRefreshKeepsOtherStatusesVisible(t *testing.T) {
	s, a := newManagedScreen(t)
	a.jobs = func(context.Context) (source.JobSnapshot, error) {
		failed, complete := jobFixture("FAILED"), jobFixture("COMPLETE")
		complete.ID = "other-complete-job"
		return source.JobSnapshot{Jobs: []source.Job{failed, complete}}, errors.New("one Google status check failed; recorded statuses remain available")
	}
	s.loadJobs(true)
	awaitJob(t, s)
	if !s.jobsUI.loaded || len(s.jobsUI.visible) != 2 || s.jobsUI.visible[1].State != "COMPLETE" || !s.banner.Visible() {
		t.Fatal("one refresh error hid independent job states")
	}
	tapDialogButton(t, s, "Close")
}

func TestCanceledMutationRequiresReloadBeforeAnotherAction(t *testing.T) {
	s, a := newManagedScreen(t)
	a.snapshot.Jobs = []source.Job{jobFixture("IN_PROGRESS")}
	s.loadJobs(false)
	awaitJob(t, s)
	s.jobsUI.list.Select(0)
	entered := make(chan struct{})
	a.action = func(ctx context.Context, _, _ string) error {
		close(entered)
		<-ctx.Done()
		a.snapshot.Intent = &source.PendingAction{Kind: "cancel", JobID: a.snapshot.Jobs[0].ID, Message: "Cancellation outcome is unknown"}
		return ctx.Err()
	}
	test.Tap(s.jobsUI.buttons["cancel"])
	tapDialogButton(t, s, "Confirm action")
	receive(t, entered)
	test.Tap(s.cancelButton)
	awaitJob(t, s)
	if s.jobsUI.loaded || !s.jobsUI.buttons["cancel"].Disabled() || !strings.Contains(s.status.Text, "remote request may have reached Google") {
		t.Fatal("canceled mutation allowed replay or claimed only a local stop")
	}
	s.loadJobs(false)
	awaitJob(t, s)
	if s.jobsUI.snapshot.Intent == nil || !strings.Contains(s.jobsUI.pending.Text, "unresolved cancel") {
		t.Fatal("explicit reload lost durable cancellation ambiguity")
	}
}

func TestResourceCheckboxesAndAdvancedEntryStayExactAcrossFiltering(t *testing.T) {
	s, _ := newManagedScreen(t)
	v := s.resourceUI
	if len(v.checks) != len(source.Resources()) || !v.checks["myactivity.search"].Checked {
		t.Fatal("chooser did not use published catalog or existing selection")
	}
	s.resources.SetText("chrome.bookmarks, myactivity.youtube")
	if !v.checks["chrome.bookmarks"].Checked || !v.checks["myactivity.youtube"].Checked || v.checks["myactivity.search"].Checked {
		t.Fatal("advanced resource entry did not update checkbox selection")
	}
	v.filter.SetText("chrome")
	if v.checks["myactivity.youtube"].Visible() {
		t.Fatal("resource filter ignored search")
	}
	test.Tap(v.checks["chrome.bookmarks"])
	if s.resources.Text != "myactivity.youtube" || !strings.Contains(s.disclosure.Text, "only for these resource groups: myactivity.youtube") {
		t.Fatalf("filter changed hidden selection or stale disclosure: %s", s.resources.Text)
	}
	s.resources.SetText("unknown.group, myactivity.search")
	if !strings.Contains(v.summary.Text, "unpublished") {
		t.Fatal("advanced unknown resource silently presented as published")
	}
	v.filter.SetText("")
	test.Tap(v.checks["myactivity.youtube"])
	if s.resources.Text != "myactivity.search, myactivity.youtube" {
		t.Fatal("checkbox choice did not match selected published groups exactly")
	}
	s.tabs.SelectIndex(2)
	s.sourceTabs.SelectIndex(0)
	for _, size := range []fyne.Size{fyne.NewSize(960, 760), fyne.NewSize(800, 640)} {
		s.w.Resize(size)
		s.content.Refresh()
		assertWindowFits(t, s, size)
		captureScreen(t, s, fmt.Sprintf("gui-resource-chooser-%d.png", int(size.Width)))
		assertButtonFits(t, s, v.manual, size)
		assertScrollPaintContained(t, s, v.scroll)
	}
	s.setBusy(true)
	if !v.filter.Disabled() || !v.checks["myactivity.search"].Disabled() {
		t.Fatal("resource controls remained editable during work")
	}
	s.setBusy(false)
	s.resources.SetText(strings.Join(source.Resources(), ", "))
	s.content.Refresh()
	assertWindowFits(t, s, fyne.NewSize(800, 640))
	if len(v.summary.Text) > 240 || v.scroll.Size().Height < 75 {
		t.Fatal("a large resource selection displaced its chooser")
	}
}

func TestFolderPickerBrowseIsExplicitAndSelectionKeepsOpaqueIDInternal(t *testing.T) {
	s, base := newTestScreen(t, nil)
	a := &managedAccount{screenAccount: base}
	s.a.Account = a
	var parents []string
	a.browse = func(_ context.Context, parent string) ([]source.Folder, error) {
		parents = append(parents, parent)
		if parent == "" {
			return []source.Folder{{ID: "opaque-folder", Name: "Takeout archives"}}, nil
		}
		return nil, nil
	}
	if len(parents) != 0 {
		t.Fatal("folder browser contacted Google during setup")
	}
	s.browseFolders()
	awaitJob(t, s)
	p := s.folderUI
	if len(parents) != 1 || parents[0] != "" || !p.choose.Disabled() {
		t.Fatal("root browse or selection guard is wrong")
	}
	open := findDialogButton(s.w.Canvas().Overlays().Top(), "Open")
	if open == nil {
		t.Fatal("folder navigation missing")
	}
	test.Tap(open)
	awaitJob(t, s)
	if len(parents) != 2 || parents[1] != "opaque-folder" || p.choose.Disabled() || strings.Contains(p.location.Text, "opaque-folder") {
		t.Fatal("folder path did not use friendly names or explicit navigation")
	}
	for _, size := range []fyne.Size{fyne.NewSize(960, 760), fyne.NewSize(800, 640)} {
		s.w.Resize(size)
		s.content.Refresh()
		assertButtonFits(t, s, p.choose, size)
		assertButtonFits(t, s, p.up, size)
		captureScreen(t, s, fmt.Sprintf("gui-drive-folder-picker-%d.png", int(size.Width)))
	}
	test.Tap(p.choose)
	if s.driveFolder.Text != "opaque-folder" || !strings.Contains(s.status.Text, "Takeout archives") {
		t.Fatal("folder selection lost ID or friendly receipt")
	}
}

func TestFolderPickerCloseSourceSwitchAndFailure(t *testing.T) {
	for _, mode := range []string{"close", "source", "failure"} {
		t.Run(mode, func(t *testing.T) {
			s, base := newTestScreen(t, nil)
			a := &managedAccount{screenAccount: base}
			s.a.Account = a
			entered, release := make(chan struct{}), make(chan struct{})
			a.browse = func(context.Context, string) ([]source.Folder, error) {
				close(entered)
				<-release
				if mode == "failure" {
					return nil, errors.New("Drive access denied; authorize again")
				}
				return []source.Folder{{ID: "late-folder", Name: "Late folder"}}, nil
			}
			s.browseFolders()
			receive(t, entered)
			if mode == "close" {
				tapDialogButton(t, s, "Close")
			}
			if mode == "source" {
				s.source.SetSelected(portabilitySourceLabel)
			}
			close(release)
			awaitJob(t, s)
			if s.driveFolder.Text != "archive-folder" {
				t.Fatal("picker changed selection without an explicit choice")
			}
			if mode == "failure" {
				if !dialogHasText(s.w.Canvas().Overlays().Top(), "Drive access denied") || !s.folderUI.choose.Disabled() {
					t.Fatal("picker failure is hidden or permits selection")
				}
				tapDialogButton(t, s, "Close")
			} else if !s.folderUI.closed || s.w.Canvas().Overlays().Top() != nil {
				t.Fatal("canceled picker reopened after late response")
			}
		})
	}
}

func TestAccessStatusDistinguishesConfiguredFromReportedGrant(t *testing.T) {
	s, a := newManagedScreen(t)
	if !strings.Contains(s.accessStatus.Text, "not a verified identity") {
		t.Fatal("configured hint was presented as verified")
	}
	a.access = func(context.Context) (source.AccessInfo, error) {
		return source.AccessInfo{OneTimeResources: []string{"chrome.bookmarks"}, TimeBasedResources: []string{"myactivity.search"}}, nil
	}
	test.Tap(s.checkAccess)
	awaitJob(t, s)
	if !strings.Contains(s.accessStatus.Text, "One-time: chrome.bookmarks") || !strings.Contains(s.accessStatus.Text, "Time-based: myactivity.search") || !strings.Contains(s.accessStatus.Text, "does not verify the identity") {
		t.Fatal("actual reported grant is missing or misrepresents identity")
	}
	a.access = func(context.Context) (source.AccessInfo, error) {
		return source.AccessInfo{}, errors.New("grant expired; authorize again")
	}
	test.Tap(s.checkAccess)
	awaitJob(t, s)
	if !strings.Contains(s.accessStatus.Text, "could not be verified") {
		t.Fatal("failed access check retained prior verified-looking grant")
	}
	tapDialogButton(t, s, "Close")
}

func TestDiagnosticsAreExplicitOptInCopyableAndCancelable(t *testing.T) {
	s, _ := newManagedScreen(t)
	v := s.diagnosticsUI
	var options []app.DiagnosticOptions
	v.diagnose = func(_ context.Context, opts app.DiagnosticOptions) (app.DiagnosticReport, error) {
		options = append(options, opts)
		return app.DiagnosticReport{Checks: []app.DiagnosticCheck{{Name: "preferences", Status: "pass", Message: "Valid settings"}, {Name: "Google access", Status: "warning", Message: "Not contacted. Enable online diagnostics."}, {Name: "destination", Status: "fail", Message: "Check trusted host key and free space."}}}, nil
	}
	if len(options) != 0 || v.online.Checked || v.storage.Checked {
		t.Fatal("diagnostics started or opted into remote checks automatically")
	}
	s.tabs.SelectIndex(5)
	test.Tap(v.run)
	awaitJob(t, s)
	if len(options) != 1 || options[0].Online || options[0].Storage || !strings.Contains(v.result.Text, "FAIL · destination") {
		t.Fatal("local diagnostics options or per-check result is wrong")
	}
	s.password.SetText("session-secret")
	test.Tap(v.online)
	test.Tap(v.storage)
	test.Tap(v.run)
	awaitJob(t, s)
	if len(options) != 2 || !options[1].Online || !options[1].Storage || options[1].StoragePassword != "session-secret" || strings.Contains(v.result.Text, "session-secret") {
		t.Fatal("explicit diagnostic options were lost or password leaked")
	}
	test.Tap(v.copy)
	if fyne.CurrentApp().Clipboard().Content() != v.result.Text {
		t.Fatal("diagnostics results are not copyable")
	}
	for _, size := range []fyne.Size{fyne.NewSize(960, 760), fyne.NewSize(800, 640)} {
		s.w.Resize(size)
		s.content.Refresh()
		assertWindowFits(t, s, size)
		captureScreen(t, s, fmt.Sprintf("gui-diagnostics-%d.png", int(size.Width)))
		assertButtonFits(t, s, v.run, size)
		assertButtonFits(t, s, v.copy, size)
	}
	entered := make(chan struct{})
	v.diagnose = func(ctx context.Context, _ app.DiagnosticOptions) (app.DiagnosticReport, error) {
		close(entered)
		<-ctx.Done()
		return app.DiagnosticReport{}, ctx.Err()
	}
	before := v.result.Text
	test.Tap(v.run)
	receive(t, entered)
	if !v.run.Disabled() || !v.online.Disabled() || !v.storage.Disabled() {
		t.Fatal("diagnostic controls remained active during work")
	}
	test.Tap(s.cancelButton)
	awaitJob(t, s)
	if v.result.Text != before || s.busy || v.run.Disabled() {
		t.Fatal("canceled diagnostics overwrote results or kept controls locked")
	}
}

func assertButtonFits(t *testing.T, s *screen, b *widget.Button, size fyne.Size) {
	t.Helper()
	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(b)
	if pos.X < 0 || pos.Y < 0 || pos.X+b.Size().Width > size.Width || pos.Y+b.Size().Height > size.Height || b.Size().Width < b.MinSize().Width {
		t.Fatalf("management button %q clipped at %v size %v in %v", b.Text, pos, b.Size(), size)
	}
}

func assertScrollPaintContained(t *testing.T, s *screen, scroll *container.Scroll) {
	t.Helper()
	scroll.ScrollToTop()
	s.content.Refresh()
	before := s.w.Canvas().Capture()
	pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(scroll)
	bottom := pos.Y + scroll.Size().Height
	scroll.ScrollToBottom()
	after := s.w.Canvas().Capture()
	for y := before.Bounds().Min.Y; y < before.Bounds().Max.Y; y++ {
		if float32(y) >= pos.Y && float32(y) < bottom {
			continue
		}
		for x := before.Bounds().Min.X; x < before.Bounds().Max.X; x++ {
			if before.At(x, y) != after.At(x, y) {
				t.Fatalf("resource scroll painted outside viewport at (%d,%d)", x, y)
			}
		}
	}
	scroll.ScrollToTop()
}

func TestOperationFailureBannerGetsOwnLayoutSpace(t *testing.T) {
	s, _ := newManagedScreen(t)
	s.run("Reading history…", func(context.Context, sourceSnapshot) error {
		return errors.New("Google status unavailable. Check diagnostics, then refresh recorded jobs.")
	}, nil, nil)
	awaitJob(t, s)
	tapDialogButton(t, s, "Close")
	if !s.banner.Visible() || s.banner.Position().Y <= 0 {
		t.Fatal("operation error banner overlaps application heading")
	}
}

func TestStartupPreviewStatusAndConfirmedChangesUseInjectedAPI(t *testing.T) {
	s, _ := newManagedScreen(t)
	v := s.startupUI
	type call struct {
		action    string
		part, dry bool
	}
	var calls []call
	v.configure = func(_ context.Context, action string, part, dry bool) (string, error) {
		calls = append(calls, call{action, part, dry})
		return "User login watcher definition/status. No secrets are included.", nil
	}
	if len(calls) != 0 {
		t.Fatal("startup inspected or changed automatically")
	}
	s.partByPart.SetChecked(true)
	s.tabs.SelectIndex(5)
	s.tabs.Items[5].Content.(*container.AppTabs).SelectIndex(1)
	test.Tap(v.buttons["preview"])
	awaitJob(t, s)
	test.Tap(v.buttons["status"])
	awaitJob(t, s)
	if len(calls) != 2 || calls[0] != (call{"install", true, true}) || calls[1] != (call{"status", true, false}) {
		t.Fatalf("startup read actions=%+v", calls)
	}
	for _, action := range []string{"install", "remove"} {
		before := len(calls)
		test.Tap(v.buttons[action])
		if len(calls) != before {
			t.Fatal("startup changed before confirmation")
		}
		if action == "install" {
			for _, text := range []string{"begin downloading ready archives immediately", "does not store the session password", "No Google sign-in or export request"} {
				if !dialogHasText(s.w.Canvas().Overlays().Top(), text) {
					t.Fatalf("install confirmation missing %q", text)
				}
			}
		}
		tapDialogButton(t, s, "Cancel")
		if len(calls) != before {
			t.Fatal("dismissed startup change executed")
		}
		test.Tap(v.buttons[action])
		tapDialogButton(t, s, "Confirm startup change")
		awaitJob(t, s)
		if len(calls) != before+1 || calls[len(calls)-1] != (call{action, true, false}) {
			t.Fatal("confirmed startup action did not use exact API options")
		}
	}
	for _, size := range []fyne.Size{fyne.NewSize(960, 760), fyne.NewSize(800, 640)} {
		s.w.Resize(size)
		s.content.Refresh()
		assertWindowFits(t, s, size)
		for _, key := range []string{"preview", "status", "install", "remove"} {
			assertButtonFits(t, s, v.buttons[key], size)
		}
		captureScreen(t, s, fmt.Sprintf("gui-startup-%d.png", int(size.Width)))
	}
}

func TestStartupConfirmationRejectsChangedSettingsAndFailuresAreCopyable(t *testing.T) {
	for _, field := range []string{"source", "destination", "workers", "part-by-part"} {
		t.Run(field, func(t *testing.T) {
			s, _ := newManagedScreen(t)
			var calls atomic.Int32
			s.startupUI.configure = func(context.Context, string, bool, bool) (string, error) { calls.Add(1); return "", nil }
			test.Tap(s.startupUI.buttons["install"])
			switch field {
			case "source":
				s.source.SetSelected(driveSourceLabel)
			case "destination":
				s.destination.SetText(t.TempDir())
			case "workers":
				s.workers.SetSelected("4")
			case "part-by-part":
				s.partByPart.SetChecked(true)
			}
			tapDialogButton(t, s, "Confirm startup change")
			if calls.Load() != 0 || s.busy {
				t.Fatal("stale startup confirmation applied new preferences")
			}
		})
	}
	s, _ := newManagedScreen(t)
	s.startupUI.configure = func(context.Context, string, bool, bool) (string, error) {
		return "", errors.New("Store OAuth authorization and destination credentials before installing")
	}
	test.Tap(s.startupUI.buttons["install"])
	tapDialogButton(t, s, "Confirm startup change")
	awaitJob(t, s)
	if !strings.Contains(s.banner.Text, "Store OAuth authorization") {
		t.Fatal("startup prerequisite failure is hidden")
	}
	tapDialogButton(t, s, "Copy error details")
	if !strings.Contains(fyne.CurrentApp().Clipboard().Content(), "destination credentials") {
		t.Fatal("startup failure is not copyable")
	}
	tapDialogButton(t, s, "Close")
}

func TestLargeResourceDisclosureKeepsConfirmationControlsVisible(t *testing.T) {
	s, account := newTestScreen(t, func(cfg *config.Settings) { cfg.Source = "portability"; cfg.PortabilityResources = source.Resources() })
	var logins atomic.Int32
	account.login = func(context.Context) error { logins.Add(1); return nil }
	s.w.Resize(fyne.NewSize(800, 640))
	s.login()
	scroll := findDialogScroll(s.w.Canvas().Overlays().Top())
	if scroll == nil {
		t.Fatal("large OAuth disclosure has no scrollable viewport; safety text can be hidden behind confirmation controls")
	}
	for _, name := range []string{"Cancel", "Continue to Google"} {
		b := findDialogButton(s.w.Canvas().Overlays().Top(), name)
		if b == nil {
			t.Fatalf("missing consent control %q", name)
		}
		assertButtonFits(t, s, b, fyne.NewSize(800, 640))
	}
	if logins.Load() != 0 {
		t.Fatal("large disclosure started OAuth without confirmation")
	}
	assertScrollPaintContained(t, s, scroll)
	captureScreen(t, s, "gui-many-resources-oauth-800.png")
	tapDialogButton(t, s, "Cancel")
}

func findDialogScroll(obj fyne.CanvasObject) *container.Scroll {
	if scroll, ok := obj.(*container.Scroll); ok {
		return scroll
	}
	switch obj := obj.(type) {
	case *widget.PopUp:
		return findDialogScroll(obj.Content)
	case *fyne.Container:
		for _, child := range obj.Objects {
			if scroll := findDialogScroll(child); scroll != nil {
				return scroll
			}
		}
	case fyne.Widget:
		for _, child := range test.WidgetRenderer(obj).Objects() {
			if scroll := findDialogScroll(child); scroll != nil {
				return scroll
			}
		}
	}
	return nil
}

func TestDestinationPasswordSaveForgetAreExplicitAndNeverExposeSecret(t *testing.T) {
	s, _ := newManagedScreen(t)
	v := s.startupUI
	s.destination.SetText("smb://operator@nas.invalid/backups")
	s.password.SetText("session-only-password")
	var saves, forgets int
	v.savePassword = func(dir, destination, password string) error {
		if dir != s.a.DataDir || destination != "smb://operator@nas.invalid/backups" || password != "session-only-password" {
			t.Error("vault save received unexpected binding")
		}
		saves++
		return nil
	}
	v.forgetPassword = func(dir, destination string) error {
		if dir != s.a.DataDir || destination != "smb://operator@nas.invalid/backups" {
			t.Error("vault forget received unexpected binding")
		}
		forgets++
		return nil
	}
	if saves != 0 || forgets != 0 || !s.password.Password {
		t.Fatal("credential action ran automatically or field is unmasked")
	}
	for _, action := range []string{"save-password", "forget-password"} {
		test.Tap(v.buttons[action])
		if dialogHasText(s.w.Canvas().Overlays().Top(), "session-only-password") {
			t.Fatal("credential confirmation exposed the password")
		}
		if action == "save-password" && !dialogHasText(s.w.Canvas().Overlays().Top(), "not written to preferences or the service definition") {
			t.Fatal("save confirmation omitted storage boundary")
		}
		if action == "forget-password" && !dialogHasText(s.w.Canvas().Overlays().Top(), "archive files and partial downloads are kept") {
			t.Fatal("forget confirmation omitted file preservation")
		}
		beforeSaves, beforeForgets := saves, forgets
		tapDialogButton(t, s, "Cancel")
		if saves != beforeSaves || forgets != beforeForgets {
			t.Fatal("dismissed credential action reached the vault")
		}
		test.Tap(v.buttons[action])
		tapDialogButton(t, s, "Confirm destination authorization")
		awaitJob(t, s)
	}
	if saves != 1 || forgets != 1 || s.password.Text != "" || strings.Contains(v.credentialResult.Text, "session-only-password") {
		t.Fatal("confirmed credential actions, session clearing, or secret boundary failed")
	}
	configBytes, err := os.ReadFile(filepath.Join(s.a.DataDir, "config.json"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(configBytes), "session-only-password") {
		t.Fatal("saved preferences contain destination password")
	}
	s.tabs.SelectIndex(5)
	s.tabs.Items[5].Content.(*container.AppTabs).SelectIndex(2)
	for _, size := range []fyne.Size{fyne.NewSize(960, 760), fyne.NewSize(800, 640)} {
		s.w.Resize(size)
		s.content.Refresh()
		assertWindowFits(t, s, size)
		for _, key := range []string{"save-password", "forget-password"} {
			assertButtonFits(t, s, v.buttons[key], size)
		}
		captureScreen(t, s, fmt.Sprintf("gui-destination-access-%d.png", int(size.Width)))
	}
}

func TestDestinationPasswordConfirmationRejectsStaleBindingsAndShowsVaultFailure(t *testing.T) {
	for _, action := range []string{"save-password", "forget-password"} {
		for _, field := range []string{"source", "destination", "password"} {
			t.Run(action+"/"+field, func(t *testing.T) {
				s, _ := newManagedScreen(t)
				v := s.startupUI
				s.destination.SetText("sftp://operator@archive.invalid/backups")
				s.password.SetText("session-secret")
				var calls atomic.Int32
				v.savePassword = func(string, string, string) error { calls.Add(1); return nil }
				v.forgetPassword = func(string, string) error { calls.Add(1); return nil }
				test.Tap(v.buttons[action])
				switch field {
				case "source":
					s.source.SetSelected(driveSourceLabel)
				case "destination":
					s.destination.SetText("sftp://operator@other.invalid/backups")
				case "password":
					s.password.SetText("changed-secret")
				}
				tapDialogButton(t, s, "Confirm destination authorization")
				if calls.Load() != 0 || s.busy {
					t.Fatal("credential confirmation acted on edited source, destination, or password")
				}
			})
		}
	}
	s, _ := newManagedScreen(t)
	v := s.startupUI
	s.destination.SetText("smb://operator@nas.invalid/backups")
	s.password.SetText("session-secret")
	v.savePassword = func(string, string, string) error {
		return errors.New("OS keyring unavailable; storage password was not saved")
	}
	v.forgetPassword = func(string, string) error {
		return errors.New("OS keyring unavailable; storage password could not be removed")
	}
	for _, action := range []string{"save-password", "forget-password"} {
		test.Tap(v.buttons[action])
		tapDialogButton(t, s, "Confirm destination authorization")
		awaitJob(t, s)
		if !strings.Contains(s.banner.Text, "OS keyring unavailable") || s.password.Text != "session-secret" {
			t.Fatal("vault failure was hidden or session password was lost")
		}
		tapDialogButton(t, s, "Close")
	}
}
