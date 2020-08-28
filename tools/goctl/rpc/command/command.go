package command

import (
	"github.com/urfave/cli"

	"github.com/tal-tech/go-zero/tools/goctl/rpc/ctx"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/gen"
)

func Rpc(c *cli.Context) error {
	rpcCtx := ctx.MustCreateRpcContextFromCli(c)
	generator := gogen.NewDefaultRpcGenerator(rpcCtx)
	rpcCtx.Must(generator.Generate())
	return nil
}

func RpcTemplate(c *cli.Context) error {
	out := c.String("out")
	idea := c.Bool("idea")
	generator := gogen.NewRpcTemplate(out, idea)
	generator.MustGenerate()
	return nil
}
