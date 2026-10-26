#!/usr/bin/env bash
set -euo pipefail
exec bash "$(dirname -- "$0")/apple-release-state.sh" save "$@"
