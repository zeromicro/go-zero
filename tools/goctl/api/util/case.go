package util

import (
	"strings"
	"unicode"
)

func IsUpperCase(r rune) bool {
	if r >= 'A' && r <= 'Z' {
		return true
	}
	return false
}

func IsLowerCase(r rune) bool {
	if r >= 'a' && r <= 'z' {
		return true
	}
	return false
}

func ToSnakeCase(s string) string {
	var out []rune
	for index, r := range s {
		if index == 0 {
			out = append(out, ToLowerCase(r))
			continue
		}

		if IsUpperCase(r) && index != 0 {
			if IsLowerCase(rune(s[index-1])) {
				out = append(out, '_', ToLowerCase(r))
				continue
			}
			if index < len(s)-1 && IsLowerCase(rune(s[index+1])) {
				out = append(out, '_', ToLowerCase(r))
				continue
			}
			out = append(out, ToLowerCase(r))
			continue
		}
		out = append(out, r)
	}
	return string(out)
}

func ToCamelCase(s string) string {
	s = ToLower(s)
	out := []rune{}
	for index, r := range s {
		if r == '_' {
			continue
		}
		if index == 0 {
			out = append(out, ToUpperCase(r))
			continue
		}

		if index > 0 && s[index-1] == '_' {
			out = append(out, ToUpperCase(r))
			continue
		}

		out = append(out, r)
	}
	return string(out)
}

func ToLowerCase(r rune) rune {
	dx := 'A' - 'a'
	if IsUpperCase(r) {
		return r - dx
	}
	return r
}
func ToUpperCase(r rune) rune {
	dx := 'A' - 'a'
	if IsLowerCase(r) {
		return r + dx
	}
	return r
}

func ToLower(s string) string {
	var out []rune
	for _, r := range s {
		out = append(out, ToLowerCase(r))
	}
	return string(out)
}

func ToUpper(s string) string {
	var out []rune
	for _, r := range s {
		out = append(out, ToUpperCase(r))
	}
	return string(out)
}

func LowerFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return ToLower(s[:1]) + s[1:]
}

func UpperFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return ToUpper(s[:1]) + s[1:]
}

func UnExport(text string) bool {
	var flag bool
	str := strings.Map(func(r rune) rune {
		if flag {
			return r
		}
		if unicode.IsLetter(r) {
			flag = true
			return unicode.ToLower(r)
		}
		return r
	}, text)
	return str == text
}
