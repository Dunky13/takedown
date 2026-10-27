#!/usr/bin/env bash
# Final distribution acceptance, without signing credentials or global installs.
set -euo pipefail
umask 077
fail() { printf 'verify-macos-package: %s\n' "$*" >&2; exit 1; }
[[ $(uname -s) == Darwin ]] || fail 'macOS is required'
: "${VERSION:?VERSION is required}"
: "${RELEASE_VERSION:?RELEASE_VERSION is required}"
: "${NOTARY_ARCH:?NOTARY_ARCH is required}"
: "${APPLE_TEAM_ID:?APPLE_TEAM_ID is required}"
[[ $VERSION =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]] || fail 'invalid numeric version'
[[ $RELEASE_VERSION =~ ^[0-9A-Za-z.+-]+$ && ${RELEASE_VERSION%%[-+]*} == "$VERSION" ]] || fail 'invalid release version'
[[ $APPLE_TEAM_ID =~ ^[A-Z0-9]{10}$ ]] || fail 'invalid Apple team'
case $NOTARY_ARCH in arm64) slice=arm64 ;; amd64) slice=x86_64 ;; *) fail 'invalid architecture' ;; esac
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
PACKAGES="$ROOT/dist/packages/darwin-$NOTARY_ARCH-release"
BASE="takedown-$RELEASE_VERSION-macos-$NOTARY_ARCH-release"
[[ -d $PACKAGES && ! -L $PACKAGES ]] || fail 'signed package directory missing'
if [[ -n ${NOTARY_CHECKPOINT_HOOK:-} ]]; then
  [[ $NOTARY_CHECKPOINT_HOOK == /* && -f $NOTARY_CHECKPOINT_HOOK && ! -L $NOTARY_CHECKPOINT_HOOK ]] || fail 'invalid acceptance checkpoint hook'
fi
WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-acceptance.XXXXXX")
MOUNTED=false
# shellcheck source=scripts/macos-ticket-policy.sh
source "$ROOT/scripts/macos-ticket-policy.sh"
STATE_ACTIVE=false CHECKPOINT_FAILED=false RESERVED=false
ATTEMPTS=0 FAILURE_CLASS=''
checkpoint() {
  local status=0
  if [[ -n ${NOTARY_CHECKPOINT_HOOK:-} ]]; then
    bash "$NOTARY_CHECKPOINT_HOOK" "$NOTARY_EVIDENCE_DIR" >/dev/null 2>&1 || status=$?
    [[ $status == 0 ]] || { CHECKPOINT_FAILED=true; rm -f "$NOTARY_EVIDENCE_DIR/acceptance.json"; fail 'acceptance checkpoint failed; private output withheld'; }
  fi
}
acceptance_state() {
  jq --arg state "$1" --argjson attempts "$ATTEMPTS" --arg failure "$FAILURE_CLASS" \
    --argjson binding "$BINDING_JSON" -n '$binding + {schema:1,state:$state,attempts:$attempts,failure_class:$failure}' > "$NOTARY_EVIDENCE_DIR/.acceptance-state.tmp"
  mv "$NOTARY_EVIDENCE_DIR/.acceptance-state.tmp" "$NOTARY_EVIDENCE_DIR/acceptance-state.json"
  checkpoint
}
cleanup() {
  result=$?
  trap - EXIT
  if [[ $MOUNTED == true ]]; then
    if hdiutil detach -quiet "$WORK/mount"; then MOUNTED=false; else result=1; fi
  fi
  # Never recurse into an owned mount that could not be detached.
  if [[ -n $WORK && $MOUNTED == false ]]; then rm -rf -- "$WORK" || result=1; fi
  # Handled interruption retains the reserved pending attempt. Cleanup failure
  # overrides pending and is terminal; neither can retain an acceptance receipt.
  if [[ $STATE_ACTIVE == true && $result != 0 && $result != 75 && $result != 130 && $result != 143 && $CHECKPOINT_FAILED == false ]]; then
    [[ -n $FAILURE_CLASS ]] || FAILURE_CLASS=unclassified
    ( acceptance_state blocked ) || result=1
  fi
  exit "$result"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
for extension in zip dmg; do
  [[ -f $PACKAGES/$BASE.$extension && ! -L $PACKAGES/$BASE.$extension ]] || fail 'signed asset missing'
done
[[ -f $PACKAGES/SHA256SUMS && ! -L $PACKAGES/SHA256SUMS ]] || fail 'checksum manifest missing'
printf '%s\n' "$BASE.zip" "$BASE.dmg" | LC_ALL=C sort > "$WORK/expected"
awk '/^[0-9a-fA-F]+  [^\/]+$/ && length($1) == 64 { print $2 }' "$PACKAGES/SHA256SUMS" | LC_ALL=C sort > "$WORK/actual"
cmp -s "$WORK/expected" "$WORK/actual" || fail 'manifest must cover exactly the final ZIP and DMG'
(cd "$PACKAGES" && shasum -a 256 -c SHA256SUMS)
hash() { shasum -a 256 "$1" | awk '{print $1}'; }
ZIP_SHA=$(hash "$PACKAGES/$BASE.zip")
DMG_SHA=$(hash "$PACKAGES/$BASE.dmg")
if [[ -n ${NOTARY_EVIDENCE_DIR:-} ]]; then
  [[ -d $NOTARY_EVIDENCE_DIR && ! -L $NOTARY_EVIDENCE_DIR ]] || fail 'invalid acceptance evidence directory'
  for name in input.json app/receipt.json dmg/receipt.json; do
    [[ -f $NOTARY_EVIDENCE_DIR/$name && ! -L $NOTARY_EVIDENCE_DIR/$name ]] || fail 'missing acceptance binding'
  done
  bash "$ROOT/scripts/apple-stages-verified.sh" "$NOTARY_EVIDENCE_DIR" || fail 'signer stages are not verified'
  jq -n --arg source "${NOTARY_SOURCE_SHA:?}" --arg version "$RELEASE_VERSION" --arg arch "$NOTARY_ARCH" \
    --arg run "${GITHUB_RUN_ID:?}" --arg attempt "${GITHUB_RUN_ATTEMPT:?}" --arg team "$APPLE_TEAM_ID" \
    --arg zip "$ZIP_SHA" --arg dmg "$DMG_SHA" \
    --arg input "$(hash "$NOTARY_EVIDENCE_DIR/input.json")" \
    --arg app_receipt "$(hash "$NOTARY_EVIDENCE_DIR/app/receipt.json")" \
    --arg dmg_receipt "$(hash "$NOTARY_EVIDENCE_DIR/dmg/receipt.json")" \
    '{source_commit:$source,version:$version,architecture:$arch,ci_run_id:$run,ci_run_attempt:$attempt,
      team_id:$team,zip_sha256:$zip,dmg_sha256:$dmg,input_sha256:$input,
      app_receipt_sha256:$app_receipt,dmg_receipt_sha256:$dmg_receipt}' > "$WORK/binding.json"
  BINDING_JSON=$(cat "$WORK/binding.json")
  if [[ -e $NOTARY_EVIDENCE_DIR/acceptance-state.json || -L $NOTARY_EVIDENCE_DIR/acceptance-state.json ]]; then
    [[ -f $NOTARY_EVIDENCE_DIR/acceptance-state.json && ! -L $NOTARY_EVIDENCE_DIR/acceptance-state.json ]] || fail 'invalid acceptance state file'
    macos_acceptance_state_valid "$NOTARY_EVIDENCE_DIR/acceptance-state.json" || fail 'malformed acceptance state'
    jq -e --slurpfile binding "$WORK/binding.json" '. as $s | $binding[0] |
      to_entries | all(.[]; $s[.key] == .value)' "$NOTARY_EVIDENCE_DIR/acceptance-state.json" >/dev/null || fail 'acceptance state binding mismatch'
    ATTEMPTS=$(jq -r .attempts "$NOTARY_EVIDENCE_DIR/acceptance-state.json")
    prior=$(jq -r .state "$NOTARY_EVIDENCE_DIR/acceptance-state.json")
    [[ $prior != blocked ]] || fail 'acceptance is terminal; no automatic retry'
    if [[ $prior == pending ]]; then
      [[ $ATTEMPTS -lt 3 ]] || { FAILURE_CLASS=unclassified; acceptance_state blocked; fail 'acceptance recovery budget exhausted'; }
      ATTEMPTS=$((ATTEMPTS + 1))
      RESERVED=true
    fi
  fi
  # Invalidate stale success before native operations. The state checkpoint also
  # carries exact finals even on the first acceptance failure.
  rm -f "$NOTARY_EVIDENCE_DIR/acceptance.json" "$NOTARY_EVIDENCE_DIR/fresh-acceptance-$NOTARY_ARCH.json"
  STATE_ACTIVE=true
  acceptance_state pending
fi
validate_ticket() {
  local result=0
  xcrun stapler validate "$1" > "$WORK/ticket.log" 2>&1 || result=$?
  [[ $result != 0 ]] || return 0
  FAILURE_CLASS=$(macos_ticket_failure_class "$WORK/ticket.log" "$result")
  printf 'verify-macos-package: stapler-validate exit=%s cause=%s (private output withheld)\n' "$result" "$FAILURE_CLASS" >&2
  if [[ $STATE_ACTIVE == true ]]; then
    if [[ $RESERVED == false && $ATTEMPTS -lt 3 ]]; then ATTEMPTS=$((ATTEMPTS + 1)); fi
    if [[ $FAILURE_CLASS == cloudkit-transport && $ATTEMPTS -lt 3 ]]; then
      acceptance_state pending
      exit 75
    fi
  fi
  fail 'ticket validation failed permanently or exhausted its bounded recovery budget'
}
check_app() {
  local app=$1
  [[ -d $app && ! -L $app && -f $app/Contents/MacOS/takedown && ! -L $app/Contents/MacOS/takedown ]] || fail 'invalid copied application'
  [[ $(plutil -extract CFBundleIdentifier raw -o - "$app/Contents/Info.plist") == io.github.dunky13.takedown ]] || fail 'wrong bundle identifier'
  [[ $(plutil -extract CFBundleExecutable raw -o - "$app/Contents/Info.plist") == takedown ]] || fail 'wrong bundle executable'
  [[ $(plutil -extract CFBundleShortVersionString raw -o - "$app/Contents/Info.plist") == "$VERSION" ]] || fail 'wrong bundle version'
  [[ $(plutil -extract CFBundleVersion raw -o - "$app/Contents/Info.plist") == "$VERSION" ]] || fail 'wrong bundle build'
  [[ $(lipo -archs "$app/Contents/MacOS/takedown") == "$slice" ]] || fail 'wrong binary slices'
  [[ $("$app/Contents/MacOS/takedown" --version) == "takedown version $RELEASE_VERSION" ]] || fail 'wrong executable version'
  codesign --verify --deep --strict --verbose=2 "$app"
  codesign --display --verbose=4 "$app" 2> "$WORK/signature"
  grep -qx "TeamIdentifier=$APPLE_TEAM_ID" "$WORK/signature" || fail 'wrong signed team'
  grep -q '^Authority=Developer ID Application: ' "$WORK/signature" || fail 'Developer ID signature missing'
  validate_ticket "$app"
  spctl --assess --type execute --verbose=2 "$app"
}
mkdir "$WORK/zip" "$WORK/mount" "$WORK/Applications"
ditto -x -k "$PACKAGES/$BASE.zip" "$WORK/zip"
check_app "$WORK/zip/takedown.app"
codesign --verify --strict --verbose=2 "$PACKAGES/$BASE.dmg"
codesign --display --verbose=4 "$PACKAGES/$BASE.dmg" 2> "$WORK/signature"
grep -qx "TeamIdentifier=$APPLE_TEAM_ID" "$WORK/signature" || fail 'wrong DMG team'
grep -q '^Authority=Developer ID Application: ' "$WORK/signature" || fail 'DMG Developer ID signature missing'
validate_ticket "$PACKAGES/$BASE.dmg"
spctl --assess --type open --context context:primary-signature --verbose=2 "$PACKAGES/$BASE.dmg"
# Set before attach so cleanup also attempts to detach a partially mounted image.
MOUNTED=true
hdiutil attach -quiet -readonly -nobrowse -mountpoint "$WORK/mount" "$PACKAGES/$BASE.dmg"
ditto "$WORK/mount/takedown.app" "$WORK/Applications/takedown.app"
check_app "$WORK/Applications/takedown.app"
cmp -s "$WORK/zip/takedown.app/Contents/MacOS/takedown" "$WORK/Applications/takedown.app/Contents/MacOS/takedown" || fail 'ZIP and installed DMG executables differ'
[[ $(hash "$PACKAGES/$BASE.zip") == "$ZIP_SHA" && $(hash "$PACKAGES/$BASE.dmg") == "$DMG_SHA" ]] || fail 'final archive digest changed during acceptance'
# Cleanup is part of acceptance, including detach after a copied-app failure.
hdiutil detach -quiet "$WORK/mount"
MOUNTED=false
rm -rf -- "$WORK"
WORK=''
if [[ -n ${NOTARY_EVIDENCE_DIR:-} ]]; then
  jq -n --arg version "$RELEASE_VERSION" --arg arch "$NOTARY_ARCH" \
    --arg source "${NOTARY_SOURCE_SHA:?NOTARY_SOURCE_SHA is required}" --arg run "${GITHUB_RUN_ID:?}" --arg attempt "${GITHUB_RUN_ATTEMPT:?}" \
    --arg zip "$ZIP_SHA" --arg dmg "$DMG_SHA" \
    '{source_commit:$source,ci_run_id:$run,ci_run_attempt:$attempt,version:$version,architecture:$arch,zip_sha256:$zip,dmg_sha256:$dmg,checks:["checksums","Developer ID","stapling","Gatekeeper","copied app metadata and version","native slices","ZIP/DMG executable equality"]}' \
    > "$NOTARY_EVIDENCE_DIR/.acceptance.tmp"
  mv "$NOTARY_EVIDENCE_DIR/.acceptance.tmp" "$NOTARY_EVIDENCE_DIR/acceptance.json"
  FAILURE_CLASS=''
  acceptance_state verified
  STATE_ACTIVE=false
fi
echo 'PASS: final ZIP, mounted DMG and copied application signatures, tickets, Gatekeeper, metadata and native version'
