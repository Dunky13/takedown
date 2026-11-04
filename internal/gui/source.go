//go:build !headless

package gui

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Dunky13/takedown/internal/app"
)

const (
	driveSourceLabel       = "Takeout in Google Drive"
	portabilitySourceLabel = "Data Portability API"
)

// sourceContent is step 1. It also creates the fields that later steps share.
func (s *screen) sourceContent() fyne.CanvasObject {
	s.source = widget.NewRadioGroup([]string{driveSourceLabel, portabilitySourceLabel}, nil)
	s.source.Required = true
	s.source.Horizontal = true
	if s.a.Settings.EffectiveSource() == "portability" {
		s.source.SetSelected(portabilitySourceLabel)
	} else {
		s.source.SetSelected(driveSourceLabel)
	}
	s.email = widget.NewEntry()
	s.email.SetText(s.a.Settings.Email)
	s.email.SetPlaceHolder("Optional account hint; confirm the account in your browser")
	s.oauthClient = widget.NewEntry()
	s.oauthClient.SetText(s.a.Settings.OAuthClientFile)
	s.oauthClient.SetPlaceHolder("Select your OAuth desktop client JSON file")
	s.driveFolder = widget.NewEntry()
	s.driveFolder.SetText(s.a.Settings.DriveFolderID)
	s.driveFolder.SetPlaceHolder("Folder ID or https://drive.google.com/drive/folders/…")
	s.resources = widget.NewEntry()
	s.resources.SetText(strings.Join(s.a.Settings.PortabilityResources, ", "))
	s.resources.SetPlaceHolder("Explicit resource groups, e.g. myactivity.search")
	s.pollInterval = widget.NewEntry()
	s.pollInterval.SetText(strconv.FormatInt(int64(s.a.Settings.PollInterval().Seconds()), 10))
	s.sourceSummary = widget.NewLabel("")
	s.sourceSummary.Truncation = fyne.TextTruncateEllipsis
	s.sourceBlurb = wrappedLabel("")
	s.sourceHelp = wrappedLabel("Creates one API export job for your selected data when you press the button. Watching only checks existing jobs; it does not request more exports.")
	s.disclosure = wrappedLabel("")
	s.signInHelp = wrappedLabel("")
	s.accessStatus = wrappedLabel("Configured source only. Authorization has not been checked in this session; the account email is a sign-in hint, not a verified identity.")
	s.checkAccess = widget.NewButton("Check authorized access", s.inspectAccess)
	s.folderBrowse = widget.NewButtonWithIcon("Choose Drive folder…", theme.FolderOpenIcon(), s.browseFolders)
	s.create = widget.NewButtonWithIcon("Create Data Portability export", theme.ContentAddIcon(), s.createExport)
	s.openTakeout = widget.NewButton("Open Google Takeout", s.createExport)
	s.openTakeout.Importance = widget.HighImportance
	s.buttons = append(s.buttons, s.create, s.folderBrowse, s.checkAccess, s.openTakeout)
	s.driveOptions = container.NewVScroll(container.NewVBox(
		widget.NewLabelWithStyle("Export only what you need", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		wrappedLabel("1. Open Google Takeout."),
		container.NewHBox(s.openTakeout),
		wrappedLabel("2. Select Deselect all, then tick only the products you want."),
		wrappedLabel("3. Choose Add to Drive as the delivery method, then select Create export."),
		wrappedLabel("4. Google emails you when the export is ready, which can take hours or days. You can do Sign in and Save to now."),
		docLink("Google's Takeout instructions", "https://support.google.com/accounts/answer/3024190"),
	))
	s.portabilityOptions = s.resourceContent()
	return container.NewBorder(container.NewVBox(
		pageTitle("Choose what to download"),
		wrappedLabel("Pick where your data comes from, then choose the part of it you want."),
		s.source, s.sourceBlurb,
	), nil, nil, nil, container.NewStack(s.driveOptions, s.portabilityOptions))
}

// authorizeContent is step 2.
func (s *screen) authorizeContent() fyne.CanvasObject {
	clientBrowse := widget.NewButtonWithIcon("Select JSON…", theme.FileIcon(), func() {
		if s.busy {
			return
		}
		picker := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				s.errorDialog(err)
				return
			}
			if reader == nil {
				return
			}
			uri := reader.URI()
			if err := reader.Close(); err != nil {
				s.errorDialog(err)
				return
			}
			if s.busy {
				return
			}
			if uri.Scheme() != "file" {
				s.errorDialog(errors.New("select a local OAuth desktop client JSON file"))
				return
			}
			s.oauthClient.SetText(uri.Path())
		}, s.w)
		picker.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
		picker.Show()
	})
	s.buttons = append(s.buttons, clientBrowse)
	s.driveDocLink = docLink("How to create a Desktop app OAuth client", "https://developers.google.com/identity/protocols/oauth2/native-app")
	s.portabilityDocLink = docLink("Data Portability API setup", "https://developers.google.com/data-portability/user-guide/setup")
	// Keep form controls fixed while only the disclosure scrolls. This avoids
	// nested entry scroll clips in Fyne's software screenshot painter.
	s.disclosureScroll = container.NewVScroll(container.NewVBox(
		s.accessStatus,
		wrappedLabel("Use an OAuth client of type Desktop app for the selected Google API. takedown opens your system browser for OAuth with PKCE; complete Google sign-in and consent there."),
		widget.NewLabelWithStyle("Access and data handling", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		s.disclosure,
		docLink("Revoke remote access in Google Account permissions", "https://myaccount.google.com/permissions"),
	))
	return container.NewBorder(container.NewVBox(
		pageTitle("Sign in to Google"),
		wrappedLabel("takedown needs your own OAuth Desktop app client. Select its JSON file, then sign in with Google in your browser."),
		s.signInHelp,
		s.driveDocLink, s.portabilityDocLink,
		widget.NewForm(widget.NewFormItem("OAuth desktop client JSON", container.NewBorder(nil, nil, nil, clientBrowse, s.oauthClient)), widget.NewFormItem("Account email hint (optional)", s.email)),
		container.NewHBox(s.authorize, s.checkAccess, s.disconnect),
	), nil, nil, nil, s.disclosureScroll)
}

