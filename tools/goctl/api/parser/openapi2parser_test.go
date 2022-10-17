/*
Copyright 2022 Sangfor co.ltd.  All rights reserved.
*/

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func TestOpenapi2Parser(T *testing.T) {
	parser, err := newOpenApi2Parser("testdata/openapi2test.yaml")
	assert.Nil(T, err)
	api, err := parser.parse()
	assert.Nil(T, err)
	assert.NotNil(T, api)
	assert.NotZero(T, len(api.Service.Groups))
	assert.Equal(T, len(api.Service.Groups[0].Routes), 1)
	fieldNames := []string{"field1", "field2", "field3"}
	route := api.Service.Groups[0].Routes[0]
	assert.Equal(T, "get", route.Method)

	assert.Equal(T, route.Path, "/api/v1/test")
	assert.Equal(T, route.Handler, "getTestOperationId")
	assert.IsType(T, spec.DefineStruct{}, route.RequestType)
	assert.Equal(T, len(route.RequestType.(spec.DefineStruct).Members), 3)
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

}
