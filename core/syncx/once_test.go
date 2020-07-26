package syncx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnce(t *testing.T) {
	var v int
	add := Once(func() {
		v++
	})

	for i := 0; i < 5; i++ {
		add()
	}

	assert.Equal(t, 1, v)
}
