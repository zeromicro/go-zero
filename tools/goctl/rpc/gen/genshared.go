package gogen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

const (
	sharedTemplateText = `{{.head}}

//go:generate mockgen -destination ./{{.name}}model_mock.go -package {{.filePackage}} -source $GOFILE

package {{.filePackage}}

import (
	"context"

	{{.package}}

	"github.com/tal-tech/go-zero/core/jsonx"
	"github.com/tal-tech/go-zero/rpcx"
)

type (
	{{.serviceName}}Model interface {
		{{.interface}}
	}

	default{{.serviceName}}Model struct {
		cli rpcx.Client
	}
)

func New{{.serviceName}}Model(cli rpcx.Client) {{.serviceName}}Model {
	return &default{{.serviceName}}Model{
		cli: cli,
	}
}

{{.functions}}
`
	sharedTemplateTypes = `{{.head}}

package {{.filePackage}}

import "errors"

var errJsonConvert = errors.New("json convert error")

{{.types}}
`
	sharedInterfaceFunctionTemplate = `{{if .hasComment}}{{.comment}}
{{end}}{{.method}}(ctx context.Context,in *{{.pbRequest}}) {{if .hasResponse}}(*{{.pbResponse}},{{end}} error{{if .hasResponse}}){{end}}`
	sharedFunctionTemplate = `
{{if .hasComment}}{{.comment}}{{end}}
func (m *default{{.rpcServiceName}}Model) {{.method}}(ctx context.Context,in *{{.pbRequest}}) {{if .hasResponse}}(*{{.pbResponse}},{{end}} error{{if .hasResponse}}){{end}} {
	client := {{.package}}.New{{.rpcServiceName}}Client(m.cli.Conn())
	var request {{.package}}.{{.pbRequest}}
	bts, err := jsonx.Marshal(in)
	if err != nil {
		return {{if .hasResponse}}nil, {{end}}errJsonConvert
	}

	err = jsonx.Unmarshal(bts, &request)
	if err != nil {
		return {{if .hasResponse}}nil, {{end}}errJsonConvert
	}

	{{if .hasResponse}}resp, err := {{else}}_, err = {{end}}client.{{.method}}(ctx, &request)
	{{if .hasResponse}}if err != nil{
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

	return &ret, nil{{else}}if err != nil {
		return err
	}

	return nil{{end}}
}
`
)

func (g *defaultRpcGenerator) genShared() error {
	sharePackage := filepath.Base(g.Ctx.SharedDir)
	file := g.ast
	typeCode, err := file.GenTypesCode()
	if err != nil {
		return err
	}

	pbPkg := file.Package
	remotePackage := fmt.Sprintf(`%v "%v"`, pbPkg, g.mustGetPackage(dirPb))
	filename := filepath.Join(g.Ctx.SharedDir, "types.go")
	head := util.GetHead(g.Ctx.ProtoSource)
	err = util.With("types").GoFmt(true).Parse(sharedTemplateTypes).SaveTo(map[string]interface{}{
		"head":                  head,
		"filePackage":           sharePackage,
		"pbPkg":                 pbPkg,
		"serviceName":           g.Ctx.ServiceName.Title(),
		"lowerStartServiceName": g.Ctx.ServiceName.UnTitle(),
		"types":                 typeCode,
	}, filename, true)

	for _, service := range file.Service {
		filename := filepath.Join(g.Ctx.SharedDir, fmt.Sprintf("%smodel.go", service.Name.Lower()))
		functions, err := g.getFuncs(service)
		if err != nil {
			return err
		}
		iFunctions, err := g.getInterfaceFuncs(service)
		if err != nil {
			return err
		}
		mockFile := filepath.Join(g.Ctx.SharedDir, fmt.Sprintf("%smodel_mock.go", service.Name.Lower()))
		os.Remove(mockFile)
		err = util.With("shared").GoFmt(true).Parse(sharedTemplateText).SaveTo(map[string]interface{}{
			"name":        service.Name.Lower(),
			"head":        head,
			"filePackage": sharePackage,
			"pbPkg":       pbPkg,
			"package":     remotePackage,
			"serviceName": service.Name.Title(),
			"functions":   strings.Join(functions, "\n"),
			"interface":   strings.Join(iFunctions, "\n"),
		}, filename, true)
		if err != nil {
			return err
		}
	}

	// if mockgen is already installed, it will generate code of gomock for shared files
	_, err = exec.LookPath("mockgen")
	if err != nil {
		g.Ctx.Warning("warning:mockgen is not found")
	} else {
		execx.Run(fmt.Sprintf("cd %s \ngo generate", g.Ctx.SharedDir))
	}
	return nil
}

func (g *defaultRpcGenerator) getFuncs(service *parser.RpcService) ([]string, error) {
	file := g.ast
	pkgName := file.Package
	functions := make([]string, 0)
	for _, method := range service.Funcs {
		data, found := file.Strcuts[strings.ToLower(method.OutType)]
		if found {
			found = len(data.Field) > 0
		}
		var comment string
		if len(method.Document) > 0 {
			comment = method.Document[0]
		}
		buffer, err := util.With("sharedFn").Parse(sharedFunctionTemplate).Execute(map[string]interface{}{
			"rpcServiceName": service.Name.Title(),
			"method":         method.Name.Title(),
			"package":        pkgName,
			"pbRequest":      method.InType,
			"pbResponse":     method.OutType,
			"hasResponse":    found,
			"hasComment":     len(method.Document) > 0,
			"comment":        comment,
		})
		if err != nil {
			return nil, err
		}

		functions = append(functions, buffer.String())
	}
	return functions, nil
}

func (g *defaultRpcGenerator) getInterfaceFuncs(service *parser.RpcService) ([]string, error) {
	file := g.ast
	functions := make([]string, 0)
	for _, method := range service.Funcs {
		data, found := file.Strcuts[strings.ToLower(method.OutType)]
		if found {
			found = len(data.Field) > 0
		}
		var comment string
		if len(method.Document) > 0 {
			comment = method.Document[0]
		}
		buffer, err := util.With("interfaceFn").Parse(sharedInterfaceFunctionTemplate).Execute(map[string]interface{}{
			"hasComment":  len(method.Document) > 0,
			"comment":     comment,
			"method":      method.Name.Title(),
			"pbRequest":   method.InType,
			"pbResponse":  method.OutType,
			"hasResponse": found,
		})
		if err != nil {
			return nil, err
		}

		functions = append(functions, buffer.String())
	}
	return functions, nil
}
