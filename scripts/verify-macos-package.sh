#!/usr/bin/env bash
# Final distribution acceptance, without signing credentials or global installs.
set -euo pipefail
fail() { printf 'verify-macos-package: %s\n' "$*" >&2; exit 1; }
[[ $(uname -s) == Darwin ]] || fail 'macOS is required'
: "${VERSION:?VERSION is required}"
: "${RELEASE_VERSION:?RELEASE_VERSION is required}"
: "${NOTARY_ARCH:?NOTARY_ARCH is required}"
: "${APPLE_TEAM_ID:?APPLE_TEAM_ID is required}"
[[ $VERSION =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]] || fail 'invalid numeric version'
[[ $RELEASE_VERSION =~ ^[0-9A-Za-z.+-]+$ && ${RELEASE_VERSION%%[-+]*} == "$VERSION" ]] || fail 'invalid release version'
[[ $APPLE_TEAM_ID =~ ^[A-Z0-9]{10}$ ]] || fail 'invalid Apple team'
case $NOTARY_ARCH in arm64) slice=arm64 ;; amd64) slice=x86_64 ;; *) fail 'invalid architecture' ;; esac
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
PACKAGES="$ROOT/dist/packages/darwin-$NOTARY_ARCH-release"
BASE="takedown-$RELEASE_VERSION-macos-$NOTARY_ARCH-release"
[[ -d $PACKAGES && ! -L $PACKAGES ]] || fail 'signed package directory missing'
WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-acceptance.XXXXXX")
MOUNTED=false
cleanup() {
  result=$?
  trap - EXIT
  if [[ $MOUNTED == true ]]; then hdiutil detach -quiet "$WORK/mount" || result=1; fi
  rm -rf -- "$WORK" || result=1
  exit "$result"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
for extension in zip dmg; do
  [[ -f $PACKAGES/$BASE.$extension && ! -L $PACKAGES/$BASE.$extension ]] || fail 'signed asset missing'
done
[[ -f $PACKAGES/SHA256SUMS && ! -L $PACKAGES/SHA256SUMS ]] || fail 'checksum manifest missing'
printf '%s\n' "$BASE.zip" "$BASE.dmg" | LC_ALL=C sort > "$WORK/expected"
awk '/^[0-9a-fA-F]+  [^\/]+$/ && length($1) == 64 { print $2 }' "$PACKAGES/SHA256SUMS" | LC_ALL=C sort > "$WORK/actual"
cmp -s "$WORK/expected" "$WORK/actual" || fail 'manifest must cover exactly the final ZIP and DMG'
(cd "$PACKAGES" && shasum -a 256 -c SHA256SUMS)
check_app() {
  local app=$1
  [[ -d $app && ! -L $app && -f $app/Contents/MacOS/takedown && ! -L $app/Contents/MacOS/takedown ]] || fail 'invalid copied application'
  [[ $(plutil -extract CFBundleIdentifier raw -o - "$app/Contents/Info.plist") == io.github.dunky13.takedown ]] || fail 'wrong bundle identifier'
  [[ $(plutil -extract CFBundleExecutable raw -o - "$app/Contents/Info.plist") == takedown ]] || fail 'wrong bundle executable'
  [[ $(plutil -extract CFBundleShortVersionString raw -o - "$app/Contents/Info.plist") == "$VERSION" ]] || fail 'wrong bundle version'
  [[ $(plutil -extract CFBundleVersion raw -o - "$app/Contents/Info.plist") == "$VERSION" ]] || fail 'wrong bundle build'
  [[ $(lipo -archs "$app/Contents/MacOS/takedown") == "$slice" ]] || fail 'wrong binary slices'
  [[ $("$app/Contents/MacOS/takedown" --version) == "takedown version $RELEASE_VERSION" ]] || fail 'wrong executable version'
  codesign --verify --deep --strict --verbose=2 "$app"
  codesign --display --verbose=4 "$app" 2> "$WORK/signature"
  grep -qx "TeamIdentifier=$APPLE_TEAM_ID" "$WORK/signature" || fail 'wrong signed team'
  grep -q '^Authority=Developer ID Application: ' "$WORK/signature" || fail 'Developer ID signature missing'
  xcrun stapler validate "$app"
  spctl --assess --type execute --verbose=2 "$app"
}
mkdir "$WORK/zip" "$WORK/mount" "$WORK/Applications"
ditto -x -k "$PACKAGES/$BASE.zip" "$WORK/zip"
check_app "$WORK/zip/takedown.app"
codesign --verify --strict --verbose=2 "$PACKAGES/$BASE.dmg"
codesign --display --verbose=4 "$PACKAGES/$BASE.dmg" 2> "$WORK/signature"
grep -qx "TeamIdentifier=$APPLE_TEAM_ID" "$WORK/signature" || fail 'wrong DMG team'
grep -q '^Authority=Developer ID Application: ' "$WORK/signature" || fail 'DMG Developer ID signature missing'
xcrun stapler validate "$PACKAGES/$BASE.dmg"
spctl --assess --type open --context context:primary-signature --verbose=2 "$PACKAGES/$BASE.dmg"
# Set before attach so cleanup also attempts to detach a partially mounted image.
MOUNTED=true
hdiutil attach -quiet -readonly -nobrowse -mountpoint "$WORK/mount" "$PACKAGES/$BASE.dmg"
ditto "$WORK/mount/takedown.app" "$WORK/Applications/takedown.app"
check_app "$WORK/Applications/takedown.app"
cmp -s "$WORK/zip/takedown.app/Contents/MacOS/takedown" "$WORK/Applications/takedown.app/Contents/MacOS/takedown" || fail 'ZIP and installed DMG executables differ'
if [[ -n ${NOTARY_EVIDENCE_DIR:-} ]]; then
  jq -n --arg version "$RELEASE_VERSION" --arg arch "$NOTARY_ARCH" \
    --arg zip "$(shasum -a 256 "$PACKAGES/$BASE.zip" | awk '{print $1}')" \
    --arg dmg "$(shasum -a 256 "$PACKAGES/$BASE.dmg" | awk '{print $1}')" \
    '{version:$version,architecture:$arch,zip_sha256:$zip,dmg_sha256:$dmg,checks:["checksums","Developer ID","stapling","Gatekeeper","copied app metadata and version","native slices","ZIP/DMG executable equality"]}' \
    > "$NOTARY_EVIDENCE_DIR/acceptance.json"
fi
echo 'PASS: final ZIP, mounted DMG and copied application signatures, tickets, Gatekeeper, metadata and native version'
