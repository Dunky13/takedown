//go:build !headless

package gui

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"github.com/Dunky13/takedown/internal/auth"
	"github.com/Dunky13/takedown/internal/config"
	"github.com/Dunky13/takedown/internal/downloader"
	"github.com/Dunky13/takedown/internal/model"
	"github.com/Dunky13/takedown/internal/source"
)

func TestGuideBackAndNextWalkTheFourSteps(t *testing.T) {
	s, _ := newTestScreen(t, func(cfg *config.Settings) { cfg.OAuthClientFile = "" })
	if s.tabs.Selected() != s.sourceTab || s.back.Visible() || s.next.Text != "Next: Sign in" {
		t.Fatalf("first step: back visible=%v, next=%q", s.back.Visible(), s.next.Text)
	}
	want := [][2]string{{"Back: Choose data", "Next: Save to"}, {"Back: Sign in", "Next: Download"}}
	for i, labels := range want {
		test.Tap(s.next)
		if s.tabs.Selected() != s.guide[i+1] || s.back.Text != labels[0] || s.next.Text != labels[1] {
			t.Fatalf("step %d: back=%q next=%q", i+2, s.back.Text, s.next.Text)
		}
	}
	test.Tap(s.next)
	if s.tabs.Selected() != s.archivesTab || s.next.Visible() || s.back.Text != "Back: Save to" {
		t.Fatalf("last step: next visible=%v back=%q", s.next.Visible(), s.back.Text)
	}
	test.Tap(s.back)
	test.Tap(s.back)
	test.Tap(s.back)
	if s.tabs.Selected() != s.sourceTab || s.back.Visible() {
		t.Fatal("Back did not return to the first step")
	}
	for _, tab := range []*container.TabItem{s.transfersTab, s.jobsTab, s.diagnosticsTab} {
		s.tabs.Select(tab)
		if s.back.Visible() || s.next.Visible() {
			t.Fatalf("guide buttons stayed visible on %q", tab.Text)
		}
	}
}

func TestJobsTabFollowsSource(t *testing.T) {
	s, _ := newTestScreen(t, nil)
	if !s.jobsTab.Disabled() {
		t.Fatal("Jobs tab is usable for Drive")
	}
	s.source.SetSelected(portabilitySourceLabel)
	if s.jobsTab.Disabled() {
		t.Fatal("Jobs tab is disabled for Data Portability")
	}
	s.source.SetSelected(driveSourceLabel)
	if !s.jobsTab.Disabled() {
		t.Fatal("Jobs tab stayed enabled after returning to Drive")
	}
}

// partServer serves the same small zip for every part path.
func partServer(t *testing.T) (*httptest.Server, []byte) {
	t.Helper()
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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/zip")
		http.ServeContent(w, r, "part.zip", time.Time{}, bytes.NewReader(content))
	}))
	t.Cleanup(server.Close)
	return server, content
}

func discoverThreeParts(t *testing.T, s *screen, account *screenAccount, server *httptest.Server, size int64) {
	t.Helper()
	account.client = server.Client()
	account.list = func(context.Context) ([]model.Export, error) {
		export := model.Export{ID: "multi", CreatedAt: time.Date(2026, 10, 25, 10, 0, 0, 0, time.UTC)}
		for i := 1; i <= 3; i++ {
			export.Parts = append(export.Parts, model.Part{ExportID: "multi", Index: i, Filename: fmt.Sprintf("part-%d.zip", i), Size: size, URL: fmt.Sprintf("%s/%d", server.URL, i)})
		}
		return []model.Export{export}, nil
	}
	s.discover()
	awaitJob(t, s)
}

func TestUntickedPartsAreNotDownloaded(t *testing.T) {
	s, account := newTestScreen(t, func(cfg *config.Settings) { cfg.SafetyBuffer = 1 })
	server, content := partServer(t)
	discoverThreeParts(t, s, account, server, int64(len(content)))
	exportCheck := s.selectionChecks[0]
	if len(s.selectionChecks) != 4 || !exportCheck.Checked || exportCheck.Partial || !strings.Contains(s.selectionSummary.Text, "Selected 3 of 3 archive files") {
		t.Fatalf("fresh discovery selection: %d checks, summary %q", len(s.selectionChecks), s.selectionSummary.Text)
	}
	test.Tap(s.selectionChecks[2]) // part-2.zip
	if !exportCheck.Partial || exportCheck.Checked || !s.selected["multi"] || !strings.HasPrefix(s.selectionSummary.Text, "Selected 2 of 3 archive files") {
		t.Fatalf("one unticked part: partial=%v checked=%v summary %q", exportCheck.Partial, exportCheck.Checked, s.selectionSummary.Text)
	}
	s.tabs.Select(s.archivesTab)
	captureScreen(t, s, "gui-step-4-selection.png")
	s.download()
	awaitJob(t, s)
	exists := func(name string) bool {
		_, err := os.Stat(filepath.Join(s.a.Settings.Destination, filepath.FromSlash(model.NamespacePath("multi", name))))
		return err == nil
	}
	if !exists("part-1.zip") || exists("part-2.zip") || !exists("part-3.zip") {
		t.Fatalf("downloaded parts: 1=%v 2=%v 3=%v", exists("part-1.zip"), exists("part-2.zip"), exists("part-3.zip"))
	}
}

