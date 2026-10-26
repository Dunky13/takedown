//go:build ignore

// Run from the repository root with Go 1.26.10:
//
//	go run -mod=readonly -p 2 scripts/generate-third-party-notices.go
//	go run -mod=readonly -p 2 scripts/generate-third-party-notices.go -check
//
// No native compiler, display, credentials, or additional tools are needed.
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"

	"golang.org/x/image/font/sfnt"
)

const outputFile = "THIRD_PARTY_NOTICES.md"

var targets = []string{"darwin/amd64", "darwin/arm64", "linux/amd64", "windows/amd64"}
var legalName = regexp.MustCompile(`(?i)^(licen[cs]e|copying|notice|copyright|patents|authors|ofl)([._-].*)?$`)
var legalText = regexp.MustCompile(`(?i)copyright|SPDX-License-Identifier|permission is (hereby )?granted|redistribution and use|licensed under|public domain`)
var cComments = regexp.MustCompile(`(?s)/\*.*?\*/|(?m)(^[\t ]*//[^\n]*(\n|$))+`)
var cIncludes = regexp.MustCompile(`(?m)^[\t ]*#[\t ]*include[\t ]*[<"]([^>"\n]+)[>"]`)

type module struct {
	Path, Version, Dir, Sum string
	Main                    bool
	Replace                 *module
}

type goPackage struct {
	Dir, ImportPath, ForTest                            string
	Goroot, Standard, Incomplete                        bool
	Module                                              *module
	GoFiles, CgoFiles, CFiles, CXXFiles, MFiles, HFiles []string
	FFiles, SFiles, SwigFiles, SwigCXXFiles, EmbedFiles []string
	CgoCFLAGS, CgoCPPFLAGS                              []string
}

type component struct {
	ID, Root, Sum, Source string
	Targets               map[string]bool
	Sources, Fonts        map[string]bool
	Legal                 map[string]bool
	IncludeDirs           map[string]bool
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "notices:", err)
		os.Exit(1)
	}
}

func goCommand(target string, args ...string) []byte {
	cmd := exec.Command("go", args...)
	// Ignore local build flags/workspaces and prohibit downloading missing source.
	// Native runtime graphs always use cgo and the default build tags.
	cmd.Env = append(os.Environ(), "GOFLAGS=", "GOWORK=off", "GOPROXY=off", "CGO_ENABLED=1", "GOEXPERIMENT=", "GOAMD64=v1", "GOARM64=v8.0", "GOFIPS140=off", "GOTOOLCHAIN=go1.26.10")
	parts := strings.Split(target, "/")
	cmd.Env = append(cmd.Env, "GOOS="+parts[0], "GOARCH="+parts[1])
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		must(fmt.Errorf("go %s: %w\n%s", strings.Join(args, " "), err, stderr.String()))
	}
	return out
}

func newComponent(id, root, sum, source string) *component {
	return &component{id, root, sum, source, map[string]bool{}, map[string]bool{}, map[string]bool{}, map[string]bool{}, map[string]bool{}}
}

func portablePath(c *component, path string) string {
	rel, err := filepath.Rel(c.Root, path)
	must(err)
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
		must(fmt.Errorf("source escapes component %s", c.ID))
	}
	return c.ID + "/" + filepath.ToSlash(rel)
}

func readFile(path string) []byte {
	info, err := os.Lstat(path)
	must(err)
	if !info.Mode().IsRegular() {
		must(fmt.Errorf("expected regular source file: %s", filepath.Base(path)))
	}
	data, err := os.ReadFile(path)
	must(err)
	return data
}

func sortedKeys[V any](values map[string]V) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func sortedPaths(values map[string]bool) []string {
	keys := sortedKeys(values)
	sort.Slice(keys, func(i, j int) bool { return filepath.ToSlash(keys[i]) < filepath.ToSlash(keys[j]) })
	return keys
}

func collectLegal(c *component) {
	must(filepath.WalkDir(c.Root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if legalName.MatchString(entry.Name()) && !strings.HasSuffix(entry.Name(), ".go") {
			c.Legal[path] = true
		}
		return nil
	}))
	if len(c.Legal) == 0 {
		must(fmt.Errorf("no original license files found for %s", c.ID))
	}
}

