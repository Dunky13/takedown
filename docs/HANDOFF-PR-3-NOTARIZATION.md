# PR #3 follow-up: macOS submission evidence and acceptance

The secret-forwarding repair in [PR #3](https://github.com/Dunky13/takedown/pull/3)
is merged. The owner authorized a new `v0.0.1` release. Before dispatch, the newly
supplied macOS distribution policy required changes to the existing submission
helper. No real Apple submission or public release occurred during this repair.

## Resulting flow

1. Credential-free Native CI builds, tests and packages the selected version.
2. The protected macOS job selects Xcode 26.3 and checks its installed tools.
   It downloads the exact native artifact from the same workflow run, verifies
   GitHub's archive digest, package checksums, source and clean-build metadata,
   bundle identity, numeric and CLI versions, and actual architecture slices.
3. Only then does the signing step receive Apple inputs. It reuses the verified
   binary, signs the executable and application, and uploads an app ZIP whose
   basename identifies product, version, architecture, source, run, attempt,
   and stage. The DMG receives the same naming treatment before submission.
4. Each stage reserves persistent evidence, retains its signed upload bytes,
   submits once, records the UUID and submitted hash, and waits finitely. The
   matching Apple log is retrieved even on acceptance. The response and log must
   both say `Accepted`; IDs and supplied filename/hash fields must match.
5. Stapling, strict signatures and Gatekeeper must pass. Final acceptance checks
   the ZIP, mounted DMG, and copied app, including executable equality, metadata,
   architecture and CLI version. Cleanup failures remain job failures.
6. Only complete successful candidates reach the existing publication job.

The artifact-only manual signer uses the same handoff and acceptance contracts.
Unsigned scheduled nightlies retain their existing behavior. Environment approval
and protected-main restrictions remain intact.

## Credential and recovery boundaries

The persistent P12 password enters OpenSSL through its environment, with private
temporary rewrapping for import. Temporary keychain passwords enter separate
`security -i` commands through private stdin, with each command's exit status
propagated. The exact original keychain search list is restored and verified.
Credential scratch and raw client output are not uploaded.

Allowlisted CI evidence includes `input.json`, `acceptance.json`, per-stage
`receipt.json`, sanitized `apple-log.json`, and descriptive signed ZIP/DMG uploads.
Submitted hashes differ from final download hashes after stapling/repacking.
Evidence is retained privately for seven days on failure as well as success.

An existing local stage refuses another submission. A later workflow attempt
stops before signing because a new runner cannot assume the old submission did
not occur. Reconcile retained UUIDs with Apple before retrying. A lost UUID requires
history reconciliation against the retained filename/hash. Pending, rejected,
missing or mismatched evidence never passes. See
[recovery](RELEASE-FLOW.md#failure-recovery).

## Evidence gathered

- 19 synthetic signer scenarios: descriptive app/DMG uploads, retained bytes and
  diagnostics, duplicate prevention, pending/rejected/unknown outcomes, wait
  timeouts, missing/malformed/mismatched logs, stapling, Gatekeeper, signature,
  private stdin command failure and keychain cleanup failure.
- 22 handoff scenarios: source/run/artifact binding, digest and package corruption,
  archive traversal/symlink rejection, bundle fields, native architecture and
  clean executable metadata.
- 13 final acceptance scenarios: correct outputs and rejection of changed copied
  apps, bad checksums, extra manifest entries, differing executables, signature,
  ticket, Gatekeeper, attach and detach failures. No quarantine removal occurs.
- Existing offline release-version, publication and R1 packaging regressions pass.
  Actionlint 1.7.12, ShellCheck, Bash syntax and whitespace checks pass.
- A prior trusted release's actual 55 MB ARM64 artifact matched GitHub's digest,
  expected archive layout, source revision and clean-build metadata.
- An independent combined review found a temporary-keychain password in argv;
  the private-stdin repair and its failure regression address that finding.

Hikyo's signing failure history and current workflow explain the explicit Xcode
26.3 selection and temporary-keychain registration. Both current macOS 15 runner
manifests list that toolchain. Agent Job Queue's evidence and recovery design
informed submit-once receipts and the separate acceptance boundary.

## Remaining work

Hosted CI must qualify this follow-up before it reaches protected `main`.
After merge, start the already authorized release with:

```sh
gh workflow run release.yml --repo Dunky13/takedown --ref main -f bump=patch
```

The existing environment reviewer must approve credential use. Confirm planning
still selects `0.0.1` before approval. Check both architectures' real Developer ID
signatures, Apple acceptance/logs, stapled tickets, Gatekeeper and copied-app
acceptance, then the seven platform assets and combined checksums on the published
release. Local synthetic results do not qualify any of those hosted trust gates.
