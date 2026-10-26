package assets

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestGeneratorRejectsUnpinnedCompilerSettingsBeforeWriting(t *testing.T) {
	const command = "go run -mod=readonly -p 2 -gcflags=all=-d=fmahash=n ./assets/generate.go"
	for _, tc := range []struct {
		name  string
		flags []string
		check bool
	}{
		{name: "plain generation"},
		{name: "plain check", check: true},
		{name: "only one dependency pinned", flags: []string{"-gcflags=github.com/srwiley/rasterx=-d=fmahash=n"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			before := generatedAssetDigests(t)
			args := append([]string{"run", "-mod=readonly", "-p", "2"}, tc.flags...)
			args = append(args, "./assets/generate.go")
			wantCommand := command
			if tc.check {
				args = append(args, "-check")
				wantCommand += " -check"
			}
			cmd := exec.Command("go", args...)
			cmd.Dir = ".."
			// The tested invocations deliberately omit the required all-packages
			// setting, including when the outer test runner uses GOFLAGS.
			cmd.Env = append(os.Environ(), "GOFLAGS=")
			output, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatalf("generator accepted unpinned compilation: %s", output)
			}
			if !strings.Contains(string(output), "icon generation requires -gcflags=all=-d=fmahash=n") ||
				!strings.Contains(string(output), wantCommand) {
				t.Fatalf("expected compiler-setting refusal and actionable command, got: %s", output)
			}
			after := generatedAssetDigests(t)
			for path, digest := range before {
				if after[path] != digest {
					t.Errorf("refused generator modified %s", path)
				}
			}
		})
	}
}

func generatedAssetDigests(t *testing.T) map[string][sha256.Size]byte {
	t.Helper()
	paths := []string{"assets/icon.png", "assets/icon.ico", "assets/icon.icns", "cmd/takedown/icon_windows_amd64.syso"}
	for _, size := range []int{16, 24, 32, 48, 64, 128, 256, 512, 1024} {
		paths = append(paths, fmt.Sprintf("assets/sizes/icon-%d.png", size))
	}
	digests := make(map[string][sha256.Size]byte, len(paths))
	for _, path := range paths {
		data, err := os.ReadFile(filepath.Join("..", path))
		if err != nil {
			t.Fatal(err)
		}
		digests[path] = sha256.Sum256(data)
	}
	return digests
}
