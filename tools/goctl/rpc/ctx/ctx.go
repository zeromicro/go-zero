package ctx

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/urfave/cli"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/tools/goctl/util"
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
		console.Console
	}
)

func MustCreateRpcContext(protoSrc, targetDir, sharedDir, serviceName string, idea bool) *RpcContext {
	log := console.NewConsole(idea)
	info, err := prepare(log)
	log.Must(err)

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
	if stringx.From(sharedDir).IsEmptyOrSpace() {
		sharedDir = filepath.Join(current, "shared")
	}
	targetDirFp, err := filepath.Abs(targetDir)
	log.Must(err)

	sharedFp, err := filepath.Abs(sharedDir)
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
		Module:       info.GoMod.Module,
		ProtoFileSrc: srcFp,
		ProtoSource:  filepath.Base(srcFp),
		TargetDir:    targetDirFp,
		SharedDir:    sharedFp,
		Console:      log,
	}
}
func MustCreateRpcContextFromCli(ctx *cli.Context) *RpcContext {
	os := runtime.GOOS
	switch os {
	case "darwin":
	case "windows":
		logx.Must(fmt.Errorf("windows will support soon"))
	default:
		logx.Must(fmt.Errorf("unexpected os: %s", os))
	}
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
