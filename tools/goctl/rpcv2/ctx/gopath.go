package ctx

import (
	"errors"
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/util"
)

func GOPATH(workDir string) (ProjectContext, error) {
	var ret ProjectContext

	if len(workDir) == 0 {
		return ret, errors.New("the work directory is not found")
	}
	if _, err := os.Stat(workDir); err != nil {
		return ret, err
	}

	buildContext := build.Default
	GOPATH := buildContext.GOPATH
	GOSRC := filepath.Join(GOPATH, "src")
	if !util.FileExists(GOSRC) {
		return ret, fmt.Errorf("the $GOPTAH/src:  %s is not found", GOSRC)
	}

	wd, err := filepath.Abs(workDir)
	if err != nil {
		return ret, err
	}

	if !strings.HasPrefix(wd, GOSRC) {
		return ret, fmt.Errorf("the work directory must in the $GOPATH/src")
	}

	projectName := strings.TrimPrefix(wd, GOSRC+string(filepath.Separator))

	return ProjectContext{
		WorkDir: workDir,
		Name:    projectName,
		Path:    projectName,
		Dir:     filepath.Join(GOSRC, projectName),
	}, nil
}
