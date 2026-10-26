#!/usr/bin/env bash
set -euo pipefail
# shellcheck source=scripts/package-common.sh
source "$(dirname -- "$0")/package-common.sh"

# Release mode signs and notarizes a binary from the credential-free build.
if [[ ${MODE:-local} == release ]]; then
  [[ -n ${BINARY:-} ]] || fail 'release signing requires a verified prebuilt BINARY'
  "$(dirname -- "$0")/sign-macos.sh" --check
fi
package_init darwin
for tool in ditto hdiutil plutil shasum xattr; do need "$tool"; done
APP="$WORK/output/takedown.app"
mkdir -p "$APP/Contents/MacOS" "$APP/Contents/Resources"
cp "$WORK/takedown" "$APP/Contents/MacOS/takedown"
cp "$ROOT/assets/icon.icns" "$APP/Contents/Resources/icon.icns"
cp "$ROOT/build/macos/Info.plist" "$APP/Contents/Info.plist"
plutil -replace CFBundleShortVersionString -string "$VERSION" "$APP/Contents/Info.plist"
plutil -replace CFBundleVersion -string "$VERSION" "$APP/Contents/Info.plist"
plutil -lint "$APP/Contents/Info.plist"
package_docs "$APP/Contents/Resources/docs"
# Drop removable checkout extended attributes before signing/archiving; normalize
# access so installation in a shared Applications directory remains readable.
xattr -cr "$APP"
chmod -R u+rwX,go+rX,go-w "$APP"
if [[ $MODE == release ]]; then "$ROOT/scripts/sign-macos.sh" "$APP"; fi
BASE="takedown-$RELEASE_VERSION-macos-$ARCH-$MODE"
ditto -c -k --sequesterRsrc --keepParent "$APP" "$WORK/output/$BASE.zip"
mkdir "$WORK/dmg"
ditto "$APP" "$WORK/dmg/takedown.app"
ln -s /Applications "$WORK/dmg/Applications"
hdiutil create -quiet -nospotlight -fs HFS+ -format UDZO -volname takedown \
  -srcfolder "$WORK/dmg" "$WORK/output/$BASE.dmg"
if [[ $MODE == release ]]; then "$ROOT/scripts/sign-macos.sh" "$WORK/output/$BASE.dmg"; fi
(
  cd "$WORK/output"
  shasum -a 256 "$BASE.zip" "$BASE.dmg" > SHA256SUMS
)
package_finish
