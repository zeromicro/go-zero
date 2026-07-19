package swagger

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-openapi/spec"
	apiSpec "github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func jsonResponseFromType(ctx Context, atDoc apiSpec.AtDoc, tp apiSpec.Type) *spec.Responses {
	statusCode := responseStatusCode(atDoc)
	var response spec.Response
	if tp == nil {
		response = spec.Response{
			ResponseProps: spec.ResponseProps{
				Description: "",
			},
		}
	} else {
		props := spec.SchemaProps{
			AdditionalProperties: mapFromGoType(ctx, tp),
			Items:                itemsFromGoType(ctx, tp),
		}
		if ctx.UseDefinitions {
			structName, ok := containsStruct(tp)
			if ok {
				props.Ref = spec.MustCreateRef(getRefName(structName))
				response = spec.Response{
					ResponseProps: spec.ResponseProps{
						Schema: &spec.Schema{
							SchemaProps: wrapCodeMsgProps(ctx, props, atDoc),
						},
					},
				}
				return responsesFromStatusCode(atDoc, statusCode, response)
			}
		}

		p, _ := propertiesFromType(ctx, tp)
		props.Type = typeFromGoType(ctx, tp)
		props.Properties = p
		response = spec.Response{
			ResponseProps: spec.ResponseProps{
				Schema: &spec.Schema{
					SchemaProps: wrapCodeMsgProps(ctx, props, atDoc),
				},
			},
		}
	}

	return responsesFromStatusCode(atDoc, statusCode, response)
}

func responsesFromStatusCode(atDoc apiSpec.AtDoc, statusCode int, response spec.Response) *spec.Responses {
	statusCodeResponses := map[int]spec.Response{
		statusCode: response,
	}
	for code, description := range responseDescriptions(atDoc) {
		if code == statusCode {
			responseWithDescription := response
			responseWithDescription.Description = description
			statusCodeResponses[code] = responseWithDescription
		} else {
			statusCodeResponses[code] = spec.Response{
				ResponseProps: spec.ResponseProps{
					Description: description,
				},
			}
		}
	}

	return &spec.Responses{
		ResponsesProps: spec.ResponsesProps{
			StatusCodeResponses: statusCodeResponses,
		},
	}
}

func responseStatusCode(atDoc apiSpec.AtDoc) int {
	return getOrDefault(atDoc.Properties, propertyKeyRespCode, http.StatusOK, func(str string, def int) int {
		statusCode, err := strconv.Atoi(str)
		if err != nil || statusCode < http.StatusContinue || statusCode > 599 {
			return def
		}

		return statusCode
	})
}

func responseDescriptions(atDoc apiSpec.AtDoc) map[int]string {
	descriptions := make(map[int]string)
	for _, item := range strings.Split(getStringFromKVOrDefault(atDoc.Properties, propertyKeyResponses, ""), "<br>") {
		codeText, description, ok := strings.Cut(item, "-")
		if !ok {
			continue
		}

		code, err := strconv.Atoi(strings.TrimSpace(codeText))
		if err != nil || code < http.StatusContinue || code > 599 {
			continue
		}

		descriptions[code] = strings.TrimSpace(description)
	}

	return descriptions
}