func docLink(label, target string) *widget.Hyperlink {
	uri, err := url.ParseRequestURI(target)
	if err != nil {
		panic("invalid built-in documentation URL: " + err.Error())
	}
	return widget.NewHyperlink(label, uri)
}

func (s *screen) selectedSource() string {
	if s.source.Selected == portabilitySourceLabel {
		return "portability"
	}
	return "drive"
}

func (s *screen) accessDisclosure() string {
	access := "Drive permission: read-only access can view and download all files in your Google Drive. takedown scans only the selected archive folder and its subfolders, but the OAuth permission is broader than that folder."
	if s.selectedSource() == "portability" {
		groups := strings.TrimSpace(s.resources.Text)
		if groups == "" {
			groups = "none selected; enter resource groups before authorizing"
		}
		access = "Data Portability archives require encryption in transit and at rest. takedown does not encrypt destination files. SFTP encrypts transit only; SMB encryption depends on your server/share configuration. You must establish suitable protected storage before downloading. takedown does not verify these protections."
		access += "\n\nData Portability permission: access is requested only for these resource groups: " + groups + ". takedown uses that access to request the exports you start and download their archives."
	}
	return access + "\n\nOAuth tokens are stored in your OS keyring. Archives go only to your chosen destination, including a network destination if you select one. takedown sends no developer telemetry. You control deletion and can remove the local token. This does not revoke Google's remote consent; revoke it in Google Account permissions. Downloading does not delete cloud data. Google's consent screen handles the actual permission grant."
}

