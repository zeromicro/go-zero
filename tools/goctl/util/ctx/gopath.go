package ctx

import (
	"errors"
	"go/build"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

// projectFromGoPath is used to find the main module and project file path
// the workDir flag specifies which folder we need to detect based on
// only valid for go mod project
func projectFromGoPath(workDir string) (*ProjectContext, error) {
	if len(workDir) == 0 {
		return nil, errors.New("the work directory is not found")
	}
	if _, err := os.Stat(workDir); err != nil {
		return nil, err
	}

	workDir, err := pathx.ReadLink(workDir)
	if err != nil {
		return nil, err
	}

	buildContext := build.Default
	goPath := buildContext.GOPATH
	goPath, err = pathx.ReadLink(goPath)
	if err != nil {
		return nil, err
	}

	goSrc := filepath.Join(goPath, "src")
	if !pathx.FileExists(goSrc) {
		return nil, errModuleCheck
	}

	wd, err := filepath.Abs(workDir)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(wd, goSrc) {
		return nil, errModuleCheck
	}

	projectName := strings.TrimPrefix(wd, goSrc+string(filepath.Separator))
	return &ProjectContext{
		WorkDir: workDir,
		Name:    projectName,
		Path:    projectName,
		Dir:     filepath.Join(goSrc, projectName),
	}, nil
}
