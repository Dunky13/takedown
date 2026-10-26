#!/usr/bin/env bash
set -euo pipefail
ROOT=$(cd -P -- "$(dirname -- "$0")" && pwd -P)
# shellcheck source=build/linux/install-common.sh
source "$ROOT/install-common.sh"
init_prefix
# POSIX link takes exactly two paths and refuses an existing destination rather
# than treating it as a target directory. No tool is installed by this script.
LINK_TOOL=$(command -v link) || fail 'the POSIX link utility is required for safe installation'
FILES=(bin/takedown share/applications/takedown.desktop share/icons/hicolor/scalable/apps/takedown.svg
  share/doc/takedown/NOTICE share/doc/takedown/THIRD_PARTY_NOTICES.md)
for name in README.md PRIVACY.md PRIVACY.html BUILD-INFO.txt LICENSE LICENSE.md LICENSE.txt COPYING; do
  src="$ROOT/share/doc/takedown/$name"
  [[ ! -L $src ]] || fail "symlink in package payload: $name"
  [[ ! -e $src ]] || FILES+=("share/doc/takedown/$name")
done
check_conflicts() {
  local rel path
  for path in "$PREFIX/share/doc/takedown" "$PREFIX/share/takedown"; do
    check_directories "$(dirname -- "$path")"
    [[ ! -e $path && ! -L $path ]] || fail "unowned or existing installation path: $path (uninstall an owned installation first)"
  done
  for rel in "${FILES[@]}"; do
    path="$PREFIX/$rel"
    check_file_path "$path"
    [[ ! -e $path ]] || fail "conflicting installation file: $path"
    check_file_path "$ROOT/$rel"
    [[ -f $ROOT/$rel ]] || fail "missing package payload: $rel"
  done
}
# Check every destination BEFORE writing, and again under the exclusive lock.
check_conflicts
lock_prefix
CREATED=()
DONE=false
DOCS_DIR_CREATED=false
MANIFEST_DIR_CREATED=false
cleanup() {
  result=$?
  trap - EXIT
  if [[ $DONE != true ]]; then
    for rel in ${CREATED[@]+"${CREATED[@]}"}; do
      path="$PREFIX/$rel"
      # Roll back only our unchanged files. Leave unexpected replacements alone.
      if [[ -f $path && ! -L $path && $(hash_file "$path") == "$(hash_file "$LOCK/stage/$rel")" ]]; then rm -f -- "$path" || result=1; fi
    done
    if [[ $DOCS_DIR_CREATED == true ]]; then rmdir "$PREFIX/share/doc/takedown" 2>/dev/null || true; fi
    if [[ $MANIFEST_DIR_CREATED == true ]]; then rmdir "$PREFIX/share/takedown" 2>/dev/null || true; fi
  fi
  rm -rf -- "$LOCK" || result=1
  exit "$result"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
check_conflicts
mkdir -p "$LOCK/stage"
for rel in "${FILES[@]}"; do
  mkdir -p "$LOCK/stage/$(dirname -- "$rel")"
  if [[ $rel == share/applications/takedown.desktop ]]; then
    while IFS= read -r line; do
      case $line in
        Exec=*) printf 'Exec="%s/bin/takedown" gui\n' "$PREFIX" ;;
        TryExec=*) printf 'TryExec=%s/bin/takedown\n' "$PREFIX" ;;
        *) printf '%s\n' "$line" ;;
      esac
    done < "$ROOT/$rel" > "$LOCK/stage/$rel"
  else cp -- "$ROOT/$rel" "$LOCK/stage/$rel"; fi
  chmod 644 "$LOCK/stage/$rel"
done
chmod 755 "$LOCK/stage/bin/takedown"
printf 'takedown-install-v1\n' > "$LOCK/manifest"
for rel in "${FILES[@]}"; do
  printf '%s  %s\n' "$(hash_file "$LOCK/stage/$rel")" "$rel" >> "$LOCK/manifest"
done
mkdir -p "$PREFIX/share/doc"
mkdir "$PREFIX/share/doc/takedown"
DOCS_DIR_CREATED=true
mkdir "$PREFIX/share/takedown"
MANIFEST_DIR_CREATED=true
for rel in "${FILES[@]}"; do
  mkdir -p "$PREFIX/$(dirname -- "$rel")"
  "$LINK_TOOL" "$LOCK/stage/$rel" "$PREFIX/$rel"
  CREATED+=("$rel")
done
# Commit ownership only after every payload link succeeds; never adopt old files.
"$LINK_TOOL" "$LOCK/manifest" "$MANIFEST"
DONE=true
if command -v update-desktop-database >/dev/null; then update-desktop-database "$PREFIX/share/applications"; fi
printf 'Installed to %s. Keep the extracted uninstall-user.sh and install-common.sh. Add bin to PATH for CLI use.\n' "$PREFIX"
