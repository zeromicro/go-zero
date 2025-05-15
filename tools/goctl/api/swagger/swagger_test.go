package swagger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_pathVariable2SwaggerVariable(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{input: "/api/:id", expected: "/api/{id}"},
		{input: "/api/:id/details", expected: "/api/{id}/details"},
		{input: "/:version/api/:id", expected: "/{version}/api/{id}"},
		{input: "/api/v1", expected: "/api/v1"},
		{input: "/api/:id/:action", expected: "/api/{id}/{action}"},
	}

	for _, tc := range testCases {
		result := pathVariable2SwaggerVariable(testingContext(t), tc.input)
		assert.Equal(t, tc.expected, result)
	}
}
