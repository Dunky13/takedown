//go:build !headless

package gui

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strconv"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Dunky13/takedown/internal/app"
	"github.com/Dunky13/takedown/internal/events"
	"github.com/Dunky13/takedown/internal/model"
	"github.com/Dunky13/takedown/internal/space"
	"github.com/Dunky13/takedown/internal/state"
)

var errSourceChanged = errors.New("source settings or authorization changed; discover archives again")

type sourceSnapshot struct {
	identity string
	revision uint64
	account  app.Account
}

func (b sourceSnapshot) matches(s *screen) bool {
	// Implementations are normally pointers. Fail closed for an account value
	// that cannot be compared rather than trusting URLs from another account.
	return b.account != nil && reflect.TypeOf(b.account).Comparable() && b.account == s.a.Account &&
		b.identity == s.a.Settings.SourceIdentity() && b.revision == s.sourceRevision
}

type screen struct {
	ctx                                                        context.Context
	cancel                                                     context.CancelFunc
	wg                                                         sync.WaitGroup
	w                                                          fyne.Window
	a                                                          *app.App
	disclosureScroll                                           *container.Scroll
	content                                                    fyne.CanvasObject
	destination, email, password, buffer, sshKey, knownHosts   *widget.Entry
	oauthClient, driveFolder, resources, pollInterval          *widget.Entry
	source, workers                                            *widget.Select
	partByPart, unknown, watchSelected                         *widget.Check
	status, banner, sourceSummary, sourceHelp, disclosure      *widget.Label
	progress                                                   *widget.ProgressBar
	exports                                                    *fyne.Container
	transferList                                               *widget.List
	transfers                                                  []events.Event
	selected                                                   map[string]bool
	ready                                                      []model.Export
	buttons                                                    []*widget.Button
	authorize, disconnect, create, find, watch, downloadButton *widget.Button
	pause, resume, cancelButton                                *widget.Button
	driveOptions, portabilityOptions                           fyne.CanvasObject
	jobCancel                                                  context.CancelFunc
	jobDone                                                    <-chan struct{}
	controlPending                                             bool
	afterControl                                               func()
	busy, watching                                             bool
	stopEvents                                                 func()
	tabs, sourceTabs                                           *container.AppTabs
	jobGeneration                                              uint64
	sourceRevision                                             uint64
	discovery                                                  sourceSnapshot
	discoveryBound                                             bool
	jobsUI                                                     *jobsView
	resourceUI                                                 *resourceChooser
	folderBrowse                                               *widget.Button
	folderUI                                                   *folderPicker
	accessStatus                                               *widget.Label
	checkAccess                                                *widget.Button
	diagnosticsUI                                              *diagnosticsView
	startupUI                                                  *startupView
}

