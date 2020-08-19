package stringx

import (
	"bytes"
	"strings"
	"unicode"
)

type (
	String struct {
		source string
	}
)

func From(data string) String {
	return String{source: data}
}

func (s String) IsEmptyOrSpace() bool {
	if len(s.source) == 0 {
		return true
	}
	if strings.TrimSpace(s.source) == "" {
		return true
	}
	return false
}

func (s String) Lower() string {
	return strings.ToLower(s.source)
}
func (s String) Upper() string {
	return strings.ToUpper(s.source)
}
func (s String) Title() string {
	if s.IsEmptyOrSpace() {
		return s.source
	}
	return strings.Title(s.source)
}

// snake->camel(upper start)
func (s String) ToCamel() string {
	list := s.splitBy(func(r rune) bool {
		return r == '_'
	}, true)
	var target []string
	for _, item := range list {
		target = append(target, From(item).Title())
	}
	return strings.Join(target, "")
}

// camel->snake
func (s String) ToSnake() string {
	list := s.splitBy(func(r rune) bool {
		return unicode.IsUpper(r)
	}, false)
	var target []string
	for _, item := range list {
		target = append(target, From(item).Lower())
	}
	return strings.Join(target, "_")
}

// return original string if rune is not letter at index 0
func (s String) UnTitle() string {
	if s.IsEmptyOrSpace() {
		return s.source
	}
	r := rune(s.source[0])
	if !unicode.IsUpper(r) && !unicode.IsLower(r) {
		return s.source
	}
	return string(unicode.ToLower(r)) + s.source[1:]
}

// it will not ignore spaces
func (s String) splitBy(fn func(r rune) bool, remove bool) []string {
	if s.IsEmptyOrSpace() {
		return nil
	}
	var list []string
	buffer := new(bytes.Buffer)
	for _, r := range s.source {
		if fn(r) {
			if buffer.Len() != 0 {
				list = append(list, buffer.String())
				buffer.Reset()
			}
			if !remove {
				buffer.WriteRune(r)
			}
			continue
		}
		buffer.WriteRune(r)
	}
	if buffer.Len() != 0 {
		list = append(list, buffer.String())
	}
	return list
}

func (s String) ReplaceAll(old, new string) string {
	return strings.ReplaceAll(s.source, old, new)
}

func (s String) Source() string {
	return s.source
}
