package upgrade

import (
	"fmt"
	"runtime"

	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
	"github.com/urfave/cli"
)

// Upgrade gets the latest goctl by
// go get -u github.com/tal-tech/go-zero/tools/goctl
func Upgrade(_ *cli.Context) error {
	cmd := `GO111MODULE=on GOPROXY=https://goproxy.cn/,direct go install github.com/tal-tech/go-zero/tools/goctl@cli`
	if runtime.GOOS == "windows" {
		cmd = `set GOPROXY=https://goproxy.cn,direct && go install github.com/tal-tech/go-zero/tools/goctl@cli`
	}
	info, err := execx.Run(cmd, "")
	if err != nil {
		return err
	}

	fmt.Print(info)
	return nil
}
