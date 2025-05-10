package swagger

import (
	"github.com/go-openapi/spec"
	apiSpec "github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func jsonResponseFromType(ctx Context, atDoc apiSpec.AtDoc, tp apiSpec.Type) *spec.Responses {
	if tp == nil {
		return nil
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
					Default: &spec.Response{
						ResponseProps: spec.ResponseProps{
							Schema: &spec.Schema{
								SchemaProps: wrapCodeMsgProps(ctx, props, atDoc),
							},
						},
					},
				},
			}
		}
	}

	p, _ := propertiesFromType(ctx, tp)
	props.Type = typeFromGoType(ctx, tp)
	props.Properties = p
	return &spec.Responses{
		ResponsesProps: spec.ResponsesProps{
			Default: &spec.Response{
				ResponseProps: spec.ResponseProps{
					Schema: &spec.Schema{
						SchemaProps: wrapCodeMsgProps(ctx, props, atDoc),
					},
				},
			},
		},
	}
}
