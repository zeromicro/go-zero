package builderx

import (
	"fmt"
	"testing"

	"github.com/go-xorm/builder"
	"github.com/stretchr/testify/assert"
)

type mockedUser struct {
	// 自增id
	Id string `db:"id" json:"id,omitempty"`
	// 姓名
	UserName string `db:"user_name" json:"userName,omitempty"`
	// 1男,2女
	Sex  int    `db:"sex" json:"sex,omitempty"`
	Uuid string `db:"uuid" uuid:"uuid,omitempty"`
	Age  int    `db:"age" json:"age"`
}

var (
	userFieldsWithRawStringQuote    = RawFieldNames(mockedUser{})
	userFieldsWithoutRawStringQuote = FieldNames(mockedUser{})
)

func TestFieldNames(t *testing.T) {
	t.Run("old", func(t *testing.T) {
		var u mockedUser
		out := FieldNames(&u)
		expected := []string{"id", "user_name", "sex", "uuid", "age"}
		assert.Equal(t, expected, out)
	})

	t.Run("new", func(t *testing.T) {
		var u mockedUser
		out := RawFieldNames(&u)
		expected := []string{"`id`", "`user_name`", "`sex`", "`uuid`", "`age`"}
		assert.Equal(t, expected, out)
	})
}

func TestNewEq(t *testing.T) {
	u := &mockedUser{
		Id:       "123456",
		UserName: "wahaha",
	}
	out := NewEq(u)
	fmt.Println(out)
	actual := builder.Eq{"id": "123456", "user_name": "wahaha"}
	assert.Equal(t, out, actual)
}

// @see https://github.com/go-xorm/builder
func TestBuilderSql(t *testing.T) {
	u := &mockedUser{
		Id: "123123",
	}
	fields := RawFieldNames(u)
	eq := NewEq(u)
	sql, args, err := builder.Select(fields...).From("user").Where(eq).ToSQL()
	fmt.Println(sql, args, err)

	actualSql := "SELECT `id`,`user_name`,`sex`,`uuid`,`age` FROM user WHERE id=?"
	actualArgs := []interface{}{"123123"}
	assert.Equal(t, sql, actualSql)
	assert.Equal(t, args, actualArgs)
}

func TestBuildSqlDefaultValue(t *testing.T) {
	var eq = builder.Eq{}
	eq["age"] = 0
	eq["user_name"] = ""

	t.Run("raw", func(t *testing.T) {
		sql, args, err := builder.Select(userFieldsWithRawStringQuote...).From("user").Where(eq).ToSQL()
		fmt.Println(sql, args, err)

		actualSql := "SELECT `id`,`user_name`,`sex`,`uuid`,`age` FROM user WHERE age=? AND user_name=?"
		actualArgs := []interface{}{0, ""}
		assert.Equal(t, sql, actualSql)
		assert.Equal(t, args, actualArgs)
	})

	t.Run("withour raw quote", func(t *testing.T) {
		sql, args, err := builder.Select(userFieldsWithoutRawStringQuote...).From("user").Where(eq).ToSQL()
		fmt.Println(sql, args, err)

		actualSql := "SELECT id,user_name,sex,uuid,age FROM user WHERE age=? AND user_name=?"
		actualArgs := []interface{}{0, ""}
		assert.Equal(t, sql, actualSql)
		assert.Equal(t, args, actualArgs)
	})
}

func TestBuilderSqlIn(t *testing.T) {
	u := &mockedUser{
		Age: 18,
	}
	gtU := NewGt(u)
	in := builder.In("id", []string{"1", "2", "3"})
	sql, args, err := builder.Select(userFieldsWithRawStringQuote...).From("user").Where(in).And(gtU).ToSQL()
	fmt.Println(sql, args, err)

	actualSql := "SELECT `id`,`user_name`,`sex`,`uuid`,`age` FROM user WHERE id IN (?,?,?) AND age>?"
	actualArgs := []interface{}{"1", "2", "3", 18}
	assert.Equal(t, sql, actualSql)
	assert.Equal(t, args, actualArgs)
}

func TestBuildSqlLike(t *testing.T) {
	like := builder.Like{"name", "wang"}
	sql, args, err := builder.Select(userFieldsWithRawStringQuote...).From("user").Where(like).ToSQL()
	fmt.Println(sql, args, err)

	actualSql := "SELECT `id`,`user_name`,`sex`,`uuid`,`age` FROM user WHERE name LIKE ?"
	actualArgs := []interface{}{"%wang%"}
	assert.Equal(t, sql, actualSql)
	assert.Equal(t, args, actualArgs)
}
