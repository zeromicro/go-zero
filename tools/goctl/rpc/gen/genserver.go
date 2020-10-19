package gen

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/core/collection"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

const (
	serverTemplate = `{{.head}}

package server

import (
	"context"

	{{.imports}}
)

type {{.types}}

func New{{.server}}Server(svcCtx *svc.ServiceContext) *{{.server}}Server {
	return &{{.server}}Server{
		svcCtx: svcCtx,
	}
}

{{.funcs}}
`
	functionTemplate = `
{{if .hasComment}}{{.comment}}{{end}}
func (s *{{.server}}Server) {{.method}} (ctx context.Context, in {{.request}}) ({{.response}}, error) {
	l := logic.New{{.logicName}}(ctx,s.svcCtx)
	return l.{{.method}}(in)
}
`
	typeFmt = `%sServer struct {
		svcCtx *svc.ServiceContext
	}`
)

func (g *defaultRpcGenerator) genHandler() error {
	serverPath := g.dirM[dirServer]
	file := g.ast
	logicImport := fmt.Sprintf(`"%v"`, g.mustGetPackage(dirLogic))
	svcImport := fmt.Sprintf(`"%v"`, g.mustGetPackage(dirSvc))
	imports := collection.NewSet()
	imports.AddStr(logicImport, svcImport)

	head := util.GetHead(g.Ctx.ProtoSource)
	for _, service := range file.Service {
		filename := fmt.Sprintf("%vserver.go", service.Name.Lower())
		serverFile := filepath.Join(serverPath, filename)
		funcList, importList, err := g.genFunctions(service)
		if err != nil {
			return err
		}
		imports.AddStr(importList...)
		err = util.With("server").GoFmt(true).Parse(serverTemplate).SaveTo(map[string]interface{}{
			"head":    head,
			"types":   fmt.Sprintf(typeFmt, service.Name.Title()),
			"server":  service.Name.Title(),
			"imports": strings.Join(imports.KeysStr(), util.NL),
			"funcs":   strings.Join(funcList, util.NL),
		}, serverFile, true)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *defaultRpcGenerator) genFunctions(service *parser.RpcService) ([]string, []string, error) {
	file := g.ast
	pkg := file.Package
	var functionList []string
	imports := collection.NewSet()
	for _, method := range service.Funcs {
		if method.ParameterIn.Package == pkg || method.ParameterOut.Package == pkg {
			imports.AddStr(fmt.Sprintf(`%v "%v"`, pkg, g.mustGetPackage(dirPb)))
		}
		imports.AddStr(g.ast.Imports[method.ParameterIn.Package])
		imports.AddStr(g.ast.Imports[method.ParameterOut.Package])
		buffer, err := util.With("func").Parse(functionTemplate).Execute(map[string]interface{}{
			"server":     service.Name.Title(),
			"logicName":  fmt.Sprintf("%sLogic", method.Name.Title()),
			"method":     method.Name.Title(),
			"package":    pkg,
			"request":    method.ParameterIn.StarExpression,
			"response":   method.ParameterOut.StarExpression,
			"hasComment": method.HaveDoc(),
			"comment":    method.GetDoc(),
		})
		if err != nil {
			return nil, nil, err
		}
		functionList = append(functionList, buffer.String())
	}
	return functionList, imports.KeysStr(), nil
}
