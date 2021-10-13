package util

import (
	"strings"

	"github.com/tal-tech/go-zero/core/lang"
)

const escapePrefix = "es_"

var goKeyword = map[string]lang.PlaceholderType{
	"var":         lang.Placeholder,
	"const":       lang.Placeholder,
	"package":     lang.Placeholder,
	"func":        lang.Placeholder,
	"return":      lang.Placeholder,
	"defer":       lang.Placeholder,
	"go":          lang.Placeholder,
	"select":      lang.Placeholder,
	"struct":      lang.Placeholder,
	"interface":   lang.Placeholder,
	"chan":        lang.Placeholder,
	"type":        lang.Placeholder,
	"map":         lang.Placeholder,
	"range":       lang.Placeholder,
	"break":       lang.Placeholder,
	"case":        lang.Placeholder,
	"continue":    lang.Placeholder,
	"for":         lang.Placeholder,
	"fallthrough": lang.Placeholder,
	"else":        lang.Placeholder,
	"if":          lang.Placeholder,
	"switch":      lang.Placeholder,
	"goto":        lang.Placeholder,
	"default":     lang.Placeholder,
}

// Title returns a string value with s[0] which has been convert into upper case that
// there are not empty input text
func Title(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]
}

// Untitle returns a string value with s[0] which has been convert into lower case that
// there are not empty input text
func Untitle(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToLower(s[:1]) + s[1:]
}

// Index returns the index where the item equal,it will return -1 if mismatched
func Index(slice []string, item string) int {
	for i := range slice {
		if slice[i] == item {
			return i
		}
	}

	return -1
}

// SafeString converts the input string into a safe naming style in golang
func SafeString(in string) string {
	if len(in) == 0 {
		return in
	}

	data := strings.Map(func(r rune) rune {
		if isSafeRune(r) {
			return r
		}
		return '_'
	}, in)

	headRune := rune(data[0])
	if isNumber(headRune) {
		return "_" + data
	}
	return data
}

func isSafeRune(r rune) bool {
	return isLetter(r) || isNumber(r) || r == '_'
}

func isLetter(r rune) bool {
	return 'A' <= r && r <= 'z'
}

func isNumber(r rune) bool {
	return '0' <= r && r <= '9'
}

func EscapeGolangKeyword(s string) string {
	if !isGolangKeyword(s) {
		return s
	}
	return escapePrefix + s
}

func isGolangKeyword(s string) bool {
	_, ok := goKeyword[s]
	return ok
}
