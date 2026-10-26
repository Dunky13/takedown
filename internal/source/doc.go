// Package source discovers Google Takeout archives through the official Drive
// and Data Portability SDKs. Constructing a Manager never logs in or starts an
// export. Login and mutations are explicit operator actions.
//
// Drive: BrowseFolders(ctx, "") browses only real direct children of My Drive
// root. Pass a returned ID to browse its children; equal folder names have
// distinct IDs. Browsing does not change the configured archive subtree. List
// discovers archives only inside that subtree. Each immutable Drive file
// revision is a separate one-part export, so different exports may have equal
// filenames. File identity, safe filenames, custody and traversal caps remain
// enforced. Local output namespacing is owned by the downloader, not this API.
//
// Data Portability: Resources returns the supported explicit resource catalog.
// CheckAccess reads which selected resources have one-time or time-based grants;
// it is not evidence that a one-time grant is unused. CreateExport is explicit,
// records intent before POST, then records Google's actual job ID. No mutation
// POST is automatically replayed, including after a lost response or restart.
// Jobs is a local read; RefreshJobs reads at most 100 due jobs. Status attempts
// persist a five-minute minimum cadence and provider backoff before requests.
// Signed links stay only in memory, last six hours, and data is conservatively
// treated as expired fourteen days after the recorded request. Time-based
// creation retains the twenty-four-hour resource cadence. Auth, quota, backoff,
// transport and local journal failures are never treated as job-local failures.
// IsJobFailure allows an unattended watcher to retain a specific job's hold
// while checking again later for unrelated jobs that become ready.
//
// RetryJob accepts a recorded FAILED job without a successor, with a maximum
// of three attempts across its lineage. It returns a new actual Google ID.
// CancelJob accepts only an in-progress, unexpired time-based job. Neither
// action runs automatically during discovery, refresh, or a watch loop.
//
// Recovery: read Jobs().Intent first. For a create or retry with a lost ID,
// ReconcileIntent accepts an evidence-backed ID from the original response or
// other trustworthy operator evidence. The caller must establish that the ID
// belongs to the exact pending request/account: Google's status endpoint cannot
// prove its resources, creation time or original grant, and Google provides no
// job enumeration API. The original durable intent supplies time and resources.
// Legacy unknown access types remain conservatively one-time. Lost cancellation
// responses reconcile using the original ID, only after terminal state evidence.
// ResolveExpiredIntent explicitly retires only create/retry attempts older than
// fourteen days plus five minutes, historically recorded time-based access and
// exclusively current time-based scope evidence, retaining
// tombstones. It never clears known one-time holds. It cannot resolve an unknown
// reset; login, a checkbox, file deletion and HTTP 401 are not proof of reset.
//
// ArchiveJob hides terminal history from discovery and frees a slot in the
// 1000-active-job limit. It preserves local files, immutable metadata, grant
// guards, cadence, retry budgets and unresolved-intent receipts. The journal
// additionally caps total retained jobs at 10000, expired-intent receipts at
// 1000 and its encoded size at 4 MiB; reaching a safety cap fails closed.
//
// ResetAuthorization is destructive to remote authorization and archive access.
// The caller must first obtain explicit consent to those consequences and pass
// WithResetConfirmation(ctx). The source independently checks every known job
// in the current generation, including archived jobs: all must be COMPLETE, with
// exact immutable part identities, full offsets and durable SHA256 download
// completion receipts. Any pending mutation or expired unknown-intent tombstone
// blocks reset. Only confirmed remote success advances the grant generation;
// old receipts and local files remain. Fresh consent is checked before another
// export. Unknown reset outcomes remain held indefinitely, without blind retry.
//
// Official method and policy references:
// https://developers.google.com/data-portability/user-guide/methods
// https://developers.google.com/data-portability/reference/rest/v1/authorization/reset
package source
