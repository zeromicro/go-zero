package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/ddl-parser/parser"
)

func TestConvertDataType(t *testing.T) {
	v, err := ConvertDataType(parser.TinyInt, false, false)
	assert.Nil(t, err)
	assert.Equal(t, "int64", v)

	v, err = ConvertDataType(parser.TinyInt, false, true)
	assert.Nil(t, err)
	assert.Equal(t, "uint64", v)

	v, err = ConvertDataType(parser.TinyInt, true, false)
	assert.Nil(t, err)
	assert.Equal(t, "sql.NullInt64", v)

	v, err = ConvertDataType(parser.Timestamp, false, false)
	assert.Nil(t, err)
	assert.Equal(t, "time.Time", v)

	v, err = ConvertDataType(parser.Timestamp, true, false)
	assert.Nil(t, err)
	assert.Equal(t, "sql.NullTime", v)
}
