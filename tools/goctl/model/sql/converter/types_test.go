package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertDataType(t *testing.T) {
	v, err := ConvertDataType("tinyint", false)
	assert.Nil(t, err)
	assert.Equal(t, "int64", v)

	v, err = ConvertDataType("tinyint", true)
	assert.Nil(t, err)
	assert.Equal(t, "sql.NullInt64", v)

	v, err = ConvertDataType("timestamp", false)
	assert.Nil(t, err)
	assert.Equal(t, "time.Time", v)

	v, err = ConvertDataType("timestamp", true)
	assert.Nil(t, err)
	assert.Equal(t, "sql.NullTime", v)

	_, err = ConvertDataType("float32", false)
	assert.NotNil(t, err)
}
