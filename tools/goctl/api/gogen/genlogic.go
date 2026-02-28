package gogen

import (
	_ "embed"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/parser/g4/gen/api"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/internal/version"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"github.com/zeromicro/go-zero/tools/goctl/vars"
)

var (
	//go:embed logic.tpl
	logicTemplate string

	//go:embed sse_logic.tpl
	sseLogicTemplate string
)

func genLogic(dir, rootPkg, projectPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	type fileKey struct {
		subdir   string
		filename string
	}

	type logicFile struct {
		pkgName   string
		imports   []string
		logics    []map[string]any
		sseEnable bool
	}

	files := make(map[fileKey]*logicFile)

	for _, group := range api.Service.Groups {
		sse := group.GetAnnotation("sse") == "true"
		for _, route := range group.Routes {
			logic := getLogicName(route)

			fileBase := route.GetAnnotation(filenameProperty)
			if len(fileBase) == 0 {
				fileBase = group.GetAnnotation(filenameProperty)
			}
			if len(fileBase) == 0 {
				fileBase = logic
			}

			goFile, err := format.FileNamingFormat(cfg.NamingFormat, fileBase)
			if err != nil {
				return err
			}

			subDir := getLogicFolderPath(group, route)
			key := fileKey{
				subdir:   subDir,
				filename: goFile + ".go",
			}

			f, ok := files[key]
			if !ok {
				f = &logicFile{
					pkgName:   subDir[strings.LastIndex(subDir, "/")+1:],
					imports:   nil,
					logics:    nil,
					sseEnable: sse,
				}
				files[key] = f
			}

			imports := genLogicImports(route, rootPkg)
			importMap := make(map[string]bool)
			for _, existing := range f.imports {
				importMap[existing] = true
			}
			for _, imp := range strings.Split(imports, "\n\t") {
				imp = strings.TrimSpace(imp)
				if len(imp) > 0 && !importMap[imp] {
					importMap[imp] = true
					f.imports = append(f.imports, imp)
				}
			}

			var responseString string
			var returnString string
			var requestString string
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
			}

			if sse {
				responseString = "error"
				returnString = "return nil"
				resp := responseGoTypeName(route, typesPacket)
				if len(requestString) == 0 {
					requestString = "client chan<- " + resp
				} else {
					requestString += ", client chan<- " + resp
				}
			}

			logicData := map[string]any{
				"logic":        strings.Title(logic),
				"function":     strings.Title(strings.TrimSuffix(logic, "Logic")),
				"responseType": responseString,
				"returnString": returnString,
				"request":      requestString,
				"hasDoc":       len(route.JoinedDoc()) > 0,
				"doc":          getDoc(route.JoinedDoc()),
			}

			f.logics = append(f.logics, logicData)
		}
	}

	for key, f := range files {
		importsJoined := strings.Join(f.imports, "\n\t")

		builtinTemplate := logicTemplate
		templateFile := logicTemplateFile
		if f.sseEnable {
			builtinTemplate = sseLogicTemplate
			templateFile = sseLogicTemplateFile
		}

		if err := genFile(fileGenConfig{
			dir:             dir,
			subdir:          key.subdir,
			filename:        key.filename,
			templateName:    "logicTemplate",
			category:        category,
			templateFile:    templateFile,
			builtinTemplate: builtinTemplate,
			data: map[string]any{
				"pkgName":    f.pkgName,
				"imports":    importsJoined,
				"Logics":     f.logics,
				"projectPkg": projectPkg,
				"version":    version.BuildVersion,
			},
		}); err != nil {
			return err
		}
	}

	return nil
}

func getLogicFolderPath(group spec.Group, route spec.Route) string {
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
