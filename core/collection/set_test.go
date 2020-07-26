package collection

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkRawSet(b *testing.B) {
	m := make(map[interface{}]struct{})
	for i := 0; i < b.N; i++ {
		m[i] = struct{}{}
		_ = m[i]
	}
}

func BenchmarkUnmanagedSet(b *testing.B) {
	s := NewUnmanagedSet()
	for i := 0; i < b.N; i++ {
		s.Add(i)
		_ = s.Contains(i)
	}
}

func BenchmarkSet(b *testing.B) {
	s := NewSet()
	for i := 0; i < b.N; i++ {
		s.AddInt(i)
		_ = s.Contains(i)
	}
}

func TestAdd(t *testing.T) {
	// given
	set := NewUnmanagedSet()
	values := []interface{}{1, 2, 3}

	// when
	set.Add(values...)

	// then
	assert.True(t, set.Contains(1) && set.Contains(2) && set.Contains(3))
	assert.Equal(t, len(values), len(set.Keys()))
}

func TestAddInt(t *testing.T) {
	// given
	set := NewSet()
	values := []int{1, 2, 3}

	// when
	set.AddInt(values...)

	// then
	assert.True(t, set.Contains(1) && set.Contains(2) && set.Contains(3))
	keys := set.KeysInt()
	sort.Ints(keys)
	assert.EqualValues(t, values, keys)
}

func TestAddInt64(t *testing.T) {
	// given
	set := NewSet()
	values := []int64{1, 2, 3}

	// when
	set.AddInt64(values...)

	// then
	assert.True(t, set.Contains(int64(1)) && set.Contains(int64(2)) && set.Contains(int64(3)))
	assert.Equal(t, len(values), len(set.KeysInt64()))
}

func TestAddUint(t *testing.T) {
	// given
	set := NewSet()
	values := []uint{1, 2, 3}

	// when
	set.AddUint(values...)

	// then
	assert.True(t, set.Contains(uint(1)) && set.Contains(uint(2)) && set.Contains(uint(3)))
	assert.Equal(t, len(values), len(set.KeysUint()))
}

func TestAddUint64(t *testing.T) {
	// given
	set := NewSet()
	values := []uint64{1, 2, 3}

	// when
	set.AddUint64(values...)

	// then
	assert.True(t, set.Contains(uint64(1)) && set.Contains(uint64(2)) && set.Contains(uint64(3)))
	assert.Equal(t, len(values), len(set.KeysUint64()))
}

func TestAddStr(t *testing.T) {
	// given
	set := NewSet()
	values := []string{"1", "2", "3"}

	// when
	set.AddStr(values...)

	// then
	assert.True(t, set.Contains("1") && set.Contains("2") && set.Contains("3"))
	assert.Equal(t, len(values), len(set.KeysStr()))
}

func TestContainsWithoutElements(t *testing.T) {
	// given
	set := NewSet()

	// then
	assert.False(t, set.Contains(1))
}

func TestContainsUnmanagedWithoutElements(t *testing.T) {
	// given
	set := NewUnmanagedSet()

	// then
	assert.False(t, set.Contains(1))
}

func TestRemove(t *testing.T) {
	// given
	set := NewSet()
	set.Add([]interface{}{1, 2, 3}...)

	// when
	set.Remove(2)

	// then
	assert.True(t, set.Contains(1) && !set.Contains(2) && set.Contains(3))
}

func TestCount(t *testing.T) {
	// given
	set := NewSet()
	set.Add([]interface{}{1, 2, 3}...)

	// then
	assert.Equal(t, set.Count(), 3)
}
