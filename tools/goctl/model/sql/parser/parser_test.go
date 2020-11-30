package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/model"
)

func TestParsePlainText(t *testing.T) {
	_, err := Parse("plain text")
	assert.NotNil(t, err)
}

func TestParseSelect(t *testing.T) {
	_, err := Parse("select * from user")
	assert.Equal(t, unSupportDDL, err)
}

func TestParseCreateTable(t *testing.T) {
	table, err := Parse("CREATE TABLE `user_snake` (\n  `id` bigint(10) NOT NULL AUTO_INCREMENT,\n  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名称',\n  `password` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户密码',\n  `mobile` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '手机号',\n  `gender` char(5) COLLATE utf8mb4_general_ci NOT NULL COMMENT '男｜女｜未公开',\n  `nickname` varchar(255) COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '用户昵称',\n  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,\n  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,\n  PRIMARY KEY (`id`),\n  UNIQUE KEY `name_index` (`name`),\n  KEY `mobile_index` (`mobile`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;")
	assert.Nil(t, err)
	assert.Equal(t, "user_snake", table.Name.Source())
	assert.Equal(t, "id", table.PrimaryKey.Name.Source())
	assert.Equal(t, true, table.ContainsTime())
}

func TestConvertColumn(t *testing.T) {
	_, err := ConvertColumn("user", "user", []*model.Column{
		{
			Name:     "id",
			DataType: "bigint",
			Key:      "",
			Extra:    "",
			Comment:  "",
		},
	})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "missing primary key")

	_, err = ConvertColumn("user", "user", []*model.Column{
		{
			Name:     "id",
			DataType: "bigint",
			Key:      "PRI",
			Extra:    "",
			Comment:  "",
		},
		{
			Name:     "mobile",
			DataType: "varchar",
			Key:      "PRI",
			Extra:    "",
			Comment:  "手机号",
		},
	})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "only one primary key expected")

	table, err := ConvertColumn("user", "user", []*model.Column{
		{
			Name:     "id",
			DataType: "bigint",
			Key:      "PRI",
			Extra:    "auto_increment",
			Comment:  "",
		},
		{
			Name:     "mobile",
			DataType: "varchar",
			Key:      "UNI",
			Extra:    "",
			Comment:  "手机号",
		},
	})
	assert.Nil(t, err)
	assert.True(t, table.PrimaryKey.AutoIncrement && table.PrimaryKey.IsPrimaryKey)
	assert.Equal(t, "id", table.PrimaryKey.Name.Source())
	for _, item := range table.Fields {
		if item.Name.Source() == "mobile" {
			assert.True(t, item.IsUniqueKey)
			break
		}
	}
}
