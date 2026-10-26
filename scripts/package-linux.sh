#!/usr/bin/env bash
set -euo pipefail
# shellcheck source=scripts/package-common.sh
source "$(dirname -- "$0")/package-common.sh"
[[ ${MODE:-local} == local ]] || fail 'Linux release signing has no configured trust policy; use MODE=local'
package_init linux
need tar
need gzip
need sha256sum
BASE="takedown-$VERSION-linux-$ARCH-$MODE"
STAGE="$WORK/$BASE"
mkdir -p "$STAGE/bin" "$STAGE/share/applications" "$STAGE/share/icons/hicolor/scalable/apps"
cp "$WORK/takedown" "$STAGE/bin/takedown"
cp "$ROOT/build/linux/takedown.desktop" "$STAGE/share/applications/"
cp "$ROOT/assets/icon.svg" "$STAGE/share/icons/hicolor/scalable/apps/takedown.svg"
cp "$ROOT/build/linux/install-user.sh" "$ROOT/build/linux/uninstall-user.sh" "$ROOT/build/linux/install-common.sh" "$STAGE/"
package_docs "$STAGE/share/doc/takedown"
# GNU tar's stable ordering, owner, mode, and timestamp + gzip without a header time.
# SOURCE_DATE_EPOCH defaults to zero; signed/platform outputs are not bit-reproducible.
SOURCE_DATE_EPOCH=${SOURCE_DATE_EPOCH:-0}
[[ $SOURCE_DATE_EPOCH =~ ^[0-9]+$ ]] || fail 'SOURCE_DATE_EPOCH must be a nonnegative integer'
tar --version | grep -q 'GNU tar' || fail 'Linux packaging requires GNU tar'
tar --sort=name --mtime="@$SOURCE_DATE_EPOCH" --owner=0 --group=0 --numeric-owner \
  --mode='u+rwX,go+rX,go-w' -C "$WORK" -cf - "$BASE" | gzip -n > "$WORK/output/$BASE.tar.gz"
(cd "$WORK/output" && sha256sum "$BASE.tar.gz" > SHA256SUMS)
package_finish
