#!/usr/bin/env bash
# Synthetic final-package acceptance; no signing, mounts, account access or build.
set -euo pipefail
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
TEST_WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-acceptance-tests.XXXXXX")
trap 'rm -rf -- "$TEST_WORK"' EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
mkdir -p "$TEST_WORK/tools" "$TEST_WORK/scratch" "$TEST_WORK/repo/scripts"
cp "$ROOT/scripts/verify-macos-package.sh" "$TEST_WORK/repo/scripts/"
cat > "$TEST_WORK/tools/uname" <<'MOCK'
#!/usr/bin/env bash
echo Darwin
MOCK
cat > "$TEST_WORK/tools/plutil" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
[[ $1 == -extract && $3 == raw && $4 == -o && $5 == - ]] || exit 99
jq -er --arg key "$2" '.[$key] // empty' "$6"
MOCK
cat > "$TEST_WORK/tools/lipo" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
[[ $1 == -archs ]] || exit 99
jq -r .fixtureArch "$(dirname -- "$2")/../Info.plist"
MOCK
cat > "$TEST_WORK/tools/codesign" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
if [[ $1 == --display ]]; then
  target=${!#}
  team=$APPLE_TEAM_ID
  if [[ $target == *.app ]]; then team=$(jq -r .fixtureTeam "$target/Contents/Info.plist"); fi
  printf 'TeamIdentifier=%s\nAuthority=Developer ID Application: Fixture\n' "$team" >&2
elif [[ $1 == --verify ]]; then
  [[ $MOCK_SCENARIO != codesign-failure || ${!#} != */Applications/* ]] || { echo 'fixture codesign failure' >&2; exit 1; }
else exit 99
fi
MOCK
cat > "$TEST_WORK/tools/ditto" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
if [[ $1 == -x && $2 == -k ]]; then
  cp -R "$MOCK_STATE/zip-app" "$4/takedown.app"
else
  [[ $# == 2 ]] || exit 99
  cp -R "$1" "$2"
fi
MOCK
cat > "$TEST_WORK/tools/hdiutil" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
case $1 in
  attach)
    printf 'attach\n' >> "$MOCK_STATE/mount-calls"
    mountpoint=''
    while (($#)); do
      if [[ $1 == -mountpoint ]]; then mountpoint=$2; break; fi
      shift
    done
    [[ -n $mountpoint ]]
    cp -R "$MOCK_STATE/dmg-app" "$mountpoint/takedown.app"
    [[ $MOCK_SCENARIO != attach-failure ]] || { echo 'fixture attach failure' >&2; exit 1; }
    ;;
  detach)
    printf 'detach\n' >> "$MOCK_STATE/mount-calls"
    [[ $MOCK_SCENARIO != detach-failure ]] || { echo 'fixture detach failure' >&2; exit 1; }
    ;;
  *) exit 99 ;;
esac
MOCK
cat > "$TEST_WORK/tools/xcrun" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
[[ $1 == stapler && $2 == validate ]] || exit 99
printf 'stapler %s\n' "$3" >> "$MOCK_STATE/trust-calls"
[[ $MOCK_SCENARIO != stapler-failure || $3 != */Applications/* ]] || { echo 'fixture stapler failure' >&2; exit 1; }
MOCK
cat > "$TEST_WORK/tools/spctl" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
[[ $1 == --assess ]] || exit 99
printf 'spctl %s\n' "${!#}" >> "$MOCK_STATE/trust-calls"
[[ $MOCK_SCENARIO != gatekeeper-failure || ${!#} != */Applications/* ]] || { echo 'fixture Gatekeeper failure' >&2; exit 1; }
MOCK
cat > "$TEST_WORK/tools/xattr" <<'MOCK'
#!/usr/bin/env bash
printf 'QUARANTINE REMOVAL FORBIDDEN\n' >> "$MOCK_STATE/forbidden"
exit 99
MOCK
chmod +x "$TEST_WORK/tools/"*
export PATH="$TEST_WORK/tools:$PATH" TMPDIR="$TEST_WORK/scratch"
export VERSION=0.0.1 RELEASE_VERSION=0.0.1 NOTARY_ARCH=arm64 APPLE_TEAM_ID=ABCDEFGHIJ
export NOTARY_SOURCE_SHA=1111111111111111111111111111111111111111 GITHUB_RUN_ID=123 GITHUB_RUN_ATTEMPT=1
export MOCK_STATE MOCK_SCENARIO NOTARY_EVIDENCE_DIR
PACKAGES="$TEST_WORK/repo/dist/packages/darwin-arm64-release"
BASE=takedown-0.0.1-macos-arm64-release
new_case() {
  MOCK_SCENARIO=$1
  MOCK_STATE="$TEST_WORK/$1"
  NOTARY_EVIDENCE_DIR="$MOCK_STATE/evidence"
  mkdir -p "$MOCK_STATE/zip-app/Contents/MacOS" "$NOTARY_EVIDENCE_DIR"
  printf '#!/bin/sh\nprintf "takedown version 0.0.1\\n"\n' > "$MOCK_STATE/zip-app/Contents/MacOS/takedown"
  chmod +x "$MOCK_STATE/zip-app/Contents/MacOS/takedown"
  jq -n '{CFBundleIdentifier:"io.github.dunky13.takedown",CFBundleExecutable:"takedown",CFBundleShortVersionString:"0.0.1",CFBundleVersion:"0.0.1",fixtureArch:"arm64",fixtureTeam:"ABCDEFGHIJ"}' > "$MOCK_STATE/zip-app/Contents/Info.plist"
  cp -R "$MOCK_STATE/zip-app" "$MOCK_STATE/dmg-app"
  rm -rf "$PACKAGES"
  mkdir -p "$PACKAGES"
  printf 'synthetic final ZIP\n' > "$PACKAGES/$BASE.zip"
  printf 'synthetic final DMG\n' > "$PACKAGES/$BASE.dmg"
  (cd "$PACKAGES" && shasum -a 256 "$BASE.zip" "$BASE.dmg" > SHA256SUMS)
  case $1 in
    wrong-copied-version) jq '.CFBundleShortVersionString="0.0.2"' "$MOCK_STATE/dmg-app/Contents/Info.plist" > "$MOCK_STATE/changed"; mv "$MOCK_STATE/changed" "$MOCK_STATE/dmg-app/Contents/Info.plist" ;;
    wrong-copied-team) jq '.fixtureTeam="ZZZZZZZZZZ"' "$MOCK_STATE/dmg-app/Contents/Info.plist" > "$MOCK_STATE/changed"; mv "$MOCK_STATE/changed" "$MOCK_STATE/dmg-app/Contents/Info.plist" ;;
    wrong-copied-architecture) jq '.fixtureArch="x86_64"' "$MOCK_STATE/dmg-app/Contents/Info.plist" > "$MOCK_STATE/changed"; mv "$MOCK_STATE/changed" "$MOCK_STATE/dmg-app/Contents/Info.plist" ;;
    wrong-copied-executable-version) printf '#!/bin/sh\necho "takedown version 0.0.2"\n' > "$MOCK_STATE/dmg-app/Contents/MacOS/takedown" ;;
    executable-mismatch) printf '# different signed candidate\n' >> "$MOCK_STATE/dmg-app/Contents/MacOS/takedown" ;;
    checksum-corruption) printf 'corrupted\n' >> "$PACKAGES/$BASE.zip" ;;
    manifest-extra) printf '%064d  unexpected.zip\n' 0 >> "$PACKAGES/SHA256SUMS" ;;
  esac
}
verify() { bash "$TEST_WORK/repo/scripts/verify-macos-package.sh" > "$MOCK_STATE/output" 2>&1; }
assert_cleanup() {
  [[ ! -f $MOCK_STATE/forbidden ]]
  [[ -z $(find "$TMPDIR" -mindepth 1 -maxdepth 1 -type d -name 'takedown-acceptance.*' -print) ]]
  if [[ -f $MOCK_STATE/mount-calls ]]; then
    [[ $(cat "$MOCK_STATE/mount-calls") == $'attach\ndetach' ]]
  fi
}
new_case success
verify || { cat "$MOCK_STATE/output"; exit 1; }
assert_cleanup
jq -e '.version=="0.0.1" and .architecture=="arm64" and (.zip_sha256|test("^[0-9a-f]{64}$")) and (.dmg_sha256|test("^[0-9a-f]{64}$")) and (.checks|index("ZIP/DMG executable equality")!=null)' "$NOTARY_EVIDENCE_DIR/acceptance.json" >/dev/null
grep -q 'stapler .*/Applications/takedown.app' "$MOCK_STATE/trust-calls"
grep -q 'spctl .*/Applications/takedown.app' "$MOCK_STATE/trust-calls"
echo 'PASS: checksums, ZIP and copied DMG app acceptance, evidence, mount cleanup, quarantine preserved'
for scenario in wrong-copied-version wrong-copied-team wrong-copied-architecture wrong-copied-executable-version checksum-corruption manifest-extra executable-mismatch gatekeeper-failure stapler-failure codesign-failure attach-failure detach-failure; do
  new_case "$scenario"
  if verify; then cat "$MOCK_STATE/output"; echo "expected rejection: $scenario" >&2; exit 1; fi
  case $scenario in
    wrong-copied-version) expected='wrong bundle version' ;;
    wrong-copied-team) expected='wrong signed team' ;;
    wrong-copied-architecture) expected='wrong binary slices' ;;
    wrong-copied-executable-version) expected='wrong executable version' ;;
    checksum-corruption) expected='FAILED' ;;
    manifest-extra) expected='manifest must cover exactly' ;;
    executable-mismatch) expected='ZIP and installed DMG executables differ' ;;
    gatekeeper-failure) expected='fixture Gatekeeper failure' ;;
    stapler-failure) expected='fixture stapler failure' ;;
    codesign-failure) expected='fixture codesign failure' ;;
    attach-failure) expected='fixture attach failure' ;;
    detach-failure) expected='fixture detach failure' ;;
  esac
  grep -Fq "$expected" "$MOCK_STATE/output" || { cat "$MOCK_STATE/output"; echo "wrong rejection for $scenario" >&2; exit 1; }
  assert_cleanup
  if [[ $scenario == checksum-corruption || $scenario == manifest-extra ]]; then
    [[ ! -f $MOCK_STATE/mount-calls ]]
  else
    [[ -f $MOCK_STATE/mount-calls ]]
  fi
  [[ $scenario == detach-failure || ! -e $NOTARY_EVIDENCE_DIR/acceptance.json ]]
  echo "PASS: $scenario rejected, owned mount cleanup attempted, quarantine preserved"
done
echo 'PASS: synthetic final-package acceptance contracts; no Apple, real mounts or signing'
