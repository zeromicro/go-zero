package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareVersions(t *testing.T) {
	cases := []struct {
		ver1     string
		ver2     string
		operator string
		out      bool
	}{
		{"1", "1.0.1", ">", false},
		{"1.0.1", "1.0", "<", false},
		{"1", "0.9.9", ">", true},
		{"1", "1.0-1", "<", true},
		{"1", "1.0-1", "!", false},
		{"1.0.1", "1-0.1", "<", false},
		{"1.0.1", "1.0.1", "==", true},
		{"1.0.1", "1.0.2", "==", false},
		{"1.1-1", "1.0.2", "==", false},
		{"1.0.1", "1.0.2", ">=", false},
		{"1.0.2", "1.0.2", ">=", true},
		{"1.0.3", "1.0.2", ">=", true},
		{"1.0.4", "1.0.2", "<=", false},
		{"1.0.4", "1.0.6", "<=", true},
		{"1.0.4", "1.0.4", "<=", true},
	}

	for _, each := range cases {
		each := each
		t.Run(each.ver1, func(t *testing.T) {
			actual := CompareVersions(each.ver1, each.operator, each.ver2)
			assert.Equal(t, each.out, actual, fmt.Sprintf("%s vs %s", each.ver1, each.ver2))
		})
	}
}

func TestStrsToInts(t *testing.T) {
	testCases := []struct {
		input    []string
		expected []int64
	}{
		{[]string{}, nil},
		{[]string{"1", "2", "3"}, []int64{1, 2, 3}},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run("", func(t *testing.T) {
			actual := strsToInts(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
