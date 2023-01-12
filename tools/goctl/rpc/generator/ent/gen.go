package ent

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/entc/load"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/ctx"
	"github.com/zeromicro/go-zero/tools/goctl/util/entx"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

const regularPerm = 0o666

type RpcLogicData struct {
	LogicName string
	LogicCode string
}

type GenEntLogicContext struct {
	Schema       string
	Output       string
	ServiceName  string
	Style        string
	ModelName    string
	Multiple     bool
	SearchKeyNum int
	ModuleName   string
	GroupName    string
	UseUUID      bool
}

// GenEntLogic generates the ent CRUD logic files of the rpc service.
func GenEntLogic(g *GenEntLogicContext) error {
	if !g.Multiple {
		return genEntLogicInCompatibility(g)
	}

	return errors.New("does not support multiple")
	// Todo: in the future may add this function
	// return GenEntLogicGroup(ctx, proto, cfg)
}

func genEntLogicInCompatibility(g *GenEntLogicContext) error {
	outputDir, err := filepath.Abs(g.Output)
	if err != nil {
		return err
	}

	logicDir := path.Join(outputDir, "internal/logic")

	schemas, err := entc.LoadGraph(g.Schema, &gen.Config{})
	if err != nil {
		return err
	}

	workDir, err := filepath.Abs("./")
	if err != nil {
		return err
	}

	projectCtx, err := ctx.Prepare(workDir)
	if err != nil {
		return err
	}

	for _, s := range schemas.Schemas {
		if g.ModelName == s.Name || g.ModelName == "" {
			// generate logic file
			rpcLogicData := GenCRUDData(g, projectCtx, s)

			for _, v := range rpcLogicData {
				logicFilename, err := format.FileNamingFormat(g.Style, v.LogicName)
				if err != nil {
					return err
				}

				// group
				var filename string
				if g.GroupName != "" {
					if err = pathx.MkdirIfNotExist(filepath.Join(logicDir, g.GroupName)); err != nil {
						return err
					}

					filename = filepath.Join(logicDir, g.GroupName, logicFilename+".go")
				} else {
					filename = filepath.Join(logicDir, logicFilename+".go")
				}

				if pathx.FileExists(filename) {
					continue
				}

				err = os.WriteFile(filename, []byte(v.LogicCode), regularPerm)
				if err != nil {
					return err
				}
			}

			// generate proto file
			protoMessage, protoFunctions, err := GenProtoData(s, g)
			if err != nil {
				return err
			}

			protoFileName := filepath.Join(outputDir, g.ServiceName+".proto")
			if !pathx.FileExists(protoFileName) {
				continue
			}

			protoFileData, err := os.ReadFile(protoFileName)
			if err != nil {
				return err
			}

			protoDataString := string(protoFileData)

			if strings.Contains(protoDataString, protoMessage) || strings.Contains(protoDataString, protoFunctions) {
				continue
			}

			// generate new proto file
			newProtoData := strings.Builder{}
			serviceIndex := strings.Index(protoDataString, "service")
			if serviceIndex == -1 {
				continue
			}
			newProtoData.WriteString(protoDataString[:serviceIndex])
			newProtoData.WriteString(fmt.Sprintf("\n// %s message\n\n", g.ModelName))
			newProtoData.WriteString(fmt.Sprintf("%s\n", protoMessage))
			newProtoData.WriteString(protoDataString[serviceIndex : len(protoDataString)-2])
			newProtoData.WriteString(fmt.Sprintf("\n  // %s management\n", g.ModelName))
			newProtoData.WriteString(fmt.Sprintf("%s\n}", protoFunctions))

			err = os.WriteFile(protoFileName, []byte(newProtoData.String()), regularPerm)
			if err != nil {
				return err
			}

		}
	}
	return nil
}

