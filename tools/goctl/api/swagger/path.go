package swagger

import (
	"encoding/json"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-openapi/spec"
	apiSpec "github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func spec2Paths(info apiSpec.Info, srv apiSpec.Service) *spec.Paths {
	paths := &spec.Paths{
		Paths: make(map[string]spec.PathItem),
	}
	for _, group := range srv.Groups {
		for _, route := range group.Routes {
			path := pathVariable2SwaggerVariable(route.Path)
			paths.Paths[path] = spec2Path(info, group, route)
		}
	}
	return paths
}

func getExample(properties map[string]string) map[int]map[string]any {
	var example = map[int]map[string]any{}
	for k, v := range properties {
		exampleVal := util.Unquote(v)
		if !strings.HasPrefix(k, "status") {
			continue
		}
		statusCode, _ := strconv.ParseInt(strings.TrimPrefix(k, "status"), 10, 32)
		if statusCode == 0 {
			continue
		}
		code := int(statusCode)
		text := http.StatusText(code)
		if len(text) == 0 {
			continue
		}
		if len(exampleVal) == 0 {
			example[code] = map[string]any{}
			continue
		}
		var v map[string]any
		if err := json.Unmarshal([]byte(exampleVal), &v); err != nil {
			example[code] = map[string]any{}
			continue
		}
		example[code] = v
	}
	return example
}

func spec2Path(info apiSpec.Info, group apiSpec.Group, route apiSpec.Route) spec.PathItem {
	globalExample := getExample(info.Properties)
	pathExample := getExample(route.AtDoc.Properties)
	// add global example to path example
	for k, v := range globalExample {
		if _, ok := pathExample[k]; !ok {
			pathExample[k] = v
		}
	}
	op := &spec.Operation{
		OperationProps: spec.OperationProps{
			Description: getStringFromKVOrDefault(route.AtDoc.Properties, "description", ""),
			Consumes:    consumesFromTypeOrDef(route.Method, route.RequestType),
			Produces:    getListFromInfoOrDefault(route.AtDoc.Properties, "produces", []string{applicationJson}),
			Schemes:     getListFromInfoOrDefault(route.AtDoc.Properties, "schemes", []string{schemeHttps}),
			Tags:        getListFromInfoOrDefault(group.Annotation.Properties, "tags", []string{""}),
			Summary:     getStringFromKVOrDefault(route.AtDoc.Properties, "summary", ""),
			Deprecated:  getBoolFromKVOrDefault(route.AtDoc.Properties, "deprecated", false),
			Parameters:  parametersFromType(route.Method, route.RequestType),
			Responses:   jsonResponseFromType(route.ResponseType, pathExample),
		},
	}
	externalDocsDescription := getStringFromKVOrDefault(route.AtDoc.Properties, "externalDocsDescription", "")
	externalDocsURL := getStringFromKVOrDefault(route.AtDoc.Properties, "externalDocsURL", "")
	if len(externalDocsDescription) > 0 || len(externalDocsURL) > 0 {
		op.ExternalDocs = &spec.ExternalDocumentation{
			Description: externalDocsDescription,
			URL:         externalDocsURL,
		}

	}
	item := spec.PathItem{}
	switch strings.ToUpper(route.Method) {
	case http.MethodGet:
		item.Get = op
	case http.MethodHead:
		item.Head = op
	case http.MethodPost:
		item.Post = op
	case http.MethodPut:
		item.Put = op
	case http.MethodPatch:
		item.Patch = op
	case http.MethodDelete:
		item.Delete = op
	case http.MethodOptions:
		item.Options = op
	default: // [http.MethodConnect,http.MethodTrace] not supported
	}
	return item
}
