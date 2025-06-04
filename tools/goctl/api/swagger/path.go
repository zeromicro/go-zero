package swagger

import (
	"net/http"
	"path"
	"strings"

	"github.com/go-openapi/spec"
	apiSpec "github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"
)

func spec2Paths(ctx Context, srv apiSpec.Service) *spec.Paths {
	paths := &spec.Paths{
		Paths: make(map[string]spec.PathItem),
	}
	for _, group := range srv.Groups {
		prefix := path.Clean(strings.TrimPrefix(group.GetAnnotation(propertyKeyPrefix), "/"))
		for _, route := range group.Routes {
			routPath := pathVariable2SwaggerVariable(ctx, route.Path)
			if len(prefix) > 0 && prefix != "." {
				routPath = "/" + path.Clean(prefix) + routPath
			}
			pathItem := spec2Path(ctx, group, route)
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

func spec2Path(ctx Context, group apiSpec.Group, route apiSpec.Route) spec.PathItem {
	authType := getStringFromKVOrDefault(group.Annotation.Properties, propertyKeyAuthType, "")
	var security []map[string][]string
	if len(authType) > 0 {
		security = []map[string][]string{
			{
				authType: []string{},
			},
		}
	}
	groupName := getStringFromKVOrDefault(group.Annotation.Properties, propertyKeyGroup, "")
	operationId := route.Handler
	if len(groupName) > 0 {
		operationId = stringx.From(groupName + "_" + route.Handler).ToCamel()
	}
	operationId = stringx.From(operationId).Untitle()
	op := &spec.Operation{
		OperationProps: spec.OperationProps{
			Description: getStringFromKVOrDefault(route.AtDoc.Properties, propertyKeyDescription, ""),
			Consumes:    consumesFromTypeOrDef(ctx, route.Method, route.RequestType),
			Produces:    getListFromInfoOrDefault(route.AtDoc.Properties, propertyKeyProduces, []string{applicationJson}),
			Schemes:     getListFromInfoOrDefault(route.AtDoc.Properties, propertyKeySchemes, []string{schemeHttps}),
			Tags:        getListFromInfoOrDefault(group.Annotation.Properties, propertyKeyTags, getListFromInfoOrDefault(group.Annotation.Properties, propertyKeySummary, []string{})),
			Summary:     getStringFromKVOrDefault(route.AtDoc.Properties, propertyKeySummary, getFirstUsableString(route.AtDoc.Text, route.Handler)),
			ID:          operationId,
			Deprecated:  getBoolFromKVOrDefault(route.AtDoc.Properties, propertyKeyDeprecated, false),
			Parameters:  parametersFromType(ctx, route.Method, route.RequestType),
			Security:    security,
			Responses:   jsonResponseFromType(ctx, route.AtDoc, route.ResponseType),
		},
	}
	externalDocsDescription := getStringFromKVOrDefault(route.AtDoc.Properties, propertyKeyExternalDocsDescription, "")
	externalDocsURL := getStringFromKVOrDefault(route.AtDoc.Properties, propertyKeyExternalDocsURL, "")
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
