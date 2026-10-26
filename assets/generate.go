//go:build ignore

// Run from the repository root:
// go run -mod=readonly -p 2 -tags=noasm -gcflags=all=-d=fmahash=n ./assets/generate.go
package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"

	"github.com/fyne-io/oksvg"
	"github.com/srwiley/rasterx"
	"golang.org/x/image/draw"
)

var sizes = []int{16, 24, 32, 48, 64, 128, 256, 512, 1024}

const (
	generatorTags          = "noasm"
	generatorCompilerFlags = "all=-d=fmahash=n"
)

func main() {
	check := flag.Bool("check", false, "verify generated files without modifying the checkout")
	flag.Parse()
	if err := generate(*check); err != nil {
		fmt.Fprintln(os.Stderr, "icons:", err)
		os.Exit(1)
	}
}

func generate(check bool) error {
	if runtime.Version() != "go1.26.10" {
		return fmt.Errorf("Go 1.26.10 is required for reproducible icon generation (running %s)", runtime.Version())
	}
	// rasterx's implicit FMA and x/image/vector's handwritten SIMD accumulation
	// independently change pixels across architectures. Require generic raster
	// code and non-fused operations before reading or writing any assets.
	info, ok := debug.ReadBuildInfo()
	if !ok || !hasGeneratorBuildSettings(info) {
		return fmt.Errorf("icon generation requires -tags=%s and -gcflags=%s; run from the repository root:\n"+
			"go run -mod=readonly -p 2 -tags=%s -gcflags=%s ./assets/generate.go%s",
			generatorTags, generatorCompilerFlags, generatorTags, generatorCompilerFlags, checkArgument(check))
	}
	svg, err := os.ReadFile("assets/icon.svg")
	if err != nil {
		return fmt.Errorf("run from the repository root: %w", err)
	}
	icon, err := oksvg.ReadIconStream(bytes.NewReader(svg), oksvg.StrictErrorMode)
	if err != nil {
		return fmt.Errorf("parse SVG: %w", err)
	}
	if icon.ViewBox.W != 512 || icon.ViewBox.H != 512 {
		return fmt.Errorf("expected a 512 by 512 SVG viewBox")
	}
	// Render once at 2x the largest output, then use the same deterministic
	// Catmull-Rom filter for each size. No platform fonts or native tools are used.
	const renderSize = 2048
	master := image.NewRGBA(image.Rect(0, 0, renderSize, renderSize))
	icon.SetTarget(0, 0, renderSize, renderSize)
	icon.Draw(rasterx.NewDasher(renderSize, renderSize,
		rasterx.NewScannerGV(renderSize, renderSize, master, master.Bounds())), 1)
	images := make(map[int]*image.NRGBA)
	pngs := make(map[int][]byte)
	files := make(map[string][]byte)
	for _, size := range sizes {
		img := image.NewNRGBA(image.Rect(0, 0, size, size))
		draw.CatmullRom.Scale(img, img.Bounds(), master, master.Bounds(), draw.Src, nil)
		var encoded bytes.Buffer
		if err := png.Encode(&encoded, img); err != nil {
			return err
		}
		images[size], pngs[size] = img, encoded.Bytes()
		files[fmt.Sprintf("assets/sizes/icon-%d.png", size)] = encoded.Bytes()
	}
	files["assets/icon.png"] = pngs[512]
	files["assets/icon.ico"] = encodeICO(images, pngs[256])
	files["assets/icon.icns"] = encodeICNS(pngs)

	// rsrc writes a COFF resource object for Go's linker. It never opens or
	// patches an executable. The explicit tool version is outside go.mod.
	work, err := os.MkdirTemp("", "takedown-icons-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(work)
	icoPath := filepath.Join(work, "icon.ico")
	if err := os.WriteFile(icoPath, files["assets/icon.ico"], 0644); err != nil {
		return err
	}
	sysoPath := filepath.Join(work, "icon_windows_amd64.syso")
	cmd := exec.Command("go", "run", "-mod=readonly", "-p", "2",
		"github.com/akavel/rsrc@v0.10.2", "-arch", "amd64", "-ico", icoPath, "-o", sysoPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("rsrc v0.10.2: %w\n%s", err, output)
	}
	files["cmd/takedown/icon_windows_amd64.syso"], err = os.ReadFile(sysoPath)
	if err != nil {
		return err
	}
	paths := []string{"assets/icon.png", "assets/icon.ico", "assets/icon.icns", "cmd/takedown/icon_windows_amd64.syso"}
	for _, size := range sizes {
		paths = append(paths, fmt.Sprintf("assets/sizes/icon-%d.png", size))
	}
	for _, path := range paths {
		if check {
			current, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if !bytes.Equal(current, files[path]) {
				return fmt.Errorf("%s is stale; regenerate the icons", path)
			}
		} else {
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return err
			}
			if err := os.WriteFile(path, files[path], 0644); err != nil {
				return err
			}
		}
	}
	if check {
		fmt.Printf("Verified %d generated icon assets byte for byte.\n", len(paths))
	} else {
		fmt.Printf("Generated %d icon assets from assets/icon.svg.\n", len(paths))
	}
	return nil
}

