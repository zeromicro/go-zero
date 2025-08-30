package rabbitmqgen

import (
	_ "embed"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/lerity-yao/go-zero/tools/cztctl/api/gogen"
	"github.com/lerity-yao/go-zero/tools/cztctl/api/parser/g4/gen/api"
	"github.com/lerity-yao/go-zero/tools/cztctl/api/spec"
	"github.com/lerity-yao/go-zero/tools/cztctl/config"
	"github.com/lerity-yao/go-zero/tools/cztctl/util/format"
	"github.com/lerity-yao/go-zero/tools/cztctl/util/pathx"
	"github.com/lerity-yao/go-zero/tools/cztctl/vars"
)

var (
	//go:embed logic.tpl
	logicTemplate string
)

func genLogic(dir, rootPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	for _, g := range api.Service.Groups {
		for _, r := range g.Routes {
			err := genLogicByRoute(dir, rootPkg, cfg, g, r)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func genLogicByRoute(dir, rootPkg string, cfg *config.Config, group spec.Group, route spec.Route) error {
	logic := getLogicName(route)
	goFile, err := format.FileNamingFormat(cfg.NamingFormat, logic)
	if err != nil {
		return err
	}

	imports := genLogicImports(route, rootPkg)

	subDir := GetLogicFolderPath(group, route)
	builtinTemplate := logicTemplate

	return gogen.GenFile(gogen.FileGenConfig{
		Dir:             dir,
		Subdir:          subDir,
		Filename:        goFile + ".go",
		TemplateName:    "logicTemplate",
		Category:        category,
		TemplateFile:    logicTemplateFile,
		BuiltinTemplate: builtinTemplate,
		Data: map[string]any{
			"pkgName": subDir[strings.LastIndex(subDir, "/")+1:],
			"imports": imports,
			"logic":   strings.Title(logic),
			"hasDoc":  len(route.JoinedDoc()) > 0,
			"doc":     GetDoc(route.JoinedDoc()),
		},
	})
}

func GetLogicFolderPath(group spec.Group, route spec.Route) string {
	folder := route.GetAnnotation(groupProperty)
	if len(folder) == 0 {
		folder = group.GetAnnotation(groupProperty)
		if len(folder) == 0 {
			return logicDir
		}
	}
	folder = strings.TrimPrefix(folder, "/")
	folder = strings.TrimSuffix(folder, "/")
	return path.Join(logicDir, folder)
}

func genLogicImports(route spec.Route, parentPkg string) string {
	var imports []string
	imports = append(imports, `"context"`+"\n")
	imports = append(imports, fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, contextDir)))
	if shallImportTypesPackage(route) {
		imports = append(imports, fmt.Sprintf("\"%s\"\n", pathx.JoinPackages(parentPkg, typesDir)))
	}
	imports = append(imports, fmt.Sprintf("\"%s/core/logx\"", vars.ProjectOpenSourceURL))
	return strings.Join(imports, "\n\t")
}

func onlyPrimitiveTypes(val string) bool {
	fields := strings.FieldsFunc(val, func(r rune) bool {
		return r == '[' || r == ']' || r == ' '
	})

	for _, field := range fields {
		if field == "map" {
			continue
		}
		// ignore array dimension number, like [5]int
		if _, err := strconv.Atoi(field); err == nil {
			continue
		}
		if !api.IsBasicType(field) {
			return false
		}
	}

	return true
}

func shallImportTypesPackage(route spec.Route) bool {
	if len(route.RequestTypeName()) > 0 {
		return true
	}

	respTypeName := route.ResponseTypeName()
	if len(respTypeName) == 0 {
		return false
	}

	if onlyPrimitiveTypes(respTypeName) {
		return false
	}

	return true
}
