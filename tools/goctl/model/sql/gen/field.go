package gen

import (
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/model/sql/parser"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	"github.com/tal-tech/go-zero/tools/goctl/util/templatex"
)

func genFields(fields []parser.Field) (string, error) {
	var list []string
	for _, field := range fields {
		result, err := genField(field)
		if err != nil {
			return "", err
		}
		list = append(list, result)
	}
	return strings.Join(list, "\n"), nil
}

func genField(field parser.Field) (string, error) {
	output, err := templatex.With("types").
		Parse(template.Field).
		Execute(map[string]string{
			"name":    field.Name.Snake2Camel(),
			"type":    field.DataType,
			"tag":     field.Name.Source(),
			"comment": field.Comment,
		})
	if err != nil {
		return "", err
	}
	return output.String(), nil
}
