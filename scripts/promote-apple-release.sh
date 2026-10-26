#!/usr/bin/env bash
set -euo pipefail
umask 077
: "${ORIGIN_RUN:?}" "${ORIGIN_ATTEMPT:?}" "${EXPECTED_SHA:?}" "${RELEASE_VERSION:?}"
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
STATE="$ROOT/scripts/apple-release-state.sh"
WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-promote.XXXXXX")
trap 'rm -rf -- "$WORK"' EXIT
mkdir "$WORK/plan" "$WORK/artifacts"
for file in plan.json notes.md native.json candidate.json; do bash "$STATE" get "$file" "$WORK/plan/$file"; done
for file in plan.json notes.md native.json; do
  key=${file%%.*}_sha256
  [[ $(shasum -a 256 "$WORK/plan/$file" | awk '{print $1}') == "$(jq -r --arg key "$key" '.[$key]' "$WORK/plan/candidate.json")" ]] || { echo 'promotion candidate manifest mismatch' >&2; exit 1; }
done
export NOTARY_SOURCE_SHA=$EXPECTED_SHA
for arch in arm64 amd64; do
  env NOTARY_ARCH="$arch" GITHUB_RUN_ID="$ORIGIN_RUN" GITHUB_RUN_ATTEMPT="$ORIGIN_ATTEMPT" bash "$STATE" restore "$WORK/artifacts/macos-evidence-$arch"
  bash "$STATE" get "fresh-acceptance-$arch-$GITHUB_RUN_ID-$GITHUB_RUN_ATTEMPT.json" "$WORK/artifacts/macos-evidence-$arch/fresh-acceptance.json"
done
for platform in Linux-X64 Windows-X64; do
  asset="unsigned-native-$platform-$EXPECTED_SHA.zip"
  bash "$STATE" get "$asset" "$WORK/native.zip"
  expected=$(jq -er --arg name "${asset%.zip}" '[.[] | select(.name == $name)] | if length == 1 then .[0].digest else error("invalid native inventory") end' "$WORK/plan/native.json")
  [[ sha256:$(shasum -a 256 "$WORK/native.zip" | awk '{print $1}') == "$expected" ]] || { echo 'native publication archive mismatch' >&2; exit 1; }
  # Only original digest-verified archives produced by the credential-free job.
  mkdir "$WORK/artifacts/$platform"
  unzip -q "$WORK/native.zip" -d "$WORK/artifacts/$platform"
done
bash "$ROOT/scripts/publish-release.sh" "$WORK/plan" "$WORK/artifacts"
jq -n --arg source "$EXPECTED_SHA" --arg version "$RELEASE_VERSION" --arg run "$ORIGIN_RUN" '{source:$source,version:$version,run_id:$run}' > "$WORK/published.json"
bash "$STATE" put published.json "$WORK/published.json"
