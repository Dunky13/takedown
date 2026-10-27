#!/usr/bin/env bash
# Isolated PF model only: never invokes sudo or the host packet filter.
set -euo pipefail
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
TEST_WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-transport-tests.XXXXXX")
trap 'rm -rf -- "$TEST_WORK"' EXIT
export MACOS_TICKET_TRANSPORT_TEST_ROOT="$TEST_WORK/mock" RUNNER_TEMP="$TEST_WORK/runner"
export GITHUB_REPOSITORY=Dunky13/takedown GITHUB_RUN_ID=123 GITHUB_RUN_ATTEMPT=2 GITHUB_JOB=sign
export GITHUB_SHA=1111111111111111111111111111111111111111
unset GITHUB_ACTIONS
mkdir -p "$MACOS_TICKET_TRANSPORT_TEST_ROOT" "$RUNNER_TEMP"
printf 'isolated-ticket-transport-mock-v1\n' > "$MACOS_TICKET_TRANSPORT_TEST_ROOT/isolated-mock"
cat > "$MACOS_TICKET_TRANSPORT_TEST_ROOT/pfctl" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
m=$MACOS_TICKET_TRANSPORT_TEST_ROOT
printf '%s\n' "$*" >> "$m/calls"
anchor=''
if [[ ${1:-} == -a ]]; then anchor=$2; shift 2; fi
case "$*" in
  '-s rules')
    if [[ -z $anchor ]]; then
      cat "$m/root"
      if [[ ${MOCK_LOAD_ROOT_DRIFT:-false} == true && -f $m/owned ]]; then echo 'pass out quick all'; fi
    else
      if [[ ${MOCK_RULE_CAPTURE_FAIL:-false} == true && -f $m/rules ]]; then
        echo 'mock-secret partial readback'; exit 7
      fi
      [[ ! -f $m/rules ]] || cat "$m/rules"
      if [[ ${MOCK_BAD_READBACK:-false} == true && -f $m/rules ]]; then echo 'pass out all'; fi
    fi ;;
  '-s Anchors')
    if [[ $anchor == com.apple ]]; then
      cat "$m/anchors"
      if [[ -s $m/rules || ( ${MOCK_KEEP_EMPTY_ANCHOR:-false} == true && -f $m/owned ) ]]; then
        if [[ ${MOCK_QUALIFIED_ANCHORS:-false} == true ]]; then
          sed 's|^|com.apple/|' "$m/owned"
        else
          cat "$m/owned"
        fi
      fi
    fi ;;
  '-s info') printf 'Status: %s                              Debug: Urgent\n' "$(cat "$m/enabled")" ;;
  '-s References')
    printf 'TOKENS:\n'
    printf 'PID      Process Name                 TOKEN                    TIMESTAMP\n'
    cat "$m/refs"
    [[ ! -f $m/token ]] || printf '999      pfctl                       %s     0 days 00:00:01\n' "$(cat "$m/token")" ;;
  '-n -f '*) [[ -n $anchor ]] ;;
  '-f '*)
    [[ $anchor == com.apple/000.takedown-* ]] || exit 90
    [[ ${MOCK_LOAD_FAIL:-false} != true ]] || exit 1
    sed 's/port 443 /port = 443 /' "$2" > "$m/rules"
    printf '%s\n' "${anchor#com.apple/}" > "$m/owned"
    if [[ ${MOCK_SIGNAL_LOAD:-false} == true ]]; then kill -TERM "$PPID"; fi ;;
  '-E')
    [[ -z $anchor ]] || exit 90
    printf '987654321123456789\n' > "$m/token"
    printf 'Enabled\n' > "$m/enabled"
    [[ ${MOCK_TOKEN_MISSING:-false} == true ]] || printf 'pf enabled\nToken : 987654321123456789\n' >&2
    if [[ ${MOCK_SIGNAL_ENABLE:-false} == true ]]; then kill -TERM "$PPID"; fi
    [[ ${MOCK_ENABLE_FAIL:-false} != true ]] ;;
  '-F rules')
    [[ $anchor == com.apple/000.takedown-* ]] || exit 90
    [[ ${MOCK_FLUSH_FAIL:-false} != true ]] || exit 1
    rm -f "$m/rules" ;;
  '-X '*)
    [[ -z $anchor && $2 == "$(cat "$m/token")" ]] || exit 90
    [[ ${MOCK_RELEASE_FAIL:-false} != true ]] || exit 1
    rm "$m/token"
    if [[ -s $m/refs ]]; then echo Enabled > "$m/enabled"; else echo Disabled > "$m/enabled"; fi ;;
  *) exit 91 ;;
