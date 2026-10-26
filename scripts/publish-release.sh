#!/usr/bin/env bash
# Publish only a complete, verified four-target set. A failed upload stays draft;
# rerunning the failed job may resume that draft but never rewrites a release.
set -euo pipefail
PLAN=${1:?usage: publish-release.sh plan-directory downloaded-artifacts}
ARTIFACTS=${2:?usage: publish-release.sh plan-directory downloaded-artifacts}
: "${GH_REPO:?GH_REPO is required}"
: "${EXPECTED_SHA:?EXPECTED_SHA is required}"
: "${RELEASE_CHANNEL:?RELEASE_CHANNEL must be manual or nightly}"
fail() { printf 'publish-release: %s\n' "$*" >&2; exit 1; }
[[ $RELEASE_CHANNEL == manual || $RELEASE_CHANNEL == nightly ]] || fail 'invalid channel'
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
  mac_mode=local
else
  [[ $precedence_version != *-nightly.* ]] || fail 'nightly namespace is reserved'
  mac_mode=release
fi
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
cp -- "$PLAN/notes.md" "$WORK/notes.md"
if [[ $RELEASE_CHANNEL == nightly ]]; then
  printf '\nAll nightly packages are unsigned development builds.\n' >> "$WORK/notes.md"
else
  printf '\nmacOS packages are Developer ID signed and notarized. Windows and Linux packages are unsigned.\n' >> "$WORK/notes.md"
fi

# Distinguish a genuine missing object from authentication/network/API errors.
api_get() {
  local endpoint=$1 output=$2
  if gh api "$endpoint" > "$output" 2> "$WORK/api-error"; then return 0; fi
  if grep -Fq '(HTTP 404)' "$WORK/api-error"; then return 4; fi
  cat "$WORK/api-error" >&2
  return 1
}
# The release-by-tag REST endpoint omits drafts. The authenticated releases list
# includes drafts, so failed uploads can be resumed rather than recreated.
gh api --paginate "repos/$GH_REPO/releases?per_page=100" > "$WORK/releases.json"
jq -se --arg tag "$tag" 'add | map(select(.tag_name == $tag)) |
  if length > 1 then error("multiple releases reserve this tag") else . end' "$WORK/releases.json" > "$WORK/matching-releases.json"
if [[ $(jq length "$WORK/matching-releases.json") == 1 ]]; then
  jq '.[0]' "$WORK/matching-releases.json" > "$WORK/release.json"
  [[ $(jq -r .draft "$WORK/release.json") == true ]] || fail 'published releases are immutable; refusing overwrite'
  resume=true
else
  resume=false
fi
if api_get "repos/$GH_REPO/git/ref/tags/$tag" "$WORK/ref.json"; then
  [[ $(jq -r .object.type "$WORK/ref.json") == commit && $(jq -r .object.sha "$WORK/ref.json") == "$sha" ]] || fail 'existing tag does not match planned source'
else
  status=$?
  [[ $status == 4 ]] || exit "$status"
  # Atomic ref creation rejects a concurrent external publisher. No forced tags.
  gh api --method POST "repos/$GH_REPO/git/refs" -f "ref=refs/tags/$tag" -f "sha=$sha" > /dev/null
fi
if [[ $resume == true ]]; then
  gh release upload "$tag" "$WORK/assets/"* --clobber
else
  gh release create "$tag" "$WORK/assets/"* --verify-tag --draft --latest=false --title "$tag" --notes-file "$WORK/notes.md" --prerelease="$prerelease"
fi
# A resumed draft may contain stale assets. Refuse to publish anything beyond
# this exact verified inventory, and catch incomplete uploads before publication.
for file in "$WORK/assets/"*; do basename -- "$file"; done | LC_ALL=C sort > "$WORK/expected-assets"
gh release view "$tag" --json assets --jq '.assets[].name' | LC_ALL=C sort > "$WORK/remote-assets"
cmp "$WORK/expected-assets" "$WORK/remote-assets" || fail 'remote draft asset inventory differs from verified set'
latest=true
[[ $prerelease == false ]] || latest=false
gh release edit "$tag" --draft=false --latest="$latest" --prerelease="$prerelease" --notes-file "$WORK/notes.md"
printf 'Published %s from %s\n' "$tag" "$sha"
