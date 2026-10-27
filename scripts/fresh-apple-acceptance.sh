#!/usr/bin/env bash
# Runs on a fresh native runner without Apple credentials or signing permission.
set -euo pipefail
: "${NOTARY_EVIDENCE_DIR:?}" "${NOTARY_ARCH:?}" "${ORIGIN_RUN:?}" "${ORIGIN_ATTEMPT:?}"
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
# shellcheck source=scripts/macos-ticket-policy.sh
source "$ROOT/scripts/macos-ticket-policy.sh"
export NOTARY_CHECKPOINT_HOOK="$ROOT/scripts/apple-checkpoint.sh"
verification_run=$GITHUB_RUN_ID
export NOTARY_ACCEPTANCE_VERIFICATION_RUN_ID=$GITHUB_RUN_ID
export NOTARY_ACCEPTANCE_VERIFICATION_RUN_ATTEMPT=$GITHUB_RUN_ATTEMPT
export NOTARY_ACCEPTANCE_VERIFICATION_WORKFLOW_SHA=${RELEASE_AUTOMATION_SHA:-}
env GITHUB_RUN_ID="$ORIGIN_RUN" GITHUB_RUN_ATTEMPT="$ORIGIN_ATTEMPT" bash "$ROOT/scripts/apple-release-state.sh" restore "$NOTARY_EVIDENCE_DIR"
macos_acceptance_custody "$NOTARY_EVIDENCE_DIR" || { echo 'Acceptance custody is incomplete or invalid; no mutation permitted.' >&2; exit 1; }
rm -f "$NOTARY_EVIDENCE_DIR/fresh-acceptance-$NOTARY_ARCH.json" "$NOTARY_EVIDENCE_DIR/fresh-acceptance-$NOTARY_ARCH.json.pending"
exception=false
if [[ ! -e $NOTARY_EVIDENCE_DIR/acceptance-state.json && ! -L $NOTARY_EVIDENCE_DIR/acceptance-state.json &&
    ( -e $NOTARY_EVIDENCE_DIR/blocked-acceptance-27b7a083cc181ed6def7479ec99b25479d7b11d9c3de36bbbd9facfbacfc60fe.json ||
      -L $NOTARY_EVIDENCE_DIR/blocked-acceptance-27b7a083cc181ed6def7479ec99b25479d7b11d9c3de36bbbd9facfbacfc60fe.json ) ]]; then
  echo 'Acceptance state missing beside immutable exception parent.' >&2
  exit 1
fi
if [[ -e $NOTARY_EVIDENCE_DIR/acceptance-state.json || -L $NOTARY_EVIDENCE_DIR/acceptance-state.json ]]; then
  macos_acceptance_state_valid "$NOTARY_EVIDENCE_DIR/acceptance-state.json" || exit 1
  macos_acceptance_exception_audit "$NOTARY_EVIDENCE_DIR" || exit 1
  if jq -e '.state == "verified" and has("acceptance_exception_parent_sha256")' "$NOTARY_EVIDENCE_DIR/acceptance-state.json" >/dev/null; then
    # A publication/peer retry cannot buy a fifth native qualification. Restore
    # the completed fresh proof from the same immutable checkpoint instead.
    macos_acceptance_ready "$NOTARY_EVIDENCE_DIR" || exit 1
    cp "$NOTARY_EVIDENCE_DIR/acceptance.json" "$NOTARY_EVIDENCE_DIR/fresh-acceptance-$NOTARY_ARCH.json"
    echo 'Retained exact completed ARM fresh qualification; original verification identity preserved, no new native attempt.'
    exit 0
  fi
  if jq -e '.state == "blocked" or has("acceptance_exception_parent_sha256")' "$NOTARY_EVIDENCE_DIR/acceptance-state.json" >/dev/null; then
    macos_acceptance_exception_eligible "$NOTARY_EVIDENCE_DIR" || { echo 'Terminal acceptance is not eligible for this approval.' >&2; exit 1; }
    exception=true
    peer_work=$(mktemp -d "${TMPDIR:-/tmp}/takedown-counterpart.XXXXXX")
    trap 'if [[ -e $peer_work ]]; then rm -rf -- "$peer_work"; fi' EXIT
    export NOTARY_ACCEPTANCE_PEER_DIR="$peer_work/evidence"
    export NOTARY_ACCEPTANCE_PEER_CLEANUP_DIR="$peer_work"
    env NOTARY_ARCH=amd64 GITHUB_RUN_ID="$ORIGIN_RUN" GITHUB_RUN_ATTEMPT="$ORIGIN_ATTEMPT" \
      bash "$ROOT/scripts/apple-release-state.sh" restore "$NOTARY_ACCEPTANCE_PEER_DIR"
    peer_status=0
    macos_acceptance_exception_peer_status "$NOTARY_ACCEPTANCE_PEER_DIR" "$NOTARY_EVIDENCE_DIR" || peer_status=$?
    if [[ $peer_status == 75 ]]; then
      echo 'Intel signer acceptance is pending; approved ARM exception remains unspent.'
      exit 0
    fi
    [[ $peer_status == 0 ]] || { echo 'Intel counterpart custody is invalid or terminal; ARM remains blocked.' >&2; exit 1; }
  fi
fi
if [[ ( $exception == false && ! -f $NOTARY_EVIDENCE_DIR/acceptance.json ) || ! -f $NOTARY_EVIDENCE_DIR/final/SHA256SUMS ]] ||
  ! bash "$ROOT/scripts/apple-stages-verified.sh" "$NOTARY_EVIDENCE_DIR"; then
  echo 'Original candidate is still pending; no fresh acceptance claimed.'
  exit 0
fi
APPLE_TEAM_ID=$(jq -er .team_id "$NOTARY_EVIDENCE_DIR/app/receipt.json")
export APPLE_TEAM_ID
jq -e --arg team "$APPLE_TEAM_ID" '.team_id == $team and .state == "verified"' "$NOTARY_EVIDENCE_DIR/dmg/receipt.json" >/dev/null
packages="$ROOT/dist/packages/darwin-$NOTARY_ARCH-release"
mkdir -p "$packages"
cp "$NOTARY_EVIDENCE_DIR/final/"* "$packages/"
status=0
env NOTARY_ACCEPTANCE_FRESH=true GITHUB_RUN_ID="$ORIGIN_RUN" GITHUB_RUN_ATTEMPT="$ORIGIN_ATTEMPT" bash "$ROOT/scripts/verify-macos-package.sh" || status=$?
if [[ $status == 75 ]]; then
  echo 'Fresh native acceptance pending; no fresh acceptance claimed.'
  exit 0
fi
[[ $status == 0 ]] || exit "$status"
if [[ $exception == true ]]; then
  cp "$NOTARY_EVIDENCE_DIR/acceptance.json" "$NOTARY_EVIDENCE_DIR/fresh-acceptance-$NOTARY_ARCH.json"
else
  jq --arg run "$verification_run" --arg attempt "$GITHUB_RUN_ATTEMPT" '. + {verification_run_id:$run,verification_run_attempt:$attempt,credential_free_runner:true}' "$NOTARY_EVIDENCE_DIR/acceptance.json" > "$NOTARY_EVIDENCE_DIR/fresh-acceptance-$NOTARY_ARCH.json"
fi
