package swagger

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/go-openapi/spec"
	apiSpec "github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/internal/version"
)

func spec2Swagger(api *apiSpec.ApiSpec) (*spec.Swagger, error) {
	ctx := contextFromApi(api.Info)
	extensions, info := specExtensions(api.Info)
	var securityDefinitions spec.SecurityDefinitions
	securityDefinitionsFromJson := getStringFromKVOrDefault(api.Info.Properties, "securityDefinitionsFromJson", `{}`)
	_ = json.Unmarshal([]byte(securityDefinitionsFromJson), &securityDefinitions)
	swagger := &spec.Swagger{
		VendorExtensible: spec.VendorExtensible{
			Extensions: extensions,
		},
		SwaggerProps: spec.SwaggerProps{
			Definitions:         definitionsFromTypes(ctx, api.Types),
			Consumes:            getListFromInfoOrDefault(api.Info.Properties, propertyKeyConsumes, []string{applicationJson}),
			Produces:            getListFromInfoOrDefault(api.Info.Properties, propertyKeyProduces, []string{applicationJson}),
			Schemes:             getListFromInfoOrDefault(api.Info.Properties, propertyKeySchemes, []string{schemeHttps}),
			Swagger:             swaggerVersion,
			Info:                info,
			Host:                getStringFromKVOrDefault(api.Info.Properties, propertyKeyHost, ""),
			BasePath:            getStringFromKVOrDefault(api.Info.Properties, propertyKeyBasePath, defaultBasePath),
			Paths:               spec2Paths(ctx, api.Service),
			SecurityDefinitions: securityDefinitions,
		},
	}

	return swagger, nil
}

func formatComment(comment string) string {
	s := strings.TrimPrefix(comment, "//")
	return strings.TrimSpace(s)
}

func sampleItemsFromGoType(ctx Context, tp apiSpec.Type) *spec.Items {
	val, ok := tp.(apiSpec.ArrayType)
	if !ok {
		return nil
	}
	item := val.Value
	switch item.(type) {
	case apiSpec.PrimitiveType:
		return &spec.Items{
			SimpleSchema: spec.SimpleSchema{
				Type: sampleTypeFromGoType(ctx, item),
			},
		}
	case apiSpec.ArrayType:
		return &spec.Items{
			SimpleSchema: spec.SimpleSchema{
				Type:  sampleTypeFromGoType(ctx, item),
				Items: sampleItemsFromGoType(ctx, item),
			},
		}
	default: // unsupported type
	}
	return nil
}

// itemsFromGoType returns the schema or array of the type, just for non json body parameters.
func itemsFromGoType(ctx Context, tp apiSpec.Type) *spec.SchemaOrArray {
	array, ok := tp.(apiSpec.ArrayType)
	if !ok {
		return nil
	}
	return itemFromGoType(ctx, array.Value)
}

func mapFromGoType(ctx Context, tp apiSpec.Type) *spec.SchemaOrBool {
	mapType, ok := tp.(apiSpec.MapType)
	if !ok {
		return nil
	}
	var schema = &spec.Schema{
		SchemaProps: spec.SchemaProps{
			Type:                 typeFromGoType(ctx, mapType.Value),
			AdditionalProperties: mapFromGoType(ctx, mapType.Value),
		},
	}
	switch sampleTypeFromGoType(ctx, mapType.Value) {
	case swaggerTypeArray:
		schema.Items = itemsFromGoType(ctx, mapType.Value)
	case swaggerTypeObject:
		p, r := propertiesFromType(ctx, mapType.Value)
		schema.Properties = p
		schema.Required = r
	}
	return &spec.SchemaOrBool{
		Allows: true,
		Schema: schema,
	}
}

// itemFromGoType returns the schema or array of the type, just for non json body parameters.
func itemFromGoType(ctx Context, tp apiSpec.Type) *spec.SchemaOrArray {
	switch itemType := tp.(type) {
	case apiSpec.PrimitiveType:
		return &spec.SchemaOrArray{
			Schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: typeFromGoType(ctx, tp),
				},
			},
		}
	case apiSpec.DefineStruct, apiSpec.NestedStruct, apiSpec.MapType:
		properties, requiredFields := propertiesFromType(ctx, itemType)
		return &spec.SchemaOrArray{
			Schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type:                 typeFromGoType(ctx, itemType),
					Items:                itemsFromGoType(ctx, itemType),
					Properties:           properties,
					Required:             requiredFields,
					AdditionalProperties: mapFromGoType(ctx, itemType),
				},
			},
		}
	case apiSpec.PointerType:
		return itemFromGoType(ctx, itemType.Type)
	case apiSpec.ArrayType:
		return &spec.SchemaOrArray{
			Schema: &spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type:  typeFromGoType(ctx, itemType),
					Items: itemsFromGoType(ctx, itemType),
				},
			},
		}
	}
	return nil
}

func typeFromGoType(ctx Context, tp apiSpec.Type) []string {
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
		return typeFromGoType(ctx, val.Type)
	}
	return nil
}

