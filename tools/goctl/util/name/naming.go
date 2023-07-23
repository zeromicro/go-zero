// Package name provides methods to verify naming style and format naming style
// See the method IsNamingValid, FormatFilename
package name

import (
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"
)

// NamingStyle the type of string
type NamingStyle = string

const (
	// NamingLower defines the lower spell case
	NamingLower NamingStyle = "lower"
	// NamingCamel defines the camel spell case
	NamingCamel NamingStyle = "camel"
	// NamingSnake defines the snake spell case
	NamingSnake NamingStyle = "snake"
)

// IsNamingValid validates whether the namingStyle is valid or not,return
// namingStyle and true if it is valid, or else return empty string
// and false, and it is a valid value even namingStyle is empty string
func IsNamingValid(namingStyle string) (NamingStyle, bool) {
	if len(namingStyle) == 0 {
		namingStyle = NamingLower
	}
	switch namingStyle {
	case NamingLower, NamingCamel, NamingSnake:
		return namingStyle, true
	default:
		return "", false
	}
}

// FormatFilename converts the filename string to the target
// naming style by calling method of stringx
func FormatFilename(filename string, style NamingStyle) string {
	switch style {
	case NamingCamel:
		return stringx.From(filename).ToCamel()
	case NamingSnake:
		return stringx.From(filename).ToSnake()
	default:
		return strings.ToLower(stringx.From(filename).ToCamel())
	}
}
