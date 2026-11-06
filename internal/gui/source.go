//go:build !headless

package gui

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Dunky13/takedown/internal/app"
)

const (
	driveSourceLabel       = "Takeout in Google Drive"
	portabilitySourceLabel = "Data Portability API"
	driveBlurb             = "Best for most people, and the route for Google Photos. You pick products on Google's Takeout page, Google delivers the archives to your Drive, and takedown downloads them."
	portabilityBlurb       = "For eligible accounts with their own Google Cloud project, billing and encrypted storage. takedown requests only the data you tick below. Google does not offer Google Photos here."
)

// sourceContent is step 1. It also creates the fields that later steps share.
func (s *screen) sourceContent() fyne.CanvasObject {
	s.source = widget.NewRadioGroup([]string{driveSourceLabel, portabilitySourceLabel}, nil)
	if s.a.Settings.EffectiveSource() == "portability" {
		s.source.SetSelected(portabilitySourceLabel)
	} else {
		s.source.SetSelected(driveSourceLabel)
	}
	s.email = newEntry()
	s.email.SetText(s.a.Settings.Email)
	s.email.SetPlaceHolder("Optional account hint; confirm the account in your browser")
	s.oauthClient = newEntry()
	s.oauthClient.SetText(s.a.Settings.OAuthClientFile)
	s.oauthClient.SetPlaceHolder("Select your OAuth desktop client JSON file")
	s.driveFolder = newEntry()
	s.driveFolder.SetText(s.a.Settings.DriveFolderID)
	s.driveFolder.SetPlaceHolder("Folder ID or https://drive.google.com/drive/folders/…")
	s.resources = newEntry()
	s.resources.SetText(strings.Join(s.a.Settings.PortabilityResources, ", "))
	s.resources.SetPlaceHolder("Explicit resource groups, e.g. myactivity.search")
	s.pollInterval = newEntry()
	s.pollInterval.SetText(strconv.FormatInt(int64(s.a.Settings.PollInterval().Seconds()), 10))
	s.sourceSummary = widget.NewLabel("")
	s.sourceSummary.Truncation = fyne.TextTruncateEllipsis
	s.sourceHelp = helperText("Creates one API export job for your selected data when you press the button. Watching only checks existing jobs; it does not request more exports.")
	s.disclosure = wrappedLabel("")
	s.signInHelp = helperText("")
	s.accessStatus = helperText("Configured source only. Authorization has not been checked in this session; the account email is a sign-in hint, not a verified identity.")
	s.checkAccess = newButton("Check authorized access", nil, s.inspectAccess)
	s.folderBrowse = newButton("Choose Drive folder…", theme.FolderOpenIcon(), s.browseFolders)
	s.create = newButton("Create Data Portability export", theme.ContentAddIcon(), s.createExport)
	s.openTakeout = newButton("Open Google Takeout", nil, s.createExport)
	s.buttons = append(s.buttons, s.create, s.folderBrowse, s.checkAccess, s.openTakeout)
	s.driveOptions = scrolled(container.NewVBox(
		flush(widget.NewLabelWithStyle("Export only what you need", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
		instructionRow(1, "Open Google Takeout.", s.openTakeout),
		instructionRow(2, "Select Deselect all, then tick only the products you want, such as Google Photos.", nil),
		instructionRow(3, "Choose Add to Drive as the delivery method, then select Create export.", nil),
		newCallout("Takeout archives count toward your Google storage until you delete them."),
		instructionRow(4, "Google emails you when the export is ready, which can take hours or days. You can do Sign in and Save to now.", nil),
		flush(docLink("Google's Takeout instructions", "https://support.google.com/accounts/answer/3024190")),
	))
	s.portabilityOptions = s.resourceContent()
	s.driveCard = newChoiceCard(driveSourceLabel, driveBlurb, func() { s.source.SetSelected(driveSourceLabel) })
	s.portabilityCard = newChoiceCard(portabilitySourceLabel, portabilityBlurb, func() { s.source.SetSelected(portabilitySourceLabel) })
	return border(spaced(
		pageHeader("Choose what to download", "Pick where your data comes from, then choose the part of it you want."),
		container.NewGridWithColumns(2, container.New(layout.NewCustomPaddedLayout(0, 0, 0, theme.Padding()), s.driveCard), container.New(layout.NewCustomPaddedLayout(0, 0, theme.Padding(), 0), s.portabilityCard)),
		gap(0),
	), nil, nil, nil, container.NewStack(s.driveOptions, s.portabilityOptions))
}

// authorizeContent is step 2.
func (s *screen) authorizeContent() fyne.CanvasObject {
	clientBrowse := newButton("Select JSON…", theme.FileIcon(), s.pickOAuthClient)
	s.buttons = append(s.buttons, clientBrowse)
	s.driveDocLink = docLink("How to create a Desktop app OAuth client", "https://developers.google.com/identity/protocols/oauth2/native-app")
	s.portabilityDocLink = docLink("Data Portability API setup", "https://developers.google.com/data-portability/user-guide/setup")
	s.driveLink, s.portabilityLink = flush(s.driveDocLink), flush(s.portabilityDocLink)
	// Keep form controls fixed while only the disclosure scrolls. This avoids
	// nested entry scroll clips in Fyne's software screenshot painter.
	s.disclosureAccordion = widget.NewAccordion(widget.NewAccordionItem("What takedown can access", container.NewVBox(
		flush(helperText("Use an OAuth client of type Desktop app for the selected Google API. takedown opens your system browser for OAuth with PKCE; complete Google sign-in and consent there.")),
		s.disclosure,
	)))
	s.disclosureScroll = scrolled(container.NewVBox(
		flush(s.accessStatus),
		s.disclosureAccordion,
	))
	return border(spaced(
		container.NewVBox(
			pageHeader("Sign in to Google", ""),
			flush(s.signInHelp),
			s.driveLink, s.portabilityLink,
		),
		field("OAuth desktop client JSON", border(nil, nil, nil, clientBrowse, s.oauthClient), ""),
		field("Account email hint (optional)", s.email, ""),
		container.NewVBox(
			container.NewHBox(s.authorize, s.checkAccess, s.disconnect),
			flush(docLink("Revoke remote access in Google Account permissions", "https://myaccount.google.com/permissions")),
		),
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
		{s.portabilityOptions, portability}, {s.portabilityPrepare, portability}, {s.portabilityLink, portability}, {s.portabilityWarning, portability}, {s.checkAccess, portability},
		{s.driveOptions, !portability}, {s.drivePrepare, !portability}, {s.driveLink, !portability},
	} {
		setVisible(shown.obj, shown.visible)
	}
	s.driveCard.setSelected(!portability)
	s.portabilityCard.setSelected(portability)
	if portability {
		s.signInHelp.SetText("takedown needs your own OAuth Desktop app client. In a Google Cloud project with billing, enable the Data Portability API (Google's verification requirements apply), create a Desktop app client, download its JSON and select it here. Then sign in with Google in your browser.")
	} else {
		s.signInHelp.SetText("takedown needs your own OAuth Desktop app client. In a Google Cloud project, enable the Google Drive API, add the drive.readonly scope to the consent screen and, while the app is in Testing, your account as a test user. Create a Desktop app client, download its JSON and select it here. Then sign in with Google in your browser.")
	}
	s.updateSummary()
	s.disclosure.SetText(s.accessDisclosure())
	relayout(s.content)
}

func (s *screen) updateSummary() {
	destination := strings.TrimSpace(s.destination.Text)
	if destination == "" {
		destination = "no destination yet"
	}
	s.sourceSummary.SetText(s.source.Selected + " · " + endOfPath(destination, 28))
}

// endOfPath keeps the end of a long path, where the folder name is.
func endOfPath(path string, limit int) string {
	runes := []rune(path)
	if len(runes) <= limit {
		return path
	}
	return "…" + string(runes[len(runes)-limit+1:])
}

func (s *screen) bindSourceChanges() {
	changed := func(string) {
		s.sourceRevision++
		s.invalidateExports()
		s.updateSourceLabels()
		if s.resourceUI != nil {
			s.resourceUI.sync()
		}
		s.setAccessStatus("Source settings changed. Configured choices are not proof of authorized access. Sign in again if needed; the email remains a sign-in hint.", widget.LowImportance)
		s.invalidateManagement()
		if s.folderUI != nil && !s.folderUI.closed {
			s.folderUI.dialog.Hide()
		}
		if s.jobCancel != nil {
			s.jobCancel()
			s.notice("Source changed. Stopping work; partial files will be kept.")
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
	for _, field := range []*entry{s.email, s.oauthClient, s.driveFolder, s.resources} {
		field.OnChanged = changed
	}
}

// setAccessStatus writes the access line; failures are not muted.
func (s *screen) setAccessStatus(text string, importance widget.Importance) {
	s.accessStatus.Importance = importance
	s.accessStatus.SetText(text)
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
		s.inform("Select an OAuth client", "Select your OAuth desktop client JSON file before authorizing.")
		return
	}
	if s.selectedSource() == "portability" && len(s.a.Settings.PortabilityResources) == 0 {
		s.tabs.Select(s.sourceTab)
		s.inform("Choose your data", "Tick the data you want to export before signing in.")
		return
	}
	binding := sourceSnapshot{identity: s.a.Settings.SourceIdentity(), revision: s.sourceRevision, account: s.a.Account}
	s.tabs.Select(s.authorizeTab)
	disclosure := container.NewVScroll(flush(archiveLabel(s.accessDisclosure() + "\n\nContinue to Google's consent screen in your system browser?")))
	size := s.w.Canvas().Size()
	s.confirmDialog("Before you authorize takedown", disclosure, "Continue to Google", "Cancel", fyne.NewSize(fyne.Min(680, size.Width-32), fyne.Min(480, size.Height-64)), func(yes bool) {
		if !yes || s.busy || s.ctx.Err() != nil {
			return
		}
		if !binding.matches(s) {
			s.notice("Source changed. Review its access disclosure before authorizing again.")
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
				s.setAccessStatus(text, widget.LowImportance)
				s.status.SetText("Authorization complete. Discover archives or start watching.")
			}
		}, nil)
	})
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
			s.setAccessStatus("Local OAuth token removed. Google's remote consent may remain.", widget.LowImportance)
			s.notice("Local OAuth token removed; Google's remote consent remains. Revoke it in Google Account permissions if needed. Existing archives and partial files were kept.")
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
			s.inform("Select an export", "Select at least one export, or turn off Only watch selected exports to watch all archives from this source.")
			return
		}
	}
	opts := app.WatchOptions{DownloadOptions: app.DownloadOptions{Destination: s.a.Settings.Destination, PartByPart: s.partByPart.Checked, StoragePassword: s.password.Text}, Interval: s.a.Settings.PollInterval(), ExportIDs: ids}
	s.run(fmt.Sprintf("Watching for ready archives every %s…", opts.Interval), func(ctx context.Context, _ sourceSnapshot) error { return s.a.Watch(ctx, opts) }, nil, func() {
		s.watchDialog.Hide()
		s.watching = true
		s.cancelButton.SetText("Stop watching")
		s.startTransfers()
		relayout(s.content)
	})
}
