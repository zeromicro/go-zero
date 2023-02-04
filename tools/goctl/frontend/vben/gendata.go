package vben

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/util"
)

func genData(g *GenContext) error {
	var basicData, searchFormData, formData strings.Builder
	var useBaseInfo bool
	var statusBasicColumnData, statusFormColumnData string
	// generate basic and search form data
	for _, v := range g.ApiSpec.Types {
		if v.Name() == fmt.Sprintf("%sInfo", g.ModelName) {
			specData, ok := v.(spec.DefineStruct)
			if !ok {
				return errors.New("cannot get the field")
			}

			for _, val := range specData.Members {
				if val.Name == "" {
					tmpType, _ := val.Type.(spec.DefineStruct)
					if tmpType.Name() == "BaseInfo" {
						useBaseInfo = true
					}
				} else if val.Name == "Status" {
					statusRenderData := bytes.NewBufferString("")
					protoTmpl, _ := template.New("proto").Parse(statusRenderTpl)
					protoTmpl.Execute(statusRenderData, map[string]any{
						"modelName": strings.TrimSuffix(specData.RawName, "Info"),
					})
					statusBasicColumnData = statusRenderData.String()

					statusFormColumnData = fmt.Sprintf("\n  {\n    field: '%s',\n    label: t('%s'),\n    component: 'RadioButtonGroup',\n"+
						"    defaultValue: 1,\n    componentProps: {\n      options: [\n        { label: t('common.on'), value: 1 },\n    "+
						"    { label: t('common.off'), value: 0 },\n      ],\n    },\n  },",
						strcase.ToLowerCamel(val.Name),
						fmt.Sprintf("%s.%s.%s", g.FolderName,
							strcase.ToLowerCamel(strings.TrimSuffix(specData.RawName, "Info")),
							strcase.ToLowerCamel(val.Name)),
					)
				} else {
					basicData.WriteString(fmt.Sprintf("\n  {\n    title: t('%s'),\n    dataIndex: '%s',\n    width: 100,\n  },",
						fmt.Sprintf("%s.%s.%s", g.FolderName,
							strcase.ToLowerCamel(strings.TrimSuffix(specData.RawName, "Info")),
							strcase.ToLowerCamel(val.Name)), strcase.ToLowerCamel(val.Name)))

					formData.WriteString(fmt.Sprintf("\n  {\n    field: '%s',\n    label: t('%s'),\n    %s\n  },",
						strcase.ToLowerCamel(val.Name),
						fmt.Sprintf("%s.%s.%s", g.FolderName,
							strcase.ToLowerCamel(strings.TrimSuffix(specData.RawName, "Info")),
							strcase.ToLowerCamel(val.Name)),
						getComponent(val.Type.Name()),
					))
				}
			}

			// put here in order to put status in the end
			if g.HasStatus {
				basicData.WriteString(statusBasicColumnData)
				formData.WriteString(statusFormColumnData)
			}
		}

		if v.Name() == fmt.Sprintf("%sListReq", g.ModelName) {
			specData, ok := v.(spec.DefineStruct)
			if !ok {
				return errors.New("cannot get field")
			}

			for _, val := range specData.Members {
				if val.Name != "" {
					searchFormData.WriteString(fmt.Sprintf("\n  {\n    field: '%s',\n    label: t('%s'),\n    component: 'Input',\n    colProps: { span: 8 },\n  },",
						strcase.ToLowerCamel(val.Name),
						fmt.Sprintf("%s.%s.%s", g.FolderName,
							strcase.ToLowerCamel(strings.TrimSuffix(specData.RawName, "ListReq")),
							strcase.ToLowerCamel(val.Name)),
					))
				}
			}
		}
	}

	if err := util.With("dataTpl").Parse(dataTpl).SaveTo(map[string]any{
		"modelName":          g.ModelName,
		"modelNameLowerCase": strings.ToLower(g.ModelName),
		"folderName":         g.FolderName,
		"basicData":          basicData.String(),
		"searchFormData":     searchFormData.String(),
		"formData":           formData.String(),
		"useBaseInfo":        useBaseInfo,
		"useUUID":            g.UseUUID,
		"hasStatus":          g.HasStatus,
	},
		filepath.Join(g.ViewDir, fmt.Sprintf("%s.data.ts", strings.ToLower(g.ModelName))), false); err != nil {
		return err
	}
	return nil
}

func getComponent(dataType string) string {
	switch dataType {
	case "string":
		return "component: 'Input',"
	case "int32", "int64", "uint32", "uint64":
		return "component: 'InputNumber',"
	case "bool":
		return "component: 'RadioButtonGroup',\n    defaultValue: false,\n    componentProps: {\n      options: [\n        { label: t('common.on'), value: false },\n        { label: t('common.off'), value: true },\n      ],\n    },"
	default:
		return "component: 'Input',"
	}
}
