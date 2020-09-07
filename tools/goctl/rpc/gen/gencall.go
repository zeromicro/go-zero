package gen

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
	callTemplateText = `{{.head}}

//go:generate mockgen -destination ./{{.name}}_mock.go -package {{.filePackage}} -source $GOFILE

package {{.filePackage}}

import (
	"context"

	{{.package}}

	"github.com/tal-tech/go-zero/core/jsonx"
	"github.com/tal-tech/go-zero/rpcx"
)

type (
	{{.serviceName}} interface {
		{{.interface}}
	}

	default{{.serviceName}} struct {
		cli rpcx.Client
	}
)

func New{{.serviceName}}(cli rpcx.Client) {{.serviceName}} {
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

{{.types}}
`
	callInterfaceFunctionTemplate = `{{if .hasComment}}{{.comment}}
{{end}}{{.method}}(ctx context.Context,in *{{.pbRequest}}) {{if .hasResponse}}(*{{.pbResponse}},{{end}} error{{if .hasResponse}}){{end}}`
	callFunctionTemplate = `
{{if .hasComment}}{{.comment}}{{end}}
func (m *default{{.rpcServiceName}}) {{.method}}(ctx context.Context,in *{{.pbRequest}}) {{if .hasResponse}}(*{{.pbResponse}},{{end}} error{{if .hasResponse}}){{end}} {
	var request {{.package}}.{{.pbRequest}}
	bts, err := jsonx.Marshal(in)
	if err != nil {
		return {{if .hasResponse}}nil, {{end}}errJsonConvert
	}

	err = jsonx.Unmarshal(bts, &request)
	if err != nil {
		return {{if .hasResponse}}nil, {{end}}errJsonConvert
	}

	client := {{.package}}.New{{.rpcServiceName}}Client(m.cli.Conn())
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

	service := file.Service[0]
	callPath, err := filepath.Abs(service.Name.Lower())
	if err != nil {
		return err
	}

	if err = util.MkdirIfNotExist(callPath); err != nil {
		return err
	}

	pbPkg := file.Package
	remotePackage := fmt.Sprintf(`%v "%v"`, pbPkg, g.mustGetPackage(dirPb))
	filename := filepath.Join(callPath, "types.go")
	head := util.GetHead(g.Ctx.ProtoSource)
	err = util.With("types").GoFmt(true).Parse(callTemplateTypes).SaveTo(map[string]interface{}{
		"head":                  head,
		"filePackage":           service.Name.Lower(),
		"pbPkg":                 pbPkg,
		"serviceName":           g.Ctx.ServiceName.Title(),
		"lowerStartServiceName": g.Ctx.ServiceName.UnTitle(),
		"types":                 typeCode,
	}, filename, true)
	if err != nil {
		return err
	}

	_, err = exec.LookPath("mockgen")
	mockGenInstalled := err == nil
	filename = filepath.Join(callPath, fmt.Sprintf("%s.go", service.Name.Lower()))
	functions, err := g.getFuncs(service)
	if err != nil {
		return err
	}

	iFunctions, err := g.getInterfaceFuncs(service)
	if err != nil {
		return err
	}

	mockFile := filepath.Join(callPath, fmt.Sprintf("%s_mock.go", service.Name.Lower()))
	os.Remove(mockFile)
	err = util.With("shared").GoFmt(true).Parse(callTemplateText).SaveTo(map[string]interface{}{
		"name":        service.Name.Lower(),
		"head":        head,
		"filePackage": service.Name.Lower(),
		"pbPkg":       pbPkg,
		"package":     remotePackage,
		"serviceName": service.Name.Title(),
		"functions":   strings.Join(functions, "\n"),
		"interface":   strings.Join(iFunctions, "\n"),
	}, filename, true)
	if err != nil {
		return err
	}
	// if mockgen is already installed, it will generate code of gomock for shared files
	_, err = exec.LookPath("mockgen")
	if mockGenInstalled {
		execx.Run(fmt.Sprintf("go generate %s", filename))
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
		buffer, err := util.With("sharedFn").Parse(callFunctionTemplate).Execute(map[string]interface{}{
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
		buffer, err := util.With("interfaceFn").Parse(callInterfaceFunctionTemplate).Execute(
			map[string]interface{}{
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
