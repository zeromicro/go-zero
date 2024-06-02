package stringx

import (
	"bytes"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var WhiteSpace = []rune{'\n', '\t', '\f', '\v', ' '}

// String  provides for converting the source text into other spell case,like lower,snake,camel
type String struct {
	source string
}

// From converts the input text to String and returns it
func From(data string) String {
	return String{source: data}
}

// IsEmptyOrSpace returns true if the length of the string value is 0 after call strings.TrimSpace, or else returns false
func (s String) IsEmptyOrSpace() bool {
	if len(s.source) == 0 {
		return true
	}
	if strings.TrimSpace(s.source) == "" {
		return true
	}
	return false
}

// Lower calls the strings.ToLower
func (s String) Lower() string {
	return strings.ToLower(s.source)
}

// Upper calls the strings.ToUpper
func (s String) Upper() string {
	return strings.ToUpper(s.source)
}

// ReplaceAll calls the strings.ReplaceAll
func (s String) ReplaceAll(old, new string) string {
	return strings.Replace(s.source, old, new, -1)
}

// Source returns the source string value
func (s String) Source() string {
	return s.source
}

// Title calls the cases.Title
func (s String) Title() string {
	if s.IsEmptyOrSpace() {
		return s.source
	}
	return cases.Title(language.English, cases.NoLower).String(s.source)
}

// ToCamel converts the input text into camel case
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

// ToSnake converts the input text into snake case
func (s String) ToSnake() string {
	list := s.splitBy(unicode.IsUpper, false)
	var target []string
	for _, item := range list {
		target = append(target, From(item).Lower())
	}
	return strings.Join(target, "_")
}

// Untitle return the original string if rune is not letter at index 0
func (s String) Untitle() string {
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

func ContainsAny(s string, runes ...rune) bool {
	if len(runes) == 0 {
		return true
	}
	tmp := make(map[rune]struct{}, len(runes))
	for _, r := range runes {
		tmp[r] = struct{}{}
	}

	for _, r := range s {
		if _, ok := tmp[r]; ok {
			return true
		}
	}
	return false
}

func ContainsWhiteSpace(s string) bool {
	return ContainsAny(s, WhiteSpace...)
}

func IsWhiteSpace(text string) bool {
	if len(text) == 0 {
		return true
	}
	for _, r := range text {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}
