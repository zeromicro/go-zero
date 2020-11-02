package generator

import (
	"path/filepath"

	"github.com/tal-tech/go-zero/tools/goctl/rpcv2/ctx"
	"github.com/tal-tech/go-zero/tools/goctl/rpcv2/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

type (
	RpcGenerator struct {
		g Generator
	}
)

func NewDefaultRpcGenerator() *RpcGenerator {
	return NewRpcGenerator(NewDefaultGenerator())
}

func NewRpcGenerator(g Generator) *RpcGenerator {
	return &RpcGenerator{
		g: g,
	}
}

func (g *RpcGenerator) Generate(src, target string, IPATH []string) error {
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

	projectCtx, err := ctx.Background(abs)
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

	err = g.g.GenEtc(dirCtx, dirCtx.GetEtc(), proto)
	if err != nil {
		return err
	}

	err = g.g.GenPb(dirCtx, IPATH, dirCtx.GetPb(), proto)
	if err != nil {
		return err
	}

	err = g.g.GenConfig(dirCtx, dirCtx.GetConfig(), proto)
	if err != nil {
		return err
	}

	err = g.g.GenSvc(dirCtx, dirCtx.GetSvc(), proto)
	if err != nil {
		return err
	}

	err = g.g.GenLogic(dirCtx, dirCtx.GetLogic(), proto)
	if err != nil {
		return err
	}

	err = g.g.GenServer(dirCtx, dirCtx.GetServer(), proto)
	if err != nil {
		return err
	}

	err = g.g.GenMain(dirCtx, dirCtx.GetWorkDir(), proto)
	if err != nil {
		return err
	}

	err = g.g.GenCall(dirCtx, dirCtx.GetCall(), proto)
	return err
}
