package parser

import (
	"context"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

var testOpenapi = "testdata/openapi3test.yaml"

func TestOpenapi3Parser(T *testing.T) {
	parser, err := newOpenApi3Parser(context.Background(), testOpenapi)
	assert.Nil(T, err, "")
	api, err := parser.parse()
	assert.Nil(T, err)
	assert.NotNil(T, api)
	assert.NotZero(T, len(api.Service.Groups))
	assert.Equal(T, len(api.Service.Groups[0].Routes), 2)
	fieldNames := []string{"field1", "field2", "field3"}
	for _, route := range api.Service.Groups[0].Routes {
		if route.Method == "get" {
			assert.Equal(T, route.Path, "/api/v1/gettest/:field1")
			assert.Equal(T, route.Handler, "testGetOperationId")
			assert.IsType(T, spec.DefineStruct{}, route.RequestType)
			assert.Equal(T, len(route.RequestType.(spec.DefineStruct).Members), 4)
			for _, fieldName := range fieldNames {
				exist := false
				for _, member := range route.RequestType.(spec.DefineStruct).Members {
					if member.Name == fieldName {
						exist = true
						break
					}
				}
				assert.Equal(T, exist, true)
			}

			assert.IsType(T, spec.DefineStruct{}, route.ResponseType)
			for _, fieldName := range fieldNames {
				exist := false
				for _, member := range route.ResponseType.(spec.DefineStruct).Members {
					if member.Name == fieldName {
						exist = true
						break
					}
				}
				assert.Equal(T, exist, true)
			}
		} else {
			assert.Equal(T, route.Path, "/api/v1/posttest")
			assert.Equal(T, route.Handler, "postGetOperationId")
			assert.IsType(T, spec.DefineStruct{}, route.RequestType)
			assert.Equal(T, len(route.RequestType.(spec.DefineStruct).Members), 4)
			for _, fieldName := range fieldNames {
				exist := false
				for _, member := range route.RequestType.(spec.DefineStruct).Members {
					if member.Name == fieldName {
						exist = true
						break
					}
				}
				assert.Equal(T, exist, true)
			}
		}
	}
}
