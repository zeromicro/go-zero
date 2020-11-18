package generator

import (
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

func formatFilename(filename string, style NamingStyle) string {
	switch style {
	case namingCamel:
		return stringx.From(filename).ToCamel()
	case namingSnake:
		return stringx.From(filename).ToSnake()
	default:
		return strings.ToLower(stringx.From(filename).ToCamel())
	}
}
