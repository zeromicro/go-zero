package gogen

import (
	"fmt"
	"path/filepath"
	"strings"

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

{{if .hasComment}}{{.comment}}{{end}}
func (s *{{.server}}Server) {{.method}} (ctx context.Context, in *{{.package}}.{{.request}}) (*{{.package}}.{{.response}}, error) {
	l := logic.New{{.logicName}}(ctx,s.svcCtx)
	return l.{{.method}}(in)
}
`
	functionTemplate = `{{.head}}

package handler

import (
	"context"

	{{.imports}}
)

{{if .hasComment}}{{.comment}}{{end}}
func (s *{{.server}}Server) {{.method}} (ctx context.Context, in *{{.package}}.{{.request}}) (*{{.package}}.{{.response}}, error) {
	l := logic.New{{.logicName}}(ctx,s.svcCtx)
	return l.{{.method}}(in)
}
`
	typeFmt = `%sServer struct {
		svcCtx *svc.ServiceContext
	}`
	newFuncFmt = `func New%sServer(svcCtx *svc.ServiceContext) *%sServer {
	return &%sServer{
		svcCtx: svcCtx,
	}
}`
)

func (g *defaultRpcGenerator) genHandler() error {
	handlerPath := g.dirM[dirHandler]
	filename := fmt.Sprintf("%vhandler.go", g.Ctx.ServiceName.Lower())
	handlerFile := filepath.Join(handlerPath, filename)
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
	types := make([]string, 0)
	newFuncs := make([]string, 0)
	head := util.GetHead(g.Ctx.ProtoSource)
	for _, service := range file.Service {
		types = append(types, fmt.Sprintf(typeFmt, service.Name.Title()))
		newFuncs = append(newFuncs, fmt.Sprintf(newFuncFmt, service.Name.Title(),
			service.Name.Title(), service.Name.Title()))
	}

	return util.With("server").GoFmt(true).Parse(handlerTemplate).SaveTo(map[string]interface{}{
		"head":     head,
		"types":    strings.Join(types, "\n"),
		"newFuncs": strings.Join(newFuncs, "\n"),
		"imports":  strings.Join(imports, "\n\t"),
	}, handlerFile, true)
}

func (g *defaultRpcGenerator) genFunctions() error {
	handlerPath := g.dirM[dirHandler]
	file := g.ast
	pkg := file.Package

	head := util.GetHead(g.Ctx.ProtoSource)
	handlerImports := make([]string, 0)
	pbImport := fmt.Sprintf(`%v "%v"`, pkg, g.mustGetPackage(dirPb))
	handlerImports = append(handlerImports, pbImport, fmt.Sprintf(`"%v"`, g.mustGetPackage(dirLogic)))
	for _, service := range file.Service {
		for _, method := range service.Funcs {
			handlerName := fmt.Sprintf("%shandler.go", method.Name.Lower())
			filename := filepath.Join(handlerPath, handlerName)
			// override
			err := util.With("func").GoFmt(true).Parse(functionTemplate).SaveTo(map[string]interface{}{
				"head":       head,
				"server":     service.Name.Title(),
				"imports":    strings.Join(handlerImports, "\n"),
				"logicName":  fmt.Sprintf("%sLogic", method.Name.Title()),
				"method":     method.Name.Title(),
				"package":    pkg,
				"request":    method.InType,
				"response":   method.OutType,
				"hasComment": len(method.Document),
				"comment":    strings.Join(method.Document, "\n"),
			}, filename, true)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
