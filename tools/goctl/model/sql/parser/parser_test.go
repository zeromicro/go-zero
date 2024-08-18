package parser

import (
	_ "embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/model"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

func TestParsePlainText(t *testing.T) {
	sqlFile := filepath.Join(pathx.MustTempDir(), "tmp.sql")
	err := os.WriteFile(sqlFile, []byte("plain text"), 0o777)
	assert.Nil(t, err)

	_, err = Parse(sqlFile, "go_zero", false)
	assert.NotNil(t, err)
}

func TestParseSelect(t *testing.T) {
	sqlFile := filepath.Join(pathx.MustTempDir(), "tmp.sql")
	err := os.WriteFile(sqlFile, []byte("select * from user"), 0o777)
	assert.Nil(t, err)

	tables, err := Parse(sqlFile, "go_zero", false)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(tables))
}

//go:embed testdata/user.sql
var user string

func TestParseCreateTable(t *testing.T) {
	sqlFile := filepath.Join(pathx.MustTempDir(), "tmp.sql")
	err := os.WriteFile(sqlFile, []byte(user), 0o777)
	assert.Nil(t, err)

	tables, err := Parse(sqlFile, "go_zero", false)
	assert.Equal(t, 1, len(tables))
	table := tables[0]
	assert.Nil(t, err)
	assert.Equal(t, "test_user", table.Name.Source())
	assert.Equal(t, "id", table.PrimaryKey.Name.Source())
	assert.Equal(t, true, table.ContainsTime())
	assert.Equal(t, 2, len(table.UniqueIndex))
	assert.True(t, func() bool {
		for _, e := range table.Fields {
			if e.Comment != util.TrimNewLine(e.Comment) {
				return false
			}
		}

		return true
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
