package vben

import (
	"fmt"
	"path/filepath"

	"github.com/iancoleman/strcase"

	"github.com/zeromicro/go-zero/tools/goctl/util"
)

func genDrawer(g *GenContext) error {
	if err := util.With("drawerTpl").Parse(drawerTpl).SaveTo(map[string]any{
		"modelName":           g.ModelName,
		"modelNameLowerCamel": strcase.ToLowerCamel(g.ModelName),
		"folderName":          g.FolderName,
		"useUUID":             g.UseUUID,
	},
		filepath.Join(g.ViewDir, fmt.Sprintf("%sDrawer.vue", g.ModelName)), false); err != nil {
		return err
	}
	return nil
}
