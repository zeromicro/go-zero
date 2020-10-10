package ctx

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/console"
	"github.com/tal-tech/go-zero/tools/goctl/util/project"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/vars"
	"github.com/urfave/cli"
)

const (
	flagSrc     = "src"
	flagDir     = "dir"
	flagService = "service"
	flagIdea    = "idea"
)

type RpcContext struct {
	ProjectPath  string
	ProjectName  stringx.String
	ServiceName  stringx.String
	CurrentPath  string
	Module       string
	ProtoFileSrc string
	ProtoSource  string
	TargetDir    string
	IsInGoEnv    bool
	console.Console
}

func MustCreateRpcContext(protoSrc, targetDir, serviceName string, idea bool) *RpcContext {
	log := console.NewConsole(idea)

	if stringx.From(protoSrc).IsEmptyOrSpace() {
		log.Fatalln("expected proto source, but nothing found")
	}
	srcFp, err := filepath.Abs(protoSrc)
	log.Must(err)

	if !util.FileExists(srcFp) {
		log.Fatalln("%s is not exists", srcFp)
	}
	current := filepath.Dir(srcFp)
	if stringx.From(targetDir).IsEmptyOrSpace() {
		targetDir = current
	}
	targetDirFp, err := filepath.Abs(targetDir)
	log.Must(err)

	if stringx.From(serviceName).IsEmptyOrSpace() {
		serviceName = getServiceFromRpcStructure(targetDirFp)
	}
	serviceNameString := stringx.From(serviceName)
	if serviceNameString.IsEmptyOrSpace() {
		log.Fatalln("service name not found")
	}

	info, err := project.Prepare(targetDir, true)
	log.Must(err)

	return &RpcContext{
		ProjectPath:  info.Path,
		ProjectName:  stringx.From(info.Name),
		ServiceName:  serviceNameString,
		CurrentPath:  current,
		Module:       info.GoMod.Module,
		ProtoFileSrc: srcFp,
		ProtoSource:  filepath.Base(srcFp),
		TargetDir:    targetDirFp,
		IsInGoEnv:    info.IsInGoEnv,
		Console:      log,
	}
}
func MustCreateRpcContextFromCli(ctx *cli.Context) *RpcContext {
	os := runtime.GOOS
	switch os {
	case vars.OsMac, vars.OsLinux, vars.OsWindows:
	default:
		logx.Must(fmt.Errorf("unexpected os: %s", os))
	}
	protoSrc := ctx.String(flagSrc)
	targetDir := ctx.String(flagDir)
	serviceName := ctx.String(flagService)
	idea := ctx.Bool(flagIdea)
	return MustCreateRpcContext(protoSrc, targetDir, serviceName, idea)
}

func getServiceFromRpcStructure(targetDir string) string {
	targetDir = filepath.Clean(targetDir)
	suffix := filepath.Join("cmd", "rpc")
	return filepath.Base(strings.TrimSuffix(targetDir, suffix))
}
