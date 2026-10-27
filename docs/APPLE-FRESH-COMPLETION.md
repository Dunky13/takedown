# Authoritative completion of the exceptional fresh ARM job

`bash scripts/verify-apple-fresh-completion.sh EVIDENCE_DIR` is a finite,
read-only completion gate for the single candidate `425202908` ARM exception.
It returns 0 only when the original fresh ARM job completed successfully and
its exact proof is independently recoverable. Every other outcome returns 1
with `HOLD`. A hold never permits another qualification, reservation, reset,
repin, signing operation or Apple submission. Default acceptance budgets remain
the caller's responsibility; this helper accepts only exception4.

An immutable verified4 checkpoint is provisional: cancellation can race with
its upload. Neither that checkpoint nor a local native receipt proves that
the original job completed. The helper requires the exact attempt-specific
Actions job, including successful automation qualification, transport setup,
native acceptance, cleanup and receipt promotion steps. A cancelled, failed,
pending, missing, ambiguous or unreadable original job cannot pass even when
both the Actions artifact and durable candidate receipt exist.

## Integration and ordering

Exceptional readiness invokes this helper after offline custody and byte audits.
Producer continuation, fresh replay and publication share that readiness gate.

1. Before creating the first exceptional native proof, capture the actual hosted
   verification run ID and attempt separately from original build
   `40470756538/1`. Add `verification_workflow_sha` from the already qualified
   `RELEASE_AUTOMATION_SHA` / `GITHUB_SHA`, with the existing equality check to
   the workflow SHA. Preserve it in checkpoint custody. No default or inferred
   SHA is accepted by this helper.
2. For this first exceptional qualification, copy `acceptance.json` into
   `fresh-acceptance-arm64.json` **byte for byte**. Do not serialize the fresh
   JSON again after the native proof is sealed. Preserve the existing deferred
   promotion, unconditional cleanup and upload gates. Keep the existing artifact
   name `macos-fresh-acceptance-arm64-${GITHUB_RUN_ATTEMPT}`.
3. After strict offline state, immutable parent, input, stage, Accepted-log and
   final-byte audits, invoke the helper in exceptional completed-proof replay
   and policy readiness. Any nonzero result propagates as a hold. Preserve the
   spent fourth attempt and every existing byte and UUID. Do not route a hold
   into a legacy/native qualification path.
4. **Do not invoke this completion gate inside the original job's first proof
   production.** The original job cannot be completed while it is producing
   that proof. Invoke it only from later consumers, replay and readiness after
   the fresh jobs finish. In particular, examine any shared readiness function's
   producer call sites when wiring it to avoid self-dependency.
5. Retain the exact native/fresh bytes through the existing bound state `put`
   operation after successful fresh jobs, under
   `fresh-acceptance-arm64-<original verification run>-<original verification attempt>.json`.
   Readiness derives this original identity from the proof, including on later
   replay, and also retains the current run's alias for existing publication
   consumers. A later publication run preserves the original asset. A newly
   named receipt belonging to the later run cannot replace it for this gate.
6. Include `verification_workflow_sha` in native/fresh equality audits and
   require the new helper in every exceptional success consumer. Root must add
   combined wrapper/readiness regression coverage after wiring. These focused
   tests exercise the helper boundary, including restored late-commit custody;
   they do not claim that an unwired wrapper has been repaired.

Required environment pins are `GH_REPO=Dunky13/takedown`,
`NOTARY_RELEASE_ID=425202908`, `NOTARY_RELEASE_SOURCE_SHA` equal to product
`ef82381d509aac8b5b02da0f784a8176b1c5b9c0`,
`NOTARY_STATE_TAG=apple-candidate-v0.0.1` and the already validated
`NOTARY_RELEASE_BODY_SHA256`. An existing GitHub token with Actions/Contents read
access suffices. No new secret or write permission is required by the helper.

The native JSON must carry the original product source, version `0.0.1`,
architecture `arm64`, original build string IDs, ZIP/DMG SHA256 values,
`acceptance_attempts:4`, `credential_free_runner:true`, fixed parent
`27b7a083cc181ed6def7479ec99b25479d7b11d9c3de36bbbd9facfbacfc60fe`,
positive string `verification_run_id` and `verification_run_attempt`, and
40-character lowercase hexadecimal `verification_workflow_sha`.
Caller offline audits remain mandatory; this helper does not replace them.

## Authoritative reads and proof binding

All calls explicitly use `gh api --hostname github.com --method GET`, REST
version `2022-11-28`, and the exact repository. There is no current/latest run
lookup, automatic retry, write, native tool, PF operation or production test
override. The whole read sequence has a 180-second time budget and each call
has at most 20 seconds. Inventories use explicit pages of 100, at most 20 pages,
with duplicate IDs and inconsistent Actions totals rejected. Larger inventories
hold; they never silently omit evidence. JSON duplicate keys are rejected.

The helper reads these endpoints:

