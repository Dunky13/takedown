# Apple provider import and custody repair

Local implementation handoff, 2026-10-26. Branch `fix/apple-provider-qualification`,
base `22c78b82b9b7f4c732c885007f76ff5afe73bdb2`. Root owns review, delivery,
hosted qualification, Apple operations, and publication.

## Result and boundary

The original password-empty PKCS#12 transfer fails reproducibly on the local Apple
importer. A nonempty password alone imports the identity, but the local explicit
`-f pkcs12` path ignores `-x` and makes the synthetic key extractable. The final
repair uses a fresh random transfer password, explicit interoperable PKCS#12
encoding, and format autodetection, whose native import API honors `-x`.

The final shared production helper imports the synthetic certificate/private-key
identity and its native key metadata reports **nonextractable**. These are local
compatibility findings, not proof of the historical hosted cause: that private
error output was destroyed, and the hosted and local macOS versions differ. Do
not infer incorrect Developer ID credentials from exit 1 or a generic authentication
status; the supplied PKCS#12 had already decoded successfully in the hosted run.

## Historical input evidence

- Exact automation base CI `40387142306` and Conventional `40387142493` succeeded
  (reported delivery evidence; not rerun here).
- Automatic continuation `40391233601`, attempt 1, retained original product
  `5b4bf4dea8f93a3a3d23467f990bb3310a3fea9b`, version `0.0.1`, original release
  `40372460129`, attempt 1, and draft candidate `424598089`.
- Native build skipped; durable custody retained the same four artifact IDs and
  digests (reported delivery evidence).
- Retained ARM job log `121813037370` confirms `stage=app command=keychain-import
  exit=1`. Runner: macOS `15.7.9` (`24G830`), image
  `macos-15-arm64/20260907.0337.1`, selected Xcode `26.3`.
- The exact hosted OpenSSL executable/version was not logged. The image inventory
  lists OpenSSL `1.1.1w`; that does not establish the executable selected by PATH.
- Latest safe receipt remains `prepared`, with no UUID or submitted hash; no Apple
  upload occurred (reported checkpoint evidence). No retained receipt/state was
  edited, copied into another candidate, or deleted by this implementation.
- Intel inventory failure remains separate and root-owned.

## Differential probe evidence

Only newly generated RSA-2048 keys and a one-day self-signed code-signing
certificate were used. Probes exercised supplied-PKCS#12 decode, temporary export,
and native import independently. Each case used a freshly created isolated
keychain with a strong random password passed only through private stdin.

Local host: ARM64 macOS `27.0.1` (`26A434`), Bash `3.2.57`. Export tools tested:
OpenSSL `4.0.3`, OpenSSL `3.5.9`, and system LibreSSL `3.3.6`.

| Explicit `-f pkcs12` path | Password | Result with all three export tools |
| --- | --- | --- |
| Default export encoding | Empty | Exit 1, mapped authentication status -25293; no certificate/key identity |
| SHA-1 MAC only | Empty | Same failure |
| PBE-SHA1-3DES for both bags only | Empty | Same failure |
| PBE-SHA1-3DES plus SHA-1 MAC | Empty | Same failure |
| Default export encoding | Fresh random nonempty | Exit 0; certificate, private key and code-signing identity present |
| Default encoding, deliberately wrong import password | Nonempty | Exit 1; no identity; private stdin propagates the failure |

A subsequent native attribute check against the nonempty explicit-format import
found `kSecKeyExtractable` equal to 1, despite `-x`. Current Apple public CLI source
corroborates this: its explicit PKCS#12 branch calls `SecPKCS12Import` without
forwarding the nonextractable attributes. Format autodetection instead calls
`SecKeychainItemImport` with those attributes. No private-key export was attempted.

| Autodetection path, nonempty password; OpenSSL 4.0.3 | Result |
| --- | --- |
| Default encoding | Exit 1; mapped MAC-verification status -25264 |
| SHA-1 MAC only | Exit 1; unknown-format status -25257 |
| PBE-SHA1-3DES for both bags only | Exit 1; MAC-verification status -25264 |
| PBE-SHA1-3DES for both bags plus SHA-1 MAC | Exit 0; identity present; native extractable attribute 0 |