func newScreen(parent context.Context, desktop fyne.App, w fyne.Window, a *app.App) *screen {
	ctx, cancel := context.WithCancel(parent)
	s := &screen{ctx: ctx, cancel: cancel, w: w, a: a, selected: make(map[string]bool)}
	s.destination = widget.NewEntry()
	s.destination.SetText(a.Settings.Destination)
	s.destination.SetPlaceHolder("Folder, Windows UNC, smb://user@host/share/path or sftp://user@host/path")
	s.password = widget.NewPasswordEntry()
	s.password.SetPlaceHolder("Storage password; kept in memory for this session")
	s.buffer = widget.NewEntry()
	s.buffer.SetText(strconv.FormatUint(a.Settings.SafetyBuffer, 10))
	s.workers = widget.NewSelect([]string{"1", "2", "3", "4", "8", "16"}, nil)
	s.workers.SetSelected(strconv.Itoa(a.Settings.Concurrency))
	if s.workers.Selected == "" {
		s.workers.Options = append(s.workers.Options, strconv.Itoa(a.Settings.Concurrency))
		s.workers.SetSelected(strconv.Itoa(a.Settings.Concurrency))
	}
	s.sshKey = widget.NewEntry()
	s.sshKey.SetText(a.Settings.PrivateKeyPath)
	s.knownHosts = widget.NewEntry()
	s.knownHosts.SetText(a.Settings.KnownHostsPath)
	s.partByPart = widget.NewCheck("Download one part at a time", nil)
	s.unknown = widget.NewCheck("Allow SFTP without free-space reporting", nil)
	s.unknown.SetChecked(a.Settings.AllowUnknownSpace)
	s.status = wrappedLabel("Set up a source, authorize in your browser, then discover archives.")
	s.banner = wrappedLabel("")
	s.banner.Importance = widget.WarningImportance
	s.banner.Hide()
	s.progress = widget.NewProgressBar()
	s.exports = container.NewVBox(wrappedLabel("No archives discovered. Set up your source, then choose Discover archives."))
	s.transferList = widget.NewList(func() int { return len(s.transfers) }, func() fyne.CanvasObject {
		return archiveLabel("Archive part\nWaiting for progress")
	}, func(i widget.ListItemID, obj fyne.CanvasObject) {
		label := obj.(*widget.Label)
		label.SetText(transferText(s.transfers[i]))
		width := s.transferList.Size().Width - theme.ScrollBarSize() - 2*theme.Padding()
		if width > 0 {
			label.Resize(fyne.NewSize(width, label.MinSize().Height))
			s.transferList.SetItemHeight(i, label.MinSize().Height)
		}
	})
	s.authorize = widget.NewButtonWithIcon("Authorize in browser", theme.LoginIcon(), s.login)
	s.disconnect = widget.NewButton("Remove local token", s.forget)
	s.find = widget.NewButtonWithIcon("Discover archives", theme.ViewRefreshIcon(), s.discover)
	s.downloadButton = widget.NewButtonWithIcon("Download selected", theme.DownloadIcon(), s.download)
	s.watch = widget.NewButtonWithIcon("Start watching", theme.MediaPlayIcon(), s.startWatch)
	s.watchSelected = widget.NewCheck("Only watch selected exports", nil)
	save := widget.NewButtonWithIcon("Save preferences", theme.DocumentSaveIcon(), func() {
		s.run("Saving preferences…", func(context.Context, sourceSnapshot) error { return nil }, func(err error, _ sourceSnapshot) {
			if !s.showError(err) {
				s.status.SetText("Preferences saved.")
			}
		}, nil)
	})
	folder := widget.NewButtonWithIcon("Browse…", theme.FolderOpenIcon(), func() {
		if s.busy {
			return
		}
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				s.errorDialog(err)
				return
			}
			if uri != nil && !s.busy {
				s.destination.SetText(uri.Path())
			}
		}, w)
	})
	s.pause = widget.NewButtonWithIcon("Pause downloads", theme.MediaPauseIcon(), func() { s.changeDownloadPause(true) })
	s.resume = widget.NewButtonWithIcon("Resume / retry space", theme.MediaPlayIcon(), func() { s.changeDownloadPause(false) })
	s.cancelButton = widget.NewButtonWithIcon("Stop", theme.MediaStopIcon(), func() {
		if s.jobCancel != nil {
			s.jobCancel()
			s.status.SetText("Stopping; partial files will be kept.")
		}
	})
	s.buttons = []*widget.Button{s.authorize, s.find, s.disconnect, save, folder, s.downloadButton, s.watch}
	sourceContent := s.sourceContent()
	destination := container.NewVBox(
		widget.NewForm(widget.NewFormItem("Network storage password", s.password)),
		wrappedLabel("Used only for SMB or SFTP storage. It is not a Google password and is not saved in preferences."),
		s.partByPart,
		widget.NewForm(widget.NewFormItem("Free-space buffer (bytes)", s.buffer), widget.NewFormItem("Concurrent parts", s.workers), widget.NewFormItem("SFTP private key", s.sshKey), widget.NewFormItem("SFTP known_hosts", s.knownHosts), widget.NewFormItem("Free-space fallback", s.unknown)),
		wrappedLabel("Without free-space reporting, the space guard cannot protect the SFTP destination."),
	)
	exportScroll := container.NewVScroll(s.exports)
	exportScroll.SetMinSize(fyne.NewSize(300, 150))
	archives := container.NewBorder(nil, container.NewVBox(s.watchSelected, container.NewHBox(s.downloadButton, s.watch), wrappedLabel("Watching downloads ready archives from this source, including new exports. It does not create or schedule exports.")), nil, nil, exportScroll)
	transfers := container.NewBorder(nil, s.progress, nil, nil, s.transferList)
	s.tabs = container.NewAppTabs(container.NewTabItem("Archives", archives), container.NewTabItem("Transfers", transfers), container.NewTabItem("Source", sourceContent), container.NewTabItem("Destination", container.NewVScroll(destination)), container.NewTabItem("Jobs", s.jobsContent()), container.NewTabItem("Diagnostics", s.diagnosticsContent()))
	s.content = container.NewBorder(
		container.NewVBox(widget.NewLabelWithStyle("takedown", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), s.banner, s.sourceSummary, container.NewHBox(s.authorize, s.find, s.disconnect), container.NewBorder(nil, nil, widget.NewLabel("Save archives to"), folder, s.destination), widget.NewSeparator()),
		container.NewVBox(widget.NewSeparator(), s.status, container.NewHBox(save, s.pause, s.resume, s.cancelButton)), nil, nil, s.tabs,
	)
	s.setBusy(false)
	s.bindSourceChanges()
	if a.Settings.OAuthClientFile == "" {
		s.tabs.SelectIndex(2)
	}
	ch, stop := a.Bus.Subscribe(512)
	s.stopEvents = stop
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		lastNotice := make(map[events.Kind]time.Time)
		for {
			select {
			case <-ctx.Done():
				return
			case e, ok := <-ch:
				if !ok {
					return
				}
				fyne.Do(func() {
					if ctx.Err() == nil {
						event := e
						if !s.busy {
							if event.Kind != events.Progress && event.Kind != events.PartCompleted {
								return
							}
							event.Message = ""
						}
						s.onEvent(event)
					}
				})
				if (e.Kind == events.SpaceWarning || e.Kind == events.ReAuthRequired || e.Kind == events.PartCompleted) && time.Since(lastNotice[e.Kind]) > 15*time.Second {
					lastNotice[e.Kind] = time.Now()
					fyne.Do(func() {
						if ctx.Err() == nil {
							desktop.SendNotification(fyne.NewNotification("takedown", e.Message))
						}
					})
				}
			}
		}
	}()
	if records, err := state.ReadSnapshot(filepath.Join(a.DataDir, "takedown-state.json")); err == nil {
		for _, r := range records {
			message := "Partial file saved"
			if r.Complete {
				message = "Archive complete"
			}
			s.transfers = append(s.transfers, events.Event{Kind: events.Status, PartKey: r.Part.Key(), Filename: r.Part.Filename, Bytes: r.Offset, Total: r.Part.Size, Message: message})
		}
	}
	return s
}