```text
repos/Dunky13/takedown/actions/runs/<verification run>/attempts/<verification attempt>
repos/Dunky13/takedown/actions/workflows/<workflow id from original attempt>
repos/Dunky13/takedown/branches/main
repos/Dunky13/takedown/compare/<verification workflow SHA>...<current main SHA>?per_page=1
repos/Dunky13/takedown/actions/runs/<verification run>/attempts/<verification attempt>/jobs?per_page=100&page=N
repos/Dunky13/takedown/actions/runs/<verification run>/artifacts?name=macos-fresh-acceptance-arm64-<verification attempt>&per_page=100&page=N
repos/Dunky13/takedown/actions/artifacts/<exact artifact id>/zip
```

The attempt must originate in the same non-fork repository on main, from
`workflow_dispatch`, `workflow_run` or `schedule`, with no PRs. Its workflow
must be `.github/workflows/release.yml`, its head SHA must equal the retained
`verification_workflow_sha`, and its referenced `release-macos.yml` must bind
that same SHA and main ref. Protected main and ancestry are rechecked live.
This establishes current protection/history and the run's original main
provenance. Historical branch protection and original automation qualification
remain covered by the caller's existing trusted release gates.

The only qualifying job name is
`macos / acceptance (macos-15, arm64)`, with labels exactly `["macos-15"]`.
Its run/attempt/SHA/branch/API URLs must match; it must be `completed/success`.
All five named critical steps must occur exactly once, in order, with unique
step numbers and `completed/success`. A different green job cannot substitute.
The whole run need not be complete or successful: readiness can run while
publication is pending, and a later publication failure does not invalidate
a successful qualification job.

The Actions artifact must be unique, belong to the original run/repository/SHA,
and have been created within that original job's start/end interval. The
attempt is bound by its exact artifact name and the exact attempt-specific job;
Actions artifact metadata itself has no run-attempt field. The helper checks
the outer downloaded ZIP's size and SHA256 against authoritative metadata.
Its only entry must be the regular file `fresh-acceptance-arm64.json`, with
bytes identical to `acceptance.json`. Foreign paths, duplicate/extra entries,
directories, symlinks, encryption, corrupt ZIP/CRC and oversized evidence hold.
The archive is inspected in memory and never extracted. ZIPs without Unix type
bits are accepted only as non-directory files; explicit non-regular types fail.

For Actions ZIP downloads, the helper uses `--allow-escape-sequences` with the
normal JSON Accept header, never an octet-stream header. GitHub CLI follows
the download redirect. Raw gh stderr is discarded so private SAS URLs cannot
reach logs. Provider errors report a fixed safe hold message.

## Durable candidate fallback and limits

Only an empty successful Actions inventory or a correctly bound expired
artifact enables the durable fallback. Actions API errors, duplicate/foreign
metadata, download failures or corrupt bytes hold immediately. They never
silently select a different evidence source.

The fallback repeats all original-job success/provenance checks. It directly
GETs candidate release ID `425202908` and validates draft/tag/product source and
the caller's exact body SHA256. It explicitly inventories that release's assets,
rejects duplicate names/IDs, and selects only the original verification run's
`fresh-acceptance-arm64-<run>-<attempt>.json`. The asset must be uploaded, have
the expected size/API URL, and have been created after original job completion.
Its digest, when supplied by GitHub, must match the proof. Before download,
the helper rechecks the same release seal and ID. The exact asset-ID download
uses `--allow-escape-sequences -H 'Accept: application/octet-stream'` and must
reproduce the native proof bytes exactly. It performs no candidate discovery,
cache write, put/save/reset or receipt creation.

If the original workflow/job API history itself is unavailable, completion
holds even if candidate assets survive. This fallback extends proof retention;
it cannot replace authoritative original job success. Missing retained proof
also holds and never buys another native attempt.

## Focused local verification

```sh
bash scripts/test-apple-fresh-completion.sh
bash -n scripts/verify-apple-fresh-completion.sh scripts/test-apple-fresh-completion.sh
shellcheck scripts/verify-apple-fresh-completion.sh scripts/test-apple-fresh-completion.sh
```

Fixtures intercept every GitHub call, enforce exact GET commands and binary
headers, refuse native/Apple/PF/signing commands, and assert evidence/counter
bytes remain unchanged. Hosted qualification, Apple trust, real transport,
integration activation and publication remain **NOTRUN** by this worker.

## Primary API references

- [Get a workflow run attempt](https://docs.github.com/en/rest/actions/workflow-runs#get-a-workflow-run-attempt)
- [List jobs for a workflow run attempt](https://docs.github.com/en/rest/actions/workflow-jobs#list-jobs-for-a-workflow-run-attempt)
- [Actions artifact metadata and download](https://docs.github.com/en/rest/actions/artifacts)
- [Get a branch](https://docs.github.com/en/rest/branches/branches#get-a-branch)
- [Compare two commits](https://docs.github.com/en/rest/commits/commits#compare-two-commits)
- [Release assets](https://docs.github.com/en/rest/releases/assets)
- [GitHub CLI API flags and pagination](https://cli.github.com/manual/gh_api)
