package gogen

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

//go:embed logic_test.tpl
var logicTestTemplate string

func genLogicTest(dir, rootPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	for _, g := range api.Service.Groups {
		for _, r := range g.Routes {
			err := genLogicTestByRoute(dir, rootPkg, cfg, g, r)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func genLogicTestByRoute(dir, rootPkg string, cfg *config.Config, group spec.Group, route spec.Route) error {
	logic := getLogicName(route)
	goFile, err := format.FileNamingFormat(cfg.NamingFormat, logic)
	if err != nil {
		return err
	}

	imports := genLogicTestImports(route, rootPkg)
	var responseString string
	var returnString string
	var requestString string
	var requestType string
	if len(route.ResponseTypeName()) > 0 {
		resp := responseGoTypeName(route, typesPacket)
		responseString = "(resp " + resp + ", err error)"
		returnString = "return"
	} else {
		responseString = "error"
		returnString = "return nil"
	}
	if len(route.RequestTypeName()) > 0 {
		requestString = "req *" + requestGoTypeName(route, typesPacket)
		requestType = requestGoTypeName(route, typesPacket)
	}

	subDir := getLogicFolderPath(group, route)
	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          subDir,
		filename:        goFile + "_test.go",
		templateName:    "logicTestTemplate",
		category:        category,
		templateFile:    logicTestTemplateFile,
		builtinTemplate: logicTestTemplate,
		data: map[string]any{
			"pkgName":      subDir[strings.LastIndex(subDir, "/")+1:],
			"imports":      imports,
			"logic":        strings.Title(logic),
			"function":     strings.Title(strings.TrimSuffix(logic, "Logic")),
			"responseType": responseString,
			"returnString": returnString,
			"request":      requestString,
			"hasRequest":   len(requestType) > 0,
			"hasResponse":  len(route.ResponseTypeName()) > 0,
			"requestType":  requestType,
			"hasDoc":       len(route.JoinedDoc()) > 0,
			"doc":          getDoc(route.JoinedDoc()),
		},
	})
}

func genLogicTestImports(route spec.Route, parentPkg string) string {
	var imports []string
	//imports = append(imports, `"context"`+"\n")
	imports = append(imports, fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, contextDir)))
	imports = append(imports, fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, configDir)))
	if shallImportTypesPackage(route) {
		imports = append(imports, fmt.Sprintf("\"%s\"\n", pathx.JoinPackages(parentPkg, typesDir)))
	}
	//imports = append(imports, fmt.Sprintf("\"%s/core/logx\"", vars.ProjectOpenSourceURL))
	return strings.Join(imports, "\n\t")
}