func sourceNotices(path string) []string {
	data := readFile(path)
	var comments []string
	if filepath.Ext(path) == ".go" {
		positions := token.NewFileSet()
		file, err := parser.ParseFile(positions, "source.go", data, parser.ParseComments|parser.SkipObjectResolution)
		must(err)
		for _, group := range file.Comments {
			start := positions.Position(group.Pos()).Offset
			last := group.List[len(group.List)-1]
			end := positions.Position(last.Pos()).Offset
			if strings.HasPrefix(last.Text, "/*") {
				end += bytes.Index(data[end:], []byte("*/")) + 2
			} else if newline := bytes.IndexByte(data[end:], '\n'); newline >= 0 {
				end += newline
			} else {
				end = len(data)
			}
			comments = append(comments, string(data[start:end]))
		}
	} else {
		comments = cComments.FindAllString(string(data), -1)
	}
	var result []string
	for _, comment := range comments {
		if legalText.MatchString(comment) {
			result = append(result, comment)
		}
	}
	return result
}

// Follow literal native includes shipped inside each component. Conditional
// includes are collected conservatively without preprocessing or reading system
// headers. This captures GLFW's vendored sources without unrelated examples.
func collectIncludes(c *component) {
	queue := sortedPaths(c.Sources)
	for i := 0; i < len(queue); i++ {
		path := queue[i]
		for _, include := range cIncludes.FindAllSubmatch(readFile(path), -1) {
			name := filepath.FromSlash(string(include[1]))
			dirs := append([]string{filepath.Dir(path)}, sortedPaths(c.IncludeDirs)...)
			for _, dir := range dirs {
				candidate := filepath.Clean(filepath.Join(dir, name))
				rel, err := filepath.Rel(c.Root, candidate)
				must(err)
				if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
					continue
				}
				info, err := os.Lstat(candidate)
				if errors.Is(err, os.ErrNotExist) {
					continue
				}
				must(err)
				if !info.Mode().IsRegular() {
					must(fmt.Errorf("non-regular native include in %s", c.ID))
				}
				if !c.Sources[candidate] {
					c.Sources[candidate] = true
					queue = append(queue, candidate)
				}
				break
			}
		}
	}
}

func literal(out *bytes.Buffer, data []byte) {
	// Preserve the original bytes inside a fence longer than any upstream fence.
	fence := "````"
	for bytes.Contains(data, []byte(fence)) {
		fence += "`"
	}
	fmt.Fprintln(out, fence+"text")
	out.Write(data)
	if !bytes.HasSuffix(data, []byte("\n")) {
		out.WriteByte('\n')
	}
	fmt.Fprintln(out, fence)
	fmt.Fprintln(out)
}

func fontNotices(path string) []byte {
	font, err := sfnt.Parse(readFile(path))
	must(err)
	var out bytes.Buffer
	for _, field := range []struct {
		ID   sfnt.NameID
		Name string
	}{{sfnt.NameIDCopyright, "Copyright"}, {sfnt.NameIDTrademark, "Trademark"}, {sfnt.NameIDLicense, "License"}, {sfnt.NameIDLicenseURL, "License URL"}} {
		value, err := font.Name(nil, field.ID)
		if errors.Is(err, sfnt.ErrNotFound) {
			continue
		}
		must(err)
		if value != "" {
			fmt.Fprintf(&out, "%s:\n%s\n\n", field.Name, value)
		}
	}
	if out.Len() == 0 {
		must(fmt.Errorf("no legal name records in embedded font %s", filepath.Base(path)))
	}
	return out.Bytes()
}

