#!/usr/bin/env bash
# Offline GitHub custody fixture, including stale reads and interrupted writes.
set -euo pipefail
printf '%s\n' "$*" >> "$MOCK_STATE/calls"
receipt() {
  local file=$1 name=${1##*/} id
  id=$(printf '%s' "$name" | cksum | awk '{print $1+100}')
  jq -n --arg name "$name" --argjson id "$id" --argjson size "$(wc -c < "$file")" '{id:$id,name:$name,size:$size,state:"uploaded"}'
}
# Legacy upload path exists to reproduce the old implementation's failure.
if [[ $1 == release && $2 == upload ]]; then
  file=$4; name=${file##*/}
  [[ ! -e $MOCK_STATE/assets/$name && $name != "${MOCK_FAIL_ASSET:-}" ]] || exit 88
  cp "$file" "$MOCK_STATE/assets/$name"
  touch "$MOCK_STATE/wrote-asset"
  exit 0
fi
[[ $1 == api ]] || exit 99
shift
method=GET; input=''; binary=false; endpoint=''; content_type=''
while [[ $# -gt 0 ]]; do
  case $1 in
    --method|-X) method=$2; shift 2 ;;
    --input) input=$2; shift 2 ;;
    -H|--header) [[ $2 != Content-Type:* ]] || content_type=$2; shift 2 ;;
    --allow-escape-sequences) binary=true; shift ;;
    --paginate) shift ;;
    *) [[ -z $endpoint ]] || exit 99; endpoint=$1; shift ;;
  esac
done
case "$method $endpoint" in
  'POST repos/'*'/releases')
    [[ ! -e $MOCK_STATE/release.json ]] || exit 9
    jq '. + {id:7}' "$input" > "$MOCK_STATE/release.json"
    touch "$MOCK_STATE/created"
    jq "${MOCK_CREATE_FILTER:-.}" "$MOCK_STATE/release.json" ;;
  'GET repos/'*'/releases?per_page=100')
    if [[ ${MOCK_STALE_CREATE:-false} == true && -f $MOCK_STATE/created ]]; then echo '[]'
    elif [[ -f $MOCK_STATE/release.json ]]; then
      if [[ ${MOCK_DUPLICATE_RELEASE:-false} == true ]]; then jq -s . "$MOCK_STATE/release.json" "$MOCK_STATE/release.json"
      else jq -s . "$MOCK_STATE/release.json"; fi
    else echo '[]'; fi ;;
  'GET repos/'*'/releases/7') jq "${MOCK_GET_FILTER:-.}" "$MOCK_STATE/release.json" ;;
  'GET repos/'*'/releases/7/assets?per_page=100')
    if [[ ${MOCK_STALE_UPLOAD:-false} == true && -f $MOCK_STATE/wrote-asset ]]; then
      rm "$MOCK_STATE/wrote-asset"; echo '[]'
    else
      find "$MOCK_STATE/assets/" -type f | LC_ALL=C sort | while IFS= read -r file; do receipt "$file"; done | jq -s .
    fi ;;
  'POST https://uploads.github.com/repos/'*'/releases/7/assets?name='*)
    [[ $content_type == 'Content-Type: application/octet-stream' ]] || exit 91
    name=${endpoint##*name=}
    name=${name//%2B/+}
    # Test names use the unreserved URI alphabet.
    [[ $name =~ ^[A-Za-z0-9._+-]+$ && ! -e $MOCK_STATE/assets/$name && $name != "${MOCK_FAIL_ASSET:-}" ]] || exit 88
    cp "$input" "$MOCK_STATE/assets/$name"
    [[ ${MOCK_CHANGED_BYTES:-false} != true ]] || printf X | dd of="$MOCK_STATE/assets/$name" bs=1 conv=notrunc 2>/dev/null
    touch "$MOCK_STATE/wrote-asset"
    [[ ${MOCK_LOST_RESPONSE:-false} != true ]] || exit 89
    receipt "$MOCK_STATE/assets/$name" | jq "${MOCK_UPLOAD_FILTER:-.}" ;;
  'GET repos/'*'/releases/assets/'*)
    [[ $binary == true || ${MOCK_REQUIRE_BINARY:-false} != true ]] || { echo 'binary response requires --allow-escape-sequences' >&2; exit 90; }
    id=${endpoint##*/}
    found=false
    for file in "$MOCK_STATE/assets/"*; do
      [[ -f $file ]] || continue
      if [[ $(receipt "$file" | jq -r .id) == "$id" ]]; then cat "$file"; found=true; break; fi
    done
    [[ $found == true ]] ;;
  'GET repos/'*'/branches/main') echo '{"protected":true}' ;;
  'GET repos/'*'/compare/'*)
    jq -n --arg sha "$MOCK_AUTOMATION_SHA" '{base_commit:{sha:$sha},merge_base_commit:{sha:$sha},status:"identical"}' ;;
  'GET repos/'*'/commits/'*'/check-runs?per_page=100')
    [[ $endpoint == */commits/"$MOCK_AUTOMATION_SHA"/* ]]
    jq -n --arg sha "$MOCK_AUTOMATION_SHA" --arg conclusion "${MOCK_CHECK:-success}" '{check_runs:[{name:"ci-required",app:{id:15368},check_suite:{id:10},head_sha:$sha,status:"completed",conclusion:$conclusion},{name:"conventional-commits",app:{id:15368},check_suite:{id:11},head_sha:$sha,status:"completed",conclusion:"success"}]}' ;;
  'GET repos/'*'/actions/workflows/'*'/runs?'*)
    [[ $endpoint == *"head_sha=$MOCK_AUTOMATION_SHA"* ]]
    workflow=${endpoint#*/actions/workflows/}; workflow=${workflow%%/*}
    suite=10; [[ $workflow != conventional-commits.yml ]] || suite=11
    jq -n --arg sha "$MOCK_AUTOMATION_SHA" --arg path ".github/workflows/$workflow" --arg repo "$GH_REPO" --argjson suite "$suite" '{workflow_runs:[{id:100,check_suite_id:$suite,head_sha:$sha,head_branch:"main",event:"push",path:$path,repository:{full_name:$repo},head_repository:{full_name:$repo},status:"completed",conclusion:"success"}]}' ;;
  'GET repos/'*'/actions/runs/123/artifacts?per_page=100') jq '{artifacts:.}' "$MOCK_TEST_WORK/metadata.json" ;;
  'GET repos/'*'/actions/artifacts/'*'/zip')
    [[ ${MOCK_ORIGINALS_MISSING:-false} != true ]] || { echo 'Not Found (HTTP 404)' >&2; exit 1; }
    [[ $binary == true || ${MOCK_REQUIRE_BINARY:-false} != true ]] || exit 90
    id=${endpoint%/zip}; id=${id##*/}; cat "$MOCK_TEST_WORK/originals/$id.zip" ;;
  *) echo "unexpected API $method $endpoint" >&2; exit 99 ;;
esac