A separate autodetection probe with **empty** passwords failed for all four
encoding combinations. Encoding and password handling therefore both need to be
explicit. Final production-helper tests passed separately with all three listed
OpenSSL/LibreSSL tools and verified native nonextractability. These results apply
to the local importer; the actual macOS 15 hosted import remains a provider gate.

## Exact implementation paths

- `scripts/macos-certificate-import.sh`: shared production/test seam. Generate a
  fresh 32-byte random transfer password in a private file; export with explicit
  `-keypbe PBE-SHA1-3DES -certpbe PBE-SHA1-3DES -macalg sha1`; import through private
  `security -i` stdin, with format autodetection and `-x -T /usr/bin/codesign`.
  Validate the password and quoted paths. No password appears in process argv.
- `scripts/sign-macos.sh`: call the helper, remove the transfer-password file
  alongside certificate scratch files, and emit safe import classification before
  private-log destruction. Existing failure phase labels and receipt bindings stay.
- `scripts/test-sign-macos.sh`: exact diagnostic allowlist tests, private-stdin and
  import-flag contracts, authentication/unclassified failures, cleanup/no upload on
  failure, and inclusion of the helper in the disposable packaging fixture.
- `scripts/test-macos-certificate-import.sh`: finite native regression exercising
  the real helper with synthetic input, identity and native nonextractability
  checks, wrong-password failure classification, and exact cleanup verification.
- `docs/APPLE-PROVIDER-REPAIR.md`: this handoff.

The transfer encoding is used only for private, temporary local rewrapping. It
changes no supplied certificate, Developer ID identity, team, trust root, release
artifact binding, or signing algorithm. No `-legacy` provider dependency, `-A`, empty
keychain password, permission bypass, default/system keychain write, or interactive
prompt was added. The final synthetic key attribute check confirms nonextractability
on the local OS; the codesign-only trusted-application flag remains enforced at the
command boundary. Hosted native behavior still requires qualification.

The signer implementation left `scripts/keychain-lifecycle.sh`, workflows, source
bindings, and root's inventory repair files unchanged. During integration, root
added the native synthetic import regression to the existing credential-free
macOS CI step. Both macOS 15 architectures must therefore verify the helper,
native nonextractability, wrong-password refusal and cleanup before merge.

## Concurrent private custody inventory

The combined patch also permits GitHub's documented `starter` state in the
private draft inventory while another architecture uploads an unrelated asset.
Every requested asset still requires exactly one matching `uploaded` entry.
Incomplete required assets cannot be downloaded, reused or overwritten; unknown
states and duplicate names/IDs still fail. Upload receipts and downloaded bytes
retain their existing strict checks. Public publication inventory is unchanged.

The focused inventory fixture failed with the original global uploaded-only
check, then passed with the scoped repair. It also verifies that a required
incomplete asset causes no new upload. This reproduces a failure mode; the exact
historical Intel response was not retained, so its cause remains unproved.
Candidate `424598089`, original source/run/artifact bindings and existing receipts
are preserved.

The integrated custody paths are `scripts/apple-release-state.sh`,
`scripts/test-apple-release-state.sh`, and `scripts/testdata/apple-custody-gh.sh`.
The macOS CI addition is in `.github/workflows/ci.yml`. No live environment or
branch policy is changed by this patch.

## Safe failure classification

Only complete known `security: SecKeychainItemImport:` messages, optionally prefixed
by the tool prompt, map to fixed statuses and constant causes:

| Mapped OSStatus | Constant cause |
| --- | --- |
| -25293 | authentication-failed |
| -25264 | pkcs12-mac-verification-failed |
| -26275 | decode-failed |
| -25257 | unknown-format |
| -25308 | interaction-not-allowed |
| -50 | invalid-parameters |
| unknown | unclassified |

The numeric value is mapped from a known constant message, never extracted from
arbitrary private text. Extra private suffixes and numeric lookalikes remain
unclassified. The original line, key material, and passwords never leave scratch.
Diagnostics supplement the existing safe command/exit job-log label; receipt
schema and artifact allowlists are unchanged. Localized or new messages can remain
unclassified. Authentication/MAC status alone does not identify incorrect supplied
credentials or distinguish all native format/policy failures.

## Red/green, validation, and cleanup

