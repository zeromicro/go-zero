package gogen

import (
	_ "embed"
	"strings"

	"github.com/lerity-yao/go-zero/tools/cztctl/api/spec"
	"github.com/lerity-yao/go-zero/tools/cztctl/config"
	"github.com/lerity-yao/go-zero/tools/cztctl/internal/version"
	"github.com/lerity-yao/go-zero/tools/cztctl/util/format"
)

//go:embed middleware.tpl
var middlewareImplementCode string

func genMiddleware(dir string, cfg *config.Config, api *spec.ApiSpec) error {
	middlewares := getMiddleware(api)
	for _, item := range middlewares {
		middlewareFilename := strings.TrimSuffix(strings.ToLower(item), "middleware") + "_middleware"
		filename, err := format.FileNamingFormat(cfg.NamingFormat, middlewareFilename)
		if err != nil {
			return err
		}

		name := strings.TrimSuffix(item, "Middleware") + "Middleware"
		err = GenFile(FileGenConfig{
			Dir:             dir,
			Subdir:          middlewareDir,
			Filename:        filename + ".go",
			TemplateName:    "contextTemplate",
			Category:        category,
			TemplateFile:    middlewareImplementCodeFile,
			BuiltinTemplate: middlewareImplementCode,
			Data: map[string]string{
				"name":    strings.Title(name),
				"version": version.BuildVersion,
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}
