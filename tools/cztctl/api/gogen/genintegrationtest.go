package gogen

import (
	_ "embed"

	"github.com/zeromicro/go-zero/tools/cztctl/api/spec"
	"github.com/zeromicro/go-zero/tools/cztctl/config"
	"github.com/zeromicro/go-zero/tools/cztctl/internal/version"
	"github.com/zeromicro/go-zero/tools/cztctl/util/format"
)

//go:embed integration_test.tpl
var integrationTestTemplate string

func genIntegrationTest(dir, rootPkg, projectPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	serviceName := api.Service.Name
	if len(serviceName) == 0 {
		serviceName = "server"
	}

	filename, err := format.FileNamingFormat(cfg.NamingFormat, serviceName)
	if err != nil {
		return err
	}

	return GenFile(FileGenConfig{
		Dir:             dir,
		Subdir:          "",
		Filename:        filename + "_test.go",
		TemplateName:    "integrationTestTemplate",
		Category:        category,
		TemplateFile:    integrationTestTemplateFile,
		BuiltinTemplate: integrationTestTemplate,
		Data: map[string]any{
			"projectPkg":  projectPkg,
			"serviceName": serviceName,
			"version":     version.BuildVersion,
			"hasRoutes":   len(api.Service.Routes()) > 0,
			"routes":      api.Service.Routes(),
		},
	})
}
