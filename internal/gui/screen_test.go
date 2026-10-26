//go:build !headless

package gui

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/Dunky13/takedown/internal/app"
	"github.com/Dunky13/takedown/internal/auth"
	"github.com/Dunky13/takedown/internal/config"
	"github.com/Dunky13/takedown/internal/events"
	"github.com/Dunky13/takedown/internal/model"
)

// Queue background UI callbacks so the Fyne test driver has the same single
// UI execution context as the desktop driver. Its default DoFromGoroutine runs
// immediately on the worker, which races with taps and screenshot rendering.
type queuedApp struct {
	fyne.App
	driver *queuedDriver
}

func (a *queuedApp) Driver() fyne.Driver { return a.driver }

type queuedDriver struct {
	fyne.Driver
	pending chan func()
}

func (d *queuedDriver) DoFromGoroutine(work func(), wait bool) {
	if !wait {
		d.pending <- work
		return
	}
	finished := make(chan struct{})
	d.pending <- func() { defer close(finished); work() }
	<-finished
}

type screenAccount struct {
	login  func(context.Context) error
	list   func(context.Context) ([]model.Export, error)
	forget func(context.Context) error
	create func(context.Context) (string, error)
	client *http.Client
}

func (a *screenAccount) Login(ctx context.Context) error {
	if a.login != nil {
		return a.login(ctx)
	}
	return nil
}
func (a *screenAccount) List(ctx context.Context) ([]model.Export, error) {
	if a.list != nil {
		return a.list(ctx)
	}
	return nil, nil
}
func (*screenAccount) RefreshLink(context.Context, string, int) (string, error) { return "", nil }
func (a *screenAccount) HTTPClient() *http.Client {
	if a.client != nil {
		return a.client
	}
	return http.DefaultClient
}
func (a *screenAccount) Forget(ctx context.Context) error {
	if a.forget != nil {
		return a.forget(ctx)
	}
	return nil
}
func (a *screenAccount) CreateExport(ctx context.Context) (string, error) {
	if a.create != nil {
		return a.create(ctx)
	}
	return "export-job", nil
}
func (*screenAccount) Close() error { return nil }

func newTestScreen(t *testing.T, change func(*config.Settings)) (*screen, *screenAccount) {
	t.Helper()
	base := test.NewTempApp(t)
	desktop := &queuedApp{App: base, driver: &queuedDriver{Driver: base.Driver(), pending: make(chan func(), 1024)}}
	fyne.SetCurrentApp(desktop)
	w := desktop.NewWindow("takedown")
	dir := t.TempDir()
	settings := config.Settings{Destination: filepath.Join(dir, "archives"), SafetyBuffer: 2_000_000_000, Concurrency: 2, Source: "drive", Email: "first@example.com", OAuthClientFile: filepath.Join(dir, "desktop-client.json"), DriveFolderID: "archive-folder", PollIntervalSeconds: 300}
	if change != nil {
		change(&settings)
	}
	account := &screenAccount{}
	a := &app.App{DataDir: dir, Settings: settings, Bus: events.New(), Account: account}
	s := newScreen(context.Background(), desktop, w, a)
	w.Resize(fyne.NewSize(960, 760))
	w.SetContent(s.content)
	w.Show()
	t.Cleanup(func() {
		s.close()
		if err := a.Close(); err != nil {
			t.Error(err)
		}
		w.Close()
	})
	return s, account
}

func bindTestDiscovery(s *screen) {
	s.ready = []model.Export{{ID: "first-export", Parts: []model.Part{{ExportID: "first-export", Index: 1, Filename: "first.zip", Size: 100, URL: "https://example.invalid/archive?signed=first"}}}}
	s.selected["first-export"] = true
	s.discoveryBound = true
	s.discovery = sourceSnapshot{identity: s.a.Settings.SourceIdentity(), revision: s.sourceRevision, account: s.a.Account}
}

func awaitJob(t *testing.T, s *screen) {
	t.Helper()
	if s.jobDone == nil {
		t.Fatal("operation did not start")
	}
	pending := fyne.CurrentApp().Driver().(*queuedDriver).pending
	finished := s.jobDone
	deadline := time.NewTimer(5 * time.Second)
	defer deadline.Stop()
	for finished != nil || s.busy {
		select {
		case work := <-pending:
			work()
		case <-finished:
			finished = nil
		case <-deadline.C:
			t.Fatal("GUI operation did not finish")
		}
	}
	// Finish already queued telemetry before the next interaction.
	for {
		select {
		case work := <-pending:
			work()
		default:
			return
		}
	}
}

