package mathx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnstable_AroundDuration(t *testing.T) {
	unstable := NewUnstable(0.05)
	for i := 0; i < 1000; i++ {
		val := unstable.AroundDuration(time.Second)
		assert.True(t, float64(time.Second)*0.95 <= float64(val))
		assert.True(t, float64(val) <= float64(time.Second)*1.05)
	}
}

func TestUnstable_AroundInt(t *testing.T) {
	const target = 10000
	unstable := NewUnstable(0.05)
	for i := 0; i < 1000; i++ {
		val := unstable.AroundInt(target)
		assert.True(t, float64(target)*0.95 <= float64(val))
		assert.True(t, float64(val) <= float64(target)*1.05)
	}
}

func TestUnstable_AroundIntLarge(t *testing.T) {
	const target int64 = 10000
	unstable := NewUnstable(5)
	for i := 0; i < 1000; i++ {
		val := unstable.AroundInt(target)
		assert.True(t, 0 <= val)
		assert.True(t, val <= 2*target)
	}
}

func TestUnstable_AroundIntNegative(t *testing.T) {
	const target int64 = 10000
	unstable := NewUnstable(-0.05)
	for i := 0; i < 1000; i++ {
		val := unstable.AroundInt(target)
		assert.Equal(t, target, val)
	}
}

func TestUnstable_Distribution(t *testing.T) {
	const (
		seconds = 10000
		total   = 10000
	)

	m := make(map[int]int)
	expiry := NewUnstable(0.05)
	for i := 0; i < total; i++ {
		val := int(expiry.AroundInt(seconds))
		m[val]++
	}

	_, ok := m[0]
	assert.False(t, ok)

	mi := make(map[any]int, len(m))
	for k, v := range m {
		mi[k] = v
	}
	entropy := CalcEntropy(mi)
	assert.True(t, len(m) > 1)
	assert.True(t, entropy > 0.95)
}
