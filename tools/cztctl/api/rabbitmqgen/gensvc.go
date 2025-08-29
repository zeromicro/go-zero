package rabbitmqgen

import (
	_ "embed"

	"github.com/lerity-yao/go-zero/tools/cztctl/config"
	"github.com/lerity-yao/go-zero/tools/cztctl/util/format"
	"github.com/lerity-yao/go-zero/tools/cztctl/util/pathx"

	"github.com/lerity-yao/go-zero/tools/cztctl/api/gogen"
)

const contextFilename = "service_context"

//go:embed svc.tpl
var contextTemplate string

func genServiceContext(dir, rootPkg string, cfg *config.Config) error {
	filename, err := format.FileNamingFormat(cfg.NamingFormat, contextFilename)
	if err != nil {
		return err
	}

	importPackages := "\"" + pathx.JoinPackages(rootPkg, configDir) + "\""

	return gogen.GenFile(gogen.FileGenConfig{
		Dir:             dir,
		Subdir:          contextDir,
		Filename:        filename + ".go",
		TemplateName:    "contextTemplate",
		Category:        category,
		TemplateFile:    contextTemplateFile,
		BuiltinTemplate: contextTemplate,
		Data: map[string]string{
			"importPackages": importPackages,
			"config":         "config.Config",
		},
	})
}
