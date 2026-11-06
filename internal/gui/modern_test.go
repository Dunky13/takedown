//go:build !headless

package gui

import (
	"context"
	"errors"
	"fmt"
	"image/color"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Dunky13/takedown/internal/events"
	"github.com/Dunky13/takedown/internal/model"
	"github.com/Dunky13/takedown/internal/source"
	"github.com/Dunky13/takedown/internal/state"
	"github.com/ncruces/zenity"
)

// variantTheme pins appTheme to one variant, because the test app has no OS preference.
type variantTheme struct {
	appTheme
	variant fyne.ThemeVariant
}

func (v variantTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	return v.appTheme.Color(name, v.variant)
}

func TestModernCapturesInBothThemes(t *testing.T) {
	for _, variant := range []struct {
		name string
		v    fyne.ThemeVariant
	}{{"light", theme.VariantLight}, {"dark", theme.VariantDark}} {
		t.Run(variant.name, func(t *testing.T) {
			s, account := newTestScreen(t, nil)
			fyne.CurrentApp().Settings().SetTheme(variantTheme{variant: variant.v})
			server, content := partServer(t)
			discoverThreeParts(t, s, account, server, int64(len(content)))
			shoot := func(name string, size fyne.Size, tab *container.TabItem) {
				s.w.Resize(size)
				s.tabs.Select(tab)
				s.content.Refresh()
				assertWindowFits(t, s, size)
				captureScreen(t, s, "modern-"+name+"-"+variant.name+".png")
			}
			large, small := fyne.NewSize(960, 760), fyne.NewSize(800, 640)
			shoot("transfers-empty", large, s.transfersTab)
			s.onEvent(events.Event{Kind: events.Progress, PartKey: "multi/1", Filename: "part-1.zip", Bytes: 50_000_000, Total: 100_000_000, BytesPerSecond: 10_000_000})
			shoot("1-choose", large, s.sourceTab)
			shoot("2-signin", large, s.authorizeTab)
			shoot("3-save", large, s.destinationTab)
			shoot("4-download", large, s.archivesTab)
			shoot("transfers", large, s.transfersTab)
			shoot("diagnostics", large, s.diagnosticsTab)
			s.diagnosticsTabs.SelectIndex(1)
			shoot("startup", large, s.diagnosticsTab)
			s.diagnosticsTabs.SelectIndex(0)
			shoot("1-choose-800", small, s.sourceTab)
			shoot("4-download-800", small, s.archivesTab)
			s.source.SetSelected(portabilitySourceLabel)
			s.a.Account = &managedAccount{screenAccount: account} // lets Jobs show its real, enabled state
			s.jobsUI.setBusy(false)
			shoot("jobs", large, s.jobsTab)
			shoot("1-choose-portability", large, s.sourceTab)
			shoot("1-choose-portability-800", small, s.sourceTab)
			shoot("3-save-portability", large, s.destinationTab)
			s.destination.SetText("sftp://me@nas.local/archives")
			s.unknown.SetChecked(true)
			shoot("3-save-sftp", large, s.destinationTab)
			shoot("3-save-sftp-800", small, s.destinationTab)
			s.banner.SetText("Free space, then retry.")
			s.banner.Show()
			s.disclosureAccordion.Open(0)
			shoot("2-signin-open", large, s.authorizeTab)
			s.tabs.Select(s.sourceTab)
			test.Tap(s.driveCard)
			if s.w.Canvas().Focused() != nil || s.driveCard.focused {
				t.Fatal("a click left focus on the card")
			}
			s.w.Canvas().Focus(s.portabilityCard) // keyboard focus
			if !s.portabilityCard.focused {
				t.Fatal("keyboard focus did not reach the card")
			}
			shoot("1-choose-focus", large, s.sourceTab)
		})
	}
}

