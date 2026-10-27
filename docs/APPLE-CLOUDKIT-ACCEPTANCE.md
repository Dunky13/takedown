# Bounded final Apple acceptance recovery

Implementation handoff for `fix/apple-cloudkit-acceptance`, based on
`e1d634bdd2d02a34ba5f1bb7ad7bdf86b6620402`. Root owns independent review,
commits, delivery, hosted qualification, Apple operations and publication.
This implementation performed local source edits and synthetic checks only.

## Diagnosis and evidence

The supplied hosted log for scheduled run `40431870962`, attempt `1`, ARM job
`121951526426`, lines 379–388, shows the retained final ZIP and DMG passing their
checksums, followed by successful ZIP-app codesign verification. Native stapler
then reports `NSURLErrorDomain Code=-1005`, “The network connection was lost,”
with `NSErrorFailingURLStringKey` naming the exact Apple ticket-delivery lookup
endpoint, followed by `CloudKit's response is inconsistent with expections:
(null)` and validation error `68`.

The fixture `scripts/testdata/stapler-cloudkit-lost.txt` preserves that multiline
shape and message spelling. Paths, object pointers, task identities and network
address data are synthetic. Original captured evidence was not edited.

Ranked hypotheses were:

1. Retained-final reuse bypasses the signer's bounded transport recovery. Its
   call chain would terminate directly at final stapler validation.
2. The shared classifier misses native output because it recognizes four short
   synthetic CloudKit lines. The native fixture would classify as unclassified.
3. A bad ticket causes the failure. The new log establishes a transport failure;
   it does not establish ticket validity or diagnose any earlier incident.

Both first predictions reproduced before repair. The actual continuation wrapper
invoking the real final verifier exited `68` instead of checkpointing pending.
The real signer restoring its previously verified archive classified the same
native fixture as unclassified, wrote a terminal receipt and exited `1` instead
of `75`. Retained local red witnesses document both failures.

The earlier Intel exit `68` remains **unproved**: its private native output was
lost. This repair never relabels that failure or changes its exact receipt pin.

The frozen product remains source `5b4bf4dea8f93a3a3d23467f990bb3310a3fea9b`,
version `0.0.1`, build run `40372460129`, attempt `1`, draft `424598089`.
Root supplied that all four original UUIDs were Accepted. This implementation
made no live query and claims no updated run, Apple or publication state.

## Native diagnostic policy

`scripts/macos-ticket-policy.sh` retains its previous four exact short-line
classifications and additionally recognizes the observed native combination:

- Stapler exit must be `68`.
- The NSError must report `NSURLErrorDomain Code=-1005` and the observed lost
  connection message.
- `NSErrorFailingURLStringKey` must equal
  `https://api.apple-cloudkit.com/database/1/com.apple.gk.ticket-delivery/production/public/records/lookup`,
  followed by the native field delimiter. Another host, endpoint or suffix fails.
- The observed complete CloudKit response line must also be present.
- Rejection, invalidity, signature/mismatch, absent-ticket, non-notarized,
  revocation and denial diagnostics veto transport classification.

The emitted classification is only `cloudkit-transport` or `unclassified`.
Raw stapler output remains private scratch and is never checkpointed or printed.
Other NSError codes and arbitrary exit `68` remain terminal. This evidence-based
allowlist deliberately makes no claim about other native service failures.
The signer consumes this shared classifier with its existing receipt semantics,
unchanged ticket/validation ceilings and unchanged blocked-parent approval.

## Final acceptance receipt and bindings

A new allowlisted `acceptance-state.json` records schema `1`, state
`pending` / `blocked` / `verified`, persistent `attempts` in `0..3`, and a safe
`failure_class`. It is separate from both original stage receipts.

Its immutable binding includes:

| Field | What it binds |
| --- | --- |
| `source_commit`, `version`, `architecture` | Frozen product identity |
| `ci_run_id`, `ci_run_attempt` | Original build identity, including on fresh runners |
| `team_id` | Expected Apple team |
| `zip_sha256`, `dmg_sha256` | Exact final download archives |
| `input_sha256` | Entire existing input receipt, including original artifact ID/digest and bundle metadata |
| `app_receipt_sha256`, `dmg_receipt_sha256` | Exact verified stage receipts, including UUIDs, archive identities, counters and any pinned parent |

