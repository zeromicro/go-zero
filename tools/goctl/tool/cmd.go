package tool

import (
	"github.com/zeromicro/go-zero/tools/goctl/internal/cobrax"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web"
)

var (
	portFlag int

	// Cmd describes a rpc command.
	Cmd = cobrax.NewCommand("tool")
)

func init() {
	Cmd.AddCommand(web.Cmd)
}