func TestNavigatorKeepsMinSizeAndRefusesDisabledPages(t *testing.T) {
	s, _ := newTestScreen(t, nil)
	want := s.tabs.stack.MinSize()
	for _, item := range []*container.TabItem{s.sourceTab, s.authorizeTab, s.destinationTab, s.archivesTab, s.transfersTab, s.diagnosticsTab} {
		s.tabs.Select(item)
		if s.tabs.Selected() != item || s.tabs.stack.MinSize() != want {
			t.Fatalf("%q: selected=%v, minimum %v, want %v on every page", item.Text, s.tabs.Selected() == item, s.tabs.stack.MinSize(), want)
		}
	}
	shown := 0
	for i, page := range s.tabs.stack.Objects {
		if page.Visible() {
			shown++
			if s.tabs.items[i] != s.tabs.Selected() {
				t.Fatal("a page other than the selected one is shown")
			}
		}
	}
	if shown != 1 {
		t.Fatalf("%d pages shown, want 1", shown)
	}
	// Drive has no Jobs page to open.
	selections := 0
	s.tabs.OnSelected = func(*container.TabItem) { selections++ }
	before := s.tabs.Selected()
	s.tabs.Select(s.jobsTab)
	if s.tabs.Selected() != before || selections != 0 || s.jobsTab.Content.Visible() {
		t.Fatal("a disabled page was selected")
	}
	if !s.toolButtons[s.jobsTab].Disabled() {
		t.Fatal("the Jobs button did not follow the disabled page")
	}
	s.tabs.EnableItem(s.jobsTab)
	s.tabs.Select(s.jobsTab)
	if s.tabs.Selected() != s.jobsTab || selections != 1 || s.toolButtons[s.jobsTab].Disabled() {
		t.Fatal("an enabled page could not be selected")
	}
}

func TestChoiceCardsSetSourceAndAreInertWhileBusy(t *testing.T) {
	s, _ := newTestScreen(t, nil)
	if !s.driveCard.selected || s.portabilityCard.selected {
		t.Fatal("cards do not show the Drive source")
	}
	test.Tap(s.portabilityCard)
	if s.source.Selected != portabilitySourceLabel || !s.portabilityCard.selected || s.driveCard.selected {
		t.Fatalf("tapping Data Portability left source %q", s.source.Selected)
	}
	s.setBusy(true)
	if !s.driveCard.Disabled() || !s.portabilityCard.Disabled() {
		t.Fatal("cards stayed enabled while busy")
	}
	test.Tap(s.driveCard) // taking focus first, as a click does
	s.driveCard.TypedKey(&fyne.KeyEvent{Name: fyne.KeySpace})
	s.driveCard.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	if s.source.Selected != portabilitySourceLabel || s.driveCard.selected {
		t.Fatal("a focused card changed the source while busy")
	}
	s.setBusy(false)
	s.driveCard.TypedKey(&fyne.KeyEvent{Name: fyne.KeySpace})
	if s.source.Selected != driveSourceLabel || !s.driveCard.selected {
		t.Fatal("Space did not select the card")
	}
	s.portabilityCard.TypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	if s.source.Selected != portabilitySourceLabel {
		t.Fatal("Enter did not select the card")
	}
}

func TestNextIsPrimaryExceptOnSignIn(t *testing.T) {
	s, _ := newTestScreen(t, nil)
	for _, step := range []struct {
		tab        *container.TabItem
		importance widget.Importance
	}{{s.sourceTab, widget.HighImportance}, {s.authorizeTab, widget.MediumImportance}, {s.destinationTab, widget.HighImportance}} {
		s.tabs.Select(step.tab)
		if s.next.Importance != step.importance || s.stepper.current != s.tabs.selected {
			t.Fatalf("%q: Next importance %v, stepper at %d", step.tab.Text, s.next.Importance, s.stepper.current)
		}
	}
	s.tabs.Select(s.transfersTab)
	if s.stepper.current != -1 {
		t.Fatal("a tool page marked a guide step as current")
	}
}

func TestCalloutShowsAndHides(t *testing.T) {
	s, _ := newTestScreen(t, nil)
	if s.banner.Visible() {
		t.Fatal("banner starts visible")
	}
	hidden := s.content.MinSize()
	s.banner.SetText("Authorize your account in your system browser.")
	s.banner.Show()
	s.content.Refresh()
	if !s.banner.Visible() || s.banner.Text == "" || s.banner.MinSize().Height <= 0 || s.content.MinSize().Height <= hidden.Height {
		t.Fatal("a shown callout took no space")
	}
	s.banner.Hide()
	s.content.Refresh()
	if s.banner.Visible() || s.content.MinSize() != hidden {
		t.Fatalf("a hidden callout kept its space: %v, want %v", s.content.MinSize(), hidden)
	}
}

