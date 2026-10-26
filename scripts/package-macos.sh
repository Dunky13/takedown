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
for tool in ditto hdiutil plutil shasum; do need "$tool"; done
APP="$WORK/output/takedown.app"
mkdir -p "$APP/Contents/MacOS" "$APP/Contents/Resources"
cp "$WORK/takedown" "$APP/Contents/MacOS/takedown"
cp "$ROOT/assets/icon.icns" "$APP/Contents/Resources/icon.icns"
cp "$ROOT/build/macos/Info.plist" "$APP/Contents/Info.plist"
plutil -replace CFBundleShortVersionString -string "$VERSION" "$APP/Contents/Info.plist"
plutil -replace CFBundleVersion -string "$VERSION" "$APP/Contents/Info.plist"
plutil -lint "$APP/Contents/Info.plist"
package_docs "$APP/Contents/Resources/docs"
# Preserve extended attributes, including quarantine. Normalize access so
# installation in a shared Applications directory remains readable.
chmod -R u+rwX,go+rX,go-w "$APP"
if [[ $MODE == release ]]; then "$ROOT/scripts/sign-macos.sh" "$APP"; fi
BASE="takedown-$RELEASE_VERSION-macos-$ARCH-$MODE"
if [[ $MODE == release ]]; then
  cp "$NOTARY_EVIDENCE_DIR/app/stapled-app.zip" "$WORK/output/$BASE.zip"
else
  ditto -c -k --sequesterRsrc --keepParent "$APP" "$WORK/output/$BASE.zip"
fi
mkdir "$WORK/dmg"
ditto "$APP" "$WORK/dmg/takedown.app"
ln -s /Applications "$WORK/dmg/Applications"
# A retained DMG must be recovered from its bound receipt by the signer. Never
# rebuild a signed outer stage on continuation (even if its UUID is still pending).
if [[ $MODE == release && -f ${NOTARY_EVIDENCE_DIR:-}/dmg/receipt.json ]] && \
  [[ $(jq -r .state "$NOTARY_EVIDENCE_DIR/dmg/receipt.json") != prepared ]]; then
  : # sign-macos restores and verifies the exact retained disk image below.
else
  hdiutil create -quiet -nospotlight -fs HFS+ -format UDZO -volname takedown \
    -srcfolder "$WORK/dmg" "$WORK/output/$BASE.dmg"
fi
if [[ $MODE == release ]]; then "$ROOT/scripts/sign-macos.sh" "$WORK/output/$BASE.dmg"; fi
(
  cd "$WORK/output"
  shasum -a 256 "$BASE.zip" "$BASE.dmg" > SHA256SUMS
)
package_finish
