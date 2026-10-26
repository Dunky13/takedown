#!/usr/bin/env bash
set -euo pipefail
cd -- "$(dirname -- "$0")/.."
[[ $(go env GOVERSION) == go1.26.10 ]] || { echo 'Go 1.26.10 is required' >&2; exit 1; }
[[ ${CGO_ENABLED:-} == 1 ]] || { echo 'Native CI requires CGO_ENABLED=1' >&2; exit 1; }
goos=$(go env GOOS)
goarch=$(go env GOARCH)
[[ $goos == "$(go env GOHOSTOS)" && $goarch == "$(go env GOHOSTARCH)" ]] || { echo 'Native CI must build for the host OS and architecture' >&2; exit 1; }
[[ $goos == "${CI_EXPECTED_GOOS:-$goos}" && $goarch == "${CI_EXPECTED_GOARCH:-$goarch}" ]] || { echo 'Native CI runner does not match the required OS and architecture' >&2; exit 1; }
[[ $(go list -m -mod=readonly -f '{{.Version}}' fyne.io/fyne/v2) == v2.8.2 ]] || { echo 'Fyne v2.8.2 is required' >&2; exit 1; }
unformatted=$(gofmt -l cmd internal assets scripts/*.go)
[[ -z $unformatted ]] || { printf 'gofmt required:\n%s\n' "$unformatted" >&2; exit 1; }
go mod verify
# Integrity alone does not detect stale requirements or missing checksum entries.
# Report the required changes without rewriting the checked-out module files.
go mod tidy -diff
# Fyne's ci driver makes tests independent of a display server. The separate
# untagged build/vet below MUST compile/link the actual native OpenGL GUI.
go test -mod=readonly -p 2 -parallel 2 -race -tags ci ./...
go vet -mod=readonly -p 2 -tags ci ./...
go vet -mod=readonly -p 2 ./...
go list -mod=readonly -deps ./cmd/takedown | grep -Fx 'fyne.io/fyne/v2/internal/driver/glfw'
mkdir -p bin
suffix=''
[[ $goos != windows ]] || suffix=.exe
go build -mod=readonly -p 2 -trimpath \
  -ldflags '-X github.com/Dunky13/takedown/internal/cli.Version=0.0.0' \
  -o "bin/takedown$suffix" ./cmd/takedown
"bin/takedown$suffix" --version
"bin/takedown$suffix" --help
