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
		rangeMemberAndDo(ctx, val, func(tag *apiSpec.Tags, required bool, member apiSpec.Member) {
			var (
				jsonTagString                      = member.Name
				minimum, maximum                   *float64
				exclusiveMinimum, exclusiveMaximum bool
				example, defaultValue              any
				enum                               []any
			)
			pathTag, _ := tag.Get(tagPath)
			if pathTag != nil {
				return
			}
			formTag, _ := tag.Get(tagForm)
			if formTag != nil {
				return
			}
			headerTag, _ := tag.Get(tagHeader)
			if headerTag != nil {
				return
			}

			jsonTag, _ := tag.Get(tagJson)
			if jsonTag != nil {
				jsonTagString = jsonTag.Name
				minimum, maximum, exclusiveMinimum, exclusiveMaximum = rangeValueFromOptions(jsonTag.Options)
				example = exampleValueFromOptions(ctx, jsonTag.Options, member.Type)
				defaultValue = defValueFromOptions(ctx, jsonTag.Options, member.Type)
				enum = enumsValueFromOptions(jsonTag.Options)
			}

			if required {
				requiredFields = append(requiredFields, jsonTagString)
			}

			if ctx.UseDefinitions {
				schema := buildSchemaWithDefinitions(ctx, member.Type, example, member.Comment, minimum, maximum, exclusiveMinimum, exclusiveMaximum, defaultValue, enum)
				if schema != nil {
					properties[jsonTagString] = *schema
					return
				}
			}

			schema := spec.Schema{
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					Example: example,
				},
				SchemaProps: spec.SchemaProps{
					Description:          formatComment(member.Comment),
					Type:                 typeFromGoType(ctx, member.Type),
					Default:              defaultValue,
					Maximum:              maximum,
					ExclusiveMaximum:     exclusiveMaximum,
					Minimum:              minimum,
					ExclusiveMinimum:     exclusiveMinimum,
					Enum:                 enum,
					AdditionalProperties: mapFromGoType(ctx, member.Type),
				},
			}

			switch sampleTypeFromGoType(ctx, member.Type) {
			case swaggerTypeArray:
				schema.Items = itemsFromGoType(ctx, member.Type)
			case swaggerTypeObject:
				p, r := propertiesFromType(ctx, member.Type)
				schema.Properties = p
				schema.Required = r
			}

			properties[jsonTagString] = schema
		})
	}

	return properties, requiredFields
}

// buildSchemaWithDefinitions Correctly handle $ref references for composite types
func buildSchemaWithDefinitions(ctx Context, tp apiSpec.Type, example any, comment string, minimum, maximum *float64, exclusiveMinimum, exclusiveMaximum bool, defaultValue any, enum []any) *spec.Schema {
	switch val := tp.(type) {
	case apiSpec.DefineStruct:
		return &spec.Schema{
			SwaggerSchemaProps: spec.SwaggerSchemaProps{
				Example: example,
			},
			SchemaProps: spec.SchemaProps{
				Description: formatComment(comment),
				Ref:         spec.MustCreateRef(getRefName(val.RawName)),
			},
		}
	case apiSpec.PointerType:
		return buildSchemaWithDefinitions(ctx, val.Type, example, comment, minimum, maximum, exclusiveMinimum, exclusiveMaximum, defaultValue, enum)
	case apiSpec.ArrayType:
		itemSchema := buildSchemaWithDefinitions(ctx, val.Value, nil, "", nil, nil, false, false, nil, nil)
		if itemSchema != nil {
			return &spec.Schema{
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					Example: example,
				},
				SchemaProps: spec.SchemaProps{
					Description: formatComment(comment),
					Type:        spec.StringOrArray{"array"},
					Items: &spec.SchemaOrArray{
						Schema: itemSchema,
					},
				},
			}
		}

		structName, containsStruct := containsStruct(val.Value)
		if containsStruct {
			return &spec.Schema{
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					Example: example,
				},
				SchemaProps: spec.SchemaProps{
					Description: formatComment(comment),
					Type:        spec.StringOrArray{"array"},
					Items: &spec.SchemaOrArray{
						Schema: &spec.Schema{
							SchemaProps: spec.SchemaProps{
								Ref: spec.MustCreateRef(getRefName(structName)),
							},
						},
					},
				},
			}
		}
	case apiSpec.MapType:
		valueSchema := buildSchemaWithDefinitions(ctx, val.Value, nil, "", nil, nil, false, false, nil, nil)
		if valueSchema != nil {
			return &spec.Schema{
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					Example: example,
				},
				SchemaProps: spec.SchemaProps{
					Description: formatComment(comment),
					Type:        spec.StringOrArray{"object"},
					AdditionalProperties: &spec.SchemaOrBool{
						Schema: valueSchema,
					},
				},
			}
		}
		structName, containsStruct := containsStruct(val.Value)
		if containsStruct {
			return &spec.Schema{
				SwaggerSchemaProps: spec.SwaggerSchemaProps{
					Example: example,
				},
				SchemaProps: spec.SchemaProps{
					Description: formatComment(comment),
					Type:        spec.StringOrArray{"object"},
					AdditionalProperties: &spec.SchemaOrBool{
						Schema: &spec.Schema{
							SchemaProps: spec.SchemaProps{
								Ref: spec.MustCreateRef(getRefName(structName)),
							},
						},
					},
				},
			}
		}
	}

	return nil
}

func containsStruct(tp apiSpec.Type) (string, bool) {
	switch val := tp.(type) {
	case apiSpec.PointerType:
		return containsStruct(val.Type)
	case apiSpec.ArrayType:
		return containsStruct(val.Value)
	case apiSpec.DefineStruct:
		return val.RawName, true
	case apiSpec.MapType:
		return containsStruct(val.Value)
	default:
		return "", false
	}
}

func getRefName(typeName string) string {
	return "#/definitions/" + typeName
}
