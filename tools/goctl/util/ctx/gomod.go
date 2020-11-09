package ctx

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/tal-tech/go-zero/core/jsonx"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
)

type Module struct {
	Path      string
	Main      bool
	Dir       string
	GoMod     string
	GoVersion string
}

// projectFromGoMod is used to find the go module and project file path
// the workDir flag specifies which folder we need to detect based on
// only valid for go mod project
func projectFromGoMod(workDir string) (*ProjectContext, error) {
	if len(workDir) == 0 {
		return nil, errors.New("the work directory is not found")
	}
	if _, err := os.Stat(workDir); err != nil {
		return nil, err
	}

	data, err := execx.Run("go list -json -m", workDir)
	if err != nil {
		return nil, err
	}

	var m Module
	err = jsonx.Unmarshal([]byte(data), &m)
	if err != nil {
		return nil, err
	}
	var ret ProjectContext
	ret.WorkDir = workDir
	ret.Name = filepath.Base(m.Dir)
	ret.Dir = m.Dir
	ret.Path = m.Path
	return &ret, nil
}
