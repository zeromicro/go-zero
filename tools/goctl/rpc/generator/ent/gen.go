package ent

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/entc/load"
	"github.com/iancoleman/strcase"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/ctx"
	"github.com/zeromicro/go-zero/tools/goctl/util/entx"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"github.com/zeromicro/go-zero/tools/goctl/util/protox"
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
	ProjectName  string
	Style        string
	ModelName    string
	Multiple     bool
	SearchKeyNum int
	ModuleName   string
	GroupName    string
	UseUUID      bool
	ProtoOut     string
}

// GenEntLogic generates the ent CRUD logic files of the rpc service.
func GenEntLogic(g *GenEntLogicContext) error {
	return genEntLogic(g)
}

func genEntLogic(g *GenEntLogicContext) error {
	outputDir, err := filepath.Abs(g.Output)
	if err != nil {
		return err
	}

	var logicDir string

	if g.Multiple {
		logicDir = path.Join(outputDir, "internal/logic", g.ServiceName)
		err = pathx.MkdirIfNotExist(logicDir)
		if err != nil {
			return err
		}
	} else {
		logicDir = path.Join(outputDir, "internal/logic")
	}

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

			var protoFileName string
			if g.ProtoOut == "" {
				protoFileName = filepath.Join(outputDir, g.ProjectName+".proto")
				if !pathx.FileExists(protoFileName) {
					continue
				}
			} else {
				protoFileName, err = filepath.Abs(g.ProtoOut)
				if err != nil {
					return err
				}
				if !pathx.FileExists(protoFileName) {
					err = os.WriteFile(protoFileName, []byte(fmt.Sprintf("syntax = \"proto3\";\n\nservice %s {\n}",
						strcase.ToCamel(g.ServiceName))), os.ModePerm)
					if err != nil {
						return fmt.Errorf("fail to create proto file : %s", err.Error())
					}
				}
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
			serviceBeginIndex, _, serviceEndIndex := protox.FindBeginEndOfService(protoDataString, strcase.ToCamel(g.ServiceName))
			if serviceBeginIndex == -1 {
				continue
			}
			newProtoData.WriteString(protoDataString[:serviceBeginIndex-1])
			newProtoData.WriteString(fmt.Sprintf("\n// %s message\n\n", g.ModelName))
			newProtoData.WriteString(fmt.Sprintf("%s\n", protoMessage))
			newProtoData.WriteString(protoDataString[serviceBeginIndex-1 : serviceEndIndex-1])
			newProtoData.WriteString(fmt.Sprintf("\n\n  // %s management\n", g.ModelName))
			newProtoData.WriteString(fmt.Sprintf("%s\n", protoFunctions))
			newProtoData.WriteString(protoDataString[serviceEndIndex-1:])

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
	hasTime, hasUUID := false, false
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
			if v.Name == "id" && entx.IsUUIDType(v.Info.Type.String()) {
				g.UseUUID = true
			}
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
	createLogicTmpl, _ := template.New("create").Parse(createTpl)
	createLogicTmpl.Execute(createLogic, map[string]any{
		"hasTime":     hasTime,
		"hasUUID":     hasUUID,
		"setLogic":    setLogic.String(),
		"modelName":   schema.Name,
		"projectName": g.ProjectName,
		"projectPath": projectCtx.Path,
		"packageName": packageName,
		"useUUID":     g.UseUUID, // UUID primary key
	})

	data = append(data, &RpcLogicData{
		LogicName: fmt.Sprintf("Create%sLogic", schema.Name),
		LogicCode: createLogic.String(),
	})

	updateLogic := bytes.NewBufferString("")
	updateLogicTmpl, _ := template.New("update").Parse(updateTpl)
	updateLogicTmpl.Execute(updateLogic, map[string]any{
		"hasTime":     hasTime,
		"hasUUID":     hasUUID,
		"setLogic":    strings.Replace(setLogic.String(), "Set", "SetNotEmpty", -1),
		"modelName":   schema.Name,
		"projectName": g.ProjectName,
		"projectPath": projectCtx.Path,
		"packageName": packageName,
		"useUUID":     g.UseUUID, // UUID primary key
	})

	data = append(data, &RpcLogicData{
		LogicName: fmt.Sprintf("Update%sLogic", schema.Name),
		LogicCode: updateLogic.String(),
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
				listData.WriteString(fmt.Sprintf("\t\t\t%s:\tuint32(v.%s),%s", nameCamelCase,
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
	getListLogicTmpl, _ := template.New("getList").Parse(getListLogicTpl)
	getListLogicTmpl.Execute(getListLogic, map[string]any{
		"predicateData":      predicateData.String(),
		"modelName":          schema.Name,
		"listData":           listData.String(),
		"projectName":        g.ProjectName,
		"projectPath":        projectCtx.Path,
		"modelNameLowerCase": strings.ToLower(schema.Name),
		"packageName":        packageName,
		"useUUID":            g.UseUUID,
	})

	data = append(data, &RpcLogicData{
		LogicName: fmt.Sprintf("Get%sListLogic", schema.Name),
		LogicCode: getListLogic.String(),
	})

	getByIdLogic := bytes.NewBufferString("")
	getByIdLogicTmpl, _ := template.New("getById").Parse(getByIdLogicTpl)
	getByIdLogicTmpl.Execute(getByIdLogic, map[string]any{
		"modelName":          schema.Name,
		"listData":           strings.Replace(listData.String(), "v.", "result.", -1),
		"projectName":        g.ProjectName,
		"projectPath":        projectCtx.Path,
		"modelNameLowerCase": strings.ToLower(schema.Name),
		"packageName":        packageName,
		"useUUID":            g.UseUUID,
	})

	data = append(data, &RpcLogicData{
		LogicName: fmt.Sprintf("Get%sByIdLogic", schema.Name),
		LogicCode: getByIdLogic.String(),
	})

	deleteLogic := bytes.NewBufferString("")
	deleteLogicTmpl, _ := template.New("delete").Parse(deleteLogicTpl)
	deleteLogicTmpl.Execute(deleteLogic, map[string]any{
		"modelName":          schema.Name,
		"modelNameLowerCase": strings.ToLower(schema.Name),
		"projectName":        g.ProjectName,
		"projectPath":        projectCtx.Path,
		"packageName":        packageName,
		"useUUID":            g.UseUUID,
	})

	data = append(data, &RpcLogicData{
		LogicName: fmt.Sprintf("Delete%sLogic", schema.Name),
		LogicCode: deleteLogic.String(),
	})

	return data
}

func GenProtoData(schema *load.Schema, g *GenEntLogicContext) (string, string, error) {
	var protoMessage strings.Builder
	schemaNameCamelCase := parser.CamelCase(schema.Name)
	// hasStatus means it has status field
	hasStatus := false
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
			protoMessage.WriteString(fmt.Sprintf("  uint32 %s = %d;\n", v.Name, index))
			hasStatus = true
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
	protoMessage.WriteString(fmt.Sprintf("message %sListReq {\n  uint64 page = 1;\n  uint64 page_size = 2;\n",
		schemaNameCamelCase))
	count := 0
	index = 3

	for i, v := range schema.Fields {
		if v.Info.Type.String() == "string" && !strings.Contains(strings.ToLower(v.Name), "uuid") && count < g.SearchKeyNum {
			if i < len(schema.Fields) && count < g.SearchKeyNum {
				protoMessage.WriteString(fmt.Sprintf("  %s %s = %d;\n", entx.ConvertEntTypeToProtoType(v.Info.Type.String()),
					v.Name, index))
				index++
				count++
			}
		}

		if i == (len(schema.Fields) - 1) {
			protoMessage.WriteString("}\n")
		}
	}

	// group
	var groupName string
	if g.GroupName != "" {
		groupName = fmt.Sprintf("  // group: %s\n", g.GroupName)
	} else {
		groupName = ""
	}

	protoRpcFunction := bytes.NewBufferString("")
	protoTmpl, err := template.New("proto").Parse(protoTpl)
	err = protoTmpl.Execute(protoRpcFunction, map[string]any{
		"modelName": schema.Name,
		"groupName": groupName,
		"useUUID":   g.UseUUID,
		"hasStatus": hasStatus,
	})

	if err != nil {
		logx.Error(err)
		return "", "", err
	}

	return protoMessage.String(), protoRpcFunction.String(), nil
}
