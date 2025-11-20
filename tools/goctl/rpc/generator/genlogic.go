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
	// service := proto.Service[0].Service.Name
	for _, rpc := range proto.Service[0].RPC {
		logicName := fmt.Sprintf("%sLogic", stringx.From(rpc.Name).ToCamel())
		logicFilename, err := format.FileNamingFormat(cfg.NamingFormat, rpc.Name+"_logic")
		if err != nil {
			return err
		}

		filename := filepath.Join(dir.Filename, logicFilename+".go")
		functions, err := g.genLogicFunction(proto, logicName, rpc)
		if err != nil {
			return err
		}

		imports := collection.NewSet[string]()
		imports.Add(fmt.Sprintf(`"%v"`, ctx.GetSvc().Package))
		funcImports, err := getImports(proto, rpc)
		if err != nil {
			return err
		}
		imports.Add(funcImports...)
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
			functions, err := g.genLogicFunction(proto, logicName, rpc)
			if err != nil {
				return err
			}

			imports := collection.NewSet[string]()
			imports.Add(fmt.Sprintf(`"%v"`, ctx.GetSvc().Package))
			funcImports, err := getImports(proto, rpc)
			if err != nil {
				return err
			}
			imports.Add(funcImports...)
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

func (g *Generator) genLogicFunction(proto parser.Proto, logicName string, rpc *parser.RPC) (string,
	error) {
	serviceName := proto.Service[0].Service.Name
	goPackage := proto.PbPackage
	functions := make([]string, 0)
	text, err := pathx.LoadTemplate(category, logicFuncTemplateFileFile, logicFunctionTemplate)
	if err != nil {
		return "", err
	}

	comment := parser.GetComment(rpc.Doc())
	requestMsg, existed := proto.GetImportMessage(rpc.RequestType)
	if !existed {
		err = fmt.Errorf("request type %s is invalid", rpc.RequestType)
		return "", err
	}
	responseMsg, existed := proto.GetImportMessage(rpc.ReturnsType)
	if !existed {
		err = fmt.Errorf("response type %s is invalid", rpc.ReturnsType)
		return "", err
	}

	// streamServer := fmt.Sprintf("%s.%s_%s%s", goPackage, parser.CamelCase(serviceName),
	// 	parser.CamelCase(rpc.Name), "Server")
	streamServer := buildPackageArg(goPackage, fmt.Sprintf("%s_%s%s", parser.CamelCase(serviceName), parser.CamelCase(rpc.Name), "Server"), true)
	buffer, err := util.With("fun").Parse(text).Execute(map[string]any{
		"logicName":    logicName,
		"method":       parser.CamelCase(rpc.Name),
		"hasReq":       !rpc.StreamsRequest,
		"request":      buildPackageArg(requestMsg.PbPackage, requestMsg.Name, true),
		"hasReply":     !rpc.StreamsRequest && !rpc.StreamsReturns,
		"response":     buildPackageArg(responseMsg.PbPackage, responseMsg.Name, true),
		"responseType": buildPackageArg(responseMsg.PbPackage, responseMsg.Name, false),
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

func getImports(proto parser.Proto, rpc *parser.RPC) ([]string,
	error) {
	var err error
	requestMsg, existed := proto.GetImportMessage(rpc.RequestType)
	if !existed {
		err = fmt.Errorf("request type %s is invalid", rpc.RequestType)
		return nil, err
	}
	responseMsg, existed := proto.GetImportMessage(rpc.ReturnsType)
	if !existed {
		err = fmt.Errorf("response type %s is invalid", rpc.ReturnsType)
		return nil, err
	}
	imports := make([]string, 0)

	imports = append(imports,
		fmt.Sprintf(`"%s"`, requestMsg.GoPackage),
		fmt.Sprintf(`"%s"`, responseMsg.GoPackage),
	)
	return imports, nil
}

func getImport(packageName string, proto parser.Proto) string {
	var importPackage string
	if packageName == proto.PbPackage {
		importPackage = proto.GoPackage
		return importPackage
	}
	for _, item := range proto.Import {
		if item.Proto.PbPackage == packageName {
			importPackage = item.Proto.GoPackage
			break
		}
	}
	return importPackage
}

func buildPackageArg(goPackage string, arg string, isPtr bool) string {
	if strings.Contains(arg, ".") {
		arg = fmt.Sprintf("%s", arg)
	} else {
		arg = fmt.Sprintf("%s.%s", goPackage, parser.CamelCase(arg))
	}

	if isPtr {
		arg = fmt.Sprintf("*%s", arg)
	}

	return arg
}
