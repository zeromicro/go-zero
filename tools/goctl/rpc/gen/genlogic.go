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
func (l *{{.logicName}}) {{.method}} (in *{{.package}}.{{.request}}) (*{{.package}}.{{.response}}, error) {
	// todo: add your logic here and delete this line
	
	return &{{.package}}.{{.response}}{}, nil
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
			functions, err := genLogicFunction(protoPkg, method)
			if err != nil {
				return err
			}
			imports := collection.NewSet()
			pbImport := fmt.Sprintf(`%v "%v"`, protoPkg, g.mustGetPackage(dirPb))
			svcImport := fmt.Sprintf(`"%v"`, g.mustGetPackage(dirSvc))
			imports.AddStr(pbImport, svcImport)
			err = util.With("logic").GoFmt(true).Parse(logicTemplate).SaveTo(map[string]interface{}{
				"logicName": fmt.Sprintf("%sLogic", method.Name.Title()),
				"functions": functions,
				"imports":   strings.Join(imports.KeysStr(), "\n"),
			}, filename, false)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func genLogicFunction(packageName string, method *parser.Func) (string, error) {
	var functions = make([]string, 0)
	buffer, err := util.With("fun").Parse(logicFunctionTemplate).Execute(map[string]interface{}{
		"logicName":  fmt.Sprintf("%sLogic", method.Name.Title()),
		"method":     method.Name.Title(),
		"package":    packageName,
		"request":    method.InType,
		"response":   method.OutType,
		"hasComment": len(method.Document) > 0,
		"comment":    strings.Join(method.Document, "\n"),
	})
	if err != nil {
		return "", err
	}
	functions = append(functions, buffer.String())
	return strings.Join(functions, "\n"), nil
}
