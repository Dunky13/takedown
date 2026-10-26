#!/usr/bin/env bash
# Explicit release operation. Only bound receipts/allowlisted diagnostics leave scratch.
{ set +x; } 2>/dev/null
set -euo pipefail
umask 077
fail() { printf 'sign-macos: %s\n' "$*" >&2; exit 1; }
[[ $# == 1 ]] || fail 'usage: sign-macos.sh --check | path.app | path.dmg'
[[ $(uname -s) == Darwin ]] || fail 'macOS is required'
for name in MACOS_CERTIFICATE_P12_BASE64 MACOS_CERTIFICATE_PASSWORD MACOS_SIGN_IDENTITY \
  APPLE_TEAM_ID APPLE_API_KEY_BASE64 APPLE_API_KEY_ID APPLE_API_ISSUER; do
  [[ -n ${!name:-} ]] || fail "missing required secret/configuration: $name"
done
[[ $APPLE_TEAM_ID =~ ^[A-Z0-9]{10}$ ]] || fail 'invalid APPLE_TEAM_ID'
[[ $MACOS_SIGN_IDENTITY == "Developer ID Application: "*" ($APPLE_TEAM_ID)" ]] || fail 'identity must be a Developer ID Application identity for APPLE_TEAM_ID'
for tool in security codesign xcrun openssl ditto plutil jq shasum spctl lipo; do
  command -v "$tool" >/dev/null || fail "missing prerequisite: $tool"
done
xcrun --find notarytool >/dev/null 2>&1 || fail 'prerequisite notarytool lookup failed'
xcrun --find stapler >/dev/null 2>&1 || fail 'prerequisite stapler lookup failed'
[[ $1 != --check ]] || { echo 'macOS signing prerequisites present (credentials not authenticated).'; exit 0; }
[[ ${NOTARY_SOURCE_SHA:-} =~ ^[0-9a-f]{40}$ ]] || fail 'NOTARY_SOURCE_SHA must be an exact source commit'
[[ ${NOTARY_ARCH:-} == arm64 || ${NOTARY_ARCH:-} == amd64 ]] || fail 'NOTARY_ARCH must be arm64 or amd64'
[[ ${RELEASE_VERSION:-} =~ ^[0-9]+\.[0-9]+\.[0-9]+([+-][0-9A-Za-z.+-]+)?$ ]] || fail 'RELEASE_VERSION must be a safe release version'
[[ ${GITHUB_RUN_ID:-0} =~ ^[0-9]+$ && ${GITHUB_RUN_ATTEMPT:-0} =~ ^[0-9]+$ ]] || fail 'invalid CI run identity'
[[ -n ${NOTARY_EVIDENCE_DIR:-} && $NOTARY_EVIDENCE_DIR == /* && ! -L $NOTARY_EVIDENCE_DIR ]] || fail 'NOTARY_EVIDENCE_DIR must be an absolute persistent directory'
if [[ -n ${NOTARY_CHECKPOINT_HOOK:-} ]]; then
  [[ $NOTARY_CHECKPOINT_HOOK == /* && -f $NOTARY_CHECKPOINT_HOOK && ! -L $NOTARY_CHECKPOINT_HOOK ]] || fail 'NOTARY_CHECKPOINT_HOOK must be an absolute trusted regular script'
fi
TARGET=$1
[[ $TARGET == /* ]] || TARGET="$PWD/$TARGET"
[[ ! -L $TARGET ]] || fail 'target must not be a symlink'
case $TARGET in
  *.app) STAGE=app; EXT=zip; VERIFIED_NAME=stapled-app.zip ;;
  *.dmg) STAGE=dmg; EXT=dmg; VERIFIED_NAME=stapled.dmg ;;
  *) fail 'expected an app bundle or DMG' ;;
esac
mkdir -p "$NOTARY_EVIDENCE_DIR"
NOTARY_EVIDENCE_DIR=$(cd -P -- "$NOTARY_EVIDENCE_DIR" && pwd -P)
EVIDENCE="$NOTARY_EVIDENCE_DIR/$STAGE"
[[ ! -L $EVIDENCE ]] || fail 'symlinked stage evidence'
mkdir -p "$EVIDENCE"
# A restored durable snapshot contains no lock. Never steal a live/local stale lock.
mkdir "$EVIDENCE/.lock" 2>/dev/null || fail "concurrent or interrupted local $STAGE operation; reconcile its lock before continuing"
WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-sign.XXXXXX")
KEYCHAIN="$WORK/signing.keychain-db"
# shellcheck source=scripts/keychain-lifecycle.sh
source "$(dirname -- "$0")/keychain-lifecycle.sh"
KEYCHAIN_MUTATION_STARTED=false
BASENAME="takedown-$RELEASE_VERSION-macos-$NOTARY_ARCH-${NOTARY_SOURCE_SHA:0:12}-${GITHUB_RUN_ID:-0}-${GITHUB_RUN_ATTEMPT:-0}-$STAGE.$EXT"
SUBMISSION_ID='' SUBMITTED_SHA='' VERIFIED_SHA='' STATE='' FAILURE_PHASE='' FAILURE_EXIT='' TICKET_ATTEMPTS=0
CHECKPOINT_FAILED=false
hash() { shasum -a 256 "$1" | cut -d ' ' -f 1; }
checkpoint() {
  local result=0
  if [[ -n ${NOTARY_CHECKPOINT_HOOK:-} ]]; then
    bash "$NOTARY_CHECKPOINT_HOOK" "$NOTARY_EVIDENCE_DIR" >/dev/null 2>&1 || result=$?
    if [[ $result != 0 ]]; then
      CHECKPOINT_FAILED=true
      printf 'sign-macos: stage=%s command=checkpoint exit=%s; stopping before further release mutation\n' "$STAGE" "$result" >&2
      exit 1
    fi
  fi
}
receipt() {
  STATE=$1
  jq -n --arg state "$STATE" --arg id "$SUBMISSION_ID" --arg filename "$BASENAME" \
    --arg sha256 "$SUBMITTED_SHA" --arg source "$NOTARY_SOURCE_SHA" \
    --arg version "$RELEASE_VERSION" --arg architecture "$NOTARY_ARCH" --arg stage "$STAGE" \
    --arg run "${GITHUB_RUN_ID:-}" --arg attempt "${GITHUB_RUN_ATTEMPT:-}" \
    --arg team "$APPLE_TEAM_ID" --arg identity "$MACOS_SIGN_IDENTITY" \
    --arg verified "$VERIFIED_SHA" --arg verified_name "$VERIFIED_NAME" \
    --arg phase "$FAILURE_PHASE" --arg result "$FAILURE_EXIT" \
    --argjson ticket_attempts "$TICKET_ATTEMPTS" \
    '{schema:2,state:$state,submission_id:(if $id=="" then null else $id end),filename:$filename,
      submitted_sha256:$sha256,source_commit:$source,version:$version,architecture:$architecture,
      stage:$stage,ci_run_id:$run,ci_run_attempt:$attempt,team_id:$team,signing_identity:$identity,
      verified_filename:$verified_name,verified_sha256:$verified,failure_phase:$phase,failure_exit:$result,ticket_attempts:$ticket_attempts}' \
    > "$EVIDENCE/.receipt.tmp"
  mv "$EVIDENCE/.receipt.tmp" "$EVIDENCE/receipt.json"
  checkpoint
}
command_failed() {
  FAILURE_PHASE=$1 FAILURE_EXIT=$2
  printf 'sign-macos: stage=%s command=%s exit=%s (private output withheld)\n' "$STAGE" "$FAILURE_PHASE" "$FAILURE_EXIT" >&2
  [[ -z $STATE ]] || receipt "$STATE"
  exit 1
}
run() {
  local label=$1 result=0
  shift
  "$@" > "$WORK/command.log" 2>&1 || result=$?
  [[ $result == 0 ]] || command_failed "$label" "$result"
}
cleanup_keychain() {
  if [[ $KEYCHAIN_MUTATION_STARTED == true ]]; then
    if ! keychain_restore_search_list; then return 1; fi
    KEYCHAIN_MUTATION_STARTED=false
  fi
}
cleanup() {
  local result=$?
  trap - EXIT
  if ! cleanup_keychain; then
    printf 'sign-macos: stage=%s command=keychain-cleanup exit=1\n' "$STAGE" >&2
    result=1
    if [[ -n $STATE && $CHECKPOINT_FAILED == false ]]; then
      FAILURE_PHASE=keychain-cleanup FAILURE_EXIT=1
      ( receipt cleanup-failed ) || result=1
    fi
  fi
  rm -rf -- "$WORK" || result=1
  rmdir "$EVIDENCE/.lock" || result=1
  exit "$result"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM

[[ -z $(find "$EVIDENCE" -type l -print) ]] || fail 'symlink in stage evidence'
for retained in "$EVIDENCE"/* "$EVIDENCE"/.[!.]* "$EVIDENCE"/..?*; do
  [[ -e $retained ]] || continue
  case ${retained##*/} in
    receipt.json|apple-log.json|"$BASENAME"|"$VERIFIED_NAME"|.lock) ;;
    *) fail 'unexpected unbound stage evidence' ;;
  esac
done
if [[ -e $EVIDENCE/receipt.json || -L $EVIDENCE/receipt.json ]]; then
  [[ -f $EVIDENCE/receipt.json && ! -L $EVIDENCE/receipt.json ]] || fail 'invalid receipt file'
  jq -e --arg filename "$BASENAME" --arg source "$NOTARY_SOURCE_SHA" --arg version "$RELEASE_VERSION" \
    --arg arch "$NOTARY_ARCH" --arg stage "$STAGE" --arg run "${GITHUB_RUN_ID:-}" --arg attempt "${GITHUB_RUN_ATTEMPT:-}" \
    --arg team "$APPLE_TEAM_ID" --arg identity "$MACOS_SIGN_IDENTITY" --arg verified_name "$VERIFIED_NAME" \
    '.schema == 2 and .filename == $filename and .source_commit == $source and .version == $version and
     .architecture == $arch and .stage == $stage and .ci_run_id == $run and .ci_run_attempt == $attempt and
     .team_id == $team and .signing_identity == $identity and .verified_filename == $verified_name and
     (.state | IN("prepared","submission-unknown","submitted","pending","accepted","verified")) and
     (.ticket_attempts | type == "number" and floor == . and . >= 0 and . <= 3) and
     (.submitted_sha256 | type == "string") and (.verified_sha256 | type == "string")' \
    "$EVIDENCE/receipt.json" >/dev/null 2>&1 || fail 'unbound, mismatched or blocked stage receipt; no reuse permitted'
  STATE=$(jq -r .state "$EVIDENCE/receipt.json")
  SUBMISSION_ID=$(jq -r '.submission_id // ""' "$EVIDENCE/receipt.json")
  SUBMITTED_SHA=$(jq -r .submitted_sha256 "$EVIDENCE/receipt.json")
  VERIFIED_SHA=$(jq -r .verified_sha256 "$EVIDENCE/receipt.json")
  TICKET_ATTEMPTS=$(jq -r .ticket_attempts "$EVIDENCE/receipt.json")
  if [[ $STATE == prepared ]]; then
    [[ ! -e $EVIDENCE/apple-log.json && -z $SUBMISSION_ID && -z $SUBMITTED_SHA && -z $VERIFIED_SHA && ! -e $EVIDENCE/$BASENAME && ! -L $EVIDENCE/$BASENAME && ! -e $EVIDENCE/$VERIFIED_NAME && ! -L $EVIDENCE/$VERIFIED_NAME ]] || fail 'prepared stage contains upload intent or retained bytes; reconcile before retrying'
  else
    [[ $SUBMITTED_SHA =~ ^[0-9a-f]{64}$ && -f $EVIDENCE/$BASENAME && ! -L $EVIDENCE/$BASENAME ]] || fail 'missing bound submitted bytes'
    [[ $(hash "$EVIDENCE/$BASENAME") == "$SUBMITTED_SHA" ]] || fail 'submitted archive hash mismatch'
    if [[ $STATE == submission-unknown ]]; then
      [[ -z $SUBMISSION_ID ]] || fail 'ambiguous unknown-submission receipt'
    else
      [[ $SUBMISSION_ID =~ ^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$ ]] || fail 'invalid retained UUID'
    fi
    if [[ $STATE == verified ]]; then
      [[ $VERIFIED_SHA =~ ^[0-9a-f]{64}$ && -f $EVIDENCE/$VERIFIED_NAME && ! -L $EVIDENCE/$VERIFIED_NAME ]] || fail 'missing verified bytes'
      [[ $(hash "$EVIDENCE/$VERIFIED_NAME") == "$VERIFIED_SHA" ]] || fail 'verified archive hash mismatch'
    else
      [[ -z $VERIFIED_SHA ]] || fail 'unexpected verified hash before verification'
    fi
  fi
else
  # No legacy/unbound files may silently become a new submission.
  [[ -z $(find "$EVIDENCE" -mindepth 1 -maxdepth 1 ! -name .lock -print) ]] || fail 'unbound stage files without receipt'
  receipt prepared
fi
# Re-checkpoint a restored state before any local signing or remote operation.
checkpoint

verify_signature() {
  local item=$1
  run signature-verify codesign --verify --deep --strict --verbose=2 "$item"
  local result=0
  codesign --display --verbose=4 "$item" > "$WORK/signature.txt" 2>&1 || result=$?
  [[ $result == 0 ]] || command_failed signature-display "$result"
  grep -qx "TeamIdentifier=$APPLE_TEAM_ID" "$WORK/signature.txt" || fail 'signed artifact has the wrong TeamIdentifier'
  grep -Fqx "Authority=$MACOS_SIGN_IDENTITY" "$WORK/signature.txt" || fail 'signed artifact has the wrong Developer ID authority'
  grep -Eq '^Timestamp=.+$' "$WORK/signature.txt" || fail 'signed artifact lacks a secure timestamp'
  if [[ $STAGE == app ]]; then
    grep -Eq '^CodeDirectory .*flags=.*\(.*runtime.*\)' "$WORK/signature.txt" || fail 'signed app lacks hardened runtime'
  fi
}
# Recovery needs no certificate import and never signs retained artifacts again.
result=0
printf '%s' "$APPLE_API_KEY_BASE64" | openssl base64 -d -A > "$WORK/key.p8" 2> "$WORK/command.log" || result=$?
[[ $result == 0 ]] || command_failed api-key-base64-decode "$result"
[[ -s $WORK/key.p8 ]] || command_failed api-key-base64-decode 1
AUTH=(--key "$WORK/key.p8" --key-id "$APPLE_API_KEY_ID" --issuer "$APPLE_API_ISSUER")

CANDIDATE="$WORK/takedown.$STAGE"
if [[ $STATE == prepared ]]; then
  if [[ $STAGE == app ]]; then
    [[ -d $TARGET && -f $TARGET/Contents/MacOS/takedown && ! -L $TARGET/Contents/MacOS/takedown ]] || fail 'invalid app bundle'
    run copy-app ditto "$TARGET" "$CANDIDATE"
    BINARY_ARCH=$(lipo -archs "$CANDIDATE/Contents/MacOS/takedown" 2>/dev/null) || fail 'cannot inspect app architecture'
    EXPECTED_ARCH=$NOTARY_ARCH; [[ $EXPECTED_ARCH != amd64 ]] || EXPECTED_ARCH=x86_64
    [[ $BINARY_ARCH == "$EXPECTED_ARCH" ]] || fail 'app architecture does not match candidate'
    BUNDLE_VERSION=$(plutil -extract CFBundleShortVersionString raw -o - "$CANDIDATE/Contents/Info.plist" 2>/dev/null) || fail 'cannot inspect app version'
    [[ $BUNDLE_VERSION == "${RELEASE_VERSION%%[-+]*}" ]] || fail 'app version does not match candidate'
  else
    [[ -f $TARGET ]] || fail 'disk image does not exist'
    run copy-dmg cp "$TARGET" "$CANDIDATE"
  fi
  # Decode privately; each operation has a safe label and captures its own status.
  result=0
  printf '%s' "$MACOS_CERTIFICATE_P12_BASE64" | openssl base64 -d -A > "$WORK/certificate.p12" 2> "$WORK/command.log" || result=$?
  [[ $result == 0 ]] || command_failed certificate-base64-decode "$result"
  [[ -s $WORK/certificate.p12 ]] || command_failed certificate-base64-decode 1
  run certificate-decode openssl pkcs12 -in "$WORK/certificate.p12" -passin env:MACOS_CERTIFICATE_PASSWORD -nodes -out "$WORK/import.pem"
  run certificate-rewrap openssl pkcs12 -export -in "$WORK/import.pem" -out "$WORK/import.p12" -passout pass:
  result=0
  keychain_capture_search_list > "$WORK/command.log" 2>&1 || result=$?
  [[ $result == 0 ]] || command_failed "$KEYCHAIN_CAPTURE_COMMAND" "$result"
  KEYCHAIN_MUTATION_STARTED=true
  # Each secret-bearing command gets its own private stdin invocation, with
  # its status checked separately. Never pass the strong password through argv.
  result=0
  KEYCHAIN_PASSWORD=$(openssl rand -hex 32 2> "$WORK/command.log") || result=$?
  [[ $result == 0 ]] || command_failed keychain-password "$result"
  [[ $KEYCHAIN_PASSWORD =~ ^[a-f0-9]{64}$ ]] || command_failed keychain-password 1
  case $KEYCHAIN in *'"'*|*\\*|*$'\n'*|*$'\r'*) fail 'unsafe temporary keychain path' ;; esac
  keychain_secret() {
    local operation=$1
    case $operation in
      create|unlock) printf '%s-keychain -p "%s" "%s"\n' "$operation" "$KEYCHAIN_PASSWORD" "$KEYCHAIN" | security -i ;;
      partitions) printf 'set-key-partition-list -S apple-tool:,apple:,codesign: -s -k "%s" "%s"\n' "$KEYCHAIN_PASSWORD" "$KEYCHAIN" | security -i ;;
    esac
  }
  run keychain-create keychain_secret create
  run keychain-permissions chmod 600 "$KEYCHAIN"
  run keychain-settings security set-keychain-settings -lut 21600 "$KEYCHAIN"
  run keychain-unlock keychain_secret unlock
  run keychain-import security import "$WORK/import.p12" -k "$KEYCHAIN" -f pkcs12 -x -P '' -T /usr/bin/codesign
  run keychain-partitions keychain_secret partitions
  run keychain-search-list security list-keychains -d user -s "$KEYCHAIN" ${ORIGINAL_KEYCHAINS[@]+"${ORIGINAL_KEYCHAINS[@]}"}
  result=0
  security find-identity -v -p codesigning "$KEYCHAIN" > "$WORK/identities.txt" 2> "$WORK/command.log" || result=$?
  [[ $result == 0 ]] || command_failed keychain-find-identity "$result"
  grep -Fq "\"$MACOS_SIGN_IDENTITY\"" "$WORK/identities.txt" || fail 'configured Developer ID identity is not available in the temporary keychain'
  run certificate-remove rm -f "$WORK/certificate.p12" "$WORK/import.pem" "$WORK/import.p12"
  if [[ $STAGE == app ]]; then
    run executable-sign codesign --force --timestamp --options runtime --keychain "$KEYCHAIN" --sign "$MACOS_SIGN_IDENTITY" "$CANDIDATE/Contents/MacOS/takedown"
    run app-sign codesign --force --timestamp --options runtime --keychain "$KEYCHAIN" --sign "$MACOS_SIGN_IDENTITY" "$CANDIDATE"
  else
    run dmg-sign codesign --force --timestamp --keychain "$KEYCHAIN" --sign "$MACOS_SIGN_IDENTITY" "$CANDIDATE"
  fi
  verify_signature "$CANDIDATE"
  # Credential cleanup must verify before uploading or promoting anything.
  cleanup_keychain || command_failed keychain-cleanup 1
  if [[ $STAGE == app ]]; then
    run submitted-archive ditto -c -k --sequesterRsrc --keepParent "$CANDIDATE" "$WORK/$BASENAME"
  else
    run submitted-archive cp "$CANDIDATE" "$WORK/$BASENAME"
  fi
  run retain-submitted cp "$WORK/$BASENAME" "$EVIDENCE/$BASENAME"
  SUBMITTED_SHA=$(hash "$EVIDENCE/$BASENAME")
  # Bytes and intent are one complete checkpoint, acknowledged BEFORE submit.
  receipt submission-unknown
