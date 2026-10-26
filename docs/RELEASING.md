# Native packaging and signing

This tooling builds the single CLI + native Fyne GUI executable locally and in
GitHub Actions. Go is pinned to **1.26.10**, Fyne to **v2.8.2**, and native builds
require Cgo. The optional `make headless` output is never a desktop package input.
No source dependencies or licensing decisions are added by this tooling.

## Scope and ownership

The tooling lives in `.github/`, `build/`, `scripts/`, `FyneApp.toml`, and
`Makefile`. This document is the entry point.
Existing build/test/check/cross-compile/clean Make targets are retained.
The application code, dependency files, OAuth implementation, and policy hosting
have separate responsibilities. Native CI produces unsigned packages for review;
it does not publish releases or use signing credentials. Paid Windows signing is
deferred. The manual signing workflow remains a separate, explicitly authorized
operation. See [CI builds](CI-BUILDS.md) for
build targets, artifact retrieval, and validation boundaries.

[Release automation](RELEASE-FLOW.md) adds automatic protected-main patch releases,
manual minor/major/custom versioning, Conventional Commit changelogs, signed macOS
publication, and change-aware signed macOS nightlies. Its activation and recovery requirements
are documented separately from these packaging primitives.

## Prerequisites

Run from the repository root. Install prerequisites explicitly; package scripts
do not silently download tools or substitute a headless build.

| Host | Requirements |
| --- | --- |
| macOS arm64/amd64 | Go 1.26.10, Apple C compiler + SDK, Bash, `plutil`, `ditto`, `xattr`, `shasum`, `hdiutil`. Signing additionally needs `security`, `codesign`, `openssl`, and Xcode's `notarytool`/`stapler` through `xcrun`. |
| Linux arm64/amd64 | Go 1.26.10, GCC, pkg-config, GNU tar, coreutils (`link`, sha256sum), gzip, Bash, X11/Wayland and OpenGL/EGL development files. The pinned GLFW module enables both Linux backends by default. Ubuntu 24.04 graphics/build packages: `gcc pkg-config libgl1-mesa-dev libegl-dev xorg-dev libwayland-dev libxkbcommon-dev libdbus-1-dev`. A graphical session and distro runtime libraries are needed to launch the GUI; the OS keyring needs a working session service. |
| Windows amd64 | Go 1.26.10, PowerShell 7 (`pwsh`), native MinGW-w64 GCC on PATH, Unicode NSIS 3.x (`makensis.exe`; CI pins Chocolatey NSIS 3.11). Signing additionally requires Windows SDK `signtool.exe` on PATH. |

CI uses `ubuntu-24.04`, `macos-15`, and `windows-2025`. The Windows runner must
provide its native MinGW-w64 compiler; a missing compiler is an explicit failure.
Runner images and OS package repositories are mutable; Action commits and Go/Fyne
versions are pinned, but this is not a hermetically pinned SDK image. NSIS and
PowerShell flags were checked against primary documentation; Windows execution
still needs the Windows CI run described below. No Fyne CLI tool is installed or
used. `FyneApp.toml` is informational; `PACKAGE_VERSION`/`-Version` drives packages.

## Unsigned local packages

These are repeatable commands with fixed artifact names and an explicit input
inventory. A fresh output directory is required to prevent stale artifacts being
mistaken for successful results. Unix packaging reserves that directory with an
exclusive `mkdir` before building; a concurrent caller fails without modifying
the winner. A `.package-incomplete` marker remains until all artifacts are placed
in the reserved directory. Consumers must wait for successful command completion
and marker removal. This is exclusive publication, not an atomic multi-file
filesystem transaction. Ordinary failure cleans only the operation's own
reservation; an interrupted host may leave an incomplete directory for inspection.
The default package version is `0.0.0`, meaning
local development; select a numeric `x.y.z` explicitly for a candidate. Windows
version components must be at most 65535. This does not choose a release version.

```sh
# macOS: reuse the existing default native binary; no build required.
make package-macos BINARY=bin/takedown PACKAGE_VERSION=0.0.0

# macOS: build native CLI + GUI and package it instead.
make package-macos PACKAGE_VERSION=0.0.0

# Linux, on a Linux host with the graphics development packages installed:
make package-linux PACKAGE_VERSION=0.0.0
```

```powershell
# Windows, from the repository root in PowerShell 7:
./scripts/package-windows.ps1 -Version 0.0.0
# Or reuse an existing native Windows binary:
./scripts/package-windows.ps1 -Version 0.0.0 -Binary bin/takedown.exe
# GNU Make on Windows also provides: make package-windows
```

