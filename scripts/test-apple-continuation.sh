#!/usr/bin/env bash
# Wrapper integration using synthetic local package/assessment tools only.
set -euo pipefail
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
TEST_WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-continuation-tests.XXXXXX")
trap 'rm -rf -- "$TEST_WORK"' EXIT
mkdir -p "$TEST_WORK/repo/scripts" "$TEST_WORK/state"
cp "$ROOT/scripts/continue-apple-release.sh" "$ROOT/scripts/fresh-apple-acceptance.sh" "$ROOT/scripts/apple-stages-verified.sh" "$ROOT/scripts/macos-ticket-policy.sh" "$TEST_WORK/repo/scripts/"
export FIXTURE_STATE="$TEST_WORK/state" VERSION=1.2.3 RELEASE_VERSION=1.2.3 NOTARY_ARCH=arm64
export NOTARY_SOURCE_SHA=1111111111111111111111111111111111111111 GITHUB_RUN_ID=123 GITHUB_RUN_ATTEMPT=1 ORIGIN_RUN=123 ORIGIN_ATTEMPT=1
export NOTARY_EVIDENCE_DIR="$TEST_WORK/evidence" GITHUB_STEP_SUMMARY="$TEST_WORK/summary" MOCK_PACKAGE_STATUS=0
cat > "$TEST_WORK/repo/scripts/package-macos.sh" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
printf 'package\n' >> "$FIXTURE_STATE/package-calls"
[[ $MOCK_PACKAGE_STATUS == 0 ]] || exit "$MOCK_PACKAGE_STATUS"
root=$(cd -- "$(dirname -- "$0")/.." && pwd)
packages="$root/dist/packages/darwin-arm64-release"
mkdir -p "$packages"
printf 'original candidate ZIP\n' > "$packages/takedown-1.2.3-macos-arm64-release.zip"
printf 'original candidate DMG\n' > "$packages/takedown-1.2.3-macos-arm64-release.dmg"
(cd "$packages" && shasum -a 256 ./*.zip ./*.dmg | sed 's@  \./@  @' > SHA256SUMS)
MOCK
cat > "$TEST_WORK/repo/scripts/verify-macos-package.sh" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
[[ $GITHUB_RUN_ID == 123 && $GITHUB_RUN_ATTEMPT == 1 ]]
root=$(cd -- "$(dirname -- "$0")/.." && pwd)
packages="$root/dist/packages/darwin-arm64-release"
(cd "$packages" && shasum -a 256 -c SHA256SUMS) >/dev/null
jq -n --arg source "$NOTARY_SOURCE_SHA" --arg version "$RELEASE_VERSION" --arg run "$GITHUB_RUN_ID" --arg attempt "$GITHUB_RUN_ATTEMPT" \
  --arg zip "$(shasum -a 256 "$packages/"*.zip | awk '{print $1}')" --arg dmg "$(shasum -a 256 "$packages/"*.dmg | awk '{print $1}')" \
  '{source_commit:$source,version:$version,architecture:"arm64",ci_run_id:$run,ci_run_attempt:$attempt,zip_sha256:$zip,dmg_sha256:$dmg}' > "$NOTARY_EVIDENCE_DIR/acceptance.json"
MOCK
cat > "$TEST_WORK/repo/scripts/apple-checkpoint.sh" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
[[ -f $1/final/SHA256SUMS && -f $1/acceptance.json ]]
printf 'checkpoint\n' >> "$FIXTURE_STATE/checkpoint-calls"
MOCK
cat > "$TEST_WORK/repo/scripts/apple-release-state.sh" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
[[ $1 == restore && ! -e $2 && $GITHUB_RUN_ID == 123 && $GITHUB_RUN_ATTEMPT == 1 ]]
mkdir "$2"
cp -R "$FIXTURE_STATE/retained/." "$2/"
MOCK
mkdir "$NOTARY_EVIDENCE_DIR"
mkdir -p "$NOTARY_EVIDENCE_DIR/app" "$NOTARY_EVIDENCE_DIR/dmg"
printf '{"team_id":"ABCDEFGHIJ","state":"verified"}\n' > "$NOTARY_EVIDENCE_DIR/app/receipt.json"
cp "$NOTARY_EVIDENCE_DIR/app/receipt.json" "$NOTARY_EVIDENCE_DIR/dmg/receipt.json"
continue_candidate() { bash "$TEST_WORK/repo/scripts/continue-apple-release.sh" > "$TEST_WORK/log" 2>&1; }
MOCK_PACKAGE_STATUS=75 continue_candidate
[[ ! -d $NOTARY_EVIDENCE_DIR/final && ! -f $FIXTURE_STATE/checkpoint-calls ]]
if MOCK_PACKAGE_STATUS=1 continue_candidate; then echo 'permanent package failure passed' >&2; exit 1; fi
[[ ! -d $NOTARY_EVIDENCE_DIR/final && ! -f $FIXTURE_STATE/checkpoint-calls ]]
continue_candidate
cp "$NOTARY_EVIDENCE_DIR/final/"*.zip "$TEST_WORK/original.zip"
cp "$NOTARY_EVIDENCE_DIR/final/"*.dmg "$TEST_WORK/original.dmg"
count=$(wc -l < "$FIXTURE_STATE/package-calls")
rm -rf "$TEST_WORK/repo/dist"
MOCK_PACKAGE_STATUS=1 continue_candidate
[[ $(wc -l < "$FIXTURE_STATE/package-calls") == "$count" ]]
cmp "$TEST_WORK/original.zip" "$NOTARY_EVIDENCE_DIR/final/"*.zip
cmp "$TEST_WORK/original.dmg" "$NOTARY_EVIDENCE_DIR/final/"*.dmg
mkdir -p "$NOTARY_EVIDENCE_DIR/app" "$NOTARY_EVIDENCE_DIR/dmg"
printf '{"team_id":"ABCDEFGHIJ","state":"verified"}\n' > "$NOTARY_EVIDENCE_DIR/app/receipt.json"
cp "$NOTARY_EVIDENCE_DIR/app/receipt.json" "$NOTARY_EVIDENCE_DIR/dmg/receipt.json"
cp -R "$NOTARY_EVIDENCE_DIR" "$FIXTURE_STATE/retained"
rm -rf "$TEST_WORK/repo/dist"
NOTARY_EVIDENCE_DIR="$TEST_WORK/fresh" GITHUB_RUN_ID=777 GITHUB_RUN_ATTEMPT=2 bash "$TEST_WORK/repo/scripts/fresh-apple-acceptance.sh" > "$TEST_WORK/log" 2>&1
jq -e '.ci_run_id=="123" and .ci_run_attempt=="1" and .verification_run_id=="777" and .verification_run_attempt=="2" and .credential_free_runner==true' "$TEST_WORK/fresh/fresh-acceptance-arm64.json" >/dev/null
cp -R "$FIXTURE_STATE/retained" "$TEST_WORK/full-retained"
for stage in app dmg; do
  for state in validation-pending ticket-blocked; do
    jq --arg state "$state" '.state=$state' "$TEST_WORK/full-retained/$stage/receipt.json" > "$FIXTURE_STATE/retained/$stage/receipt.json"
    if bash "$TEST_WORK/repo/scripts/apple-stages-verified.sh" "$FIXTURE_STATE/retained"; then echo 'stale finals made blocked stage ready' >&2; exit 1; fi
    NOTARY_EVIDENCE_DIR="$TEST_WORK/fresh-$stage-$state" GITHUB_RUN_ID=777 bash "$TEST_WORK/repo/scripts/fresh-apple-acceptance.sh" > "$TEST_WORK/log" 2>&1
    [[ ! -f $TEST_WORK/fresh-$stage-$state/fresh-acceptance-arm64.json ]]
    cp "$FIXTURE_STATE/retained/$stage/receipt.json" "$NOTARY_EVIDENCE_DIR/$stage/receipt.json"
    count=$(wc -l < "$FIXTURE_STATE/package-calls")
    if MOCK_PACKAGE_STATUS=1 continue_candidate; then echo 'stale finals bypassed signing recovery' >&2; exit 1; fi
    [[ $(wc -l < "$FIXTURE_STATE/package-calls") -eq $((count + 1)) ]]
    cp "$TEST_WORK/full-retained/$stage/receipt.json" "$FIXTURE_STATE/retained/$stage/receipt.json"
    cp "$TEST_WORK/full-retained/$stage/receipt.json" "$NOTARY_EVIDENCE_DIR/$stage/receipt.json"
  done
done
echo 'PASS: pending/blocked app or DMG cannot reuse stale finals, claim fresh acceptance or pass readiness'
rm "$FIXTURE_STATE/retained/acceptance.json"
NOTARY_EVIDENCE_DIR="$TEST_WORK/pending" GITHUB_RUN_ID=777 bash "$TEST_WORK/repo/scripts/fresh-apple-acceptance.sh" > "$TEST_WORK/log" 2>&1
[[ ! -f $TEST_WORK/pending/fresh-acceptance-arm64.json ]]
checkpoints=$(wc -l < "$FIXTURE_STATE/checkpoint-calls")
printf 'corrupt\n' >> "$NOTARY_EVIDENCE_DIR/final/"*.zip
if continue_candidate; then echo 'corrupt retained final candidate passed' >&2; exit 1; fi
[[ $(wc -l < "$FIXTURE_STATE/checkpoint-calls") == "$checkpoints" ]]
cat > "$TEST_WORK/releasectl" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
[[ $* == 'validate-version --version 1.2.3' ]]
printf 'validated\n' >> "$FIXTURE_STATE/validator-calls"
exit "${MOCK_VALIDATE_STATUS:-0}"
MOCK
chmod +x "$TEST_WORK/releasectl"
validate_prebuilt() (
  # shellcheck source=scripts/package-common.sh
  source "$ROOT/scripts/package-common.sh"
  RELEASECTL="$TEST_WORK/releasectl" VERSION=1.2.3 RELEASE_VERSION=1.2.3
  package_validate_release_version
)
validate_prebuilt
[[ $(wc -l < "$FIXTURE_STATE/validator-calls") -eq 1 ]]
if MOCK_VALIDATE_STATUS=17 validate_prebuilt > "$TEST_WORK/log" 2>&1; then echo 'failed prebuilt validator was bypassed' >&2; exit 1; fi
echo 'PASS: synthetic wrappers block pending/failure promotion, reuse exact finals, bind fresh receipts and consume a prebuilt version validator'
