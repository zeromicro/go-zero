package gen

import (
	"database/sql"
	_ "embed"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/builder"
	corstringx "github.com/zeromicro/go-zero/core/stringx"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"github.com/zeromicro/go-zero/tools/goctl/util/stringx"
)

//go:embed testdata/user.sql
var source string

func TestCacheModel(t *testing.T) {
	logx.Disable()
	_ = Clean()

	sqlFile := filepath.Join(pathx.MustTempDir(), "tmp.sql")
	err := os.WriteFile(sqlFile, []byte(source), 0o777)
	assert.Nil(t, err)

	dir := filepath.Join(pathx.MustTempDir(), "./testmodel")
	cacheDir := filepath.Join(dir, "cache")
	noCacheDir := filepath.Join(dir, "nocache")
	g, err := NewDefaultGenerator("cache", cacheDir, &config.Config{
		NamingFormat: "GoZero",
	})
	assert.Nil(t, err)

	err = g.StartFromDDL(sqlFile, true, false, "go_zero")
	assert.Nil(t, err)
	assert.True(t, func() bool {
		_, err := os.Stat(filepath.Join(cacheDir, "TestUserModel.go"))
		return err == nil
	}())
	g, err = NewDefaultGenerator("cache", noCacheDir, &config.Config{
		NamingFormat: "gozero",
	})
	assert.Nil(t, err)

	err = g.StartFromDDL(sqlFile, false, false, "go_zero")
	assert.Nil(t, err)
	assert.True(t, func() bool {
		_, err := os.Stat(filepath.Join(noCacheDir, "testusermodel.go"))
		return err == nil
	}())
}

func TestNamingModel(t *testing.T) {
	logx.Disable()
	_ = Clean()

	sqlFile := filepath.Join(pathx.MustTempDir(), "tmp.sql")
	err := os.WriteFile(sqlFile, []byte(source), 0o777)
	assert.Nil(t, err)

	dir, _ := filepath.Abs("./testmodel")
	camelDir := filepath.Join(dir, "camel")
	snakeDir := filepath.Join(dir, "snake")
	defer func() {
		_ = os.RemoveAll(dir)
	}()
	g, err := NewDefaultGenerator("cache", camelDir, &config.Config{
		NamingFormat: "GoZero",
	})
	assert.Nil(t, err)

	err = g.StartFromDDL(sqlFile, true, false, "go_zero")
	assert.Nil(t, err)
	assert.True(t, func() bool {
		_, err := os.Stat(filepath.Join(camelDir, "TestUserModel.go"))
		return err == nil
	}())
	g, err = NewDefaultGenerator("cache", snakeDir, &config.Config{
		NamingFormat: "go_zero",
	})
	assert.Nil(t, err)

	err = g.StartFromDDL(sqlFile, true, false, "go_zero")
	assert.Nil(t, err)
	assert.True(t, func() bool {
		_, err := os.Stat(filepath.Join(snakeDir, "test_user_model.go"))
		return err == nil
	}())
}

func TestFolderName(t *testing.T) {
	logx.Disable()
	_ = Clean()

	sqlFile := filepath.Join(pathx.MustTempDir(), "tmp.sql")
	err := os.WriteFile(sqlFile, []byte(source), 0o777)
	assert.Nil(t, err)

	dir, _ := filepath.Abs("./testmodel")
	camelDir := filepath.Join(dir, "go-camel")
	snakeDir := filepath.Join(dir, "go-snake")
	defer func() {
		_ = os.RemoveAll(dir)
	}()
	g, err := NewDefaultGenerator("cache", camelDir, &config.Config{
		NamingFormat: "GoZero",
	})
	assert.Nil(t, err)

	pkg := g.pkg

	err = g.StartFromDDL(sqlFile, true, true, "go_zero")
	assert.Nil(t, err)
	assert.True(t, func() bool {
		_, err := os.Stat(filepath.Join(camelDir, "TestUserModel.go"))
		return err == nil
	}())
	assert.Equal(t, pkg, g.pkg)

	g, err = NewDefaultGenerator("cache", snakeDir, &config.Config{
		NamingFormat: "go_zero",
	})
	assert.Nil(t, err)

	err = g.StartFromDDL(sqlFile, true, true, "go_zero")
	assert.Nil(t, err)
	assert.True(t, func() bool {
		_, err := os.Stat(filepath.Join(snakeDir, "test_user_model.go"))
		return err == nil
	}())
}

func TestWrapWithRawString(t *testing.T) {
	assert.Equal(t, "``", wrapWithRawString("", false))
	assert.Equal(t, "``", wrapWithRawString("``", false))
	assert.Equal(t, "`a`", wrapWithRawString("a", false))
	assert.Equal(t, "a", wrapWithRawString("a", true))
	assert.Equal(t, "`   `", wrapWithRawString("   ", false))
}

