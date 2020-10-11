package gen

import (
	"fmt"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

func genFindOneByField(table Table, withCache bool) (string, string, error) {
	t := util.With("findOneByField").Parse(template.FindOneByField)
	var list []string
	camelTableName := table.Name.ToCamel()
	for _, field := range table.Fields {
		if field.IsPrimaryKey || !field.IsUniqueKey {
			continue
		}
		camelFieldName := field.Name.ToCamel()
		output, err := t.Execute(map[string]interface{}{
			"upperStartCamelObject":     camelTableName,
			"upperField":                camelFieldName,
			"in":                        fmt.Sprintf("%s %s", stringx.From(camelFieldName).UnTitle(), field.DataType),
			"withCache":                 withCache,
			"cacheKey":                  table.CacheKey[field.Name.Source()].KeyExpression,
			"cacheKeyVariable":          table.CacheKey[field.Name.Source()].Variable,
			"primaryKeyLeft":            table.CacheKey[table.PrimaryKey.Name.Source()].Left,
			"lowerStartCamelObject":     stringx.From(camelTableName).UnTitle(),
			"lowerStartCamelField":      stringx.From(camelFieldName).UnTitle(),
			"upperStartCamelPrimaryKey": table.PrimaryKey.Name.ToCamel(),
			"originalField":             field.Name.Source(),
		})
		if err != nil {
			return "", "", err
		}
		list = append(list, output.String())
	}
	if withCache {
		out, err := util.With("findOneByFieldExtraMethod").Parse(template.FindOneByFieldExtraMethod).Execute(map[string]interface{}{
			"upperStartCamelObject": camelTableName,
			"lowerStartCamelObject": stringx.From(camelTableName).UnTitle(),
			"originalPrimaryField":  table.PrimaryKey.Name.Source(),
		})
		if err != nil {
			return "", "", err
		}
		return strings.Join(list, "\n"), out.String(), nil
	}
	return strings.Join(list, "\n"), "", nil

}
