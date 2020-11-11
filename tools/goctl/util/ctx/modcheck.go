package ctx

import (
	"errors"
	"os"

	"github.com/tal-tech/go-zero/core/jsonx"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
)

// IsGoMod is used to determine whether workDir is a go module project through command `go list -json -m`
func IsGoMod(workDir string) (bool, error) {
	if len(workDir) == 0 {
		return false, errors.New("the work directory is not found")
	}
	if _, err := os.Stat(workDir); err != nil {
		return false, err
	}

	data, err := execx.Run("go list -json -m", workDir)
	if err != nil {
		return false, nil
	}

	var m Module
	err = jsonx.Unmarshal([]byte(data), &m)
	if err != nil {
		return false, err
	}

	return len(m.GoMod) > 0, nil
}
