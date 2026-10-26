#!/usr/bin/env bash
# Offline publication boundary tests. All payloads, Apple evidence and gh API
# responses below are SYNTHETIC ONLY; they establish no real Apple acceptance.
set -euo pipefail
cd -- "$(dirname -- "$0")/.."
WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-publish-tests.XXXXXX")
cleanup() {
  local result=$?
  if [[ $result != 0 && -f $WORK/log ]]; then cat "$WORK/log" >&2; fi
  rm -rf -- "$WORK"
}
trap cleanup EXIT
export GH_REPO=fixture/takedown EXPECTED_SHA MOCK_STATE
EXPECTED_SHA=$(git rev-parse HEAD)
MOCK_STATE="$WORK/state"
mkdir "$WORK/mock-bin" "$WORK/plan"
cat > "$WORK/mock-bin/gh" <<'MOCK'
#!/usr/bin/env bash
set -euo pipefail
printf '%s\n' "$*" >> "$MOCK_STATE/calls"
# Fixture also enforces the no-overwrite contract on every command.
for arg in "$@"; do [[ $arg != --clobber* ]] || { echo 'overwrite forbidden' >&2; exit 1; }; done
if [[ $1 == api ]]; then
  if [[ $2 == --method ]]; then
    [[ $3 == POST && $4 == */git/refs && ! -f $MOCK_STATE/ref ]] || exit 1
    if [[ ${MOCK_TAG_COLLISION:-false} == true ]]; then
      printf '{"object":{"type":"commit","sha":"0000000000000000000000000000000000000000"}}\n' > "$MOCK_STATE/ref"
      echo 'gh: Reference already exists (HTTP 422)' >&2
      exit 1
    fi
    [[ $5 == -f && $6 == ref=refs/tags/* && $7 == -f && $8 == "sha=$EXPECTED_SHA" ]] || exit 1
    printf '{"object":{"type":"commit","sha":"%s"}}\n' "$EXPECTED_SHA" > "$MOCK_STATE/ref"
    cat "$MOCK_STATE/ref"
  elif [[ $2 == --paginate && $3 == */releases\?* ]]; then
    if [[ ${MOCK_API_ERROR:-false} == true ]]; then echo 'gh: Forbidden (HTTP 403)' >&2; exit 1; fi
    # Two pages exercise the authenticated draft discovery path.
    echo '[]'
    if [[ -f $MOCK_STATE/release ]]; then
      jq -s '.' "$MOCK_STATE/release"
      [[ ${MOCK_DUPLICATE_RELEASE:-false} != true ]] || jq -s '.' "$MOCK_STATE/release"
    else
      echo '[]'
    fi
  elif [[ $2 == --paginate && $3 == */releases/1/assets\?* ]]; then
    echo '[]'
    cat "$MOCK_STATE/assets"
    if [[ ${MOCK_EXTRA_ASSET:-false} == true ]]; then
      echo '[{"id":999,"name":"unexpected.exe","state":"uploaded"}]'
    fi
  elif [[ $2 == --allow-escape-sequences && $3 == -H && $4 == 'Accept: application/octet-stream' ]]; then
    id=${5##*/}
    [[ -f $MOCK_STATE/bytes/$id ]] || exit 1
    cat "$MOCK_STATE/bytes/$id"
  elif [[ $2 == */git/ref/tags/* && -f $MOCK_STATE/ref ]]; then
    cat "$MOCK_STATE/ref"
  elif [[ ${MOCK_AUTOMATION_SHA:-} != '' && ( $2 == */branches/main || $2 == */compare/* || $2 == --paginate ) ]]; then
    exec bash "$MOCK_AUTOMATION_FIXTURE" "$@"
  else
    printf 'gh: Not Found (HTTP 404)\n' >&2
    exit 1
  fi
elif [[ $1 == release && $2 == create ]]; then
  [[ ! -f $MOCK_STATE/release && -f $MOCK_STATE/ref ]] || exit 1
  tag=$3; shift 3
  prerelease=''; target=''
  while [[ $# -gt 0 ]]; do
    case $1 in
      --prerelease=*) prerelease=${1#*=} ;;
      --target) target=$2; shift ;;
      --notes-file) cp "$2" "$MOCK_STATE/notes"; shift ;;
    esac
    shift
  done
  [[ $target == "$EXPECTED_SHA" && ( $prerelease == true || $prerelease == false ) ]] || exit 1
  jq -n --arg tag "$tag" --arg target "$target" --argjson prerelease "$prerelease" \
    '{id:1,draft:true,tag_name:$tag,target_commitish:$target,prerelease:$prerelease}' > "$MOCK_STATE/release"
elif [[ $1 == release && $2 == upload ]]; then
  [[ $(jq -r .draft "$MOCK_STATE/release") == true && $# == 4 ]] || exit 1
  item=$4; name=${item##*/}
  if jq -e --arg name "$name" 'any(.[]; .name == $name)' "$MOCK_STATE/assets" >/dev/null; then
    echo 'gh: asset already exists (HTTP 422)' >&2; exit 1
  fi
  [[ ${MOCK_DROP_UPLOAD:-false} != true ]] || exit 0
  id=$((100 + $(jq length "$MOCK_STATE/assets")))
  cp "$item" "$MOCK_STATE/bytes/$id"
  jq --arg name "$name" --argjson id "$id" '. + [{id:$id,name:$name,state:"uploaded"}]' "$MOCK_STATE/assets" > "$MOCK_STATE/assets.tmp"
  mv "$MOCK_STATE/assets.tmp" "$MOCK_STATE/assets"
  if [[ ${MOCK_CORRUPT_UPLOAD:-false} == true || ${MOCK_UPLOAD_COLLISION:-false} == true ]]; then
    printf 'corrupt remote bytes\n' > "$MOCK_STATE/bytes/$id"
  fi
  [[ ${MOCK_FAIL_UPLOAD:-false} != true && ${MOCK_UPLOAD_COLLISION:-false} != true ]] || exit 1
elif [[ $1 == release && $2 == edit ]]; then
  [[ $(jq -r .draft "$MOCK_STATE/release") == true ]] || exit 1
  jq '.draft=false' "$MOCK_STATE/release" > "$MOCK_STATE/release.tmp"
  mv "$MOCK_STATE/release.tmp" "$MOCK_STATE/release"
else
  echo 'unexpected gh fixture command' >&2
  exit 1
fi
MOCK
chmod +x "$WORK/mock-bin/gh"
export PATH="$WORK/mock-bin:$PATH"
# macOS provides shasum rather than sha256sum; emulate only the required interface.
if ! command -v sha256sum >/dev/null; then
  cat > "$WORK/mock-bin/sha256sum" <<'SUM'
#!/usr/bin/env bash
[[ ${1:-} != -- ]] || shift
exec shasum -a 256 "$@"
SUM
  chmod +x "$WORK/mock-bin/sha256sum"
fi
hash() { shasum -a 256 "$1" | awk '{print $1}'; }
reset_fixture() {
  rm -rf "$MOCK_STATE" "$WORK/artifacts"
  mkdir -p "$MOCK_STATE/bytes" "$WORK/artifacts"
  printf '[]\n' > "$MOCK_STATE/assets"
  RELEASE_CHANNEL=$1
  export RELEASE_CHANNEL
  version=${2:-1.2.0-rc.1}
  if [[ $RELEASE_CHANNEL == nightly ]]; then
    version="1.2.1-nightly.20261026022300.${EXPECTED_SHA:0:12}"
  fi
  jq -n --arg version "$version" --arg sha "$EXPECTED_SHA" \
    '{version:$version,package_version:($version|split("+")[0]|split("-")[0]),tag:("v"+$version),sha:$sha,prerelease:($version|split("+")[0]|contains("-")),skip:false}' > "$WORK/plan/plan.json"
  printf '# Release fixture\n' > "$WORK/plan/notes.md"
  names=("linux-amd64-local.tar.gz" "windows-amd64-local.zip" "windows-amd64-local-setup.exe"
    "macos-arm64-release.zip" "macos-arm64-release.dmg" "macos-amd64-release.zip" "macos-amd64-release.dmg")
  for name in "${names[@]}"; do
    directory="$WORK/artifacts/$name"
    mkdir "$directory"
    printf 'SYNTHETIC payload for %s\n' "$name" > "$directory/takedown-$version-$name"
    (cd "$directory" && shasum -a 256 "takedown-$version-$name" > SHA256SUMS)
  done
  # Original build identity is deliberately unrelated to any current job env.
  jq -n --arg source "$EXPECTED_SHA" --arg version "$version" --arg repo "$GH_REPO" \
    '{schema:1,source:$source,version:$version,repository:$repo,run_id:"12345",run_attempt:"2"}' > "$WORK/plan/candidate.json"
  printf '[]\n' > "$WORK/plan/native.json"
  for arch in arm64 amd64; do
    evidence="$WORK/artifacts/macos-evidence-$arch"
    mkdir -p "$evidence/app" "$evidence/dmg"
    runner_arch=ARM64; artifact_id=101
    [[ $arch != amd64 ]] || { runner_arch=X64; artifact_id=102; }
    digest="sha256:$(printf 'synthetic-native-%s' "$arch" | shasum -a 256 | awk '{print $1}')"
    jq --arg name "unsigned-native-macOS-$runner_arch-$EXPECTED_SHA" --arg source "$EXPECTED_SHA" \
      --arg digest "$digest" --argjson id "$artifact_id" \
      '. + [{id:$id,name:$name,digest:$digest,workflow_run:{id:12345,head_sha:$source}}]' "$WORK/plan/native.json" > "$WORK/native.tmp"
    mv "$WORK/native.tmp" "$WORK/plan/native.json"
    jq -n --arg source "$EXPECTED_SHA" --arg version "$version" --arg arch "$arch" --arg repo "$GH_REPO" \
      --arg id "$artifact_id" --arg digest "$digest" \
      '{source:$source,version:$version,package_version:($version|split("+")[0]|split("-")[0]),architecture:$arch,repository:$repo,bundle_id:"io.github.dunky13.takedown",run_id:"12345",run_attempt:"2",artifact_id:$id,artifact_digest:$digest}' > "$evidence/input.json"
    for stage in app dmg; do
      extension=zip; [[ $stage != dmg ]] || extension=dmg
      filename="takedown-$version-macos-$arch-${EXPECTED_SHA:0:12}-12345-2-$stage.$extension"
      submitted=$(printf 'SYNTHETIC submission %s' "$filename" | shasum -a 256 | awk '{print $1}')
      jq -n --arg source "$EXPECTED_SHA" --arg version "$version" --arg arch "$arch" --arg stage "$stage" \
        --arg filename "$filename" --arg submitted "$submitted" \
        '{source_commit:$source,version:$version,architecture:$arch,stage:$stage,ci_run_id:"12345",ci_run_attempt:"2",state:"verified",filename:$filename,submitted_sha256:$submitted,submission_id:"01234567-89ab-cdef-0123-456789abcdef"}' > "$evidence/$stage/receipt.json"
      jq '{status:"Accepted",jobId:.submission_id,archiveFilename:.filename,sha256:.submitted_sha256}' "$evidence/$stage/receipt.json" > "$evidence/$stage/apple-log.json"
    done
    jq -n --arg source "$EXPECTED_SHA" --arg version "$version" --arg arch "$arch" \
      --arg zip "$(hash "$WORK/artifacts/macos-$arch-release.zip/takedown-$version-macos-$arch-release.zip")" \
      --arg dmg "$(hash "$WORK/artifacts/macos-$arch-release.dmg/takedown-$version-macos-$arch-release.dmg")" \
      '{source_commit:$source,version:$version,architecture:$arch,ci_run_id:"12345",ci_run_attempt:"2",zip_sha256:$zip,dmg_sha256:$dmg}' > "$evidence/acceptance.json"
    jq '. + {credential_free_runner:true,verification_run_id:"23456",verification_run_attempt:"1"}' "$evidence/acceptance.json" > "$evidence/fresh-acceptance.json"
  done
}
publish() { bash scripts/publish-release.sh "$WORK/plan" "$WORK/artifacts" > "$WORK/log" 2>&1; }
reject() { if publish; then cat "$WORK/log"; echo 'expected publication rejection' >&2; exit 1; fi; }
no_writes() { if grep -Eq '^api --method|^release (create|upload|edit)' "$MOCK_STATE/calls" 2>/dev/null; then echo 'unexpected remote write' >&2; exit 1; fi; }
no_publish() { if grep -q '^release edit' "$MOCK_STATE/calls"; then echo 'unexpected publication' >&2; exit 1; fi; }
clear_calls() { : > "$MOCK_STATE/calls"; }
mutate() { jq "$2" "$1" > "$WORK/mutated.json"; mv "$WORK/mutated.json" "$1"; }
partial_draft() {
  export MOCK_FAIL_UPLOAD=true
  reject
  unset MOCK_FAIL_UPLOAD
  [[ $(jq -r .draft "$MOCK_STATE/release") == true && $(jq length "$MOCK_STATE/assets") == 1 ]]
  no_publish
  clear_calls
}

reset_fixture manual
publish
grep -Fq -- '--latest=false' "$MOCK_STATE/calls"
grep -Fq -- "--target $EXPECTED_SHA" "$MOCK_STATE/calls"
[[ $(jq -r .draft "$MOCK_STATE/release") == false && $(jq length "$MOCK_STATE/assets") == 8 ]]
[[ $(grep -c '^api --allow-escape-sequences -H' "$MOCK_STATE/calls") == 8 ]]
# An identical published release is a read-only success, including byte downloads.
clear_calls
publish
no_writes
[[ $(grep -c '^api --allow-escape-sequences -H' "$MOCK_STATE/calls") == 8 ]]
grep -q 'Already published' "$WORK/log"
# A published release with changed bytes must fail without mutation.
printf corruption >> "$MOCK_STATE/bytes/100"
clear_calls
reject
no_writes
grep -q 'remote asset bytes differ' "$WORK/log"

reset_fixture nightly
publish
grep -q 'unsigned development' "$WORK/plan/notes.md" && exit 1 # source notes remain untouched
grep -q 'macOS packages are Developer ID signed and notarized' "$MOCK_STATE/notes"
if grep -q 'unsigned development' "$MOCK_STATE/notes"; then exit 1; fi

reset_fixture manual '1.2.0+build-meta'
publish
grep -q -- '--latest=true' "$MOCK_STATE/calls"
reset_fixture auto '1.2.1'
publish
grep -q -- '--latest=true' "$MOCK_STATE/calls"

reset_fixture manual
file="$WORK/artifacts/${names[0]}/takedown-$version-${names[0]}"
printf corrupted >> "$file"
reject
no_writes

reset_fixture manual
rm -rf "$WORK/artifacts/${names[0]}"
reject
no_writes

reset_fixture manual
cp -R "$WORK/artifacts/${names[0]}" "$WORK/artifacts/duplicate"
reject
no_writes

reset_fixture manual
printf '{"object":{"type":"commit","sha":"0000000000000000000000000000000000000000"}}\n' > "$MOCK_STATE/ref"
reject
no_writes

# A partial upload retains its original byte, then only missing names are appended.
reset_fixture manual
partial_draft
before=$(hash "$MOCK_STATE/bytes/100")
existing=$(jq -r '.[0].name' "$MOCK_STATE/assets")
GITHUB_RUN_ID=98765 GITHUB_RUN_ATTEMPT=9 publish
grep -q '^release upload' "$MOCK_STATE/calls"
[[ $(jq -r .draft "$MOCK_STATE/release") == false && $(jq length "$MOCK_STATE/assets") == 8 ]]
[[ $(hash "$MOCK_STATE/bytes/100") == "$before" && $(grep -c '^release upload' "$MOCK_STATE/calls") == 7 ]]
if grep '^release upload' "$MOCK_STATE/calls" | grep -Fq "/$existing"; then exit 1; fi
[[ $(grep -c '^api --allow-escape-sequences -H' "$MOCK_STATE/calls") == 9 ]] # existing byte, then all eight

reset_fixture manual
partial_draft
printf collision >> "$MOCK_STATE/bytes/100"
before=$(hash "$MOCK_STATE/bytes/100")
reject
no_writes
[[ $(hash "$MOCK_STATE/bytes/100") == "$before" ]]

reset_fixture manual
export MOCK_EXTRA_ASSET=true
reject
no_publish
unset MOCK_EXTRA_ASSET

reset_fixture manual
partial_draft
export MOCK_EXTRA_ASSET=true
reject
no_writes
unset MOCK_EXTRA_ASSET

# Successful API uploads still need exact complete downloaded bytes.
for behavior in MOCK_DROP_UPLOAD MOCK_CORRUPT_UPLOAD MOCK_UPLOAD_COLLISION; do
  reset_fixture manual
  export "$behavior=true"
  reject
  no_publish
  [[ $(jq -r .draft "$MOCK_STATE/release") == true ]]
  unset "$behavior"
done

# Atomic tag collisions are terminal; the winning ref and bytes are preserved.
reset_fixture manual
export MOCK_TAG_COLLISION=true
reject
[[ ! -f $MOCK_STATE/release && $(jq -r .object.sha "$MOCK_STATE/ref") == 0000000000000000000000000000000000000000 ]]
unset MOCK_TAG_COLLISION

reset_fixture manual
export MOCK_API_ERROR=true
reject
no_writes
unset MOCK_API_ERROR

# Existing releases must bind the exact target and channel metadata.
for mutation in '.target_commitish="main"' '.prerelease=false'; do
  reset_fixture manual
  partial_draft
  mutate "$MOCK_STATE/release" "$mutation"
  reject
  no_writes
done

reset_fixture manual
partial_draft
export MOCK_DUPLICATE_RELEASE=true
reject
no_writes
unset MOCK_DUPLICATE_RELEASE

for mutation in '. + [.[0]]' '.[0].state="starter"'; do
  reset_fixture manual
  partial_draft
  mutate "$MOCK_STATE/assets" "$mutation"
  reject
  no_writes
done

reset_fixture manual
mutate "$WORK/plan/plan.json" '.tag="v9.9.9"'
reject
no_writes

reset_fixture manual
publish
mutate "$MOCK_STATE/release" '.target_commitish="main"'
clear_calls
reject
no_writes
reset_fixture manual
publish
rm "$MOCK_STATE/ref"
clear_calls
reject
no_writes
reset_fixture manual
publish
mutate "$MOCK_STATE/assets" '. + [{id:999,name:"unexpected.exe",state:"uploaded"}]'
clear_calls
reject
no_writes
reset_fixture manual
publish
mutate "$MOCK_STATE/assets" '.[1:]'
clear_calls
reject
no_writes
reset_fixture manual
publish
file="$WORK/artifacts/${names[0]}/takedown-$version-${names[0]}"
printf 'different candidate bytes\n' >> "$file"
(cd "${file%/*}" && shasum -a 256 "${file##*/}" > SHA256SUMS)
clear_calls
reject
no_writes
grep -q 'remote asset bytes differ' "$WORK/log"

for spec in \
  'candidate.json|.run_id="invalid"' \
  'candidate.json|.run_attempt="3"' \
  'native.json|.[0].id=0' \
  'native.json|.[0].digest="sha256:bad"' \
  'native.json|.[0].workflow_run.head_sha="bad"' \
  'native.json|.[0].workflow_run.id=999'; do
  reset_fixture manual
  mutate "$WORK/plan/${spec%%|*}" "${spec#*|}"
  reject
  no_writes
done

# Both architectures, both stages, original build/run metadata and final bytes
# must bind. These receipts and logs remain synthetic even in the success cases.
for arch in arm64 amd64; do
  for spec in \
    'input.json|.source="bad"' \
    'input.json|.version="9.9.9"' \
    'input.json|.architecture="other"' \
    'input.json|.repository="other/repo"' \
    'input.json|.run_id="not-numeric"' \
    'input.json|.run_attempt="3"' \
    'input.json|.artifact_id="999"' \
    'input.json|.artifact_digest="sha256:bad"' \
    'app/receipt.json|.source_commit="bad"' \
    'app/receipt.json|.version="9.9.9"' \
    'app/receipt.json|.architecture="other"' \
    'app/receipt.json|.state="accepted"' \
    'app/receipt.json|.submission_id=null' \
    'app/receipt.json|.ci_run_id="999"' \
    'app/receipt.json|.ci_run_attempt="3"' \
    'app/apple-log.json|.status="Rejected"' \
    'app/apple-log.json|.jobId="bad"' \
    'app/apple-log.json|.archiveFilename="bad"' \
    'app/apple-log.json|.sha256="bad"' \
    'dmg/receipt.json|.state="pending"' \
    'dmg/apple-log.json|.jobId="bad"' \
    'acceptance.json|.source_commit="bad"' \
    'acceptance.json|.ci_run_id="999"' \
    'acceptance.json|.ci_run_attempt="3"' \
    'acceptance.json|.version="9.9.9"' \
    'acceptance.json|.architecture="other"' \
    'acceptance.json|.zip_sha256="bad"' \
    'acceptance.json|.dmg_sha256="bad"' \
    'fresh-acceptance.json|.credential_free_runner=false' \
    'fresh-acceptance.json|.source_commit="bad"' \
    'fresh-acceptance.json|.verification_run_id="invalid"' \
    'fresh-acceptance.json|.zip_sha256="bad"'; do
    reset_fixture manual
    mutate "$WORK/artifacts/macos-evidence-$arch/${spec%%|*}" "${spec#*|}"
    reject
    no_writes
  done
done
reset_fixture nightly
rm "$WORK/artifacts/macos-evidence-arm64/acceptance.json"
reject
no_writes
reset_fixture nightly
mv "$WORK/artifacts/macos-arm64-release.zip/takedown-$version-macos-arm64-release.zip" "$WORK/artifacts/macos-arm64-release.zip/takedown-$version-macos-arm64-local.zip"
(cd "$WORK/artifacts/macos-arm64-release.zip" && shasum -a 256 "takedown-$version-macos-arm64-local.zip" > SHA256SUMS)
reject
no_writes

# The marker is the newest published nightly, across API pages. A newer failed
# draft, stable release, or custom prerelease must not suppress the next nightly.
cat > "$WORK/releases.json" <<'JSON'
[{"tag_name":"v1.2.1-nightly.20261025022300.0123456789ab","draft":false,"prerelease":true,"published_at":"2026-10-25T03:00:00Z"}]
[{"tag_name":"v1.2.1-nightly.20261026022300.0123456789ab","draft":false,"prerelease":true,"published_at":"2026-10-26T03:00:00Z"},
 {"tag_name":"v1.2.1-nightly.20261027022300.0123456789ab","draft":true,"prerelease":true,"published_at":null},
 {"tag_name":"v1.3.0-rc.1","draft":false,"prerelease":true,"published_at":"2026-10-28T03:00:00Z"},
 {"tag_name":"v1.3.0","draft":false,"prerelease":false,"published_at":"2026-10-29T03:00:00Z"}]
JSON
[[ $(bash scripts/select-last-nightly.sh "$WORK/releases.json") == v1.2.1-nightly.20261026022300.0123456789ab ]]
printf '[]\n' > "$WORK/releases.json"
[[ -z $(bash scripts/select-last-nightly.sh "$WORK/releases.json") ]]
printf 'invalid JSON\n' > "$WORK/releases.json"
if bash scripts/select-last-nightly.sh "$WORK/releases.json" >/dev/null 2>&1; then exit 1; fi
printf 'PASS: synthetic publication integrity, append-only recovery, remote bytes, collisions, immutable idempotence, bound Apple evidence, signed nightlies, and nightly marker selection\n'

reset_fixture manual
export MOCK_AUTOMATION_SHA=bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb
export RELEASE_AUTOMATION_SHA=$MOCK_AUTOMATION_SHA RELEASE_WORKFLOW_SHA=$MOCK_AUTOMATION_SHA
export MOCK_AUTOMATION_FIXTURE="$PWD/scripts/testdata/apple-custody-gh.sh" MOCK_REAL_GIT
MOCK_REAL_GIT=$(command -v git)
cat > "$WORK/mock-bin/git" <<'MOCK'
#!/usr/bin/env bash
case "$1 $2" in
  'rev-parse HEAD') echo "$MOCK_AUTOMATION_SHA" ;;
  'status --porcelain') exit 0 ;;
  'merge-base --is-ancestor') exit 0 ;;
  *) exec "$MOCK_REAL_GIT" "$@" ;;
esac
MOCK
chmod +x "$WORK/mock-bin/git"
publish
jq -e --arg sha "$EXPECTED_SHA" '.target_commitish==$sha and .draft==false' "$MOCK_STATE/release" >/dev/null
if MOCK_CHECK=failure publish; then echo 'publisher accepted unqualified automation' >&2; exit 1; fi
echo 'PASS: publisher retains original product source under separately qualified automation'
