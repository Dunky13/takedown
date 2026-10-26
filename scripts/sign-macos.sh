#!/usr/bin/env bash
# Explicit, networked release operation. Never invoked by PR checks/local packaging.
# Disable tracing before reading any secret, including when called with bash -x.
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
for tool in security codesign xcrun openssl ditto plutil; do
  command -v "$tool" >/dev/null || fail "missing prerequisite: $tool"
done
xcrun --find notarytool >/dev/null
xcrun --find stapler >/dev/null
[[ $1 != --check ]] || { echo 'macOS signing prerequisites present (credentials not authenticated).'; exit 0; }
TARGET=$1
[[ $TARGET == /* ]] || TARGET="$PWD/$TARGET"
case $TARGET in
  *.app) [[ -d $TARGET && -f $TARGET/Contents/MacOS/takedown ]] || fail 'invalid app bundle' ;;
  *.dmg) [[ -f $TARGET ]] || fail 'disk image does not exist' ;;
  *) fail 'expected an app bundle or DMG' ;;
esac
WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-sign.XXXXXX")
KEYCHAIN="$WORK/signing.keychain-db"
# shellcheck source=scripts/keychain-lifecycle.sh
source "$(dirname -- "$0")/keychain-lifecycle.sh"
KEYCHAIN_MUTATION_STARTED=false
cleanup() {
  result=$?
  trap - EXIT
  if [[ $KEYCHAIN_MUTATION_STARTED == true ]]; then keychain_restore_search_list || result=1; fi
  rm -rf -- "$WORK" || result=1
  exit "$result"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
printf '%s' "$MACOS_CERTIFICATE_P12_BASE64" | openssl base64 -d -A > "$WORK/certificate.p12"
printf '%s' "$APPLE_API_KEY_BASE64" | openssl base64 -d -A > "$WORK/key.p8"
[[ -s $WORK/certificate.p12 && -s $WORK/key.p8 ]] || fail 'empty decoded signing material'
KEYCHAIN_PASSWORD=$(openssl rand -hex 32)
keychain_capture_search_list
KEYCHAIN_MUTATION_STARTED=true
security create-keychain -p "$KEYCHAIN_PASSWORD" "$KEYCHAIN"
security set-keychain-settings -lut 21600 "$KEYCHAIN"
security unlock-keychain -p "$KEYCHAIN_PASSWORD" "$KEYCHAIN"
# A dedicated keychain, nonextractable key, codesign-only trusted application.
# create-keychain transiently changes the user search list; cleanup restores and
# verifies its captured list. Use an isolated signing runner with no other writers.
# The default keychain is never explicitly reassigned.
security import "$WORK/certificate.p12" -k "$KEYCHAIN" -f pkcs12 -x \
  -P "$MACOS_CERTIFICATE_PASSWORD" -T /usr/bin/codesign >/dev/null
security set-key-partition-list -S apple-tool:,apple:,codesign: -s \
  -k "$KEYCHAIN_PASSWORD" "$KEYCHAIN" >/dev/null
rm -f "$WORK/certificate.p12"
case $TARGET in
  *.app)
    codesign --force --timestamp --options runtime --keychain "$KEYCHAIN" \
      --sign "$MACOS_SIGN_IDENTITY" "$TARGET/Contents/MacOS/takedown"
    codesign --force --timestamp --options runtime --keychain "$KEYCHAIN" \
      --sign "$MACOS_SIGN_IDENTITY" "$TARGET"
    codesign --verify --deep --strict --verbose=2 "$TARGET"
    ditto -c -k --sequesterRsrc --keepParent "$TARGET" "$WORK/submit.zip"
    SUBMIT="$WORK/submit.zip"
    ;;
  *.dmg)
    codesign --force --timestamp --keychain "$KEYCHAIN" --sign "$MACOS_SIGN_IDENTITY" "$TARGET"
    codesign --verify --strict --verbose=2 "$TARGET"
    SUBMIT="$TARGET"
    ;;
esac
codesign --display --verbose=4 "$TARGET" 2> "$WORK/signature.txt"
grep -qx "TeamIdentifier=$APPLE_TEAM_ID" "$WORK/signature.txt" || fail 'signed artifact has the wrong TeamIdentifier'
# Plist output allows an explicit Accepted check without a JSON parser dependency.
# A timeout remains a failure; do not resubmit automatically or claim acceptance.
NOTARY_RESULT=0
xcrun notarytool submit "$SUBMIT" --key "$WORK/key.p8" --key-id "$APPLE_API_KEY_ID" \
  --issuer "$APPLE_API_ISSUER" --wait --timeout 30m --output-format plist > "$WORK/notary.plist" || NOTARY_RESULT=$?
# Preserve the submission handle in the caller/CI log even on a finite wait failure.
# Only this non-secret identifier is printed; never dump signing material or argv.
if [[ -s $WORK/notary.plist ]]; then
  SUBMISSION_ID=$(plutil -extract id raw -o - "$WORK/notary.plist" 2>/dev/null || true)
  [[ -z $SUBMISSION_ID ]] || printf 'Notarization submission: %s\n' "$SUBMISSION_ID"
fi
[[ $NOTARY_RESULT == 0 ]] || fail "notarytool exited $NOTARY_RESULT; reconcile any submitted request before retrying"
STATUS=$(plutil -extract status raw -o - "$WORK/notary.plist")
[[ $STATUS == Accepted ]] || fail "notarization was not Accepted: $STATUS"
xcrun stapler staple "$TARGET"
xcrun stapler validate "$TARGET"
echo 'Developer ID signature and stapled notarization ticket verified.'
