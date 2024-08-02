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

const functionTemplate = `
{{if .hasComment}}{{.comment}}{{end}}
func (s *{{.server}}Server) {{.method}} ({{if .notStream}}ctx context.Context,{{if .hasReq}} in {{.request}}{{end}}{{else}}{{if .hasReq}} in {{.request}},{{end}}stream {{.streamBody}}{{end}}) ({{if .notStream}}{{.response}},{{end}}error) {
	l := {{.logicPkg}}.New{{.logicName}}({{if .notStream}}ctx,{{else}}stream.Context(),{{end}}s.svcCtx)
	return l.{{.method}}({{if .hasReq}}in{{if .stream}} ,stream{{end}}{{else}}{{if .stream}}stream{{end}}{{end}})
}
`

//go:embed server.tpl
var serverTemplate string

// GenServer generates rpc server file, which is an implementation of rpc server
func (g *Generator) GenServer(ctx DirContext, proto parser.Proto, cfg *conf.Config,
	c *ZRpcContext) error {
	if !c.Multiple {
		return g.genServerInCompatibility(ctx, proto, cfg, c)
	}

	return g.genServerGroup(ctx, proto, cfg)
}

func (g *Generator) genServerGroup(ctx DirContext, proto parser.Proto, cfg *conf.Config) error {
	dir := ctx.GetServer()
	for _, service := range proto.Service {
		var (
			serverFile  string
			logicImport string
		)

		serverFilename, err := format.FileNamingFormat(cfg.NamingFormat, service.Name+"_server")
		if err != nil {
			return err
		}

		serverChildPkg, err := dir.GetChildPackage(service.Name)
		if err != nil {
			return err
		}

		logicChildPkg, err := ctx.GetLogic().GetChildPackage(service.Name)
		if err != nil {
			return err
		}

		serverDir := filepath.Base(serverChildPkg)
		logicImport = fmt.Sprintf(`"%v"`, logicChildPkg)
		serverFile = filepath.Join(dir.Filename, serverDir, serverFilename+".go")

		svcImport := fmt.Sprintf(`"%v"`, ctx.GetSvc().Package)
		pbImport := fmt.Sprintf(`"%v"`, ctx.GetPb().Package)

		imports := collection.NewSet()
		imports.AddStr(logicImport, svcImport, pbImport)

		head := util.GetHead(proto.Name)

		funcList, err := g.genFunctions(proto.PbPackage, service, true)
		if err != nil {
			return err
		}

		text, err := pathx.LoadTemplate(category, serverTemplateFile, serverTemplate)
		if err != nil {
			return err
		}

		notStream := false
		for _, rpc := range service.RPC {
			if !rpc.StreamsRequest && !rpc.StreamsReturns {
				notStream = true
				break
			}
		}

		if err = util.With("server").GoFmt(true).Parse(text).SaveTo(map[string]any{
			"head": head,
			"unimplementedServer": fmt.Sprintf("%s.Unimplemented%sServer", proto.PbPackage,
				parser.CamelCase(service.Name)),
			"server":    stringx.From(service.Name).ToCamel(),
			"imports":   strings.Join(imports.KeysStr(), pathx.NL),
			"funcs":     strings.Join(funcList, pathx.NL),
			"notStream": notStream,
		}, serverFile, true); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) genServerInCompatibility(ctx DirContext, proto parser.Proto,
	cfg *conf.Config, c *ZRpcContext) error {
	dir := ctx.GetServer()
	logicImport := fmt.Sprintf(`"%v"`, ctx.GetLogic().Package)
	svcImport := fmt.Sprintf(`"%v"`, ctx.GetSvc().Package)
	pbImport := fmt.Sprintf(`"%v"`, ctx.GetPb().Package)

	imports := collection.NewSet()
	imports.AddStr(logicImport, svcImport, pbImport)

	head := util.GetHead(proto.Name)
	service := proto.Service[0]
	serverFilename, err := format.FileNamingFormat(cfg.NamingFormat, service.Name+"_server")
	if err != nil {
		return err
	}

	serverFile := filepath.Join(dir.Filename, serverFilename+".go")
	funcList, err := g.genFunctions(proto.PbPackage, service, false)
	if err != nil {
		return err
	}

	text, err := pathx.LoadTemplate(category, serverTemplateFile, serverTemplate)
	if err != nil {
		return err
	}

	notStream := false
	for _, rpc := range service.RPC {
		if !rpc.StreamsRequest && !rpc.StreamsReturns {
			notStream = true
			break
		}
	}

	return util.With("server").GoFmt(true).Parse(text).SaveTo(map[string]any{
		"head": head,
		"unimplementedServer": fmt.Sprintf("%s.Unimplemented%sServer", proto.PbPackage,
			parser.CamelCase(service.Name)),
		"server":    stringx.From(service.Name).ToCamel(),
		"imports":   strings.Join(imports.KeysStr(), pathx.NL),
		"funcs":     strings.Join(funcList, pathx.NL),
		"notStream": notStream,
	}, serverFile, true)
}

func (g *Generator) genFunctions(goPackage string, service parser.Service, multiple bool) ([]string, error) {
	var (
		functionList []string
		logicPkg     string
	)
	for _, rpc := range service.RPC {
		text, err := pathx.LoadTemplate(category, serverFuncTemplateFile, functionTemplate)
		if err != nil {
			return nil, err
		}

		var logicName string
		if !multiple {
			logicPkg = "logic"
			logicName = fmt.Sprintf("%sLogic", stringx.From(rpc.Name).ToCamel())
		} else {
			nameJoin := fmt.Sprintf("%s_logic", service.Name)
			logicPkg = strings.ToLower(stringx.From(nameJoin).ToCamel())
			logicName = fmt.Sprintf("%sLogic", stringx.From(rpc.Name).ToCamel())
		}

		comment := parser.GetComment(rpc.Doc())
		streamServer := fmt.Sprintf("%s.%s_%s%s", goPackage, parser.CamelCase(service.Name),
			parser.CamelCase(rpc.Name), "Server")
		buffer, err := util.With("func").Parse(text).Execute(map[string]any{
			"server":     stringx.From(service.Name).ToCamel(),
			"logicName":  logicName,
			"method":     parser.CamelCase(rpc.Name),
			"request":    fmt.Sprintf("*%s.%s", goPackage, parser.CamelCase(rpc.RequestType)),
			"response":   fmt.Sprintf("*%s.%s", goPackage, parser.CamelCase(rpc.ReturnsType)),
			"hasComment": len(comment) > 0,
			"comment":    comment,
			"hasReq":     !rpc.StreamsRequest,
			"stream":     rpc.StreamsRequest || rpc.StreamsReturns,
			"notStream":  !rpc.StreamsRequest && !rpc.StreamsReturns,
			"streamBody": streamServer,
			"logicPkg":   logicPkg,
		})
		if err != nil {
			return nil, err
		}

		functionList = append(functionList, buffer.String())
	}
	return functionList, nil
}
