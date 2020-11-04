package generator

import (
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

func formatFilename(filename string) string {
	return strings.ToLower(stringx.From(filename).ToCamel())
}
