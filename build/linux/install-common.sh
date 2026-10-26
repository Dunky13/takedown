#!/usr/bin/env bash
# Shared by the extracted archive's install/uninstall commands.
set -euo pipefail
fail() { printf 'takedown install: %s\n' "$*" >&2; exit 1; }
init_prefix() {
  PREFIX=${INSTALL_PREFIX:-"$HOME/.local"}
  [[ $PREFIX == /* && $PREFIX != / && $PREFIX != */ && $PREFIX != *//* ]] || fail 'INSTALL_PREFIX must be absolute and non-root, without trailing/repeated slashes'
  [[ $PREFIX != *'"'* && $PREFIX != *'`'* && $PREFIX != *'$'* && $PREFIX != *'%'* && $PREFIX != *\\* && $PREFIX != *$'\n'* && $PREFIX != *$'\r'* && $PREFIX != *$'\t'* ]] || fail 'INSTALL_PREFIX contains unsupported desktop-entry characters'
  check_directories "$PREFIX"
  MANIFEST="$PREFIX/share/takedown/install-manifest.sha256"
  check_file_path "$MANIFEST"
  LOCK="$PREFIX/.takedown-install-lock"
  if command -v sha256sum >/dev/null; then HASH_TOOL=sha256sum
  elif command -v shasum >/dev/null; then HASH_TOOL=shasum
  else fail 'sha256sum (or shasum) is required'; fi
}
check_directories() {
  local path=$1 part current=''
  local -a parts
  IFS=/ read -r -a parts <<< "$path"
  for part in "${parts[@]}"; do
    [[ -n $part ]] || continue
    [[ $part != . && $part != .. ]] || fail 'dot path components are not supported'
    current="$current/$part"
    [[ ! -L $current ]] || fail "symlink in destination ancestry: $current"
    [[ ! -e $current || -d $current ]] || fail "non-directory ancestor: $current"
  done
}
check_file_path() {
  check_directories "$(dirname -- "$1")"
  [[ ! -L $1 ]] || fail "symlink payload path: $1"
  [[ ! -e $1 || -f $1 ]] || fail "non-file payload path: $1"
}
hash_file() {
  if [[ $HASH_TOOL == sha256sum ]]; then sha256sum < "$1"
  else shasum -a 256 < "$1"; fi | cut -d ' ' -f 1
}
allowed_file() {
  case $1 in
    bin/takedown|share/applications/takedown.desktop|share/icons/hicolor/scalable/apps/takedown.svg) return 0 ;;
    share/doc/takedown/README.md|share/doc/takedown/PRIVACY.md|share/doc/takedown/PRIVACY.html|share/doc/takedown/BUILD-INFO.txt|share/doc/takedown/LICENSE|share/doc/takedown/LICENSE.md|share/doc/takedown/LICENSE.txt|share/doc/takedown/COPYING) return 0 ;;
    share/doc/takedown/NOTICE|share/doc/takedown/THIRD_PARTY_NOTICES.md) return 0 ;;
    *) return 1 ;;
  esac
}
lock_prefix() {
  mkdir -p -- "$PREFIX"
  mkdir -- "$LOCK" || fail 'another operation or an unresolved interrupted operation owns the prefix lock'
  chmod 700 "$LOCK"
}
