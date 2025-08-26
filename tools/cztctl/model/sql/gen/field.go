package gen

import (
	"strings"

	"github.com/lerity-yao/go-zero/tools/cztctl/model/sql/parser"
	"github.com/lerity-yao/go-zero/tools/cztctl/model/sql/template"
	"github.com/lerity-yao/go-zero/tools/cztctl/util"
	"github.com/lerity-yao/go-zero/tools/cztctl/util/pathx"
)

func genFields(table Table, fields []*parser.Field) (string, error) {
	var list []string

	for _, field := range fields {
		result, err := genField(table, field)
		if err != nil {
			return "", err
		}

		list = append(list, result)
	}

	return strings.Join(list, "\n"), nil
}

func genField(table Table, field *parser.Field) (string, error) {
	tag, err := genTag(table, field.NameOriginal)
	if err != nil {
		return "", err
	}

	text, err := pathx.LoadTemplate(category, fieldTemplateFile, template.Field)
	if err != nil {
		return "", err
	}

	output, err := util.With("types").
		Parse(text).
		Execute(map[string]any{
			"name":       util.SafeString(field.Name.ToCamel()),
			"type":       field.DataType,
			"tag":        tag,
			"hasComment": field.Comment != "",
			"comment":    field.Comment,
			"data":       table,
		})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
