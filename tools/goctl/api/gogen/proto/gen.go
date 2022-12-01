package proto

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"entgo.io/ent/entc/load"
	"github.com/emicklei/proto"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/ctx"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

const regularPerm = 0o666

var protoField *protoFieldData

type protoFieldData struct {
	Name string
	Type string
}

type ApiLogicData struct {
	LogicName string
	LogicCode string
}

// GenLogicByProtoContext describe the data used for logic generation with proto file
type GenLogicByProtoContext struct {
	ProtoDir     string
	OutputDir    string
	ServiceName  string
	Style        string
	ModelName    string
	SearchKeyNum int
	RpcName      string
	GrpcPackage  string
}

func GenLogicByProto(p *GenLogicByProtoContext) error {
	outputDir, err := filepath.Abs(p.OutputDir)
	if err != nil {
		return err
	}

	logicDir := path.Join(outputDir, "internal/logic")

	protoParser := parser.NewDefaultProtoParser()
	protoData, err := protoParser.Parse(p.ProtoDir, false)
	if err != nil {
		return err
	}

	protoField = &protoFieldData{}

	workDir, err := filepath.Abs("./")
	if err != nil {
		return err
	}

	projectCtx, err := ctx.Prepare(workDir)
	if err != nil {
		return err
	}

	// generate logic file
	rpcLogicData := GenCRUDData(p, &protoData, projectCtx)

	for _, v := range rpcLogicData {
		logicFilename, err := format.FileNamingFormat(p.Style, v.LogicName)
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

	//// generate api file
	//apiMessage, apiFunctions, err := GenApiData(s, searchKeyNum)
	//if err != nil {
	//	return err
	//}
	//
	//apiFileName := filepath.Join(outputDir, serviceName+".api")
	//if !pathx.FileExists(apiFileName) {
	//	continue
	//}
	//
	//apiFileData, err := os.ReadFile(apiFileName)
	//if err != nil {
	//	return err
	//}
	//
	//apiDataString := string(apiFileData)
	//
	//if strings.Contains(apiDataString, apiMessage) || strings.Contains(apiDataString, apiFunctions) {
	//	continue
	//}
	//
	//// generate new api file
	//newProtoData := strings.Builder{}
	//serviceIndex := strings.Index(apiDataString, "service")
	//if serviceIndex == -1 {
	//	continue
	//}
	//newProtoData.WriteString(apiDataString[:serviceIndex])
	//newProtoData.WriteString(fmt.Sprintf("\n// %s message\n\n", modelName))
	//newProtoData.WriteString(fmt.Sprintf("%s\n", apiMessage))
	//newProtoData.WriteString(apiDataString[serviceIndex : len(apiDataString)-2])
	//newProtoData.WriteString(fmt.Sprintf("\n  // %s management\n", modelName))
	//newProtoData.WriteString(fmt.Sprintf("%s\n}", apiFunctions))
	//
	//err = os.WriteFile(apiFileName, []byte(newProtoData.String()), regularPerm)
	//if err != nil {
	//	return err
	//}

	return nil
}

type MessageVisitor struct {
	proto.NoopVisitor
}

func (m MessageVisitor) VisitNormalField(i *proto.NormalField) {
	protoField.Name = i.Field.Name
	protoField.Type = i.Field.Type
}

func GenCRUDData(ctx *GenLogicByProtoContext, p *parser.Proto, projectCtx *ctx.ProjectContext) []*ApiLogicData {
	var data []*ApiLogicData
	setLogic := strings.Builder{}

	for _, v := range p.Message {
		if strings.Contains(v.Name, ctx.ModelName) {
			if fmt.Sprintf("%sInfo", ctx.ModelName) == v.Name {
				for _, field := range v.Elements {
					field.Accept(MessageVisitor{})
					if protoField.Name == "id" || protoField.Name == "created_at" || protoField.Name == "updated_at" || protoField.Name == "deleted_at" {
						continue
					}
					setLogic.WriteString(fmt.Sprintf("\n\t\t\t%s: req.%s,", parser.CamelCase(protoField.Name),
						parser.CamelCase(protoField.Name)))
				}
				createLogic := bytes.NewBufferString("")
				createLogicTmpl, _ := template.New("createOrUpdate").Parse(createOrUpdateTpl)
				logx.Must(createLogicTmpl.Execute(createLogic, map[string]interface{}{
					"setLogic":           setLogic.String(),
					"modelName":          ctx.ModelName,
					"modelNameLowerCase": strings.ToLower(ctx.ModelName),
					"projectPackage":     projectCtx.Path,
					"rpcPackage":         ctx.GrpcPackage,
					"rpcName":            ctx.RpcName,
					"rpcNameLowerCase":   strings.ToLower(ctx.RpcName),
				}))

				data = append(data, &ApiLogicData{
					LogicName: fmt.Sprintf("CreateOrUpdate%sLogic", ctx.ModelName),
					LogicCode: createLogic.String(),
				})

				// delete logic
				deleteLogic := bytes.NewBufferString("")
				deleteLogicTmpl, _ := template.New("delete").Parse(deleteLogicTpl)
				logx.Must(deleteLogicTmpl.Execute(deleteLogic, map[string]interface{}{
					"setLogic":           setLogic.String(),
					"modelName":          ctx.ModelName,
					"modelNameLowerCase": strings.ToLower(ctx.ModelName),
					"projectPackage":     projectCtx.Path,
					"rpcPackage":         ctx.GrpcPackage,
					"rpcName":            ctx.RpcName,
					"rpcNameLowerCase":   strings.ToLower(ctx.RpcName),
				}))

				data = append(data, &ApiLogicData{
					LogicName: fmt.Sprintf("Delete%sLogic", ctx.ModelName),
					LogicCode: deleteLogic.String(),
				})

				// batch delete logic
				batchDeleteLogic := bytes.NewBufferString("")
				batchDeleteLogicTmpl, _ := template.New("batchDelete").Parse(batchDeleteLogicTpl)
				logx.Must(batchDeleteLogicTmpl.Execute(batchDeleteLogic, map[string]interface{}{
					"setLogic":           setLogic.String(),
					"modelName":          ctx.ModelName,
					"modelNameLowerCase": strings.ToLower(ctx.ModelName),
					"projectPackage":     projectCtx.Path,
					"rpcPackage":         ctx.GrpcPackage,
					"rpcName":            ctx.RpcName,
					"rpcNameLowerCase":   strings.ToLower(ctx.RpcName),
				}))

				data = append(data, &ApiLogicData{
					LogicName: fmt.Sprintf("BatchDelete%sLogic", ctx.ModelName),
					LogicCode: batchDeleteLogic.String(),
				})
			}

			if fmt.Sprintf("%sPageReq", ctx.ModelName) == v.Name {
				searchLogic := strings.Builder{}
				for _, field := range v.Elements {
					field.Accept(MessageVisitor{})
					if protoField.Name == "page" || protoField.Name == "page_size" {
						continue
					}
					searchLogic.WriteString(fmt.Sprintf("\n\t\t\t%s: req.%s,", parser.CamelCase(protoField.Name),
						parser.CamelCase(protoField.Name)))
				}

				getListLogic := bytes.NewBufferString("")
				getListLogicTmpl, _ := template.New("getList").Parse(getListLogicTpl)
				logx.Must(getListLogicTmpl.Execute(getListLogic, map[string]interface{}{
					"setLogic":           strings.Replace(setLogic.String(), "req.", "v.", -1),
					"modelName":          ctx.ModelName,
					"modelNameLowerCase": strings.ToLower(ctx.ModelName),
					"projectPackage":     projectCtx.Path,
					"rpcPackage":         ctx.GrpcPackage,
					"rpcName":            ctx.RpcName,
					"rpcNameLowerCase":   strings.ToLower(ctx.RpcName),
					"searchKeys":         searchLogic.String(),
				}))

				data = append(data, &ApiLogicData{
					LogicName: fmt.Sprintf("Get%sListLogic", ctx.ModelName),
					LogicCode: getListLogic.String(),
				})
			}

		}
	}

	return data
}

func GenApiData(schema *load.Schema, searchKeyNum int) (string, string, error) {
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
					typeName := v.Info.Type.String()
					if typeName == "float32" || typeName == "float64" {
						typeName = "float"
					}
					protoMessage.WriteString(fmt.Sprintf("  %s %s = %d;\n", typeName, v.Name, index))
				}
			} else {
				if strings.Contains(v.Name, "at") {
					protoMessage.WriteString(fmt.Sprintf("  int64  %s = %d;\n}\n\n", v.Name, index))
				} else {
					typeName := v.Info.Type.String()
					if typeName == "float32" || typeName == "float64" {
						typeName = "float"
					}
					protoMessage.WriteString(fmt.Sprintf("  %s %s = %d;\n}\n\n", typeName, v.Name, index))
				}
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
		if v.Info.Type.String() == "string" && !strings.Contains(strings.ToLower(v.Name), "uuid") && count <= searchKeyNum {
			if i < (len(schema.Fields)-1) && count < (searchKeyNum-1) {
				protoMessage.WriteString(fmt.Sprintf("  %s %s = %d;\n", v.Info.Type.String(), v.Name, index))
			}
			index++
			count++
		}

		if i == (len(schema.Fields) - 1) {
			protoMessage.WriteString("}\n")
		}
	}

	protoRpcFunction := bytes.NewBufferString("")
	protoTmpl, err := template.New("proto").Parse(apiTpl)
	err = protoTmpl.Execute(protoRpcFunction, map[string]interface{}{
		"modelName": schema.Name,
	})

	if err != nil {
		logx.Error(err)
		return "", "", err
	}

	return protoMessage.String(), protoRpcFunction.String(), nil
}