func awaitUIState(t *testing.T, ready func() bool) {
	t.Helper()
	pending := fyne.CurrentApp().Driver().(*queuedDriver).pending
	deadline := time.NewTimer(5 * time.Second)
	defer deadline.Stop()
	for !ready() {
		select {
		case work := <-pending:
			work()
		case <-deadline.C:
			t.Fatal("UI did not reach the expected state")
		}
	}
}

func receive[T any](t *testing.T, ch <-chan T) T {
	t.Helper()
	select {
	case value := <-ch:
		return value
	case <-time.After(5 * time.Second):
		t.Fatal("operation did not reach the expected step")
		var zero T
		return zero
	}
}

func captureScreen(t *testing.T, s *screen, filename string) {
	t.Helper()
	root := os.Getenv("TAKEDOWN_TEST_ARTIFACT_DIR")
	if root == "" {
		return
	}
	if err := os.MkdirAll(root, 0700); err != nil {
		t.Fatal(err)
	}
	f, err := os.Create(filepath.Join(root, filename))
	if err != nil {
		t.Fatal(err)
	}
	err = png.Encode(f, s.w.Canvas().Capture())
	closeErr := f.Close()
	if err := errors.Join(err, closeErr); err != nil {
		t.Fatal(err)
	}
}

func TestScreenRendersProgressAndSpaceWarning(t *testing.T) {
	s, _ := newTestScreen(t, nil)
	s.onEvent(events.Event{Kind: events.Progress, PartKey: "export/1", Bytes: 50_000_000, Total: 100_000_000, BytesPerSecond: 10_000_000})
	if s.progress.Value != 0.5 {
		t.Fatalf("progress = %f", s.progress.Value)
	}
	s.onEvent(events.Event{Kind: events.ReAuthRequired, PartKey: "export/1", Total: 100_000_000, Message: "Authorize your account in your system browser."})
	if s.progress.Value != 0.5 {
		t.Fatal("reauthorization reset retained byte progress")
	}
	s.onEvent(events.Event{Kind: events.SpaceWarning, Message: "Free space, then retry."})
	if !s.banner.Visible() || s.banner.Text != "Free space, then retry." {
		t.Fatal("space warning not visible")
	}
	if s.banner.Position().Y <= 0 {
		t.Fatal("warning overlaps the application heading")
	}
	s.setBusy(true)
	if !s.destination.Disabled() || !s.source.Disabled() || !s.oauthClient.Disabled() || s.cancelButton.Disabled() {
		t.Fatal("running job controls are unsafe")
	}
	s.setBusy(false)
	if s.destination.Disabled() || s.source.Disabled() || !s.cancelButton.Disabled() {
		t.Fatal("stopped job controls did not recover")
	}
	s.tabs.SelectIndex(1)
	captureScreen(t, s, "gui-preview.png")
}

func TestSourceChangesInvalidateSelectedURLsBeforeDownload(t *testing.T) {
	for _, name := range []string{"email", "client", "folder", "provider", "resources", "round-trip", "account-replaced"} {
		t.Run(name, func(t *testing.T) {
			s, _ := newTestScreen(t, func(settings *config.Settings) {
				if name == "resources" {
					settings.Source = "portability"
					settings.PortabilityResources = []string{"myactivity.search"}
				}
			})
			bindTestDiscovery(s)
			oldIdentity := s.a.Settings.SourceIdentity()
			switch name {
			case "email":
				s.email.SetText("second@example.com")
			case "client":
				s.oauthClient.SetText(filepath.Join(s.a.DataDir, "another-client.json"))
			case "folder":
				s.driveFolder.SetText("different-folder")
			case "provider":
				s.source.SetSelected(portabilitySourceLabel)
			case "resources":
				s.resources.SetText("myactivity.youtube")
			case "round-trip":
				s.email.SetText("second@example.com")
				s.email.SetText("first@example.com")
			case "account-replaced":
				s.a.Account = &screenAccount{}
			}
			s.download()
			if len(s.ready) != 0 || len(s.selected) != 0 || s.discoveryBound || s.busy || s.jobDone != nil {
				t.Fatal("source change scheduled stale archive URLs")
			}
			if name != "round-trip" && name != "account-replaced" && s.a.Settings.SourceIdentity() == oldIdentity {
				t.Fatal("new source preferences were not applied")
			}
		})
	}
}

