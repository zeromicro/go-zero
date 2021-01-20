package gen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/model/sql/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

type Key struct {
	// cacheUserIdPrefix
	VarLeft string
	// "cache#user#id#"
	VarRight string
	// cacheUserIdPrefix = "cache#user#id#"
	VarExpression string

	// userKey
	KeyLeft string
	// fmt.Sprintf("%s%v", cacheUserPrefix, user)
	KeyRight string
	// fmt.Sprintf("%s%v", cacheUserPrefix, data.User)
	DataKeyRight string
	// userKey := fmt.Sprintf("%s%v", cacheUserPrefix, user)
	KeyExpression string
	// userKey := fmt.Sprintf("%s%v", cacheUserPrefix, data.User)
	DataKeyExpression string
	FieldNameJoin     Join
	Fields            []*parser.Field
}

type Join []string

func genCacheKeys(table parser.Table) (Key, []Key) {
	var primaryKey Key
	var uniqueKey []Key
	primaryKey = genCacheKey(table.Name, []*parser.Field{&table.PrimaryKey.Field})
	for _, each := range table.UniqueIndex {
		uniqueKey = append(uniqueKey, genCacheKey(table.Name, each))
	}
	sort.Slice(uniqueKey, func(i, j int) bool {
		return uniqueKey[i].VarLeft < uniqueKey[j].VarLeft
	})

	return primaryKey, uniqueKey
}

func genCacheKey(table stringx.String, in []*parser.Field) Key {
	var (
		varLeftJoin, varRightJon, fieldNameJoin Join
		varLeft, varRight, varExpression        string

		keyLeftJoin, keyRightJoin, keyRightArgJoin, dataRightJoin         Join
		keyLeft, keyRight, dataKeyRight, keyExpression, dataKeyExpression string
	)

	varLeftJoin = append(varLeftJoin, "cache", table.Source())
	varRightJon = append(varRightJon, "cache", table.Source())
	keyLeftJoin = append(keyLeftJoin, table.Source())

	for _, each := range in {
		varLeftJoin = append(varLeftJoin, each.Name.Source())
		varRightJon = append(varRightJon, each.Name.Source())
		keyLeftJoin = append(keyLeftJoin, each.Name.Source())
		keyRightJoin = append(keyRightJoin, stringx.From(each.Name.ToCamel()).Untitle())
		keyRightArgJoin = append(keyRightArgJoin, "%v")
		dataRightJoin = append(dataRightJoin, "data."+each.Name.ToCamel())
		fieldNameJoin = append(fieldNameJoin, each.Name.Source())
	}
	varLeftJoin = append(varLeftJoin, "prefix")
	keyLeftJoin = append(keyLeftJoin, "key")

	varLeft = varLeftJoin.Camel().With("").Untitle()
	varRight = fmt.Sprintf(`"%s"`, varRightJon.Camel().Untitle().With("#").Source()+"#")
	varExpression = fmt.Sprintf(`%s = %s`, varLeft, varRight)

	keyLeft = keyLeftJoin.Camel().With("").Untitle()
	keyRight = fmt.Sprintf(`fmt.Sprintf("%s%s", %s, %s)`, "%s", keyRightArgJoin.With("").Source(), varLeft, keyRightJoin.With(", ").Source())
	dataKeyRight = fmt.Sprintf(`fmt.Sprintf("%s%s", %s, %s)`, "%s", keyRightArgJoin.With("").Source(), varLeft, dataRightJoin.With(", ").Source())
	keyExpression = fmt.Sprintf("%s := %s", keyLeft, keyRight)
	dataKeyExpression = fmt.Sprintf("%s := %s", keyLeft, dataKeyRight)

	return Key{
		VarLeft:           varLeft,
		VarRight:          varRight,
		VarExpression:     varExpression,
		KeyLeft:           keyLeft,
		KeyRight:          keyRight,
		DataKeyRight:      dataKeyRight,
		KeyExpression:     keyExpression,
		DataKeyExpression: dataKeyExpression,
		Fields:            in,
		FieldNameJoin:     fieldNameJoin,
	}
}

func (j Join) Title() Join {
	var join Join
	for _, each := range j {
		join = append(join, stringx.From(each).Title())
	}

	return join
}

func (j Join) Camel() Join {
	var join Join
	for _, each := range j {
		join = append(join, stringx.From(each).ToCamel())
	}
	return join
}

func (j Join) Snake() Join {
	var join Join
	for _, each := range j {
		join = append(join, stringx.From(each).ToSnake())
	}

	return join
}

func (j Join) Untitle() Join {
	var join Join
	for _, each := range j {
		join = append(join, stringx.From(each).Untitle())
	}

	return join
}

func (j Join) Upper() Join {
	var join Join
	for _, each := range j {
		join = append(join, stringx.From(each).Upper())
	}

	return join
}

func (j Join) Lower() Join {
	var join Join
	for _, each := range j {
		join = append(join, stringx.From(each).Lower())
	}

	return join
}

func (j Join) With(sep string) stringx.String {
	return stringx.From(strings.Join(j, sep))
}
