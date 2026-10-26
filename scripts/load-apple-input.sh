#!/usr/bin/env bash
set -euo pipefail
umask 077
: "${NOTARY_EVIDENCE_DIR:?}" "${NOTARY_STAGING_DIR:?}" "${NOTARY_ARCH:?}" "${NOTARY_SOURCE_SHA:?}"
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
STATE="$ROOT/scripts/apple-release-state.sh"
scratch=$(mktemp -d "${TMPDIR:-/tmp}/takedown-input.XXXXXX")
trap 'rm -rf -- "$scratch"' EXIT
bash "$STATE" restore "$NOTARY_EVIDENCE_DIR"
bash "$STATE" get candidate.json "$scratch/candidate.json"
bash "$STATE" get custody-complete.json "$scratch/custody-complete.json"
cmp -s "$scratch/candidate.json" "$scratch/custody-complete.json" || { echo 'candidate custody is incomplete' >&2; exit 1; }
bash "$STATE" get native.json "$scratch/native.json"
[[ $(shasum -a 256 "$scratch/native.json" | awk '{print $1}') == "$(jq -r .native_sha256 "$scratch/candidate.json")" ]] || { echo 'native manifest hash mismatch' >&2; exit 1; }
case $NOTARY_ARCH in arm64) runner_arch=ARM64 ;; amd64) runner_arch=X64 ;; *) exit 1 ;; esac
asset="unsigned-native-macOS-$runner_arch-$NOTARY_SOURCE_SHA.zip"
bash "$STATE" get "$asset" "$scratch/native.zip"
export NOTARY_NATIVE_METADATA="$scratch/native.json" NOTARY_NATIVE_ARCHIVE="$scratch/native.zip"
bash "$ROOT/scripts/prepare-macos-signing.sh"
if [[ -f $NOTARY_EVIDENCE_DIR/input.json ]]; then
  cmp -s "$NOTARY_STAGING_DIR/input.json" "$NOTARY_EVIDENCE_DIR/input.json" || { echo 'restored input receipt differs from original build' >&2; exit 1; }
else
  cp "$NOTARY_STAGING_DIR/input.json" "$NOTARY_EVIDENCE_DIR/input.json"
fi
