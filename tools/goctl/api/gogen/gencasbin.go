package gogen

import (
	_ "embed"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/config"
)

//go:embed authortymiddleware.tpl
var authorityMiddlewareTemplate string

func genCasbin(dir string, cfg *config.Config, api *spec.ApiSpec) error {
	err := genFile(fileGenConfig{
		dir:             dir,
		subdir:          middlewareDir,
		filename:        "authority_middleware.go",
		templateName:    "authorityTemplate",
		category:        category,
		templateFile:    authorityTemplateFile,
		builtinTemplate: authorityMiddlewareTemplate,
		data:            map[string]string{},
	})
	if err != nil {
		return err
	}

	return nil
}
