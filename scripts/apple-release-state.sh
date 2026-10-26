#!/usr/bin/env bash
# Durable, append-only candidate custody. Draft assets outlive Actions retention.
set -euo pipefail
umask 077
fail() { printf 'apple-release-state: %s\n' "$*" >&2; exit 1; }
: "${GH_REPO:?GH_REPO is required}"
: "${NOTARY_STATE_TAG:?NOTARY_STATE_TAG is required}"
[[ $GH_REPO =~ ^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$ ]] || fail 'invalid repository'
[[ $NOTARY_STATE_TAG =~ ^apple-candidate-v[0-9A-Za-z.+-]+$ ]] || fail 'invalid candidate tag'
WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-state.XXXXXX")
trap 'rm -rf -- "$WORK"' EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
hash() { shasum -a 256 "$1" | awk '{print $1}'; }
release_identity() {
  gh api "repos/$GH_REPO/releases/$release_id" > "$WORK/exact-release.json"
  jq -e --argjson id "$release_id" --arg tag "$NOTARY_STATE_TAG" --arg source "$release_source" \
    '.id == $id and .draft == true and .tag_name == $tag and .target_commitish == $source and (.body|type == "string")' \
    "$WORK/exact-release.json" >/dev/null || fail 'exact draft identity mismatch'
  jq -j .body "$WORK/exact-release.json" > "$WORK/exact-body"
  [[ $(hash "$WORK/exact-body") == "$release_body_sha" ]] || fail 'exact draft body mismatch'
}
cache_assets() {
  [[ -n ${NOTARY_ASSET_CACHE:-} ]] || return 0
  jq -n --argjson id "$release_id" --arg source "$release_source" --arg body "$release_body_sha" --slurpfile assets "$WORK/assets.json" \
    '{release_id:$id,source:$source,body_sha256:$body,assets:$assets[0]}' > "$NOTARY_ASSET_CACHE.tmp"
  mv "$NOTARY_ASSET_CACHE.tmp" "$NOTARY_ASSET_CACHE"
}
inventory() {
  release_id=${NOTARY_RELEASE_ID:-}
  release_source=${NOTARY_RELEASE_SOURCE_SHA:-}
  release_body_sha=${NOTARY_RELEASE_BODY_SHA256:-}
  if [[ -z $release_id ]]; then
    gh api --paginate "repos/$GH_REPO/releases?per_page=100" > "$WORK/releases.json"
    jq -se --arg tag "$NOTARY_STATE_TAG" 'add | map(select(.tag_name == $tag)) |
      if length == 1 and .[0].draft == true then .[0] else error("candidate draft missing or ambiguous") end' "$WORK/releases.json" > "$WORK/release.json"
    release_id=$(jq -r .id "$WORK/release.json")
    release_source=$(jq -r .target_commitish "$WORK/release.json")
    jq -j .body "$WORK/release.json" > "$WORK/body"
    release_body_sha=$(hash "$WORK/body")
    [[ -z ${NOTARY_RELEASE_BODY_SHA256:-} || $release_body_sha == "$NOTARY_RELEASE_BODY_SHA256" ]] || fail 'draft body pin mismatch'
    [[ -z ${NOTARY_RELEASE_SOURCE_SHA:-} || $release_source == "$NOTARY_RELEASE_SOURCE_SHA" ]] || fail 'draft source pin mismatch'
  fi
  [[ $release_id =~ ^[1-9][0-9]*$ && $release_source =~ ^[0-9a-f]{40}$ && $release_body_sha =~ ^[0-9a-f]{64}$ ]] || fail 'invalid draft identity pins'
  release_identity
  if [[ -n ${NOTARY_ASSET_CACHE:-} && -e $NOTARY_ASSET_CACHE ]]; then
    [[ -f $NOTARY_ASSET_CACHE && ! -L $NOTARY_ASSET_CACHE ]] || fail 'invalid asset cache'
    jq -e --argjson id "$release_id" --arg source "$release_source" --arg body "$release_body_sha" \
      '.release_id == $id and .source == $source and .body_sha256 == $body' "$NOTARY_ASSET_CACHE" >/dev/null || fail 'asset cache binding mismatch'
    jq .assets "$NOTARY_ASSET_CACHE" > "$WORK/assets.json"
  else
    gh api --paginate "repos/$GH_REPO/releases/$release_id/assets?per_page=100" | jq -s 'add' > "$WORK/assets.json"
  fi
  jq -e 'type == "array" and (map(.name) | length == (unique | length)) and
    (map(.id) | length == (unique | length)) and all(.[];
      (.id | type == "number" and . > 0 and floor == .) and
      (.name | type == "string") and (.state == "uploaded" or .state == "starter") and
      (.size | type == "number" and . >= 0 and floor == .))' "$WORK/assets.json" >/dev/null || fail 'invalid or duplicate draft assets'
  cache_assets
}
asset_receipt() {
  local name=$1 size=$2 receipt=$3
  jq -e --arg name "$name" --argjson size "$size" '
    (.id | type == "number" and . > 0 and floor == .) and
    .name == $name and .size == $size and .state == "uploaded"' "$receipt" >/dev/null || fail "invalid upload receipt: $name"
}
download() {
  local name=$1 output=$2 id size
  # Other architectures can have unrelated uploads in flight. The requested
  # asset still needs an acknowledged upload before any of its bytes are used.
  jq -e --arg name "$name" '[.[] | select(.name == $name)] | if length == 1 and .[0].state == "uploaded" then .[0] else error("missing, incomplete or ambiguous asset") end' "$WORK/assets.json" > "$WORK/asset.json"
  id=$(jq -r .id "$WORK/asset.json"); size=$(jq -r .size "$WORK/asset.json")
  release_identity
  gh api --allow-escape-sequences -H 'Accept: application/octet-stream' "repos/$GH_REPO/releases/assets/$id" > "$output"
  [[ $(wc -c < "$output") -eq $size ]] || fail "remote asset size differs: $name"
}
put() {
  local name=$1 file=$2 size encoded
  [[ $name =~ ^[A-Za-z0-9._+-]+$ && -f $file && ! -L $file ]] || fail 'invalid asset name or file'
  if jq -e --arg name "$name" 'any(.[]; .name == $name)' "$WORK/assets.json" >/dev/null; then
    download "$name" "$WORK/existing"
    cmp -s "$file" "$WORK/existing" || fail "immutable asset collision: $name"
  else
    # Use the exact release ID. Never retry a write with an uncertain result.
    release_identity
    size=$(wc -c < "$file")
    encoded=$(jq -nr --arg name "$name" '$name | @uri')
    gh api --method POST "https://uploads.github.com/repos/$GH_REPO/releases/$release_id/assets?name=$encoded" \
      -H 'Content-Type: application/octet-stream' --input "$file" > "$WORK/upload-receipt.json"
    asset_receipt "$name" "$size" "$WORK/upload-receipt.json"
    jq --slurpfile asset "$WORK/upload-receipt.json" '. + $asset | if (map(.id)|length == (unique|length)) then . else error("upload ID collision") end' "$WORK/assets.json" > "$WORK/assets.tmp"
    mv "$WORK/assets.tmp" "$WORK/assets.json"
    download "$name" "$WORK/uploaded"
    cmp -s "$file" "$WORK/uploaded" || fail "remote bytes differ: $name"
    cache_assets
  fi
}
allowed() {
  case $1 in
    input.json|acceptance.json|app/receipt.json|app/apple-log.json|app/*.zip|dmg/receipt.json|dmg/apple-log.json|dmg/*.dmg|final/*.zip|final/*.dmg|final/SHA256SUMS) ;;
    *) return 1 ;;
  esac
  [[ $1 != *..* && $1 != *\\* && $1 != *$'\n'* && $1 != *$'\r'* && $1 != *//* ]]
}
binding() {
  : "${NOTARY_SOURCE_SHA:?}" "${RELEASE_VERSION:?}" "${NOTARY_ARCH:?}" "${GITHUB_RUN_ID:?}" "${GITHUB_RUN_ATTEMPT:?}"
  [[ $NOTARY_SOURCE_SHA =~ ^[0-9a-f]{40}$ && $GITHUB_RUN_ID =~ ^[1-9][0-9]*$ && $GITHUB_RUN_ATTEMPT =~ ^[1-9][0-9]*$ ]] || fail 'invalid candidate identity'
  [[ $NOTARY_ARCH == arm64 || $NOTARY_ARCH == amd64 ]] || fail 'invalid architecture'
  [[ $NOTARY_STATE_TAG == "apple-candidate-v$RELEASE_VERSION" ]] || fail 'candidate/version mismatch'
  download candidate.json "$WORK/candidate.json"
  jq -e --arg source "$NOTARY_SOURCE_SHA" --arg version "$RELEASE_VERSION" --arg repo "$GH_REPO" --arg run "$GITHUB_RUN_ID" --arg attempt "$GITHUB_RUN_ATTEMPT" \
    '.schema == 1 and .source == $source and .version == $version and .repository == $repo and .run_id == $run and .run_attempt == $attempt' "$WORK/candidate.json" >/dev/null || fail 'candidate custody mismatch'
}
operation=${1:?expected save, restore, put, get, or inventory}
shift
inventory
case $operation in
  put) [[ $# == 2 && $1 != */* ]] || fail 'put requires asset basename and file'; put "$1" "$2" ;;
  get) [[ $# == 2 && $1 != */* ]] || fail 'get requires asset basename and output'; download "$1" "$2" ;;
  inventory) cat "$WORK/assets.json" ;;
  save)
    [[ $# == 1 && -d $1 && ! -L $1 ]] || fail 'save requires evidence directory'
    evidence=$1
    binding
    sequence=0
    [[ ! -f $evidence/checkpoint.json ]] || sequence=$(jq -er .sequence "$evidence/checkpoint.json")
    [[ $sequence =~ ^[0-9]+$ ]] || fail 'invalid checkpoint sequence'
    sequence=$((sequence + 1))
    printf -v serial '%06d' "$sequence"
    # The whole release graph is serialized; a second writer for this architecture
    # must not fork a sequence even when its contents differ.
    prefix="checkpoint-$NOTARY_ARCH-$serial-"
    jq -e --arg prefix "$prefix" 'all(.[]; .name | startswith($prefix) | not)' "$WORK/assets.json" >/dev/null || fail 'checkpoint sequence already reserved'
    printf '[]\n' > "$WORK/files.json"
    while IFS= read -r -d '' file; do
      relative=${file#"$evidence/"}
      [[ $relative != checkpoint.json ]] || continue
      allowed "$relative" || fail "unallowlisted evidence: $relative"
      [[ -f $file && ! -L $file ]] || fail 'evidence must contain regular files'
      sha=$(hash "$file")
      asset="blob-$sha-${relative##*/}"
      put "$asset" "$file"
      jq --arg path "$relative" --arg sha "$sha" --arg asset "$asset" '. + [{path:$path,sha256:$sha,asset:$asset}]' "$WORK/files.json" > "$WORK/files.tmp"
      mv "$WORK/files.tmp" "$WORK/files.json"
    done < <(find "$evidence" -type f -print0)
    [[ -z $(find "$evidence" -type l -print) ]] || fail 'symlink in evidence'
    jq -e 'any(.[]; .path == "input.json")' "$WORK/files.json" >/dev/null || fail 'missing input binding'
    jq -n --argjson sequence "$sequence" --arg source "$NOTARY_SOURCE_SHA" --arg version "$RELEASE_VERSION" --arg arch "$NOTARY_ARCH" \
      --arg run "$GITHUB_RUN_ID" --arg attempt "$GITHUB_RUN_ATTEMPT" --slurpfile files "$WORK/files.json" \
      '{schema:1,sequence:$sequence,source:$source,version:$version,architecture:$arch,run_id:$run,run_attempt:$attempt,files:($files[0]|sort_by(.path))}' > "$WORK/checkpoint.json"
    sha=$(hash "$WORK/checkpoint.json")
    put "$prefix$sha.json" "$WORK/checkpoint.json"
    cp "$WORK/checkpoint.json" "$evidence/checkpoint.json"
    ;;
  restore)
    [[ $# == 1 && ! -e $1 ]] || fail 'restore requires a fresh evidence directory'
    evidence=$1
    binding
    jq --arg arch "$NOTARY_ARCH" '[.[] | select(.name | test("^checkpoint-" + $arch + "-[0-9]{6}-[0-9a-f]{64}\\.json$"))] | sort_by(.name)' "$WORK/assets.json" > "$WORK/checkpoints.json"
    mkdir -p "$evidence"
    if [[ $(jq length "$WORK/checkpoints.json") == 0 ]]; then exit 0; fi
    name=$(jq -r '.[-1].name' "$WORK/checkpoints.json")
    serial=${name#"checkpoint-$NOTARY_ARCH-"}; serial=${serial%%-*}
    jq -e --arg prefix "checkpoint-$NOTARY_ARCH-$serial-" '[.[] | select(.name | startswith($prefix))] | length == 1' "$WORK/assets.json" >/dev/null || fail 'forked checkpoint sequence'
    download "$name" "$WORK/checkpoint.json"
    expected=${name##*-}; expected=${expected%.json}
    [[ $(hash "$WORK/checkpoint.json") == "$expected" ]] || fail 'checkpoint hash mismatch'
    jq -e --arg source "$NOTARY_SOURCE_SHA" --arg version "$RELEASE_VERSION" --arg arch "$NOTARY_ARCH" --arg run "$GITHUB_RUN_ID" --arg attempt "$GITHUB_RUN_ATTEMPT" \
      '.schema == 1 and .source == $source and .version == $version and .architecture == $arch and .run_id == $run and .run_attempt == $attempt and (.files|type == "array") and (.files|map(.path)|length == (unique|length))' "$WORK/checkpoint.json" >/dev/null || fail 'checkpoint binding mismatch'
    while IFS= read -r row; do
      relative=$(jq -r .path <<< "$row"); sha=$(jq -r .sha256 <<< "$row"); asset=$(jq -r .asset <<< "$row")
      allowed "$relative" || fail 'unsafe checkpoint path'
      [[ $sha =~ ^[0-9a-f]{64}$ && $asset == "blob-$sha-${relative##*/}" ]] || fail 'invalid blob identity'
      mkdir -p "$(dirname -- "$evidence/$relative")"
      download "$asset" "$evidence/$relative"
      [[ $(hash "$evidence/$relative") == "$sha" ]] || fail 'restored blob hash mismatch'
    done < <(jq -c '.files[]' "$WORK/checkpoint.json")
    cp "$WORK/checkpoint.json" "$evidence/checkpoint.json"
    ;;
  *) fail 'unknown operation' ;;
esac
