package command

import (
	"github.com/urfave/cli"

	"github.com/tal-tech/go-zero/tools/goctl/rpc/ctx"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/goen"
)

func Rpc(c *cli.Context) error {
	rpcCtx := ctx.MustCreateRpcContextFromCli(c)
	generator := gogen.NewDefaultRpcGenerator(rpcCtx)
	rpcCtx.Must(generator.Generate())
	return nil
}
