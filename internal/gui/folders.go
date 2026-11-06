//go:build !headless

package gui

import (
	"context"
	"errors"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Dunky13/takedown/internal/source"
)

type folderBrowser interface {
	BrowseFolders(context.Context, string) ([]source.Folder, error)
}
type accessInspector interface {
	CheckAccess(context.Context) (source.AccessInfo, error)
}

type folderPicker struct {
	s               *screen
	dialog          *dialog.CustomDialog
	binding         sourceSnapshot
	path            []source.Folder
	rows            *fyne.Container
	location        *widget.Label
	up, choose      *button
	closed, loading bool
}

func (s *screen) browseFolders() {
	if s.busy || s.selectedSource() != "drive" {
		return
	}
	if err := s.save(); err != nil {
		s.showError(err)
		return
	}
	if _, ok := s.a.Account.(folderBrowser); !ok {
		s.showError(errors.New("this source does not support folder browsing; paste the archive folder ID or link instead"))
		return
	}
	p := &folderPicker{s: s, binding: sourceSnapshot{identity: s.a.Settings.SourceIdentity(), revision: s.sourceRevision, account: s.a.Account}}
	s.folderUI = p
	p.rows = container.NewVBox()
	p.location = archiveLabel("My Drive")
	p.up = newButton("Up one folder", theme.NavigateBackIcon(), func() {
		if len(p.path) > 0 {
			p.path = p.path[:len(p.path)-1]
			p.load()
		}
	})
	p.choose = newButton("Use this archive folder", theme.ConfirmIcon(), func() {
		if p.loading || p.closed || s.busy || len(p.path) == 0 || !p.binding.matches(s) {
			return
		}
		folder := p.path[len(p.path)-1]
		p.dialog.Hide()
		s.driveFolder.SetText(folder.ID)
		s.status.SetText("Archive folder selected: " + folder.Name + ". Discover archives to verify its contents.")
	})
	p.choose.Importance = widget.HighImportance
	closeButton := newButton("Close", nil, func() { p.dialog.Hide() })
	content := border(container.NewVBox(flush(helperText("Browse explicitly, then choose the folder containing your Takeout archives. This selection does not download or scan unrelated folders.")), flush(p.location), container.NewHBox(p.up)), nil, nil, nil, scrolled(p.rows))
	p.dialog = s.dialogWithButtons("Choose Drive archive folder", content, closeButton, p.choose)
	p.dialog.Resize(fyne.NewSize(620, 380))
	p.dialog.SetOnClosed(func() {
		p.closed = true
		if p.loading && s.jobCancel != nil {
			s.jobCancel()
		}
	})
	p.dialog.Show()
	p.load()
}

func (p *folderPicker) load() {
	if p.closed || p.loading || p.s.busy || !p.binding.matches(p.s) {
		return
	}
	parent := ""
	names := []string{"My Drive"}
	for _, folder := range p.path {
		parent = folder.ID
		names = append(names, folder.Name)
	}
	p.location.SetText(strings.Join(names, " / "))
	p.rows.RemoveAll()
	p.rows.Add(wrappedLabel("Loading folders…"))
	p.up.Disable()
	p.choose.Disable()
	var folders []source.Folder
	p.loading = true
	started := p.s.run("Reading folders you requested…", func(ctx context.Context, b sourceSnapshot) error {
		var err error
		folders, err = b.account.(folderBrowser).BrowseFolders(ctx, parent)
		return err
	}, func(err error, _ sourceSnapshot) {
		p.loading = false
		if p.closed {
			return
		}
		if errors.Is(err, errSourceChanged) {
			p.dialog.Hide()
			return
		}
		p.rows.RemoveAll()
		if err != nil {
			p.rows.Add(archiveLabel("Could not list folders. " + err.Error() + "\nClose this picker and check authorization or paste the folder ID into the Drive archive folder field."))
		} else {
			for _, folder := range folders {
				// A leading, low-importance button reads as a list row and opens the folder.
				row := newButton(folder.Name, theme.FolderIcon(), func() {
					if !p.s.busy && !p.closed {
						p.path = append(p.path, folder)
						p.load()
					}
				})
				row.Importance, row.Alignment = widget.LowImportance, widget.ButtonAlignLeading
				p.rows.Add(row)
			}
			if len(folders) == 0 {
				p.rows.Add(wrappedLabel("No child folders here. You can choose this folder or go up."))
			}
		}
		if len(p.path) > 0 {
			p.up.Enable()
			if err == nil {
				p.choose.Enable()
			}
		}
		p.s.status.SetText("Choose an archive folder, or close the picker to keep the existing selection.")
	}, nil)
	if !started {
		p.loading = false
		p.dialog.Hide()
	}
}

func (s *screen) inspectAccess() {
	var info source.AccessInfo
	s.run("Checking the current Google grant without starting sign-in…", func(ctx context.Context, b sourceSnapshot) error {
		inspector, ok := b.account.(accessInspector)
		if !ok {
			return errors.New("this source cannot inspect authorized access; configured settings do not verify identity or consent")
		}
		var err error
		info, err = inspector.CheckAccess(ctx)
		return err
	}, func(err error, _ sourceSnapshot) {
		if err != nil {
			if !errors.Is(err, errSourceChanged) && !errors.Is(err, context.Canceled) {
				s.setAccessStatus("Access could not be verified. Check authorization and diagnostics; configured choices remain unverified.", widget.WarningImportance)
			}
			s.showError(err)
			return
		}
		one, timed := strings.Join(info.OneTimeResources, ", "), strings.Join(info.TimeBasedResources, ", ")
		if one == "" {
			one = "none reported"
		}
		if timed == "" {
			timed = "none reported"
		}
		s.setAccessStatus("Authorized resource groups reported by Google:\nOne-time: "+one+"\nTime-based: "+timed+"\nThis does not verify the identity of the account email hint.", widget.LowImportance)
		s.status.SetText("Access check finished without starting an OAuth flow.")
	}, nil)
}