Unix options are `VERSION`, `BINARY`, `MODE`, and `OUT_DIR` environment variables;
Make exposes `PACKAGE_VERSION`, `BINARY`, `MODE`, and `OUT_DIR`. PowerShell options
are `-Version`, `-Binary`, `-Mode`, and `-OutputDirectory`. Choose another output
directory for a subsequent run. Default directories are
`dist/packages/<os>-<arch>-local`; the macOS OS label is `darwin`.

Reused binaries are checked with `go version -m`: exact Go and Fyne versions,
target OS/architecture, Cgo enabled, and no `headless` or Fyne `ci` build tags.
Packaging does not claim the reused binary matches today's uncommitted source.
The binary's recorded VCS revision, dirty flag, and compiler metadata are retained
in `BUILD-INFO.txt`; its local filesystem path is omitted. The existing binary's
CLI version remains its own build version, even if the package version differs.
Fresh builds embed the requested package version. Source and reused builds must
be qualified separately by the integrating parent.

| Platform | Output and installation |
| --- | --- |
| macOS | `takedown.app`, ZIP containing the app, compressed DMG with an Applications shortcut, and `SHA256SUMS`. Copy the app to Applications. CLI commands use the same `takedown.app/Contents/MacOS/takedown` executable. The bundle currently uses the OS generic app icon. |
| Windows | Portable ZIP, per-user NSIS `*-setup.exe`, and `SHA256SUMS`. Installer writes only the current user's Programs directory, Start Menu shortcut, and HKCU uninstall registration; it requests no elevation and makes no machine PATH change. The same console-subsystem executable serves CLI and GUI, so launching the GUI may show a console. |
| Linux | TAR.GZ with `bin/`, desktop entry, SVG icon, docs, `install-user.sh`, `uninstall-user.sh`, `install-common.sh`, and `SHA256SUMS`. Extract it, then run `bash install-user.sh`. Keep all three scripts together. Default prefix is the user's `.local`; use the same `INSTALL_PREFIX` for install and uninstall. A prefix containing spaces works; characters requiring ambiguous desktop escaping are rejected. Add its `bin` to PATH for CLI use. |

Linux installation checks every destination ancestor and payload path before
writing. Prefixes must have real directory ancestors without symlinks, dot path
components, repeated slashes, or a trailing slash; select their physical path
when the default path has a symlinked ancestor. An existing executable, launcher,
icon, application docs directory, or ownership directory is a conflict: the tool
never adopts or overwrites an old installation. Uninstall an owned installation
first; there is no in-place upgrade mode. Cooperative install/uninstall operations
share an exclusive prefix lock. After a crash, investigate a leftover lock rather
than deleting it while another operation may be active.

Payloads are staged in the selected prefix and published with the POSIX `link`
utility, which refuses existing destinations, including symlinks and directories,
instead of interpreting a destination directory as a place to append a basename.
The target filesystem must support hard links from the prefix staging area.
Only after all payload files succeed is a versioned SHA256 ownership manifest
published at `share/takedown/install-manifest.sha256`. Uninstall validates the
whole manifest, allowed relative paths, destination ancestry, and all extant
owned-file hashes before removing any file. Missing/invalid ownership, symlinks,
or modified owned files cause a refusal with the manifest retained. Only recorded,
unchanged files are removed; unknown docs and other files are preserved and
application directories are removed only if empty. A stale pre-manifest install
must be resolved manually; this tool does not guess its ownership.

Windows uninstall checks executable deletion before removing recovery metadata.
Interactive removal allows at most two explicit close-and-retry attempts after
the first failure; silent removal fails immediately with worker exit code 10.
The uninstaller, shortcut and registry entry remain available when the executable
is locked. Close the process or release the lock, then rerun the registered
uninstaller. No automatic reboot removal is scheduled.

Installers retain archives, partial files, job receipts, preferences, and keyring
tokens. Uninstalling is not OAuth revocation or data erasure; see the included
privacy policy. No `.deb`, AppImage, automatic updater, or privileged installation
is claimed. Linux release mode fails closed because no Linux signing trust policy
has been selected; Linux local archives remain available for independent review.

Verify checksums **from the artifact directory**, for example:

```sh
cd dist/packages/darwin-arm64-local
shasum -a 256 -c SHA256SUMS
# Linux artifact directory: sha256sum -c SHA256SUMS
```

