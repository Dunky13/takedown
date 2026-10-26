#!/usr/bin/env bash
# Credential-free handoff from this release run's immutable native artifact.
set -euo pipefail
umask 077
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
cd "$ROOT"
fail() { printf 'prepare-macos-signing: %s\n' "$*" >&2; exit 1; }
for name in INPUT_SOURCE_SHA GITHUB_SHA GITHUB_REPOSITORY GITHUB_RUN_ID \
  GITHUB_RUN_ATTEMPT GITHUB_OUTPUT VERSION RELEASE_VERSION NOTARY_ARCH NOTARY_STAGING_DIR; do
  [[ -n ${!name:-} ]] || fail "missing configuration: $name"
done
[[ $INPUT_SOURCE_SHA =~ ^[0-9a-f]{40}$ && $INPUT_SOURCE_SHA == "$GITHUB_SHA" ]] || fail 'source SHA must match the release run'
[[ $GITHUB_RUN_ID =~ ^[1-9][0-9]*$ && $GITHUB_RUN_ATTEMPT =~ ^[1-9][0-9]*$ ]] || fail 'invalid run identity'
[[ $VERSION =~ ^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$ ]] || fail 'invalid package version'
[[ $RELEASE_VERSION =~ ^[0-9A-Za-z.+-]+$ && ${RELEASE_VERSION%%[-+]*} == "$VERSION" ]] || fail 'release and package versions differ'
bash "$ROOT/scripts/release-source.sh" verify "$INPUT_SOURCE_SHA" || fail 'checkout must be clean at the release source SHA or qualified automation SHA'
[[ -z $(git status --porcelain) ]] || fail 'checkout must be clean at the release source SHA'
case $NOTARY_ARCH in
  arm64) runner_arch=ARM64; slice=arm64 ;;
  amd64) runner_arch=X64; slice=x86_64 ;;
  *) fail 'unsupported architecture' ;;
