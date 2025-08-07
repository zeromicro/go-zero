package swagger

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	apiSpec "github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func TestIsPostJson(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		hasJson  bool
		expected bool
	}{
		{"POST with JSON", http.MethodPost, true, true},
		{"POST without JSON", http.MethodPost, false, false},
		{"GET with JSON", http.MethodGet, true, false},
		{"PUT with JSON", http.MethodPut, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStruct := createTestStruct("TestStruct", tt.hasJson)
			_, result := isPostJson(testingContext(t), tt.method, testStruct)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParametersFromType(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		useDefinitions bool
		hasJson        bool
		expectedCount  int
		expectedBody   bool
	}{
		{"POST JSON with definitions", http.MethodPost, true, true, 1, true},
		{"POST JSON without definitions", http.MethodPost, false, true, 1, true},
		{"GET with form", http.MethodGet, false, false, 1, false},
		{"POST with form", http.MethodPost, false, false, 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := Context{UseDefinitions: tt.useDefinitions}
			testStruct := createTestStruct("TestStruct", tt.hasJson)
			params := parametersFromType(ctx, tt.method, testStruct)

			assert.Equal(t, tt.expectedCount, len(params))
			if tt.expectedBody {
				assert.Equal(t, paramsInBody, params[0].In)
			} else if len(params) > 0 {
				assert.NotEqual(t, paramsInBody, params[0].In)
			}
		})
	}
}

func TestParametersFromType_EdgeCases(t *testing.T) {
	ctx := testingContext(t)

	params := parametersFromType(ctx, http.MethodPost, nil)
	assert.Empty(t, params)

	primitiveType := apiSpec.PrimitiveType{RawName: "string"}
	params = parametersFromType(ctx, http.MethodPost, primitiveType)
	assert.Empty(t, params)
}

func createTestStruct(name string, hasJson bool) apiSpec.DefineStruct {
	tag := `form:"username"`
	if hasJson {
		tag = `json:"username"`
	}

	return apiSpec.DefineStruct{
		RawName: name,
		Members: []apiSpec.Member{
			{
				Name: "Username",
				Type: apiSpec.PrimitiveType{RawName: "string"},
				Tag:  tag,
			},
		},
	}
}