esac
MOCK
chmod +x "$MACOS_TICKET_TRANSPORT_TEST_ROOT/pfctl"
helper() {
  helper_result=0
  bash "$ROOT/scripts/macos-ticket-transport.sh" "$@" > "$TEST_WORK/log" 2>&1 || helper_result=$?
  if grep -Eq '987654321123456789|111222333444555666|mock-secret|Process Name|000[.]takedown-|900[.]service' "$TEST_WORK/log"; then
    echo 'FAIL: private PF/token output escaped' >&2; exit 1
  fi
  return "$helper_result"
}
refuse() { if helper "$@"; then echo "FAIL: unexpectedly allowed $*" >&2; exit 1; fi; }
reset_fixture() {
  rm -rf "$RUNNER_TEMP"; mkdir "$RUNNER_TEMP"
  rm -f "$MACOS_TICKET_TRANSPORT_TEST_ROOT/"{rules,owned,token,calls}
  printf 'scrub-anchor "com.apple/*" all fragment reassemble\nanchor "com.apple/*" all\n' > "$MACOS_TICKET_TRANSPORT_TEST_ROOT/root"
  printf '900.service\n' > "$MACOS_TICKET_TRANSPORT_TEST_ROOT/anchors"
  printf '42       service                     111222333444555666       0 days 00:00:01\n' > "$MACOS_TICKET_TRANSPORT_TEST_ROOT/refs"
  echo Enabled > "$MACOS_TICKET_TRANSPORT_TEST_ROOT/enabled"
}
state_dir() { printf '%s/takedown-ticket-transport-%s' "$RUNNER_TEMP" "$GITHUB_JOB"; }
preserved() {
  [[ ! -f $MACOS_TICKET_TRANSPORT_TEST_ROOT/token && ! -s $MACOS_TICKET_TRANSPORT_TEST_ROOT/rules ]]
  [[ $(cat "$MACOS_TICKET_TRANSPORT_TEST_ROOT/enabled") == "$1" ]]
  cmp "$MACOS_TICKET_TRANSPORT_TEST_ROOT/root" "$(state_dir)/root.before"
  cmp "$MACOS_TICKET_TRANSPORT_TEST_ROOT/refs" "$TEST_WORK/refs.before"
}
# These are hypothesis injections, not captured hosted PF output. The real
# lifecycle must identify its failing predicate without accepting new formats
# or exposing private readbacks. Run alone for a small red/green feedback loop.
diagnostic() {
  if ! grep -Fxq "ticket-transport: phase=$1 $2=false" "$TEST_WORK/log"; then
    printf 'FAIL: missing readback diagnostic phase=%s predicate=%s\n' "$1" "$2" >&2
    exit 1
  fi
  # Only fixed, allowlisted phase/predicate booleans may use the new prefix.
  if grep '^ticket-transport: phase=' "$TEST_WORK/log" | grep -Ev '^ticket-transport: phase=(setup-load|active|cleanup|post-cleanup) ((rule|root|anchors|references|status|owned-empty|token)=(true|false)|shape=(rule|root|anchors) available=false|shape=(rule|root) available=true lines=[0-9]+ sha256=[0-9a-f]{64}|shape=anchors available=true lines=[0-9]+ owned-leaf=[0-9]+ owned-qualified=[0-9]+ remaining-before=(true|false))$'; then
    echo 'FAIL: non-allowlisted readback diagnostic' >&2; exit 1
  fi
}
shape_diagnostic() {
  if ! grep -Fxq "$1" "$TEST_WORK/log"; then
    printf 'FAIL: missing safe shape diagnostic: %s\n' "$1" >&2
    exit 1
  fi
}
for predicate in rule root anchors; do
  reset_fixture
  cp "$MACOS_TICKET_TRANSPORT_TEST_ROOT/refs" "$TEST_WORK/refs.before"
  case $predicate in
    rule) MOCK_BAD_READBACK=true refuse setup ;;
    root) MOCK_LOAD_ROOT_DRIFT=true refuse setup ;;
    anchors) MOCK_QUALIFIED_ANCHORS=true refuse setup ;;
  esac
  diagnostic setup-load "$predicate"
  if [[ $predicate == anchors ]]; then
    shape_diagnostic 'ticket-transport: phase=setup-load shape=anchors available=true lines=2 owned-leaf=0 owned-qualified=1 remaining-before=true'
  fi
  [[ ! -f $(state_dir)/enable-attempted && ! -f $MACOS_TICKET_TRANSPORT_TEST_ROOT/token ]]
  helper cleanup
  preserved Enabled
