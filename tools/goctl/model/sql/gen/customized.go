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

// Field describes a table field
type Field struct {
	NameOriginal    string
	UpperName       string
	LowerName       string
	DataType        string
	Comment         string
	SeqInIndex      int
	OrdinalPosition int
}

func genCustomized(table Table, withCache, postgreSql bool) (string, error) {
	expressions := make([]string, 0)
	expressionValues := make([]string, 0)
	fields := make([]Field, 0)
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

		f := Field{
			NameOriginal:    field.NameOriginal,
			UpperName:       camel,
			LowerName:       stringx.From(camel).Untitle(),
			DataType:        field.DataType,
			Comment:         field.Comment,
			SeqInIndex:      field.SeqInIndex,
			OrdinalPosition: field.OrdinalPosition,
		}
		fields = append(fields, f)
	}

	keySet := collection.NewSet()
	keyVariableSet := collection.NewSet()
	keySet.AddStr(table.PrimaryCacheKey.KeyExpression)
	keyVariableSet.AddStr(table.PrimaryCacheKey.KeyLeft)
	for _, key := range table.UniqueCacheKey {
		keySet.AddStr(key.DataKeyExpression)
		keyVariableSet.AddStr(key.KeyLeft)
	}
	keys := keySet.KeysStr()
	sort.Strings(keys)
	keyVars := keyVariableSet.KeysStr()
	sort.Strings(keyVars)

	camel := table.Name.ToCamel()
	text, err := pathx.LoadTemplate(category, customizedTemplateFile, template.Customized)
	if err != nil {
		return "", err
	}

	output, err := util.With("customized").
		Parse(text).
		Execute(map[string]any{
			"withCache":                 withCache,
			"containsIndexCache":        table.ContainsUniqueCacheKey,
			"upperStartCamelObject":     camel,
			"lowerStartCamelObject":     stringx.From(camel).Untitle(),
			"lowerStartCamelPrimaryKey": util.EscapeGolangKeyword(stringx.From(table.PrimaryKey.Name.ToCamel()).Untitle()),
			"upperStartCamelPrimaryKey": table.PrimaryKey.Name.ToCamel(),
			"primaryKeyDataType":        table.PrimaryKey.DataType,
			"originalPrimaryKey":        wrapWithRawString(table.PrimaryKey.Name.Source(), postgreSql),
			"primaryCacheKey":           table.PrimaryCacheKey.DataKeyExpression,
			"primaryKeyVariable":        table.PrimaryCacheKey.KeyLeft,
			"keys":                      strings.Join(keys, "\n"),
			"keyValues":                 strings.Join(keyVars, ", "),
			"expression":                strings.Join(expressions, ", "),
			"expressionValues":          strings.Join(expressionValues, ", "),
			"postgreSql":                postgreSql,
			"fields":                    fields,
			"data":                      table,
		})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
