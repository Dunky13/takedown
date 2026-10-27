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
The inventory supports a consistent set of direct child leaves, or a consistent
set of exactly `com.apple/LEAF` paths. Only leading horizontal display indentation
is removed. Every leaf must match `[A-Za-z0-9_][A-Za-z0-9_.-]*` and fit the PF
64-byte name buffer including its terminator. Empty inventories are supported;
blank rows, duplicates, mixed spellings, other prefixes, trailing whitespace,
control characters (including embedded NUL bytes before `awk` parsing), relative
traversal and deeper paths are refused. Ordering is
checked on the exact child component in either representation. The sorted
nonowned rows retain their original spelling and must compare byte for byte;
changing a sibling from leaf to qualified form is drift, even if it denotes the
same path. Unsupported root layouts, collisions, status/reference combinations
or readback formats fail setup without a global-rules fallback.

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
Before releasing its enable token, it requires empty owned filter rules, NAT
rules, tables and direct-child inventory, with successful fresh captures for each.
Unexpected auxiliary contents cause refusal without flushing those contents.
It then requires unchanged root and nonowned sibling inventories, original
reference identities, and the original enabled state before recording completion.
The owned row must occur exactly once during setup and active checks, as either
the exact leaf or the exact qualified custody path in the inventory's consistent
representation. Completed cleanup permits that same row to disappear or remain
once, provided a private `load-attempted` journal binds it to this helper and all
the empty-namespace and preservation checks pass again. No other row is removed
from comparison. An empty pre-existing collision is still refused.

This is an explicit namespace-lifetime distinction: cleanup removes the helper's
packet-filtering effect and restores original nonowned state; it does not promise
that PF deallocates its exact empty owned namespace. The namespace is retained
only on disposable hosted VMs and is never reused. A replay still refuses rules,
NAT, tables, descendants, duplicates, mixed aliases or foreign inventory drift.
There is no root reload, recursive flush, broad `-F all`, state flush, or attempt
to force kernel reference release to delete a name.
A replay after reference release can recognize exact original references/status
without releasing the token twice. Repeated completed cleanup performs readback
without further mutation. Missing, corrupt, symlinked, public or differently
bound custody fails closed; no token is inferred from a reference-list delta.
Concurrent external PF changes cause a failed preservation check, not an attempt
to overwrite or release another owner's state.

Readback gates emit only fixed phase and predicate booleans, for example
`ticket-transport: phase=setup-load anchors=false`. Phases are `setup-load`,
`active`, `cleanup`, and `post-cleanup`; predicates are `rule`, `root`, `anchors`,
`references`, `status`, `owned-empty`, and `token`. Evaluation still stops at the
first failed predicate and retains the aggregate failure message. A predicate
not printed after a failure was not evaluated. These facts distinguish the
failed gate without printing raw PF output, anchor names, reference identities,
tokens, process names, or network diagnostics. Private readbacks stay private.

Failed rule, root and anchor gates also emit a fixed `shape` record, using only
the successful capture from that predicate. Failed captures report
`available=false`, so partial or stale output cannot be presented as fresh
evidence. Rule/root records contain only line counts and a SHA-256 digest.
Anchor records contain line counts, exact owned-leaf and owned-qualified match
counts, and whether the remaining sorted inventory matches the original bytes.
These classifications do not relax the acceptance comparisons. They distinguish
an owned-name representation or lingering owned child from a foreign inventory
change without disclosing anchor names. An auxiliary namespace capture failing
before the parent inventory is captured reports `available=false`.
No diagnostic performs another PF command or changes cleanup decisions.
Diagnostic writes are best effort and preserve the underlying predicate status,
including when stdout is closed. They cannot skip later successful preservation
checks or the cleanup journal. The pre-existing final success-message write can
still fail the overall command if its output is unavailable.
Each diagnostic write uses a child stdio context: Bash 3.2 can otherwise retain
a failed builtin write and flush its buffered text into a later private-file
redirection. Diagnostic output cannot contaminate subsequent rule comparisons.

### Evidence for the supported inventory and namespace lifetime

Hosted macOS 15 sign jobs `122392707722` (scheduled ARM), `122388659916`
(workflow-run ARM) and `122388659873` (workflow-run Intel), observed on
2026-10-27, reported successful rule/root setup readback followed by inventory
failure: three rows, zero exact owned leaves, one exact qualified owned path,
and remaining rows equal to the original inventory. After scoped cleanup they
reported unchanged root/references/status and empty owned filter rules, followed
by the same qualified-row inventory failure. These safe classifications establish
the observed representation and retained owned row; they do not disclose raw PF
state or prove why that kernel retained it. Apple credentials remained skipped.

The installed macOS `pfctl(8)` manual specifies direct children for
`-a ANCHOR -s Anchors` and scoped `-a` rule operations; it does not promise that
printed children are unqualified leaves or that a rule flush deletes a namespace.
Apple's tagged [ruleset implementation](https://github.com/apple-oss-distributions/xnu/blob/xnu-11417.101.15/bsd/net/pf_ruleset.c#L178-L228)
separates finding/referencing an anchor from its active rules. Its
[empty-ruleset removal conditions](https://github.com/apple-oss-distributions/xnu/blob/xnu-11417.101.15/bsd/net/pf_ruleset.c#L417-L474)
include reference count, children, tables, pending transactions, and active and
inactive rule queues. Consequently, an empty active filter ruleset alone does
not imply namespace removal. The supported scoped flush is retained: neither
the manual nor those conditions provides a safe standalone name-delete operation
that justifies a wider flush or another ruleset reload.

The [Apple evaluator](https://github.com/apple-oss-distributions/xnu/blob/xnu-11417.101.15/bsd/net/pf.c#L3114-L3185)
walks direct children and reads each child's active queue for the evaluated rule
category; an empty queue supplies no rule. This supports the narrow distinction
between a retained name and an active filtering effect. The helper additionally
checks empty NAT, table and child readbacks rather than treating any retained
namespace as safe. The source tag is an authoritative implementation reference,
not a claim of a byte match to the hosted kernel or userland. No cause such as a
specific reference leak is asserted. Private kernel references and inactive
transactions are not measured by these readbacks; they are not force-released.

The previous aggregate failures alone were insufficient evidence. These hosted
predicate/shape observations now support only the exact owned spellings and
lifecycle described above. Parser-only `-n` output and mock success remain
insufficient to establish native PF/TCP behavior, Apple acceptance or publication.
The repaired helper still requires new hosted lifecycle qualification on both
architectures, including failure and handled cancellation, before that claim.

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

For the targeted diagnostic regression scenarios only:

```sh
bash scripts/test-macos-ticket-transport.sh --readback-diagnostics
```

These scenarios inject each failed readback predicate at the real helper
lifecycle seam, including mixed leaf/qualified aliases and foreign inventory
changes. To select the new inventory and empty-namespace regressions only:

```sh
bash scripts/test-macos-ticket-transport.sh --anchor-custody
```

The new fixtures reproduce the observed three-row qualified shape using synthetic
identities, then exercise setup, active checks, guarded pre-success cleanup,
repeated cleanup and later namespace disappearance. They preserve the leaf-format
compatibility case and reject duplicates, mixed aliases, malformed/hierarchical
paths, sibling drift, unexpected owned contents and failed captures. Supported
qualified output and an exactly bound empty namespace replace the earlier
hypothesis refusals explicitly; their adversarial boundaries remain tested.
The full suite retains its original negative, cancellation, credential, custody,
diagnostic-output-I/O and workflow assertions. These fixtures are mock regressions,
not native qualification or proof of a kernel root cause.

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
