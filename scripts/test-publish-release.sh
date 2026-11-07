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
if [[ ${1:-} == api && ${2:-} == --hostname ]]; then exec "$FIXTURE/bin/gh" "$@"; fi
# Fixture also enforces the no-overwrite contract on every command.
for arg in "$@"; do [[ $arg != --clobber* ]] || { echo 'overwrite forbidden' >&2; exit 1; }; done
if [[ $1 == api ]]; then
  if [[ $2 == --method ]]; then
    if [[ $3 == POST && $4 == */git/refs ]]; then
      [[ ! -f $MOCK_STATE/ref ]] || exit 1
      if [[ ${MOCK_TAG_COLLISION:-false} == true ]]; then
        printf '{"object":{"type":"commit","sha":"0000000000000000000000000000000000000000"}}\n' > "$MOCK_STATE/ref"
        echo 'gh: Reference already exists (HTTP 422)' >&2
        exit 1
      fi
      [[ $5 == -f && $6 == ref=refs/tags/* && $7 == -f && $8 == "sha=$EXPECTED_SHA" ]] || exit 1
      printf '{"object":{"type":"commit","sha":"%s"}}\n' "$EXPECTED_SHA" > "$MOCK_STATE/ref"
      cat "$MOCK_STATE/ref"
    elif [[ $3 == POST && $4 == "repos/$GH_REPO/releases" ]]; then
      [[ ! -f $MOCK_STATE/release && -f $MOCK_STATE/ref && $5 == --input ]] || exit 1
      jq -e --arg sha "$EXPECTED_SHA" '.target_commitish == $sha and .draft == true and .make_latest == "false"' "$6" >/dev/null
      jq '. + {id:1}' "$6" > "$MOCK_STATE/release"
      jq -j .body "$6" > "$MOCK_STATE/notes"
      if [[ ${MOCK_LOST_CREATE:-false} == true ]]; then exit 1; fi
      if [[ -n ${MOCK_CREATE_RECEIPT_MUTATION:-} ]]; then jq "$MOCK_CREATE_RECEIPT_MUTATION" "$MOCK_STATE/release"; else cat "$MOCK_STATE/release"; fi
    elif [[ $3 == POST && $4 == "https://uploads.github.com/repos/$GH_REPO/releases/1/assets?name="* ]]; then
      [[ $(jq -r .draft "$MOCK_STATE/release") == true && $5 == -H && $6 == 'Content-Type: application/octet-stream' && $7 == --input ]] || exit 1
      item=$8; name=${item##*/}
      encoded=$(jq -nr --arg name "$name" '$name | @uri')
      [[ $4 == *"?name=$encoded" ]] || exit 1
      if jq -e --arg name "$name" 'any(.[]; .name == $name)' "$MOCK_STATE/assets" >/dev/null; then
        echo 'gh: asset already exists (HTTP 422)' >&2; exit 1
      fi
      [[ ${MOCK_DROP_UPLOAD:-false} != true ]] || exit 0
      id=$((100 + $(jq length "$MOCK_STATE/assets")))
      cp "$item" "$MOCK_STATE/bytes/$id"
      jq -n --arg name "$name" --argjson id "$id" --argjson size "$(wc -c < "$item")" \
        '{id:$id,name:$name,state:"uploaded",size:$size}' > "$MOCK_STATE/upload"
      jq --slurpfile asset "$MOCK_STATE/upload" '. + $asset' "$MOCK_STATE/assets" > "$MOCK_STATE/assets.tmp"
      mv "$MOCK_STATE/assets.tmp" "$MOCK_STATE/assets"
      if [[ ${MOCK_CORRUPT_UPLOAD:-false} == true || ${MOCK_UPLOAD_COLLISION:-false} == true ]]; then
        printf 'corrupt remote bytes\n' > "$MOCK_STATE/bytes/$id"
      fi
      [[ ${MOCK_FAIL_UPLOAD:-false} != true && ${MOCK_UPLOAD_COLLISION:-false} != true ]] || exit 1
      if [[ -n ${MOCK_UPLOAD_RECEIPT_MUTATION:-} ]]; then jq "$MOCK_UPLOAD_RECEIPT_MUTATION" "$MOCK_STATE/upload"; else cat "$MOCK_STATE/upload"; fi
    elif [[ $3 == PATCH && $4 == "repos/$GH_REPO/releases/1" ]]; then
      [[ $(jq -r .draft "$MOCK_STATE/release") == true && $5 == --input ]] || exit 1
      jq -e '.draft == false' "$6" >/dev/null
      jq --slurpfile update "$6" '. + $update[0]' "$MOCK_STATE/release" > "$MOCK_STATE/release.tmp"
      mv "$MOCK_STATE/release.tmp" "$MOCK_STATE/release"
      jq -j .body "$6" > "$MOCK_STATE/notes"
      [[ ${MOCK_LOST_PUBLISH:-false} != true ]] || exit 1
      if [[ -n ${MOCK_PUBLISH_RECEIPT_MUTATION:-} ]]; then jq "$MOCK_PUBLISH_RECEIPT_MUTATION" "$MOCK_STATE/release"; else cat "$MOCK_STATE/release"; fi
    else
      echo 'unexpected API mutation' >&2; exit 1
    fi
  elif [[ $2 == --paginate && $3 == */releases\?* ]]; then
    if [[ ${MOCK_API_ERROR:-false} == true ]]; then echo 'gh: Forbidden (HTTP 403)' >&2; exit 1; fi
    # Two pages exercise the authenticated draft discovery path.
    echo '[]'
    if [[ -f $MOCK_STATE/release && ${MOCK_STALE_RELEASE_LIST:-false} != true ]]; then
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
  elif [[ $2 == "repos/$GH_REPO/releases/1" && -f $MOCK_STATE/release ]]; then
    if [[ -n ${MOCK_EXACT_RELEASE_MUTATION:-} ]]; then jq "$MOCK_EXACT_RELEASE_MUTATION" "$MOCK_STATE/release"; else cat "$MOCK_STATE/release"; fi
  elif [[ ${MOCK_AUTOMATION_SHA:-} != '' && ( $2 == */branches/main || $2 == */compare/* || $2 == --paginate ) ]]; then
    exec bash "$MOCK_AUTOMATION_FIXTURE" "$@"
  else
    printf 'gh: Not Found (HTTP 404)\n' >&2
    exit 1
  fi
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
    jq --arg input "$(hash "$evidence/input.json")" --arg app "$(hash "$evidence/app/receipt.json")" --arg dmg "$(hash "$evidence/dmg/receipt.json")" \
      '. + {schema:1,state:"verified",attempts:0,failure_class:"",team_id:"ABCDEFGHIJ",input_sha256:$input,app_receipt_sha256:$app,dmg_receipt_sha256:$dmg}' \
      "$evidence/acceptance.json" > "$evidence/acceptance-state.json"
  done
}
publish() { bash scripts/publish-release.sh "$WORK/plan" "$WORK/artifacts" > "$WORK/log" 2>&1; }
reject() { if publish; then cat "$WORK/log"; echo 'expected publication rejection' >&2; exit 1; fi; }
no_writes() { if grep -Eq '^api --method|^release (create|upload|edit)' "$MOCK_STATE/calls" 2>/dev/null; then echo 'unexpected remote write' >&2; exit 1; fi; }
no_publish() { if grep -q '^api --method PATCH' "$MOCK_STATE/calls"; then echo 'unexpected publication' >&2; exit 1; fi; }
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

# A successful write may precede release-list visibility. Keep that list stale
# for the whole invocation: the write receipt and exact ID must be sufficient.
reset_fixture manual
MOCK_STALE_RELEASE_LIST=true publish
[[ $(jq -r .draft "$MOCK_STATE/release") == false && $(jq length "$MOCK_STATE/assets") == 8 ]]
[[ $(grep -c '^api --paginate .*releases?per_page=' "$MOCK_STATE/calls") == 1 ]]
echo 'PASS: publication does not rediscover a newly created draft through a stale list'

reset_fixture nightly
MOCK_STALE_RELEASE_LIST=true publish
jq -e '.prerelease == true and .make_latest == "false" and .draft == false' "$MOCK_STATE/release" >/dev/null
[[ $(jq length "$MOCK_STATE/assets") == 8 ]]
echo 'PASS: signed nightly publication uses the same ID-bound path and preserves prerelease/Latest policy'

# Creation and upload failures leave the original draft intact for reconciliation;
# the same invocation never repeats an uncertain write.
for operation in CREATE PUBLISH; do
  reset_fixture manual
  export "MOCK_LOST_$operation=true"
  reject
  unset "MOCK_LOST_$operation"
  [[ -f $MOCK_STATE/release ]]
  clear_calls
  publish
  [[ $(jq -r .draft "$MOCK_STATE/release") == false && $(jq length "$MOCK_STATE/assets") == 8 ]]
  if [[ $operation == PUBLISH ]]; then no_writes; fi
  if grep -q "^api --method POST repos/$GH_REPO/releases --input" "$MOCK_STATE/calls"; then exit 1; fi
done
echo 'PASS: lost creation/publication responses reconcile the retained ID without duplicate releases or overwrites'

for mutation in '.id=0' '.tag_name="v9.9.9"' '.target_commitish="main"' '.draft=false' '.prerelease=false'; do
  reset_fixture manual
  MOCK_CREATE_RECEIPT_MUTATION=$mutation reject
  no_publish
  [[ $(jq length "$MOCK_STATE/assets") == 0 ]]
done
for mutation in '.id=2' '.tag_name="v9.9.9"' '.target_commitish="main"' '.draft=false' '.prerelease=false'; do
  reset_fixture manual
  MOCK_EXACT_RELEASE_MUTATION=$mutation reject
  no_publish
  [[ $(jq length "$MOCK_STATE/assets") == 0 ]]
done
for mutation in '.id=0' '.name="unexpected.exe"' '.size+=1' '.state="starter"'; do
  reset_fixture manual
  MOCK_UPLOAD_RECEIPT_MUTATION=$mutation reject
  no_publish
  [[ $(jq length "$MOCK_STATE/assets") == 1 ]]
done
for mutation in '.id=2' '.target_commitish="main"' '.draft=true'; do
  reset_fixture manual
  MOCK_PUBLISH_RECEIPT_MUTATION=$mutation reject
  [[ $(jq -r .draft "$MOCK_STATE/release") == false ]]
  clear_calls
  publish
  no_writes
done
echo 'PASS: invalid write receipts and changed exact-ID objects fail closed; publication retries reconcile read-only'
if [[ ${1:-} == --publication-recovery-only ]]; then exit 0; fi

reset_fixture manual
publish
jq -e --arg sha "$EXPECTED_SHA" '.make_latest == "false" and .target_commitish == $sha' "$MOCK_STATE/release" >/dev/null
[[ $(jq -r .draft "$MOCK_STATE/release") == false && $(jq length "$MOCK_STATE/assets") == 8 ]]
[[ $(grep -c '^api --allow-escape-sequences -H' "$MOCK_STATE/calls") == 16 ]]
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
jq -e '.make_latest == "true"' "$MOCK_STATE/release" >/dev/null
reset_fixture auto '1.2.1'
publish
jq -e '.make_latest == "true"' "$MOCK_STATE/release" >/dev/null

# Pre-policy candidates remain publishable only through every original
# signer/Accepted-log/native/fresh-acceptance binding and remote-byte gate.
reset_fixture manual
for arch in arm64 amd64; do
  rm "$WORK/artifacts/macos-evidence-$arch/acceptance-state.json"
done
publish
[[ $(jq -r .draft "$MOCK_STATE/release") == false && $(jq length "$MOCK_STATE/assets") == 8 ]]
echo 'PASS: legacy evidence without final acceptance state retains all original publication gates'

# Exercise the actual publisher with a complete exception receipt and immutable
# parent. Pin substitutions exist only in a private policy copy for synthetic
# candidate/run/repository bytes; the shipped policy cannot authorize fixtures.
mkdir "$WORK/exception-scripts"
cp scripts/publish-release.sh scripts/release-source.sh "$WORK/exception-scripts/"
exception_fixture() {
  reset_fixture manual 0.0.1
  evidence="$WORK/artifacts/macos-evidence-arm64"
  mkdir "$evidence/final"
  for stage in app dmg; do
    extension=zip; [[ $stage != dmg ]] || extension=dmg
    file="$WORK/artifacts/macos-arm64-release.$extension/takedown-$version-macos-arm64-release.$extension"
    cp "$file" "$evidence/final/"
    verified_file=stapled-app.zip; [[ $stage != dmg ]] || verified_file=stapled.dmg
    cp "$file" "$evidence/$stage/$verified_file"
    mutate "$evidence/$stage/receipt.json" ". + {schema:2,team_id:\"ABCDEFGHIJ\",signing_identity:\"Developer ID Application: Marc Went (ABCDEFGHIJ)\",verified_filename:\"$verified_file\",verified_sha256:\"$(hash "$file")\",ticket_attempts:1,validation_attempts:0}"
  done
  (cd "$evidence/final" && shasum -a 256 ./*.zip ./*.dmg | sed 's@  \./@  @' > SHA256SUMS)
  rm -rf "$WORK/artifacts/macos-arm64-release.zip" "$WORK/artifacts/macos-arm64-release.dmg"
  mutate "$evidence/acceptance-state.json" ".state=\"blocked\" | .attempts=3 | .failure_class=\"cloudkit-transport\" | .app_receipt_sha256=\"$(hash "$evidence/app/receipt.json")\" | .dmg_receipt_sha256=\"$(hash "$evidence/dmg/receipt.json")\""
  exception_pin=$(hash "$evidence/acceptance-state.json")
  cp "$evidence/acceptance-state.json" "$evidence/blocked-acceptance-$exception_pin.json"
  sed -e "s/27b7a083cc181ed6def7479ec99b25479d7b11d9c3de36bbbd9facfbacfc60fe/$exception_pin/g" \
    -e "s/ef82381d509aac8b5b02da0f784a8176b1c5b9c0/$EXPECTED_SHA/g" \
    -e 's/40470756538/12345/g' -e 's/\.ci_run_attempt == "1"/.ci_run_attempt == "2"/g' \
    -e 's@Dunky13/takedown@fixture/takedown@g' scripts/macos-ticket-policy.sh > "$WORK/exception-scripts/macos-ticket-policy.sh"
  mutate "$evidence/acceptance-state.json" ".state=\"verified\" | .attempts=4 | .failure_class=\"\" | .acceptance_exception_parent_sha256=\"$exception_pin\""
  for receipt in acceptance.json fresh-acceptance.json; do mutate "$evidence/$receipt" ".acceptance_attempts=4 | .acceptance_exception_parent_sha256=\"$exception_pin\" | .credential_free_runner=true | .verification_run_id=\"23456\" | .verification_run_attempt=\"1\" | .verification_workflow_sha=\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\""; done
  cp "$evidence/acceptance.json" "$evidence/fresh-acceptance.json"
  sed -e "s/27b7a083cc181ed6def7479ec99b25479d7b11d9c3de36bbbd9facfbacfc60fe/$exception_pin/g" \
    -e "s/ef82381d509aac8b5b02da0f784a8176b1c5b9c0/$EXPECTED_SHA/g" \
    -e 's/40470756538/12345/g' -e 's/"ci_run_attempt": "1"/"ci_run_attempt": "2"/g' \
    -e 's@Dunky13/takedown@fixture/takedown@g' scripts/verify-apple-fresh-completion.sh > "$WORK/exception-scripts/verify-apple-fresh-completion.sh"
  bash scripts/test-apple-fresh-completion.sh --emit-fixture "$evidence/acceptance.json" "$MOCK_STATE/completion-provider" fixture/takedown
  export FIXTURE="$MOCK_STATE/completion-provider" NOTARY_RELEASE_ID=425202908 NOTARY_RELEASE_SOURCE_SHA="$EXPECTED_SHA" NOTARY_STATE_TAG=apple-candidate-v0.0.1
  NOTARY_RELEASE_BODY_SHA256=$(jq -r .NOTARY_RELEASE_BODY_SHA256 "$FIXTURE/pins.json")
  export NOTARY_RELEASE_BODY_SHA256
}
exception_publish() { bash "$WORK/exception-scripts/publish-release.sh" "$WORK/plan" "$WORK/artifacts" > "$WORK/log" 2>&1; }
exception_fixture
exception_publish
[[ $(jq -r .draft "$MOCK_STATE/release") == false ]]
for exception_case in missing-parent changed-parent missing-state-and-parent fresh-only-marker stale-signer stale-fresh missing-proof foreign-verification pending blocked counter foreign-parent; do
  exception_fixture
  case $exception_case in
    missing-parent) rm "$evidence/blocked-acceptance-$exception_pin.json" ;;
    changed-parent) echo ' ' >> "$evidence/blocked-acceptance-$exception_pin.json" ;;
    missing-state-and-parent|fresh-only-marker)
      rm "$evidence/acceptance-state.json" "$evidence/blocked-acceptance-$exception_pin.json"
      if [[ $exception_case == fresh-only-marker ]]; then mutate "$evidence/acceptance.json" 'del(.acceptance_attempts,.acceptance_exception_parent_sha256)'; fi ;;
    stale-signer) mutate "$evidence/acceptance.json" 'del(.acceptance_attempts,.acceptance_exception_parent_sha256)' ;;
    stale-fresh) mutate "$evidence/fresh-acceptance.json" 'del(.acceptance_attempts,.acceptance_exception_parent_sha256)' ;;
    missing-proof) mutate "$evidence/acceptance.json" 'del(.verification_run_id)' ;;
    foreign-verification) mutate "$evidence/fresh-acceptance.json" '.verification_run_id="999"' ;;
    pending) mutate "$evidence/acceptance-state.json" '.state="pending"' ;;
    blocked) mutate "$evidence/acceptance-state.json" '.state="blocked" | .failure_class="cloudkit-transport"' ;;
    counter) mutate "$evidence/acceptance-state.json" '.attempts=3' ;;
    foreign-parent) mutate "$evidence/acceptance-state.json" '.acceptance_exception_parent_sha256="aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"' ;;
  esac
  if exception_publish; then echo "publisher accepted invalid exception $exception_case" >&2; exit 1; fi
  no_writes
done
echo 'PASS: publisher accepts bound verified4 only with immutable eligible parent and new native/fresh receipts; stale/pending/blocked/foreign/missing/changed evidence cannot promote'

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
grep -q '^api --method POST https://uploads.github.com/' "$MOCK_STATE/calls"
[[ $(jq -r .draft "$MOCK_STATE/release") == false && $(jq length "$MOCK_STATE/assets") == 8 ]]
[[ $(hash "$MOCK_STATE/bytes/100") == "$before" && $(grep -c '^api --method POST https://uploads.github.com/' "$MOCK_STATE/calls") == 7 ]]
if grep '^api --method POST https://uploads.github.com/' "$MOCK_STATE/calls" | grep -Fq "/$existing"; then exit 1; fi
[[ $(grep -c '^api --allow-escape-sequences -H' "$MOCK_STATE/calls") == 16 ]] # existing byte, seven upload receipts, then all eight

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
    'app/receipt.json|.state="validation-pending"' \
    'app/receipt.json|.state="ticket-blocked"' \
    'app/receipt.json|.requalification_parent_sha256="aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"' \
    'app/receipt.json|.validation_attempts=4' \
    'app/receipt.json|.failure_phase="stapler-validate"' \
    'app/receipt.json|.failure_exit="68"' \
    'app/receipt.json|.failure_class="unclassified"' \
    'app/receipt.json|.submission_id=null' \
    'app/receipt.json|.ci_run_id="999"' \
    'app/receipt.json|.ci_run_attempt="3"' \
    'app/apple-log.json|.status="Rejected"' \
    'app/apple-log.json|.jobId="bad"' \
    'app/apple-log.json|.archiveFilename="bad"' \
    'app/apple-log.json|.sha256="bad"' \
    'dmg/receipt.json|.state="pending"' \
    'dmg/receipt.json|.state="validation-pending"' \
    'dmg/receipt.json|.state="ticket-blocked"' \
    'dmg/receipt.json|.failure_phase="stapler-validate"' \
    'dmg/receipt.json|.failure_exit="68"' \
    'dmg/receipt.json|.failure_class="unclassified"' \
    'dmg/apple-log.json|.jobId="bad"' \
    'acceptance-state.json|.state="pending"' \
    'acceptance-state.json|.state="blocked"|.failure_class="cloudkit-transport"' \
    'acceptance-state.json|.attempts=4' \
    'acceptance-state.json|.attempts=4|.acceptance_exception_parent_sha256="27b7a083cc181ed6def7479ec99b25479d7b11d9c3de36bbbd9facfbacfc60fe"' \
    'acceptance-state.json|.attempts=-1' \
    'acceptance-state.json|.failure_class="unclassified"' \
    'acceptance-state.json|.source_commit="bad"' \
    'acceptance-state.json|.zip_sha256="bad"' \
    'acceptance-state.json|.dmg_sha256="bad"' \
    'acceptance-state.json|.input_sha256="bad"' \
    'acceptance-state.json|.app_receipt_sha256="bad"' \
    'acceptance-state.json|.dmg_receipt_sha256="bad"' \
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
