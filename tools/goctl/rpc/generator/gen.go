package generator

import (
	"path/filepath"

	conf "github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"
	"github.com/zeromicro/go-zero/tools/goctl/util/ctx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

// RPCGenerator defines a generator and configure
type RPCGenerator struct {
	g   Generator
	cfg *conf.Config
	ctx *ZRpcContext
}

type RPCGeneratorOption func(g *RPCGenerator)

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

// NewDefaultRPCGenerator wraps Generator with configure
func NewDefaultRPCGenerator(style string, options ...RPCGeneratorOption) (*RPCGenerator, error) {
	cfg, err := conf.NewConfig(style)
	if err != nil {
		return nil, err
	}
	return NewRPCGenerator(NewDefaultGenerator(), cfg, options...), nil
}

// NewRPCGenerator creates an instance for RPCGenerator
func NewRPCGenerator(g Generator, cfg *conf.Config, options ...RPCGeneratorOption) *RPCGenerator {
	out := &RPCGenerator{
		g:   g,
		cfg: cfg,
	}
	for _, opt := range options {
		opt(out)
	}
	return out
}

func WithZRpcContext(c *ZRpcContext) RPCGeneratorOption {
	return func(g *RPCGenerator) {
		g.ctx = c
	}
}

// Generate generates an rpc service, through the proto file,
// code storage directory, and proto import parameters to control
// the source file and target location of the rpc service that needs to be generated
func (g *RPCGenerator) Generate(src, target string, protoImportPath []string, goOptions ...string) error {
	abs, err := filepath.Abs(target)
	if err != nil {
		return err
	}

	err = pathx.MkdirIfNotExist(abs)
	if err != nil {
		return err
	}

	err = g.g.Prepare()
	if err != nil {
		return err
	}

	projectCtx, err := ctx.Prepare(abs)
	if err != nil {
		return err
	}

	p := parser.NewDefaultProtoParser()
	proto, err := p.Parse(src)
	if err != nil {
		return err
	}

	dirCtx, err := mkdir(projectCtx, proto, g.cfg, g.ctx)
	if err != nil {
		return err
	}

	err = g.g.GenEtc(dirCtx, proto, g.cfg)
	if err != nil {
		return err
	}

	err = g.g.GenPb(dirCtx, protoImportPath, proto, g.cfg, g.ctx, goOptions...)
	if err != nil {
		return err
	}

	err = g.g.GenConfig(dirCtx, proto, g.cfg)
	if err != nil {
		return err
	}

	err = g.g.GenSvc(dirCtx, proto, g.cfg)
	if err != nil {
		return err
	}

	err = g.g.GenLogic(dirCtx, proto, g.cfg)
	if err != nil {
		return err
	}

	err = g.g.GenServer(dirCtx, proto, g.cfg)
	if err != nil {
		return err
	}

	err = g.g.GenMain(dirCtx, proto, g.cfg)
	if err != nil {
		return err
	}

	err = g.g.GenCall(dirCtx, proto, g.cfg)

	console.NewColorConsole().MarkDone()

	return err
}
