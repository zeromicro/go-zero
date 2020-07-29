package util

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func GetFullPackage(pkg string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	pkgPath := path.Join(dir, pkg)
	info, err := os.Stat(pkgPath)
	if err != nil {
		return "", err
	}

	if !info.IsDir() {
		return "", fmt.Errorf("%s is not a directory", pkg)
	}

	gopath := os.Getenv("GOPATH")
	parent := path.Join(gopath, "src")
	pos := strings.Index(pkgPath, parent)
	if pos < 0 {
		return "", fmt.Errorf("%s is not a correct package", pkg)
	}

	// skip slash
	return pkgPath[len(parent)+1:], nil
}
