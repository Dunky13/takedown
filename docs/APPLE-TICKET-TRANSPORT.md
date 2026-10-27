# Temporary Apple ticket transport preference

The automatic macOS release signing and fresh credential-free acceptance jobs
use `scripts/macos-ticket-transport.sh` to temporarily reject outgoing UDP/443.
This lets native Apple ticket clients attempt HTTPS over TCP. Ordinary TCP/TLS,
certificate validation, Gatekeeper, quarantine, notarization topology and all
existing custody and acceptance gates keep their existing behavior.

The motivation is the runner-image contributor's reproduced HTTP/3 connection
loss and successful TCP workaround in
[actions/runner-images issue 14924](https://github.com/actions/runner-images/issues/14924#issuecomment-6188826956).
The incident identifies the Intel macOS 15 image. It does not establish the
cause of Takedown's ARM failure. Applying the same temporary transport preference
to both release architectures is an implementation choice, not ARM attribution.
A passing mock suite does not prove native PF compatibility, Apple acceptance,
or fresh hosted macOS acceptance.

## Supported environment and PF contract

Production execution requires GitHub Actions, a GitHub-hosted macOS 15 image,
protected main in `Dunky13/takedown`, the release workflow, and a signing or
acceptance job. The workflow supplies the actual `runner.environment` context
to each helper step. The helper uses fixed system `sudo` and `pfctl` paths.
Normal local and self-hosted execution fails before accessing PF. Environment
checks prevent accidental local use; they are not an authentication boundary
against someone deliberately forging all CI variables on their own machine.

Setup runs after pinned automation/product qualification and before Apple
credential exposure or native ticket operations. It refuses any of the seven
Apple signing secret variables. Fresh acceptance receives no Apple secrets.
Tool discovery and help/version checks remain credential-free.

The helper requires the existing root filter rules to contain exactly one
unconditioned `anchor "com.apple/*" all`, with only the standard optional
`scrub-anchor` rule alongside it. It never inserts a root attachment or reloads
the root ruleset. A uniquely named direct child of `com.apple` carries one
`block return out quick proto udp ... port 443` rule. Existing child anchors must
sort after the owned child, so an earlier quick pass cannot skip the guard.
Unsupported root layouts, collisions, status/reference combinations or readback
formats fail setup without a global-rules fallback.

The child name contains the actual workflow run, attempt and job discriminator,
plus 96 random bits, and is shorter than 64 characters. Private custody binds it
to repository, actual workflow run/attempt/job, workflow SHA, and hashes of the
original root filter rules, child inventory and enable-reference identities.
The display timestamp/age column is excluded; PID, process name and exact
decimal token remain unchanged. Token comparisons use strings to avoid uint64
rounding in `awk`. A private manifest digest detects accidental custody changes.
Original artifact run/source bindings are separately preserved by the existing
release scripts; the wrapper passes their original run variables only to its
child command.
The guarded child also receives separate actual workflow run/attempt values.
Only cleanup consumes those values, preserving PF custody when a native verifier
temporarily uses the original artifact build's run/attempt for release evidence.

PF is enabled with `-E`, retaining its returned reference token privately. The
helper requires exactly one matching new reference and unchanged original
reference identities. An already-enabled filter without a pre-existing reference
is refused: releasing a new reference cannot safely restore legacy `-e`
enablement. No global `-e`, `-d`, `-f`, `-F`, state killing, proxy configuration,
trust adjustment or persistent configuration-file edit is used.

The installed macOS `pfctl(8)` manual documents anchor-scoped `-a` loads and
flushes, `-E`, `-X token`, and `-s References`; `pf.conf(5)` documents wildcard
direct-child evaluation and `quick` propagation to the root. The primary
[OpenBSD anchor documentation](https://www.openbsd.org/faq/pf/anchors.html)
also describes scoped rule loads/flushes. macOS-specific reference operations
must follow the installed macOS manual, not an OpenBSD flag assumption.
Apple's [XNU PF implementation](https://github.com/apple-oss-distributions/xnu/blob/main/bsd/net/pf_ioctl.c)
also documents token creation and reference-counted release in source.

## Lifecycle and receipt gate

`setup` privately snapshots and verifies supported state, journals each possible
mutation, loads/readbacks the child, acquires the reference, and verifies the
active guard. Setup failures and handled signals attempt owned cleanup.

`run COMMAND ...` verifies the active binding and launches the release command
in a separately verified owned process group. On failure, SIGINT or SIGTERM it
records cancellation, terminates every wrapper/native descendant and awaits
their cleanup before PF teardown. Process readback ignores inert zombies;
unknown or still-running process state withholds teardown. A bounded graceful
drain can force-kill remaining owned processes, which fails the job. It never
signals the caller's process group or unrelated processes. Native acceptance
checks the guard's cancellation marker before and after checkpointing and before
success; a cancelled reservation cannot write a late verified proof. Existing
child scripts keep their own credential and mount cleanup. A workflow `always()` step repeats cleanup/readback
after an attempted setup, including handled cancellation. Abrupt runner loss or
SIGKILL cannot execute traps; the jobs are restricted to disposable hosted VMs.

Cleanup flushes only the owned child rules and releases only the acquired token
with `-X`. It never restores a saved root ruleset over another writer's changes.
It requires empty owned rules, unchanged root and sibling inventories, original
reference identities, and the original enabled state before recording completion.
A replay after reference release can recognize exact original references/status
without releasing the token twice. Repeated completed cleanup performs readback
without further mutation. Missing, corrupt, symlinked, public or differently
bound custody fails closed; no token is inferred from a reference-list delta.
Concurrent external PF changes cause a failed preservation check, not an attempt
to overwrite or release another owner's state.

Cleanup failure fails the operation/job even when native verification succeeded.
The native verifier calls the same cleanup helper before it writes success or a
verified checkpoint. The outer guard and workflow then repeat its idempotent
preservation readback. Failure before success is recorded leaves acceptance
blocked rather than retaining a premature verified checkpoint.
The guarded fresh acceptance command moves its new receipt to `.json.pending`
before transport cleanup. The workflow deletes
pending/final fresh receipts on cleanup failure or failed native acceptance.
Only a successful job with verified cleanup may rename it to the normal receipt
name and upload it. Existing readiness/publication dependencies therefore cannot
retain or use a new fresh receipt from a failed transport lifecycle. Signing
failure receipts remain allowlisted and available for diagnosis; private PF
state, tokens and raw command output are never included in uploads.

Transport setup does not reopen terminal acceptance, replenish an attempt budget,
change counters, import superseded evidence, resubmit a UUID, or authorize a new
qualification. A blocked receipt remains blocked even if transport setup passes.

## Local validation and operator handoff

Run only the synthetic suite locally:

```sh
bash scripts/test-macos-ticket-transport.sh
```

It creates an owned temporary model directory with an explicit mock marker and
PF driver, and executes the real helper lifecycle against that model. It also
parses the workflow and executes its cleanup/promotion shell bodies against
isolated fixture files, including deferred receipt bytes and propagation of
native failure. Mocks are refused when `GITHUB_ACTIONS=true`; offline CI explicitly
unsets that variable for this synthetic suite. It never invokes the host's PF,
changes networking or uses real Apple credentials.

Relevant checks include ShellCheck, Bash syntax, Actionlint, and the existing
Apple custody, source qualification, signing and final-verifier shell contracts.
Native CI runs the transport suite on Linux and both macOS architectures.

Root owns source review and delivery, hosted PF/root-anchor compatibility,
success/failure/cancellation cleanup qualification on both runner architectures,
real TCP fallback, native Apple operations, and release acceptance/publication.
No local mock result substitutes for those checks or changes a terminal hold.