// TestModernDialogsFitInBothThemes captures each dialog and checks it fits its window.
func TestModernDialogsFitInBothThemes(t *testing.T) {
	for _, variant := range []struct {
		name string
		v    fyne.ThemeVariant
	}{{"light", theme.VariantLight}, {"dark", theme.VariantDark}} {
		t.Run(variant.name, func(t *testing.T) {
			for _, size := range []fyne.Size{fyne.NewSize(960, 760), fyne.NewSize(800, 640)} {
				// A fresh screen per size: a dialog removed behind its own back will not reopen.
				s, base := newTestScreen(t, nil)
				a := &managedAccount{screenAccount: base}
				a.browse = func(context.Context, string) ([]source.Folder, error) {
					return []source.Folder{{ID: "f1", Name: "Takeout archives"}, {ID: "f2", Name: "gui-long-drive-folder-name-that-keeps-going-and-going-and-going"}}, nil
				}
				fyne.CurrentApp().Settings().SetTheme(variantTheme{variant: variant.v})
				s.w.Resize(size)
				suffix := variant.name
				if size.Width < 900 {
					suffix += "-800"
				}
				shoot := func(name string) {
					top := s.w.Canvas().Overlays().Top()
					if top == nil {
						t.Fatalf("%s: no dialog open", name)
					}
					if top.Size().Width > size.Width || top.Size().Height > size.Height {
						t.Fatalf("%s: dialog %v exceeds window %v", name, top.Size(), size)
					}
					captureScreen(t, s, "modern-dialog-"+name+"-"+suffix+".png")
					s.w.Canvas().Overlays().Remove(top)
				}
				s.oauthClient.SetText("/tmp/client.json") // may rebuild the account, so swap it after
				s.login()
				shoot("authorize")
				s.showWatchDialog()
				shoot("watch")
				s.a.Account = a
				s.browseFolders()
				awaitJob(t, s)
				shoot("folders")
				s.errorDialog(errors.New("googleapi: Error 403: the request is missing a valid API key, which is a long message that must wrap inside the dialog without clipping"))
				shoot("error")
			}
		})
	}
}

func TestStepperPaintsPositionNeverCompletion(t *testing.T) {
	s, _ := newTestScreen(t, nil)
	r, ok := test.WidgetRenderer(s.stepper).(*stepperRenderer)
	if !ok {
		t.Fatal("stepper renderer has an unexpected type")
	}
	for at, step := range s.guide {
		s.tabs.Select(step)
		for i, b := range s.stepper.buttons {
			want := numberIcon(i+1, numberGrey).Name()
			if i == at {
				want = numberIcon(i+1, theme.Color(theme.ColorNamePrimary)).Name()
			}
			if b.Icon.Name() != want {
				t.Errorf("on step %d, badge %d is %s, want %s", at+1, i+1, b.Icon.Name(), want)
			}
		}
		for i, c := range r.connectors {
			if c.FillColor != theme.Color(theme.ColorNameSeparator) {
				t.Errorf("on step %d, connector %d is %v, want the neutral separator", at+1, i, c.FillColor)
			}
		}
	}
}

func TestPartByPartStaysVisibleBeforeDiscovery(t *testing.T) {
	s, account := newTestScreen(t, nil)
	s.tabs.Select(s.archivesTab)
	inRow := func(row *fyne.Container, obj fyne.CanvasObject) bool {
		for _, o := range row.Objects {
			if o == obj {
				return true
			}
		}
		return false
	}
	if !s.partByPart.Visible() || inRow(s.downloadRow, s.partByPart) || s.downloadRow.Visible() {
		t.Fatal("before discovery the part-by-part choice must show and Download selected must not")
	}
	server, content := partServer(t)
	discoverThreeParts(t, s, account, server, int64(len(content)))
	if !s.partByPart.Visible() || !s.downloadRow.Visible() {
		t.Fatal("after discovery both the choice and Download selected must show")
	}
}

