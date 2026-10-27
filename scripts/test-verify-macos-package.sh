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
export REAL_RM REAL_JQ REAL_MV MOCK_WRITE_FAILURE MOCK_DETACH_FAILURE
REAL_RM=$(command -v rm)
REAL_JQ=$(command -v jq)
REAL_MV=$(command -v mv)
MOCK_WRITE_FAILURE='' MOCK_DETACH_FAILURE=false
export MOCK_NATIVE_OUTPUT="$ROOT/scripts/testdata/stapler-cloudkit-lost.txt"
cat > "$TEST_WORK/tools/jq" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
if [[ ( ${MOCK_WRITE_FAILURE:-} == receipt-write && ${!#} == *'checks:['* ) ||
      ( ${MOCK_WRITE_FAILURE:-} == reservation-write && $* == *'--arg state pending'* && $* == *'--argjson attempts 4'* ) ||
      ( ${MOCK_WRITE_FAILURE:-} == terminal-write && $* == *'--arg state blocked'* ) ]]; then
  printf '{partial'
  exit 17
fi
exec "$REAL_JQ" "$@"
MOCK
cat > "$TEST_WORK/tools/mv" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
if [[ ${MOCK_WRITE_FAILURE:-} == receipt-rename && $1 == */.acceptance.tmp ]]; then exit 18; fi
if [[ ${MOCK_WRITE_FAILURE:-} == reservation-rename && $1 == */.acceptance-state.tmp ]] &&
  grep -q '"attempts": 4' "$1" && grep -q '"state": "pending"' "$1"; then exit 18; fi
if [[ ${MOCK_WRITE_FAILURE:-} == terminal-rename && $1 == */.acceptance-state.tmp ]] &&
  grep -q '"state": "blocked"' "$1"; then exit 19; fi
exec "$REAL_MV" "$@"
MOCK
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
  [[ $MOCK_SCENARIO == missing-timestamp ]] || echo 'Timestamp=Oct 27, 2026 at 12:00:00' >&2
  [[ $MOCK_SCENARIO == missing-runtime ]] || echo 'CodeDirectory v=20500 size=100 flags=0x10000(runtime)' >&2
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
    [[ $MOCK_SCENARIO != detach-failure && ${MOCK_DETACH_FAILURE:-false} != true ]] || { echo 'fixture detach failure' >&2; exit 1; }
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
  if [[ $MOCK_SCENARIO == cancelled-late ]]; then : > "${NOTARY_ACCEPTANCE_CANCEL_FILE:?}"; fi
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
if [[ ${MOCK_SCENARIO:-} == peer-cleanup-failure && $* == *takedown-counterpart.* && ! -f $MOCK_STATE/peer-cleanup-failed ]]; then
  touch "$MOCK_STATE/peer-cleanup-failed"
  exit 19
fi
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
state=$(jq -r .state "$1/acceptance-state.json")
printf '%s\n' "$state" >> "$MOCK_STATE/checkpoint-attempts"
number=$(wc -l < "$MOCK_STATE/checkpoint-attempts" | tr -d ' ')
[[ ${MOCK_CHECKPOINT_FAIL_NUMBER:-0} != "$number" ]] || exit 88
[[ ${MOCK_CHECKPOINT_FAIL:-false} != true ]] || exit 88
[[ ${MOCK_CHECKPOINT_FAIL:-false} != verified || $state != verified ]] || exit 88
cp "$1/acceptance-state.json" "$MOCK_STATE/checkpoint-state.json"
cp "$1/acceptance-state.json" "$MOCK_STATE/checkpoint-$number.json"
if [[ ${MOCK_CHECKPOINT_LATE_CANCEL:-false} == true && $state == verified ]]; then
  : > "${NOTARY_ACCEPTANCE_CANCEL_FILE:?}"
  cp -R "$1" "$MOCK_STATE/durable-late"
fi
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
if [[ $NOTARY_ARCH == amd64 && -n ${MOCK_PEER_RESTORE:-} ]]; then cp -R "$MOCK_PEER_RESTORE" "$2"; else cp -R "$MOCK_RESTORE" "$2"; fi
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
write_boundary_cases() {
  local mode
  for mode in final-checkpoint receipt-write receipt-rename terminal-write terminal-rename; do
    seed_native "write-$mode"
    MOCK_SCENARIO=success MOCK_WRITE_FAILURE='' MOCK_CHECKPOINT_FAIL=false
    case $mode in
      final-checkpoint) MOCK_CHECKPOINT_FAIL=verified ;;
      terminal-*) MOCK_WRITE_FAILURE=$mode MOCK_SCENARIO=codesign-failure ;;
      *) MOCK_WRITE_FAILURE=$mode ;;
    esac
    if continue_native; then echo "unexpected acceptance after $mode failure" >&2; exit 1; fi
    MOCK_WRITE_FAILURE='' MOCK_CHECKPOINT_FAIL=false
    [[ ! -f $NOTARY_EVIDENCE_DIR/acceptance.json && ! -f $NOTARY_EVIDENCE_DIR/fresh-acceptance-$NOTARY_ARCH.json ]]
    if [[ $mode == terminal-* ]]; then
      # Conditional cleanup must not move a failed/partial state write or
      # acknowledge an old pending receipt after a failed terminal rename.
      jq -e '.state=="pending"' "$NOTARY_EVIDENCE_DIR/acceptance-state.json" >/dev/null
      [[ $(wc -l < "$MOCK_STATE/checkpoint-attempts") -eq 1 ]]
    elif [[ $mode == final-checkpoint ]]; then
      jq -e '.state=="pending"' "$MOCK_STATE/checkpoint-state.json" >/dev/null
      if bash "$TEST_WORK/repo/scripts/apple-acceptance-ready.sh" "$NOTARY_EVIDENCE_DIR"; then exit 1; fi
    else
      jq -e '.state=="blocked"' "$NOTARY_EVIDENCE_DIR/acceptance-state.json" >/dev/null
    fi
    assert_cleanup
    echo "PASS: $mode failure cannot issue acceptance or acknowledge a failed state transition"
  done
}
write_boundary_cases
[[ ${ACCEPTANCE_TEST_GROUP:-all} != write-boundaries ]] || exit 0

