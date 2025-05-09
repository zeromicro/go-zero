//go:build linux || darwin || freebsd

package pathx

import (
	"os"
	"path/filepath"
)

// ReadLink returns the destination of the named symbolic link recursively.
func ReadLink(name string) (string, error) {
	name, err := filepath.Abs(name)
	if err != nil {
		return "", err
	}

	if _, err := os.Lstat(name); err != nil {
		return name, nil
	}

	// uncheck condition: ignore file path /var, maybe be temporary file path
	if name == "/" || name == "/var" {
		return name, nil
	}

	isLink, err := isLink(name)
	if err != nil {
		return "", err
	}

	if !isLink {
		dir, base := filepath.Split(name)
		dir = filepath.Clean(dir)
		dir, err := ReadLink(dir)
		if err != nil {
			return "", err
		}

		return filepath.Join(dir, base), nil
	}

	link, err := os.Readlink(name)
	if err != nil {
		return "", err
	}

	dir, base := filepath.Split(link)
	dir = filepath.Dir(dir)
	dir, err = ReadLink(dir)
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, base), nil
}
