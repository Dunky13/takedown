# takedown: documented automation handoff

Updated 26 October 2026. This is the entry point for the OAuth/API redesign,
superseding the earlier browser-automation implementation. The authorized endpoint
is a local implementation and validation. Live Google account testing remains
with the user. No Google approval, public publication, or deployment is implied.

## What is automated

| Workflow | Application behavior | User or operator step |
| --- | --- | --- |
| Google authorization | System browser, desktop OAuth with PKCE and a loopback callback; automatic token refresh using the official OAuth library. | Configure a Desktop app OAuth client, enable the relevant API, complete Google sign-in/MFA and consent. |
| Broad Takeout coverage through Drive | Read-only recursive discovery of archive files in the selected folder, exact sizes/checksums, automatic downloads and a persistent watch loop. | Choose Add to Drive and any recurring schedule in Google's normal Takeout interface; choose the resulting folder in takedown. |
| Supported Data Portability products | Explicit API export creation and job management, durable receipts, bounded polling, signed-link renewal and automatic downloading. | Configure selected resource groups and an eligible, appropriately verified API project/grant. |
| Transfer recovery | Resume retained bytes, validate source identity, retry temporary failures with backoff, refresh valid OAuth tokens, check archive integrity. | Restore revoked permission or resolve a non-temporary quota, permission or integrity failure. |
| Storage safety | Combined reservations, preflight capacity checks, space checks at most every 50 MB per worker, pause/checkpoint and resume after space recovers. | Select a trusted destination; change destination or explicitly select part-by-part mode when a batch does not fit. |

`watch` runs until stopped. An explicit `service install` can register it in the
current user's login session; the GUI exposes the same action. Neither mode opens
sign-in or creates exports implicitly. Google can revoke authorization and change
eligibility or quotas.

## Recovery and job management

The Jobs view and `jobs list` expose local IN_PROGRESS, FAILED, CANCELLED,
COMPLETE and EXPIRED states, archived receipts, retry counts and held actions.
`jobs list --refresh` checks at most 100 due jobs per invocation, retaining the
minimum five-minute status cadence and provider retry delays. Ready jobs remain
available when an unrelated job fails. Job-local inaccessible/invalid metadata
stays held; watch continues checking for other jobs. Authorization, quota,
provider backoff, transport and journal failures retain their separate handling.

Retry and cancellation are explicit remote mutations, each preceded by a durable
intent and attempted once. Retry is restricted to FAILED jobs and capped at three
attempts per lineage, including reserved attempts with uncertain outcomes.
Cancellation requires an in-progress time-based job. Archiving hides eligible
terminal jobs while preserving their grant, cadence and retry receipts; it does
not remove downloaded data. The journal permits 1,000 active jobs, 10,000 retained
jobs, 1,000 expired-intent receipts and 4 MiB total. Retained-history and byte caps
fail closed; archival does not erase safety evidence to evade them.

For a lost create/retry response, `jobs reconcile JOB_ID` requires a receipt or
other trustworthy evidence matching the exact pending request, account, resource
groups and creation time. Google's status API proves accessible job state, not
that association, and it provides no job-enumeration endpoint. Guessed IDs are
unsafe. The original intent supplies retained time and resources; legacy unknown
grant types remain conservatively one-time.

`jobs resolve-expired` is deliberately narrower: only a grant recorded as
TIME_BASED before the request can qualify, after fourteen days plus five minutes
and matching current time-based grant evidence. Slow intent persistence reduces
the allowed send window; it cannot extend this recovery boundary. Unknown or
one-time historical grants remain held without evidence-backed known-ID recovery.
Expired-intent tombstones remain in history.

