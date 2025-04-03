package swagger

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/go-openapi/spec"
	apiSpec "github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/internal/version"
	"github.com/zeromicro/go-zero/tools/goctl/util"
)

func spec2Swagger(api *apiSpec.ApiSpec) (*spec.Swagger, error) {
	extensions, info := specExtensions(api.Info)
	swagger := &spec.Swagger{
		VendorExtensible: spec.VendorExtensible{
			Extensions: extensions,
		},
		SwaggerProps: spec.SwaggerProps{
			Consumes: getListFromInfoOrDefault(api.Info.Properties, "consumes", []string{applicationJson}),
			Produces: getListFromInfoOrDefault(api.Info.Properties, "produces", []string{applicationJson}),
			Schemes:  getListFromInfoOrDefault(api.Info.Properties, "schemes", []string{schemeHttps}),
			Swagger:  swaggerVersion,
			Info:     info,
			Host:     getStringFromKVOrDefault(api.Info.Properties, "host", defaultHost),
			BasePath: getStringFromKVOrDefault(api.Info.Properties, "basePath", defaultBasePath),
			Paths:    spec2Paths(api.Info, api.Service),
		},
	}

	return swagger, nil
}

func formatComment(comment string) string {
	s := strings.TrimPrefix(comment, "//")
	return strings.TrimSpace(s)
}

func sampleItemsFromGoType(tp apiSpec.Type) *spec.Items {
	val, ok := tp.(apiSpec.ArrayType)
	if !ok {
		return nil
	}
	item := val.Value
	switch item.(type) {
	case apiSpec.PrimitiveType:
		return &spec.Items{
			SimpleSchema: spec.SimpleSchema{
				Type: sampleTypeFromGoType(item),
			},
		}
	case apiSpec.ArrayType:
		return &spec.Items{
			SimpleSchema: spec.SimpleSchema{
				Type:  sampleTypeFromGoType(item),
				Items: sampleItemsFromGoType(item),
			},
		}
	default: // unsupported type
	}
	return nil
}

// itemsFromGoType returns the schema or array of the type, just for non json body parameters.
func itemsFromGoType(tp apiSpec.Type) *spec.SchemaOrArray {
	array, ok := tp.(apiSpec.ArrayType)
	if !ok {
		return nil
	}
	return itemFromGoType(array)
}

func mapFromGoType(tp apiSpec.Type) *spec.SchemaOrBool {
	mapType, ok := tp.(apiSpec.MapType)
	if !ok {
		return nil
	}
	p, r := propertiesFromType(mapType.Value)
	return &spec.SchemaOrBool{
		Allows: true,
		Schema: &spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type:                 typeFromGoType(mapType.Value),
				Items:                itemsFromGoType(mapType.Value),
				Properties:           p,
				Required:             r,
				AdditionalProperties: mapFromGoType(mapType.Value),
			},
		},
	}
}

// itemFromGoType returns the schema or array of the type, just for non json body parameters.
func itemFromGoType(tp apiSpec.Type) *spec.SchemaOrArray {
	switch itemType := tp.(type) {
	case apiSpec.PrimitiveType:
		return &spec.SchemaOrArray{
			Schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: typeFromGoType(tp),
				},
			},
		}
	case apiSpec.DefineStruct:
		var (
			properties     = map[string]spec.Schema{}
			requiredFields []string
		)
		rangeMemberAndDo(itemType, func(tag *apiSpec.Tags, required bool, member apiSpec.Member) {
			jsonTag, _ := tag.Get(tagJson)
			if jsonTag == nil {
				return
			}
			minimum, maximum, exclusiveMinimum, exclusiveMaximum := rangeValueFromOptions(jsonTag.Options)
			if required {
				requiredFields = append(requiredFields, jsonTag.Name)
			}
			p, r := propertiesFromType(member.Type)
			properties[jsonTag.Name] = spec.Schema{
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					Example: exampleValueFromOptions(jsonTag.Options, member.Type),
				},
				SchemaProps: spec.SchemaProps{
					Description:          formatComment(member.Comment),
					Type:                 typeFromGoType(member.Type),
					Default:              defValueFromOptions(jsonTag.Options, member.Type),
					Maximum:              maximum,
					ExclusiveMaximum:     exclusiveMaximum,
					Minimum:              minimum,
					ExclusiveMinimum:     exclusiveMinimum,
					Enum:                 enumsValueFromOptions(jsonTag.Options),
					Items:                itemsFromGoType(member.Type),
					Properties:           p,
					Required:             r,
					AdditionalProperties: mapFromGoType(member.Type),
				},
			}
		})
		return &spec.SchemaOrArray{
			Schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type:                 typeFromGoType(itemType),
					Items:                itemsFromGoType(itemType),
					Properties:           properties,
					Required:             requiredFields,
					AdditionalProperties: mapFromGoType(itemType),
				},
			},
		}
	case apiSpec.PointerType:
		return itemsFromGoType(itemType.Type)
	case apiSpec.ArrayType:
		p, r := propertiesFromType(itemType.Value)
		return &spec.SchemaOrArray{
			Schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type:       typeFromGoType(itemType.Value),
					Items:      itemsFromGoType(itemType.Value),
					Properties: p,
					Required:   r,
				},
			},
		}
	}
	return nil
}

