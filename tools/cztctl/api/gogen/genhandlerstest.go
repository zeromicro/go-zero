package gogen

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/lerity-yao/go-zero/tools/cztctl/api/spec"
	"github.com/lerity-yao/go-zero/tools/cztctl/config"
	"github.com/lerity-yao/go-zero/tools/cztctl/internal/version"
	"github.com/lerity-yao/go-zero/tools/cztctl/util"
	"github.com/lerity-yao/go-zero/tools/cztctl/util/format"
	"github.com/lerity-yao/go-zero/tools/cztctl/util/pathx"
)

//go:embed handler_test.tpl
var handlerTestTemplate string

func genHandlerTest(dir, rootPkg, projectPkg string, cfg *config.Config, group spec.Group, route spec.Route) error {
	handler := getHandlerName(route)
	handlerPath := getHandlerFolderPath(group, route)
	pkgName := handlerPath[strings.LastIndex(handlerPath, "/")+1:]
	logicName := defaultLogicPackage
	if handlerPath != handlerDir {
		handler = strings.Title(handler)
		logicName = pkgName
	}
	filename, err := format.FileNamingFormat(cfg.NamingFormat, handler)
	if err != nil {
		return err
	}

	return GenFile(FileGenConfig{
		Dir:             dir,
		Subdir:          getHandlerFolderPath(group, route),
		Filename:        filename + "_test.go",
		TemplateName:    "handlerTestTemplate",
		Category:        category,
		TemplateFile:    handlerTestTemplateFile,
		BuiltinTemplate: handlerTestTemplate,
		Data: map[string]any{
			"PkgName":        pkgName,
			"ImportPackages": genHandlerTestImports(group, route, rootPkg),
			"HandlerName":    handler,
			"RequestType":    util.Title(route.RequestTypeName()),
			"ResponseType":   util.Title(route.ResponseTypeName()),
			"LogicName":      logicName,
			"LogicType":      strings.Title(getLogicName(route)),
			"Call":           strings.Title(strings.TrimSuffix(handler, "Handler")),
			"HasResp":        len(route.ResponseTypeName()) > 0,
			"HasRequest":     len(route.RequestTypeName()) > 0,
			"HasDoc":         len(route.JoinedDoc()) > 0,
			"Doc":            GetDoc(route.JoinedDoc()),
			"projectPkg":     projectPkg,
			"version":        version.BuildVersion,
		},
	})
}

func genHandlersTest(dir, rootPkg, projectPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	for _, group := range api.Service.Groups {
		for _, route := range group.Routes {
			if err := genHandlerTest(dir, rootPkg, projectPkg, cfg, group, route); err != nil {
				return err
			}
		}
	}

	return nil
}

func genHandlerTestImports(group spec.Group, route spec.Route, parentPkg string) string {
	imports := []string{
		//fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, getLogicFolderPath(group, route))),
		fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, contextDir)),
		fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, configDir)),
	}
	if len(route.RequestTypeName()) > 0 {
		imports = append(imports, fmt.Sprintf("\"%s\"\n", pathx.JoinPackages(parentPkg, typesDir)))
	}

	return strings.Join(imports, "\n\t")
}
