package config

import (
	"path/filepath"
	"strings"
)

// AbsoluteRuntimePaths fixes paths to the caller's current working directory
// before they are persisted for a login task with a different working directory.
// Remote destination URLs and empty optional paths retain their exact values.
func (s Settings) AbsoluteRuntimePaths() (Settings, error) {
	for _, path := range []*string{&s.OAuthClientFile, &s.PrivateKeyPath, &s.KnownHostsPath} {
		if *path == "" {
			continue
		}
		absolute, err := filepath.Abs(*path)
		if err != nil {
			return Settings{}, err
		}
		*path = absolute
	}
	if s.Destination != "" && !strings.Contains(s.Destination, "://") {
		absolute, err := filepath.Abs(s.Destination)
		if err != nil {
			return Settings{}, err
		}
		s.Destination = absolute
	}
	return s, nil
}
