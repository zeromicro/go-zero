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
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

const regularPerm = 0o666

type RpcLogicData struct {
	LogicName string
	LogicCode string
}

// GenEntLogic generates the ent CRUD logic files of the rpc service.
func GenEntLogic(schema, output, serviceName, style, modelName string, multiple bool, searchKeyNum int) error {
	if !multiple {
		return genEntLogicInCompatibility(schema, output, serviceName, style, modelName, searchKeyNum)
	}

	return errors.New("does not support multiple")
	// Todo: in the future may add this function
	//return GenEntLogicGroup(ctx, proto, cfg)
}

func genEntLogicInCompatibility(schema, output, serviceName, style, modelName string, searchKeyNum int) error {
	outputDir, err := filepath.Abs(output)
	if err != nil {
		return err
	}

	logicDir := path.Join(outputDir, "internal/logic")

	schemas, err := entc.LoadGraph(schema, &gen.Config{})
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
		if modelName == s.Name || modelName == "" {
			// generate logic file
			rpcLogicData := GenCRUDData(serviceName, projectCtx, s, searchKeyNum)

			for _, v := range rpcLogicData {
				logicFilename, err := format.FileNamingFormat(style, v.LogicName)
				if err != nil {
					return err
				}

				filename := filepath.Join(logicDir, logicFilename+".go")
				if pathx.FileExists(filename) {
					continue
				}

				err = os.WriteFile(filename, []byte(v.LogicCode), regularPerm)
				if err != nil {
					return err
				}
			}

			// generate proto file
			protoMessage, protoFunctions, err := GenProtoData(s, searchKeyNum)
			if err != nil {
				return err
			}

			protoFileName := filepath.Join(outputDir, serviceName+".proto")
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
			newProtoData.WriteString(fmt.Sprintf("\n// %s message\n\n", modelName))
			newProtoData.WriteString(fmt.Sprintf("%s\n", protoMessage))
			newProtoData.WriteString(protoDataString[serviceIndex : len(protoDataString)-2])
			newProtoData.WriteString(fmt.Sprintf("\n\n  // %s management\n", modelName))
			newProtoData.WriteString(fmt.Sprintf("%s\n}", protoFunctions))

			err = os.WriteFile(protoFileName, []byte(newProtoData.String()), regularPerm)
			if err != nil {
				return err
			}

		}

	}
	return nil
}

