package mcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParseArguments_MapStringString tests parsing map[string]string arguments
func TestParseArguments_MapStringString(t *testing.T) {
	// Sample request struct to populate
	type requestStruct struct {
		Name    string `json:"name"`
		Message string `json:"message"`
		Count   int    `json:"count"`
		Enabled bool   `json:"enabled"`
	}

	// Create test arguments
	args := map[string]string{
		"name":    "test-name",
		"message": "hello world",
		"count":   "42",
		"enabled": "true",
	}

	// Create a target object to populate
	var req requestStruct

	// Parse the arguments
	err := ParseArguments(args, &req)

	// Verify results
	assert.NoError(t, err, "Should parse map[string]string without error")
	assert.Equal(t, "test-name", req.Name, "Name should be correctly parsed")
	assert.Equal(t, "hello world", req.Message, "Message should be correctly parsed")
	assert.Equal(t, 42, req.Count, "Count should be correctly parsed to int")
	assert.True(t, req.Enabled, "Enabled should be correctly parsed to bool")
}

// TestParseArguments_MapStringAny tests parsing map[string]any arguments
func TestParseArguments_MapStringAny(t *testing.T) {
	// Sample request struct to populate
	type requestStruct struct {
		Name     string            `json:"name"`
		Message  string            `json:"message"`
		Count    int               `json:"count"`
		Enabled  bool              `json:"enabled"`
		Tags     []string          `json:"tags"`
		Metadata map[string]string `json:"metadata"`
	}

	// Create test arguments with mixed types
	args := map[string]any{
		"name":    "test-name",
		"message": "hello world",
		"count":   42,   // note: this is already an int
		"enabled": true, // note: this is already a bool
		"tags":    []string{"tag1", "tag2"},
		"metadata": map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	// Create a target object to populate
	var req requestStruct

	// Parse the arguments
	err := ParseArguments(args, &req)

	// Verify results
	assert.NoError(t, err, "Should parse map[string]any without error")
	assert.Equal(t, "test-name", req.Name, "Name should be correctly parsed")
	assert.Equal(t, "hello world", req.Message, "Message should be correctly parsed")
	assert.Equal(t, 42, req.Count, "Count should be correctly parsed")
	assert.True(t, req.Enabled, "Enabled should be correctly parsed")
	assert.Equal(t, []string{"tag1", "tag2"}, req.Tags, "Tags should be correctly parsed")
	assert.Equal(t, map[string]string{
		"key1": "value1",
		"key2": "value2",
	}, req.Metadata, "Metadata should be correctly parsed")
}

// TestParseArguments_UnsupportedType tests parsing with an unsupported type
func TestParseArguments_UnsupportedType(t *testing.T) {
	// Sample request struct to populate
	type requestStruct struct {
		Name    string `json:"name"`
		Message string `json:"message"`
	}

	// Use an unsupported argument type (slice)
	args := []string{"not", "a", "map"}

	// Create a target object to populate
	var req requestStruct

	// Parse the arguments
	err := ParseArguments(args, &req)

	// Verify error is returned with correct message
	assert.Error(t, err, "Should return error for unsupported type")
	assert.Contains(t, err.Error(), "unsupported argument type", "Error should mention unsupported type")
	assert.Contains(t, err.Error(), "[]string", "Error should include the actual type")
}

// TestParseArguments_EmptyMap tests parsing with empty maps
func TestParseArguments_EmptyMap(t *testing.T) {
	// Sample request struct to populate
	type requestStruct struct {
		Name    string `json:"name,optional"`
		Message string `json:"message,optional"`
	}

	// Test empty map[string]string
	t.Run("EmptyMapStringString", func(t *testing.T) {
		args := map[string]string{}
		var req requestStruct

		err := ParseArguments(args, &req)

		assert.NoError(t, err, "Should parse empty map[string]string without error")
		assert.Empty(t, req.Name, "Name should be empty string")
		assert.Empty(t, req.Message, "Message should be empty string")
	})

	// Test empty map[string]any
	t.Run("EmptyMapStringAny", func(t *testing.T) {
		args := map[string]any{}
		var req requestStruct

		err := ParseArguments(args, &req)

		assert.NoError(t, err, "Should parse empty map[string]any without error")
		assert.Empty(t, req.Name, "Name should be empty string")
		assert.Empty(t, req.Message, "Message should be empty string")
	})
}
