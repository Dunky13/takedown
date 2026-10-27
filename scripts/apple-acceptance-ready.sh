#!/usr/bin/env bash
# Readiness supplements the original signer and native acceptance gates.
set -euo pipefail
# shellcheck source=scripts/macos-ticket-policy.sh
source "$(dirname -- "$0")/macos-ticket-policy.sh"
macos_acceptance_ready "${1:?evidence directory required}"
