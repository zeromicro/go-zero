package spec

import (
	"errors"
	"strings"

	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

const bodyTagKey = "json"

var definedKeys = []string{"json", "form", "path"}

func (s Service) Routes() []Route {
	var result []Route
	for _, group := range s.Groups {
		result = append(result, group.Routes...)
	}
	return result
}

func (m Member) Tags() []*Tag {
	tags, err := Parse(m.Tag)
	if err != nil {
		panic(m.Tag + ", " + err.Error())
	}

	return tags.Tags()
}

func (m Member) IsOptional() bool {
	if !m.IsBodyMember() {
		return false
	}

	tag := m.Tags()
	for _, item := range tag {
		if item.Key == bodyTagKey {
			if stringx.Contains(item.Options, "optional") {
				return true
			}
		}
	}
	return false
}

func (m Member) IsOmitempty() bool {
	if !m.IsBodyMember() {
		return false
	}

	tag := m.Tags()
	for _, item := range tag {
		if item.Key == bodyTagKey {
			if stringx.Contains(item.Options, "omitempty") {
				return true
			}
		}
	}
	return false
}

func (m Member) GetPropertyName() (string, error) {
	tags := m.Tags()
	for _, tag := range tags {
		if stringx.Contains(definedKeys, tag.Key) {
			if tag.Name == "-" {
				return util.Untitle(m.Name), nil
			}
			return tag.Name, nil
		}
	}

	return "", errors.New("json property name not exist, member: " + m.Name)
}

func (m Member) GetComment() string {
	return strings.TrimSpace(strings.Join(m.Comments, "; "))
}

func (m Member) IsBodyMember() bool {
	if m.IsInline {
		return true
	}

	tags := m.Tags()
	for _, tag := range tags {
		if tag.Key == bodyTagKey {
			return true
		}
	}
	return false
}

func (t Type) GetBodyMembers() []Member {
	var result []Member
	for _, member := range t.Members {
		if member.IsBodyMember() {
			result = append(result, member)
		}
	}
	return result
}

func (t Type) GetNonBodyMembers() []Member {
	var result []Member
	for _, member := range t.Members {
		if !member.IsBodyMember() {
			result = append(result, member)
		}
	}
	return result
}