func TestDiscoveryRejectsLateResultAfterSourceEdit(t *testing.T) {
	s, account := newTestScreen(t, nil)
	entered, release := make(chan struct{}), make(chan struct{})
	account.list = func(context.Context) ([]model.Export, error) {
		close(entered)
		<-release // A provider may finish after cancellation; still reject its result.
		return []model.Export{{ID: "old-source"}}, nil
	}
	s.discover()
	receive(t, entered)
	if !s.busy || !s.email.Disabled() {
		t.Fatal("discovery did not lock source preferences")
	}
	s.email.SetText("other@example.com")
	close(release)
	awaitJob(t, s)
	if s.discoveryBound || len(s.ready) != 0 || s.busy {
		t.Fatal("late discovery rebound stale archives")
	}
	if !strings.Contains(s.status.Text, "Source changed") {
		t.Fatalf("status = %q", s.status.Text)
	}
}

func TestSourceOnboardingAndDisclosure(t *testing.T) {
	s, _ := newTestScreen(t, nil)
	s.tabs.SelectIndex(2)
	if !s.driveOptions.Visible() || s.portabilityOptions.Visible() || !strings.Contains(s.disclosure.Text, "all files in your Google Drive") {
		t.Fatal("Drive setup or permission breadth missing")
	}
	for _, expected := range []string{"OS keyring", "chosen destination", "network destination", "no developer telemetry", "control deletion", "local token", "remote consent", "Google Account permissions", "Google's consent screen"} {
		if !strings.Contains(s.disclosure.Text, expected) {
			t.Errorf("disclosure missing %q", expected)
		}
	}
	captureScreen(t, s, "gui-source-drive.png")
	s.sourceTabs.SelectIndex(1)
	captureScreen(t, s, "gui-source-drive-access.png")
	s.sourceTabs.SelectIndex(0)
	s.source.SetSelected(portabilitySourceLabel)
	s.resources.SetText("myactivity.search, myactivity.youtube")
	if s.driveOptions.Visible() || !s.portabilityOptions.Visible() {
		t.Fatal("provider switch did not change setup fields")
	}
	if !strings.Contains(s.disclosure.Text, "only for these resource groups: myactivity.search, myactivity.youtube") || strings.Contains(s.disclosure.Text, "all files in your Google Drive") {
		t.Fatal("Data Portability disclosure did not follow the resource selection")
	}
	if s.create.Text != "Create Data Portability export" {
		t.Fatalf("export action = %q", s.create.Text)
	}
	captureScreen(t, s, "gui-source-portability.png")
	s.sourceTabs.SelectIndex(1)
	captureScreen(t, s, "gui-source-portability-access.png")
	s.sourceTabs.SelectIndex(0)
	cfg, err := s.preferences()
	if err != nil || len(cfg.PortabilityResources) != 2 {
		t.Fatalf("resource groups = %v, %v", cfg.PortabilityResources, err)
	}
	s.source.SetSelected(driveSourceLabel)
	cfg, err = s.preferences()
	if err != nil || len(cfg.PortabilityResources) != 0 {
		t.Fatal("Drive preferences retained incompatible portability scopes")
	}
}

func findDialogButton(obj fyne.CanvasObject, label string) *widget.Button {
	switch object := obj.(type) {
	case *widget.Button:
		if object.Text == label {
			return object
		}
	case *widget.PopUp:
		return findDialogButton(object.Content, label)
	case *fyne.Container:
		for _, child := range object.Objects {
			if button := findDialogButton(child, label); button != nil {
				return button
			}
		}
	case fyne.Widget:
		for _, child := range test.WidgetRenderer(object).Objects() {
			if button := findDialogButton(child, label); button != nil {
				return button
			}
		}
	}
	return nil
}

func tapDialogButton(t *testing.T, s *screen, label string) {
	t.Helper()
	button := findDialogButton(s.w.Canvas().Overlays().Top(), label)
	if button == nil {
		t.Fatalf("dialog button %q not found", label)
	}
	test.Tap(button)
}