fi

# Only this invocation's newly prepared intent may submit. Restored unknown intent
# always uses history, including when a prior checkpoint failed before the upload.
if [[ -f $WORK/$BASENAME ]]; then
  result=0
  xcrun notarytool submit "$EVIDENCE/$BASENAME" "${AUTH[@]}" --no-wait --output-format json > "$WORK/submit.json" 2> "$WORK/notary.stderr" || result=$?
  SUBMISSION_ID=$(jq -er '.id' "$WORK/submit.json" 2>/dev/null || true)
  [[ $SUBMISSION_ID =~ ^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$ ]] || {
    SUBMISSION_ID=''
    [[ $result != 0 ]] || result=1
    command_failed notary-submit-uuid "$result"
  }
  receipt submitted
  if [[ ${NOTARY_YIELD_AFTER_SUBMIT:-false} == true ]]; then
    printf 'sign-macos: stage=%s new UUID durably checkpointed; qualification continuation required (exit 75)\n' "$STAGE" >&2
    exit 75
  fi
  # Even a failed client with a valid UUID continues by that UUID, never uploads.
fi
if [[ $STATE == submission-unknown ]]; then
  run notary-history xcrun notarytool history "${AUTH[@]}" --output-format json
  cp "$WORK/command.log" "$WORK/history.json"
  jq -e '.history | type == "array"' "$WORK/history.json" >/dev/null 2>&1 || fail 'malformed notarization history; upload remains unresolved'
  MATCHES=$(jq --arg filename "$BASENAME" '[.history[] | select(.name == $filename)] | length' "$WORK/history.json")
  [[ $MATCHES == 1 ]] || fail 'upload remains unresolved: history filename absent or ambiguous; automatic resubmission prohibited'
  SUBMISSION_ID=$(jq -r --arg filename "$BASENAME" '.history[] | select(.name == $filename) | .id' "$WORK/history.json")
  [[ $SUBMISSION_ID =~ ^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$ ]] || fail 'invalid history UUID'
  receipt submitted
