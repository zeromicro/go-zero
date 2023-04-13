package gogen

import (
	_ "embed"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
)

//go:embed makefile.tpl
var makefileTemplate string

func genMakefile(dir string, cfg *config.Config, api *spec.ApiSpec, g *GenContext) error {
	service := api.Service

	serviceNameStyle, err := format.FileNamingFormat(cfg.NamingFormat, service.Name)
	if err != nil {
		return err
	}

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          "",
		filename:        "Makefile",
		templateName:    "makefileTemplate",
		category:        category,
		templateFile:    makefileTemplateFile,
		builtinTemplate: makefileTemplate,
		data: map[string]any{
			"serviceName":      strcase.ToCamel(service.Name),
			"useEnt":           g.UseEnt,
			"serviceNameStyle": serviceNameStyle,
			"serviceNameLower": strings.ToLower(service.Name),
			"serviceNameSnake": strcase.ToSnake(service.Name),
			"serviceNameDash":  strings.ReplaceAll(strcase.ToSnake(service.Name), "_", "-"),
		},
	})
}
