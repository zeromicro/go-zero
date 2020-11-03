package generator

import (
	"path/filepath"

	"github.com/tal-tech/go-zero/tools/goctl/rpcv2/ctx"
	"github.com/tal-tech/go-zero/tools/goctl/rpcv2/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/console"
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

	err = g.g.GenEtc(dirCtx, proto)
	if err != nil {
		return err
	}

	err = g.g.GenPb(dirCtx, IPATH, proto)
	if err != nil {
		return err
	}

	err = g.g.GenConfig(dirCtx, proto)
	if err != nil {
		return err
	}

	err = g.g.GenSvc(dirCtx, proto)
	if err != nil {
		return err
	}

	err = g.g.GenLogic(dirCtx, proto)
	if err != nil {
		return err
	}

	err = g.g.GenServer(dirCtx, proto)
	if err != nil {
		return err
	}

	err = g.g.GenMain(dirCtx, proto)
	if err != nil {
		return err
	}

	err = g.g.GenCall(dirCtx, proto)

	console.NewColorConsole().MarkDone()

	return err
}
