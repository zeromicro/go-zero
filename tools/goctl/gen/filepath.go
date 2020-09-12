package gen

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/tal-tech/go-zero/tools/goctl/util"
)

func getFilePath(file string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	projPath, ok := util.FindGoModPath(filepath.Join(wd, file))
	if !ok {
		projPath, err = util.PathFromGoSrc()
		if err != nil {
			return "", errors.New("no go.mod found, or not in GOPATH")
		}
	}

	return projPath, nil
}
