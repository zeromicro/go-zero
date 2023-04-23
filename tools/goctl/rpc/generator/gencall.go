package generator

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/emicklei/proto"
	"github.com/zeromicro/go-zero/core/collection"
	conf "github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"
)

const (
	callInterfaceFunctionTemplate = `{{if .hasComment}}{{.comment}}
{{end}}{{.method}}(ctx context.Context{{if .hasReq}}, in *{{.pbRequest}}{{end}}, opts ...grpc.CallOption) ({{if .notStream}}*{{.pbResponse}}, {{else}}{{.streamBody}},{{end}} error)`

	callFunctionTemplate = `
{{if .hasComment}}{{.comment}}{{end}}
func (m *default{{.serviceName}}) {{.method}}(ctx context.Context{{if .hasReq}}, in *{{.pbRequest}}{{end}}, opts ...grpc.CallOption) ({{if .notStream}}*{{.pbResponse}}, {{else}}{{.streamBody}},{{end}} error) {
	client := {{if .isCallPkgSameToGrpcPkg}}{{else}}{{.package}}.{{end}}New{{.rpcServiceName}}Client(m.cli.Conn())
	return client.{{.method}}(ctx{{if .hasReq}}, in{{end}}, opts...)
}
`
)

//go:embed call.tpl
var callTemplateText string

// GenCall generates the rpc client code, which is the entry point for the rpc service call.
// It is a layer of encapsulation for the rpc client and shields the details in the pb.
func (g *Generator) GenCall(ctx DirContext, proto parser.Proto, cfg *conf.Config,
	c *ZRpcContext) error {
	if !c.Multiple {
		return g.genCallInCompatibility(ctx, proto, cfg)
	}

	return g.genCallGroup(ctx, proto, cfg)
}

