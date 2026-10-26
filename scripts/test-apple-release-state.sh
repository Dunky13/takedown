#!/usr/bin/env bash
# Synthetic durable-custody failures. No GitHub or Apple calls.
set -euo pipefail
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
TEST_WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-state-tests.XXXXXX")
cleanup() {
  result=$?
  if [[ $result != 0 ]]; then
    [[ ! -f ${MOCK_STATE:-}/log ]] || cat "$MOCK_STATE/log" >&2
    [[ ! -f $TEST_WORK/log ]] || cat "$TEST_WORK/log" >&2
  fi
  rm -rf -- "$TEST_WORK"
}
trap cleanup EXIT
mkdir "$TEST_WORK/bin" "$TEST_WORK/remote" "$TEST_WORK/evidence"
export MOCK_REMOTE="$TEST_WORK/remote" GH_REPO=fixture/takedown NOTARY_STATE_TAG=apple-candidate-v1.2.3
export NOTARY_SOURCE_SHA=1111111111111111111111111111111111111111 RELEASE_VERSION=1.2.3 NOTARY_ARCH=arm64 GITHUB_RUN_ID=123 GITHUB_RUN_ATTEMPT=1
export MOCK_STATE="$TEST_WORK/state"
mkdir -p "$MOCK_STATE"
ln -s "$MOCK_REMOTE" "$MOCK_STATE/assets"
jq -n --arg tag "$NOTARY_STATE_TAG" --arg source "$NOTARY_SOURCE_SHA" '{id:7,draft:true,tag_name:$tag,target_commitish:$source,body:"sealed"}' > "$MOCK_STATE/release.json"
cp "$ROOT/scripts/testdata/apple-custody-gh.sh" "$TEST_WORK/bin/gh"
chmod +x "$TEST_WORK/bin/gh"
export PATH="$TEST_WORK/bin:$PATH"
jq -n --arg source "$NOTARY_SOURCE_SHA" --arg version "$RELEASE_VERSION" --arg repo "$GH_REPO" --arg run "$GITHUB_RUN_ID" --arg attempt "$GITHUB_RUN_ATTEMPT" '{schema:1,source:$source,version:$version,repository:$repo,run_id:$run,run_attempt:$attempt}' > "$MOCK_REMOTE/candidate.json"
printf '{"synthetic":true}\n' > "$TEST_WORK/evidence/input.json"
state() { bash "$ROOT/scripts/apple-release-state.sh" "$@" > "$TEST_WORK/log" 2>&1; }
reject() { if state "$@"; then echo 'expected custody failure' >&2; exit 1; fi; }
state save "$TEST_WORK/evidence"
state restore "$TEST_WORK/restored"
cmp "$TEST_WORK/evidence/input.json" "$TEST_WORK/restored/input.json"
printf '{"synthetic":true,"stage":"prepared"}\n' > "$TEST_WORK/evidence/input.json"
state save "$TEST_WORK/evidence"
state restore "$TEST_WORK/restored-new"
cmp "$TEST_WORK/evidence/input.json" "$TEST_WORK/restored-new/input.json"
[[ $(jq -r .sequence "$TEST_WORK/restored-new/checkpoint.json") == 2 ]]
GITHUB_RUN_ID=456 reject restore "$TEST_WORK/wrong-run"
NOTARY_ARCH=amd64 state restore "$TEST_WORK/other-arch"
[[ ! -f $TEST_WORK/other-arch/input.json ]]
# Durable bytes are immutable even when a draft is only partially populated.
printf 'collision\n' > "$TEST_WORK/collision"
reject put candidate.json "$TEST_WORK/collision"
# A changed sequence cannot be forked, including after a lost upload response.
cp "$MOCK_REMOTE/"checkpoint-arm64-000002-*.json "$MOCK_REMOTE/checkpoint-arm64-000002-aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.json"
reject restore "$TEST_WORK/forked"
rm "$MOCK_REMOTE/checkpoint-arm64-000002-aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.json"
blob=$(jq -r '.files[0].asset' "$TEST_WORK/evidence/checkpoint.json")
printf 'corrupt\n' >> "$MOCK_REMOTE/$blob"
reject restore "$TEST_WORK/corrupt"
printf 'PRIVATE-SENTINEL\n' > "$TEST_WORK/evidence/credential.p8"
reject save "$TEST_WORK/evidence"
[[ -z $(find "$MOCK_REMOTE" -name '*credential*' -print) ]]
printf 'same bytes\n' > "$TEST_WORK/upload"
MOCK_STALE_UPLOAD=true state put stale.txt "$TEST_WORK/upload"
for filter in '.id=0' '.name="wrong"' '.size=0' '.state="starter"'; do
  name="bad-$(printf '%s' "$filter" | cksum | awk '{print $1}').txt"
  MOCK_UPLOAD_FILTER="$filter" reject put "$name" "$TEST_WORK/upload"
done
MOCK_CHANGED_BYTES=true reject put changed.txt "$TEST_WORK/upload"
MOCK_LOST_RESPONSE=true reject put interrupted.txt "$TEST_WORK/upload"
state put interrupted.txt "$TEST_WORK/upload"
# Successful recovery reads existing bytes and never retries the POST.
[[ $(grep -c 'POST .*name=interrupted.txt' "$MOCK_STATE/calls") == 1 ]]
NOTARY_RELEASE_ID=8 reject inventory
NOTARY_RELEASE_BODY_SHA256=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa reject inventory
MOCK_REQUIRE_BINARY=true state get interrupted.txt "$TEST_WORK/recovered"
cmp "$TEST_WORK/upload" "$TEST_WORK/recovered"
echo 'PASS: offline custody restore/sequence/run/arch, collision/fork/corruption/allowlist, stale uploads, receipt identity, interrupted writes and binary downloads'
