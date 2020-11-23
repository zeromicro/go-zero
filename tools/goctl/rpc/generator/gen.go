package generator

import (
	"path/filepath"

	conf "github.com/tal-tech/go-zero/tools/goctl/config"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/console"
	"github.com/tal-tech/go-zero/tools/goctl/util/ctx"
)

type RpcGenerator struct {
	g   Generator
	cfg *conf.Config
}

func NewDefaultRpcGenerator(style string) (*RpcGenerator, error) {
	cfg, err := conf.NewConfig(style)
	if err != nil {
		return nil, err
	}
	return NewRpcGenerator(NewDefaultGenerator(), cfg), nil
}

func NewRpcGenerator(g Generator, cfg *conf.Config) *RpcGenerator {
	return &RpcGenerator{
		g:   g,
		cfg: cfg,
	}
}

func (g *RpcGenerator) Generate(src, target string, protoImportPath []string) error {
	abs, err := filepath.Abs(target)
	if err != nil {
		return err
	}

	err = util.MkdirIfNotExist(abs)
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

	dirCtx, err := mkdir(projectCtx, proto)
	if err != nil {
		return err
	}

	err = g.g.GenEtc(dirCtx, proto, g.cfg)
	if err != nil {
		return err
	}

	err = g.g.GenPb(dirCtx, protoImportPath, proto, g.cfg)
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
