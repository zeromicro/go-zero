package gogen

import (
	"fmt"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"github.com/tal-tech/go-zero/tools/goctl/config"
	"github.com/tal-tech/go-zero/tools/goctl/util/format"
	"github.com/tal-tech/go-zero/tools/goctl/vars"
)

const (
	configFile     = "config"
	configTemplate = `package config

import {{.authImport}}

type Config struct {
	rest.RestConf
	{{.auth}}
}
`

	jwtTemplate = ` struct {
		AccessSecret string
		AccessExpire int64
	}
`
)

func genConfig(dir string, cfg *config.Config, api *spec.ApiSpec) error {
	filename, err := format.FileNamingFormat(cfg.NamingFormat, configFile)
	if err != nil {
		return err
	}

	var authNames = getAuths(api)
	var auths []string
	for _, item := range authNames {
		auths = append(auths, fmt.Sprintf("%s %s", item, jwtTemplate))
	}
	var authImportStr = fmt.Sprintf("\"%s/rest\"", vars.ProjectOpenSourceUrl)

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          configDir,
		filename:        filename + ".go",
		templateName:    "configTemplate",
		category:        category,
		templateFile:    configTemplateFile,
		builtinTemplate: configTemplate,
		data: map[string]string{
			"authImport": authImportStr,
			"auth":       strings.Join(auths, "\n"),
		},
	})
}
