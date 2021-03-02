package util

import (
	"strings"
	"unicode"
)

// IsUpperCase returns true if the rune in A-Z
func IsUpperCase(r rune) bool {
	if r >= 'A' && r <= 'Z' {
		return true
	}
	return false
}

// IsLowerCase returns true if the rune in a-z
func IsLowerCase(r rune) bool {
	if r >= 'a' && r <= 'z' {
		return true
	}
	return false
}

// ToSnakeCase returns a copy string by converting camel case into snake case
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

// ToCamelCase returns a copy string by converting snake case into camel case
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

// ToLowerCase converts rune into lower case
func ToLowerCase(r rune) rune {
	dx := 'A' - 'a'
	if IsUpperCase(r) {
		return r - dx
	}
	return r
}

// ToUpperCase converts rune into upper case
func ToUpperCase(r rune) rune {
	dx := 'A' - 'a'
	if IsLowerCase(r) {
		return r + dx
	}
	return r
}

// ToLower returns a copy string by converting it into lower
func ToLower(s string) string {
	var out []rune
	for _, r := range s {
		out = append(out, ToLowerCase(r))
	}
	return string(out)
}

// ToUpper returns a copy string by converting it into upper
func ToUpper(s string) string {
	var out []rune
	for _, r := range s {
		out = append(out, ToUpperCase(r))
	}
	return string(out)
}

// UpperFirst converts s[0] into upper case
func UpperFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return ToUpper(s[:1]) + s[1:]
}

// UnExport converts the first letter into lower case
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
