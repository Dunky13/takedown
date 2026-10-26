#!/usr/bin/env bash
# Synthetic origin/check/duplicate planning tests, no release or Apple operations.
set -euo pipefail
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
cd "$ROOT"
TEST_WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-plan-tests.XXXXXX")
trap 'rm -rf -- "$TEST_WORK"' EXIT
mkdir "$TEST_WORK/bin"
export GH_REPO=fixture/takedown GITHUB_RUN_ID=123 GITHUB_RUN_ATTEMPT=1 SOURCE_REF=refs/heads/main EVENT=workflow_run
export TRIGGER_BRANCH=main TRIGGER_EVENT=push TRIGGER_CONCLUSION=success TRIGGER_REPOSITORY=$GH_REPO TRIGGER_PATH=.github/workflows/ci.yml
TRIGGER_SOURCE=$(git rev-parse HEAD)
export TRIGGER_SOURCE MOCK_SOURCE=$TRIGGER_SOURCE
export TRIGGER_RUN_ID=100 MOCK_SPOOF=false
export MOCK_TEST_WORK=$TEST_WORK MOCK_PROTECTED=true MOCK_CHECK=success MOCK_PUBLISHED=false CONTINUATION_ONLY=false
export MOCK_ACTIVE=false MOCK_ORIGIN_BRANCH=main
cat > "$TEST_WORK/bin/gh" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
endpoint=${!#}
printf '%s\n' "$*" >> "$MOCK_TEST_WORK/calls"
case $endpoint in
  */branches/main) jq -n --arg sha "$MOCK_SOURCE" --argjson protected "$MOCK_PROTECTED" '{protected:$protected,commit:{sha:$sha}}' ;;
  *'/releases?per_page=100')
    if [[ $MOCK_ACTIVE == true ]]; then jq -n --arg source "$MOCK_SOURCE" --rawfile body "$MOCK_TEST_WORK/reservation.json" '[{id:7,draft:true,tag_name:"apple-candidate-v0.0.1",target_commitish:$source,body:$body}]'
    elif [[ $MOCK_PUBLISHED == true ]]; then jq -n --arg source "$MOCK_SOURCE" '[{draft:false,prerelease:false,tag_name:"v0.0.1",target_commitish:$source}]'; else echo '[]'; fi ;;
  */releases/7/assets\?per_page=100)
    find "$MOCK_TEST_WORK/remote" -type f | LC_ALL=C sort | while IFS= read -r file; do name=${file##*/}; jq -n --arg name "$name" '{id:$name,name:$name}'; done | jq -s . ;;
  */releases/assets/*) cat "$MOCK_TEST_WORK/remote/${endpoint##*/}" ;;
  */actions/runs/456)
    jq -n --arg source "$MOCK_SOURCE" --arg repo "$GH_REPO" --arg branch "$MOCK_ORIGIN_BRANCH" '{head_sha:$source,head_branch:$branch,repository:{full_name:$repo},event:"workflow_run",path:".github/workflows/release.yml"}' ;;
  */check-runs\?per_page=100)
    jq -n --arg check "$MOCK_CHECK" --arg source "$MOCK_SOURCE" --argjson spoof "$MOCK_SPOOF" '{check_runs:([{id:1,name:"ci-required",app:{id:15368},check_suite:{id:10},head_sha:$source,status:"completed",conclusion:$check},{id:2,name:"conventional-commits",app:{id:15368},check_suite:{id:11},head_sha:$source,status:"completed",conclusion:"success"}] + (if $spoof then [{id:999,name:"ci-required",app:{id:15368},check_suite:{id:99},head_sha:$source,status:"completed",conclusion:"success"}] else [] end))}' ;;
  */actions/workflows/*/runs\?*)
    workflow=${endpoint#*/actions/workflows/}; workflow=${workflow%%/*}
    suite=10; id=100
    [[ $workflow != conventional-commits.yml ]] || { suite=11; id=101; }
    jq -n --arg source "$MOCK_SOURCE" --arg repo "$GH_REPO" --arg path ".github/workflows/$workflow" --argjson suite "$suite" --argjson id "$id" '{workflow_runs:[{id:$id,check_suite_id:$suite,head_sha:$source,head_branch:"main",event:"push",path:$path,repository:{full_name:$repo},head_repository:{full_name:$repo},status:"completed",conclusion:"success",run_started_at:"2026-10-26T00:00:00Z",run_attempt:1}]}' ;;
  *) echo 'unexpected GitHub operation' >&2; exit 99 ;;
esac
MOCK
cat > "$TEST_WORK/bin/go" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
[[ $* == *'./cmd/releasectl plan'* ]]
mkdir "$PLAN_DIR"
jq -n --arg source "$MOCK_SOURCE" '{version:"0.0.1",package_version:"0.0.1",sha:$source,skip:false}' > "$PLAN_DIR/plan.json"
printf 'skip=false\n' > "$PLAN_DIR/github-output"
printf 'synthetic notes\n' > "$PLAN_DIR/notes.md"
printf 'planned\n' >> "$MOCK_TEST_WORK/plans"
MOCK
chmod +x "$TEST_WORK/bin/gh" "$TEST_WORK/bin/go"
export PATH="$TEST_WORK/bin:$PATH"
plan() {
  rm -rf "$TEST_WORK/plan" "$TEST_WORK/plans"
  export PLAN_DIR="$TEST_WORK/plan" GITHUB_OUTPUT="$TEST_WORK/output"
  : > "$GITHUB_OUTPUT"
  bash scripts/prepare-apple-release.sh plan > "$TEST_WORK/log" 2>&1
}
reject() { if plan; then echo 'expected rejection' >&2; exit 1; fi; [[ ! -f $TEST_WORK/plans ]]; }
skip() { plan; grep -qx skip=true "$GITHUB_OUTPUT"; [[ ! -f $TEST_WORK/plans ]]; }
plan
grep -qx resume=false "$GITHUB_OUTPUT"
grep -qx origin_run=123 "$GITHUB_OUTPUT"
TRIGGER_EVENT=pull_request skip
TRIGGER_REPOSITORY=attacker/takedown skip
TRIGGER_PATH=.github/workflows/untrusted.yml skip
TRIGGER_CONCLUSION=failure skip
TRIGGER_SOURCE=1111111111111111111111111111111111111111 skip
MOCK_PROTECTED=false reject
MOCK_CHECK=failure reject
MOCK_CHECK=neutral reject
MOCK_CHECK=failure MOCK_SPOOF=true reject
TRIGGER_RUN_ID=999 reject
SOURCE_REF=refs/tags/v0.0.1 reject
MOCK_PUBLISHED=true skip
EVENT=schedule CONTINUATION_ONLY=true MOCK_PUBLISHED=true skip
EVENT=schedule CONTINUATION_ONLY=true plan
mkdir "$TEST_WORK/remote"
cp "$TEST_WORK/plan/plan.json" "$TEST_WORK/remote/plan.json"
printf 'original candidate notes\n' > "$TEST_WORK/remote/notes.md"
printf '[]\n' > "$TEST_WORK/remote/native.json"
jq -n --arg repo "$GH_REPO" --arg source "$MOCK_SOURCE" \
  --arg plan "$(shasum -a 256 "$TEST_WORK/remote/plan.json" | awk '{print $1}')" \
  --arg notes "$(shasum -a 256 "$TEST_WORK/remote/notes.md" | awk '{print $1}')" \
  --arg native "$(shasum -a 256 "$TEST_WORK/remote/native.json" | awk '{print $1}')" \
  '{schema:1,repository:$repo,source:$source,version:"0.0.1",run_id:"456",run_attempt:"2",channel:"manual",plan_sha256:$plan,notes_sha256:$notes,native_sha256:$native}' > "$TEST_WORK/remote/candidate.json"
jq -n --slurpfile candidate "$TEST_WORK/remote/candidate.json" --slurpfile plan "$TEST_WORK/remote/plan.json" --slurpfile native "$TEST_WORK/remote/native.json" --rawfile notes "$TEST_WORK/remote/notes.md" '{schema:1,candidate:$candidate[0],plan:$plan[0],native:$native[0],notes:$notes}' > "$TEST_WORK/reservation.json"
MOCK_ACTIVE=true plan
grep -qx resume=true "$GITHUB_OUTPUT"
grep -qx origin_run=456 "$GITHUB_OUTPUT"
grep -qx origin_attempt=2 "$GITHUB_OUTPUT"
[[ ! -f $TEST_WORK/plans ]]
MOCK_ACTIVE=true MOCK_ORIGIN_BRANCH=feature reject
cp "$TEST_WORK/reservation.json" "$TEST_WORK/good-reservation.json"
jq '.plan.sha="corrupt"' "$TEST_WORK/good-reservation.json" > "$TEST_WORK/reservation.json"
MOCK_ACTIVE=true reject
cp "$TEST_WORK/good-reservation.json" "$TEST_WORK/reservation.json"
rm "$TEST_WORK/remote/candidate.json"
MOCK_ACTIVE=true plan
grep -qx resume=true "$GITHUB_OUTPUT"
[[ ! -f $TEST_WORK/plans ]]
echo 'PASS: 19 synthetic planning cases (trusted run suites, spoof refusal, original-candidate continuation, atomic reservation and corruption)'
