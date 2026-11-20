package gogen

import (
	_ "embed"
	"fmt"
	"path"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/internal/version"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

const defaultLogicPackage = "logic"

var (
	//go:embed handler.tpl
	handlerTemplate string
	//go:embed sse_handler.tpl
	sseHandlerTemplate string
)

func genHandlers(dir, rootPkg, projectPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	type fileKey struct {
		subdir   string
		filename string
	}

	type handlerFile struct {
		pkgName   string
		imports   []string
		handlers  []map[string]any
		sseEnable bool
	}

	files := make(map[fileKey]*handlerFile)

	for _, group := range api.Service.Groups {
		sse := group.GetAnnotation("sse") == "true"
		for _, route := range group.Routes {
			handlerName := getHandlerName(route)
			handlerPath := getHandlerFolderPath(group, route)
			pkgName := handlerPath[strings.LastIndex(handlerPath, "/")+1:]
			logicName := defaultLogicPackage
			if handlerPath != handlerDir {
				handlerName = strings.Title(handlerName)
				logicName = pkgName
			}

			fileBase := route.GetAnnotation(filenameProperty)
			if len(fileBase) == 0 {
				fileBase = group.GetAnnotation(filenameProperty)
			}
			if len(fileBase) == 0 {
				fileBase = handlerName
			}

			fileBase, err := format.FileNamingFormat(cfg.NamingFormat, fileBase)
			if err != nil {
				return err
			}

			subdir := getHandlerFolderPath(group, route)
			key := fileKey{
				subdir:   subdir,
				filename: fileBase + ".go",
			}

			f, ok := files[key]
			if !ok {
				f = &handlerFile{
					pkgName:   pkgName,
					imports:   nil,
					handlers:  nil,
					sseEnable: sse,
				}
				files[key] = f
			}

			importsStr := genHandlerImports(group, route, rootPkg)
			importMap := make(map[string]bool)
			for _, existing := range f.imports {
				importMap[existing] = true
			}
			for _, imp := range strings.Split(importsStr, "\n\t") {
				imp = strings.TrimSpace(imp)
				if len(imp) > 0 && !importMap[imp] {
					importMap[imp] = true
					f.imports = append(f.imports, imp)
				}
			}

			handlerData := map[string]any{
				"HandlerName":  handlerName,
				"RequestType":  util.Title(route.RequestTypeName()),
				"ResponseType": responseGoTypeName(route, typesPacket),
				"LogicName":    logicName,
				"LogicType":    strings.Title(getLogicName(route)),
				"Call":         strings.Title(strings.TrimSuffix(handlerName, "Handler")),
				"HasResp":      len(route.ResponseTypeName()) > 0,
				"HasRequest":   len(route.RequestTypeName()) > 0,
				"HasDoc":       len(route.JoinedDoc()) > 0,
				"Doc":          getDoc(route.JoinedDoc()),
			}

			f.handlers = append(f.handlers, handlerData)
		}
	}

	for key, f := range files {
		importsJoined := strings.Join(f.imports, "\n\t")

		builtinTemplate := handlerTemplate
		templateFile := handlerTemplateFile
		if f.sseEnable {
			builtinTemplate = sseHandlerTemplate
			templateFile = sseHandlerTemplateFile
		}

		if err := genFile(fileGenConfig{
			dir:             dir,
			subdir:          key.subdir,
			filename:        key.filename,
			templateName:    "handlerTemplate",
			category:        category,
			templateFile:    templateFile,
			builtinTemplate: builtinTemplate,
			data: map[string]any{
				"PkgName":        f.pkgName,
				"ImportPackages": importsJoined,
				"Handlers":       f.handlers,
				"projectPkg":     projectPkg,
				"version":        version.BuildVersion,
			},
		}); err != nil {
			return err
		}
	}

	return nil
}

func genHandlerImports(group spec.Group, route spec.Route, parentPkg string) string {
	imports := []string{
		fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, getLogicFolderPath(group, route))),
		fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, contextDir)),
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

	return handler + "Logic"
}