func wrappedLabel(text string) *widget.Label {
	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapWord
	return label
}

func (s *screen) close() { s.cancel(); s.stopEvents(); s.wg.Wait() }

func (s *screen) invalidateExports() {
	s.ready = nil
	s.selected = make(map[string]bool)
	s.discoveryBound = false
	s.discovery = sourceSnapshot{}
	s.exports.RemoveAll()
	s.exports.Add(wrappedLabel("Discover archives for the current source and authorization before downloading."))
}

func (s *screen) setBusy(busy bool) {
	s.busy = busy
	for _, b := range s.buttons {
		if busy {
			b.Disable()
		} else {
			b.Enable()
		}
	}
	for _, e := range []*widget.Entry{s.destination, s.email, s.password, s.buffer, s.sshKey, s.knownHosts, s.oauthClient, s.driveFolder, s.resources, s.pollInterval} {
		if busy {
			e.Disable()
		} else {
			e.Enable()
		}
	}
	for _, selectBox := range []*widget.Select{s.source, s.workers} {
		if busy {
			selectBox.Disable()
		} else {
			selectBox.Enable()
		}
	}
	for _, check := range []*widget.Check{s.partByPart, s.unknown, s.watchSelected} {
		if busy {
			check.Disable()
		} else {
			check.Enable()
		}
	}
	if s.resourceUI != nil {
		s.resourceUI.setBusy(busy)
	}
	if s.jobsUI != nil {
		s.jobsUI.setBusy(busy)
	}
	if s.diagnosticsUI != nil {
		s.diagnosticsUI.setBusy(busy)
	}
	if s.startupUI != nil {
		s.startupUI.setBusy(busy)
	}
	for _, obj := range s.exports.Objects {
		if check, ok := obj.(*widget.Check); ok {
			if busy {
				check.Disable()
			} else {
				check.Enable()
			}
		}
	}
	if busy {
		s.cancelButton.Enable()
	} else {
		s.watching = false
		s.pause.Disable()
		s.resume.Disable()
		s.cancelButton.Disable()
		s.cancelButton.SetText("Stop")
	}
	if s.content != nil {
		s.content.Refresh()
	}
}

