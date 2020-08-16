package gen

import (
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
)

func genImports(withCache bool) string {
	if withCache {
		return template.Imports
	} else {
		return template.ImportsNoCache
	}
}
