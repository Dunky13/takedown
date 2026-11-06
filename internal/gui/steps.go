//go:build !headless

package gui

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// stepBy moves along the four guide steps. It only navigates; nothing is saved or gated.
func (s *screen) stepBy(delta int) {
	for i, item := range s.guide {
		if item == s.tabs.Selected() && i+delta >= 0 && i+delta < len(s.guide) {
			s.tabs.Select(s.guide[i+delta])
			return
		}
	}
}

// syncGuide names the next step on Next, hides Back and Next outside the guide
// or at its ends, and moves the stepper and toolbar highlights.
func (s *screen) syncGuide() {
	at := -1
	for i, item := range s.guide {
		if item == s.tabs.Selected() {
			at = i
		}
	}
	if at < 0 || at == 0 {
		s.back.Hide()
	} else {
		s.back.Show()
	}
	if at < 0 || at+1 >= len(s.guide) {
		s.next.Hide()
	} else {
		s.next.SetText("Next: " + s.guide[at+1].Text)
		// Authorize is the one primary on Sign in.
		s.next.Importance = widget.HighImportance
		if s.guide[at] == s.authorizeTab {
			s.next.Importance = widget.MediumImportance
		}
		s.next.Refresh()
		s.next.Show()
	}
	s.stepper.setCurrent(at)
	setVisible(s.savePrefs, at >= 0)
	s.syncToolbar()
	relayout(s.content)
}

// syncToolbar marks the open tool page and mirrors which pages are usable.
func (s *screen) syncToolbar() {
	for item, button := range s.toolButtons {
		button.Importance = widget.LowImportance
		if item == s.tabs.Selected() {
			button.Importance = widget.MediumImportance
		}
		if s.tabs.Disabled(item) {
			button.Disable()
		} else {
			button.Enable()
		}
		button.Refresh()
	}
}

