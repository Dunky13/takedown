<img src="assets/icon.svg" width="64" height="64" alt="">

# takedown

A Go application for downloading Google Takeout archives delivered to Drive and
exports created with the Data Portability API, using documented Google APIs.
Destinations include local folders, mounted shares, Windows UNC paths, SMB2/3,
and SFTP. The default binary contains a
Cobra CLI and a native Fyne desktop interface. No arguments opens the desktop app.

## Build and run

Use **Go 1.26.10 or newer**. Current SSH/network dependencies require Go 1.26,
and the module selects a patched toolchain.
A desktop build also needs a C compiler and native graphics development libraries.
Google sign-in and consent use your normal system browser. No Chromium automation,
password entry, TOTP entry, cookies, or stealth are used.

```sh
go mod download
make build
./bin/takedown
```

On macOS, install the Xcode command line tools (`xcode-select --install`). On
Windows, use a supported MinGW-w64 C compiler on `PATH`. On Debian/Ubuntu Linux:

```sh
sudo apt-get install build-essential pkg-config libgl1-mesa-dev libegl-dev xorg-dev libwayland-dev libxkbcommon-dev
```

On Fedora, install `gcc`, `pkgconf-pkg-config`, `mesa-libGL-devel`, `mesa-libEGL-devel`,
`libXcursor-devel`, `libXrandr-devel`, `libXinerama-devel`, `libXi-devel`,
`libXxf86vm-devel`, `wayland-devel`, and `libxkbcommon-devel`. On Arch, install
`base-devel`, `pkgconf`, `mesa`, `libglvnd`, `libxcursor`, `libxrandr`, `libxinerama`,
`libxi`, `libxxf86vm`, `wayland`, and `libxkbcommon`. These supply OpenGL/EGL and
both X11 and Wayland development files. The pinned GLFW module enables both Linux
backends by default. Fyne does not use WebKit2GTK. See
[GLFW's X11 and Wayland requirements](https://www.glfw.org/docs/3.4/compile.html#compile_deps)
and [Fyne's build requirements](https://docs.fyne.io/started/).

For a server without native graphics libraries, `make headless` produces the
optional `bin/takedown-cli` binary. The default build supports both interfaces.

## Set up documented access

1. In Google Cloud, create a **Desktop app** OAuth client, configure the consent
   screen and enable the Google Drive API. Download the client JSON file into
   a private local file, such as `./oauth-client.json`. Do not commit it.
2. Create an export in Google Takeout with **Add to Drive**. You can select
   Google's recurring export option. `takedown export` opens the documented
   Takeout setup URL in your browser; you complete its sign-in and export choices.
3. Copy the Drive folder link containing your Takeout archives. Authorize
   takedown once, then watch for and download ready archives:

```sh
./bin/takedown auth --oauth-client ./oauth-client.json --drive-folder 'https://drive.google.com/drive/folders/FOLDER_ID'
./bin/takedown watch --destination ./archives
```

Replace `FOLDER_ID` with the actual folder ID from Drive. A successful `auth`
saves those preferences. **Drive access is read-only but broad:** Google's
`drive.readonly` scope permits viewing and downloading all Drive files. The app
limits discovery to your selected folder and its descendants. The narrower
`drive.file` scope does not automatically grant access to all Google-created
archives and future files in a selected folder. Public distribution using the
restricted scope requires Google's applicable verification process; source code
and local tests do not establish approval. See [Drive scope requirements](https://developers.google.com/workspace/drive/api/guides/api-specific-auth).

The OS keyring stores OAuth tokens; a functioning, unlocked Windows Credential
Manager, macOS Keychain, or Linux Secret Service is required. Token refresh is
automatic while Google's grant remains valid. A revoked or expired grant stops
background work and asks you to authorize again. External OAuth apps left in
Testing can receive refresh tokens that expire after seven days; publication
status and Google's verification rules remain operator setup, not an app bypass.
See [OAuth token expiration](https://developers.google.com/identity/protocols/oauth2#expiration).

```sh
./bin/takedown export                  # Opens normal Takeout setup for Drive delivery
./bin/takedown list --json
./bin/takedown download --destination ./archives
./bin/takedown download --export EXPORT_ID --destination ./archives --part-by-part
./bin/takedown watch --poll-interval 5m --destination ./archives
./bin/takedown status --json
./bin/takedown auth --forget           # Deletes this profile's local authorization
./bin/takedown gui
```

`auth --forget` removes local authorization, not Google's remote grant. Revoke
remote access in [Google Account permissions](https://myaccount.google.com/permissions).
The optional `--email` is a browser login hint and local profile slot, not a claim
that the selected account has been verified. Use the same hint for later runs.
Help and status need neither Google credentials nor keyring access.

### Data Portability exports

For eligible accounts and supported products, this API can create exports without
operating the Takeout web UI. It has different product coverage and cannot list
or retrieve existing email-link Takeout exports. Configure/enable the Data
Portability API and its required billing and app verification, then select only
the resource groups you need. It uses a separate OAuth grant from Drive.

Data Portability archives require encryption in transit and at rest. Establish
protected storage before downloading: takedown does not encrypt destination
files or verify their storage encryption. SFTP encrypts transit only; direct
SMB encryption depends on the server and share configuration. An encrypted OS
volume or appropriately protected network destination is an operator
prerequisite, not a protection this app automatically supplies.

```sh
./bin/takedown auth --source portability --oauth-client ./portability-client.json --resources myactivity.search
./bin/takedown export
./bin/takedown watch --destination ./archives
```

The consent screen lets eligible users choose one-time or time-based access.
Export creation is explicit; `list` and `watch` never create new exports. The app
persists creation intent before sending a request, prevents blind duplication
when Google's response is lost, and observes the five-minute polling cadence.
For time-based grants, separate export invocations must respect the 24-hour limit.
An ambiguous creation outcome requires reconciliation rather than an automatic
second POST. See [Google's API workflow](https://developers.google.com/data-portability/user-guide/introduction).

The GUI includes a resource catalog, explicit Drive folder navigation, grant
inspection, and a Jobs view. The same job controls are available in the CLI:

```sh
./bin/takedown jobs list --source portability --json
./bin/takedown jobs list --source portability --refresh
./bin/takedown jobs access --source portability
./bin/takedown jobs retry JOB_ID --source portability --confirm
./bin/takedown jobs cancel JOB_ID --source portability --confirm
./bin/takedown jobs archive JOB_ID --source portability --confirm
```

Retry is limited to three attempts per lineage. Cancellation requires an
in-progress time-based job. Archiving hides terminal history while retaining
grant, cadence, retry, and download receipts. Failed or unavailable jobs do not
prevent unrelated ready jobs from downloading; watch retains their hold.

For an uncertain request, `jobs reconcile JOB_ID` requires independent evidence
that the ID belongs to the exact pending request, account, resources and time.
API accessibility alone cannot prove that association. Do not guess an ID.
`jobs resolve-expired` only releases eligible historical time-based requests
after fourteen days plus five minutes and matching current grant evidence.
Unknown or one-time historical grants stay held without a known ID.

`jobs reset-authorization` requires explicit confirmation and exact durable
download completion for every known job. It revokes Google's authorization and
allows old remote archives to be removed. It cannot clear an uncertain reset;
that outcome remains held. Reauthorizing does not erase the safety history.
See [Google's reset contract](https://developers.google.com/data-portability/reference/rest/v1/authorization/reset).

### Background behavior and preferences

`download` selects all ready archives unless `--export` narrows the selection.
`watch` stays active until canceled, discovers new ready archives, avoids repeating
completed downloads, and retries temporary service/network failures with backoff
that respects server retry instructions. It never opens a browser automatically.
Default polling is five minutes; the minimum is one minute for Drive and five
minutes for Data Portability. Drive source identities include file revisions;
changed remote contents cannot silently resume into an old archive.

Default download concurrency is two. `--workers` accepts 1 through 16, and
`--safety-buffer` sets retained free bytes (default 2,000,000,000). If a whole batch
does not fit, interactive `download` users can choose part-by-part mode.
Noninteractive jobs require an explicit `--part-by-part` choice. Ctrl+C cancels
work while retaining partial files and durable checkpoints.

Preferences and resumable state live in the platform's user configuration
directory under `takedown`. Override with `--data-dir ./private-data`. The GUI,
successful `auth`, and export setup save preferences; other flags apply to that
invocation. Changing account/source/client/folder/permissions invalidates old
archive selections. Status reads the atomic manifest without taking a writer lock.

Diagnostics are available in the GUI and with `takedown doctor`. By default they
inspect local setup; `--online` explicitly probes Google and `--storage` connects
to the destination and checks free space. A storage probe may create its
directory. Diagnostics do not verify actual destination encryption.

### Start watching at user login

After saving preferences and authorizing explicitly, preview a startup definition:

```sh
./bin/takedown service install --dry-run
./bin/takedown service install
./bin/takedown service status
./bin/takedown service remove
```

Installation explicitly starts the watch and registers it for the current user's
login session: launchd on macOS, systemd user services on Linux, or an interactive
Scheduled Task on Windows. It requires the saved OAuth refresh grant, an unlocked
keyring, a stable binary path, and any network destination authorization. It does
not start sign-in or create exports. Use `--part-by-part` at installation if that
space policy is desired. Remove the watch before moving or uninstalling the
binary. Installing an already loaded watcher keeps its current run; remove and
reinstall to apply changed preferences immediately. Native startup registration
and reboot behavior require platform testing.

## Network destinations

```sh
./bin/takedown download --destination 'smb://operator@nas.example/backups/takeout'
./bin/takedown download --destination 'sftp://operator@nas.example/exports/takeout' --ssh-key ./keys/archive --known-hosts ./keys/known_hosts
```

On Windows, an already accessible share can be supplied directly:

```powershell
.\bin\takedown.exe download --destination '\\server\share\takeout'
```

For direct SMB or password-based SFTP, supply `TAKEDOWN_STORAGE_PASSWORD` through
your environment or use the GUI's masked network-password field. For unattended
use, explicitly save it with `takedown storage auth --destination 'smb://operator@nas.example/backups/takeout'`
or Diagnostics > Destination access in the GUI. `takedown storage forget` removes
the selected destination's saved password. These use the OS keyring. Do not embed
passwords in destination URLs. Remote passwords are not written into preferences
or the state manifest. SFTP requires a trusted host key in the configured
`known_hosts` file. Establish that trust independently before downloading.

The SFTP `statvfs` extension and SMB filesystem information provide available
space. Unsupported SFTP free-space reporting fails closed by default.
`--allow-unknown-space` explicitly permits a warning-only policy when the server
cannot report free space. Under that policy, disk capacity cannot be guaranteed.
Network storage must support seeking, regular files, and rename operations.
SFTP downloads also require `fsync@openssh.com` for durable checkpoints. A
server that cannot sync archive bytes fails explicitly before a newer checkpoint
is published. Canceling remote I/O may leave the last durable offset behind the
file length; the next run reconciles and verifies the retained prefix.

## Resume, space, and verification

New archives are written under `exports/<export-identity-hash>/filename.part`.
Different exports can share the original filename without overwriting each
other. Existing checkpointed flat paths continue in place. Both interfaces use the same job
engine and atomic JSON manifest. Process locking prevents simultaneous writers
from using the same manifest. Signed URLs and OAuth tokens are excluded from that manifest.

Before a batch starts, remaining archive bytes plus the safety buffer must fit.
Active workers reserve their combined remaining bytes. Space is checked before
each part and at most every 50,000,000 bytes written per worker. Low space pauses
the job, flushes partial files, and reports a warning. Free space and resume, or
stop the job and choose a different destination. Changing destinations starts
there; it does not silently move previously downloaded bytes.

HTTP resume validates the requested byte range and the expected archive size.
When the server supplies no strong ETag, the engine compares the retained prefix
against the same complete response before appending. That preserves identity
safety but re-reads the prefix over the network.
OAuth access tokens refresh automatically. Expired API archive URLs are renewed
through the source API while preserving partial files. Permissions, revoked
consent, quotas, and archive integrity failures are not treated as reasons to
bypass a Google check. Long-running watches retry classified transient failures;
other failures stop with an actionable error.

Archives are verified before their `.part` suffix is removed. ZIP structure and
entry checksums and complete gzip streams are checked. Drive MD5 checksums, when
provided, are verified over the archive bytes as well. Verification of remote
archives reads them through the storage connection. ZIP/gzip verification detects
corruption; it does not establish a cryptographic Google provenance guarantee.
Successful completion records store a SHA-256 digest of the compressed archive.
Restarts recheck that digest before accepting an existing file, including a
complete `.part`; a different valid archive of the same size is rejected. Older
checkpoints without a digest or upstream checksum must prove their contents
against the source before they become trusted completions. Version 2 manifests
and version 3 manifests are read and upgraded to version 4 on the next durable write.
Archive files are never extracted automatically.

## Cross-platform builds

Fyne uses CGo. Setting `GOOS` alone cannot cross-compile the desktop binary:
you also need a C compiler, SDK, and graphics libraries for the target.

```sh
make windows WINDOWS_CC=x86_64-w64-mingw32-gcc
make linux LINUX_CC=x86_64-linux-gnu-gcc
make macos MACOS_CC=clang
```

The macOS command assumes a macOS host and valid SDK; on another host, supply
an appropriate SDK/toolchain compiler. The targets produce Windows amd64,
Linux amd64, and macOS arm64 binaries. For another architecture, build directly:

```sh
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -p 2 -trimpath -o bin/takedown-macos-amd64 ./cmd/takedown
```

Run that example on macOS with its native C compiler. For repeatable
containerized cross-compilation, see [fyne-cross](https://github.com/fyne-io/fyne-cross).
Build and runtime acceptance on each target remain separate checks. Unsigned
binaries may require the platform's normal local application trust approval.

Native CI and package tooling produce macOS app/ZIP/DMG, Windows ZIP/installer,
and Linux tar packages. See [CI builds and downloads](docs/CI-BUILDS.md) for
the native target matrix and how to retrieve unsigned test packages.
Signing is explicit and fails closed when its required
credentials or verification are missing. See the [release guide](docs/RELEASING.md)
and its [offline HTML version](docs/RELEASING.html). Hosted CI, native target
runtime acceptance, signing and notarization are separate qualification steps.

## Development and acceptance

```sh
make check
make test
make build
make headless
```

Tests use HTTP servers, temporary storage, a fake credential store, and Fyne's
software test driver. They do not sign into Google or modify the real OS keyring.
The `ci` build tag is for Fyne's test driver; do not distribute a `ci` build as
the desktop application.

The application no longer scrapes Takeout pages or captures Google browser
sessions. It uses documented OAuth, Drive, and Data Portability interfaces.
OAuth client setup, consent, account eligibility, applicable Google verification,
and live provider acceptance remain separate from local implementation tests.
Live Google, SMB, and cross-platform GUI acceptance are operator checks.

See the [automation and policy handoff](docs/AUTOMATION.md), its
[self-contained HTML](docs/AUTOMATION.html), and the [privacy policy](PRIVACY.md).
The [original terms research](docs/TOS-RESEARCH.md) describes the earlier browser
implementation and is retained as historical evidence.
