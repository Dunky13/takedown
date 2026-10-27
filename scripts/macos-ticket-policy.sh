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
    (.attempts | type == "number" and floor == . and . >= 0 and . <= 4) and
    (if has("acceptance_exception_parent_sha256") then
      .attempts == 4 and .acceptance_exception_parent_sha256 == "27b7a083cc181ed6def7479ec99b25479d7b11d9c3de36bbbd9facfbacfc60fe"
     else .attempts <= 3 end) and
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

# This is a single approved receipt, never a configurable retry ceiling. Its
# digest seals the original input, stage UUIDs/receipts, team and final bytes.
macos_acceptance_exception_parent() {
  local file=$1
  [[ -f $file && ! -L $file && $(shasum -a 256 "$file" | awk '{print $1}') == 27b7a083cc181ed6def7479ec99b25479d7b11d9c3de36bbbd9facfbacfc60fe ]] || return 1
  macos_acceptance_state_valid "$file" || return 1
  jq -e '.state == "blocked" and .attempts == 3 and .failure_class == "cloudkit-transport" and
    (has("acceptance_exception_parent_sha256") | not) and
    .architecture == "arm64" and .version == "0.0.1" and
    .source_commit == "ef82381d509aac8b5b02da0f784a8176b1c5b9c0" and
    .ci_run_id == "40470756538" and .ci_run_attempt == "1"' "$file" >/dev/null
}

macos_acceptance_exception_audit() {
  local directory=$1 file="$1/acceptance-state.json" parent
  parent="$directory/blocked-acceptance-27b7a083cc181ed6def7479ec99b25479d7b11d9c3de36bbbd9facfbacfc60fe.json"
  if ! jq -e 'has("acceptance_exception_parent_sha256")' "$file" >/dev/null; then
    [[ -e $parent || -L $parent ]] || return 0
    # The only pre-reservation parent checkpoint retains the exact blocked3.
    macos_acceptance_exception_parent "$file" && cmp -s "$file" "$parent" && [[ ! -L $parent ]]
    return $?
  fi
  macos_acceptance_exception_parent "$parent" || return 1
  macos_acceptance_state_valid "$file" || return 1
  jq -e --slurpfile parent "$parent" '. as $s | $parent[0] |
    del(.state,.attempts,.failure_class) | to_entries | all(.[]; $s[.key] == .value)' "$file" >/dev/null
}

