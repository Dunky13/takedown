VERSION ?= $(shell git describe --tags --always --dirty)
LDFLAGS = -X github.com/Dunky13/takedown/internal/cli.Version=$(VERSION)
GO_BUILD = go build -p 2 -trimpath -ldflags "$(LDFLAGS)"
PACKAGE_VERSION ?= 0.0.0
MODE ?= local
# Optional reuse of an existing native GUI binary. Empty means build from source.
BINARY ?=
OUT_DIR ?=

.PHONY: build headless test check windows linux macos clean package-macos package-linux package-windows ci-native

build:
	mkdir -p bin
	$(GO_BUILD) -o bin/takedown ./cmd/takedown

headless:
	mkdir -p bin
	CGO_ENABLED=0 $(GO_BUILD) -tags headless -o bin/takedown-cli ./cmd/takedown

test:
	go test -p 2 -parallel 2 -race -tags ci ./...

check:
	@test -z "$$(gofmt -l cmd internal assets scripts/*.go)" || { gofmt -l cmd internal assets scripts/*.go; exit 1; }
	go vet -p 2 -tags ci ./...

# Fyne needs target C compilers and graphics headers. These targets deliberately
# require explicit compiler choices instead of silently producing a CLI-only app.
windows:
	@test -n "$(WINDOWS_CC)" || { echo 'Set WINDOWS_CC to a Windows C cross compiler'; exit 1; }
	mkdir -p dist
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC="$(WINDOWS_CC)" $(GO_BUILD) -o dist/takedown-windows-amd64.exe ./cmd/takedown

linux:
	@test -n "$(LINUX_CC)" || { echo 'Set LINUX_CC to a Linux C compiler with X11 and OpenGL headers'; exit 1; }
	mkdir -p dist
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC="$(LINUX_CC)" $(GO_BUILD) -o dist/takedown-linux-amd64 ./cmd/takedown

macos:
	@test -n "$(MACOS_CC)" || { echo 'Set MACOS_CC to a macOS C compiler with a valid macOS SDK'; exit 1; }
	mkdir -p dist
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 CC="$(MACOS_CC)" $(GO_BUILD) -o dist/takedown-macos-arm64 ./cmd/takedown

# These package on the native host and fail explicitly for missing prerequisites.
# PACKAGE_VERSION is numeric x.y.z; VERSION above remains build provenance.
package-macos:
	VERSION="$(PACKAGE_VERSION)" MODE="$(MODE)" BINARY="$(BINARY)" OUT_DIR="$(OUT_DIR)" bash scripts/package-macos.sh

package-linux:
	VERSION="$(PACKAGE_VERSION)" MODE="$(MODE)" BINARY="$(BINARY)" OUT_DIR="$(OUT_DIR)" bash scripts/package-linux.sh

package-windows:
	pwsh -NoLogo -NoProfile -File scripts/package-windows.ps1 -Version "$(PACKAGE_VERSION)" -Mode "$(MODE)" -Binary "$(BINARY)" -OutputDirectory "$(OUT_DIR)"

ci-native:
	CGO_ENABLED=1 bash scripts/ci-native.sh

clean:
	rm -rf bin dist
