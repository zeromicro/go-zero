package rabbitmqgen

import (
	_ "embed"
	"fmt"
	"path"
	"strings"

	"github.com/lerity-yao/go-zero/tools/cztctl/api/gogen"
	"github.com/lerity-yao/go-zero/tools/cztctl/api/spec"
	"github.com/lerity-yao/go-zero/tools/cztctl/config"
	"github.com/lerity-yao/go-zero/tools/cztctl/util/format"
	"github.com/lerity-yao/go-zero/tools/cztctl/util/pathx"
	"github.com/lerity-yao/go-zero/tools/cztctl/vars"
)

const defaultLogicPackage = "logic"

var (
	//go:embed handler.tpl
	handlerTemplate string
	////go:embed sse_handler.tpl
	//sseHandlerTemplate string
)

func genHandler(dir, rootPkg, projectPkg string, cfg *config.Config, group spec.Group, route spec.Route) error {
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

	var builtinTemplate = handlerTemplate

	return gogen.GenFile(gogen.FileGenConfig{
		Dir:             dir,
		Subdir:          getHandlerFolderPath(group, route),
		Filename:        filename + ".go",
		TemplateName:    "handlerTemplate",
		Category:        category,
		TemplateFile:    handlerTemplateFile,
		BuiltinTemplate: builtinTemplate,
		Data: map[string]any{
			"PkgName":          pkgName,
			"ImportPackages":   genHandlerImports(group, route, rootPkg),
			"HandlerName":      handler,
			"RabbitmqConfName": fmt.Sprintf("%s%s", strings.TrimSuffix(handler, "Handler"), "RabbitmqConf"),
			"LogicName":        logicName,
			"LogicType":        getLogicName(route),
			"HasDoc":           len(route.JoinedDoc()) > 0,
			"Doc":              GetDoc(route.JoinedDoc()),
		},
	})
}

func genHandlers(dir, rootPkg, projectPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	for _, group := range api.Service.Groups {
		for _, route := range group.Routes {
			if err := genHandler(dir, rootPkg, projectPkg, cfg, group, route); err != nil {
				return err
			}
		}
	}

	return nil
}

func genHandlerImports(group spec.Group, route spec.Route, parentPkg string) string {
	imports := []string{
		fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, gogen.GetLogicFolderPath(group, route))),
		fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, contextDir)),
		fmt.Sprintf("\"%s\"", pathx.JoinPackages(vars.RabbitmqProjectOpenSourceURL, "go-mq/rabbitmq")),
		fmt.Sprintf("\"%s\"", pathx.JoinPackages(vars.ProjectOpenSourceURL, "core/service")),
	}
	sse := group.GetAnnotation("sse")
	if len(route.RequestTypeName()) > 0 || sse == "true" {
		imports = append(imports, fmt.Sprintf("\"%s\"\n", pathx.JoinPackages(parentPkg, typesDir)))
	}

	return strings.Join(imports, "\n\t")
}

func getHandlerBaseName(route spec.Route) (string, error) {
	handler := route.Handler
	handler = strings.TrimSpace(handler)
	handler = strings.TrimSuffix(handler, "handler")
	handler = strings.TrimSuffix(handler, "Handler")

	return handler, nil
}

func getHandlerFolderPath(group spec.Group, route spec.Route) string {
	folder := route.GetAnnotation(groupProperty)
	if len(folder) == 0 {
		folder = group.GetAnnotation(groupProperty)
		if len(folder) == 0 {
			return handlerDir
		}
	}

	folder = strings.TrimPrefix(folder, "/")
	folder = strings.TrimSuffix(folder, "/")

	return path.Join(handlerDir, folder)
}

func getHandlerName(route spec.Route) string {
	handler, err := getHandlerBaseName(route)
	if err != nil {
		panic(err)
	}

	return handler + "Handler"
}

func getLogicName(route spec.Route) string {
	handler, err := getHandlerBaseName(route)
	if err != nil {
		panic(err)
	}

	return strings.TrimSuffix(handler, "Handler") + "Logic"
}
