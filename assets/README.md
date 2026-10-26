# Application icon

`icon.svg` is the canonical, original takedown artwork, covered by the repository's
Apache-2.0 license (see [LICENSE](../LICENSE)). The archive box and downward arrow
represent taking an archive out and downloading it. The navy, indigo, and azure
palette uses no third-party marks, fonts, or external images.

## Generate and verify

From the repository root with Go 1.26.10:

```sh
go run -mod=readonly -p 2 ./assets/generate.go
go run -mod=readonly -p 2 ./assets/generate.go -check
```

The first command regenerates the assets. The second renders into memory and a
temporary directory, compares all 13 generated files byte for byte, and exits
with an error if a file is missing or stale. It does not modify the checkout.

The generator uses the SVG rasterizer and image scaler already pinned in
`go.mod` and `go.sum`. It rasterizes at 2048 pixels, then downsamples each output
with Catmull-Rom filtering. Explicit rectangle radii and relative gradients keep
the SVG compatible with that rasterizer. PNG, ICO, and ICNS encoding has no
timestamps, fonts, platform image APIs, or native converter prerequisites.

Windows resource generation invokes the standard
[`github.com/akavel/rsrc`](https://github.com/akavel/rsrc) tool at **v0.10.2**, with
`-arch amd64`. Go fetches and verifies that pinned tool if it is not cached;
first use requires access to the configured Go module proxy and checksum
database. This build-only tool does not change the application's dependencies
and does not patch an executable. Commit regenerated assets together with any
SVG changes so normal builds need no generator or converter.

## Files and consumers

| File | Purpose |
| --- | --- |
| `icon.svg` | Canonical vector; Linux scalable application icon |
| `icon.png` | 512-pixel PNG; embedded GUI resource and README artwork |
| `sizes/icon-{16,24,32,48,64,128,256,512,1024}.png` | Raster representations for inspection and native containers |
| `icon.ico` | Windows installer and uninstaller; 16, 24, 32, 48, 64, 128, and 256 pixels |
| `icon.icns` | macOS bundle; standard and Retina representations through 1024 pixels |
| `../cmd/takedown/icon_windows_amd64.syso` | Generated amd64 COFF resources, linked automatically by Go on Windows |

ICO entries below 256 pixels use 32-bit BGRA bitmaps and transparency masks for
NSIS compatibility; the 256-pixel entry uses PNG. ICNS contains PNG chunks for
standard and Retina resolutions. macOS packaging copies `icon.icns` to
`Contents/Resources/icon.icns`; `CFBundleIconFile` names it in the bundle plist.
NSIS uses the same ICO for both its installer and uninstaller.

The `assets` Go package embeds only `icon.png`. The native desktop sets that
resource as the app icon before creating its first window. It does not read
artwork from disk at runtime, and headless builds remain independent of the GUI.
`FyneApp.toml` also records the PNG as informational metadata.

For a quick macOS container check, run:

```sh
iconutil -c iconset assets/icon.icns -o icon-check.iconset
```

The resulting directory contains the standard and Retina PNGs and can be
removed after inspection. Regeneration does not require `iconutil`.
