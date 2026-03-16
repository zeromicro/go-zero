package swagger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getBoolFromKVOrDefault(t *testing.T) {
	properties := map[string]string{
		"enabled":     `"true"`,
		"disabled":    `"false"`,
		"invalid":     `"notabool"`,
		"empty_value": `""`,
	}

	assert.True(t, getBoolFromKVOrDefault(properties, "enabled", false))
	assert.False(t, getBoolFromKVOrDefault(properties, "disabled", true))
	assert.False(t, getBoolFromKVOrDefault(properties, "invalid", false))
	assert.True(t, getBoolFromKVOrDefault(properties, "missing", true))
	assert.False(t, getBoolFromKVOrDefault(properties, "empty_value", false))
	assert.False(t, getBoolFromKVOrDefault(nil, "nil", false))
	assert.False(t, getBoolFromKVOrDefault(map[string]string{}, "empty", false))

	// Test with unquoted values (as stored by RawText())
	unquotedProperties := map[string]string{
		"enabled":     "true",
		"disabled":    "false",
		"invalid":     "notabool",
		"empty_value": "",
	}

	assert.True(t, getBoolFromKVOrDefault(unquotedProperties, "enabled", false))
	assert.False(t, getBoolFromKVOrDefault(unquotedProperties, "disabled", true))
	assert.False(t, getBoolFromKVOrDefault(unquotedProperties, "invalid", false))
	assert.False(t, getBoolFromKVOrDefault(unquotedProperties, "empty_value", false))
}

func Test_getStringFromKVOrDefault(t *testing.T) {
	properties := map[string]string{
		"name":  `"example"`,
		"empty": `""`,
	}

	assert.Equal(t, "example", getStringFromKVOrDefault(properties, "name", "default"))
	assert.Equal(t, "default", getStringFromKVOrDefault(properties, "empty", "default"))
	assert.Equal(t, "default", getStringFromKVOrDefault(properties, "missing", "default"))
	assert.Equal(t, "default", getStringFromKVOrDefault(nil, "nil", "default"))
	assert.Equal(t, "default", getStringFromKVOrDefault(map[string]string{}, "empty", "default"))

	// Test with unquoted values (as stored by RawText())
	unquotedProperties := map[string]string{
		"name":  "example",
		"title": "Demo API",
		"empty": "",
	}

	assert.Equal(t, "example", getStringFromKVOrDefault(unquotedProperties, "name", "default"))
	assert.Equal(t, "Demo API", getStringFromKVOrDefault(unquotedProperties, "title", "default"))
	assert.Equal(t, "default", getStringFromKVOrDefault(unquotedProperties, "empty", "default"))
}

func Test_getListFromInfoOrDefault(t *testing.T) {
	properties := map[string]string{
		"list":  `"a, b, c"`,
		"empty": `""`,
	}

	assert.Equal(t, []string{"a", " b", " c"}, getListFromInfoOrDefault(properties, "list", []string{"default"}))
	assert.Equal(t, []string{"default"}, getListFromInfoOrDefault(properties, "empty", []string{"default"}))
	assert.Equal(t, []string{"default"}, getListFromInfoOrDefault(properties, "missing", []string{"default"}))
	assert.Equal(t, []string{"default"}, getListFromInfoOrDefault(nil, "nil", []string{"default"}))
	assert.Equal(t, []string{"default"}, getListFromInfoOrDefault(map[string]string{}, "empty", []string{"default"}))
	assert.Equal(t, []string{"default"}, getListFromInfoOrDefault(map[string]string{
		"foo": ",,",
	}, "foo", []string{"default"}))

	// Test with unquoted values (as stored by RawText())
	unquotedProperties := map[string]string{
		"list":    "a, b, c",
		"schemes": "http,https",
		"tags":    "query",
		"empty":   "",
	}

	// Note: FieldsAndTrimSpace doesn't actually trim the spaces from returned values
	assert.Equal(t, []string{"a", " b", " c"}, getListFromInfoOrDefault(unquotedProperties, "list", []string{"default"}))
	assert.Equal(t, []string{"http", "https"}, getListFromInfoOrDefault(unquotedProperties, "schemes", []string{"default"}))
	assert.Equal(t, []string{"query"}, getListFromInfoOrDefault(unquotedProperties, "tags", []string{"default"}))
	assert.Equal(t, []string{"default"}, getListFromInfoOrDefault(unquotedProperties, "empty", []string{"default"}))
}