func GenCRUDData(g *GenEntLogicContext, projectCtx *ctx.ProjectContext, schema *load.Schema) []*RpcLogicData {
	var data []*RpcLogicData
	hasTime := false
	hasUUID := false
	// end string means whether to use \n
	endString := ""
	var packageName string
	if g.GroupName != "" {
		packageName = g.GroupName
	} else {
		packageName = "logic"
	}

	setLogic := strings.Builder{}
	for _, v := range schema.Fields {
		if entx.IsBaseProperty(v.Name) {
			continue
		} else if v.Name == "status" {
			setLogic.WriteString(fmt.Sprintf("\t\t\tSet%s(uint8(in.%s)).\n", parser.CamelCase(v.Name),
				parser.CamelCase(v.Name)))
		} else {
			if entx.IsTimeProperty(v.Name) {
				hasTime = true
				setLogic.WriteString(fmt.Sprintf("\t\t\tSet%s(time.Unix(in.%s, 0)).\n", parser.CamelCase(v.Name),
					parser.CamelCase(v.Name)))
			} else if entx.IsUpperProperty(v.Name) {
				if entx.IsGoTypeNotPrototype(v.Info.Type.String()) {
					if v.Info.Type.String() == "[16]byte" {
						setLogic.WriteString(fmt.Sprintf("\t\t\tSet%s(uuidx.ParseUUIDString(in.%s)).\n", entx.ConvertSpecificNounToUpper(v.Name),
							parser.CamelCase(v.Name)))
						hasUUID = true
					} else {
						setLogic.WriteString(fmt.Sprintf("\t\t\tSet%s(%s(in.%s)).\n", entx.ConvertSpecificNounToUpper(v.Name),
							v.Info.Type.String(), parser.CamelCase(v.Name)))
					}
				} else {
					setLogic.WriteString(fmt.Sprintf("\t\t\tSet%s(in.%s).\n", entx.ConvertSpecificNounToUpper(v.Name),
						parser.CamelCase(v.Name)))
				}
			} else {
				if entx.IsGoTypeNotPrototype(v.Info.Type.String()) {
					setLogic.WriteString(fmt.Sprintf("\t\t\tSet%s(%s(in.%s)).\n", parser.CamelCase(v.Name),
						v.Info.Type.String(), parser.CamelCase(v.Name)))
				} else {
					setLogic.WriteString(fmt.Sprintf("\t\t\tSet%s(in.%s).\n", parser.CamelCase(v.Name),
						parser.CamelCase(v.Name)))
				}
			}
		}
	}
	setLogic.WriteString("\t\t\tExec(l.ctx)")

	createLogic := bytes.NewBufferString("")
	createLogicTmpl, err := template.New("createOrUpdate").Parse(createOrUpdateTpl)
	err = createLogicTmpl.Execute(createLogic, map[string]any{
		"hasTime":     hasTime,
		"hasUUID":     hasUUID,
		"setLogic":    setLogic.String(),
		"modelName":   schema.Name,
		"serviceName": g.ServiceName,
		"projectPath": projectCtx.Path,
		"packageName": packageName,
		"useUUID":     g.UseUUID, // UUID primary key
	})

	if err != nil {
		logx.Error(err)
		return nil
	}

	data = append(data, &RpcLogicData{
		LogicName: fmt.Sprintf("CreateOrUpdate%sLogic", schema.Name),
		LogicCode: createLogic.String(),
	})

	predicateData := strings.Builder{}
	predicateData.WriteString(fmt.Sprintf("\tvar predicates []predicate.%s\n", schema.Name))
	count := 0
	for _, v := range schema.Fields {
		if v.Info.Type.String() == "string" && !strings.Contains(strings.ToLower(v.Name), "uuid") &&
			count < g.SearchKeyNum && !entx.IsBaseProperty(v.Name) {
			camelName := parser.CamelCase(v.Name)
			predicateData.WriteString(fmt.Sprintf("\tif in.%s != \"\" {\n\t\tpredicates = append(predicates, %s.%sContains(in.%s))\n\t}\n",
				camelName, strings.ToLower(schema.Name), entx.ConvertSpecificNounToUpper(v.Name), camelName))
			count++
		}
	}
	predicateData.WriteString(fmt.Sprintf("\tresult, err := l.svcCtx.DB.%s.Query().Where(predicates...).Page(l.ctx, in.Page, in.PageSize)",
		schema.Name))

	listData := strings.Builder{}

	for i, v := range schema.Fields {
		if entx.IsBaseProperty(v.Name) {
			continue
		} else {
			nameCamelCase := parser.CamelCase(v.Name)

			if i < (len(schema.Fields) - 1) {
				endString = "\n"
			} else {
				endString = ""
			}

			if entx.IsUUIDType(v.Info.Type.String()) {
				listData.WriteString(fmt.Sprintf("\t\t\t%s:\tv.%s.String(),%s", nameCamelCase,
					entx.ConvertSpecificNounToUpper(nameCamelCase), endString))
			} else if v.Name == "status" {
				listData.WriteString(fmt.Sprintf("\t\t\t%s:\tuint64(v.%s),%s", nameCamelCase,
					entx.ConvertSpecificNounToUpper(nameCamelCase), endString))
			} else if entx.IsTimeProperty(v.Name) {
				listData.WriteString(fmt.Sprintf("\t\t\t%s:\tv.%s.UnixMilli(),%s", nameCamelCase,
					entx.ConvertSpecificNounToUpper(nameCamelCase), endString))
			} else {
				if entx.IsUpperProperty(v.Name) {
					listData.WriteString(fmt.Sprintf("\t\t\t%s:\tv.%s,%s", nameCamelCase,
						entx.ConvertSpecificNounToUpper(v.Name), endString))
				} else {
					if entx.IsGoTypeNotPrototype(v.Info.Type.String()) {
						listData.WriteString(fmt.Sprintf("\t\t\t%s:\t%s(v.%s),%s", nameCamelCase,
							entx.ConvertEntTypeToGotype(v.Info.Type.String()), nameCamelCase, endString))
					} else {
						listData.WriteString(fmt.Sprintf("\t\t\t%s:\tv.%s,%s", nameCamelCase,
							nameCamelCase, endString))
					}
				}
			}
		}
	}

	getListLogic := bytes.NewBufferString("")
	getListLogicTmpl, err := template.New("getList").Parse(getListLogicTpl)
	getListLogicTmpl.Execute(getListLogic, map[string]any{
		"predicateData":      predicateData.String(),
		"modelName":          schema.Name,
		"listData":           listData.String(),
		"serviceName":        g.ServiceName,
		"projectPath":        projectCtx.Path,
		"modelNameLowerCase": strings.ToLower(schema.Name),
		"packageName":        packageName,
		"useUUID":            g.UseUUID,
	})

	data = append(data, &RpcLogicData{
		LogicName: fmt.Sprintf("Get%sListLogic", schema.Name),
		LogicCode: getListLogic.String(),
	})

	deleteLogic := bytes.NewBufferString("")
	deleteLogicTmpl, err := template.New("delete").Parse(deleteLogicTpl)
	deleteLogicTmpl.Execute(deleteLogic, map[string]any{
		"modelName":   schema.Name,
		"serviceName": g.ServiceName,
		"projectPath": projectCtx.Path,
		"packageName": packageName,
		"useUUID":     g.UseUUID,
	})

	data = append(data, &RpcLogicData{
		LogicName: fmt.Sprintf("Delete%sLogic", schema.Name),
		LogicCode: deleteLogic.String(),
	})

	batchDeleteLogic := bytes.NewBufferString("")
	batchDeleteLogicTmpl, err := template.New("batchDelete").Parse(batchDeleteLogicTpl)
	batchDeleteLogicTmpl.Execute(batchDeleteLogic, map[string]any{
		"modelName":          schema.Name,
		"serviceName":        g.ServiceName,
		"projectPath":        projectCtx.Path,
		"modelNameLowerCase": strings.ToLower(schema.Name),
		"packageName":        packageName,
		"useUUID":            g.UseUUID,
	})

	data = append(data, &RpcLogicData{
		LogicName: fmt.Sprintf("BatchDelete%sLogic", schema.Name),
		LogicCode: batchDeleteLogic.String(),
	})

	return data
}

