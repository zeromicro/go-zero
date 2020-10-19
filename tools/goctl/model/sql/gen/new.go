package gen

import (
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	"github.com/tal-tech/go-zero/tools/goctl/templatex"
)

func genNew(table Table, withCache bool) (string, error) {
	text, err := templatex.LoadTemplate(category, modelNewTemplateFile, template.New)
	if err != nil {
		return "", err
	}
	output, err := templatex.With("new").
		Parse(text).
		Execute(map[string]interface{}{
			"withCache":             withCache,
			"upperStartCamelObject": table.Name.ToCamel(),
		})
	if err != nil {
		return "", err
	}
	return output.String(), nil
}
