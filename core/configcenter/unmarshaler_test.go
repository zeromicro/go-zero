package configurator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUnmarshaler(t *testing.T) {
	RegisterUnmarshaler("test", func(data []byte, v interface{}) error {
		return nil
	})

	_, ok := Unmarshaler("test")
	assert.True(t, ok)

	_, ok = Unmarshaler("test2")
	assert.False(t, ok)

	_, ok = Unmarshaler("json")
	assert.True(t, ok)

	_, ok = Unmarshaler("toml")
	assert.True(t, ok)

	_, ok = Unmarshaler("yaml")
	assert.True(t, ok)
}