Restored state must match all bindings before any native operation. A malformed,
symlinked, changed, terminal or exhausted state cannot be used to reset recovery.
No signer receipt, blocked parent, UUID, submitted archive or final archive is
modified by final acceptance recovery. There is no new approval variable.

The continuation retains final ZIP/DMG/checksums before the first verifier call.
The pending checkpoint therefore includes the archives even when their first
final acceptance fails. Later runs reuse those exact files when acceptance is
pending, without requiring a stale success receipt. Missing bound final custody
or nonverified signer stages cause refusal before packaging can run.

## Persistent budget and cleanup

The final acceptance budget is per architecture and exact final-archive pair,
shared by the signing-runner final verifier and credential-free fresh verifier.
It does not replace or extend either existing signer budget.

| Transition | Counter behavior |
| --- | --- |
| Ordinary initial/repeated check | Persist pending before native operations; keep counter unchanged |
| First diagnosed failure during an ordinary check | Spend one attempt; exit `75` only while below `3` |
| Resume a pending check | Reserve the next attempt and checkpoint before native operations |
| Successful check | Preserve spent/reserved count; clear failure class; mark verified only after complete trust checks and cleanup |
| Failure on attempt `3` | Terminal blocked; no automatic retry |
| Interrupted reserved attempt `3` | Remains pending locally/durably; next invocation blocks before native operations |
| Unknown/ticket/signature/Gatekeeper/metadata/cleanup failure | Terminal; cannot become a transport retry |

An interrupted ordinary check starts recovery on the next invocation. Handled
interruptions preserve the existing pending reservation; abrupt runner loss uses
the last acknowledged checkpoint. Failed checkpoint acknowledgement stops the
job before further native operations and cannot issue acceptance. A successful
ordinary recheck never replenishes the counter, including after attempt `3`.
A third reserved check may succeed, but a later failure cannot reopen its budget.

All retries repeat final checksums, app/DMG signatures and teams, ticket validation
at all three sites, Gatekeeper, metadata, architecture, executable version,
read-only DMG mount, copied-app checks and executable equality. Final archive
hashes are checked again. Acceptance is written only after owned mount detach
and private scratch removal succeed. Cleanup failures override pending success.
Quarantine and Gatekeeper policy are preserved. Final recovery performs no
signing, submission, stapling, repacking or credential import.

Root's integration review retains mount ownership until a successful detach.
If the final detach fails, EXIT cleanup attempts that owned detach again and
the original failure still blocks acceptance. If the mount remains attached,
cleanup leaves its private directory in place rather than recursing into it.
Focused fixtures cover both a recovered detach and persistent detach failure.

## Producer and consumer integration

`apple-release-state.sh` now permits only the additional exact filename
`acceptance-state.json` through its existing content-addressed blob/checkpoint
custody. The native verifier checkpoints pending, reservations, diagnosed failure,
terminal failure and complete success. Both wrappers handle exit `75` as pending,
without creating fresh acceptance.

`fresh-apple-acceptance.sh` uses the original build run/attempt when checkpointing,
while its fresh receipt retains the current verification run/attempt. Both native
architecture matrix entries use this same path without Apple credentials.

`apple-acceptance-ready.sh` and publication call the common acceptance gate. When
the new state exists it must be verified, structurally valid, match the retained
input/stage receipt hashes, and match the acceptance receipt's candidate identity
and final hashes. Pending/blocked state defeats stale acceptance and final files.
Original signer gates and all publisher artifact, Accepted-log and original-run
checks remain mandatory. Legacy evidence without the new state keeps its prior
checks; every actual new native check creates state before provider operations.
The workflow also preserves the new safe receipt in failure artifacts.

## Local validation

The following use synthetic tools only, with no real credentials, signing,
Apple requests or mounts:

```sh
ACCEPTANCE_TEST_GROUP=cloudkit bash scripts/test-verify-macos-package.sh
SIGN_MACOS_TEST_GROUP=ticket-recovery bash scripts/test-sign-macos.sh
bash scripts/test-sign-macos.sh
bash scripts/test-verify-macos-package.sh
bash scripts/test-apple-continuation.sh
bash scripts/test-apple-release-state.sh
bash scripts/test-publish-release.sh
bash scripts/test-apple-custody-recovery.sh
shellcheck -x scripts/*.sh build/linux/*.sh
for script in scripts/*.sh build/linux/*.sh; do bash -n "$script"; done
actionlint .github/workflows/release-macos.yml
git diff --check
```

