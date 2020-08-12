package gen

import (
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	"github.com/tal-tech/go-zero/tools/goctl/util/templatex"
)

func genImports(withCache bool) (string, error) {
	output, err := templatex.With("import").
		Parse(template.Imports).
		Execute(map[string]interface{}{
			"withCache": withCache,
		})
	if err != nil {
		return "", err
	}
	return output.String(), nil
}
