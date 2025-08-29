package rabbitmqgen

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/lerity-yao/go-zero/tools/cztctl/api/gogen"
	"github.com/lerity-yao/go-zero/tools/cztctl/api/spec"
	"github.com/lerity-yao/go-zero/tools/cztctl/config"
	"github.com/lerity-yao/go-zero/tools/cztctl/util/format"
	"github.com/lerity-yao/go-zero/tools/cztctl/util/pathx"
	"github.com/lerity-yao/go-zero/tools/cztctl/vars"
)

//go:embed main.tpl
var mainTemplate string

func genMain(dir, rootPkg, projectPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	filename, err := format.FileNamingFormat(cfg.NamingFormat, api.Service.Name)
	if err != nil {
		return err
	}

	configName := filename
	if strings.HasSuffix(filename, "-api") {
		filename = strings.ReplaceAll(filename, "-api", "")
	}

	return gogen.GenFile(gogen.FileGenConfig{
		Dir:             dir,
		Subdir:          "",
		Filename:        filename + ".go",
		TemplateName:    "mainTemplate",
		Category:        category,
		TemplateFile:    mainTemplateFile,
		BuiltinTemplate: mainTemplate,
		Data: map[string]string{
			"importPackages": genMainImports(rootPkg),
			"serviceName":    configName,
			"projectPkg":     projectPkg,
		},
	})
}

func genMainImports(parentPkg string) string {
	var imports []string
	imports = append(imports, fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, configDir)))
	imports = append(imports, fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, handlerDir)))
	imports = append(imports, fmt.Sprintf("\"%s\"\n", pathx.JoinPackages(parentPkg, contextDir)))
	imports = append(imports, fmt.Sprintf("\"%s/core/conf\"", vars.ProjectOpenSourceURL))
	return strings.Join(imports, "\n\t")
}
