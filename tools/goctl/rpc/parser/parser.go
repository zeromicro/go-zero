package parser

import (
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/core/lang"
	"github.com/tal-tech/go-zero/tools/goctl/util/console"
)

func Transfer(proto, target string, externalImport []*Import, console console.Console) (*PbAst, error) {
	messageM := make(map[string]lang.PlaceholderType)
	enumM := make(map[string]*Enum)
	protoAst, err := parseProto(proto, messageM, enumM)
	if err != nil {
		return nil, err
	}
	for _, item := range externalImport {
		err = checkImport(item.OriginalProtoPath)
		if err != nil {
			return nil, err
		}
		innerAst, err := parseProto(item.OriginalProtoPath, protoAst.Message, protoAst.Enum)
		if err != nil {
			return nil, err
		}
		for k, v := range innerAst.Message {
			protoAst.Message[k] = v
		}
		for k, v := range innerAst.Enum {
			protoAst.Enum[k] = v
		}
	}
	protoAst.Import = externalImport
	protoAst.PbSrc = filepath.Join(target, strings.TrimSuffix(filepath.Base(proto), ".proto")+".pb.go")
	return transfer(protoAst, console)
}

func transfer(proto *Proto, console console.Console) (*PbAst, error) {
	parser := MustNewAstParser(proto, console)
	parse, err := parser.Parse()
	if err != nil {
		return nil, err
	}
	return parse, nil
}
