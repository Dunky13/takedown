#!/usr/bin/env bash
# Safe constants only; raw stapler output stays in private scratch.
macos_ticket_failure_class() {
  local output=$1 result=$2
  if [[ $result == 68 ]] && {
    grep -Eq '^CloudKit (network connection lost|connection reset|connection timed out|service unavailable)$' "$output" || {
      # Native stapler emits a multiline NSError, not the short fixture marker.
      # Match the observed error AND the exact Apple ticket-delivery endpoint.
      grep -Fq 'error is Error Domain=NSURLErrorDomain Code=-1005 "The network connection was lost." UserInfo={' "$output" &&
      grep -Eq 'NSErrorFailingURLStringKey=https://api\.apple-cloudkit\.com/database/1/com\.apple\.gk\.ticket-delivery/production/public/records/lookup[,}]' "$output" &&
      grep -Fxq "CloudKit's response is inconsistent with expections: (null)" "$output"
    }
  } &&
    ! grep -Eiq '(reject|invalid|signature|mismatch|not found|does not match|could not validate|could not find|no ticket|not stapled|not notarized|revoked|denied)' "$output"; then
    printf 'cloudkit-transport\n'
  else
    printf 'unclassified\n'
  fi
}

# Optional for legacy evidence. New final acceptance is independently bounded;
# it must never reopen a blocked signer receipt or reset its ticket budgets.
macos_acceptance_state_valid() {
  jq -e '.schema == 1 and (.state | IN("pending","blocked","verified")) and
    (.attempts | type == "number" and floor == . and . >= 0 and . <= 3) and
    (.failure_class | IN("","unclassified","cloudkit-transport")) and
    (.source_commit | type == "string" and test("^[0-9a-f]{40}$")) and
    (.version | type == "string") and (.architecture | IN("arm64","amd64")) and
    (.ci_run_id | type == "string" and test("^[1-9][0-9]*$")) and
    (.ci_run_attempt | type == "string" and test("^[1-9][0-9]*$")) and
    (.team_id | type == "string" and test("^[A-Z0-9]{10}$")) and
    ([.zip_sha256,.dmg_sha256,.input_sha256,.app_receipt_sha256,.dmg_receipt_sha256] |
      all(.[]; type == "string" and test("^[0-9a-f]{64}$"))) and
    (if .state == "verified" then .failure_class == ""
     elif .state == "pending" then .failure_class == "" or
       (.failure_class == "cloudkit-transport" and .attempts > 0 and .attempts < 3)
     else .failure_class != "" end)' "$1" >/dev/null
}

macos_acceptance_ready() {
  local directory=$1 file="$1/acceptance-state.json" name key sha
  [[ -e $file || -L $file ]] || return 0
  [[ -f $file && ! -L $file ]] || return 1
  macos_acceptance_state_valid "$file" || return 1
  jq -e '.state == "verified"' "$file" >/dev/null || return 1
  for name in input.json app/receipt.json dmg/receipt.json; do
    case $name in input.json) key=input_sha256 ;; app/*) key=app_receipt_sha256 ;; dmg/*) key=dmg_receipt_sha256 ;; esac
    [[ -f $directory/$name && ! -L $directory/$name ]] || return 1
    sha=$(shasum -a 256 "$directory/$name" | awk '{print $1}') || return 1
    jq -e --arg key "$key" --arg sha "$sha" '.[$key] == $sha' "$file" >/dev/null || return 1
  done
  [[ -f $directory/acceptance.json && ! -L $directory/acceptance.json ]] || return 1
  jq -e --slurpfile state "$file" '. as $a | $state[0] as $s |
    ["source_commit","version","architecture","ci_run_id","ci_run_attempt","zip_sha256","dmg_sha256"] |
    all(.[]; . as $key | $a[$key] == $s[$key])' "$directory/acceptance.json" >/dev/null
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
