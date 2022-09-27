package collection

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stringx"
)

func TestSafeMap(t *testing.T) {
	tests := []struct {
		size      int
		exception int
	}{
		{
			100000,
			2000,
		},
		{
			100000,
			50,
		},
	}
	for _, test := range tests {
		t.Run(stringx.Rand(), func(t *testing.T) {
			testSafeMapWithParameters(t, test.size, test.exception)
		})
	}
}

func TestSafeMap_CopyNew(t *testing.T) {
	const (
		size       = 100000
		exception1 = 5
		exception2 = 500
	)
	m := NewSafeMap()

	for i := 0; i < size; i++ {
		m.Set(i, i)
	}
	for i := 0; i < size; i++ {
		if i%exception1 == 0 {
			m.Del(i)
		}
	}

	for i := size; i < size<<1; i++ {
		m.Set(i, i)
	}
	for i := size; i < size<<1; i++ {
		if i%exception2 != 0 {
			m.Del(i)
		}
	}

	for i := 0; i < size; i++ {
		val, ok := m.Get(i)
		if i%exception1 != 0 {
			assert.True(t, ok)
			assert.Equal(t, i, val.(int))
		} else {
			assert.False(t, ok)
		}
	}
	for i := size; i < size<<1; i++ {
		val, ok := m.Get(i)
		if i%exception2 == 0 {
			assert.True(t, ok)
			assert.Equal(t, i, val.(int))
		} else {
			assert.False(t, ok)
		}
	}
}

func testSafeMapWithParameters(t *testing.T, size, exception int) {
	m := NewSafeMap()

	for i := 0; i < size; i++ {
		m.Set(i, i)
	}
	for i := 0; i < size; i++ {
		if i%exception != 0 {
			m.Del(i)
		}
	}

	assert.Equal(t, size/exception, m.Size())

	for i := size; i < size<<1; i++ {
		m.Set(i, i)
	}
	for i := size; i < size<<1; i++ {
		if i%exception != 0 {
			m.Del(i)
		}
	}

	for i := 0; i < size<<1; i++ {
		val, ok := m.Get(i)
		if i%exception == 0 {
			assert.True(t, ok)
			assert.Equal(t, i, val.(int))
		} else {
			assert.False(t, ok)
		}
	}
}

func TestSafeMap_Range(t *testing.T) {
	const (
		size       = 100000
		exception1 = 5
		exception2 = 500
	)

	m := NewSafeMap()
	newMap := NewSafeMap()

	for i := 0; i < size; i++ {
		m.Set(i, i)
	}
	for i := 0; i < size; i++ {
		if i%exception1 == 0 {
			m.Del(i)
		}
	}

	for i := size; i < size<<1; i++ {
		m.Set(i, i)
	}
	for i := size; i < size<<1; i++ {
		if i%exception2 != 0 {
			m.Del(i)
		}
	}

	var count int32
	m.Range(func(k, v interface{}) bool {
		atomic.AddInt32(&count, 1)
		newMap.Set(k, v)
		return true
	})
	assert.Equal(t, int(atomic.LoadInt32(&count)), m.Size())
	assert.Equal(t, m.dirtyNew, newMap.dirtyNew)
	assert.Equal(t, m.dirtyOld, newMap.dirtyOld)
}
