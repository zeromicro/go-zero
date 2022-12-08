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
						strcase.ToCamel(val.Name), strcase.ToLowerCamel(val.Name)))
					formData.WriteString(fmt.Sprintf("\n  {\n    field: '%s',\n    label: t('%s'),\n    component: 'Input',\n  },",
						strcase.ToCamel(val.Name), strcase.ToLowerCamel(val.Name)))
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
						strcase.ToCamel(val.Name), strcase.ToLowerCamel(val.Name)))
				}
			}
		}
	}

	if err := util.With("dataTpl").Parse(dataTpl).SaveTo(map[string]interface{}{
		"modelName":      g.ModelName,
		"basicData":      basicData.String(),
		"searchFormData": searchFormData.String(),
		"formData":       formData.String(),
		"useBaseInfo":    useBaseInfo,
	},
		filepath.Join(g.ViewDir, fmt.Sprintf("%s.data.ts", strings.ToLower(g.ModelName))), false); err != nil {
		return err
	}
	return nil
}
