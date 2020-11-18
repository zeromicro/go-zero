package name

import (
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

type NamingStyle = string

const (
	NamingLower NamingStyle = "lower"
	NamingCamel NamingStyle = "camel"
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
