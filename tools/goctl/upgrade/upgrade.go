package upgrade

import (
	"fmt"
	"github.com/tal-tech/go-zero/tools/goctl/internal/errorx"

	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
	"github.com/urfave/cli"
)

// Upgrade gets the latest goctl by
// go get -u github.com/tal-tech/go-zero/tools/goctl
func Upgrade(_ *cli.Context) error {
	info, err := execx.Run("GO111MODULE=on GOPROXY=https://goproxy.cn/,direct go get -u github.com/tal-tech/go-zero/tools/goctl", "")
	errorx.Must(err)

	fmt.Print(info)
	return nil
}
