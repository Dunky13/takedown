#!/usr/bin/env bash
# Credential-free planning and durable custody for the original build.
set -euo pipefail
umask 077
fail() { printf 'prepare-apple-release: %s\n' "$*" >&2; exit 1; }
: "${GH_REPO:?}" "${GITHUB_OUTPUT:?}" "${PLAN_DIR:?}" "${GITHUB_RUN_ID:?}" "${GITHUB_RUN_ATTEMPT:?}"
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
STATE="$ROOT/scripts/apple-release-state.sh"
mkdir -p "$PLAN_DIR"
WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-plan.XXXXXX")
trap 'rm -rf -- "$WORK"' EXIT
hash() { shasum -a 256 "$1" | awk '{print $1}'; }
outputs() {
  {
    jq -r '"version=" + .version, "package_version=" + .package_version, "sha=" + .sha' "$PLAN_DIR/plan.json"
    printf 'state_tag=%s\norigin_run=%s\norigin_attempt=%s\nresume=%s\nskip=false\n' "$NOTARY_STATE_TAG" "$origin_run" "$origin_attempt" "$resume"
    printf 'channel=%s\n' "$channel"
  } >> "$GITHUB_OUTPUT"
}
trusted_checks() {
  local workflow check suite
  gh api --paginate "repos/$GH_REPO/commits/$source/check-runs?per_page=100" > "$WORK/checks.json"
  for workflow in ci.yml conventional-commits.yml; do
    gh api --paginate "repos/$GH_REPO/actions/workflows/$workflow/runs?branch=main&event=push&head_sha=$source&per_page=100" > "$WORK/runs.json"
    jq -s '[.[].workflow_runs[]] | sort_by(.run_started_at,.run_attempt) | last' "$WORK/runs.json" > "$WORK/trusted-run.json"
    jq -e --arg source "$source" --arg repo "$GH_REPO" --arg path ".github/workflows/$workflow" '
      .head_sha == $source and .head_branch == "main" and .event == "push" and
      .repository.full_name == $repo and .head_repository.full_name == $repo and .path == $path and
      .status == "completed" and .conclusion == "success" and
      (.id | type == "number" and . > 0 and floor == .) and
      (.check_suite_id | type == "number" and . > 0 and floor == .)
    ' "$WORK/trusted-run.json" >/dev/null || fail "trusted main workflow has not succeeded: $workflow"
    if [[ ${EVENT:-} == workflow_run && $workflow == ci.yml ]]; then
      [[ ${TRIGGER_RUN_ID:-} == "$(jq -r .id "$WORK/trusted-run.json")" ]] || fail 'completion is not the current trusted native run'
    fi
    suite=$(jq -r .check_suite_id "$WORK/trusted-run.json")
    check=ci-required; [[ $workflow != conventional-commits.yml ]] || check=conventional-commits
    jq -se --arg source "$source" --arg name "$check" --argjson suite "$suite" '
      [.[].check_runs[] | select(.name == $name and .app.id == 15368 and .check_suite.id == $suite and .head_sha == $source)] |
      sort_by(.started_at) | last | .status == "completed" and .conclusion == "success"
    ' "$WORK/checks.json" >/dev/null || fail "required check is not successful in its trusted run suite: $check"
  done
}
case ${1:?expected plan or persist} in
  plan)
    [[ ${SOURCE_REF:-} == refs/heads/main ]] || fail 'release orchestration requires main'
    if [[ ${EVENT:-} == workflow_run ]]; then
      [[ ${TRIGGER_BRANCH:-} == main && ${TRIGGER_EVENT:-} == push && ${TRIGGER_CONCLUSION:-} == success && ${TRIGGER_REPOSITORY:-} == "$GH_REPO" && ${TRIGGER_PATH:-} == .github/workflows/ci.yml && ${TRIGGER_RUN_ID:-} =~ ^[1-9][0-9]*$ ]] || { printf 'skip=true\n' >> "$GITHUB_OUTPUT"; exit 0; }
    fi
    gh api "repos/$GH_REPO/branches/main" > "$WORK/main.json"
    jq -e '.protected == true' "$WORK/main.json" >/dev/null || fail 'main must remain protected'
    gh api --paginate "repos/$GH_REPO/releases?per_page=100" > "$WORK/releases.json"
    jq -s 'add | [.[] | select(.draft == true and (.tag_name | startswith("apple-candidate-v")))]' "$WORK/releases.json" > "$WORK/candidates.json"
    active=()
    while IFS= read -r candidate; do
      export NOTARY_STATE_TAG=$candidate
      bash "$STATE" inventory > "$WORK/inventory.json"
      if ! jq -e 'any(.[]; .name == "published.json")' "$WORK/inventory.json" >/dev/null; then active+=("$candidate"); fi
    done < <(jq -r '.[].tag_name' "$WORK/candidates.json")
    [[ ${#active[@]} -le 1 ]] || fail 'multiple unresolved candidates; reconcile custody'
    if [[ ${#active[@]} == 1 ]]; then
      export NOTARY_STATE_TAG=${active[0]}
      # The reservation atomically retains all original metadata in its body.
      # Recovery therefore also works before the first asset upload succeeds.
      jq -er --arg tag "$NOTARY_STATE_TAG" '.[] | select(.tag_name == $tag) | .body' "$WORK/candidates.json" > "$WORK/body.json"
      jq -e '.schema == 1 and (.candidate|type == "object") and (.plan|type == "object") and (.native|type == "array") and (.notes|type == "string")' "$WORK/body.json" >/dev/null || fail 'unbound candidate reservation'
      jq .candidate "$WORK/body.json" > "$WORK/candidate.json"
      jq .plan "$WORK/body.json" > "$PLAN_DIR/plan.json"
      jq .native "$WORK/body.json" > "$PLAN_DIR/native.json"
      jq -j .notes "$WORK/body.json" > "$PLAN_DIR/notes.md"
      jq -e --arg repo "$GH_REPO" '.schema == 1 and .repository == $repo and (.source|test("^[0-9a-f]{40}$")) and (.run_id|test("^[1-9][0-9]*$")) and (.run_attempt|test("^[1-9][0-9]*$")) and (.channel == "manual" or .channel == "nightly")' "$WORK/candidate.json" >/dev/null || fail 'invalid candidate'
      origin_run=$(jq -r .run_id "$WORK/candidate.json"); origin_attempt=$(jq -r .run_attempt "$WORK/candidate.json"); channel=$(jq -r .channel "$WORK/candidate.json")
      source=$(jq -r .source "$WORK/candidate.json")
      [[ $NOTARY_STATE_TAG == "apple-candidate-v$(jq -r .version "$WORK/candidate.json")" ]] || fail 'candidate tag mismatch'
      jq -e --arg tag "$NOTARY_STATE_TAG" --arg source "$source" 'any(.[]; .tag_name == $tag and .target_commitish == $source)' "$WORK/candidates.json" >/dev/null || fail 'candidate target mismatch'
      git merge-base --is-ancestor "$source" HEAD || fail 'candidate source is not on main'
      gh api "repos/$GH_REPO/actions/runs/$origin_run" > "$WORK/origin.json"
      jq -e --arg source "$source" --arg repo "$GH_REPO" '.head_sha == $source and .head_branch == "main" and .repository.full_name == $repo and .path == ".github/workflows/release.yml" and (.event == "workflow_run" or .event == "schedule" or .event == "workflow_dispatch")' "$WORK/origin.json" >/dev/null || fail 'untrusted original run'
      for file in plan.json notes.md native.json; do
        key=${file%%.*}_sha256
        [[ $(hash "$PLAN_DIR/$file") == "$(jq -r --arg key "$key" '.[$key]' "$WORK/candidate.json")" ]] || fail 'candidate manifest hash mismatch'
      done
      jq -e --arg sha "$source" --arg version "${NOTARY_STATE_TAG#apple-candidate-v}" '.sha == $sha and .version == $version and .skip == false' "$PLAN_DIR/plan.json" >/dev/null || fail 'candidate plan mismatch'
      resume=true
      outputs
      exit 0
    fi
    source=$(git rev-parse HEAD)
    [[ $source == "$(jq -r .commit.sha "$WORK/main.json")" ]] || fail 'planning checkout is not current main'
    if [[ ${EVENT:-} == workflow_run && ${TRIGGER_SOURCE:-} != "$source" ]]; then printf 'skip=true\n' >> "$GITHUB_OUTPUT"; exit 0; fi
    trusted_checks
    kind=${BUMP:-patch}; channel=manual; last=''
    if [[ ${EVENT:-} == schedule && ${CONTINUATION_ONLY:-false} != true ]]; then
      if jq -se --arg sha "$source" 'add | any(.[]; .draft == false and .prerelease == false and .target_commitish == $sha)' "$WORK/releases.json" >/dev/null; then
        kind=nightly; channel=nightly
        last=$(bash scripts/select-last-nightly.sh "$WORK/releases.json")
      fi
    fi
    if [[ ${EVENT:-} != workflow_dispatch && $kind != nightly ]] && jq -se --arg sha "$source" 'add | any(.[]; .draft == false and .prerelease == false and .target_commitish == $sha)' "$WORK/releases.json" >/dev/null; then printf 'skip=true\n' >> "$GITHUB_OUTPUT"; exit 0; fi
    rmdir "$PLAN_DIR"
    go run -mod=readonly -p 2 ./cmd/releasectl plan --kind "$kind" --custom "${CUSTOM:-}" --last-nightly "$last" --output "$PLAN_DIR"
    if grep -qx skip=true "$PLAN_DIR/github-output"; then printf 'skip=true\n' >> "$GITHUB_OUTPUT"; exit 0; fi
    NOTARY_STATE_TAG="apple-candidate-v$(jq -r .version "$PLAN_DIR/plan.json")"
    export NOTARY_STATE_TAG
    origin_run=$GITHUB_RUN_ID; origin_attempt=$GITHUB_RUN_ATTEMPT; resume=false
    outputs
    ;;
  persist)
    : "${NOTARY_STATE_TAG:?}" "${RELEASE_CHANNEL:?}"
    source=$(jq -r .sha "$PLAN_DIR/plan.json"); version=$(jq -r .version "$PLAN_DIR/plan.json")
    [[ $(git rev-parse HEAD) == "$source" ]] || fail 'wrong persistence checkout'
    if [[ ! -f $PLAN_DIR/native.json ]]; then
      gh api --paginate "repos/$GH_REPO/actions/runs/$GITHUB_RUN_ID/artifacts?per_page=100" > "$WORK/artifacts.json"
      jq -s '[.[].artifacts[] | select(.name | startswith("unsigned-native-")) | {id,name,digest,expired,workflow_run:{id:.workflow_run.id,head_sha:.workflow_run.head_sha}}] | sort_by(.name)' "$WORK/artifacts.json" > "$PLAN_DIR/native.json"
    fi
    jq -e --arg sha "$source" --argjson run "$GITHUB_RUN_ID" 'length == 4 and all(.[]; .expired == false and .workflow_run.id == $run and .workflow_run.head_sha == $sha and (.id|type == "number" and . > 0 and floor == .) and (.name|test("^unsigned-native-(Linux-X64|Windows-X64|macOS-ARM64|macOS-X64)-" + $sha + "$")) and (.digest|test("^sha256:[0-9a-f]{64}$"))) and (map(.name)|length == (unique|length))' "$PLAN_DIR/native.json" >/dev/null || fail 'expected exact original native artifacts'
    jq -n --arg source "$source" --arg version "$version" --arg repo "$GH_REPO" --arg run "$GITHUB_RUN_ID" --arg attempt "$GITHUB_RUN_ATTEMPT" --arg channel "$RELEASE_CHANNEL" \
      --arg plan "$(hash "$PLAN_DIR/plan.json")" --arg notes "$(hash "$PLAN_DIR/notes.md")" --arg native "$(hash "$PLAN_DIR/native.json")" \
      '{schema:1,source:$source,version:$version,repository:$repo,run_id:$run,run_attempt:$attempt,channel:$channel,plan_sha256:$plan,notes_sha256:$notes,native_sha256:$native}' > "$WORK/candidate.json"
    jq -n --slurpfile candidate "$WORK/candidate.json" --slurpfile plan "$PLAN_DIR/plan.json" --slurpfile native "$PLAN_DIR/native.json" --rawfile notes "$PLAN_DIR/notes.md" '{schema:1,candidate:$candidate[0],plan:$plan[0],native:$native[0],notes:$notes}' > "$WORK/body.json"
    [[ $(wc -c < "$WORK/body.json") -le 60000 ]] || fail 'candidate reservation metadata exceeds supported release body size'
    gh api --paginate "repos/$GH_REPO/releases?per_page=100" | jq -s --arg tag "$NOTARY_STATE_TAG" 'add | map(select(.tag_name == $tag))' > "$WORK/matches.json"
    count=$(jq length "$WORK/matches.json")
    if [[ $count == 0 ]]; then
      jq -n --arg tag "$NOTARY_STATE_TAG" --arg source "$source" --rawfile body "$WORK/body.json" '{tag_name:$tag,target_commitish:$source,name:$tag,body:$body,draft:true,prerelease:true}' > "$WORK/create.json"
      gh api --method POST "repos/$GH_REPO/releases" --input "$WORK/create.json" >/dev/null
    else
      jq -e --arg source "$source" --rawfile body "$WORK/body.json" 'length == 1 and .[0].draft == true and .[0].target_commitish == $source and .[0].body == $body' "$WORK/matches.json" >/dev/null || fail 'candidate reservation collision'
    fi
    for file in candidate.json plan.json notes.md native.json; do
      path="$PLAN_DIR/$file"; [[ $file != candidate.json ]] || path="$WORK/candidate.json"
      bash "$STATE" put "$file" "$path"
    done
    bash "$STATE" inventory > "$WORK/inventory.json"
    while IFS= read -r row; do
      id=$(jq -r .id <<< "$row"); name=$(jq -r .name <<< "$row"); digest=$(jq -r .digest <<< "$row")
      if jq -e --arg name "$name.zip" 'any(.[]; .name == $name)' "$WORK/inventory.json" >/dev/null; then
        bash "$STATE" get "$name.zip" "$WORK/native.zip"
      else
        gh api "repos/$GH_REPO/actions/artifacts/$id/zip" > "$WORK/native.zip"
      fi
      [[ sha256:$(hash "$WORK/native.zip") == "$digest" ]] || fail 'original artifact digest mismatch'
      bash "$STATE" put "$name.zip" "$WORK/native.zip"
    done < <(jq -c '.[]' "$PLAN_DIR/native.json")
    bash "$STATE" put custody-complete.json "$WORK/candidate.json"
    ;;
  *) fail 'unknown operation' ;;
esac
