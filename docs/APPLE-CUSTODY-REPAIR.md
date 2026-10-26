# Apple custody identity repair

Local implementation handoff. Root owns review, delivery, hosted qualification,
Apple operations and publication. This work made no live mutations, commits,
pushes, PRs, secret requests, Apple calls or release dispatches.

## Preserved candidate and observed failure

- Product source: `5b4bf4dea8f93a3a3d23467f990bb3310a3fea9b`.
- Original release run: `40372460129`, attempt `1`; its four versioned native jobs
  and aggregate succeeded in the saved evidence.
- Custody job `121757462337` created draft `424598089`,
  `apple-candidate-v0.0.1`, then exited `5` with
  `jq: candidate draft missing or ambiguous`. The later saved inventory is empty.
- The failed runner's actual list-response body was not retained. Staleness is a
  reproduced failure mode, not a proved explanation for GitHub's original response.
- The frozen reservation reconstructs byte-identically, including the original
  candidate receipt and all three sealed manifest hashes. Reservation SHA-256:
  `fdf2f7aecf6d47ba2eafdeaf81e1fbb1ae4607be0e97a1ea85273769e00f0716`.

## Implementation

`prepare-apple-release.sh` retains and validates the creation response, then
exports the numeric release ID, original target SHA and sealed-body digest.
Existing reservations must match exactly one full tag and the exact original body.
`apple-release-state.sh` checks exact GET-by-ID identity before custody reads and
writes. Direct-ID REST octet-stream uploads validate receipt name, positive asset
ID, size and uploaded state, then download that ID and compare bytes. A private
inventory bound to the reservation is extended by verified upload receipts across
initial-custody subprocesses. No immediate release/asset relist is required after a
known write. Collisions refuse; uncertain writes are never blindly retried.
A later invocation reconciles retained bytes before appending anything missing.

All actual binary API downloads, including original Actions ZIPs, state assets,
signing input and publication verification, use the installed/supported
`gh api --allow-escape-sequences` flag. Installed CLI verified locally: `2.103.0`.

The automatic release workflow executes planning, custody, signing, fresh
acceptance, readiness and publishing from immutable workflow `github.sha` with full
history. Native builds still use the original planned product SHA. The central
`release-source.sh` requires an explicit automation/workflow pin match, clean
checkout, original-product ancestry, protected-main membership, exact successful
main-push `ci.yml` and `conventional-commits.yml` runs and their required-check
suites. Requalification also applies on resume. Legacy calls without an automation
pin still require the exact product checkout. A separate append-only
`automation-<SHA>.json` binds automation to the original product; original receipt
bytes, versions, run/attempts, artifact IDs and digests are preserved.

Custody ID/body pins propagate through reusable macOS inputs and publisher inputs.
Only credential-free qualification adds `actions: read` and `checks: read` where
needed. Signing secrets, cleanup, Apple Accepted diagnostics, staple/Gatekeeper,
fresh acceptance and publisher gates remain enforced. Automatic `workflow_run`,
daily and twenty-minute triggers and original-candidate priority remain intact.

## Regression evidence

Tests were written before the custody fix. Baseline stale-create reproduction
exited `5` with `candidate draft missing or ambiguous`; baseline stale-upload
reproduction exited `5` with `missing asset`. The binary fixture also refused the
baseline download without `--allow-escape-sequences` (exit `90`).

Coverage includes stale create/upload lists, creation and exact-read ID/tag/body/
target/draft mismatches, upload receipt ID/name/size/state, same-size changed
remote bytes, immutable collisions, duplicate reservations, lost upload response,
interrupted initial custody, expired original artifacts after complete custody,
old-product/new-automation continuation, original receipt equality, stale-event
priority, failed/spoofed CI suites and legacy exact-source enforcement.

The finite serial batch completed with exit `0`. Every suite below exited `0`:

| Command | Result |
| --- | --- |
| `bash scripts/test-apple-release-state.sh` | PASS: custody identity, immutable state and interrupted-write regressions |
| `bash scripts/test-apple-custody-recovery.sh` | PASS: original reservation and partial/complete archive recovery |
| `bash scripts/test-prepare-apple-release.sh` | PASS: exact suites, origin, sealed metadata and stale-event continuation |
| `bash scripts/test-release-source.sh` | PASS: 16 dual-pin qualification scenarios |
| `bash scripts/test-apple-continuation.sh` | PASS: pending/failure blocking, final-byte reuse and fresh receipts |
| `bash scripts/test-prepare-macos-signing.sh` | PASS: 29 original-artifact/receipt/provenance scenarios |
| `bash scripts/test-publish-release.sh` | PASS: original publication suite plus separately qualified automation |
| `bash scripts/test-release-r1.sh` | PASS: existing loss/collision/security-cleanup invariants |
| `bash scripts/test-sign-macos.sh` | PASS: original synthetic submit-once, diagnostic, cleanup and promotion contracts |
| `bash scripts/test-verify-macos-package.sh` | PASS: original synthetic copied-app, mount-cleanup and trust-assessment contracts |

ShellCheck passed for all 12 changed/new shell files. `bash -n` passed for the
same files. Actionlint passed for `ci.yml`, `release.yml` and `release-macos.yml`.
`git diff --check` passed. These are local results only.

## Changed paths

- Custody: `scripts/apple-release-state.sh`, `scripts/prepare-apple-release.sh`.
- Shared qualification and consumers: `scripts/release-source.sh`,
  `scripts/prepare-macos-signing.sh`, `scripts/publish-release.sh`.
- Tests: `scripts/test-apple-release-state.sh`,
  `scripts/test-apple-custody-recovery.sh`, `scripts/test-prepare-apple-release.sh`,
  `scripts/test-release-source.sh`, `scripts/test-prepare-macos-signing.sh`,
  `scripts/test-publish-release.sh`, `scripts/testdata/apple-custody-gh.sh`.
- Workflows: `.github/workflows/ci.yml`, `.github/workflows/release.yml`,
  `.github/workflows/release-macos.yml`.
- Documentation: `docs/RELEASE-FLOW.md`, this handoff.

The HTML handoff and retained local validation/red/green logs are local chat
artifacts, excluded from version control.

## Delivery and qualification

Root should review the complete delta, deliver it to main through normal checks,
and retain the immutable merged automation SHA. Both exact main-push workflow
suites must succeed before automatic continuation can use that SHA. Existing
candidate `424598089` must resume from its sealed original metadata and archives;
there is no replacement candidate or manual dispatch in this repair.

Hosted custody recovery, credential cleanup, real Apple Accepted diagnostics,
stapling, Gatekeeper, fresh native install acceptance and publication/download
verification remain operational gates. Local fixtures establish none of them.
