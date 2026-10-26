#!/usr/bin/env bash
# Offline tests of the real handoff verifier; no credentials, builds or Apple API.
set -euo pipefail
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
TEST_ROOT=$(mktemp -d "${TMPDIR:-/tmp}/takedown-signing-input.XXXXXX")
trap 'rm -rf -- "$TEST_ROOT"' EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
mkdir -p "$TEST_ROOT/tools" "$TEST_ROOT/app/takedown.app/Contents/MacOS" "$TEST_ROOT/outer/darwin-arm64-local"
export INPUT_SOURCE_SHA=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
export GITHUB_SHA=$INPUT_SOURCE_SHA GITHUB_REPOSITORY=fixture/takedown
export GITHUB_RUN_ID=123 GITHUB_RUN_ATTEMPT=1 VERSION=0.0.1 RELEASE_VERSION=0.0.1
export NOTARY_ARCH=arm64 TEST_ROOT
cat > "$TEST_ROOT/app/takedown.app/Contents/Info.plist" <<'PLIST'
CFBundleIdentifier=io.github.dunky13.takedown
CFBundleExecutable=takedown
CFBundlePackageType=APPL
CFBundleShortVersionString=0.0.1
CFBundleVersion=0.0.1
PLIST
cat > "$TEST_ROOT/app/takedown.app/Contents/MacOS/takedown" <<'BINARY'
#!/usr/bin/env bash
[[ $1 == --version ]] || exit 1
printf 'takedown version %s\n' "${FIXTURE_APP_VERSION:-0.0.1}"
BINARY
chmod +x "$TEST_ROOT/app/takedown.app/Contents/MacOS/takedown"
(cd "$TEST_ROOT/app" && zip -qr "$TEST_ROOT/outer/darwin-arm64-local/takedown-0.0.1-macos-arm64-local.zip" takedown.app)
printf 'unsigned dmg fixture\n' > "$TEST_ROOT/outer/darwin-arm64-local/takedown-0.0.1-macos-arm64-local.dmg"
(cd "$TEST_ROOT/outer/darwin-arm64-local" && shasum -a 256 ./*.zip ./*.dmg | sed 's@  \./@  @' > SHA256SUMS)
(cd "$TEST_ROOT/outer" && zip -qr "$TEST_ROOT/artifact.zip" darwin-arm64-local)
export FIXTURE_ARTIFACT_SHA
FIXTURE_ARTIFACT_SHA=$(shasum -a 256 "$TEST_ROOT/artifact.zip" | awk '{print $1}')
cat > "$TEST_ROOT/tools/gh" <<'GH'
#!/usr/bin/env bash
set -euo pipefail
case ${!#} in
  */artifacts\?per_page=100)
    jq -n --arg source "$INPUT_SOURCE_SHA" --arg digest "$FIXTURE_ARTIFACT_SHA" \
      --arg mode "${FIXTURE_ARTIFACT_CASE:-valid}" '
      {id:456,name:("unsigned-native-macOS-ARM64-"+$source),expired:false,
       digest:("sha256:"+$digest),workflow_run:{id:123,head_sha:$source}} as $artifact |
      {artifacts:(if $mode == "missing" then []
        elif $mode == "duplicate" then [$artifact,$artifact]
        elif $mode == "expired" then [$artifact | .expired=true]
        elif $mode == "wrong-source" then [$artifact | .workflow_run.head_sha="bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"]
        elif $mode == "wrong-run" then [$artifact | .workflow_run.id=124]
        elif $mode == "bad-digest" then [$artifact | .digest="sha256:bad"]
        else [$artifact] end)}'
    ;;
  */artifacts/456/zip)
    if [[ ${FIXTURE_ARTIFACT_CASE:-valid} == corrupt ]]; then printf 'corrupt archive\n'
    else cat "$TEST_ROOT/artifact.zip"
    fi
    ;;
  *) printf 'unexpected gh invocation\n' >&2; exit 1 ;;
esac
GH
cat > "$TEST_ROOT/tools/git" <<'GIT'
#!/usr/bin/env bash
case "$1 $2" in
  'rev-parse HEAD') printf '%s\n' "${FIXTURE_CHECKOUT_SHA:-$INPUT_SOURCE_SHA}" ;;
  'status --porcelain') [[ ${FIXTURE_CHECKOUT_DIRTY:-false} != true ]] || printf ' M tracked-file\n' ;;
  *) exit 1 ;;
esac
exit 0
GIT
cat > "$TEST_ROOT/tools/ditto" <<'DITTO'
#!/usr/bin/env bash
[[ $1 == -x && $2 == -k ]] || exit 1
exec unzip -q "$3" -d "$4"
DITTO
cat > "$TEST_ROOT/tools/plutil" <<'PLUTIL'
#!/usr/bin/env bash
if [[ ${FIXTURE_METADATA_FIELD:-} == "$2" ]]; then echo incorrect
else awk -F= -v field="$2" '$1 == field { print $2 }' "${!#}"
fi
PLUTIL
cat > "$TEST_ROOT/tools/lipo" <<'LIPO'
#!/usr/bin/env bash
printf '%s\n' "${FIXTURE_SLICE:-arm64}"
LIPO
cat > "$TEST_ROOT/tools/go" <<'GO'
#!/usr/bin/env bash
[[ $1 == version && $2 == -m ]] || exit 1
printf '%s: go1.26.10\n\tbuild\tvcs.revision=%s\n\tbuild\tvcs.modified=%s\n' \
  "$3" "${FIXTURE_BINARY_SOURCE:-$INPUT_SOURCE_SHA}" "${FIXTURE_BINARY_DIRTY:-false}"