done
printf 'PASS: setup readback failures identify exact predicate before reference acquisition\n'
reset_fixture
cp "$MACOS_TICKET_TRANSPORT_TEST_ROOT/refs" "$TEST_WORK/refs.before"
MOCK_RULE_CAPTURE_FAIL=true refuse setup
diagnostic setup-load rule
shape_diagnostic 'ticket-transport: phase=setup-load shape=rule available=false'
[[ ! -f $(state_dir)/enable-attempted && ! -f $MACOS_TICKET_TRANSPORT_TEST_ROOT/token ]]
helper cleanup
preserved Enabled
printf 'PASS: failed capture exposes no partial or stale shape data\n'
for predicate in rule root anchors references status token; do
  reset_fixture; helper setup
  case $predicate in
    rule) MOCK_BAD_READBACK=true refuse check ;;
    root)
      echo 'pass out quick all' >> "$MACOS_TICKET_TRANSPORT_TEST_ROOT/root"
      refuse check
      cp "$(state_dir)/root.before" "$MACOS_TICKET_TRANSPORT_TEST_ROOT/root" ;;
    anchors) MOCK_QUALIFIED_ANCHORS=true refuse check ;;
    references)
      echo 'unsupported mock-secret reference format' >> "$MACOS_TICKET_TRANSPORT_TEST_ROOT/refs"
      refuse check
      cp "$TEST_WORK/refs.before" "$MACOS_TICKET_TRANSPORT_TEST_ROOT/refs" ;;
    status)
      echo Disabled > "$MACOS_TICKET_TRANSPORT_TEST_ROOT/enabled"
      refuse check
      echo Enabled > "$MACOS_TICKET_TRANSPORT_TEST_ROOT/enabled" ;;
    token)
      cp "$(state_dir)/enable-output" "$TEST_WORK/enable.saved"
      echo 'mock-secret unsupported token' > "$(state_dir)/enable-output"
      refuse check
      cp "$TEST_WORK/enable.saved" "$(state_dir)/enable-output" ;;
  esac
  diagnostic active "$predicate"
  [[ ! -f $(state_dir)/cleaned ]]
  helper cleanup
done
printf 'PASS: active readbacks identify predicate and withhold success\n'
for predicate in root references status owned-empty anchors; do
  reset_fixture; helper setup; helper cleanup
  # Re-enter restoration without mutation to isolate each preservation gate.
  mv "$(state_dir)/load-attempted" "$TEST_WORK/load.saved"
  for phase in cleanup post-cleanup; do
    if [[ $phase == cleanup ]]; then rm "$(state_dir)/cleaned"; fi
    case $predicate in
      root) echo 'pass out quick all' >> "$MACOS_TICKET_TRANSPORT_TEST_ROOT/root" ;;
      references) echo 'unsupported mock-secret reference format' >> "$MACOS_TICKET_TRANSPORT_TEST_ROOT/refs" ;;
      status) echo Disabled > "$MACOS_TICKET_TRANSPORT_TEST_ROOT/enabled" ;;
      owned-empty)
        printf 'pass out all\n' > "$MACOS_TICKET_TRANSPORT_TEST_ROOT/rules" ;;
      anchors) printf '950.foreign-mock-secret\n' >> "$MACOS_TICKET_TRANSPORT_TEST_ROOT/anchors" ;;
    esac
    calls=$(wc -l < "$MACOS_TICKET_TRANSPORT_TEST_ROOT/calls")
    refuse cleanup
    diagnostic "$phase" "$predicate"
    if [[ $predicate == anchors ]]; then
      shape_diagnostic "ticket-transport: phase=$phase shape=anchors available=true lines=2 owned-leaf=0 owned-qualified=0 remaining-before=false"
    fi
    tail -n "+$((calls + 1))" "$MACOS_TICKET_TRANSPORT_TEST_ROOT/calls" > "$TEST_WORK/diagnostic.calls"
    if grep -Eq -- '(^| )(-f|-F|-E|-X)( |$)' "$TEST_WORK/diagnostic.calls"; then
      echo 'FAIL: restoration diagnostic caused a mutation' >&2; exit 1
    fi
    case $predicate in
      root) cp "$(state_dir)/root.before" "$MACOS_TICKET_TRANSPORT_TEST_ROOT/root" ;;
      references) cp "$TEST_WORK/refs.before" "$MACOS_TICKET_TRANSPORT_TEST_ROOT/refs" ;;
      status) echo Enabled > "$MACOS_TICKET_TRANSPORT_TEST_ROOT/enabled" ;;
      owned-empty) rm "$MACOS_TICKET_TRANSPORT_TEST_ROOT/rules" ;;
      anchors) printf '900.service\n' > "$MACOS_TICKET_TRANSPORT_TEST_ROOT/anchors" ;;
    esac
    helper cleanup
  done
  mv "$TEST_WORK/load.saved" "$(state_dir)/load-attempted"
done
reset_fixture; helper setup
MOCK_KEEP_EMPTY_ANCHOR=true refuse cleanup
diagnostic cleanup anchors
shape_diagnostic 'ticket-transport: phase=cleanup shape=anchors available=true lines=2 owned-leaf=1 owned-qualified=0 remaining-before=true'
[[ ! -f $(state_dir)/cleaned && ! -s $MACOS_TICKET_TRANSPORT_TEST_ROOT/rules ]]
helper cleanup
printf 'PASS: restoration diagnostics preserve strict sibling and empty-child refusal\n'
closed_output() {
  closed_status=0
  bash "$ROOT/scripts/macos-ticket-transport.sh" "$@" 1>&- 2> "$TEST_WORK/closed-output.err" || closed_status=$?
  if grep -Eq '987654321123456789|111222333444555666|mock-secret|Process Name|000[.]takedown-|900[.]service' "$TEST_WORK/closed-output.err"; then
    echo 'FAIL: private PF data escaped through a failed diagnostic write' >&2; exit 1
  fi
}
reset_fixture; helper setup
calls=$(wc -l < "$MACOS_TICKET_TRANSPORT_TEST_ROOT/calls")
helper check
tail -n "+$((calls + 1))" "$MACOS_TICKET_TRANSPORT_TEST_ROOT/calls" > "$TEST_WORK/writable-check.calls"
calls=$(wc -l < "$MACOS_TICKET_TRANSPORT_TEST_ROOT/calls")
closed_output check
if ! cmp -s "$(state_dir)/rule.now" "$(state_dir)/rule.expected"; then
  echo 'FAIL: diagnostic output error contaminated private rule readback files' >&2; exit 1
