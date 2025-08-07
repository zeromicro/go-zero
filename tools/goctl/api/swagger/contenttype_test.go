package swagger

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func TestConsumesFromTypeOrDef(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		tp       spec.Type
		expected []string
	}{
		{
			name:     "GET method with nil type",
			method:   http.MethodGet,
			tp:       nil,
			expected: []string{},
		},
		{
			name:     "post nil",
			method:   http.MethodPost,
			tp:       nil,
			expected: []string{},
		},
		{
			name:   "json tag",
			method: http.MethodPost,
			tp: spec.DefineStruct{
				Members: []spec.Member{
					{
						Tag: `json:"example"`,
					},
				},
			},
			expected: []string{applicationJson},
		},
		{
			name:   "form tag",
			method: http.MethodPost,
			tp: spec.DefineStruct{
				Members: []spec.Member{
					{
						Tag: `form:"example"`,
					},
				},
			},
			expected: []string{applicationForm},
		},
		{
			name:     "Non struct type",
			method:   http.MethodPost,
			tp:       spec.ArrayType{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := consumesFromTypeOrDef(testingContext(t), tt.method, tt.tp)
			assert.Equal(t, tt.expected, result)
		})
	}
}
