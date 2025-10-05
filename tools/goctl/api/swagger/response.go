package swagger

import (
	"net/http"

	"github.com/go-openapi/spec"
	apiSpec "github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func jsonResponseFromType(ctx Context, atDoc apiSpec.AtDoc, tp apiSpec.Type) *spec.Responses {
	if tp == nil {
		return &spec.Responses{
			ResponsesProps: spec.ResponsesProps{
				StatusCodeResponses: map[int]spec.Response{
					http.StatusOK: {
						ResponseProps: spec.ResponseProps{
							Description: "",
							Schema:      &spec.Schema{},
						},
					},
				},
			},
		}
	}
	props := spec.SchemaProps{
		AdditionalProperties: mapFromGoType(ctx, tp),
		Items:                itemsFromGoType(ctx, tp),
	}
	
	// Handle arrays with useDefinitions
	if arrayType, isArray := tp.(apiSpec.ArrayType); isArray && ctx.UseDefinitions {
		if structName, containsStruct := containsStruct(arrayType.Value); containsStruct {
			// For arrays, set $ref inside items, not at schema level
			props.Items = &spec.SchemaOrArray{
				Schema: &spec.Schema{
					SchemaProps: spec.SchemaProps{
						Ref: spec.MustCreateRef(getRefName(structName)),
					},
				},
			}
		}
	}
	
	if ctx.UseDefinitions {
		// For non-array types containing structs, use $ref at schema level
		if _, isArray := tp.(apiSpec.ArrayType); !isArray {
			structName, ok := containsStruct(tp)
			if ok {
				props.Ref = spec.MustCreateRef(getRefName(structName))
				return &spec.Responses{
					ResponsesProps: spec.ResponsesProps{
						StatusCodeResponses: map[int]spec.Response{
							http.StatusOK: {
								ResponseProps: spec.ResponseProps{
									Schema: &spec.Schema{
										SchemaProps: wrapCodeMsgProps(ctx, props, atDoc),
									},
								},
							},
						},
					},
				}
			}
		}
	}

	props.Type = typeFromGoType(ctx, tp)
	
	// For array types with useDefinitions, items are already set correctly above
	// For non-array types, we need to set properties
	if _, isArray := tp.(apiSpec.ArrayType); !isArray {
		p, _ := propertiesFromType(ctx, tp)
		props.Properties = p
	}
	
	return &spec.Responses{
		ResponsesProps: spec.ResponsesProps{
			StatusCodeResponses: map[int]spec.Response{
				http.StatusOK: {
					ResponseProps: spec.ResponseProps{
						Schema: &spec.Schema{
							SchemaProps: wrapCodeMsgProps(ctx, props, atDoc),
						},
					},
				},
			},
		},
	}
}
