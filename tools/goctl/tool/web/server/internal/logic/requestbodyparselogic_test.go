package logic

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/types"
)

func Test_parseJSON(t *testing.T) {
	var testData = []struct {
		name     string
		input    string
		expected []*types.FormItem
		err      bool
	}{
		{
			name:     "empty",
			input:    ``,
			expected: []*types.FormItem{},
		},
		{
			name:  "invalidJson",
			input: `{]`,
			err:   true,
		},
		{
			name:  "invalidEOF",
			input: `{`,
			err:   true,
		},
		{
			name:  "invalidTypeSlice",
			input: `[]`,
			err:   true,
		},
		{
			name:  "invalidTypeStruct",
			input: `{"foo":{}}`,
			err:   true,
		},
		{
			name:  "invalidElemType",
			input: `{"foo":[{}]}`,
			err:   true,
		},
		{
			name:  "invalidElemSliceType",
			input: `{"foo":[[]]}`,
			err:   true,
		},
		{
			name:  "success",
			input: `{"id":1,"name":"foo","age":18,"score":99.5,"active":true,"intList":[1,2,3],"anyList":[],"boolList":[true,false],"floatList":[1.1,1.2],"stringList":["a","b"]}`,
			expected: []*types.FormItem{
				{
					Name: "id",
					Type: "int64",
				},
				{
					Name: "name",
					Type: "string",
				},
				{
					Name: "age",
					Type: "int64",
				},
				{
					Name: "score",
					Type: "float64",
				},
				{
					Name: "active",
					Type: "bool",
				},
				{
					Name: "intList",
					Type: "[]int64",
				},
				{
					Name: "floatList",
					Type: "[]float64",
				},
				{
					Name: "anyList",
					Type: "[]interface{}",
				},
				{
					Name: "boolList",
					Type: "[]bool",
				},
				{
					Name: "stringList",
					Type: "[]string",
				},
			},
		},
	}

	for _, data := range testData {
		t.Run(data.name, func(t *testing.T) {
			got, err := parseJSON(data.input)
			if data.err {
				assert.Error(t, err)
				fmt.Println(err.Error())
				return
			}
			assert.NoError(t, err)
			sort.SliceStable(data.expected, func(i, j int) bool {
				return data.expected[i].Name < data.expected[j].Name
			})
			assert.Equal(t, data.expected, got)
		})
	}
}