func main() {
	check := flag.Bool("check", false, "compare the generated bundle without changing files")
	flag.Parse()
	if flag.NArg() != 0 {
		must(errors.New("unexpected positional arguments"))
	}
	if runtime.Version() != "go1.26.10" {
		must(errors.New("Go 1.26.10 is required"))
	}
	// Both cache archives and extracted module trees must match their downloaded
	// hashes. This verifies cache integrity; go.sum records module provenance.
	goCommand(targets[0], "mod", "verify")
	components := map[string]*component{}
	goRoot := strings.TrimSpace(string(goCommand(targets[0], "env", "GOROOT")))
	goRuntime := newComponent("go@go1.26.10", goRoot, "", "https://go.googlesource.com/go/+/refs/tags/go1.26.10/")
	components[goRuntime.ID] = goRuntime
	for _, name := range []string{"LICENSE", "PATENTS"} {
		goRuntime.Legal[filepath.Join(goRoot, name)] = true
	}
	for _, target := range targets {
		data := goCommand(target, "list", "-mod=readonly", "-p", "2", "-deps", "-json", "./cmd/takedown")
		decoder := json.NewDecoder(bytes.NewReader(data))
		for {
			var pkg goPackage
			err := decoder.Decode(&pkg)
			if err == io.EOF {
				break
			}
			must(err)
			if pkg.Incomplete || pkg.ForTest != "" {
				must(fmt.Errorf("incomplete or test-only package: %s", pkg.ImportPath))
			}
			c := goRuntime
			if pkg.Module != nil {
				m := pkg.Module
				if m.Main {
					continue
				}
				if m.Replace != nil || m.Version == "" || m.Sum == "" {
					must(fmt.Errorf("dependency must have an unreplaced version and checksum: %s", m.Path))
				}
				id := m.Path + "@" + m.Version
				if components[id] == nil {
					components[id] = newComponent(id, m.Dir, m.Sum, "https://pkg.go.dev/"+id)
				}
				c = components[id]
			} else if !pkg.Standard && !pkg.Goroot {
				must(fmt.Errorf("package has no source provenance: %s", pkg.ImportPath))
			}
			c.Targets[target] = true
			if pkg.Dir == "" {
				continue // Compiler-provided package unsafe has no source directory.
			}
			c.IncludeDirs[pkg.Dir] = true
			flags := append(append([]string{}, pkg.CgoCFLAGS...), pkg.CgoCPPFLAGS...)
			for i := 0; i < len(flags); i++ {
				var includeDir string
				if flags[i] == "-I" && i+1 < len(flags) {
					i++
					includeDir = flags[i]
				} else if strings.HasPrefix(flags[i], "-I") {
					includeDir = strings.TrimPrefix(flags[i], "-I")
				}
				if includeDir != "" {
					includeDir = strings.ReplaceAll(includeDir, "${SRCDIR}", pkg.Dir)
					if !filepath.IsAbs(includeDir) {
						includeDir = filepath.Join(pkg.Dir, includeDir)
					}
					c.IncludeDirs[includeDir] = true
				}
			}
			for _, files := range [][]string{pkg.GoFiles, pkg.CgoFiles, pkg.CFiles, pkg.CXXFiles, pkg.MFiles, pkg.HFiles, pkg.FFiles, pkg.SFiles, pkg.SwigFiles, pkg.SwigCXXFiles} {
				for _, name := range files {
					c.Sources[filepath.Join(pkg.Dir, name)] = true
				}
			}
			for _, name := range pkg.EmbedFiles {
				if ext := strings.ToLower(filepath.Ext(name)); ext == ".ttf" || ext == ".otf" {
					c.Fonts[filepath.Join(pkg.Dir, name)] = true
				}
			}
			// The Go distribution has additional licenses above its vendored
			// runtime packages. Do not collect unrelated compiler/tool licenses.
			if c == goRuntime {
				for dir := pkg.Dir; dir != goRoot; dir = filepath.Dir(dir) {
					portablePath(c, dir)
					entries, err := os.ReadDir(dir)
					must(err)
					for _, entry := range entries {
						if !entry.IsDir() && legalName.MatchString(entry.Name()) && !strings.HasSuffix(entry.Name(), ".go") {
							c.Legal[filepath.Join(dir, entry.Name())] = true
						}
					}
				}
			}
		}
	}
	for _, c := range components {
		if c != goRuntime {
			collectLegal(c)
		}
		collectIncludes(c)
	}
	var out bytes.Buffer
	fmt.Fprintln(&out, "# Third-party notices")
	fmt.Fprintln(&out, `
This bundle preserves original upstream license, copyright, patent, attribution,
and notice text. Each component retains its own terms; the application's license
does not relicense these components. This inventory is not legal certification.

Scope: the union of the default, cgo-enabled native runtime dependency graphs for
./cmd/takedown on darwin/amd64, darwin/arm64, linux/amd64, and windows/amd64,
using Go 1.26.10. Dependencies reached only through tests or tools are excluded.
Original standalone legal files are included conservatively from each selected
module, so some files may describe unused subcomponents or test assets. Original
source comment notices and embedded font legal name records are also preserved.
Literal native includes shipped within each component are followed recursively.
Vendored OpenGL/GLFW native sources and platform headers are included conservatively;
their inclusion does not assert that each header is linked on every target.
System libraries, SDKs, compiler support libraries, and packaging tools supplied
by build hosts are outside this Go source inventory. Binary distributors must
also account for any separately bundled native libraries or other added payloads.

Regenerate from the repository root after downloading the pinned modules:

    go mod download
    go run -mod=readonly -p 2 scripts/generate-third-party-notices.go

Check without modifying the bundle:

    go run -mod=readonly -p 2 scripts/generate-third-party-notices.go -check

Generation requires Go 1.26.10 and the pinned module cache. The helper verifies
the module cache with go mod verify, prohibits source downloads while generating,
and queries go list -mod=readonly -p 2 -deps -json without -test for all four targets.
Default architecture feature levels, no GOEXPERIMENT, and GOFIPS140=off are used.
No cross compiler or native linking is needed. Targets, components, files, and
deduplicated source notices use stable ordering; output has no timestamp or local
cache paths. Module sums and SHA-256 hashes identify the original source files.
The same pinned inputs produce the same bytes on macOS, Linux, and Windows.

Method references: https://go.dev/ref/mod#go-mod-verify and
https://go.dev/cmd/go/#hdr-List_packages_or_modules.

## Components`)
	for _, id := range sortedKeys(components) {
		c := components[id]
		fmt.Fprintf(&out, "\n- `%s` (%s)\n", id, strings.Join(sortedKeys(c.Targets), ", "))
	}
	for _, id := range sortedKeys(components) {
		c := components[id]
		fmt.Fprintf(&out, "\n## %s\n\nSource: %s\n\nTargets: %s\n\n", id, c.Source, strings.Join(sortedKeys(c.Targets), ", "))
		if c.Sum != "" {
			fmt.Fprintf(&out, "Module checksum: `%s`\n\n", c.Sum)
		}
		for _, path := range sortedPaths(c.Legal) {
			data := readFile(path)
			fmt.Fprintf(&out, "### %s\n\nSHA-256: `%x`\n\n", portablePath(c, path), sha256.Sum256(data))
			literal(&out, data)
		}
		groups := map[string][]string{}
		for _, path := range sortedPaths(c.Sources) {
			data := readFile(path)
			for _, notice := range sourceNotices(path) {
				provenance := fmt.Sprintf("%s (source SHA-256: `%x`)", portablePath(c, path), sha256.Sum256(data))
				groups[notice] = append(groups[notice], provenance)
			}
		}
		for _, notice := range sortedKeys(groups) {
			fmt.Fprintln(&out, "### Original source notice")
			fmt.Fprintln(&out)
			for _, path := range groups[notice] {
				fmt.Fprintf(&out, "- %s\n", path)
			}
			fmt.Fprintln(&out)
			literal(&out, []byte(notice))
		}
		for _, path := range sortedPaths(c.Fonts) {
			fmt.Fprintf(&out, "### %s (font legal name records)\n\nSource SHA-256: `%x`\n\n", portablePath(c, path), sha256.Sum256(readFile(path)))
			literal(&out, fontNotices(path))
		}
	}
	if *check {
		if !bytes.Equal(readFile(outputFile), out.Bytes()) {
			must(errors.New(outputFile + " is stale; regenerate with go run -mod=readonly -p 2 scripts/generate-third-party-notices.go"))
		}
		fmt.Printf("Verified %s: %d runtime modules plus Go 1.26.10, four native targets\n", outputFile, len(components)-1)
		return
	}
	must(os.WriteFile(outputFile, out.Bytes(), 0644))
	fmt.Printf("Generated %s: %d runtime modules plus Go 1.26.10, four native targets\n", outputFile, len(components)-1)
}
