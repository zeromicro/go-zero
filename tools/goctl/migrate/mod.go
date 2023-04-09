package migrate

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/ctx"
)

const (
	goZeroMod = "github.com/zeromicro/go-zero"
	adminTool = "github.com/suyuan32/simple-admin-tools"
)

var errInvalidGoMod = errors.New("it's only working for go module")

func editMod(zeroVersion, toolVersion string, verbose bool) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	isGoMod, _ := ctx.IsGoMod(wd)
	if !isGoMod {
		return nil
	}

	mod := fmt.Sprintf("%s@%s", goZeroMod, zeroVersion)

	err = addRequire(mod, verbose)
	if err != nil {
		return err
	}

	// add replace
	mod = fmt.Sprintf("%s@%s=%s@%s", goZeroMod, zeroVersion, adminTool, toolVersion)

	err = addReplace(mod, verbose)
	if err != nil {
		return err
	}

	return nil
}

func addRequire(mod string, verbose bool) error {
	if verbose {
		console.Info("adding require %s ...", mod)
		time.Sleep(200 * time.Millisecond)
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	isGoMod, _ := ctx.IsGoMod(wd)
	if !isGoMod {
		return errInvalidGoMod
	}

	_, err = execx.Run("go mod edit -require "+mod, wd)
	return err
}

func addReplace(mod string, verbose bool) error {
	if verbose {
		console.Info("adding replace %s ...", mod)
		time.Sleep(200 * time.Millisecond)
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	isGoMod, _ := ctx.IsGoMod(wd)
	if !isGoMod {
		return errInvalidGoMod
	}

	_, err = execx.Run("go mod edit -replace "+mod, wd)
	return err
}

func tidy(verbose bool) error {
	if verbose {
		console.Info("go mod tidy ...")
		time.Sleep(200 * time.Millisecond)
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	isGoMod, _ := ctx.IsGoMod(wd)
	if !isGoMod {
		return nil
	}

	_, err = execx.Run("go mod tidy", wd)
	return err
}
