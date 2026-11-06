//go:build !headless

package gui

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"github.com/ncruces/zenity"
)

// fileRequest describes one OS file dialog: a folder unless jsonOnly or anyFile
// is set. anyFile shows hidden entries, so a key under ~/.ssh is reachable.
type fileRequest struct {
	// start is a hint, resolved to an existing folder off the UI thread; fallback is tried when it has none.
	title, start, fallback string
	jsonOnly, anyFile      bool
}

// osFileDialog blocks until the OS file dialog closes or ctx ends.
func osFileDialog(ctx context.Context, r fileRequest) (string, error) {
	options := []zenity.Option{zenity.Context(ctx), zenity.Title(r.title)}
	if r.start != "" {
		options = append(options, zenity.Filename(r.start+string(filepath.Separator)))
	}
	switch {
	case r.jsonOnly:
		options = append(options, zenity.FileFilter{Name: "JSON files", Patterns: []string{"*.json"}, CaseFold: true})
	case r.anyFile:
		options = append(options, zenity.ShowHidden())
	default:
		options = append(options, zenity.Directory())
	}
	return zenity.SelectFile(options...)
}

// existingFolder finds the nearest local folder to open a dialog in. Network
// destinations have none, and macOS rejects a start folder that does not exist.
func existingFolder(path string) string {
	path = strings.TrimSpace(path)
	if path == "" || strings.Contains(path, "://") || strings.HasPrefix(path, `\\`) {
		return ""
	}
	for {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			return path
		}
		parent := filepath.Dir(path)
		if parent == path {
			return ""
		}
		path = parent
	}
}

// pickFile runs the OS dialog off the UI thread, because it blocks until closed.
func (s *screen) pickFile(r fileRequest, chosen func(string)) {
	if s.busy || s.ctx.Err() != nil {
		return
	}
	if s.picking {
		// An OS dialog belongs to another process, so it cannot be raised from here.
		s.status.SetText("A file dialog is already open. Finish or close it first.")
		return
	}
	s.picking = true
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		// A stale network mount can block a stat, so this stays off the UI thread.
		if r.start = existingFolder(r.start); r.start == "" {
			r.start = existingFolder(r.fallback)
		}
		path, err := s.openFileDialog(s.ctx, r)
		fyne.Do(func() {
			s.picking = false
			if s.ctx.Err() != nil {
				return
			}
			s.w.RequestFocus()
			switch {
			case errors.Is(err, zenity.ErrCanceled):
			case errors.Is(err, exec.ErrNotFound):
				// Linux needs zenity, qarma or matedialog for a native dialog.
				s.fyneFileDialog(r, chosen)
			case err != nil:
				s.fileDialogFailed(err)
			case !s.busy:
				chosen(path)
			}
		})
	}()
}

// fileDialogFailed reports a dialog that could not open or complete. Typing the
// path into the field beside the button always works.
func (s *screen) fileDialogFailed(err error) {
	s.errorDialogAdvice(err, "The file dialog could not open. Type the path into the field instead.")
}

// locate opens a Fyne dialog in the requested folder when there is one.
func locate(r fileRequest, set func(fyne.ListableURI)) {
	if r.start == "" {
		return
	}
	if dir, err := storage.ListerForURI(storage.NewFileURI(r.start)); err == nil {
		set(dir)
	}
}

func (s *screen) fyneFileDialog(r fileRequest, chosen func(string)) {
	if !r.jsonOnly && !r.anyFile {
		picker := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				s.fileDialogFailed(err)
			} else if uri != nil && !s.busy {
				chosen(uri.Path())
			}
		}, s.w)
		locate(r, picker.SetLocation)
		picker.Show()
		return
	}
	picker := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			s.fileDialogFailed(err)
			return
		}
		if reader == nil {
			return
		}
		uri := reader.URI()
		if err := reader.Close(); err != nil {
			s.fileDialogFailed(err)
			return
		}
		if s.busy {
			return
		}
		if uri.Scheme() != "file" {
			s.fileDialogFailed(errors.New("select a local file"))
			return
		}
		chosen(uri.Path())
	}, s.w)
	if r.jsonOnly {
		picker.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
	}
	locate(r, picker.SetLocation)
	picker.Show()
}

func (s *screen) pickDestination() {
	s.pickFile(fileRequest{title: "Choose where to save archives", start: s.destination.Text}, s.destination.SetText)
}

func (s *screen) pickOAuthClient() {
	start := ""
	if client := strings.TrimSpace(s.oauthClient.Text); client != "" {
		start = filepath.Dir(client)
	}
	s.pickFile(fileRequest{title: "Select your OAuth desktop client JSON file", start: start, jsonOnly: true}, s.oauthClient.SetText)
}

// pickSSHFile chooses a private key or known_hosts file, opening beside the
// current value, else in ~/.ssh when that folder exists.
func (s *screen) pickSSHFile(title string, field *entry) {
	start, fallback := "", ""
	if text := expandHome(field.Text); filepath.IsAbs(text) {
		start = filepath.Dir(text)
	}
	if home, err := os.UserHomeDir(); err == nil {
		fallback = filepath.Join(home, ".ssh")
	}
	s.pickFile(fileRequest{title: title, start: start, fallback: fallback, anyFile: true}, field.SetText)
}
