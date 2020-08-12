package gen

import (
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/util/templatex"
)

func genFindOne(table Table, withCache bool) (string, error) {
	camel := table.Name.Snake2Camel()
	output, err := templatex.With("findOne").
		Parse(template.FindOne).
		Execute(map[string]interface{}{
			"withCache":                 withCache,
			"upperStartCamelObject":     camel,
			"lowerStartCamelObject":     stringx.From(camel).LowerStart(),
			"originalPrimaryKey":        table.PrimaryKey.Name.Source(),
			"lowerStartCamelPrimaryKey": stringx.From(table.PrimaryKey.Name.Snake2Camel()).LowerStart(),
			"dataType":                  table.PrimaryKey.DataType,
			"cacheKey":                  table.CacheKey[table.PrimaryKey.Name.Source()].KeyExpression,
			"cacheKeyVariable":          table.CacheKey[table.PrimaryKey.Name.Source()].Variable,
		})
	if err != nil {
		return "", err
	}
	return output.String(), nil
}
