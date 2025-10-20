package gogen

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/lerity-yao/go-zero/tools/cztctl/api/spec"
	"github.com/lerity-yao/go-zero/tools/cztctl/config"
	"github.com/zeromicro/go-zero/tools/cztctl/internal/version"
	"github.com/lerity-yao/go-zero/tools/cztctl/util/format"
	"github.com/lerity-yao/go-zero/tools/cztctl/vars"
)

const (
	configFile = "config"

	jwtTemplate = ` struct {
		AccessSecret string
		AccessExpire int64
	}
`
	jwtTransTemplate = ` struct {
		Secret     string
		PrevSecret string
	}
`
)

//go:embed config.tpl
var configTemplate string

func genConfig(dir, projectPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	filename, err := format.FileNamingFormat(cfg.NamingFormat, configFile)
	if err != nil {
		return err
	}

	authNames := getAuths(api)
	auths := make([]string, 0, len(authNames))
	for _, item := range authNames {
		auths = append(auths, fmt.Sprintf("%s %s", item, jwtTemplate))
	}

	jwtTransNames := getJwtTrans(api)
	jwtTransList := make([]string, 0, len(jwtTransNames))
	for _, item := range jwtTransNames {
		jwtTransList = append(jwtTransList, fmt.Sprintf("%s %s", item, jwtTransTemplate))
	}
	authImportStr := fmt.Sprintf("\"%s/rest\"", vars.ProjectOpenSourceURL)

	return GenFile(FileGenConfig{
		Dir:             dir,
		Subdir:          configDir,
		Filename:        filename + ".go",
		TemplateName:    "configTemplate",
		Category:        category,
		TemplateFile:    configTemplateFile,
		BuiltinTemplate: configTemplate,
		Data: map[string]string{
			"authImport": authImportStr,
			"auth":       strings.Join(auths, "\n"),
			"jwtTrans":   strings.Join(jwtTransList, "\n"),
			"projectPkg": projectPkg,
			"version":    version.BuildVersion,
		},
	})
}
