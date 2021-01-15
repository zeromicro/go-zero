package spec

import (
	"errors"
	"strings"

	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

const (
	bodyTagKey        = "json"
	formTagKey        = "form"
	defaultSummaryKey = "summary"
)

var definedKeys = []string{bodyTagKey, formTagKey, "path"}

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

func (m Member) IsOmitEmpty() bool {
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
	return strings.TrimSpace(m.Comment)
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

func (m Member) IsFormMember() bool {
	if m.IsInline {
		return false
	}

	tags := m.Tags()
	for _, tag := range tags {
		if tag.Key == formTagKey {
			return true
		}
	}
	return false
}

func (t DefineStruct) GetBodyMembers() []Member {
	var result []Member
	for _, member := range t.Members {
		if member.IsBodyMember() {
			result = append(result, member)
		}
	}
	return result
}

func (t DefineStruct) GetFormMembers() []Member {
	var result []Member
	for _, member := range t.Members {
		if member.IsFormMember() {
			result = append(result, member)
		}
	}
	return result
}

func (t DefineStruct) GetNonBodyMembers() []Member {
	var result []Member
	for _, member := range t.Members {
		if !member.IsBodyMember() {
			result = append(result, member)
		}
	}
	return result
}

func (r Route) JoinedDoc() string {
	doc := r.AtDoc.Text
	if r.AtDoc.Properties != nil {
		doc += r.AtDoc.Properties[defaultSummaryKey]
	}
	doc += strings.Join(r.Docs, " ")
	return strings.TrimSpace(doc)
}

func (r Route) GetAnnotation(key string) string {
	if r.Annotation.Properties == nil {
		return ""
	}

	return r.Annotation.Properties[key]
}

func (g Group) GetAnnotation(key string) string {
	if g.Annotation.Properties == nil {
		return ""
	}

	return g.Annotation.Properties[key]
}

func (r Route) ResponseTypeName() string {
	if r.ResponseType == nil {
		return ""
	}

	return r.ResponseType.Name()
}

func (r Route) RequestTypeName() string {
	if r.RequestType == nil {
		return ""
	}

	return r.RequestType.Name()
}