func TestExportToggleSelectsAllPartsAndSourceEditsResetSkips(t *testing.T) {
	s, account := newTestScreen(t, nil)
	server, _ := partServer(t)
	discoverThreeParts(t, s, account, server, 1_000_000_000)
	exportCheck := s.selectionChecks[0]
	test.Tap(exportCheck) // all -> none
	if s.selected["multi"] || exportCheck.Checked || exportCheck.Partial || len(s.skippedParts) != 3 || s.selectionChecks[1].Checked {
		t.Fatal("unticking the export did not skip every part")
	}
	test.Tap(s.selectionChecks[3])
	if !exportCheck.Partial || !s.selected["multi"] || !strings.Contains(s.selectionSummary.Text, "Selected 1 of 3 archive files  ·  1.00 GB") {
		t.Fatalf("one ticked part: partial=%v summary %q", exportCheck.Partial, s.selectionSummary.Text)
	}
	test.Tap(exportCheck) // partial -> all
	if !exportCheck.Checked || len(s.skippedParts) != 0 {
		t.Fatal("ticking a partial export did not select every part")
	}
	test.Tap(s.selectionChecks[1])
	s.email.SetText("other@example.com")
	if len(s.skippedParts) != 0 || len(s.selectionChecks) != 0 || s.selectionRow.Visible() || s.discoveryBound {
		t.Fatal("a source edit kept a stale part skip")
	}
}

// Drive delivers each file as its own one-part export.
func discoverDriveFiles(t *testing.T, s *screen, account *screenAccount, server *httptest.Server, size int64, names ...string) {
	t.Helper()
	account.client = server.Client()
	account.list = func(context.Context) ([]model.Export, error) {
		var exports []model.Export
		for _, name := range names {
			exports = append(exports, model.Export{ID: name, CreatedAt: time.Date(2026, 10, 25, 10, 0, 0, 0, time.UTC), Parts: []model.Part{{ExportID: name, Index: 1, Filename: name, Size: size, URL: server.URL + "/" + name}}})
		}
		return exports, nil
	}
	s.discover()
	awaitJob(t, s)
}

func TestDriveFilesGetOneCheckAndSelectAllOrNone(t *testing.T) {
	s, account := newTestScreen(t, nil)
	server, _ := partServer(t)
	discoverDriveFiles(t, s, account, server, 1_000_000_000, "a.zip", "b.zip")
	if len(s.selectionChecks) != 2 {
		t.Fatalf("two one-part exports rendered %d checks", len(s.selectionChecks))
	}
	s.tabs.Select(s.archivesTab)
	captureScreen(t, s, "gui-step-4-drive-files.png")
	s.selectAllParts(false)
	if len(s.selectedParts()) != 0 || s.selectionChecks[0].Checked || !strings.Contains(s.selectionSummary.Text, "Selected 0 of 2") {
		t.Fatalf("select none left %q", s.selectionSummary.Text)
	}
	test.Tap(s.selectionChecks[1])
	if parts := s.selectedParts(); len(parts) != 1 || parts[0].Filename != "b.zip" || s.selected["a.zip"] || !s.selected["b.zip"] {
		t.Fatalf("ticking one file selected %v", parts)
	}
	s.selectAllParts(true)
	if len(s.selectedParts()) != 2 || !s.selectionChecks[0].Checked {
		t.Fatal("select all did not restore every file")
	}
	s.setBusy(true)
	s.selectAllParts(false)
	s.setBusy(false)
	if len(s.selectedParts()) != 2 {
		t.Fatal("selection changed while a job was running")
	}
}

func TestLaterSessionStartsFromWhatIsLeft(t *testing.T) {
	s, account := newTestScreen(t, func(cfg *config.Settings) { cfg.SafetyBuffer = 1 })
	server, content := partServer(t)
	discoverThreeParts(t, s, account, server, int64(len(content)))
	test.Tap(s.selectionChecks[2]) // keep part-2.zip for later
	s.download()
	awaitJob(t, s)
	s.discover()
	awaitJob(t, s)
	if parts := s.selectedParts(); len(parts) != 1 || parts[0].Filename != "part-2.zip" {
		t.Fatalf("second discovery selected %v", parts)
	}
	if !s.selectionChecks[0].Partial || s.selectionChecks[1].Checked || !s.selectionChecks[2].Checked {
		t.Fatal("downloaded parts were not shown as already handled")
	}
	test.Tap(s.selectionChecks[1])
	if len(s.selectedParts()) != 2 {
		t.Fatal("a downloaded part could not be ticked again")
	}
}

