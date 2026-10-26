//go:build !headless

package gui

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/Dunky13/takedown/internal/source"
)

// Keep management optional: account implementations used for downloading do
// not need to implement Google-specific job mutations.
type jobManager interface {
	Jobs(context.Context) (source.JobSnapshot, error)
	RefreshJobs(context.Context) (source.JobSnapshot, error)
	RetryJob(context.Context, string) (string, error)
	CancelJob(context.Context, string) error
	ReconcileIntent(context.Context, string) error
	ResolveExpiredIntent(context.Context) error
	ArchiveJob(context.Context, string) error
	ResetAuthorization(context.Context) error
}

type jobsView struct {
	s                *screen
	snapshot         source.JobSnapshot
	binding          sourceSnapshot
	loaded           bool
	selectedID       string
	visible          []source.Job
	list             *widget.List
	details, pending *widget.Label
	knownID          *widget.Entry
	archived         *widget.Check
	buttons          map[string]*widget.Button
}

func (s *screen) jobsContent() fyne.CanvasObject {
	v := &jobsView{s: s, buttons: make(map[string]*widget.Button)}
	s.jobsUI = v
	v.details = archiveLabel("Load recorded jobs, then select a row for details and available actions.")
	v.pending = archiveLabel("No job history loaded. Loading reads local receipts; refreshing checks Google without creating exports.")
	v.knownID = widget.NewEntry()
	v.knownID.SetPlaceHolder("Known job ID from Google's response or your saved receipt")
	v.knownID.OnChanged = func(string) { v.setBusy(s.busy) }
	v.archived = widget.NewCheck("Include archived history", func(bool) { v.render() })
	v.list = widget.NewList(func() int { return len(v.visible) }, func() fyne.CanvasObject { return archiveLabel("Job\nState") }, func(i widget.ListItemID, obj fyne.CanvasObject) {
		j := v.visible[i]
		label := obj.(*widget.Label)
		label.SetText(jobTitle(j) + "\n" + strings.Join(j.Resources, ", "))
		width := v.list.Size().Width - 24
		if width > 0 {
			label.Resize(fyne.NewSize(width, label.MinSize().Height))
			v.list.SetItemHeight(i, label.MinSize().Height)
		}
	})
	v.list.OnSelected = func(id widget.ListItemID) { v.selectedID = v.visible[id].ID; v.showDetails(); v.setBusy(s.busy) }
	button := func(key, label string, action func()) *widget.Button {
		b := widget.NewButton(label, action)
		v.buttons[key] = b
		return b
	}
	load := button("load", "Load recorded jobs", func() { s.loadJobs(false) })
	refresh := button("refresh", "Refresh Google status", func() { s.loadJobs(true) })
	retry := button("retry", "Retry failed job", func() { v.mutate("retry") })
	cancel := button("cancel", "Cancel remote job", func() { v.mutate("cancel") })
	archive := button("archive", "Archive history", func() { v.mutate("archive") })
	reset := button("reset", "Reset Google access…", func() { v.mutate("reset") })
	reconcile := button("reconcile", "Reconcile known job", func() { v.mutate("reconcile") })
	expire := button("expire", "Resolve expired hold…", func() { v.mutate("expire") })
	v.knownID.Hide()
	reconcile.Hide()
	expire.Hide()
	return container.NewBorder(
		container.NewHBox(load, refresh, v.archived),
		container.NewVBox(widget.NewSeparator(), v.pending, v.knownID, container.NewHBox(reconcile, expire), container.NewHBox(retry, cancel, archive, reset)),
		nil, nil, container.NewHSplit(v.list, container.NewVScroll(v.details)),
	)
}

func jobTitle(j source.Job) string {
	created := "Date unavailable"
	if !j.CreatedAt.IsZero() {
		created = j.CreatedAt.Local().Format("02 Jan 2006 15:04")
	}
	state := strings.ToLower(strings.ReplaceAll(j.State, "_", " "))
	if state == "" {
		state = "status unknown"
	}
	if j.Archived {
		state += " · archived"
	}
	return created + " · " + state
}

func (v *jobsView) selected() *source.Job {
	for i := range v.snapshot.Jobs {
		if v.snapshot.Jobs[i].ID == v.selectedID {
			return &v.snapshot.Jobs[i]
		}
	}
	return nil
}

