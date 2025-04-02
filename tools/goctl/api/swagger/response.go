package swagger

import (
	"github.com/go-openapi/spec"
	apiSpec "github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func jsonResponseFromType(tp apiSpec.Type) *spec.Responses {
	p, _ := propertiesFromType(tp)
	return &spec.Responses{
		ResponsesProps: spec.ResponsesProps{
			Default: &spec.Response{
				ResponseProps: spec.ResponseProps{
					Schema: &spec.Schema{
						SchemaProps: spec.SchemaProps{
							Type:                 typeFromGoType(tp),
							Properties:           p,
							AdditionalProperties: mapFromGoType(tp),
							Items:                itemsFromGoType(tp),
						},
					},
				},
			},
		},
	}
}
