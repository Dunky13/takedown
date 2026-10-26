#!/usr/bin/env bash
# Offline publication boundary tests; gh is a stateful file-only fixture.
set -euo pipefail
cd -- "$(dirname -- "$0")/.."
WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-publish-tests.XXXXXX")
cleanup() {
  local result=$?
  if [[ $result != 0 && -f $WORK/log ]]; then cat "$WORK/log" >&2; fi
  rm -rf -- "$WORK"
}
trap cleanup EXIT
export GH_REPO=fixture/takedown EXPECTED_SHA MOCK_STATE
EXPECTED_SHA=$(git rev-parse HEAD)
MOCK_STATE="$WORK/state"
mkdir "$WORK/mock-bin" "$WORK/plan"
cat > "$WORK/mock-bin/gh" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
printf '%s\n' "$*" >> "$MOCK_STATE/calls"
if [[ $1 == api ]]; then
  if [[ $2 == --method ]]; then
    printf '{"object":{"type":"commit","sha":"%s"}}\n' "$EXPECTED_SHA" > "$MOCK_STATE/ref"
    cat "$MOCK_STATE/ref"
  elif [[ $2 == --paginate ]]; then
    if [[ -f $MOCK_STATE/release ]]; then jq -s '.' "$MOCK_STATE/release"; else echo '[]'; fi
  elif [[ $2 == */git/ref/tags/* && -f $MOCK_STATE/ref ]]; then
    cat "$MOCK_STATE/ref"
  else
    printf 'gh: Not Found (HTTP 404)\n' >&2
    exit 1
  fi
elif [[ $1 == release && $2 == create ]]; then
  printf '{"draft":true,"tag_name":"%s"}\n' "$3" > "$MOCK_STATE/release"
  [[ ${MOCK_FAIL_UPLOAD:-false} != true ]] || exit 1
  shift 3
  for item in "$@"; do [[ $item != */assets/* ]] || basename "$item"; done > "$MOCK_STATE/assets"
elif [[ $1 == release && $2 == upload ]]; then
  shift 3
  for item in "$@"; do [[ $item != */assets/* ]] || basename "$item"; done > "$MOCK_STATE/assets"
elif [[ $1 == release && $2 == view ]]; then
  cat "$MOCK_STATE/assets"
  [[ ${MOCK_EXTRA_ASSET:-false} != true ]] || echo unexpected.exe
elif [[ $1 == release && $2 == edit ]]; then
  printf '{"draft":false,"tag_name":"%s"}\n' "$3" > "$MOCK_STATE/release"
else
  echo 'unexpected gh fixture command' >&2
  exit 1
fi
MOCK
chmod +x "$WORK/mock-bin/gh"
export PATH="$WORK/mock-bin:$PATH"
reset_fixture() {
  rm -rf "$MOCK_STATE" "$WORK/artifacts"
  mkdir "$MOCK_STATE" "$WORK/artifacts"
  RELEASE_CHANNEL=$1
  export RELEASE_CHANNEL
  version=${2:-1.2.0-rc.1}
  mode=release
  if [[ $RELEASE_CHANNEL == nightly ]]; then
    version="1.2.1-nightly.20261026022300.${EXPECTED_SHA:0:12}"
    mode=local
  fi
  jq -n --arg version "$version" --arg sha "$EXPECTED_SHA" \
    '{version:$version,package_version:($version|split("+")[0]|split("-")[0]),tag:("v"+$version),sha:$sha,prerelease:($version|split("+")[0]|contains("-")),skip:false}' > "$WORK/plan/plan.json"
  printf '# Release fixture\n' > "$WORK/plan/notes.md"
  names=("linux-amd64-local.tar.gz" "windows-amd64-local.zip" "windows-amd64-local-setup.exe"
    "macos-arm64-$mode.zip" "macos-arm64-$mode.dmg" "macos-amd64-$mode.zip" "macos-amd64-$mode.dmg")
  for name in "${names[@]}"; do
    directory="$WORK/artifacts/$name"
    mkdir "$directory"
    printf 'payload for %s\n' "$name" > "$directory/takedown-$version-$name"
    (cd "$directory" && shasum -a 256 "takedown-$version-$name" > SHA256SUMS)
  done
}
# macOS provides shasum rather than sha256sum; emulate only the required interface.
if ! command -v sha256sum >/dev/null; then
  cat > "$WORK/mock-bin/sha256sum" <<'SUM'
#!/usr/bin/env bash
[[ ${1:-} != -- ]] || shift
exec shasum -a 256 "$@"
SUM
  chmod +x "$WORK/mock-bin/sha256sum"
fi
publish() { bash scripts/publish-release.sh "$WORK/plan" "$WORK/artifacts" > "$WORK/log" 2>&1; }
reject() { if publish; then cat "$WORK/log"; echo 'expected publication rejection' >&2; exit 1; fi; }
no_writes() { if grep -Eq '^api --method|^release (create|upload|edit)' "$MOCK_STATE/calls" 2>/dev/null; then echo 'unexpected remote write' >&2; exit 1; fi; }

reset_fixture manual
publish
grep -Fq -- '--latest=false' "$MOCK_STATE/calls"
[[ $(jq -r .draft "$MOCK_STATE/release") == false ]]
# Published versions are immutable, even on a retry.
cp "$MOCK_STATE/calls" "$WORK/previous-calls"
reject
grep -q 'published releases are immutable' "$WORK/log"

reset_fixture nightly
publish
grep -q 'unsigned development' "$WORK/plan/notes.md" && exit 1 # source notes remain untouched

reset_fixture manual '1.2.0+build-meta'
publish
grep -q -- '--latest=true' "$MOCK_STATE/calls"

reset_fixture manual
file="$WORK/artifacts/${names[0]}/takedown-$version-${names[0]}"
printf corrupted >> "$file"
reject
no_writes

reset_fixture manual
rm -rf "$WORK/artifacts/${names[0]}"
reject
no_writes

reset_fixture manual
cp -R "$WORK/artifacts/${names[0]}" "$WORK/artifacts/duplicate"
reject
no_writes

reset_fixture manual
printf '{"object":{"type":"commit","sha":"0000000000000000000000000000000000000000"}}\n' > "$MOCK_STATE/ref"
reject
no_writes

reset_fixture manual
export MOCK_FAIL_UPLOAD=true
reject
[[ $(jq -r .draft "$MOCK_STATE/release") == true ]]
if grep -q '^release edit' "$MOCK_STATE/calls"; then exit 1; fi
unset MOCK_FAIL_UPLOAD
publish
grep -q '^release upload' "$MOCK_STATE/calls"
[[ $(jq -r .draft "$MOCK_STATE/release") == false ]]

reset_fixture manual
export MOCK_EXTRA_ASSET=true
reject
if grep -q '^release edit' "$MOCK_STATE/calls"; then exit 1; fi
unset MOCK_EXTRA_ASSET

# The marker is the newest published nightly, across API pages. A newer failed
# draft, stable release, or custom prerelease must not suppress the next nightly.
cat > "$WORK/releases.json" <<'JSON'
[{"tag_name":"v1.2.1-nightly.20261025022300.0123456789ab","draft":false,"prerelease":true,"published_at":"2026-10-25T03:00:00Z"}]
[{"tag_name":"v1.2.1-nightly.20261026022300.0123456789ab","draft":false,"prerelease":true,"published_at":"2026-10-26T03:00:00Z"},
 {"tag_name":"v1.2.1-nightly.20261027022300.0123456789ab","draft":true,"prerelease":true,"published_at":null},
 {"tag_name":"v1.3.0-rc.1","draft":false,"prerelease":true,"published_at":"2026-10-28T03:00:00Z"},
 {"tag_name":"v1.3.0","draft":false,"prerelease":false,"published_at":"2026-10-29T03:00:00Z"}]
JSON
[[ $(bash scripts/select-last-nightly.sh "$WORK/releases.json") == v1.2.1-nightly.20261026022300.0123456789ab ]]
printf '[]\n' > "$WORK/releases.json"
[[ -z $(bash scripts/select-last-nightly.sh "$WORK/releases.json") ]]
printf 'invalid JSON\n' > "$WORK/releases.json"
if bash scripts/select-last-nightly.sh "$WORK/releases.json" >/dev/null 2>&1; then exit 1; fi
printf 'PASS: complete publication, integrity, collisions, draft recovery, immutable releases, and successful-nightly marker selection\n'
