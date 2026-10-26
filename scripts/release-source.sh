#!/usr/bin/env bash
# Validate product provenance separately from explicitly pinned release automation.
set -euo pipefail
fail() { printf 'release-source: %s\n' "$*" >&2; exit 1; }
mode=${1:?expected verify or qualify}; source=${2:?expected product or qualification SHA}
[[ $source =~ ^[0-9a-f]{40}$ ]] || fail 'invalid source SHA'
if [[ $mode == verify ]]; then
  if [[ -z ${RELEASE_AUTOMATION_SHA:-} ]]; then
    [[ $(git rev-parse HEAD) == "$source" ]] || fail 'checkout must be at the exact product source SHA'
    exit 0
  fi
  [[ $RELEASE_AUTOMATION_SHA =~ ^[0-9a-f]{40}$ && $RELEASE_AUTOMATION_SHA == "${RELEASE_WORKFLOW_SHA:-}" ]] || fail 'automation must match the explicit workflow SHA'
  [[ $(git rev-parse HEAD) == "$RELEASE_AUTOMATION_SHA" && -z $(git status --porcelain) ]] || fail 'automation checkout must be clean at its pinned SHA'
  git merge-base --is-ancestor "$source" "$RELEASE_AUTOMATION_SHA" || fail 'product source is not an ancestor of automation'
  source=$RELEASE_AUTOMATION_SHA
elif [[ $mode != qualify ]]; then fail 'unknown operation'; fi
: "${GH_REPO:?}"
[[ $GH_REPO =~ ^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$ ]] || fail 'invalid repository'
WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-source.XXXXXX")
trap 'rm -rf -- "$WORK"' EXIT
# The exact automation must belong to protected main, including on continuation.
gh api "repos/$GH_REPO/branches/main" > "$WORK/main.json"
jq -e '.protected == true' "$WORK/main.json" >/dev/null || fail 'main must remain protected'
gh api "repos/$GH_REPO/compare/$source...main" > "$WORK/compare.json"
jq -e --arg source "$source" '.base_commit.sha == $source and .merge_base_commit.sha == $source and (.status == "ahead" or .status == "identical")' "$WORK/compare.json" >/dev/null || fail 'qualified source is not on main'
gh api --paginate "repos/$GH_REPO/commits/$source/check-runs?per_page=100" > "$WORK/checks.json"
for workflow in ci.yml conventional-commits.yml; do
  gh api --paginate "repos/$GH_REPO/actions/workflows/$workflow/runs?branch=main&event=push&head_sha=$source&per_page=100" > "$WORK/runs.json"
  jq -s '[.[].workflow_runs[]] | sort_by(.run_started_at,.run_attempt) | last' "$WORK/runs.json" > "$WORK/run.json"
  jq -e --arg source "$source" --arg repo "$GH_REPO" --arg path ".github/workflows/$workflow" '
    .head_sha == $source and .head_branch == "main" and .event == "push" and
    .repository.full_name == $repo and .head_repository.full_name == $repo and .path == $path and
    .status == "completed" and .conclusion == "success" and
    (.id | type == "number" and . > 0 and floor == .) and
    (.check_suite_id | type == "number" and . > 0 and floor == .)
  ' "$WORK/run.json" >/dev/null || fail "trusted main workflow has not succeeded: $workflow"
  # Trigger freshness applies when planning a new product, never when prioritizing
  # an existing candidate with an older source/trigger.
  if [[ ${RELEASE_REQUIRE_TRIGGER:-false} == true && ${EVENT:-} == workflow_run && $workflow == ci.yml ]]; then
    [[ ${TRIGGER_RUN_ID:-} == "$(jq -r .id "$WORK/run.json")" ]] || fail 'completion is not the current trusted native run'
  fi
  suite=$(jq -r .check_suite_id "$WORK/run.json")
  check=ci-required; [[ $workflow != conventional-commits.yml ]] || check=conventional-commits
  jq -se --arg source "$source" --arg name "$check" --argjson suite "$suite" '
    [.[].check_runs[] | select(.name == $name and .app.id == 15368 and .check_suite.id == $suite and .head_sha == $source)] |
    sort_by(.started_at) | last | .status == "completed" and .conclusion == "success"
  ' "$WORK/checks.json" >/dev/null || fail "required check is not successful in its trusted run suite: $check"
done
