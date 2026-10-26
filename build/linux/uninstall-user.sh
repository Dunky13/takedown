#!/usr/bin/env bash
set -euo pipefail
ROOT=$(cd -P -- "$(dirname -- "$0")" && pwd -P)
# shellcheck source=build/linux/install-common.sh
source "$ROOT/install-common.sh"
init_prefix
check_file_path "$MANIFEST"
[[ -f $MANIFEST ]] || fail 'no ownership manifest; refusing to remove unowned files'
lock_prefix
trap 'rmdir -- "$LOCK"' EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
# Validate the COMPLETE manifest and every owned file before removing anything.
FILES=()
HASHES=()
{
  IFS= read -r header
  [[ $header == takedown-install-v1 ]] || fail 'invalid ownership manifest header'
  while IFS= read -r record; do
    [[ $record =~ ^([0-9a-f]{64})\ \ (.+)$ ]] || fail 'invalid ownership manifest record'
    digest=${BASH_REMATCH[1]}
    rel=${BASH_REMATCH[2]}
    allowed_file "$rel" || fail "unrecognized ownership path: $rel"
    for previous in ${FILES[@]+"${FILES[@]}"}; do [[ $previous != "$rel" ]] || fail 'duplicate ownership path'; done
    FILES+=("$rel")
    HASHES+=("$digest")
  done
} < "$MANIFEST"
[[ ${#FILES[@]} -ge 3 ]] || fail 'incomplete ownership manifest'
for required in bin/takedown share/applications/takedown.desktop share/icons/hicolor/scalable/apps/takedown.svg; do
  found=false
  for rel in "${FILES[@]}"; do [[ $rel != "$required" ]] || found=true; done
  [[ $found == true ]] || fail "ownership manifest lacks $required"
done
for ((i=0; i<${#FILES[@]}; i++)); do
  path="$PREFIX/${FILES[i]}"
  check_file_path "$path"
  [[ ! -e $path || $(hash_file "$path") == "${HASHES[i]}" ]] || fail "owned file changed; preserved installation and manifest: $path"
done
for rel in "${FILES[@]}"; do rm -f -- "$PREFIX/$rel"; done
rm -- "$MANIFEST"
# Unknown files in these directories belong to the user and are never recursive removals.
rmdir "$PREFIX/share/doc/takedown" "$PREFIX/share/takedown" 2>/dev/null || true
if command -v update-desktop-database >/dev/null; then update-desktop-database "$PREFIX/share/applications"; fi
echo 'Removed recorded, unchanged application files. Unknown files, archives, state and tokens were retained.'
