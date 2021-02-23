package gen

import (
	"fmt"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/model/sql/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

// Key defines cache key variable for generating code
type Key struct {
	// VarExpression likes cacheUserIdPrefix = "cache#User#id#"
	VarExpression string
	// Left likes cacheUserIdPrefix
	Left string
	// Right likes  cache#user#id#
	Right string
	// Variable likes userIdKey
	Variable string
	// KeyExpression likes userIdKey: = fmt.Sprintf("cache#user#id#%v", userId)
	KeyExpression string
	// DataKeyExpression likes userIdKey: = fmt.Sprintf("cache#user#id#%v", data.userId)
	DataKeyExpression string
	// RespKeyExpression likes userIdKey: = fmt.Sprintf("cache#user#id#%v", resp.userId)
	RespKeyExpression string
}

// key-数据库原始字段名,value-缓存key相关数据
func genCacheKeys(table parser.Table) (map[string]Key, error) {
	fields := table.Fields
	m := make(map[string]Key)
	camelTableName := table.Name.ToCamel()
	lowerStartCamelTableName := stringx.From(camelTableName).Untitle()
	for _, field := range fields {
		if field.IsUniqueKey || field.IsPrimaryKey {
			camelFieldName := field.Name.ToCamel()
			lowerStartCamelFieldName := stringx.From(camelFieldName).Untitle()
			left := fmt.Sprintf("cache%s%sPrefix", camelTableName, camelFieldName)
			if strings.ToLower(camelFieldName) == strings.ToLower(camelTableName) {
				left = fmt.Sprintf("cache%sPrefix", camelTableName)
			}
			right := fmt.Sprintf("cache#%s#%s#", camelTableName, lowerStartCamelFieldName)
			variable := fmt.Sprintf("%s%sKey", lowerStartCamelTableName, camelFieldName)
			if strings.ToLower(lowerStartCamelTableName) == strings.ToLower(camelFieldName) {
				variable = fmt.Sprintf("%sKey", lowerStartCamelTableName)
			}

			m[field.Name.Source()] = Key{
				VarExpression:     fmt.Sprintf(`%s = "%s"`, left, right),
				Left:              left,
				Right:             right,
				Variable:          variable,
				KeyExpression:     fmt.Sprintf(`%s := fmt.Sprintf("%s%s", %s,%s)`, variable, "%s", "%v", left, lowerStartCamelFieldName),
				DataKeyExpression: fmt.Sprintf(`%s := fmt.Sprintf("%s%s",%s, data.%s)`, variable, "%s", "%v", left, camelFieldName),
				RespKeyExpression: fmt.Sprintf(`%s := fmt.Sprintf("%s%s", %s,resp.%s)`, variable, "%s", "%v", left, camelFieldName),
			}
		}
	}

	return m, nil
}
