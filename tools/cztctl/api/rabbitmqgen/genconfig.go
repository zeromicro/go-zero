package rabbitmqgen

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/lerity-yao/go-zero/tools/cztctl/api/gogen"
	"github.com/lerity-yao/go-zero/tools/cztctl/api/spec"
	"github.com/lerity-yao/go-zero/tools/cztctl/config"
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

func genConfig(dir string, cfg *config.Config, api *spec.ApiSpec) error {
	filename, err := format.FileNamingFormat(cfg.NamingFormat, configFile)
	if err != nil {
		return err
	}

	authImportStr := fmt.Sprintf("\"%s/rest\"", vars.ProjectOpenSourceURL)
	authImportStr = authImportStr + fmt.Sprintf("\n\"%s/go-mq/rabbitmq\"", vars.RabbitmqProjectOpenSourceURL)
	configNames := generateRabbitmqConfigNames(api)

	return gogen.GenFile(gogen.FileGenConfig{
		Dir:             dir,
		Subdir:          configDir,
		Filename:        filename + ".go",
		TemplateName:    "configTemplate",
		Category:        category,
		TemplateFile:    configTemplateFile,
		BuiltinTemplate: configTemplate,
		Data: map[string]string{
			"authImport":     authImportStr,
			"rabbitmqConfig": strings.Join(configNames, "\n"),
		},
	})
}
