package swagger

import (
	"net/http"
	"path"
	"strings"

	"github.com/go-openapi/spec"
	apiSpec "github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func spec2Paths(info apiSpec.Info, srv apiSpec.Service) *spec.Paths {
	paths := &spec.Paths{
		Paths: make(map[string]spec.PathItem),
	}
	for _, group := range srv.Groups {
		prefix := path.Clean(strings.TrimPrefix(group.GetAnnotation("prefix"), "/"))
		for _, route := range group.Routes {
			routPath := pathVariable2SwaggerVariable(route.Path)
			if len(prefix) > 0 && prefix != "." {
				routPath = "/" + path.Clean(prefix) + routPath
			}
			pathItem := spec2Path(info, group, route)
			existPathItem, ok := paths.Paths[routPath]
			if !ok {
				paths.Paths[routPath] = pathItem
			} else {
				paths.Paths[routPath] = mergePathItem(existPathItem, pathItem)
			}
		}
	}
	return paths
}

func mergePathItem(old, new spec.PathItem) spec.PathItem {
	if new.Get != nil {
		old.Get = new.Get
	}
	if new.Put != nil {
		old.Put = new.Put
	}
	if new.Post != nil {
		old.Post = new.Post
	}
	if new.Delete != nil {
		old.Delete = new.Delete
	}
	if new.Options != nil {
		old.Options = new.Options
	}
	if new.Head != nil {
		old.Head = new.Head
	}
	if new.Patch != nil {
		old.Patch = new.Patch
	}
	if new.Parameters != nil {
		old.Parameters = new.Parameters
	}
	return old
}

func spec2Path(info apiSpec.Info, group apiSpec.Group, route apiSpec.Route) spec.PathItem {
	needJwt := hasKey(group.Annotation.Properties, "jwt")
	var security []map[string][]string
	if needJwt {
		security = []map[string][]string{
			{
				swaggerSecurityDefinitionBearerAuth: []string{},
			},
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
			Responses:   jsonResponseFromType(info, route.ResponseType),
			Security:    security,
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