fi
printf 'Notarization submission (%s): %s\n' "$STAGE" "$SUBMISSION_ID"
run notary-info xcrun notarytool info "$SUBMISSION_ID" "${AUTH[@]}" --output-format json
cp "$WORK/command.log" "$WORK/info.json"
jq -e --arg id "$SUBMISSION_ID" --arg filename "$BASENAME" '.id == $id and .name == $filename and (.status | type == "string")' "$WORK/info.json" >/dev/null 2>&1 || fail 'Apple info UUID/filename/status mismatch'
STATUS=$(jq -r .status "$WORK/info.json")
WAIT_RESULT=0
if [[ $STATUS == 'In Progress' ]]; then
  WAIT_RESULT=0
  xcrun notarytool wait "$SUBMISSION_ID" "${AUTH[@]}" --timeout 30m --output-format json > "$WORK/wait.json" 2> "$WORK/notary.stderr" || WAIT_RESULT=$?
  jq -e --arg id "$SUBMISSION_ID" '.id == $id and (.status | type == "string")' "$WORK/wait.json" >/dev/null 2>&1 || command_failed notary-wait-response "$WAIT_RESULT"
  STATUS=$(jq -r .status "$WORK/wait.json")
  if [[ $STATUS == 'In Progress' ]]; then
    receipt pending
    printf 'sign-macos: stage=%s pending UUID=%s; continue this candidate (exit 75)\n' "$STAGE" "$SUBMISSION_ID" >&2
    exit 75
  fi
