package gen

import (
	"fmt"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/util/templatex"
)

func genFineOneByField(table Table, withCache bool) (string, error) {
	t := templatex.With("findOneByField").Parse(template.FindOneByField)
	var list []string
	camelTableName := table.Name.Snake2Camel()
	for _, field := range table.Fields {
		if field.IsPrimaryKey || !field.IsKey {
			continue
		}
		camelFieldName := field.Name.Snake2Camel()
		output, err := t.Execute(map[string]interface{}{
			"upperStartCamelObject":     camelTableName,
			"upperField":                camelFieldName,
			"in":                        fmt.Sprintf("%s %s", stringx.From(camelFieldName).LowerStart(), field.DataType),
			"withCache":                 withCache,
			"cacheKey":                  table.CacheKey[field.Name.Source()].KeyExpression,
			"cacheKeyVariable":          table.CacheKey[field.Name.Source()].Variable,
			"primaryKeyLeft":            table.CacheKey[table.PrimaryKey.Name.Source()].Left,
			"lowerStartCamelObject":     stringx.From(camelTableName).LowerStart(),
			"lowerStartCamelField":      stringx.From(camelFieldName).LowerStart(),
			"upperStartCamelPrimaryKey": table.PrimaryKey.Name.Snake2Camel(),
			"originalField":             field.Name.Source(),
			"originalPrimaryField":      table.PrimaryKey.Name.Source(),
		})
		if err != nil {
			return "", err
		}
		list = append(list, output.String())
	}
	return strings.Join(list, "\n"), nil
}
