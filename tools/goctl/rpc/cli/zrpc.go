package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/emicklei/proto"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/generator"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/pathx"
	"github.com/urfave/cli"
)

var (
	errInvalidGrpcOutput = errors.New("ZRPC: missing grpc output")
	errInvalidZrpcOutput = errors.New("ZRPC: missing zrpc output, please use --zrpc_out to specify the output")
	errInvalidInput      = errors.New("ZRPC: missing source")
	errMultiInput        = errors.New("ZRPC: only one source is expected")
)

const (
	optImport         = "import"
	optSourceRelative = "source_relative"
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
	src := filepath.Dir(source)
	goPackage, protoPkg, err := getGoPackage(source)
	if err != nil {
		return err
	}

	grpcOut := c.String("go-grpc_out")
	goOut := c.String("go_out")
	goOpt := c.String("go_opt")
	grpcOpt := c.String("go-grpc_opt")
	zrpcOut := c.String("zrpc_out")
	style := c.String("style")
	home := c.String("home")
	remote := c.String("remote")
	if len(remote) > 0 {
		repo, _ := util.CloneIntoGitHome(remote)
		if len(repo) > 0 {
			home = repo
		}
	}

	if len(home) > 0 {
		pathx.RegisterGoctlHome(home)
	}
	if len(goOut) == 0 {
		return errInvalidGrpcOutput
	}
	if len(zrpcOut) == 0 {
		return errInvalidZrpcOutput
	}
	if !filepath.IsAbs(zrpcOut) {
		zrpcOut = filepath.Join(pwd, zrpcOut)
	}

	goOut = removePluginFlag(goOut)
	goOut, err = parseOutOut(src, goOut, goOpt, goPackage)
	if err != nil {
		return err
	}

	var isGoolePlugin = len(grpcOut) > 0
	// If grpcOut is not empty means that user generates grpc code by
	// https://google.golang.org/protobuf/cmd/protoc-gen-go and
	// https://google.golang.org/grpc/cmd/protoc-gen-go-grpc,
	// for details please see https://grpc.io/docs/languages/go/quickstart/
	if isGoolePlugin {
		grpcOut, err = parseOutOut(src, grpcOut, grpcOpt, goPackage)
		if err != nil {
			return err
		}
	} else {
		// Else it means that user generates grpc code by
		// https://github.com/golang/protobuf/tree/master/protoc-gen-go
		grpcOut = goOut
	}

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

	if isGoolePlugin && grpcOut != goOut {
		return fmt.Errorf("the --go_out and --go-grpc_out must be the same")
	}

	if goOut == zrpcOut || grpcOut == zrpcOut {
		recommendName := goPackage
		if len(recommendName) == 0 {
			recommendName = protoPkg
		}
		return fmt.Errorf("the zrpc and grpc output can not be the same, it is recommended to output grpc to the %q",
			filepath.Join(goOut, recommendName))
	}

	var ctx generator.ZRpcContext
	ctx.Src = source
	ctx.ProtoGenGoDir = goOut
	ctx.ProtoGenGrpcDir = grpcOut
	ctx.Output = zrpcOut
	ctx.ProtocCmd = strings.Join(protocArgs, " ")
	g, err := generator.NewDefaultRPCGenerator(style, generator.WithZRpcContext(&ctx))
	if err != nil {
		return err
	}

	return g.Generate(source, zrpcOut, nil)
}

// parseOutOut calculates the output place to grpc code, about to calculate logic for details
// please see https://developers.google.com/protocol-buffers/docs/reference/go-generated#invocation.
func parseOutOut(sourceDir, grpcOut, grpcOpt, goPackage string) (string, error) {
	if !filepath.IsAbs(grpcOut) {
		grpcOut = filepath.Join(sourceDir, grpcOut)
	}
	switch grpcOpt {
	case "", optImport:
		grpcOut = filepath.Join(grpcOut, goPackage)
	case optSourceRelative:
		grpcOut = filepath.Join(grpcOut)
	default:
		return "", fmt.Errorf("parseAndSetGrpcOut:  unknown path type %q: want %q or %q",
			grpcOpt, optImport, optSourceRelative)
	}

	return grpcOut, nil
}

func getGoPackage(source string) (string, string, error) {
	r, err := os.Open(source)
	if err != nil {
		return "", "", err
	}
	defer func() {
		_ = r.Close()
	}()

	parser := proto.NewParser(r)
	set, err := parser.Parse()
	if err != nil {
		return "", "", err
	}

	var goPackage, protoPkg string
	proto.Walk(set, proto.WithOption(func(option *proto.Option) {
		if option.Name == "go_package" {
			goPackage = option.Constant.Source
		}
	}), proto.WithPackage(func(p *proto.Package) {
		protoPkg = p.Name
	}))

	return goPackage, protoPkg, nil
}

func removeGoctlFlag(args []string) []string {
	var ret []string
	var step int
	for step < len(args) {
		arg := args[step]
		switch {
		case arg == "--style", arg == "--home", arg == "--zrpc_out":
			step += 2
			continue
		case strings.HasPrefix(arg, "--style="),
			strings.HasPrefix(arg, "--home="),
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