func TestLongErrorKeepsEveryPageWithinTheWindow(t *testing.T) {
	sentence := "The remote service refused the request and said why at length. "
	for _, length := range []int{340, 500, 2000} {
		long := strings.Repeat(sentence, length/len(sentence)+1)[:length]
		for _, source := range []string{driveSourceLabel, portabilitySourceLabel} {
			t.Run(fmt.Sprintf("%d/%s", length, source), func(t *testing.T) {
				s, _ := newTestScreen(t, nil)
				s.source.SetSelected(source)
				size := fyne.NewSize(800, 640)
				s.w.Resize(size)
				s.run("Working…", func(context.Context, sourceSnapshot) error { return errors.New(long) }, func(error, sourceSnapshot) {}, nil)
				awaitJob(t, s)
				if !s.banner.Visible() || len([]rune(s.banner.Text)) > bannerLimit+1 || !strings.HasPrefix(long, strings.TrimSuffix(s.banner.Text, "…")) {
					t.Fatalf("banner %q (visible %v)", s.banner.Text, s.banner.Visible())
				}
				if length <= bannerLimit && s.banner.Text != long {
					t.Fatalf("a short message was altered: %q", s.banner.Text)
				}
				assertStepsFit(t, s, size)
				for _, tab := range []*container.TabItem{s.transfersTab, s.diagnosticsTab} {
					s.tabs.Select(tab)
					assertWindowFits(t, s, size)
				}
			})
		}
	}
}

func TestFooterStatusStaysOneLineWhateverItSays(t *testing.T) {
	s, _ := newTestScreen(t, nil)
	size := fyne.NewSize(800, 640)
	s.w.Resize(size)
	s.tabs.Select(s.sourceTab)
	footerHeight := func() float32 {
		objects := s.content.(*fyne.Container).Objects
		return objects[len(objects)-1].Size().Height
	}
	s.status.SetText("Short.")
	s.content.Refresh()
	want := footerHeight()
	for _, length := range []int{77, 157, 600} {
		s.status.SetText(strings.Repeat("Waiting for the destination to answer. ", length/38+1)[:length])
		s.content.Refresh()
		assertWindowFits(t, s, size)
		if got := footerHeight(); got != want {
			t.Errorf("%d characters of status: footer %v high, want %v", length, got, want)
		}
	}
	s.notice("Source changed. Stopping work; partial files will be kept.")
	if s.banner.Text == s.status.Text || strings.Contains(s.status.Text, "Source changed") {
		t.Errorf("notice repeated in the footer: banner %q, status %q", s.banner.Text, s.status.Text)
	}
}

func TestFooterKeepsItsHeightAcrossPages(t *testing.T) {
	s, _ := newTestScreen(t, nil)
	height := func() float32 {
		return s.content.(*fyne.Container).Objects[len(s.content.(*fyne.Container).Objects)-1].MinSize().Height
	}
	s.tabs.Select(s.sourceTab)
	want := height()
	for _, tab := range []*container.TabItem{s.authorizeTab, s.archivesTab, s.transfersTab, s.diagnosticsTab} {
		s.tabs.Select(tab)
		if got := height(); got != want {
			t.Errorf("%s: footer minimum height %v, want %v", tab.Text, got, want)
		}
	}
	if s.next.Size().Height != s.next.MinSize().Height {
		t.Errorf("Next stretched to %v, natural %v", s.next.Size().Height, s.next.MinSize().Height)
	}
}

func TestKeyboardOrderFollowsTheScreen(t *testing.T) {
	s, _ := newTestScreen(t, nil)
	s.tabs.Select(s.sourceTab)
	canvas := s.w.Canvas()
	walker, ok := canvas.(interface{ FocusNext() })
	if !ok {
		t.Fatal("the test canvas cannot walk focus")
	}
	canvas.Unfocus()
	order := map[fyne.Focusable]int{}
	for i := 0; i < 40; i++ {
		walker.FocusNext()
		if f := canvas.Focused(); f != nil {
			if _, seen := order[f]; seen {
				break
			}
			order[f] = i
		}
	}
	at := func(name string, f fyne.Focusable) int {
		i, ok := order[f]
		if !ok {
			t.Fatalf("%s was never focused", name)
		}
		return i
	}
	toolbar := at("Transfers", s.toolButtons[s.transfersTab])
	step := at("stepper", s.stepper.buttons[0])
	cards := at("Drive card", s.driveCard)
	takeout := at("Open Google Takeout", s.openTakeout)
	footer := at("Save preferences", s.savePrefs)
	if !(toolbar < step && step < cards && cards < takeout && takeout < footer) {
		t.Fatalf("focus order toolbar=%d stepper=%d cards=%d takeout=%d footer=%d", toolbar, step, cards, takeout, footer)
	}
	for f, i := range order {
		if i > footer && f != fyne.Focusable(s.next) {
			t.Errorf("%T focused after the footer", f)
		}
	}
}

