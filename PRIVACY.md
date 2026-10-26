# Privacy policy

Last updated: 25 October 2026. This policy describes the local takedown application
in this repository. A distributor must identify its operator and publish this
policy at a public URL in the OAuth consent configuration before public release.

## Access and purpose

The Drive source requests read-only permission to view and download Drive files.
Although Google grants broad Drive access, takedown scans only the folder you
select and its descendants to find ZIP/TGZ archives. The Data Portability source
requests only the resource groups you select and starts export jobs only when
you choose Create export or run the export command. These sources use separate
OAuth grants. The app uses the data to download and verify archives at your
chosen destination. It does not extract archive contents automatically.

Google handles passwords, passkeys, MFA and consent in your system browser.
takedown does not collect Google passwords, TOTP secrets, or browser cookies.
The optional email hint selects a local authorization profile; Google controls
which account you authorize in its consent screen.

## Storage, sharing and retention

OAuth access and refresh tokens are stored in the OS keyring.
Destination passwords are also stored in the OS keyring only when you explicitly
choose Save destination authorization or run `storage auth`. Otherwise a supplied
destination password remains a session/environment credential. Passwords are not
embedded in login-task definitions, preferences or job/download manifests.
Client configuration is read from the Desktop app JSON file you provide.
Preferences contain account
hints, folder IDs, selected resource groups, destination paths and transfer
settings. Local state contains archive metadata, byte offsets, completion records,
and Data Portability job receipts, including mutation holds and retained grant,
cadence and retry history. Archiving a receipt hides history without deleting its
safety evidence. It excludes OAuth tokens and signed download
URLs. Your configuration, client JSON and job state remain private local files.

Archives and partial downloads are written to the destination you choose. A
network destination sends those archive bytes to that server using its configured
protocol and credentials. Use trusted storage with appropriate encryption and
access controls for your data. The app does not automatically encrypt archive
contents. Local/OS-mounted destinations inherit the operating system's protections;
SMB signing is not a promise of transport encryption. Use SFTP or a separately
secured mount when encrypted network transport is required. Compliance with any
restricted-scope security requirements must be established for the actual storage
and deployment before use or distribution.

There is no developer-operated backend, telemetry, advertising, sale of data,
model training, or human review of your archives. Requests go to Google APIs,
Google archive storage, and the network destination you explicitly configure.
Google and that destination's operator retain their own responsibilities and
policies. Downloaded data and local metadata are kept until you remove them;
Google controls its own export retention periods.

## Control and deletion

Stop or cancel a job to retain resumable partial files. Delete archives and their
`.part` files directly at the chosen destination when you no longer need them.
Remove the local application configuration/state directory to erase preferences
and transfer/job history after all running jobs are stopped. Preserve a job
receipt with an unresolved creation intent until it has been reconciled, to avoid
starting a duplicate export by losing its history.

Use `takedown auth --forget` with the same client/profile settings to remove its
locally stored OAuth tokens and any supported legacy secrets for that profile.
Other profiles must be forgotten separately. This action does not revoke Google's
remote grant. Revoke remote access at [Google Account permissions](https://myaccount.google.com/permissions).
Saved destination passwords are separate: use `takedown storage forget` or the
GUI's destination-authorization removal control for each saved destination and
data-directory profile. Remove these keyring entries before deleting preferences
needed to identify their slots. Removing a configuration directory alone does
not erase OS keyring entries. Remove the user-login watch before uninstalling
or moving its executable.

Data Portability's explicit authorization-reset control also revokes its Google
grant and allows Google to remove previous remote archives. It requires user
confirmation and exact durable completion evidence for every known archive;
it is not an automatic cleanup or retry action. Retry and cancellation are
separate explicit Google API operations. Local downloaded files are retained.
Browser downloads, exported source data, and copies on remote storage are not
deleted by forgetting authorization or uninstalling the app.

## Google Limited Use

The use and transfer of information received from Google APIs adhere to the
[Google API Services User Data Policy](https://developers.google.com/terms/api-services-user-data-policy),
including its Limited Use requirements. For the Data Portability source, the same
commitment applies to the
[Data Portability API user data and developer policy](https://developers.google.com/data-portability/policy),
including its Limited Use requirements. These commitments do not claim that Google
has reviewed, verified, or approved a particular build or deployment.

For questions about this source implementation, use the repository's
[issue tracker](https://github.com/Dunky13/takedown/issues). Do not post tokens,
credentials, signed URLs, or archive contents in an issue.
