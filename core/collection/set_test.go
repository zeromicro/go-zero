package collection

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
)

func init() {
	logx.Disable()
}

func BenchmarkRawSet(b *testing.B) {
	m := make(map[any]struct{})
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
	values := []any{1, 2, 3}

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
	set.Add([]any{1, 2, 3}...)

	// when
	set.Remove(2)

	// then
	assert.True(t, set.Contains(1) && !set.Contains(2) && set.Contains(3))
}

func TestCount(t *testing.T) {
	// given
	set := NewSet()
	set.Add([]any{1, 2, 3}...)

	// then
	assert.Equal(t, set.Count(), 3)
}

func TestKeysIntMismatch(t *testing.T) {
	set := NewSet()
	set.add(int64(1))
	set.add(2)
	vals := set.KeysInt()
	assert.EqualValues(t, []int{2}, vals)
}

func TestKeysInt64Mismatch(t *testing.T) {
	set := NewSet()
	set.add(1)
	set.add(int64(2))
	vals := set.KeysInt64()
	assert.EqualValues(t, []int64{2}, vals)
}

func TestKeysUintMismatch(t *testing.T) {
	set := NewSet()
	set.add(1)
	set.add(uint(2))
	vals := set.KeysUint()
	assert.EqualValues(t, []uint{2}, vals)
}

func TestKeysUint64Mismatch(t *testing.T) {
	set := NewSet()
	set.add(1)
	set.add(uint64(2))
	vals := set.KeysUint64()
	assert.EqualValues(t, []uint64{2}, vals)
}

func TestKeysStrMismatch(t *testing.T) {
	set := NewSet()
	set.add(1)
	set.add("2")
	vals := set.KeysStr()
	assert.EqualValues(t, []string{"2"}, vals)
}

func TestSetType(t *testing.T) {
	set := NewUnmanagedSet()
	set.add(1)
	set.add("2")
	vals := set.Keys()
	assert.ElementsMatch(t, []any{1, "2"}, vals)
}
