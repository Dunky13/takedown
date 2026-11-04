# PR #15: guided setup and partial downloads

PR: <https://github.com/Dunky13/takedown/pull/15>

## Outcome and scope

The desktop GUI is now a four-step guide in a left rail: Choose data, Sign in,
Save to and Download. Transfers, Jobs and Diagnostics follow in the same rail.
Back and Next only navigate; they do not save or block anything.

Partial importing works at two levels:

- **Choosing data:** a Takeout walkthrough, or the Data Portability catalog
  grouped by product with readable names. Exact API names remain under Advanced
  and in the access disclosure.
- **Choosing files:** individual archive files can be ticked. Parts complete at
  the current destination start unticked, so later sessions fetch what is left.

Unchanged: `accessDisclosure()` text, the pre-OAuth confirmation, source binding
and invalidation, job and diagnostics behaviour, and saved file formats. Outside
`internal/gui`, only one status string in `internal/source/source.go` and two
README lines changed.

## Behaviour worth knowing

- Drive delivers each file as its own one-part export, so Drive rows have one
  check. Only multi-part exports have a parent row with a mixed state.
- The Drive folder picker is on step 4 because browsing needs a token.
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
  sources.
- The real macOS binary was launched and captured on Drive and Data Portability
  with throwaway profiles. Navigation inside it was exercised only through
  Fyne's test driver, not real clicks.
- Set `TAKEDOWN_TEST_ARTIFACT_DIR` when running the GUI tests to regenerate
  screenshots. They stay out of Git.

## Remaining work

1. Run a live sign-in, Discover and partial download on real Drive and Data
   Portability accounts. No live Google flow has been exercised.
2. Check runtime rendering on Windows and Linux.
3. Attach before and after screenshots to the PR if wanted. No approved host was
   available when it was opened.
4. Consider later: step completion marks, which need a truthful local
   authorization check, and a free-space buffer entered in GB instead of bytes.
