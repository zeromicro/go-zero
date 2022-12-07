package vben

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/util"
)

func genApi(g *GenContext) error {
	if err := util.With("apiTpl").Parse(apiTpl).SaveTo(map[string]interface{}{
		"modelName":          g.ModelName,
		"modelNameLowerCase": strings.ToLower(g.ModelName),
	},
		filepath.Join(g.ApiDir, fmt.Sprintf("%s.ts", strings.ToLower(g.ModelName))), false); err != nil {
		return err
	}
	return nil
}
