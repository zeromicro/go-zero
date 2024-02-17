package mapping

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapValuerWithInherit_Value(t *testing.T) {
	input := map[string]any{
		"discovery": map[string]any{
			"host": "localhost",
			"port": 8080,
		},
		"component": map[string]any{
			"name": "test",
		},
	}
	valuer := recursiveValuer{
		current: mapValuer(input["component"].(map[string]any)),
		parent: simpleValuer{
			current: mapValuer(input),
		},
	}

	val, ok := valuer.Value("discovery")
	assert.True(t, ok)

	m, ok := val.(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, "localhost", m["host"])
	assert.Equal(t, 8080, m["port"])
}

func TestRecursiveValuer_Value(t *testing.T) {
	input := map[string]any{
		"component": map[string]any{
			"name": "test",
			"foo": map[string]any{
				"bar": "baz",
			},
		},
		"foo": "value",
	}
	valuer := recursiveValuer{
		current: mapValuer(input["component"].(map[string]any)),
		parent: simpleValuer{
			current: mapValuer(input),
		},
	}

	val, ok := valuer.Value("foo")
	assert.True(t, ok)
	assert.EqualValues(t, map[string]any{
		"bar": "baz",
	}, val)
}
