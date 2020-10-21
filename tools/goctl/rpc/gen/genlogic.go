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
	logicTemplate = `package logic

import (
	"context"

	{{.imports}}

	"github.com/tal-tech/go-zero/core/logx"
)

type {{.logicName}} struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func New{{.logicName}}(ctx context.Context,svcCtx *svc.ServiceContext) *{{.logicName}} {
	return &{{.logicName}}{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}
{{.functions}}
`
	logicFunctionTemplate = `{{if .hasComment}}{{.comment}}{{end}}
func (l *{{.logicName}}) {{.method}} (in {{.request}}) ({{.response}}, error) {
	// todo: add your logic here and delete this line
	
	return &{{.responseType}}{}, nil
}
`
)

func (g *defaultRpcGenerator) genLogic() error {
	logicPath := g.dirM[dirLogic]
	protoPkg := g.ast.Package
	service := g.ast.Service
	for _, item := range service {
		for _, method := range item.Funcs {
			logicName := fmt.Sprintf("%slogic.go", method.Name.Lower())
			filename := filepath.Join(logicPath, logicName)
			functions, importList, err := g.genLogicFunction(protoPkg, method)
			if err != nil {
				return err
			}
			imports := collection.NewSet()
			svcImport := fmt.Sprintf(`"%v"`, g.mustGetPackage(dirSvc))
			imports.AddStr(svcImport)
			imports.AddStr(importList...)
			text, err := util.LoadTemplate(category, logicTemplateFileFile, logicTemplate)
			if err != nil {
				return err
			}
			err = util.With("logic").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
				"logicName": fmt.Sprintf("%sLogic", method.Name.Title()),
				"functions": functions,
				"imports":   strings.Join(imports.KeysStr(), util.NL),
			}, filename, false)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *defaultRpcGenerator) genLogicFunction(packageName string, method *parser.Func) (string, []string, error) {
	var functions = make([]string, 0)
	var imports = collection.NewSet()
	if method.ParameterIn.Package == packageName || method.ParameterOut.Package == packageName {
		imports.AddStr(fmt.Sprintf(`%v "%v"`, packageName, g.mustGetPackage(dirPb)))
	}
	imports.AddStr(g.ast.Imports[method.ParameterIn.Package])
	imports.AddStr(g.ast.Imports[method.ParameterOut.Package])
	text, err := util.LoadTemplate(category, logicFuncTemplateFileFile, logicFunctionTemplate)
	if err != nil {
		return "", nil, err
	}

	buffer, err := util.With("fun").Parse(text).Execute(map[string]interface{}{
		"logicName":    fmt.Sprintf("%sLogic", method.Name.Title()),
		"method":       method.Name.Title(),
		"request":      method.ParameterIn.StarExpression,
		"response":     method.ParameterOut.StarExpression,
		"responseType": method.ParameterOut.Expression,
		"hasComment":   method.HaveDoc(),
		"comment":      method.GetDoc(),
	})
	if err != nil {
		return "", nil, err
	}

	functions = append(functions, buffer.String())
	return strings.Join(functions, util.NL), imports.KeysStr(), nil
}
