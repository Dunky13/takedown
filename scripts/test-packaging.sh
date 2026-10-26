#!/usr/bin/env bash
# Focused offline contract checks. No certificates, account access, GUI launches,
# compilation, or claims about Windows/macOS signature acceptance.
set -euo pipefail
ROOT=$(cd -- "$(dirname -- "$0")/.." && pwd)
cd "$ROOT"
WORK=$(mktemp -d "${TMPDIR:-/tmp}/takedown-package-tests.XXXXXX")
WORK=$(cd -P "$WORK" && pwd -P)
trap 'rm -rf -- "$WORK"' EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
expect_failure() {
  local expected=$1
  shift
  if "$@" > "$WORK/failure.log" 2>&1; then
    echo "expected failure: $expected" >&2; exit 1
  fi
  grep -Fq "$expected" "$WORK/failure.log" || { cat "$WORK/failure.log"; exit 1; }
}
expect_failure 'VERSION must be a numeric x.y.z' env MODE=local VERSION='bad;value' bash scripts/package-macos.sh
expect_failure 'Linux release signing has no configured trust policy' env MODE=release bash scripts/package-linux.sh
if [[ $(uname -s) == Darwin ]]; then
  expect_failure 'missing required secret/configuration: MACOS_CERTIFICATE_P12_BASE64' \
    env -i PATH="$PATH" MODE=release bash scripts/package-macos.sh
fi
# Test actual installer integration using a harmless payload and a prefix with spaces.
FIXTURE="$WORK/payload"
PREFIX="$WORK/install prefix"
mkdir -p "$FIXTURE/bin" "$FIXTURE/share/applications" \
  "$FIXTURE/share/icons/hicolor/scalable/apps" "$FIXTURE/share/doc/takedown"
printf '#!/bin/sh\nexit 0\n' > "$FIXTURE/bin/takedown"
cp build/linux/takedown.desktop "$FIXTURE/share/applications/"
cp assets/icon.svg "$FIXTURE/share/icons/hicolor/scalable/apps/takedown.svg"
printf 'test docs\n' > "$FIXTURE/share/doc/takedown/README.md"
printf 'fixture application license\n' > "$FIXTURE/share/doc/takedown/LICENSE"
printf 'fixture application notice\n' > "$FIXTURE/share/doc/takedown/NOTICE"
printf 'fixture dependency notices\n' > "$FIXTURE/share/doc/takedown/THIRD_PARTY_NOTICES.md"
cp build/linux/install-user.sh build/linux/uninstall-user.sh build/linux/install-common.sh "$FIXTURE/"
INSTALL_PREFIX="$PREFIX" bash "$FIXTURE/install-user.sh"
[[ -x $PREFIX/bin/takedown && -f $PREFIX/share/doc/takedown/README.md ]]
for name in LICENSE NOTICE THIRD_PARTY_NOTICES.md; do
  cmp "$FIXTURE/share/doc/takedown/$name" "$PREFIX/share/doc/takedown/$name"
  grep -Eq "^[0-9a-f]{64}  share/doc/takedown/$name$" "$PREFIX/share/takedown/install-manifest.sha256"
done
grep -Fxq "Exec=\"$PREFIX/bin/takedown\" gui" "$PREFIX/share/applications/takedown.desktop"
printf 'retain me\n' > "$PREFIX/user-state"
INSTALL_PREFIX="$PREFIX" bash build/linux/uninstall-user.sh
[[ ! -e $PREFIX/bin/takedown && ! -e $PREFIX/share/applications/takedown.desktop && -f $PREFIX/user-state ]]
[[ ! -e $PREFIX/share/doc/takedown/LICENSE && ! -e $PREFIX/share/doc/takedown/NOTICE && ! -e $PREFIX/share/doc/takedown/THIRD_PARTY_NOTICES.md ]]
expect_failure 'INSTALL_PREFIX must be absolute' env INSTALL_PREFIX=relative bash "$FIXTURE/install-user.sh"
expect_failure 'unsupported desktop-entry characters' env INSTALL_PREFIX="$WORK/bad%prefix" bash "$FIXTURE/install-user.sh"
if [[ $(uname -s) == Darwin && -n ${NATIVE_BINARY:-} ]]; then
  # No build or DMG generation occurs on these rejection paths.
  mkdir "$WORK/existing-output"
  expect_failure 'OUT_DIR already exists' env MODE=local BINARY="$NATIVE_BINARY" OUT_DIR="$WORK/existing-output" bash scripts/package-macos.sh
  if [[ -n ${HEADLESS_BINARY:-} ]]; then
    expect_failure 'binary lacks pinned Fyne dependency' env MODE=local BINARY="$HEADLESS_BINARY" OUT_DIR="$WORK/rejected" bash scripts/package-macos.sh
    [[ ! -e $WORK/rejected ]]
  fi
fi
# Exercise the actual Unix documentation payload shared by macOS and Linux.
# A disposable source fixture also works before the root license owner lands.
DOC_SOURCE="$WORK/docs-source"
mkdir "$DOC_SOURCE"
cp README.md PRIVACY.md PRIVACY.html THIRD_PARTY_NOTICES.md "$DOC_SOURCE/"
printf 'fixture license\n' > "$DOC_SOURCE/LICENSE"
printf 'fixture notice\n' > "$DOC_SOURCE/NOTICE"
printf '/fixture/takedown: go1.26.10\n\tbuild\tGOOS=fixture\n' > "$WORK/build-info"
cat > "$WORK/copy-docs.sh" <<'DOCS'
#!/usr/bin/env bash
set -euo pipefail
source "$1/scripts/package-common.sh"
cd "$2"
WORK=$3 VERSION=0.0.0 MODE=local ARCH=fixture
package_docs "$4"
DOCS
bash "$WORK/copy-docs.sh" "$ROOT" "$DOC_SOURCE" "$WORK" "$WORK/docs-stage"
for name in LICENSE NOTICE THIRD_PARTY_NOTICES.md; do
  cmp "$DOC_SOURCE/$name" "$WORK/docs-stage/$name"
  mv "$DOC_SOURCE/$name" "$WORK/held-doc"
  expect_failure "missing regular license payload: $name" bash "$WORK/copy-docs.sh" "$ROOT" "$DOC_SOURCE" "$WORK" "$WORK/rejected-docs"
  [[ ! -e $WORK/rejected-docs ]]
  ln -s "$WORK/held-doc" "$DOC_SOURCE/$name"
  expect_failure "missing regular license payload: $name" bash "$WORK/copy-docs.sh" "$ROOT" "$DOC_SOURCE" "$WORK" "$WORK/rejected-docs"
  rm "$DOC_SOURCE/$name"
  mv "$WORK/held-doc" "$DOC_SOURCE/$name"
done
echo 'PASS: invalid version, fail-closed release, installer paths/removal, and optional native-input rejection checks'
echo 'PASS: required license/notice payload copied intact; missing and symlinked legal files rejected'
