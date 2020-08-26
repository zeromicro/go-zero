package ctx

import (
	"path/filepath"
	"strings"

	"github.com/urfave/cli"

	"github.com/tal-tech/go-zero/tools/goctl/rpc/project"
	"github.com/tal-tech/go-zero/tools/goctl/util/console"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

const (
	flagSrc     = "src"
	flagDir     = "dir"
	flagShared  = "shared"
	flagService = "service"
	flagIdea    = "idea"
)

var (
	tips = `warning: this operation will overwrite the protoc-gen-plugin in gopath
protoc-gen-go: switch to %s`
)

type (
	RpcContext struct {
		ProjectPath  string
		ProjectName  stringx.String
		ServiceName  stringx.String
		CurrentPath  string
		Module       string
		ProtoFileSrc string
		ProtoSource  string
		TargetDir    string
		SharedDir    string
		GoPath       string
		console.Console
	}
)

func MustCreateRpcContext(protoSrc, targetDir, sharedDir, serviceName string, idea bool) *RpcContext {
	log := console.NewConsole(idea)
	goMod, err := prepare()
	log.Must(err)
	log.Info(tips, goMod.Protobuf())

	if stringx.From(protoSrc).IsEmptyOrSpace() {
		log.Fatalln("expected proto source, but nothing found")
	}
	srcFp, err := filepath.Abs(protoSrc)
	log.Must(err)
	current := filepath.Dir(srcFp)
	if stringx.From(targetDir).IsEmptyOrSpace() {
		targetDir = current
	}
	if stringx.From(sharedDir).IsEmptyOrSpace() {
		sharedDir = filepath.Join(current, "shared")
	}
	targetDirFp, err := filepath.Abs(targetDir)
	log.Must(err)
	sharedFp, err := filepath.Abs(sharedDir)
	log.Must(err)
	info, err := project.Info()
	log.Must(err)
	if stringx.From(serviceName).IsEmptyOrSpace() {
		serviceName = getServiceFromRpcStructure(targetDirFp)
	}
	serviceNameString := stringx.From(serviceName)
	if serviceNameString.IsEmptyOrSpace() {
		log.Fatalln("service name is not found")
	}
	return &RpcContext{
		ProjectPath:  info.Path,
		ProjectName:  stringx.From(info.Name),
		ServiceName:  serviceNameString,
		CurrentPath:  current,
		Module:       goMod.module,
		ProtoFileSrc: srcFp,
		ProtoSource:  filepath.Base(srcFp),
		TargetDir:    targetDirFp,
		SharedDir:    sharedFp,
		GoPath:       info.GoPath,
		Console:      log,
	}
}
func MustCreateRpcContextFromCli(ctx *cli.Context) *RpcContext {
	protoSrc := ctx.String(flagSrc)
	targetDir := ctx.String(flagDir)
	sharedDir := ctx.String(flagShared)
	serviceName := ctx.String(flagService)
	idea := ctx.Bool(flagIdea)
	return MustCreateRpcContext(protoSrc, targetDir, sharedDir, serviceName, idea)
}

func getServiceFromRpcStructure(targetDir string) string {
	targetDir = filepath.Clean(targetDir)
	suffix := filepath.Join("cmd", "rpc")
	return filepath.Base(strings.TrimSuffix(targetDir, suffix))
}
