package ctx

import (
	"errors"
	"os"

	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
)

// UpdateGoWorkIfExist updates go work if workDir is in a go workspace
func UpdateGoWorkIfExist(workDir string) error {
	if isGoWork, err := isGoWork(workDir); !isGoWork || err != nil {
		return err
	}

	_, err := execx.Run("go work use .", workDir)
	return err
}

// isGoWork detect if the workDir is in a go workspace
func isGoWork(workDir string) (bool, error) {
	if len(workDir) == 0 {
		return false, errors.New("the work directory is not found")
	}
	if _, err := os.Stat(workDir); err != nil {
		return false, err
	}
	goWorkPath, err := execx.Run("go env GOWORK", workDir)
	if err != nil {
		return false, err
	}
	return len(goWorkPath) > 0, nil
}
