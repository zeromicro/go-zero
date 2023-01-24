package gen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/template"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"
)

func genInsert(table Table, withCache, postgreSql bool) (string, string, error) {
	keySet := collection.NewSet()
	keyVariableSet := collection.NewSet()
	keySet.AddStr(table.PrimaryCacheKey.DataKeyExpression)
	keyVariableSet.AddStr(table.PrimaryCacheKey.KeyLeft)
	for _, key := range table.UniqueCacheKey {
		keySet.AddStr(key.DataKeyExpression)
		keyVariableSet.AddStr(key.KeyLeft)
	}
	keys := keySet.KeysStr()
	sort.Strings(keys)
	keyVars := keyVariableSet.KeysStr()
	sort.Strings(keyVars)

	expressions := make([]string, 0)
	expressionValues := make([]string, 0)
	var count int
	for _, field := range table.Fields {
		camel := util.SafeString(field.Name.ToCamel())
		if table.isIgnoreColumns(field.Name.Source()) {
			continue
		}

		if field.Name.Source() == table.PrimaryKey.Name.Source() {
			if table.PrimaryKey.AutoIncrement {
				continue
			}
		}

		count += 1
		if postgreSql {
			expressions = append(expressions, fmt.Sprintf("$%d", count))
		} else {
			expressions = append(expressions, "?")
		}
		expressionValues = append(expressionValues, "data."+camel)
	}

	camel := table.Name.ToCamel()
	text, err := pathx.LoadTemplate(category, insertTemplateFile, template.Insert)
	if err != nil {
		return "", "", err
	}

	output, err := util.With("insert").
		Parse(text).
		Execute(map[string]any{
			"withCache":             withCache,
			"upperStartCamelObject": camel,
			"lowerStartCamelObject": stringx.From(camel).Untitle(),
			"expression":            strings.Join(expressions, ", "),
			"expressionValues":      strings.Join(expressionValues, ", "),
			"keys":                  strings.Join(keys, "\n"),
			"keyValues":             strings.Join(keyVars, ", "),
			"data":                  table,
		})
	if err != nil {
		return "", "", err
	}

	// interface method
	text, err = pathx.LoadTemplate(category, insertTemplateMethodFile, template.InsertMethod)
	if err != nil {
		return "", "", err
	}

	insertMethodOutput, err := util.With("insertMethod").Parse(text).Execute(map[string]any{
		"upperStartCamelObject": camel,
		"data":                  table,
	})
	if err != nil {
		return "", "", err
	}

	return output.String(), insertMethodOutput.String(), nil
}
