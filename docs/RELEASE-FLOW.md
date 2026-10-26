# Release automation

`Release and nightly` freezes one protected-main source and builds a versioned
candidate without Apple credentials. Eligible successful `Native CI` main pushes
publish the next patch stable release automatically. Manual patch, minor, major,
and custom SemVer dispatch remains available. The daily 02:23 UTC schedule
produces a signed macOS nightly if main has changed since the last published
nightly; it catches up a missing stable release first. The twenty-minute schedule
resumes unresolved candidates and catches up eligible main after missed events.
Schedules may be delayed by GitHub.

## Authority and supported targets

The release owner authorized automatic stable/nightly macOS signing, Apple
submissions and publication, including removal of the recurring `release-signing`
reviewer. Keep that environment restricted to branch `main`, with administrator
bypass disabled and its seven Apple inputs isolated there. Do not copy credentials
to repository secrets or expose them to PR/native-build jobs. The separate manual
`Signed packages` workflow keeps its explicit authorization input and does not
publish releases. Windows signing remains outside this automatic release graph.

Protected main requires `ci-required` and `conventional-commits`, both from GitHub
Actions. The four-target native matrix and all existing checks remain required.
Automatic candidates reject PR/fork completions, wrong workflow paths, stale
source events, unprotected main and incomplete/unsuccessful exact-source checks.
Required checks are bound to the suites of successful `ci.yml` and
`conventional-commits.yml` main-push runs, including the triggering run ID.
All action references remain pinned and checkout credentials are not persisted.
GitHub draft discovery requires push-capable access, so the custody reader jobs
receive `contents: write` on protected main. They receive no Apple secrets.
The native build matrix retains `contents: read` and no Apple credentials.

The project adapter is Go 1.26.10, Fyne v2.8.2, Cgo and the untagged native GUI
binary. Linux amd64 and Windows amd64 remain unsigned. macOS arm64 and amd64
are separate apps, not a universal binary. Xcode 26.3 is selected and installed
notarytool flags are inspected before credential use. Icons retain both `noasm`
and `-gcflags=all=-d=fmahash=n`. Headless/test-driver binaries cannot be promoted.

Takedown has no automatic updater, signed update inventory or package-manager
publication contract. This change adds no replacement trust root. Developer ID,
Apple diagnostics, tickets and Gatekeeper remain the macOS distribution authority;
SHA256SUMS identifies the final download bytes.

## Version and channel policy

Tags reachable from frozen main are the version source. Automatic stable releases
use the existing patch policy: no tags means `0.0.1`, and `v1.2.3` means `1.2.4`.
Commit types organize notes; they do not silently allocate breaking/minor versions.
Manual `minor`, `major`, and `custom` remain explicit choices. Custom versions must
be valid SemVer greater than the stable base; the nightly namespace is reserved.
Native metadata uses the numeric core (at most 65535 per component for Windows);
filenames, CLI version and build information retain the full SemVer.

```sh
gh workflow run release.yml --ref main -f bump=minor
gh workflow run release.yml --ref main -f bump=custom -f custom_version=1.0.0-rc.1
```

Run one dispatch for the intended version. An unresolved candidate always resumes
before another candidate is allocated, regardless of new version inputs or new
main commits. Duplicate successful-main completions are no-ops after publication.
The nightly version includes UTC planning time and the frozen source prefix;
nightlies/custom prereleases do not become Latest or advance the stable base.

## Build once and retain custody

The release graph serializes planning through publication across both channels
and never cancels an active release. Native CI builds the versioned candidate once
without signing secrets. Routine source CI may separately validate that commit;
its development-version packages are not the release inputs.

After the native aggregate succeeds, a credential-free custody job creates a
private draft named `apple-candidate-v<VERSION>`. It retains `candidate.json`,
`plan.json`, notes, original `native.json` metadata and the four original Actions
artifact ZIPs. The manifest binds repository, clean source, original run/attempt,
version/channel and metadata hashes. Native artifact IDs and SHA-256 digests
remain available after Actions artifact retention expires. A candidate draft is
never published; its assets contain no credentials or raw client output.
The draft reservation atomically includes the original plan, native inventory,
notes and binding metadata in its body. If any initial asset upload is interrupted,
the next run reconstructs that metadata and appends missing original archives
with digest checks, without rebuilding. A completion marker is uploaded last;
signing cannot start without complete custody. If an original archive was never
retained and has already expired from Actions, the candidate stops explicitly.
Reservation metadata is limited to 60,000 bytes and fails before creation if
that bound is exceeded.

The signer verifies the original outer archive digest, its exact package entry
set, inner checksums, regular paths, bundle identity/version, binary slices,
Go build revision and clean-tree metadata before loading credentials. It does
not build the executable again. Both source and original run identities are
preserved during later continuation runs.

