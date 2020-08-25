package command

import (
	"github.com/urfave/cli"

	"github.com/tal-tech/go-zero/tools/goctl/rpc/ctx"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/goen"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/protoc"
)

var (
	tips = `warning: this operation will overwrite the protoc-gen-plugin in gopath

protoc-gen-go: switch to %s`
)

func Rpc(c *cli.Context) error {
	rpcCtx := ctx.MustCreateRpcContext(c)
	plugin, err := protoc.Prepare()
	rpcCtx.Must(err)
	rpcCtx.Info(tips, plugin)
	generator := gogen.NewDefaultRpcGenerator(rpcCtx)
	rpcCtx.Must(generator.Generate())
	return nil
}
