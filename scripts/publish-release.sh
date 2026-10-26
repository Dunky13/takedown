#!/usr/bin/env bash
# Publish only a complete, verified four-target set. A failed upload stays draft;
# rerunning the failed job may resume that draft but never rewrites a release.
set -euo pipefail
PLAN=${1:?usage: publish-release.sh plan-directory downloaded-artifacts}
ARTIFACTS=${2:?usage: publish-release.sh plan-directory downloaded-artifacts}
: "${GH_REPO:?GH_REPO is required}"
: "${EXPECTED_SHA:?EXPECTED_SHA is required}"
: "${RELEASE_CHANNEL:?RELEASE_CHANNEL must be manual, auto or nightly}"
fail() { printf 'publish-release: %s\n' "$*" >&2; exit 1; }
[[ $GH_REPO =~ ^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$ ]] || fail 'invalid repository'
[[ $RELEASE_CHANNEL == manual || $RELEASE_CHANNEL == auto || $RELEASE_CHANNEL == nightly ]] || fail 'invalid channel'
jq -e '(.version | type == "string") and (.package_version | type == "string") and
  (.sha | type == "string") and (.tag | type == "string") and
  (.prerelease | type == "boolean") and .skip == false' "$PLAN/plan.json" >/dev/null
version=$(jq -r .version "$PLAN/plan.json")
package_version=$(jq -r .package_version "$PLAN/plan.json")
sha=$(jq -r .sha "$PLAN/plan.json")
tag=$(jq -r .tag "$PLAN/plan.json")
prerelease=$(jq -r .prerelease "$PLAN/plan.json")
go run -mod=readonly -p 2 ./cmd/releasectl validate-version --version "$version"
precedence_version=${version%%+*}
[[ ${version%%[-+]*} == "$package_version" && $tag == "v$version" ]] || fail 'inconsistent version plan'
[[ $sha =~ ^[0-9a-f]{40}$ && $sha == "$EXPECTED_SHA" && $(git rev-parse HEAD) == "$sha" ]] || fail 'source SHA does not match plan and checkout'
if [[ $RELEASE_CHANNEL == nightly ]]; then
  [[ $precedence_version == *-nightly.* && $prerelease == true ]] || fail 'nightly must be a nightly prerelease'
else
  [[ $precedence_version != *-nightly.* ]] || fail 'nightly namespace is reserved'
fi
mac_mode=release
if [[ $precedence_version == *-* ]]; then
  [[ $prerelease == true ]] || fail 'prerelease flag does not match version'
else
  [[ $prerelease == false ]] || fail 'prerelease flag does not match version'
fi
WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-publish.XXXXXX")
trap 'rm -rf -- "$WORK"' EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
mkdir "$WORK/assets"

