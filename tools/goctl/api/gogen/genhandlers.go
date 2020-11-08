package gogen

import (
	"bytes"
	"fmt"
	"path"
	"strings"
	"text/template"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	apiutil "github.com/tal-tech/go-zero/tools/goctl/api/util"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/vars"
)

const handlerTemplate = `package handler

import (
	"net/http"

	{{.ImportPackages}}
)

func {{.HandlerName}}(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		{{if .HasRequest}}var req types.{{.RequestType}}
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}{{end}}

		l := logic.New{{.LogicType}}(r.Context(), ctx)
		{{if .HasResp}}resp, {{end}}err := l.{{.Call}}({{if .HasRequest}}req{{end}})
		if err != nil {
			httpx.Error(w, err)
		} else {
			{{if .HasResp}}httpx.OkJson(w, resp){{else}}httpx.Ok(w){{end}}
		}
	}
}
`

type Handler struct {
	ImportPackages string
	HandlerName    string
	RequestType    string
	LogicType      string
	Call           string
	HasResp        bool
	HasRequest     bool
}

func genHandler(dir string, group spec.Group, route spec.Route) error {
	handler, ok := apiutil.GetAnnotationValue(route.Annotations, "server", "handler")
	if !ok {
		return fmt.Errorf("missing handler annotation for %q", route.Path)
	}

	handler = getHandlerName(handler)
	if getHandlerFolderPath(group, route) != handlerDir {
		handler = strings.Title(handler)
	}
	parentPkg, err := getParentPackage(dir)
	if err != nil {
		return err
	}

	return doGenToFile(dir, handler, group, route, Handler{
		ImportPackages: genHandlerImports(group, route, parentPkg),
		HandlerName:    handler,
		RequestType:    util.Title(route.RequestType.Name),
		LogicType:      strings.TrimSuffix(strings.Title(handler), "Handler") + "Logic",
		Call:           strings.Title(strings.TrimSuffix(handler, "Handler")),
		HasResp:        len(route.ResponseType.Name) > 0,
		HasRequest:     len(route.RequestType.Name) > 0,
	})
}

func doGenToFile(dir, handler string, group spec.Group, route spec.Route, handleObj Handler) error {
	if getHandlerFolderPath(group, route) != handlerDir {
		handler = strings.Title(handler)
	}
	filename := strings.ToLower(handler)
	if strings.HasSuffix(filename, "handler") {
		filename = filename + ".go"
	} else {
		filename = filename + "handler.go"
	}
	fp, created, err := apiutil.MaybeCreateFile(dir, getHandlerFolderPath(group, route), filename)
	if err != nil {
		return err
	}
	if !created {
		return nil
	}
	defer fp.Close()

	text, err := util.LoadTemplate(category, handlerTemplateFile, handlerTemplate)
	if err != nil {
		return err
	}

	buffer := new(bytes.Buffer)
	err = template.Must(template.New("handlerTemplate").Parse(text)).Execute(buffer, handleObj)
	if err != nil {
		return err
	}

	formatCode := formatCode(buffer.String())
	_, err = fp.WriteString(formatCode)
	return err
}

func genHandlers(dir string, api *spec.ApiSpec) error {
	for _, group := range api.Service.Groups {
		for _, route := range group.Routes {
			if err := genHandler(dir, group, route); err != nil {
				return err
			}
		}
	}

	return nil
}

func genHandlerImports(group spec.Group, route spec.Route, parentPkg string) string {
	var imports []string
	imports = append(imports, fmt.Sprintf("\"%s\"",
		util.JoinPackages(parentPkg, getLogicFolderPath(group, route))))
	imports = append(imports, fmt.Sprintf("\"%s\"", util.JoinPackages(parentPkg, contextDir)))
	if len(route.RequestType.Name) > 0 {
		imports = append(imports, fmt.Sprintf("\"%s\"\n", util.JoinPackages(parentPkg, typesDir)))
	}
	imports = append(imports, fmt.Sprintf("\"%s/rest/httpx\"", vars.ProjectOpenSourceUrl))

	return strings.Join(imports, "\n\t")
}

func getHandlerBaseName(handler string) string {
	handlerName := util.Untitle(handler)
	if strings.HasSuffix(handlerName, "handler") {
		handlerName = strings.ReplaceAll(handlerName, "handler", "")
	} else if strings.HasSuffix(handlerName, "Handler") {
		handlerName = strings.ReplaceAll(handlerName, "Handler", "")
	}
	return handlerName
}

func getHandlerFolderPath(group spec.Group, route spec.Route) string {
	folder, ok := apiutil.GetAnnotationValue(route.Annotations, "server", groupProperty)
	if !ok {
		folder, ok = apiutil.GetAnnotationValue(group.Annotations, "server", groupProperty)
		if !ok {
			return handlerDir
		}
	}
	folder = strings.TrimPrefix(folder, "/")
	folder = strings.TrimSuffix(folder, "/")
	return path.Join(handlerDir, folder)
}

func getHandlerName(handler string) string {
	return getHandlerBaseName(handler) + "Handler"
}