fi
# Always retrieve terminal diagnostics, including rejection; unknown/malformed
# status cannot pass. Pending jobs do not yet have a diagnostic log at Apple.
run notary-log xcrun notarytool log "$SUBMISSION_ID" "$WORK/apple-log.json" "${AUTH[@]}"
jq '{logFormatVersion,jobId,status,statusSummary,statusCode,archiveFilename,uploadDate,sha256,
  ticketContents:(if .ticketContents == null then null else [.ticketContents[] | {path,digestAlgorithm,cdhash,arch}] end),
  issues:(if .issues == null then null else [.issues[] | {severity,code,path,message,docUrl,architecture}] end)}' \
  "$WORK/apple-log.json" > "$WORK/safe-log.json" 2> "$WORK/command.log" || fail 'malformed Apple diagnostic log'
mv "$WORK/safe-log.json" "$EVIDENCE/apple-log.json"
checkpoint
jq -e --arg id "$SUBMISSION_ID" --arg filename "$BASENAME" --arg sha256 "$SUBMITTED_SHA" \
  '.jobId == $id and .archiveFilename == $filename and .sha256 == $sha256' \
  "$EVIDENCE/apple-log.json" >/dev/null 2>&1 || { receipt diagnostic-mismatch; fail 'Apple diagnostic UUID, filename or hash mismatch'; }