func dialogHasText(obj fyne.CanvasObject, text string) bool {
	switch object := obj.(type) {
	case *widget.Label:
		return strings.Contains(object.Text, text)
	case *widget.PopUp:
		return dialogHasText(object.Content, text)
	case *fyne.Container:
		for _, child := range object.Objects {
			if dialogHasText(child, text) {
				return true
			}
		}
	case fyne.Widget:
		for _, child := range test.WidgetRenderer(object).Objects() {
			if dialogHasText(child, text) {
				return true
			}
		}
	}
	return false
}

func TestPortabilityEncryptionDisclosurePrecedesOAuth(t *testing.T) {
	s, account := newTestScreen(t, func(cfg *config.Settings) {
		cfg.Source = "portability"
		cfg.PortabilityResources = []string{"myactivity.search"}
	})
	var calls atomic.Int32
	account.login = func(context.Context) error { calls.Add(1); return nil }
	s.w.Resize(fyne.NewSize(800, 640))
	s.login()
	for _, required := range []string{
		"Data Portability archives require encryption in transit and at rest",
		"takedown does not encrypt destination files",
		"SFTP encrypts transit only",
		"SMB encryption depends on your server/share configuration",
		"You must establish suitable protected storage before downloading",
		"takedown does not verify these protections",
	} {
		if !strings.Contains(s.disclosure.Text, required) || !dialogHasText(s.w.Canvas().Overlays().Top(), required) {
			t.Errorf("source disclosure or pre-OAuth confirmation missing %q", required)
		}
	}
	if calls.Load() != 0 || s.busy || s.jobDone != nil {
		t.Fatal("OAuth started before the storage requirements were acknowledged")
	}
	captureScreen(t, s, "gui-portability-encryption-disclosure-800.png")
	tapDialogButton(t, s, "Cancel")
	if calls.Load() != 0 {
		t.Fatal("canceling the storage disclosure started OAuth")
	}
	s.source.SetSelected(driveSourceLabel)
	if strings.Contains(s.disclosure.Text, "Data Portability archives require") {
		t.Fatal("Data Portability storage requirements appeared in the Drive disclosure")
	}
}

func TestAuthorizeRequiresDisclosureConfirmation(t *testing.T) {
	s, account := newTestScreen(t, nil)
	var calls atomic.Int32
	account.login = func(ctx context.Context) error {
		calls.Add(1)
		return nil
	}
	s.login()
	if calls.Load() != 0 || s.busy || s.w.Canvas().Overlays().Top() == nil {
		t.Fatal("authorization started before disclosure confirmation")
	}
	captureScreen(t, s, "gui-oauth-disclosure.png")
	tapDialogButton(t, s, "Cancel")
	if calls.Load() != 0 || s.jobDone != nil {
		t.Fatal("canceling disclosure started OAuth")
	}
	s.login()
	tapDialogButton(t, s, "Continue to Google")
	awaitJob(t, s)
	if calls.Load() != 1 || s.busy || !strings.Contains(s.status.Text, "Authorization complete") {
		t.Fatal("explicit confirmation did not authorize exactly once")
	}
}

func TestAuthorizeRejectsChangedSourceWhileDisclosureIsOpen(t *testing.T) {
	s, account := newTestScreen(t, nil)
	var calls atomic.Int32
	account.login = func(context.Context) error { calls.Add(1); return nil }
	s.login()
	s.driveFolder.SetText("another-folder")
	tapDialogButton(t, s, "Continue to Google")
	if calls.Load() != 0 || s.busy || s.jobDone != nil {
		t.Fatal("a changed source reused an earlier authorization disclosure")
	}
}

func TestDisconnectInvalidatesDiscoveryWithoutCredentialsForm(t *testing.T) {
	s, account := newTestScreen(t, nil)
	bindTestDiscovery(s)
	var calls atomic.Int32
	account.forget = func(context.Context) error { calls.Add(1); return nil }
	test.Tap(s.disconnect)
	awaitJob(t, s)
	if calls.Load() != 1 || s.discoveryBound || len(s.ready) != 0 || len(s.selected) != 0 {
		t.Fatal("disconnect retained old discovery or did not remove the OAuth token")
	}
	if !s.password.Password || !strings.Contains(s.password.PlaceHolder, "Storage password") {
		t.Fatal("network storage password field was removed or repurposed")
	}
}