func TestKeepWatchingDialogHasOnePrimary(t *testing.T) {
	s, _ := newTestScreen(t, nil)
	s.showWatchDialog()
	primaries := 0
	var walk func(fyne.CanvasObject)
	walk = func(obj fyne.CanvasObject) {
		switch o := obj.(type) {
		case *widget.Button:
			if o.Importance == widget.HighImportance {
				primaries++
			}
		case *button:
			if o.Importance == widget.HighImportance {
				primaries++
			}
		case *widget.PopUp:
			walk(o.Content)
		case *fyne.Container:
			for _, c := range o.Objects {
				walk(c)
			}
		case fyne.Widget:
			for _, c := range test.WidgetRenderer(o).Objects() {
				walk(c)
			}
		}
	}
	walk(s.w.Canvas().Overlays().Top())
	if primaries != 1 || s.watch.Importance != widget.HighImportance {
		t.Fatalf("%d primary buttons in the Keep watching dialog, Start watching importance %v", primaries, s.watch.Importance)
	}
}

func TestEmptyHintsWrapInsideTheirPane(t *testing.T) {
	s, _ := newTestScreen(t, nil)
	s.source.SetSelected(portabilitySourceLabel)
	s.a.Account = &managedAccount{screenAccount: &screenAccount{}}
	size := fyne.NewSize(620, 640)
	s.w.Resize(size)
	s.tabs.Select(s.jobsTab)
	s.content.Refresh()
	pane := s.jobsUI.list.Size().Width
	label := s.jobsUI.empty
	right := fyne.CurrentApp().Driver().AbsolutePositionForObject(label).X + label.Size().Width
	left := fyne.CurrentApp().Driver().AbsolutePositionForObject(s.jobsUI.list).X
	if label.Size().Width < 60 || right > left+pane {
		t.Fatalf("jobs hint is %v wide and ends at %v, the list pane ends at %v", label.Size().Width, right, left+pane)
	}
	if label.MinSize().Height < 2*theme.Size(theme.SizeNameCaptionText) {
		t.Fatalf("jobs hint did not wrap: height %v", label.MinSize().Height)
	}
}

func TestTransfersShowSavedProgressAndOnlyRelevantControls(t *testing.T) {
	s, _ := newTestScreen(t, nil)
	if s.transferControls.Visible() || !s.transferEmpty.Visible() {
		t.Fatal("an empty Transfers page shows controls")
	}
	// Sign-in, discovery and similar work is not a transfer.
	s.setBusy(true)
	if s.transferControls.Visible() || s.pauseRow.Visible() {
		t.Fatal("unrelated busy work showed the transfer controls")
	}
	s.setBusy(false)

	store, err := state.Open(filepath.Join(s.a.DataDir, "takedown-state.json"))
	if err != nil {
		t.Fatal(err)
	}
	part := model.Part{ExportID: "e", Index: 1, Filename: "a.zip", Size: 200, URL: "https://example.invalid/a"}
	if err := store.Put(state.Record{Part: part, Destination: "/dest", Offset: 50}); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	s.loadSavedTransfers()
	if s.progress.Value != 0.25 || s.progressText.Text != "25%" {
		t.Fatalf("saved files show progress %v (%q), want 25%%", s.progress.Value, s.progressText.Text)
	}
	if !s.transferControls.Visible() || s.transferEmpty.Visible() || s.pauseRow.Visible() {
		t.Fatal("saved files should show progress without Pause and Resume while nothing runs")
	}
	s.startTransfers()
	if !s.pauseRow.Visible() {
		t.Fatal("a running download hides Pause and Resume")
	}
	s.setBusy(false)
	if s.pauseRow.Visible() {
		t.Fatal("Pause and Resume stayed after the work ended")
	}
}

