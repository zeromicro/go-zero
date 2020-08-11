package util

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/vars"
)

const pkgSep = "/"

func JoinPackages(pkgs ...string) string {
	return strings.Join(pkgs, pkgSep)
}

func MkdirIfNotExist(dir string) error {
	if len(dir) == 0 {
		return nil
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}

	return nil
}

func PathFromGoSrc() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	gopath := os.Getenv("GOPATH")
	parent := path.Join(gopath, "src", vars.ProjectName)
	pos := strings.Index(dir, parent)
	if pos < 0 {
		return "", fmt.Errorf("%s is not in GOPATH", dir)
	}

	// skip slash
	return dir[len(parent)+1:], nil
}
