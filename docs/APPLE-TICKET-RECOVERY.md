# Retained Apple ticket requalification

Local implementation handoff for `fix/apple-ticket-recovery`, based on
`edd4da3edfc4f0788b0fed4405ea27a4097e0a03`. Root owns independent review,
signed/DCO delivery, hosted qualification, release configuration and publication.
This implementation used synthetic providers only. It did not query or interfere
with the active hosted retry, access real credentials, sign, submit, trigger a
workflow, commit, push or change the candidate.

## Evidence and preserved uncertainty

The original private candidate is draft `424598089`, source
`5b4bf4dea8f93a3a3d23467f990bb3310a3fea9b`, version `0.0.1`, build run
`40372460129`, attempt `1`. Its native inputs and original uploaded artifacts
remain immutable. The scheduled run `40407849027`, attempt `1`, demonstrated
automatic continuation of that original candidate and its known UUIDs. Attempt
`2` failed Intel signing job `121896466520` while validating an already verified
app. The supplied task evidence says the user started attempt `3`; its current
state was deliberately not queried by this delegate.

The retained Intel terminal app receipt has these exact identities, independently
read and hashed locally during this implementation:

| Field | Value |
| --- | --- |
| Receipt SHA-256 | `d808e6b66d8ad89a7eb690a21b33b7428c0bacb8401e22447c1f35ce2c9927ae` |
| Receipt asset ID | `679738980` (supplied root evidence) |
| Latest checkpoint in supplied evidence | `checkpoint-amd64-000024-3bf18e3eaff57d613789562fc114a4aaedf3e6938aeffe7091b526d16d49bc3d.json` |
| State/schema | `ticket-blocked` / `2` |
| Submission UUID | `946c8a8c-7c00-447d-8404-67952086f794` |
| Submitted SHA-256 | `1f11f175292cfeccf8121b000e594529498996f06d37ebb14f78925564c2469d` |
| Verified ZIP SHA-256 | `f2c384d41ff9398f763fa89876e4268927e04afa950945fdc889df60ed41483e` |
| Failure | `stapler-validate`, exit `68`, `ticket_attempts: 1` |
| Bound origin | Original product source/version/run/attempt above, `amd64`, stage `app` |

Root's supplied evidence reports both app UUIDs and both outer DMG UUIDs Accepted,
with no publication. Root also reports that verified ZIP asset `679571511`, with
the hash above, passed real signature verification, ticket validation and
Gatekeeper assessment after extraction on local ARM macOS 27.0.1 / Xcode 27.1.
That is evidence that the retained bytes/ticket worked locally; it does not prove
the historical hosted failure's cause or qualify macOS 15 Intel.

The failed runner's private stapler output was destroyed. Exit `68` alone does
not identify CloudKit or any other cause. Its cause remains **unproved**.

## Reproduced seam

Before repair, the synthetic real-signer fixture first produced a verified app,
then supplied a diagnosed transient validation failure while restoring its exact
verified ZIP. The signer exited `1` and wrote `ticket-blocked`, retaining its
nonempty verified hash and `ticket_attempts: 1`. A regression expecting exit `75`
failed before implementation. The existing guard excluded `STATE=verified` from
ticket continuation, and the next invocation refused every blocked receipt.

The repaired regression exits `75`, persists `validation-pending`, and later
revalidates the same ZIP without another signature, key import, staple, archive
creation or submission. Existing rejection, diagnostic mismatch, Gatekeeper,
cleanup and transactional packaging failure tests remain present.

## Receipt safety contract

`scripts/macos-ticket-policy.sh` emits only `cloudkit-transport` or `unclassified`
from failed stapler output. Automatic continuation requires exit `68` plus a
complete allowlisted diagnostic line: `CloudKit network connection lost`,
`CloudKit connection reset`, `CloudKit connection timed out`, or
`CloudKit service unavailable`. Recognized rejection, invalidity, signature,
mismatch and missing-ticket diagnostics veto that classification. Generic
connection failures, an ambiguous CloudKit query, arbitrary exit `68`, and
unrecognized output stay terminal. Raw output remains private and is destroyed;
only the safe classification, phase and exit enter receipts/logs. These markers
define a conservative synthetic-tested policy; they do not classify the lost
historical output or assert what the next Apple runner will emit.

Schema `2` adds optional, backward-compatible fields:

