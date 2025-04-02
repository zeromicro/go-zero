package swagger

import (
	"github.com/go-openapi/spec"
	apiSpec "github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func jsonResponseFromType(tp apiSpec.Type) *spec.Responses {
	return &spec.Responses{
		ResponsesProps: spec.ResponsesProps{
			Default: &spec.Response{
				ResponseProps: spec.ResponseProps{
					Description: "",
					Schema:      &spec.Schema{},
					Examples: map[string]any{
						applicationJson: `{"example":true}`,
					},
				},
			},
		},
	}
}