func (s *screen) updateSourceLabels() {
	portability := s.selectedSource() == "portability"
	// The Jobs tab only manages Data Portability jobs; every control in it is unusable for Drive.
	if portability {
		s.tabs.EnableItem(s.jobsTab)
	} else {
		s.tabs.DisableItem(s.jobsTab)
	}
	for _, shown := range []struct {
		obj     fyne.CanvasObject
		visible bool
	}{
		{s.portabilityOptions, portability}, {s.portabilityPrepare, portability}, {s.portabilityDocLink, portability}, {s.portabilityWarning, portability}, {s.checkAccess, portability},
		{s.driveOptions, !portability}, {s.drivePrepare, !portability}, {s.driveDocLink, !portability},
	} {
		setVisible(shown.obj, shown.visible)
	}
	if portability {
		s.sourceBlurb.SetText("For eligible accounts with their own Google Cloud project, billing and encrypted storage. takedown requests only the data you tick below.")
		s.signInHelp.SetText("In a Google Cloud project with billing, enable the Data Portability API and meet Google's verification requirements. Then create a Desktop app OAuth client and download its JSON.")
	} else {
		s.sourceBlurb.SetText("Best for most people. You pick products on Google's Takeout page, Google delivers the archives to your Drive, and takedown downloads them.")
		s.signInHelp.SetText("In a Google Cloud project, enable the Google Drive API and add the drive.readonly scope to the consent screen. While the app is in Testing, add your account as a test user. Then create a Desktop app OAuth client and download its JSON.")
	}
	s.updateSummary()
	s.disclosure.SetText(s.accessDisclosure())
	s.content.Refresh()
}

func (s *screen) updateSummary() {
	destination := strings.TrimSpace(s.destination.Text)
	if destination == "" {
		destination = "no destination yet"
	}
	s.sourceSummary.SetText(s.source.Selected + " → " + destination)
}

func (s *screen) bindSourceChanges() {
	changed := func(string) {
		s.sourceRevision++
		s.invalidateExports()
		s.updateSourceLabels()
		if s.resourceUI != nil {
			s.resourceUI.sync()
		}
		s.accessStatus.SetText("Source settings changed. Configured choices are not proof of authorized access. Sign in again if needed; the email remains a sign-in hint.")
		s.invalidateManagement()
		if s.folderUI != nil && !s.folderUI.closed {
			s.folderUI.dialog.Hide()
		}
		if s.jobCancel != nil {
			s.jobCancel()
			s.status.SetText("Source changed. Stopping work; partial files will be kept.")
		} else {
			s.status.SetText("Settings updated. Sign in if these choices are new, then discover archives.")
		}
	}
	s.source.OnChanged = func(selected string) {
		// Drive allows faster polling. A saved Drive interval must not block every
		// Data Portability action from a field folded away on step 4.
		if interval, err := strconv.Atoi(strings.TrimSpace(s.pollInterval.Text)); err == nil && s.selectedSource() == "portability" && interval != 0 && interval < 300 {
			s.pollInterval.SetText("300")
		}
		changed(selected)
	}
	for _, entry := range []*widget.Entry{s.email, s.oauthClient, s.driveFolder, s.resources} {
		entry.OnChanged = changed
	}
}

