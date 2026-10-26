package autostart

import (
	"encoding/xml"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestLoginDefinitionsKeepArgumentsAndUserBoundary(t *testing.T) {
	for _, platform := range []string{"darwin", "linux", "windows"} {
		t.Run(platform, func(t *testing.T) {
			root := t.TempDir()
			o := Options{
				Platform:   platform,
				Home:       filepath.Join(root, "user home"),
				ConfigHome: filepath.Join(root, "user home", ".config"),
				Executable: filepath.Join(root, "app $special%", "takedown"),
				DataDir:    filepath.Join(root, "data & <profiles>"),
				UID:        "501", User: "DOMAIN\\operator", PartByPart: true,
			}
			d, err := Build(o)
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(d.Content, "TAKEDOWN_STORAGE_PASSWORD") || strings.Contains(d.Content, "RefreshToken") || strings.Contains(d.Content, "<RunLevel>Highest") {
				t.Fatal("definition includes secrets or privilege elevation")
			}
			if !strings.Contains(d.Content, "watch") || !strings.Contains(d.Content, "--part-by-part") {
				t.Fatal("watch policy missing")
			}
			if platform == "linux" {
				if !strings.Contains(d.Content, "ExecStart=:") || !strings.Contains(d.Content, "%%") || strings.Contains(d.Content, "Restart=always") {
					t.Fatal("systemd argument expansion or retry policy is unsafe")
				}
			} else {
				decoder := xml.NewDecoder(strings.NewReader(d.Content))
				for {
					_, err = decoder.Token()
					if err == io.EOF {
						break
					}
					if err != nil {
						t.Fatal(err)
					}
				}
			}
			if platform == "windows" && (!strings.Contains(d.Content, "InteractiveToken") || !strings.Contains(d.Disable[0].Args[len(d.Disable[0].Args)-1], "Stop-ScheduledTask")) {
				t.Fatal("Windows user boundary or running-task teardown missing")
			}
		})
	}
}

func TestNativeDefinitionsRejectRelativePaths(t *testing.T) {
	root := t.TempDir()
	valid := Options{Home: root, Executable: filepath.Join(root, "takedown"), DataDir: filepath.Join(root, "data"), UID: "501", User: "DOMAIN\\operator"}
	if _, err := Build(valid); err != nil {
		t.Fatalf("native absolute paths rejected: %v", err)
	}
	paths := []string{"relative/path"}
	if runtime.GOOS == "windows" {
		paths = append(paths, "/private/path", `C:relative\path`, `\rooted\path`)
	}
	for _, field := range []string{"home", "executable", "data-dir"} {
		for _, path := range paths {
			t.Run(field+"/"+path, func(t *testing.T) {
				o := valid
				switch field {
				case "home":
					o.Home = path
				case "executable":
					o.Executable = path
				case "data-dir":
					o.DataDir = path
				}
				if _, err := Build(o); err == nil || !strings.Contains(err.Error(), "absolute runtime paths") {
					t.Fatalf("native non-absolute path must be rejected: %v", err)
				}
			})
		}
	}
}

func TestNativePlistParsesAndNoInstallOccurs(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Apple's native plist validator is Darwin-only; XML is checked on all hosts")
	}
	plutil, err := exec.LookPath("plutil")
	if err != nil {
		t.Skip("native plist validator unavailable")
	}
	d, err := Build(Options{Platform: "darwin", Home: t.TempDir(), Executable: "/Applications/with spaces/takedown", DataDir: "/private/data&name", UID: "501"})
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "watch.plist")
	if err = os.WriteFile(path, []byte(d.Content), 0600); err != nil {
		t.Fatal(err)
	}
	if out, err := exec.Command(plutil, "-lint", path).CombinedOutput(); err != nil {
		t.Fatalf("native plist validation: %v: %s", err, out)
	}
	if _, err = os.Stat(d.Path); !os.IsNotExist(err) {
		t.Fatal("rendering installed a login task")
	}
}

func TestRejectControlCharactersAndIsolateProfiles(t *testing.T) {
	root := t.TempDir()
	o := Options{Platform: "linux", Home: filepath.Join(root, "operator"), Executable: filepath.Join(root, "takedown"), DataDir: filepath.Join(root, "profile-a")}
	a, err := Build(o)
	if err != nil {
		t.Fatal(err)
	}
	o.DataDir = filepath.Join(root, "profile-b")
	b, err := Build(o)
	if err != nil {
		t.Fatal(err)
	}
	if a.Name == b.Name || a.Path == b.Path {
		t.Fatal("different profiles share a login task")
	}
	o.DataDir = filepath.Join(root, "profile\nExecStart=/other")
	if _, err = Build(o); err == nil {
		t.Fatal("control-character injection accepted")
	}
}