fi
tail -n "+$((calls + 1))" "$MACOS_TICKET_TRANSPORT_TEST_ROOT/calls" > "$TEST_WORK/closed-check.calls"
if ! cmp -s "$TEST_WORK/writable-check.calls" "$TEST_WORK/closed-check.calls"; then
  echo 'FAIL: diagnostic output error changed successful PF readback sequence' >&2; exit 1
fi
calls=$(wc -l < "$MACOS_TICKET_TRANSPORT_TEST_ROOT/calls")
MOCK_BAD_READBACK=true refuse check
tail -n "+$((calls + 1))" "$MACOS_TICKET_TRANSPORT_TEST_ROOT/calls" > "$TEST_WORK/refused-check.calls"
calls=$(wc -l < "$MACOS_TICKET_TRANSPORT_TEST_ROOT/calls")
MOCK_BAD_READBACK=true closed_output check
[[ $closed_status != 0 && ! -f $(state_dir)/cleaned ]]
tail -n "+$((calls + 1))" "$MACOS_TICKET_TRANSPORT_TEST_ROOT/calls" > "$TEST_WORK/closed-refused.calls"
if ! cmp -s "$TEST_WORK/refused-check.calls" "$TEST_WORK/closed-refused.calls"; then
  echo 'FAIL: diagnostic output error changed failed PF predicate ordering' >&2; exit 1
fi
closed_output cleanup
# The pre-existing final success-message write may still fail, but successful
# predicates must complete every preservation check and commit the journal.
[[ -f $(state_dir)/cleaned ]] || { echo 'FAIL: diagnostic output error skipped cleanup journal' >&2; exit 1; }
preserved Enabled
helper cleanup
printf 'PASS: diagnostic output errors preserve readback ordering, refusal and cleanup journal\n'
case ${1:-} in
  --readback-diagnostics) exit 0 ;;
  '') ;;
  *) echo 'unsupported test selection' >&2; exit 1 ;;
esac
reset_fixture
# These refusals must happen before even a read-only PF call.
cp "$MACOS_TICKET_TRANSPORT_TEST_ROOT/refs" "$TEST_WORK/refs.before"
env -u MACOS_TICKET_TRANSPORT_TEST_ROOT bash "$ROOT/scripts/macos-ticket-transport.sh" setup > "$TEST_WORK/log" 2>&1 && exit 1
[[ ! -f $MACOS_TICKET_TRANSPORT_TEST_ROOT/calls ]]
GITHUB_ACTIONS=true refuse setup
[[ ! -f $MACOS_TICKET_TRANSPORT_TEST_ROOT/calls ]]
APPLE_API_KEY_BASE64=mock-secret refuse setup
[[ ! -f $MACOS_TICKET_TRANSPORT_TEST_ROOT/calls ]]
printf 'PASS: local and CI mock masquerade refused before PF\n'
helper setup
helper check
sed 's/0 days 00:00:01/0 days 00:00:02/' "$MACOS_TICKET_TRANSPORT_TEST_ROOT/refs" > "$TEST_WORK/aged-refs"
cp "$TEST_WORK/aged-refs" "$MACOS_TICKET_TRANSPORT_TEST_ROOT/refs"
cp "$TEST_WORK/aged-refs" "$TEST_WORK/refs.before"
helper check
anchor1=$(jq -r .anchor "$(state_dir)/state.json")
refuse setup
helper run bash -c 'exit 0'
preserved Enabled
helper cleanup
refuse check
printf 'PASS: success cleanup preserves existing rules and enable references; stale reuse refused\n'
reset_fixture
: > "$MACOS_TICKET_TRANSPORT_TEST_ROOT/refs"; : > "$TEST_WORK/refs.before"
echo Disabled > "$MACOS_TICKET_TRANSPORT_TEST_ROOT/enabled"
helper setup
anchor2=$(jq -r .anchor "$(state_dir)/state.json")
[[ $anchor1 != "$anchor2" ]]
helper cleanup
preserved Disabled
printf 'PASS: unique anchor, originally disabled PF restored by owned reference release\n'
reset_fixture
printf '42       service                     987654321123456788       0 days 00:00:01\n' > "$MACOS_TICKET_TRANSPORT_TEST_ROOT/refs"
cp "$MACOS_TICKET_TRANSPORT_TEST_ROOT/refs" "$TEST_WORK/refs.before"
helper setup; helper cleanup
preserved Enabled
printf 'PASS: adjacent uint64 reference tokens retain exact string identity\n'
reset_fixture
cp "$MACOS_TICKET_TRANSPORT_TEST_ROOT/refs" "$TEST_WORK/refs.before"
helper setup
GITHUB_RUN_ID=124 refuse cleanup
[[ -f $MACOS_TICKET_TRANSPORT_TEST_ROOT/token ]]
helper cleanup
preserved Enabled
printf 'PASS: wrong-run cleanup refused\n'
reset_fixture
helper setup
helper run env GITHUB_RUN_ID=456 GITHUB_RUN_ATTEMPT=5 bash "$ROOT/scripts/macos-ticket-transport.sh" cleanup
preserved Enabled
printf 'PASS: pre-success cleanup retains actual workflow identity across original-build handoff\n'
for mode in load readback enable signal-load signal-enable; do
  reset_fixture
  case $mode in
    load) MOCK_LOAD_FAIL=true refuse setup ;;
    readback) MOCK_BAD_READBACK=true refuse setup ;;
    enable) MOCK_ENABLE_FAIL=true refuse setup ;;
    signal-load) MOCK_SIGNAL_LOAD=true refuse setup ;;
    signal-enable) MOCK_SIGNAL_ENABLE=true refuse setup ;;
  esac
  preserved Enabled
 done
