package cli

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/urfave/cli"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/generator"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

// RPCNew is to generate rpc greet service, this greet service can speed
// up your understanding of the zrpc service structure
func RPCNew(c *cli.Context) error {
	if c.NArg() == 0 {
		cli.ShowCommandHelpAndExit(c, "new", 1)
	}

	rpcname := c.Args().First()
	ext := filepath.Ext(rpcname)
	if len(ext) > 0 {
		return fmt.Errorf("unexpected ext: %s", ext)
	}
	style := c.String("style")
	home := c.String("home")
	remote := c.String("remote")
	branch := c.String("branch")
	verbose := c.Bool("verbose")
	if len(remote) > 0 {
		repo, _ := util.CloneIntoGitHome(remote, branch)
		if len(repo) > 0 {
			home = repo
		}
	}
	if len(home) > 0 {
		pathx.RegisterGoctlHome(home)
	}

	protoName := rpcname + ".proto"
	filename := filepath.Join(".", rpcname, protoName)
	src, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	err = generator.ProtoTmpl(src)
	if err != nil {
		return err
	}

	var ctx generator.ZRpcContext
	ctx.Src = src
	ctx.GoOutput = filepath.Dir(src)
	ctx.GrpcOutput = filepath.Dir(src)
	ctx.IsGooglePlugin = true
	ctx.Output = filepath.Dir(src)
	ctx.ProtocCmd = fmt.Sprintf("protoc -I=%s %s --go_out=%s --go-grpc_out=%s", filepath.Dir(src), filepath.Base(src), filepath.Dir(src), filepath.Dir(src))
	g := generator.NewGenerator(style, verbose)
	return g.Generate(&ctx)
}

// RPCTemplate is the entry for generate rpc template
func RPCTemplate(c *cli.Context) error {
	console.Warning("deprecated: goctl rpc template -o is deprecated and will be removed in the future, use goctl rpc -o instead")

	if c.NumFlags() == 0 {
		cli.ShowCommandHelpAndExit(c, "template", 1)
	}

	protoFile := c.String("o")
	home := c.String("home")
	remote := c.String("remote")
	branch := c.String("branch")
	if len(remote) > 0 {
		repo, _ := util.CloneIntoGitHome(remote, branch)
		if len(repo) > 0 {
			home = repo
		}
	}
	if len(home) > 0 {
		pathx.RegisterGoctlHome(home)
	}

	if len(protoFile) == 0 {
		return errors.New("missing -o")
	}

	return generator.ProtoTmpl(protoFile)
}
