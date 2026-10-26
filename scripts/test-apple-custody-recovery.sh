#!/usr/bin/env bash
# Partial initial uploads recover the exact original build, with no Apple calls.
set -euo pipefail
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
cd "$ROOT"
TEST_WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-custody-tests.XXXXXX")
trap 'rm -rf -- "$TEST_WORK"' EXIT
mkdir "$TEST_WORK/bin" "$TEST_WORK/originals"
export GH_REPO=fixture/takedown GITHUB_RUN_ID=123 GITHUB_RUN_ATTEMPT=1 RELEASE_CHANNEL=manual NOTARY_STATE_TAG=apple-candidate-v0.0.1
export GITHUB_OUTPUT="$TEST_WORK/output" MOCK_ORIGINALS_MISSING=false MOCK_FAIL_ASSET=''
source_sha=$(git rev-parse HEAD)
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
cat > "$TEST_WORK/bin/gh" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
printf '%s\n' "$*" >> "$MOCK_STATE/calls"
if [[ $1 == api && $2 == --method ]]; then
  [[ $3 == POST && $4 == */releases && ! -f $MOCK_STATE/release.json ]]
  jq '. + {id:7}' "$6" > "$MOCK_STATE/release.json"
  cat "$MOCK_STATE/release.json"
elif [[ $1 == api ]]; then
  endpoint=${!#}
  case $endpoint in
    *'/releases?per_page=100') if [[ -f $MOCK_STATE/release.json ]]; then jq -s . "$MOCK_STATE/release.json"; else echo '[]'; fi ;;
    */releases/7/assets\?per_page=100)
      find "$MOCK_STATE/assets" -type f | LC_ALL=C sort | while IFS= read -r file; do name=${file##*/}; jq -n --arg name "$name" '{id:$name,name:$name}'; done | jq -s . ;;
    */releases/assets/*) cat "$MOCK_STATE/assets/${endpoint##*/}" ;;
    */actions/runs/123/artifacts\?per_page=100) jq '{artifacts:.}' "$MOCK_TEST_WORK/metadata.json" ;;
    */actions/artifacts/*/zip)
      [[ $MOCK_ORIGINALS_MISSING != true ]] || { echo 'Not Found (HTTP 404)' >&2; exit 1; }
      id=${endpoint%/zip}; id=${id##*/}
      cat "$MOCK_TEST_WORK/originals/$id.zip" ;;
    *) echo "unexpected API $endpoint" >&2; exit 99 ;;
  esac
elif [[ $1 == release && $2 == upload ]]; then
  file=$4; name=${file##*/}
  [[ $name != "$MOCK_FAIL_ASSET" ]] || exit 88
  [[ ! -e $MOCK_STATE/assets/$name ]]
  cp "$file" "$MOCK_STATE/assets/$name"
else exit 99
fi
MOCK
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
  [[ $(grep -c '^api --method POST' "$MOCK_STATE/calls") == 1 && -f $MOCK_STATE/assets/custody-complete.json ]]
  cmp "$MOCK_STATE/assets/candidate.json" "$MOCK_STATE/assets/custody-complete.json"
  for platform in Linux-X64 Windows-X64 macOS-ARM64 macOS-X64; do [[ -f $MOCK_STATE/assets/unsigned-native-$platform-$source_sha.zip ]]; done
  # Completed durable custody must not depend on expired Actions bytes.
  MOCK_ORIGINALS_MISSING=true
  persist
  MOCK_ORIGINALS_MISSING=false
  [[ $(grep -c '^api --method POST' "$MOCK_STATE/calls") == 1 ]]
  jq '.version="9.9.9"' "$PLAN_DIR/plan.json" > "$MOCK_STATE/changed-plan.json"
  mv "$MOCK_STATE/changed-plan.json" "$PLAN_DIR/plan.json"
  if persist; then echo 'changed candidate reused a reservation' >&2; exit 1; fi
  [[ -f $MOCK_STATE/assets/custody-complete.json ]]
done
echo 'PASS: synthetic reservation-before-first-asset and partial-archive failures resume original custody without rebuild, overwrite or Apple access'
