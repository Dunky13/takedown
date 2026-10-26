//go:build !headless

package gui

import (
	"context"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/Dunky13/takedown/internal/app"
)

type diagnosticsView struct {
	s               *screen
	online, storage *widget.Check
	run, copy       *widget.Button
	result          *widget.Label
	diagnose        func(context.Context, app.DiagnosticOptions) (app.DiagnosticReport, error)
}

func (s *screen) diagnosticsContent() fyne.CanvasObject {
	v := &diagnosticsView{s: s, diagnose: s.a.Diagnose}
	s.diagnosticsUI = v
	v.online = widget.NewCheck("Check Google access online (no sign-in prompt)", nil)
	v.storage = widget.NewCheck("Test destination access (may create the selected directory)", nil)
	v.result = archiveLabel("Diagnostics have not run. Start with local configuration checks. Optional online and storage checks run only when selected; they never start OAuth or create exports.")
	v.run = widget.NewButton("Run diagnostics", func() {
		opts := app.DiagnosticOptions{Online: v.online.Checked, Storage: v.storage.Checked, StoragePassword: s.password.Text}
		var report app.DiagnosticReport
		s.run("Running the diagnostics you selected…", func(ctx context.Context, _ sourceSnapshot) error {
			var err error
			report, err = v.diagnose(ctx, opts)
			return err
		}, func(err error, _ sourceSnapshot) {
			if err != nil {
				s.showError(err)
				return
			}
			var lines []string
			for _, check := range report.Checks {
				lines = append(lines, fmt.Sprintf("%s · %s\n%s", strings.ToUpper(check.Status), check.Name, check.Message))
			}
			if len(lines) == 0 {
				lines = append(lines, "No diagnostic checks were reported.")
			}
			v.result.SetText(strings.Join(lines, "\n\n"))
			s.status.SetText("Diagnostics completed. Review each check and its suggested action.")
		}, nil)
	})
	v.copy = widget.NewButton("Copy diagnostic results", func() { fyne.CurrentApp().Clipboard().SetContent(v.result.Text) })
	checks := container.NewBorder(container.NewVBox(wrappedLabel("Local checks inspect configuration and stored authorization. Online checks use existing Google access. Storage checks use this session's password and may create the selected directory, but no archive files. No check starts sign-in, creates an export, or verifies destination encryption."), v.online, v.storage, container.NewHBox(v.run, v.copy)), nil, nil, nil, container.NewVScroll(v.result))
	startup := s.startupContent()
	return container.NewAppTabs(container.NewTabItem("Checks", checks), container.NewTabItem("Startup", startup), container.NewTabItem("Destination access", s.startupUI.credentialsContent()))
}

func (v *diagnosticsView) setBusy(busy bool) {
	if busy {
		v.online.Disable()
		v.storage.Disable()
		v.run.Disable()
	} else {
		v.online.Enable()
		v.storage.Enable()
		v.run.Enable()
	}
}

func (s *screen) errorDialog(err error) {
	text := err.Error() + "\n\nCheck the selected source and authorization. Diagnostics can check configuration, Google access, or the destination without starting sign-in. Retry mutations only after reviewing recorded job history."
	copy := widget.NewButton("Copy error details", func() { fyne.CurrentApp().Clipboard().SetContent(err.Error()) })
	content := container.NewBorder(nil, copy, nil, nil, container.NewVScroll(archiveLabel(text)))
	d := dialog.NewCustom("Action could not complete", "Close", content, s.w)
	d.Resize(fyne.NewSize(600, 340))
	d.Show()
}