func TestCreateExportIsExplicitAndDoesNotStartWatching(t *testing.T) {
	for _, provider := range []string{"drive", "portability"} {
		t.Run(provider, func(t *testing.T) {
			s, account := newTestScreen(t, func(cfg *config.Settings) {
				cfg.Source = provider
				if provider == "portability" {
					cfg.PortabilityResources = []string{"myactivity.search"}
				}
			})
			var calls, discoveries atomic.Int32
			message := "Review and confirm Create export in Google Takeout."
			if provider == "portability" {
				message = "Data Portability archive job explicit-job created and recorded."
			}
			account.create = func(context.Context) (string, error) { calls.Add(1); return message, nil }
			account.list = func(context.Context) ([]model.Export, error) { discoveries.Add(1); return nil, nil }
			test.Tap(s.create)
			awaitJob(t, s)
			if calls.Load() != 1 || discoveries.Load() != 0 || s.watching || s.busy || s.status.Text != message {
				t.Fatal("export creation did not remain one explicit source action")
			}
		})
	}
}

func TestWatchControlsAndCancellationRetainPartials(t *testing.T) {
	s, account := newTestScreen(t, nil)
	listed := make(chan struct{}, 1)
	account.list = func(context.Context) ([]model.Export, error) { listed <- struct{}{}; return nil, nil }
	if err := os.MkdirAll(s.a.Settings.Destination, 0700); err != nil {
		t.Fatal(err)
	}
	partial := filepath.Join(s.a.Settings.Destination, "retained.zip.part")
	if err := os.WriteFile(partial, []byte("partial archive"), 0600); err != nil {
		t.Fatal(err)
	}
	test.Tap(s.watch)
	receive(t, listed)
	if !s.busy || !s.watching || !s.source.Disabled() || !s.driveFolder.Disabled() || !s.resources.Disabled() || !s.pollInterval.Disabled() || !s.watch.Disabled() || !s.find.Disabled() || !s.create.Disabled() || !s.disconnect.Disabled() || !s.watchSelected.Disabled() {
		t.Fatal("watch did not freeze preferences and competing actions")
	}
	if s.pause.Disabled() || s.resume.Disabled() || s.cancelButton.Disabled() || s.cancelButton.Text != "Stop watching" {
		t.Fatal("watch transfer controls unavailable")
	}
	captureScreen(t, s, "gui-watch.png")
	s.w.Resize(fyne.NewSize(800, 640))
	s.content.Refresh()
	assertWindowFits(t, s, fyne.NewSize(800, 640))
	captureScreen(t, s, "gui-watch-800.png")
	for _, control := range []*widget.Button{s.pause, s.resume} {
		test.Tap(control)
		if !s.pause.Disabled() || !s.resume.Disabled() || s.cancelButton.Disabled() {
			t.Fatal("pause/resume requests were not serialized with Stop available")
		}
		awaitUIState(t, func() bool { return !s.pause.Disabled() && !s.resume.Disabled() })
	}
	if err := s.save(); err == nil {
		t.Fatal("busy watch allowed a preference rebind")
	}
	test.Tap(s.pause)
	test.Tap(s.cancelButton)
	awaitJob(t, s)
	if s.controlPending || s.busy || s.watching || s.source.Disabled() || s.watch.Disabled() || !s.pause.Disabled() || !s.resume.Disabled() || !s.cancelButton.Disabled() {
		t.Fatal("watch controls did not recover after stop")
	}
	if got, err := os.ReadFile(partial); err != nil || string(got) != "partial archive" {
		t.Fatal("stopping watch removed or changed a partial archive")
	}
	if !strings.Contains(s.status.Text, "Stopped") {
		t.Fatalf("status = %q", s.status.Text)
	}
}

