# Unattended Apple release handoff

The owner authorized implementation through merge, real hosted Apple qualification
and release, including removal of the recurring `release-signing` reviewer.
Source was frozen at `84167fb93abf7090b7095b77b60fff1b67311d47` before edits.
Live environment `25335501562` had reviewer rule `70123082`, main-only deployment
policy `64384235`, disabled administrator bypass and seven environment secret
names. Main required strict `ci-required` and `conventional-commits` from GitHub
Actions, with administrator enforcement and no force-push/deletion permission.
The initial release/tag inventories were empty. Values were not read or exported.

[Release automation](RELEASE-FLOW.md) is the implementation/operator entry point.
The existing first-version policy produces `0.0.1`; automatic main releases use
patch bumps. Both Mac architectures retain app-ZIP then DMG notarization. Windows
and Linux remain unsigned. No updater trust authority or package manager is added.

The implementation separates credential-free versioned native builds, durable
candidate custody, environment-isolated Apple stages, fresh credential-free ARM
and Intel acceptance, and immutable publication. Original source/run/attempt and
Actions artifact ID/digest are retained across scheduled continuations. Private
candidate drafts retain append-only checkpoints and signed bytes beyond Actions
artifact expiry. Unknown submissions reconcile identity rather than resubmit.
Final download bytes remain frozen across architecture and publication retries.

Focused entry points:

```sh
bash scripts/test-apple-release-state.sh
bash scripts/test-apple-custody-recovery.sh
bash scripts/test-prepare-apple-release.sh
bash scripts/test-apple-continuation.sh
bash scripts/test-sign-macos.sh
bash scripts/test-prepare-macos-signing.sh
bash scripts/test-verify-macos-package.sh
bash scripts/test-publish-release.sh
bash scripts/test-packaging.sh
bash scripts/test-release-r1.sh
bash scripts/test-release-version.sh
shellcheck -x scripts/*.sh build/linux/*.sh
actionlint
git diff --check
```

Run the complete native CI script with `CGO_ENABLED=1` under the shared heavy-job
queue when using this agent environment. Source checks and synthetic fixtures do
not establish Developer ID, Apple acceptance, tickets, installed-app acceptance
or publication. Hosted qualification must retain each original run, submitted
hash/UUID/log, final package hash, fresh acceptance and public asset verification.
Use `APPLE_QUALIFY_CONTINUATION=true` only to yield after retaining each new real
UUID and prove continuation; it cannot produce Accepted evidence. Remove the
variable after qualification. Keep unresolved candidate drafts and their bytes.

At the initial handoff commit, real qualification remains pending. Its exact
results belong in the release's operational evidence, not a claim derived from
mock success. Existing failed run `40324314083` remains the pre-change witness:
its two signers failed before Apple submission and retained only prepared app
receipts. No particular command is established as its root cause.