func (s *screen) login() {
	if s.busy {
		return
	}
	if err := s.save(); err != nil {
		s.errorDialog(err)
		return
	}
	if strings.TrimSpace(s.oauthClient.Text) == "" {
		s.tabs.Select(s.authorizeTab)
		dialog.ShowInformation("Select an OAuth client", "Select your OAuth desktop client JSON file before authorizing.", s.w)
		return
	}
	if s.selectedSource() == "portability" && len(s.a.Settings.PortabilityResources) == 0 {
		s.tabs.Select(s.sourceTab)
		dialog.ShowInformation("Choose your data", "Tick the data you want to export before signing in.", s.w)
		return
	}
	binding := sourceSnapshot{identity: s.a.Settings.SourceIdentity(), revision: s.sourceRevision, account: s.a.Account}
	s.tabs.Select(s.authorizeTab)
	disclosure := container.NewVScroll(archiveLabel(s.accessDisclosure() + "\n\nContinue to Google's consent screen in your system browser?"))
	confirmation := dialog.NewCustomConfirm("Before you authorize takedown", "Continue to Google", "Cancel", disclosure, func(yes bool) {
		if !yes || s.busy || s.ctx.Err() != nil {
			return
		}
		if !binding.matches(s) {
			s.status.SetText("Source changed. Review its access disclosure before authorizing again.")
			return
		}
		s.invalidateExports()
		s.sourceRevision++
		s.invalidateManagement()
		s.run("Complete Google sign-in and consent in your system browser…", func(ctx context.Context, b sourceSnapshot) error { return b.account.Login(ctx) }, func(err error, _ sourceSnapshot) {
			if !s.showError(err) {
				text := "OAuth flow completed for the selected source. The account hint is not a verified identity."
				if s.selectedSource() == "portability" {
					text += " Use Check authorized access to inspect the current grant."
				}
				s.accessStatus.SetText(text)
				s.status.SetText("Authorization complete. Discover archives or start watching.")
			}
		}, nil)
	}, s.w)
	size := s.w.Canvas().Size()
	confirmation.Resize(fyne.NewSize(fyne.Min(680, size.Width-32), fyne.Min(480, size.Height-64)))
	confirmation.Show()
}

func (s *screen) forget() {
	if s.busy {
		return
	}
	s.sourceRevision++
	s.invalidateExports()
	s.invalidateManagement()
	s.run("Removing this source's OAuth token…", func(ctx context.Context, b sourceSnapshot) error { return b.account.Forget(ctx) }, func(err error, _ sourceSnapshot) {
		if !s.showError(err) {
			s.accessStatus.SetText("Local OAuth token removed. Google's remote consent may remain.")
			s.status.SetText("Local OAuth token removed; Google's remote consent remains. Revoke it in Google Account permissions if needed. Existing archives and partial files were kept.")
		}
	}, nil)
}

func (s *screen) createExport() {
	var result string
	s.run("Preparing your export request…", func(ctx context.Context, b sourceSnapshot) error {
		var err error
		result, err = b.account.CreateExport(ctx)
		return err
	}, func(err error, _ sourceSnapshot) {
		if s.showError(err) {
			return
		}
		s.status.SetText(result)
	}, nil)
}

func (s *screen) startWatch() {
	if s.busy {
		return
	}
	if err := s.save(); err != nil {
		s.errorDialog(err)
		return
	}
	var ids []string
	if s.watchSelected.Checked {
		if !s.requireDiscovery() {
			return
		}
		for _, export := range s.ready {
			if s.selected[export.ID] {
				ids = append(ids, export.ID)
			}
		}
		if len(ids) == 0 {
			dialog.ShowInformation("Select an export", "Select at least one export, or turn off Only watch selected exports to watch all archives from this source.", s.w)
			return
		}
	}
	opts := app.WatchOptions{DownloadOptions: app.DownloadOptions{Destination: s.a.Settings.Destination, PartByPart: s.partByPart.Checked, StoragePassword: s.password.Text}, Interval: s.a.Settings.PollInterval(), ExportIDs: ids}
	s.run(fmt.Sprintf("Watching for ready archives every %s…", opts.Interval), func(ctx context.Context, _ sourceSnapshot) error { return s.a.Watch(ctx, opts) }, nil, func() {
		s.watchDialog.Hide()
		s.watching = true
		s.tabs.Select(s.transfersTab)
		s.cancelButton.SetText("Stop watching")
		s.pause.Enable()
		s.resume.Enable()
		s.content.Refresh()
	})
}
