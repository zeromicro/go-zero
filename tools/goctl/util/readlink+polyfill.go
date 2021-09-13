//go:build windows
// +build windows

package util

func ReadLink(name string) (string, error) {
	return name, nil
}
