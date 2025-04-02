package swagger

import (
	"github.com/go-openapi/spec"
	apiSpec "github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"net/http"
)

func jsonResponseFromType(tp apiSpec.Type, statusResponse map[int]map[string]any) *spec.Responses {
	p, _ := propertiesFromType(tp)
	resp := map[int]spec.Response{
		http.StatusOK: {
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
	}
	for httpStatus, example := range statusResponse {
		if httpStatus == http.StatusOK {
			continue
		}
		response:=spec.NewResponse().AddExample(applicationJson,example)
		resp[httpStatus] = *response
	}
	return &spec.Responses{
		ResponsesProps: spec.ResponsesProps{
			StatusCodeResponses: resp,
		},
	}
}
