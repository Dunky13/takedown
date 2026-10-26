#!/usr/bin/env bash
# gh api --paginate emits one JSON array per page; slurp and combine every page.
set -euo pipefail
jq -sr 'add | map(select(.draft == false and .prerelease == true and
  (.tag_name | test("^v[0-9]+\\.[0-9]+\\.[0-9]+-nightly\\.[0-9]{14}\\.[0-9a-f]{12}$")))) |
  sort_by(.published_at) | last | .tag_name // ""' "${1:?usage: select-last-nightly.sh paginated-releases.json}"