# Validate every downloaded checksum manifest before selecting assets. In
# particular, never flatten the per-target SHA256SUMS files over one another.
manifests=0
while IFS= read -r -d '' manifest; do
  manifests=$((manifests + 1))
  entries=0
  while IFS= read -r line; do
    line=${line%$'\r'}
    [[ $line =~ ^[0-9a-fA-F]{64}\ \ (takedown-[A-Za-z0-9.+_-]+)$ ]] || fail 'invalid checksum entry'
    name=${BASH_REMATCH[1]}
    [[ -f ${manifest%/*}/$name && ! -L ${manifest%/*}/$name ]] || fail 'checksum payload missing or symlinked'
    entries=$((entries + 1))
  done < "$manifest"
  [[ $entries -gt 0 ]] || fail 'empty checksum manifest'
  (cd "${manifest%/*}" && sha256sum -c SHA256SUMS)
done < <(find "$ARTIFACTS" -type f -name SHA256SUMS -print0)
[[ $manifests -gt 0 ]] || fail 'no checksum manifests'

expected=(
  "takedown-$version-linux-amd64-local.tar.gz"
  "takedown-$version-windows-amd64-local.zip"
  "takedown-$version-windows-amd64-local-setup.exe"
)
for arch in arm64 amd64; do
  expected+=("takedown-$version-macos-$arch-$mac_mode.zip" "takedown-$version-macos-$arch-$mac_mode.dmg")
done
for name in "${expected[@]}"; do
  matches=()
  while IFS= read -r -d '' file; do matches+=("$file"); done < <(find "$ARTIFACTS" -type f -name "$name" -print0)
  [[ ${#matches[@]} == 1 ]] || fail "expected exactly one asset: $name"
  file=${matches[0]}
  [[ ! -L $file && -s $file ]] || fail "invalid asset: $name"
  # Require the selected asset to be covered by its verified original manifest.
  grep -Eq "^[0-9a-fA-F]{64}  " "${file%/*}/SHA256SUMS" || fail 'asset has no manifest'
  found=false
  while IFS= read -r record; do
    record=${record%$'\r'}
    [[ ${record:66} != "$name" ]] || found=true
  done < "${file%/*}/SHA256SUMS"
  [[ $found == true ]] || fail "asset is not checksummed: $name"
  cp -- "$file" "$WORK/assets/$name"
done
(
  cd "$WORK/assets"
  sha256sum -- takedown-* > SHA256SUMS
  sha256sum -c SHA256SUMS
)
# Custody metadata comes from the original candidate, including on later
# continuation runs. Never bind acceptance to the continuation job's identity.
for metadata in candidate.json native.json; do
  [[ -f $PLAN/$metadata && ! -L $PLAN/$metadata ]] || fail "missing original metadata: $metadata"
done
jq -e --arg source "$sha" --arg version "$version" --arg repo "$GH_REPO" '
  .schema == 1 and .source == $source and .version == $version and .repository == $repo and
  (.run_id | type == "string" and test("^[1-9][0-9]*$")) and
  (.run_attempt | type == "string" and test("^[1-9][0-9]*$"))
' "$PLAN/candidate.json" >/dev/null || fail 'invalid original candidate binding'
origin_run=$(jq -r .run_id "$PLAN/candidate.json")
origin_attempt=$(jq -r .run_attempt "$PLAN/candidate.json")
[[ ${ORIGIN_RUN:-$origin_run} == "$origin_run" && ${ORIGIN_ATTEMPT:-$origin_attempt} == "$origin_attempt" ]] || fail 'original run identity differs from candidate'
for arch in arm64 amd64; do
  evidence="$ARTIFACTS/macos-evidence-$arch"
  for path in input.json acceptance.json fresh-acceptance.json app/receipt.json app/apple-log.json dmg/receipt.json dmg/apple-log.json; do
    [[ -f $evidence/$path && ! -L $evidence/$path ]] || fail "missing macOS evidence: $arch/$path"
  done
  [[ -z $(find "$evidence" -type l -print) ]] || fail 'symlink in macOS evidence'
  runner_arch=ARM64
  [[ $arch != amd64 ]] || runner_arch=X64
  artifact_name="unsigned-native-macOS-$runner_arch-$sha"
  jq -e --arg name "$artifact_name" --arg source "$sha" --argjson run "$origin_run" '
    [.[] | select(.name == $name)] | length == 1 and (.[0] |
      (.id | type == "number" and . > 0 and floor == .) and
      (.digest | type == "string" and test("^sha256:[0-9a-f]{64}$")) and
      .workflow_run.id == $run and .workflow_run.head_sha == $source)
  ' "$PLAN/native.json" >/dev/null || fail "invalid original native artifact: $arch"
  artifact_id=$(jq -r --arg name "$artifact_name" '.[] | select(.name == $name) | .id' "$PLAN/native.json")
  artifact_digest=$(jq -r --arg name "$artifact_name" '.[] | select(.name == $name) | .digest' "$PLAN/native.json")
  jq -e --arg source "$sha" --arg version "$version" --arg package "$package_version" --arg arch "$arch" \
    --arg repo "$GH_REPO" --arg run "$origin_run" --arg attempt "$origin_attempt" \
    --arg id "$artifact_id" --arg digest "$artifact_digest" '
    .source == $source and .version == $version and .package_version == $package and
    .architecture == $arch and .repository == $repo and .bundle_id == "io.github.dunky13.takedown" and
    .run_id == $run and .run_attempt == $attempt and .artifact_id == $id and .artifact_digest == $digest
  ' "$evidence/input.json" >/dev/null || fail "macOS input binding mismatch: $arch"
  for stage in app dmg; do
    extension=zip
    [[ $stage != dmg ]] || extension=dmg
    filename="takedown-$version-macos-$arch-${sha:0:12}-$origin_run-$origin_attempt-$stage.$extension"
    jq -e --arg source "$sha" --arg version "$version" --arg arch "$arch" --arg stage "$stage" \
      --arg run "$origin_run" --arg attempt "$origin_attempt" --arg filename "$filename" '
      .source_commit == $source and .version == $version and .architecture == $arch and .stage == $stage and
      .ci_run_id == $run and .ci_run_attempt == $attempt and .state == "verified" and .filename == $filename and
      (.submission_id | type == "string" and test("^[0-9a-fA-F]{8}(-[0-9a-fA-F]{4}){3}-[0-9a-fA-F]{12}$")) and
      (.submitted_sha256 | type == "string" and test("^[0-9a-f]{64}$"))
    ' "$evidence/$stage/receipt.json" >/dev/null || fail "unverified macOS receipt: $arch/$stage"
    jq -e --slurpfile receipt "$evidence/$stage/receipt.json" '
      .status == "Accepted" and .jobId == $receipt[0].submission_id and
      .archiveFilename == $receipt[0].filename and .sha256 == $receipt[0].submitted_sha256
    ' "$evidence/$stage/apple-log.json" >/dev/null || fail "Apple log binding mismatch: $arch/$stage"
  done
  zip_sha=$(sha256sum "$WORK/assets/takedown-$version-macos-$arch-release.zip" | awk '{print $1}')
  dmg_sha=$(sha256sum "$WORK/assets/takedown-$version-macos-$arch-release.dmg" | awk '{print $1}')
  jq -e --arg source "$sha" --arg version "$version" --arg arch "$arch" \
    --arg run "$origin_run" --arg attempt "$origin_attempt" --arg zip "$zip_sha" --arg dmg "$dmg_sha" '
    .source_commit == $source and .version == $version and .architecture == $arch and
    .ci_run_id == $run and .ci_run_attempt == $attempt and .zip_sha256 == $zip and .dmg_sha256 == $dmg
  ' "$evidence/acceptance.json" >/dev/null || fail "macOS acceptance binding mismatch: $arch"
  jq -e --arg source "$sha" --arg version "$version" --arg arch "$arch" \
    --arg run "$origin_run" --arg attempt "$origin_attempt" --arg zip "$zip_sha" --arg dmg "$dmg_sha" '
    .credential_free_runner == true and
    (.verification_run_id | type == "string" and test("^[1-9][0-9]*$")) and
    (.verification_run_attempt | type == "string" and test("^[1-9][0-9]*$")) and
    .source_commit == $source and .version == $version and .architecture == $arch and
    .ci_run_id == $run and .ci_run_attempt == $attempt and .zip_sha256 == $zip and .dmg_sha256 == $dmg
  ' "$evidence/fresh-acceptance.json" >/dev/null || fail "fresh macOS acceptance binding mismatch: $arch"
done
cp -- "$PLAN/notes.md" "$WORK/notes.md"
printf '\nmacOS packages are Developer ID signed and notarized. Windows and Linux packages are unsigned.\n' >> "$WORK/notes.md"

# Distinguish a genuine missing object from authentication/network/API errors.
api_get() {
  local endpoint=$1 output=$2
  if gh api "$endpoint" > "$output" 2> "$WORK/api-error"; then return 0; fi
  if grep -Fq '(HTTP 404)' "$WORK/api-error"; then return 4; fi
  cat "$WORK/api-error" >&2
  return 1
}
for file in "$WORK/assets/"*; do basename -- "$file"; done | LC_ALL=C sort > "$WORK/expected-assets"
release_lookup() {
  # The authenticated list includes drafts, unlike the release-by-tag endpoint.
  gh api --paginate "repos/$GH_REPO/releases?per_page=100" > "$WORK/releases.json"
  jq -se --arg tag "$tag" 'add | map(select(.tag_name == $tag)) |
    if length > 1 then error("multiple releases reserve this tag") else . end' "$WORK/releases.json" > "$WORK/matching-releases.json"
  jq '.[0] // null' "$WORK/matching-releases.json" > "$WORK/release.json"
}
validate_release() {
  jq -e --arg tag "$tag" --arg sha "$sha" --argjson prerelease "$prerelease" '
    .tag_name == $tag and .target_commitish == $sha and .prerelease == $prerelease and
    (.draft | type == "boolean") and (.id | type == "number" and . > 0 and floor == .)
  ' "$WORK/release.json" >/dev/null || fail 'release tag, target or version metadata differs from plan'
}
validate_tag() {
  [[ $(jq -r .object.type "$WORK/ref.json") == commit && $(jq -r .object.sha "$WORK/ref.json") == "$sha" ]] || fail 'existing tag does not match planned source'
}
verify_remote() {
  local complete=$1 name id expected_hash remote_hash
  gh api --paginate "repos/$GH_REPO/releases/$release_id/assets?per_page=100" > "$WORK/asset-pages.json"
  jq -se 'add | if type == "array" and all(.[];
      (.name | type == "string") and .state == "uploaded" and
      (.id | type == "number" and . > 0 and floor == .)) and
      (map(.name) | length == (unique | length)) and (map(.id) | length == (unique | length))
    then . else error("invalid or duplicate remote assets") end' "$WORK/asset-pages.json" > "$WORK/remote.json"
  jq -r '.[].name' "$WORK/remote.json" | LC_ALL=C sort > "$WORK/remote-assets"
  # Check the whole inventory before downloading or appending anything.
  while IFS= read -r name; do
    grep -Fxq -- "$name" "$WORK/expected-assets" || fail "unexpected remote asset: $name"
  done < "$WORK/remote-assets"
  if [[ $complete == true ]]; then
    cmp -s "$WORK/expected-assets" "$WORK/remote-assets" || fail 'remote asset inventory differs from verified set'
  fi
  while IFS= read -r name; do
    id=$(jq -r --arg name "$name" '.[] | select(.name == $name) | .id' "$WORK/remote.json")
    gh api -H 'Accept: application/octet-stream' "repos/$GH_REPO/releases/assets/$id" > "$WORK/download"
    expected_hash=$(sha256sum "$WORK/assets/$name" | awk '{print $1}')
    remote_hash=$(sha256sum "$WORK/download" | awk '{print $1}')
    if [[ $expected_hash != "$remote_hash" ]] || ! cmp -s "$WORK/assets/$name" "$WORK/download"; then
      fail "remote asset bytes differ: $name"
    fi
  done < "$WORK/remote-assets"
}
release_lookup
exists=false
published=false
if [[ $(jq length "$WORK/matching-releases.json") == 1 ]]; then
  exists=true
  validate_release
  release_id=$(jq -r .id "$WORK/release.json")
  [[ $(jq -r .draft "$WORK/release.json") == true ]] || published=true
fi
tag_exists=false
if api_get "repos/$GH_REPO/git/ref/tags/$tag" "$WORK/ref.json"; then
  validate_tag
  tag_exists=true
else
  status=$?
  [[ $status == 4 ]] || exit "$status"
fi
if [[ $published == true ]]; then
  [[ $tag_exists == true ]] || fail 'published release has no matching tag'
  verify_remote true
  printf 'Already published %s from %s; exact remote bytes verified.\n' "$tag" "$sha"
  exit 0
fi
# Verify every existing draft byte before any remote mutation, including a tag.
[[ $exists == false ]] || verify_remote false
if [[ $tag_exists == false ]]; then
  # Atomic creation rejects a concurrent external publisher; never force a ref.
  gh api --method POST "repos/$GH_REPO/git/refs" -f "ref=refs/tags/$tag" -f "sha=$sha" > /dev/null
fi
# Re-read after creation to verify the server's actual ref, including races.
gh api "repos/$GH_REPO/git/ref/tags/$tag" > "$WORK/ref.json"
validate_tag
if [[ $exists == false ]]; then
  gh release create "$tag" --verify-tag --target "$sha" --draft --latest=false --title "$tag" --notes-file "$WORK/notes.md" --prerelease="$prerelease"
  release_lookup
  validate_release
  release_id=$(jq -r .id "$WORK/release.json")
  [[ $(jq -r .draft "$WORK/release.json") == true ]] || fail 'new release is not a draft'
  verify_remote false
fi
for name in "${expected[@]}" SHA256SUMS; do
  if ! grep -Fxq -- "$name" "$WORK/remote-assets"; then
    # Append only. A concurrent name collision is terminal, preserving its bytes.
    gh release upload "$tag" "$WORK/assets/$name"
  fi
done
# Re-read the complete inventory and download every byte before publication.
verify_remote true
original_id=$release_id
release_lookup
validate_release
[[ $(jq -r .id "$WORK/release.json") == "$original_id" && $(jq -r .draft "$WORK/release.json") == true ]] || fail 'draft changed during publication'
gh api "repos/$GH_REPO/git/ref/tags/$tag" > "$WORK/ref.json"
validate_tag
latest=true
[[ $prerelease == false ]] || latest=false
gh release edit "$tag" --draft=false --latest="$latest" --prerelease="$prerelease" --notes-file "$WORK/notes.md"
printf 'Published %s from %s\n' "$tag" "$sha"
