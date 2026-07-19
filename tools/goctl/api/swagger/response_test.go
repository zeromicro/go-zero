package swagger

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func TestJsonResponseFromTypeStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		properties map[string]string
		response   spec.Type
		want       int
	}{
		{
			name:     "defaults to ok",
			response: spec.PrimitiveType{RawName: "string"},
			want:     http.StatusOK,
		},
		{
			name: "uses custom status code",
			properties: map[string]string{
				propertyKeyRespCode: "201",
			},
			response: spec.PrimitiveType{RawName: "string"},
			want:     http.StatusCreated,
		},
		{
			name: "supports quoted custom status code",
			properties: map[string]string{
				propertyKeyRespCode: `"204"`,
			},
			want: http.StatusNoContent,
		},
		{
			name: "defaults for invalid status code",
			properties: map[string]string{
				propertyKeyRespCode: "600",
			},
			response: spec.PrimitiveType{RawName: "string"},
			want:     http.StatusOK,
		},
		{
			name: "defaults for non-numeric status code",
			properties: map[string]string{
				propertyKeyRespCode: "created",
			},
			response: spec.PrimitiveType{RawName: "string"},
			want:     http.StatusOK,
		},
		{
			name: "uses custom status code without response body",
			properties: map[string]string{
				propertyKeyRespCode: "204",
			},
			want: http.StatusNoContent,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			responses := jsonResponseFromType(testingContext(t), spec.AtDoc{
				Properties: test.properties,
			}, test.response)

			assert.Len(t, responses.StatusCodeResponses, 1)
			assert.Contains(t, responses.StatusCodeResponses, test.want)
			if test.response == nil {
				assert.Nil(t, responses.StatusCodeResponses[test.want].Schema)
			}
		})
	}
}

func TestJsonResponseFromTypeMultipleStatusCodes(t *testing.T) {
	responses := jsonResponseFromType(testingContext(t), spec.AtDoc{
		Properties: map[string]string{
			propertyKeyResponses: "200-OK<br>401-Unauthorized<br>404-User not found",
		},
	}, spec.PrimitiveType{RawName: "string"})

	assert.Len(t, responses.StatusCodeResponses, 3)
	assert.Equal(t, "OK", responses.StatusCodeResponses[http.StatusOK].Description)
	assert.NotNil(t, responses.StatusCodeResponses[http.StatusOK].Schema)
	assert.Equal(t, "Unauthorized", responses.StatusCodeResponses[http.StatusUnauthorized].Description)
	assert.Nil(t, responses.StatusCodeResponses[http.StatusUnauthorized].Schema)
	assert.Equal(t, "User not found", responses.StatusCodeResponses[http.StatusNotFound].Description)
	assert.Nil(t, responses.StatusCodeResponses[http.StatusNotFound].Schema)
}