GO
chmod +x "$TEST_ROOT/tools/"*
export PATH="$TEST_ROOT/tools:$PATH"
case_number=0
run_case() {
  case_number=$((case_number + 1))
  export NOTARY_STAGING_DIR="$TEST_ROOT/stage-$case_number" GITHUB_OUTPUT="$TEST_ROOT/output-$case_number"
  bash "$ROOT/scripts/prepare-macos-signing.sh"
}
reject() {
  local expected=$1
  shift
  if run_case > "$TEST_ROOT/failure.log" 2>&1; then echo "expected rejection: $expected" >&2; exit 1; fi
  grep -Fq "$expected" "$TEST_ROOT/failure.log" || { cat "$TEST_ROOT/failure.log"; exit 1; }
  [[ ! -s $GITHUB_OUTPUT ]] || { echo 'failed verifier exported signing inputs' >&2; exit 1; }
}
run_case
jq -e --arg source "$INPUT_SOURCE_SHA" '.source==$source and .artifact_id=="456" and
  .architecture=="arm64" and .bundle_id=="io.github.dunky13.takedown" and
  (.binary_sha256|test("^[0-9a-f]{64}$"))' "$NOTARY_STAGING_DIR/input.json" >/dev/null
grep -Fxq "binary=$NOTARY_STAGING_DIR/extracted/takedown.app/Contents/MacOS/takedown" "$GITHUB_OUTPUT"
for mode in missing duplicate expired wrong-source wrong-run bad-digest; do
  export FIXTURE_ARTIFACT_CASE=$mode
  reject 'expected one unexpired artifact bound to this source and run'
done
export FIXTURE_ARTIFACT_CASE=corrupt
reject 'artifact archive digest mismatch'
unset FIXTURE_ARTIFACT_CASE
export FIXTURE_CHECKOUT_DIRTY=true
reject 'checkout must be clean'
unset FIXTURE_CHECKOUT_DIRTY
export FIXTURE_CHECKOUT_SHA=bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
reject 'checkout must be clean'
unset FIXTURE_CHECKOUT_SHA
export FIXTURE_SLICE=x86_64
reject 'binary architecture does not match target'
unset FIXTURE_SLICE
for field in CFBundleIdentifier CFBundleExecutable CFBundlePackageType CFBundleShortVersionString CFBundleVersion; do
  export FIXTURE_METADATA_FIELD=$field
  reject "application metadata mismatch: $field"
done
unset FIXTURE_METADATA_FIELD
export FIXTURE_BINARY_SOURCE=bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
reject 'binary source revision mismatch'
unset FIXTURE_BINARY_SOURCE
export FIXTURE_BINARY_DIRTY=true
reject 'binary was not built from clean source'
unset FIXTURE_BINARY_DIRTY
export FIXTURE_APP_VERSION=0.0.2
reject 'binary release version mismatch'
unset FIXTURE_APP_VERSION
# A correct artifact-level digest cannot hide corrupt native package bytes.
printf corrupt >> "$TEST_ROOT/outer/darwin-arm64-local/takedown-0.0.1-macos-arm64-local.dmg"
(cd "$TEST_ROOT/outer" && zip -qr "$TEST_ROOT/artifact.zip" darwin-arm64-local)
FIXTURE_ARTIFACT_SHA=$(shasum -a 256 "$TEST_ROOT/artifact.zip" | awk '{print $1}')
reject 'native package checksum mismatch'
printf 'unsigned dmg fixture\n' > "$TEST_ROOT/outer/darwin-arm64-local/takedown-0.0.1-macos-arm64-local.dmg"
rebuild_archive() {
  (cd "$TEST_ROOT/outer/darwin-arm64-local" && shasum -a 256 ./*.zip ./*.dmg | sed 's@  \./@  @' > SHA256SUMS)
  (cd "$TEST_ROOT/outer" && zip -qr "$TEST_ROOT/artifact.zip" darwin-arm64-local)
  FIXTURE_ARTIFACT_SHA=$(shasum -a 256 "$TEST_ROOT/artifact.zip" | awk '{print $1}')
}
printf 'must never escape staging\n' > "$TEST_ROOT/outside"
(cd "$TEST_ROOT/app" && zip -q "$TEST_ROOT/outer/darwin-arm64-local/takedown-0.0.1-macos-arm64-local.zip" ../outside)
rebuild_archive
reject 'unexpected application archive path'
rm "$TEST_ROOT/outer/darwin-arm64-local/takedown-0.0.1-macos-arm64-local.zip"
ln -s /fixture/outside "$TEST_ROOT/app/takedown.app/Contents/link"
(cd "$TEST_ROOT/app" && zip -qry "$TEST_ROOT/outer/darwin-arm64-local/takedown-0.0.1-macos-arm64-local.zip" takedown.app)
rebuild_archive
reject 'application archive must not contain symbolic links'
printf 'PASS: immutable artifact handoff, clean source, checksums, bundle identity, architecture and version (%s scenarios)\n' "$case_number"
