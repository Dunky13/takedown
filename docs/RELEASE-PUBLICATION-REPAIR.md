# Cross-platform release publication repair

Repair: [PR #14](https://github.com/Dunky13/takedown/pull/14).

## Failure and impact

[Release run 40574763682](https://github.com/Dunky13/takedown/actions/runs/40574763682)
at automation source `3f106784ce62c7bad9afbd16fbf24d4ba0a107ce` passed both
macOS signing jobs, both fresh native acceptance jobs and readiness. Publication
failed with `release tag, target or version metadata differs from plan` immediately
after creating the version draft. Native CI for that automation source passed.

Read-only inspection found version draft ID `425937228`, tag `v0.0.1`, target
`ef82381d509aac8b5b02da0f784a8176b1c5b9c0`, zero release assets and
`prerelease=false`. The candidate remains draft ID `425202908`,
`apple-candidate-v0.0.1`, with original run `40470756538`, attempt `1`.
These observations are a snapshot, not permission to change those objects.

The publisher discarded the creation response and immediately queried the release
list to rediscover the draft. A deterministic fixture that omits the newly created
draft from list responses reproduces the same error. The original failed list
response was not retained, so its exact contents cannot be recovered from the log.
The later direct-ID response has the expected identity. This repair eliminates
the reproduced visibility dependency without relaxing metadata validation.

All seven platform downloads and SHA256SUMS pass through this publication job.
Its failure therefore prevents Linux and Windows downloads as well as both Mac
architectures. The unresolved stable candidate takes priority over new nightly
planning, so the same failure also holds back nightlies.

## Repair

`scripts/publish-release.sh` retains the REST creation response, validates it and
reads the exact release ID. Existing drafts are discovered once, then checked by
ID. Uploads use that ID, validate each returned asset ID/name/size/state and
compare its downloaded bytes. The final full inventory and checksum checks still
precede publication. PATCH publication and subsequent GET both validate the
original numeric release ID, tag, source, prerelease flag and published state.

Uncertain write responses stop the invocation. A later invocation discovers and
verifies retained state, appends only missing assets or verifies an already
published release read-only. Conflicting metadata, bytes, asset names and tags
remain terminal. No overwrite, forced tag, new Apple submission, signing change
or notarization bypass is introduced.

Stable, manual and scheduled nightly versions share the repaired publisher.
Nightlies continue to require signed/notarized macOS release packages and remain
prereleases outside Latest. Linux TAR.GZ and Windows ZIP/setup EXE remain unsigned.
No schedule or version policy changes are needed for this failure.

## Verification

Run these from the repository root with Go 1.26.10, Bash, jq, Ruby, ShellCheck,
Actionlint and the existing packaging prerequisites available:

```sh
bash scripts/test-publish-release.sh --publication-recovery-only
bash scripts/test-publish-release.sh
bash scripts/test-prepare-apple-release.sh
bash scripts/test-release-source.sh
env -u VERSION -u RELEASE_VERSION bash scripts/test-release-version.sh
bash scripts/test-ci-topology.sh
go test -p 2 -parallel 2 ./internal/releaseflow ./cmd/releasectl
shellcheck -x scripts/*.sh build/linux/*.sh
actionlint
git diff --check
```

The focused regression failed before the repair with the same metadata error.
The tests exercise the real publisher with synthetic GitHub responses and Apple
evidence. Coverage includes persistently stale draft lists for stable/prerelease
and nightly publication, lost create/publication responses, malformed write
receipts, changed direct-ID identities, partial uploads, corrupt remote bytes,
collisions, nightly channel metadata and read-only replay. Existing Apple evidence
refusal cases remain in the complete suite. Synthetic success establishes neither
new Apple acceptance nor successful hosted publication.

All commands above passed locally, including the complete publication suite and
the CI topology mutation checks. The independent source review found no
actionable issues. The queue client lost its final response; the retained job
result records exit code zero and the complete log contains every suite result.
Hosted PR checks and protected-main publication remain separate gates.

## Operational handoff

Merge remains an owner approval gate. Once the repair is on protected main with
successful exact-source Native CI and Conventional commits, the existing release
workflow can continue the unresolved candidate. Keep the original source, version,
run/attempt, archive bytes, Apple UUIDs and retained proofs. Do not delete either
draft or rerun a build to recover publication. The next publisher should discover
draft `425937228` and append its eight verified assets by ID.

The public outcome remains unverified until the hosted publication job succeeds,
the public tag resolves to the original source, the complete public asset set
matches SHA256SUMS and downloaded bytes match the retained final artifacts.
Nightly runtime qualification requires its own later eligible scheduled run;
passing a fixture does not establish that result.
