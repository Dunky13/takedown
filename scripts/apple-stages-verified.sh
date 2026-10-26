#!/usr/bin/env bash
# State gate only. Callers still perform all identity, trust and acceptance checks.
set -euo pipefail
# shellcheck source=scripts/macos-ticket-policy.sh
source "$(dirname -- "$0")/macos-ticket-policy.sh"
for stage in app dmg; do
  receipt="${1:?evidence directory required}/$stage/receipt.json"
  macos_ticket_verified_receipt "$receipt"
done
