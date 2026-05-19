package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	assert.False(t, matchRegex("v1ddd", versionRegex))
}

func TestImportRegex(t *testing.T) {
	tests := []struct {
		value   string
		matched bool
	}{
		{`"bar.api"`, true},
		{`"foo/bar.api"`, true},
		{`"/foo/bar.api"`, true},
		{`"../foo/bar.api"`, true},
		{`"../../foo/bar.api"`, true},

		{`"//bar.api"`, false},
		{`"/foo/foo_bar.api"`, true},
	}
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			assert.Equal(t, tt.matched, matchRegex(tt.value, importValueRegex))
		})
	}
}

func TestIsBasicType(t *testing.T) {
	tests := []struct {
		value    string
		expected bool
	}{
		{"string", true},
		{"int", true},
		{"bool", true},
		{"float64", true},
		{"byte", true},
		{"rune", true},
		{"File", false}, // File is not a basic type
		{"User", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsBasicType(tt.value))
		})
	}
}

func TestIsFileType(t *testing.T) {
	tests := []struct {
		value    string
		expected bool
	}{
		{"File", true},
		{"file", false}, // case sensitive
		{"FILE", false}, // case sensitive
		{"string", false},
		{"int", false},
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsFileType(tt.value))
		})
	}
}
