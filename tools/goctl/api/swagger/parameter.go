package swagger

import (
	"net/http"
	"strings"

	"github.com/go-openapi/spec"
	apiSpec "github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func parametersFromType(method string, tp apiSpec.Type) []spec.Parameter {
	if tp == nil {
		return []spec.Parameter{}
	}
	structType, ok := tp.(apiSpec.DefineStruct)
	if !ok {
		return []spec.Parameter{}
	}
	var (
		resp           []spec.Parameter
		properties     = map[string]spec.Schema{}
		requiredFields []string
	)
	rangeMemberAndDo(structType, func(tag *apiSpec.Tags, required bool, member apiSpec.Member) {
		headerTag, _ := tag.Get(tagHeader)
		hasHeader := headerTag != nil

		pathParameterTag, _ := tag.Get(tagPath)
		hasPathParameter := pathParameterTag != nil

		formTag, _ := tag.Get(tagForm)
		hasForm := formTag != nil

		jsonTag, _ := tag.Get(tagJson)
		hasJson := jsonTag != nil
		if hasHeader {
			minimum, maximum, exclusiveMinimum, exclusiveMaximum := rangeValueFromOptions(headerTag.Options)
			resp = append(resp, spec.Parameter{
				CommonValidations: spec.CommonValidations{
					Maximum:          maximum,
					ExclusiveMaximum: exclusiveMaximum,
					Minimum:          minimum,
					ExclusiveMinimum: exclusiveMinimum,
					Enum:             enumsValueFromOptions(headerTag.Options),
				},
				SimpleSchema: spec.SimpleSchema{
					Type:    sampleTypeFromGoType(member.Type),
					Default: defValueFromOptions(headerTag.Options, member.Type),
					Example: exampleValueFromOptions(headerTag.Options, member.Type),
					Items:   sampleItemsFromGoType(member.Type),
				},
				ParamProps: spec.ParamProps{
					In:          paramsInHeader,
					Name:        headerTag.Name,
					Description: formatComment(member.Comment),
					Required:    required,
				},
			})
		}
		if hasPathParameter {
			minimum, maximum, exclusiveMinimum, exclusiveMaximum := rangeValueFromOptions(pathParameterTag.Options)
			resp = append(resp, spec.Parameter{
				CommonValidations: spec.CommonValidations{
					Maximum:          maximum,
					ExclusiveMaximum: exclusiveMaximum,
					Minimum:          minimum,
					ExclusiveMinimum: exclusiveMinimum,
					Enum:             enumsValueFromOptions(pathParameterTag.Options),
				},
				SimpleSchema: spec.SimpleSchema{
					Type:    sampleTypeFromGoType(member.Type),
					Default: defValueFromOptions(pathParameterTag.Options, member.Type),
					Example: exampleValueFromOptions(pathParameterTag.Options, member.Type),
					Items:   sampleItemsFromGoType(member.Type),
				},
				ParamProps: spec.ParamProps{
					In:          paramsInPath,
					Name:        pathParameterTag.Name,
					Description: formatComment(member.Comment),
					Required:    required,
				},
			})
		}
		if hasForm {
			minimum, maximum, exclusiveMinimum, exclusiveMaximum := rangeValueFromOptions(formTag.Options)
			if strings.EqualFold(method, http.MethodGet) {
				resp = append(resp, spec.Parameter{
					CommonValidations: spec.CommonValidations{
						Maximum:          maximum,
						ExclusiveMaximum: exclusiveMaximum,
						Minimum:          minimum,
						ExclusiveMinimum: exclusiveMinimum,
						Enum:             enumsValueFromOptions(formTag.Options),
					},
					SimpleSchema: spec.SimpleSchema{
						Type:    sampleTypeFromGoType(member.Type),
						Default: defValueFromOptions(formTag.Options, member.Type),
						Example: exampleValueFromOptions(formTag.Options, member.Type),
						Items:   sampleItemsFromGoType(member.Type),
					},
					ParamProps: spec.ParamProps{
						In:              paramsInQuery,
						Name:            formTag.Name,
						Description:     formatComment(member.Comment),
						Required:        required,
						AllowEmptyValue: !required,
					},
				})
			} else {
				resp = append(resp, spec.Parameter{
					CommonValidations: spec.CommonValidations{
						Maximum:          maximum,
						ExclusiveMaximum: exclusiveMaximum,
						Minimum:          minimum,
						ExclusiveMinimum: exclusiveMinimum,
						Enum:             enumsValueFromOptions(formTag.Options),
					},
					SimpleSchema: spec.SimpleSchema{
						Type:    sampleTypeFromGoType(member.Type),
						Default: defValueFromOptions(formTag.Options, member.Type),
						Example: exampleValueFromOptions(formTag.Options, member.Type),
						Items:   sampleItemsFromGoType(member.Type),
					},
					ParamProps: spec.ParamProps{
						In:              paramsInForm,
						Name:            formTag.Name,
						Description:     formatComment(member.Comment),
						Required:        required,
						AllowEmptyValue: !required,
					},
				})
			}

		}
		if hasJson {
			minimum, maximum, exclusiveMinimum, exclusiveMaximum := rangeValueFromOptions(jsonTag.Options)
			if required {
				requiredFields = append(requiredFields, jsonTag.Name)
			}
			var schema = spec.Schema{
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
					AdditionalProperties: mapFromGoType(member.Type),
				},
			}
			switch sampleTypeFromGoType(member.Type) {
			case swaggerTypeArray:
				schema.Items = itemsFromGoType(member.Type)
			case swaggerTypeObject:
				p, r := propertiesFromType(member.Type)
				schema.Properties = p
				schema.Required = r
			}
			properties[jsonTag.Name] = schema
		}
	})
	if len(properties) > 0 {
		resp = append(resp, spec.Parameter{
			ParamProps: spec.ParamProps{
				In:       paramsInBody,
				Name:     paramsInBody,
				Required: true,
				Schema: &spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type:       typeFromGoType(structType),
						Properties: properties,
						Required:   requiredFields,
					},
				},
			},
		})
	}
	return resp
}
