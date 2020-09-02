package util

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/vars"
)

const (
	pkgSep           = "/"
	goModeIdentifier = "go.mod"
)

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

func FindGoModPath(dir string) (string, bool) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", false
	}

	absDir = strings.ReplaceAll(absDir, `\`, `/`)
	var rootPath string
	var tempPath = absDir
	var hasGoMod = false
	for {
		if FileExists(filepath.Join(tempPath, goModeIdentifier)) {
			tempPath = filepath.Dir(tempPath)
			rootPath = strings.TrimPrefix(absDir[len(tempPath):], "/")
			hasGoMod = true
			break
		}

		if tempPath == filepath.Dir(tempPath) {
			break
		}

		tempPath = filepath.Dir(tempPath)
		if tempPath == string(filepath.Separator) {
			break
		}
	}
	if hasGoMod {
		return rootPath, true
	}
	return "", false
}
