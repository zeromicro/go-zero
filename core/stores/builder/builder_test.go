package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockedUser struct {
	ID       string `db:"id" json:"id,omitempty"`
	UserName string `db:"user_name" json:"userName,omitempty"`
	Sex      int    `db:"sex" json:"sex,omitempty"`
	UUID     string `db:"uuid" uuid:"uuid,omitempty"`
	Age      int    `db:"age" json:"age"`
}

func TestFieldNames(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		var u mockedUser
		out := RawFieldNames(&u)
		expected := []string{"`id`", "`user_name`", "`sex`", "`uuid`", "`age`"}
		assert.Equal(t, expected, out)
	})
}

type mockedUserWithOptions struct {
	ID       string `db:"id" json:"id,omitempty"`
	UserName string `db:"user_name,type=varchar,length=255" json:"userName,omitempty"`
	Sex      int    `db:"sex" json:"sex,omitempty"`
	UUID     string `db:",type=varchar,length=16" uuid:"uuid,omitempty"`
	Age      int    `db:"age" json:"age"`
}

func TestFieldNamesWithTagOptions(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		var u mockedUserWithOptions
		out := RawFieldNames(&u)
		expected := []string{"`id`", "`user_name`", "`sex`", "`UUID`", "`age`"}
		assert.Equal(t, expected, out)
	})
}

type mockedUserWithDashTag struct {
	ID       string `db:"id" json:"id,omitempty"`
	UserName string `db:"user_name" json:"userName,omitempty"`
	Mobile   string `db:"-" json:"mobile,omitempty"`
}

func TestFieldNamesWithDashTag(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		var u mockedUserWithDashTag
		out := RawFieldNames(&u)
		expected := []string{"`id`", "`user_name`"}
		assert.Equal(t, expected, out)
	})
}

type mockedUserWithDashTagAndOptions struct {
	ID       string `db:"id" json:"id,omitempty"`
	UserName string `db:"user_name,type=varchar,length=255" json:"userName,omitempty"`
	Mobile   string `db:"-,type=varchar,length=255" json:"mobile,omitempty"`
}

func TestFieldNamesWithDashTagAndOptions(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		var u mockedUserWithDashTagAndOptions
		out := RawFieldNames(&u)
		expected := []string{"`id`", "`user_name`"}
		assert.Equal(t, expected, out)
	})
}

func TestPostgreSqlJoin(t *testing.T) {
	// Test with empty input array
	var input []string
	var expectedOutput string
	assert.Equal(t, expectedOutput, PostgreSqlJoin(input))

	// Test with single element input array
	input = []string{"foo"}
	expectedOutput = "foo = $2"
	assert.Equal(t, expectedOutput, PostgreSqlJoin(input))

	// Test with multiple elements input array
	input = []string{"foo", "bar", "baz"}
	expectedOutput = "foo = $2, bar = $3, baz = $4"
	assert.Equal(t, expectedOutput, PostgreSqlJoin(input))
}

type testStruct struct {
	Foo string `db:"foo"`
	Bar int    `db:"bar"`
	Baz bool   `db:"-"`
}

func TestRawFieldNames(t *testing.T) {
	// Test with a struct without tags
	in := struct {
		Foo string
		Bar int
	}{}
	expectedOutput := []string{"`Foo`", "`Bar`"}
	assert.ElementsMatch(t, expectedOutput, RawFieldNames(in))

	// Test pg without db tag
	expectedOutput = []string{"Foo", "Bar"}
	assert.ElementsMatch(t, expectedOutput, RawFieldNames(in, true))

	// Test with a struct with tags
	input := testStruct{}
	expectedOutput = []string{"`foo`", "`bar`"}
	assert.ElementsMatch(t, expectedOutput, RawFieldNames(input))

	// Test with nil input (pointer)
	var nilInput *testStruct
	assert.Panics(t, func() {
		RawFieldNames(nilInput)
	}, "RawFieldNames should panic with nil input")

	// Test with non-struct input
	inputInt := 42
	assert.Panics(t, func() {
		RawFieldNames(inputInt)
	}, "RawFieldNames should panic with non-struct input")

	// Test with PostgreSQL flag
	input = testStruct{}
	expectedOutput = []string{"foo", "bar"}
	assert.ElementsMatch(t, expectedOutput, RawFieldNames(input, true))
}