func GenProtoData(schema *load.Schema, g *GenEntLogicContext) (string, string, error) {
	var protoMessage strings.Builder
	schemaNameCamelCase := parser.CamelCase(schema.Name)
	// end string means whether to use \n
	endString := ""
	// info message
	protoMessage.WriteString(fmt.Sprintf("message %sInfo {\n  %s id = 1;\n  int64 created_at = 2;\n  int64 updated_at = 3;\n",
		schemaNameCamelCase, entx.ConvertIDType(g.UseUUID)))
	index := 4
	for i, v := range schema.Fields {
		if entx.IsBaseProperty(v.Name) {
			continue
		} else if v.Name == "status" {
			protoMessage.WriteString(fmt.Sprintf("  uint64 %s = %d;\n", v.Name, index))
			index++
		} else {
			if i < (len(schema.Fields) - 1) {
				endString = "\n"
			} else {
				endString = ""
			}

			if entx.IsTimeProperty(v.Name) {
				protoMessage.WriteString(fmt.Sprintf("  int64  %s = %d;%s", v.Name, index, endString))
			} else {
				protoMessage.WriteString(fmt.Sprintf("  %s %s = %d;%s", entx.ConvertEntTypeToProtoType(v.Info.Type.String()),
					v.Name, index, endString))
			}

			if i == (len(schema.Fields) - 1) {
				protoMessage.WriteString("\n}\n\n")
			}

			index++
		}
	}

	// List message
	protoMessage.WriteString(fmt.Sprintf("message %sListResp {\n  uint64 total = 1;\n  repeated %sInfo data = 2;\n}\n\n",
		schemaNameCamelCase, schemaNameCamelCase))

	// List Request message
	protoMessage.WriteString(fmt.Sprintf("message %sPageReq {\n  uint64 page = 1;\n  uint64 page_size = 2;\n",
		schemaNameCamelCase))
	count := 0
	index = 3

	for i, v := range schema.Fields {
		if v.Info.Type.String() == "string" && !strings.Contains(strings.ToLower(v.Name), "uuid") && count < g.SearchKeyNum {
			if i < (len(schema.Fields)-1) && count < g.SearchKeyNum {
				protoMessage.WriteString(fmt.Sprintf("  %s %s = %d;\n", entx.ConvertEntTypeToProtoType(v.Info.Type.String()),
					v.Name, index))
			}
			index++
			count++
		}

		if i == (len(schema.Fields) - 1) {
			protoMessage.WriteString("}\n")
		}
	}

	// group
	var groupName string
	if g.GroupName != "" {
		groupName = fmt.Sprintf("  // group: %s\n", groupName)
	} else {
		groupName = ""
	}

	protoRpcFunction := bytes.NewBufferString("")
	protoTmpl, err := template.New("proto").Parse(protoTpl)
	err = protoTmpl.Execute(protoRpcFunction, map[string]any{
		"modelName": schema.Name,
		"groupName": groupName,
		"useUUID":   g.UseUUID,
	})

	if err != nil {
		logx.Error(err)
		return "", "", err
	}

	return protoMessage.String(), protoRpcFunction.String(), nil
}

