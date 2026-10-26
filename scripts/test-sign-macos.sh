#!/usr/bin/env bash
# Execute the real signer with synthetic tools; never touch an Apple account/keychain.
set -euo pipefail
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
TEST_WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-sign-tests.XXXXXX")
cleanup_tests() {
  local result=$?
  if [[ $result != 0 ]]; then
    printf 'FAIL: case=%s exit=%s\n' "${MOCK_STATE:-setup}" "$result" >&2
    [[ ! -f ${MOCK_STATE:-}/output ]] || cat "$MOCK_STATE/output" >&2
    [[ ! -f ${MOCK_STATE:-}/package-output ]] || cat "$MOCK_STATE/package-output" >&2
  fi
  rm -rf -- "$TEST_WORK"
}
trap cleanup_tests EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
mkdir "$TEST_WORK/tools" "$TEST_WORK/scratch"
export REAL_OPENSSL REAL_DITTO REAL_PLUTIL REAL_RM
# shellcheck source=scripts/macos-certificate-import.sh
source "$ROOT/scripts/macos-certificate-import.sh"
# shellcheck source=scripts/macos-ticket-policy.sh
source "$ROOT/scripts/macos-ticket-policy.sh"
# Classify only complete known constants, never arbitrary numeric/private text.
while IFS='|' read -r message status cause; do
  printf 'security> security: SecKeychainItemImport: %s\n' "$message" > "$TEST_WORK/diagnostic-input"
  diagnostic=$(macos_certificate_import_diagnostic "$TEST_WORK/diagnostic-input" 2>&1)
  [[ $diagnostic == "sign-macos: keychain-import security_status=$status cause=$cause (allowlisted classification; raw output withheld)" ]]
done <<'CASES'
The user name or passphrase you entered is not correct.|-25293|authentication-failed
MAC verification failed during PKCS12 import (wrong password?)|-25264|pkcs12-mac-verification-failed
Unable to decode the provided data.|-26275|decode-failed
Unknown format in import.|-25257|unknown-format
User interaction is not allowed.|-25308|interaction-not-allowed
One or more parameters passed to a function were not valid.|-50|invalid-parameters
The user name or passphrase you entered is not correct. fixture-p12-password|unknown|unclassified
arbitrary -25293 MUST-NOT-RETAIN|unknown|unclassified
CASES
echo 'PASS: import diagnostics emit only exact allowlisted Security constants/status mappings'
REAL_OPENSSL=$(command -v openssl)
REAL_DITTO=$(command -v ditto)
REAL_PLUTIL=$(command -v plutil)
REAL_RM=$(command -v rm)
cat > "$TEST_WORK/tools/rm" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
if [[ ${MOCK_SCENARIO:-} == private-cleanup-failure && $1 == -rf && -d ${!#} && ${!#} == */takedown-sign.* && ! -f $MOCK_STATE/private-remove-failed ]]; then
  touch "$MOCK_STATE/private-remove-failed"
  exit 55
fi
exec "$REAL_RM" "$@"
MOCK
cat > "$TEST_WORK/tools/uname" <<'MOCK'
#!/usr/bin/env bash
echo Darwin
MOCK
cat > "$TEST_WORK/tools/openssl" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
for argument in "$@"; do [[ $argument != "$MACOS_CERTIFICATE_PASSWORD" ]] || exit 91; done
case $1 in
  base64) exec "$REAL_OPENSSL" "$@" ;;
  rand) printf '%064d\n' 0 ;;
  pkcs12)
    [[ $MOCK_SCENARIO != certificate-decode-failure ]] || { echo 'fixture-p12-password raw error' >&2; exit 42; }
    while (($#)); do
      if [[ $1 == -out ]]; then printf 'fixture certificate\n' > "$2"; exit 0; fi
      shift
    done
    exit 1
    ;;
  *) exit 99 ;;
esac
MOCK
cat > "$TEST_WORK/tools/security" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
for argument in "$@"; do [[ $argument != "$MACOS_CERTIFICATE_PASSWORD" ]] || exit 91; done
for argument in "$@"; do [[ $argument != "$(printf '%064d' 0)" ]] || exit 92; done
case $1 in
  -i)
    [[ $MOCK_SCENARIO != security-stdin-failure ]] || exit 77
    IFS= read -r line
    printf 'private-stdin-command\n' >> "$MOCK_STATE/private-commands"
    [[ $MOCK_SCENARIO != keychain-unlock-failure || $line != unlock-keychain* ]] || { echo 'fixture-p12-password raw security error' >&2; exit 78; }
    [[ $MOCK_SCENARIO != keychain-partitions-failure || $line != set-key-partition-list* ]] || { echo 'fixture-p12-password raw security error' >&2; exit 79; }
    if [[ $line == import* ]]; then
      pattern='^import "[^"]+" -k "[^"]+" -x -P "[a-f0-9]{64}" -T /usr/bin/codesign$'
      [[ $line =~ $pattern ]] || exit 99
      case $MOCK_SCENARIO in
        keychain-import-auth-failure) echo 'security: SecKeychainItemImport: The user name or passphrase you entered is not correct.' >&2; echo 'fixture-p12-password MUST-NOT-RETAIN' >&2; exit 1 ;;
        keychain-import-unknown-failure) echo 'security: SecKeychainItemImport: arbitrary fixture-p12-password MUST-NOT-RETAIN -25293' >&2; exit 1 ;;
      esac
      exit 0
    fi
    pattern='^(create-keychain|unlock-keychain) -p "([a-f0-9]{64})" "([^"]+)"$'
    if [[ $line =~ $pattern ]]; then
      if [[ ${BASH_REMATCH[1]} == create-keychain ]]; then
        touch "${BASH_REMATCH[3]}"
        printf '    "%s"\n' "${BASH_REMATCH[3]}" >> "$MOCK_STATE/keychains"
      fi
    else
      pattern='^set-key-partition-list -S apple-tool:,apple:,codesign: -s -k "[a-f0-9]{64}" "[^"]+"$'
      [[ $line =~ $pattern ]] || exit 99
    fi
    ;;
  default-keychain) printf '    "/fixture/login.keychain-db"\n' ;;
  create-keychain)
    [[ $MOCK_SCENARIO != security-stdin-failure ]] || exit 77
    touch "${!#}"
    printf '    "%s"\n' "${!#}" >> "$MOCK_STATE/keychains"
    ;;
  delete-keychain)
    [[ $MOCK_SCENARIO != cleanup-failure ]] || exit 1
    rm "$2"
    ;;
  list-keychains)
    if [[ $# == 3 ]]; then cat "$MOCK_STATE/keychains"
    elif [[ $4 == -s ]]; then
      shift 4
      : > "$MOCK_STATE/keychains"
      for path in "$@"; do printf '    "%s"\n' "$path" >> "$MOCK_STATE/keychains"; done
    else exit 99
    fi
    ;;
  set-keychain-settings|unlock-keychain|import|set-key-partition-list) : ;;
  find-identity) printf '  1) FIXTURE "%s"\n' "$MACOS_SIGN_IDENTITY" ;;
  *) exit 99 ;;
