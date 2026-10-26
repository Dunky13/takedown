#!/usr/bin/env bash
set -euo pipefail
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
export NOTARY_CHECKPOINT_HOOK="$ROOT/scripts/apple-checkpoint.sh"
status=0
packages="$ROOT/dist/packages/darwin-$NOTARY_ARCH-release"
if [[ -f $NOTARY_EVIDENCE_DIR/acceptance.json && -f $NOTARY_EVIDENCE_DIR/final/SHA256SUMS ]] &&
  bash "$ROOT/scripts/apple-stages-verified.sh" "$NOTARY_EVIDENCE_DIR"; then
  # Successful architectures retain their final ZIP bytes across continuations,
  # including a publication retry. Repacking must never mutate a reserved asset.
  mkdir -p "$packages"
  cp "$NOTARY_EVIDENCE_DIR/final/"* "$packages/"
else
  bash "$ROOT/scripts/package-macos.sh" || status=$?
fi
if [[ $status == 75 ]]; then
  printf 'Apple candidate pending; the scheduled continuation retains the original UUID.\n' >> "${GITHUB_STEP_SUMMARY:?}"
  exit 0
fi
[[ $status == 0 ]] || exit "$status"
bash "$ROOT/scripts/verify-macos-package.sh"
mkdir -p "$NOTARY_EVIDENCE_DIR/final"
for file in "$packages/"*.zip "$packages/"*.dmg "$packages/SHA256SUMS"; do
  retained="$NOTARY_EVIDENCE_DIR/final/${file##*/}"
  if [[ -e $retained ]]; then cmp -s "$file" "$retained" || { echo 'final candidate bytes changed' >&2; exit 1; }; else cp "$file" "$retained"; fi
done
bash "$NOTARY_CHECKPOINT_HOOK" "$NOTARY_EVIDENCE_DIR"
