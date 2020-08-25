package gogen

import (
	"github.com/tal-tech/go-zero/tools/goctl/rpc/ctx"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/parser"
)

const (
	dirTarget          = "dirTarget"
	dirConfig          = "config"
	dirEtc             = "etc"
	dirSvc             = "svc"
	dirHandler         = "handler"
	dirLogic           = "logic"
	dirPb              = "pb"
	dirInternal        = "internal"
	fileConfig         = "config.go"
	fileServiceContext = "servicecontext.go"
)

type (
	defaultRpcGenerator struct {
		dirM map[string]string
		Ctx  *ctx.RpcContext
		ast  *parser.PbAst
	}
)

func NewDefaultRpcGenerator(ctx *ctx.RpcContext) *defaultRpcGenerator {
	return &defaultRpcGenerator{
		Ctx: ctx,
	}
}

func (g *defaultRpcGenerator) Generate() error {
	ctx := g.Ctx
	ctx.Must(g.createDir())
	ctx.Must(g.genEtc())
	ctx.Must(g.genPb())
	ctx.Must(g.genConfig())
	ctx.Must(g.genSvc())
	ctx.Must(g.genLogic())
	ctx.Must(g.genRemoteHandler())
	ctx.Must(g.genMain())
	ctx.Must(g.genShared())
	return nil
}
