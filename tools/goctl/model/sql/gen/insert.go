package gen

import (
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/util/templatex"
)

func genInsert(table Table, withCache bool) (string, error) {
	expressions := make([]string, 0)
	expressionValues := make([]string, 0)
	for _, filed := range table.Fields {
		camel := filed.Name.Snake2Camel()
		if camel == "CreateTime" || camel == "UpdateTime" {
			continue
		}
		if filed.IsPrimaryKey && table.PrimaryKey.AutoIncrement {
			continue
		}
		expressions = append(expressions, "?")
		expressionValues = append(expressionValues, "data."+camel)
	}
	camel := table.Name.Snake2Camel()
	output, err := templatex.With("insert").
		Parse(template.Insert).
		Execute(map[string]interface{}{
			"withCache":             withCache,
			"upperStartCamelObject": camel,
			"lowerStartCamelObject": stringx.From(camel).LowerStart(),
			"expression":            strings.Join(expressions, ", "),
			"expressionValues":      strings.Join(expressionValues, ", "),
		})
	if err != nil {
		return "", err
	}
	return output.String(), nil
}
