package vben

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/util"
)

func genData(g *GenContext) error {
	var basicData, searchFormData, formData strings.Builder
	var useBaseInfo bool
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
		"modelName":      g.ModelName,
		"basicData":      basicData.String(),
		"searchFormData": searchFormData.String(),
		"formData":       formData.String(),
		"useBaseInfo":    useBaseInfo,
		"useUUID":        g.UseUUID,
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
	default:
		return "component: 'Input',"
	}
}
