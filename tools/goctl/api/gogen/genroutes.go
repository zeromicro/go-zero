package gogen

import (
	"bytes"
	"fmt"
	"path"
	"sort"
	"strings"
	"text/template"

	"github.com/tal-tech/go-zero/core/collection"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	apiutil "github.com/tal-tech/go-zero/tools/goctl/api/util"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/vars"
)

const (
	routesFilename = "routes.go"
	routesTemplate = `// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	{{.importPackages}}
)

func RegisterHandlers(engine *rest.Server, serverCtx *svc.ServiceContext) {
	{{.routesAdditions}}
}
`
	routesAdditionTemplate = `
	engine.AddRoutes(
		{{.routes}}
	{{.jwt}}{{.signature}})
`
)

var mapping = map[string]string{
	"delete": "http.MethodDelete",
	"get":    "http.MethodGet",
	"head":   "http.MethodHead",
	"post":   "http.MethodPost",
	"put":    "http.MethodPut",
	"patch":  "http.MethodPatch",
}

type (
	group struct {
		routes           []route
		jwtEnabled       bool
		signatureEnabled bool
		authName         string
		middleware       []string
	}
	route struct {
		method  string
		path    string
		handler string
	}
)

func genRoutes(dir string, api *spec.ApiSpec, force bool) error {
	var builder strings.Builder
	groups, err := getRoutes(api)
	if err != nil {
		return err
	}

	gt := template.Must(template.New("groupTemplate").Parse(routesAdditionTemplate))
	for _, g := range groups {
		var gbuilder strings.Builder
		for _, r := range g.routes {
			fmt.Fprintf(&gbuilder, `
		{
			Method:  %s,
			Path:    "%s",
			Handler: %s,
		},`,
				r.method, r.path, r.handler)
		}
		var jwt string
		if g.jwtEnabled {
			jwt = fmt.Sprintf(", rest.WithJwt(serverCtx.Config.%s.AccessSecret)", g.authName)
		}
		var signature string
		if g.signatureEnabled {
			signature = fmt.Sprintf(", rest.WithSignature(serverCtx.Config.%s.Signature)", g.authName)
		}

		var routes string
		if len(g.middleware) > 0 {
			var params = g.middleware
			for i := range params {
				params[i] = "serverCtx." + params[i]
			}
			var middlewareStr = strings.Join(params, ", ")
			routes = fmt.Sprintf("rest.WithMiddlewares(\n[]rest.Middleware{ %s }, \n[]rest.Route{\n %s \n}...,\n),",
				middlewareStr, strings.TrimSpace(gbuilder.String()))
		} else {
			routes = fmt.Sprintf("[]rest.Route{\n %s \n},", strings.TrimSpace(gbuilder.String()))
		}

		if err := gt.Execute(&builder, map[string]string{
			"routes":    routes,
			"jwt":       jwt,
			"signature": signature,
		}); err != nil {
			return err
		}
	}

	parentPkg, err := getParentPackage(dir)
	if err != nil {
		return err
	}

	filename := path.Join(dir, handlerDir, routesFilename)
	if !force {
		if err := util.RemoveOrQuit(filename); err != nil {
			return err
		}
	}

	fp, created, err := apiutil.MaybeCreateFile(dir, handlerDir, routesFilename)
	if err != nil {
		return err
	}
	if !created {
		return nil
	}
	defer fp.Close()

	t := template.Must(template.New("routesTemplate").Parse(routesTemplate))
	buffer := new(bytes.Buffer)
	err = t.Execute(buffer, map[string]string{
		"importPackages":  genRouteImports(parentPkg, api),
		"routesAdditions": strings.TrimSpace(builder.String()),
	})
	if err != nil {
		return nil
	}
	formatCode := formatCode(buffer.String())
	_, err = fp.WriteString(formatCode)
	return err
}

func genRouteImports(parentPkg string, api *spec.ApiSpec) string {
	var importSet = collection.NewSet()
	importSet.AddStr(fmt.Sprintf("\"%s\"", util.JoinPackages(parentPkg, contextDir)))
	for _, group := range api.Service.Groups {
		for _, route := range group.Routes {
			folder, ok := apiutil.GetAnnotationValue(route.Annotations, "server", groupProperty)
			if !ok {
				folder, ok = apiutil.GetAnnotationValue(group.Annotations, "server", groupProperty)
				if !ok {
					continue
				}
			}
			importSet.AddStr(fmt.Sprintf("%s \"%s\"", folder,
				util.JoinPackages(parentPkg, handlerDir, folder)))
		}
	}
	imports := importSet.KeysStr()
	sort.Strings(imports)
	projectSection := strings.Join(imports, "\n\t")
	depSection := fmt.Sprintf("\"%s/rest\"", vars.ProjectOpenSourceUrl)
	return fmt.Sprintf("%s\n\n\t%s", projectSection, depSection)
}

func getRoutes(api *spec.ApiSpec) ([]group, error) {
	var routes []group

	for _, g := range api.Service.Groups {
		var groupedRoutes group
		for _, r := range g.Routes {
			handler, ok := apiutil.GetAnnotationValue(r.Annotations, "server", "handler")
			if !ok {
				return nil, fmt.Errorf("missing handler annotation for route %q", r.Path)
			}
			handler = getHandlerBaseName(handler) + "Handler(serverCtx)"
			folder, ok := apiutil.GetAnnotationValue(r.Annotations, "server", groupProperty)
			if ok {
				handler = folder + "." + strings.ToUpper(handler[:1]) + handler[1:]
			} else {
				folder, ok = apiutil.GetAnnotationValue(g.Annotations, "server", groupProperty)
				if ok {
					handler = folder + "." + strings.ToUpper(handler[:1]) + handler[1:]
				}
			}
			groupedRoutes.routes = append(groupedRoutes.routes, route{
				method:  mapping[r.Method],
				path:    r.Path,
				handler: handler,
			})
		}

		if value, ok := apiutil.GetAnnotationValue(g.Annotations, "server", "jwt"); ok {
			groupedRoutes.authName = value
			groupedRoutes.jwtEnabled = true
		}
		if value, ok := apiutil.GetAnnotationValue(g.Annotations, "server", "middleware"); ok {
			for _, item := range strings.Split(value, ",") {
				groupedRoutes.middleware = append(groupedRoutes.middleware, item)
			}
		}
		routes = append(routes, groupedRoutes)
	}

	return routes, nil
}
