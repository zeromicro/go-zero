package vben

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/zeromicro/go-zero/tools/goctl/util"
)

func genApi(g *GenContext) error {
	if err := util.With("apiTpl").Parse(apiTpl).SaveTo(map[string]any{
		"modelName":           g.ModelName,
		"modelNameLowerCase":  strings.Replace(strcase.ToSnake(g.ModelName), "_", " ", -1),
		"modelNameLowerCamel": strcase.ToLowerCamel(g.ModelName),
		"modelNameSnake":      strcase.ToSnake(g.ModelName),
		"prefix":              g.Prefix,
		"useUUID":             g.UseUUID,
		"hasStatus":           g.HasStatus,
	},
		filepath.Join(g.ApiDir, fmt.Sprintf("%s.ts", strcase.ToLowerCamel(g.ModelName))), false); err != nil {
		return err
	}
	return nil
}
