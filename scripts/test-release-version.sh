#!/usr/bin/env bash
# Exercise real packagers with disposable native-tool fixtures. No GUI builds,
# signing, disk-image mounts, or account access occur here.
set -euo pipefail
unset VERSION RELEASE_VERSION BINARY MODE OUT_DIR
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
cd "$ROOT"
TEST_WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-release-version.XXXXXX")
trap 'rm -rf -- "$TEST_WORK"' EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
export REAL_GO_PATH
REAL_GO_PATH=$(command -v go)
export RELEASE_VERSION_TRACE="$TEST_WORK/trace"
mkdir "$TEST_WORK/tools"
cat > "$TEST_WORK/tools/go" <<'GO'
#!/usr/bin/env bash
set -euo pipefail
case "$1 $2" in
  'run -mod=readonly') exec "$REAL_GO_PATH" "$@" ;;
  'env GOVERSION') echo go1.26.10 ;;
  'env GOHOSTOS') echo "$FIXTURE_GOOS" ;;
  'env GOHOSTARCH') echo amd64 ;;
  'list -m') echo v2.8.2 ;;
  'version -m')
    printf '%s: go1.26.10\n\tdep\tfyne.io/fyne/v2\tv2.8.2\tfixture\n\tbuild\tCGO_ENABLED=1\n\tbuild\tGOOS=%s\n\tbuild\tGOARCH=amd64\n' "$3" "$FIXTURE_GOOS"
    ;;
  'build -mod=readonly')
    version='' output=''
    while (($#)); do
      case $1 in
        -ldflags) version=${2##*=}; shift ;;
        -o) output=$2; shift ;;
      esac
      shift
    done
    [[ -n $version && -n $output ]]
    printf 'build-version=%s\n' "$version" >> "$RELEASE_VERSION_TRACE"
    printf '#!/bin/sh\nprintf "takedown version %%s\\n" "%s"\n' "$version" > "$output"
    chmod +x "$output"
    ;;
  *) printf 'unexpected fixture go invocation\n' >&2; exit 99 ;;
esac
GO
cat > "$TEST_WORK/tools/ditto" <<'DITTO'
#!/usr/bin/env bash
set -euo pipefail
if [[ $1 == -c ]]; then printf 'zip fixture\n' > "${!#}"
else cp -R "$1" "$2"
fi
DITTO
cat > "$TEST_WORK/tools/hdiutil" <<'DMG'
#!/usr/bin/env bash
set -euo pipefail
printf 'dmg fixture\n' > "${!#}"
DMG
cat > "$TEST_WORK/tools/plutil" <<'PLIST'
#!/usr/bin/env bash
set -euo pipefail
printf '%s\n' "$*" >> "$RELEASE_VERSION_TRACE"
PLIST
cat > "$TEST_WORK/tools/xattr" <<'XATTR'
#!/usr/bin/env bash
exit 0
XATTR
cat > "$TEST_WORK/tools/tar" <<'TAR'
#!/usr/bin/env bash
set -euo pipefail
if [[ $1 == --version ]]; then echo 'GNU tar fixture'; else printf 'tar fixture\n'; fi
TAR
cat > "$TEST_WORK/tools/sha256sum" <<'SUM'
#!/usr/bin/env bash
exec shasum -a 256 "$@"
SUM
chmod +x "$TEST_WORK/tools/"*
export PATH="$TEST_WORK/tools:$PATH"
expect_failure() {
  local expected=$1
  shift
  if "$@" > "$TEST_WORK/failure.log" 2>&1; then echo "expected rejection: $expected" >&2; exit 1; fi
  grep -Fq "$expected" "$TEST_WORK/failure.log" || { cat "$TEST_WORK/failure.log"; exit 1; }
}

# Full prerelease reaches the build, filenames and bundled metadata; native
# package fields remain numeric. Real checksum generation covers both assets.
export FIXTURE_GOOS=darwin
VERSION=1.2.3 RELEASE_VERSION=1.2.3-rc.1+build.7 MODE=local OUT_DIR="$TEST_WORK/macos" bash scripts/package-macos.sh
[[ -f $TEST_WORK/macos/takedown-1.2.3-rc.1+build.7-macos-amd64-local.zip ]]
[[ -f $TEST_WORK/macos/takedown-1.2.3-rc.1+build.7-macos-amd64-local.dmg ]]
grep -Fxq 'build-version=1.2.3-rc.1+build.7' "$RELEASE_VERSION_TRACE"
grep -Fq -- '-replace CFBundleShortVersionString -string 1.2.3 ' "$RELEASE_VERSION_TRACE"
grep -Fq -- '-replace CFBundleVersion -string 1.2.3 ' "$RELEASE_VERSION_TRACE"
grep -Fxq 'Release version: 1.2.3-rc.1+build.7' "$TEST_WORK/macos/takedown.app/Contents/Resources/docs/BUILD-INFO.txt"
(cd "$TEST_WORK/macos" && shasum -a 256 -c SHA256SUMS)

export FIXTURE_GOOS=linux
VERSION=1.2.3 RELEASE_VERSION=1.2.3-nightly.20261026030000.abcdef0 MODE=local OUT_DIR="$TEST_WORK/linux" bash scripts/package-linux.sh
[[ -f $TEST_WORK/linux/takedown-1.2.3-nightly.20261026030000.abcdef0-linux-amd64-local.tar.gz ]]
(cd "$TEST_WORK/linux" && shasum -a 256 -c SHA256SUMS)

export FIXTURE_GOOS=darwin
for invalid in '1.2.3-rc.01' '1.2.3-rc/1' '1.2.3;touch sentinel' ''; do
  expect_failure 'RELEASE_VERSION must be a supported SemVer version' env VERSION=1.2.3 RELEASE_VERSION="$invalid" MODE=local OUT_DIR="$TEST_WORK/invalid" bash scripts/package-macos.sh
  [[ ! -e $TEST_WORK/invalid ]]
done
expect_failure 'RELEASE_VERSION numeric core must match VERSION' env VERSION=1.2.3 RELEASE_VERSION=1.2.4-rc.1 MODE=local OUT_DIR="$TEST_WORK/mismatch" bash scripts/package-macos.sh
[[ ! -e $TEST_WORK/mismatch ]]

# Explicit versions reject incorrectly labelled reused binaries. Omission keeps
# the existing local BINARY contract, whose metadata is checked separately.
printf '#!/bin/sh\necho "takedown version 1.2.3-rc.2"\n' > "$TEST_WORK/reused"
chmod +x "$TEST_WORK/reused"
expect_failure 'binary version does not match RELEASE_VERSION' env VERSION=1.2.3 RELEASE_VERSION=1.2.3-rc.1 BINARY="$TEST_WORK/reused" MODE=local OUT_DIR="$TEST_WORK/rejected-binary" bash scripts/package-macos.sh
[[ ! -e $TEST_WORK/rejected-binary ]]
env -u RELEASE_VERSION VERSION=1.2.3 BINARY="$TEST_WORK/reused" MODE=local OUT_DIR="$TEST_WORK/legacy" bash scripts/package-macos.sh
[[ -f $TEST_WORK/legacy/takedown-1.2.3-macos-amd64-local.zip ]]
echo 'PASS: full release labels, numeric native metadata, checksums, unsafe/mismatched versions, reused binary version, and legacy BINARY compatibility'