func hasGeneratorBuildSettings(info *debug.BuildInfo) bool {
	var tags, compilerFlags string
	for _, setting := range info.Settings {
		switch setting.Key {
		case "-tags":
			tags = setting.Value
		case "-gcflags":
			compilerFlags = setting.Value
		}
	}
	return tags == generatorTags && compilerFlags == generatorCompilerFlags
}

func checkArgument(check bool) string {
	if check {
		return " -check"
	}
	return ""
}

// ICO uses 32-bit bottom-up BGRA DIBs with transparency masks at sizes below
// 256, for compatibility with NSIS, and a PNG at 256 for Windows Explorer.
func encodeICO(images map[int]*image.NRGBA, largePNG []byte) []byte {
	iconSizes := []int{16, 24, 32, 48, 64, 128, 256}
	data := make([][]byte, len(iconSizes))
	for i, size := range iconSizes {
		if size == 256 {
			data[i] = largePNG
		} else {
			data[i] = encodeDIB(images[size])
		}
	}
	result := make([]byte, 6+16*len(data))
	binary.LittleEndian.PutUint16(result[2:], 1)
	binary.LittleEndian.PutUint16(result[4:], uint16(len(data)))
	offset := len(result)
	for i, size := range iconSizes {
		entry := result[6+16*i:]
		entry[0], entry[1] = byte(size%256), byte(size%256)
		binary.LittleEndian.PutUint16(entry[4:], 1)
		binary.LittleEndian.PutUint16(entry[6:], 32)
		binary.LittleEndian.PutUint32(entry[8:], uint32(len(data[i])))
		binary.LittleEndian.PutUint32(entry[12:], uint32(offset))
		offset += len(data[i])
	}
	for _, payload := range data {
		result = append(result, payload...)
	}
	return result
}

func encodeDIB(img *image.NRGBA) []byte {
	size := img.Bounds().Dx()
	maskStride := ((size + 31) / 32) * 4
	result := make([]byte, 40+size*size*4+maskStride*size)
	binary.LittleEndian.PutUint32(result, 40)
	binary.LittleEndian.PutUint32(result[4:], uint32(size))
	binary.LittleEndian.PutUint32(result[8:], uint32(size*2))
	binary.LittleEndian.PutUint16(result[12:], 1)
	binary.LittleEndian.PutUint16(result[14:], 32)
	binary.LittleEndian.PutUint32(result[20:], uint32(size*size*4))
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			pixel := img.NRGBAAt(x, size-1-y)
			i := 40 + (y*size+x)*4
			result[i], result[i+1], result[i+2], result[i+3] = pixel.B, pixel.G, pixel.R, pixel.A
			if pixel.A == 0 {
				result[40+size*size*4+y*maskStride+x/8] |= 0x80 >> (x % 8)
			}
		}
	}
	return result
}

// Modern ICNS chunks contain PNGs. Include standard and Retina representations
// explicitly, without iconutil or any platform-dependent image conversion.
func encodeICNS(pngs map[int][]byte) []byte {
	entries := []struct {
		kind string
		size int
	}{
		{"icp4", 16}, {"icp5", 32}, {"icp6", 64}, {"ic07", 128},
		{"ic08", 256}, {"ic09", 512}, {"ic10", 1024},
		{"ic11", 32}, {"ic12", 64}, {"ic13", 256}, {"ic14", 512},
	}
	result := []byte{'i', 'c', 'n', 's', 0, 0, 0, 0}
	for _, entry := range entries {
		chunk := make([]byte, 8)
		copy(chunk, entry.kind)
		binary.BigEndian.PutUint32(chunk[4:], uint32(8+len(pngs[entry.size])))
		result = append(result, chunk...)
		result = append(result, pngs[entry.size]...)
	}
	binary.BigEndian.PutUint32(result[4:], uint32(len(result)))
	return result
}
