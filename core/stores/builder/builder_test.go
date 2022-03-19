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
