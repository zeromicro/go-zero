package gen

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

func TestGenCacheKeys(t *testing.T) {
	m, err := genCacheKeys(parser.Table{
		Name: stringx.From("user"),
		PrimaryKey: parser.Primary{
			Field: parser.Field{
				Name:         stringx.From("id"),
				DataBaseType: "bigint",
				DataType:     "int64",
				IsPrimaryKey: true,
				IsUniqueKey:  false,
				Comment:      "自增id",
			},
			AutoIncrement: true,
		},
		Fields: []parser.Field{
			{
				Name:         stringx.From("mobile"),
				DataBaseType: "varchar",
				DataType:     "string",
				IsPrimaryKey: false,
				IsUniqueKey:  true,
				Comment:      "手机号",
			},
			{
				Name:         stringx.From("name"),
				DataBaseType: "varchar",
				DataType:     "string",
				IsPrimaryKey: false,
				IsUniqueKey:  true,
				Comment:      "姓名",
			},
			{
				Name:         stringx.From("createTime"),
				DataBaseType: "timestamp",
				DataType:     "time.Time",
				IsPrimaryKey: false,
				IsUniqueKey:  false,
				Comment:      "创建时间",
			},
			{
				Name:         stringx.From("updateTime"),
				DataBaseType: "timestamp",
				DataType:     "time.Time",
				IsPrimaryKey: false,
				IsUniqueKey:  false,
				Comment:      "更新时间",
			},
		},
	})
	assert.Nil(t, err)

	for fieldName, key := range m {
		name := stringx.From(fieldName)
		assert.Equal(t, fmt.Sprintf(`cacheUser%sPrefix = "cache#User#%s#"`, name.ToCamel(), name.Untitle()), key.VarExpression)
		assert.Equal(t, fmt.Sprintf(`cacheUser%sPrefix`, name.ToCamel()), key.Left)
		assert.Equal(t, fmt.Sprintf(`cache#User#%s#`, name.Untitle()), key.Right)
		assert.Equal(t, fmt.Sprintf(`user%sKey`, name.ToCamel()), key.Variable)
		assert.Equal(t, `user`+name.ToCamel()+`Key := fmt.Sprintf("%s%v", cacheUser`+name.ToCamel()+`Prefix,`+name.Untitle()+`)`, key.KeyExpression)
	}
}
