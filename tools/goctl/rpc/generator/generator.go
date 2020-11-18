package generator

import (
	"github.com/tal-tech/go-zero/tools/goctl/rpc/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util/name"
)

type Generator interface {
	Prepare() error
	GenMain(ctx DirContext, proto parser.Proto, namingStyle name.NamingStyle) error
	GenCall(ctx DirContext, proto parser.Proto, namingStyle name.NamingStyle) error
	GenEtc(ctx DirContext, proto parser.Proto, namingStyle name.NamingStyle) error
	GenConfig(ctx DirContext, proto parser.Proto, namingStyle name.NamingStyle) error
	GenLogic(ctx DirContext, proto parser.Proto, namingStyle name.NamingStyle) error
	GenServer(ctx DirContext, proto parser.Proto, namingStyle name.NamingStyle) error
	GenSvc(ctx DirContext, proto parser.Proto, namingStyle name.NamingStyle) error
	GenPb(ctx DirContext, protoImportPath []string, proto parser.Proto, namingStyle name.NamingStyle) error
}