func (s *screen) run(message string, work func(context.Context, sourceSnapshot) error, done func(error, sourceSnapshot), started func()) bool {
	if s.busy || s.ctx.Err() != nil {
		return false
	}
	if err := s.save(); err != nil {
		s.errorDialog(err)
		return false
	}
	binding := sourceSnapshot{identity: s.a.Settings.SourceIdentity(), revision: s.sourceRevision, account: s.a.Account}
	ctx, cancel := context.WithCancel(s.ctx)
	finished := make(chan struct{})
	s.jobDone = finished
	s.jobCancel = cancel
	s.setBusy(true)
	s.jobGeneration++
	s.banner.Hide()
	s.status.SetText(message)
	if started != nil {
		started()
	}
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer close(finished)
		err := work(ctx, binding)
		if err == nil && ctx.Err() != nil {
			err = ctx.Err()
		}
		cancel()
		fyne.Do(func() {
			if s.ctx.Err() != nil {
				return
			}
			finish := func() {
				if !binding.matches(s) {
					err = errSourceChanged
					s.invalidateExports()
				}
				s.jobCancel = nil
				s.setBusy(false)
				if errors.Is(err, errSourceChanged) {
					s.status.SetText("Source changed. Authorize if needed, then discover archives again.")
				} else if errors.Is(err, context.Canceled) {
					s.status.SetText("Stopped. Partial files are available for the next run.")
				} else if err != nil {
					s.status.SetText(err.Error())
					s.banner.SetText(err.Error())
					s.banner.Show()
				} else {
					s.status.SetText("Ready.")
				}
				s.content.Refresh()
				if done != nil {
					done(err, binding)
				} else {
					s.showError(err)
				}
			}
			if s.controlPending {
				s.afterControl = finish
			} else {
				finish()
			}
		})
	}()
	return true
}

// showError returns true for any unsuccessful operation, including an expected
// cancellation, so completion callbacks never replace a stopped status.
func (s *screen) showError(err error) bool {
	if err == nil {
		return false
	}
	if !errors.Is(err, context.Canceled) && !errors.Is(err, errSourceChanged) {
		s.errorDialog(err)
	}
	return true
}

func (s *screen) discover() {
	if s.busy {
		return
	}
	s.invalidateExports()
	var exports []model.Export
	s.run("Discovering ready archives and exact sizes…", func(ctx context.Context, binding sourceSnapshot) error {
		var err error
		exports, err = binding.account.List(ctx)
		return err
	}, func(err error, binding sourceSnapshot) {
		if s.showError(err) {
			return
		}
		s.ready = exports
		s.discovery = binding
		s.discoveryBound = true
		s.selected = make(map[string]bool)
		s.exports.RemoveAll()
		if len(exports) == 0 {
			s.exports.Add(wrappedLabel("No ready archives in this source. Create an export, then discover again or start watching."))
			return
		}
		for _, export := range exports {
			var total int64
			for _, part := range export.Parts {
				total += part.Size
			}
			s.selected[export.ID] = true
			created := "date unavailable"
			if !export.CreatedAt.IsZero() {
				created = export.CreatedAt.Format("2006-01-02")
			}
			partCount := fmt.Sprintf("%d parts", len(export.Parts))
			if len(export.Parts) == 1 {
				partCount = "1 part"
			}
			check := widget.NewCheck(fmt.Sprintf("%s  ·  %s  ·  %.2f GB", created, partCount, float64(total)/1e9), func(value bool) {
				if !s.busy && binding.matches(s) {
					s.selected[export.ID] = value
				}
			})
			check.SetChecked(true)
			s.exports.Add(check)
			for _, part := range export.Parts {
				s.exports.Add(archiveLabel(fmt.Sprintf("%s  (%d bytes)", part.Filename, part.Size)))
			}
		}
		s.tabs.SelectIndex(0)
		s.status.SetText("Select exports, confirm the destination, and start downloading.")
	}, nil)
}

