package collection

import (
	"github.com/zeromicro/go-zero/core/lang"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	unmanaged = iota
	untyped
	intType
	int64Type
	uintType
	uint64Type
	stringType
)

// TypedSet is a type-safe generic set collection. It's not thread-safe,
// use with synchronization for concurrent access.
//
// Advantages over the legacy Set:
// - Compile-time type safety (no runtime type validation needed)
// - Better performance (no type assertions or reflection overhead)
// - Cleaner API (single Add method instead of multiple type-specific methods)
// - No need for type-specific Keys methods (KeysInt, KeysStr, etc.)
// - Zero-allocation for empty checks and direct type access
type TypedSet[T comparable] struct {
	data map[T]lang.PlaceholderType
}

// NewTypedSet returns a new type-safe set.
func NewTypedSet[T comparable]() *TypedSet[T] {
	return &TypedSet[T]{
		data: make(map[T]lang.PlaceholderType),
	}
}

// NewIntSet returns a new int-typed set.
func NewIntSet() *TypedSet[int] {
	return NewTypedSet[int]()
}

// NewInt64Set returns a new int64-typed set.
func NewInt64Set() *TypedSet[int64] {
	return NewTypedSet[int64]()
}

// NewUintSet returns a new uint-typed set.
func NewUintSet() *TypedSet[uint] {
	return NewTypedSet[uint]()
}

// NewUint64Set returns a new uint64-typed set.
func NewUint64Set() *TypedSet[uint64] {
	return NewTypedSet[uint64]()
}

// NewStringSet returns a new string-typed set.
func NewStringSet() *TypedSet[string] {
	return NewTypedSet[string]()
}

// Add adds items to the set. Duplicates are automatically ignored.
func (s *TypedSet[T]) Add(items ...T) {
	for _, item := range items {
		s.data[item] = lang.Placeholder
	}
}

// Contains checks if an item exists in the set.
func (s *TypedSet[T]) Contains(item T) bool {
	_, ok := s.data[item]
	return ok
}

// Remove removes an item from the set.
func (s *TypedSet[T]) Remove(item T) {
	delete(s.data, item)
}

// Keys returns all elements in the set as a slice.
func (s *TypedSet[T]) Keys() []T {
	keys := make([]T, 0, len(s.data))
	for key := range s.data {
		keys = append(keys, key)
	}
	return keys
}

// Count returns the number of items in the set.
func (s *TypedSet[T]) Count() int {
	return len(s.data)
}

// Clear removes all items from the set.
func (s *TypedSet[T]) Clear() {
	s.data = make(map[T]lang.PlaceholderType)
}

// Set is not thread-safe, for concurrent use, make sure to use it with synchronization.
// Deprecated: Use TypedSet[T] instead for better type safety and performance.
// TypedSet provides compile-time type checking and eliminates the need for type-specific methods.
type Set struct {
	data map[any]lang.PlaceholderType
	tp   int
}

// NewSet returns a managed Set, can only put the values with the same type.
// Deprecated: Use NewTypedSet[T]() instead for better type safety and performance.
// Example: NewIntSet() instead of NewSet() with AddInt()
func NewSet() *Set {
	return &Set{
		data: make(map[any]lang.PlaceholderType),
		tp:   untyped,
	}
}

// NewUnmanagedSet returns an unmanaged Set, which can put values with different types.
// Deprecated: Use TypedSet[any] or multiple TypedSet instances for different types instead.
// If you really need mixed types, consider using map[any]struct{} directly.
func NewUnmanagedSet() *Set {
	return &Set{
		data: make(map[any]lang.PlaceholderType),
		tp:   unmanaged,
	}
}

// Add adds i into s.
// Deprecated: Use TypedSet[T].Add() instead for better type safety and performance.
func (s *Set) Add(i ...any) {
	for _, each := range i {
		s.add(each)
	}
}

// AddInt adds int values ii into s.
// Deprecated: Use NewIntSet().Add() instead for better type safety and performance.
// Example: intSet := NewIntSet(); intSet.Add(1, 2, 3)
func (s *Set) AddInt(ii ...int) {
	for _, each := range ii {
		s.add(each)
	}
}

// AddInt64 adds int64 values ii into s.
// Deprecated: Use NewInt64Set().Add() instead for better type safety and performance.
// Example: int64Set := NewInt64Set(); int64Set.Add(1, 2, 3)
func (s *Set) AddInt64(ii ...int64) {
	for _, each := range ii {
		s.add(each)
	}
}

// AddUint adds uint values ii into s.
// Deprecated: Use NewUintSet().Add() instead for better type safety and performance.
// Example: uintSet := NewUintSet(); uintSet.Add(1, 2, 3)
func (s *Set) AddUint(ii ...uint) {
	for _, each := range ii {
		s.add(each)
	}
}

