package gen

import (
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

func genTypes(table Table, methods string, withCache bool) (string, error) {
	fields := table.Fields
	fieldsString, err := genFields(fields)
	if err != nil {
		return "", err
	}

	text, err := util.LoadTemplate(category, typesTemplateFile, template.Types)
	if err != nil {
		return "", err
	}

	output, err := util.With("types").
		Parse(text).
		Execute(map[string]interface{}{
			"withCache":             withCache,
			"method":                methods,
			"upperStartCamelObject": table.Name.ToCamel(),
			"fields":                fieldsString,
		})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
