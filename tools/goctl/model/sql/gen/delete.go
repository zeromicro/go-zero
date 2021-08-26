package gen

import (
	"strings"

	"github.com/tal-tech/go-zero/core/collection"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

func genDelete(table Table, withCache, postgreSql bool) (string, string, error) {
	keySet := collection.NewSet()
	keyVariableSet := collection.NewSet()
	keySet.AddStr(table.PrimaryCacheKey.KeyExpression)
	keyVariableSet.AddStr(table.PrimaryCacheKey.KeyLeft)
	for _, key := range table.UniqueCacheKey {
		keySet.AddStr(key.DataKeyExpression)
		keyVariableSet.AddStr(key.KeyLeft)
	}

	camel := table.Name.ToCamel()
	text, err := util.LoadTemplate(category, deleteTemplateFile, template.Delete)
	if err != nil {
		return "", "", err
	}

	output, err := util.With("delete").
		Parse(text).
		Execute(map[string]interface{}{
			"upperStartCamelObject":     camel,
			"withCache":                 withCache,
			"containsIndexCache":        table.ContainsUniqueCacheKey,
			"lowerStartCamelPrimaryKey": stringx.From(table.PrimaryKey.Name.ToCamel()).Untitle(),
			"dataType":                  table.PrimaryKey.DataType,
			"keys":                      strings.Join(keySet.KeysStr(), "\n"),
			"originalPrimaryKey":        wrapWithRawString(table.PrimaryKey.Name.Source(), postgreSql),
			"keyValues":                 strings.Join(keyVariableSet.KeysStr(), ", "),
			"postgreSql":                postgreSql,
		})
	if err != nil {
		return "", "", err
	}

	// interface method
	text, err = util.LoadTemplate(category, deleteMethodTemplateFile, template.DeleteMethod)
	if err != nil {
		return "", "", err
	}

	deleteMethodOut, err := util.With("deleteMethod").
		Parse(text).
		Execute(map[string]interface{}{
			"lowerStartCamelPrimaryKey": stringx.From(table.PrimaryKey.Name.ToCamel()).Untitle(),
			"dataType":                  table.PrimaryKey.DataType,
		})
	if err != nil {
		return "", "", err
	}

	return output.String(), deleteMethodOut.String(), nil
}
