package spec

import (
	"errors"
	"strings"

	"github.com/fatih/structtag"
)

var errTagNotExist = errors.New("tag does not exist")

type (
	// Tag defines a tag for structure filed
	Tag struct {
		// Key is the tag key, such as json, xml, etc..
		// i.e: `json:"foo,omitempty". Here key is: "json"
		Key string

		// Name is a part of the value
		// i.e: `json:"foo,omitempty". Here name is: "foo"
		Name string

		// Options is a part of the value. It contains a slice of tag options i.e:
		// `json:"foo,omitempty". Here options is: ["omitempty"]
		Options []string
	}

	// Tags defines a slice for Tag
	Tags struct {
		tags []*Tag
	}
)

// Parse converts tag string into Tag
func Parse(tag string) (*Tags, error) {
	tag = strings.TrimPrefix(tag, "`")
	tag = strings.TrimSuffix(tag, "`")
	tags, err := structtag.Parse(tag)
	if err != nil {
		return nil, err
	}

	var result Tags
	for _, item := range tags.Tags() {
		result.tags = append(result.tags, &Tag{Key: item.Key, Name: item.Name, Options: item.Options})
	}
	return &result, nil
}

// Get gets tag value by specified key
func (t *Tags) Get(key string) (*Tag, error) {
	for _, tag := range t.tags {
		if tag.Key == key {
			return tag, nil
		}
	}

	return nil, errTagNotExist
}

// Keys returns all keys in Tags
func (t *Tags) Keys() []string {
	var keys []string
	for _, tag := range t.tags {
		keys = append(keys, tag.Key)
	}
	return keys
}

// Tags returns all tags in Tags
func (t *Tags) Tags() []*Tag {
	return t.tags
}
