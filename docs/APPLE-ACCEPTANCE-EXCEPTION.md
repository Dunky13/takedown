# One-time ARM final acceptance exception

This implements the user's approved extra credential-free qualification for
candidate `425202908`, version `0.0.1`, product source
`ef82381d509aac8b5b02da0f784a8176b1c5b9c0`, original build `40470756538`, attempt
`1`. Implementation is local only. Root owns review, workflow integration,
delivery, exact activation, hosted qualification and publication.

## Authority and immutable custody

The nonsecret approval input is `NOTARY_ACCEPTANCE_REQUALIFY_SHA256`. The only
eligible value is:

```text
27b7a083cc181ed6def7479ec99b25479d7b11d9c3de36bbbd9facfbacfc60fe
```

This seals the exact original ARM `acceptance-state.json`: schema 1, blocked,
spent count 3, `cloudkit-transport`, with no previous exception parent. The
approval also requires `NOTARY_RELEASE_ID=425202908`, repository
`Dunky13/takedown`, ARM architecture and the original source/version/run identity.
No second digest is eligible. Intel and future candidates keep ceiling 3.
The legacy signer `NOTARY_REQUALIFY_RECEIPT_SHA256` input is unchanged.

The original state is retained byte-for-byte as
`blocked-acceptance-27b7a083cc181ed6def7479ec99b25479d7b11d9c3de36bbbd9facfbacfc60fe.json`.
Only this extra filename is allowed through existing content-addressed checkpoint
custody. The checkpoint containing that parent retains the unchanged blocked3
state and must be acknowledged before reserved4 is written/checkpointed. Failure
of either acknowledgement stops before native trust checks. Unknown upload
outcomes retain the existing custody collision/fork protections.

Before handoff or reservation, the common policy audits the exact parent, input
and stage receipt hashes, original identities, schema, both verified signer
stages, matching Accepted logs, stage UUID/filename/submitted-hash bindings,
retained stapled stage files, exact final ZIP/DMG hashes and checksum manifest.
Missing custody, symlinks, changed bytes, wrong pins and foreign evidence refuse
before native qualification. Nothing signs, submits, staples or repacks these
retained archives.

## Producer and consumer

The signing continuation handles the eligible blocked3 receipt by retaining the
parent and checkpointing the exact old state/finals. It removes stale native and
fresh success receipts, performs no native qualification, spends no attempt and
issues no acceptance receipt. It imports no credentials. A completed verified4
checkpoint can pass an audited signer continuation without a fifth native check.

The fresh wrapper tolerates absent signer acceptance only through the fully
validated approved blocked-parent path. It restores Intel custody using the same
candidate ID/body/source pins and original run/attempt. Intel's original app UUID
`88951ab7-f126-485b-945c-a11fdbdee902` and DMG UUID
`70d28837-3ac6-4de3-86b8-a07c1813af89` must be preserved. Known pending signer or
final-acceptance work returns pending with ARM still unspent3. Rejected, unknown,
foreign or malformed counterpart evidence blocks; ready evidence must pass the
complete byte/identity/Accepted-log/native-acceptance audit. ARM's original app
UUID `4549ac29-2d27-4cae-b118-d0ce238b1e0c` and DMG UUID
`31f135e5-19a0-4dff-b393-f1cd16b72878` remain sealed by its stage receipt hashes.

Only the fresh wrapper marks the verifier as a fresh consumer. The verifier
independently repeats eligibility and counterpart audits, requires a checkpoint
hook and refuses nonempty signing certificate/identity/API credential inputs.
It checkpoints blocked3 plus parent, then reserves pending4 plus parent before
the first native check. It never resets the counter or relabels the old failure.

Every native checksum, Developer ID/team/timestamp, app hardened-runtime, ticket,
Gatekeeper, metadata, slice and executable-version check remains mandatory. The
DMG is mounted read-only; its copied app must match the ZIP executable. Owned
detach, native scratch and fresh counterpart-directory cleanup must finish before success. No quarantine or trust
authority is changed.