// AddUint64 adds uint64 values ii into s.
// Deprecated: Use NewUint64Set().Add() instead for better type safety and performance.
// Example: uint64Set := NewUint64Set(); uint64Set.Add(1, 2, 3)
func (s *Set) AddUint64(ii ...uint64) {
	for _, each := range ii {
		s.add(each)
	}
}

// AddStr adds string values ss into s.
// Deprecated: Use NewStringSet().Add() instead for better type safety and performance.
// Example: stringSet := NewStringSet(); stringSet.Add("a", "b", "c")
func (s *Set) AddStr(ss ...string) {
	for _, each := range ss {
		s.add(each)
	}
}

// Contains checks if i is in s.
// Deprecated: Use TypedSet[T].Contains() instead for better type safety and performance.
func (s *Set) Contains(i any) bool {
	if len(s.data) == 0 {
		return false
	}

	s.validate(i)
	_, ok := s.data[i]
	return ok
}

// Keys returns the keys in s.
// Deprecated: Use TypedSet[T].Keys() instead for better type safety and performance.
func (s *Set) Keys() []any {
	var keys []any

	for key := range s.data {
		keys = append(keys, key)
	}

	return keys
}

// KeysInt returns the int keys in s.
// Deprecated: Use NewIntSet().Keys() instead for better type safety and performance.
// The TypedSet version returns []int directly without type casting.
func (s *Set) KeysInt() []int {
	var keys []int

	for key := range s.data {
		if intKey, ok := key.(int); ok {
			keys = append(keys, intKey)
		}
	}

	return keys
}

// KeysInt64 returns int64 keys in s.
// Deprecated: Use NewInt64Set().Keys() instead for better type safety and performance.
// The TypedSet version returns []int64 directly without type casting.
func (s *Set) KeysInt64() []int64 {
	var keys []int64

	for key := range s.data {
		if intKey, ok := key.(int64); ok {
			keys = append(keys, intKey)
		}
	}

	return keys
}

// KeysUint returns uint keys in s.
// Deprecated: Use NewUintSet().Keys() instead for better type safety and performance.
// The TypedSet version returns []uint directly without type casting.
func (s *Set) KeysUint() []uint {
	var keys []uint

	for key := range s.data {
		if intKey, ok := key.(uint); ok {
			keys = append(keys, intKey)
		}
	}

	return keys
}

// KeysUint64 returns uint64 keys in s.
//
// Deprecated: Use NewUint64Set().Keys() instead for better type safety and performance.
// The TypedSet version returns []uint64 directly without type casting.
func (s *Set) KeysUint64() []uint64 {
	var keys []uint64

	for key := range s.data {
		if intKey, ok := key.(uint64); ok {
			keys = append(keys, intKey)
		}
	}

	return keys
}

// KeysStr returns string keys in s.
// Deprecated: Use NewStringSet().Keys() instead for better type safety and performance.
// The TypedSet version returns []string directly without type casting.
func (s *Set) KeysStr() []string {
	var keys []string

	for key := range s.data {
		if strKey, ok := key.(string); ok {
			keys = append(keys, strKey)
		}
	}

	return keys
}

// Remove removes i from s.
// Deprecated: Use TypedSet[T].Remove() instead for better type safety and performance.
func (s *Set) Remove(i any) {
	s.validate(i)
	delete(s.data, i)
}

// Count returns the number of items in s.
// Deprecated: Use TypedSet[T].Count() instead for better type safety and performance.
func (s *Set) Count() int {
	return len(s.data)
}

func (s *Set) add(i any) {
	switch s.tp {
	case unmanaged:
		// do nothing
	case untyped:
		s.setType(i)
	default:
		s.validate(i)
	}
	s.data[i] = lang.Placeholder
}

func (s *Set) setType(i any) {
	// s.tp can only be untyped here
	switch i.(type) {
	case int:
		s.tp = intType
	case int64:
		s.tp = int64Type
	case uint:
		s.tp = uintType
	case uint64:
		s.tp = uint64Type
	case string:
		s.tp = stringType
	}
}

func (s *Set) validate(i any) {
	if s.tp == unmanaged {
		return
	}

	switch i.(type) {
	case int:
		if s.tp != intType {
			logx.Errorf("element is int, but set contains elements with type %d", s.tp)
		}
	case int64:
		if s.tp != int64Type {
			logx.Errorf("element is int64, but set contains elements with type %d", s.tp)
		}
	case uint:
		if s.tp != uintType {
			logx.Errorf("element is uint, but set contains elements with type %d", s.tp)
		}
	case uint64:
		if s.tp != uint64Type {
			logx.Errorf("element is uint64, but set contains elements with type %d", s.tp)
		}
	case string:
		if s.tp != stringType {
			logx.Errorf("element is string, but set contains elements with type %d", s.tp)
		}
	}
}