// Todo: in the future
//func GenEntLogicGroup(ctx DirContext, proto parser.Proto, cfg *conf.Config) error {
//	dir := ctx.GetLogic()
//	for _, item := range proto.Service {
//		serviceName := item.Name
//		for _, rpc := range item.RPC {
//			var (
//				err           error
//				filename      string
//				logicName     string
//				logicFilename string
//				packageName   string
//			)
//
//			logicName = fmt.Sprintf("%sLogic", stringx.From(rpc.Name).ToCamel())
//			childPkg, err := dir.GetChildPackage(serviceName)
//			if err != nil {
//				return err
//			}
//
//			serviceDir := filepath.Base(childPkg)
//			nameJoin := fmt.Sprintf("%s_logic", serviceName)
//			packageName = strings.ToLower(stringx.From(nameJoin).ToCamel())
//			logicFilename, err = format.FileNamingFormat(cfg.NamingFormat, rpc.Name+"_logic")
//			if err != nil {
//				return err
//			}
//
//			filename = filepath.Join(dir.Filename, serviceDir, logicFilename+".go")
//			functions, err := g.genLogicFunction(serviceName, proto.PbPackage, logicName, rpc)
//			if err != nil {
//				return err
//			}
//
//			imports := collection.NewSet()
//			imports.AddStr(fmt.Sprintf(`"%v"`, ctx.GetSvc().Package))
//			imports.AddStr(fmt.Sprintf(`"%v"`, ctx.GetPb().Package))
//			text, err := pathx.LoadTemplate(category, logicTemplateFileFile, logicTemplate)
//			if err != nil {
//				return err
//			}
//
//			if err = util.With("logic").GoFmt(true).Parse(text).SaveTo(map[string]any{
//				"logicName":   logicName,
//				"functions":   functions,
//				"packageName": packageName,
//				"imports":     strings.Join(imports.KeysStr(), pathx.NL),
//			}, filename, false); err != nil {
//				return err
//			}
//		}
//	}
//	return nil
//}
//