// layoutContent is the window: a header with the tool pages at the top right,
// the progress stepper, one centred column for the page, and the footer.
func (s *screen) layoutContent(save *button) fyne.CanvasObject {
	s.toolButtons = map[*container.TabItem]*button{}
	var tools []fyne.CanvasObject
	for _, tool := range []struct {
		item *container.TabItem
		icon fyne.Resource
	}{{s.transfersTab, theme.DownloadIcon()}, {s.jobsTab, theme.ListIcon()}, {s.diagnosticsTab, theme.InfoIcon()}} {
		button := newButton(tool.item.Text, tool.icon, func() { s.tabs.Select(tool.item) })
		s.toolButtons[tool.item] = button
		tools = append(tools, button)
	}
	s.sourceSummary.Importance = widget.LowImportance
	s.sourceSummary.SizeName = theme.SizeNameCaptionText
	title := widget.NewLabelWithStyle("takedown", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	s.header = border(nil, nil, title, container.NewHBox(tools...), container.NewVBox(layout.NewSpacer(), s.sourceSummary, layout.NewSpacer()))
	padded := func(o fyne.CanvasObject, vertical float32) fyne.CanvasObject {
		return container.New(layout.NewCustomPaddedLayout(vertical, vertical, 16, 16), o)
	}
	// A hidden banner must not leave its margin behind.
	top := newSurface(container.NewVBox(padded(s.header, focusRingSpace+focusRingWidth), container.New(maxWidth{max: pageWidth, margin: 16, top: 4}, s.banner), container.New(maxWidth{max: pageWidth, margin: 16}, s.stepper), widget.NewSeparator()), theme.ColorNameHeaderBackground)
	// Buttons keep their natural height at the top. The strut holds the footer's
	// height steady on tool pages, where the buttons are hidden.
	strut := canvas.NewRectangle(color.Transparent)
	strut.SetMinSize(fyne.NewSize(1, s.next.MinSize().Height))
	footer := border(nil, nil, strut, container.NewVBox(container.NewHBox(save, s.cancelButton, s.back, s.next)), container.New(layout.NewCustomPaddedLayout(theme.Padding()*2, 0, 0, 0), s.status))
	bottom := newSurface(container.NewVBox(widget.NewSeparator(), padded(footer, 4)), theme.ColorNameHeaderBackground)
	return border(top, bottom, nil, nil, container.New(maxWidth{max: pageWidth, margin: 20, top: 12}, s.tabs.stack))
}

// pageWidth keeps text at a comfortable line length on wide windows.
const pageWidth = 760

// transfersContent lists downloads, with overall progress and the controls on top.
func (s *screen) transfersContent() fyne.CanvasObject {
	s.transferEmpty = hint(helperText("No downloads yet. Pick archive files on the Download step."))
	// The bar sits in its own fixed-height row; the percentage is a caption
	// beside it, because text straddling the fill reads poorly.
	s.progress.TextFormatter = func() string { return "" }
	s.progressText = widget.NewLabel("0%")
	s.pauseRow = container.NewHBox(s.pause, s.resume)
	s.transferControls = spaced(border(nil, nil, nil, s.progressText, container.New(slimBar{height: 10}, s.progress)), s.pauseRow)
	s.syncTransfers()
	return border(spaced(pageHeader("Transfers", ""), s.transferControls), nil, nil, nil, container.NewStack(s.transferList, s.transferEmpty))
}

// syncTransfers shows the hint, or the progress, never both. Pause and Resume
// show only while a download or watch runs; other work, such as sign-in, does not count.
func (s *screen) syncTransfers() {
	if s.transferEmpty == nil {
		return
	}
	active := len(s.transfers) > 0 || s.transferring
	setVisible(s.transferEmpty, !active)
	setVisible(s.transferControls, active)
	setVisible(s.pauseRow, s.transferring)
	s.transferEmpty.Refresh()
}

// destinationContent is step 3. Network and transfer settings stay folded away
// until a network destination is entered or they are needed.
func (s *screen) destinationContent(browse *button) fyne.CanvasObject {
	s.portabilityWarning = newCallout("Data Portability archives must be stored encrypted. takedown does not encrypt or verify the destination, so choose an encrypted volume or protected share.")
	keyBrowse := newButton("Browse…", theme.FileIcon(), func() { s.pickSSHFile("Select your SFTP private key", s.sshKey) })
	hostsBrowse := newButton("Browse…", theme.FileIcon(), func() { s.pickSSHFile("Select your known_hosts file", s.knownHosts) })
	s.buttons = append(s.buttons, keyBrowse, hostsBrowse)
	spaceWarning := newCallout("Without free-space reporting, the space guard cannot protect the SFTP destination.")
	s.unknown.OnChanged = func(on bool) {
		setVisible(spaceWarning, on)
		s.storageDetails.Refresh()
	}
	setVisible(spaceWarning, s.unknown.Checked)
	network := widget.NewAccordionItem("Network storage (SMB / SFTP)", container.NewVBox(
		field("Network storage password", s.password, "Used only for SMB or SFTP storage. It is not a Google password and is not saved in preferences."),
		field("SFTP private key", border(nil, nil, nil, keyBrowse, s.sshKey), "An unencrypted key file, or one already unlocked. Leave it empty to sign in with the password above."),
		field("SFTP known_hosts", border(nil, nil, nil, hostsBrowse, s.knownHosts), "The file that lists trusted server keys. Leave it empty to use ~/.ssh/known_hosts. A path starting with ~/ is expanded."),
		field("Free-space fallback", s.unknown, ""),
		spaceWarning,
	))
	network.Open = strings.Contains(s.destination.Text, "://")
	limits := widget.NewAccordionItem("Transfer limits", container.NewVBox(
		field("Free-space buffer (bytes)", s.buffer, ""),
		field("Concurrent parts", s.workers, ""),
	))
	s.storageDetails = widget.NewAccordion(network, limits)
	s.storageDetails.MultiOpen = true
	// As on step 2, only the optional settings scroll; the destination stays fixed.
	return border(spaced(
		pageHeader("Choose where to save", "Archives are written to this folder or network location. Nothing is deleted from Google."),
		field("Save archives to", border(nil, nil, nil, browse, s.destination), "A folder, Windows UNC path, smb://user@host/share/path or sftp://user@host/path"),
		s.portabilityWarning,
	), nil, nil, nil, scrolled(s.storageDetails))
}

func (s *screen) showWatchDialog() {
	size := s.w.Canvas().Size()
	s.watchDialog.Resize(fyne.NewSize(fyne.Min(620, size.Width-32), fyne.Min(420, size.Height-64)))
	s.watchDialog.Show()
}

// archivesContent is step 4: find ready archives, pick parts, download.
func (s *screen) archivesContent() fyne.CanvasObject {
	s.drivePrepare = field("Drive archive folder", border(nil, nil, nil, s.folderBrowse, s.driveFolder), "Browsing Drive needs sign-in (step 2). Archives appear once Google has delivered them.")
	s.portabilityPrepare = container.NewVBox(container.NewHBox(s.create), s.sourceHelp)
	s.selectionSummary = widget.NewLabel("")
	selectAll := newButton("Select all", nil, func() { s.selectAllParts(true) })
	selectNone := newButton("Select none", nil, func() { s.selectAllParts(false) })
	selectAll.Importance, selectNone.Importance = widget.LowImportance, widget.LowImportance
	s.buttons = append(s.buttons, selectAll, selectNone)
	s.selectionRow = border(nil, nil, nil, container.NewHBox(selectAll, selectNone), s.selectionSummary)
	s.selectionRow.Hide()
	exportScroll := container.NewVScroll(s.exports)
	exportScroll.SetMinSize(fyne.NewSize(300, 40)) // small, so a banner still fits 640 px; it grows with the window
	// A dialog, because folding this into the page outgrows small screens.
	closeWatch := newButton("Close", nil, func() { s.watchDialog.Hide() })
	s.watchDialog = s.dialogWithButtons("Keep watching for new archives", container.NewVBox(
		s.watchSelected,
		field("Check every (seconds)", s.pollInterval, "Default: 300 seconds. Minimum: 60 for Drive, 300 for Data Portability. Stop watching keeps partial files; Pause downloads keeps the watch running."),
		flush(helperText("Watching downloads ready archives from this source, including new exports. It does not create or schedule exports. Everything is included unless you turn on Only watch selected exports.")),
		newCallout("A partly ticked export is downloaded whole while watching."),
	), closeWatch, s.watch)
	keepWatching := newButton("Keep watching…", nil, s.showWatchDialog)
	s.buttons = append(s.buttons, keepWatching)
	s.find.Importance = widget.HighImportance // until something is discovered
	// The part-by-part choice also drives Watch and the login watcher, so it never hides.
	s.downloadRow = container.NewHBox(s.downloadButton)
	s.downloadRow.Hide()
	return border(
		spaced(
			pageHeader("Download your archives", "Find the archives that are ready, tick the files you want, then download."),
			container.NewStack(s.drivePrepare, s.portabilityPrepare),
			container.NewHBox(s.find, keepWatching, layout.NewSpacer(), s.partByPart),
			gap(0),
		),
		container.NewVBox(gap(theme.Padding()), s.downloadRow, gap(theme.Padding())), nil, nil, newPanel(border(s.selectionRow, nil, nil, nil, exportScroll)),
	)
}

func setVisible(obj fyne.CanvasObject, shown bool) {
	if shown {
		obj.Show()
	} else {
		obj.Hide()
	}
}
