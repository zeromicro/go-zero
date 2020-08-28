package gogen

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/util"
)

const (
	remoteTemplate = `{{.head}}

package handler

import {{.imports}}

type {{.types}}

{{.newFuncs}}
`
	functionTemplate = `{{.head}}

package handler

import (
	"context"

	{{.imports}}
)

type {{.server}}Server struct{}

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

func (g *defaultRpcGenerator) genRemoteHandler() error {
	handlerPath := g.dirM[dirHandler]
	serverGo := fmt.Sprintf("%vhandler.go", g.Ctx.ServiceName.Lower())
	fileName := filepath.Join(handlerPath, serverGo)
	file := g.ast
	svcImport := fmt.Sprintf(`"%v"`, g.mustGetPackage(dirSvc))
	types := make([]string, 0)
	newFuncs := make([]string, 0)
	head := util.GetHead(g.Ctx.ProtoSource)
	for _, service := range file.Service {
		types = append(types, fmt.Sprintf(typeFmt, service.Name.Title()))
		newFuncs = append(newFuncs, fmt.Sprintf(newFuncFmt, service.Name.Title(), service.Name.Title(), service.Name.Title()))
	}
	err := util.With("server").GoFmt(true).Parse(remoteTemplate).SaveTo(map[string]interface{}{
		"head":     head,
		"types":    strings.Join(types, "\n"),
		"newFuncs": strings.Join(newFuncs, "\n"),
		"imports":  svcImport,
	}, fileName, true)
	if err != nil {
		return err
	}
	return g.genFunctions()
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
				"imports":    strings.Join(handlerImports, "\r\n"),
				"logicName":  fmt.Sprintf("%sLogic", method.Name.Title()),
				"method":     method.Name.Title(),
				"package":    pkg,
				"request":    method.InType,
				"response":   method.OutType,
				"hasComment": len(method.Document),
				"comment":    strings.Join(method.Document, "\r\n"),
			}, filename, true)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