esac
MOCK
cat > "$TEST_WORK/tools/codesign" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
[[ $MOCK_SCENARIO != signature-failure ]] || exit 1
printf '%s\n' "$*" >> "$MOCK_STATE/codesigns"
if [[ $1 == --force ]]; then
  item=${!#}
  if [[ -d $item ]]; then printf 'signed bundle\n' > "$item/.synthetic-signature"; else printf 'signed code\n' >> "$item"; fi
fi
if [[ $1 == --display ]]; then
  printf 'TeamIdentifier=%s\nAuthority=%s\nTimestamp=fixture secure timestamp\nCodeDirectory v=20500 flags=0x10000(runtime)\n' "$APPLE_TEAM_ID" "$MACOS_SIGN_IDENTITY" >&2
fi
MOCK
cat > "$TEST_WORK/tools/ditto" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
exec "$REAL_DITTO" "$@"
MOCK
cat > "$TEST_WORK/tools/plutil" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
[[ $1 == -extract ]] || exec "$REAL_PLUTIL" "$@"
[[ $3 == raw && $4 == -o && $5 == - ]] || exit 99
if [[ $2 == CFBundleShortVersionString ]]; then echo 0.0.1; else
  jq -er --arg key "$2" '.[$key] // empty' "$6"
fi
MOCK
cat > "$TEST_WORK/tools/lipo" <<'MOCK'
#!/usr/bin/env bash
echo arm64
MOCK
cat > "$TEST_WORK/tools/spctl" <<'MOCK'
#!/usr/bin/env bash
[[ $MOCK_SCENARIO != gatekeeper-failure ]]
MOCK
cat > "$TEST_WORK/tools/xcrun" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
id=12345678-1234-1234-1234-123456789abc
other=87654321-1234-1234-1234-123456789abc
if [[ $1 == --find ]]; then echo "/fixture/$2"; exit 0; fi
if [[ $1 == stapler ]]; then
  printf '%s\n' "$2" >> "$MOCK_STATE/tickets"
  if [[ $2 == validate && $MOCK_SCENARIO == validation-transient ]]; then echo 'CloudKit network connection lost' >&2; echo 'MUST-NOT-RETAIN fixture-p12-password' >&2; exit 68; fi
  if [[ $2 == validate && $MOCK_SCENARIO == validation-mutate-archive ]]; then printf 'changed\n' >> "$NOTARY_EVIDENCE_DIR/app/stapled-app.zip"; fi
  if [[ $2 == staple && $MOCK_SCENARIO == ticket-transient ]]; then echo 'CloudKit network connection lost' >&2; exit 68; fi
  [[ $2 != staple || $MOCK_SCENARIO != staple-failure ]] || exit 68
  [[ $2 != validate || $MOCK_SCENARIO != validation-failure ]] || exit 68
  if [[ $2 == staple && -f $3 ]]; then printf 'stapled ticket\n' >> "$3"; fi
  exit 0
fi
[[ $1 == notarytool ]] || exit 99
if [[ $2 == submit && $3 == *.dmg ]]; then id=22345678-1234-1234-1234-123456789abc; fi
if [[ $2 == info || $2 == wait || $2 == log ]]; then
  [[ $3 == "$id" || $3 == 22345678-1234-1234-1234-123456789abc ]] || exit 99
  id=$3
fi
case $2 in
  submit)
    printf '%s\n' "$3" >> "$MOCK_STATE/submits"
    [[ $* == *'--no-wait'* && $* != *'--wait'* ]] || exit 99
    basename "$3" > "$MOCK_STATE/filename"
    shasum -a 256 "$3" | cut -d ' ' -f 1 > "$MOCK_STATE/hash"
    cp "$MOCK_STATE/filename" "$MOCK_STATE/$id.filename"
    cp "$MOCK_STATE/hash" "$MOCK_STATE/$id.hash"
    if [[ $MOCK_SCENARIO == unknown-submission ]]; then echo '{}'; exit 1; fi
    jq -n --arg id "$id" '{id:$id}'
    ;;
  history)
    printf 'history\n' >> "$MOCK_STATE/queries"
    count=1
    [[ $MOCK_SCENARIO != history-absent ]] || count=0
    [[ $MOCK_SCENARIO != history-ambiguous ]] || count=2
    jq -n --arg id "$id" --arg name "$(cat "$MOCK_STATE/filename")" --argjson count "$count" '{history:[range($count) | {id:$id,name:$name,status:"Accepted"}]}'
    ;;
  info)
    printf 'info %s\n' "$3" >> "$MOCK_STATE/queries"
    [[ $3 == "$id" ]] || exit 99
    status='In Progress'
    [[ $MOCK_SCENARIO != accepted-info ]] || status=Accepted
    jq -n --arg id "$id" --arg name "$(cat "$MOCK_STATE/$id.filename")" --arg status "$status" '{id:$id,name:$name,status:$status}'
    ;;
  wait)
    printf 'wait %s\n' "$3" >> "$MOCK_STATE/queries"
    [[ $3 == "$id" && $* == *'--timeout 30m'* ]] || exit 99
    status=Accepted
    if [[ $MOCK_SCENARIO == package-dmg-pending && $(cat "$MOCK_STATE/filename") == *-dmg.dmg ]]; then status='In Progress'; fi
    [[ $MOCK_SCENARIO != rejected ]] || status=Invalid
    [[ $MOCK_SCENARIO != pending ]] || status='In Progress'
    [[ $MOCK_SCENARIO != wait-id-mismatch ]] || id=$other
    jq -n --arg id "$id" --arg status "$status" '{id:$id,status:$status}'
    [[ $MOCK_SCENARIO != wait-timeout && $status != 'In Progress' ]] || exit 69
    ;;
  log)
    [[ $3 == "$id" ]] || exit 99
    [[ $MOCK_SCENARIO != missing-log ]] || exit 1
    if [[ $MOCK_SCENARIO == malformed-log ]]; then printf '{broken' > "$4"; exit 0; fi
    filename=$(cat "$MOCK_STATE/$id.filename")
    hash=$(cat "$MOCK_STATE/$id.hash")
    status=Accepted
    [[ $MOCK_SCENARIO != log-id-mismatch ]] || id=$other
    [[ $MOCK_SCENARIO != log-name-mismatch ]] || filename=wrong.zip
    [[ $MOCK_SCENARIO != log-hash-mismatch ]] || hash=wrong
    [[ $MOCK_SCENARIO != rejected && $MOCK_SCENARIO != log-status-mismatch ]] || status=Invalid
    [[ $MOCK_SCENARIO != pending ]] || status='In Progress'
    jq -n --arg id "$id" --arg filename "$filename" --arg sha256 "$hash" --arg status "$status" \
      '{logFormatVersion:1,jobId:$id,status:$status,archiveFilename:$filename,sha256:$sha256,privateField:"MUST-NOT-RETAIN",issues:[{severity:"warning",message:"fixture warning",privateField:"MUST-NOT-RETAIN"}]}' > "$4"
    ;;
  *) exit 99 ;;
