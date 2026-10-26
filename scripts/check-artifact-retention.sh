#!/usr/bin/env bash
# Check literal upload-artifact globs against every final checksum entry.
# Does not upload, publish, inspect account state, or need a YAML dependency.
set -euo pipefail
WORKFLOW=${1:?usage: check-artifact-retention.sh workflow.yml [package-root]}
ARTIFACT_ROOT=${2:-dist/packages}
ARTIFACT_ROOT=${ARTIFACT_ROOT%/}
[[ -d $ARTIFACT_ROOT ]] || { echo 'artifact root is missing' >&2; exit 1; }
PATTERNS=()
while IFS= read -r pattern; do PATTERNS+=("$pattern"); done < <(
  awk '
    /uses: actions\/upload-artifact@/ { upload=1 }
    upload && /^[[:space:]]+path: \|/ { paths=1; next }
    paths && /^[[:space:]]+dist\/packages\// { sub(/^[[:space:]]+/, ""); print; next }
    paths { exit }
  ' "$WORKFLOW"
)
[[ ${#PATTERNS[@]} -gt 0 ]] || { echo 'no supported literal artifact upload paths' >&2; exit 1; }
manifests=0
while IFS= read -r -d '' manifest; do
  manifests=$((manifests + 1))
  entries=0
  directory=${manifest%/*}
  while IFS= read -r record; do
    record=${record%$'\r'}
    [[ $record =~ ^[0-9a-fA-F]{64}\ \ ([^/]+)$ ]] || { echo 'invalid checksum record' >&2; exit 1; }
    name=${BASH_REMATCH[1]}
    [[ $name != . && $name != .. && $name != *\\* ]] || { echo 'invalid checksum filename' >&2; exit 1; }
    file="$directory/$name"
    [[ -f $file && ! -L $file ]] || { printf 'checksum artifact is missing: %s\n' "$name" >&2; exit 1; }
    logical="dist/packages/${file#"$ARTIFACT_ROOT/"}"
    selected=false
    for pattern in "${PATTERNS[@]}"; do
      # Dynamic RHS is intentionally a glob, like upload-artifact's path list.
      # shellcheck disable=SC2254
      case $logical in $pattern) selected=true ;; esac
    done
    [[ $selected == true ]] || { printf 'checksum artifact is not retained: %s\n' "$logical" >&2; exit 1; }
    entries=$((entries + 1))
  done < "$manifest"
  [[ $entries -gt 0 ]] || { echo 'empty checksum manifest' >&2; exit 1; }
done < <(find "$ARTIFACT_ROOT" -type f -name SHA256SUMS -print0)
[[ $manifests -gt 0 ]] || { echo 'no checksum manifests found' >&2; exit 1; }
printf 'PASS: every entry in %s checksum manifest(s) is retained\n' "$manifests"
