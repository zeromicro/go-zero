package sqlx

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscape(t *testing.T) {
	s := "a\x00\n\r\\'\"\x1ab"

	out := escape(s)

	assert.Equal(t, `a\x00\n\r\\\'\"\x1ab`, out)
}

func TestDesensitize(t *testing.T) {
	datasource := "user:pass@tcp(111.222.333.44:3306)/any_table?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"
	datasource = desensitize(datasource)
	assert.False(t, strings.Contains(datasource, "user"))
	assert.False(t, strings.Contains(datasource, "pass"))
	assert.True(t, strings.Contains(datasource, "tcp(111.222.333.44:3306)"))
}

func TestDesensitize_WithoutAccount(t *testing.T) {
	datasource := "tcp(111.222.333.44:3306)/any_table?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"
	datasource = desensitize(datasource)
	assert.True(t, strings.Contains(datasource, "tcp(111.222.333.44:3306)"))
}

func TestFormat(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		args   []interface{}
		expect string
		hasErr bool
	}{
		{
			name:   "mysql normal",
			query:  "select name, age from users where bool=? and phone=?",
			args:   []interface{}{true, "133"},
			expect: "select name, age from users where bool=1 and phone='133'",
		},
		{
			name:   "mysql normal",
			query:  "select name, age from users where bool=? and phone=?",
			args:   []interface{}{false, "133"},
			expect: "select name, age from users where bool=0 and phone='133'",
		},
		{
			name:   "pg normal",
			query:  "select name, age from users where bool=$1 and phone=$2",
			args:   []interface{}{true, "133"},
			expect: "select name, age from users where bool=1 and phone='133'",
		},
		{
			name:   "pg normal reverse",
			query:  "select name, age from users where bool=$2 and phone=$1",
			args:   []interface{}{"133", false},
			expect: "select name, age from users where bool=0 and phone='133'",
		},
		{
			name:   "pg error not number",
			query:  "select name, age from users where bool=$a and phone=$1",
			args:   []interface{}{"133", false},
			hasErr: true,
		},
		{
			name:   "pg error more args",
			query:  "select name, age from users where bool=$2 and phone=$1 and nickname=$3",
			args:   []interface{}{"133", false},
			hasErr: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual, err := format(test.query, test.args...)
			if test.hasErr {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, test.expect, actual)
			}
		})
	}
}
