// Copyright 2023 The Ryan SU Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

	"github.com/logrusorgru/aurora"

	"github.com/zeromicro/go-zero/tools/goctl/api/gogen/proto"
	"github.com/zeromicro/go-zero/tools/goctl/util/console"

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
)

const regularPerm = 0o666

type ApiLogicData struct {
	LogicName string
	LogicCode string
}

type GenEntLogicContext struct {
	Schema       string
	Output       string
	ServiceName  string
	Style        string
	ModelName    string
	SearchKeyNum int
	GroupName    string
	UseUUID      bool
	JSONStyle    string
	Overwrite    bool
}

func (g GenEntLogicContext) Validate() error {
	if g.Schema == "" {
		return errors.New("the schema dir cannot be empty ")
	} else if !strings.HasSuffix(g.Schema, "schema") {
		return errors.New("please input correct schema directory e.g. ./ent/schema ")
	} else if g.ServiceName == "" {
		return errors.New("please set the API service name via --api_service_name")
	} else if g.ModelName == "" {
		return errors.New("please set the model name via --model ")
	}
	return nil
}

// GenEntLogic generates the ent CRUD logic files of the api service.
func GenEntLogic(g *GenEntLogicContext) error {
	return genEntLogic(g)
}

func genEntLogic(g *GenEntLogicContext) error {
	console.NewColorConsole(true).Info(aurora.Green("Generating...").String())

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
			apiLogicData := GenCRUDData(g, projectCtx, s)

			for _, v := range apiLogicData {
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

				if pathx.FileExists(filename) && !g.Overwrite {
					continue
				}

				err = os.WriteFile(filename, []byte(v.LogicCode), regularPerm)
				if err != nil {
					return err
				}
			}

			// generate api file
			apiData, err := GenApiData(s, g)
			if err != nil {
				return err
			}

			apiFilePath := filepath.Join(workDir, "desc", fmt.Sprintf("%s.api", strcase.ToSnake(g.ModelName)))

			if pathx.FileExists(apiFilePath) && !g.Overwrite {
				return nil
			}

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

			if !strings.Contains(allApiString, fmt.Sprintf("%s.api", strcase.ToSnake(g.ModelName))) {
				allApiString += fmt.Sprintf("\nimport \"%s\"", fmt.Sprintf("%s.api", strcase.ToSnake(g.ModelName)))
			}

			err = os.WriteFile(allApiFile, []byte(allApiString), regularPerm)
			if err != nil {
				return err
			}
		}
	}

	console.NewColorConsole().Success(aurora.Green("Generate Ent Logic files for API successfully").String())

	return nil
}