# Exact-receipt exception through the real producer -> fresh consumer -> native
# verifier seam. Only the policy's sealed digest is replaced in this PRIVATE
# script copy to authorize synthetic bytes. Production has no test override.
exception_cases() {
  local original_pin=27b7a083cc181ed6def7479ec99b25479d7b11d9c3de36bbbd9facfbacfc60fe stage extension uuid sha parent_pin case_name calls field mutation
  export NOTARY_ACCEPTANCE_REQUALIFY_SHA256 NOTARY_RELEASE_ID=425202908 GH_REPO=Dunky13/takedown MOCK_PEER_RESTORE MOCK_CHECKPOINT_FAIL_NUMBER=0 NOTARY_ACCEPTANCE_CLEANUP_HOOK=''
  NOTARY_SOURCE_SHA=ef82381d509aac8b5b02da0f784a8176b1c5b9c0 GITHUB_RUN_ID=40470756538
  export ORIGIN_RUN=40470756538 ORIGIN_ATTEMPT=1
  export GITHUB_SHA=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa RELEASE_AUTOMATION_SHA=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa RELEASE_WORKFLOW_SHA=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
  exception_seed() {
    cp "$ROOT/scripts/macos-ticket-policy.sh" "$TEST_WORK/repo/scripts/"
    seed_native "$1"
    MOCK_SCENARIO=success MOCK_CHECKPOINT_FAIL=false MOCK_CHECKPOINT_FAIL_NUMBER=0 MOCK_WRITE_FAILURE=''
    jq -n --arg source "$NOTARY_SOURCE_SHA" --arg arch arm64 --arg run "$GITHUB_RUN_ID" \
      '{source:$source,version:"0.0.1",package_version:"0.0.1",architecture:$arch,run_id:$run,run_attempt:"1",bundle_id:"io.github.dunky13.takedown",repository:"Dunky13/takedown"}' > "$NOTARY_EVIDENCE_DIR/input.json"
    for stage in app dmg; do
      extension=zip uuid=4549ac29-2d27-4cae-b118-d0ce238b1e0c
      [[ $stage != dmg ]] || { extension=dmg; uuid=31f135e5-19a0-4dff-b393-f1cd16b72878; }
      sha=$(shasum -a 256 "$PACKAGES/$BASE.$extension" | awk '{print $1}')
      jq -n --arg stage "$stage" --arg sha "$sha" --arg uuid "$uuid" --arg source "$NOTARY_SOURCE_SHA" --arg run "$GITHUB_RUN_ID" --arg team "$APPLE_TEAM_ID" \
        '{schema:2,state:"verified",source_commit:$source,version:"0.0.1",architecture:"arm64",ci_run_id:$run,ci_run_attempt:"1",team_id:$team,stage:$stage,submission_id:$uuid,submitted_sha256:$sha,filename:("takedown-0.0.1-macos-arm64-"+$source[0:12]+"-"+$run+"-1-"+$stage+(if $stage=="app" then ".zip" else ".dmg" end)),verified_filename:(if $stage=="app" then "stapled-app.zip" else "stapled.dmg" end),verified_sha256:$sha,signing_identity:("Developer ID Application: Marc Went ("+$team+")"),ticket_attempts:1,validation_attempts:0,failure_class:"",failure_exit:"",failure_phase:""}' > "$NOTARY_EVIDENCE_DIR/$stage/receipt.json"
      jq '{status:"Accepted",jobId:.submission_id,archiveFilename:.filename,sha256:.submitted_sha256}' "$NOTARY_EVIDENCE_DIR/$stage/receipt.json" > "$NOTARY_EVIDENCE_DIR/$stage/apple-log.json"
      cp "$PACKAGES/$BASE.$extension" "$NOTARY_EVIDENCE_DIR/$stage/$(jq -r .verified_filename "$NOTARY_EVIDENCE_DIR/$stage/receipt.json")"
    done
    # Seed an acknowledged pre-terminal state inside this synthetic fixture.
    # Production refuses a missing state for this already-spent immutable build.
    jq -n --arg source "$NOTARY_SOURCE_SHA" --arg run "$GITHUB_RUN_ID" --arg team "$APPLE_TEAM_ID" \
      --arg zip "$(shasum -a 256 "$PACKAGES/$BASE.zip" | awk '{print $1}')" --arg dmg "$(shasum -a 256 "$PACKAGES/$BASE.dmg" | awk '{print $1}')" \
      --arg input "$(shasum -a 256 "$NOTARY_EVIDENCE_DIR/input.json" | awk '{print $1}')" \
      --arg app_receipt "$(shasum -a 256 "$NOTARY_EVIDENCE_DIR/app/receipt.json" | awk '{print $1}')" \
      --arg dmg_receipt "$(shasum -a 256 "$NOTARY_EVIDENCE_DIR/dmg/receipt.json" | awk '{print $1}')" \
      '{schema:1,state:"pending",attempts:0,failure_class:"",source_commit:$source,version:"0.0.1",architecture:"arm64",ci_run_id:$run,ci_run_attempt:"1",team_id:$team,
        zip_sha256:$zip,dmg_sha256:$dmg,input_sha256:$input,app_receipt_sha256:$app_receipt,dmg_receipt_sha256:$dmg_receipt}' > "$NOTARY_EVIDENCE_DIR/acceptance-state.json"
    verify
    jq '.state="blocked" | .attempts=3 | .failure_class="cloudkit-transport"' "$NOTARY_EVIDENCE_DIR/acceptance-state.json" > "$MOCK_STATE/parent.json"
    cp "$MOCK_STATE/parent.json" "$NOTARY_EVIDENCE_DIR/acceptance-state.json"
    parent_pin=$(shasum -a 256 "$MOCK_STATE/parent.json" | awk '{print $1}')
    sed "s/$original_pin/$parent_pin/g" "$ROOT/scripts/macos-ticket-policy.sh" > "$TEST_WORK/repo/scripts/macos-ticket-policy.sh"
    for file in verify-macos-package.sh fresh-apple-acceptance.sh; do
      sed "s/$original_pin/$parent_pin/g" "$ROOT/scripts/$file" > "$TEST_WORK/repo/scripts/$file"
    done
    sed "s/$original_pin/$parent_pin/g" "$ROOT/scripts/verify-apple-fresh-completion.sh" > "$TEST_WORK/repo/scripts/verify-apple-fresh-completion.sh"
    NOTARY_ACCEPTANCE_REQUALIFY_SHA256=$parent_pin
    MOCK_PEER_RESTORE="$MOCK_STATE/peer"
    cp -R "$NOTARY_EVIDENCE_DIR" "$MOCK_PEER_RESTORE"
    jq '.architecture="amd64"' "$MOCK_PEER_RESTORE/input.json" > "$MOCK_STATE/change"; mv "$MOCK_STATE/change" "$MOCK_PEER_RESTORE/input.json"
    for stage in app dmg; do
      uuid=88951ab7-f126-485b-945c-a11fdbdee902
      [[ $stage != dmg ]] || uuid=70d28837-3ac6-4de3-86b8-a07c1813af89
      jq --arg uuid "$uuid" '.architecture="amd64" | .submission_id=$uuid | .filename |= sub("arm64";"amd64")' "$MOCK_PEER_RESTORE/$stage/receipt.json" > "$MOCK_STATE/change"; mv "$MOCK_STATE/change" "$MOCK_PEER_RESTORE/$stage/receipt.json"
      jq '{status:"Accepted",jobId:.submission_id,archiveFilename:.filename,sha256:.submitted_sha256}' "$MOCK_PEER_RESTORE/$stage/receipt.json" > "$MOCK_PEER_RESTORE/$stage/apple-log.json"
    done
    for extension in zip dmg; do mv "$MOCK_PEER_RESTORE/final/$BASE.$extension" "$MOCK_PEER_RESTORE/final/${BASE/arm64/amd64}.$extension"; done
    sed 's/arm64/amd64/g' "$NOTARY_EVIDENCE_DIR/final/SHA256SUMS" > "$MOCK_PEER_RESTORE/final/SHA256SUMS"
    jq --arg input "$(shasum -a 256 "$MOCK_PEER_RESTORE/input.json" | awk '{print $1}')" --arg app "$(shasum -a 256 "$MOCK_PEER_RESTORE/app/receipt.json" | awk '{print $1}')" --arg dmg "$(shasum -a 256 "$MOCK_PEER_RESTORE/dmg/receipt.json" | awk '{print $1}')" \
      '.architecture="amd64" | .state="verified" | .attempts=0 | .failure_class="" | .input_sha256=$input | .app_receipt_sha256=$app | .dmg_receipt_sha256=$dmg' "$MOCK_STATE/parent.json" > "$MOCK_PEER_RESTORE/acceptance-state.json"
    jq '.architecture="amd64"' "$NOTARY_EVIDENCE_DIR/acceptance.json" > "$MOCK_PEER_RESTORE/acceptance.json"
    rm -f "$MOCK_STATE/trust-calls" "$MOCK_STATE/mount-calls" "$MOCK_STATE/checkpoint-attempts" "$MOCK_STATE/checkpoint-state.json" "$MOCK_STATE/"checkpoint-[0-9]*.json
  }
  exception_fresh() {
    MOCK_RESTORE="$NOTARY_EVIDENCE_DIR"
    NOTARY_EVIDENCE_DIR="$MOCK_STATE/fresh-$1" GITHUB_RUN_ID=${MOCK_FRESH_RUN:-777} GITHUB_RUN_ATTEMPT=2 \
      bash "$TEST_WORK/repo/scripts/fresh-apple-acceptance.sh" > "$MOCK_STATE/output" 2>&1
  }
  exception_seed exception-success
  continue_native
  [[ ! -f $MOCK_STATE/trust-calls && ! -f $NOTARY_EVIDENCE_DIR/acceptance.json ]]
  cmp "$MOCK_STATE/parent.json" "$NOTARY_EVIDENCE_DIR/acceptance-state.json"
  cmp "$MOCK_STATE/parent.json" "$NOTARY_EVIDENCE_DIR/blocked-acceptance-$parent_pin.json"
  exception_fresh success
  bash "$ROOT/scripts/test-apple-fresh-completion.sh" --emit-fixture "$MOCK_STATE/fresh-success/acceptance.json" "$MOCK_STATE/completion-provider"
  export FIXTURE="$MOCK_STATE/completion-provider"
  export NOTARY_RELEASE_SOURCE_SHA=$NOTARY_SOURCE_SHA NOTARY_STATE_TAG=apple-candidate-v0.0.1
  NOTARY_RELEASE_BODY_SHA256=$(jq -r .NOTARY_RELEASE_BODY_SHA256 "$FIXTURE/pins.json")
  export NOTARY_RELEASE_BODY_SHA256
  cat > "$TEST_WORK/tools/gh" <<'MOCK'
