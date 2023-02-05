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

func genModel(g *GenContext) error {
	var infoData strings.Builder
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
						infoData.WriteString("  id: number;\n  createdAt?: number;\n")
					} else if tmpType.Name() == "BaseUUIDInfo" {
						infoData.WriteString("  id: string;\n  createdAt?: number;\n")
						g.UseUUID = true
					}
				} else if val.Name == "Status" {
					g.HasStatus = true
				} else {
					infoData.WriteString(fmt.Sprintf("  %s: %s;\n", strcase.ToLowerCamel(val.Name),
						ConvertGoTypeToTsType(val.Type.Name())))
				}
			}

		}
	}
	if err := util.With("modelTpl").Parse(modelTpl).SaveTo(map[string]any{
		"modelName": g.ModelName,
		"infoData":  infoData.String(),
	},
		filepath.Join(g.ModelDir, fmt.Sprintf("%sModel.ts", strcase.ToLowerCamel(g.ModelName))), false); err != nil {
		return err
	}
	return nil
}
