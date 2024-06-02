package pathx

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/vars"
)

const (
	pkgSep           = "/"
	goModeIdentifier = "go.mod"
)

// JoinPackages calls strings.Join and returns
func JoinPackages(pkgs ...string) string {
	return strings.Join(pkgs, pkgSep)
}

// MkdirIfNotExist makes directories if the input path is not exists
func MkdirIfNotExist(dir string) error {
	if len(dir) == 0 {
		return nil
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}

	return nil
}

// PathFromGoSrc returns the path without slash where has been trim the prefix $GOPATH
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

// FindGoModPath returns the path in project where has file go.mod,
// it returns empty string if there is no go.mod file in project.
func FindGoModPath(dir string) (string, bool) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", false
	}

	absDir = strings.ReplaceAll(absDir, `\`, `/`)
	var rootPath string
	tempPath := absDir
	hasGoMod := false
	for {
		if FileExists(filepath.Join(tempPath, goModeIdentifier)) {
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

// FindProjectPath returns the parent directory where has file go.mod in project
func FindProjectPath(loc string) (string, bool) {
	var dir string
	if strings.IndexByte(loc, '/') == 0 {
		dir = loc
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return "", false
		}

		dir = filepath.Join(wd, loc)
	}

	for {
		if FileExists(filepath.Join(dir, goModeIdentifier)) {
			return dir, true
		}

		dir = filepath.Dir(dir)
		if dir == "/" {
			break
		}
	}

	return "", false
}

func isLink(name string) (bool, error) {
	fi, err := os.Lstat(name)
	if err != nil {
		return false, err
	}

	return fi.Mode()&os.ModeSymlink != 0, nil
}
