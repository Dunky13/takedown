#!/usr/bin/env bash
# Synthetic durable-custody failures. No GitHub or Apple calls.
set -euo pipefail
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
TEST_WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-state-tests.XXXXXX")
trap 'rm -rf -- "$TEST_WORK"' EXIT
mkdir "$TEST_WORK/bin" "$TEST_WORK/remote" "$TEST_WORK/evidence"
export MOCK_REMOTE="$TEST_WORK/remote" GH_REPO=fixture/takedown NOTARY_STATE_TAG=apple-candidate-v1.2.3
export NOTARY_SOURCE_SHA=1111111111111111111111111111111111111111 RELEASE_VERSION=1.2.3 NOTARY_ARCH=arm64 GITHUB_RUN_ID=123 GITHUB_RUN_ATTEMPT=1
cat > "$TEST_WORK/bin/gh" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
if [[ $1 == release && $2 == upload ]]; then
  [[ $* != *--clobber* ]]
  file=$4; name=${file##*/}
  [[ ! -e $MOCK_REMOTE/$name ]] || exit 9
  cp "$file" "$MOCK_REMOTE/$name"
elif [[ $1 == api ]]; then
  endpoint=${!#}
  case $endpoint in
    *'/releases?per_page=100') printf '[{"id":7,"draft":true,"tag_name":"%s"}]\n' "$NOTARY_STATE_TAG" ;;
    */releases/7/assets\?per_page=100)
      find "$MOCK_REMOTE" -type f | LC_ALL=C sort | while IFS= read -r file; do name=${file##*/}; jq -n --arg name "$name" '{id:$name,name:$name}'; done | jq -s . ;;
    */releases/assets/*) cat "$MOCK_REMOTE/${endpoint##*/}" ;;
    *) echo "unexpected API: $endpoint" >&2; exit 99 ;;
  esac
else exit 99
fi
MOCK
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
echo 'PASS: 10 synthetic custody cases (restore, sequence, run/arch binding, collision, fork, corruption, allowlist)'
