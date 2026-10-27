#!/usr/bin/env bash
# Temporary UDP/443 preference for qualified, ephemeral macOS release jobs only.
# No root ruleset loads, global enable/disable, state flushes, or trust changes.
set -euo pipefail
umask 077
export LC_ALL=C
fail() { printf 'ticket-transport: %s\n' "$1" >&2; exit 1; }
mode=${1:-}; shift || true
case $mode in setup|check|cleanup|run) ;; *) fail 'expected setup, check, cleanup or run' ;; esac
# Native artifact consumers temporarily use the original build's run identity.
# A cleanup hook must retain the enclosing guard's actual workflow identity.
if [[ $mode == cleanup && -n ${MACOS_TICKET_TRANSPORT_WORKFLOW_RUN_ID:-} ]]; then
  [[ $MACOS_TICKET_TRANSPORT_WORKFLOW_RUN_ID =~ ^[1-9][0-9]*$ &&
     ${MACOS_TICKET_TRANSPORT_WORKFLOW_RUN_ATTEMPT:-} =~ ^[1-9][0-9]*$ ]] || fail 'invalid guarded workflow identity'
  export GITHUB_RUN_ID=$MACOS_TICKET_TRANSPORT_WORKFLOW_RUN_ID
  export GITHUB_RUN_ATTEMPT=$MACOS_TICKET_TRANSPORT_WORKFLOW_RUN_ATTEMPT
fi
if [[ -n ${MACOS_TICKET_TRANSPORT_TEST_ROOT:-} ]]; then
  mock=$MACOS_TICKET_TRANSPORT_TEST_ROOT
  [[ ${GITHUB_ACTIONS:-false} != true && -d $mock && ! -L $mock && -O $mock &&
     -f $mock/isolated-mock && ! -L $mock/isolated-mock &&
     $(cat "$mock/isolated-mock") == isolated-ticket-transport-mock-v1 &&
     -x $mock/pfctl && -f $mock/pfctl && ! -L $mock/pfctl && -O $mock/pfctl ]] || fail 'isolated mock required; mocks are forbidden in Actions'
  pf() { "$mock/pfctl" "$@"; }
else
  [[ ${GITHUB_ACTIONS:-} == true && ${RUNNER_OS:-} == macOS &&
     ${MACOS_TICKET_TRANSPORT_RUNNER_ENVIRONMENT:-} == github-hosted &&
     $(/usr/bin/uname -s) == Darwin && ${GITHUB_REPOSITORY:-} == Dunky13/takedown &&
     ${GITHUB_REF:-} == refs/heads/main && ${GITHUB_REF_PROTECTED:-} == true &&
     ${ImageOS:-} == macos15 && -n ${ImageVersion:-} ]] || fail 'only the hosted macOS release environment is supported'
  case ${GITHUB_WORKFLOW_REF:-} in
    Dunky13/takedown/.github/workflows/release.yml@refs/heads/main|Dunky13/takedown/.github/workflows/release-macos.yml@refs/heads/main) ;;
    *) fail 'unexpected release workflow' ;;
  esac
  pf() { /usr/bin/sudo -n /sbin/pfctl "$@"; }
fi
[[ ${GITHUB_RUN_ID:-} =~ ^[1-9][0-9]*$ && ${GITHUB_RUN_ATTEMPT:-} =~ ^[1-9][0-9]*$ &&
   ${GITHUB_SHA:-} =~ ^[0-9a-f]{40}$ && ${GITHUB_JOB:-} =~ ^(sign|acceptance)$ &&
   ${GITHUB_REPOSITORY:-} == Dunky13/takedown ]] || fail 'invalid job identity'
