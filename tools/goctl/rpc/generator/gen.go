package generator

import (
	"path/filepath"

	"github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/ctx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

type ZRpcContext struct {
	Src             string
	ProtocCmd       string
	ProtoGenGrpcDir string
	ProtoGenGoDir   string
	IsGooglePlugin  bool
	GoOutput        string
	GrpcOutput      string
	Output          string
}

// Generate generates an rpc service, through the proto file,
// code storage directory, and proto import parameters to control
// the source file and target location of the rpc service that needs to be generated
func (g *Generator) Generate(zctx *ZRpcContext) error {
	abs, err := filepath.Abs(zctx.Output)
	if err != nil {
		return err
	}

	err = pathx.MkdirIfNotExist(abs)
	if err != nil {
		return err
	}

	err = g.Prepare()
	if err != nil {
		return err
	}

	projectCtx, err := ctx.Prepare(abs)
	if err != nil {
		return err
	}

	p := parser.NewDefaultProtoParser()
	proto, err := p.Parse(zctx.Src)
	if err != nil {
		return err
	}

	dirCtx, err := mkdir(projectCtx, proto, g.cfg, zctx)
	if err != nil {
		return err
	}

	err = g.GenEtc(dirCtx, proto, g.cfg)
	if err != nil {
		return err
	}

	err = g.GenPb(dirCtx, zctx)
	if err != nil {
		return err
	}

	err = g.GenConfig(dirCtx, proto, g.cfg)
	if err != nil {
		return err
	}

	err = g.GenSvc(dirCtx, proto, g.cfg)
	if err != nil {
		return err
	}

	err = g.GenLogic(dirCtx, proto, g.cfg)
	if err != nil {
		return err
	}

	err = g.GenServer(dirCtx, proto, g.cfg)
	if err != nil {
		return err
	}

	err = g.GenMain(dirCtx, proto, g.cfg)
	if err != nil {
		return err
	}

	err = g.GenCall(dirCtx, proto, g.cfg)

	console.NewColorConsole().MarkDone()

	return err
}
