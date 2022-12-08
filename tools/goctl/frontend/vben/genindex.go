package vben

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/util"
)

func genIndex(g *GenContext) error {
	if err := util.With("indexTpl").Parse(indexTpl).SaveTo(map[string]interface{}{
		"modelName":          g.ModelName,
		"modelNameLowerCase": strings.ToLower(g.ModelName),
		"folderName":         g.FolderName,
		"addButtonTitle":     fmt.Sprintf("{{ t('%s.%s.add%s') }}", g.FolderName, strings.ToLower(g.ModelName), g.ModelName),
		"deleteButtonTitle":  "{{ t('common.delete') }}",
	},
		filepath.Join(g.ViewDir, "index.vue"), false); err != nil {
		return err
	}
	return nil
}
