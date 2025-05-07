package swagger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_generateDefaultOperationId(t *testing.T) {
	testCases := []struct {
		groupName   string
		handlerName string
		expected    string
	}{
		{
			groupName:   "group",
			handlerName: "handler",
			expected:    "group_handler",
		},
		{
			groupName:   "",
			handlerName: "handler",
			expected:    "handler",
		},
	}

	for _, tc := range testCases {
		require.Equal(t, tc.expected, generateDefaultOperationId(tc.groupName, tc.handlerName))
	}
}