esac
MOCK
chmod +x "$TEST_WORK/tools/"*
export PATH="$TEST_WORK/tools:$PATH" TMPDIR="$TEST_WORK/scratch"
export MACOS_CERTIFICATE_P12_BASE64=Zml4dHVyZQ== MACOS_CERTIFICATE_PASSWORD=fixture-p12-password
export MACOS_SIGN_IDENTITY='Developer ID Application: Fixture (ABCDEFGHIJ)' APPLE_TEAM_ID=ABCDEFGHIJ
export APPLE_API_KEY_BASE64=Zml4dHVyZQ== APPLE_API_KEY_ID=FIXTUREKEY APPLE_API_ISSUER=fixture-issuer
export NOTARY_SOURCE_SHA=0123456789abcdef0123456789abcdef01234567 NOTARY_ARCH=arm64 RELEASE_VERSION=0.0.1
export GITHUB_RUN_ID=123 GITHUB_RUN_ATTEMPT=1 MOCK_SCENARIO MOCK_STATE NOTARY_EVIDENCE_DIR
CASE_SEQUENCE=0
new_case() {
  CASE_SEQUENCE=$((CASE_SEQUENCE + 1))
  MOCK_SCENARIO=$1
  MOCK_STATE="$TEST_WORK/$CASE_SEQUENCE-$1-$2"
  mkdir -p "$MOCK_STATE/takedown.app/Contents/MacOS" "$MOCK_STATE/evidence"
  printf 'fixture executable\n' > "$MOCK_STATE/takedown.app/Contents/MacOS/takedown"
  printf 'fixture signed disk image\n' > "$MOCK_STATE/takedown.dmg"
  printf '    "/fixture/login.keychain-db"\n' > "$MOCK_STATE/keychains"
  cp "$MOCK_STATE/keychains" "$MOCK_STATE/original-keychains"
  NOTARY_EVIDENCE_DIR="$MOCK_STATE/evidence"
  TARGET="$MOCK_STATE/takedown.$2"
}
run_signer() { bash "$ROOT/scripts/sign-macos.sh" "$TARGET" > "$MOCK_STATE/output" 2>&1; }
assert_private() {
  [[ -z $(find "$TMPDIR" -mindepth 1 -maxdepth 1 -type d -name 'takedown-sign.*' -print) ]]
  local result=0
  local private_paths=("$NOTARY_EVIDENCE_DIR")
  [[ ! -f $MOCK_STATE/output ]] || private_paths+=("$MOCK_STATE/output")
  [[ ! -f $MOCK_STATE/package-output ]] || private_paths+=("$MOCK_STATE/package-output")
  grep -REq "fixture-p12-password|MUST-NOT-RETAIN|$(printf '%064d' 0)" "${private_paths[@]}" || result=$?
  [[ $result == 1 ]] || { echo 'private material escaped or privacy inspection failed' >&2; exit 1; }
}
verified_recovery_repro() {
  new_case success verified-recovery
  TARGET="$MOCK_STATE/takedown.app"
  run_signer
  retained=$(jq -r .verified_sha256 "$NOTARY_EVIDENCE_DIR/app/receipt.json")
  MOCK_SCENARIO=validation-transient
  result=0; run_signer || result=$?
  printf 'Verified-restore transient validation result: %s (expected 75)\n' "$result"
  if [[ $result != 75 ]]; then cat "$MOCK_STATE/output"; cat "$NOTARY_EVIDENCE_DIR/app/receipt.json"; fi
  [[ $result == 75 ]]
  jq -e --arg sha "$retained" '.state=="validation-pending" and .verified_sha256==$sha and .failure_class=="cloudkit-transport" and .validation_attempts==1' "$NOTARY_EVIDENCE_DIR/app/receipt.json" >/dev/null
  MOCK_SCENARIO=success
  run_signer
  jq -e --arg sha "$retained" '.state=="verified" and .verified_sha256==$sha and .validation_attempts==2 and .ticket_attempts==1 and .failure_phase=="" and .failure_class==""' "$NOTARY_EVIDENCE_DIR/app/receipt.json" >/dev/null
  [[ $(grep -c '^staple$' "$MOCK_STATE/tickets") == 1 && $(wc -l < "$MOCK_STATE/submits" | tr -d ' ') == 1 && $(grep -c -- --force "$MOCK_STATE/codesigns") == 2 ]]
  assert_private
  echo 'PASS: verified restore yields safely on diagnosed validation failure and recovers exact bytes without restapling/signing/upload'
}
verified_recovery_repro
[[ ${SIGN_MACOS_TEST_GROUP:-all} != verified-recovery ]] || exit 0
ticket_recovery_contracts() {
  local message expected seed original pin queries signs imports staples case_id stage parent snapshot
  while IFS='|' read -r message expected; do
    printf '%s\n' "$message" > "$TEST_WORK/ticket-diagnostic"
    [[ $(macos_ticket_failure_class "$TEST_WORK/ticket-diagnostic" 68) == "$expected" ]]
    [[ $(macos_ticket_failure_class "$TEST_WORK/ticket-diagnostic" 1) == unclassified ]]
  done <<'CASES'
CloudKit network connection lost|cloudkit-transport
CloudKit connection reset|cloudkit-transport
CloudKit connection timed out|cloudkit-transport
CloudKit service unavailable|cloudkit-transport
connection lost|unclassified
CloudKit query failed due to (null)|unclassified
CloudKit network connection lost fixture-p12-password|unclassified
68|unclassified
CASES
  printf 'CloudKit network connection lost\nInvalid signature\n' > "$TEST_WORK/ticket-diagnostic"
  [[ $(macos_ticket_failure_class "$TEST_WORK/ticket-diagnostic" 68) == unclassified ]]
  echo 'PASS: ticket classification uses safe allowlisted constants; exit 68, generic transport and contradictory diagnostics are terminal'

  # Generate the actual verified -> terminal seam, then remove only extension
  # fields to reproduce the legacy schema-2 receipt shape. No real metadata used.
  new_case success legacy-blocked
  TARGET="$MOCK_STATE/takedown.app"
  run_signer
  MOCK_SCENARIO=validation-failure
  if run_signer; then exit 1; fi
  jq -e '.state=="ticket-blocked" and .failure_class=="unclassified" and .validation_attempts==1 and .verified_sha256!=""' "$NOTARY_EVIDENCE_DIR/app/receipt.json" >/dev/null
  jq 'del(.failure_class,.validation_attempts,.requalification_parent_sha256)' "$NOTARY_EVIDENCE_DIR/app/receipt.json" > "$MOCK_STATE/legacy.json"
  cp "$MOCK_STATE/legacy.json" "$NOTARY_EVIDENCE_DIR/app/receipt.json"
  seed="$MOCK_STATE/blocked-seed"
  cp -R "$NOTARY_EVIDENCE_DIR" "$seed"
  original="$MOCK_STATE/legacy.json"
  queries=$(wc -l < "$MOCK_STATE/queries")
  signs=$(grep -c -- --force "$MOCK_STATE/codesigns")
  imports=$(wc -l < "$MOCK_STATE/private-commands")
  staples=$(grep -c '^staple$' "$MOCK_STATE/tickets")
  pin=$(shasum -a 256 "$original" | awk '{print $1}')
  MOCK_SCENARIO=success
  if run_signer; then echo 'unapproved terminal receipt reopened' >&2; exit 1; fi
  cmp "$original" "$NOTARY_EVIDENCE_DIR/app/receipt.json"
  [[ $(wc -l < "$MOCK_STATE/queries") == "$queries" ]]
  export NOTARY_REQUALIFY_RECEIPT_SHA256
  for case_id in wrong-pin malformed-pin source run attempt arch stage filename version uuid missing-verified symlink-verified verified-digest submitted-digest non-validation non-68 classified rejected diagnostic missing-log exhausted-ticket exhausted-validation zero-budget parent; do
    NOTARY_EVIDENCE_DIR="$MOCK_STATE/refuse-$case_id"
    cp -R "$seed" "$NOTARY_EVIDENCE_DIR"
    case $case_id in
      source) message='.source_commit="aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"' ;;
      run) message='.ci_run_id="999"' ;;
      attempt) message='.ci_run_attempt="2"' ;;
      arch) message='.architecture="amd64"' ;;
      stage) message='.stage="dmg"' ;;
      filename) message='.filename="wrong.zip"' ;;
      version) message='.version="0.0.2"' ;;
      uuid) message='.submission_id="bad"' ;;
      non-validation) message='.failure_phase="gatekeeper"' ;;
      non-68) message='.failure_exit="1"' ;;
      classified) message='.failure_class="cloudkit-transport"' ;;
      exhausted-ticket) message='.ticket_attempts=3' ;;
      exhausted-validation) message='.validation_attempts=3' ;;
      zero-budget) message='.validation_attempts=0' ;;
      parent) message='.requalification_parent_sha256="aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"' ;;
      *) message='.' ;;
    esac
    jq "$message" "$NOTARY_EVIDENCE_DIR/app/receipt.json" > "$MOCK_STATE/mutated.json"
    mv "$MOCK_STATE/mutated.json" "$NOTARY_EVIDENCE_DIR/app/receipt.json"
    NOTARY_REQUALIFY_RECEIPT_SHA256=$(shasum -a 256 "$NOTARY_EVIDENCE_DIR/app/receipt.json" | awk '{print $1}')
    case $case_id in
      wrong-pin) NOTARY_REQUALIFY_RECEIPT_SHA256=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa ;;
      malformed-pin) NOTARY_REQUALIFY_RECEIPT_SHA256=not-a-digest ;;
      missing-verified) rm "$NOTARY_EVIDENCE_DIR/app/stapled-app.zip" ;;
      symlink-verified) rm "$NOTARY_EVIDENCE_DIR/app/stapled-app.zip"; ln -s "$seed/app/stapled-app.zip" "$NOTARY_EVIDENCE_DIR/app/stapled-app.zip" ;;
      verified-digest) printf 'corrupt\n' >> "$NOTARY_EVIDENCE_DIR/app/stapled-app.zip" ;;
      submitted-digest) printf 'corrupt\n' >> "$NOTARY_EVIDENCE_DIR/app/$(jq -r .filename "$original")" ;;
      rejected|diagnostic)
        message='.status="Rejected"'; [[ $case_id != diagnostic ]] || message='.sha256="wrong"'
        jq "$message" "$NOTARY_EVIDENCE_DIR/app/apple-log.json" > "$MOCK_STATE/log.json"
        mv "$MOCK_STATE/log.json" "$NOTARY_EVIDENCE_DIR/app/apple-log.json" ;;
      missing-log) rm "$NOTARY_EVIDENCE_DIR/app/apple-log.json" ;;
    esac
    cp "$NOTARY_EVIDENCE_DIR/app/receipt.json" "$MOCK_STATE/before.json"
    if run_signer; then echo "unsafe requalification: $case_id" >&2; exit 1; fi
    cmp "$MOCK_STATE/before.json" "$NOTARY_EVIDENCE_DIR/app/receipt.json"
    [[ $(wc -l < "$MOCK_STATE/queries") == "$queries" && ! -d $NOTARY_EVIDENCE_DIR/app/.lock ]]
    assert_private
  done
  echo 'PASS: legacy terminal receipt remains byte-identical without exact approval; 24 pin/identity/bytes/phase/Accepted-log/budget refusals make no provider query'
  NOTARY_EVIDENCE_DIR="$MOCK_STATE/pinned-success"
  cp -R "$seed" "$NOTARY_EVIDENCE_DIR"
  NOTARY_REQUALIFY_RECEIPT_SHA256=$pin
  run_signer
  jq -e --arg pin "$pin" '.state=="verified" and .requalification_parent_sha256==$pin and .validation_attempts==2 and .ticket_attempts==1 and .failure_class==""' "$NOTARY_EVIDENCE_DIR/app/receipt.json" >/dev/null
  parent="$NOTARY_EVIDENCE_DIR/app/blocked-receipt-$pin.json"
  cmp "$original" "$parent"
  cmp "$seed/app/stapled-app.zip" "$NOTARY_EVIDENCE_DIR/app/stapled-app.zip"
  [[ $(grep -c '^staple$' "$MOCK_STATE/tickets") == "$staples" && $(grep -c -- --force "$MOCK_STATE/codesigns") == "$signs" && $(wc -l < "$MOCK_STATE/private-commands") == "$imports" && $(wc -l < "$MOCK_STATE/submits" | tr -d ' ') == 1 ]]
  # A later normal recheck spends no new recovery budget, even with the pin unset.
  unset NOTARY_REQUALIFY_RECEIPT_SHA256
  run_signer
  jq -e '.validation_attempts==2 and .state=="verified"' "$NOTARY_EVIDENCE_DIR/app/receipt.json" >/dev/null
  queries=$(wc -l < "$MOCK_STATE/queries")
  snapshot="$MOCK_STATE/verified.json"
  cp "$NOTARY_EVIDENCE_DIR/app/receipt.json" "$snapshot"
  printf 'tampered\n' >> "$parent"
  if run_signer; then echo 'tampered parent reused' >&2; exit 1; fi
  cmp "$snapshot" "$NOTARY_EVIDENCE_DIR/app/receipt.json"
  [[ $(wc -l < "$MOCK_STATE/queries") == "$queries" ]]
  assert_private
  echo 'PASS: exact pin preserves terminal parent, UUID, archives and counters; subsequent rechecks audit the parent without restaple, signing, key import or submit'

  for case_id in validation-transient gatekeeper-failure signature-failure log-hash-mismatch rejected private-cleanup-failure validation-mutate-archive; do
    NOTARY_EVIDENCE_DIR="$MOCK_STATE/pinned-$case_id"
    cp -R "$seed" "$NOTARY_EVIDENCE_DIR"
    export NOTARY_REQUALIFY_RECEIPT_SHA256=$pin
    MOCK_SCENARIO=$case_id
    result=0; run_signer || result=$?
    [[ $result != 0 ]]
    cmp "$original" "$NOTARY_EVIDENCE_DIR/app/blocked-receipt-$pin.json"
    jq -e --arg pin "$pin" '.state!="verified" and .requalification_parent_sha256==$pin and .validation_attempts==2' "$NOTARY_EVIDENCE_DIR/app/receipt.json" >/dev/null
    if [[ $case_id == validation-transient ]]; then
      [[ $result == 75 ]]
      MOCK_SCENARIO=success
      run_signer
      jq -e '.state=="verified" and .validation_attempts==3' "$NOTARY_EVIDENCE_DIR/app/receipt.json" >/dev/null
    else
      queries=$(wc -l < "$MOCK_STATE/queries")
      NOTARY_REQUALIFY_RECEIPT_SHA256=$(shasum -a 256 "$NOTARY_EVIDENCE_DIR/app/receipt.json" | awk '{print $1}')
      MOCK_SCENARIO=success
      if run_signer; then echo "terminal second pin bypassed: $case_id" >&2; exit 1; fi
      [[ $(wc -l < "$MOCK_STATE/queries") == "$queries" ]]
    fi
    assert_private
  done
  unset NOTARY_REQUALIFY_RECEIPT_SHA256
  echo 'PASS: pinned recovery still requires full Accepted diagnostics, signatures, ticket, Gatekeeper, cleanup and exact verified digest; a second pin cannot reset it'

  cat > "$TEST_WORK/requalification-checkpoint.sh" <<'HOOK'
