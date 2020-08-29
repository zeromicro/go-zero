package command

import (
	"github.com/tal-tech/go-zero/tools/goctl/rpc/ctx"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/gen"
	"github.com/urfave/cli"
)

func Rpc(c *cli.Context) error {
	rpcCtx := ctx.MustCreateRpcContextFromCli(c)
	generator := gen.NewDefaultRpcGenerator(rpcCtx)
	rpcCtx.Must(generator.Generate())
	return nil
}

func RpcTemplate(c *cli.Context) error {
	out := c.String("out")
	idea := c.Bool("idea")
	generator := gen.NewRpcTemplate(out, idea)
	generator.MustGenerate()
	return nil
}
