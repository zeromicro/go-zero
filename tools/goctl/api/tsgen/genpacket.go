package tsgen

import (
	"fmt"
	"path"
	"strings"
	"text/template"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	apiutil "github.com/tal-tech/go-zero/tools/goctl/api/util"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

const (
	handlerTemplate = `{{.imports}}

{{.apis}}
`
)

func genHandler(dir, webApi, caller string, api *spec.ApiSpec, unwrapApi bool) error {
	filename := strings.Replace(api.Service.Name, "-api", "", 1) + ".ts"
	if err := util.RemoveIfExist(path.Join(dir, filename)); err != nil {
		return err
	}
	fp, created, err := apiutil.MaybeCreateFile(dir, "", filename)
	if err != nil {
		return err
	}
	if !created {
		return nil
	}
	defer fp.Close()

	imports := ""
	if len(caller) == 0 {
		caller = "webapi"
	}
	importCaller := caller
	if unwrapApi {
		importCaller = "{ " + importCaller + " }"
	}
	if len(webApi) > 0 {
		imports += `import ` + importCaller + ` from ` + "\"" + webApi + "\""
	}

	if len(api.Types) != 0 {
		if len(imports) > 0 {
			imports += util.NL
		}
		outputFile := apiutil.ComponentName(api)
		imports += fmt.Sprintf(`import * as components from "%s"`, "./"+outputFile)
		imports += fmt.Sprintf(`%sexport * from "%s"`, util.NL, "./"+outputFile)
	}

	apis, err := genApi(api, caller)
	if err != nil {
		return err
	}

	t := template.Must(template.New("handlerTemplate").Parse(handlerTemplate))
	return t.Execute(fp, map[string]string{
		"imports": imports,
		"apis":    strings.TrimSpace(apis),
	})
}

func genApi(api *spec.ApiSpec, caller string) (string, error) {
	var builder strings.Builder
	for _, group := range api.Service.Groups {
		for _, route := range group.Routes {
			handler := route.Handler
			if len(handler) == 0 {
				return "", fmt.Errorf("missing handler annotation for route %q", route.Path)
			}

			handler = util.Untitle(handler)
			handler = strings.Replace(handler, "Handler", "", 1)
			comment := commentForRoute(route)
			if len(comment) > 0 {
				fmt.Fprintf(&builder, "%s\n", comment)
			}
			fmt.Fprintf(&builder, "export function %s(%s) {\n", handler, paramsForRoute(route))
			writeIndent(&builder, 1)
			responseGeneric := "<null>"
			if len(route.ResponseTypeName()) > 0 {
				val, err := goTypeToTs(route.ResponseType, true)
				if err != nil {
					return "", err
				}

				responseGeneric = fmt.Sprintf("<%s>", val)
			}
			fmt.Fprintf(&builder, `return %s.%s%s(%s)`, caller, strings.ToLower(route.Method),
				util.Title(responseGeneric), callParamsForRoute(route, group))
			builder.WriteString("\n}\n\n")
		}
	}

	apis := builder.String()
	return apis, nil
}

func paramsForRoute(route spec.Route) string {
	if route.RequestType == nil {
		return ""
	}
	hasParams := pathHasParams(route)
	hasBody := hasRequestBody(route)
	rt, err := goTypeToTs(route.RequestType, true)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	if hasParams && hasBody {
		return fmt.Sprintf("params: %s, req: %s", rt+"Params", rt)
	} else if hasParams {
		return fmt.Sprintf("params: %s", rt+"Params")
	} else if hasBody {
		return fmt.Sprintf("req: %s", rt)
	}
	return ""
}

func commentForRoute(route spec.Route) string {
	var builder strings.Builder
	comment := route.JoinedDoc()
	builder.WriteString("/**")
	builder.WriteString("\n * @description " + comment)
	hasParams := pathHasParams(route)
	hasBody := hasRequestBody(route)
	if hasParams && hasBody {
		builder.WriteString("\n * @param params")
		builder.WriteString("\n * @param req")
	} else if hasParams {
		builder.WriteString("\n * @param params")
	} else if hasBody {
		builder.WriteString("\n * @param req")
	}
	builder.WriteString("\n */")
	return builder.String()
}

func callParamsForRoute(route spec.Route, group spec.Group) string {
	hasParams := pathHasParams(route)
	hasBody := hasRequestBody(route)
	if hasParams && hasBody {
		return fmt.Sprintf("%s, %s, %s", pathForRoute(route, group), "params", "req")
	} else if hasParams {
		return fmt.Sprintf("%s, %s", pathForRoute(route, group), "params")
	} else if hasBody {
		return fmt.Sprintf("%s, %s", pathForRoute(route, group), "req")
	}

	return pathForRoute(route, group)
}

func pathForRoute(route spec.Route, group spec.Group) string {
	prefix := group.GetAnnotation("pathPrefix")
	if len(prefix) == 0 {
		return "\"" + route.Path + "\""
	} else {
		prefix = strings.TrimPrefix(prefix, `"`)
		prefix = strings.TrimSuffix(prefix, `"`)
		return fmt.Sprintf(`"%s/%s"`, prefix, strings.TrimPrefix(route.Path, "/"))
	}
}

func pathHasParams(route spec.Route) bool {
	ds, ok := route.RequestType.(spec.DefineStruct)
	if !ok {
		return false
	}

	return len(ds.Members) != len(ds.GetBodyMembers())
}

func hasRequestBody(route spec.Route) bool {
	ds, ok := route.RequestType.(spec.DefineStruct)
	if !ok {
		return false
	}

	return len(route.RequestTypeName()) > 0 && len(ds.GetBodyMembers()) > 0
}
