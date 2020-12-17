package gen

import (
	"fmt"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

type findOneCode struct {
	findOneMethod          string
	findOneInterfaceMethod string
	cacheExtra             string
}

func genFindOneByField(table Table, withCache bool) (*findOneCode, error) {
	text, err := util.LoadTemplate(category, findOneByFieldTemplateFile, template.FindOneByField)
	if err != nil {
		return nil, err
	}

	t := util.With("findOneByField").Parse(text)
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
			"in":                        fmt.Sprintf("%s %s", stringx.From(camelFieldName).Untitle(), field.DataType),
			"withCache":                 withCache,
			"cacheKey":                  table.CacheKey[field.Name.Source()].KeyExpression,
			"cacheKeyVariable":          table.CacheKey[field.Name.Source()].Variable,
			"lowerStartCamelObject":     stringx.From(camelTableName).Untitle(),
			"lowerStartCamelField":      stringx.From(camelFieldName).Untitle(),
			"upperStartCamelPrimaryKey": table.PrimaryKey.Name.ToCamel(),
			"originalField":             wrapWithRawString(field.Name.Source()),
		})
		if err != nil {
			return nil, err
		}

		list = append(list, output.String())
	}

	text, err = util.LoadTemplate(category, findOneByFieldMethodTemplateFile, template.FindOneByFieldMethod)
	if err != nil {
		return nil, err
	}

	t = util.With("findOneByFieldMethod").Parse(text)
	var listMethod []string
	for _, field := range table.Fields {
		if field.IsPrimaryKey || !field.IsUniqueKey {
			continue
		}
		camelFieldName := field.Name.ToCamel()
		output, err := t.Execute(map[string]interface{}{
			"upperStartCamelObject": camelTableName,
			"upperField":            camelFieldName,
			"in":                    fmt.Sprintf("%s %s", stringx.From(camelFieldName).Untitle(), field.DataType),
		})
		if err != nil {
			return nil, err
		}

		listMethod = append(listMethod, output.String())
	}

	if withCache {
		text, err := util.LoadTemplate(category, findOneByFieldExtraMethodTemplateFile, template.FindOneByFieldExtraMethod)
		if err != nil {
			return nil, err
		}

		out, err := util.With("findOneByFieldExtraMethod").Parse(text).Execute(map[string]interface{}{
			"upperStartCamelObject": camelTableName,
			"primaryKeyLeft":        table.CacheKey[table.PrimaryKey.Name.Source()].Left,
			"lowerStartCamelObject": stringx.From(camelTableName).Untitle(),
			"originalPrimaryField":  wrapWithRawString(table.PrimaryKey.Name.Source()),
		})
		if err != nil {
			return nil, err
		}

		return &findOneCode{
			findOneMethod:          strings.Join(list, util.NL),
			findOneInterfaceMethod: strings.Join(listMethod, util.NL),
			cacheExtra:             out.String(),
		}, nil
	}

	return &findOneCode{
		findOneMethod:          strings.Join(list, util.NL),
		findOneInterfaceMethod: strings.Join(listMethod, util.NL),
	}, nil
}
