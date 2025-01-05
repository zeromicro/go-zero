package ctx

import (
	"errors"
	"os"

	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
)

// IsGoMod is used to determine whether workDir is a go module project through command `go env GOMOD`
func IsGoMod(workDir string) (bool, error) {
	if len(workDir) == 0 {
		return false, errors.New("the work directory is not found")
	}
	if _, err := os.Stat(workDir); err != nil {
		return false, err
	}

	data, err := execx.Run("go env GOMOD", workDir)
	if err != nil || data == "/dev/null" {
		return false, nil
	}

	return true, nil
}
