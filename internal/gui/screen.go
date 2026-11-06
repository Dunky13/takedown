//go:build !headless

package gui

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
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
	destination, email, password, buffer, sshKey, knownHosts   *entry
	oauthClient, driveFolder, resources, pollInterval          *entry
	source                                                     *widget.RadioGroup
	workers                                                    *widget.Select
	partByPart, unknown, watchSelected                         *widget.Check
	status, sourceSummary, sourceHelp, disclosure              *widget.Label
	signInHelp, selectionSummary                               *widget.Label
	banner, portabilityWarning                                 *callout
	driveCard, portabilityCard                                 *choiceCard
	disclosureAccordion                                        *widget.Accordion
	stepper                                                    *stepper
	toolButtons                                                map[*container.TabItem]*button
	header                                                     *fyne.Container
	downloadRow                                                *fyne.Container
	driveDocLink, portabilityDocLink                           *widget.Hyperlink
	driveLink, portabilityLink                                 fyne.CanvasObject // the flush wrappers that show and hide
	savePrefs                                                  *button
	transferEmpty                                              fyne.CanvasObject
	transferControls, pauseRow                                 fyne.CanvasObject
	progressText                                               *widget.Label
	progress                                                   *widget.ProgressBar
	exports                                                    *fyne.Container
	transferList                                               *widget.List
	transfers                                                  []events.Event
	selected, skippedParts, doneParts                          map[string]bool
	selectionChecks                                            []*widget.Check
	setAllParts                                                []func(bool)
	ready                                                      []model.Export
	buttons                                                    []*button
	authorize, disconnect, create, find, watch, downloadButton *button
	pause, resume, cancelButton                                *button
	openTakeout, back, next                                    *button
	driveOptions, portabilityOptions                           fyne.CanvasObject
	drivePrepare, portabilityPrepare                           fyne.CanvasObject
	storageDetails                                             *widget.Accordion
	selectionRow                                               *fyne.Container
	watchDialog                                                *dialog.CustomDialog
	jobCancel                                                  context.CancelFunc
	jobDone                                                    <-chan struct{}
	controlPending                                             bool
	afterControl                                               func()
	busy, watching, picking, transferring                      bool
	openFileDialog                                             func(context.Context, fileRequest) (string, error)
	stopEvents                                                 func()
	tabs                                                       *navigator
	sourceTab, authorizeTab, destinationTab, archivesTab       *container.TabItem
	transfersTab, jobsTab, diagnosticsTab                      *container.TabItem
	guide                                                      []*container.TabItem
	jobGeneration                                              uint64
	sourceRevision                                             uint64
	discovery                                                  sourceSnapshot
	discoveryBound                                             bool
	jobsUI                                                     *jobsView
	resourceUI                                                 *resourceChooser
	folderBrowse                                               *button
	folderUI                                                   *folderPicker
	accessStatus                                               *widget.Label
	checkAccess                                                *button
	diagnosticsUI                                              *diagnosticsView
	diagnosticsTabs                                            *subTabs
	startupUI                                                  *startupView
}

