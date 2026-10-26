#!/usr/bin/env bash
# Safe constants only; raw stapler output stays in private scratch.
macos_ticket_failure_class() {
  local output=$1 result=$2
  if [[ $result == 68 ]] &&
    grep -Eq '^CloudKit (network connection lost|connection reset|connection timed out|service unavailable)$' "$output" &&
    ! grep -Eiq '(reject|invalid|signature|mismatch|not found|does not match|could not validate)' "$output"; then
    printf 'cloudkit-transport\n'
  else
    printf 'unclassified\n'
  fi
}

# An explicit requalification carries the original terminal receipt forever.
# Check both its exact bytes and its candidate/stage/UUID/archive identity.
macos_ticket_parent_valid() {
  local receipt=$1 directory=$2 parent file
  parent=$(jq -er 'if has("requalification_parent_sha256") then
    .requalification_parent_sha256 | if type == "string" then . else error("invalid parent") end
    else "" end' "$receipt") || return 1
  [[ -n $parent ]] || return 0
  [[ $parent =~ ^[0-9a-f]{64}$ ]] || return 1
  file="$directory/blocked-receipt-$parent.json"
  [[ -f $file && ! -L $file && $(shasum -a 256 "$file" | awk '{print $1}') == "$parent" ]] || return 1
  jq -e --slurpfile parent "$file" '
    . as $current | $parent[0] as $p |
    $p.schema == 2 and $p.state == "ticket-blocked" and
    $p.failure_phase == "stapler-validate" and $p.failure_exit == "68" and
    (($p.failure_class // "unclassified") == "unclassified") and
    (($p.requalification_parent_sha256 // "") == "") and
    ($p.verified_sha256 | type == "string" and test("^[0-9a-f]{64}$")) and
    ($p.ticket_attempts | type == "number" and floor == . and . >= 1 and . < 3) and
    (($p.validation_attempts // $p.ticket_attempts) | type == "number" and floor == . and . >= 1 and . < 3) and
    (.validation_attempts | type == "number" and floor == . and . <= 3) and
    (.validation_attempts > ($p.validation_attempts // $p.ticket_attempts)) and
    (["schema","submission_id","filename","submitted_sha256","source_commit","version",
      "architecture","stage","ci_run_id","ci_run_attempt","team_id","signing_identity",
      "verified_filename","verified_sha256","ticket_attempts"] | all(. as $key | $p[$key] == $current[$key]))
  ' "$receipt" >/dev/null 2>&1
}

macos_ticket_verified_receipt() {
  local receipt=$1
  [[ -f $receipt && ! -L $receipt ]] || return 1
  jq -e '.state == "verified" and
    (if has("failure_phase") then .failure_phase == "" else true end) and
    (if has("failure_exit") then .failure_exit == "" else true end) and
    (if has("failure_class") then .failure_class == "" else true end) and
    (if has("validation_attempts") then
      (.validation_attempts | type == "number" and floor == . and . >= 0 and . <= 3)
     else true end)' "$receipt" >/dev/null || return 1
  macos_ticket_parent_valid "$receipt" "${receipt%/*}"
}
