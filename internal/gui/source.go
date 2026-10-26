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
	portabilitySourceLabel = "Data Portability API (optional)"
)

func (s *screen) sourceContent() fyne.CanvasObject {
	s.source = widget.NewSelect([]string{driveSourceLabel, portabilitySourceLabel}, nil)
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
	s.sourceSummary = wrappedLabel("")
	s.sourceHelp = wrappedLabel("")
	s.disclosure = wrappedLabel("")
	s.accessStatus = wrappedLabel("Configured source only. Authorization has not been checked in this session; the account email is a sign-in hint, not a verified identity.")
	s.checkAccess = widget.NewButton("Check authorized access", s.inspectAccess)
	s.folderBrowse = widget.NewButtonWithIcon("Choose Drive folder…", theme.FolderOpenIcon(), s.browseFolders)
	s.create = widget.NewButtonWithIcon("", theme.ContentAddIcon(), s.createExport)
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
	s.buttons = append(s.buttons, clientBrowse, s.create, s.folderBrowse, s.checkAccess)
	s.driveOptions = container.NewBorder(container.NewVBox(
		wrappedLabel("In Google Takeout, choose products and Add to Drive delivery, then create an export or schedule. After delivery, choose its archive folder. Folder browsing runs only when you request it; watching scans only the chosen archive tree."),
		s.folderBrowse,
		wrappedLabel("Advanced: paste a folder ID or link"), s.driveFolder,
	), nil, nil, nil, container.NewVScroll(container.NewVBox(docLink("Google's Takeout instructions", "https://support.google.com/accounts/answer/3024190"))))
	s.portabilityOptions = s.resourceContent()
	s.updateSourceLabels()
	sourcePage := container.NewBorder(widget.NewForm(widget.NewFormItem("Source", s.source)), nil, nil, nil, container.NewStack(s.driveOptions, s.portabilityOptions))
	// Keep form controls fixed while only the disclosure scrolls. This avoids
	// nested entry scroll clips in Fyne's software screenshot painter.
	s.disclosureScroll = container.NewVScroll(container.NewVBox(
		s.accessStatus,
		wrappedLabel("Use an OAuth client of type Desktop app for the selected Google API. takedown opens your system browser for OAuth with PKCE; complete Google sign-in and consent there."),
		widget.NewLabelWithStyle("Access and data handling", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		s.disclosure,
		docLink("Revoke remote access in Google Account permissions", "https://myaccount.google.com/permissions"),
	))
	authorizationPage := container.NewBorder(container.NewVBox(
		wrappedLabel("Select the downloaded JSON for your OAuth Desktop app client."),
		widget.NewForm(widget.NewFormItem("OAuth desktop client JSON", container.NewBorder(nil, nil, nil, clientBrowse, s.oauthClient)), widget.NewFormItem("Account email hint (optional)", s.email)),
		s.checkAccess,
	), nil, nil, nil, s.disclosureScroll)
	exportPage := container.NewVBox(
		s.create, s.sourceHelp,
		widget.NewForm(widget.NewFormItem("Check every (seconds)", s.pollInterval)),
		wrappedLabel("Default: 300 seconds. Minimum: 60 for Drive, 300 for Data Portability. Discover and Start watching are in the Archives tab. Stop watching keeps partial files; Pause downloads keeps the watch running."),
	)
	s.sourceTabs = container.NewAppTabs(
		container.NewTabItem("1. Source", sourcePage),
		container.NewTabItem("2. Authorization", authorizationPage),
		container.NewTabItem("3. Export / watch", exportPage),
	)
	return s.sourceTabs
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
	if s.selectedSource() == "portability" {
		s.driveOptions.Hide()
		s.portabilityOptions.Show()
		s.sourceSummary.SetText("Source: Data Portability API · account email is a hint, not a verified identity")
		s.create.SetText("Create Data Portability export")
		s.checkAccess.Show()
		s.sourceHelp.SetText("Creates one API export job for your selected resource groups when you press the button. Watching only checks existing jobs; it does not request more exports.")
	} else {
		s.portabilityOptions.Hide()
		s.driveOptions.Show()
		s.sourceSummary.SetText("Source: Takeout delivered to Google Drive · account email is a sign-in hint")
		s.create.SetText("Open Takeout to create / schedule")
		s.checkAccess.Hide()
		s.sourceHelp.SetText("Opens Google Takeout in your normal browser. Review the products, Add to Drive delivery, and schedule there, then confirm Create export yourself. After delivery, configure the archive folder above.")
	}
	s.disclosure.SetText(s.accessDisclosure())
}

func (s *screen) bindSourceChanges() {
	changed := func(string) {
		s.sourceRevision++
		s.invalidateExports()
		s.updateSourceLabels()
		if s.resourceUI != nil {
			s.resourceUI.sync()
		}
		s.accessStatus.SetText("Source settings changed. Configured choices are not proof of authorized access. Check access again; the email remains a sign-in hint.")
		s.invalidateManagement()
		if s.folderUI != nil && !s.folderUI.closed {
			s.folderUI.dialog.Hide()
		}
		if s.jobCancel != nil {
			s.jobCancel()
			s.status.SetText("Source changed. Stopping work; partial files will be kept.")
		} else {
			s.status.SetText("Source settings changed. Authorize if needed, then discover archives again.")
		}
	}
	s.source.OnChanged = changed
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
		s.tabs.SelectIndex(2)
		s.sourceTabs.SelectIndex(1)
		dialog.ShowInformation("Select an OAuth client", "Select your OAuth desktop client JSON file before authorizing.", s.w)
		return
	}
	if s.selectedSource() == "portability" && len(s.a.Settings.PortabilityResources) == 0 {
		s.tabs.SelectIndex(2)
		dialog.ShowInformation("Choose resource groups", "Enter the Data Portability resource groups you want to authorize.", s.w)
		return
	}
	binding := sourceSnapshot{identity: s.a.Settings.SourceIdentity(), revision: s.sourceRevision, account: s.a.Account}
	s.tabs.SelectIndex(2)
	s.sourceTabs.SelectIndex(1)
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
				s.accessStatus.SetText("OAuth flow completed for the selected source. The account hint is not a verified identity; Check authorized access to inspect the current grant.")
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
		s.watching = true
		s.tabs.SelectIndex(1)
		s.cancelButton.SetText("Stop watching")
		s.pause.Enable()
		s.resume.Enable()
		s.content.Refresh()
	})
}