func TestProgressBarIsSlim(t *testing.T) {
	s, _ := newTestScreen(t, nil)
	s.onEvent(events.Event{Kind: events.Progress, PartKey: "multi/1", Filename: "part-1.zip", Bytes: 50, Total: 100})
	s.tabs.Select(s.transfersTab)
	if h := s.progress.Size().Height; h > 12 {
		t.Fatalf("progress bar is %v high, want a slim bar", h)
	}
}

func TestErrorEventsAndNoticesReachTheBanner(t *testing.T) {
	s, _ := newTestScreen(t, nil)
	s.setBusy(true)
	for _, e := range []events.Event{
		{Kind: events.Error, Message: "Pause could not flush download checkpoints"},
		{Kind: events.Warning, Message: "An export job is held for review; other jobs remain eligible and watching will continue"},
		{Kind: events.SpaceWarning, Message: "Not enough free space"},
	} {
		s.banner.Hide()
		s.onEvent(e)
		if !s.banner.Visible() || s.banner.Text != e.Message || strings.Contains(s.status.Text, e.Message) {
			t.Errorf("%q: banner %q (visible %v), status %q", e.Message, s.banner.Text, s.banner.Visible(), s.status.Text)
		}
	}
	s.banner.Hide()
	s.onEvent(events.Event{Kind: events.Status, Message: "Paused; partial archives checkpointed"})
	if s.banner.Visible() || s.status.Text != "Paused; partial archives checkpointed" {
		t.Error("a routine status went to the banner")
	}
}

func TestClipTextAndEndOfPath(t *testing.T) {
	if got := clipText("  short \n text ", 50); got != "short text" {
		t.Errorf("clipText collapsed to %q", got)
	}
	got := clipText(strings.Repeat("word ", 100), 20)
	if r := []rune(got); len(r) > 21 || !strings.HasSuffix(got, "…") || strings.Contains(got, "wor…") {
		t.Errorf("clipText cut mid-word or overran: %q", got)
	}
	if got := endOfPath("/a/very/long/path/to/the/archives/folder-name", 20); !strings.HasSuffix(got, "rchives/folder-name") || !strings.HasPrefix(got, "…") || len([]rune(got)) != 20 {
		t.Errorf("endOfPath dropped the folder name: %q", got)
	}
}

func TestSFTPFieldsHaveChoosersHelpAndHomeExpansion(t *testing.T) {
	s, _ := newTestScreen(t, nil)
	if s.sshKey.PlaceHolder == "" || s.knownHosts.PlaceHolder == "" {
		t.Fatal("SFTP key and known_hosts fields have no placeholder")
	}
	var seen []fileRequest
	s.openFileDialog = func(_ context.Context, r fileRequest) (string, error) {
		seen = append(seen, r)
		return "/chosen/id_ed25519", nil
	}
	s.pickSSHFile("Select your SFTP private key", s.sshKey)
	awaitUIState(t, func() bool { return !s.picking })
	if len(seen) != 1 || !seen[0].anyFile || seen[0].jsonOnly || s.sshKey.Text != "/chosen/id_ed25519" {
		t.Fatalf("requests %+v, key %q", seen, s.sshKey.Text)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("no home folder")
	}
	s.sshKey.SetText("~/.ssh/id_ed25519")
	s.knownHosts.SetText("")
	cfg, err := s.preferences()
	if err != nil || cfg.PrivateKeyPath != filepath.Join(home, ".ssh", "id_ed25519") || cfg.KnownHostsPath != "" {
		t.Fatalf("key %q, known_hosts %q, err %v", cfg.PrivateKeyPath, cfg.KnownHostsPath, err)
	}
}

func TestBrowseWhileADialogIsOpenSaysSo(t *testing.T) {
	s, _ := newTestScreen(t, nil)
	release := make(chan struct{})
	s.openFileDialog = func(context.Context, fileRequest) (string, error) { <-release; return "", zenity.ErrCanceled }
	s.pickDestination()
	s.pickDestination()
	if !strings.Contains(s.status.Text, "already open") {
		t.Errorf("status %q does not say a dialog is open", s.status.Text)
	}
	release <- struct{}{}
	awaitUIState(t, func() bool { return !s.picking })
}

