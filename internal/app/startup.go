package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/Dunky13/takedown/internal/auth"
	"github.com/Dunky13/takedown/internal/autostart"
	"github.com/Dunky13/takedown/internal/config"
	"github.com/Dunky13/takedown/internal/source"
)

// ConfigureStartup manages only this data directory's user-login watch. Actual
// installation is explicit and starts a watch, never OAuth or export creation.
func (a *App) ConfigureStartup(ctx context.Context, action string, partByPart, dryRun bool) (string, error) {
	if err := a.beginJob(); err != nil {
		return "", err
	}
	defer a.endJob()
	dir, err := filepath.Abs(a.DataDir)
	if err != nil {
		return "", err
	}
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return "", err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configHome, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	identity, err := user.Current()
	if err != nil {
		return "", errors.New("cannot determine the current login user")
	}
	opts := autostart.Options{Home: home, ConfigHome: configHome, Executable: exe, DataDir: dir, User: identity.Username, UID: identity.Uid, PartByPart: partByPart}
	d, err := autostart.Build(opts)
	if err != nil {
		return "", err
	}
	if dryRun {
		return fmt.Sprintf("Definition: %s\n%s", d.Path, d.Content), nil
	}
	switch action {
	case "install":
		cfg, err := a.Settings.AbsoluteRuntimePaths()
		if err != nil {
			return "", err
		}
		scopes, scopeErr := source.Scopes(cfg.EffectiveSource(), cfg.PortabilityResources)
		if scopeErr != nil {
			return "", scopeErr
		}
		m, authErr := auth.New(auth.Options{DataDir: dir, ClientFile: cfg.OAuthClientFile, Email: cfg.Email, Scopes: scopes})
		if authErr != nil {
			return "", authErr
		}
		status, authErr := m.InspectAuthorization(ctx)
		_ = m.Close()
		if authErr != nil {
			return "", authErr
		}
		if !status.TokenStored || !status.RefreshAvailable {
			return "", errors.New("authorize explicitly and save a refresh grant before installing a background watch")
		}
		if strings.HasPrefix(cfg.Destination, "smb://") || (strings.HasPrefix(cfg.Destination, "sftp://") && cfg.PrivateKeyPath == "") {
			password, credentialErr := config.StoragePassword(dir, cfg.Destination)
			if credentialErr != nil {
				return "", credentialErr
			}
			if password == "" {
				return "", errors.New("save network authorization with storage auth, or configure an SFTP key/OS-mounted share, before installing a background watch")
			}
		}
		if err = config.Save(dir, cfg); err != nil {
			return "", err
		}
		if _, err = autostart.Install(ctx, opts); err != nil {
			return "", err
		}
		return "User-login watch registered. An already loaded watch keeps its current run; remove and reinstall to restart it with changed preferences. It needs the user session and unlocked keyring; it never opens sign-in or creates exports. Definition: " + d.Path, nil
	case "remove":
		if err = autostart.Remove(ctx, opts); err != nil {
			return "", err
		}
		return "User-login watch removed. Archives and authorization remain.", nil
	case "status":
		if _, err = os.Lstat(d.Path); errors.Is(err, os.ErrNotExist) {
			return "No user-login watch is installed for this data directory.", nil
		} else if err != nil {
			return "", err
		}
		return autostart.Run(ctx, d.Query)
	default:
		return "", errors.New("unknown startup action")
	}
}
