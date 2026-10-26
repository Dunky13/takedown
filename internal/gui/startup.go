//go:build !headless

package gui

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/Dunky13/takedown/internal/config"
)

type startupView struct {
	s                *screen
	configure        func(context.Context, string, bool, bool) (string, error)
	result           *widget.Label
	buttons          map[string]*widget.Button
	savePassword     func(string, string, string) error
	forgetPassword   func(string, string) error
	credentialResult *widget.Label
}

func (s *screen) startupContent() fyne.CanvasObject {
	v := &startupView{s: s, configure: s.a.ConfigureStartup, buttons: make(map[string]*widget.Button), savePassword: config.SaveStoragePassword, forgetPassword: config.ForgetStoragePassword}
	s.startupUI = v
	v.result = archiveLabel("Startup has not been inspected. Preview the definition before installing. These controls manage only this user's login watcher.")
	for _, item := range []struct{ key, label string }{{"preview", "Preview definition"}, {"status", "Check startup status"}, {"install", "Install login watcher…"}, {"remove", "Remove login watcher…"}} {
		v.buttons[item.key] = widget.NewButton(item.label, func() { v.request(item.key) })
	}
	v.buttons["save-password"] = widget.NewButton("Save destination password…", func() { v.credentials(false) })
	v.buttons["forget-password"] = widget.NewButton("Forget destination password…", func() { v.credentials(true) })
	v.credentialResult = archiveLabel("No destination authorization action has been requested in this session.")
	copy := widget.NewButton("Copy startup result", func() { fyne.CurrentApp().Clipboard().SetContent(v.result.Text) })
	return container.NewBorder(container.NewVBox(wrappedLabel("A login watcher downloads ready archives using saved preferences. Installation starts it for your user and may download immediately. It does not sign in or create exports. OAuth and destination credentials must already be stored."), container.NewHBox(v.buttons["preview"], v.buttons["status"]), container.NewHBox(v.buttons["install"], v.buttons["remove"])), copy, nil, nil, container.NewVScroll(v.result))
}

func (v *startupView) credentialsContent() fyne.CanvasObject {
	return container.NewBorder(container.NewVBox(wrappedLabel("Use the masked Network storage password field in Destination, then save it here for unattended SMB or password-based SFTP downloads. Save stores it in the OS keyring, never preferences or the service definition. Forget removes stored authorization and keeps archives. Nothing is stored automatically."), container.NewHBox(v.buttons["save-password"], v.buttons["forget-password"])), nil, nil, nil, container.NewVScroll(v.credentialResult))
}

func (v *startupView) credentials(forget bool) {
	if v.s.busy || v.s.ctx.Err() != nil {
		return
	}
	if err := v.s.save(); err != nil {
		v.s.showError(err)
		return
	}
	settings := v.s.a.Settings
	destination := settings.Destination
	password := v.s.password.Text
	if !strings.HasPrefix(destination, "smb://") && !strings.HasPrefix(destination, "sftp://") {
		v.s.showError(errors.New("use an smb://user@host/share or sftp://user@host/path destination before managing a stored password"))
		return
	}
	if !forget && password == "" {
		v.s.showError(errors.New("enter the destination password in the masked field in Destination before saving it"))
		return
	}
	binding := sourceSnapshot{identity: settings.SourceIdentity(), revision: v.s.sourceRevision, account: v.s.a.Account}
	text := "Save the masked destination password in your OS keyring for the selected destination? Future login watches can use it without this window. The password is not written to preferences or the service definition."
	if forget {
		text = "Forget the selected destination's password from the OS keyring? Future automatic downloads may need credentials again. Existing archive files and partial downloads are kept. The masked session password is also cleared."
	}
	d := dialog.NewConfirm("Confirm destination authorization", text, func(yes bool) {
		if !yes || v.s.busy || v.s.ctx.Err() != nil {
			return
		}
		current, err := v.s.preferences()
		if err != nil || !binding.matches(v.s) || !reflect.DeepEqual(settings, current) || password != v.s.password.Text {
			v.s.status.SetText("Settings or password changed. Review destination authorization and confirm again.")
			return
		}
		v.s.run("Updating the confirmed destination authorization…", func(context.Context, sourceSnapshot) error {
			if forget {
				return v.forgetPassword(v.s.a.DataDir, destination)
			}
			return v.savePassword(v.s.a.DataDir, destination, password)
		}, func(err error, _ sourceSnapshot) {
			if v.s.showError(err) {
				return
			}
			if forget {
				v.s.password.SetText("")
				v.credentialResult.SetText("Destination password removed from the OS keyring and this session. Archive files were kept.")
			} else {
				v.credentialResult.SetText("Destination password saved in the OS keyring. Preferences and the service definition contain no password.")
			}
			v.s.status.SetText("Destination authorization updated.")
		}, nil)
	}, v.s.w)
	d.SetConfirmText("Confirm destination authorization")
	d.SetDismissText("Cancel")
	d.Show()
}

