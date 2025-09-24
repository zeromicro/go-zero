package swagger

import (
	"net/http"

	apiSpec "github.com/dropship-dev/go-zero/tools/goctl/api/spec"
	"github.com/go-openapi/spec"
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
	if ctx.UseDefinitions {
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

	p, r := propertiesFromType(ctx, tp)
	props.Required = r
	props.Type = typeFromGoType(ctx, tp)
	props.Properties = p
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