#!/usr/bin/env bash
exec "$FIXTURE/bin/gh" "$@"
MOCK
  chmod +x "$TEST_WORK/tools/gh"
  jq -e '.state=="blocked" and .attempts==3' "$MOCK_STATE/checkpoint-2.json" >/dev/null
  jq -e '.state=="pending" and .attempts==4' "$MOCK_STATE/checkpoint-3.json" >/dev/null
  jq -e '.state=="verified" and .attempts==4' "$MOCK_STATE/checkpoint-4.json" >/dev/null
  cmp "$MOCK_STATE/parent.json" "$MOCK_STATE/fresh-success/blocked-acceptance-$parent_pin.json"
  cmp "$NOTARY_EVIDENCE_DIR/final/$BASE.zip" "$PACKAGES/$BASE.zip"
  cmp "$NOTARY_EVIDENCE_DIR/final/$BASE.dmg" "$PACKAGES/$BASE.dmg"
  jq -e '.credential_free_runner==true and .ci_run_id=="40470756538" and .verification_run_id=="777"' "$MOCK_STATE/fresh-success/fresh-acceptance-arm64.json" >/dev/null
  bash "$TEST_WORK/repo/scripts/apple-acceptance-ready.sh" "$MOCK_STATE/fresh-success"
  for mutation in missing-parent changed-parent changed-final missing-state missing-state-and-parent wrong-counter changed-binding blocked pending stale-signer missing-fresh-proof; do
    cp -R "$MOCK_STATE/fresh-success" "$MOCK_STATE/audit-$mutation"
    case $mutation in
      missing-parent) rm "$MOCK_STATE/audit-$mutation/blocked-acceptance-$parent_pin.json" ;;
      changed-parent) echo ' ' >> "$MOCK_STATE/audit-$mutation/blocked-acceptance-$parent_pin.json" ;;
      changed-final) echo changed >> "$MOCK_STATE/audit-$mutation/final/$BASE.zip" ;;
      missing-state) rm "$MOCK_STATE/audit-$mutation/acceptance-state.json" ;;
      missing-state-and-parent) rm "$MOCK_STATE/audit-$mutation/acceptance-state.json" "$MOCK_STATE/audit-$mutation/blocked-acceptance-$parent_pin.json" ;;
      stale-signer) cp "$NOTARY_EVIDENCE_DIR/acceptance-state.json" "$MOCK_STATE/audit-$mutation/acceptance.json" ;;
      missing-fresh-proof) jq 'del(.credential_free_runner,.verification_run_id,.verification_run_attempt)' "$MOCK_STATE/audit-$mutation/acceptance.json" > "$MOCK_STATE/change"; mv "$MOCK_STATE/change" "$MOCK_STATE/audit-$mutation/acceptance.json" ;;
      *)
        case $mutation in wrong-counter) field='.attempts=3' ;; changed-binding) field='.ci_run_id="99"' ;; blocked) field='.state="blocked"|.failure_class="cloudkit-transport"' ;; pending) field='.state="pending"' ;; esac
        jq "$field" "$MOCK_STATE/audit-$mutation/acceptance-state.json" > "$MOCK_STATE/change"; mv "$MOCK_STATE/change" "$MOCK_STATE/audit-$mutation/acceptance-state.json" ;;
    esac
    if bash "$TEST_WORK/repo/scripts/apple-acceptance-ready.sh" "$MOCK_STATE/audit-$mutation"; then echo "exception audit passed $mutation" >&2; exit 1; fi
  done
  calls=$(wc -l < "$MOCK_STATE/trust-calls")
  NOTARY_EVIDENCE_DIR="$MOCK_STATE/audit-missing-state"
  if exception_fresh missing-state; then echo 'missing exception state reset fresh budget' >&2; exit 1; fi
  [[ $(wc -l < "$MOCK_STATE/trust-calls") == "$calls" ]]
  for mutation in missing-state-and-parent fresh-marker-only stripped-markers; do
    directory="$MOCK_STATE/wrapper-$mutation"
    cp -R "$MOCK_STATE/fresh-success" "$directory"
    rm -f "$directory/acceptance-state.json" "$directory/blocked-acceptance-$parent_pin.json"
    if [[ $mutation != missing-state-and-parent ]]; then
      jq 'del(.acceptance_attempts,.acceptance_exception_parent_sha256)' "$directory/acceptance.json" > "$MOCK_STATE/change"; mv "$MOCK_STATE/change" "$directory/acceptance.json"
    fi
    if [[ $mutation == stripped-markers ]]; then
      jq 'del(.acceptance_attempts,.acceptance_exception_parent_sha256)' "$directory/fresh-acceptance-arm64.json" > "$MOCK_STATE/change"; mv "$MOCK_STATE/change" "$directory/fresh-acceptance-arm64.json"
    fi
    cp "$directory/acceptance.json" "$MOCK_STATE/receipt-before"
    cp "$directory/fresh-acceptance-arm64.json" "$MOCK_STATE/fresh-before"
    NOTARY_EVIDENCE_DIR=$directory
    if continue_native; then echo "producer reset missing custody $mutation" >&2; exit 1; fi
    if bash "$TEST_WORK/repo/scripts/verify-macos-package.sh" > "$MOCK_STATE/output" 2>&1; then echo "direct verifier reset missing custody $mutation" >&2; exit 1; fi
    if exception_fresh "wrapper-$mutation"; then echo "fresh wrapper reset missing custody $mutation" >&2; exit 1; fi
    cmp "$MOCK_STATE/receipt-before" "$directory/acceptance.json"
    cmp "$MOCK_STATE/fresh-before" "$directory/fresh-acceptance-arm64.json"
    cmp "$MOCK_STATE/receipt-before" "$MOCK_STATE/fresh-wrapper-$mutation/acceptance.json"
    cmp "$MOCK_STATE/fresh-before" "$MOCK_STATE/fresh-wrapper-$mutation/fresh-acceptance-arm64.json"
    [[ ! -e $directory/acceptance-state.json && $(wc -l < "$MOCK_STATE/trust-calls") == "$calls" ]]
  done
  NOTARY_EVIDENCE_DIR="$MOCK_STATE/fresh-success"
  continue_native
  if bash "$TEST_WORK/repo/scripts/verify-macos-package.sh" > "$MOCK_STATE/output" 2>&1; then echo 'verified exception qualified again' >&2; exit 1; fi
  MOCK_FRESH_RUN=778 exception_fresh replay
  cmp "$MOCK_STATE/fresh-success/acceptance.json" "$MOCK_STATE/fresh-replay/fresh-acceptance-arm64.json"
  jq -e '.verification_run_id=="777" and .verification_run_attempt=="2" and .acceptance_attempts==4' "$MOCK_STATE/fresh-replay/fresh-acceptance-arm64.json" >/dev/null
  [[ $(wc -l < "$MOCK_STATE/trust-calls") == "$calls" ]]
  NOTARY_EVIDENCE_DIR="$MOCK_STATE/audit-missing-fresh-proof"
  if exception_fresh stale-proof; then echo 'stale native proof substituted for completed fresh qualification' >&2; exit 1; fi
  [[ $(wc -l < "$MOCK_STATE/trust-calls") == "$calls" ]]
  echo 'PASS: exact approved producer handoff spends zero; fresh blocked3 -> durable reserved4 -> verified4; exact completed fresh proof survives publication retry without a fifth native check'
  for case_name in peer-pending peer-rejected peer-foreign peer-unknown missing-pin malformed-pin wrong-pin foreign-candidate wrong-source wrong-input wrong-uuid wrong-log wrong-schema wrong-identity changed-archive missing-custody changed-parent used-parent unknown-failure credential-present; do
    exception_seed "exception-$case_name"
    case $case_name in
      peer-pending) jq '.state="validation-pending"' "$MOCK_PEER_RESTORE/app/receipt.json" > "$MOCK_STATE/change"; mv "$MOCK_STATE/change" "$MOCK_PEER_RESTORE/app/receipt.json" ;;
      peer-rejected|peer-unknown|peer-foreign)
        case $case_name in peer-rejected) field='.state="rejected"' ;; peer-unknown) field='.state="unknown"' ;; peer-foreign) field='.submission_id="4549ac29-2d27-4cae-b118-d0ce238b1e0c"' ;; esac
        jq "$field" "$MOCK_PEER_RESTORE/app/receipt.json" > "$MOCK_STATE/change"; mv "$MOCK_STATE/change" "$MOCK_PEER_RESTORE/app/receipt.json" ;;
      credential-present) export APPLE_API_KEY_ID=SYNTHETIC-NOT-A-SECRET ;;
      missing-pin) NOTARY_ACCEPTANCE_REQUALIFY_SHA256='' ;;
      malformed-pin) NOTARY_ACCEPTANCE_REQUALIFY_SHA256=invalid ;;
      wrong-pin) NOTARY_ACCEPTANCE_REQUALIFY_SHA256=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa ;;
      foreign-candidate) NOTARY_RELEASE_ID=424598089 ;;
      wrong-source) NOTARY_SOURCE_SHA=2222222222222222222222222222222222222222 ;;
      wrong-input) echo ' ' >> "$NOTARY_EVIDENCE_DIR/input.json" ;;
      wrong-uuid) jq '.submission_id="88951ab7-f126-485b-945c-a11fdbdee902"' "$NOTARY_EVIDENCE_DIR/app/receipt.json" > "$MOCK_STATE/change"; mv "$MOCK_STATE/change" "$NOTARY_EVIDENCE_DIR/app/receipt.json" ;;
      wrong-log) jq '.status="Rejected"' "$NOTARY_EVIDENCE_DIR/app/apple-log.json" > "$MOCK_STATE/change"; mv "$MOCK_STATE/change" "$NOTARY_EVIDENCE_DIR/app/apple-log.json" ;;
      wrong-schema|unknown-failure|used-parent)
        case $case_name in wrong-schema) field='.schema=2' ;; unknown-failure) field='.failure_class="unclassified"' ;; used-parent) field='.acceptance_exception_parent_sha256="bad"' ;; esac
        jq "$field" "$NOTARY_EVIDENCE_DIR/acceptance-state.json" > "$MOCK_STATE/change"; mv "$MOCK_STATE/change" "$NOTARY_EVIDENCE_DIR/acceptance-state.json" ;;
      wrong-identity) jq '.team_id="ZZZZZZZZZZ"' "$NOTARY_EVIDENCE_DIR/dmg/receipt.json" > "$MOCK_STATE/change"; mv "$MOCK_STATE/change" "$NOTARY_EVIDENCE_DIR/dmg/receipt.json" ;;
      changed-archive) echo changed >> "$NOTARY_EVIDENCE_DIR/final/$BASE.zip" ;;
      missing-custody) rm "$NOTARY_EVIDENCE_DIR/app/stapled-app.zip" ;;
      changed-parent) echo '{}' > "$NOTARY_EVIDENCE_DIR/blocked-acceptance-$parent_pin.json" ;;
    esac
    if exception_fresh reject; then [[ $case_name == peer-pending ]] || { echo "exception unexpectedly passed $case_name" >&2; exit 1; }; fi
    [[ ! -f $MOCK_STATE/trust-calls && ! -f $MOCK_STATE/fresh-reject/fresh-acceptance-arm64.json ]]
    if [[ $case_name == peer-pending ]]; then cmp "$MOCK_STATE/parent.json" "$MOCK_STATE/fresh-reject/acceptance-state.json"; [[ ! -f $MOCK_STATE/checkpoint-attempts ]]; fi
    NOTARY_SOURCE_SHA=ef82381d509aac8b5b02da0f784a8176b1c5b9c0 NOTARY_RELEASE_ID=425202908
    unset APPLE_API_KEY_ID
  done
  echo 'PASS: peer-pending preserves unspent3; missing/wrong pin, foreign/unknown/schema/identity/UUID/log/hash/parent/custody refuse before trust calls'
  cat > "$TEST_WORK/repo/scripts/acceptance-cleanup-fixture.sh" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
