package gen

import (
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

func genTypes(table Table, withCache bool) (string, error) {
	fields := table.Fields
	fieldsString, err := genFields(fields)
	if err != nil {
		return "", err
	}
	output, err := util.With("types").
		Parse(template.Types).
		Execute(map[string]interface{}{
			"withCache":             withCache,
			"upperStartCamelObject": table.Name.ToCamel(),
			"fields":                fieldsString,
		})
	if err != nil {
		return "", err
	}
	return output.String(), nil
}
