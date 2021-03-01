package parser

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/model"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

func TestParsePlainText(t *testing.T) {
	_, err := Parse("plain text")
	assert.NotNil(t, err)
}

func TestParseSelect(t *testing.T) {
	_, err := Parse("select * from user")
	assert.Equal(t, errUnsupportDDL, err)
}

func TestParseCreateTable(t *testing.T) {
	table, err := Parse("CREATE TABLE `test_user` (\n  `id` bigint NOT NULL AUTO_INCREMENT,\n  `mobile` varchar(255) COLLATE utf8mb4_bin NOT NULL,\n  `class` bigint NOT NULL,\n  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,\n  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,\n  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,\n  PRIMARY KEY (`id`),\n  UNIQUE KEY `mobile_unique` (`mobile`),\n  UNIQUE KEY `class_name_unique` (`class`,`name`),\n  KEY `create_index` (`create_time`),\n  KEY `name_index` (`name`)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;")
	assert.Nil(t, err)
	assert.Equal(t, "test_user", table.Name.Source())
	assert.Equal(t, "id", table.PrimaryKey.Name.Source())
	assert.Equal(t, true, table.ContainsTime())
	assert.Equal(t, true, func() bool {
		mobileUniqueIndex, ok := table.UniqueIndex["mobile_unique"]
		if !ok {
			return false
		}

		classNameUniqueIndex, ok := table.UniqueIndex["class_name_unique"]
		if !ok {
			return false
		}

		equal := func(f1, f2 []*Field) bool {
			sort.Slice(f1, func(i, j int) bool {
				return f1[i].Name.Source() < f1[j].Name.Source()
			})
			sort.Slice(f2, func(i, j int) bool {
				return f2[i].Name.Source() < f2[j].Name.Source()
			})

			if len(f2) != len(f2) {
				return false
			}

			for index, f := range f1 {
				if f1[index].Name.Source() != f.Name.Source() {
					return false
				}
			}
			return true
		}

		if !equal(mobileUniqueIndex, []*Field{
			{
				Name:         stringx.From("mobile"),
				DataBaseType: "varchar",
				DataType:     "string",
				SeqInIndex:   1,
			},
		}) {
			return false
		}

		return equal(classNameUniqueIndex, []*Field{
			{
				Name:         stringx.From("class"),
				DataBaseType: "bigint",
				DataType:     "int64",
				SeqInIndex:   1,
			},
			{
				Name:         stringx.From("name"),
				DataBaseType: "varchar",
				DataType:     "string",
				SeqInIndex:   2,
			},
		})
	}())
}

func TestConvertColumn(t *testing.T) {
	t.Run("missingPrimaryKey", func(t *testing.T) {
		columnData := model.ColumnData{
			Db:    "user",
			Table: "user",
			Columns: []*model.Column{
				{
					DbColumn: &model.DbColumn{
						Name:     "id",
						DataType: "bigint",
					},
				},
			},
		}
		_, err := columnData.Convert()
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "missing primary key")
	})

	t.Run("jointPrimaryKey", func(t *testing.T) {
		columnData := model.ColumnData{
			Db:    "user",
			Table: "user",
			Columns: []*model.Column{
				{
					DbColumn: &model.DbColumn{
						Name:     "id",
						DataType: "bigint",
					},
					Index: &model.DbIndex{
						IndexName: "PRIMARY",
					},
				},
				{
					DbColumn: &model.DbColumn{
						Name:     "mobile",
						DataType: "varchar",
						Comment:  "手机号",
					},
					Index: &model.DbIndex{
						IndexName: "PRIMARY",
					},
				},
			},
		}
		_, err := columnData.Convert()
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "joint primary key is not supported")
	})

	t.Run("normal", func(t *testing.T) {
		columnData := model.ColumnData{
			Db:    "user",
			Table: "user",
			Columns: []*model.Column{
				{
					DbColumn: &model.DbColumn{
						Name:     "id",
						DataType: "bigint",
						Extra:    "auto_increment",
					},
					Index: &model.DbIndex{
						IndexName:  "PRIMARY",
						SeqInIndex: 1,
					},
				},
				{
					DbColumn: &model.DbColumn{
						Name:     "mobile",
						DataType: "varchar",
						Comment:  "手机号",
					},
					Index: &model.DbIndex{
						IndexName:  "mobile_unique",
						SeqInIndex: 1,
					},
				},
			},
		}

		table, err := columnData.Convert()
		assert.Nil(t, err)
		assert.True(t, table.PrimaryKey.Index.IndexName == "PRIMARY" && table.PrimaryKey.Name == "id")
		for _, item := range table.Columns {
			if item.Name == "mobile" {
				assert.True(t, item.Index.NonUnique == 0)
				break
			}
		}
	})
}