[[ $1 == cleanup && ! -e $NOTARY_EVIDENCE_DIR/acceptance.json ]]
jq -e '.state=="pending" and .attempts==4' "$NOTARY_EVIDENCE_DIR/acceptance-state.json" >/dev/null
echo cleanup >> "$MOCK_STATE/cleanup-hook"
exit 73
MOCK
  for case_name in first-checkpoint reservation-checkpoint last-checkpoint interrupted late-cancellation cloudkit receipt-write receipt-rename reservation-write reservation-rename terminal-write terminal-rename cleanup cleanup-hook peer-cleanup; do
    exception_seed "exception-$case_name"
    case $case_name in
      first-checkpoint) MOCK_CHECKPOINT_FAIL_NUMBER=1 ;;
      reservation-checkpoint) MOCK_CHECKPOINT_FAIL_NUMBER=2 ;;
      last-checkpoint) MOCK_CHECKPOINT_FAIL_NUMBER=3 ;;
      interrupted|cloudkit) MOCK_SCENARIO=$case_name ;;
      late-cancellation) MOCK_SCENARIO=cancelled-late; export NOTARY_ACCEPTANCE_CANCEL_FILE="$MOCK_STATE/cancelled" ;;
      receipt-*|reservation-write|reservation-rename) MOCK_WRITE_FAILURE=$case_name ;;
      terminal-*) MOCK_WRITE_FAILURE=$case_name MOCK_SCENARIO=codesign-failure ;;
      cleanup) MOCK_SCENARIO=cleanup-failure ;;
      cleanup-hook) NOTARY_ACCEPTANCE_CLEANUP_HOOK="$TEST_WORK/repo/scripts/acceptance-cleanup-fixture.sh" ;;
      peer-cleanup) MOCK_SCENARIO=peer-cleanup-failure ;;
    esac
    if exception_fresh failure; then echo "exception failure passed $case_name" >&2; exit 1; fi
    [[ ! -f $MOCK_STATE/fresh-failure/acceptance.json && ! -f $MOCK_STATE/fresh-failure/fresh-acceptance-arm64.json ]]
    if [[ $case_name == first-checkpoint || $case_name == reservation-checkpoint || $case_name == reservation-write || $case_name == reservation-rename ]]; then [[ ! -f $MOCK_STATE/trust-calls ]]; fi
    if [[ $case_name == cleanup-hook ]]; then [[ -f $MOCK_STATE/cleanup-hook ]]; fi
    if [[ $case_name == late-cancellation ]]; then
      [[ -f $MOCK_STATE/cancelled ]]
      jq -e '.state=="pending" and .attempts==4' "$MOCK_STATE/fresh-failure/acceptance-state.json" >/dev/null
      [[ ! -f $MOCK_STATE/checkpoint-3.json ]] || ! jq -e '.state=="verified"' "$MOCK_STATE/checkpoint-3.json" >/dev/null
    fi
    if [[ $case_name != first-checkpoint ]]; then
      jq -e '.attempts==4' "$MOCK_STATE/fresh-failure/acceptance-state.json" >/dev/null
      calls=0; [[ ! -f $MOCK_STATE/trust-calls ]] || calls=$(wc -l < "$MOCK_STATE/trust-calls")
      NOTARY_EVIDENCE_DIR="$MOCK_STATE/fresh-failure"
      MOCK_SCENARIO=success MOCK_WRITE_FAILURE='' MOCK_CHECKPOINT_FAIL_NUMBER=0
      if exception_fresh retry; then echo 'spent exception retried' >&2; exit 1; fi
      [[ ! -f $MOCK_STATE/trust-calls || $(wc -l < "$MOCK_STATE/trust-calls") == "$calls" ]]
      NOTARY_ACCEPTANCE_REQUALIFY_SHA256=$(shasum -a 256 "$NOTARY_EVIDENCE_DIR/acceptance-state.json" | awk '{print $1}')
      if exception_fresh repin; then echo 'spent exception repinned' >&2; exit 1; fi
    fi
    NOTARY_ACCEPTANCE_CLEANUP_HOOK=''
    unset NOTARY_ACCEPTANCE_CANCEL_FILE
  done
  echo 'PASS: first/reservation/final acknowledgement, receipt/terminal writes/renames, interruption, transport and cleanup fail closed; spent4 never retries or repins'
  exception_seed late-checkpoint-cancellation
  export MOCK_CHECKPOINT_LATE_CANCEL=true NOTARY_ACCEPTANCE_CANCEL_FILE="$MOCK_STATE/cancelled"
  if exception_fresh cancelled-checkpoint; then echo 'cancelled final checkpoint passed' >&2; exit 1; fi
  unset MOCK_CHECKPOINT_LATE_CANCEL NOTARY_ACCEPTANCE_CANCEL_FILE
  [[ ! -f $MOCK_STATE/fresh-cancelled-checkpoint/acceptance.json && -f $MOCK_STATE/durable-late/acceptance.json ]]
  jq -e '.state=="verified" and .attempts==4' "$MOCK_STATE/durable-late/acceptance-state.json" >/dev/null
  bash "$ROOT/scripts/test-apple-fresh-completion.sh" --emit-fixture "$MOCK_STATE/durable-late/acceptance.json" "$MOCK_STATE/cancelled-provider"
  export FIXTURE="$MOCK_STATE/cancelled-provider"
  python3 - "$FIXTURE" <<'PY'
