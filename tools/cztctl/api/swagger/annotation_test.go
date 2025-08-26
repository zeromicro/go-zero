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
}
