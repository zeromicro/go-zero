package generator

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/rpcv2/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

const (
	callTemplateText = `{{.head}}

//go:generate mockgen -destination ./{{.name}}_mock.go -package {{.filePackage}} -source $GOFILE

package {{.filePackage}}

import (
	"context"

	{{.package}}

	"github.com/tal-tech/go-zero/zrpc"
)

type (
	{{.alias}}

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

	callInterfaceFunctionTemplate = `{{if .hasComment}}{{.comment}}
{{end}}{{.method}}(ctx context.Context,in *{{.pbRequest}}) (*{{.pbResponse}},error)`

	callFunctionTemplate = `
{{if .hasComment}}{{.comment}}{{end}}
func (m *default{{.rpcServiceName}}) {{.method}}(ctx context.Context,in *{{.pbRequest}}) (*{{.pbResponse}}, error) {
	client := {{.package}}.New{{.rpcServiceName}}Client(m.cli.Conn())
	return client.{{.method}}(ctx, in)
}
`
)

func (g *defaultGenerator) GenCall(ctx DirContext, dir Dir, proto parser.Proto) error {
	service := proto.Service
	head := util.GetHead(proto.Name)

	filename := filepath.Join(dir.Filename, fmt.Sprintf("%s.go", formatFilename(service.Name)))
	functions, err := g.genFunction(proto.PbPackage, service)
	if err != nil {
		return err
	}

	iFunctions, err := g.getInterfaceFuncs(proto.PbPackage, service)
	if err != nil {
		return err
	}
	text, err := util.LoadTemplate(category, callTemplateFile, callTemplateText)
	if err != nil {
		return err
	}
	var alias []string
	for _, item := range service.RPC {
		alias = append(alias, fmt.Sprintf("%s = %s", item.RequestType, fmt.Sprintf("%s.%s", proto.PbPackage, item.RequestType)))
		alias = append(alias, fmt.Sprintf("%s = %s", item.ReturnsType, fmt.Sprintf("%s.%s", proto.PbPackage, item.ReturnsType)))
	}

	err = util.With("shared").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"name":        formatFilename(service.Name),
		"alias":       strings.Join(alias, util.NL),
		"head":        head,
		"filePackage": formatFilename(service.Name),
		"package":     fmt.Sprintf(`"%s"`, ctx.GetPb().Package),
		"serviceName": parser.CamelCase(service.Name),
		"functions":   strings.Join(functions, util.NL),
		"interface":   strings.Join(iFunctions, util.NL),
	}, filename, true)
	return err
}

func (g *defaultGenerator) genFunction(goPackage string, service parser.Service) ([]string, error) {
	functions := make([]string, 0)
	for _, rpc := range service.RPC {
		text, err := util.LoadTemplate(category, callFunctionTemplateFile, callFunctionTemplate)
		if err != nil {
			return nil, err
		}
		comment := parser.GetComment(rpc.Doc())
		buffer, err := util.With("sharedFn").Parse(text).Execute(map[string]interface{}{
			"rpcServiceName": stringx.From(service.Name).Title(),
			"method":         stringx.From(rpc.Name).Title(),
			"package":        goPackage,
			"pbRequest":      rpc.RequestType,
			"pbResponse":     rpc.ReturnsType,
			"hasComment":     len(comment) > 0,
			"comment":        comment,
		})
		if err != nil {
			return nil, err
		}

		functions = append(functions, buffer.String())
	}
	return functions, nil
}

func (g *defaultGenerator) getInterfaceFuncs(goPackage string, service parser.Service) ([]string, error) {
	functions := make([]string, 0)

	for _, rpc := range service.RPC {
		text, err := util.LoadTemplate(category, callInterfaceFunctionTemplateFile, callInterfaceFunctionTemplate)
		if err != nil {
			return nil, err
		}

		comment := parser.GetComment(rpc.Doc())
		buffer, err := util.With("interfaceFn").Parse(text).Execute(
			map[string]interface{}{
				"hasComment": len(comment) > 0,
				"comment":    comment,
				"method":     stringx.From(rpc.Name).Title(),
				"pbRequest":  rpc.RequestType,
				"pbResponse": rpc.ReturnsType,
			})
		if err != nil {
			return nil, err
		}

		functions = append(functions, buffer.String())
	}

	return functions, nil
}