func (v *jobsView) showDetails() {
	j := v.selected()
	if j == nil {
		v.details.SetText("Select a job to view its status. Job failures do not hide other jobs.")
		return
	}
	text := jobTitle(*j) + "\nResource groups: " + strings.Join(j.Resources, ", ")
	access := "Access type unknown"
	if j.AccessType == "ACCESS_TYPE_ONE_TIME" {
		access = "One-time access; archiving does not clear its usage guard."
	}
	if j.AccessType == "ACCESS_TYPE_TIME_BASED" {
		access = "Time-based access"
	}
	text += "\n" + access + fmt.Sprintf("\nRetries used: %d of 3", j.RetryCount)
	for _, item := range []struct {
		name string
		at   time.Time
	}{{"Last checked", j.LastPoll}, {"Archive expires", j.ExpiresAt}, {"Retry available after", j.RetryNotBefore}} {
		if !item.at.IsZero() {
			text += "\n" + item.name + ": " + item.at.Local().Format("02 Jan 2006 15:04:05 MST")
		}
	}
	if j.Message != "" {
		text += "\n" + j.Message
	}
	v.details.SetText(text)
}

func (v *jobsView) render() {
	v.visible = nil
	for _, j := range v.snapshot.Jobs {
		if !j.Archived || v.archived.Checked {
			v.visible = append(v.visible, j)
		}
	}
	v.list.UnselectAll()
	v.list.Refresh()
	selectedVisible := false
	for i, j := range v.visible {
		if j.ID == v.selectedID {
			selectedVisible = true
			v.list.Select(i)
			break
		}
	}
	if !selectedVisible {
		v.selectedID = ""
	}
	v.showDetails()
	if p := v.snapshot.Intent; p != nil {
		text := "Held: unresolved " + p.Kind + " request. " + p.Message + "\nDo not create another export. Reconcile only a known Google job ID; an unknown result stays held."
		if p.JobID != "" {
			text += "\nRecorded recovery job ID: " + p.JobID
		}
		v.pending.SetText(text)
		v.knownID.Show()
		v.buttons["reconcile"].Show()
		v.buttons["expire"].Show()
	} else {
		v.pending.SetText(fmt.Sprintf("%d recorded jobs. Archived history keeps local files and one-time access guards. No unresolved request.", len(v.snapshot.Jobs)))
		v.knownID.Hide()
		v.buttons["reconcile"].Hide()
		v.buttons["expire"].Hide()
	}
	v.setBusy(v.s.busy)
}

func (s *screen) invalidateManagement() {
	v := s.jobsUI
	if v == nil {
		return
	}
	v.loaded = false
	v.snapshot = source.JobSnapshot{}
	v.binding = sourceSnapshot{}
	v.selectedID = ""
	v.visible = nil
	v.list.UnselectAll()
	v.list.Refresh()
	v.knownID.SetText("")
	v.knownID.Hide()
	v.buttons["reconcile"].Hide()
	v.buttons["expire"].Hide()
	v.details.SetText("Source or authorization changed. Load recorded jobs for this source.")
	v.pending.SetText("Previous history is hidden until reloaded.")
	v.setBusy(s.busy)
}

func (v *jobsView) setBusy(busy bool) {
	_, capable := v.s.a.Account.(jobManager)
	ready := !busy && capable && v.s.selectedSource() == "portability"
	bound := ready && v.loaded && v.binding.matches(v.s)
	allowed := map[string]bool{"load": ready, "refresh": ready}
	j := v.selected()
	if bound && v.snapshot.Intent == nil {
		if j != nil && !j.Archived {
			allowed["retry"] = j.State == "FAILED" && j.RetryCount < 3 && !time.Now().Before(j.RetryNotBefore)
			allowed["cancel"] = j.State == "IN_PROGRESS" && j.AccessType == "ACCESS_TYPE_TIME_BASED"
			allowed["archive"] = j.State == "COMPLETE" || j.State == "FAILED" || j.State == "EXPIRED" || j.State == "CANCELLED"
		}
		allowed["reset"] = true
		for _, job := range v.snapshot.Jobs {
			if job.State != "COMPLETE" {
				allowed["reset"] = false
			}
		}
	}
	if bound && v.snapshot.Intent != nil {
		allowed["reconcile"] = strings.TrimSpace(v.knownID.Text) != ""
		allowed["expire"] = true // The source verifies expiry and whether resolution is safe.
	}
	for key, b := range v.buttons {
		if allowed[key] {
			b.Enable()
		} else {
			b.Disable()
		}
	}
	if busy {
		v.knownID.Disable()
		v.archived.Disable()
	} else {
		v.knownID.Enable()
		v.archived.Enable()
	}
}

