package swagger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func TestRangeValueFromOptions(t *testing.T) {
	tests := []struct {
		name            string
		options         []string
		expectedMin     *float64
		expectedMax     *float64
		expectedExclMin bool
		expectedExclMax bool
	}{
		{
			name:            "Valid range with inclusive bounds",
			options:         []string{"range=[1.0:10.0]"},
			expectedMin:     floatPtr(1.0),
			expectedMax:     floatPtr(10.0),
			expectedExclMin: false,
			expectedExclMax: false,
		},
		{
			name:            "Valid range with exclusive bounds",
			options:         []string{"range=(1.0:10.0)"},
			expectedMin:     floatPtr(1.0),
			expectedMax:     floatPtr(10.0),
			expectedExclMin: true,
			expectedExclMax: true,
		},
		{
			name:            "Invalid range format",
			options:         []string{"range=1.0:10.0"},
			expectedMin:     nil,
			expectedMax:     nil,
			expectedExclMin: false,
			expectedExclMax: false,
		},
		{
			name:            "Invalid range start",
			options:         []string{"range=[a:1.0)"},
			expectedMin:     nil,
			expectedMax:     nil,
			expectedExclMin: false,
			expectedExclMax: false,
		},
		{
			name:            "Missing range end",
			options:         []string{"range=[1.0:)"},
			expectedMin:     floatPtr(1.0),
			expectedMax:     nil,
			expectedExclMin: false,
			expectedExclMax: true,
		},
		{
			name:            "Missing range start and end",
			options:         []string{"range=[:)"},
			expectedMin:     nil,
			expectedMax:     nil,
			expectedExclMin: false,
			expectedExclMax: true,
		},
		{
			name:            "Missing range start",
			options:         []string{"range=[:1.0)"},
			expectedMin:     nil,
			expectedMax:     floatPtr(1.0),
			expectedExclMin: false,
			expectedExclMax: true,
		},
		{
			name:            "Invalid range end",
			options:         []string{"range=[1.0:b)"},
			expectedMin:     nil,
			expectedMax:     nil,
			expectedExclMin: false,
			expectedExclMax: false,
		},
		{
			name:            "Empty options",
			options:         []string{},
			expectedMin:     nil,
			expectedMax:     nil,
			expectedExclMin: false,
			expectedExclMax: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			min, max, exclMin, exclMax := rangeValueFromOptions(tt.options)
			assert.Equal(t, tt.expectedMin, min)
			assert.Equal(t, tt.expectedMax, max)
			assert.Equal(t, tt.expectedExclMin, exclMin)
			assert.Equal(t, tt.expectedExclMax, exclMax)
		})
	}
}

func TestEnumsValueFromOptions(t *testing.T) {
	tests := []struct {
		name     string
		options  []string
		expected []any
	}{
		{
			name:     "Valid enums",
			options:  []string{"options=a|b|c"},
			expected: []any{"a", "b", "c"},
		},
		{
			name:     "Empty enums",
			options:  []string{"options="},
			expected: []any{},
		},
		{
			name:     "No enum option",
			options:  []string{},
			expected: []any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := enumsValueFromOptions(tt.options)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefValueFromOptions(t *testing.T) {
	tests := []struct {
		name     string
		options  []string
		apiType  spec.Type
		expected any
	}{
		{
			name:     "Default integer value",
			options:  []string{"default=42"},
			apiType:  spec.PrimitiveType{RawName: "int"},
			expected: int64(42),
		},
		{
			name:     "Default string value",
			options:  []string{"default=hello"},
			apiType:  spec.PrimitiveType{RawName: "string"},
			expected: "hello",
		},
		{
			name:     "No default value",
			options:  []string{},
			apiType:  spec.PrimitiveType{RawName: "string"},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := defValueFromOptions(testingContext(t), tt.options, tt.apiType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExampleValueFromOptions(t *testing.T) {
	tests := []struct {
		name     string
		options  []string
		apiType  spec.Type
		expected any
	}{
		{
			name:     "Example value present",
			options:  []string{"example=3.14"},
			apiType:  spec.PrimitiveType{RawName: "float"},
			expected: 3.14,
		},
		{
			name:     "Fallback to default value",
			options:  []string{"default=42"},
			apiType:  spec.PrimitiveType{RawName: "int"},
			expected: int64(42),
		},
		{
			name:     "Fallback to default value",
			options:  []string{"default="},
			apiType:  spec.PrimitiveType{RawName: "int"},
			expected: int64(0),
		},
		{
			name:     "No example or default value",
			options:  []string{},
			apiType:  spec.PrimitiveType{RawName: "string"},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exampleValueFromOptions(testingContext(t), tt.options, tt.apiType)
		})
	}
}

func TestValueFromOptions(t *testing.T) {
	tests := []struct {
		name     string
		options  []string
		key      string
		tp       string
		expected any
	}{
		{
			name:     "Integer value",
			options:  []string{"default=42"},
			key:      "default=",
			tp:       "integer",
			expected: int64(42),
		},
		{
			name:     "Boolean value",
			options:  []string{"default=true"},
			key:      "default=",
			tp:       "boolean",
			expected: true,
		},
		{
			name:     "Number value",
			options:  []string{"default=1.1"},
			key:      "default=",
			tp:       "number",
			expected: 1.1,
		},
		{
			name:     "No matching key",
			options:  []string{"example=42"},
			key:      "default=",
			tp:       "integer",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := valueFromOptions(testingContext(t), tt.options, tt.key, tt.tp)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func floatPtr(f float64) *float64 {
	return &f
}