1. **Original chain red:** introduce shared helpers retaining original export/import
   behavior; the native test fails at import with exit 1.
2. **Attribute red:** nonempty password with explicit format imports an identity,
   but the native metadata check proves it extractable and fails the test.
3. **Encoding red:** autodetection with default encoding fails at native import;
   the differential matrix separates MAC, bag encryption, and password behavior.
4. **Final green:** production helper imports a nonextractable synthetic identity;
   wrong-password native import fails with an allowlisted classification. All
   three listed OpenSSL/LibreSSL tools pass against the local importer.
5. The complete synthetic signer/packager suite passes after the final format and
   encoding changes. Its first run found a missing helper copy in the disposable
   packaging fixture; that dependency was repaired and the failure witness retained.
6. Focused ShellCheck, Bash syntax checks, and `git diff --check` pass. No unrelated
   Go/native builds or broad analyzers were rerun.

Every native run verifies the exact original user search list/order and default
keychain byte-for-byte, deletes its owned keychain, and checks removal of private
synthetic material. Traps handle success, failure, INT, and TERM. Each native
Security command has a 15-second alarm. No real credential or user-keychain private
key was accessed. No Apple request, manual workflow dispatch, live repository or
release-environment mutation, commit, or push occurred.

Run from the repository root on macOS, one native probe at a time to avoid search
list restoration races with another signer:

```sh
bash scripts/test-macos-certificate-import.sh
bash scripts/test-sign-macos.sh
shellcheck -x scripts/macos-certificate-import.sh scripts/test-macos-certificate-import.sh scripts/sign-macos.sh scripts/test-sign-macos.sh
for script in scripts/macos-certificate-import.sh scripts/test-macos-certificate-import.sh scripts/sign-macos.sh scripts/test-sign-macos.sh; do bash -n "$script" || exit; done
git diff --check
```

Expected native output includes tool/OS versions, nonextractable metadata,
identity/negative-case passes, exact keychain restoration, and private-material
removal; expected exit is 0. The signer suite reports synthetic release/recovery
passes; it does not establish Apple acceptance. Queue the signer suite through the
environment's shared heavy-command wrapper when required. The native regression is
also runs in the existing macOS CI regression step on both native architectures.

Root's integrated ten-suite batch also passed: native certificate import, signer,
private custody state, original custody recovery, planning, source qualification,
automatic continuation, macOS artifact handoff, final package acceptance and
immutable publication. Focused ShellCheck, syntax, Actionlint and diff checks
passed. Root source review covered the nine integrated paths with no unresolved
finding. This is bounded source/local qualification; hosted CI, Apple proof and
publication remain separate gates.

## Remaining provider gates

Root independently reviews the combined signer and inventory patch, delivers it
through CI, and resumes the **same candidate** with the original product/run/artifact
provenance. The prepared receipt remains reusable only under its existing binding
and absence-of-upload checks. Submitted/unknown UUID states retain the existing
recovery path; this patch adds no automatic resubmission.

Both architectures still require real Developer ID import/signing, secure timestamps,
submit-once UUID capture/recovery, explicit Apple Accepted diagnostics, app and
outer-DMG stapling/validation, Gatekeeper, copied/installed package and fresh
acceptance checks, credential cleanup before upload/promotion, and publication.
Hosted success, provider acceptance, and publication remain unqualified by these
local tests. Signing-reviewer environment policy changes belong to root.

## Primary references

- [OpenSSL PKCS#12 options and interoperability notes](https://docs.openssl.org/3.5/man1/openssl-pkcs12/)
- [Apple Security CLI import implementation](https://github.com/apple-oss-distributions/Security/blob/main/SecurityTool/macOS/keychain_import.c)
- [Apple PKCS#12 decoder](https://github.com/apple-oss-distributions/Security/blob/main/OSX/libsecurity_pkcs12/lib/pkcs12Decode.cpp)
- [Exact hosted runner image inventory](https://github.com/actions/runner-images/blob/macos-15-arm64/20260907.0337/images/macos/macos-15-arm64-Readme.md)

Installed `security help import`, `security help create-keychain`, `security help
error`, `security error`, `openssl pkcs12 -help`, and SDK `SecKey.h` attribute
constants were consulted. Current public source is explanatory; it is not asserted
to match the hosted OS binary exactly.
