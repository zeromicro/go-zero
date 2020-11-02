package cli

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/tal-tech/go-zero/tools/goctl/rpcv2/execx"
	"github.com/tal-tech/go-zero/tools/goctl/rpcv2/generator"
	"github.com/urfave/cli"
)

func Rpc(c *cli.Context) error {
	out := c.String("out")
	IPATH := c.StringSlice("proto_path")
	src := c.Args().First()
	if len(src) == 0 {
		return errors.New("the proto source can not be nil")
	}

	if len(out) == 0 {
		out = filepath.Dir(src)
	}
	g := generator.NewDefaultRpcGenerator()
	return g.Generate(src, out, IPATH)
}

func RpcNew(c *cli.Context) error {
	name := c.Args().First()
	ext := filepath.Ext(name)
	if len(ext) > 0 {
		return fmt.Errorf("unexpected ext: %s", ext)
	}
	protoName := name + ".proto"
	filename := filepath.Join(".", name, protoName)
	src, err := filepath.Abs(filename)
	if err != nil {
		return err
	}

	err = generator.ProtoTmpl(src)
	if err != nil {
		return err
	}

	workDir := filepath.Dir(src)
	_, err = execx.Run("go mod init "+name, workDir)
	if err != nil {
		return err
	}

	g := generator.NewDefaultRpcGenerator()
	return g.Generate(src, filepath.Dir(src), nil)
}

func RpcTemplate(c *cli.Context) error {
	name := c.Args().First()
	if len(name) == 0 {
		name = "greet.proto"
	}
	return generator.ProtoTmpl(name)
}