esac
[[ $NOTARY_STAGING_DIR == /* && $NOTARY_STAGING_DIR != *$'\n'* && $NOTARY_STAGING_DIR != *$'\r'* ]] || fail 'staging path must be absolute and single-line'
mkdir "$NOTARY_STAGING_DIR" || fail 'staging directory must be fresh'
stage=$NOTARY_STAGING_DIR
artifact_name="unsigned-native-macOS-$runner_arch-$INPUT_SOURCE_SHA"
if [[ -n ${NOTARY_NATIVE_METADATA:-} ]]; then
  [[ -f $NOTARY_NATIVE_METADATA && ! -L $NOTARY_NATIVE_METADATA && -f ${NOTARY_NATIVE_ARCHIVE:-} && ! -L $NOTARY_NATIVE_ARCHIVE ]] || fail 'durable native input missing'
  # The candidate manifest hash and immutable draft asset bytes are checked by
  # the orchestrator before this original Actions metadata is consumed.
  jq '{artifacts:.}' "$NOTARY_NATIVE_METADATA" > "$stage/artifacts.json"
else
  gh api --paginate "repos/$GITHUB_REPOSITORY/actions/runs/$GITHUB_RUN_ID/artifacts?per_page=100" > "$stage/artifacts.json"
fi
jq -s --arg name "$artifact_name" '[.[].artifacts[] | select(.name == $name)]' \
  "$stage/artifacts.json" > "$stage/matches.json"
jq -e --arg sha "$INPUT_SOURCE_SHA" --argjson run "$GITHUB_RUN_ID" '
  length == 1 and (.[0] | .expired == false and
    (.id | type == "number" and . > 0 and floor == .) and
    (.digest | type == "string" and test("^sha256:[0-9a-f]{64}$")) and
    .workflow_run.id == $run and .workflow_run.head_sha == $sha)
' "$stage/matches.json" >/dev/null || fail 'expected one unexpired artifact bound to this source and run'
artifact_id=$(jq -r '.[0].id' "$stage/matches.json")
artifact_digest=$(jq -r '.[0].digest' "$stage/matches.json")
if [[ -n ${NOTARY_NATIVE_METADATA:-} ]]; then
  cp "$NOTARY_NATIVE_ARCHIVE" "$stage/native-artifact.zip"
else
  gh api --allow-escape-sequences "repos/$GITHUB_REPOSITORY/actions/artifacts/$artifact_id/zip" > "$stage/native-artifact.zip"
fi
archive_sha=$(shasum -a 256 "$stage/native-artifact.zip" | awk '{print $1}')
[[ sha256:$archive_sha == "$artifact_digest" ]] || fail 'artifact archive digest mismatch'

base="takedown-$RELEASE_VERSION-macos-$NOTARY_ARCH-local"
directory="darwin-$NOTARY_ARCH-local"
# Select only the expected regular payload bytes instead of extracting the
# outer artifact's paths. Reject extra or duplicate entries before reading any.
unzip -Z -1 "$stage/native-artifact.zip" > "$stage/outer-entries.txt"
zip_count=0; dmg_count=0; manifest_count=0
while IFS= read -r entry; do
  case $entry in
    "$directory/") ;;
    "$directory/$base.zip") zip_count=$((zip_count + 1)) ;;
    "$directory/$base.dmg") dmg_count=$((dmg_count + 1)) ;;
    "$directory/SHA256SUMS") manifest_count=$((manifest_count + 1)) ;;
    *) fail 'unexpected artifact archive entry' ;;
  esac
done < "$stage/outer-entries.txt"
[[ $zip_count == 1 && $dmg_count == 1 && $manifest_count == 1 ]] || fail 'artifact must contain exactly one ZIP, DMG and checksum manifest'
mkdir "$stage/packages" "$stage/extracted"
for name in "$base.zip" "$base.dmg" SHA256SUMS; do
  unzip -p "$stage/native-artifact.zip" "$directory/$name" > "$stage/packages/$name"
done
records=0
zip_sha=''; dmg_sha=''
while IFS= read -r record; do
  [[ $record =~ ^([0-9a-f]{64})\ \ (.+)$ ]] || fail 'invalid native checksum record'
  checksum=${BASH_REMATCH[1]}; name=${BASH_REMATCH[2]}
  case $name in
    "$base.zip") [[ -z $zip_sha ]] || fail 'duplicate ZIP checksum'; zip_sha=$checksum ;;
    "$base.dmg") [[ -z $dmg_sha ]] || fail 'duplicate DMG checksum'; dmg_sha=$checksum ;;
    *) fail 'unexpected native checksum filename' ;;
  esac
  records=$((records + 1))
done < "$stage/packages/SHA256SUMS"
[[ $records == 2 && -n $zip_sha && -n $dmg_sha ]] || fail 'native checksum manifest must cover ZIP and DMG'
(cd "$stage/packages" && shasum -a 256 -c SHA256SUMS) || fail 'native package checksum mismatch'
unzip -Z -1 "$stage/packages/$base.zip" > "$stage/app-entries.txt"
[[ -z $(sort "$stage/app-entries.txt" | uniq -d) ]] || fail 'duplicate application archive entry'
while IFS= read -r entry; do
  case $entry in
    takedown.app|takedown.app/|takedown.app/*|__MACOSX/|__MACOSX/._takedown.app|__MACOSX/takedown.app|__MACOSX/takedown.app/|__MACOSX/takedown.app/*) ;;
    *) fail 'unexpected application archive path' ;;
  esac
  case /$entry/ in
    */../*|*/./*|*\\*|*$'\r'*) fail 'unsafe application archive path' ;;
  esac
done < "$stage/app-entries.txt"
unzip -Z -l "$stage/packages/$base.zip" > "$stage/app-entry-modes.txt"
if awk 'substr($1, 1, 1) == "l" { found=1 } END { exit !found }' "$stage/app-entry-modes.txt"; then
  fail 'application archive must not contain symbolic links'
fi
ditto -x -k "$stage/packages/$base.zip" "$stage/extracted"
app="$stage/extracted/takedown.app"
plist="$app/Contents/Info.plist"
binary="$app/Contents/MacOS/takedown"
[[ -f $plist && ! -L $plist && -f $binary && ! -L $binary ]] || fail 'missing regular app metadata or executable'
[[ -z $(find "$stage/extracted" -type l -print) ]] || fail 'extracted application contains a symbolic link'
for field in CFBundleIdentifier CFBundleExecutable CFBundlePackageType CFBundleShortVersionString CFBundleVersion; do
  case $field in
    CFBundleIdentifier) expected=io.github.dunky13.takedown ;;
    CFBundleExecutable) expected=takedown ;;
    CFBundlePackageType) expected=APPL ;;
    *) expected=$VERSION ;;
  esac
  [[ $(plutil -extract "$field" raw -o - "$plist") == "$expected" ]] || fail "application metadata mismatch: $field"
done
[[ $(lipo -archs "$binary") == "$slice" ]] || fail 'binary architecture does not match target'
go version -m "$binary" > "$stage/build-info.txt"
grep -Eq "build[[:space:]]+vcs.revision=$INPUT_SOURCE_SHA$" "$stage/build-info.txt" || fail 'binary source revision mismatch'
grep -Eq 'build[[:space:]]+vcs.modified=false$' "$stage/build-info.txt" || fail 'binary was not built from clean source'
[[ $("$binary" --version) == "takedown version $RELEASE_VERSION" ]] || fail 'binary release version mismatch'
binary_sha=$(shasum -a 256 "$binary" | awk '{print $1}')
jq -n --arg source "$INPUT_SOURCE_SHA" --arg version "$RELEASE_VERSION" \
  --arg package_version "$VERSION" --arg architecture "$NOTARY_ARCH" \
  --arg repository "$GITHUB_REPOSITORY" --arg run "$GITHUB_RUN_ID" \
  --arg attempt "$GITHUB_RUN_ATTEMPT" --arg artifact_id "$artifact_id" \
  --arg artifact_digest "$artifact_digest" --arg zip_sha "$zip_sha" \
  --arg dmg_sha "$dmg_sha" --arg binary_sha "$binary_sha" \
  '{source:$source,version:$version,package_version:$package_version,
    architecture:$architecture,bundle_id:"io.github.dunky13.takedown",
    repository:$repository,run_id:$run,run_attempt:$attempt,
    artifact_id:$artifact_id,artifact_digest:$artifact_digest,
    native_zip_sha256:$zip_sha,native_dmg_sha256:$dmg_sha,binary_sha256:$binary_sha}' > "$stage/input.json"
printf 'binary=%s\ninput_receipt=%s\n' "$binary" "$stage/input.json" >> "$GITHUB_OUTPUT"
printf 'Verified native artifact %s for %s (%s).\n' "$artifact_id" "$RELEASE_VERSION" "$NOTARY_ARCH"