func GenCRUDData(g *GenEntLogicContext, projectCtx *ctx.ProjectContext, schema *load.Schema) []*ApiLogicData {
	var data []*ApiLogicData
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
		} else {
			if entx.IsTimeProperty(v.Info.Type.String()) {
				hasTime = true
				setLogic.WriteString(fmt.Sprintf("\t\t\tSet%s(time.Unix(req.%s, 0)).\n", parser.CamelCase(v.Name),
					parser.CamelCase(v.Name)))
			} else if entx.IsUpperProperty(v.Name) {
				if entx.IsUUIDType(v.Info.Type.String()) {
					setLogic.WriteString(fmt.Sprintf("\t\t\tSet%s(uuidx.ParseUUIDString(req.%s)).\n", entx.ConvertSpecificNounToUpper(v.Name),
						parser.CamelCase(v.Name)))
					hasUUID = true
				} else {
					setLogic.WriteString(fmt.Sprintf("\t\t\tSet%s(req.%s).\n", entx.ConvertSpecificNounToUpper(v.Name),
						parser.CamelCase(v.Name)))
				}
			} else {
				setLogic.WriteString(fmt.Sprintf("\t\t\tSet%s(req.%s).\n", parser.CamelCase(v.Name),
					parser.CamelCase(v.Name)))
			}
		}
	}
	setLogic.WriteString("\t\t\tExec(l.ctx)")

	createLogic := bytes.NewBufferString("")
	createLogicTmpl, _ := template.New("create").Parse(createTpl)
	_ = createLogicTmpl.Execute(createLogic, map[string]any{
		"hasTime":     hasTime,
		"hasUUID":     hasUUID,
		"setLogic":    strings.ReplaceAll(setLogic.String(), "Exec", "Save"),
		"modelName":   schema.Name,
		"projectPath": projectCtx.Path,
		"packageName": packageName,
		"useUUID":     g.UseUUID, // UUID primary key
	})

	data = append(data, &ApiLogicData{
		LogicName: fmt.Sprintf("Create%sLogic", schema.Name),
		LogicCode: createLogic.String(),
	})

	updateLogic := bytes.NewBufferString("")
	updateLogicTmpl, _ := template.New("update").Parse(updateTpl)
	_ = updateLogicTmpl.Execute(updateLogic, map[string]any{
		"hasTime":     hasTime,
		"hasUUID":     hasUUID,
		"setLogic":    strings.Replace(setLogic.String(), "Set", "SetNotEmpty", -1),
		"modelName":   schema.Name,
		"projectPath": projectCtx.Path,
		"packageName": packageName,
		"useUUID":     g.UseUUID, // UUID primary key
	})

	data = append(data, &ApiLogicData{
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
			predicateData.WriteString(fmt.Sprintf("\tif req.%s != \"\" {\n\t\tpredicates = append(predicates, %s.%sContains(req.%s))\n\t}\n",
				camelName, strings.ToLower(schema.Name), entx.ConvertSpecificNounToUpper(v.Name), camelName))
			count++
		}
	}
	predicateData.WriteString(fmt.Sprintf("\tdata, err := l.svcCtx.DB.%s.Query().Where(predicates...).Page(l.ctx, req.Page, req.PageSize)",
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
			} else if entx.IsTimeProperty(v.Info.Type.String()) {
				listData.WriteString(fmt.Sprintf("\t\t\t%s:\tv.%s.UnixMilli(),%s", nameCamelCase,
					entx.ConvertSpecificNounToUpper(nameCamelCase), endString))
			} else {
				if entx.IsUpperProperty(v.Name) {
					listData.WriteString(fmt.Sprintf("\t\t\t%s:\tv.%s,%s", nameCamelCase,
						entx.ConvertSpecificNounToUpper(v.Name), endString))
				} else {
					listData.WriteString(fmt.Sprintf("\t\t\t%s:\tv.%s,%s", nameCamelCase,
						nameCamelCase, endString))
				}
			}
		}
	}

	getListLogic := bytes.NewBufferString("")
	getListLogicTmpl, _ := template.New("getList").Parse(getListLogicTpl)
	_ = getListLogicTmpl.Execute(getListLogic, map[string]any{
		"predicateData":      predicateData.String(),
		"modelName":          schema.Name,
		"listData":           listData.String(),
		"projectPath":        projectCtx.Path,
		"modelNameLowerCase": strings.ToLower(schema.Name),
		"packageName":        packageName,
		"useUUID":            g.UseUUID,
	})

	data = append(data, &ApiLogicData{
		LogicName: fmt.Sprintf("Get%sListLogic", schema.Name),
		LogicCode: getListLogic.String(),
	})

	getByIdLogic := bytes.NewBufferString("")
	getByIdLogicTmpl, _ := template.New("getById").Parse(getByIdLogicTpl)
	_ = getByIdLogicTmpl.Execute(getByIdLogic, map[string]any{
		"modelName":          schema.Name,
		"listData":           strings.Replace(listData.String(), "v.", "data.", -1),
		"projectPath":        projectCtx.Path,
		"modelNameLowerCase": strings.ToLower(schema.Name),
		"packageName":        packageName,
		"useUUID":            g.UseUUID,
	})

	data = append(data, &ApiLogicData{
		LogicName: fmt.Sprintf("Get%sByIdLogic", schema.Name),
		LogicCode: getByIdLogic.String(),
	})

	deleteLogic := bytes.NewBufferString("")
	deleteLogicTmpl, _ := template.New("delete").Parse(deleteLogicTpl)
	_ = deleteLogicTmpl.Execute(deleteLogic, map[string]any{
		"modelName":          schema.Name,
		"modelNameLowerCase": strings.ToLower(schema.Name),
		"projectPath":        projectCtx.Path,
		"packageName":        packageName,
		"useUUID":            g.UseUUID,
	})

	data = append(data, &ApiLogicData{
		LogicName: fmt.Sprintf("Delete%sLogic", schema.Name),
		LogicCode: deleteLogic.String(),
	})

	return data
}

func GenApiData(schema *load.Schema, ctx *GenEntLogicContext) (string, error) {
	infoData := strings.Builder{}
	listData := strings.Builder{}
	searchKeyNum := ctx.SearchKeyNum
	var data string

	for _, v := range schema.Fields {
		if entx.IsBaseProperty(v.Name) {
			continue
		}

		var structData string

		jsonTag, err := format.FileNamingFormat(ctx.JSONStyle, v.Name)
		if err != nil {
			return "", err
		}

		structData = fmt.Sprintf("\n\n        // %s\n        %s  %s `json:\"%s,optional\"`",
			parser.CamelCase(v.Name),
			parser.CamelCase(v.Name),
			entx.ConvertEntTypeToGotypeInSingleApi(v.Info.Type.String()),
			jsonTag)

		infoData.WriteString(structData)

		if v.Info.Type.String() == "string" && searchKeyNum > 0 {
			listData.WriteString(structData)
			searchKeyNum--
		}
	}

	apiTemplateData := bytes.NewBufferString("")
	apiTmpl, _ := template.New("apiTpl").Parse(proto.ApiTpl)
	logx.Must(apiTmpl.Execute(apiTemplateData, map[string]any{
		"infoData":           infoData.String(),
		"modelName":          ctx.ModelName,
		"modelNameSpace":     strings.Replace(strcase.ToSnake(ctx.ModelName), "_", " ", -1),
		"modelNameLowerCase": strings.ToLower(ctx.ModelName),
		"modelNameSnake":     strcase.ToSnake(ctx.ModelName),
		"listData":           listData.String(),
		"apiServiceName":     strcase.ToCamel(ctx.ServiceName),
		"useUUID":            ctx.UseUUID,
	}))
	data = apiTemplateData.String()

	return data, nil
}