func Test_getFirstUsableString(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		result := getFirstUsableString()
		assert.Equal(t, "", result, "should return empty string for no arguments")
	})

	t.Run("single plain string", func(t *testing.T) {
		result := getFirstUsableString("Check server health status.")
		assert.Equal(t, "Check server health status.", result)
	})

	t.Run("single quoted string", func(t *testing.T) {
		// This is how Go would represent a quoted string literal
		result := getFirstUsableString(`"Check server health status."`)
		assert.Equal(t, "Check server health status.", result, "should unquote quoted strings")
	})

	t.Run("multiple plain strings", func(t *testing.T) {
		result := getFirstUsableString("", "second", "third")
		assert.Equal(t, "second", result, "should return first non-empty string")
	})

	t.Run("handler name fallback", func(t *testing.T) {
		// Simulates the real use case: @doc text, handler name
		result := getFirstUsableString("", "HealthCheck")
		assert.Equal(t, "HealthCheck", result, "should fallback to handler name")
	})

	t.Run("doc text over handler name", func(t *testing.T) {
		// Simulates the real use case with @doc text
		result := getFirstUsableString("Check server health status.", "HealthCheck")
		assert.Equal(t, "Check server health status.", result, "should use doc text over handler name")
	})

	t.Run("empty strings before valid", func(t *testing.T) {
		result := getFirstUsableString("", "", "valid")
		assert.Equal(t, "valid", result, "should skip empty strings")
	})

	t.Run("all empty strings", func(t *testing.T) {
		result := getFirstUsableString("", "", "")
		assert.Equal(t, "", result, "should return empty if all are empty")
	})

	t.Run("quoted then plain", func(t *testing.T) {
		result := getFirstUsableString(`"quoted"`, "plain")
		assert.Equal(t, "quoted", result, "should unquote first quoted string")
	})

	t.Run("plain then quoted", func(t *testing.T) {
		result := getFirstUsableString("plain", `"quoted"`)
		assert.Equal(t, "plain", result, "should use first plain string")
	})

	t.Run("invalid quoted string", func(t *testing.T) {
		// String that looks quoted but isn't valid Go syntax
		result := getFirstUsableString(`"incomplete`, "fallback")
		assert.Equal(t, `"incomplete`, result, "should use as-is if unquote fails but not empty")
	})

	t.Run("whitespace only", func(t *testing.T) {
		result := getFirstUsableString("   ", "fallback")
		assert.Equal(t, "   ", result, "should not trim whitespace, return as-is")
	})

	t.Run("real world API doc scenario", func(t *testing.T) {
		// This is the actual bug scenario from issue #5229
		atDocText := "Check server health status."
		handlerName := "HealthCheck"
		
		result := getFirstUsableString(atDocText, handlerName)
		assert.Equal(t, "Check server health status.", result, 
			"should use @doc text for API summary")
	})

	t.Run("real world with empty doc", func(t *testing.T) {
		// When @doc is empty, should fall back to handler name
		atDocText := ""
		handlerName := "HealthCheck"
		
		result := getFirstUsableString(atDocText, handlerName)
		assert.Equal(t, "HealthCheck", result, 
			"should fallback to handler name when @doc is empty")
	})

	t.Run("complex summary with special characters", func(t *testing.T) {
		result := getFirstUsableString("Get user by ID: /users/{id}")
		assert.Equal(t, "Get user by ID: /users/{id}", result, 
			"should handle special characters in plain strings")
	})

	t.Run("multiline string", func(t *testing.T) {
		result := getFirstUsableString("Line 1\nLine 2")
		assert.Equal(t, "Line 1\nLine 2", result, 
			"should handle multiline strings")
	})

	t.Run("unicode characters", func(t *testing.T) {
		result := getFirstUsableString("健康检查", "HealthCheck")
		assert.Equal(t, "健康检查", result, 
			"should handle unicode characters")
	})
}