func (g *Generator) genCallGroup(ctx DirContext, proto parser.Proto, cfg *conf.Config) error {
	dir := ctx.GetCall()
	head := util.GetHead(proto.Name)
	for _, service := range proto.Service {
		childPkg, err := dir.GetChildPackage(service.Name)
		if err != nil {
			return err
		}

		callFilename, err := format.FileNamingFormat(cfg.NamingFormat, service.Name)
		if err != nil {
			return err
		}

		childDir := filepath.Base(childPkg)
		filename := filepath.Join(dir.Filename, childDir, fmt.Sprintf("%s.go", callFilename))
		isCallPkgSameToPbPkg := childDir == ctx.GetProtoGo().Filename
		isCallPkgSameToGrpcPkg := childDir == ctx.GetProtoGo().Filename

		serviceName := stringx.From(service.Name).ToCamel()
		alias := collection.NewSet()
		var hasSameNameBetweenMessageAndService bool
		for _, item := range proto.Message {
			msgName := getMessageName(*item.Message)
			if serviceName == msgName {
				hasSameNameBetweenMessageAndService = true
			}
			if !isCallPkgSameToPbPkg {
				alias.AddStr(fmt.Sprintf("%s = %s", parser.CamelCase(msgName),
					fmt.Sprintf("%s.%s", proto.PbPackage, parser.CamelCase(msgName))))
			}
		}
		if hasSameNameBetweenMessageAndService {
			serviceName = stringx.From(service.Name + "_zrpc_client").ToCamel()
		}

		functions, err := g.genFunction(proto.PbPackage, serviceName, service, isCallPkgSameToGrpcPkg)
		if err != nil {
			return err
		}

		iFunctions, err := g.getInterfaceFuncs(proto.PbPackage, service, isCallPkgSameToGrpcPkg)
		if err != nil {
			return err
		}

		text, err := pathx.LoadTemplate(category, callTemplateFile, callTemplateText)
		if err != nil {
			return err
		}

		pbPackage := fmt.Sprintf(`"%s"`, ctx.GetPb().Package)
		protoGoPackage := fmt.Sprintf(`"%s"`, ctx.GetProtoGo().Package)
		if isCallPkgSameToGrpcPkg {
			pbPackage = ""
			protoGoPackage = ""
		}

		aliasKeys := alias.KeysStr()
		sort.Strings(aliasKeys)
		if err = util.With("shared").GoFmt(true).Parse(text).SaveTo(map[string]any{
			"name":           callFilename,
			"alias":          strings.Join(aliasKeys, pathx.NL),
			"head":           head,
			"filePackage":    childDir,
			"pbPackage":      pbPackage,
			"protoGoPackage": protoGoPackage,
			"serviceName":    serviceName,
			"functions":      strings.Join(functions, pathx.NL),
			"interface":      strings.Join(iFunctions, pathx.NL),
		}, filename, true); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) genCallInCompatibility(ctx DirContext, proto parser.Proto,
	cfg *conf.Config) error {
	dir := ctx.GetCall()
	service := proto.Service[0]
	head := util.GetHead(proto.Name)
	isCallPkgSameToPbPkg := ctx.GetCall().Filename == ctx.GetPb().Filename
	isCallPkgSameToGrpcPkg := ctx.GetCall().Filename == ctx.GetProtoGo().Filename

	callFilename, err := format.FileNamingFormat(cfg.NamingFormat, service.Name)
	if err != nil {
		return err
	}

	serviceName := stringx.From(service.Name).ToCamel()
	alias := collection.NewSet()
	var hasSameNameBetweenMessageAndService bool
	for _, item := range proto.Message {
		msgName := getMessageName(*item.Message)
		if serviceName == msgName {
			hasSameNameBetweenMessageAndService = true
		}
		if !isCallPkgSameToPbPkg {
			alias.AddStr(fmt.Sprintf("%s = %s", parser.CamelCase(msgName),
				fmt.Sprintf("%s.%s", proto.PbPackage, parser.CamelCase(msgName))))
		}
	}

	if hasSameNameBetweenMessageAndService {
		serviceName = stringx.From(service.Name + "_zrpc_client").ToCamel()
	}

	filename := filepath.Join(dir.Filename, fmt.Sprintf("%s.go", callFilename))
	functions, err := g.genFunction(proto.PbPackage, serviceName, service, isCallPkgSameToGrpcPkg)
	if err != nil {
		return err
	}

	iFunctions, err := g.getInterfaceFuncs(proto.PbPackage, service, isCallPkgSameToGrpcPkg)
	if err != nil {
		return err
	}

	text, err := pathx.LoadTemplate(category, callTemplateFile, callTemplateText)
	if err != nil {
		return err
	}

	pbPackage := fmt.Sprintf(`"%s"`, ctx.GetPb().Package)
	protoGoPackage := fmt.Sprintf(`"%s"`, ctx.GetProtoGo().Package)
	if isCallPkgSameToGrpcPkg {
		pbPackage = ""
		protoGoPackage = ""
	}
	aliasKeys := alias.KeysStr()
	sort.Strings(aliasKeys)
	return util.With("shared").GoFmt(true).Parse(text).SaveTo(map[string]any{
		"name":           callFilename,
		"alias":          strings.Join(aliasKeys, pathx.NL),
		"head":           head,
		"filePackage":    dir.Base,
		"pbPackage":      pbPackage,
		"protoGoPackage": protoGoPackage,
		"serviceName":    serviceName,
		"functions":      strings.Join(functions, pathx.NL),
		"interface":      strings.Join(iFunctions, pathx.NL),
	}, filename, true)
}

func getMessageName(msg proto.Message) string {
	list := []string{msg.Name}

	for {
		parent := msg.Parent
		if parent == nil {
			break
		}

		parentMsg, ok := parent.(*proto.Message)
		if !ok {
			break
		}

		tmp := []string{parentMsg.Name}
		list = append(tmp, list...)
		msg = *parentMsg
	}

	return strings.Join(list, "_")
}

func (g *Generator) genFunction(goPackage string, serviceName string, service parser.Service,
	isCallPkgSameToGrpcPkg bool) ([]string, error) {
	functions := make([]string, 0)

	for _, rpc := range service.RPC {
		text, err := pathx.LoadTemplate(category, callFunctionTemplateFile, callFunctionTemplate)
		if err != nil {
			return nil, err
		}

		comment := parser.GetComment(rpc.Doc())
		streamServer := fmt.Sprintf("%s.%s_%s%s", goPackage, parser.CamelCase(service.Name),
			parser.CamelCase(rpc.Name), "Client")
		if isCallPkgSameToGrpcPkg {
			streamServer = fmt.Sprintf("%s_%s%s", parser.CamelCase(service.Name),
				parser.CamelCase(rpc.Name), "Client")
		}
		buffer, err := util.With("sharedFn").Parse(text).Execute(map[string]any{
			"serviceName":            serviceName,
			"rpcServiceName":         parser.CamelCase(service.Name),
			"method":                 parser.CamelCase(rpc.Name),
			"package":                goPackage,
			"pbRequest":              parser.CamelCase(rpc.RequestType),
			"pbResponse":             parser.CamelCase(rpc.ReturnsType),
			"hasComment":             len(comment) > 0,
			"comment":                comment,
			"hasReq":                 !rpc.StreamsRequest,
			"notStream":              !rpc.StreamsRequest && !rpc.StreamsReturns,
			"streamBody":             streamServer,
			"isCallPkgSameToGrpcPkg": isCallPkgSameToGrpcPkg,
		})
		if err != nil {
			return nil, err
		}

		functions = append(functions, buffer.String())
	}

	return functions, nil
}

func (g *Generator) getInterfaceFuncs(goPackage string, service parser.Service,
	isCallPkgSameToGrpcPkg bool) ([]string, error) {
	functions := make([]string, 0)

	for _, rpc := range service.RPC {
		text, err := pathx.LoadTemplate(category, callInterfaceFunctionTemplateFile,
			callInterfaceFunctionTemplate)
		if err != nil {
			return nil, err
		}

		comment := parser.GetComment(rpc.Doc())
		streamServer := fmt.Sprintf("%s.%s_%s%s", goPackage, parser.CamelCase(service.Name),
			parser.CamelCase(rpc.Name), "Client")
		if isCallPkgSameToGrpcPkg {
			streamServer = fmt.Sprintf("%s_%s%s", parser.CamelCase(service.Name),
				parser.CamelCase(rpc.Name), "Client")
		}
		buffer, err := util.With("interfaceFn").Parse(text).Execute(
			map[string]any{
				"hasComment": len(comment) > 0,
				"comment":    comment,
				"method":     parser.CamelCase(rpc.Name),
				"hasReq":     !rpc.StreamsRequest,
				"pbRequest":  parser.CamelCase(rpc.RequestType),
				"notStream":  !rpc.StreamsRequest && !rpc.StreamsReturns,
				"pbResponse": parser.CamelCase(rpc.ReturnsType),
				"streamBody": streamServer,
			})
		if err != nil {
			return nil, err
		}

		functions = append(functions, buffer.String())
	}

	return functions, nil
}