[[ ${RUNNER_TEMP:-} == /* && -d $RUNNER_TEMP && ! -L $RUNNER_TEMP && -O $RUNNER_TEMP ]] || fail 'private runner temporary directory required'
state="$RUNNER_TEMP/takedown-ticket-transport-$GITHUB_JOB"
hash() { shasum -a 256 "$1" | awk '{print $1}'; }
private_mode() {
  if [[ $(uname -s) == Darwin ]]; then
    [[ $(/usr/bin/stat -f '%Lp' "$1") == "$2" ]]
  else
    [[ $(stat -c '%a' "$1") == "$2" ]]
  fi
}
# All PF stdout/stderr, including enable-reference tokens, stays private.
capture() { pf "$@" > "$state/readback" 2> "$state/pf-error"; }
references() {
  capture -s References || return 1
  # Keep PID/name/token identities. The human-readable timestamp/age column is
  # not reference identity and must not make a later cleanup compare unequal.
  # Unknown output is unsupported, rather than a guessed token format.
  awk '
    /^[[:space:]]*$/ || /^TOKENS:$/ || /^No pf_enabled references$/ || /^No pf starter references held$/ {next}
    /^[[:space:]]*PID[[:space:]]+Process Name[[:space:]]+TOKEN[[:space:]]+TIMESTAMP[[:space:]]*$/ {next}
    /^[[:space:]]*[0-9]+[[:space:]]/ {
      sub(/[[:space:]]+$/, "")
      if (!sub(/[[:space:]]+[0-9]+ days [0-9][0-9]:[0-9][0-9]:[0-9][0-9]$/, "")) bad=1
      if ($NF !~ /^[1-9][0-9]*$/) bad=1
      print; next
    }
    {bad=1} END {exit bad}
  ' "$state/readback" | sort > "$state/refs.now"
}
status_read() {
  capture -s info || return 1
  enabled=$(awk '$1=="Status:" && ($2=="Enabled" || $2=="Disabled") {print $2}' "$state/readback")
  [[ $enabled == Enabled || $enabled == Disabled ]]
}
anchors_read() {
  capture -a com.apple -s Anchors || return 1
  sed 's/^[[:space:]]*//' "$state/readback" | sort > "$state/anchors.now"
}
load_state() {
  if ! [[ -d $state && ! -L $state && -O $state ]] || ! private_mode "$state" 700; then fail 'missing or unsafe owned state directory'; fi
  for file in state.json state.sha256 root.before refs.before anchors.before; do
    [[ -f $state/$file && ! -L $state/$file && -O $state/$file ]] || fail 'missing or unsafe custody file'
    private_mode "$state/$file" 600 || fail 'custody permissions must be private'
  done
  [[ $(cat "$state/state.sha256") =~ ^[0-9a-f]{64}$ && $(hash "$state/state.json") == "$(cat "$state/state.sha256")" ]] || fail 'transport custody digest mismatch'
  jq -e --arg repo "$GITHUB_REPOSITORY" --arg run "$GITHUB_RUN_ID" --arg attempt "$GITHUB_RUN_ATTEMPT" \
    --arg job "$GITHUB_JOB" --arg sha "$GITHUB_SHA" \
    --arg root "$(hash "$state/root.before")" --arg refs "$(hash "$state/refs.before")" --arg anchors "$(hash "$state/anchors.before")" '
    .schema == 1 and .repository == $repo and .run == $run and .attempt == $attempt and .job == $job and .sha == $sha and
    .root_sha256 == $root and .refs_sha256 == $refs and .anchors_sha256 == $anchors and
    (.enabled == "Enabled" or .enabled == "Disabled") and
    (.anchor | type == "string" and test("^com.apple/000\\.takedown-" + $run + "-" + $attempt + "-" + $job[0:1] + "-[A-Fa-f0-9]{24}$"))
  ' "$state/state.json" > /dev/null 2>&1 || fail 'corrupt or differently bound transport custody'
  anchor=$(jq -r .anchor "$state/state.json")
  before_enabled=$(jq -r .enabled "$state/state.json")
  for file in load-attempted enable-attempted released cleaned enable-output cancelled run-group; do
    [[ ! -e $state/$file && ! -L $state/$file ]] ||
      [[ -f $state/$file && ! -L $state/$file && -O $state/$file ]] || fail 'unsafe lifecycle journal'
    if [[ -e $state/$file ]]; then private_mode "$state/$file" 600 || fail 'lifecycle journal permissions must be private'; fi
  done
}
root_unchanged() {
  capture -s rules && cmp -s "$state/root.before" "$state/readback"
}
owned_token() {
  [[ -f $state/enable-output && ! -L $state/enable-output ]] || return 1
  token=$(sed -nE 's/^Token : ([0-9]+)$/\1/p' "$state/enable-output")
  [[ $token =~ ^[1-9][0-9]*$ && ${#token} -le 20 ]]
}
references_owned() {
  references || return 1
  # The token must occur in exactly one reference, and removing that record
  # must leave pre-existing PID/name/token identities byte-for-byte unchanged.
  # Prefixes force string comparison, avoiding rounding of uint64 tokens.
  awk -v token="$token" '{if("t:" $NF == "t:" token) count++; else print} END {if(count!=1) exit 1}' \
    "$state/refs.now" > "$state/refs.other" && cmp -s "$state/refs.before" "$state/refs.other"
}
rule_text() { printf 'block return out quick proto udp from any to any port 443 label "takedown-ticket-transport"\n'; }
rule_readback() {
  capture -a "$anchor" -s rules || return 1
  # pfctl prints a numeric port as either 443 or its standard service name.
  sed 's/port = https /port = 443 /' "$state/readback" > "$state/rule.now"
  printf 'block return out quick proto udp from any to any port = 443 label "takedown-ticket-transport"\n' > "$state/rule.expected"
  cmp -s "$state/rule.expected" "$state/rule.now"
}
anchors_owned() {
  anchors_read || return 1
  awk -v own="${anchor#com.apple/}" '$0==own {count++; next} {print} END {if(count!=1) exit 1}' \
    "$state/anchors.now" > "$state/anchors.other" && cmp -s "$state/anchors.before" "$state/anchors.other"
}
restored() {
  root_unchanged && references && cmp -s "$state/refs.before" "$state/refs.now" &&
    status_read && [[ $enabled == "$before_enabled" ]] &&
    capture -a "$anchor" -s rules && [[ ! -s $state/readback ]] &&
    anchors_read && cmp -s "$state/anchors.before" "$state/anchors.now"
}
owned_group_live() {
  # Zombies cannot execute or checkpoint. Unknown process readback is treated
  # as still active, never as permission to tear down the transport guard.
  /bin/ps -axo pgid=,stat= > "$state/process-readback" 2> "$state/process-error" || return 0
  awk -v own="$1" '"g:" $1 == "g:" own && $2 !~ /^Z/ {live=1} END {exit !live}' "$state/process-readback"
}
cleanup() {
  load_state
  if [[ -f $state/cancelled && -f $state/run-group ]]; then
    group=$(cat "$state/run-group")
    [[ $group =~ ^[1-9][0-9]*$ ]] || fail 'invalid owned process group custody'
    ! owned_group_live "$group" || fail 'cancelled command still active; transport teardown withheld'
  fi
  if [[ -f $state/cleaned ]]; then
    restored || fail 'post-cleanup preservation readback failed'
    printf 'ticket-transport: cleanup verified\n'
    return
  fi
  # Flush only our unique child. Do not release PF while owned rules remain.
  if [[ -f $state/load-attempted ]]; then
    pf -a "$anchor" -F rules > "$state/flush-output" 2>&1 || fail 'owned anchor cleanup failed'
    capture -a "$anchor" -s rules && [[ ! -s $state/readback ]] || fail 'owned anchor cleanup readback failed'
  fi
  if [[ -f $state/enable-attempted && ! -f $state/released ]]; then
    owned_token || fail 'owned enable token unavailable; no guessed reference released'
    # A handled interruption may land after -X but before its journal write.
    # Accept that only when the complete original references/status read back.
    references || fail 'enable-reference readback failed'
    if cmp -s "$state/refs.before" "$state/refs.now"; then
      status_read && [[ $enabled == "$before_enabled" ]] || fail 'released reference status mismatch'
    else
      references_owned || fail 'enable-reference custody changed; no reference released'
      pf -X "$token" > "$state/release-output" 2>&1 || fail 'owned reference release failed'
    fi
    : > "$state/released"
  fi
  restored || fail 'PF preservation readback failed'
  : > "$state/cleaned"
  printf 'ticket-transport: cleanup verified\n'
}
check() {
  load_state
  [[ -f $state/load-attempted && -f $state/enable-attempted && ! -e $state/released && ! -e $state/cleaned ]] || fail 'transport guard is not active'
  owned_token && root_unchanged && rule_readback && anchors_owned && references_owned &&
    status_read && [[ $enabled == Enabled ]] || fail 'active transport readback failed'
  printf 'ticket-transport: active readback verified\n'
}
case $mode in
  setup)
    # Setup is a credential-free step; never run it with Apple secrets loaded.
    for secret in MACOS_CERTIFICATE_P12_BASE64 MACOS_CERTIFICATE_PASSWORD MACOS_SIGN_IDENTITY APPLE_TEAM_ID APPLE_API_KEY_BASE64 APPLE_API_KEY_ID APPLE_API_ISSUER; do
      [[ -z ${!secret:-} ]] || fail 'setup must run before Apple credential exposure'
    done
    mkdir -m 700 "$state" 2>/dev/null || fail 'transport state already exists'
    # Before mutation, a failed preflight needs no PF cleanup. After the journal
    # is committed, every failure/signal attempts cleanup and retains custody.
    armed=false
    setup_exit() {
      result=$?
      trap - EXIT INT TERM
      if [[ $result != 0 && $armed == true ]]; then
        bash "$0" cleanup || result=1
      fi
      exit "$result"
    }
    trap setup_exit EXIT
    trap 'exit 130' INT
    trap 'exit 143' TERM
    capture -s rules || fail 'root rules unavailable'
    cp "$state/readback" "$state/root.before"
    # An unconditioned wildcard must be reachable, with no earlier pass/quick
    # rule. Unknown root layouts are unsupported; never attach a new root rule.
    awk '
      /^scrub-anchor "com.apple\/\*" all fragment reassemble$/ {next}
      /^anchor "com.apple\/\*" all$/ {count++; next}
      {bad=1} END {exit (bad || count!=1)}
    ' "$state/root.before" || fail 'unsupported existing root anchor rules'
    references || fail 'unsupported PF reference readback'; cp "$state/refs.now" "$state/refs.before"
    status_read || fail 'PF status unavailable'
    # Legacy -e enablement without a reference cannot be preserved by -X.
    [[ ( $enabled == Enabled && -s $state/refs.before ) || ( $enabled == Disabled && ! -s $state/refs.before ) ]] || fail 'unsupported PF enable/reference state'
    anchors_read || fail 'anchor inventory unavailable'; cp "$state/anchors.now" "$state/anchors.before"
    if [[ -n ${mock:-} ]]; then
      nonce=$(od -An -N 16 -tx1 /dev/urandom | tr -d ' \n')
      [[ $nonce =~ ^[a-f0-9]{32}$ ]] || fail 'mock anchor identity unavailable'
    else
      nonce=$(/usr/bin/uuidgen)
      [[ $nonce =~ ^[A-Fa-f0-9-]{36}$ ]] || fail 'unique anchor identity unavailable'
      nonce=${nonce//-/}
    fi
    anchor="com.apple/000.takedown-$GITHUB_RUN_ID-$GITHUB_RUN_ATTEMPT-${GITHUB_JOB:0:1}-${nonce:0:24}"
    [[ ${#anchor} -lt 64 ]] || fail 'job identity exceeds supported anchor length'
    while IFS= read -r sibling; do
      [[ -z $sibling || $sibling > ${anchor#com.apple/} ]] || fail 'earlier sibling anchor may bypass transport guard'
    done < "$state/anchors.before"
    capture -a "$anchor" -s rules && [[ ! -s $state/readback ]] || fail 'owned anchor collision'
    jq -n --arg repo "$GITHUB_REPOSITORY" --arg run "$GITHUB_RUN_ID" --arg attempt "$GITHUB_RUN_ATTEMPT" \
      --arg job "$GITHUB_JOB" --arg sha "$GITHUB_SHA" --arg anchor "$anchor" --arg enabled "$enabled" \
      --arg root "$(hash "$state/root.before")" --arg refs "$(hash "$state/refs.before")" --arg anchors "$(hash "$state/anchors.before")" \
      '{schema:1,repository:$repo,run:$run,attempt:$attempt,job:$job,sha:$sha,anchor:$anchor,enabled:$enabled,root_sha256:$root,refs_sha256:$refs,anchors_sha256:$anchors}' > "$state/state.json"
    hash "$state/state.json" > "$state/state.sha256"
    armed=true
    rule_text > "$state/rule.input"
    pf -a "$anchor" -n -f "$state/rule.input" > "$state/parse-output" 2>&1 || fail 'owned rule syntax unsupported'
    : > "$state/load-attempted"
    pf -a "$anchor" -f "$state/rule.input" > "$state/load-output" 2>&1 || fail 'owned rule load failed'
    if ! rule_readback || ! root_unchanged || ! anchors_owned; then fail 'owned rule load readback failed'; fi
    : > "$state/enable-attempted"
    pf -E > "$state/enable-output" 2>&1 || fail 'owned enable reference acquisition failed'
    check
    ;;
  check) check ;;
  cleanup) cleanup ;;
  run)
    [[ $# -gt 0 ]] || fail 'guarded command required'
    check
    child=''
    # Invoked by the EXIT trap after the guarded command or a handled signal.
    # shellcheck disable=SC2329
    finish() {
      result=$?
      trap - EXIT INT TERM
      if [[ -n $child ]]; then
        if [[ $result != 0 ]] || owned_group_live "$child"; then
          : > "$state/cancelled"
          [[ $result != 0 ]] || result=1
          # The child verifies it owns this group before executing release
          # code. Signal every wrapper/native descendant together; traps may
          # then finish credential and mount cleanup before PF is released.
          kill -TERM -- "-$child" 2>/dev/null || true
          for ((drain=0; drain<150; drain++)); do
            owned_group_live "$child" || break
            sleep 0.1
          done
          if owned_group_live "$child"; then
            result=1
            kill -KILL -- "-$child" 2>/dev/null || true
            for ((drain=0; drain<50; drain++)); do
              owned_group_live "$child" || break
              sleep 0.1
            done
          fi
        fi
        wait "$child" 2>/dev/null || true
        if owned_group_live "$child"; then
          printf 'ticket-transport: owned command did not stop; no teardown claimed\n' >&2
          exit 1
        fi
      fi
      bash "$0" cleanup || result=1
      exit "$result"
    }
    trap finish EXIT
    # Keep a startup signal from landing between spawn and PID custody.
    cancel_status=0
    trap 'cancel_status=130; : > "$state/cancelled"' INT
    trap 'cancel_status=143; : > "$state/cancelled"' TERM
    export MACOS_TICKET_TRANSPORT_WORKFLOW_RUN_ID=$GITHUB_RUN_ID
    export MACOS_TICKET_TRANSPORT_WORKFLOW_RUN_ATTEMPT=$GITHUB_RUN_ATTEMPT
    export NOTARY_ACCEPTANCE_CANCEL_FILE="$state/cancelled"
    # Bash job control creates a separate owned process group on both macOS
    # Bash 3.2 and Linux. Verify it in the child before executing any operation.
    set -m
    bash -c '
      own_group=$(/bin/ps -o pgid= -p "$$")
      own_group=${own_group//[[:space:]]/}
      [[ $own_group == "$$" ]] || exit 125
      exec "$@"
    ' ticket-transport-child "$@" & child=$!
    set +m
    printf '%s\n' "$child" > "$state/run-group"
    trap 'exit 130' INT
    trap 'exit 143' TERM
    [[ $cancel_status == 0 ]] || exit "$cancel_status"
    command_result=0
    wait "$child" || command_result=$?
    exit "$command_result"
    ;;
esac