#!/usr/bin/env bash
set -euo pipefail
state=$(jq -r .state "$1/app/receipt.json")
if [[ $state == "$MOCK_FAIL_RECOVERY_CHECKPOINT" ]]; then exit 88; fi
HOOK
  export MOCK_FAIL_RECOVERY_CHECKPOINT NOTARY_CHECKPOINT_HOOK NOTARY_REQUALIFY_RECEIPT_SHA256=$pin
  NOTARY_CHECKPOINT_HOOK="$TEST_WORK/requalification-checkpoint.sh"
  for case_id in ticket-blocked validation-pending; do
    NOTARY_EVIDENCE_DIR="$MOCK_STATE/interrupted-$case_id"
    cp -R "$seed" "$NOTARY_EVIDENCE_DIR"
    MOCK_FAIL_RECOVERY_CHECKPOINT=$case_id MOCK_SCENARIO=success
    queries=$(wc -l < "$MOCK_STATE/queries")
    if run_signer; then echo 'failed pin checkpoint passed' >&2; exit 1; fi
    [[ $(wc -l < "$MOCK_STATE/queries") == "$queries" ]]
    cmp "$original" "$NOTARY_EVIDENCE_DIR/app/blocked-receipt-$pin.json"
    if [[ $case_id == ticket-blocked ]]; then
      cmp "$original" "$NOTARY_EVIDENCE_DIR/app/receipt.json"
    else
      jq -e '.state=="validation-pending" and .validation_attempts==2' "$NOTARY_EVIDENCE_DIR/app/receipt.json" >/dev/null
    fi
    MOCK_FAIL_RECOVERY_CHECKPOINT=never
    run_signer
    jq -e '.state=="verified" and .validation_attempts>=2 and .validation_attempts<=3' "$NOTARY_EVIDENCE_DIR/app/receipt.json" >/dev/null
    assert_private
  done
  unset MOCK_FAIL_RECOVERY_CHECKPOINT NOTARY_CHECKPOINT_HOOK NOTARY_REQUALIFY_RECEIPT_SHA256
  echo 'PASS: parent/pending checkpoint failures stop before provider operations and restart without losing original failure or resetting reserved attempts'

  for stage in app dmg; do
    new_case success "verified-budget-$stage"
    TARGET="$MOCK_STATE/takedown.$stage"
    run_signer
    MOCK_SCENARIO=validation-transient
    for attempt in 1 2 3; do
      result=0; run_signer || result=$?
      if [[ $attempt -lt 3 ]]; then [[ $result == 75 ]]; else [[ $result == 1 ]]; fi
      jq -e --argjson count "$attempt" '.validation_attempts==$count and .ticket_attempts==1' "$NOTARY_EVIDENCE_DIR/$stage/receipt.json" >/dev/null
    done
    jq -e '.state=="ticket-blocked" and .failure_class=="cloudkit-transport"' "$NOTARY_EVIDENCE_DIR/$stage/receipt.json" >/dev/null
    queries=$(wc -l < "$MOCK_STATE/queries")
    if run_signer; then exit 1; fi
    [[ $(wc -l < "$MOCK_STATE/queries") == "$queries" && $(grep -c '^staple$' "$MOCK_STATE/tickets") == 1 && $(wc -l < "$MOCK_STATE/submits" | tr -d ' ') == 1 ]]
    assert_private
  done
  echo 'PASS: app and DMG validation recovery budgets persist across restarts and terminal exhaustion never becomes verified'
}
ticket_recovery_contracts
[[ ${SIGN_MACOS_TEST_GROUP:-all} != ticket-recovery ]] || exit 0
for stage in app dmg; do
  new_case success "$stage"
  run_signer
  jq -e --arg stage "$stage" --arg source "$NOTARY_SOURCE_SHA" '.state=="verified" and .stage==$stage and .source_commit==$source and .version=="0.0.1" and .architecture=="arm64" and .ci_run_id=="123" and .ci_run_attempt=="1" and (.submitted_sha256 | test("^[0-9a-f]{64}$"))' "$NOTARY_EVIDENCE_DIR/$stage/receipt.json" >/dev/null
  filename=$(jq -r .filename "$NOTARY_EVIDENCE_DIR/$stage/receipt.json")
  [[ $filename == "takedown-0.0.1-macos-arm64-0123456789ab-123-1-$stage."* ]]
  [[ -f $NOTARY_EVIDENCE_DIR/$stage/$filename ]]
  [[ $(shasum -a 256 "$NOTARY_EVIDENCE_DIR/$stage/$filename" | cut -d ' ' -f 1) == "$(jq -r .submitted_sha256 "$NOTARY_EVIDENCE_DIR/$stage/receipt.json")" ]]
  jq -e '.issues[0].message=="fixture warning"' "$NOTARY_EVIDENCE_DIR/$stage/apple-log.json" >/dev/null
  cmp "$MOCK_STATE/keychains" "$MOCK_STATE/original-keychains"
  retained=$(jq -r .verified_sha256 "$NOTARY_EVIDENCE_DIR/$stage/receipt.json")
  run_signer
  [[ $(jq -r .verified_sha256 "$NOTARY_EVIDENCE_DIR/$stage/receipt.json") == "$retained" ]]
  [[ $(grep -c -- --force "$MOCK_STATE/codesigns") == $([[ $stage == app ]] && echo 2 || echo 1) ]]
  [[ $(wc -l < "$MOCK_STATE/submits" | tr -d ' ') == 1 ]]
  assert_private
  echo "PASS: $stage descriptive upload, bound receipt, retained signed bytes/log warnings, cleanup, duplicate prevention"