func TestChangingDestinationRederivesDownloadedParts(t *testing.T) {
	s, account := newTestScreen(t, func(cfg *config.Settings) { cfg.SafetyBuffer = 1 })
	server, content := partServer(t)
	discoverDriveFiles(t, s, account, server, int64(len(content)), "a.zip")
	first := s.destination.Text
	s.download()
	awaitJob(t, s)
	s.discover()
	awaitJob(t, s)
	if len(s.selectedParts()) != 0 {
		t.Fatal("completed file was selected again at its own destination")
	}
	s.destination.SetText(filepath.Join(t.TempDir(), "elsewhere"))
	if len(s.selectedParts()) != 1 {
		t.Fatal("a file absent from the new destination stayed unticked")
	}
	s.destination.SetText(first)
	if len(s.selectedParts()) != 0 {
		t.Fatal("returning to the first destination lost its completed state")
	}
}

func TestChangingDestinationKeepsManualUnticks(t *testing.T) {
	s, account := newTestScreen(t, nil)
	server, _ := partServer(t)
	discoverThreeParts(t, s, account, server, 1000)
	test.Tap(s.selectionChecks[2]) // part-2.zip
	s.destination.SetText(s.destination.Text + "-copy")
	parts := s.selectedParts()
	if len(parts) != 2 || parts[0].Filename != "part-1.zip" || parts[1].Filename != "part-3.zip" || s.selectionChecks[2].Checked || !s.selectionChecks[0].Partial {
		t.Fatalf("destination edit changed manual selection: %v", parts)
	}
}

func TestKeepWatchingOpensWithoutGrowingTheWindow(t *testing.T) {
	s, _ := newTestScreen(t, nil)
	s.w.Resize(fyne.NewSize(800, 640))
	s.tabs.Select(s.archivesTab)
	s.showWatchDialog()
	if minimum := s.content.MinSize(); minimum.Height > 640 {
		t.Fatalf("watch controls grew the page to %v", minimum)
	}
	closeButton := findDialogButton(s.w.Canvas().Overlays().Top(), "Close")
	if closeButton == nil {
		t.Fatal("watch dialog has no Close button")
	}
	position := fyne.CurrentApp().Driver().AbsolutePositionForObject(closeButton)
	if position.Y < 0 || position.Y+closeButton.Size().Height > 640 {
		t.Fatalf("Close button at %v does not fit the window", position)
	}
	captureScreen(t, s, "gui-keep-watching.png")
	s.watchDialog.Hide()
}

func TestDriveIntervalDoesNotBlockDataPortability(t *testing.T) {
	s, _ := newTestScreen(t, func(cfg *config.Settings) { cfg.PollIntervalSeconds = 60 })
	s.source.SetSelected(portabilitySourceLabel)
	if _, err := s.preferences(); err != nil || s.pollInterval.Text != "300" {
		t.Fatalf("interval %q after switching source: %v", s.pollInterval.Text, err)
	}
}

func TestExpiredAuthorizationReturnsToSignIn(t *testing.T) {
	s, account := newTestScreen(t, nil)
	account.list = func(context.Context) ([]model.Export, error) {
		return nil, fmt.Errorf("list: %w", auth.ErrAuthorizationRequired)
	}
	s.tabs.Select(s.archivesTab)
	s.discover()
	awaitJob(t, s)
	if s.tabs.Selected() != s.authorizeTab {
		t.Fatalf("authorization failure stayed on %q", s.tabs.Selected().Text)
	}
	if !dialogHasText(s.w.Canvas().Overlays().Top(), "Authorize in browser on step 2") {
		t.Fatal("authorization failure did not say how to sign in")
	}
	tapDialogButton(t, s, "Close")
	s.tabs.Select(s.transfersTab)
	s.errorDialog(fmt.Errorf("transfer: %w", downloader.ErrReauthentication))
	if s.tabs.Selected() != s.authorizeTab {
		t.Fatal("a transfer needing sign-in did not return to the Sign in step")
	}
}

func TestStepsRenderForCapture(t *testing.T) {
	s, _ := newTestScreen(t, nil)
	for _, shot := range []struct {
		tab  *container.TabItem
		name string
	}{{s.sourceTab, "gui-step-1-drive.png"}, {s.authorizeTab, "gui-step-2-drive.png"}, {s.destinationTab, "gui-step-3.png"}, {s.archivesTab, "gui-step-4-drive.png"}} {
		s.tabs.Select(shot.tab)
		assertWindowFits(t, s, fyne.NewSize(960, 760))
		captureScreen(t, s, shot.name)
	}
	s.source.SetSelected(portabilitySourceLabel)
	s.resources.SetText(strings.Join(source.Resources()[:3], ", "))
	s.tabs.Select(s.sourceTab)
	captureScreen(t, s, "gui-step-1-portability.png")
}
