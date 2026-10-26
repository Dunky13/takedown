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
			o := Options{Platform: platform, Home: "/private/user home", ConfigHome: "/private/user home/.config", Executable: "/private/app $special%/takedown", DataDir: "/private/data & <profiles>", UID: "501", User: "DOMAIN\\operator", PartByPart: true}
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
	o := Options{Platform: "linux", Home: "/private/operator", Executable: "/opt/takedown", DataDir: "/private/profile-a"}
	a, err := Build(o)
	if err != nil {
		t.Fatal(err)
	}
	o.DataDir = "/private/profile-b"
	b, err := Build(o)
	if err != nil {
		t.Fatal(err)
	}
	if a.Name == b.Name || a.Path == b.Path {
		t.Fatal("different profiles share a login task")
	}
	o.DataDir = "/private/profile\nExecStart=/other"
	if _, err = Build(o); err == nil {
		t.Fatal("control-character injection accepted")
	}
}
