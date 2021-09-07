package load

import (
	"testing"
	"time"

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

func TestLoopTrue(t *testing.T) {
	ch := make(chan time.Time, 1)
	ch <- time.Now()
	close(ch)
	st := new(SheddingStat)
	logEnabled.Set(true)
	st.loop(ch)
}

func TestLoopTrueAndDrop(t *testing.T) {
	ch := make(chan time.Time, 1)
	ch <- time.Now()
	close(ch)
	st := new(SheddingStat)
	st.IncrementDrop()
	logEnabled.Set(true)
	st.loop(ch)
}

func TestLoopFalseAndDrop(t *testing.T) {
	ch := make(chan time.Time, 1)
	ch <- time.Now()
	close(ch)
	st := new(SheddingStat)
	st.IncrementDrop()
	logEnabled.Set(false)
	st.loop(ch)
}