func (s *screen) loadJobs(refresh bool) {
	if s.busy || s.selectedSource() != "portability" {
		return
	}
	var snapshot source.JobSnapshot
	s.run("Reading recorded export jobs…", func(ctx context.Context, b sourceSnapshot) error {
		manager, ok := b.account.(jobManager)
		if !ok {
			return errors.New("job management is available for Data Portability sources")
		}
		var err error
		if refresh {
			snapshot, err = manager.RefreshJobs(ctx)
		} else {
			snapshot, err = manager.Jobs(ctx)
		}
		return err
	}, func(err error, b sourceSnapshot) {
		if !errors.Is(err, errSourceChanged) && !errors.Is(err, context.Canceled) && (err == nil || snapshot.Jobs != nil || snapshot.Intent != nil) {
			s.jobsUI.snapshot, s.jobsUI.binding, s.jobsUI.loaded = snapshot, b, true
			s.jobsUI.render()
		}
		if s.showError(err) {
			return
		}
		s.status.SetText("Job history loaded. Select a job to manage it; refreshing never creates exports.")
	}, nil)
}

func (v *jobsView) mutate(kind string) {
	if v.s.busy || !v.loaded || !v.binding.matches(v.s) {
		return
	}
	if b := v.buttons[kind]; b == nil || b.Disabled() {
		return
	}
	binding := v.binding
	id := v.selectedID
	knownID := strings.TrimSpace(v.knownID.Text)
	text := map[string]string{
		"retry":     "Retry this failed job? This sends one explicit retry request to Google. At most three retries are allowed; one-time access restrictions still apply.",
		"cancel":    "Cancel this in-progress, time-based Google job? Google may stop preparing its archive. Local downloads and files are kept.",
		"archive":   "Archive this job's history? Local files are kept. Archiving does not reset one-time access guards or make the grant reusable.",
		"reconcile": "Reconcile only with a Google response or saved receipt that belongs to this EXACT pending request, Google account, resource groups, and request time. The API can prove that a job is accessible and report its state, but cannot prove its association with this pending request. Do not guess a job ID or infer that no job was created.",
		"expire":    "Ask the source to resolve the expired request hold? It must verify that expiry is safe. An unknown outcome remains held when it cannot be resolved safely.",
		"reset":     "Reset Google Data Portability access? This revokes the Google grant and permits Google to remove old remote archives. This cannot be undone. Every job must be complete and its archives downloaded and verified by the source guard. Local files are kept; you must authorize again for future exports.",
	}[kind]
	if kind == "reconcile" {
		text += "\n\nKnown receipt ID: " + knownID
	}
	confirmation := dialog.NewConfirm("Confirm export job action", text, func(yes bool) {
		if !yes || v.s.busy || v.s.ctx.Err() != nil || !binding.matches(v.s) {
			return
		}
		var snapshot source.JobSnapshot
		var gotSnapshot bool
		v.s.run("Applying the confirmed job action…", func(ctx context.Context, b sourceSnapshot) error {
			manager := b.account.(jobManager)
			var err error
			switch kind {
			case "retry":
				_, err = manager.RetryJob(ctx, id)
			case "cancel":
				err = manager.CancelJob(ctx, id)
			case "archive":
				err = manager.ArchiveJob(ctx, id)
			case "reconcile":
				err = manager.ReconcileIntent(ctx, knownID)
			case "expire":
				err = manager.ResolveExpiredIntent(ctx)
			case "reset":
				err = manager.ResetAuthorization(source.WithResetConfirmation(ctx))
			}
			if ctx.Err() == nil {
				var readErr error
				snapshot, readErr = manager.Jobs(ctx)
				gotSnapshot = readErr == nil
				err = errors.Join(err, readErr)
			}
			return err
		}, func(err error, b sourceSnapshot) {
			if !errors.Is(err, errSourceChanged) && !errors.Is(err, context.Canceled) && gotSnapshot {
				v.snapshot, v.binding = snapshot, b
				v.render()
			}
			if errors.Is(err, context.Canceled) {
				v.loaded = false
				v.setBusy(false)
				v.s.status.SetText("Stopped. A remote request may have reached Google. Load recorded jobs before another action; unknown outcomes stay held.")
				return
			}
			if v.s.showError(err) {
				return
			}
			v.s.invalidateExports()
			if kind == "reset" {
				v.s.sourceRevision++
				v.s.invalidateManagement()
				v.s.accessStatus.SetText("Google access reset. Authorize again before requesting another export.")
			}
			v.s.status.SetText("Confirmed job action completed. Local files were kept.")
		}, nil)
	}, v.s.w)
	confirmation.SetConfirmText("Confirm action")
	confirmation.SetDismissText("Cancel")
	confirmation.Show()
}
