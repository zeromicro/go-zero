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

func genDefine(table Table, withCache, postgreSql bool) (string, error) {
	expressionValues := make([]string, 0)
	fields := make([]Field, 0)
	for _, field := range table.Fields {
		camel := util.SafeString(field.Name.ToCamel())
		if camel == "CreateTime" || camel == "UpdateTime" {
			continue
		}

		if field.Name.Source() == table.PrimaryKey.Name.Source() {
			continue
		}

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
		expressionValues = append(expressionValues, "data."+camel)
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
		expressionValues = append([]string{"data." + table.PrimaryKey.Name.ToCamel()}, expressionValues...)
	} else {
		expressionValues = append(expressionValues, "data."+table.PrimaryKey.Name.ToCamel())
	}
	camelTableName := table.Name.ToCamel()
	text, err := pathx.LoadTemplate(category, defineTemplateFile, template.Define)
	if err != nil {
		return "", err
	}

	t := util.With("define").
		Parse(text).
		GoFmt(true)
	output, err := t.Execute(map[string]interface{}{
		"withCache":                 withCache,
		"containsIndexCache":        table.ContainsUniqueCacheKey,
		"upperStartCamelObject":     camelTableName,
		"lowerStartCamelObject":     stringx.From(camelTableName).Untitle(),
		"upperStartCamelPrimaryKey": table.PrimaryKey.Name.ToCamel(),
		"lowerStartCamelPrimaryKey": stringx.From(table.PrimaryKey.Name.ToCamel()).Untitle(),
		"primaryKeyDataType":        table.PrimaryKey.DataType,
		"originalPrimaryKey":        wrapWithRawString(table.PrimaryKey.Name.Source(), postgreSql),
		"primaryCacheKey":           table.PrimaryCacheKey.DataKeyExpression,
		"primaryKeyVariable":        table.PrimaryCacheKey.KeyLeft,
		"keys":                      strings.Join(keys, "\n"),
		"keyValues":                 strings.Join(keyVars, ", "),
		"fields":                    fields,
		"expressionValues":          strings.Join(expressionValues, ", "),
		"postgreSql":                postgreSql,
		"data":                      table,
	})
	if err != nil {
		return "", nil
	}

	return output.String(), nil
}
