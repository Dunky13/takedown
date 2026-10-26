#!/usr/bin/env bash
# Partial initial uploads recover the exact original build, with no Apple calls.
set -euo pipefail
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
cd "$ROOT"
TEST_WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-custody-tests.XXXXXX")
cleanup() {
  result=$?
  if [[ $result != 0 ]]; then
    [[ ! -f ${MOCK_STATE:-}/log ]] || cat "$MOCK_STATE/log" >&2
    [[ ! -f $TEST_WORK/log ]] || cat "$TEST_WORK/log" >&2
  fi
  rm -rf -- "$TEST_WORK"
}
trap cleanup EXIT
mkdir "$TEST_WORK/bin" "$TEST_WORK/originals"
export GH_REPO=fixture/takedown GITHUB_RUN_ID=123 GITHUB_RUN_ATTEMPT=1 RELEASE_CHANNEL=manual NOTARY_STATE_TAG=apple-candidate-v0.0.1
export GITHUB_OUTPUT="$TEST_WORK/output" MOCK_ORIGINALS_MISSING=false MOCK_FAIL_ASSET=''
source_sha=$(git rev-parse HEAD)
export MOCK_REQUIRE_BINARY=${MOCK_REQUIRE_BINARY:-false}
printf '[]\n' > "$TEST_WORK/metadata.json"
number=0
for platform in Linux-X64 Windows-X64 macOS-ARM64 macOS-X64; do
  number=$((number + 1))
  printf 'synthetic original native archive %s\n' "$platform" > "$TEST_WORK/originals/$number.zip"
  digest=$(shasum -a 256 "$TEST_WORK/originals/$number.zip" | awk '{print $1}')
  jq --arg name "unsigned-native-$platform-$source_sha" --arg digest "sha256:$digest" --arg source "$source_sha" --argjson id "$number" '. + [{id:$id,name:$name,digest:$digest,expired:false,workflow_run:{id:123,head_sha:$source}}]' "$TEST_WORK/metadata.json" > "$TEST_WORK/metadata.tmp"
  mv "$TEST_WORK/metadata.tmp" "$TEST_WORK/metadata.json"
done
export MOCK_TEST_WORK=$TEST_WORK MOCK_STATE
cp "$ROOT/scripts/testdata/apple-custody-gh.sh" "$TEST_WORK/bin/gh"
chmod +x "$TEST_WORK/bin/gh"
export PATH="$TEST_WORK/bin:$PATH"
new_case() {
  MOCK_STATE="$TEST_WORK/$1"
  export PLAN_DIR="$MOCK_STATE/plan"
  mkdir -p "$PLAN_DIR" "$MOCK_STATE/assets"
  jq -n --arg source "$source_sha" '{version:"0.0.1",package_version:"0.0.1",sha:$source,skip:false}' > "$PLAN_DIR/plan.json"
  printf 'original notes\n' > "$PLAN_DIR/notes.md"
}
persist() { bash scripts/prepare-apple-release.sh persist > "$MOCK_STATE/log" 2>&1; }
for failed_asset in candidate.json "unsigned-native-macOS-ARM64-$source_sha.zip"; do
  new_case "case-${failed_asset%%.*}"
  MOCK_FAIL_ASSET=$failed_asset
  if persist; then echo 'injected custody failure passed' >&2; exit 1; fi
  jq -e '.body | fromjson | .schema==1 and .candidate.run_id=="123" and (.native|length)==4 and .notes=="original notes\n"' "$MOCK_STATE/release.json" >/dev/null
  [[ ! -f $MOCK_STATE/assets/custody-complete.json ]]
  MOCK_FAIL_ASSET=''
  persist
  [[ $(grep -c '^api --method POST repos/fixture/takedown/releases --input' "$MOCK_STATE/calls") == 1 && -f $MOCK_STATE/assets/custody-complete.json ]]
  cmp "$MOCK_STATE/assets/candidate.json" "$MOCK_STATE/assets/custody-complete.json"
  for platform in Linux-X64 Windows-X64 macOS-ARM64 macOS-X64; do [[ -f $MOCK_STATE/assets/unsigned-native-$platform-$source_sha.zip ]]; done
  # Completed durable custody must not depend on expired Actions bytes.
  MOCK_ORIGINALS_MISSING=true
  persist
  MOCK_ORIGINALS_MISSING=false
  [[ $(grep -c '^api --method POST repos/fixture/takedown/releases --input' "$MOCK_STATE/calls") == 1 ]]
  jq '.version="9.9.9"' "$PLAN_DIR/plan.json" > "$MOCK_STATE/changed-plan.json"
  mv "$MOCK_STATE/changed-plan.json" "$PLAN_DIR/plan.json"
  if persist; then echo 'changed candidate reused a reservation' >&2; exit 1; fi
  [[ -f $MOCK_STATE/assets/custody-complete.json ]]