done
for scenario in rejected pending wait-timeout unknown-submission missing-log malformed-log log-id-mismatch log-name-mismatch log-hash-mismatch log-status-mismatch wait-id-mismatch staple-failure validation-failure gatekeeper-failure cleanup-failure signature-failure security-stdin-failure; do
  new_case "$scenario" app
  if run_signer; then cat "$MOCK_STATE/output"; echo "expected rejection: $scenario" >&2; exit 1; fi
  [[ -f $NOTARY_EVIDENCE_DIR/app/receipt.json ]]
  [[ ! -e $TARGET/.synthetic-signature ]]
  [[ $(cat "$TARGET/Contents/MacOS/takedown") == 'fixture executable' ]]
  if [[ $scenario == unknown-submission ]]; then
    jq -e '.state=="submission-unknown" and .submission_id==null and (.submitted_sha256|length)==64' "$NOTARY_EVIDENCE_DIR/app/receipt.json" >/dev/null
  elif [[ $scenario == cleanup-failure ]]; then
    jq -e '.state=="cleanup-failed"' "$NOTARY_EVIDENCE_DIR/app/receipt.json" >/dev/null
  elif [[ $scenario == rejected || $scenario == pending || $scenario == wait-timeout ]]; then
    if [[ $scenario != pending ]]; then [[ -s $NOTARY_EVIDENCE_DIR/app/apple-log.json ]]; fi
  fi
  if [[ $scenario != signature-failure && $scenario != security-stdin-failure && $scenario != cleanup-failure ]]; then
    [[ $(wc -l < "$MOCK_STATE/submits" | tr -d ' ') == 1 ]]
    if [[ $scenario == unknown-submission ]]; then
      run_signer # A unique history match may recover, but must not upload again.
    elif run_signer; then echo "unexpected retry success: $scenario" >&2; exit 1; fi
    [[ $(wc -l < "$MOCK_STATE/submits" | tr -d ' ') == 1 ]]
  fi
  [[ $scenario == cleanup-failure ]] || cmp "$MOCK_STATE/keychains" "$MOCK_STATE/original-keychains"
  assert_private
  echo "PASS: $scenario fails closed and retains evidence without resubmission"
