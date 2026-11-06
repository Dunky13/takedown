# PR #15: guided setup and partial downloads

PR: <https://github.com/Dunky13/takedown/pull/15>

## Outcome and scope

The desktop GUI is now a four-step guide: a horizontal stepper over one centred
column, for Choose data, Sign in, Save to and Download. Transfers, Jobs and
Diagnostics sit in a top-right toolbar. Back and Next only navigate; they do not
save or block anything. A custom theme provides tinted neutrals, one blue accent
and light and dark variants that follow the OS.

Save to and Sign in now use the OS file dialogs, through
`github.com/ncruces/zenity`, a dependency the owner's request for native dialogs
required.

Partial importing works at two levels:

- **Choosing data:** a Takeout walkthrough, or the Data Portability catalog
  grouped by product with readable names. Exact API names remain under Advanced
  and in the access disclosure.
- **Choosing files:** individual archive files can be ticked. Parts complete at
  the current destination start unticked, so later sessions fetch what is left.

Unchanged: `accessDisclosure()` text, the pre-OAuth confirmation, source binding
and invalidation, and saved file formats. Outside `internal/gui`, the changes are
one status string in `internal/source/source.go`, a new `events.Warning` kind in
`internal/events/events.go` used for held jobs in `internal/app/watch.go` (with
its test in `internal/app/job_watch_test.go`), README lines and the dependency
files.

## Behaviour worth knowing

- Drive delivers each file as its own one-part export, so Drive rows have one
  check. Only multi-part exports have a parent row with a mixed state.
- The Drive folder picker is on step 4 because browsing needs a token. It lists
  Google Drive, so it stays an in-app dialog and has no OS equivalent.
- The stepper shows position only: the accent marks the current step. It never
  implies a step is complete, because sign-in state is not truthfully known.
- Linux needs `zenity`, `qarma` or `matedialog` for native file dialogs. Without
  one, `pickFile` falls back to Fyne's own dialog.
- On macOS the dialog runs in `osascript`, so it is not attached to the takedown
  window as a sheet.
- Focus: Fyne's `Button.Tapped` calls `Canvas.Focus(nil)`, so a clicked button
  never keeps focus. Keyboard focus shows as a 2 px ring that the app's own
  `button` type draws (`widgets.go`), because Fyne's focus tint is blended into
  the fill and cannot show on a primary button. Choice cards also drop focus on
  click. Entries draw the same ring (`entry`). Every app dialog, including the
  confirmations, uses `dialogWithButtons`, so its buttons show the ring too.
- Messages that affect safety (errors, space and sign-in warnings, held jobs,
  stale confirmations) go to the banner, which is capped at 200 characters. The
  footer is one line with an ellipsis and only points at the banner. The error
  dialog has the full text and Copy.
- Transfers shows saved progress after a restart. Pause and Resume show only
  while a download or watch runs.
- SFTP private key and known_hosts have Browse buttons (hidden files shown, start
  in `~/.ssh`), placeholders and help. A leading `~/` is expanded when
  preferences are read. Empty known_hosts means `~/.ssh/known_hosts`.
- Diagnostics sub-tabs are buttons in the normal theme, so they take keyboard
  focus. The earlier AppTabs override recoloured every button inside them.
- Browse while an OS file dialog is open only reports that one is open: the
  dialog belongs to another process and cannot be raised.
- Editing the destination re-derives "downloaded earlier" skips and keeps manual
  unticks.
- `downloadedParts` trusts recorded completion. A file moved or deleted from the
  destination still shows as downloaded earlier; ticking it fetches it again.
- Watching cannot honour part selection. It downloads everything unless Only
  watch selected exports is on, and a partly ticked export downloads whole. The
  Keep watching dialog says so.
- Switching from Drive to Data Portability raises a poll interval below 300 to
  300, the Data Portability minimum.
- Authorization failures select the Sign in step. Messages from `internal/auth`
  still say "run Login", because the CLI shares them.

## Evidence

- `make check`, the full `-race -tags ci` suite and the headless build pass.
- Fit tests walk all four steps at 960x760 and 800x640, idle and busy, for both
  sources. A contrast test covers both theme variants.
- A real macOS binary was clicked through on both sources and both themes,
  before the final repair pass. The final tree renders through the test driver.
- zenity opened the real macOS folder panel from a small probe, not through the
  app's Browse button.
- Set `TAKEDOWN_TEST_ARTIFACT_DIR` when running the GUI tests to regenerate
  screenshots. They stay out of Git.

## Review repair pass

Fixes for the independent spec review: focus ring on every app button,
Diagnostics tab theming and keyboard access, SFTP key and known_hosts choosers,
Transfers progress and controls, banner cap, footer one-liner, visible warnings,
native-style dialog buttons (natural width, trailing edge, one primary), folder
picker rows with icons, panels on Jobs and results, triangle expanders on the
resource list, and truncation that keeps the end of the destination path.
Tests: `TestFocusIsVisibleAndKeepsTextReadable` (ring and label contrast over
real button fills), `TestButtonAndEntryShowKeyboardFocusNotClickFocus`,
`TestDiagnosticsButtonsKeepTheNormalTheme`, `TestLongErrorKeepsEveryPageWithinTheWindow`
(340, 500 and 2000 characters), `TestFooterStatusStaysOneLineWhateverItSays`,
`TestTransfersShowSavedProgressAndOnlyRelevantControls` and the other tests in
`modern_test.go`. The held-job watch message in `internal/app/watch.go` is now
a new `events.Warning` event: the CLI prints it as `warning: ...` where it
printed `status: ...` before, and the GUI banners it by kind (`TestWatchKeepsPollingWhenOnlyCurrentJobIsHeld`
pins the kind). The Drive storage note is a warning callout.

Known limits: entries still scroll long paths from the start, so the filename
of a very long path can sit out of view until the field is focused. At 800x640
the Data Portability catalog shows about 4 to 5 rows (tighter rows gained a
half row); a banner takes about one more.

## Remaining work

1. Run a live sign-in, Discover and partial download on real Drive and Data
   Portability accounts. No live Google flow has been exercised.
2. Check runtime rendering on Windows and Linux.
3. Attach before and after screenshots to the PR if wanted. No approved host was
   available when it was opened.
4. Try the Browse and Select JSON buttons on a signed and notarized build, and on
   Windows and Linux. Only the unsigned macOS panel has been seen.
5. Consider later: step completion marks, which need a truthful local
   authorization check, and a free-space buffer entered in GB instead of bytes.

## Test time under the race detector

Step changes and busy toggles used `Container.Refresh` on the whole window,
which refreshed every widget on every page and re-rendered every checkbox icon.
On the macOS Intel CI runner that pushed `internal/gui` past Go's 10-minute test
timeout under `-race`. `relayout` in `widgets.go` now re-runs only the container
layouts. Locally the GUI package under `-race` went from 580 s to 204 s.
