#!/usr/bin/env bash
set -euo pipefail
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
export NOTARY_CHECKPOINT_HOOK="$ROOT/scripts/apple-checkpoint.sh"
# shellcheck source=scripts/macos-ticket-policy.sh
source "$ROOT/scripts/macos-ticket-policy.sh"
macos_acceptance_custody "$NOTARY_EVIDENCE_DIR" || { echo 'Acceptance custody is incomplete or invalid; no mutation permitted.' >&2; exit 1; }
export NOTARY_ACCEPTANCE_FRESH=false
# The signer only hands off this exact approved terminal candidate. It cannot
# consume the fresh runner's reservation or manufacture an acceptance receipt.
if macos_acceptance_exception_eligible "$NOTARY_EVIDENCE_DIR"; then
  parent="$NOTARY_EVIDENCE_DIR/blocked-acceptance-$NOTARY_ACCEPTANCE_REQUALIFY_SHA256.json"
  if [[ ! -e $parent ]]; then cp -p "$NOTARY_EVIDENCE_DIR/acceptance-state.json" "$parent"; fi
  cmp -s "$parent" "$NOTARY_EVIDENCE_DIR/acceptance-state.json"
  rm -f "$NOTARY_EVIDENCE_DIR/acceptance.json" "$NOTARY_EVIDENCE_DIR/fresh-acceptance-$NOTARY_ARCH.json" "$NOTARY_EVIDENCE_DIR/fresh-acceptance-$NOTARY_ARCH.json.pending"
  bash "$NOTARY_CHECKPOINT_HOOK" "$NOTARY_EVIDENCE_DIR"
  echo 'Exact blocked ARM finals retained for approved credential-free qualification; no native attempt spent.'
  exit 0
fi
# Completed exceptional evidence remains reusable for publication, without
# running a fifth check on a signer. Its fresh receipt is still mandatory.
if [[ -f $NOTARY_EVIDENCE_DIR/acceptance-state.json ]] &&
  jq -e 'has("acceptance_exception_parent_sha256")' "$NOTARY_EVIDENCE_DIR/acceptance-state.json" >/dev/null; then
  macos_acceptance_ready "$NOTARY_EVIDENCE_DIR" || exit 1
  bash "$NOTARY_CHECKPOINT_HOOK" "$NOTARY_EVIDENCE_DIR"
  exit 0
fi
status=0
if [[ -e $NOTARY_EVIDENCE_DIR/acceptance-state.json || -L $NOTARY_EVIDENCE_DIR/acceptance-state.json ]]; then
  if ! [[ -f $NOTARY_EVIDENCE_DIR/acceptance-state.json && ! -L $NOTARY_EVIDENCE_DIR/acceptance-state.json &&
     -f $NOTARY_EVIDENCE_DIR/final/SHA256SUMS ]] ||
    ! bash "$ROOT/scripts/apple-stages-verified.sh" "$NOTARY_EVIDENCE_DIR"; then
      echo 'bound final acceptance requires retained finals and verified signer stages; no repack permitted' >&2
      exit 1
  fi
fi
packages="$ROOT/dist/packages/darwin-$NOTARY_ARCH-release"
if [[ ( -f $NOTARY_EVIDENCE_DIR/acceptance.json || -f $NOTARY_EVIDENCE_DIR/acceptance-state.json ) && -f $NOTARY_EVIDENCE_DIR/final/SHA256SUMS ]] &&
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
mkdir -p "$NOTARY_EVIDENCE_DIR/final"
for file in "$packages/"*.zip "$packages/"*.dmg "$packages/SHA256SUMS"; do
  retained="$NOTARY_EVIDENCE_DIR/final/${file##*/}"
  if [[ -e $retained ]]; then cmp -s "$file" "$retained" || { echo 'final candidate bytes changed' >&2; exit 1; }; else cp "$file" "$retained"; fi
done
status=0
bash "$ROOT/scripts/verify-macos-package.sh" || status=$?
if [[ $status == 75 ]]; then
  printf 'Final Apple acceptance pending; retained bytes and persistent budget will resume.\n' >> "${GITHUB_STEP_SUMMARY:?}"
  exit 0
fi
[[ $status == 0 ]] || exit "$status"
bash "$NOTARY_CHECKPOINT_HOOK" "$NOTARY_EVIDENCE_DIR"
