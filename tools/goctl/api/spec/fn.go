package spec

import (
	"errors"
	"regexp"
	"strings"

	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

const (
	TagKey    = "tag"
	NameKey   = "name"
	OptionKey = "option"
	BodyTag   = "json"
)

var (
	TagRe       = regexp.MustCompile(`(?P<tag>\w+):"(?P<name>[^,"]+)[,]?(?P<option>[^"]*)"`)
	TagSubNames = TagRe.SubexpNames()
	definedTags = []string{TagKey, NameKey, OptionKey}
)

type Attribute struct {
	Key   string
	value string
}

func (m Member) IsOptional() bool {
	var option string

	matches := TagRe.FindStringSubmatch(m.Tag)
	for i := range matches {
		name := TagSubNames[i]
		if name == OptionKey {
			option = matches[i]
		}
	}

	if len(option) == 0 {
		return false
	}

	fields := strings.Split(option, ",")
	for _, field := range fields {
		if field == "optional" || strings.HasPrefix(field, "default=") {
			return true
		}
	}

	return false
}

func (m Member) IsOmitempty() bool {
	var option string

	matches := TagRe.FindStringSubmatch(m.Tag)
	for i := range matches {
		name := TagSubNames[i]
		if name == OptionKey {
			option = matches[i]
		}
	}

	if len(option) == 0 {
		return false
	}

	fields := strings.Split(option, ",")
	for _, field := range fields {
		if field == "omitempty" {
			return true
		}
	}

	return false
}

func (m Member) GetAttributes() []Attribute {
	matches := TagRe.FindStringSubmatch(m.Tag)
	var result []Attribute
	for i := range matches {
		name := TagSubNames[i]
		if stringx.Contains(definedTags, name) {
			result = append(result, Attribute{
				Key:   name,
				value: matches[i],
			})
		}
	}
	return result
}

func (m Member) GetPropertyName() (string, error) {
	attrs := m.GetAttributes()
	for _, attr := range attrs {
		if attr.Key == NameKey && len(attr.value) > 0 {
			if attr.value == "-" {
				return util.Untitle(m.Name), nil
			}
			return attr.value, nil
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
	attrs := m.GetAttributes()
	for _, attr := range attrs {
		if attr.value == BodyTag {
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
