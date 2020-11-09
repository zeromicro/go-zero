package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertDataType(t *testing.T) {
	v, err := ConvertDataType("tinyint")
	assert.Nil(t, err)
	assert.Equal(t, "int64", v)

	v, err = ConvertDataType("timestamp")
	assert.Nil(t, err)
	assert.Equal(t, "time.Time", v)

	_, err = ConvertDataType("float32")
	assert.NotNil(t, err)
}