printf 'PASS: load/readback/enable failures and setup interruptions roll back owned changes\n'
reset_fixture
printf 'pass out quick all\n' > "$MACOS_TICKET_TRANSPORT_TEST_ROOT/root"
refuse setup
[[ ! -f $MACOS_TICKET_TRANSPORT_TEST_ROOT/token && ! -f $MACOS_TICKET_TRANSPORT_TEST_ROOT/rules ]]
reset_fixture
printf '000.ahead-of-guard\n' > "$MACOS_TICKET_TRANSPORT_TEST_ROOT/anchors"
refuse setup
printf 'PASS: unsupported root and earlier sibling anchor refused\n'
reset_fixture
: > "$MACOS_TICKET_TRANSPORT_TEST_ROOT/refs"
refuse setup
[[ ! -f $MACOS_TICKET_TRANSPORT_TEST_ROOT/token ]]
printf 'PASS: legacy enabled PF without an existing reference refused\n'
reset_fixture
helper setup
cp "$(state_dir)/root.before" "$TEST_WORK/root.saved"
echo 'pass out quick all' >> "$MACOS_TICKET_TRANSPORT_TEST_ROOT/root"
refuse check
refuse cleanup
[[ -f $MACOS_TICKET_TRANSPORT_TEST_ROOT/token || -f $(state_dir)/released ]]
[[ $(tail -1 "$MACOS_TICKET_TRANSPORT_TEST_ROOT/root") == 'pass out quick all' ]]
cp "$TEST_WORK/root.saved" "$MACOS_TICKET_TRANSPORT_TEST_ROOT/root"
helper cleanup
preserved Enabled
printf 'PASS: root drift fails readback without overwriting external rules\n'
reset_fixture
helper setup
printf '43       service2                    111222333444555667       0 days 00:00:01\n' >> "$MACOS_TICKET_TRANSPORT_TEST_ROOT/refs"
refuse cleanup
[[ -f $MACOS_TICKET_TRANSPORT_TEST_ROOT/token && ! -f $(state_dir)/cleaned ]]
cp "$TEST_WORK/refs.before" "$MACOS_TICKET_TRANSPORT_TEST_ROOT/refs"
helper cleanup
preserved Enabled
printf 'PASS: changed references fail custody readback without releasing another owner\n'
reset_fixture; helper setup; helper cleanup
rm "$(state_dir)/released" "$(state_dir)/cleaned"
releases=$(grep -c '^-X ' "$MACOS_TICKET_TRANSPORT_TEST_ROOT/calls")
helper cleanup
[[ $(grep -c '^-X ' "$MACOS_TICKET_TRANSPORT_TEST_ROOT/calls") == "$releases" ]]
preserved Enabled
printf 'PASS: interruption after reference release replays readback without a second release\n'
for failure in missing corrupt json-digest snapshot symlink permissions; do
  reset_fixture; helper setup
  case $failure in
    missing) mv "$(state_dir)/state.json" "$TEST_WORK/saved.json" ;;
    corrupt) cp "$(state_dir)/state.json" "$TEST_WORK/saved.json"; echo invalid > "$(state_dir)/state.json" ;;
    json-digest) cp "$(state_dir)/state.json" "$TEST_WORK/saved.json"; jq '.enabled="Disabled"' "$TEST_WORK/saved.json" > "$(state_dir)/state.json" ;;
    snapshot) cp "$(state_dir)/root.before" "$TEST_WORK/saved.root"; echo invalid >> "$(state_dir)/root.before" ;;
    symlink) mv "$(state_dir)/state.json" "$TEST_WORK/saved.json"; ln -s "$TEST_WORK/saved.json" "$(state_dir)/state.json" ;;
    permissions) chmod 644 "$(state_dir)/state.json" ;;
  esac
  calls=$(wc -l < "$MACOS_TICKET_TRANSPORT_TEST_ROOT/calls")
  refuse cleanup
  [[ $(wc -l < "$MACOS_TICKET_TRANSPORT_TEST_ROOT/calls") == "$calls" ]]
  case $failure in
    snapshot) cp "$TEST_WORK/saved.root" "$(state_dir)/root.before" ;;
    permissions) chmod 600 "$(state_dir)/state.json" ;;
    *) rm -f "$(state_dir)/state.json"; mv "$TEST_WORK/saved.json" "$(state_dir)/state.json" ;;
  esac
  helper cleanup
 done
