package gen

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/rpc/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

const (
	handlerTemplate = `{{.head}}

package handler

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
func (s *{{.server}}Server) {{.method}} (ctx context.Context, in *{{.package}}.{{.request}}) (*{{.package}}.{{.response}}, error) {
	l := logic.New{{.logicName}}(ctx,s.svcCtx)
	return l.{{.method}}(in)
}
`
	typeFmt = `%sServer struct {
		svcCtx *svc.ServiceContext
	}`
)

func (g *defaultRpcGenerator) genHandler() error {
	handlerPath := g.dirM[dirHandler]
	file := g.ast
	pkg := file.Package
	pbImport := fmt.Sprintf(`%v "%v"`, pkg, g.mustGetPackage(dirPb))
	logicImport := fmt.Sprintf(`"%v"`, g.mustGetPackage(dirLogic))
	svcImport := fmt.Sprintf(`"%v"`, g.mustGetPackage(dirSvc))
	imports := []string{
		pbImport,
		logicImport,
		svcImport,
	}
	head := util.GetHead(g.Ctx.ProtoSource)
	for _, service := range file.Service {
		filename := fmt.Sprintf("%vhandler.go", service.Name.Lower())
		handlerFile := filepath.Join(handlerPath, filename)
		funcList, err := g.genFunctions(service)
		if err != nil {
			return err
		}
		err = util.With("server").GoFmt(true).Parse(handlerTemplate).SaveTo(map[string]interface{}{
			"head":    head,
			"types":   fmt.Sprintf(typeFmt, service.Name.Title()),
			"server":  service.Name.Title(),
			"imports": strings.Join(imports, "\n\t"),
			"funcs":   strings.Join(funcList, "\n"),
		}, handlerFile, true)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *defaultRpcGenerator) genFunctions(service *parser.RpcService) ([]string, error) {
	file := g.ast
	pkg := file.Package
	var functionList []string
	for _, method := range service.Funcs {
		buffer, err := util.With("func").Parse(functionTemplate).Execute(map[string]interface{}{
			"server":     service.Name.Title(),
			"logicName":  fmt.Sprintf("%sLogic", method.Name.Title()),
			"method":     method.Name.Title(),
			"package":    pkg,
			"request":    method.InType,
			"response":   method.OutType,
			"hasComment": len(method.Document),
			"comment":    strings.Join(method.Document, "\n"),
		})
		if err != nil {
			return nil, err
		}
		functionList = append(functionList, buffer.String())
	}
	return functionList, nil
}
