#!/usr/bin/env bash
# Runs on a fresh native runner without Apple credentials or signing permission.
set -euo pipefail
: "${NOTARY_EVIDENCE_DIR:?}" "${NOTARY_ARCH:?}" "${ORIGIN_RUN:?}" "${ORIGIN_ATTEMPT:?}"
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
verification_run=$GITHUB_RUN_ID
env GITHUB_RUN_ID="$ORIGIN_RUN" GITHUB_RUN_ATTEMPT="$ORIGIN_ATTEMPT" bash "$ROOT/scripts/apple-release-state.sh" restore "$NOTARY_EVIDENCE_DIR"
if [[ ! -f $NOTARY_EVIDENCE_DIR/acceptance.json || ! -f $NOTARY_EVIDENCE_DIR/final/SHA256SUMS ]]; then
  echo 'Original candidate is still pending; no fresh acceptance claimed.'
  exit 0
fi
APPLE_TEAM_ID=$(jq -er .team_id "$NOTARY_EVIDENCE_DIR/app/receipt.json")
export APPLE_TEAM_ID
jq -e --arg team "$APPLE_TEAM_ID" '.team_id == $team and .state == "verified"' "$NOTARY_EVIDENCE_DIR/dmg/receipt.json" >/dev/null
packages="$ROOT/dist/packages/darwin-$NOTARY_ARCH-release"
mkdir -p "$packages"
cp "$NOTARY_EVIDENCE_DIR/final/"* "$packages/"
env GITHUB_RUN_ID="$ORIGIN_RUN" GITHUB_RUN_ATTEMPT="$ORIGIN_ATTEMPT" bash "$ROOT/scripts/verify-macos-package.sh"
jq --arg run "$verification_run" --arg attempt "$GITHUB_RUN_ATTEMPT" '. + {verification_run_id:$run,verification_run_attempt:$attempt,credential_free_runner:true}' "$NOTARY_EVIDENCE_DIR/acceptance.json" > "$NOTARY_EVIDENCE_DIR/fresh-acceptance-$NOTARY_ARCH.json"
