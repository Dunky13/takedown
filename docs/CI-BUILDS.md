# CI builds and app packages

The `Native CI` workflow builds the same Go executable for CLI and native Fyne
GUI use. Pull requests, pushes to `main`, and explicit manual runs execute the
workflow. No Google account, OAuth client file, production keyring, signing
credential, or release publication is needed for these builds.

## Native targets and checks

The target matrix covers Linux amd64 (`ubuntu-24.04`), Windows amd64
(`windows-2025`), macOS Apple Silicon (`macos-15`), and macOS Intel
(`macos-15-intel`). Every target uses Go 1.26.10 and Fyne v2.8.2.
Repository attributes enforce LF line endings for text files on every target,
including Windows checkouts with `core.autocrlf=true`; CI checks formatting
without rewriting source files.
The generated dependency-notice bundle is exempt from newline and whitespace
normalization so quoted upstream legal text retains its original bytes. Its
generator check verifies the complete bundle instead.

Every native job checks formatting and module integrity, runs race-enabled tests
with Fyne's software test driver, vets both the test and native configurations,
and builds the actual Cgo/OpenGL desktop executable. CLI help and version smoke
checks use that desktop executable. The software test driver and optional
headless executable are never desktop package inputs.

The workflow also checks shell and PowerShell tooling, exercises offline package
contracts, and packages the native executable. Linux and Windows run temporary
per-user installation and uninstallation checks; Windows additionally verifies
that a locked executable retains recovery metadata and can be removed on retry.
These checks do not establish interactive GUI behavior on a clean machine.

## Download unsigned packages

Open the [Native CI workflow](https://github.com/Dunky13/takedown/actions/workflows/ci.yml),
select the run for the intended commit, and download its `unsigned-native-*`
artifacts. Artifact names include the platform, architecture, and workflow
revision. Artifacts are kept for seven days and access follows repository
permissions. A pull-request run can build GitHub's temporary merge revision;
check the run and the packaged `BUILD-INFO.txt` rather than assuming the filename
uses the pull-request head revision.

The retained packages are:

- macOS: ZIP containing `takedown.app`, plus a DMG.
- Windows: portable ZIP and a per-user `*-setup.exe` installer.
- Linux: TAR.GZ containing the executable and per-user installation scripts.
- Every target: `SHA256SUMS` covering its final archives and installers.

With GitHub CLI, replace `RUN_ID` with the selected run's numeric ID:

```sh
gh run view RUN_ID --repo Dunky13/takedown
gh run download RUN_ID --repo Dunky13/takedown --dir ci-artifacts
```

Verify each `SHA256SUMS` from its own package directory, using
`shasum -a 256 -c SHA256SUMS` on macOS or `sha256sum -c SHA256SUMS` on Linux.
Checksums detect changed bytes; they do not establish publisher identity.
The package version `0.0.0` denotes a development build, not a release.

## Icons and license files

The GUI and native packages use the original archive/download icon in
[`assets/`](../assets/). The SVG is the editable source; PNG, ICO, ICNS, and the
Windows resource object are generated assets. See the asset instructions there
before regenerating them. The README uses the same PNG as the GUI.
Every native CI target regenerates the assets in memory and compares their bytes
without rewriting the checkout. Generator compilation pins floating-point
operation behavior and selects the upstream generic raster path as well as
pinning the Go version. This prevents fused multiply-add instructions and
architecture-specific SIMD accumulation order from changing rasterized pixels.
These settings apply only to artwork generation; normal application builds use
their usual compiler settings.

Packages carry the project's Apache 2.0 `LICENSE` and `NOTICE`, together with
`THIRD_PARTY_NOTICES.md` for dependencies that retain their own licenses. Changes
to dependencies or bundled assets require checking those notices as well as the
application. These files do not indicate Google app verification or signed
publisher status.
Linux CI also verifies the generated notice bundle against all four pinned
native runtime dependency graphs. The bundle documents regeneration and scope.

## Local equivalent

Run from the repository root with the native compiler, graphics headers, and
Go 1.26.10 installed as described in [the release guide](RELEASING.md):

```sh
make ci-native
make package-macos BINARY=bin/takedown
# On Linux instead:
make package-linux BINARY=bin/takedown
```

On Windows, after the native build, use PowerShell 7:

```powershell
./scripts/package-windows.ps1 -Binary bin/takedown.exe
```

Packaging requires a fresh output directory. Use `OUT_DIR` or
`-OutputDirectory` for a subsequent package operation. Each host produces its
own native architecture; CI supplies the separate target hosts.

## Qualification and deferred signing

A successful CI run qualifies only the checks that actually ran for that exact
revision. Hosted results, local builds, installer checks, interactive GUI
acceptance, live Google access, and production network storage are separate
evidence. Live Google acceptance remains the user's check. Google app
verification and destination encryption requirements still apply.

Paid Windows signing is deferred. CI creates unsigned review packages without
enrolling in a signing service or using credentials. macOS notarization and the
release workflows require separately authorized configuration and runs. The
owner-authorized unattended release graph is described in [release automation](RELEASE-FLOW.md).
Unsigned downloads can trigger normal operating-system trust prompts.
No public release, automatic update, or deployment is performed by Native CI.

Source and workflow review remain required before merging. The `ci-required`
aggregate must pass together with every native matrix job. Check the repository's
current branch rules separately; this guide does not establish their enforcement
or grant approval to merge, sign, or publish a release.