printf 'PASS: missing/corrupt/symlinked custody refused without guessed cleanup\n'
reset_fixture; helper setup
refuse run bash -c 'exit 23'
preserved Enabled
reset_fixture; helper setup
# Expansion belongs to the child, which signals the real helper's run trap.
# shellcheck disable=SC2016
refuse run bash -c 'kill -TERM "$PPID"; exit 0'
preserved Enabled
printf 'PASS: failed and interrupted guarded commands clean up\n'
for cancel_phase in startup armed native-failure; do
  reset_fixture; helper setup
  ruby - "$ROOT" "$TEST_WORK" "$(state_dir)" "$cancel_phase" <<'RUBY'
require 'timeout'
root, work, state, phase = ARGV
raise 'unknown cancellation phase' unless %w[startup armed native-failure].include?(phase)
fixture = File.join(work, "nested-cancellation-#{phase}")
Dir.mkdir(fixture)
File.write(File.join(fixture, 'leaf.sh'), <<~'SH')
  #!/usr/bin/env bash
  set -euo pipefail
  echo "$$" > "$NESTED_FIXTURE/leaf.pid"
  sleep 30 & native=$!
  echo "$native" > "$NESTED_FIXTURE/native.pid"
  # PID files alone do not establish that the cleanup handler is installed.
  : > "$NESTED_FIXTURE/startup-ready"
  while [[ ! -f $NESTED_FIXTURE/allow-trap ]]; do sleep 0.01; done
  cleanup_leaf() {
    status=$?
    trap - EXIT INT TERM
    kill -TERM "$native" 2>/dev/null || true
    wait "$native" 2>/dev/null || true
    if [[ -f $MACOS_TICKET_TRANSPORT_TEST_ROOT/token && -s $MACOS_TICKET_TRANSPORT_TEST_ROOT/rules ]]; then
      echo ordered > "$NESTED_FIXTURE/cleanup"
    fi
    exit "$status"
  }
  # Child failure can make wait exit under -e before the leaf receives TERM.
  # Release wrappers clean up on EXIT as well as handled signals.
  trap cleanup_leaf EXIT
  trap 'exit 143' TERM
  trap 'exit 130' INT
  : > "$NESTED_FIXTURE/armed"
  wait "$native"
  echo invalid > "$NESTED_FIXTURE/late-success"
SH
File.write(File.join(fixture, 'wrapper.sh'), <<~'SH')
  #!/usr/bin/env bash
  bash "$NESTED_FIXTURE/leaf.sh" & wait "$!"
SH
guard = nil
live = ->(pid) { !IO.popen(['ps', '-o', 'stat=', '-p', pid.to_s], &:read).strip.match?(/\A(?:Z.*)?\z/) }
begin
  guard = Process.spawn({'NESTED_FIXTURE'=>fixture}, 'bash', File.join(root, 'scripts/macos-ticket-transport.sh'), 'run', 'bash', File.join(fixture, 'wrapper.sh'), out: File.join(fixture, 'guard.log'), err: [:child, :out])
  # Publish complete PID custody before testing startup cancellation. The armed
  # case additionally waits for the real TERM handler, rather than racing it.
  Timeout.timeout(10) { sleep 0.02 until File.exist?(File.join(fixture, 'startup-ready')) }
  if phase != 'startup'
    File.write(File.join(fixture, 'allow-trap'), '')
    Timeout.timeout(10) { sleep 0.02 until File.exist?(File.join(fixture, 'armed')) }
  end
  if phase == 'native-failure'
    Process.kill('TERM', Integer(File.read(File.join(fixture, 'native.pid'))))
  else
    Process.kill('TERM', guard)
  end
  _, status = Timeout.timeout(10) { Process.wait2(guard) }
  raise 'guard cancellation did not fail' unless status.exitstatus == 143
  if phase != 'startup'
    raise 'nested cleanup did not complete before PF teardown' unless File.read(File.join(fixture, 'cleanup')) == "ordered\n"
  else
    raise 'startup cancellation claimed an uninstalled cleanup handler' if File.exist?(File.join(fixture, 'cleanup'))
  end
  %w[leaf native].each do |kind|
    raise 'nested command survived guard cleanup' if live.call(Integer(File.read(File.join(fixture, kind + '.pid'))))
  end
  raise 'cancelled command wrote late success' if File.exist?(File.join(fixture, 'late-success'))
  raise 'PF cleanup did not complete' unless File.exist?(File.join(state, 'cleaned'))
  puts "PASS: #{phase} cancellation drains owned descendants before PF teardown"
