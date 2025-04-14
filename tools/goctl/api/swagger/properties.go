package swagger

import (
	"github.com/go-openapi/spec"
	apiSpec "github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func propertiesFromType(tp apiSpec.Type) (spec.SchemaProperties, []string) {
	var (
		properties     = map[string]spec.Schema{}
		requiredFields []string
	)
	switch val := tp.(type) {
	case apiSpec.PointerType:
		return propertiesFromType(val.Type)
	case apiSpec.ArrayType:
		return propertiesFromType(val.Value)
	case apiSpec.DefineStruct,apiSpec.NestedStruct:
		rangeMemberAndDo(val, func(tag *apiSpec.Tags, required bool, member apiSpec.Member) {
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
	}

	return properties, requiredFields
}