func sampleTypeFromGoType(ctx Context, tp apiSpec.Type) string {
	switch val := tp.(type) {
	case apiSpec.PrimitiveType:
		return tpMapper[val.RawName]
	case apiSpec.ArrayType:
		return swaggerTypeArray
	case apiSpec.DefineStruct, apiSpec.MapType, apiSpec.NestedStruct:
		return swaggerTypeObject
	case apiSpec.PointerType:
		return sampleTypeFromGoType(ctx, val.Type)
	default:
		return ""
	}
}

func typeContainsTag(_ Context, structType apiSpec.DefineStruct, tag string) bool {
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

func expandMembers(ctx Context, tp apiSpec.Type) []apiSpec.Member {
	var members []apiSpec.Member
	switch val := tp.(type) {
	case apiSpec.DefineStruct:
		for _, v := range val.Members {
			if v.IsInline {
				members = append(members, expandMembers(ctx, v.Type)...)
				continue
			}
			members = append(members, v)
		}
	case apiSpec.NestedStruct:
		for _, v := range val.Members {
			if v.IsInline {
				members = append(members, expandMembers(ctx, v.Type)...)
				continue
			}
			members = append(members, v)
		}
	}

	return members
}

func rangeMemberAndDo(ctx Context, structType apiSpec.Type, do func(tag *apiSpec.Tags, required bool, member apiSpec.Member)) {
	var members = expandMembers(ctx, structType)

	for _, field := range members {
		tags, _ := apiSpec.Parse(field.Tag)
		required := isRequired(ctx, tags)
		do(tags, required, field)
	}
}

func isRequired(ctx Context, tags *apiSpec.Tags) bool {
	tag, err := tags.Get(tagJson)
	if err == nil {
		return !isOptional(ctx, tag.Options)
	}
	tag, err = tags.Get(tagForm)
	if err == nil {
		return !isOptional(ctx, tag.Options)
	}
	tag, err = tags.Get(tagPath)
	if err == nil {
		return !isOptional(ctx, tag.Options)
	}
	return false
}

func isOptional(_ Context, options []string) bool {
	for _, option := range options {
		if option == optionalFlag {
			return true
		}
	}
	return false
}

func pathVariable2SwaggerVariable(_ Context, path string) string {
	pathItems := strings.FieldsFunc(path, slashRune)
	var resp []string
	for _, v := range pathItems {
		if strings.HasPrefix(v, ":") {
			resp = append(resp, "{"+v[1:]+"}")
		} else {
			resp = append(resp, v)
		}
	}
	return "/" + strings.Join(resp, "/")
}

func wrapCodeMsgProps(ctx Context, properties spec.SchemaProps, atDoc apiSpec.AtDoc) spec.SchemaProps {
	if !ctx.WrapCodeMsg {
		return properties
	}
	globalCodeDesc := ctx.BizCodeEnumDescription
	methodCodeDesc := getStringFromKVOrDefault(atDoc.Properties, propertyKeyBizCodeEnumDescription, globalCodeDesc)
	return spec.SchemaProps{
		Type: []string{swaggerTypeObject},
		Properties: spec.SchemaProperties{
			"code": {
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					Example: 0,
				},
				SchemaProps: spec.SchemaProps{
					Type:        []string{swaggerTypeInteger},
					Description: methodCodeDesc,
				},
			},
			"msg": {
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					Example: "ok",
				},
				SchemaProps: spec.SchemaProps{
					Type:        []string{swaggerTypeString},
					Description: "business message",
				},
			},
			"data": {
				SchemaProps: properties,
			},
		},
	}
}

func specExtensions(api apiSpec.Info) (spec.Extensions, *spec.Info) {
	ext := spec.Extensions{}
	ext.Add("x-goctl-version", version.BuildVersion)
	ext.Add("x-description", "This is a goctl generated swagger file.")
	ext.Add("x-date", time.Now().Format(time.DateTime))
	ext.Add("x-github", "https://github.com/zeromicro/go-zero")
	ext.Add("x-go-zero-doc", "https://go-zero.dev/")

	info := &spec.Info{}
	info.Title = getStringFromKVOrDefault(api.Properties, propertyKeyTitle, "")
	info.Description = getStringFromKVOrDefault(api.Properties, propertyKeyDescription, "")
	info.TermsOfService = getStringFromKVOrDefault(api.Properties, propertyKeyTermsOfService, "")
	info.Version = getStringFromKVOrDefault(api.Properties, propertyKeyVersion, "1.0")

	contactInfo := spec.ContactInfo{}
	contactInfo.Name = getStringFromKVOrDefault(api.Properties, propertyKeyContactName, "")
	contactInfo.URL = getStringFromKVOrDefault(api.Properties, propertyKeyContactURL, "")
	contactInfo.Email = getStringFromKVOrDefault(api.Properties, propertyKeyContactEmail, "")
	if len(contactInfo.Name) > 0 || len(contactInfo.URL) > 0 || len(contactInfo.Email) > 0 {
		info.Contact = &contactInfo
	}

	license := &spec.License{}
	license.Name = getStringFromKVOrDefault(api.Properties, propertyKeyLicenseName, "")
	license.URL = getStringFromKVOrDefault(api.Properties, propertyKeyLicenseURL, "")
	if len(license.Name) > 0 || len(license.URL) > 0 {
		info.License = license
	}
	return ext, info
}
