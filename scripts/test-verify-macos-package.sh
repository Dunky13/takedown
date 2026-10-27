#!/usr/bin/env bash
# Synthetic final-package acceptance; no signing, mounts, account access or build.
set -euo pipefail
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
TEST_WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-acceptance-tests.XXXXXX")
cleanup_tests() {
  local status=$?
  if [[ $status != 0 && -f ${MOCK_STATE:-}/output ]]; then
    printf 'FAIL: scenario=%s\n' "${MOCK_STATE:-setup}" >&2
    cat "$MOCK_STATE/output" >&2
  fi
  rm -rf -- "$TEST_WORK"
}
trap cleanup_tests EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
mkdir -p "$TEST_WORK/tools" "$TEST_WORK/scratch" "$TEST_WORK/repo/scripts"
cp "$ROOT/scripts/verify-macos-package.sh" "$ROOT/scripts/macos-ticket-policy.sh" "$ROOT/scripts/continue-apple-release.sh" "$ROOT/scripts/apple-stages-verified.sh" "$ROOT/scripts/apple-acceptance-ready.sh" "$ROOT/scripts/fresh-apple-acceptance.sh" "$TEST_WORK/repo/scripts/"
export REAL_RM
REAL_RM=$(command -v rm)
export MOCK_NATIVE_OUTPUT="$ROOT/scripts/testdata/stapler-cloudkit-lost.txt"
cat > "$TEST_WORK/tools/uname" <<'MOCK'
#!/usr/bin/env bash
echo Darwin
MOCK
cat > "$TEST_WORK/tools/plutil" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
[[ $1 == -extract && $3 == raw && $4 == -o && $5 == - ]] || exit 99
jq -er --arg key "$2" '.[$key] // empty' "$6"
MOCK
cat > "$TEST_WORK/tools/lipo" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
[[ $1 == -archs ]] || exit 99
jq -r .fixtureArch "$(dirname -- "$2")/../Info.plist"
MOCK
cat > "$TEST_WORK/tools/codesign" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
if [[ $1 == --display ]]; then
  target=${!#}
  team=$APPLE_TEAM_ID
  if [[ $target == *.app ]]; then team=$(jq -r .fixtureTeam "$target/Contents/Info.plist"); fi
  printf 'TeamIdentifier=%s\nAuthority=Developer ID Application: Fixture\n' "$team" >&2
elif [[ $1 == --verify ]]; then
  [[ $MOCK_SCENARIO != codesign-failure || ${!#} != */Applications/* ]] || { echo 'fixture codesign failure' >&2; exit 1; }
else exit 99
fi
MOCK
cat > "$TEST_WORK/tools/ditto" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
if [[ $1 == -x && $2 == -k ]]; then
  cp -R "$MOCK_STATE/zip-app" "$4/takedown.app"
else
  [[ $# == 2 ]] || exit 99
  cp -R "$1" "$2"
fi
MOCK
cat > "$TEST_WORK/tools/hdiutil" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
case $1 in
  attach)
    printf 'attach\n' >> "$MOCK_STATE/mount-calls"
    mountpoint=''
    while (($#)); do
      if [[ $1 == -mountpoint ]]; then mountpoint=$2; break; fi
      shift
    done
    [[ -n $mountpoint ]]
    cp -R "$MOCK_STATE/dmg-app" "$mountpoint/takedown.app"
    [[ $MOCK_SCENARIO != attach-failure ]] || { echo 'fixture attach failure' >&2; exit 1; }
    ;;
  detach)
    printf 'detach\n' >> "$MOCK_STATE/mount-calls"
    if [[ $MOCK_SCENARIO == detach-once && ! -f $MOCK_STATE/detach-failed ]]; then
      touch "$MOCK_STATE/detach-failed"
      echo 'fixture detach failure' >&2
      exit 1
    fi
    [[ $MOCK_SCENARIO != detach-failure ]] || { echo 'fixture detach failure' >&2; exit 1; }
    ;;
  *) exit 99 ;;
esac
MOCK
cat > "$TEST_WORK/tools/xcrun" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
[[ $1 == stapler && $2 == validate ]] || exit 99
printf 'stapler %s\n' "$3" >> "$MOCK_STATE/trust-calls"
target=zip
[[ $3 != *.dmg ]] || target=dmg
[[ $3 != */Applications/* ]] || target=copied
if [[ $target == "${MOCK_FAIL_TARGET:-zip}" ]]; then
  if [[ $MOCK_SCENARIO == cloudkit ]]; then cat "$MOCK_NATIVE_OUTPUT"; exit "${MOCK_STAPLER_EXIT:-68}"; fi
  if [[ $MOCK_SCENARIO == interrupted ]]; then kill -TERM "$PPID"; exit 1; fi
fi
[[ $MOCK_SCENARIO != stapler-failure || $3 != */Applications/* ]] || { echo 'fixture stapler failure' >&2; exit 1; }
MOCK
cat > "$TEST_WORK/tools/spctl" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
[[ $1 == --assess ]] || exit 99
printf 'spctl %s\n' "${!#}" >> "$MOCK_STATE/trust-calls"
[[ $MOCK_SCENARIO != gatekeeper-failure || ${!#} != */Applications/* ]] || { echo 'fixture Gatekeeper failure' >&2; exit 1; }
MOCK
cat > "$TEST_WORK/tools/rm" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
if [[ ${MOCK_SCENARIO:-} == cleanup-failure && $* == *takedown-acceptance.* && ! -f $MOCK_STATE/cleanup-failed ]]; then
  touch "$MOCK_STATE/cleanup-failed"
  exit 19
fi
exec "$REAL_RM" "$@"
MOCK
cat > "$TEST_WORK/tools/xattr" <<'MOCK'
#!/usr/bin/env bash
printf 'QUARANTINE REMOVAL FORBIDDEN\n' >> "$MOCK_STATE/forbidden"
exit 99
MOCK
chmod +x "$TEST_WORK/tools/"*
export PATH="$TEST_WORK/tools:$PATH" TMPDIR="$TEST_WORK/scratch"
export VERSION=0.0.1 RELEASE_VERSION=0.0.1 NOTARY_ARCH=arm64 APPLE_TEAM_ID=ABCDEFGHIJ
export NOTARY_SOURCE_SHA=1111111111111111111111111111111111111111 GITHUB_RUN_ID=123 GITHUB_RUN_ATTEMPT=1
export MOCK_STATE MOCK_SCENARIO NOTARY_EVIDENCE_DIR
PACKAGES="$TEST_WORK/repo/dist/packages/darwin-arm64-release"
BASE=takedown-0.0.1-macos-arm64-release
new_case() {
  MOCK_SCENARIO=$1
  MOCK_STATE="$TEST_WORK/$1"
  NOTARY_EVIDENCE_DIR="$MOCK_STATE/evidence"
  mkdir -p "$MOCK_STATE/zip-app/Contents/MacOS" "$NOTARY_EVIDENCE_DIR/app" "$NOTARY_EVIDENCE_DIR/dmg"
  printf '{"synthetic":true}\n' > "$NOTARY_EVIDENCE_DIR/input.json"
  printf '{"state":"verified","team_id":"ABCDEFGHIJ"}\n' > "$NOTARY_EVIDENCE_DIR/app/receipt.json"
  cp "$NOTARY_EVIDENCE_DIR/app/receipt.json" "$NOTARY_EVIDENCE_DIR/dmg/receipt.json"
  printf '#!/bin/sh\nprintf "takedown version 0.0.1\\n"\n' > "$MOCK_STATE/zip-app/Contents/MacOS/takedown"
  chmod +x "$MOCK_STATE/zip-app/Contents/MacOS/takedown"
  jq -n --arg arch "${FIXTURE_SLICE:-arm64}" '{CFBundleIdentifier:"io.github.dunky13.takedown",CFBundleExecutable:"takedown",CFBundleShortVersionString:"0.0.1",CFBundleVersion:"0.0.1",fixtureArch:$arch,fixtureTeam:"ABCDEFGHIJ"}' > "$MOCK_STATE/zip-app/Contents/Info.plist"
  cp -R "$MOCK_STATE/zip-app" "$MOCK_STATE/dmg-app"
  rm -rf "$PACKAGES"
  mkdir -p "$PACKAGES"
  printf 'synthetic final ZIP\n' > "$PACKAGES/$BASE.zip"
  printf 'synthetic final DMG\n' > "$PACKAGES/$BASE.dmg"
  (cd "$PACKAGES" && shasum -a 256 "$BASE.zip" "$BASE.dmg" > SHA256SUMS)
  case $1 in
    wrong-copied-version) jq '.CFBundleShortVersionString="0.0.2"' "$MOCK_STATE/dmg-app/Contents/Info.plist" > "$MOCK_STATE/changed"; mv "$MOCK_STATE/changed" "$MOCK_STATE/dmg-app/Contents/Info.plist" ;;
    wrong-copied-team) jq '.fixtureTeam="ZZZZZZZZZZ"' "$MOCK_STATE/dmg-app/Contents/Info.plist" > "$MOCK_STATE/changed"; mv "$MOCK_STATE/changed" "$MOCK_STATE/dmg-app/Contents/Info.plist" ;;
    wrong-copied-architecture) jq '.fixtureArch="x86_64"' "$MOCK_STATE/dmg-app/Contents/Info.plist" > "$MOCK_STATE/changed"; mv "$MOCK_STATE/changed" "$MOCK_STATE/dmg-app/Contents/Info.plist" ;;
    wrong-copied-executable-version) printf '#!/bin/sh\necho "takedown version 0.0.2"\n' > "$MOCK_STATE/dmg-app/Contents/MacOS/takedown" ;;
    executable-mismatch) printf '# different signed candidate\n' >> "$MOCK_STATE/dmg-app/Contents/MacOS/takedown" ;;
    checksum-corruption) printf 'corrupted\n' >> "$PACKAGES/$BASE.zip" ;;
    manifest-extra) printf '%064d  unexpected.zip\n' 0 >> "$PACKAGES/SHA256SUMS" ;;
  esac
}
verify() { bash "$TEST_WORK/repo/scripts/verify-macos-package.sh" > "$MOCK_STATE/output" 2>&1; }
assert_cleanup() {
  [[ ! -f $MOCK_STATE/forbidden ]]
  [[ -z $(find "$TMPDIR" -mindepth 1 -maxdepth 1 -type d -name 'takedown-acceptance.*' -print) ]]
  if [[ -f $MOCK_STATE/mount-calls ]]; then
    awk 'NR%2 == 1 && $0 != "attach" {exit 1} NR%2 == 0 && $0 != "detach" {exit 1} END {if (NR%2) exit 1}' "$MOCK_STATE/mount-calls"
  fi
}
# Faithful native output through the actual continuation -> final verifier seam.
new_case cloudkit
mkdir -p "$NOTARY_EVIDENCE_DIR/final" "$NOTARY_EVIDENCE_DIR/app" "$NOTARY_EVIDENCE_DIR/dmg"
cp "$PACKAGES/"* "$NOTARY_EVIDENCE_DIR/final/"
printf '{"state":"verified","team_id":"ABCDEFGHIJ"}\n' > "$NOTARY_EVIDENCE_DIR/app/receipt.json"
cp "$NOTARY_EVIDENCE_DIR/app/receipt.json" "$NOTARY_EVIDENCE_DIR/dmg/receipt.json"
printf '{"synthetic":true}\n' > "$NOTARY_EVIDENCE_DIR/input.json"
printf '{"stale":true}\n' > "$NOTARY_EVIDENCE_DIR/acceptance.json"
cat > "$TEST_WORK/repo/scripts/apple-checkpoint.sh" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
cp "$1/acceptance-state.json" "$MOCK_STATE/checkpoint-state.json"
[[ ${MOCK_CHECKPOINT_FAIL:-false} != true ]] || exit 88
MOCK
cat > "$TEST_WORK/repo/scripts/package-macos.sh" <<'MOCK'
#!/usr/bin/env bash
echo 'forbidden repack' >&2
exit 99
MOCK
export GITHUB_STEP_SUMMARY="$MOCK_STATE/summary"
if ! bash "$TEST_WORK/repo/scripts/continue-apple-release.sh" > "$MOCK_STATE/output" 2>&1; then
  cat "$MOCK_STATE/output"
  echo 'FAIL: native CloudKit transport must checkpoint pending through retained-final continuation' >&2
  exit 1
fi
jq -e '.state=="pending" and .attempts==1 and .failure_class=="cloudkit-transport"' "$MOCK_STATE/checkpoint-state.json" >/dev/null
[[ ! -f $NOTARY_EVIDENCE_DIR/acceptance.json ]]
assert_cleanup
echo 'PASS: native CloudKit final reuse checkpoints pending without repacking'
# Resume exact finals successfully; the budget cannot reset after success.
MOCK_SCENARIO=success
bash "$TEST_WORK/repo/scripts/continue-apple-release.sh" > "$MOCK_STATE/output" 2>&1
jq -e '.state=="verified" and .attempts==2 and .failure_class==""' "$NOTARY_EVIDENCE_DIR/acceptance-state.json" >/dev/null
bash "$TEST_WORK/repo/scripts/apple-acceptance-ready.sh" "$NOTARY_EVIDENCE_DIR"
cp "$NOTARY_EVIDENCE_DIR/acceptance.json" "$MOCK_STATE/good-acceptance.json"
  MOCK_SCENARIO=cloudkit
  if bash "$TEST_WORK/repo/scripts/continue-apple-release.sh" > "$MOCK_STATE/output" 2>&1; then
    # A verified ordinary recheck first consumes its last attempt on failure.
    echo 'exhausted final validation passed' >&2; exit 1
  fi
  jq -e '.state=="blocked" and .attempts==3 and .failure_class=="cloudkit-transport"' "$NOTARY_EVIDENCE_DIR/acceptance-state.json" >/dev/null
calls=$(wc -l < "$MOCK_STATE/trust-calls")
MOCK_SCENARIO=success
if bash "$TEST_WORK/repo/scripts/continue-apple-release.sh" >> "$MOCK_STATE/output" 2>&1; then exit 1; fi
[[ $(wc -l < "$MOCK_STATE/trust-calls") == "$calls" ]]
cp "$MOCK_STATE/good-acceptance.json" "$NOTARY_EVIDENCE_DIR/acceptance.json"
if bash "$TEST_WORK/repo/scripts/apple-acceptance-ready.sh" "$NOTARY_EVIDENCE_DIR"; then exit 1; fi
echo 'PASS: exact final recovery succeeds, preserves spent budget, and terminal exhaustion rejects stale success before native checks'

# Save/restore a pending checkpoint at the fresh-runner wrapper seam.
cat > "$TEST_WORK/repo/scripts/apple-release-state.sh" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
[[ $1 == restore && ! -e $2 ]]
cp -R "$MOCK_RESTORE" "$2"
MOCK
export MOCK_RESTORE MOCK_FAIL_TARGET MOCK_CHECKPOINT_FAIL
seed_native() {
  new_case "$1"
  mkdir -p "$NOTARY_EVIDENCE_DIR/final"
  cp "$PACKAGES/"* "$NOTARY_EVIDENCE_DIR/final/"
  printf '{"stale":true}\n' > "$NOTARY_EVIDENCE_DIR/acceptance.json"
  GITHUB_STEP_SUMMARY="$MOCK_STATE/summary"
}
continue_native() { bash "$TEST_WORK/repo/scripts/continue-apple-release.sh" > "$MOCK_STATE/output" 2>&1; }
fresh_native() {
  MOCK_RESTORE="$NOTARY_EVIDENCE_DIR"
  NOTARY_EVIDENCE_DIR="$MOCK_STATE/fresh-$1" ORIGIN_RUN=123 ORIGIN_ATTEMPT=1 GITHUB_RUN_ID=777 GITHUB_RUN_ATTEMPT=2 \
    bash "$TEST_WORK/repo/scripts/fresh-apple-acceptance.sh" > "$MOCK_STATE/output" 2>&1
}
for variant in arm64:zip arm64:dmg arm64:copied amd64:zip amd64:dmg amd64:copied; do
  NOTARY_ARCH=${variant%%:*}
  point=${variant#*:}
  FIXTURE_SLICE=arm64
  [[ $NOTARY_ARCH != amd64 ]] || FIXTURE_SLICE=x86_64
  PACKAGES="$TEST_WORK/repo/dist/packages/darwin-$NOTARY_ARCH-release"
  BASE="takedown-0.0.1-macos-$NOTARY_ARCH-release"
  seed_native "fresh-$NOTARY_ARCH-$point"
  MOCK_FAIL_TARGET=$point MOCK_SCENARIO=cloudkit
  fresh_native pending
  jq -e '.state=="pending" and .attempts==1' "$MOCK_STATE/fresh-pending/acceptance-state.json" >/dev/null
  [[ ! -f $MOCK_STATE/fresh-pending/acceptance.json && ! -f $MOCK_STATE/fresh-pending/fresh-acceptance-$NOTARY_ARCH.json ]]
  # A new scheduled signing continuation resumes this exact checkpoint.
  cp "$MOCK_STATE/fresh-pending/acceptance-state.json" "$NOTARY_EVIDENCE_DIR/acceptance-state.json"
  rm "$NOTARY_EVIDENCE_DIR/acceptance.json"
  MOCK_SCENARIO=success
  continue_native
  fresh_native success
  jq -e '.credential_free_runner==true and .ci_run_id=="123" and .verification_run_id=="777" and .verification_run_attempt=="2"' "$MOCK_STATE/fresh-success/fresh-acceptance-$NOTARY_ARCH.json" >/dev/null
  jq -e '.state=="verified" and .attempts==2' "$MOCK_STATE/fresh-success/acceptance-state.json" >/dev/null
  assert_cleanup
done
echo 'PASS: ZIP, DMG and copied-app native transport failures persist on credential-free runners and fully recheck exact finals before fresh success'
NOTARY_ARCH=arm64 FIXTURE_SLICE=arm64
PACKAGES="$TEST_WORK/repo/dist/packages/darwin-arm64-release"
BASE=takedown-0.0.1-macos-arm64-release
MOCK_FAIL_TARGET=zip
seed_native repeated
MOCK_SCENARIO=cloudkit
for count in 1 2; do
  continue_native
  jq -e --argjson count "$count" '.state=="pending" and .attempts==$count' "$NOTARY_EVIDENCE_DIR/acceptance-state.json" >/dev/null
done
if continue_native; then exit 1; fi
jq -e '.state=="blocked" and .attempts==3' "$NOTARY_EVIDENCE_DIR/acceptance-state.json" >/dev/null
assert_cleanup
seed_native interrupted
MOCK_SCENARIO=interrupted
if continue_native; then exit 1; fi
jq -e '.state=="pending" and .attempts==0' "$NOTARY_EVIDENCE_DIR/acceptance-state.json" >/dev/null
if continue_native; then exit 1; fi
jq -e '.state=="pending" and .attempts==1' "$NOTARY_EVIDENCE_DIR/acceptance-state.json" >/dev/null
MOCK_SCENARIO=success
continue_native
jq -e '.state=="verified" and .attempts==2' "$NOTARY_EVIDENCE_DIR/acceptance-state.json" >/dev/null
assert_cleanup
seed_native interrupted-last
MOCK_SCENARIO=cloudkit
continue_native
continue_native
MOCK_SCENARIO=interrupted
if continue_native; then exit 1; fi
jq -e '.state=="pending" and .attempts==3' "$NOTARY_EVIDENCE_DIR/acceptance-state.json" >/dev/null
calls=$(wc -l < "$MOCK_STATE/trust-calls")
MOCK_SCENARIO=success
if continue_native; then exit 1; fi
[[ $(wc -l < "$MOCK_STATE/trust-calls") == "$calls" ]]
jq -e '.state=="blocked" and .attempts==3' "$NOTARY_EVIDENCE_DIR/acceptance-state.json" >/dev/null
seed_native checkpoint-failure
MOCK_SCENARIO=success MOCK_CHECKPOINT_FAIL=true
if continue_native; then exit 1; fi
[[ ! -f $MOCK_STATE/trust-calls && ! -f $NOTARY_EVIDENCE_DIR/acceptance.json ]]
MOCK_CHECKPOINT_FAIL=false
continue_native
jq -e '.state=="verified" and .attempts==1' "$NOTARY_EVIDENCE_DIR/acceptance-state.json" >/dev/null
seed_native cleanup-failure
MOCK_SCENARIO=cleanup-failure
if continue_native; then exit 1; fi
[[ ! -f $NOTARY_EVIDENCE_DIR/acceptance.json ]]
jq -e '.state=="blocked"' "$NOTARY_EVIDENCE_DIR/acceptance-state.json" >/dev/null
assert_cleanup
echo 'PASS: finite failure and interruption budgets, pre-native checkpoint refusal, last-attempt interruption, and cleanup refusal'

# A failed final detach must retain ownership for EXIT cleanup. Even when that
# cleanup succeeds, the original detach failure cannot issue acceptance.
seed_native detach-once
MOCK_SCENARIO=detach-once
if continue_native; then exit 1; fi
[[ $(cat "$MOCK_STATE/mount-calls") == $'attach\ndetach\ndetach' ]]
[[ ! -f $NOTARY_EVIDENCE_DIR/acceptance.json ]]
jq -e '.state=="blocked"' "$NOTARY_EVIDENCE_DIR/acceptance-state.json" >/dev/null
[[ -z $(find "$TMPDIR" -mindepth 1 -maxdepth 1 -type d -name 'takedown-acceptance.*' -print) ]]
echo 'PASS: final detach failure retains mount ownership, retries cleanup, and cannot accept'

# Binding changes cannot reset pending recovery or cause any native operation.
for field in source_commit version architecture ci_run_id ci_run_attempt team_id zip_sha256 dmg_sha256 input_sha256 app_receipt_sha256 dmg_receipt_sha256; do
  seed_native "binding-$field"
  MOCK_SCENARIO=cloudkit
  continue_native
  calls=$(wc -l < "$MOCK_STATE/trust-calls")
  jq --arg key "$field" '.[$key]="wrong"' "$NOTARY_EVIDENCE_DIR/acceptance-state.json" > "$MOCK_STATE/change.json"
  mv "$MOCK_STATE/change.json" "$NOTARY_EVIDENCE_DIR/acceptance-state.json"
  MOCK_SCENARIO=success
  if continue_native; then exit 1; fi
  [[ $(wc -l < "$MOCK_STATE/trust-calls") == "$calls" && ! -f $NOTARY_EVIDENCE_DIR/acceptance.json ]]
done
seed_native missing-finals
MOCK_SCENARIO=cloudkit
continue_native
rm "$NOTARY_EVIDENCE_DIR/final/SHA256SUMS"
if continue_native; then exit 1; fi
grep -q 'no repack permitted' "$MOCK_STATE/output"
echo 'PASS: all candidate/archive/input/stage-receipt binding mismatches refuse native recovery; missing bound finals cannot repack'
MOCK_FAIL_TARGET=zip
original_native=$MOCK_NATIVE_OUTPUT
export MOCK_STAPLER_EXIT
for negative in rejected invalid signature ticket-absent other-endpoint other-domain other-code unknown68 wrong-exit; do
  seed_native "negative-$negative"
  MOCK_NATIVE_OUTPUT="$MOCK_STATE/diagnostic.txt"
  cp "$original_native" "$MOCK_NATIVE_OUTPUT"
  case $negative in
    rejected) echo 'Rejected' >> "$MOCK_NATIVE_OUTPUT" ;;
    invalid) echo 'Invalid ticket' >> "$MOCK_NATIVE_OUTPUT" ;;
    signature) echo 'Signature mismatch' >> "$MOCK_NATIVE_OUTPUT" ;;
    ticket-absent) echo 'Could not find ticket' >> "$MOCK_NATIVE_OUTPUT" ;;
    other-endpoint) sed 's@/records/lookup,@/records/lookup.evil,@g' "$original_native" > "$MOCK_NATIVE_OUTPUT" ;;
    other-domain) sed 's@api.apple-cloudkit.com@elsewhere.invalid@g' "$original_native" > "$MOCK_NATIVE_OUTPUT" ;;
    other-code) sed 's/Code=-1005/Code=-1001/g' "$original_native" > "$MOCK_NATIVE_OUTPUT" ;;
    unknown68) echo 'The validate action failed! Error 68.' > "$MOCK_NATIVE_OUTPUT" ;;
    wrong-exit) MOCK_STAPLER_EXIT=1 ;;
  esac
  MOCK_SCENARIO=cloudkit
  if continue_native; then exit 1; fi
  jq -e '.state=="blocked" and .failure_class=="unclassified"' "$NOTARY_EVIDENCE_DIR/acceptance-state.json" >/dev/null
  [[ ! -e $NOTARY_EVIDENCE_DIR/acceptance.json ]]
  calls=$(wc -l < "$MOCK_STATE/trust-calls")
  MOCK_SCENARIO=success
  if continue_native; then exit 1; fi
  [[ $(wc -l < "$MOCK_STATE/trust-calls") == "$calls" ]]
  MOCK_STAPLER_EXIT=68
  assert_cleanup
done
MOCK_NATIVE_OUTPUT=$original_native
echo 'PASS: native contradictory/absent-ticket/unknown/other-endpoint/other-code/non-68 diagnostics remain terminal'

[[ ${ACCEPTANCE_TEST_GROUP:-all} != cloudkit ]] || exit 0
new_case success
verify || { cat "$MOCK_STATE/output"; exit 1; }
assert_cleanup
jq -e '.version=="0.0.1" and .architecture=="arm64" and (.zip_sha256|test("^[0-9a-f]{64}$")) and (.dmg_sha256|test("^[0-9a-f]{64}$")) and (.checks|index("ZIP/DMG executable equality")!=null)' "$NOTARY_EVIDENCE_DIR/acceptance.json" >/dev/null
grep -q 'stapler .*/Applications/takedown.app' "$MOCK_STATE/trust-calls"
grep -q 'spctl .*/Applications/takedown.app' "$MOCK_STATE/trust-calls"
echo 'PASS: checksums, ZIP and copied DMG app acceptance, evidence, mount cleanup, quarantine preserved'
for scenario in wrong-copied-version wrong-copied-team wrong-copied-architecture wrong-copied-executable-version checksum-corruption manifest-extra executable-mismatch gatekeeper-failure stapler-failure codesign-failure attach-failure detach-failure; do
  new_case "$scenario"
  if verify; then cat "$MOCK_STATE/output"; echo "expected rejection: $scenario" >&2; exit 1; fi
  case $scenario in
    wrong-copied-version) expected='wrong bundle version' ;;
    wrong-copied-team) expected='wrong signed team' ;;
    wrong-copied-architecture) expected='wrong binary slices' ;;
    wrong-copied-executable-version) expected='wrong executable version' ;;
    checksum-corruption) expected='FAILED' ;;
    manifest-extra) expected='manifest must cover exactly' ;;
    executable-mismatch) expected='ZIP and installed DMG executables differ' ;;
    gatekeeper-failure) expected='fixture Gatekeeper failure' ;;
    stapler-failure) expected='ticket validation failed permanently' ;;
    codesign-failure) expected='fixture codesign failure' ;;
    attach-failure) expected='fixture attach failure' ;;
    detach-failure) expected='fixture detach failure' ;;
  esac
  grep -Fq "$expected" "$MOCK_STATE/output" || { cat "$MOCK_STATE/output"; echo "wrong rejection for $scenario" >&2; exit 1; }
  if [[ $scenario == detach-failure ]]; then
    [[ $(cat "$MOCK_STATE/mount-calls") == $'attach\ndetach\ndetach' ]]
    [[ ! -f $NOTARY_EVIDENCE_DIR/acceptance.json ]]
    # The mock mount cannot detach. Keep its directory until explicit test
    # teardown, proving production cleanup never recurses into that mount.
    retained_work=$(find "$TMPDIR" -mindepth 1 -maxdepth 1 -type d -name 'takedown-acceptance.*' -print)
    [[ -n $retained_work && -d $retained_work/mount/takedown.app ]]
    "$REAL_RM" -rf -- "$retained_work"
  else
    assert_cleanup
  fi
  if [[ $scenario == checksum-corruption || $scenario == manifest-extra ]]; then
    [[ ! -f $MOCK_STATE/mount-calls ]]
  else
    [[ -f $MOCK_STATE/mount-calls ]]
  fi
  [[ ! -e $NOTARY_EVIDENCE_DIR/acceptance.json ]]
  echo "PASS: $scenario rejected, owned mount cleanup attempted, quarantine preserved"
done
echo 'PASS: synthetic final-package acceptance contracts; no Apple, real mounts or signing'
