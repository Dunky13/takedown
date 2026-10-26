#!/usr/bin/env bash
# Offline regressions for R1; all writes are disposable, no app or keyring use.
set -euo pipefail
ROOT=$(cd -P -- "$(dirname -- "$0")/.." && pwd -P)
cd "$ROOT"
WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-r1-tests.XXXXXX")
WORK=$(cd -P "$WORK" && pwd -P)
trap 'rm -rf -- "$WORK"' EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
expect_failure() {
  local expected=$1
  shift
  if "$@" > "$WORK/failure.log" 2>&1; then echo "expected rejection: $expected" >&2; exit 1; fi
  grep -Fq "$expected" "$WORK/failure.log" || { cat "$WORK/failure.log"; exit 1; }
}
FIXTURE="$WORK/payload"
mkdir -p "$FIXTURE/bin" "$FIXTURE/share/applications" "$FIXTURE/share/icons/hicolor/scalable/apps" "$FIXTURE/share/doc/takedown"
printf '#!/bin/sh\nexit 0\n' > "$FIXTURE/bin/takedown"
printf 'fixture docs\n' > "$FIXTURE/share/doc/takedown/README.md"
printf 'fixture application notice\n' > "$FIXTURE/share/doc/takedown/NOTICE"
printf 'fixture dependency notices\n' > "$FIXTURE/share/doc/takedown/THIRD_PARTY_NOTICES.md"
cp assets/icon.svg "$FIXTURE/share/icons/hicolor/scalable/apps/takedown.svg"
cp build/linux/takedown.desktop "$FIXTURE/share/applications/"
cp build/linux/*.sh "$FIXTURE/"
printf 'KEEP THIS UNRELATED FILE\n' > "$WORK/sentinel"
cp "$WORK/sentinel" "$WORK/expected-sentinel"
# Exercise the installed publication primitive against directory/symlink targets.
mkdir "$WORK/link-directory"
if link "$WORK/sentinel" "$WORK/link-directory" 2>/dev/null; then echo 'link merged into an existing directory' >&2; exit 1; fi
ln -s "$WORK/link-directory" "$WORK/link-directory-symlink"
if link "$WORK/sentinel" "$WORK/link-directory-symlink" 2>/dev/null; then echo 'link followed a destination symlink' >&2; exit 1; fi
[[ ! -e $WORK/link-directory/sentinel ]]
# R1-01: reproduce each old loss witness in a NEW fixture, never mutate R1 evidence.
PREFIX="$WORK/symlink-prefix"
mkdir -p "$PREFIX/share/applications"
ln -s "$WORK/sentinel" "$PREFIX/share/applications/takedown.desktop"
expect_failure 'symlink payload path' env INSTALL_PREFIX="$PREFIX" bash "$FIXTURE/install-user.sh"
cmp "$WORK/sentinel" "$WORK/expected-sentinel"
[[ ! -e $PREFIX/bin/takedown ]]
PREFIX="$WORK/ancestor-prefix"
mkdir "$PREFIX" "$WORK/outside"
ln -s "$WORK/outside" "$PREFIX/share"
expect_failure 'symlink in destination ancestry' env INSTALL_PREFIX="$PREFIX" bash "$FIXTURE/install-user.sh"
[[ ! -e $PREFIX/bin && ! -e $WORK/outside/applications ]]
ln -s "$WORK/outside" "$WORK/prefix-link"
expect_failure 'symlink in destination ancestry' env INSTALL_PREFIX="$WORK/prefix-link" bash "$FIXTURE/install-user.sh"
PREFIX="$WORK/conflict-prefix"
mkdir -p "$PREFIX/bin" "$PREFIX/share/doc/takedown"
cp "$WORK/sentinel" "$PREFIX/bin/takedown"
cp "$WORK/sentinel" "$PREFIX/share/doc/takedown/local-notes.txt"
expect_failure 'unowned or existing installation path' env INSTALL_PREFIX="$PREFIX" bash "$FIXTURE/install-user.sh"
expect_failure 'no ownership manifest' env INSTALL_PREFIX="$PREFIX" bash "$FIXTURE/uninstall-user.sh"
cmp "$WORK/sentinel" "$PREFIX/bin/takedown"
cmp "$WORK/sentinel" "$PREFIX/share/doc/takedown/local-notes.txt"
PREFIX="$WORK/owned prefix"
INSTALL_PREFIX="$PREFIX" bash "$FIXTURE/install-user.sh"
expect_failure 'unowned or existing installation path' env INSTALL_PREFIX="$PREFIX" bash "$FIXTURE/install-user.sh"
cp "$WORK/sentinel" "$PREFIX/share/doc/takedown/local-notes.txt"
printf 'modified docs\n' > "$PREFIX/share/doc/takedown/README.md"
expect_failure 'owned file changed' env INSTALL_PREFIX="$PREFIX" bash "$FIXTURE/uninstall-user.sh"
[[ -f $PREFIX/bin/takedown && -f $PREFIX/share/takedown/install-manifest.sha256 ]]
cp "$FIXTURE/share/doc/takedown/README.md" "$PREFIX/share/doc/takedown/README.md"
for name in NOTICE THIRD_PARTY_NOTICES.md; do
  cmp "$FIXTURE/share/doc/takedown/$name" "$PREFIX/share/doc/takedown/$name"
  printf 'modified notice\n' > "$PREFIX/share/doc/takedown/$name"
  expect_failure 'owned file changed' env INSTALL_PREFIX="$PREFIX" bash "$FIXTURE/uninstall-user.sh"
  [[ -f $PREFIX/bin/takedown && -f $PREFIX/share/takedown/install-manifest.sha256 ]]
  cp "$FIXTURE/share/doc/takedown/$name" "$PREFIX/share/doc/takedown/$name"
done
cp "$PREFIX/share/takedown/install-manifest.sha256" "$WORK/good-manifest"
printf '%064d  ../../sentinel\n' 0 >> "$PREFIX/share/takedown/install-manifest.sha256"
expect_failure 'unrecognized ownership path' env INSTALL_PREFIX="$PREFIX" bash "$FIXTURE/uninstall-user.sh"
cmp "$WORK/sentinel" "$WORK/expected-sentinel"
cp "$WORK/good-manifest" "$PREFIX/share/takedown/install-manifest.sha256"
printf '%064d  share/doc/takedown/NOTICE.extra\n' 0 >> "$PREFIX/share/takedown/install-manifest.sha256"
expect_failure 'unrecognized ownership path' env INSTALL_PREFIX="$PREFIX" bash "$FIXTURE/uninstall-user.sh"
[[ -f $PREFIX/bin/takedown ]]
cp "$WORK/good-manifest" "$PREFIX/share/takedown/install-manifest.sha256"
mv "$PREFIX/share/applications/takedown.desktop" "$WORK/owned-desktop"
ln -s "$WORK/sentinel" "$PREFIX/share/applications/takedown.desktop"
expect_failure 'symlink payload path' env INSTALL_PREFIX="$PREFIX" bash "$FIXTURE/uninstall-user.sh"
[[ -f $PREFIX/bin/takedown ]]
rm "$PREFIX/share/applications/takedown.desktop"
mv "$WORK/owned-desktop" "$PREFIX/share/applications/takedown.desktop"
INSTALL_PREFIX="$PREFIX" bash "$FIXTURE/uninstall-user.sh"
[[ ! -e $PREFIX/bin/takedown && ! -e $PREFIX/share/takedown/install-manifest.sha256 ]]
[[ ! -e $PREFIX/share/doc/takedown/NOTICE && ! -e $PREFIX/share/doc/takedown/THIRD_PARTY_NOTICES.md ]]
cmp "$WORK/sentinel" "$PREFIX/share/doc/takedown/local-notes.txt"
for name in NOTICE THIRD_PARTY_NOTICES.md; do
  mv "$FIXTURE/share/doc/takedown/$name" "$WORK/held-notice"
  expect_failure 'missing package payload' env INSTALL_PREFIX="$WORK/missing-notice-prefix" bash "$FIXTURE/install-user.sh"
  [[ ! -e $WORK/missing-notice-prefix ]]
  ln -s "$WORK/held-notice" "$FIXTURE/share/doc/takedown/$name"
  expect_failure 'symlink payload path' env INSTALL_PREFIX="$WORK/missing-notice-prefix" bash "$FIXTURE/install-user.sh"
  [[ ! -e $WORK/missing-notice-prefix ]]
  rm "$FIXTURE/share/doc/takedown/$name"
  mv "$WORK/held-notice" "$FIXTURE/share/doc/takedown/$name"
done
echo 'PASS R1-01: symlink/ancestor/conflict rejection, manifest validation, modified/unknown file preservation'

# R1-02: apply the real workflow path patterns to ZIP + setup + checksum fixtures.
PACKAGES="$WORK/upload/dist/packages"
mkdir -p "$PACKAGES/windows-amd64-local"
printf 'zip fixture\n' > "$PACKAGES/windows-amd64-local/takedown-0.0.0-windows-amd64-local.zip"
printf 'setup fixture\n' > "$PACKAGES/windows-amd64-local/takedown-0.0.0-windows-amd64-local-setup.exe"
(cd "$PACKAGES/windows-amd64-local" && shasum -a 256 ./*.zip ./*.exe | sed 's@  \./@  @' > SHA256SUMS)
bash scripts/check-artifact-retention.sh .github/workflows/ci.yml "$PACKAGES"
bash scripts/check-artifact-retention.sh .github/workflows/signed-packages.yml "$PACKAGES"
sed '/dist\/packages\/\*\*\/\*-setup.exe/d' .github/workflows/ci.yml > "$WORK/old-upload.yml"
expect_failure 'checksum artifact is not retained' bash scripts/check-artifact-retention.sh "$WORK/old-upload.yml" "$PACKAGES"
rm "$PACKAGES/windows-amd64-local/"*-setup.exe
expect_failure 'checksum artifact is missing' bash scripts/check-artifact-retention.sh .github/workflows/ci.yml "$PACKAGES"
echo 'PASS R1-02: retained setup/checksum closure, old-glob and missing-file rejection'

# R1-03: run two real reservation helpers concurrently with an explicit barrier.
cat > "$WORK/publisher.sh" <<'PUBLISH'
#!/usr/bin/env bash
set -euo pipefail
source "$1/scripts/package-common.sh"
OUT_DIR=$2
WORK=$(mktemp -d "$3/stage.XXXXXX")
OUTPUT_RESERVED=false
trap package_cleanup EXIT
if [[ $4 == unreserved ]]; then package_finish; exit; fi
package_reserve_output
mkdir "$WORK/output"
printf '%s\n' "$4" > "$WORK/output/SHA256SUMS"
touch "$3/ready"
for ((attempt=0; attempt<200; attempt++)); do
  [[ ! -e $3/release ]] || break
  sleep 0.02
done
[[ -e $3/release ]] || fail 'test publication barrier timed out'
package_finish
PUBLISH
mkdir "$WORK/race"
bash "$WORK/publisher.sh" "$ROOT" "$WORK/final" "$WORK/race" first > "$WORK/first.log" 2>&1 &
publisher=$!
for ((attempt=0; attempt<200; attempt++)); do [[ ! -e $WORK/race/ready ]] || break; sleep 0.02; done
[[ -e $WORK/race/ready ]] || { cat "$WORK/first.log"; exit 1; }
expect_failure 'OUT_DIR already exists or is reserved' bash "$WORK/publisher.sh" "$ROOT" "$WORK/final" "$WORK/race" second
[[ -f $WORK/final/.package-incomplete ]]
touch "$WORK/race/release"
wait "$publisher"
[[ $(cat "$WORK/final/SHA256SUMS") == first && ! -e $WORK/final/output && ! -e $WORK/final/.package-incomplete ]]
# The old source-only witness may not bypass the new reservation precondition.
expect_failure 'owned reservation' bash "$WORK/publisher.sh" "$ROOT" "$WORK/final" "$WORK/race" unreserved
[[ $(cat "$WORK/final/SHA256SUMS") == first ]]
echo 'PASS R1-03: simultaneous reservation loser fails, winner remains intact, unreserved publication rejected'

# R1-05: source the actual lifecycle helper and shadow security with a file-only
# mock. These tests NEVER invoke the installed security tool or any keyring API.
(
  WORK="$WORK/keychain-mock"
  mkdir "$WORK"
  KEYCHAIN="$WORK/temporary.keychain-db"
  MOCK_LIST="$WORK/list"
  printf '    "/fixture/login.keychain-db"\n    "/fixture/other keychain.keychain-db"\n' > "$MOCK_LIST"
  MOCK_DELETE_FAIL=false
  MOCK_RESTORE_FAIL=false
  MOCK_IGNORE_SET=false
  security() {
    case $1 in
      delete-keychain) [[ $MOCK_DELETE_FAIL == false ]] || return 1; rm -- "$2" ;;
      list-keychains)
        [[ $2 == -d && $3 == user ]] || return 99
        if [[ $# == 3 ]]; then cat "$MOCK_LIST"
        else
          [[ $4 == -s && $MOCK_RESTORE_FAIL == false ]] || return 1
          [[ $MOCK_IGNORE_SET == false ]] || return 0
          shift 4
          : > "$MOCK_LIST"
          for path in "$@"; do printf '    "%s"\n' "$path" >> "$MOCK_LIST"; done
        fi ;;
      *) return 99 ;;
    esac
  }
  # shellcheck source=scripts/keychain-lifecycle.sh
  source "$ROOT/scripts/keychain-lifecycle.sh"
  keychain_capture_search_list
  touch "$KEYCHAIN"
  printf '    "%s"\n' "$KEYCHAIN" >> "$MOCK_LIST"
  keychain_restore_search_list
  cmp "$WORK/keychains-before.txt" "$MOCK_LIST"
  [[ ! -e $KEYCHAIN ]]
  touch "$KEYCHAIN"
  MOCK_DELETE_FAIL=true
  if keychain_restore_search_list; then echo 'mock delete failure was ignored' >&2; exit 1; fi
  [[ -f $KEYCHAIN ]]
  MOCK_DELETE_FAIL=false
  MOCK_RESTORE_FAIL=true
  printf '    "%s"\n' "$KEYCHAIN" >> "$MOCK_LIST"
  if keychain_restore_search_list; then echo 'mock restoration failure was ignored' >&2; exit 1; fi
  # Readback mismatch must fail even when a tool returns zero for the set call.
  MOCK_RESTORE_FAIL=false
  MOCK_IGNORE_SET=true
  if keychain_restore_search_list; then echo 'mock successful-but-ineffective restoration was ignored' >&2; exit 1; fi
  MOCK_IGNORE_SET=false
  : > "$MOCK_LIST"
  keychain_capture_search_list
  touch "$KEYCHAIN"
  printf '    "%s"\n' "$KEYCHAIN" >> "$MOCK_LIST"
  keychain_restore_search_list
  [[ ! -s $MOCK_LIST && ! -e $KEYCHAIN ]]
  printf 'invalid unquoted snapshot\n' > "$MOCK_LIST"
  if keychain_capture_search_list; then echo 'invalid snapshot was accepted' >&2; exit 1; fi
)
echo 'PASS R1-05: mock-only exact user search-list/order restoration and propagated cleanup/snapshot failures'
echo 'HELD R1-04: Windows locked-payload worker regression is wired into hosted Windows CI; not run locally'
