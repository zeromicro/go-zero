package generator

import "github.com/tal-tech/go-zero/tools/goctl/rpc/parser"

type Generator interface {
	Prepare() error
	GenMain(ctx DirContext, proto parser.Proto, namingStyle NamingStyle) error
	GenCall(ctx DirContext, proto parser.Proto, namingStyle NamingStyle) error
	GenEtc(ctx DirContext, proto parser.Proto, namingStyle NamingStyle) error
	GenConfig(ctx DirContext, proto parser.Proto, namingStyle NamingStyle) error
	GenLogic(ctx DirContext, proto parser.Proto, namingStyle NamingStyle) error
	GenServer(ctx DirContext, proto parser.Proto, namingStyle NamingStyle) error
	GenSvc(ctx DirContext, proto parser.Proto, namingStyle NamingStyle) error
	GenPb(ctx DirContext, protoImportPath []string, proto parser.Proto, namingStyle NamingStyle) error
}
