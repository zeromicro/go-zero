package gen

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/api/tsgen/template"
	apiutil "github.com/zeromicro/go-zero/tools/goctl/api/util"
	"github.com/zeromicro/go-zero/tools/goctl/util"
)

func GenHandler(dir, caller string, api *spec.ApiSpec, unwrapAPI bool, customBody bool, baseUrl string) error {
	filename := strings.Replace(api.Service.Name, "-api", "", 1)

	data := template.HandlerTemplateData{
		Caller:      caller,
		IsUnwrapAPI: unwrapAPI,
	}

	if len(api.Types) != 0 {
		data.ComponentName = apiutil.ComponentName(api)
	}

	apis, err := genAPI(api, customBody, baseUrl)
	if err != nil {
		return err
	}
	data.Routes = apis

	return template.GenTsFile(dir, filename, template.Handlers, data)
}

func genAPI(api *spec.ApiSpec, customBody bool, baseUrl string) ([]*template.HandlerRouteTemplateData, error) {
	routes := []*template.HandlerRouteTemplateData{}
	for _, group := range api.Service.Groups {
		for _, route := range group.Routes {
			hrt := template.HandlerRouteTemplateData{}
			handler := route.Handler
			if len(handler) == 0 {
				return nil, fmt.Errorf("missing handler annotation for route %q", route.Path)
			}

			handler = util.Untitle(handler)
			hrt.FuncName = strings.Replace(handler, "Handler", "", 1)
			hrt.Comment = commentForRoute(route, customBody)

			if customBody {
				hrt.GenericsTypes = "<T>"
			}
			hrt.FuncArgs = paramsForRoute(route, customBody)
			hrt.ResponseType = "null"
			if len(route.ResponseTypeName()) > 0 {
				val, err := GoTypeToTs(route.ResponseType, true)
				if err != nil {
					return nil, err
				}
				hrt.ResponseType = val
			}
			hrt.HttpMethod = strings.ToLower(route.Method)
			hrt.CallArgs = callParamsForRoute(route, group, customBody, baseUrl)
			routes = append(routes, &hrt)
		}
	}

	return routes, nil
}

func paramsForRoute(route spec.Route, customBody bool) string {
	if route.RequestType == nil {
		if customBody {
			return "body?: T"
		}
		return ""
	}
	hasParams := pathHasParams(route)
	hasBody := hasRequestBody(route)
	hasHeader := hasRequestHeader(route)
	hasPath := hasRequestPath(route)
	rt, err := GoTypeToTs(route.RequestType, true)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	var params []string
	if hasParams {
		params = append(params, fmt.Sprintf("params: %s", rt+"Params"))
	}
	if hasBody {
		params = append(params, fmt.Sprintf("req: %s", rt))
	}
	if hasHeader {
		params = append(params, fmt.Sprintf("headers: %s", rt+"Headers"))
	}
	if hasPath {
		ds, ok := route.RequestType.(spec.DefineStruct)
		if !ok {
			fmt.Printf("invalid route.RequestType: {%v}\n", route.RequestType)
		}
		members := ds.GetTagMembers(pathTagKey)
		for _, member := range members {
			tags := member.Tags()

			if len(tags) > 0 && tags[0].Key == pathTagKey {
				valueType, err := GoTypeToTs(member.Type, false)
				if err != nil {
					fmt.Println(err.Error())
					return ""
				}
				params = append(params, fmt.Sprintf("%s: %s", tags[0].Name, valueType))
			}
		}
	}

	if customBody {
		params = append(params, "body?: T")
	}
	return strings.Join(params, ", ")
}

func commentForRoute(route spec.Route, customBody bool) string {
	var builder strings.Builder
	comment := route.JoinedDoc()
	builder.WriteString("/**")
	builder.WriteString("\n * @description " + comment)
	hasParams := pathHasParams(route)
	hasBody := hasRequestBody(route)
	hasHeader := hasRequestHeader(route)
	if hasParams {
		builder.WriteString("\n * @param params")
	}
	if hasBody {
		builder.WriteString("\n * @param req")
	}
	if hasHeader {
		builder.WriteString("\n * @param headers")
	}
	if customBody {
		builder.WriteString("\n * @param body")
	}
	builder.WriteString("\n */")
	return builder.String()
}

func callParamsForRoute(route spec.Route, group spec.Group, customBody bool, baseUrl string) string {
	hasParams := pathHasParams(route)
	hasBody := hasRequestBody(route)
	hasHeader := hasRequestHeader(route)

	var params = []string{pathForRoute(route, group, baseUrl)}
	if hasParams {
		params = append(params, "params")
	} else {
		params = append(params, "null")
	}

	configParams := []string{}

	if hasBody {
		if customBody {
			configParams = append(configParams, "body: JSON.stringify(body ?? req)")
		} else {
			configParams = append(configParams, "body: JSON.stringify(req)")
		}
	} else if customBody {
		configParams = append(configParams, "body: body ? JSON.stringify(body): null")
	}
	if hasHeader {
		configParams = append(configParams, "headers: headers")
	}

	params = append(params, fmt.Sprintf("{%s}", strings.Join(configParams, ", ")))

	return strings.Join(params, ", ")
}

func pathForRoute(route spec.Route, group spec.Group, baseUrl string) string {
	prefix := group.GetAnnotation(pathPrefix)

	routePath := route.Path
	if strings.Contains(routePath, ":") {
		pathSlice := strings.Split(routePath, "/")
		for i, part := range pathSlice {
			if strings.Contains(part, ":") {
				pathSlice[i] = fmt.Sprintf("${%s}", part[1:])
			}
		}
		routePath = strings.Join(pathSlice, "/")
	}
	if len(prefix) == 0 {
		return "`" + baseUrl + routePath + "`"
	}

	prefix = strings.TrimPrefix(prefix, `"`)
	prefix = strings.TrimSuffix(prefix, `"`)
	return fmt.Sprintf("`%s%s/%s`", baseUrl, prefix, strings.TrimPrefix(routePath, "/"))
}

func pathHasParams(route spec.Route) bool {
	ds, ok := route.RequestType.(spec.DefineStruct)
	if !ok {
		return false
	}

	return len(ds.Members) != (len(ds.GetBodyMembers()) + len(ds.GetTagMembers(headerTagKey)) + len(ds.GetTagMembers(pathTagKey)))
}

func hasRequestBody(route spec.Route) bool {
	ds, ok := route.RequestType.(spec.DefineStruct)
	if !ok {
		return false
	}

	return len(route.RequestTypeName()) > 0 && len(ds.GetBodyMembers()) > 0
}

func hasRequestPath(route spec.Route) bool {
	ds, ok := route.RequestType.(spec.DefineStruct)
	if !ok {
		return false
	}

	return len(route.RequestTypeName()) > 0 && len(ds.GetTagMembers(pathTagKey)) > 0
}

func hasRequestHeader(route spec.Route) bool {
	ds, ok := route.RequestType.(spec.DefineStruct)
	if !ok {
		return false
	}

	return len(route.RequestTypeName()) > 0 && len(ds.GetTagMembers(headerTagKey)) > 0
}
