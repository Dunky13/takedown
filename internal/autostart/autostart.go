// Package autostart installs a watch in the current user's login session.
// It never installs a system-wide service or embeds passwords/tokens.
package autostart

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Options struct {
	Platform, Home, ConfigHome, Executable, DataDir, User, UID string
	PartByPart                                                 bool
}
type Command struct {
	Program string
	Args    []string
}
type Definition struct {
	Name, Path, Content string
	Enable, Disable     []Command
	Query               Command
}

func Build(o Options) (Definition, error) {
	if o.Platform == "" {
		o.Platform = runtime.GOOS
	}
	for _, s := range []string{o.Home, o.ConfigHome, o.Executable, o.DataDir, o.User, o.UID} {
		if strings.ContainsAny(s, "\x00\r\n") {
			return Definition{}, errors.New("login watch paths and user identifiers cannot contain control characters")
		}
	}
	if o.Executable == "" || o.DataDir == "" || o.Home == "" {
		return Definition{}, errors.New("executable, data directory and home directory are required")
	}
	// On the running OS, refuse relative paths that would change meaning in a
	// login task's working directory. Foreign plans are render-only for review.
	if o.Platform == runtime.GOOS && (!filepath.IsAbs(o.Executable) || !filepath.IsAbs(o.DataDir) || !filepath.IsAbs(o.Home)) {
		return Definition{}, errors.New("login watch requires absolute runtime paths")
	}
	hash := sha256.Sum256([]byte(o.DataDir))
	name := fmt.Sprintf("takedown-%x", hash[:8])
	args := []string{"watch", "--data-dir", o.DataDir}
	if o.PartByPart {
		args = append(args, "--part-by-part")
	}
	d := Definition{Name: name}
	switch o.Platform {
	case "linux":
		if o.ConfigHome == "" {
			o.ConfigHome = filepath.Join(o.Home, ".config")
		}
		d.Path = filepath.Join(o.ConfigHome, "systemd", "user", name+".service")
		argv := append([]string{o.Executable}, args...)
		for i := range argv {
			argv[i] = systemdArgument(argv[i])
		}
		d.Content = "[Unit]\nDescription=takedown archive watch\nAfter=network-online.target\n\n[Service]\nType=simple\nExecStart=:" + strings.Join(argv, " ") + "\nRestart=on-abnormal\nRestartSec=60\nUMask=0077\n\n[Install]\nWantedBy=default.target\n"
		d.Enable = []Command{{"systemctl", []string{"--user", "daemon-reload"}}, {"systemctl", []string{"--user", "enable", "--now", name + ".service"}}}
		d.Disable = []Command{{"systemctl", []string{"--user", "disable", "--now", name + ".service"}}}
		d.Query = Command{"systemctl", []string{"--user", "status", "--no-pager", name + ".service"}}
	case "darwin":
		if _, err := strconv.ParseUint(o.UID, 10, 32); err != nil {
			return d, errors.New("macOS user ID is required")
		}
		label := "io.github.dunky13." + name
		d.Name = label
		d.Path = filepath.Join(o.Home, "Library", "LaunchAgents", label+".plist")
		var b strings.Builder
		b.WriteString(xml.Header + "<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n<plist version=\"1.0\"><dict><key>Label</key><string>" + xmlText(label) + "</string><key>ProgramArguments</key><array>")
		for _, arg := range append([]string{o.Executable}, args...) {
			b.WriteString("<string>" + xmlText(arg) + "</string>")
		}
		b.WriteString("</array><key>RunAtLoad</key><true/><key>KeepAlive</key><dict><key>Crashed</key><true/></dict><key>ThrottleInterval</key><integer>60</integer><key>StandardOutPath</key><string>" + xmlText(filepath.Join(o.DataDir, "watch.log")) + "</string><key>StandardErrorPath</key><string>" + xmlText(filepath.Join(o.DataDir, "watch.log")) + "</string></dict></plist>\n")
		d.Content = b.String()
		domain := "gui/" + o.UID
		d.Enable = []Command{{"launchctl", []string{"bootstrap", domain, d.Path}}}
		d.Disable = []Command{{"launchctl", []string{"bootout", domain + "/" + label}}}
		d.Query = Command{"launchctl", []string{"print", domain + "/" + label}}
	case "windows":
		if o.User == "" {
			return d, errors.New("Windows login user is required")
		}
		d.Path = filepath.Join(o.DataDir, name+".task.xml")
		quoted := make([]string, len(args))
		for i, arg := range args {
			quoted[i] = windowsArgument(arg)
		}
		d.Content = xml.Header + "<Task version=\"1.2\" xmlns=\"http://schemas.microsoft.com/windows/2004/02/mit/task\"><Triggers><LogonTrigger><Enabled>true</Enabled><UserId>" + xmlText(o.User) + "</UserId></LogonTrigger></Triggers><Principals><Principal id=\"Author\"><UserId>" + xmlText(o.User) + "</UserId><LogonType>InteractiveToken</LogonType><RunLevel>LeastPrivilege</RunLevel></Principal></Principals><Settings><MultipleInstancesPolicy>IgnoreNew</MultipleInstancesPolicy><DisallowStartIfOnBatteries>false</DisallowStartIfOnBatteries><StopIfGoingOnBatteries>false</StopIfGoingOnBatteries><ExecutionTimeLimit>PT0S</ExecutionTimeLimit></Settings><Actions Context=\"Author\"><Exec><Command>" + xmlText(o.Executable) + "</Command><Arguments>" + xmlText(strings.Join(quoted, " ")) + "</Arguments></Exec></Actions></Task>\n"
		d.Enable = []Command{{"schtasks.exe", []string{"/Create", "/TN", name, "/XML", d.Path, "/F"}}, {"schtasks.exe", []string{"/Run", "/TN", name}}}
		// Task Scheduler deletion alone leaves an existing process running. The
		// typed cmdlets stop this exact root task if needed, then unregister it.
		script := "$ErrorActionPreference='Stop'; try { $task=Get-ScheduledTask -TaskPath '\\' -TaskName '" + name + "' } catch { if($_.CategoryInfo.Category -eq 'ObjectNotFound'){exit 0}; throw }; if($task.State -eq 'Running'){Stop-ScheduledTask -InputObject $task}; Unregister-ScheduledTask -InputObject $task -Confirm:$false"
		d.Disable = []Command{{"powershell.exe", []string{"-NoProfile", "-NonInteractive", "-Command", script}}}
		d.Query = Command{"schtasks.exe", []string{"/Query", "/TN", name, "/V", "/FO", "LIST"}}
	default:
		return d, errors.New("login watch supports Windows, macOS and Linux")
	}
	return d, nil
}

