package ctx

import (
	"errors"
	"go/build"
	"os"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/util"
)

func projectFromGoPath(workDir string) (*ProjectContext, error) {
	if len(workDir) == 0 {
		return nil, errors.New("the work directory is not found")
	}
	if _, err := os.Stat(workDir); err != nil {
		return nil, err
	}

	buildContext := build.Default
	goPath := buildContext.GOPATH
	goSrc := filepath.Join(goPath, "src")
	if !util.FileExists(goSrc) {
		return nil, moduleCheckErr
	}

	wd, err := filepath.Abs(workDir)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(wd, goSrc) {
		return nil, moduleCheckErr
	}

	projectName := strings.TrimPrefix(wd, goSrc+string(filepath.Separator))
	return &ProjectContext{
		WorkDir: workDir,
		Name:    projectName,
		Path:    projectName,
		Dir:     filepath.Join(goSrc, projectName),
	}, nil
}
