#!/usr/bin/env bash
# Sourced by sign-macos.sh; functions do no work until explicitly called.
keychain_capture_search_list() {
  local line pattern='^[[:space:]]*"([^"]+)"[[:space:]]*$'
  ORIGINAL_KEYCHAINS=()
  security list-keychains -d user > "$WORK/keychains-before.txt" || return 1
  while IFS= read -r line; do
    [[ -n $line ]] || continue
    [[ $line =~ $pattern ]] || { printf 'sign-macos: cannot safely parse user keychain search list\n' >&2; return 1; }
    [[ ${BASH_REMATCH[1]} == /* ]] || { printf 'sign-macos: non-absolute keychain path\n' >&2; return 1; }
    ORIGINAL_KEYCHAINS+=("${BASH_REMATCH[1]}")
  done < "$WORK/keychains-before.txt"
}

keychain_restore_search_list() {
  local result=0
  if [[ -e $KEYCHAIN ]]; then security delete-keychain "$KEYCHAIN" >/dev/null 2>&1 || result=1; fi
  # create-keychain transiently adds to the USER search list. On an isolated
  # signing runner restore its exact captured list/order, then compare readback.
  # No default-keychain assignment or system-domain search-list write occurs.
  security list-keychains -d user -s ${ORIGINAL_KEYCHAINS[@]+"${ORIGINAL_KEYCHAINS[@]}"} >/dev/null 2>&1 || result=1
  security list-keychains -d user > "$WORK/keychains-after.txt" || result=1
  if ! cmp -s "$WORK/keychains-before.txt" "$WORK/keychains-after.txt"; then
    printf 'sign-macos: user keychain search-list restoration did not verify\n' >&2
    result=1
  fi
  [[ ! -e $KEYCHAIN ]] || result=1
  return "$result"
}