done
# A successful POST is authoritative even when the release list remains stale.
new_case stale-create
MOCK_STALE_CREATE=true persist
[[ -f $MOCK_STATE/assets/custody-complete.json ]]
# A successful upload is verified using its receipt ID, without an asset relist.
new_case stale-upload
MOCK_STALE_UPLOAD=true persist
[[ -f $MOCK_STATE/assets/custody-complete.json ]]
[[ $(grep -c 'releases/7/assets?per_page=100' "$MOCK_STATE/calls") == 1 ]]
# Neither a create receipt nor an exact read may change the sealed reservation.
for filter in '.id=8' '.draft=false' '.tag_name="apple-candidate-v9.9.9"' '.target_commitish="bad"' '.body="changed"'; do
  new_case "create-reject-$(find "$TEST_WORK" -type d | wc -l | tr -d ' ')"
  if MOCK_CREATE_FILTER="$filter" persist; then echo 'mismatched create receipt passed' >&2; exit 1; fi
  [[ -z $(find "$MOCK_STATE/assets" -type f -print) ]]
done
new_case duplicate-reservation
MOCK_FAIL_ASSET=candidate.json persist && { echo 'injected first upload did not fail' >&2; exit 1; }
if MOCK_DUPLICATE_RELEASE=true persist; then echo 'ambiguous reservation passed' >&2; exit 1; fi
[[ -z $(find "$MOCK_STATE/assets" -type f -print) ]]
new_case changed-get
if MOCK_GET_FILTER='.body="changed"' persist; then echo 'changed exact release passed' >&2; exit 1; fi
[[ -z $(find "$MOCK_STATE/assets" -type f -print) ]]
new_case binary-downloads
MOCK_REQUIRE_BINARY=true persist
[[ -f $MOCK_STATE/assets/custody-complete.json ]]
# A repaired automation descendant may retain the original candidate without
# rewriting any original bytes. Dirty-tree simulation is confined to this mock.
export MOCK_REAL_GIT
MOCK_REAL_GIT=$(command -v git)
cat > "$TEST_WORK/bin/git" <<'MOCK'
#!/usr/bin/env bash
if [[ $1 == status && $2 == --porcelain ]]; then exit 0; fi
exec "$MOCK_REAL_GIT" "$@"
MOCK
chmod +x "$TEST_WORK/bin/git"
export MOCK_AUTOMATION_SHA RELEASE_AUTOMATION_SHA RELEASE_WORKFLOW_SHA
MOCK_AUTOMATION_SHA=$(git rev-parse HEAD)
RELEASE_AUTOMATION_SHA=$MOCK_AUTOMATION_SHA; RELEASE_WORKFLOW_SHA=$MOCK_AUTOMATION_SHA
original_source=$source_sha
source_sha=$(git rev-parse HEAD^)
jq --arg old "$original_source" --arg source "$source_sha" 'map(.name |= sub($old; $source) | .workflow_run.head_sha=$source)' "$TEST_WORK/metadata.json" > "$TEST_WORK/metadata.tmp"
mv "$TEST_WORK/metadata.tmp" "$TEST_WORK/metadata.json"
new_case original-with-repaired-automation
persist
cp "$MOCK_STATE/assets/candidate.json" "$MOCK_STATE/original-candidate.json"
jq -e --arg source "$source_sha" --arg automation "$RELEASE_AUTOMATION_SHA" '.source==$source and .automation==$automation' "$MOCK_STATE/assets/automation-$RELEASE_AUTOMATION_SHA.json" >/dev/null
persist
cmp "$MOCK_STATE/original-candidate.json" "$MOCK_STATE/assets/candidate.json"
# Requalification is required on every resume; a spoof or failed check refuses.
if MOCK_CHECK=failure persist; then echo 'unqualified automation resumed custody' >&2; exit 1; fi
cmp "$MOCK_STATE/original-candidate.json" "$MOCK_STATE/assets/candidate.json"
echo 'PASS: synthetic reservation-before-first-asset and partial-archive failures resume original custody without rebuild, overwrite or Apple access'
