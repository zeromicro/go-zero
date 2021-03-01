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
	for _, key := range table.UniqueCacheKey {
		var inJoin, paramJoin, argJoin Join
		for _, f := range key.Fields {
			param := stringx.From(f.Name.ToCamel()).Untitle()
			inJoin = append(inJoin, fmt.Sprintf("%s %s", param, f.DataType))
			paramJoin = append(paramJoin, param)
			argJoin = append(argJoin, fmt.Sprintf("%s = ?", wrapWithRawString(f.Name.Source())))
		}
		var in string
		if len(inJoin) > 0 {
			in = inJoin.With(", ").Source()
		}

		var paramJoinString string
		if len(paramJoin) > 0 {
			paramJoinString = paramJoin.With(",").Source()
		}

		var originalFieldString string
		if len(argJoin) > 0 {
			originalFieldString = argJoin.With(" and ").Source()
		}

		output, err := t.Execute(map[string]interface{}{
			"upperStartCamelObject":     camelTableName,
			"upperField":                key.FieldNameJoin.Camel().With("").Source(),
			"in":                        in,
			"withCache":                 withCache,
			"cacheKey":                  key.KeyExpression,
			"cacheKeyVariable":          key.KeyLeft,
			"lowerStartCamelObject":     stringx.From(camelTableName).Untitle(),
			"lowerStartCamelField":      paramJoinString,
			"upperStartCamelPrimaryKey": table.PrimaryKey.Name.ToCamel(),
			"originalField":             originalFieldString,
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
	for _, key := range table.UniqueCacheKey {
		var inJoin, paramJoin Join
		for _, f := range key.Fields {
			param := stringx.From(f.Name.ToCamel()).Untitle()
			inJoin = append(inJoin, fmt.Sprintf("%s %s", param, f.DataType))
			paramJoin = append(paramJoin, param)
		}

		var in string
		if len(inJoin) > 0 {
			in = inJoin.With(", ").Source()
		}
		output, err := t.Execute(map[string]interface{}{
			"upperStartCamelObject": camelTableName,
			"upperField":            key.FieldNameJoin.Camel().With("").Source(),
			"in":                    in,
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
			"primaryKeyLeft":        table.PrimaryCacheKey.VarLeft,
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
