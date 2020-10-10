package gen

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/core/collection"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/parser"
)

const (
	protocCmd     = "protoc"
	grpcPluginCmd = "--go_out=plugins=grpc"
)

func (g *defaultRpcGenerator) genPb() error {
	pbPath := g.dirM[dirPb]
	imports, containsAny, err := parser.ParseImport(g.Ctx.ProtoFileSrc)
	if err != nil {
		return err
	}

	err = g.protocGenGo(pbPath, imports)
	if err != nil {
		return err
	}
	ast, err := parser.Transfer(g.Ctx.ProtoFileSrc, pbPath, imports, g.Ctx.Console)
	if err != nil {
		return err
	}
	ast.ContainsAny = containsAny

	if len(ast.Service) == 0 {
		return fmt.Errorf("service not found")
	}
	g.ast = ast
	return nil
}

func (g *defaultRpcGenerator) protocGenGo(target string, imports []*parser.Import) error {
	dir := filepath.Dir(g.Ctx.ProtoFileSrc)
	// cmd join,see the document of proto generating class @https://developers.google.com/protocol-buffers/docs/proto3#generating
	// template: protoc -I=${import_path} -I=${other_import_path} -I=${...} --go_out=plugins=grpc,M${pb_package_kv}, M${...} :${target_dir}
	// eg: protoc -I=${GOPATH}/src -I=. example.proto --go_out=plugins=grpc,Mbase/base.proto=github.com/go-zero/base.proto:.
	// note: the external import out of the project which are found in ${GOPATH}/src so far.

	buffer := new(bytes.Buffer)
	buffer.WriteString(protocCmd + " ")
	targetImportFiltered := collection.NewSet()

	for _, item := range imports {
		buffer.WriteString(fmt.Sprintf("-I=%s ", item.OriginalDir))
		if len(item.BridgeImport) == 0 {
			continue
		}
		targetImportFiltered.AddStr(item.BridgeImport)

	}
	buffer.WriteString("-I=${GOPATH}/src ")
	buffer.WriteString(fmt.Sprintf("-I=%s %s ", dir, g.Ctx.ProtoFileSrc))

	buffer.WriteString(grpcPluginCmd)
	if targetImportFiltered.Count() > 0 {
		buffer.WriteString(fmt.Sprintf(",%v", strings.Join(targetImportFiltered.KeysStr(), ",")))
	}
	buffer.WriteString(":" + target)
	g.Ctx.Debug("-> " + buffer.String())
	stdout, err := execx.Run(buffer.String(), "")
	if err != nil {
		return err
	}

	if len(stdout) > 0 {
		g.Ctx.Info(stdout)
	}

	return nil
}