`jobs reset-authorization` requires explicit consent to revoking Google's grant
and losing old remote archive access. The source independently requires every
known job in the current generation, including archived jobs, to be COMPLETE with
matching immutable part metadata, full offsets and durable SHA-256 completion
receipts. Any unresolved intent or expired unknown-intent tombstone blocks it.
Confirmed success advances the local grant generation while retaining receipts;
authorization must then be explicit again. An uncertain reset remains held
indefinitely: a new login, checkbox, deleted file or HTTP 401 cannot prove success.
This follows [Google's reset conditions](https://developers.google.com/data-portability/reference/rest/v1/authorization/reset),
rather than treating reset as a way to bypass one-time limits.

## Integration and policy boundaries

[Google's desktop OAuth flow](https://developers.google.com/identity/protocols/oauth2/native-app)
keeps passwords and MFA in the system browser. The application no longer uses
Rod, stealth, browser profiles/cookies, DOM scraping, credential entry, or CDP
interception. Tokens are stored only in the OS keyring. A missing/unavailable
keyring fails explicitly; there is no plaintext token fallback.

[Drive's download API](https://developers.google.com/workspace/drive/api/guides/manage-downloads)
provides archive bytes using `files.get` with `alt=media`. A file's revision,
metadata, exact byte size and supplied MD5 checksum bind its identity. Each Drive
archive file is presented separately; the app does not infer that all parts of a
Google web export have arrived merely because one filename is visible.

Automatic discovery of Google-created archives and future files requires the
restricted `drive.readonly` scope. Google grants visibility to all Drive files,
even though takedown's discovery is confined to the chosen subtree. A narrower
per-file grant is not represented as recursive permission to every existing or
future child. Public distribution requires applicable Google review. See
[Drive scopes](https://developers.google.com/workspace/drive/api/guides/api-specific-auth).

The [Data Portability API workflow](https://developers.google.com/data-portability/user-guide/introduction)
creates new API jobs for explicitly selected resource groups. It does not list
existing email-link Takeout exports. Its grant is separate from Drive; extra
Drive, email, or OpenID scopes are not appended. Supported countries, account
eligibility, verification and any restricted-scope security assessment still
apply. See [API overview](https://developers.google.com/data-portability/user-guide/overview)
and [API setup](https://developers.google.com/data-portability/user-guide/setup).

The app observes the documented five-minute job-polling minimum and 24-hour
spacing for time-based export requests. Creation intent is persisted before the
non-idempotent POST; an unknown result blocks another creation attempt rather
than assuming failure. Job IDs and archive metadata are retained, with no tokens
or signed URLs in the job manifest. Link renewal uses the same API job, observes
polling cooldowns, and checks that the archive objects have not changed. Renewed
links are shared across sibling parts; a rejected current signature returns the
remaining cooldown as a temporary failure. API
signed links last six hours and data lasts fourteen days, according to
[Google's API method guide](https://developers.google.com/data-portability/user-guide/methods).

No source can bypass a revoked grant, an export quota, access restrictions,
account eligibility, or a Google verification requirement. Existing Takeout
email-link exports must be downloaded through Google's normal interface or
re-created using a supported delivery/API workflow.

## Operator setup

See [README.md](../README.md) for build requirements and copy-pasteable commands.
For Drive, configure the Desktop app OAuth JSON and selected folder, run `auth`,
then `watch`. The GUI also provides explicit folder navigation and a filterable
resource catalog with exact manual-entry synchronization. `export` opens normal
Takeout setup. For Data Portability, use a
separate OAuth profile/grant with selected `--resources`, run `auth`, then
`export` and `watch`.

The optional email is a login hint/local profile slot, not verified Google
identity. Keep it consistent across invocations. Data Portability job histories
are isolated by the validated OAuth client and local profile slot, while resource
changes within a profile retain its previous job/intention history.

`auth --forget` removes local authorization and supported legacy credentials for
the selected profile. It does not revoke Google's remote grant; do that in
[Google Account permissions](https://myaccount.google.com/permissions). Other
profiles are separate. Downloaded archives remain until you delete them.

Publish the [privacy policy](../PRIVACY.md) at a public operator-controlled URL
and configure its URL in the OAuth consent project before public distribution.
The project must accurately represent its identity, requested access, data use,
retention and deletion. The policy includes the required Google Limited Use
commitment. These requirements come from the
[Google API Services User Data Policy](https://developers.google.com/terms/api-services-user-data-policy)
and [Data Portability policy](https://developers.google.com/data-portability/policy).

Archive outputs are not encrypted by the app. Protect the destination with
appropriate encryption and permissions. SFTP provides encrypted transport;
plain direct SMB is not asserted to provide encryption. Local/OS-mounted
protection depends on the actual deployment. Restricted-scope security
requirements must be satisfied for that deployment before use/distribution;
this source implementation does not certify them.

Data Portability's encryption requirement and the destination limitations are
shown in both normal CLI authorization and the GUI's disclosure/confirmation
before OAuth. Reading or accepting that notice does not establish encryption.
The actual destination must satisfy the policy before the API is used.

## Diagnostics and unattended watching

`doctor` and the GUI Diagnostics view inspect local configuration and saved-token
availability. Google discovery is opt-in (`--online`); connecting to storage and
querying capacity is separately opt-in (`--storage`) and may create the selected
directory. These checks do not establish verified account identity, durable
remote writes or destination encryption. Grant inspection uses `jobs access` or
the GUI's explicit check; a one-time grant result does not prove it is unused.

`service install --dry-run` renders a profile-specific definition without reading
tokens, installing files or registering a task. An explicit install saves absolute
runtime paths and starts a watch in the current user's login session. The native
definitions are a macOS LaunchAgent, Linux systemd user service, and Windows
interactive Scheduled Task with least privilege. No system-wide service,
password-bearing task, automatic sign-in or implicit export creation is used.
The binary must remain at its installed path; remove the watch before moving it.
Saved OAuth refresh authorization and an unlocked user keyring are prerequisites.
Confirmed GUI installation adopts the absolute client, destination, SSH-key and
known-hosts paths before scheduling work; later saves retain them. Matching
loaded macOS definitions are already installed, rather than bootstrapped twice.
An existing run keeps its current preferences; remove and reinstall to restart
with changed choices. Removal checks OS state even if its definition file is
missing, and safely unlinks stale definitions without following their symlinks.
Native registration, restart and teardown still require platform acceptance.

For password-based SMB/SFTP, explicitly save destination authorization using
`storage auth` or Diagnostics > Destination access. It is stored in a separate OS
keyring entry bound to the data directory and exact destination, never in a
service definition, preferences or manifest. `storage forget` removes that entry;
`auth --forget` removes OAuth authorization separately. Remove saved entries and
the startup task before deleting their profile directory. SFTP keys and trusted
host-key configuration are supported instead of a saved password. Background
watch cannot satisfy an interactive key-unlock or Google consent prompt for you.

## Output layout and packaging

New downloads use `exports/<SHA-256 of export identity>/<original filename>`;
different exports can share filenames without overwriting each other. The state
records this relative path separately from immutable provider identity. Existing
flat checkpoints keep their established paths, including partial-file resume.
Conflicting legacy paths fail explicitly rather than merging unrelated bytes.

Native CI and package tooling, user installation, signing prerequisites and
checksum verification are documented in [RELEASING.md](RELEASING.md).
Local package generation, hosted CI,
native platform runtime acceptance and real signing/notarization are separate
evidence boundaries. Packaging never establishes Google project approval.

## Restart integrity and retry behavior

Successful completions store a SHA-256 digest of the compressed archive bytes in
the durable manifest. Restarts compare existing final files and complete partial
files against that digest, so replacing an archive with a different valid ZIP of
the same size fails. Versions 2 and 3 remain readable and upgrade to version
4 on the next write. A full legacy checkpoint without a local digest or upstream
checksum must prove its bytes against the source before completion is trusted.

OAuth and source retry delays survive downloader error handling. A delay longer
than the short batch retry budget is returned immediately to the watch loop;
watch waits at least the requested delay. Authorization/setup failures retain a
safe authorization classification and actionable instructions. Signed request
URLs and raw provider response bodies are not exposed in those errors.

## Local acceptance

The frozen implementation passed local checks on macOS arm64 with Go 1.26.10:

- Race suite: 541 tests/subtests across 13 packages, using
  `go test -mod=readonly -p 2 -parallel 2 -race -tags ci -json ./...`.
- `go vet -mod=readonly -p 2 -tags ci ./...` and `go mod verify` passed.
- Native macOS desktop/CLI and headless Windows amd64/Linux amd64 builds passed.
  Cross-compilation is build evidence, not runtime acceptance on those systems.
- Fake-backed job, reconciliation, grant-history, reset, independent failure,
  duplicate filename, legacy resume, diagnostics and login-definition checks
  passed. No tests registered real startup tasks or wrote real keyring entries.
- GUI race/vet, confirmations, stale settings, cancellation and clipping checks
  passed. Jobs, folder/resource selection, diagnostics, startup, destination
  authorization and consent screenshots were inspected at 960×760 and 800×640.
  The full-catalog consent regression keeps its controls visible while scrolling.
- The native desktop window stayed alive during startup, exited with code 0 on
  SIGTERM, emitted no migration/thread warnings, and wrote no application data.
- CLI startup preview, local diagnostics and an unconfirmed authorization reset
  were exercised without creating a profile or contacting Google.
- Fresh unsigned macOS app/ZIP/DMG packaging passed checksum verification, ZIP
  integrity and inventory checks, plist validation, bundled CLI help and DMG
  verification. The ZIP contains the exact tested native binary. No installation,
  signing or notarization was performed.

The final Go/module tree contains 113 files. Its SHA-256 is
`07737b9e1fa9c61d94b464401db0e28fcfd0a417d03a2731ef3360e4212d9620`,
computed over sorted lines `<file SHA256>  <repository-relative path>\n` for
`cmd/**/*.go`, `internal/**/*.go`, `go.mod` and `go.sum`. The native macOS binary
has SHA-256 `d43cc5fa1dbb2dfc0233666f17ddd364cb145ccefa6ee8753516798008092a58`.
These identify the uncommitted local implementation, not a published release.

The first independent refinement review found four integration defects: GUI path
persistence, fair refresh beyond 100 jobs, repeated startup installation and stale
definition removal. The first release-tooling review found five installer,
publication, artifact-retention and keychain-lifecycle defects. All were repaired
with the original failing evidence preserved.

Fresh independent R2 reviews are **CLEAN within their bounded local scopes**:

- Integration: all four findings resolved; 24 top-level / 44 named tests and
  subtests passed independently with race detection, plus vet. All 118 custody
  files matched before/after, and the independently recomputed 113-file Go/module
  digest matches the final tree above. All original failure evidence is intact.
- Release tooling: all five findings resolved at their tested source, mock and
  offline levels; 33 focused command invocations and 22 additional fixtures passed
  their intended assertions. All 24 reviewed files and original failure evidence
  remained unchanged. The Windows locked-executable test is wired into native CI
  but was not executed here.

The full reports and executable witnesses are private local evidence, shared
separately. These review results do not promote platform/signing/provider holds
into runtime acceptance or Google certification.

Live Google account acceptance, real SMB/SFTP/UNC storage,
destination encryption, Windows/Linux GUI runtime acceptance, Google verification
and public privacy-policy hosting remain unqualified. Native startup registration,
hosted CI, real keyring integration, signing and notarization are also unqualified.
Cross-provider review was
skipped; ordinary OpenAI independent review is separate from it. No commit,
publication, release or deployment was requested or performed.

The earlier [terms research](TOS-RESEARCH.md) and [implementation handoff](HANDOFF.md)
are retained as historical evidence. Their removed-source links are not current
implementation references.
