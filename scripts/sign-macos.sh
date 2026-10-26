#!/usr/bin/env bash
# Explicit release operation; raw Apple output and credentials stay in private scratch.
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
for tool in security codesign xcrun openssl ditto plutil jq shasum spctl; do
  command -v "$tool" >/dev/null || fail "missing prerequisite: $tool"
done
xcrun --find notarytool >/dev/null
xcrun --find stapler >/dev/null
[[ $1 != --check ]] || { echo 'macOS signing prerequisites present (credentials not authenticated).'; exit 0; }
[[ ${NOTARY_SOURCE_SHA:-} =~ ^[0-9a-f]{40}$ ]] || fail 'NOTARY_SOURCE_SHA must be an exact source commit'
[[ ${NOTARY_ARCH:-} == arm64 || ${NOTARY_ARCH:-} == amd64 ]] || fail 'NOTARY_ARCH must be arm64 or amd64'
[[ ${RELEASE_VERSION:-} =~ ^[0-9]+\.[0-9]+\.[0-9]+([+-][0-9A-Za-z.+-]+)?$ ]] || fail 'RELEASE_VERSION must be a safe release version'
[[ ${GITHUB_RUN_ID:-0} =~ ^[0-9]+$ && ${GITHUB_RUN_ATTEMPT:-0} =~ ^[0-9]+$ ]] || fail 'invalid CI run identity'
[[ -n ${NOTARY_EVIDENCE_DIR:-} && $NOTARY_EVIDENCE_DIR == /* && ! -L $NOTARY_EVIDENCE_DIR ]] || fail 'NOTARY_EVIDENCE_DIR must be an absolute persistent directory'
TARGET=$1
[[ $TARGET == /* ]] || TARGET="$PWD/$TARGET"
[[ ! -L $TARGET ]] || fail 'target must not be a symlink'
case $TARGET in
  *.app) [[ -d $TARGET && -f $TARGET/Contents/MacOS/takedown && ! -L $TARGET/Contents/MacOS/takedown ]] || fail 'invalid app bundle'; STAGE=app; EXT=zip ;;
  *.dmg) [[ -f $TARGET ]] || fail 'disk image does not exist'; STAGE=dmg; EXT=dmg ;;
  *) fail 'expected an app bundle or DMG' ;;
esac
mkdir -p "$NOTARY_EVIDENCE_DIR"
NOTARY_EVIDENCE_DIR=$(cd -P -- "$NOTARY_EVIDENCE_DIR" && pwd -P)
# The reservation survives failures, including a submit response lost before its UUID.
# A subsequent invocation must reconcile this candidate, never upload it again.
EVIDENCE="$NOTARY_EVIDENCE_DIR/$STAGE"
mkdir "$EVIDENCE" 2>/dev/null || fail "existing $STAGE evidence: reconcile its receipt/Apple history before another submission"
BASENAME="takedown-$RELEASE_VERSION-macos-$NOTARY_ARCH-${NOTARY_SOURCE_SHA:0:12}-${GITHUB_RUN_ID:-0}-${GITHUB_RUN_ATTEMPT:-0}-$STAGE.$EXT"
WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-sign.XXXXXX")
KEYCHAIN="$WORK/signing.keychain-db"
# shellcheck source=scripts/keychain-lifecycle.sh
source "$(dirname -- "$0")/keychain-lifecycle.sh"
KEYCHAIN_MUTATION_STARTED=false
SUBMISSION_ID=''
SUBMITTED_SHA=''
receipt() {
  jq -n --arg state "$1" --arg id "$SUBMISSION_ID" --arg filename "$BASENAME" \
    --arg sha256 "$SUBMITTED_SHA" --arg source "$NOTARY_SOURCE_SHA" \
    --arg version "$RELEASE_VERSION" --arg architecture "$NOTARY_ARCH" --arg stage "$STAGE" \
    --arg run "${GITHUB_RUN_ID:-}" --arg attempt "${GITHUB_RUN_ATTEMPT:-}" \
    '{state:$state,submission_id:(if $id=="" then null else $id end),filename:$filename,submitted_sha256:$sha256,source_commit:$source,version:$version,architecture:$architecture,stage:$stage,ci_run_id:$run,ci_run_attempt:$attempt}' \
    > "$EVIDENCE/.receipt.tmp"
  mv "$EVIDENCE/.receipt.tmp" "$EVIDENCE/receipt.json"
}
cleanup() {
  local result=$?
  trap - EXIT
  if [[ $KEYCHAIN_MUTATION_STARTED == true ]] && ! keychain_restore_search_list; then
    result=1
    receipt cleanup-failed || result=1
  fi
  rm -rf -- "$WORK" || result=1
  exit "$result"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
receipt prepared
printf '%s' "$MACOS_CERTIFICATE_P12_BASE64" | openssl base64 -d -A > "$WORK/certificate.p12"
printf '%s' "$APPLE_API_KEY_BASE64" | openssl base64 -d -A > "$WORK/key.p8"
[[ -s $WORK/certificate.p12 && -s $WORK/key.p8 ]] || fail 'empty decoded signing material'
# Pass the certificate password through the environment, never as a process argument.
# Rewrap privately with an empty import password for security(1)'s argv-only API.
openssl pkcs12 -in "$WORK/certificate.p12" -passin env:MACOS_CERTIFICATE_PASSWORD -nodes -out "$WORK/import.pem" > "$WORK/openssl.log" 2>&1 || fail 'cannot decode signing certificate'
openssl pkcs12 -export -in "$WORK/import.pem" -out "$WORK/import.p12" -passout pass: >> "$WORK/openssl.log" 2>&1 || fail 'cannot prepare signing certificate'
KEYCHAIN_PASSWORD=$(openssl rand -hex 32)
[[ $KEYCHAIN_PASSWORD =~ ^[a-f0-9]{64}$ ]] || fail 'invalid generated keychain password'
case $KEYCHAIN in *'"'*|*\\*|*$'\n'*|*$'\r'*) fail 'unsafe temporary keychain path' ;; esac
keychain_capture_search_list
KEYCHAIN_MUTATION_STARTED=true
{
  # security's interactive command parser accepts one private stdin command.
  # Invoke each secret-bearing command separately so its exit status propagates.
  printf 'create-keychain -p "%s" "%s"\n' "$KEYCHAIN_PASSWORD" "$KEYCHAIN" | security -i
  security set-keychain-settings -lut 21600 "$KEYCHAIN"
  printf 'unlock-keychain -p "%s" "%s"\n' "$KEYCHAIN_PASSWORD" "$KEYCHAIN" | security -i
  security import "$WORK/import.p12" -k "$KEYCHAIN" -f pkcs12 -x -P '' -T /usr/bin/codesign
  printf 'set-key-partition-list -S apple-tool:,apple:,codesign: -s -k "%s" "%s"\n' "$KEYCHAIN_PASSWORD" "$KEYCHAIN" | security -i
  security list-keychains -d user -s "$KEYCHAIN" ${ORIGINAL_KEYCHAINS[@]+"${ORIGINAL_KEYCHAINS[@]}"}
} > "$WORK/security.log" 2>&1
security find-identity -v -p codesigning "$KEYCHAIN" > "$WORK/identities.txt" 2>> "$WORK/security.log"
grep -Fq "\"$MACOS_SIGN_IDENTITY\"" "$WORK/identities.txt" || fail 'configured Developer ID identity is not available in the temporary keychain'
rm -f "$WORK/certificate.p12" "$WORK/import.pem" "$WORK/import.p12"
case $STAGE in
  app)
    codesign --force --timestamp --options runtime --keychain "$KEYCHAIN" --sign "$MACOS_SIGN_IDENTITY" "$TARGET/Contents/MacOS/takedown"
    codesign --force --timestamp --options runtime --keychain "$KEYCHAIN" --sign "$MACOS_SIGN_IDENTITY" "$TARGET"
    codesign --verify --deep --strict --verbose=2 "$TARGET"
    ditto -c -k --sequesterRsrc --keepParent "$TARGET" "$WORK/$BASENAME"
    ;;
  dmg)
    codesign --force --timestamp --keychain "$KEYCHAIN" --sign "$MACOS_SIGN_IDENTITY" "$TARGET"
    codesign --verify --strict --verbose=2 "$TARGET"
    cp "$TARGET" "$WORK/$BASENAME"
    ;;
esac
codesign --display --verbose=4 "$TARGET" 2> "$WORK/signature.txt"
grep -qx "TeamIdentifier=$APPLE_TEAM_ID" "$WORK/signature.txt" || fail 'signed artifact has the wrong TeamIdentifier'
# Retain exactly the signed upload bytes, separately from the stapled final target.
cp "$WORK/$BASENAME" "$EVIDENCE/$BASENAME"
SUBMITTED_SHA=$(shasum -a 256 "$EVIDENCE/$BASENAME" | cut -d ' ' -f 1)
receipt submission-unknown
AUTH=(--key "$WORK/key.p8" --key-id "$APPLE_API_KEY_ID" --issuer "$APPLE_API_ISSUER")
SUBMIT_RESULT=0
xcrun notarytool submit "$EVIDENCE/$BASENAME" "${AUTH[@]}" --no-wait --output-format plist > "$WORK/submit.plist" 2> "$WORK/notary.stderr" || SUBMIT_RESULT=$?
SUBMISSION_ID=$(plutil -extract id raw -o - "$WORK/submit.plist" 2>/dev/null || true)
[[ $SUBMISSION_ID =~ ^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$ ]] || { SUBMISSION_ID=''; receipt submission-unknown; fail 'submission UUID unavailable; reconcile Apple history against the retained filename/hash before retrying'; }
receipt submitted
printf 'Notarization submission (%s): %s\n' "$STAGE" "$SUBMISSION_ID"
[[ $SUBMIT_RESULT == 0 ]] || fail 'submission returned a failure with a UUID; reconcile the retained submission before retrying'
WAIT_RESULT=0
xcrun notarytool wait "$SUBMISSION_ID" "${AUTH[@]}" --timeout 30m --output-format plist > "$WORK/wait.plist" 2> "$WORK/notary.stderr" || WAIT_RESULT=$?
WAIT_ID=$(plutil -extract id raw -o - "$WORK/wait.plist" 2>/dev/null || true)
STATUS=$(plutil -extract status raw -o - "$WORK/wait.plist" 2>/dev/null || true)
receipt wait-unresolved
# Logs are required even after acceptance. Only allowlisted diagnostics leave scratch.
LOG_RESULT=0
xcrun notarytool log "$SUBMISSION_ID" "$WORK/apple-log.json" "${AUTH[@]}" > "$WORK/log.stdout" 2> "$WORK/notary.stderr" || LOG_RESULT=$?
[[ $LOG_RESULT == 0 && -s $WORK/apple-log.json ]] || fail 'Apple diagnostic log unavailable; retain receipt and reconcile the existing submission'
jq '{logFormatVersion,jobId,status,statusSummary,statusCode,archiveFilename,uploadDate,sha256,
  ticketContents:(if .ticketContents == null then null else [.ticketContents[] | {path,digestAlgorithm,cdhash,arch}] end),
  issues:(if .issues == null then null else [.issues[] | {severity,code,path,message,docUrl,architecture}] end)}' \
  "$WORK/apple-log.json" > "$EVIDENCE/.apple-log.tmp" || fail 'malformed Apple diagnostic log'
mv "$EVIDENCE/.apple-log.tmp" "$EVIDENCE/apple-log.json"
jq -e --arg id "$SUBMISSION_ID" --arg filename "$BASENAME" --arg sha256 "$SUBMITTED_SHA" \
  '.jobId == $id and (.archiveFilename == null or .archiveFilename == $filename) and (.sha256 == null or .sha256 == $sha256)' "$EVIDENCE/apple-log.json" >/dev/null || fail 'Apple diagnostic identity, filename or hash mismatch'
[[ $WAIT_ID == "$SUBMISSION_ID" ]] || fail 'Apple wait response submission ID mismatch'
[[ $WAIT_RESULT == 0 && $STATUS == Accepted ]] || { receipt not-accepted; fail 'Apple submission was not Accepted; inspect retained diagnostics and reconcile before retrying'; }
jq -e '.status == "Accepted"' "$EVIDENCE/apple-log.json" >/dev/null || { receipt not-accepted; fail 'Apple diagnostic status was not Accepted'; }
receipt accepted
xcrun stapler staple "$TARGET"
xcrun stapler validate "$TARGET"
if [[ $STAGE == app ]]; then
  codesign --verify --deep --strict --verbose=2 "$TARGET"
  spctl --assess --type execute --verbose=2 "$TARGET"
else
  codesign --verify --strict --verbose=2 "$TARGET"
  spctl --assess --type open --context context:primary-signature --verbose=2 "$TARGET"
fi
receipt verified
echo 'Developer ID, Apple diagnostics, stapled ticket and Gatekeeper assessment verified.'