Unsigned Linux archives normalize tar ordering, ownership, mode, and timestamps
(`SOURCE_DATE_EPOCH`, default zero), with a timestamp-free gzip header. Go VCS
metadata still reflects the input tree. ZIP, NSIS, DMG files and signatures carry
platform timestamps/metadata: no cross-host or signed bit-for-bit reproducibility
claim is made. SHA256SUMS records final archive/installer bytes, not a signature
or a claim of publisher identity.

## Package contents and privacy

Packages use a fixed allowlist: the native executable, platform integration files,
README, privacy Markdown/HTML, build metadata, and an existing root LICENSE or
COPYING file if one exists. The app receives no OAuth client JSON, client secret,
token, browser profile, job data, signed download URL, or archive payload. No
recursive repository bundling occurs. macOS preserves extended attributes,
including quarantine, and normalizes bundle permissions before signing.

A license is an owner decision. Missing licensing terms are not replaced with a
template license or an assertion that distribution is authorized. Artifact
creation does not settle license obligations for the source or dependencies.

Google verification and public distribution requirements remain independent of
build/signature success. The [Google API Services User Data Policy](https://developers.google.com/terms/api-services-user-data-policy)
requires accurate disclosures and a public privacy-policy URL for public OAuth
use. The [Data Portability developer policy](https://developers.google.com/data-portability/policy)
adds purpose, minimum-access, secure-handling and Limited Use restrictions. This
tooling makes no claim of Google verification, compliance approval, or live-account
acceptance. The packaged local privacy files do not host a public policy or
configure an OAuth consent screen. Those require a separately authorized task.

## Native CI and trust separation

`ci.yml` runs on PRs, main pushes, and manual dispatch with `contents: read`, no
persisted checkout credential, no publishing credentials, and no cache shared
with signing jobs. No `pull_request_target` or secret-bearing PR job exists.
Each OS verifies formatting and module integrity, runs race-enabled tests with
Fyne's display-independent `ci` driver, vets both the test and native build,
checks the default dependency graph contains Fyne's GLFW native driver, builds
the untagged native GUI binary, and runs safe CLI help/version smokes on that
same binary. It then packages the native executable and runs Linux/Windows
install/uninstall smokes on isolated runner user paths. GUI window interaction and
authenticated flows are not established by those smokes.

Only package files and checksums are retained as CI artifacts for seven days,
including the Windows `*-setup.exe`. Before either unsigned or signed upload,
`check-artifact-retention.sh` verifies every checksum entry exists and is selected
by that workflow's literal upload paths. A ZIP alone cannot hide a missing setup
artifact. This check covers selection closure, not a hosted upload/download receipt.
There is no GitHub Release, tag, package-registry push, or public hosting action.

Action pins were verified on 2026-10-25 with `git ls-remote` against the official
GitHub repositories, then their pinned action manifests were inspected:

| Official repository | Tag checked | Commit used |
| --- | --- | --- |
| [actions/checkout](https://github.com/actions/checkout) | v6 | `d23441a48e516b6c34aea4fa41551a30e30af803` |
| [actions/setup-go](https://github.com/actions/setup-go) | v6 | `924ae3a1cded613372ab5595356fb5720e22ba16` |
| [actions/upload-artifact](https://github.com/actions/upload-artifact) | v6 | `b7c566a772e6b6bfb58ed0dc250532a479d7789f` |

For each pin update, repeat, for example,
`git ls-remote https://github.com/actions/checkout.git refs/tags/v6`, inspect the
new pinned source, and change the full SHA intentionally. See GitHub's
[secure use guidance](https://docs.github.com/en/actions/reference/security/secure-use).

## Signing and notarization: separately authorized operations

`MODE=local` never requests signing material or contacts notarization services.
`MODE=release` on macOS/Windows requires all signing inputs and fails on missing
credentials, signing failure, timestamp failure, or verification failure. No
self-signed fallback or ad-hoc release signature is created. Go's existing local
Mach-O linker signature is not a Developer ID signature or release qualification.

Supply credentials through a private secret manager/environment, never command
history, committed files, workflow inputs, or chat. Do not enable shell tracing,
PowerShell transcripts, or debug logging in a signing environment. Native signing
utilities accept some passwords as process arguments; only use isolated trusted
runners/users. Cancellation traps/finally blocks clean up ordinary failures and
interruptions; a force kill or host loss still requires ephemeral-runner disposal.

| Platform | Required private configuration |
| --- | --- |
| macOS | `MACOS_CERTIFICATE_P12_BASE64`, `MACOS_CERTIFICATE_PASSWORD`, `MACOS_SIGN_IDENTITY` (Developer ID Application identity), `APPLE_TEAM_ID`, `APPLE_API_KEY_BASE64` (App Store Connect team API key), `APPLE_API_KEY_ID`, `APPLE_API_ISSUER`. |
| Windows | `WINDOWS_CERTIFICATE_PFX_BASE64`, `WINDOWS_CERTIFICATE_PASSWORD`, `WINDOWS_TIMESTAMP_URL` (HTTPS RFC3161 service approved by the certificate issuer/operator). The certificate must be a valid code-signing certificate accepted by Windows policy. |

macOS uses a random-password temporary keychain, a nonextractable imported key,
codesign access rather than broad application access, and an explicit keychain
argument. Apple's `create-keychain` transiently adds the keychain to the user's
search list. Before creating it, the tool captures the ordered user-domain list;
cleanup deletes the temporary keychain, explicitly restores that captured list,
and compares a fresh readback byte-for-byte. Deletion, restoration or readback
failure propagates as a packaging failure. The default keychain is not explicitly
reassigned and no system-domain list is written. This operation requires an
isolated signing user/runner with no concurrent keychain writers: restoring a
snapshot can displace another process's intervening search-list changes. It does
not promise that the search list is unchanged while signing is in progress.
It signs the executable and app with hardened runtime and secure timestamps,
checks the team ID and strict signatures, submits an app ZIP through `notarytool`,
requires `Accepted`, and staples/validates the app ticket. It then builds, signs,
notarizes and staples the DMG. The final ZIP is created from the stapled app.
Notarization waits up to 30 minutes per submission. Pending work resumes from
durable original UUIDs and signed bytes without resubmission. Temporary
certificate, API key and keychain are removed on exit; cleanup failure is an
error. Local synthetic qualification does not authenticate Apple credentials.
Unknown-upload identity must be reconciled against history and diagnostic hashes;
ambiguity stops automatic continuation. Submission IDs are retained in safe
receipts and printed in CI logs.

Windows restricts each temporary PFX directory to the current identity, signs the
application, uses NSIS's uninstaller-finalization hook to sign the embedded
uninstaller, and signs the completed installer. Each signature uses SHA256 and an
RFC3161 SHA256 timestamp; SignTool verification plus Authenticode status and
timestamp presence are mandatory. Temporary PFX files are removed in `finally`;
no persistent certificate-store import occurs. A software PFX workflow does not
provide hardware/HSM signing support. If the issuer requires a hardware/cloud
key, this workflow needs a separately reviewed adapter; no assurance bypass is
provided. See Microsoft's [SignTool reference](https://learn.microsoft.com/en-us/windows/win32/seccrypto/signtool)
and the [NSIS uninstaller-finalization reference](https://nsis.sourceforge.io/Docs/Chapter5.html#uninstfinalize).
Apple's [security tool source](https://github.com/apple-oss-distributions/Security/blob/main/SecurityTool/macOS/keychain_list.c)
and installed `security help list-keychains` describe user-domain list mutation.
The lifecycle regressions use a file-only mock and do not establish real Apple
keychain restoration; verification remains held for an authorized isolated signer.
Apple flags were checked using installed `security`, `codesign`, `notarytool`,
`stapler`, `plutil`, and `hdiutil` help; see also Apple's
[notarization workflow](https://developer.apple.com/documentation/security/customizing-the-notarization-workflow).

macOS release mode requires a verified prebuilt `BINARY`, `RELEASE_VERSION`,
`NOTARY_SOURCE_SHA`, `NOTARY_ARCH`, and an absolute persistent
`NOTARY_EVIDENCE_DIR`, in addition to the private signing inputs. Use the protected
workflow for the complete artifact handoff and acceptance checks; do not rebuild
with credentials loaded. The helper reserves separate app and DMG evidence
directories, retains descriptive signed upload archives, binds their hashes to
Apple UUIDs, and fetches matching diagnostic logs even on success. It never
automatically resubmits a reserved stage. See
[submission evidence and recovery](RELEASE-FLOW.md#apple-stages-and-automatic-continuation).

After separate authorization, Windows release packaging uses
`./scripts/package-windows.ps1 -Mode release -Version <numeric-version>` with the
private environment already loaded.

## Manual protected workflow setup

The owner authorized removal of the recurring `release-signing` reviewer for
unattended releases. Retain its `main` deployment restriction, disabled
administrator bypass and isolated environment secrets, together with protected
main's required checks. YAML cannot create or prove those server-side protections;
verify their live API settings separately.

`signed-packages.yml` has only `workflow_dispatch`. It requires explicit
`authorize_signing: true`, a numeric version, `refs/heads/main`, and a protected
ref. It checks out the triggering immutable SHA and uses no PR-provided ref or
release-publishing token. Secrets exist only on their platform's packaging step;
there is no build-cache restore/save in this job. Its only upload is signed
artifacts/checksums retained for review. Dispatch authorizes signing and Apple
submission when the owner deliberately runs it; **actual release publication
requires separate user authorization** and is not implemented here.

## Local qualification and remaining boundaries

Focused checks for this implementation:

```sh
shellcheck -x scripts/*.sh build/linux/*.sh
for script in scripts/*.sh build/linux/*.sh; do bash -n "$script"; done
actionlint .github/workflows/*.yml
plutil -lint build/macos/Info.plist
NATIVE_BINARY=bin/takedown HEADLESS_BINARY=bin/takedown-cli bash scripts/test-packaging.sh
# R1 repairs: synthetic payloads only; no app binary, keyring or signing execution.
bash scripts/test-release-r1.sh
```

The offline checks cover malformed versions, release without secrets, unsupported
Linux release trust, existing output rejection, actual headless-binary rejection,
desktop installation with spaces in the prefix, launcher path generation,
uninstall behavior, and retained unrelated user state. They use a harmless
installer payload and make no real-signing or application-runtime claim.

The earlier, pre-R1 local macOS arm64 qualification used the pre-existing native binary, with heavy
packaging serialized through the shared agent queue. It built an app, ZIP, DMG,
and checksums; it verified plist structure, archive contents, DMG integrity, bundle
permissions, and bundled CLI help/version without account access. Binary identity
was checked against the input. These native operations were not repeated during
the R1 repair task. That evidence does not establish GUI interaction, a fresh
integrated source build, or a signed/notarized artifact.

Not run here: hosted Actions; Windows/Linux native builds, installer execution on
those operating systems; real Developer ID/Authenticode signing; notarization or
Gatekeeper/SmartScreen acceptance; clean-machine GUI installation; OAuth flows;
public policy hosting; release publication; deployment. PowerShell syntax and
NSIS compilation run in Windows CI and remain unqualified until that run passes.
The integrating parent owns fresh combined builds/tests and final source review.

## R1 repair handoff

All original R1 witnesses and both review report files are preserved. The repair
task edits only the previously authorized release/tooling scope. No child agents,
dependencies, commits, real application installations, live keyring operations,
real signing, CI dispatches, publication, or account operations are part of this
repair round. The independent fresh R2 review remains with the parent.

| Finding | Repair and regression evidence |
| --- | --- |
| R1-01 | Linux install rejects symlinked ancestors/payloads and existing unowned conflicts, publishes files without overwrite, records ownership, and uninstalls only validated unchanged owned files. New offline fixtures preserve the unrelated sentinel, pre-existing executable, unowned notes and modified/unknown docs. Invalid/escaping ownership records and symlink substitutions are refused before deletion. |
| R1-02 | Unsigned upload retains `*-setup.exe`; both workflows check checksum-entry retention before upload. Offline fixtures pass with current real patterns and reject the old setup-omitting pattern and a missing setup file. |
| R1-03 | Exclusive output reservation precedes work. A barrier-controlled pair of helper processes confirms the losing producer fails and cannot delete/merge into the winner. Direct publication without reservation is refused. |
| R1-04 | Failed executable deletion preserves recovery metadata. `test-windows-installer.ps1` is restricted to a disposable Windows Actions runner and is wired into native CI. It holds a non-delete-sharing file handle, runs a copied uninstaller with `_?=` to observe the actual worker's exit 10, checks file hashes and registration, releases the handle, and verifies a successful retry's exit/postconditions. NSIS/PowerShell execution remains unqualified locally. |
| R1-05 | Documentation now states the transient user search-list mutation and isolated-runner requirement. The lifecycle captures/restores/verifies the ordered user list. File-only mocks cover normal and empty lists, failed deletion/restoration, ineffective successful restoration, and invalid snapshots. No real keychain state was read or changed. |

R1 offline qualification uses the installed Bash 3.2 and Bash 5 interpreters,
ShellCheck, Actionlint, Bash syntax checks, and plist validation. Tests use fresh
disposable prefixes and mock state; the old failed witnesses remain unchanged.
Source and witness SHA256 inventories and validation output are handed to the
parent separately. These are bounded author checks, not an independent CLEAN
review or platform/release qualification.
