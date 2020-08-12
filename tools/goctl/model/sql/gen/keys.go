package gen

import (
	"fmt"

	"github.com/tal-tech/go-zero/tools/goctl/model/sql/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

type (
	// tableName:user
	// {{prefix}}=cache
	// key:id
	Key struct {
		VarExpression     string // cacheUserIdPrefix="cache#user#id#"
		Left              string // cacheUserIdPrefix
		Right             string // cache#user#id#
		Variable          string // userIdKey
		KeyExpression     string // userIdKey: = fmt.Sprintf("cache#user#id#%v", userId)
		DataKeyExpression string // userIdKey: = fmt.Sprintf("cache#user#id#%v", data.userId)
		RespKeyExpression string // userIdKey: = fmt.Sprintf("cache#user#id#%v", resp.userId)
	}
)

// key-数据库原始字段名,value-缓存key相关数据
func genCacheKeys(table parser.Table) (map[string]Key, error) {
	fields := table.Fields
	m := make(map[string]Key)
	camelTableName := table.Name.Snake2Camel()
	lowerStartCamelTableName := stringx.From(camelTableName).LowerStart()
	for _, field := range fields {
		if !field.IsKey {
			continue
		}
		camelFieldName := field.Name.Snake2Camel()
		lowerStartCamelFieldName := stringx.From(camelFieldName).LowerStart()
		left := fmt.Sprintf("cache%s%sPrefix", lowerStartCamelTableName, camelFieldName)
		right := fmt.Sprintf("cache#%s#%s#", lowerStartCamelTableName, lowerStartCamelFieldName)
		variable := fmt.Sprintf("%s%sKey", lowerStartCamelTableName, camelFieldName)
		m[field.Name.Source()] = Key{
			VarExpression:     fmt.Sprintf(`%s = "%s"`, left, right),
			Left:              left,
			Right:             right,
			Variable:          variable,
			KeyExpression:     fmt.Sprintf(`%s := fmt.Sprintf("cache#%s#%s#%s", %s)`, variable, lowerStartCamelTableName, lowerStartCamelFieldName, "%v", lowerStartCamelFieldName),
			DataKeyExpression: fmt.Sprintf(`%s := fmt.Sprintf("cache#%s#%s#%s", data.%s)`, variable, lowerStartCamelTableName, lowerStartCamelFieldName, "%v", camelFieldName),
			RespKeyExpression: fmt.Sprintf(`%s := fmt.Sprintf("cache#%s#%s#%s", resp.%s)`, variable, lowerStartCamelTableName, lowerStartCamelFieldName, "%v", camelFieldName),
		}
	}
	return m, nil
}
