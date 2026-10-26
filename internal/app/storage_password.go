package app

import (
	"net/url"

	"github.com/Dunky13/takedown/internal/config"
)

func (a *App) storagePassword(destination, explicit string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}
	u, err := url.Parse(destination)
	if err != nil || (u.Scheme != "smb" && u.Scheme != "sftp") {
		return "", nil
	}
	if u.Scheme == "sftp" && a.Settings.PrivateKeyPath != "" {
		return "", nil
	}
	return config.StoragePassword(a.DataDir, destination)
}
