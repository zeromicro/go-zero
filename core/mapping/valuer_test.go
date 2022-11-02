package mapping

import (
	"fmt"
	"testing"
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
			parent:  nil,
		},
	}
	val, ok := valuer.Value("discovery")
	fmt.Println(val, ok)
}