if [[ $STATUS == Invalid || $STATUS == Rejected ]]; then receipt rejected; fail 'Apple submission rejected; candidate blocked'; fi
jq -e '.status == "Accepted"' "$EVIDENCE/apple-log.json" >/dev/null 2>&1 || { receipt diagnostic-mismatch; fail 'Apple diagnostic status was not Accepted'; }
[[ $STATUS == Accepted ]] || fail 'Apple submission was not Accepted'
[[ $WAIT_RESULT == 0 ]] || command_failed notary-wait "$WAIT_RESULT"
if [[ $STATE != verified ]]; then receipt accepted; fi

# Work from the original signed upload or the exact previously verified bytes.
# No recovered app/DMG is rebuilt or re-signed.
if [[ $STAGE == app ]]; then
  mkdir "$WORK/restore"
  ARCHIVE="$EVIDENCE/$BASENAME"
  [[ $STATE != verified ]] || ARCHIVE="$EVIDENCE/$VERIFIED_NAME"
  run restore-app ditto -x -k "$ARCHIVE" "$WORK/restore"
  CANDIDATE="$WORK/restore/takedown.app"
  [[ -d $CANDIDATE && ! -L $CANDIDATE && -f $CANDIDATE/Contents/MacOS/takedown && ! -L $CANDIDATE/Contents/MacOS/takedown ]] || fail 'retained archive does not contain takedown.app'
