package swagger

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	apiSpec "github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func TestParametersFromType_PathParameterValidation(t *testing.T) {
	tests := []struct {
		name                string
		routePath           string
		structType          apiSpec.DefineStruct
		expectedPathParams  []string
	}{
		{
			name:      "path parameter matches route placeholder",
			routePath: "/api/users/:id",
			structType: apiSpec.DefineStruct{
				RawName: "TestRequest",
				Members: []apiSpec.Member{
					{
						Name: "ID",
						Type: apiSpec.PrimitiveType{RawName: "string"},
						Tag:  `path:"id"`,
					},
				},
			},
			expectedPathParams: []string{"id"},
		},
		{
			name:      "path parameter matches {id} style placeholder",
			routePath: "/api/users/{id}",
			structType: apiSpec.DefineStruct{
				RawName: "TestRequest",
				Members: []apiSpec.Member{
					{
						Name: "ID",
						Type: apiSpec.PrimitiveType{RawName: "string"},
						Tag:  `path:"id"`,
					},
				},
			},
			expectedPathParams: []string{"id"},
		},
		{
			name:      "path parameter does NOT match route placeholder - should be filtered out",
			routePath: "/api/users",
			structType: apiSpec.DefineStruct{
				RawName: "TestRequest",
				Members: []apiSpec.Member{
					{
						Name: "ID",
						Type: apiSpec.PrimitiveType{RawName: "string"},
						Tag:  `path:"id"`,
					},
				},
			},
			expectedPathParams: []string{}, // No path params should be generated
		},
		{
			name:      "multiple path parameters, only some match route placeholders",
			routePath: "/api/:namespace/users/:id",
			structType: apiSpec.DefineStruct{
				RawName: "TestRequest",
				Members: []apiSpec.Member{
					{
						Name: "Namespace",
						Type: apiSpec.PrimitiveType{RawName: "string"},
						Tag:  `path:"namespace"`,
					},
					{
						Name: "ID",
						Type: apiSpec.PrimitiveType{RawName: "string"},
						Tag:  `path:"id"`,
					},
					{
						Name: "Extra",
						Type: apiSpec.PrimitiveType{RawName: "string"},
						Tag:  `path:"extra"`, // This should be filtered out
					},
				},
			},
			expectedPathParams: []string{"namespace", "id"},
		},
		{
			name:      "no path tag - no path params should be generated",
			routePath: "/api/users/:id",
			structType: apiSpec.DefineStruct{
				RawName: "TestRequest",
				Members: []apiSpec.Member{
					{
						Name: "ID",
						Type: apiSpec.PrimitiveType{RawName: "string"},
						Tag:  `form:"id"`, // Form param, not path param
					},
				},
			},
			expectedPathParams: []string{},
		},
		{
			name:      "mixed route path - colon and brace styles",
			routePath: "/api/:namespace/users/{id}",
			structType: apiSpec.DefineStruct{
				RawName: "TestRequest",
				Members: []apiSpec.Member{
					{
						Name: "Namespace",
						Type: apiSpec.PrimitiveType{RawName: "string"},
						Tag:  `path:"namespace"`,
					},
					{
						Name: "ID",
						Type: apiSpec.PrimitiveType{RawName: "string"},
						Tag:  `path:"id"`,
					},
				},
			},
			expectedPathParams: []string{"namespace", "id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := testingContext(t)

			params := parametersFromType(ctx, http.MethodGet, tt.structType, tt.routePath)

			// Filter only path parameters
			var pathParams []string
			for _, param := range params {
				if param.In == paramsInPath {
					pathParams = append(pathParams, param.Name)
				}
			}

			// Verify we have the expected number of path parameters
			assert.Equal(t, len(tt.expectedPathParams), len(pathParams),
				"Expected %d path parameters, got %d. Route: %s, Expected: %v, Got: %v",
				len(tt.expectedPathParams), len(pathParams), tt.routePath, tt.expectedPathParams, pathParams)

			// Verify all expected path parameters are present
			for _, expected := range tt.expectedPathParams {
				assert.Contains(t, pathParams, expected,
					"Expected path parameter '%s' to be present", expected)
			}
		})
	}
}

func TestParametersFromType_PathParameterFilteringWithQueryString(t *testing.T) {
	// Test query string parameters are NOT affected by path parameter filtering
	ctx := testingContext(t)

	structType := apiSpec.DefineStruct{
		RawName: "TestRequest",
		Members: []apiSpec.Member{
			{
				Name: "ID",
				Type: apiSpec.PrimitiveType{RawName: "string"},
				Tag:  `path:"id"`,
			},
			{
				Name: "Name",
				Type: apiSpec.PrimitiveType{RawName: "string"},
				Tag:  `form:"name"`, // Query param should still be generated
			},
		},
	}

	// Route has placeholder for id, so path param should be generated
	params := parametersFromType(ctx, http.MethodGet, structType, "/api/:id")

	var pathParams, queryParams []string
	for _, param := range params {
		if param.In == paramsInPath {
			pathParams = append(pathParams, param.Name)
		} else if param.In == paramsInQuery {
			queryParams = append(queryParams, param.Name)
		}
	}

	assert.Equal(t, []string{"id"}, pathParams)
	assert.Equal(t, []string{"name"}, queryParams)
}