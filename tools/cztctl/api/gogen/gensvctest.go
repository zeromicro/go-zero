package gogen

import (
	_ "embed"

	"github.com/lerity-yao/go-zero/tools/cztctl/api/spec"
	"github.com/lerity-yao/go-zero/tools/cztctl/config"
	"github.com/lerity-yao/go-zero/tools/cztctl/internal/version"
	"github.com/lerity-yao/go-zero/tools/cztctl/util/format"
)

//go:embed svc_test.tpl
var svcTestTemplate string

func genServiceContextTest(dir, rootPkg, projectPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	filename, err := format.FileNamingFormat(cfg.NamingFormat, contextFilename)
	if err != nil {
		return err
	}

	return GenFile(FileGenConfig{
		Dir:             dir,
		Subdir:          contextDir,
		Filename:        filename + "_test.go",
		TemplateName:    "svcTestTemplate",
		Category:        category,
		TemplateFile:    svcTestTemplateFile,
		BuiltinTemplate: svcTestTemplate,
		Data: map[string]any{
			"projectPkg": projectPkg,
			"version":    version.BuildVersion,
		},
	})
}
