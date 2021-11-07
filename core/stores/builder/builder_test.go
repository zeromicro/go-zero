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
