#!/usr/bin/env bash
# Native, synthetic-only regression. Creates/deletes an owned isolated keychain.
{ set +x; } 2>/dev/null
set -euo pipefail
umask 077
[[ $(uname -s) == Darwin ]] || { echo 'macOS required' >&2; exit 1; }
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-import-test.XXXXXX")
KEYCHAIN="$WORK/synthetic.keychain-db"
# shellcheck source=scripts/keychain-lifecycle.sh
source "$ROOT/scripts/keychain-lifecycle.sh"
# shellcheck source=scripts/macos-certificate-import.sh
source "$ROOT/scripts/macos-certificate-import.sh"
MUTATED=false
cleanup() {
  local result=$?
  trap - EXIT
  if [[ $MUTATED == true ]]; then
    if keychain_restore_search_list; then
      echo 'PASS: exact user search list/default preserved; owned keychain absent'
    else result=1; fi
  fi
  rm -rf -- "$WORK" || result=1
  if [[ ! -e $WORK ]]; then echo 'PASS: private synthetic material removed'; else result=1; fi
  exit "$result"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
# Bound native commands without exposing their stdin (including passwords).
security() { /usr/bin/perl -e 'alarm 15; exec @ARGV' /usr/bin/security "$@"; }
keychain_capture_search_list
openssl version
sw_vers -productVersion
openssl req -x509 -newkey rsa:2048 -nodes -days 1 \
  -subj /CN=Takedown-Synthetic-Import-Test \
  -addext basicConstraints=critical,CA:FALSE -addext keyUsage=critical,digitalSignature \
  -addext extendedKeyUsage=codeSigning \
  -keyout "$WORK/key.pem" -out "$WORK/cert.pem" > "$WORK/private.log" 2>&1
cat "$WORK/key.pem" "$WORK/cert.pem" > "$WORK/input.pem"
openssl rand -hex 32 > "$WORK/source-password"
openssl pkcs12 -export -in "$WORK/input.pem" -out "$WORK/supplied.p12" \
  -passout "file:$WORK/source-password" > "$WORK/private.log" 2>&1
openssl pkcs12 -in "$WORK/supplied.p12" -passin "file:$WORK/source-password" \
  -nodes -out "$WORK/import.pem" > "$WORK/private.log" 2>&1
macos_certificate_rewrap "$WORK/import.pem" "$WORK/import.p12" "$WORK/transfer-password" > "$WORK/private.log" 2>&1
PASSWORD=$(openssl rand -hex 32)
[[ $PASSWORD =~ ^[a-f0-9]{64}$ ]]
case $KEYCHAIN in *'"'*|*\\*|*$'\n'*|*$'\r'*) exit 1 ;; esac
MUTATED=true
printf 'create-keychain -p "%s" "%s"\n' "$PASSWORD" "$KEYCHAIN" | security -i > "$WORK/private.log" 2>&1
[[ -f $KEYCHAIN ]]
chmod 600 "$KEYCHAIN"
security set-keychain-settings -lut 21600 "$KEYCHAIN" > "$WORK/private.log" 2>&1
printf 'unlock-keychain -p "%s" "%s"\n' "$PASSWORD" "$KEYCHAIN" | security -i > "$WORK/private.log" 2>&1
result=0
macos_certificate_import "$WORK/import.p12" "$KEYCHAIN" "$WORK/transfer-password" > "$WORK/private.log" 2>&1 || result=$?
if [[ $result != 0 ]]; then
  macos_certificate_import_diagnostic "$WORK/private.log"
  printf 'FAIL: synthetic production export/import chain exit=%s (private output withheld)\n' "$result" >&2
  exit 1
fi
security find-identity -p codesigning "$KEYCHAIN" > "$WORK/identities" 2> "$WORK/private.log"
grep -Fq 'Takedown-Synthetic-Import-Test' "$WORK/identities"
security find-key -t private "$KEYCHAIN" > "$WORK/key-metadata" 2> "$WORK/private.log"
# SecKey.h: kSecKeyExtractable == 16; inspect attributes only, never export a key.
if grep -Eq '0x00000010[[:space:]]+<uint32>=0x00000000' "$WORK/key-metadata"; then
  echo 'PASS: owned synthetic key metadata reports nonextractable'
else
  if grep -Eq '0x00000010[[:space:]]+<uint32>=0x00000001' "$WORK/key-metadata"; then
    echo 'FAIL: native importer made the owned synthetic key extractable' >&2
  fi
  echo 'FAIL: owned synthetic key metadata does not report nonextractable' >&2
  exit 1
fi
echo 'PASS: production rewrap/import helpers retain a synthetic certificate and private-key identity'
# A known-bad password must fail through private stdin and classify without
# retaining native output. No provider credentials or signing/submission occurs.
openssl rand -hex 32 > "$WORK/wrong-password"
result=0
macos_certificate_import "$WORK/import.p12" "$KEYCHAIN" "$WORK/wrong-password" > "$WORK/private.log" 2>&1 || result=$?
[[ $result != 0 ]]
macos_certificate_import_diagnostic "$WORK/private.log" > "$WORK/diagnostic" 2>&1
grep -Eq 'security_status=-(25293|25264) cause=(authentication-failed|pkcs12-mac-verification-failed)' "$WORK/diagnostic"
echo 'PASS: native wrong-password failure returns nonzero and emits only an allowlisted classification'
