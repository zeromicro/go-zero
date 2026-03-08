package generator

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/core/collection"
	conf "github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"
)

const logicFunctionTemplate = `{{if .hasComment}}{{.comment}}{{end}}
func (l *{{.logicName}}) {{.method}} ({{if .hasReq}}in {{.request}}{{if .stream}},stream {{.streamBody}}{{end}}{{else}}stream {{.streamBody}}{{end}}) ({{if .hasReply}}{{.response}},{{end}} error) {
	// todo: add your logic here and delete this line
	
	return {{if .hasReply}}&{{.responseType}}{},{{end}} nil
}
`

//go:embed logic.tpl
var logicTemplate string

// GenLogic generates the logic file of the rpc service, which corresponds to the RPC definition items in proto.
func (g *Generator) GenLogic(ctx DirContext, proto parser.Proto, cfg *conf.Config,
	c *ZRpcContext) error {
	if !c.Multiple {
		return g.genLogicInCompatibility(ctx, proto, cfg)
	}

	return g.genLogicGroup(ctx, proto, cfg)
}

func (g *Generator) genLogicInCompatibility(ctx DirContext, proto parser.Proto,
	cfg *conf.Config) error {
	dir := ctx.GetLogic()
	service := proto.Service[0].Service.Name
	pkgMap := parser.BuildProtoPackageMap(proto.ImportedProtos)
	for _, rpc := range proto.Service[0].RPC {
		logicName := fmt.Sprintf("%sLogic", stringx.From(rpc.Name).ToCamel())
		logicFilename, err := format.FileNamingFormat(cfg.NamingFormat, rpc.Name+"_logic")
		if err != nil {
			return err
		}

		filename := filepath.Join(dir.Filename, logicFilename+".go")
		functions, err := g.genLogicFunction(service, proto.PbPackage, proto.GoPackage, logicName, rpc, pkgMap)
		if err != nil {
			return err
		}

		imports := collection.NewSet[string]()
		imports.Add(fmt.Sprintf(`"%v"`, ctx.GetSvc().Package))
		addLogicImports(imports, ctx.GetPb().Package, proto.PbPackage, proto.GoPackage, rpc, pkgMap)

		text, err := pathx.LoadTemplate(category, logicTemplateFileFile, logicTemplate)
		if err != nil {
			return err
		}
		err = util.With("logic").GoFmt(true).Parse(text).SaveTo(map[string]any{
			"logicName":   fmt.Sprintf("%sLogic", stringx.From(rpc.Name).ToCamel()),
			"functions":   functions,
			"packageName": "logic",
			"imports":     strings.Join(imports.Keys(), pathx.NL),
		}, filename, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) genLogicGroup(ctx DirContext, proto parser.Proto, cfg *conf.Config) error {
	dir := ctx.GetLogic()
	pkgMap := parser.BuildProtoPackageMap(proto.ImportedProtos)
	for _, item := range proto.Service {
		serviceName := item.Name
		for _, rpc := range item.RPC {
			var (
				err           error
				filename      string
				logicName     string
				logicFilename string
				packageName   string
			)

			logicName = fmt.Sprintf("%sLogic", stringx.From(rpc.Name).ToCamel())
			childPkg, err := dir.GetChildPackage(serviceName)
			if err != nil {
				return err
			}

			serviceDir := filepath.Base(childPkg)
			nameJoin := fmt.Sprintf("%s_logic", serviceName)
			packageName = strings.ToLower(stringx.From(nameJoin).ToCamel())
			logicFilename, err = format.FileNamingFormat(cfg.NamingFormat, rpc.Name+"_logic")
			if err != nil {
				return err
			}

			filename = filepath.Join(dir.Filename, serviceDir, logicFilename+".go")
			functions, err := g.genLogicFunction(serviceName, proto.PbPackage, proto.GoPackage, logicName, rpc, pkgMap)
			if err != nil {
				return err
			}

			imports := collection.NewSet[string]()
			imports.Add(fmt.Sprintf(`"%v"`, ctx.GetSvc().Package))
			addLogicImports(imports, ctx.GetPb().Package, proto.PbPackage, proto.GoPackage, rpc, pkgMap)

			text, err := pathx.LoadTemplate(category, logicTemplateFileFile, logicTemplate)
			if err != nil {
				return err
			}

			if err = util.With("logic").GoFmt(true).Parse(text).SaveTo(map[string]any{
				"logicName":   logicName,
				"functions":   functions,
				"packageName": packageName,
				"imports":     strings.Join(imports.Keys(), pathx.NL),
			}, filename, false); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *Generator) genLogicFunction(serviceName, goPackage, mainGoPackage, logicName string,
	rpc *parser.RPC, pkgMap map[string]parser.ImportedProto) (string, error) {
	functions := make([]string, 0)
	text, err := pathx.LoadTemplate(category, logicFuncTemplateFileFile, logicFunctionTemplate)
	if err != nil {
		return "", err
	}

	comment := parser.GetComment(rpc.Doc())
	streamServer := fmt.Sprintf("%s.%s_%s%s", goPackage, parser.CamelCase(serviceName),
		parser.CamelCase(rpc.Name), "Server")

	reqRef := resolveRPCTypeRef(rpc.RequestType, goPackage, mainGoPackage, pkgMap)
	respRef := resolveRPCTypeRef(rpc.ReturnsType, goPackage, mainGoPackage, pkgMap)

	buffer, err := util.With("fun").Parse(text).Execute(map[string]any{
		"logicName":    logicName,
		"method":       parser.CamelCase(rpc.Name),
		"hasReq":       !rpc.StreamsRequest,
		"request":      "*" + reqRef.GoRef,
		"hasReply":     !rpc.StreamsRequest && !rpc.StreamsReturns,
		"response":     "*" + respRef.GoRef,
		"responseType": respRef.GoRef,
		"stream":       rpc.StreamsRequest || rpc.StreamsReturns,
		"streamBody":   streamServer,
		"hasComment":   len(comment) > 0,
		"comment":      comment,
	})
	if err != nil {
		return "", err
	}

	functions = append(functions, buffer.String())
	return strings.Join(functions, pathx.NL), nil
}

// addLogicImports adds the correct import paths to imports for a single RPC's
// logic file. The main pb package is only included when it is actually referenced
// (i.e. when the request or response type lives in that package, or the RPC streams).
func addLogicImports(imports *collection.Set[string], pbImportPath, goPackage, mainGoPackage string,
rpc *parser.RPC, pkgMap map[string]parser.ImportedProto) {
// Streaming RPCs always reference the main pb package (for the stream type).
if rpc.StreamsRequest || rpc.StreamsReturns {
imports.Add(fmt.Sprintf(`"%s"`, pbImportPath))
return
}

reqRef := resolveRPCTypeRef(rpc.RequestType, goPackage, mainGoPackage, pkgMap)
respRef := resolveRPCTypeRef(rpc.ReturnsType, goPackage, mainGoPackage, pkgMap)

// Add main pb import if any type ref is from the main package (no extra import path).
if reqRef.ImportPath == "" || respRef.ImportPath == "" {
imports.Add(fmt.Sprintf(`"%s"`, pbImportPath))
}
if reqRef.ImportPath != "" {
imports.Add(fmt.Sprintf(`"%s"`, reqRef.ImportPath))
}
if respRef.ImportPath != "" {
imports.Add(fmt.Sprintf(`"%s"`, respRef.ImportPath))
}
}
