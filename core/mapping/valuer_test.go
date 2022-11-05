package mapping

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapValuerWithInherit_Value(t *testing.T) {
	input := map[string]interface{}{
		"discovery": map[string]interface{}{
			"host": "localhost",
			"port": 8080,
		},
		"component": map[string]interface{}{
			"name": "test",
		},
	}
	valuer := recursiveValuer{
		current: mapValuer(input["component"].(map[string]interface{})),
		parent: simpleValuer{
			current: mapValuer(input),
		},
	}

	val, ok := valuer.Value("discovery")
	assert.True(t, ok)

	m, ok := val.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "localhost", m["host"])
	assert.Equal(t, 8080, m["port"])
}