Run full suites as a finite batch through the shared heavy-command queue.
The targeted native final-reuse regression was red before repair and green after.
The signer red witness also used the real restoration seam and native fixture.
The final validation results and objections are recorded below before handoff.

## Changed paths

- `.github/workflows/release-macos.yml`
- `scripts/apple-acceptance-ready.sh`
- `scripts/apple-release-state.sh`
- `scripts/continue-apple-release.sh`
- `scripts/fresh-apple-acceptance.sh`
- `scripts/macos-ticket-policy.sh`
- `scripts/publish-release.sh`
- `scripts/verify-macos-package.sh`
- `scripts/test-apple-release-state.sh`
- `scripts/test-publish-release.sh`
- `scripts/test-sign-macos.sh`
- `scripts/test-verify-macos-package.sh`
- `scripts/testdata/stapler-cloudkit-lost.txt`
- `docs/APPLE-CLOUDKIT-ACCEPTANCE.md`

The full standalone HTML is an ignored local chat artifact and is not for Git.
Root's existing reports, live recovery, followup and rendering files are outside
this implementation's edits.

## Final results and review objections

Finite queued validation batch `2184` exited `0` on the final implementation:

| Check | Result |
| --- | --- |
| Full real-signer/packager synthetic suite | PASS, including native-output restoration recovery, existing pinned-parent refusals and immutable budgets |
| Full final-verifier synthetic suite | PASS, including both architectures at ZIP/DMG/copied-app sites, budget/interruption/cleanup/refusal cases and all original acceptance failures |
| Continuation wrapper suite | PASS, including original pending signer and stale-final refusal contracts |
| Durable state suite | PASS, including additional state-file custody round trip and original blocked-parent custody |
| Publication suite | PASS, including both architectures with pending/blocked/malformed/unbound acceptance state and stale success receipts |
| Original candidate custody recovery suite | PASS, original reservation/partial-upload recovery preserved |
| ShellCheck, all existing shell scripts in scope | PASS |
| Bash syntax, all existing shell scripts in scope | PASS |
| Actionlint, release macOS workflow | PASS |
| `git diff --check` | PASS |

Red commands/results before repair:

- `ACCEPTANCE_TEST_GROUP=cloudkit bash scripts/test-verify-macos-package.sh`:
  exit `1`; wrapper exposed native stapler exit `68` instead of pending recovery.
- `SIGN_MACOS_TEST_GROUP=verified-recovery bash scripts/test-sign-macos.sh`:
  queued job `2180`, exit `1`; real restoration seam classified the native
  diagnostic as unclassified and stopped permanently.

The final complete signer and verifier suites include these same regressions and
passed. Focused final-verifier job `2183` also passed. An intermediate regression
run found a test harness assumption that only one mount pair would occur; the
new repeated-recovery tests correctly make multiple balanced attach/detach pairs.
The assertion now validates every pair. No original failure test was removed.

No known source defect remains from this scoped implementation review. Root's
independent review remains outstanding. Material limits and tradeoffs:

1. Native classification is deliberately limited to the observed combination;
   unfamiliar failures remain terminal and require new evidence.
2. Fresh runners now durably checkpoint safe acceptance state using their
   existing GitHub draft-write permission. Extra custody reads/writes are part of
   the fail-closed contract; acknowledgement failure blocks progress.
3. A fully exhausted or unclassified final-acceptance failure has no automatic
   approval/reset mechanism. Existing signer pins cannot reopen it.
4. Synthetic success establishes no current Apple service behavior, native
   macOS 15 runner acceptance, hosted CI qualification or release publication.
   The historical Intel failure remains undiagnosed.

Root should review the final source and evidence, then own normal signed/DCO
delivery and exact-head CI. Continue only the retained original candidate under
its existing authorization. Reconcile live custody before any retry, preserve the
original four UUIDs and archives, and require real complete native acceptance on
both existing runners before publication. This implementation did not inspect,
restart, cancel or otherwise interfere with root's active failed-jobs retry.
