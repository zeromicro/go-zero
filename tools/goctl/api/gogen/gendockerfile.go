package gogen

import (
	_ "embed"
	"fmt"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
)

//go:embed dockerfile.tpl
var dockerfileTemplate string

func genDockerfile(dir string, cfg *config.Config, api *spec.ApiSpec, g *GenContext) error {

	service, err := format.FileNamingFormat(cfg.NamingFormat, api.Service.Name)
	if err != nil {
		return err
	}

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          "",
		filename:        "Dockerfile",
		templateName:    "dockerfileTemplate",
		category:        category,
		templateFile:    dockerfileTemplateFile,
		builtinTemplate: dockerfileTemplate,
		data: map[string]string{
			"serviceName": service,
			"port":        fmt.Sprint(g.Port),
			"imageTag":    "golang:1.20.2-alpine3.17",
		},
	})
}
