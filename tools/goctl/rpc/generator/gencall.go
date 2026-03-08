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
	pkgMap := parser.BuildProtoPackageMap(proto.ImportedProtos)
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
		alias := collection.NewSet[string]()
		var hasSameNameBetweenMessageAndService bool
		for _, item := range proto.Message {
			msgName := getMessageName(*item.Message)
			if serviceName == msgName {
				hasSameNameBetweenMessageAndService = true
			}
			if !isCallPkgSameToPbPkg {
				alias.Add(fmt.Sprintf("%s = %s", parser.CamelCase(msgName),
					fmt.Sprintf("%s.%s", proto.PbPackage, parser.CamelCase(msgName))))
			}
		}
		if hasSameNameBetweenMessageAndService {
			serviceName = stringx.From(service.Name + "_zrpc_client").ToCamel()
		}

		extraImports := collection.NewSet[string]()
		functions, err := g.genFunction(proto.PbPackage, proto.GoPackage, serviceName, service, isCallPkgSameToGrpcPkg, pkgMap, alias, extraImports)
		if err != nil {
			return err
		}

		iFunctions, err := g.getInterfaceFuncs(proto.PbPackage, proto.GoPackage, service, isCallPkgSameToGrpcPkg, pkgMap, extraImports)
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

		extraImportLines := buildExtraImportLines(extraImports)
		aliasKeys := alias.Keys()
		sort.Strings(aliasKeys)
		if err = util.With("shared").GoFmt(true).Parse(text).SaveTo(map[string]any{
			"name":           callFilename,
			"alias":          strings.Join(aliasKeys, pathx.NL),
			"head":           head,
			"filePackage":    childDir,
			"pbPackage":      pbPackage,
			"protoGoPackage": protoGoPackage,
			"extraImports":   extraImportLines,
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
	alias := collection.NewSet[string]()
	var hasSameNameBetweenMessageAndService bool
	for _, item := range proto.Message {
		msgName := getMessageName(*item.Message)
		if serviceName == msgName {
			hasSameNameBetweenMessageAndService = true
		}
		if !isCallPkgSameToPbPkg {
			alias.Add(fmt.Sprintf("%s = %s", parser.CamelCase(msgName),
				fmt.Sprintf("%s.%s", proto.PbPackage, parser.CamelCase(msgName))))
		}
	}

	if hasSameNameBetweenMessageAndService {
		serviceName = stringx.From(service.Name + "_zrpc_client").ToCamel()
	}

	pkgMap := parser.BuildProtoPackageMap(proto.ImportedProtos)
	extraImports := collection.NewSet[string]()
	filename := filepath.Join(dir.Filename, fmt.Sprintf("%s.go", callFilename))
	functions, err := g.genFunction(proto.PbPackage, proto.GoPackage, serviceName, service, isCallPkgSameToGrpcPkg, pkgMap, alias, extraImports)
	if err != nil {
		return err
	}

	iFunctions, err := g.getInterfaceFuncs(proto.PbPackage, proto.GoPackage, service, isCallPkgSameToGrpcPkg, pkgMap, extraImports)
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
	extraImportLines := buildExtraImportLines(extraImports)
	aliasKeys := alias.Keys()
	sort.Strings(aliasKeys)
	return util.With("shared").GoFmt(true).Parse(text).SaveTo(map[string]any{
		"name":           callFilename,
		"alias":          strings.Join(aliasKeys, pathx.NL),
		"head":           head,
		"filePackage":    dir.Base,
		"pbPackage":      pbPackage,
		"protoGoPackage": protoGoPackage,
		"extraImports":   extraImportLines,
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

func (g *Generator) genFunction(goPackage, mainGoPackage, serviceName string, service parser.Service,
	isCallPkgSameToGrpcPkg bool, pkgMap map[string]parser.ImportedProto,
	alias, extraImports *collection.Set[string]) ([]string, error) {
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

		reqName, reqAlias, reqImport := resolveCallTypeRef(rpc.RequestType, goPackage, mainGoPackage, pkgMap)
		respName, respAlias, respImport := resolveCallTypeRef(rpc.ReturnsType, goPackage, mainGoPackage, pkgMap)
		if reqAlias != "" {
			alias.Add(reqAlias)
		}
		if respAlias != "" {
			alias.Add(respAlias)
		}
		if reqImport != "" {
			extraImports.Add(reqImport)
		}
		if respImport != "" {
			extraImports.Add(respImport)
		}

		buffer, err := util.With("sharedFn").Parse(text).Execute(map[string]any{
			"serviceName":            serviceName,
			"rpcServiceName":         parser.CamelCase(service.Name),
			"method":                 parser.CamelCase(rpc.Name),
			"package":                goPackage,
			"pbRequest":              reqName,
			"pbResponse":             respName,
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

func (g *Generator) getInterfaceFuncs(goPackage, mainGoPackage string, service parser.Service,
	isCallPkgSameToGrpcPkg bool, pkgMap map[string]parser.ImportedProto,
	extraImports *collection.Set[string]) ([]string, error) {
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

		reqName, _, reqImport := resolveCallTypeRef(rpc.RequestType, goPackage, mainGoPackage, pkgMap)
		respName, _, respImport := resolveCallTypeRef(rpc.ReturnsType, goPackage, mainGoPackage, pkgMap)
		if reqImport != "" {
			extraImports.Add(reqImport)
		}
		if respImport != "" {
			extraImports.Add(respImport)
		}

		buffer, err := util.With("interfaceFn").Parse(text).Execute(
			map[string]any{
				"hasComment": len(comment) > 0,
				"comment":    comment,
				"method":     parser.CamelCase(rpc.Name),
				"hasReq":     !rpc.StreamsRequest,
				"pbRequest":  reqName,
				"notStream":  !rpc.StreamsRequest && !rpc.StreamsReturns,
				"pbResponse": respName,
				"streamBody": streamServer,
			})
		if err != nil {
			return nil, err
		}

		functions = append(functions, buffer.String())
	}

	return functions, nil
}

// buildExtraImportLines converts a set of import paths into quoted import lines
// for use in the call.tpl {{.extraImports}} placeholder.
func buildExtraImportLines(extraImports *collection.Set[string]) string {
if extraImports.Count() == 0 {
return ""
}
keys := extraImports.Keys()
sort.Strings(keys)
lines := make([]string, 0, len(keys))
for _, k := range keys {
lines = append(lines, fmt.Sprintf(`"%s"`, k))
}
return strings.Join(lines, "\n\t")
}