done

new_case rejected terminal-rejection
TARGET="$MOCK_STATE/takedown.app"
if run_signer; then exit 1; fi
jq -e '.state=="rejected"' "$NOTARY_EVIDENCE_DIR/app/receipt.json" >/dev/null
queries=$(wc -l < "$MOCK_STATE/queries")
if run_signer; then exit 1; fi
[[ $(wc -l < "$MOCK_STATE/queries") == "$queries" ]]
assert_private
echo 'PASS: terminal rejection retains a blocking receipt and makes no further provider query'

new_case ticket-transient bounded-ticket
TARGET="$MOCK_STATE/takedown.app"
for attempt in 1 2 3; do
  result=0; run_signer || result=$?
  if [[ $attempt -lt 3 ]]; then [[ $result == 75 ]]; else [[ $result == 1 ]]; fi
  jq -e --argjson attempts "$attempt" '.ticket_attempts==$attempts' "$NOTARY_EVIDENCE_DIR/app/receipt.json" >/dev/null
done
jq -e '.state=="ticket-blocked"' "$NOTARY_EVIDENCE_DIR/app/receipt.json" >/dev/null
queries=$(wc -l < "$MOCK_STATE/queries")
if run_signer; then exit 1; fi
[[ $(wc -l < "$MOCK_STATE/queries") == "$queries" && $(wc -l < "$MOCK_STATE/submits" | tr -d ' ') == 1 && $(grep -c -- --force "$MOCK_STATE/codesigns") == 2 ]]
assert_private
new_case ticket-transient recovered-ticket
TARGET="$MOCK_STATE/takedown.app"
result=0; run_signer || result=$?
[[ $result == 75 ]]
MOCK_SCENARIO=success
run_signer
jq -e '.state=="verified" and .ticket_attempts==2' "$NOTARY_EVIDENCE_DIR/app/receipt.json" >/dev/null
[[ $(wc -l < "$MOCK_STATE/submits" | tr -d ' ') == 1 && $(grep -c -- --force "$MOCK_STATE/codesigns") == 2 ]]
assert_private
echo 'PASS: only diagnosed CloudKit failures continue, with a persistent three-attempt budget and no signing/upload replay'

new_case private-cleanup-failure private-cleanup
TARGET="$MOCK_STATE/takedown.app"
if run_signer; then exit 1; fi
jq -e '.state=="cleanup-failed" and .failure_phase=="private-cleanup" and .failure_exit=="55"' "$NOTARY_EVIDENCE_DIR/app/receipt.json" >/dev/null
[[ ! -e $TARGET/.synthetic-signature && $(cat "$TARGET/Contents/MacOS/takedown") == 'fixture executable' ]]
assert_private
echo 'PASS: private credential-scratch cleanup failure blocks verified receipt and target promotion'

# Continuation of one finite pending UUID must neither import keys nor re-sign.
new_case pending app
result=0; run_signer || result=$?
[[ $result == 75 ]]
jq -e '.state=="pending"' "$NOTARY_EVIDENCE_DIR/app/receipt.json" >/dev/null
MOCK_SCENARIO=success
run_signer
[[ $(wc -l < "$MOCK_STATE/submits" | tr -d ' ') == 1 ]]
[[ $(grep -c -- --force "$MOCK_STATE/codesigns") == 2 ]]
[[ $(wc -l < "$MOCK_STATE/private-commands" | tr -d ' ') == 4 ]]
[[ $(grep -c 'info 12345678-1234-1234-1234-123456789abc' "$MOCK_STATE/queries") == 2 ]]
assert_private
echo 'PASS: pending exits 75, same UUID continues to Accepted without signing or submitting again'

