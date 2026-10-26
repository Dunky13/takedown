# takedown implementation handoff

> Historical report: this describes the original browser-automation implementation before the documented-API redesign. Current behavior and acceptance are recorded in [AUTOMATION.md](AUTOMATION.md). Some source links below refer to files removed by that redesign and are no longer available in the checkout.

All four implementation phases are present in the local checkout. The default
binary combines the CLI and native Fyne GUI. Live Google account compatibility
remains unqualified: the user explicitly reserved that acceptance check for
themselves. No commit, pull request, publication, or deployment was performed.

## Implementation

| Phase | Packages and behavior |
| --- | --- |
| 1 | `events`, `storage`, `space`: bounded telemetry; local/mounted/UNC, SMB2/3 and SFTP destinations; available-space checks; combined worker reservations and a default 2 GB buffer. |
| 2 | `state`, `downloader`: atomic JSON checkpoints and process locks; two workers by default; `.part` files; range and identity validation; bounded retries; global pause and sync; ZIP CRC and complete gzip verification before rename. |
| 3 | `auth`, `browserprofile`: Rod/stealth Google browser control; encrypted cookies and Chromium profile; OS keyring credentials and optional TOTP; complete ready-card discovery; exact byte metadata; serialized link refresh and visible challenge fallback. |
| 4 | `app`, `cli`, `gui`, `cmd/takedown`: shared job orchestration; no-argument GUI launch; auth/list/download/status commands; export selection, destination preferences, progress, speed, warnings, pause/resume and stop controls. |

The module requires Go 1.26.10 or newer. Current SSH/network dependencies require
Go 1.26; the selected patch toolchain and updated image dependency address the
issues found during dependency scanning. Native Fyne builds need a C compiler
and the target graphics development libraries. Build instructions are in README.md.

## Verification

| Check | Evidence and limits |
| --- | --- |
| Race tests | Full `go test -p 2 -parallel 2 -race -tags ci ./...` suite passed locally: 230 tests and subtests, zero failures. Fixtures include HTTP interruption, range mismatch, link expiry, corrupt archives, oversized responses, restart, space pauses, sync failure, parser ambiguity, encryption and GUI account changes. |
| Static checks | `make check`: formatting and `go vet -p 2 -tags ci ./...`. |
| macOS arm64 | Default desktop binary and optional CLI-only binary built. A native GUI startup/cancellation smoke check exited successfully; Fyne software rendering was also inspected. |
| Windows/Linux amd64 | CLI-only binaries cross-compiled. Native GUI builds and runtime acceptance on these systems were not performed; target C toolchains are required. |
| Network storage | Local storage and real SSH/SFTP test connections exercised resume and host-key checks. Production SMB, Windows UNC and production SFTP targets were not exercised. |
| Dependency scan | Native and software-driver `govulncheck` scans found zero reachable vulnerable symbols. Two unused dependency advisories remain, described below. |
| Independent review | Repeated source review drove repairs for sync ordering, complete export-card validation, unsupported layouts, database journal preservation and account-bound GUI selections. Final bounded parser review passed 30 focused checks with no additional findings. |

The SMB library has an advisory for its unused `Share.ReadDir` operation,
[GO-2026-5051](https://pkg.go.dev/vuln/GO-2026-5051). The crypto module also contains
the unimported legacy OpenPGP packages covered by
[GO-2026-5932](https://pkg.go.dev/vuln/GO-2026-5932). The application calls neither
affected path; the scan is not a claim that every dependency module is advisory-free.

## Runtime contracts

- Session and profile encryption require a functioning OS keyring. There is no
  plaintext credential fallback. Tests use fake keyring implementations.
- Chromium must already be installed. Authentication never opens the user's
  ordinary browser profile. Its private runtime profile is snapshotted and
  encrypted after confirmed graceful shutdown; forced shutdown retains the prior
  encrypted snapshot and reports an error. Abrupt process termination can leave
  private temporary runtime files and lose changes since the last snapshot.
- Profile snapshots exclude known caches and process locks, retain LevelDB
  journals, and are bounded to 256 MiB and 25,000 entries. Cookie sessions are
  bounded to 8 MiB and parsed pages to 16 MiB.
- Google page structure is not a stable public download API. Discovery requires
  an explicit empty/non-ready list or consistent complete ready cards, with two
  matching snapshots within a bounded wait. Unsupported or partial layouts fail
  explicitly. Collapsed export details may require manual expansion.
- Exact metadata can require a single one-byte range probe per archive. It never
  reads the response body, but it may count against Google's download allowance.
  Repeated standalone discovery commands can repeat probes. Google documents an
  approximate seven-day expiry and five-download allowance in its
  [Takeout documentation](https://support.google.com/accounts/answer/3024190).
- Resume uses byte ranges with strong ETags. Without a strong ETag, the engine
  compares the retained prefix against the same complete response before
  appending, which re-reads that prefix over the network.
- Three retries bound transient transfer recovery. Failed protocol, identity,
  integrity or sync checks retain the partial file and do not publish completion.
- SFTP needs trusted host keys and `fsync@openssh.com` for durable download
  checkpoints. Missing `statvfs` fails closed unless the user explicitly enables
  warning-only unknown-space mode. Canceling a remote operation can close its
  connection; restarting the job reopens it and reconciles actual file length.
- Remote directory mutation by an untrusted concurrent writer is outside the
  storage contract; remote protocols lack the local backend's handle-relative
  no-follow operations. Use a destination controlled by the operator.
- Stop a running job before changing destinations. Existing partial files stay
  at the old destination; they are not moved automatically.

## Account acceptance reserved for the user

1. Create a small ready Google Takeout export, preferably with multiple archive
   parts. Use a supported installed Chromium and an unlocked OS keyring.
2. Run `make build`, then `./bin/takedown auth --email you@example.com` and finish
   Google's own visible sign-in, passkey or phone prompts.
3. Run `./bin/takedown list`. Compare export IDs, creation dates, part counts and
   filenames with Google's page. Exact sizes come from archive response headers.
4. Replace `EXPORT_ID` with a listed ID and run
   `./bin/takedown download --export EXPORT_ID --destination ./acceptance-archives`.
   Verify complete archives appear only after integrity verification.
5. Interrupt once, inspect `./bin/takedown status`, and rerun the same download
   command. Confirm retained bytes resume safely. Keep Google's finite download
   allowance in mind when designing repeated checks.
6. In the GUI, test pause/resume, stop, account changes, and a different destination
   after stopping. Use a disposable quota-limited test volume to exercise a real
   low-space pause without filling an ordinary working volume.
7. Exercise link expiry and archive-specific re-verification. Optional saved
   password/TOTP should attempt supported challenges once; passkey or phone
   challenges should surface the browser and resume after completion.
8. Close normally and reopen. Confirm encrypted session/profile restoration with
   the real browser and OS keyring. Repeat native GUI/runtime acceptance on each
   intended Windows or Linux target and test the intended network destination.

The auth module's detailed integration reference is in internal/auth/HANDOFF.md.
The standalone HTML version of this report contains the same full handoff and
opens locally without a server or network access.
