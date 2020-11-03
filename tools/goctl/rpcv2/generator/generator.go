package generator

import (
	"github.com/tal-tech/go-zero/tools/goctl/rpcv2/parser"
)

type (
	Generator interface {
		Prepare() error
		GenMain(ctx DirContext, proto parser.Proto) error
		GenCall(ctx DirContext, proto parser.Proto) error
		GenEtc(ctx DirContext, proto parser.Proto) error
		GenConfig(ctx DirContext, proto parser.Proto) error
		GenLogic(ctx DirContext, proto parser.Proto) error
		GenServer(ctx DirContext, proto parser.Proto) error
		GenSvc(ctx DirContext, proto parser.Proto) error
		// IPATH is the native command of protoc, see $protoc -h
		GenPb(ctx DirContext, IPATH []string, proto parser.Proto) error
	}
)
