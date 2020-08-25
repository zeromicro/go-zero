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

type (
	RpcContext struct {
		ProjectName  stringx.String
		ServiceName  stringx.String
		CurrentPath  string
		ProtoFileSrc string
		ProtoSource  string
		TargetDir    string
		SharedDir    string
		console.Console
	}
)

func MustCreateRpcContext(ctx *cli.Context) *RpcContext {
	protoSrc := ctx.String(flagSrc)
	targetDir := ctx.String(flagDir)
	sharedDir := ctx.String(flagShared)
	serviceName := ctx.String(flagService)
	idea := ctx.Bool(flagIdea)
	log := console.NewConsole(idea)
	if stringx.From(protoSrc).IsEmptyOrSpace() {
		log.Fatalln("expected proto source, but nothing found")
	}
	if stringx.From(targetDir).IsEmptyOrSpace() {
		targetDir = "."
	}
	if stringx.From(sharedDir).IsEmptyOrSpace() {
		targetDir = filepath.Join(".", "shared")
	}
	current, err := filepath.Abs(".")
	log.Must(err)
	srcFp, err := filepath.Abs(protoSrc)
	log.Must(err)
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
		ProjectName:  stringx.From(info.Name),
		ServiceName:  serviceNameString,
		CurrentPath:  current,
		ProtoFileSrc: srcFp,
		ProtoSource:  filepath.Base(srcFp),
		TargetDir:    targetDirFp,
		SharedDir:    sharedFp,
	}
}

func getServiceFromRpcStructure(targetDir string) string {
	targetDir = filepath.Clean(targetDir)
	suffix := filepath.Join("cmd", "rpc")
	return strings.TrimSuffix(targetDir, suffix)
}