func TestWatchDownloadsFromSourceWithoutPriorDiscovery(t *testing.T) {
	s, account := newTestScreen(t, func(cfg *config.Settings) { cfg.SafetyBuffer = 1 })
	var archive bytes.Buffer
	zipWriter := zip.NewWriter(&archive)
	file, err := zipWriter.Create("Takeout/data.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.Write([]byte("test archive bytes")); err != nil {
		t.Fatal(err)
	}
	if err := zipWriter.Close(); err != nil {
		t.Fatal(err)
	}
	content := archive.Bytes()
	requested, release := make(chan struct{}), make(chan struct{})
	defer func() {
		select {
		case <-release:
		default:
			close(release)
		}
	}()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(requested)
		<-release
		w.Header().Set("Content-Type", "application/zip")
		http.ServeContent(w, r, "new.zip", time.Time{}, bytes.NewReader(content))
	}))
	t.Cleanup(server.Close)
	account.client = server.Client()
	account.list = func(context.Context) ([]model.Export, error) {
		return []model.Export{{ID: "new-export", Parts: []model.Part{{ExportID: "new-export", Index: 1, Filename: "new.zip", Size: int64(len(content)), URL: server.URL}}}}, nil
	}
	ch, stop := s.a.Bus.Subscribe(64)
	defer stop()
	destination := s.a.Settings.Destination
	s.startWatch()
	receive(t, requested)
	replacement := t.TempDir()
	s.destination.SetText(replacement)
	if err := s.save(); err == nil {
		t.Fatal("running watch applied an edited destination")
	}
	close(release)
	for {
		e := receive(t, ch)
		if e.Kind == events.PartCompleted {
			break
		}
	}
	test.Tap(s.cancelButton)
	awaitJob(t, s)
	archivePath := filepath.FromSlash(model.NamespacePath("new-export", "new.zip"))
	got, err := os.ReadFile(filepath.Join(destination, archivePath))
	if err != nil || string(got) != string(content) {
		t.Fatalf("watch archive = %q, %v", got, err)
	}
	if _, err := os.Stat(filepath.Join(replacement, archivePath)); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("watch did not retain its original destination snapshot")
	}
	if entries, err := os.ReadDir(replacement); err != nil || len(entries) != 0 {
		t.Fatal("running watch wrote into the edited destination")
	}
	if s.discoveryBound {
		t.Fatal("watch unexpectedly reused the manual discovery selection")
	}
}

func TestWatchRejectsSourceEditAndKeepsStopCancellation(t *testing.T) {
	s, account := newTestScreen(t, nil)
	listed := make(chan struct{}, 1)
	account.list = func(context.Context) ([]model.Export, error) { listed <- struct{}{}; return nil, nil }
	s.startWatch()
	receive(t, listed)
	s.oauthClient.SetText(filepath.Join(s.a.DataDir, "another-client.json"))
	awaitJob(t, s)
	if s.busy || s.watching || s.discoveryBound || !strings.Contains(s.status.Text, "Source changed") {
		t.Fatal("source edit did not stop the watch and invalidate its authority")
	}
}

func TestWatchAuthorizationFailureIsTerminalAndVisible(t *testing.T) {
	s, account := newTestScreen(t, nil)
	account.list = func(context.Context) ([]model.Export, error) { return nil, auth.ErrAuthorizationRequired }
	s.startWatch()
	awaitJob(t, s)
	if s.busy || s.watching || !s.banner.Visible() || !strings.Contains(s.banner.Text, "authorization required") {
		t.Fatal("authorization failure did not stop watching with a visible alert")
	}
	if s.authorize.Disabled() || !s.cancelButton.Disabled() || !s.pause.Disabled() {
		t.Fatal("terminal watch failure left unsafe or unavailable controls")
	}
	captureScreen(t, s, "gui-watch-authorization-required.png")
}

func TestCanceledDiscoveryCannotPublishAnIgnoredCancellation(t *testing.T) {
	s, account := newTestScreen(t, nil)
	entered, release := make(chan struct{}), make(chan struct{})
	account.list = func(context.Context) ([]model.Export, error) {
		close(entered)
		<-release
		return []model.Export{{ID: "after-stop"}}, nil
	}
	s.discover()
	receive(t, entered)
	test.Tap(s.cancelButton)
	close(release)
	awaitJob(t, s)
	if s.discoveryBound || len(s.ready) != 0 || !strings.Contains(s.status.Text, "Stopped") {
		t.Fatal("a canceled discovery published provider results after Stop")
	}
}

func TestWatchSelectedRejectsStaleDiscovery(t *testing.T) {
	s, _ := newTestScreen(t, nil)
	bindTestDiscovery(s)
	s.watchSelected.SetChecked(true)
	s.driveFolder.SetText("another-folder")
	s.startWatch()
	if s.busy || s.jobDone != nil || s.discoveryBound {
		t.Fatal("selected watch accepted exports from another source")
	}
}

