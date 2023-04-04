package gogen

import (
	_ "embed"
	"fmt"
	"path"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

const defaultLogicPackage = "logic"

//go:embed handler.tpl
var handlerTemplate string

type handlerInfo struct {
	PkgName            string
	ImportPackages     string
	ImportHttpxPackage string
	HandlerDoc         string
	HandlerName        string
	RequestType        string
	LogicName          string
	LogicType          string
	Call               string
	HasResp            bool
	HasRequest         bool
	TransErr           bool
}

func genHandler(dir, rootPkg string, cfg *config.Config, group spec.Group, route spec.Route, trans bool) error {
	handler := getHandlerName(route)
	handlerPath := getHandlerFolderPath(group, route)
	pkgName := handlerPath[strings.LastIndex(handlerPath, "/")+1:]
	logicName := defaultLogicPackage
	if handlerPath != handlerDir {
		handler = cases.Title(language.English, cases.NoLower).String(handler)
		logicName = pkgName
	}

	// write doc for swagger
	var handlerDoc *strings.Builder
	handlerDoc = &strings.Builder{}

	var isParameterRequest bool
	if route.RequestType != nil && strings.Contains(strings.Join(route.RequestType.Documents(), ""),
		"swagger:parameters") {
		isParameterRequest = true
	} else {
		isParameterRequest = false
	}

	var swaggerPath string
	if strings.Contains(route.Path, ":") {
		swaggerPath = ConvertRoutePathToSwagger(route.Path)
	} else {
		swaggerPath = route.Path
	}

	prefix := group.GetAnnotation(spec.RoutePrefixKey)
	prefix = strings.ReplaceAll(prefix, `"`, "")
	prefix = strings.TrimSpace(prefix)
	if len(prefix) > 0 {
		prefix = path.Join("/", prefix)
	}
	swaggerPath = path.Join("/", prefix, swaggerPath)

	handlerDoc.WriteString(fmt.Sprintf("// swagger:route %s %s %s %s \n", route.Method, swaggerPath,
		group.GetAnnotation("group"), strings.TrimSuffix(handler, "Handler")))
	handlerDoc.WriteString("//\n")
	handlerDoc.WriteString(fmt.Sprintf("%s\n", strings.Join(route.HandlerDoc, " ")))
	handlerDoc.WriteString("//\n")
	handlerDoc.WriteString(fmt.Sprintf("%s\n", strings.Join(route.HandlerDoc, " ")))
	handlerDoc.WriteString("//\n")

	// HasRequest
	if len(route.RequestTypeName()) > 0 && !isParameterRequest {
		handlerDoc.WriteString(fmt.Sprintf(`// Parameters:
			//  + name: body
			//    require: true
			//    in: %s
			//    type: %s
			//
			`, "body", route.RequestTypeName()))
	}
	// HasResp
	if len(route.ResponseTypeName()) > 0 {
		handlerDoc.WriteString(fmt.Sprintf(`// Responses:
			//  200: %s`, route.ResponseTypeName()))
	}

	return doGenToFile(dir, handler, cfg, group, route, handlerInfo{
		PkgName:        pkgName,
		ImportPackages: genHandlerImports(group, route, rootPkg),
		HandlerDoc:     handlerDoc.String(),
		HandlerName:    handler,
		RequestType:    util.Title(route.RequestTypeName()),
		LogicName:      logicName,
		LogicType:      cases.Title(language.English, cases.NoLower).String(getLogicName(route)),
		Call:           cases.Title(language.English, cases.NoLower).String(strings.TrimSuffix(handler, "Handler")),
		HasResp:        len(route.ResponseTypeName()) > 0,
		HasRequest:     len(route.RequestTypeName()) > 0,
		TransErr:       trans,
	})
}

func doGenToFile(dir, handler string, cfg *config.Config, group spec.Group,
	route spec.Route, handleObj handlerInfo,
) error {
	filename, err := format.FileNamingFormat(cfg.NamingFormat, handler)
	if err != nil {
		return err
	}

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          getHandlerFolderPath(group, route),
		filename:        filename + ".go",
		templateName:    "handlerTemplate",
		category:        category,
		templateFile:    handlerTemplateFile,
		builtinTemplate: handlerTemplate,
		data:            handleObj,
	})
}

func genHandlers(dir, rootPkg string, cfg *config.Config, api *spec.ApiSpec, g *GenContext) error {
	for _, group := range api.Service.Groups {
		for _, route := range group.Routes {
			if err := genHandler(dir, rootPkg, cfg, group, route, g.TransErr); err != nil {
				return err
			}
		}
	}

	return nil
}

func genHandlerImports(group spec.Group, route spec.Route, parentPkg string) string {
	imports := []string{
		fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, getLogicFolderPath(group, route))),
		fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, contextDir)),
	}
	if len(route.RequestTypeName()) > 0 {
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

	return handler + "Logic"
}
