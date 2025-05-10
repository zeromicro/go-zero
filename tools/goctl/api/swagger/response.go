package swagger

import (
	"github.com/go-openapi/spec"
	apiSpec "github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func jsonResponseFromType(ctx Context, atDoc apiSpec.AtDoc, tp apiSpec.Type) *spec.Responses {
	p, _ := propertiesFromType(ctx, tp)
	props := spec.SchemaProps{
		Type:                 typeFromGoType(ctx, tp),
		Properties:           p,
		AdditionalProperties: mapFromGoType(ctx, tp),
		Items:                itemsFromGoType(ctx, tp),
	}

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