# Strict evidence/byte audit for the exception and its counterpart. No native
# trust or provider call belongs here. Existing ordinary/legacy gates stay intact.
macos_acceptance_bound_finals() {
  local directory=$1 arch=$2 state="$1/acceptance-state.json" stage receipt file sha key version
  [[ -d $directory && ! -L $directory && -f $state && ! -L $state && -z $(find "$directory" -type l -print) ]] || return 1
  macos_acceptance_state_valid "$state" || return 1
  macos_acceptance_exception_audit "$directory" || return 1
  jq -e --arg arch "$arch" '.architecture == $arch and (.version | test("^[0-9A-Za-z.+-]+$"))' "$state" >/dev/null || return 1
  version=$(jq -er .version "$state") || return 1
  for file in input.json app/receipt.json dmg/receipt.json; do
    case $file in input.json) key=input_sha256 ;; app/*) key=app_receipt_sha256 ;; dmg/*) key=dmg_receipt_sha256 ;; esac
    [[ -f $directory/$file && ! -L $directory/$file ]] || return 1
    sha=$(shasum -a 256 "$directory/$file" | awk '{print $1}') || return 1
    jq -e --arg key "$key" --arg sha "$sha" '.[$key] == $sha' "$state" >/dev/null || return 1
  done
  jq -e --slurpfile state "$state" '. as $i | $state[0] as $s |
    $i.source == $s.source_commit and $i.version == $s.version and $i.package_version == $s.version and
    $i.architecture == $s.architecture and $i.run_id == $s.ci_run_id and $i.run_attempt == $s.ci_run_attempt and
    $i.bundle_id == "io.github.dunky13.takedown" and $i.repository == "Dunky13/takedown"' "$directory/input.json" >/dev/null || return 1
  for stage in app dmg; do
    receipt="$directory/$stage/receipt.json"
    macos_ticket_verified_receipt "$receipt" || return 1
    jq -e --arg stage "$stage" --slurpfile state "$state" '. as $r | $state[0] as $s |
      $r.schema == 2 and $r.stage == $stage and
      $r.filename == ("takedown-" + $s.version + "-macos-" + $s.architecture + "-" + $s.source_commit[0:12] + "-" + $s.ci_run_id + "-" + $s.ci_run_attempt + "-" + $stage + (if $stage == "app" then ".zip" else ".dmg" end)) and
      (["source_commit","version","architecture","ci_run_id","ci_run_attempt","team_id"] |
        all(.[]; . as $key | $r[$key] == $s[$key])) and
      ($r.submission_id | type == "string" and test("^[0-9a-fA-F]{8}(-[0-9a-fA-F]{4}){3}-[0-9a-fA-F]{12}$")) and
      ($r.submitted_sha256 | type == "string" and test("^[0-9a-f]{64}$")) and
      ($r.ticket_attempts | type == "number" and floor == . and . >= 0 and . <= 3) and
      $r.signing_identity == ("Developer ID Application: Marc Went (" + $s.team_id + ")") and
      $r.verified_filename == (if $stage == "app" then "stapled-app.zip" else "stapled.dmg" end) and
      $r.verified_sha256 == (if $stage == "app" then $s.zip_sha256 else $s.dmg_sha256 end)' "$receipt" >/dev/null || return 1
    [[ -f $directory/$stage/apple-log.json && ! -L $directory/$stage/apple-log.json ]] || return 1
    jq -e --slurpfile receipt "$receipt" '.status == "Accepted" and .jobId == $receipt[0].submission_id and
      .archiveFilename == $receipt[0].filename and .sha256 == $receipt[0].submitted_sha256' "$directory/$stage/apple-log.json" >/dev/null || return 1
    file=$(jq -er .verified_filename "$receipt") || return 1
    [[ -f $directory/$stage/$file && ! -L $directory/$stage/$file ]] || return 1
    sha=$(shasum -a 256 "$directory/$stage/$file" | awk '{print $1}') || return 1
    jq -e --arg sha "$sha" '.verified_sha256 == $sha' "$receipt" >/dev/null || return 1
  done
  [[ -f $directory/final/SHA256SUMS && ! -L $directory/final/SHA256SUMS ]] || return 1
  # Exact manifest bytes are reconstructed rather than trusting arbitrary paths.
  for stage in zip dmg; do
    file="takedown-$version-macos-$arch-release.$stage"
    [[ -f $directory/final/$file && ! -L $directory/final/$file ]] || return 1
    sha=$(shasum -a 256 "$directory/final/$file" | awk '{print $1}') || return 1
    jq -e --arg key "${stage}_sha256" --arg sha "$sha" '.[$key] == $sha' "$state" >/dev/null || return 1
  done
  [[ $(LC_ALL=C sort "$directory/final/SHA256SUMS") == "$(jq -r --arg base "takedown-$version-macos-$arch-release" \
    '.dmg_sha256 + "  " + $base + ".dmg", .zip_sha256 + "  " + $base + ".zip"' "$state" | LC_ALL=C sort)" ]]
}

macos_acceptance_exception_eligible() {
  local directory=$1 parent="$1/blocked-acceptance-27b7a083cc181ed6def7479ec99b25479d7b11d9c3de36bbbd9facfbacfc60fe.json"
  [[ ${NOTARY_ACCEPTANCE_REQUALIFY_SHA256:-} == 27b7a083cc181ed6def7479ec99b25479d7b11d9c3de36bbbd9facfbacfc60fe &&
     ${NOTARY_RELEASE_ID:-} == 425202908 && ${GH_REPO:-} == Dunky13/takedown && ${NOTARY_ARCH:-} == arm64 ]] || return 1
  jq -e --arg source "${NOTARY_SOURCE_SHA:-}" --arg version "${RELEASE_VERSION:-}" \
    --arg run "${ORIGIN_RUN:-${GITHUB_RUN_ID:-}}" --arg attempt "${ORIGIN_ATTEMPT:-${GITHUB_RUN_ATTEMPT:-}}" \
    '.source_commit == $source and .version == $version and .ci_run_id == $run and .ci_run_attempt == $attempt' "$directory/acceptance-state.json" >/dev/null || return 1
  macos_acceptance_exception_parent "$directory/acceptance-state.json" || return 1
  if [[ -e $parent || -L $parent ]]; then cmp -s "$parent" "$directory/acceptance-state.json" && [[ ! -L $parent ]] || return 1; fi
  macos_acceptance_bound_finals "$directory" arm64
}

macos_acceptance_exception_peer() {
  local directory=$1 current=$2
  macos_acceptance_bound_finals "$directory" amd64 || return 1
  macos_acceptance_ready "$directory" || return 1
  jq -e --slurpfile current "$current/acceptance-state.json" '. as $p | $current[0] as $s |
    ["source_commit","version","ci_run_id","ci_run_attempt","team_id"] | all(.[]; . as $key | $p[$key] == $s[$key])' \
    "$directory/acceptance-state.json" >/dev/null || return 1
  jq -e '.submission_id == "88951ab7-f126-485b-945c-a11fdbdee902"' "$directory/app/receipt.json" >/dev/null || return 1
  jq -e '.submission_id == "70d28837-3ac6-4de3-86b8-a07c1813af89"' "$directory/dmg/receipt.json" >/dev/null
}

# 75 means incomplete signer work, 1 means invalid/terminal custody. Neither
# status may spend ARM's reservation; only complete, bound native success can.
macos_acceptance_exception_peer_status() {
  local directory=$1 current=$2 stage receipt state pending=false uuid
  [[ -d $directory && ! -L $directory && -z $(find "$directory" -type l -print) ]] || return 1
  for stage in app dmg; do
    receipt="$directory/$stage/receipt.json"
    [[ -f $receipt ]] || return 75
    uuid=88951ab7-f126-485b-945c-a11fdbdee902
    [[ $stage != dmg ]] || uuid=70d28837-3ac6-4de3-86b8-a07c1813af89
    jq -e --arg stage "$stage" --arg uuid "$uuid" --slurpfile current "$current/acceptance-state.json" '. as $r | $current[0] as $s |
      $r.schema == 2 and $r.architecture == "amd64" and $r.stage == $stage and
      (["source_commit","version","ci_run_id","ci_run_attempt","team_id"] | all(.[]; . as $key | $r[$key] == $s[$key])) and
      ($r.submission_id == "" or $r.submission_id == $uuid) and
      ($r.ticket_attempts | type == "number" and floor == . and . >= 0 and . <= 3) and
      ($r.validation_attempts | type == "number" and floor == . and . >= 0 and . <= 3)' "$receipt" >/dev/null || return 1
    state=$(jq -er .state "$receipt") || return 1
    case $state in
      verified) macos_ticket_verified_receipt "$receipt" || return 1 ;;
      prepared|submission-unknown|submitted|pending|accepted|validation-pending) pending=true ;;
      *) return 1 ;;
    esac
  done
  [[ $pending == false ]] || return 75
  [[ -f $directory/acceptance-state.json ]] || return 75
  macos_acceptance_state_valid "$directory/acceptance-state.json" || return 1
  state=$(jq -er .state "$directory/acceptance-state.json") || return 1
  [[ $state != blocked ]] || return 1
  if [[ $state == pending ]]; then
    macos_acceptance_bound_finals "$directory" amd64 || return 1
    return 75
  fi
  macos_acceptance_exception_peer "$directory" "$current"
}

macos_acceptance_custody() {
  local directory=$1 file="$1/acceptance-state.json" name key sha exceptional=false
  # Call BEFORE any receipt deletion, state initialization or native operation.
  # Even pending fresh output retains evidence that the exception was spent.
  for name in acceptance.json fresh-acceptance.json fresh-acceptance-arm64.json fresh-acceptance-amd64.json fresh-acceptance-arm64.json.pending fresh-acceptance-amd64.json.pending; do
    if [[ -e $directory/$name || -L $directory/$name ]]; then
      [[ -f $directory/$name && ! -L $directory/$name ]] || return 1
      jq -e 'type == "object"' "$directory/$name" >/dev/null || return 1
      if jq -e 'has("acceptance_exception_parent_sha256") or has("acceptance_attempts")' "$directory/$name" >/dev/null; then exceptional=true; fi
    fi
  done
  if [[ $exceptional == true ]]; then
    [[ -f $file && ! -L $file ]] || return 1
    jq -e 'has("acceptance_exception_parent_sha256")' "$file" >/dev/null || return 1
  fi
  if [[ ! -e $file && ! -L $file ]]; then
    # This particular immutable build already has a recorded terminal budget;
    # removing all markers cannot turn it into an unqualified legacy build.
    if [[ -f $directory/input.json ]] && jq -e '.source == "ef82381d509aac8b5b02da0f784a8176b1c5b9c0" and
      .version == "0.0.1" and .architecture == "arm64" and .run_id == "40470756538" and .run_attempt == "1"' "$directory/input.json" >/dev/null; then return 1; fi
    [[ ! -e $directory/blocked-acceptance-27b7a083cc181ed6def7479ec99b25479d7b11d9c3de36bbbd9facfbacfc60fe.json &&
       ! -L $directory/blocked-acceptance-27b7a083cc181ed6def7479ec99b25479d7b11d9c3de36bbbd9facfbacfc60fe.json ]]
    return $?
  fi
  [[ -f $file && ! -L $file ]] || return 1
  macos_acceptance_state_valid "$file" || return 1
  macos_acceptance_exception_audit "$directory" || return 1
  if jq -e 'has("acceptance_exception_parent_sha256")' "$file" >/dev/null; then
    macos_acceptance_bound_finals "$directory" arm64 || return 1
  fi
}

macos_acceptance_ready() {
  local directory=$1 file="$1/acceptance-state.json" name key sha exceptional=false
  macos_acceptance_custody "$directory" || return 1
  [[ -f $file ]] || return 0
  if jq -e 'has("acceptance_exception_parent_sha256")' "$file" >/dev/null; then exceptional=true; fi
  jq -e '.state == "verified"' "$file" >/dev/null || return 1
  for name in input.json app/receipt.json dmg/receipt.json; do
    case $name in input.json) key=input_sha256 ;; app/*) key=app_receipt_sha256 ;; dmg/*) key=dmg_receipt_sha256 ;; esac
    [[ -f $directory/$name && ! -L $directory/$name ]] || return 1
    sha=$(shasum -a 256 "$directory/$name" | awk '{print $1}') || return 1
    jq -e --arg key "$key" --arg sha "$sha" '.[$key] == $sha' "$file" >/dev/null || return 1
  done
  [[ -f $directory/acceptance.json && ! -L $directory/acceptance.json ]] || return 1
  jq -e --slurpfile state "$file" '. as $a | $state[0] as $s |
    (if $s | has("acceptance_exception_parent_sha256") then
      $a.acceptance_attempts == 4 and $a.acceptance_exception_parent_sha256 == $s.acceptance_exception_parent_sha256 and
      $a.credential_free_runner == true and
      ($a.verification_run_id | type == "string" and test("^[1-9][0-9]*$")) and
      ($a.verification_run_attempt | type == "string" and test("^[1-9][0-9]*$")) and
      ($a.verification_workflow_sha | type == "string" and test("^[0-9a-f]{40}$"))
     else true end) and
    (["source_commit","version","architecture","ci_run_id","ci_run_attempt","zip_sha256","dmg_sha256"] |
    all(.[]; . as $key | $a[$key] == $s[$key]))' "$directory/acceptance.json" >/dev/null || return 1
  if [[ $exceptional == true ]]; then
    for name in fresh-acceptance.json fresh-acceptance-arm64.json fresh-acceptance-amd64.json; do
      [[ ! -e $directory/$name ]] || jq -e --slurpfile native "$directory/acceptance.json" '. as $a | $native[0] as $n |
        ["source_commit","version","architecture","ci_run_id","ci_run_attempt","zip_sha256","dmg_sha256",
         "acceptance_attempts","acceptance_exception_parent_sha256","credential_free_runner","verification_run_id","verification_run_attempt","verification_workflow_sha"] |
        all(.[]; . as $key | $a[$key] == $n[$key])' "$directory/$name" >/dev/null || return 1
      [[ ! -e $directory/$name ]] || cmp -s "$directory/acceptance.json" "$directory/$name" || return 1
    done
    # The immutable verified checkpoint may have committed during cancellation.
    # Only authoritative original-job success and exact proof bytes close it.
    bash "$(dirname -- "${BASH_SOURCE[0]}")/verify-apple-fresh-completion.sh" "$directory" || return 1
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
