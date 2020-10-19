package gen

import (
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

func genVars(table Table, withCache bool) (string, error) {
	keys := make([]string, 0)
	for _, v := range table.CacheKey {
		keys = append(keys, v.VarExpression)
	}
	camel := table.Name.ToCamel()
	output, err := util.With("var").
		Parse(template.Vars).
		GoFmt(true).
		Execute(map[string]interface{}{
			"lowerStartCamelObject": stringx.From(camel).UnTitle(),
			"upperStartCamelObject": camel,
			"cacheKeys":             strings.Join(keys, "\n"),
			"autoIncrement":         table.PrimaryKey.AutoIncrement,
			"originalPrimaryKey":    table.PrimaryKey.Name.Source(),
			"withCache":             withCache,
		})
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
