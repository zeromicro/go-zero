package vben

import (
	"fmt"
	"path/filepath"

	"github.com/iancoleman/strcase"

	"github.com/zeromicro/go-zero/tools/goctl/util"
)

func genIndex(g *GenContext) error {
	if err := util.With("indexTpl").Parse(indexTpl).SaveTo(map[string]any{
		"modelName":           g.ModelName,
		"modelNameLowerCamel": strcase.ToLowerCamel(g.ModelName),
		"folderName":          g.FolderName,
		"addButtonTitle":      fmt.Sprintf("{{ t('%s.%s.add%s') }}", g.FolderName, strcase.ToLowerCamel(g.ModelName), g.ModelName),
		"deleteButtonTitle":   "{{ t('common.delete') }}",
		"useUUID":             g.UseUUID,
	},
		filepath.Join(g.ViewDir, "index.vue"), false); err != nil {
		return err
	}
	return nil
}
