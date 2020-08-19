package gen

import (
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/util/templatex"
)

func genUpdate(table Table, withCache bool) (string, error) {
	expressionValues := make([]string, 0)
	for _, filed := range table.Fields {
		camel := filed.Name.Snake2Camel()
		if camel == "CreateTime" || camel == "UpdateTime" {
			continue
		}
		if filed.IsPrimaryKey {
			continue
		}
		expressionValues = append(expressionValues, "data."+camel)
	}
	expressionValues = append(expressionValues, "data."+table.PrimaryKey.Name.Snake2Camel())
	camelTableName := table.Name.Snake2Camel()
	output, err := templatex.With("update").
		Parse(template.Update).
		Execute(map[string]interface{}{
			"withCache":             withCache,
			"upperStartCamelObject": camelTableName,
			"primaryCacheKey":       table.CacheKey[table.PrimaryKey.Name.Source()].DataKeyExpression,
			"primaryKeyVariable":    table.CacheKey[table.PrimaryKey.Name.Source()].Variable,
			"lowerStartCamelObject": stringx.From(camelTableName).LowerStart(),
			"originalPrimaryKey":    table.PrimaryKey.Name.Source(),
			"expressionValues":      strings.Join(expressionValues, ", "),
		})
	if err != nil {
		return "", nil
	}
	return output.String(), nil
}