func TestFields(t *testing.T) {
	type Student struct {
		ID         int64           `db:"id"`
		Name       string          `db:"name"`
		Age        sql.NullInt64   `db:"age"`
		Score      sql.NullFloat64 `db:"score"`
		CreateTime time.Time       `db:"create_time"`
		UpdateTime sql.NullTime    `db:"update_time"`
	}
	var (
		studentFieldNames          = builder.RawFieldNames(&Student{})
		studentRows                = strings.Join(studentFieldNames, ",")
		studentRowsExpectAutoSet   = strings.Join(corstringx.Remove(studentFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
		studentRowsWithPlaceHolder = strings.Join(corstringx.Remove(studentFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"
	)

	assert.Equal(t, []string{"`id`", "`name`", "`age`", "`score`", "`create_time`", "`update_time`"}, studentFieldNames)
	assert.Equal(t, "`id`,`name`,`age`,`score`,`create_time`,`update_time`", studentRows)
	assert.Equal(t, "`name`,`age`,`score`", studentRowsExpectAutoSet)
	assert.Equal(t, "`name`=?,`age`=?,`score`=?", studentRowsWithPlaceHolder)
}

func Test_genPublicModel(t *testing.T) {
	var err error
	dir := pathx.MustTempDir()
	modelDir := path.Join(dir, "model")
	err = os.MkdirAll(modelDir, 0o777)
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	modelFilename := filepath.Join(modelDir, "foo.sql")
	err = os.WriteFile(modelFilename, []byte(source), 0o777)
	require.NoError(t, err)

	g, err := NewDefaultGenerator("cache", modelDir, &config.Config{
		NamingFormat: config.DefaultFormat,
	})
	require.NoError(t, err)

	tables, err := parser.Parse(modelFilename, "", false)
	require.Equal(t, 1, len(tables))

	code, err := g.genModelCustom(*tables[0], false)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(code, "package model"))
	assert.True(t, strings.Contains(code, `	TestUserModel interface {
		testUserModel
		withSession(session sqlx.Session) TestUserModel
	}`))
	assert.True(t, strings.Contains(code, "customTestUserModel struct {\n\t\t*defaultTestUserModel\n\t}\n"))
	assert.True(t, strings.Contains(code, "func NewTestUserModel(conn sqlx.SqlConn) TestUserModel {"))
}

func TestGenFindOneByFieldWithPartialIndex(t *testing.T) {
	primaryField := &parser.Field{
		Name:     stringx.From("id"),
		DataType: "int64",
		Comment:  "主键",
	}
	emailField := &parser.Field{
		Name:     stringx.From("email"),
		DataType: "string",
		Comment:  "邮箱",
	}
	parsedTable := parser.Table{
		Name: stringx.From("user"),
		Db:   stringx.From("go_zero"),
		PrimaryKey: parser.Primary{
			Field: *primaryField,
		},
		UniqueIndex: map[string][]*parser.Field{
			"idx_active_email": {emailField},
		},
		UniqueIndexPredicate: map[string]string{
			"idx_active_email": "(status = 1) AND (deleted_at IS NULL)",
		},
		Fields: []*parser.Field{
			primaryField,
			emailField,
		},
	}

	primaryKey, uniqueKeys := genCacheKeys("cache", parsedTable)

	table := Table{
		Table:                  parsedTable,
		PrimaryCacheKey:        primaryKey,
		UniqueCacheKey:         uniqueKeys,
		ContainsUniqueCacheKey: len(uniqueKeys) > 0,
	}

	result, err := genFindOneByField(table, false, true)
	assert.NoError(t, err)

	assert.Contains(t, result.findOneMethod, "where email = $1 and (status = 1) AND (deleted_at IS NULL)")
	assert.Contains(t, result.findOneMethod, "FindOneByEmail")

	// Verify non-partial indexes do NOT append predicate
	parsedTable2 := parser.Table{
		Name: stringx.From("user"),
		Db:   stringx.From("go_zero"),
		PrimaryKey: parser.Primary{
			Field: *primaryField,
		},
		UniqueIndex: map[string][]*parser.Field{
			"email_unique": {emailField},
		},
		UniqueIndexPredicate: map[string]string{},
		Fields: []*parser.Field{
			primaryField,
			emailField,
		},
	}
	_, uniqueKeys2 := genCacheKeys("cache", parsedTable2)
	table2 := Table{
		Table:                  parsedTable2,
		PrimaryCacheKey:        primaryKey,
		UniqueCacheKey:         uniqueKeys2,
		ContainsUniqueCacheKey: len(uniqueKeys2) > 0,
	}
	result2, err := genFindOneByField(table2, false, true)
	assert.NoError(t, err)
	assert.Contains(t, result2.findOneMethod, "where email = $1")
	assert.NotContains(t, result2.findOneMethod, " and (")
}
