package ctx

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/tal-tech/go-zero/core/jsonx"
	"github.com/tal-tech/go-zero/tools/goctl/rpcv2/execx"
)

type (
	// go list -json -m
	Module struct {
		Path      string
		Main      bool
		Dir       string
		GoMod     string
		GoVersion string
	}
)

func GOMOD(workDir string) (ProjectContext, error) {
	var ret ProjectContext

	if len(workDir) == 0 {
		return ret, errors.New("the work directory is not found")
	}
	if _, err := os.Stat(workDir); err != nil {
		return ret, err
	}

	data, err := execx.Run("go list -json -m", workDir)
	if err != nil {
		return ret, err
	}

	// go mod
	var m Module
	err = jsonx.Unmarshal([]byte(data), &m)
	if err != nil {
		return ret, err
	}

	ret.WorkDir = workDir
	ret.Name = filepath.Base(m.Dir)
	ret.Dir = m.Dir
	ret.Path = m.Path
	return ret, nil
}