func TestFileDialogFailureGivesFileDialogAdvice(t *testing.T) {
	s, _ := newTestScreen(t, nil)
	s.openFileDialog = func(context.Context, fileRequest) (string, error) { return "", errors.New("osascript failed") }
	s.pickDestination()
	awaitUIState(t, func() bool { return !s.picking })
	top := s.w.Canvas().Overlays().Top()
	if !dialogHasText(top, "Type the path into the field instead") || dialogHasText(top, "authorization") {
		t.Fatal("a failed file dialog shows advice about sources and authorization")
	}
}

func TestDiagnosticsTabsAreKeyboardReachable(t *testing.T) {
	s, _ := newTestScreen(t, nil)
	s.tabs.Select(s.diagnosticsTab)
	canvas := s.w.Canvas()
	walker := canvas.(interface{ FocusNext() })
	canvas.Unfocus()
	reached := map[fyne.Focusable]bool{}
	for i := 0; i < 40; i++ {
		walker.FocusNext()
		reached[canvas.Focused()] = true
	}
	for i, b := range s.diagnosticsTabs.buttons {
		if !reached[b] {
			t.Errorf("diagnostics tab %d cannot be reached with Tab", i)
		}
	}
	test.Tap(s.diagnosticsTabs.buttons[2])
	if s.diagnosticsTabs.nav.selected != 2 || s.diagnosticsTabs.buttons[2].Importance != widget.MediumImportance {
		t.Fatal("activating a diagnostics tab did not select it")
	}
}

// After a discovery the export list and Download row are on screen; a banner
// must still leave every page within a small window.
func TestDiscoveredStateWithBannerFitsASmallWindow(t *testing.T) {
	s, account := newTestScreen(t, nil)
	server, content := partServer(t)
	discoverThreeParts(t, s, account, server, int64(len(content)))
	size := fyne.NewSize(800, 640)
	s.w.Resize(size)
	for _, text := range []string{"Free space on the destination, then resume. The archive needs more room than is available.", strings.Repeat("A long warning that wraps over several lines. ", 10)} {
		s.notice(text)
		if !s.downloadRow.Visible() {
			t.Fatal("the Download row is not showing after discovery")
		}
		assertStepsFit(t, s, size)
		for _, tab := range []*container.TabItem{s.transfersTab, s.diagnosticsTab} {
			s.tabs.Select(tab)
			assertWindowFits(t, s, size)
		}
	}
}

func TestSSHBrowseOpensBesideTheValueEvenWithATilde(t *testing.T) {
	home := t.TempDir()
	if err := os.MkdirAll(filepath.Join(home, ".ssh"), 0o700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	s, _ := newTestScreen(t, nil)
	var seen []fileRequest
	s.openFileDialog = func(_ context.Context, r fileRequest) (string, error) {
		seen = append(seen, r)
		return "", zenity.ErrCanceled
	}
	for _, text := range []string{"~/.ssh/id_ed25519", "", "~/missing/key"} {
		s.sshKey.SetText(text)
		s.pickSSHFile("Select your SFTP private key", s.sshKey)
		awaitUIState(t, func() bool { return !s.picking })
	}
	if len(seen) != 3 {
		t.Fatalf("%d requests", len(seen))
	}
	// A missing folder falls back to its nearest existing parent, never the working directory.
	for i, want := range []string{filepath.Join(home, ".ssh"), filepath.Join(home, ".ssh"), home} {
		if seen[i].start != want {
			t.Errorf("request %d starts in %q, want %q", i, seen[i].start, want)
		}
	}
}

func TestFyneFallbackDialogOpensInTheStartFolder(t *testing.T) {
	s, _ := newTestScreen(t, nil)
	dir := filepath.Join(t.TempDir(), "locate-me-here")
	if err := os.Mkdir(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	s.destination.SetText(dir)
	s.openFileDialog = func(context.Context, fileRequest) (string, error) { return "", exec.ErrNotFound }
	s.pickDestination()
	awaitUIState(t, func() bool { return !s.picking })
	top := s.w.Canvas().Overlays().Top()
	if top == nil || !dialogHasText(top, "locate-me-here") {
		t.Fatal("the Fyne fallback did not open in the requested folder")
	}
}
