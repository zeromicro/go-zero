package cli

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/tal-tech/go-zero/tools/goctl/rpc/generator"
	"github.com/tal-tech/go-zero/tools/goctl/util/name"
	"github.com/urfave/cli"
)

// Rpc is to generate rpc service code from a proto file by specifying a proto file using flag src,
// you can specify a target folder for code generation, when the proto file has import, you can specify
// the import search directory through the proto_path command, for specific usage, please refer to protoc -h
func Rpc(c *cli.Context) error {
	src := c.String("src")
	out := c.String("dir")
	style := c.String("style")
	protoImportPath := c.StringSlice("proto_path")
	if len(src) == 0 {
		return errors.New("missing -src")
	}
	if len(out) == 0 {
		return errors.New("missing -dir")
	}

	namingStyle, valid := name.IsNamingValid(style)
	if !valid {
		return fmt.Errorf("unexpected naming style %s", style)
	}

	g := generator.NewDefaultRpcGenerator(namingStyle)
	return g.Generate(src, out, protoImportPath)
}

// RpcNew is to generate rpc greet service, this greet service can speed
// up your understanding of the zrpc service structure
func RpcNew(c *cli.Context) error {
	rpcname := c.Args().First()
	ext := filepath.Ext(rpcname)
	if len(ext) > 0 {
		return fmt.Errorf("unexpected ext: %s", ext)
	}

	style := c.String("style")
	namingStyle, valid := name.IsNamingValid(style)
	if !valid {
		return fmt.Errorf("expected naming style [lower|camel|snake], but found %s", style)
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

	g := generator.NewDefaultRpcGenerator(namingStyle)
	return g.Generate(src, filepath.Dir(src), nil)
}

func RpcTemplate(c *cli.Context) error {
	protoFile := c.String("o")
	if len(protoFile) == 0 {
		return errors.New("missing -o")
	}

	return generator.ProtoTmpl(protoFile)
}
