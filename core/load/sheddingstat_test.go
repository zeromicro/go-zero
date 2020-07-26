package load

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSheddingStat(t *testing.T) {
	st := NewSheddingStat("any")
	for i := 0; i < 3; i++ {
		st.IncrementTotal()
	}
	for i := 0; i < 5; i++ {
		st.IncrementPass()
	}
	for i := 0; i < 7; i++ {
		st.IncrementDrop()
	}
	result := st.reset()
	assert.Equal(t, int64(3), result.Total)
	assert.Equal(t, int64(5), result.Pass)
	assert.Equal(t, int64(7), result.Drop)
}
