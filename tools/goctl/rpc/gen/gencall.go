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
	typesFilename    = "types.go"
	callTemplateText = `{{.head}}

//go:generate mockgen -destination ./{{.name}}_mock.go -package {{.filePackage}} -source $GOFILE

package {{.filePackage}}

import (
	"context"

	{{.package}}

	"github.com/tal-tech/go-zero/core/jsonx"
	"github.com/tal-tech/go-zero/zrpc"
)

type (
	{{.serviceName}} interface {
		{{.interface}}
	}

	default{{.serviceName}} struct {
		cli zrpc.Client
	}
)

func New{{.serviceName}}(cli zrpc.Client) {{.serviceName}} {
	return &default{{.serviceName}}{
		cli: cli,
	}
}

{{.functions}}
`
	callTemplateTypes = `{{.head}}

package {{.filePackage}}

import "errors"

var errJsonConvert = errors.New("json convert error")

{{.const}}

{{.types}}
`
	callInterfaceFunctionTemplate = `{{if .hasComment}}{{.comment}}
{{end}}{{.method}}(ctx context.Context,in *{{.pbRequest}}) (*{{.pbResponse}},error)`

	callFunctionTemplate = `
{{if .hasComment}}{{.comment}}{{end}}
func (m *default{{.rpcServiceName}}) {{.method}}(ctx context.Context,in *{{.pbRequestName}}) (*{{.pbResponse}}, error) {
	var request {{.pbRequest}}
	bts, err := jsonx.Marshal(in)
	if err != nil {
		return nil, errJsonConvert
	}

	err = jsonx.Unmarshal(bts, &request)
	if err != nil {
		return nil, errJsonConvert
	}

	client := {{.package}}.New{{.rpcServiceName}}Client(m.cli.Conn())
	resp, err :=  client.{{.method}}(ctx, &request)
	if err != nil{
		return nil, err
	}

	var ret {{.pbResponse}}
	bts, err = jsonx.Marshal(resp)
	if err != nil{
		return nil, errJsonConvert
	}

	err = jsonx.Unmarshal(bts, &ret)
	if err != nil{
		return nil, errJsonConvert
	}

	return &ret, nil
}
`
)

func (g *defaultRpcGenerator) genCall() error {
	file := g.ast
	if len(file.Service) == 0 {
		return nil
	}
	if len(file.Service) > 1 {
		return fmt.Errorf("we recommend only one service in a proto, currently %d", len(file.Service))
	}

	typeCode, err := file.GenTypesCode()
	if err != nil {
		return err
	}

	constLit, err := file.GenEnumCode()
	if err != nil {
		return err
	}

	service := file.Service[0]
	callPath := filepath.Join(g.dirM[dirTarget], service.Name.Lower())
	if err = util.MkdirIfNotExist(callPath); err != nil {
		return err
	}

	filename := filepath.Join(callPath, typesFilename)
	head := util.GetHead(g.Ctx.ProtoSource)
	text, err := util.LoadTemplate(category, callTypesTemplateFile, callTemplateTypes)
	if err != nil {
		return err
	}
	err = util.With("types").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"head":                  head,
		"const":                 constLit,
		"filePackage":           service.Name.Lower(),
		"serviceName":           g.Ctx.ServiceName.Title(),
		"lowerStartServiceName": g.Ctx.ServiceName.UnTitle(),
		"types":                 typeCode,
	}, filename, true)
	if err != nil {
		return err
	}

	filename = filepath.Join(callPath, fmt.Sprintf("%s.go", service.Name.Lower()))
	functions, importList, err := g.genFunction(service)
	if err != nil {
		return err
	}

	iFunctions, err := g.getInterfaceFuncs(service)
	if err != nil {
		return err
	}
	text, err = util.LoadTemplate(category, callTemplateFile, callTemplateText)
	if err != nil {
		return err
	}
	err = util.With("shared").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"name":        service.Name.Lower(),
		"head":        head,
		"filePackage": service.Name.Lower(),
		"package":     strings.Join(importList, util.NL),
		"serviceName": service.Name.Title(),
		"functions":   strings.Join(functions, util.NL),
		"interface":   strings.Join(iFunctions, util.NL),
	}, filename, true)
	return err
}

func (g *defaultRpcGenerator) genFunction(service *parser.RpcService) ([]string, []string, error) {
	file := g.ast
	pkgName := file.Package
	functions := make([]string, 0)
	imports := collection.NewSet()
	imports.AddStr(fmt.Sprintf(`%v "%v"`, pkgName, g.mustGetPackage(dirPb)))
	for _, method := range service.Funcs {
		imports.AddStr(g.ast.Imports[method.ParameterIn.Package])
		text, err := util.LoadTemplate(category, callFunctionTemplateFile, callFunctionTemplate)
		if err != nil {
			return nil, nil, err
		}
		buffer, err := util.With("sharedFn").Parse(text).Execute(map[string]interface{}{
			"rpcServiceName": service.Name.Title(),
			"method":         method.Name.Title(),
			"package":        pkgName,
			"pbRequestName":  method.ParameterIn.Name,
			"pbRequest":      method.ParameterIn.Expression,
			"pbResponse":     method.ParameterOut.Name,
			"hasComment":     method.HaveDoc(),
			"comment":        method.GetDoc(),
		})
		if err != nil {
			return nil, nil, err
		}

		functions = append(functions, buffer.String())
	}
	return functions, imports.KeysStr(), nil
}

func (g *defaultRpcGenerator) getInterfaceFuncs(service *parser.RpcService) ([]string, error) {
	functions := make([]string, 0)

	for _, method := range service.Funcs {
		text, err := util.LoadTemplate(category, callInterfaceFunctionTemplateFile, callInterfaceFunctionTemplate)
		if err != nil {
			return nil, err
		}

		buffer, err := util.With("interfaceFn").Parse(text).Execute(
			map[string]interface{}{
				"hasComment": method.HaveDoc(),
				"comment":    method.GetDoc(),
				"method":     method.Name.Title(),
				"pbRequest":  method.ParameterIn.Name,
				"pbResponse": method.ParameterOut.Name,
			})
		if err != nil {
			return nil, err
		}

		functions = append(functions, buffer.String())
	}

	return functions, nil
}
