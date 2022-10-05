package spec

import (
	"errors"
	"path"
	"strings"

	"github.com/zeromicro/go-zero/core/stringx"
	"github.com/zeromicro/go-zero/tools/goctl/util"
)

const (
	bodyTagKey        = "json"
	formTagKey        = "form"
	pathTagKey        = "path"
	headerTagKey      = "header"
	defaultSummaryKey = "summary"
)

var definedKeys = []string{bodyTagKey, formTagKey, pathTagKey, headerTagKey}

func (s Service) JoinPrefix() Service {
	var groups []Group
	for _, g := range s.Groups {
		prefix := strings.TrimSpace(g.GetAnnotation(RoutePrefixKey))
		prefix = strings.ReplaceAll(prefix, `"`, "")
		var routes []Route
		for _, r := range g.Routes {
			r.Path = path.Join("/", prefix, r.Path)
			routes = append(routes, r)
		}
		g.Routes = routes
		groups = append(groups, g)
	}
	s.Groups = groups
	return s
}

// Routes returns all routes in api service
func (s Service) Routes() []Route {
	var result []Route
	for _, group := range s.Groups {
		result = append(result, group.Routes...)
	}
	return result
}

// Tags returns all tags in Member
func (m Member) Tags() []*Tag {
	tags, err := Parse(m.Tag)
	if err != nil {
		panic(m.Tag + ", " + err.Error())
	}

	return tags.Tags()
}

// IsOptional returns true if tag is optional
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

// IsOmitEmpty returns true if tag contains omitempty
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

// GetPropertyName returns json tag value
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

// GetComment returns comment value of Member
func (m Member) GetComment() string {
	return strings.TrimSpace(m.Comment)
}

// IsBodyMember returns true if contains json tag
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

// IsFormMember returns true if contains form tag
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

// IsTagMember returns true if contains given tag
func (m Member) IsTagMember(tagKey string) bool {
	if m.IsInline {
		return true
	}

	tags := m.Tags()
	for _, tag := range tags {
		if tag.Key == tagKey {
			return true
		}
	}
	return false
}

// GetEnumOptions return a slice contains all enumeration options
func (m Member) GetEnumOptions() []string {
	if !m.IsBodyMember() {
		return nil
	}

	tags := m.Tags()
	for _, tag := range tags {
		if tag.Key == bodyTagKey {
			options := tag.Options
			for _, option := range options {
				if strings.Index(option, "options=") == 0 {
					option = strings.TrimPrefix(option, "options=")
					return strings.Split(option, "|")
				}
			}
		}
	}
	return nil
}

// GetBodyMembers returns all json fields
func (t DefineStruct) GetBodyMembers() []Member {
	var result []Member
	for _, member := range t.Members {
		if member.IsBodyMember() {
			result = append(result, member)
		}
	}
	return result
}

// GetFormMembers returns all form fields
func (t DefineStruct) GetFormMembers() []Member {
	var result []Member
	for _, member := range t.Members {
		if member.IsFormMember() {
			result = append(result, member)
		}
	}
	return result
}

// GetNonBodyMembers returns all have no tag fields
func (t DefineStruct) GetNonBodyMembers() []Member {
	var result []Member
	for _, member := range t.Members {
		if !member.IsBodyMember() {
			result = append(result, member)
		}
	}
	return result
}

// GetTagMembers returns all given key fields
func (t DefineStruct) GetTagMembers(tagKey string) []Member {
	var result []Member
	for _, member := range t.Members {
		if member.IsTagMember(tagKey) {
			result = append(result, member)
		}
	}
	return result
}

// JoinedDoc joins comments and summary value in AtDoc
func (r Route) JoinedDoc() string {
	doc := r.AtDoc.Text
	if r.AtDoc.Properties != nil {
		doc += r.AtDoc.Properties[defaultSummaryKey]
	}
	doc += strings.Join(r.Docs, " ")
	return strings.TrimSpace(doc)
}

// GetAnnotation returns the value by specified key from @server
func (r Route) GetAnnotation(key string) string {
	if r.AtServerAnnotation.Properties == nil {
		return ""
	}

	return r.AtServerAnnotation.Properties[key]
}

// GetAnnotation returns the value by specified key from @server
func (g Group) GetAnnotation(key string) string {
	if g.Annotation.Properties == nil {
		return ""
	}

	return g.Annotation.Properties[key]
}

// ResponseTypeName returns response type name of route
func (r Route) ResponseTypeName() string {
	if r.ResponseType == nil {
		return ""
	}

	return r.ResponseType.Name()
}

// RequestTypeName returns request type name of route
func (r Route) RequestTypeName() string {
	if r.RequestType == nil {
		return ""
	}

	return r.RequestType.Name()
}
