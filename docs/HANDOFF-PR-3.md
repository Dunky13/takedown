# PR #3: restore macOS signing secret resolution

PR: <https://github.com/Dunky13/takedown/pull/3>

## Outcome and scope

Forward the seven named macOS signing inputs from `release.yml` to
`release-macos.yml`, and declare each input as a required workflow-call secret.
GitHub requires caller forwarding even when values exist only in the called
job's environment. Environment values take precedence after approval.

Protected `main`, the `release-signing` required reviewer, and scheduled unsigned
nightlies retain their existing behavior. No secret values or repository settings
were changed. Native CI receives no signing secrets.

## Failure witness

[Run 40305764527](https://github.com/Dunky13/takedown/actions/runs/40305764527)
was a manual release for `0.0.1`, not a scheduled nightly. Native CI passed on
all four targets. After environment approval, both macOS signing jobs failed
with `missing required secret/configuration: MACOS_CERTIFICATE_P12_BASE64`.
Their environment entries were empty; publication was skipped. The protected
environment contained all seven configured secret names before the run.

## References checked

- [GitHub reusable workflow secret contract](https://docs.github.com/en/actions/how-tos/reuse-automations/reuse-workflows#using-inputs-and-secrets-in-a-reusable-workflow).
- Hikyo's `release.yml` uses direct jobs bound to `apple-signing`; run
  [40171412797](https://github.com/Hikyo-Org/Hikyo/actions/runs/40171412797)
  reports successful Apple signing, notarization, and publication.
- Agent Job Queue's `distribution.yml` uses a direct protected signing job;
  run [40026971766](https://github.com/Dunky13/agent-job-queue/actions/runs/40026971766)
  reports successful signing/notarization, both architecture acceptance jobs,
  and publication.
- Read the T3 histories titled Add Apple Developer Signing, Separate macOS
  Notarization CI Flow, and Reduce Apple Notarization CI Failures. They explain
  separate artifact verification and recovery for accepted Apple submissions.
  Those later-stage failures do not explain this run's empty credential context.

## Validation and remaining gates

Actionlint 1.7.12 passes. In disposable copies of the workflows, removing each
of the seven caller mappings separately produces a required-secret validation
failure; restoring every mapping passes. ShellCheck and shell syntax checks pass
for the signer, macOS packager, and keychain helper. `git diff --check` passes.
These checks do not authenticate Apple credentials or prove notarization.

Hosted PR checks and review are recorded on the PR. The repaired signing path
must run from protected `main` after merge and the existing environment approval.
Re-running the old failed run keeps its old workflow definition, so it cannot
test this repair. Start a new authorized release run after merge. Successful
signing includes Developer ID verification, Apple acceptance, and stapled ticket
validation for the application and DMG on both macOS architectures.

Do not remove approval, copy credentials into repository secrets, accept unsigned
packages as signed, or suppress an Apple failure to make the jobs green. If an
Apple submission has already started, reconcile its recorded identifier before
retrying. See [release recovery](RELEASE-FLOW.md#failure-recovery).
