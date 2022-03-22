package cli

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/generator"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	errInvalidGrpcOutput = errors.New("ZRPC: missing --go-grpc_out")
	errInvalidGoOutput   = errors.New("ZRPC: missing --go_out")
	errInvalidZrpcOutput = errors.New("ZRPC: missing zrpc output, please use --zrpc_out to specify the output")
	errInvalidInput      = errors.New("ZRPC: missing source")
	errMultiInput        = errors.New("ZRPC: only one source is expected")
)

// ZRPC generates grpc code directly by protoc and generates
// zrpc code by goctl.
func ZRPC(c *cli.Context) error {
	args := c.Parent().Args()
	protocArgs := removeGoctlFlag(args)
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	source, err := getSourceProto(c.Args(), pwd)
	if err != nil {
		return err
	}

	grpcOutList := c.StringSlice("go-grpc_out")
	goOutList := c.StringSlice("go_out")
	zrpcOut := c.String("zrpc_out")
	style := c.String("style")
	home := c.String("home")
	remote := c.String("remote")
	branch := c.String("branch")
	verbose := c.Bool("verbose")
	if len(grpcOutList) == 0 {
		return errInvalidGrpcOutput
	}
	if len(goOutList) == 0 {
		return errInvalidGoOutput
	}
	goOut := goOutList[len(goOutList)-1]
	grpcOut := grpcOutList[len(grpcOutList)-1]
	if len(goOut) == 0 {
		return errInvalidGrpcOutput
	}
	if len(zrpcOut) == 0 {
		return errInvalidZrpcOutput
	}
	goOutAbs, err := filepath.Abs(goOut)
	if err != nil {
		return err
	}
	grpcOutAbs, err := filepath.Abs(grpcOut)
	if err != nil {
		return err
	}
	err = pathx.MkdirIfNotExist(goOutAbs)
	if err != nil {
		return err
	}
	err = pathx.MkdirIfNotExist(grpcOutAbs)
	if err != nil {
		return err
	}
	if len(remote) > 0 {
		repo, _ := util.CloneIntoGitHome(remote, branch)
		if len(repo) > 0 {
			home = repo
		}
	}

	if len(home) > 0 {
		pathx.RegisterGoctlHome(home)
	}
	if !filepath.IsAbs(zrpcOut) {
		zrpcOut = filepath.Join(pwd, zrpcOut)
	}

	isGooglePlugin := len(grpcOut) > 0
	goOut, err = filepath.Abs(goOut)
	if err != nil {
		return err
	}
	grpcOut, err = filepath.Abs(grpcOut)
	if err != nil {
		return err
	}
	zrpcOut, err = filepath.Abs(zrpcOut)
	if err != nil {
		return err
	}

	var ctx generator.ZRpcContext
	ctx.Src = source
	ctx.GoOutput = goOut
	ctx.GrpcOutput = grpcOut
	ctx.IsGooglePlugin = isGooglePlugin
	ctx.Output = zrpcOut
	ctx.ProtocCmd = strings.Join(protocArgs, " ")
	g := generator.NewGenerator(style, verbose)
	return g.Generate(&ctx)
}

func removeGoctlFlag(args []string) []string {
	var ret []string
	var step int
	for step < len(args) {
		arg := args[step]
		switch {
		case arg == "--style", arg == "--home", arg == "--zrpc_out", arg == "--verbose", arg == "-v":
			step += 2
			continue
		case strings.HasPrefix(arg, "--style="),
			strings.HasPrefix(arg, "--home="),
			strings.HasPrefix(arg, "--verbose="),
			strings.HasPrefix(arg, "-v="),
			strings.HasPrefix(arg, "--zrpc_out="):
			step += 1
			continue
		}
		step += 1
		ret = append(ret, arg)
	}

	return ret
}

func getSourceProto(args []string, pwd string) (string, error) {
	var source []string
	for _, p := range args {
		if strings.HasSuffix(p, ".proto") {
			source = append(source, p)
		}
	}

	switch len(source) {
	case 0:
		return "", errInvalidInput
	case 1:
		isAbs := filepath.IsAbs(source[0])
		if isAbs {
			return source[0], nil
		}

		abs := filepath.Join(pwd, source[0])
		return abs, nil
	default:
		return "", errMultiInput
	}
}

func removePluginFlag(goOut string) string {
	goOut = strings.ReplaceAll(goOut, "plugins=", "")
	index := strings.LastIndex(goOut, ":")
	if index < 0 {
		return goOut
	}
	return goOut[index+1:]
}
