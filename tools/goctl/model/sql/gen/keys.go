package gen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/model/sql/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"
)

// Key describes cache key
type Key struct {
	// VarLeft describes the variable of cache key expression which likes cacheUserIdPrefix
	VarLeft string
	// VarRight describes the value of cache key expression which likes "cache:user:id:"
	VarRight string
	// VarExpression describes the cache key expression which likes cacheUserIdPrefix = "cache:user:id:"
	VarExpression string
	// KeyLeft describes the variable of key definition expression which likes userKey
	KeyLeft string
	// KeyRight describes the value of key definition expression which likes fmt.Sprintf("%s%v", cacheUserPrefix, user)
	KeyRight string
	// DataKeyRight describes data key likes fmt.Sprintf("%s%v", cacheUserPrefix, data.User)
	DataKeyRight string
	// KeyExpression describes key expression likes userKey := fmt.Sprintf("%s%v", cacheUserPrefix, user)
	KeyExpression string
	// DataKeyExpression describes data key expression likes userKey := fmt.Sprintf("%s%v", cacheUserPrefix, data.User)
	DataKeyExpression string
	// FieldNameJoin describes the filed slice of table
	FieldNameJoin Join
	// Fields describes the fields of table
	Fields []*parser.Field
}

// Join describes an alias of string slice
type Join []string

func genCacheKeys(prefix string, table parser.Table) (Key, []Key) {
	var primaryKey Key
	var uniqueKey []Key
	primaryKey = genCacheKey(prefix, table.Db, table.Name, []*parser.Field{&table.PrimaryKey.Field})
	for _, each := range table.UniqueIndex {
		uniqueKey = append(uniqueKey, genCacheKey(prefix, table.Db, table.Name, each))
	}
	sort.Slice(uniqueKey, func(i, j int) bool {
		return uniqueKey[i].VarLeft < uniqueKey[j].VarLeft
	})

	return primaryKey, uniqueKey
}

func genCacheKey(prefix string, db, table stringx.String, in []*parser.Field) Key {
	var (
		varLeftJoin, varRightJoin, fieldNameJoin Join
		varLeft, varRight, varExpression         string

		keyLeftJoin, keyRightJoin, keyRightArgJoin, dataRightJoin         Join
		keyLeft, keyRight, dataKeyRight, keyExpression, dataKeyExpression string
	)

	dbName, tableName := util.SafeString(db.Source()), util.SafeString(table.Source())
	if len(dbName) > 0 {
		varLeftJoin = append(varLeftJoin, prefix, dbName, tableName)
		varRightJoin = append(varRightJoin, prefix, dbName, tableName)
		keyLeftJoin = append(keyLeftJoin, dbName, tableName)
	} else {
		varLeftJoin = append(varLeftJoin, prefix, tableName)
		varRightJoin = append(varRightJoin, prefix, tableName)
		keyLeftJoin = append(keyLeftJoin, tableName)
	}

	for _, each := range in {
		varLeftJoin = append(varLeftJoin, each.Name.Source())
		varRightJoin = append(varRightJoin, each.Name.Source())
		keyLeftJoin = append(keyLeftJoin, each.Name.Source())
		keyRightJoin = append(keyRightJoin, util.EscapeGolangKeyword(stringx.From(each.Name.ToCamel()).Untitle()))
		keyRightArgJoin = append(keyRightArgJoin, "%v")
		dataRightJoin = append(dataRightJoin, "data."+each.Name.ToCamel())
		fieldNameJoin = append(fieldNameJoin, each.Name.Source())
	}
	varLeftJoin = append(varLeftJoin, "prefix")
	keyLeftJoin = append(keyLeftJoin, "key")

	varLeft = util.SafeString(varLeftJoin.Camel().With("").Untitle())
	varRight = fmt.Sprintf(`"%s"`, varRightJoin.Camel().Untitle().With(":").Source()+":")
	varExpression = fmt.Sprintf(`%s = %s`, varLeft, varRight)

	keyLeft = util.SafeString(keyLeftJoin.Camel().With("").Untitle())
	keyRight = fmt.Sprintf(`fmt.Sprintf("%s%s", %s, %s)`, "%s", keyRightArgJoin.With(":").Source(), varLeft, keyRightJoin.With(", ").Source())
	dataKeyRight = fmt.Sprintf(`fmt.Sprintf("%s%s", %s, %s)`, "%s", keyRightArgJoin.With(":").Source(), varLeft, dataRightJoin.With(", ").Source())
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

// Title convert items into Title and return
func (j Join) Title() Join {
	var join Join
	for _, each := range j {
		join = append(join, stringx.From(each).Title())
	}

	return join
}

// Camel convert items into Camel and return
func (j Join) Camel() Join {
	var join Join
	for _, each := range j {
		join = append(join, stringx.From(each).ToCamel())
	}
	return join
}

// Snake convert items into Snake and return
func (j Join) Snake() Join {
	var join Join
	for _, each := range j {
		join = append(join, stringx.From(each).ToSnake())
	}

	return join
}

// Untitle converts items into Untitle and return
func (j Join) Untitle() Join {
	var join Join
	for _, each := range j {
		join = append(join, stringx.From(each).Untitle())
	}

	return join
}

// Upper convert items into Upper and return
func (j Join) Upper() Join {
	var join Join
	for _, each := range j {
		join = append(join, stringx.From(each).Upper())
	}

	return join
}

// Lower convert items into Lower and return
func (j Join) Lower() Join {
	var join Join
	for _, each := range j {
		join = append(join, stringx.From(each).Lower())
	}

	return join
}

// With convert items into With and return
func (j Join) With(sep string) stringx.String {
	return stringx.From(strings.Join(j, sep))
}