func newScreen(parent context.Context, desktop fyne.App, w fyne.Window, a *app.App) *screen {
	ctx, cancel := context.WithCancel(parent)
	s := &screen{ctx: ctx, cancel: cancel, w: w, a: a, selected: make(map[string]bool), skippedParts: make(map[string]bool), openFileDialog: osFileDialog}
	s.destination = newEntry()
	s.destination.SetText(a.Settings.Destination)
	s.destination.SetPlaceHolder("Folder, Windows UNC, smb://user@host/share/path or sftp://user@host/path")
	s.password = newPasswordEntry()
	s.password.SetPlaceHolder("Storage password; kept in memory for this session")
	s.buffer = newEntry()
	s.buffer.SetText(strconv.FormatUint(a.Settings.SafetyBuffer, 10))
	s.workers = widget.NewSelect([]string{"1", "2", "3", "4", "8", "16"}, nil)
	s.workers.SetSelected(strconv.Itoa(a.Settings.Concurrency))
	if s.workers.Selected == "" {
		s.workers.Options = append(s.workers.Options, strconv.Itoa(a.Settings.Concurrency))
		s.workers.SetSelected(strconv.Itoa(a.Settings.Concurrency))
	}
	s.sshKey = newEntry()
	s.sshKey.SetPlaceHolder("Path to a private key, e.g. ~/.ssh/id_ed25519")
	s.sshKey.SetText(a.Settings.PrivateKeyPath)
	s.knownHosts = newEntry()
	s.knownHosts.SetPlaceHolder("Default: ~/.ssh/known_hosts")
	s.knownHosts.SetText(a.Settings.KnownHostsPath)
	s.partByPart = widget.NewCheck("Download one part at a time", nil)
	s.unknown = widget.NewCheck("Allow SFTP without free-space reporting", nil)
	s.unknown.SetChecked(a.Settings.AllowUnknownSpace)
	s.status = helperText("Follow the steps above: choose data, sign in, pick a destination, then download.")
	// One line, so the footer keeps its height. The banner and error dialog hold full text.
	s.status.Wrapping, s.status.Truncation = fyne.TextWrapOff, fyne.TextTruncateEllipsis
	s.banner = newCallout("")
	s.banner.Hide()
	s.progress = widget.NewProgressBar()
	s.exports = container.NewVBox(helperText("No archives discovered. Set up your source, then choose Discover archives."))
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
	s.authorize = newButton("Authorize in browser", theme.LoginIcon(), s.login)
	s.authorize.Importance = widget.HighImportance
	s.disconnect = newButton("Remove local token", nil, s.forget)
	s.find = newButton("Discover archives", theme.ViewRefreshIcon(), s.discover)
	s.downloadButton = newButton("Download selected", theme.DownloadIcon(), s.download)
	s.downloadButton.Importance = widget.HighImportance
	s.watch = newButton("Start watching", theme.MediaPlayIcon(), s.startWatch)
	s.watch.Importance = widget.HighImportance // the one primary in its dialog
	s.watchSelected = widget.NewCheck("Only watch selected exports", nil)
	save := newButton("Save preferences", nil, func() {
		s.run("Saving preferences…", func(context.Context, sourceSnapshot) error { return nil }, func(err error, _ sourceSnapshot) {
			if !s.showError(err) {
				s.status.SetText("Preferences saved.")
			}
		}, nil)
	})
	save.Importance = widget.LowImportance
	s.savePrefs = save
	folder := newButton("Browse…", theme.FolderOpenIcon(), s.pickDestination)
	s.pause = newButton("Pause downloads", theme.MediaPauseIcon(), func() { s.changeDownloadPause(true) })
	s.resume = newButton("Resume / retry space", theme.MediaPlayIcon(), func() { s.changeDownloadPause(false) })
	s.cancelButton = newButton("Stop", theme.MediaStopIcon(), func() {
		if s.jobCancel != nil {
			s.jobCancel()
			s.status.SetText("Stopping; partial files will be kept.")
		}
	})
	s.buttons = []*button{s.authorize, s.find, s.disconnect, save, folder, s.downloadButton, s.watch}
	s.sourceTab = container.NewTabItem("Choose data", s.sourceContent())
	s.authorizeTab = container.NewTabItem("Sign in", s.authorizeContent())
	s.destinationTab = container.NewTabItem("Save to", s.destinationContent(folder))
	s.archivesTab = container.NewTabItem("Download", s.archivesContent())
	s.transfersTab = container.NewTabItem("Transfers", s.transfersContent())
	s.jobsTab = container.NewTabItem("Jobs", s.jobsContent())
	s.diagnosticsTab = container.NewTabItem("Diagnostics", s.diagnosticsContent())
	s.guide = []*container.TabItem{s.sourceTab, s.authorizeTab, s.destinationTab, s.archivesTab}
	s.back = newButton("Back", nil, func() { s.stepBy(-1) })
	s.back.Importance = widget.LowImportance
	s.next = newButton("", theme.NavigateNextIcon(), func() { s.stepBy(1) })
	s.next.IconPlacement = widget.ButtonIconTrailingText
	s.tabs = newNavigator(s.sourceTab, s.authorizeTab, s.destinationTab, s.archivesTab, s.transfersTab, s.jobsTab, s.diagnosticsTab)
	// The title sits outside Diagnostics' own tabs.
	s.tabs.setPage(s.diagnosticsTab, border(pageHeader("Diagnostics", ""), nil, nil, nil, s.diagnosticsTab.Content))
	s.tabs.OnSelected = func(*container.TabItem) { s.syncGuide() }
	s.tabs.onEnabledChanged = s.syncToolbar
	s.stepper = newStepper(s.guide, s.tabs.Select)
	s.content = s.layoutContent(save)
	w.SetPadded(false) // chrome and separators run edge to edge
	s.updateSourceLabels()
	s.destination.OnChanged = func(text string) {
		s.updateSummary()
		if strings.Contains(text, "://") {
			s.storageDetails.Open(0)
		}
		// What is already downloaded depends on the destination, so re-derive the ticks.
		if s.discoveryBound && !s.busy {
			s.showExports(s.discovery)
		}
	}
	s.setBusy(false)
	s.bindSourceChanges()
	if a.Settings.OAuthClientFile == "" {
		s.tabs.Select(s.sourceTab)
	} else {
		s.tabs.Select(s.archivesTab)
	}
	s.syncGuide()
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
	s.loadSavedTransfers()
	return s
}

// loadSavedTransfers lists the files an earlier session saved, with their progress.
func (s *screen) loadSavedTransfers() {
	records, err := state.ReadSnapshot(filepath.Join(s.a.DataDir, "takedown-state.json"))
	if err != nil {
		return
	}
	for _, r := range records {
		message := "Partial file saved"
		if r.Complete {
			message = "Archive complete"
		}
		s.transfers = append(s.transfers, events.Event{Kind: events.Status, PartKey: r.Part.Key(), Filename: r.Part.Filename, Bytes: r.Offset, Total: r.Part.Size, Message: message})
	}
	s.updateProgress()
	s.syncTransfers()
}