Small sequence-numbered checkpoint manifests reference immutable SHA-256-addressed
blob assets in that draft. Every safe receipt transition is acknowledged remotely
before the next operation. Sequence forks, changed bytes, mismatched identities,
missing checkpoints and corrupted blobs fail closed. The latest checkpoint uses
explicit sequence order, not asset-ID magnitude. Draft storage is durable beyond
Actions retention; administrators must not delete unresolved candidate drafts.

## Apple stages and automatic continuation

Each architecture preserves this topology:

1. Sign the executable and app with Developer ID, secure timestamps and hardened
   runtime. Verify expected team/authority and clean temporary credentials.
2. Persist the descriptive signed app ZIP and upload intent, then submit once.
   Persist its UUID immediately. Require matching Accepted diagnostics.
3. Staple/validate the app and assess Gatekeeper. Retain its exact stapled bytes.
4. Package that accepted app in the drag-and-drop DMG, sign and verify the image,
   persist its upload bytes/intent, submit once and retain the second UUID.
5. Require matching Accepted diagnostics, staple/validate the DMG and assess it.
   Create the final ZIP from the stapled app and checksum final distributables.

Submission basenames include product, full version, architecture, source prefix,
original run/attempt and stage. Submitted hashes differ from stapled download
hashes. The UUID, source, original run, stage, filename, submitted hash, selected
team/authority and state are bound together. Never import an unrelated receipt.

Private temporary files and an isolated temporary keychain keep credentials out
of source, argv and public logs. Restore and verify the original user keychain
search list; do not change default/system keychains. Cleanup failure blocks
release. Safe operation labels and exit codes identify the failing signing stage;
private tool output is not uploaded.

`In Progress` returns a pending result rather than publication. The next scheduled
run restores the original build and signed bytes, then queries/waits for the same
UUID. Accepted stages restore their retained app/DMG without signing or uploading
again. Missing UUIDs reconcile history by descriptive name, then require the
original diagnostic hash/name/UUID binding. Ambiguous or absent history never
authorizes another upload. Rejected, malformed or mismatched diagnostics block
promotion. A client timeout is not an Apple rejection or permission to resubmit.
Terminal rejection and invalid diagnostic identity leave blocking receipts.
Only a diagnosed CloudKit connection failure with stapler exit 68 can continue
automatically; a persistent counter permits at most three ticket attempts per
stage. Other ticket failures stop rather than retry generically.

`APPLE_QUALIFY_CONTINUATION=true` is an optional repository variable for hosted
recovery qualification. A newly captured, durably retained UUID yields before
waiting; continuation must use that UUID. This does not simulate acceptance or
bypass any Apple gate. Clear the variable after recording recovery evidence.

Successful architectures retain their exact final ZIP/DMG bytes across later
continuations. Quarantine is preserved. No `spctl` disable, ad-hoc signature,
quarantine removal or synthetic Apple response is an allowed release repair.

## Final acceptance and immutable publication

Signing runners verify final ZIP/DMG checksums, Developer ID/team, tickets,
Gatekeeper, bundle metadata, exact slices and native CLI version. They mount the
DMG read-only, copy its app to an isolated Applications directory and verify that
installed copy, including ZIP/DMG executable equality and owned-mount cleanup.
Fresh credential-free ARM and Intel runners repeat that acceptance against the
retained final bytes. Both fresh receipts are bound to the original candidate
and the verification run. Failed or incomplete acceptance cannot publish.

The publisher requires both architectures' bound input receipts, verified app/DMG
receipts, Accepted logs, signer acceptance and fresh native acceptance. It selects
exactly seven payloads and one combined checksum file. Linux TAR.GZ, Windows ZIP
and per-user setup EXE remain unsigned; both macOS ZIPs and DMGs are signed and
notarized, including nightlies.

Publication atomically creates an unmoving `v*` tag at the frozen source, stages
assets in a separate version draft and verifies every re-downloaded byte and the
complete remote inventory before making it visible. A retry appends missing
assets only. Existing different bytes, unexpected names or a changed source stop
publication. There is no clobber or forced tag. A fully published identical release
is verified read-only and succeeds idempotently. A lost response after publication
can therefore be reconciled without rewriting the release. A durable completion
marker closes candidate custody only after successful promotion.

## Evidence and operational recovery

The private candidate draft is the recovery entry point. Actions also retains safe
failure receipts for seven days. Keep raw client output, key material and private
scratch out of evidence. Permanent credential, certificate, Apple rejection or
ambiguous-submission failures require diagnosis; scheduled continuation cannot
make those successful. Inspect the safe failure phase and exact retained UUID.
Never delete a reservation or manually submit again to clear an unresolved run.

Local fixtures exercise trust-origin refusal, exact CI gates, candidate/run/hash
binding, state corruption/collision, pending/rejected diagnostics, submit-once
recovery, cleanup, stapling/Gatekeeper and partial/immutable publication. They
prove failure behavior only. Real hosted credentials, Apple Accepted diagnostics,
tickets, installed-app results and public download verification must be reported
separately with source, run, submission and asset evidence.