| Field | Contract |
| --- | --- |
| `failure_class` | Empty after complete success; otherwise an allowlisted constant. |
| `validation_attempts` | Persistent integer `0..3`, separate from unchanged staple `ticket_attempts`; never reset. |
| `requalification_parent_sha256` | Empty normally; exact immutable terminal receipt SHA after explicit requalification, retained in every subsequent receipt. |
| `state: validation-pending` | Original verified archive/hash remains bound, but revalidation is incomplete and publication is blocked. |

Each ordinary recheck of a verified archive enters `validation-pending` and
checkpoints before provider operations. A successful ordinary recheck consumes
no recovery budget. Its first failed validation consumes one attempt; later
continuations reserve an attempt and checkpoint before any provider operation.
Three attempts is the persistent ceiling. A failure on the third attempt is
terminal; an interrupted third reserved attempt cannot be retried. A successful
third attempt may finish verification, but it does not restore the budget.
An interrupted initial ordinary recheck reserves an attempt when resumed.

Legacy verified receipts without the new counter start at zero because no
validation recovery failure was recorded. Legacy blocked verified-archive
receipts migrate from their retained `ticket_attempts`; the known one starts at
one and its explicit requalification reserves attempt two. This is conservative
when a legacy receipt consumed multiple staple attempts. Existing staple budgets
remain unchanged; an exhausted staple or validation budget cannot be pinned open.

Every recovery validates original source, version, original run and attempt,
architecture, stage, filename, UUID, team, identity, submitted archive and verified
archive hashes. It restores the **verified** archive whenever its bound verified
hash exists. It never signs, submits, staples, repacks or replaces those bytes.
The verified archive digest is checked again after trust assessment.

Full matching fresh Apple `Accepted` info and diagnostic, signature identity,
secure timestamp/hardened runtime, ticket validation, a second signature check,
Gatekeeper and private credential-scratch cleanup must finish before `verified`.
Failure phase/class fields clear only at that point. Non-validation failures on
retained verified bytes become terminal, including signature, Gatekeeper,
malformed/mismatched diagnostics, changed bytes and cleanup failures.

## Explicit terminal-receipt pin

The optional **nonsecret** signer input is
`NOTARY_REQUALIFY_RECEIPT_SHA256`. Protected-main signing forwards it from
repository/environment variable `APPLE_REQUALIFY_RECEIPT_SHA256`. It adds no
workflow trigger, provider action, secret, architecture or runner change.

This input is an approval of one exact blocked receipt, not an exit-code bypass.
An eligible receipt must be `ticket-blocked`, with `stapler-validate` / `68`,
unclassified (including legacy missing classification), no earlier pinned parent,
retained verified bytes and unexhausted counters. Its retained Apple diagnostic
must already match UUID, original upload filename/hash and `Accepted`; all fresh
provider checks remain mandatory afterward. The signer checks all original
bindings/digests before changing any receipt or making any provider operation.
Wrong/malformed pins, another candidate or stage, missing/changed bytes,
non-validation failure, rejection, diagnostic mismatch and exhausted budgets
leave the terminal receipt unchanged and make no provider query.

Before entering pending revalidation, the signer copies the exact terminal
receipt to `blocked-receipt-<SHA256>.json` in that same stage, then checkpoints it
alongside the unchanged current receipt. Only then does it checkpoint a pending
receipt carrying `requalification_parent_sha256` and the reserved counter.
`apple-release-state.sh` allows only those exact stage-relative parent filenames,
and its existing content-addressed append-only blobs/checkpoints preserve the
original blob/checkpoint. Parent writes reconcile byte-identically after an
interruption; they cannot replace different bytes.

Every later signer load, readiness check and publication validates the parent
digest, terminal failure class, budget floor, original candidate/stage/UUID and
archive identities. A missing, changed or unrelated parent blocks progress.
A later terminal failure cannot be opened by pinning it again: the nonempty
parent and persistent counter prevent repeated approvals from resetting a budget.
The original historical failure is never relabeled as diagnosed transient.

## Consumer gates

`scripts/apple-stages-verified.sh` is a common state/audit gate for both app and
DMG. It supplements existing identity, checksum and native trust checks.

- Continuation cannot reuse retained finals/acceptance while either stage is
  pending or blocked; it must run the signer recovery path.
- Fresh native acceptance makes no new receipt while either stage is pending or
  blocked, even when stale final files/acceptance remain.
- Workflow readiness requires both stage receipts verified before `ready=true`.
- Publication retains its original exact-identity `verified` requirement and now
  also validates new counters, cleared failure fields and any pinned parent.

Fresh copied-app and native acceptance on both existing macOS 15 architectures
remain mandatory. No quarantine removal, keychain permission change, rejection
bypass, submission replay or trust-assessment weakening was introduced.