ensure
  %w[native leaf].each do |kind|
    file = File.join(fixture, kind + '.pid')
    begin Process.kill('KILL', Integer(File.read(file))) if File.exist?(file); rescue Errno::ESRCH; end
  end
  begin Process.kill('KILL', guard) if guard; rescue Errno::ESRCH; end
  begin Process.wait(guard) if guard; rescue Errno::ECHILD; end
end
RUBY
  preserved Enabled
done
printf 'PASS: nested cancellation finishes owned native cleanup before PF teardown and cannot write late success\n'
for failure in flush release; do
  reset_fixture; helper setup
  case $failure in
    flush) MOCK_FLUSH_FAIL=true refuse run bash -c 'exit 0' ;;
    release) MOCK_RELEASE_FAIL=true refuse run bash -c 'exit 0' ;;
  esac
  [[ ! -f $(state_dir)/cleaned ]]
  helper cleanup
  preserved Enabled
 done
printf 'PASS: cleanup failure overrides successful operation; retry verifies restoration\n'
reset_fixture
MOCK_TOKEN_MISSING=true refuse setup
[[ -f $MACOS_TICKET_TRANSPORT_TEST_ROOT/token && ! -f $(state_dir)/cleaned ]]
printf 'PASS: missing owned token fails closed without releasing guessed reference\n'
# Static workflow contract assertions plus execution of the actual cleanup and
# promotion shell blocks, parsed by Ruby stdlib YAML (no provider calls).
ruby - "$ROOT" "$TEST_WORK" <<'RUBY'
require 'yaml'
require 'fileutils'
root, work = ARGV
workflow = YAML.load_file(File.join(root, '.github/workflows/release-macos.yml'))
jobs = workflow.fetch('jobs')
assert = ->(condition, message) { raise message unless condition }
%w[sign acceptance].each do |name|
  job = jobs.fetch(name)
  steps = job.fetch('steps')
  setup = steps.index { |s| s['id'] == 'transport_setup' }
  cleanup = steps.index { |s| s['id'] == 'transport_cleanup' }
  operation = steps.index { |s| s.fetch('run', '').include?('macos-ticket-transport.sh run') }
  qualification = steps.index { |s| s.fetch('run', '').match?(/load-apple-input|release-source.sh verify/) }
  assert.call(qualification < setup && setup < operation && operation < cleanup, "#{name}: unsafe sequence")
  assert.call(steps[setup]['run'] == 'bash scripts/macos-ticket-transport.sh setup', 'unexpected setup command')
  assert.call(steps[setup].fetch('env').keys == ['MACOS_TICKET_TRANSPORT_RUNNER_ENVIRONMENT'], 'setup exposes credentials')
  assert.call(steps[operation].fetch('env').fetch('NOTARY_ACCEPTANCE_REQUALIFY_SHA256') == '${{ vars.APPLE_REQUALIFY_ACCEPTANCE_SHA256 }}', 'approval pin not forwarded')
  assert.call(steps[operation].fetch('env').fetch('NOTARY_ACCEPTANCE_CLEANUP_HOOK') == '${{ github.workspace }}/scripts/macos-ticket-transport.sh', 'pre-success cleanup hook missing')
  [setup, operation, cleanup].each do |index|
    assert.call(steps[index].fetch('env').fetch('MACOS_TICKET_TRANSPORT_RUNNER_ENVIRONMENT') == '${{ runner.environment }}', 'missing runner gate')
  end
  assert.call(steps[cleanup]['if'] == "always() && steps.transport_setup.outcome != 'skipped'", 'cleanup does not handle failure/cancellation')
end
sign = jobs.fetch('sign')
assert.call(sign['environment'] == 'release-signing' && sign['if'] == "github.ref == 'refs/heads/main' && github.ref_protected", 'signing authority changed')
acceptance = jobs.fetch('acceptance')
assert.call(!acceptance.key?('environment') && !acceptance.to_s.include?('secrets.'), 'fresh runner receives Apple secrets')
steps = acceptance.fetch('steps')
cleanup = steps.find { |s| s['id'] == 'transport_cleanup' }
promotion = steps.find { |s| s.fetch('name', '').start_with?('Promote fresh receipt') }
upload = steps.find { |s| s.fetch('uses', '').start_with?('actions/upload-artifact@') }
assert.call(steps.index(cleanup) < steps.index(promotion) && steps.index(promotion) < steps.index(upload), 'receipt promoted before cleanup')
[promotion, upload].each do |step|
  assert.call(step['if'] == "success() && steps.transport_cleanup.outcome == 'success'", 'receipt lacks cleanup success gate')
