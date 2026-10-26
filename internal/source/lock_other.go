//go:build !unix && !windows

package source

import "errors"

func lockState(string) (func(), error) {
	return nil, errors.New("durable Data Portability job locking is unsupported on this platform")
}
func replaceState(string, string) error {
	return errors.New("durable Data Portability job storage is unsupported on this platform")
}
