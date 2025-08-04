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

// Set functionality tests
func TestTypedSetInt(t *testing.T) {
	set := NewSet[int]()
	values := []int{1, 2, 3, 2, 1} // Contains duplicates

	// Test adding
	set.Add(values...)
	assert.Equal(t, 3, set.Count()) // Should only have 3 elements after deduplication

	// Test contains
	assert.True(t, set.Contains(1))
	assert.True(t, set.Contains(2))
	assert.True(t, set.Contains(3))
	assert.False(t, set.Contains(4))

	// Test getting all keys
	keys := set.Keys()
	sort.Ints(keys)
	assert.EqualValues(t, []int{1, 2, 3}, keys)

	// Test removal
	set.Remove(2)
	assert.False(t, set.Contains(2))
	assert.Equal(t, 2, set.Count())
}

func TestTypedSetStringOps(t *testing.T) {
	set := NewSet[string]()
	values := []string{"a", "b", "c", "b", "a"}

	set.Add(values...)
	assert.Equal(t, 3, set.Count())

	assert.True(t, set.Contains("a"))
	assert.True(t, set.Contains("b"))
	assert.True(t, set.Contains("c"))
	assert.False(t, set.Contains("d"))

	keys := set.Keys()
	sort.Strings(keys)
	assert.EqualValues(t, []string{"a", "b", "c"}, keys)
}

func TestTypedSetClear(t *testing.T) {
	set := NewSet[int]()
	set.Add(1, 2, 3)
	assert.Equal(t, 3, set.Count())

	set.Clear()
	assert.Equal(t, 0, set.Count())
	assert.False(t, set.Contains(1))
}

func TestTypedSetEmpty(t *testing.T) {
	set := NewSet[int]()
	assert.Equal(t, 0, set.Count())
	assert.False(t, set.Contains(1))
	assert.Empty(t, set.Keys())
}

func TestTypedSetMultipleTypes(t *testing.T) {
	// Test different typed generic sets
	intSet := NewSet[int]()
	int64Set := NewSet[int64]()
	uintSet := NewSet[uint]()
	uint64Set := NewSet[uint64]()
	stringSet := NewSet[string]()

	intSet.Add(1, 2, 3)
	int64Set.Add(1, 2, 3)
	uintSet.Add(1, 2, 3)
	uint64Set.Add(1, 2, 3)
	stringSet.Add("1", "2", "3")

	assert.Equal(t, 3, intSet.Count())
	assert.Equal(t, 3, int64Set.Count())
	assert.Equal(t, 3, uintSet.Count())
	assert.Equal(t, 3, uint64Set.Count())
	assert.Equal(t, 3, stringSet.Count())
}

// Set benchmarks
func BenchmarkTypedIntSet(b *testing.B) {
	s := NewSet[int]()
	for i := 0; i < b.N; i++ {
		s.Add(i)
		_ = s.Contains(i)
	}
}

func BenchmarkTypedStringSet(b *testing.B) {
	s := NewSet[string]()
	for i := 0; i < b.N; i++ {
		s.Add(string(rune(i)))
		_ = s.Contains(string(rune(i)))
	}
}

// Legacy tests remain unchanged for backward compatibility
func BenchmarkRawSet(b *testing.B) {
	m := make(map[any]struct{})
	for i := 0; i < b.N; i++ {
		m[i] = struct{}{}
		_ = m[i]
	}
}

func TestAdd(t *testing.T) {
	// given
	set := NewSet[int]()
	values := []int{1, 2, 3}

	// when
	set.Add(values...)

	// then
	assert.True(t, set.Contains(1) && set.Contains(2) && set.Contains(3))
	assert.Equal(t, len(values), len(set.Keys()))
}

func TestAddInt(t *testing.T) {
	// given
	set := NewSet[int]()
	values := []int{1, 2, 3}

	// when
	set.Add(values...)

	// then
	assert.True(t, set.Contains(1) && set.Contains(2) && set.Contains(3))
	keys := set.Keys()
	sort.Ints(keys)
	assert.EqualValues(t, values, keys)
}

func TestAddInt64(t *testing.T) {
	// given
	set := NewSet[int64]()
	values := []int64{1, 2, 3}

	// when
	set.Add(values...)

	// then
	assert.True(t, set.Contains(1) && set.Contains(2) && set.Contains(3))
	assert.Equal(t, len(values), len(set.Keys()))
}

func TestAddUint(t *testing.T) {
	// given
	set := NewSet[uint]()
	values := []uint{1, 2, 3}

	// when
	set.Add(values...)

	// then
	assert.True(t, set.Contains(1) && set.Contains(2) && set.Contains(3))
	assert.Equal(t, len(values), len(set.Keys()))
}

func TestAddUint64(t *testing.T) {
	// given
	set := NewSet[uint64]()
	values := []uint64{1, 2, 3}

	// when
	set.Add(values...)

	// then
	assert.True(t, set.Contains(1) && set.Contains(2) && set.Contains(3))
	assert.Equal(t, len(values), len(set.Keys()))
}

func TestAddStr(t *testing.T) {
	// given
	set := NewSet[string]()
	values := []string{"1", "2", "3"}

	// when
	set.Add(values...)

	// then
	assert.True(t, set.Contains("1") && set.Contains("2") && set.Contains("3"))
	assert.Equal(t, len(values), len(set.Keys()))
}

func TestContainsWithoutElements(t *testing.T) {
	// given
	set := NewSet[int]()

	// then
	assert.False(t, set.Contains(1))
}

func TestRemove(t *testing.T) {
	// given
	set := NewSet[int]()
	set.Add([]int{1, 2, 3}...)

	// when
	set.Remove(2)

	// then
	assert.True(t, set.Contains(1) && !set.Contains(2) && set.Contains(3))
}

func TestCount(t *testing.T) {
	// given
	set := NewSet[int]()
	set.Add([]int{1, 2, 3}...)

	// then
	assert.Equal(t, set.Count(), 3)
}
