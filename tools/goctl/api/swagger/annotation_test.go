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