func typeFromGoType(tp apiSpec.Type) []string {
	switch val := tp.(type) {
	case apiSpec.PrimitiveType:
		res, ok := tpMapper[val.RawName]
		if ok {
			return []string{res}
		}
	case apiSpec.ArrayType:
		return []string{swaggerTypeArray}
	case apiSpec.DefineStruct, apiSpec.MapType:
		return []string{swaggerTypeObject}
	case apiSpec.PointerType:
		return typeFromGoType(val.Type)
	}
	return nil
}

func sampleTypeFromGoType(tp apiSpec.Type) string {
	switch val := tp.(type) {
	case apiSpec.PrimitiveType:
		return tpMapper[val.RawName]
	case apiSpec.ArrayType:
		return swaggerTypeArray
	case apiSpec.DefineStruct, apiSpec.MapType:
		return swaggerTypeObject
	case apiSpec.PointerType:
		return sampleTypeFromGoType(val.Type)
	default:
		return ""
	}
}

func typeContainsTag(structType apiSpec.DefineStruct, tag string) bool {
	for _, field := range structType.Members {
		tags, _ := apiSpec.Parse(field.Tag)
		for _, t := range tags.Tags() {
			if t.Key == tag {
				return true
			}
		}
	}
	return false
}

func rangeMemberAndDo(structType apiSpec.DefineStruct, do func(tag *apiSpec.Tags, required bool, member apiSpec.Member)) {
	for _, field := range structType.Members {
		var required = false
		for _, t := range field.Tags() {
			required = len(t.Options) > 0 && t.Options[0] != "optional"
		}
		tags, _ := apiSpec.Parse(field.Tag)
		do(tags, required, field)

	}
}

func pathVariable2SwaggerVariable(path string) string {
	pathItems := strings.FieldsFunc(path, slashRune)
	var resp []string
	for _, v := range pathItems {
		if strings.HasPrefix(v, ":") {
			resp = append(resp, "{"+v[1:]+"}")
		} else {
			resp = append(resp, v)
		}
	}
	return "/" + filepath.Join(resp...)
}

func specExtensions(api apiSpec.Info) (spec.Extensions, *spec.Info) {
	ext := spec.Extensions{}
	ext.Add("x-goctl-version", version.BuildVersion)
	ext.Add("x-description", "This is a goctl generated swagger file.")
	ext.Add("x-date", time.Now().Format("2006-01-02 15:04:05"))
	ext.Add("x-github", "https://github.com/zeromicro/go-zero")
	ext.Add("x-go-zero-doc", "https://go-zero.dev/")

	info := &spec.Info{}
	info.Description = util.Unquote(api.Properties["description"])
	info.Title = util.Unquote(api.Properties["title"])
	info.TermsOfService = util.Unquote(api.Properties["termsOfService"])
	info.Version = util.Unquote(api.Properties["version"])

	contactInfo := spec.ContactInfo{}
	contactInfo.Name = util.Unquote(api.Properties["contactName"])
	contactInfo.URL = util.Unquote(api.Properties["contactURL"])
	contactInfo.Email = util.Unquote(api.Properties["contactEmail"])
	if len(contactInfo.Name) > 0 || len(contactInfo.URL) > 0 || len(contactInfo.Email) > 0 {
		info.Contact = &contactInfo
	}

	license := &spec.License{}
	license.Name = util.Unquote(api.Properties["licenseName"])
	license.URL = util.Unquote(api.Properties["licenseURL"])
	if len(license.Name) > 0 || len(license.URL) > 0 {
		info.License = license
	}
	return ext, info
}
