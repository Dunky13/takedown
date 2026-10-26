#!/usr/bin/env bash
# Offline dual-pin qualification; no GitHub, builds, credentials or Apple calls.
set -euo pipefail
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
TEST_WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-source-tests.XXXXXX")
trap 'rm -rf -- "$TEST_WORK"' EXIT
mkdir "$TEST_WORK/bin"
export MOCK_SOURCE=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa MOCK_HEAD=bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
export GH_REPO=fixture/takedown MOCK_CALLS="$TEST_WORK/calls"
cat > "$TEST_WORK/bin/git" <<'MOCK'
#!/usr/bin/env bash
case "$1 $2" in
  'rev-parse HEAD') echo "$MOCK_HEAD" ;;
  'status --porcelain') [[ ${MOCK_DIRTY:-false} != true ]] || echo ' M tracked'; exit 0 ;;
  'merge-base --is-ancestor') [[ ${MOCK_UNRELATED:-false} != true ]] ;;
  *) exit 99 ;;
esac
MOCK
cat > "$TEST_WORK/bin/gh" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
endpoint=${!#}
printf '%s\n' "$endpoint" >> "$MOCK_CALLS"
case $endpoint in
  */branches/main) jq -n --argjson protected "${MOCK_PROTECTED:-true}" '{protected:$protected}' ;;
  */compare/*)
    jq -n --arg source "$MOCK_HEAD" --arg status "${MOCK_MAIN_STATUS:-identical}" '{base_commit:{sha:$source},merge_base_commit:{sha:$source},status:$status}' ;;
  */commits/*/check-runs\?*)
    [[ $endpoint == */commits/"$MOCK_HEAD"/* ]]
    jq -n --arg sha "$MOCK_HEAD" --arg conclusion "${MOCK_CHECK:-success}" --argjson suite "${MOCK_SUITE:-10}" '{check_runs:[{name:"ci-required",head_sha:$sha,app:{id:15368},check_suite:{id:$suite},status:"completed",conclusion:$conclusion},{name:"conventional-commits",head_sha:$sha,app:{id:15368},check_suite:{id:11},status:"completed",conclusion:"success"}]}' ;;
  */actions/workflows/*/runs\?*)
    [[ $endpoint == *"head_sha=$MOCK_HEAD"* ]]
    workflow=${endpoint#*/actions/workflows/}; workflow=${workflow%%/*}
    suite=10; [[ $workflow != conventional-commits.yml ]] || suite=11
    jq -n --arg sha "$MOCK_HEAD" --arg repo "$GH_REPO" --arg path ".github/workflows/$workflow" --arg event "${MOCK_EVENT:-push}" --arg conclusion "${MOCK_RUN:-success}" --argjson suite "$suite" '{workflow_runs:[{id:100,check_suite_id:$suite,head_sha:$sha,head_branch:"main",event:$event,path:$path,repository:{full_name:$repo},head_repository:{full_name:$repo},status:"completed",conclusion:$conclusion}]}' ;;
  *) echo "unexpected API $endpoint" >&2; exit 99 ;;
esac
MOCK
chmod +x "$TEST_WORK/bin/"*
export PATH="$TEST_WORK/bin:$PATH"
count=0
verify() {
  count=$((count + 1))
  bash "$ROOT/scripts/release-source.sh" verify "$MOCK_SOURCE" > "$TEST_WORK/log" 2>&1
}
reject() {
  local expected=$1
  if verify; then echo "expected rejection: $expected" >&2; exit 1; fi
  grep -Fq "$expected" "$TEST_WORK/log" || { cat "$TEST_WORK/log"; exit 1; }
}
# The old exact-source contract is the default, with no network calls.
MOCK_HEAD=$MOCK_SOURCE verify
reject 'exact product source'
[[ ! -f $MOCK_CALLS ]]
export RELEASE_AUTOMATION_SHA=$MOCK_HEAD RELEASE_WORKFLOW_SHA=$MOCK_HEAD
verify
[[ $(wc -l < "$MOCK_CALLS") -eq 5 ]]
RELEASE_WORKFLOW_SHA=$MOCK_SOURCE reject 'explicit workflow SHA'
RELEASE_AUTOMATION_SHA=$MOCK_SOURCE reject 'explicit workflow SHA'
MOCK_DIRTY=true reject 'clean at its pinned SHA'
MOCK_UNRELATED=true reject 'not an ancestor'
MOCK_PROTECTED=false reject 'main must remain protected'
MOCK_MAIN_STATUS=diverged reject 'not on main'
MOCK_RUN=failure reject 'trusted main workflow has not succeeded'
MOCK_EVENT=pull_request reject 'trusted main workflow has not succeeded'
MOCK_CHECK=neutral reject 'required check is not successful'
MOCK_CHECK=failure reject 'required check is not successful'
MOCK_SUITE=99 reject 'required check is not successful'
# A stale event must never displace the original candidate; fresh events for a
# new candidate must still match the current exact CI run.
EVENT=workflow_run TRIGGER_RUN_ID=999 verify
EVENT=workflow_run TRIGGER_RUN_ID=999 RELEASE_REQUIRE_TRIGGER=true reject 'current trusted native run'
echo "PASS: $count offline dual-pin source/automation qualification scenarios"
