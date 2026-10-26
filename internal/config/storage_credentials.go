package config

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/zalando/go-keyring"
)

const storageService = "takedown.storage"

func storageAccount(dataDir, destination string) (string, error) {
	u, err := url.Parse(destination)
	if err != nil || (u.Scheme != "smb" && u.Scheme != "sftp") || u.User == nil || u.User.Username() == "" || u.Host == "" {
		return "", errors.New("storage authorization requires a valid smb:// or sftp:// destination with a username")
	}
	if _, hasPassword := u.User.Password(); hasPassword || u.RawQuery != "" || u.Fragment != "" {
		return "", errors.New("storage destination cannot contain passwords, query parameters or fragments")
	}
	dir, err := filepath.Abs(dataDir)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256([]byte(dir + "\x00" + destination))
	return fmt.Sprintf("%x", digest), nil
}

func SaveStoragePassword(dataDir, destination, password string) error {
	account, err := storageAccount(dataDir, destination)
	if err != nil {
		return err
	}
	if password == "" {
		return errors.New("storage password is empty")
	}
	if err = keyring.Set(storageService, account, password); err != nil {
		return errors.New("OS keyring unavailable; storage password was not saved")
	}
	return nil
}
func StoragePassword(dataDir, destination string) (string, error) {
	account, err := storageAccount(dataDir, destination)
	if err != nil {
		return "", err
	}
	password, err := keyring.Get(storageService, account)
	if errors.Is(err, keyring.ErrNotFound) {
		return "", nil
	}
	if err != nil {
		return "", errors.New("OS keyring unavailable; unlock/configure it before reading the saved storage password")
	}
	return password, nil
}
func ForgetStoragePassword(dataDir, destination string) error {
	account, err := storageAccount(dataDir, destination)
	if err != nil {
		return err
	}
	err = keyring.Delete(storageService, account)
	if err != nil && !errors.Is(err, keyring.ErrNotFound) {
		return errors.New("OS keyring unavailable; storage password could not be removed")
	}
	return nil
}