cat > "$TEST_WORK/yield-checkpoint.sh" <<'HOOK'
#!/usr/bin/env bash
set -euo pipefail
jq -er .state "$1/app/receipt.json" >> "$MOCK_STATE/yield-checkpoints"
HOOK
new_case success app
export NOTARY_CHECKPOINT_HOOK="$TEST_WORK/yield-checkpoint.sh"
export NOTARY_YIELD_AFTER_SUBMIT=true
result=0; run_signer || result=$?
[[ $result == 75 && ! -e $MOCK_STATE/queries && $(tail -1 "$MOCK_STATE/yield-checkpoints") == submitted ]]
jq -e '.state=="submitted" and .submission_id=="12345678-1234-1234-1234-123456789abc"' "$NOTARY_EVIDENCE_DIR/app/receipt.json" >/dev/null
run_signer # Flag deliberately stays true: an existing UUID must not yield again.
[[ $(wc -l < "$MOCK_STATE/submits" | tr -d ' ') == 1 ]]
unset NOTARY_YIELD_AFTER_SUBMIT NOTARY_CHECKPOINT_HOOK
echo 'PASS: qualification yield occurs only after a NEW UUID checkpoint; retained UUID resumes with flag still set'

for field in schema source_commit version architecture stage ci_run_id ci_run_attempt filename submitted_sha256 team_id signing_identity verified_filename; do
  new_case pending "binding-$field"
  TARGET="$MOCK_STATE/takedown.app"
  result=0; run_signer || result=$?
  [[ $result == 75 ]]
  jq --arg field "$field" '.[$field]="mismatch"' "$NOTARY_EVIDENCE_DIR/app/receipt.json" > "$MOCK_STATE/changed.json"
  mv "$MOCK_STATE/changed.json" "$NOTARY_EVIDENCE_DIR/app/receipt.json"
  if run_signer; then echo "unbound receipt reused: $field" >&2; exit 1; fi
  [[ $(wc -l < "$MOCK_STATE/submits" | tr -d ' ') == 1 ]]
  [[ $(wc -l < "$MOCK_STATE/queries" | tr -d ' ') == 2 ]]
  [[ $(grep -c -- --force "$MOCK_STATE/codesigns") == 2 ]]
done
echo 'PASS: every candidate/filename/hash/receipt-schema binding rejects reuse before remote operations'

for scenario in history-absent history-ambiguous; do
  new_case unknown-submission "$scenario"
  TARGET="$MOCK_STATE/takedown.app"
  if run_signer; then exit 1; fi
  MOCK_SCENARIO=$scenario
  if run_signer; then echo "ambiguous upload recovered: $scenario" >&2; exit 1; fi
  [[ $(wc -l < "$MOCK_STATE/submits" | tr -d ' ') == 1 ]]
  jq -e '.state=="submission-unknown" and .submission_id==null' "$NOTARY_EVIDENCE_DIR/app/receipt.json" >/dev/null
  assert_private
done
echo 'PASS: unknown upload with absent/ambiguous history blocks without resubmission'

cat > "$TEST_WORK/checkpoint.sh" <<'HOOK'
#!/usr/bin/env bash
set -euo pipefail
[[ $1 == "$NOTARY_EVIDENCE_DIR" ]]
state=$(jq -r .state "$1/app/receipt.json")
printf '%s\n' "$state" >> "$MOCK_STATE/checkpoints"
[[ $state != "${MOCK_CHECKPOINT_FAIL_STATE:-never}" ]] || { echo 'fixture-p12-password raw checkpoint error' >&2; exit 88; }
# At upload intent, prove exact signed bytes are already bound and present.
if [[ $state == submission-unknown ]]; then
  file=$(jq -r .filename "$1/app/receipt.json")
  [[ $(shasum -a 256 "$1/app/$file" | cut -d ' ' -f 1) == "$(jq -r .submitted_sha256 "$1/app/receipt.json")" ]]
fi
HOOK
export NOTARY_CHECKPOINT_HOOK="$TEST_WORK/checkpoint.sh" MOCK_CHECKPOINT_FAIL_STATE
for state in prepared submission-unknown submitted verified; do
  new_case success "checkpoint-$state"
  TARGET="$MOCK_STATE/takedown.app"
  MOCK_CHECKPOINT_FAIL_STATE=$state
  cp "$TARGET/Contents/MacOS/takedown" "$MOCK_STATE/original-binary"
  if run_signer; then echo "failed checkpoint passed: $state" >&2; exit 1; fi
  cmp "$TARGET/Contents/MacOS/takedown" "$MOCK_STATE/original-binary"
  if [[ $state == prepared || $state == submission-unknown ]]; then [[ ! -e $MOCK_STATE/submits ]]; else
    [[ $(wc -l < "$MOCK_STATE/submits" | tr -d ' ') == 1 ]]
  fi
  [[ $state != submitted || ! -e $MOCK_STATE/queries ]]
  grep -Fq 'command=checkpoint exit=88' "$MOCK_STATE/output"
  MOCK_CHECKPOINT_FAIL_STATE=never
  if [[ $state == submission-unknown ]]; then
    # An unacknowledged intent remains unknown. No history match means no upload.
    printf '%s\n' "$(jq -r .filename "$NOTARY_EVIDENCE_DIR/app/receipt.json")" > "$MOCK_STATE/filename"
    MOCK_SCENARIO=history-absent
    if run_signer; then exit 1; fi
    [[ ! -e $MOCK_STATE/submits ]]
  else
    run_signer
    [[ $(wc -l < "$MOCK_STATE/submits" | tr -d ' ') == 1 ]]
  fi
  assert_private
done
unset NOTARY_CHECKPOINT_HOOK MOCK_CHECKPOINT_FAIL_STATE
echo 'PASS: checkpoint failure stops at prepared/intent/UUID/verified boundaries, with no target promotion'