func GenCRUDData(serviceName string, projectCtx *ctx.ProjectContext, schema *load.Schema, searchKeyNum int) []*RpcLogicData {
	var data []*RpcLogicData
	hasTime := false

	setLogic := strings.Builder{}
	for _, v := range schema.Fields {
		if v.Name == "id" || v.Name == "created_at" || v.Name == "updated_at" || v.Name == "deleted_at" {
			continue
		} else if v.Name == "status" {
			setLogic.WriteString(fmt.Sprintf("\t\t\tSet%s(uint8(in.%s)).\n", parser.CamelCase(v.Name),
				parser.CamelCase(v.Name)))
		} else {
			if strings.Contains(v.Name, "at") {
				hasTime = true
				setLogic.WriteString(fmt.Sprintf("\t\t\tSet%s(time.Unix(in.%s, 0)).\n", parser.CamelCase(v.Name),
					parser.CamelCase(v.Name)))
			} else if v.Name == "uuid" || v.Name == "api" {
				setLogic.WriteString(fmt.Sprintf("\t\t\tSet%s(in.%s).\n", strings.ToUpper(v.Name),
					parser.CamelCase(v.Name)))
			} else {
				setLogic.WriteString(fmt.Sprintf("\t\t\tSet%s(in.%s).\n", parser.CamelCase(v.Name),
					parser.CamelCase(v.Name)))
			}
		}
	}
	setLogic.WriteString("\t\t\tExec(l.ctx)")

	createLogic := bytes.NewBufferString("")
	createLogicTmpl, err := template.New("createOrUpdate").Parse(createOrUpdateTpl)
	err = createLogicTmpl.Execute(createLogic, map[string]interface{}{
		"hasTime":     hasTime,
		"setLogic":    setLogic.String(),
		"modelName":   schema.Name,
		"serviceName": serviceName,
		"projectPath": projectCtx.Path,
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
		if v.Info.Type.String() == "string" && !strings.Contains(strings.ToLower(v.Name), "uuid") && count <= searchKeyNum {
			camelName := parser.CamelCase(v.Name)
			predicateData.WriteString(fmt.Sprintf("\tif in.%s != \"\" {\n\t\tpredicates = append(predicates, %s.%sContains(in.%s))\n\t}\n",
				camelName, strings.ToLower(schema.Name), camelName, camelName))
			count++
		}
	}
	predicateData.WriteString(fmt.Sprintf("\tresult, err := l.svcCtx.DB.%s.Query().Where(predicates...).Page(l.ctx, in.Page, in.PageSize)",
		schema.Name))

	listData := strings.Builder{}

	for i, v := range schema.Fields {
		if v.Name == "id" || v.Name == "created_at" || v.Name == "updated_at" || v.Name == "deleted_at" {
			continue
		} else if v.Name == "status" {
			listData.WriteString(fmt.Sprintf("\t\t\t%s:\tuint64(v.%s),\n", parser.CamelCase(v.Name),
				parser.CamelCase(v.Name)))
		} else {
			if i < (len(schema.Fields) - 1) {
				if strings.Contains(v.Name, "at") {
					listData.WriteString(fmt.Sprintf("\t\t\t%s:\tv.%s.UnixMilli(),\n", parser.CamelCase(v.Name),
						parser.CamelCase(v.Name)))
				} else {
					listData.WriteString(fmt.Sprintf("\t\t\t%s:\tv.%s,\n", parser.CamelCase(v.Name),
						parser.CamelCase(v.Name)))
				}
			} else {
				if strings.Contains(v.Name, "at") {
					listData.WriteString(fmt.Sprintf("\t\t\t%s:\tv.%s.UnixMilli(),", parser.CamelCase(v.Name),
						parser.CamelCase(v.Name)))
				} else {
					listData.WriteString(fmt.Sprintf("\t\t\t%s:\tv.%s,", parser.CamelCase(v.Name),
						parser.CamelCase(v.Name)))
				}
			}
		}
	}

	getListLogic := bytes.NewBufferString("")
	getListLogicTmpl, err := template.New("getList").Parse(getListLogicTpl)
	getListLogicTmpl.Execute(getListLogic, map[string]interface{}{
		"predicateData": predicateData.String(),
		"modelName":     schema.Name,
		"listData":      listData.String(),
		"serviceName":   serviceName,
		"projectPath":   projectCtx.Path,
	})

	data = append(data, &RpcLogicData{
		LogicName: fmt.Sprintf("Get%sListLogic", schema.Name),
		LogicCode: getListLogic.String(),
	})

	deleteLogic := bytes.NewBufferString("")
	deleteLogicTmpl, err := template.New("delete").Parse(deleteLogicTpl)
	deleteLogicTmpl.Execute(deleteLogic, map[string]interface{}{
		"modelName":   schema.Name,
		"serviceName": serviceName,
		"projectPath": projectCtx.Path,
	})

	data = append(data, &RpcLogicData{
		LogicName: fmt.Sprintf("Delete%sLogic", schema.Name),
		LogicCode: deleteLogic.String(),
	})

	batchDeleteLogic := bytes.NewBufferString("")
	batchDeleteLogicTmpl, err := template.New("delete").Parse(batchDeleteLogicTpl)
	batchDeleteLogicTmpl.Execute(batchDeleteLogic, map[string]interface{}{
		"modelName":   schema.Name,
		"serviceName": serviceName,
		"projectPath": projectCtx.Path,
	})

	data = append(data, &RpcLogicData{
		LogicName: fmt.Sprintf("BatchDelete%sLogic", schema.Name),
		LogicCode: deleteLogic.String(),
	})

	return data
}

func GenProtoData(schema *load.Schema, searchKeyNum int) (string, string, error) {
	var protoMessage strings.Builder
	schemaNameCamelCase := parser.CamelCase(schema.Name)
	// info message
	protoMessage.WriteString(fmt.Sprintf("message %sInfo {\n  uint64 id = 1;\n  int64 created_at = 2;\n  int64 updated_at = 3;\n",
		schemaNameCamelCase))
	index := 4
	for i, v := range schema.Fields {
		if v.Name == "id" || v.Name == "created_at" || v.Name == "updated_at" || v.Name == "deleted_at" {
			continue
		} else if v.Name == "status" {
			protoMessage.WriteString(fmt.Sprintf("  uint64 %s = %d;\n", v.Name, index))
			index++
		} else {
			if i < (len(schema.Fields) - 1) {
				if strings.Contains(v.Name, "at") {
					protoMessage.WriteString(fmt.Sprintf("  int64  %s = %d;\n", v.Name, index))
				} else {
					protoMessage.WriteString(fmt.Sprintf("  %s %s = %d;\n", v.Info.Type.String(), v.Name, index))
				}
			} else {
				if strings.Contains(v.Name, "at") {
					protoMessage.WriteString(fmt.Sprintf("  int64  %s = %d;\n}\n\n", v.Name, index))
				} else {
					protoMessage.WriteString(fmt.Sprintf("  %s %s = %d;\n}\n\n", v.Info.Type.String(), v.Name, index))
				}
			}
			index++
		}
	}

	// List message
	protoMessage.WriteString(fmt.Sprintf("message %sListResp {\n  uint64 total = 1;\n  repeated %sInfo data = 2;\n}\n\n",
		schemaNameCamelCase, schemaNameCamelCase))

	// List Request message
	protoMessage.WriteString(fmt.Sprintf("message %sListReq {\n  uint64 page = 1;\n  uint64 page_size = 2;\n",
		schemaNameCamelCase))
	count := 0
	index = 3

	for i, v := range schema.Fields {
		if v.Info.Type.String() == "string" && !strings.Contains(strings.ToLower(v.Name), "uuid") && count <= searchKeyNum {
			if i < (len(schema.Fields)-1) && count < (searchKeyNum-1) {
				protoMessage.WriteString(fmt.Sprintf("  %s %s = %d;\n", v.Info.Type.String(), v.Name, index))
			} else {
				protoMessage.WriteString(fmt.Sprintf("  %s %s = %d;\n}\n", v.Info.Type.String(), v.Name, index))
			}
			index++
			count++
		}
	}

	protoRpcFunction := bytes.NewBufferString("")
	protoTmpl, err := template.New("proto").Parse(protoTpl)
	err = protoTmpl.Execute(protoRpcFunction, map[string]interface{}{
		"modelName": schema.Name,
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
//			if err = util.With("logic").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
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