end
operation = steps.find { |s| s['id'] == 'native_acceptance' }
assert.call(operation.fetch('run').include?('"$receipt.pending"'), 'receipt not deferred')
fixture = File.join(work, 'workflow')
FileUtils.mkdir_p(File.join(fixture, 'scripts'))
File.write(File.join(fixture, 'scripts/macos-ticket-transport.sh'), <<~SH)
  #!/bin/bash
  if [[ $1 == run ]]; then shift; "$@"; else exit "${MOCK_CLEANUP_STATUS:?}"; fi
SH
File.write(File.join(fixture, 'scripts/fresh-apple-acceptance.sh'), <<~SH)
  #!/bin/bash
  [[ ${MOCK_NATIVE_STATUS:-0} == 0 ]] || exit "$MOCK_NATIVE_STATUS"
  printf 'fresh receipt' > "$NOTARY_EVIDENCE_DIR/fresh-acceptance-$NOTARY_ARCH.json"
SH
temp = File.join(fixture, 'temp')
evidence = File.join(temp, 'fresh-evidence')
FileUtils.mkdir_p(evidence)
receipt = File.join(evidence, 'fresh-acceptance-arm64.json')
run = ->(script, extra) { system({'RUNNER_TEMP'=>temp, 'ARCH'=>'arm64', 'NATIVE_ACCEPTANCE_OUTCOME'=>'success'}.merge(extra), 'bash', '-euo', 'pipefail', '-c', script, chdir: fixture, out: File::NULL, err: File::NULL) }
# Execute the real guarded operation body: success defers the receipt before
# cleanup, while a native failure must remain failed under the nested shell.
assert.call(run.call(operation['run'], {'NOTARY_EVIDENCE_DIR'=>evidence, 'NOTARY_ARCH'=>'arm64'}), 'guarded acceptance failed')
assert.call(!File.exist?(receipt) && File.read(receipt + '.pending') == 'fresh receipt', 'receipt not deferred within guard')
assert.call(!run.call(operation['run'], {'NOTARY_EVIDENCE_DIR'=>evidence, 'NOTARY_ARCH'=>'arm64', 'MOCK_NATIVE_STATUS'=>'23'}), 'native failure swallowed')
# The actual workflow cleanup body must remove both final and pending receipts.
File.write(receipt, 'final'); File.write(receipt + '.pending', 'pending')
assert.call(!run.call(cleanup['run'], {'MOCK_CLEANUP_STATUS'=>'1'}), 'cleanup failure passed')
assert.call(!File.exist?(receipt) && !File.exist?(receipt + '.pending'), 'cleanup failure retained receipt')
File.write(receipt + '.pending', 'pending')
assert.call(run.call(cleanup['run'], {'MOCK_CLEANUP_STATUS'=>'0', 'NATIVE_ACCEPTANCE_OUTCOME'=>'failure'}), 'cleanup failed unexpectedly')
assert.call(!File.exist?(receipt + '.pending'), 'failed operation retained receipt')
File.write(receipt + '.pending', 'pending')
assert.call(run.call(cleanup['run'], {'MOCK_CLEANUP_STATUS'=>'0'}), 'successful cleanup failed')
assert.call(run.call(promotion['run'], {}), 'promotion failed')
assert.call(File.read(receipt) == 'pending' && !File.exist?(receipt + '.pending'), 'receipt promotion wrong')
# Execute the real readiness retention body. Replayed ARM proof keeps its old
# identity while the ordinary Intel proof belongs to the current workflow.
retention = jobs.fetch('readiness').fetch('steps').find { |s| s['name'] == 'Retain both fresh native acceptance receipts' }
proofs = File.join(fixture, '.local/fresh-acceptance')
FileUtils.mkdir_p(proofs)
require 'json'
arm_proof = {verification_run_id:'777', verification_run_attempt:'2'}.to_json
File.write(File.join(proofs, 'fresh-acceptance-arm64.json'), arm_proof)
File.write(File.join(proofs, 'fresh-acceptance-amd64.json'), {verification_run_id:'888', verification_run_attempt:'3'}.to_json)
File.write(File.join(fixture, 'scripts/apple-release-state.sh'), <<~'SH')
  #!/bin/bash
  [[ $1 == put && -f $3 ]] || exit 1
  printf '%s\n' "$2" >> "$RUNNER_TEMP/retained-proof-names"
SH
assert.call(run.call(retention.fetch('run'), {'GITHUB_RUN_ID'=>'888', 'GITHUB_RUN_ATTEMPT'=>'3'}), 'proof retention failed')
expected_names = %w[fresh-acceptance-arm64-777-2.json fresh-acceptance-arm64-888-3.json fresh-acceptance-amd64-888-3.json]
assert.call(File.readlines(File.join(temp, 'retained-proof-names'), chomp:true) == expected_names, 'original proof identity lost or retention duplicated')
assert.call(File.read(File.join(proofs, 'fresh-acceptance-arm64.json')) == arm_proof, 'retention rewrote original proof')
RUBY
printf 'PASS: transport workflow sequencing, cleanup and fresh-receipt promotion\n'