func xmlText(s string) string {
	var b bytes.Buffer
	_ = xml.EscapeText(&b, []byte(s))
	return b.String()
}
func systemdArgument(s string) string { s = strings.ReplaceAll(s, "%", "%%"); return strconv.Quote(s) }
func windowsArgument(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	slashes := 0
	for _, r := range s {
		if r == '\\' {
			slashes++
			continue
		}
		if r == '"' {
			b.WriteString(strings.Repeat("\\", slashes*2+1))
			b.WriteRune(r)
		} else {
			b.WriteString(strings.Repeat("\\", slashes))
			b.WriteRune(r)
		}
		slashes = 0
	}
	b.WriteString(strings.Repeat("\\", slashes*2))
	b.WriteByte('"')
	return b.String()
}

func Run(ctx context.Context, c Command) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, c.Program, c.Args...).CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("%s failed: %w; %s", c.Program, err, strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

func Install(ctx context.Context, o Options) (Definition, error) {
	if o.Platform != "" && o.Platform != runtime.GOOS {
		return Definition{}, errors.New("foreign login definitions are render-only")
	}
	d, err := Build(o)
	if err != nil {
		return d, err
	}
	if info, e := os.Stat(o.Executable); e != nil || !info.Mode().IsRegular() {
		return d, errors.New("login watch executable must be an existing regular file")
	}
	if err = os.MkdirAll(filepath.Dir(d.Path), 0700); err != nil {
		return d, err
	}
	if err = os.MkdirAll(o.DataDir, 0700); err != nil {
		return d, err
	}
	if info, e := os.Lstat(d.Path); e == nil && !info.Mode().IsRegular() {
		return d, errors.New("existing watch definition must be a regular file, not a symlink")
	} else if e != nil && !errors.Is(e, os.ErrNotExist) {
		return d, e
	}
	if existing, e := os.ReadFile(d.Path); e == nil {
		if string(existing) != d.Content {
			return d, errors.New("a different watch definition already exists; remove it before changing installation options")
		}
	} else if !errors.Is(e, os.ErrNotExist) {
		return d, e
	} else {
		f, e := os.OpenFile(d.Path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if e != nil {
			return d, e
		}
		_, e = f.WriteString(d.Content)
		e = errors.Join(e, f.Sync(), f.Close())
		if e != nil {
			return d, e
		}
	}
	if runtime.GOOS == "darwin" {
		loaded, err := launchdLoaded(ctx, d)
		if err != nil {
			return d, err
		}
		if loaded {
			return d, nil
		}
	}
	for _, c := range d.Enable {
		if _, err = Run(ctx, c); err != nil {
			return d, err
		}
	}
	return d, nil
}

func Remove(ctx context.Context, o Options) error {
	if o.Platform != "" && o.Platform != runtime.GOOS {
		return errors.New("foreign login definitions are render-only")
	}
	d, err := Build(o)
	if err != nil {
		return err
	}
	_, statErr := os.Lstat(d.Path)
	if statErr != nil && !errors.Is(statErr, os.ErrNotExist) {
		return statErr
	}
	stop := true
	if runtime.GOOS == "darwin" {
		stop, err = launchdLoaded(ctx, d)
	} else if runtime.GOOS == "linux" {
		stop, err = systemdLoaded(ctx, d)
	}
	if err != nil {
		return err
	}
	if stop {
		for _, c := range d.Disable {
			if _, err = Run(ctx, c); err != nil {
				if runtime.GOOS == "darwin" && launchdMissing(err) {
					// It may exit between the presence query and bootout.
					continue
				}
				return err
			}
		}
	}
	if err = os.Remove(d.Path); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

// launchdLoaded distinguishes the documented missing-service status from
// permission/connection failures. File absence alone does not stop a loaded job.
func launchdLoaded(ctx context.Context, d Definition) (bool, error) {
	_, err := Run(ctx, d.Query)
	if err == nil {
		return true, nil
	}
	if launchdMissing(err) {
		return false, nil
	}
	return false, err
}

func launchdMissing(err error) bool {
	var exit *exec.ExitError
	return errors.As(err, &exit) && exit.ExitCode() == 113
}

func systemdLoaded(ctx context.Context, d Definition) (bool, error) {
	out, err := Run(ctx, Command{"systemctl", []string{"--user", "show", "--property=LoadState", "--property=ActiveState", d.Name + ".service"}})
	if err != nil {
		return false, err
	}
	values := make(map[string]string)
	for line := range strings.SplitSeq(strings.TrimSpace(out), "\n") {
		key, value, ok := strings.Cut(line, "=")
		if ok {
			values[key] = value
		}
	}
	if values["LoadState"] == "" || values["ActiveState"] == "" {
		return false, errors.New("systemd returned an unrecognized watch status; definition retained")
	}
	return values["LoadState"] != "not-found" || (values["ActiveState"] != "inactive" && values["ActiveState"] != "failed"), nil
}
