package proto

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/emicklei/proto"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/iancoleman/strcase"

	"github.com/zeromicro/go-zero/tools/goctl/rpc/execx"
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
	apiLogicData := GenCRUDData(p, &protoData, projectCtx)

	for _, v := range apiLogicData {
		logicFilename, err := format.FileNamingFormat(p.Style, v.LogicName)
		if err != nil {
			return err
		}

		filename := filepath.Join(logicDir, strings.ToLower(p.ModelName), logicFilename+".go")
		if err = pathx.MkdirIfNotExist(filepath.Join(logicDir, strings.ToLower(p.ModelName))); err != nil {
			return err
		}

		if pathx.FileExists(filename) {
			continue
		}

		err = os.WriteFile(filename, []byte(v.LogicCode), regularPerm)
		if err != nil {
			return err
		}
	}

	// generate api file
	apiData, err := GenApiData(p, &protoData, projectCtx)
	if err != nil {
		return err
	}

	apiFilePath := filepath.Join(workDir, "desc", fmt.Sprintf("%s.api", strings.ToLower(p.ModelName)))

	err = os.WriteFile(apiFilePath, []byte(apiData), regularPerm)
	if err != nil {
		return err
	}

	allApiFile := filepath.Join(workDir, "desc", "all.api")
	allApiData, err := os.ReadFile(allApiFile)
	if err != nil {
		return err
	}
	allApiString := string(allApiData)

	if !strings.Contains(allApiString, fmt.Sprintf("%s.api", strings.ToLower(p.ModelName))) {
		allApiString += fmt.Sprintf("\nimport \"%s\"", fmt.Sprintf("%s.api", strings.ToLower(p.ModelName)))
	}

	err = os.WriteFile(allApiFile, []byte(allApiString), regularPerm)
	if err != nil {
		return err
	}

	_, err = execx.Run("make gen-api", workDir)
	if err != nil {
		return err
	}

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

func GenApiData(ctx *GenLogicByProtoContext, p *parser.Proto, projectCtx *ctx.ProjectContext) (string, error) {
	infoData := strings.Builder{}
	listData := strings.Builder{}
	count := 0
	var data string

	for _, v := range p.Message {
		if strings.Contains(v.Name, ctx.ModelName) {
			if fmt.Sprintf("%sInfo", ctx.ModelName) == v.Name {
				for _, field := range v.Elements {
					field.Accept(MessageVisitor{})
					if protoField.Name == "id" || protoField.Name == "created_at" || protoField.Name == "updated_at" || protoField.Name == "deleted_at" {
						continue
					}
					var structData string

					structData += fmt.Sprintf("\n\n\t\t// %s\n\t\t%s  %s `json:\"%s\"`",
						parser.CamelCase(protoField.Name),
						parser.CamelCase(protoField.Name),
						protoField.Type,
						strcase.ToLowerCamel(protoField.Name))

					infoData.WriteString(structData)

					if count < ctx.SearchKeyNum && protoField.Type == "string" {
						listData.WriteString(structData)
					}
				}

				apiTemplateData := bytes.NewBufferString("")
				apiTmpl, _ := template.New("apiTpl").Parse(apiTpl)
				logx.Must(apiTmpl.Execute(apiTemplateData, map[string]interface{}{
					"infoData":           infoData.String(),
					"modelName":          ctx.ModelName,
					"modelNameLowerCase": strings.ToLower(ctx.ModelName),
					"listData":           listData.String(),
					"serviceName":        ctx.ServiceName,
				}))
				data = apiTemplateData.String()
			}
		}
	}
	return data, nil
}
