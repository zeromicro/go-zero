package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTags_Get(t *testing.T) {
	tags := &Tags{
		tags: []*Tag{
			{Key: "json", Name: "foo", Options: []string{"omitempty"}},
			{Key: "xml", Name: "bar", Options: nil},
		},
	}

	tag, err := tags.Get("json")
	assert.NoError(t, err)
	assert.NotNil(t, tag)
	assert.Equal(t, "json", tag.Key)
	assert.Equal(t, "foo", tag.Name)

	_, err = tags.Get("yaml")
	assert.Error(t, err)

	var nilTags *Tags
	_, err = nilTags.Get("json")
	assert.Error(t, err)
}

func TestTags_Keys(t *testing.T) {
	tags := &Tags{
		tags: []*Tag{
			{Key: "json", Name: "foo", Options: []string{"omitempty"}},
			{Key: "xml", Name: "bar", Options: nil},
		},
	}

	keys := tags.Keys()
	expected := []string{"json", "xml"}
	assert.Equal(t, expected, keys)

	var nilTags *Tags
	nilKeys := nilTags.Keys()
	assert.Empty(t, nilKeys)
}

func TestTags_Tags(t *testing.T) {
	tags := &Tags{
		tags: []*Tag{
			{Key: "json", Name: "foo", Options: []string{"omitempty"}},
			{Key: "xml", Name: "bar", Options: nil},
		},
	}

	result := tags.Tags()
	assert.Len(t, result, 2)
	assert.Equal(t, "json", result[0].Key)
	assert.Equal(t, "foo", result[0].Name)
	assert.Equal(t, []string{"omitempty"}, result[0].Options)
	assert.Equal(t, "xml", result[1].Key)
	assert.Equal(t, "bar", result[1].Name)
	assert.Nil(t, result[1].Options)

	var nilTags *Tags
	nilResult := nilTags.Tags()
	assert.Empty(t, nilResult)
}
