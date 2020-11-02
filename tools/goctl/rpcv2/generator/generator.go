package generator

import (
	"github.com/tal-tech/go-zero/tools/goctl/rpcv2/parser"
)

type (
	Generator interface {
		Prepare() error
		GenMain(ctx DirContext, dir Dir, proto parser.Proto) error
		GenCall(ctx DirContext, dir Dir, proto parser.Proto) error
		GenEtc(ctx DirContext, dir Dir, proto parser.Proto) error
		GenConfig(ctx DirContext, dir Dir, proto parser.Proto) error
		GenLogic(ctx DirContext, dir Dir, proto parser.Proto) error
		GenServer(ctx DirContext, dir Dir, proto parser.Proto) error
		GenSvc(ctx DirContext, dir Dir, proto parser.Proto) error
		// IPATH is the native command of protoc, see $protoc -h
		GenPb(ctx DirContext, IPATH string, dir Dir, proto parser.Proto) error
	}
)
