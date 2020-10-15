package gen

import (
	"strings"

	"github.com/tal-tech/go-zero/core/collection"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	"github.com/tal-tech/go-zero/tools/goctl/templatex"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

func genInsert(table Table, withCache bool) (string, error) {
	keySet := collection.NewSet()
	keyVariableSet := collection.NewSet()
	for fieldName, key := range table.CacheKey {
		if fieldName == table.PrimaryKey.Name.Source() {
			continue
		}
		keySet.AddStr(key.DataKeyExpression)
		keyVariableSet.AddStr(key.Variable)
	}

	expressions := make([]string, 0)
	expressionValues := make([]string, 0)
	for _, filed := range table.Fields {
		camel := filed.Name.ToCamel()
		if camel == "CreateTime" || camel == "UpdateTime" {
			continue
		}
		if filed.IsPrimaryKey && table.PrimaryKey.AutoIncrement {
			continue
		}
		expressions = append(expressions, "?")
		expressionValues = append(expressionValues, "data."+camel)
	}
	camel := table.Name.ToCamel()
	output, err := templatex.With("insert").
		Parse(template.Insert).
		Execute(map[string]interface{}{
			"withCache":             withCache,
			"containsIndexCache":    table.ContainsUniqueKey,
			"upperStartCamelObject": camel,
			"lowerStartCamelObject": stringx.From(camel).UnTitle(),
			"expression":            strings.Join(expressions, ", "),
			"expressionValues":      strings.Join(expressionValues, ", "),
			"keys":                  strings.Join(keySet.KeysStr(), "\n"),
			"keyValues":             strings.Join(keyVariableSet.KeysStr(), ", "),
		})
	if err != nil {
		return "", err
	}
	return output.String(), nil
}