func assertWindowFits(t *testing.T, s *screen, size fyne.Size) {
	t.Helper()
	if minimum := s.content.MinSize(); minimum.Width > size.Width || minimum.Height > size.Height {
		t.Fatalf("content minimum %v exceeds requested window %v", minimum, size)
	}
	for _, button := range []*widget.Button{s.authorize, s.find, s.disconnect, s.pause, s.resume, s.cancelButton} {
		position := fyne.CurrentApp().Driver().AbsolutePositionForObject(button)
		if position.X < 0 || position.Y < 0 || position.X+button.Size().Width > size.Width || position.Y+button.Size().Height > size.Height || button.Size().Width < button.MinSize().Width {
			t.Fatalf("button %q is clipped at %v size %v, minimum %v in %v", button.Text, position, button.Size(), button.MinSize(), size)
		}
	}
}

func TestOnboardingFitsDesktopAndSmallerWindow(t *testing.T) {
	s, _ := newTestScreen(t, func(cfg *config.Settings) { cfg.OAuthClientFile = "" })
	if s.tabs.SelectedIndex() != 2 {
		t.Fatal("unconfigured source did not open onboarding")
	}
	for _, viewport := range []struct {
		size fyne.Size
		name string
	}{{fyne.NewSize(960, 760), "gui-onboarding-960.png"}, {fyne.NewSize(800, 640), "gui-onboarding-800.png"}} {
		s.w.Resize(viewport.size)
		s.content.Refresh()
		assertWindowFits(t, s, viewport.size)
		captureScreen(t, s, viewport.name)
	}
}

func TestLongDriveIdentityAndFilenameWrapWithoutExposingIDs(t *testing.T) {
	s, account := newTestScreen(t, nil)
	filename := "takeout-" + strings.Repeat("archive_with_a_very_long_name_", 7) + ".zip"
	metadata := struct {
		File    string `json:"file"`
		Folder  string `json:"folder"`
		Name    string `json:"name"`
		Size    int64  `json:"size"`
		Version int    `json:"version"`
		MD5     string `json:"md5"`
	}{strings.Repeat("f", 100), strings.Repeat("d", 100), filename, 100_000_000, 1, "0123456789abcdef0123456789abcdef"}
	encoded, err := json.Marshal(metadata)
	if err != nil {
		t.Fatal(err)
	}
	opaqueID := "drive-" + base64.RawURLEncoding.EncodeToString(encoded)
	if len(opaqueID) < 500 || len(filename) > 255 {
		t.Fatal("fixture does not represent a long supported Drive archive")
	}
	part := model.Part{ExportID: opaqueID, Index: 1, Filename: filename, Size: metadata.Size, MD5: metadata.MD5}
	account.list = func(context.Context) ([]model.Export, error) {
		return []model.Export{{ID: opaqueID, CreatedAt: time.Date(2026, 10, 25, 10, 0, 0, 0, time.UTC), Parts: []model.Part{part}}}, nil
	}
	s.discover()
	awaitJob(t, s)
	if !s.selected[opaqueID] || len(s.ready) != 1 || s.ready[0].ID != opaqueID {
		t.Fatal("opaque identity was lost from selection state")
	}
	var nameLabel *widget.Label
	for _, object := range s.exports.Objects {
		switch item := object.(type) {
		case *widget.Check:
			if strings.Contains(item.Text, opaqueID) {
				t.Fatal("discovery exposed an opaque identity")
			}
		case *widget.Label:
			if strings.Contains(item.Text, filename) {
				nameLabel = item
			}
		}
	}
	if nameLabel == nil || nameLabel.Wrapping != fyne.TextWrapBreak {
		t.Fatal("full filename is not displayed with wrapping")
	}
	for _, viewport := range []struct {
		size   fyne.Size
		suffix string
	}{{fyne.NewSize(960, 760), "960"}, {fyne.NewSize(800, 640), "800"}} {
		s.w.Resize(viewport.size)
		s.tabs.SelectIndex(0)
		s.content.Refresh()
		assertWindowFits(t, s, viewport.size)
		captureScreen(t, s, "gui-long-drive-archives-"+viewport.suffix+".png")
		s.tabs.SelectIndex(1)
		s.onEvent(events.Event{Kind: events.Progress, PartKey: part.Key(), Filename: filename, Bytes: 50_000_000, Total: part.Size, BytesPerSecond: 10_000_000})
		s.content.Refresh()
		assertWindowFits(t, s, viewport.size)
		label := archiveLabel("")
		s.transferList.UpdateItem(0, label)
		if strings.Contains(label.Text, opaqueID) || !strings.Contains(label.Text, filename) || label.Wrapping != fyne.TextWrapBreak || label.MinSize().Width > viewport.size.Width {
			t.Fatal("transfer label exposed identity, hid the filename, or forced excess width")
		}
		captureScreen(t, s, "gui-long-drive-transfers-"+viewport.suffix+".png")
	}
	fallback := transferText(events.Event{PartKey: part.Key(), Total: 100})
	if !strings.HasPrefix(fallback, "Archive part\n") || strings.Contains(fallback, opaqueID) {
		t.Fatal("missing filename fell back to an opaque part key")
	}
	s.onEvent(events.Event{Kind: events.ReAuthRequired, PartKey: part.Key(), Message: "Authorize in your system browser."})
	if s.transfers[0].Filename != filename {
		t.Fatal("a status event lost the known archive filename")
	}
}