else
  ARCHIVE="$EVIDENCE/$BASENAME"
  [[ $STATE != verified ]] || ARCHIVE="$EVIDENCE/$VERIFIED_NAME"
  run restore-dmg cp "$ARCHIVE" "$CANDIDATE"
fi
verify_signature "$CANDIDATE"
ticket_operation() {
  local operation=$1 result=0
  xcrun stapler "$operation" "$CANDIDATE" > "$WORK/ticket.log" 2>&1 || result=$?
  [[ $result != 0 ]] || return 0
  FAILURE_PHASE="stapler-$operation" FAILURE_EXIT=$result
  # Only the diagnosed CloudKit transport failure is retryable. An Apple
  # rejection, signature error or arbitrary stapler failure is terminal.
  if [[ $STATE != verified && $result == 68 && $TICKET_ATTEMPTS -lt 3 ]] &&
    grep -Eiq '(CloudKit.*(connection|network|unavailable)|connection.*(lost|reset|timed out))' "$WORK/ticket.log"; then
    receipt accepted
    printf 'sign-macos: stage=%s command=%s exit=%s; bounded ticket continuation %s/3\n' "$STAGE" "$FAILURE_PHASE" "$result" "$TICKET_ATTEMPTS" >&2
    exit 75
  fi
  receipt ticket-blocked
  fail 'ticket operation failed permanently or exhausted its bounded continuation budget'
}
if [[ $STATE != verified ]]; then
  [[ $TICKET_ATTEMPTS -lt 3 ]] || { receipt ticket-blocked; fail 'ticket continuation budget exhausted'; }
  TICKET_ATTEMPTS=$((TICKET_ATTEMPTS + 1))
  receipt accepted
  ticket_operation staple