## Validation and root integration

The focused synthetic suite can be run without credentials:

```sh
SIGN_MACOS_TEST_GROUP=ticket-recovery bash scripts/test-sign-macos.sh
```

The full original signer/packager suite also runs these regressions by default.
Additional affected suites are `test-apple-continuation.sh`,
`test-apple-release-state.sh`, `test-publish-release.sh`, and
`test-verify-macos-package.sh`. Static checks use ShellCheck, Bash syntax,
Actionlint and `git diff --check`. The accompanying local report records exact
results and retained red/green witnesses. This is shell/YAML/documentation work;
no unrelated Go/native build or TypeScript/Python check is required.

The finite queued batch (local job `2171`) exited `0`:

| Check | Local result |
| --- | --- |
| Focused `SIGN_MACOS_TEST_GROUP=ticket-recovery` signer suite | PASS on final files, including 24 pre-provider refusals, interrupted parent/pending checkpoints and both stage budgets |
| `bash scripts/test-sign-macos.sh` | PASS, full original signer/packager suite plus new recovery regressions |
| `bash scripts/test-apple-continuation.sh` | PASS, including stale acceptance/finals with pending or blocked app/DMG |
| `bash scripts/test-apple-release-state.sh` | PASS, including exact parent receipt custody round trip |
| `bash scripts/test-publish-release.sh` | PASS, including pending/blocked receipts and invalid recovery metadata |
| `bash scripts/test-verify-macos-package.sh` | PASS, existing copied-app/trust/cleanup failure contracts retained |
| `shellcheck -x scripts/*.sh build/linux/*.sh` | PASS |
| Bash syntax, all existing shell scripts above | PASS |
| `actionlint .github/workflows/release-macos.yml` | PASS |
| `git diff --check` | PASS |

Root's integration review also found that legacy `verified` receipts without the
new counter could retain a failure marker and pass the common state gate. The
gate now requires every present failure phase, exit and class to be empty for
legacy and extended receipts. A focused before/after probe reproduces and closes
that gap; publisher regressions cover both stages and architectures. Root's
combined seven-suite batch and independent combined source review are recorded
separately from the implementation delegate's results.

Changed paths:

- Signer/policy: `scripts/sign-macos.sh`, `scripts/macos-ticket-policy.sh`.
- State/readiness: `scripts/apple-release-state.sh`,
  `scripts/apple-stages-verified.sh`, `scripts/continue-apple-release.sh`,
  `scripts/fresh-apple-acceptance.sh`, `scripts/publish-release.sh`.
- Protected-main forwarding/readiness: `.github/workflows/release-macos.yml`.
- Regression fixtures: `scripts/test-sign-macos.sh`,
  `scripts/test-apple-continuation.sh`, `scripts/test-apple-release-state.sh`,
  `scripts/test-publish-release.sh`.
- Handoff: `docs/APPLE-TICKET-RECOVERY.md`. The companion standalone HTML is
  ignored, shared through chat, and excluded from Git.

Root's concrete remaining steps, under the existing authorized endpoint:

1. Review the same-candidate policy and all changed paths, run the combined
   independent review, then own signed/DCO commit, PR, exact-head CI and merge.
2. Reconcile the newest original candidate checkpoint after the active retry.
   Preserve all existing UUIDs and uploaded/final bytes. If its terminal receipt
   changed, the historical pin below must fail; do not substitute a past receipt
   or reset counters. Review the new evidence before selecting another exact pin.
3. Only if the newest eligible receipt remains byte-identical to `d808...927ae`,
   set nonsecret `APPLE_REQUALIFY_RECEIPT_SHA256` to
   `d808e6b66d8ad89a7eb690a21b33b7428c0bacb8401e22447c1f35ce2c9927ae`
   for protected-main signing. Do not use the submitted ZIP hash, verified ZIP
   hash, checkpoint hash or asset ID as the pin. No variable was set locally.
4. Let the original-candidate continuation run only after merged automation is
   qualified by its required exact-source CI. Confirm the original receipt blob
   and checkpoint still exist, the new checkpoint retains the exact terminal
   parent and counter, and no UUID/submitted/final archive changed.
5. Require actual fresh Accepted diagnostics, signatures, ticket validation,
   Gatekeeper, credential cleanup, copied-app/native acceptance on **both** runners
   and publisher checks. Remove the one-receipt approval variable after use.
   Public release and published-byte readback remain root's separate work.

No hosted, Apple, macOS 15 qualification or publication success is claimed by
this local handoff. The active retry's result and independent combined review
remain outside this delegate's evidence.
