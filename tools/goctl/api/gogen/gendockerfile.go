package gogen

import (
	_ "embed"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

//go:embed dockerfile.tpl
var dockerfileTemplate string

func genDockerfile(dir string, api *spec.ApiSpec, g *GenContext) error {
	service := api.Service

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          "",
		filename:        "Dockerfile",
		templateName:    "dockerfileTemplate",
		category:        category,
		templateFile:    dockerfileTemplateFile,
		builtinTemplate: dockerfileTemplate,
		data: map[string]string{
			"serviceName": service.Name,
			"port":        string(g.Port),
		},
	})
}