fi
ticket_operation validate
verify_signature "$CANDIDATE"
if [[ $STAGE == app ]]; then
  run gatekeeper spctl --assess --type execute --verbose=2 "$CANDIDATE"
else
  run gatekeeper spctl --assess --type open --context context:primary-signature --verbose=2 "$CANDIDATE"
fi
if [[ $STATE != verified ]]; then
  if [[ $STAGE == app ]]; then
    run verified-archive ditto -c -k --sequesterRsrc --keepParent "$CANDIDATE" "$WORK/$VERIFIED_NAME"
  else
    run verified-archive cp "$CANDIDATE" "$WORK/$VERIFIED_NAME"
  fi
  run retain-verified cp "$WORK/$VERIFIED_NAME" "$EVIDENCE/$VERIFIED_NAME"
  VERIFIED_SHA=$(hash "$EVIDENCE/$VERIFIED_NAME")
fi
# Final credentials/scratch cleanup is acknowledged before a verified receipt
# can be published. Receipt/checkpoint writes do not depend on private scratch.
result=0
rm -rf -- "$WORK" >/dev/null 2>&1 || result=$?
if [[ $result != 0 ]]; then
  FAILURE_PHASE=private-cleanup FAILURE_EXIT=$result
  receipt cleanup-failed
  fail 'private scratch cleanup failed'
fi
receipt verified
# The original credential scratch is gone and custody acknowledged. This new
# restore-only scratch contains code, tickets and command diagnostics, no keys.
WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-restore.XXXXXX")
if [[ $STAGE == app ]]; then
  mkdir "$WORK/final"
  run restore-final-app ditto -x -k "$EVIDENCE/$VERIFIED_NAME" "$WORK/final"
  run remove-provisional rm -rf -- "$TARGET"
  run restore-verified-app ditto "$WORK/final/takedown.app" "$TARGET"
else
  run restore-verified-dmg cp "$EVIDENCE/$VERIFIED_NAME" "$TARGET"
fi
echo 'Developer ID, Apple diagnostics, stapled ticket and Gatekeeper assessment verified.'
