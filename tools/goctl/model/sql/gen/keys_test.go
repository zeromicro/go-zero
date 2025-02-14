package gen

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"
)

func TestGenCacheKeys(t *testing.T) {
	primaryField := &parser.Field{
		Name:       stringx.From("id"),
		DataType:   "int64",
		Comment:    "自增id",
		SeqInIndex: 1,
	}
	mobileField := &parser.Field{
		Name:       stringx.From("mobile"),
		DataType:   "string",
		Comment:    "手机号",
		SeqInIndex: 1,
	}
	classField := &parser.Field{
		Name:       stringx.From("class"),
		DataType:   "string",
		Comment:    "班级",
		SeqInIndex: 1,
	}
	nameField := &parser.Field{
		Name:       stringx.From("name"),
		DataType:   "string",
		Comment:    "姓名",
		SeqInIndex: 2,
	}
	primariCacheKey, uniqueCacheKey := genCacheKeys("cache", parser.Table{
		Name: stringx.From("user"),
		Db:   stringx.From("go_zero"),
		PrimaryKey: parser.Primary{
			Field:         *primaryField,
			AutoIncrement: true,
		},
		UniqueIndex: map[string][]*parser.Field{
			"mobile_unique": {
				mobileField,
			},
			"class_name_unique": {
				classField,
				nameField,
			},
		},
		Fields: []*parser.Field{
			primaryField,
			mobileField,
			classField,
			nameField,
			{
				Name:     stringx.From("createTime"),
				DataType: "time.Time",
				Comment:  "创建时间",
			},
			{
				Name:     stringx.From("updateTime"),
				DataType: "time.Time",
				Comment:  "更新时间",
			},
		},
	})

	t.Run("primaryCacheKey", func(t *testing.T) {
		assert.Equal(t, true, func() bool {
			return cacheKeyEqual(primariCacheKey, Key{
				VarLeft:           "cacheGoZeroUserIdPrefix",
				VarRight:          `"cache:goZero:user:id:"`,
				VarExpression:     `cacheGoZeroUserIdPrefix = "cache:goZero:user:id:"`,
				KeyLeft:           "goZeroUserIdKey",
				KeyRight:          `fmt.Sprintf("%s%v", cacheGoZeroUserIdPrefix, id)`,
				DataKeyRight:      `fmt.Sprintf("%s%v", cacheGoZeroUserIdPrefix, data.Id)`,
				KeyExpression:     `goZeroUserIdKey := fmt.Sprintf("%s%v", cacheGoZeroUserIdPrefix, id)`,
				DataKeyExpression: `goZeroUserIdKey := fmt.Sprintf("%s%v", cacheGoZeroUserIdPrefix, data.Id)`,
				FieldNameJoin:     []string{"id"},
			})
		}())
	})

	t.Run("uniqueCacheKey", func(t *testing.T) {
		assert.Equal(t, true, func() bool {
			expected := []Key{
				{
					VarLeft:           "cacheGoZeroUserClassNamePrefix",
					VarRight:          `"cache:goZero:user:class:name:"`,
					VarExpression:     `cacheGoZeroUserClassNamePrefix = "cache:goZero:user:class:name:"`,
					KeyLeft:           "goZeroUserClassNameKey",
					KeyRight:          `fmt.Sprintf("%s%v:%v", cacheGoZeroUserClassNamePrefix, class, name)`,
					DataKeyRight:      `fmt.Sprintf("%s%v:%v", cacheGoZeroUserClassNamePrefix, data.Class, data.Name)`,
					KeyExpression:     `goZeroUserClassNameKey := fmt.Sprintf("%s%v:%v", cacheGoZeroUserClassNamePrefix, class, name)`,
					DataKeyExpression: `goZeroUserClassNameKey := fmt.Sprintf("%s%v:%v", cacheGoZeroUserClassNamePrefix, data.Class, data.Name)`,
					FieldNameJoin:     []string{"class", "name"},
				},
				{
					VarLeft:           "cacheGoZeroUserMobilePrefix",
					VarRight:          `"cache:goZero:user:mobile:"`,
					VarExpression:     `cacheGoZeroUserMobilePrefix = "cache:goZero:user:mobile:"`,
					KeyLeft:           "goZeroUserMobileKey",
					KeyRight:          `fmt.Sprintf("%s%v", cacheGoZeroUserMobilePrefix, mobile)`,
					DataKeyRight:      `fmt.Sprintf("%s%v", cacheGoZeroUserMobilePrefix, data.Mobile)`,
					KeyExpression:     `goZeroUserMobileKey := fmt.Sprintf("%s%v", cacheGoZeroUserMobilePrefix, mobile)`,
					DataKeyExpression: `goZeroUserMobileKey := fmt.Sprintf("%s%v", cacheGoZeroUserMobilePrefix, data.Mobile)`,
					FieldNameJoin:     []string{"mobile"},
				},
			}
			sort.Slice(uniqueCacheKey, func(i, j int) bool {
				return uniqueCacheKey[i].VarLeft < uniqueCacheKey[j].VarLeft
			})

			if len(expected) != len(uniqueCacheKey) {
				return false
			}

			for index, each := range uniqueCacheKey {
				expecting := expected[index]
				if !cacheKeyEqual(expecting, each) {
					return false
				}
			}

			return true
		}())
	})
	t.Run("no database name", func(t *testing.T) {
		primariCacheKey, _ = genCacheKeys("cache", parser.Table{
			Name: stringx.From("user"),
			Db:   stringx.From(""),
			PrimaryKey: parser.Primary{
				Field:         *primaryField,
				AutoIncrement: true,
			},
			UniqueIndex: map[string][]*parser.Field{
				"mobile_unique": {
					mobileField,
				},
				"class_name_unique": {
					classField,
					nameField,
				},
			},
			Fields: []*parser.Field{
				primaryField,
				mobileField,
				classField,
				nameField,
				{
					Name:     stringx.From("createTime"),
					DataType: "time.Time",
					Comment:  "创建时间",
				},
				{
					Name:     stringx.From("updateTime"),
					DataType: "time.Time",
					Comment:  "更新时间",
				},
			},
		})

		assert.Equal(t, true, func() bool {
			return cacheKeyEqual(primariCacheKey, Key{
				VarLeft:           "cacheUserIdPrefix",
				VarRight:          `"cache:user:id:"`,
				VarExpression:     `cacheUserIdPrefix = "cache:user:id:"`,
				KeyLeft:           "userIdKey",
				KeyRight:          `fmt.Sprintf("%s%v", cacheUserIdPrefix, id)`,
				DataKeyRight:      `fmt.Sprintf("%s%v", cacheUserIdPrefix, data.Id)`,
				KeyExpression:     `userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)`,
				DataKeyExpression: `userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, data.Id)`,
				FieldNameJoin:     []string{"id"},
			})
		}())
	})
}

func cacheKeyEqual(k1, k2 Key) bool {
	k1Join := k1.FieldNameJoin
	k2Join := k2.FieldNameJoin
	sort.Strings(k1Join)
	sort.Strings(k2Join)
	if len(k1Join) != len(k2Join) {
		return false
	}

	for index, each := range k1Join {
		k2Item := k2Join[index]
		if each != k2Item {
			return false
		}
	}

	return k1.VarLeft == k2.VarLeft &&
		k1.VarRight == k2.VarRight &&
		k1.VarExpression == k2.VarExpression &&
		k1.KeyLeft == k2.KeyLeft &&
		k1.KeyRight == k2.KeyRight &&
		k1.DataKeyRight == k2.DataKeyRight &&
		k1.DataKeyExpression == k2.DataKeyExpression &&
		k1.KeyExpression == k2.KeyExpression
}
