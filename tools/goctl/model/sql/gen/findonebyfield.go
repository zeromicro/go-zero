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

func genFindOneByField(table Table, withCache, postgreSql bool) (*findOneCode, error) {
	text, err := util.LoadTemplate(category, findOneByFieldTemplateFile, template.FindOneByField)
	if err != nil {
		return nil, err
	}

	t := util.With("findOneByField").Parse(text)
	var list []string
	camelTableName := table.Name.ToCamel()
	for _, key := range table.UniqueCacheKey {
		in, paramJoinString, originalFieldString := convertJoin(key, postgreSql)

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
			"postgreSql":                postgreSql,
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
			"originalPrimaryField":  wrapWithRawString(table.PrimaryKey.Name.Source(), postgreSql),
			"postgreSql":            postgreSql,
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

func convertJoin(key Key, postgreSql bool) (in, paramJoinString, originalFieldString string) {
	var inJoin, paramJoin, argJoin Join
	for index, f := range key.Fields {
		param := stringx.From(f.Name.ToCamel()).Untitle()
		inJoin = append(inJoin, fmt.Sprintf("%s %s", param, f.DataType))
		paramJoin = append(paramJoin, param)
		if postgreSql {
			argJoin = append(argJoin, fmt.Sprintf("%s = $%d", wrapWithRawString(f.Name.Source(), postgreSql), index+1))
		} else {
			argJoin = append(argJoin, fmt.Sprintf("%s = ?", wrapWithRawString(f.Name.Source(), postgreSql)))
		}
	}
	if len(inJoin) > 0 {
		in = inJoin.With(", ").Source()
	}

	if len(paramJoin) > 0 {
		paramJoinString = paramJoin.With(",").Source()
	}

	if len(argJoin) > 0 {
		originalFieldString = argJoin.With(" and ").Source()
	}
	return in, paramJoinString, originalFieldString
}
