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

func genVars(table Table, withCache, postgreSql bool) (string, error) {
	keys := make([]string, 0)
	keys = append(keys, table.PrimaryCacheKey.VarExpression)
	for _, v := range table.UniqueCacheKey {
		keys = append(keys, v.VarExpression)
	}

	camel := table.Name.ToCamel()
	text, err := pathx.LoadTemplate(category, varTemplateFile, template.Vars)
	if err != nil {
		return "", err
	}

	output, err := util.With("var").Parse(text).
		GoFmt(true).Execute(map[string]any{
		"lowerStartCamelObject": stringx.From(camel).Untitle(),
		"upperStartCamelObject": camel,
		"cacheKeys":             strings.Join(keys, "\n"),
		"autoIncrement":         table.PrimaryKey.AutoIncrement,
		"originalPrimaryKey":    wrapWithRawString(table.PrimaryKey.Name.Source(), postgreSql),
		"withCache":             withCache,
		"postgreSql":            postgreSql,
		"data":                  table,
		"ignoreColumns": func() string {
			var set = collection.NewSet()
			for _, c := range table.ignoreColumns {
				if postgreSql {
					set.AddStr(fmt.Sprintf(`"%s"`, c))
				} else {
					set.AddStr(fmt.Sprintf("\"`%s`\"", c))
				}
			}
			list := set.KeysStr()
			sort.Strings(list)
			return strings.Join(list, ", ")
		}(),
	})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