import json, sys
from pathlib import Path
directory = Path(sys.argv[1])
responses = json.loads((directory / 'responses.json').read_text())
key = next(key for key in responses if '/attempts/2/jobs?' in key)
file = directory / responses[key]
data = json.loads(file.read_text())
data['jobs'][0]['conclusion'] = 'cancelled'
file.write_text(json.dumps(data))
PY
  calls=$(wc -l < "$MOCK_STATE/trust-calls")
  NOTARY_EVIDENCE_DIR="$MOCK_STATE/durable-late"
  if bash "$TEST_WORK/repo/scripts/apple-acceptance-ready.sh" "$NOTARY_EVIDENCE_DIR" > "$MOCK_STATE/output" 2>&1; then echo 'cancelled durable proof passed readiness' >&2; exit 1; fi
  if continue_native; then echo 'cancelled durable proof passed producer' >&2; exit 1; fi
  if exception_fresh cancelled-replay; then echo 'cancelled durable proof replayed' >&2; exit 1; fi
  [[ $(wc -l < "$MOCK_STATE/trust-calls") == "$calls" ]]
  if grep -Eq '/artifacts[/?]' "$FIXTURE/calls.jsonl"; then echo 'cancelled job reached artifact proof reads' >&2; exit 1; fi
  echo 'PASS: late committed verified4 cannot pass actual readiness/producer/fresh replay when original hosted job was cancelled; no fifth qualification'
  cp "$ROOT/scripts/macos-ticket-policy.sh" "$TEST_WORK/repo/scripts/"
  cp "$ROOT/scripts/verify-macos-package.sh" "$ROOT/scripts/fresh-apple-acceptance.sh" "$TEST_WORK/repo/scripts/"
  unset NOTARY_ACCEPTANCE_REQUALIFY_SHA256 NOTARY_RELEASE_ID GH_REPO MOCK_PEER_RESTORE MOCK_CHECKPOINT_FAIL_NUMBER ORIGIN_RUN ORIGIN_ATTEMPT NOTARY_ACCEPTANCE_CLEANUP_HOOK FIXTURE NOTARY_RELEASE_SOURCE_SHA NOTARY_STATE_TAG NOTARY_RELEASE_BODY_SHA256 GITHUB_SHA RELEASE_AUTOMATION_SHA RELEASE_WORKFLOW_SHA
  NOTARY_SOURCE_SHA=1111111111111111111111111111111111111111 GITHUB_RUN_ID=123 MOCK_SCENARIO=success MOCK_WRITE_FAILURE='' MOCK_CHECKPOINT_FAIL=false
}
exception_cases
[[ ${ACCEPTANCE_TEST_GROUP:-all} != exception ]] || exit 0
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