func (v *startupView) setBusy(busy bool) {
	for _, b := range v.buttons {
		if busy {
			b.Disable()
		} else {
			b.Enable()
		}
	}
}

func (v *startupView) request(action string) {
	if v.s.busy || v.s.ctx.Err() != nil {
		return
	}
	if action == "preview" || action == "status" {
		v.apply(action)
		return
	}
	// Bind confirmation to the same configuration and authorization revision as
	// other source mutations. Editing preferences requires a fresh confirmation.
	if err := v.s.save(); err != nil {
		v.s.showError(err)
		return
	}
	binding := sourceSnapshot{identity: v.s.a.Settings.SourceIdentity(), revision: v.s.sourceRevision, account: v.s.a.Account}
	settings := v.s.a.Settings
	partByPart := v.s.partByPart.Checked
	text := "Install and start this user's login watcher? It may begin downloading ready archives immediately and will watch again after login. The current source, destination and transfer preferences will be saved. This install action does not store the session password; use Destination access to save it explicitly if needed. OAuth authorization must already be stored. No Google sign-in or export request is started."
	if action == "remove" {
		text = "Stop and remove this user's login watcher? Automatic watching will stop. Existing archive files and partial downloads are kept."
	}
	d := dialog.NewConfirm("Confirm startup change", text, func(yes bool) {
		if !yes || v.s.busy || v.s.ctx.Err() != nil {
			return
		}
		current, err := v.s.preferences()
		if err != nil || !binding.matches(v.s) || !reflect.DeepEqual(settings, current) || partByPart != v.s.partByPart.Checked {
			v.s.status.SetText("Settings changed. Review startup settings and confirm again.")
			return
		}
		v.apply(action)
	}, v.s.w)
	d.SetConfirmText("Confirm startup change")
	d.SetDismissText("Cancel")
	d.Show()
}

func (v *startupView) apply(action string) {
	if v.s.busy || v.s.ctx.Err() != nil {
		return
	}
	partByPart := v.s.partByPart.Checked
	dryRun := action == "preview"
	if dryRun {
		action = "install"
	}
	if action == "install" && !dryRun {
		if err := v.s.adoptStartupRuntimePaths(); err != nil {
			v.s.showError(err)
			return
		}
	}
	var result string
	v.s.run("Checking or applying the startup action you selected…", func(ctx context.Context, _ sourceSnapshot) error {
		var err error
		result, err = v.configure(ctx, action, partByPart, dryRun)
		return err
	}, func(err error, _ sourceSnapshot) {
		if errors.Is(err, context.Canceled) {
			v.s.status.SetText("Startup action stopped. Check startup status before trying another change.")
			return
		}
		if v.s.showError(err) {
			return
		}
		v.result.SetText(result)
		v.s.status.SetText("Startup action completed. Review the result above.")
	}, nil)
}
