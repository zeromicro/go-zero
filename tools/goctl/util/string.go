package util

import (
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/util/console"
)

var goKeyword = map[string]string{
	"var":         "variable",
	"const":       "constant",
	"package":     "pkg",
	"func":        "function",
	"return":      "rtn",
	"defer":       "dfr",
	"go":          "goo",
	"select":      "slt",
	"struct":      "structure",
	"interface":   "itf",
	"chan":        "channel",
	"type":        "tp",
	"map":         "mp",
	"range":       "rg",
	"break":       "brk",
	"case":        "caz",
	"continue":    "ctn",
	"for":         "fr",
	"fallthrough": "fth",
	"else":        "es",
	"if":          "ef",
	"switch":      "swt",
	"goto":        "gt",
	"default":     "dft",
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

// EscapeGolangKeyword escapes the golang keywords.
func EscapeGolangKeyword(s string) string {
	if !isGolangKeyword(s) {
		return s
	}

	r := goKeyword[s]
	console.Info("[EscapeGolangKeyword]: go keyword is forbidden %q, converted into %q", s, r)
	return r
}

func isGolangKeyword(s string) bool {
	_, ok := goKeyword[s]
	return ok
}

func TrimWhiteSpace(s string) string {
	r := strings.NewReplacer(" ", "", "\t", "", "\n", "", "\f", "", "\r", "")
	return r.Replace(s)
}

func IsEmptyStringOrWhiteSpace(s string) bool {
	v := TrimWhiteSpace(s)
	return len(v) == 0
}