for scenario in certificate-decode-failure security-stdin-failure keychain-unlock-failure keychain-partitions-failure keychain-import-auth-failure keychain-import-unknown-failure; do
  new_case "$scenario" diagnostics
  TARGET="$MOCK_STATE/takedown.app"
  if run_signer; then exit 1; fi
  case $scenario in
    certificate-decode-failure) phase=certificate-decode; code=42 ;;
    security-stdin-failure) phase=keychain-create; code=77 ;;
    keychain-unlock-failure) phase=keychain-unlock; code=78 ;;
    keychain-partitions-failure) phase=keychain-partitions; code=79 ;;
    keychain-import-*) phase=keychain-import; code=1 ;;
  esac
  grep -Fq "stage=app command=$phase exit=$code" "$MOCK_STATE/output"
  jq -e --arg phase "$phase" --arg code "$code" '.failure_phase==$phase and .failure_exit==$code' "$NOTARY_EVIDENCE_DIR/app/receipt.json" >/dev/null
  [[ ! -e $MOCK_STATE/submits ]]
  cmp "$MOCK_STATE/keychains" "$MOCK_STATE/original-keychains"
  case $scenario in
    keychain-import-auth-failure) grep -Fq 'security_status=-25293 cause=authentication-failed' "$MOCK_STATE/output" ;;
    keychain-import-unknown-failure) grep -Fq 'security_status=unknown cause=unclassified' "$MOCK_STATE/output" ;;
  esac
  assert_private
  if [[ $scenario == keychain-import-auth-failure ]]; then
    jq -e '.state=="prepared" and .submission_id==null and .submitted_sha256=="" and .verified_sha256==""' "$NOTARY_EVIDENCE_DIR/app/receipt.json" >/dev/null
    MOCK_SCENARIO=success
    run_signer
    jq -e --arg source "$NOTARY_SOURCE_SHA" '.state=="verified" and .source_commit==$source and .ci_run_id=="123" and .ci_run_attempt=="1"' "$NOTARY_EVIDENCE_DIR/app/receipt.json" >/dev/null
    [[ $(wc -l < "$MOCK_STATE/submits" | tr -d ' ') == 1 ]]
    cmp "$MOCK_STATE/keychains" "$MOCK_STATE/original-keychains"
    assert_private
  fi
done
echo 'PASS: private pre-submission failures retain safe diagnostics, clean keychains and resume the same prepared receipt'

# Exercise the real packager + signer in a disposable project, using the same
# synthetic Apple/keychain tools. This proves publication remains transactional.
cat > "$TEST_WORK/tools/go" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
case $1 in
  env) case $2 in GOVERSION) echo go1.26.10 ;; GOHOSTOS) echo darwin ;; GOHOSTARCH) echo arm64 ;; *) exit 99 ;; esac ;;
  list) echo v2.8.2 ;;
  run) exit 0 ;;
  version)
    printf '%s: go1.26.10\n\tdep\tfyne.io/fyne/v2\tv2.8.2\th1:fixture\n\tbuild\tCGO_ENABLED=1\n\tbuild\tGOOS=darwin\n\tbuild\tGOARCH=arm64\n' "$3"
    ;;
  *) exit 99 ;;
esac
MOCK
cat > "$TEST_WORK/tools/hdiutil" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
printf 'create\n' >> "$MOCK_STATE/dmg-creates"
printf 'synthetic signed disk image\n' > "${!#}"
MOCK
cat > "$TEST_WORK/tools/xattr" <<'MOCK'
#!/usr/bin/env bash
echo 'quarantine stripping prohibited' >&2
exit 96
MOCK
chmod +x "$TEST_WORK/tools/go" "$TEST_WORK/tools/hdiutil" "$TEST_WORK/tools/xattr"
FIXTURE_ROOT="$TEST_WORK/package-project"
mkdir -p "$FIXTURE_ROOT/scripts" "$FIXTURE_ROOT/assets" "$FIXTURE_ROOT/build/macos"
cp "$ROOT/scripts/package-common.sh" "$ROOT/scripts/package-macos.sh" "$ROOT/scripts/sign-macos.sh" "$ROOT/scripts/keychain-lifecycle.sh" "$ROOT/scripts/macos-certificate-import.sh" "$ROOT/scripts/macos-ticket-policy.sh" "$FIXTURE_ROOT/scripts/"
cp "$ROOT/build/macos/Info.plist" "$FIXTURE_ROOT/build/macos/"
printf 'fixture icon\n' > "$FIXTURE_ROOT/assets/icon.icns"
for file in README.md PRIVACY.md PRIVACY.html LICENSE NOTICE THIRD_PARTY_NOTICES.md; do
  printf 'fixture legal documentation\n' > "$FIXTURE_ROOT/$file"
done
printf '#!/usr/bin/env bash\necho "takedown version 0.0.1"\n' > "$FIXTURE_ROOT/binary"
chmod +x "$FIXTURE_ROOT/binary"
run_package() {
  MODE=release VERSION=0.0.1 BINARY="$FIXTURE_ROOT/binary" OUT_DIR="$MOCK_STATE/packages" \
    bash "$FIXTURE_ROOT/scripts/package-macos.sh" > "$MOCK_STATE/package-output" 2>&1
}
new_case package-dmg-pending package
result=0; run_package || result=$?
[[ $result == 75 && ! -e $MOCK_STATE/packages ]]
[[ $(wc -l < "$MOCK_STATE/submits" | tr -d ' ') == 2 ]]
[[ $(wc -l < "$MOCK_STATE/dmg-creates" | tr -d ' ') == 1 ]]
jq -e '.state=="verified"' "$NOTARY_EVIDENCE_DIR/app/receipt.json" >/dev/null
jq -e '.state=="pending"' "$NOTARY_EVIDENCE_DIR/dmg/receipt.json" >/dev/null
# Simulate a new runner restoring only allowlisted evidence and original binding.
cp -R "$NOTARY_EVIDENCE_DIR" "$MOCK_STATE/restored-evidence"
NOTARY_EVIDENCE_DIR="$MOCK_STATE/restored-evidence"
MOCK_SCENARIO=accepted-info
run_package
cmp "$NOTARY_EVIDENCE_DIR/app/stapled-app.zip" "$MOCK_STATE/packages/takedown-0.0.1-macos-arm64-release.zip"
cmp "$NOTARY_EVIDENCE_DIR/dmg/stapled.dmg" "$MOCK_STATE/packages/takedown-0.0.1-macos-arm64-release.dmg"
[[ $(wc -l < "$MOCK_STATE/submits" | tr -d ' ') == 2 ]]
[[ $(wc -l < "$MOCK_STATE/dmg-creates" | tr -d ' ') == 1 ]]
[[ $(grep -c -- --force "$MOCK_STATE/codesigns") == 3 ]]
(cd "$MOCK_STATE/packages" && shasum -a 256 -c SHA256SUMS) >/dev/null
assert_private
echo 'PASS: fresh packager restores exact accepted ZIP/DMG, no DMG rebuild/re-sign/re-upload; pending promotes nothing'

for scenario in rejected log-hash-mismatch cleanup-failure gatekeeper-failure; do
  new_case "$scenario" package-failure
  if run_package; then echo "packaging failure passed: $scenario" >&2; exit 1; fi
  [[ ! -e $MOCK_STATE/packages ]]
  assert_private
done
echo 'PASS: rejected/mismatched/cleanup/assessment failures publish no package directory'
echo 'PASS: synthetic signer and packager contracts; no Apple account, real signing, or keychain access'
