package swagger

import (
	"github.com/go-openapi/spec"
	apiSpec "github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func propertiesFromType(ctx Context, tp apiSpec.Type) (spec.SchemaProperties, []string) {
	var (
		properties     = map[string]spec.Schema{}
		requiredFields []string
	)
	switch val := tp.(type) {
	case apiSpec.PointerType:
		return propertiesFromType(ctx, val.Type)
	case apiSpec.ArrayType:
		return propertiesFromType(ctx, val.Value)
	case apiSpec.DefineStruct, apiSpec.NestedStruct:
		rangeMemberAndDo(ctx,val, func(tag *apiSpec.Tags, required bool, member apiSpec.Member) {
			var (
				jsonTagString                      = member.Name
				minimum, maximum                   *float64
				exclusiveMinimum, exclusiveMaximum bool
				example, defaultValue              any
				enum                               []any
			)
			jsonTag, _ := tag.Get(tagJson)
			if jsonTag != nil {
				jsonTagString = jsonTag.Name
				minimum, maximum, exclusiveMinimum, exclusiveMaximum = rangeValueFromOptions(jsonTag.Options)
				example = exampleValueFromOptions(ctx,jsonTag.Options, member.Type)
				defaultValue = defValueFromOptions(ctx,jsonTag.Options, member.Type)
				enum = enumsValueFromOptions(jsonTag.Options)
			}

			if required {
				requiredFields = append(requiredFields, jsonTagString)
			}
			var schema = spec.Schema{
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					Example: example,
				},
				SchemaProps: spec.SchemaProps{
					Description:          formatComment(member.Comment),
					Type:                 typeFromGoType(ctx,member.Type),
					Default:              defaultValue,
					Maximum:              maximum,
					ExclusiveMaximum:     exclusiveMaximum,
					Minimum:              minimum,
					ExclusiveMinimum:     exclusiveMinimum,
					Enum:                 enum,
					AdditionalProperties: mapFromGoType(ctx,member.Type),
				},
			}
			switch sampleTypeFromGoType(ctx,member.Type) {
			case swaggerTypeArray:
				schema.Items = itemsFromGoType(ctx,member.Type)
			case swaggerTypeObject:
				p, r := propertiesFromType(ctx,member.Type)
				schema.Properties = p
				schema.Required = r
			}

			properties[jsonTagString] = schema
		})
	}

	return properties, requiredFields
}
