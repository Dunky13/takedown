#!/usr/bin/env bash
# Execute the real signer with synthetic tools; never touch an Apple account/keychain.
set -euo pipefail
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
TEST_WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-sign-tests.XXXXXX")
trap 'rm -rf -- "$TEST_WORK"' EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
mkdir "$TEST_WORK/tools" "$TEST_WORK/scratch"
export REAL_OPENSSL
REAL_OPENSSL=$(command -v openssl)
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
  create-keychain)
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
if [[ $1 == --display ]]; then printf 'TeamIdentifier=%s\n' "$APPLE_TEAM_ID" >&2; fi
MOCK
cat > "$TEST_WORK/tools/ditto" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
printf 'synthetic signed app archive\n' > "${!#}"
MOCK
cat > "$TEST_WORK/tools/plutil" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
[[ $1 == -extract && $3 == raw && $4 == -o && $5 == - ]] || exit 99
jq -er --arg key "$2" '.[$key] // empty' "$6"
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
  [[ $2 != staple || $MOCK_SCENARIO != staple-failure ]] || exit 68
  [[ $2 != validate || $MOCK_SCENARIO != validation-failure ]] || exit 68
  if [[ $2 == staple && -f $3 ]]; then printf 'stapled ticket\n' >> "$3"; fi
  exit 0
fi
[[ $1 == notarytool ]] || exit 99
case $2 in
  submit)
    printf '%s\n' "$3" >> "$MOCK_STATE/submits"
    [[ $* == *'--no-wait'* && $* != *'--wait'* ]] || exit 99
    basename "$3" > "$MOCK_STATE/filename"
    shasum -a 256 "$3" | cut -d ' ' -f 1 > "$MOCK_STATE/hash"
    if [[ $MOCK_SCENARIO == unknown-submission ]]; then echo '{}'; exit 1; fi
    jq -n --arg id "$id" '{id:$id}'
    ;;
  wait)
    [[ $3 == "$id" && $* == *'--timeout 30m'* ]] || exit 99
    status=Accepted
    [[ $MOCK_SCENARIO != rejected ]] || status=Invalid
    [[ $MOCK_SCENARIO != pending ]] || status='In Progress'
    [[ $MOCK_SCENARIO != wait-id-mismatch ]] || id=$other
    jq -n --arg id "$id" --arg status "$status" '{id:$id,status:$status}'
    [[ $MOCK_SCENARIO != wait-timeout && $MOCK_SCENARIO != pending ]] || exit 69
    ;;
  log)
    [[ $3 == "$id" ]] || exit 99
    [[ $MOCK_SCENARIO != missing-log ]] || exit 1
    if [[ $MOCK_SCENARIO == malformed-log ]]; then printf '{broken' > "$4"; exit 0; fi
    filename=$(cat "$MOCK_STATE/filename")
    hash=$(cat "$MOCK_STATE/hash")
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
new_case() {
  MOCK_SCENARIO=$1
  MOCK_STATE="$TEST_WORK/$1-$2"
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
  if grep -REq "fixture-p12-password|MUST-NOT-RETAIN|$(printf '%064d' 0)" "$NOTARY_EVIDENCE_DIR" "$MOCK_STATE/output"; then
    echo 'private material escaped' >&2; exit 1
  fi
}
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
  if run_signer; then echo 'duplicate submission accepted' >&2; exit 1; fi
  [[ $(wc -l < "$MOCK_STATE/submits" | tr -d ' ') == 1 ]]
  assert_private
  echo "PASS: $stage descriptive upload, bound receipt, retained signed bytes/log warnings, cleanup, duplicate prevention"
done
for scenario in rejected pending wait-timeout unknown-submission missing-log malformed-log log-id-mismatch log-name-mismatch log-hash-mismatch log-status-mismatch wait-id-mismatch staple-failure validation-failure gatekeeper-failure cleanup-failure signature-failure security-stdin-failure; do
  new_case "$scenario" app
  if run_signer; then cat "$MOCK_STATE/output"; echo "expected rejection: $scenario" >&2; exit 1; fi
  [[ -f $NOTARY_EVIDENCE_DIR/app/receipt.json ]]
  if [[ $scenario == unknown-submission ]]; then
    jq -e '.state=="submission-unknown" and .submission_id==null and (.submitted_sha256|length)==64' "$NOTARY_EVIDENCE_DIR/app/receipt.json" >/dev/null
  elif [[ $scenario == cleanup-failure ]]; then
    jq -e '.state=="cleanup-failed"' "$NOTARY_EVIDENCE_DIR/app/receipt.json" >/dev/null
  elif [[ $scenario == rejected || $scenario == pending || $scenario == wait-timeout ]]; then
    [[ -s $NOTARY_EVIDENCE_DIR/app/apple-log.json ]]
  fi
  if [[ $scenario != signature-failure && $scenario != security-stdin-failure ]]; then
    [[ $(wc -l < "$MOCK_STATE/submits" | tr -d ' ') == 1 ]]
    if run_signer; then echo "retry resubmitted $scenario" >&2; exit 1; fi
    [[ $(wc -l < "$MOCK_STATE/submits" | tr -d ' ') == 1 ]]
  fi
  [[ $scenario == cleanup-failure ]] || cmp "$MOCK_STATE/keychains" "$MOCK_STATE/original-keychains"
  assert_private
  echo "PASS: $scenario fails closed and retains evidence without resubmission"
done
echo 'PASS: synthetic signer contracts; no Apple account, real signing, or keychain access'
