package gen

import (
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

func genUpdate(table Table, withCache bool) (string, error) {
	expressionValues := make([]string, 0)
	for _, filed := range table.Fields {
		camel := filed.Name.ToCamel()
		if camel == "CreateTime" || camel == "UpdateTime" {
			continue
		}
		if filed.IsPrimaryKey {
			continue
		}
		expressionValues = append(expressionValues, "data."+camel)
	}
	expressionValues = append(expressionValues, "data."+table.PrimaryKey.Name.ToCamel())
	camelTableName := table.Name.ToCamel()
	output, err := util.With("update").
		Parse(template.Update).
		Execute(map[string]interface{}{
			"withCache":             withCache,
			"upperStartCamelObject": camelTableName,
			"primaryCacheKey":       table.CacheKey[table.PrimaryKey.Name.Source()].DataKeyExpression,
			"primaryKeyVariable":    table.CacheKey[table.PrimaryKey.Name.Source()].Variable,
			"lowerStartCamelObject": stringx.From(camelTableName).UnTitle(),
			"originalPrimaryKey":    table.PrimaryKey.Name.Source(),
			"expressionValues":      strings.Join(expressionValues, ", "),
		})
	if err != nil {
		return "", nil
	}
	return output.String(), nil
}
