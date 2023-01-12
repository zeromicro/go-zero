package vben

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/util"
)

func genApi(g *GenContext) error {
	if err := util.With("apiTpl").Parse(apiTpl).SaveTo(map[string]any{
		"modelName":          g.ModelName,
		"modelNameLowerCase": strings.ToLower(g.ModelName),
		"prefix":             g.Prefix,
		"useUUID":            g.UseUUID,
	},
		filepath.Join(g.ApiDir, fmt.Sprintf("%s.ts", strings.ToLower(g.ModelName))), false); err != nil {
		return err
	}
	return nil
}