Success writes verified4 only after cleanup, with an acceptance receipt carrying
`acceptance_attempts:4` and the immutable `acceptance_exception_parent_sha256`.
The same durable receipt records the actual credential-free verification run
and attempt. A later peer or publication continuation may restore this exact
completed fresh proof after all parent, stage and final-byte audits; it preserves
the original verification identity and runs no fifth native qualification.
The checkpoint remains provisional until the exact original fresh job is
authoritatively completed/success and its uploaded proof matches byte for byte.
Cancellation, failure or uncertain completion blocks replay even after a late
immutable verified checkpoint commit. See
[fresh completion authority](APPLE-FRESH-COMPLETION.md).
Missing or malformed fresh provenance cannot use that recovery path.
The fresh receipt inherits these fields. Readiness audits the parent, bindings,
retained byte custody and native receipt; publication additionally requires the
fresh receipt's same attempt/parent binding. Legacy receipts without an exception
retain all prior checks. Stale signer/fresh receipts cannot satisfy the new gate.
Exception markers in either native or fresh receipts require complete exception
state and immutable-parent custody. Omitting both custody files cannot downgrade
those receipts to the legacy acceptance path. A guard cancellation withholds
native success and checkpoint writes even if an interrupted native tool returns
success late; the reserved fourth attempt stays spent.

Any attempt-4 failure is terminal. An interrupted or unacknowledged reserved4
may retain pending4, but is spent and cannot be retried, repinned or promoted.
Later wrappers refuse it before native checks. State/receipt write and rename
failures explicitly return errors; cleanup cannot acknowledge a failed state
transition as success.

## Required root integration and activation sequence

1. Review the combined source, run combined static checks and existing CI, and
   deliver the authorized automation normally. This document authorizes no
   hosted operation or activation by this implementation worker.
2. Forward repository variable `APPLE_REQUALIFY_ACCEPTANCE_SHA256` into
   `NOTARY_ACCEPTANCE_REQUALIFY_SHA256` on both signing continuation and fresh
   acceptance steps. Retain all current candidate ID/body/source and original
   run/attempt pins; do not forward Apple credentials into fresh acceptance.
3. For the job-owned transport guard, set `NOTARY_ACCEPTANCE_CLEANUP_HOOK` to the
   absolute runtime path of `scripts/macos-ticket-transport.sh` in both native
   verifiers' environments. The verifier invokes `bash "$hook" cleanup` after
   owned native cleanup and before either success receipt or verified checkpoint.
   The hook must return nonzero on cleanup failure. The workflow's unconditional
   cleanup remains required and must be idempotent after a successful hook.
   Deferred fresh receipt promotion remains required until workflow cleanup
   succeeds. Do not activate the exception with cleanup wired only after success.
   The guard exports its actual workflow run/attempt to its child; cleanup uses
   those values even while artifact consumers retain the original build IDs.
4. Include the exact `blocked-acceptance-*.json` parent in safe failure evidence
   uploads. Its durable checkpoint path is already allowlisted by core custody.
   Preserve existing trust, native matrix, workflow/branch and permission gates.
5. Reconcile current custody without changing the candidate, UUIDs or archives.
   Allow Intel's existing signer work to become ready under its original budgets.
   ARM's producer handoff may finish while Intel is pending; the fresh consumer
   cannot spend4 until its bound counterpart checkpoint is ready.
6. Root owns setting the exact approved value, observing the one fresh ARM
   reservation and recording native results. Failure/interruption is another
   terminal hold. Require both architectures' current fresh receipts and every
   existing publication gate before any separately authorized publication.

## Local verification and limits

Owned suites are `test-verify-macos-package.sh`, `test-apple-continuation.sh`,
`test-apple-release-state.sh`, `test-publish-release.sh`,
`test-apple-custody-recovery.sh`, and `test-sign-macos.sh`. Run these sequentially
as a finite shared heavy-queue batch. Run ShellCheck and Bash syntax checks only
for the changed shell paths; root owns combined workflow validation.

`ACCEPTANCE_TEST_GROUP=exception bash scripts/test-verify-macos-package.sh` runs
the real producer/fresh/verifier seam with synthetic native/provider tools. The
private test script copies substitute a synthetic receipt digest and matching
parent filename; production has no test override. Publisher fixtures substitute synthetic frozen
provenance only in a private policy copy. These tests establish state, custody
and refusal contracts, not real Apple, hosted macOS, transport repair, first-open
or publication qualification. Existing failure tests remain present.

Detailed local results and the owned source manifest are in the accompanying
ignored implementation report. No real Apple/provider request, hosted mutation,
commit, push, PR or publication was performed by this worker. Root's independent
review and combined integration remain outstanding.
