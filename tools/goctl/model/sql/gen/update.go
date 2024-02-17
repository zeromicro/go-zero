package gen

import (
	"sort"
	"strings"

	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/template"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"
)

func genUpdate(table Table, withCache, postgreSql bool) (
	string, string, error,
) {
	expressionValues := make([]string, 0)
	pkg := "data."
	if table.ContainsUniqueCacheKey {
		pkg = "newData."
	}
	for _, field := range table.Fields {
		camel := util.SafeString(field.Name.ToCamel())
		if table.isIgnoreColumns(field.Name.Source()) {
			continue
		}

		if field.Name.Source() == table.PrimaryKey.Name.Source() {
			continue
		}

		expressionValues = append(expressionValues, pkg+camel)
	}

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

	if postgreSql {
		expressionValues = append(
			[]string{pkg + table.PrimaryKey.Name.ToCamel()},
			expressionValues...,
		)
	} else {
		expressionValues = append(
			expressionValues, pkg+table.PrimaryKey.Name.ToCamel(),
		)
	}
	camelTableName := table.Name.ToCamel()
	text, err := pathx.LoadTemplate(category, updateTemplateFile, template.Update)
	if err != nil {
		return "", "", err
	}

	output, err := util.With("update").Parse(text).Execute(
		map[string]any{
			"withCache":             withCache,
			"containsIndexCache":    table.ContainsUniqueCacheKey,
			"upperStartCamelObject": camelTableName,
			"keys":                  strings.Join(keys, "\n"),
			"keyValues":             strings.Join(keyVars, ", "),
			"primaryCacheKey":       table.PrimaryCacheKey.DataKeyExpression,
			"primaryKeyVariable":    table.PrimaryCacheKey.KeyLeft,
			"lowerStartCamelObject": stringx.From(camelTableName).Untitle(),
			"upperStartCamelPrimaryKey": util.EscapeGolangKeyword(
				stringx.From(table.PrimaryKey.Name.ToCamel()).Title(),
			),
			"originalPrimaryKey": wrapWithRawString(
				table.PrimaryKey.Name.Source(), postgreSql,
			),
			"expressionValues": strings.Join(
				expressionValues, ", ",
			),
			"postgreSql": postgreSql,
			"data":       table,
		},
	)
	if err != nil {
		return "", "", nil
	}

	// update interface method
	text, err = pathx.LoadTemplate(category, updateMethodTemplateFile, template.UpdateMethod)
	if err != nil {
		return "", "", err
	}

	updateMethodOutput, err := util.With("updateMethod").Parse(text).Execute(
		map[string]any{
			"upperStartCamelObject": camelTableName,
			"data":                  table,
		},
	)
	if err != nil {
		return "", "", nil
	}

	return output.String(), updateMethodOutput.String(), nil
}