func (s *screen) close() { s.cancel(); s.stopEvents(); s.wg.Wait() }

func (s *screen) invalidateExports() {
	s.ready = nil
	s.selected = make(map[string]bool)
	s.discoveryBound = false
	s.discovery = sourceSnapshot{}
	s.skippedParts = make(map[string]bool)
	s.doneParts = nil
	s.selectionChecks = nil
	s.setAllParts = nil
	s.exports.RemoveAll()
	s.exports.Add(helperText("Discover archives for the current source and authorization before downloading."))
	s.updateSelectionSummary()
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
	for _, e := range []*entry{s.destination, s.email, s.password, s.buffer, s.sshKey, s.knownHosts, s.oauthClient, s.driveFolder, s.resources, s.pollInterval} {
		if busy {
			e.Disable()
		} else {
			e.Enable()
		}
	}
	for _, choice := range []fyne.Disableable{s.source, s.driveCard, s.portabilityCard, s.workers} {
		if busy {
			choice.Disable()
		} else {
			choice.Enable()
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
	for _, check := range s.selectionChecks {
		if busy {
			check.Disable()
		} else {
			check.Enable()
		}
	}
	if !busy {
		s.transferring = false
	}
	s.syncTransfers()
	if busy && s.w != nil {
		s.w.Canvas().Unfocus() // a clicked button that starts work must not stay highlighted
	}
	if busy {
		s.cancelButton.Show()
		s.cancelButton.Enable()
	} else {
		s.cancelButton.Hide()
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
					s.notice("Source changed. Authorize if needed, then discover archives again.")
				} else if errors.Is(err, context.Canceled) {
					s.status.SetText("Stopped. Partial files are available for the next run.")
				} else if err != nil {
					// The banner and the dialog hold the full text; the footer stays short.
					s.notice(err.Error())
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

// bannerLimit keeps a long error from pushing the page out of the window; the
// error dialog holds the full text.
const bannerLimit = 200

// clipText collapses whitespace and cuts at a word boundary with an ellipsis.
func clipText(text string, limit int) string {
	text = strings.Join(strings.Fields(text), " ")
	runes := []rune(text)
	if len(runes) <= limit {
		return text
	}
	cut := string(runes[:limit])
	if i := strings.LastIndex(cut, " "); i > limit/2 {
		cut = cut[:i]
	}
	return cut + "…"
}

// notice puts a safety message in the banner, where it cannot be missed. The
// footer only points at it, so the text is not shown twice.
func (s *screen) notice(text string) {
	s.status.SetText("Needs your attention. See the message above.")
	s.banner.SetText(clipText(text, bannerLimit))
	s.banner.Show()
	if s.content != nil {
		s.content.Refresh()
	}
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
		s.showExports(binding)
		if len(exports) == 0 {
			return
		}
		s.tabs.Select(s.archivesTab)
		s.status.SetText("Choose the archive files to download, then select Download selected.")
	}, nil)
}

func (s *screen) requireDiscovery() bool {
	if s.discoveryBound && s.discovery.matches(s) {
		return true
	}
	s.invalidateExports()
	s.inform("Discover archives", "Discover archives for the current source and authorization before selecting downloads.")
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
	parts := s.selectedParts()
	if len(parts) == 0 {
		s.inform("Select an export", "Discover archives and select at least one archive file.")
		return
	}
	opts := app.DownloadOptions{Destination: s.a.Settings.Destination, PartByPart: s.partByPart.Checked, StoragePassword: s.password.Text}
	s.run("Starting downloads…", func(ctx context.Context, _ sourceSnapshot) error { return s.a.Download(ctx, parts, opts) }, func(err error, _ sourceSnapshot) {
		if errors.Is(err, space.ErrInsufficientSpace) && !opts.PartByPart {
			s.confirmText("Batch does not fit", "Free space or change the destination. Download one part at a time instead?", "Download one at a time", func(yes bool) {
				if yes && !s.busy {
					s.partByPart.SetChecked(true)
					s.download()
				}
			})
		} else {
			s.showError(err)
		}
	}, s.startTransfers)
}

// startTransfers opens the Transfers page for a running download or watch.
func (s *screen) startTransfers() {
	s.transferring = true
	s.tabs.Select(s.transfersTab)
	s.pause.Enable()
	s.resume.Enable()
	s.syncTransfers()
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
	label.Wrapping = fyne.TextWrapWord
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
	if e.Message != "" && (e.Kind == events.SpaceWarning || e.Kind == events.ReAuthRequired || e.Kind == events.Warning || e.Kind == events.Error) {
		s.notice(e.Message)
	} else if e.Message != "" {
		s.status.SetText(e.Message)
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
		s.syncTransfers()
	}
	s.updateProgress()
}

// updateProgress shows the bytes saved so far across every listed transfer.
func (s *screen) updateProgress() {
	var current, total int64
	for _, part := range s.transfers {
		current += part.Bytes
		total += part.Total
	}
	if total > 0 {
		s.progress.SetValue(float64(current) / float64(total))
		s.progressText.SetText(fmt.Sprintf("%d%%", int(100*current/total)))
	}
}
