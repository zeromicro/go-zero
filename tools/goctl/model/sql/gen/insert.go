package gen

import (
	"strings"

	"github.com/tal-tech/go-zero/core/collection"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

func genInsert(table Table, withCache bool) (string, string, error) {
	keySet := collection.NewSet()
	keyVariableSet := collection.NewSet()
	for _, key := range table.UniqueCacheKey {
		keySet.AddStr(key.DataKeyExpression)
		keyVariableSet.AddStr(key.KeyLeft)
	}

	expressions := make([]string, 0)
	expressionValues := make([]string, 0)
	for _, field := range table.Fields {
		camel := field.Name.ToCamel()
		if camel == "CreateTime" || camel == "UpdateTime" {
			continue
		}

		if field.Name.Source() == table.PrimaryKey.Name.Source() {
			if table.PrimaryKey.AutoIncrement {
				continue
			}
		}

		expressions = append(expressions, "?")
		expressionValues = append(expressionValues, "data."+camel)
	}

	camel := table.Name.ToCamel()
	text, err := util.LoadTemplate(category, insertTemplateFile, template.Insert)
	if err != nil {
		return "", "", err
	}

	output, err := util.With("insert").
		Parse(text).
		Execute(map[string]interface{}{
			"withCache":             withCache,
			"containsIndexCache":    table.ContainsUniqueCacheKey,
			"upperStartCamelObject": camel,
			"lowerStartCamelObject": stringx.From(camel).Untitle(),
			"expression":            strings.Join(expressions, ", "),
			"expressionValues":      strings.Join(expressionValues, ", "),
			"keys":                  strings.Join(keySet.KeysStr(), "\n"),
			"keyValues":             strings.Join(keyVariableSet.KeysStr(), ", "),
		})
	if err != nil {
		return "", "", err
	}

	// interface method
	text, err = util.LoadTemplate(category, insertTemplateMethodFile, template.InsertMethod)
	if err != nil {
		return "", "", err
	}

	insertMethodOutput, err := util.With("insertMethod").Parse(text).Execute(map[string]interface{}{
		"upperStartCamelObject": camel,
	})
	if err != nil {
		return "", "", err
	}

	return output.String(), insertMethodOutput.String(), nil
}