func (s *screen) requireDiscovery() bool {
	if s.discoveryBound && s.discovery.matches(s) {
		return true
	}
	s.invalidateExports()
	dialog.ShowInformation("Discover archives", "Discover archives for the current source and authorization before selecting downloads.", s.w)
	return false
}

func (s *screen) download() {
	if s.busy {
		return
	}
	// Apply source preferences before copying any URLs from a discovery.
	if err := s.save(); err != nil {
		s.errorDialog(err)
		return
	}
	if !s.requireDiscovery() {
		return
	}
	var parts []model.Part
	for _, export := range s.ready {
		if s.selected[export.ID] {
			parts = append(parts, export.Parts...)
		}
	}
	if len(parts) == 0 {
		dialog.ShowInformation("Select an export", "Discover archives and select at least one export.", s.w)
		return
	}
	opts := app.DownloadOptions{Destination: s.a.Settings.Destination, PartByPart: s.partByPart.Checked, StoragePassword: s.password.Text}
	s.run("Starting downloads…", func(ctx context.Context, _ sourceSnapshot) error { return s.a.Download(ctx, parts, opts) }, func(err error, _ sourceSnapshot) {
		if errors.Is(err, space.ErrInsufficientSpace) && !opts.PartByPart {
			dialog.ShowConfirm("Batch does not fit", "Free space or change the destination. Download one part at a time instead?", func(yes bool) {
				if yes && !s.busy {
					s.partByPart.SetChecked(true)
					s.download()
				}
			}, s.w)
		} else {
			s.showError(err)
		}
	}, func() {
		s.tabs.SelectIndex(1)
		s.pause.Enable()
		s.resume.Enable()
	})
}

// Pause can flush a network destination. Serialize pause/resume requests while
// keeping Stop available, and ignore a late control response from an older job.
func (s *screen) changeDownloadPause(paused bool) {
	if !s.busy || s.ctx.Err() != nil || (paused && s.pause.Disabled()) || (!paused && s.resume.Disabled()) {
		return
	}
	generation := s.jobGeneration
	s.controlPending = true
	s.pause.Disable()
	s.resume.Disable()
	if paused {
		s.status.SetText("Pausing downloads and saving partial files…")
	} else {
		s.banner.Hide()
		s.status.SetText("Rechecking space and resuming…")
	}
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if paused {
			s.a.Pause()
		} else {
			s.a.Resume()
		}
		fyne.Do(func() {
			s.controlPending = false
			if s.ctx.Err() != nil {
				s.afterControl = nil
				return
			}
			if s.afterControl != nil {
				finish := s.afterControl
				s.afterControl = nil
				finish()
				return
			}
			if s.ctx.Err() == nil && s.busy && s.jobGeneration == generation {
				s.pause.Enable()
				s.resume.Enable()
			}
		})
	}()
}

func archiveLabel(text string) *widget.Label {
	label := widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapBreak
	return label
}

func transferText(e events.Event) string {
	name := e.Filename
	if name == "" {
		name = "Archive part"
	}
	detail := e.Message
	if e.Total > 0 {
		detail = fmt.Sprintf("%.1f / %.1f MB  ·  %.1f MB/s", float64(e.Bytes)/1e6, float64(e.Total)/1e6, e.BytesPerSecond/1e6)
	}
	return name + "\n" + detail
}

func (s *screen) onEvent(e events.Event) {
	if e.Message != "" {
		s.status.SetText(e.Message)
	}
	if e.Kind == events.SpaceWarning || e.Kind == events.ReAuthRequired {
		s.banner.SetText(e.Message)
		s.banner.Show()
		s.content.Refresh()
	}
	if e.PartKey != "" {
		found := false
		for i, old := range s.transfers {
			if old.PartKey == e.PartKey {
				if e.Filename == "" {
					e.Filename = old.Filename
				}
				if e.Kind != events.Progress && e.Kind != events.PartCompleted {
					e.Bytes, e.Total, e.BytesPerSecond = old.Bytes, old.Total, old.BytesPerSecond
				}
				s.transfers[i] = e
				found = true
				break
			}
		}
		if !found {
			s.transfers = append(s.transfers, e)
		}
		s.transferList.Refresh()
	}
	var current, total int64
	for _, part := range s.transfers {
		current += part.Bytes
		total += part.Total
	}
	if total > 0 {
		s.progress.SetValue(float64(current) / float64(total))
	}
}
