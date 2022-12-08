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

func genDrawer(g *GenContext) error {
	var infoData strings.Builder
	for _, v := range g.ApiSpec.Types {
		if v.Name() == fmt.Sprintf("%sInfo", g.ModelName) {
			specData, ok := v.(spec.DefineStruct)
			if !ok {
				return errors.New("cannot get the field")
			}

			for _, val := range specData.Members {
				if val.Name != "" {
					infoData.WriteString(fmt.Sprintf("\n          %s: values['%s'],", strcase.ToLowerCamel(val.Name),
						strcase.ToLowerCamel(val.Name)))
				}
			}

		}
	}

	if err := util.With("drawerTpl").Parse(drawerTpl).SaveTo(map[string]interface{}{
		"modelName":          g.ModelName,
		"modelNameLowerCase": strings.ToLower(g.ModelName),
		"folderName":         g.FolderName,
		"infoData":           infoData.String(),
	},
		filepath.Join(g.ViewDir, fmt.Sprintf("%sDrawer.vue", g.ModelName)), false); err != nil {
		return err
	}
	return nil
}