seed_native pending-detach-failure
MOCK_SCENARIO=cloudkit MOCK_FAIL_TARGET=copied MOCK_DETACH_FAILURE=true
if continue_native; then exit 1; fi
MOCK_DETACH_FAILURE=false MOCK_FAIL_TARGET=zip
[[ $(cat "$MOCK_STATE/mount-calls") == $'attach\ndetach' ]]
[[ ! -f $NOTARY_EVIDENCE_DIR/acceptance.json ]]
jq -e '.state=="blocked" and .attempts==1' "$MOCK_STATE/checkpoint-state.json" >/dev/null
retained_work=$(find "$TMPDIR" -mindepth 1 -maxdepth 1 -type d -name 'takedown-acceptance.*' -print)
[[ -n $retained_work && -d $retained_work/mount/takedown.app ]]
"$REAL_RM" -rf -- "$retained_work"
echo 'PASS: copied-app pending transport failure becomes terminal when owned detach fails, without recursive mounted-tree cleanup'

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
for field in source_commit version architecture ci_run_id ci_run_attempt team_id zip_sha256 dmg_sha256 input_sha256 app_receipt_sha256 dmg_receipt_sha256; do
  seed_native "foreign-$field"
  MOCK_SCENARIO=cloudkit
  continue_native
  calls=$(wc -l < "$MOCK_STATE/trust-calls")
  case $field in
    source_commit) foreign=2222222222222222222222222222222222222222 ;;
    version) foreign=0.0.2 ;;
    architecture) foreign=amd64 ;;
    ci_run_id) foreign=124 ;;
    ci_run_attempt) foreign=2 ;;
    team_id) foreign=ZZZZZZZZZZ ;;
    *) foreign=2222222222222222222222222222222222222222222222222222222222222222 ;;
  esac
  jq --arg key "$field" --arg foreign "$foreign" '.[$key]=$foreign' "$NOTARY_EVIDENCE_DIR/acceptance-state.json" > "$MOCK_STATE/change.json"
  mv "$MOCK_STATE/change.json" "$NOTARY_EVIDENCE_DIR/acceptance-state.json"
  # The refusal must come from exact binding, not schema validation.
  bash -c 'source "$1"; macos_acceptance_state_valid "$2"' _ "$TEST_WORK/repo/scripts/macos-ticket-policy.sh" "$NOTARY_EVIDENCE_DIR/acceptance-state.json"
  MOCK_SCENARIO=success
  if continue_native; then exit 1; fi
  [[ $(wc -l < "$MOCK_STATE/trust-calls") == "$calls" && ! -f $NOTARY_EVIDENCE_DIR/acceptance.json ]]
done
echo 'PASS: structurally valid but foreign candidate, archive and receipt bindings refuse native recovery'
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
for scenario in wrong-copied-version wrong-copied-team wrong-copied-architecture wrong-copied-executable-version checksum-corruption manifest-extra executable-mismatch gatekeeper-failure stapler-failure codesign-failure attach-failure detach-failure missing-timestamp missing-runtime; do
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
    missing-timestamp) expected='secure signing timestamp missing' ;;
    missing-runtime) expected='hardened runtime missing' ;;
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
  if [[ $scenario == checksum-corruption || $scenario == manifest-extra || $scenario == missing-timestamp || $scenario == missing-runtime ]]; then
    [[ ! -f $MOCK_STATE/mount-calls ]]
  else
    [[ -f $MOCK_STATE/mount-calls ]]
  fi
  [[ ! -e $NOTARY_EVIDENCE_DIR/acceptance.json ]]
  echo "PASS: $scenario rejected, owned mount cleanup attempted, quarantine preserved"
done
echo 'PASS: synthetic final-package acceptance contracts; no Apple, real mounts or signing'
