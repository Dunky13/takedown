#!/usr/bin/env bash
# Sourced by native Unix packagers; no Fyne CLI or unpinned tool installation.
set -euo pipefail
umask 022

fail() { printf 'package: %s\n' "$*" >&2; exit 1; }
need() { command -v "$1" >/dev/null 2>&1 || fail "missing prerequisite: $1"; }

package_reserve_output() {
  # mkdir is the exclusive reservation. Never rename a directory onto another
  # publisher's directory: that operation can silently create OUT_DIR/output.
  mkdir -- "$OUT_DIR" || fail 'OUT_DIR already exists or is reserved by another packager'
  OUTPUT_RESERVED=true
  printf '%s\n' "$WORK" > "$OUT_DIR/.package-incomplete"
}

package_cleanup() {
  local result=$?
  trap - EXIT
  if [[ ${OUTPUT_RESERVED:-false} == true ]]; then
    # Only this operation's reservation may be removed on failure.
    if [[ -f $OUT_DIR/.package-incomplete && ! -L $OUT_DIR/.package-incomplete && $(cat "$OUT_DIR/.package-incomplete") == "$WORK" ]]; then
      rm -rf -- "$OUT_DIR" || result=1
    else
      printf 'package: reservation changed; preserving OUT_DIR for inspection\n' >&2
      result=1
    fi
  fi
  rm -rf -- "$WORK" || result=1
  exit "$result"
}

package_init() {
  need go
  ROOT=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)
  cd "$ROOT"
  VERSION=${VERSION:-0.0.0}
  MODE=${MODE:-local}
  [[ $VERSION =~ ^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$ ]] || fail 'VERSION must be a numeric x.y.z (without v)'
  [[ $MODE == local || $MODE == release ]] || fail 'MODE must be local or release'
  [[ $(go env GOVERSION) == go1.26.10 ]] || fail 'Go 1.26.10 is required'
  [[ $(go env GOHOSTOS) == "$1" ]] || fail "run this packager natively on $1"
  ARCH=$(go env GOHOSTARCH)
  [[ $ARCH == arm64 || $ARCH == amd64 ]] || fail "unsupported architecture: $ARCH"
  [[ $(go list -m -mod=readonly -f '{{.Version}}' fyne.io/fyne/v2) == v2.8.2 ]] || fail 'Fyne v2.8.2 is required'
  OUT_DIR=${OUT_DIR:-"$ROOT/dist/packages/$1-$ARCH-$MODE"}
  [[ $OUT_DIR == /* ]] || OUT_DIR="$ROOT/$OUT_DIR"
  [[ ! -e $OUT_DIR && ! -L $OUT_DIR ]] || fail 'OUT_DIR already exists; choose a fresh directory'
  mkdir -p "$(dirname -- "$OUT_DIR")"
  WORK=$(mktemp -d "$(dirname -- "$OUT_DIR")/.package.XXXXXX")
  OUTPUT_RESERVED=false
  trap package_cleanup EXIT
  trap 'exit 130' INT
  trap 'exit 143' TERM
  package_reserve_output
  mkdir "$WORK/output"
  if [[ -n ${BINARY:-} ]]; then
    [[ -f $BINARY && ! -L $BINARY ]] || fail 'BINARY must be an existing regular native Go executable'
    cp -- "$BINARY" "$WORK/takedown"
  else
    need cc
    CGO_ENABLED=1 GOOS="$1" GOARCH="$ARCH" go build -mod=readonly -p 2 -trimpath \
      -ldflags "-X github.com/Dunky13/takedown/internal/cli.Version=$VERSION" \
      -o "$WORK/takedown" ./cmd/takedown
  fi
  # Reject optional CLI-only and Fyne's software CI driver, including reused binaries.
  go version -m "$WORK/takedown" > "$WORK/build-info"
  grep -Eq ': go1\.26\.10$' "$WORK/build-info" || fail 'binary must be built with Go 1.26.10'
  grep -Eq 'dep[[:space:]]+fyne.io/fyne/v2[[:space:]]+v2.8.2[[:space:]]' "$WORK/build-info" || fail 'binary lacks pinned Fyne dependency'
  grep -q 'CGO_ENABLED=1' "$WORK/build-info" || fail 'binary must enable native Cgo'
  grep -q "GOOS=$1" "$WORK/build-info" || fail 'binary has the wrong operating system'
  grep -q "GOARCH=$ARCH" "$WORK/build-info" || fail 'binary has the wrong architecture'
  if grep -E -- '-tags=.*(headless|ci)' "$WORK/build-info" >/dev/null; then
    fail 'headless/ci tagged binaries cannot be packaged as native desktop apps'
  fi
  chmod 755 "$WORK/takedown"
}

package_docs() {
  mkdir -p "$1"
  cp README.md PRIVACY.md PRIVACY.html "$1/"
  # Preserve an existing license; never choose one on the distributor's behalf.
  for license in LICENSE LICENSE.md LICENSE.txt COPYING; do
    [[ ! -f $license ]] || cp "$license" "$1/"
  done
  {
    printf 'Package version: %s\nMode: %s\nArchitecture: %s\n' "$VERSION" "$MODE" "$ARCH"
    # First line contains the local binary path, so omit it from distributed metadata.
    sed '1d' "$WORK/build-info"
  } > "$1/BUILD-INFO.txt"
}

package_finish() {
  [[ ${OUTPUT_RESERVED:-false} == true && -f $OUT_DIR/.package-incomplete && ! -L $OUT_DIR/.package-incomplete ]] || fail 'output publication requires an owned reservation'
  [[ $(cat "$OUT_DIR/.package-incomplete") == "$WORK" ]] || fail 'output reservation owner changed'
  local entry
  # OUT_DIR is ours for the whole operation. Preflight the complete entry set;
  # never move the aggregate output directory using directory-merging semantics.
  for entry in "$WORK/output"/*; do
    [[ -e $entry ]] || fail 'no artifacts to publish'
    [[ ! -e $OUT_DIR/$(basename -- "$entry") && ! -L $OUT_DIR/$(basename -- "$entry") ]] || fail 'conflicting artifact in reserved output'
  done
  for entry in "$WORK/output"/*; do mv -- "$entry" "$OUT_DIR/"; done
  rm -- "$OUT_DIR/.package-incomplete"
  OUTPUT_RESERVED=false
  printf 'Artifacts: %s\n' "$OUT_DIR"
}