func TestDisclosureScrollDoesNotPaintOverFixedControls(t *testing.T) {
	s, _ := newTestScreen(t, func(cfg *config.Settings) {
		cfg.Source = "portability"
		cfg.PortabilityResources = []string{"myactivity.search", "myactivity.youtube"}
	})
	s.tabs.SelectIndex(2)
	s.sourceTabs.SelectIndex(1)
	if !strings.Contains(s.disclosure.Text, "Data Portability archives require encryption in transit and at rest") {
		t.Fatal("scroll clipping regression does not cover the encryption disclosure")
	}
	for _, viewport := range []struct {
		size   fyne.Size
		suffix string
	}{{fyne.NewSize(960, 760), "960"}, {fyne.NewSize(800, 640), "800"}} {
		s.w.Resize(viewport.size)
		s.disclosureScroll.ScrollToTop()
		s.content.Refresh()
		assertWindowFits(t, s, viewport.size)
		before := s.w.Canvas().Capture()
		position := fyne.CurrentApp().Driver().AbsolutePositionForObject(s.disclosureScroll)
		bottom := position.Y + s.disclosureScroll.Size().Height
		s.disclosureScroll.ScrollToBottom()
		after := s.w.Canvas().Capture()
		for y := before.Bounds().Min.Y; y < before.Bounds().Max.Y; y++ {
			if float32(y) >= position.Y && float32(y) < bottom {
				continue
			}
			for x := before.Bounds().Min.X; x < before.Bounds().Max.X; x++ {
				r1, g1, b1, a1 := before.At(x, y).RGBA()
				r2, g2, b2, a2 := after.At(x, y).RGBA()
				if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
					t.Fatalf("disclosure scroll painted outside its viewport at (%d, %d), window %v", x, y, viewport.size)
				}
			}
		}
		captureScreen(t, s, "gui-portability-disclosure-scrolled-"+viewport.suffix+".png")
	}
}

func TestWatchIntervalAndResourceValidation(t *testing.T) {
	s, _ := newTestScreen(t, nil)
	s.pollInterval.SetText("59")
	if _, err := s.preferences(); err == nil {
		t.Fatal("Drive interval below 60 seconds accepted")
	}
	s.pollInterval.SetText("60")
	if _, err := s.preferences(); err != nil {
		t.Fatalf("Drive minimum rejected: %v", err)
	}
	s.source.SetSelected(portabilitySourceLabel)
	if _, err := s.preferences(); err == nil {
		t.Fatal("Data Portability interval below 300 seconds accepted")
	}
	s.pollInterval.SetText("")
	cfg, err := s.preferences()
	if err != nil || cfg.PollInterval() != 5*time.Minute {
		t.Fatalf("default interval = %s, %v", cfg.PollInterval(), err)
	}
	s.resources.SetText("myactivity.search,,myactivity.youtube")
	if _, err := s.preferences(); err == nil {
		t.Fatal("empty resource group was silently dropped")
	}
}
