package gen

import (
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	"github.com/tal-tech/go-zero/tools/goctl/util/templatex"
)

func genTypes(table Table) (string, error) {
	fields := table.Fields
	fieldsString, err := genFields(fields)
	if err != nil {
		return "", err
	}
	output, err := templatex.With("types").
		Parse(template.Types).
		Execute(map[string]interface{}{
			"upperStartCamelObject": table.Name.Snake2Camel(),
			"fields":                fieldsString,
		})
	if err != nil {
		return "", err
	}
	return output.String(), nil
}
