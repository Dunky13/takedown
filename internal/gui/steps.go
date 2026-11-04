//go:build !headless

package gui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// pageTitle is the heading each guide step opens with.
func pageTitle(text string) *widget.Label {
	label := widget.NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	label.Wrapping = fyne.TextWrapWord
	label.SizeName = theme.SizeNameSubHeadingText
	return label
}

// stepBy moves along the four guide steps. It only navigates; nothing is saved or gated.
func (s *screen) stepBy(delta int) {
	for i, item := range s.guide {
		if item == s.tabs.Selected() && i+delta >= 0 && i+delta < len(s.guide) {
			s.tabs.Select(s.guide[i+delta])
			return
		}
	}
}

// syncGuide names the neighbouring steps on Back and Next, and hides them
// outside the guide or at its ends.
func (s *screen) syncGuide() {
	at := -1
	for i, item := range s.guide {
		if item == s.tabs.Selected() {
			at = i
		}
	}
	setNeighbour := func(button *widget.Button, prefix string, i int) {
		if at < 0 || i < 0 || i >= len(s.guide) {
			button.Hide()
			return
		}
		_, title, _ := strings.Cut(s.guide[i].Text, ". ")
		button.SetText(prefix + title)
		button.Show()
	}
	setNeighbour(s.back, "Back: ", at-1)
	setNeighbour(s.next, "Next: ", at+1)
	s.content.Refresh()
}

// destinationContent is step 3. Network and transfer settings stay folded away
// until a network destination is entered or they are needed.
func (s *screen) destinationContent(browse *widget.Button) fyne.CanvasObject {
	s.portabilityWarning = wrappedLabel("Data Portability archives must be stored encrypted. takedown does not encrypt or verify the destination, so choose an encrypted volume or protected share.")
	s.portabilityWarning.Importance = widget.WarningImportance
	network := widget.NewAccordionItem("Network storage (SMB / SFTP)", container.NewVBox(
		widget.NewForm(widget.NewFormItem("Network storage password", s.password)),
		wrappedLabel("Used only for SMB or SFTP storage. It is not a Google password and is not saved in preferences."),
		widget.NewForm(widget.NewFormItem("SFTP private key", s.sshKey), widget.NewFormItem("SFTP known_hosts", s.knownHosts), widget.NewFormItem("Free-space fallback", s.unknown)),
		wrappedLabel("Without free-space reporting, the space guard cannot protect the SFTP destination."),
	))
	network.Open = strings.Contains(s.destination.Text, "://")
	limits := widget.NewAccordionItem("Transfer limits", widget.NewForm(widget.NewFormItem("Free-space buffer (bytes)", s.buffer), widget.NewFormItem("Concurrent parts", s.workers)))
	s.storageDetails = widget.NewAccordion(network, limits)
	s.storageDetails.MultiOpen = true
	// As on step 2, only the optional settings scroll; the destination stays fixed.
	return container.NewBorder(container.NewVBox(
		pageTitle("Choose where to save"),
		wrappedLabel("Archives are written to this folder or network location. Nothing is deleted from Google."),
		container.NewBorder(nil, nil, widget.NewLabel("Save archives to"), browse, s.destination),
		wrappedLabel("A folder, Windows UNC path, smb://user@host/share/path or sftp://user@host/path"),
		s.portabilityWarning,
	), nil, nil, nil, container.NewVScroll(s.storageDetails))
}

func (s *screen) showWatchDialog() {
	size := s.w.Canvas().Size()
	s.watchDialog.Resize(fyne.NewSize(fyne.Min(620, size.Width-32), fyne.Min(420, size.Height-64)))
	s.watchDialog.Show()
}

// archivesContent is step 4: find ready archives, pick parts, download.
func (s *screen) archivesContent() fyne.CanvasObject {
	s.drivePrepare = container.NewVBox(
		widget.NewForm(widget.NewFormItem("Drive archive folder", container.NewBorder(nil, nil, nil, s.folderBrowse, s.driveFolder))),
		wrappedLabel("Browsing Drive needs sign-in (step 2). Archives appear once Google has delivered them."),
	)
	s.portabilityPrepare = container.NewVBox(s.create, s.sourceHelp)
	s.selectionSummary = widget.NewLabel("")
	selectAll := widget.NewButton("Select all", func() { s.selectAllParts(true) })
	selectNone := widget.NewButton("Select none", func() { s.selectAllParts(false) })
	selectAll.Importance, selectNone.Importance = widget.LowImportance, widget.LowImportance
	s.buttons = append(s.buttons, selectAll, selectNone)
	s.selectionRow = container.NewHBox(s.selectionSummary, layout.NewSpacer(), selectAll, selectNone)
	s.selectionRow.Hide()
	exportScroll := container.NewVScroll(s.exports)
	exportScroll.SetMinSize(fyne.NewSize(300, 100))
	// A dialog, because folding this into the page outgrows small screens.
	s.watchDialog = dialog.NewCustom("Keep watching for new archives", "Close", container.NewVBox(
		s.watchSelected,
		widget.NewForm(widget.NewFormItem("Check every (seconds)", s.pollInterval)),
		wrappedLabel("Default: 300 seconds. Minimum: 60 for Drive, 300 for Data Portability. Stop watching keeps partial files; Pause downloads keeps the watch running."),
		s.watch,
		wrappedLabel("Watching downloads ready archives from this source, including new exports. It does not create or schedule exports. Everything is included unless you turn on Only watch selected exports, and a partly ticked export is downloaded whole."),
	), s.w)
	keepWatching := widget.NewButton("Keep watching…", s.showWatchDialog)
	keepWatching.Importance = widget.LowImportance
	s.buttons = append(s.buttons, keepWatching)
	return container.NewBorder(
		container.NewVBox(
			pageTitle("Download your archives"),
			wrappedLabel("Find the archives that are ready, tick the files you want, then download."),
			container.NewStack(s.drivePrepare, s.portabilityPrepare),
			container.NewHBox(s.find),
		),
		container.NewVBox(s.selectionRow, container.NewHBox(s.downloadButton, s.partByPart, layout.NewSpacer(), keepWatching)),
		nil, nil, exportScroll,
	)
}

func setVisible(obj fyne.CanvasObject, shown bool) {
	if shown {
		obj.Show()
	} else {
		obj.Hide()
	}
}
