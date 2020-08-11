package tsgen

import (
	"errors"
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

{{.types}}

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

	var localTypes []spec.Type
	for _, route := range api.Service.Routes {
		rts := apiutil.GetLocalTypes(api, route)
		localTypes = append(localTypes, rts...)
	}

	var prefixForType = func(ty string) string {
		if _, pri := primitiveType(ty); pri {
			return ""
		}
		for _, item := range localTypes {
			if util.Title(item.Name) == ty {
				return ""
			}
		}
		return packagePrefix
	}

	types, err := genTypes(localTypes, func(name string) (*spec.Type, error) {
		for _, ty := range api.Types {
			if strings.ToLower(ty.Name) == strings.ToLower(name) {
				return &ty, nil
			}
		}
		return nil, errors.New("inline type " + name + " not exist, please correct api file")
	}, prefixForType)
	if err != nil {
		return err
	}

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
	shardTypes := apiutil.GetSharedTypes(api)
	if len(shardTypes) != 0 {
		if len(imports) > 0 {
			imports += "\n"
		}
		outputFile := apiutil.ComponentName(api)
		imports += fmt.Sprintf(`import * as components from "%s"`, "./"+outputFile)
	}

	apis, err := genApi(api, localTypes, caller, prefixForType)
	if err != nil {
		return err
	}

	t := template.Must(template.New("handlerTemplate").Parse(handlerTemplate))
	return t.Execute(fp, map[string]string{
		"webApi":  webApi,
		"types":   strings.TrimSpace(types),
		"imports": imports,
		"apis":    strings.TrimSpace(apis),
	})
}

func genTypes(localTypes []spec.Type, inlineType func(string) (*spec.Type, error), prefixForType func(string) string) (string, error) {
	var builder strings.Builder
	var first bool

	for _, tp := range localTypes {
		if first {
			first = false
		} else {
			fmt.Fprintln(&builder)
		}
		if err := writeType(&builder, tp, func(name string) (s *spec.Type, err error) {
			return inlineType(name)
		}, prefixForType); err != nil {
			return "", err
		}
	}
	types := builder.String()
	return types, nil
}

func genApi(api *spec.ApiSpec, localTypes []spec.Type, caller string, prefixForType func(string) string) (string, error) {
	var builder strings.Builder
	for _, route := range api.Service.Routes {
		handler, ok := apiutil.GetAnnotationValue(route.Annotations, "server", "handler")
		if !ok {
			return "", fmt.Errorf("missing handler annotation for route %q", route.Path)
		}
		handler = util.Untitle(handler)
		handler = strings.Replace(handler, "Handler", "", 1)
		comment := commentForRoute(route)
		if len(comment) > 0 {
			fmt.Fprintf(&builder, "%s\n", comment)
		}
		fmt.Fprintf(&builder, "export function %s(%s) {\n", handler, paramsForRoute(route, prefixForType))
		writeIndent(&builder, 1)
		responseGeneric := "<null>"
		if len(route.ResponseType.Name) > 0 {
			val, err := goTypeToTs(route.ResponseType.Name, prefixForType)
			if err != nil {
				return "", err
			}
			responseGeneric = fmt.Sprintf("<%s>", val)
		}
		fmt.Fprintf(&builder, `return %s.%s%s(%s)`, caller, strings.ToLower(route.Method),
			util.Title(responseGeneric), callParamsForRoute(route))
		builder.WriteString("\n}\n\n")
	}

	apis := builder.String()
	return apis, nil
}

func paramsForRoute(route spec.Route, prefixForType func(string) string) string {
	hasParams := pathHasParams(route)
	hasBody := hasRequestBody(route)
	rt, err := goTypeToTs(route.RequestType.Name, prefixForType)
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
	comment, _ := apiutil.GetAnnotationValue(route.Annotations, "doc", "summary")
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

func callParamsForRoute(route spec.Route) string {
	hasParams := pathHasParams(route)
	hasBody := hasRequestBody(route)
	if hasParams && hasBody {
		return fmt.Sprintf("%s, %s, %s", pathForRoute(route), "params", "req")
	} else if hasParams {
		return fmt.Sprintf("%s, %s", pathForRoute(route), "params")
	} else if hasBody {
		return fmt.Sprintf("%s, %s", pathForRoute(route), "req")
	}
	return pathForRoute(route)
}

func pathForRoute(route spec.Route) string {
	return "\"" + route.Path + "\""
}

func pathHasParams(route spec.Route) bool {
	return len(route.RequestType.Members) != len(route.RequestType.GetBodyMembers())
}

func hasRequestBody(route spec.Route) bool {
	return len(route.RequestType.Name) > 0 && len(route.RequestType.GetBodyMembers()) > 0
}
