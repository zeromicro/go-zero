package collection

import "github.com/zeromicro/go-zero/core/lang"

// Set is a type-safe generic set collection.
// It's not thread-safe, use with synchronization for concurrent access.
type Set[T comparable] struct {
	data map[T]lang.PlaceholderType
}

// NewSet returns a new type-safe set.
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		data: make(map[T]lang.PlaceholderType),
	}
}

// Add adds items to the set. Duplicates are automatically ignored.
func (s *Set[T]) Add(items ...T) {
	for _, item := range items {
		s.data[item] = lang.Placeholder
	}
}

// Clear removes all items from the set.
func (s *Set[T]) Clear() {
	clear(s.data)
}

// Contains checks if an item exists in the set.
func (s *Set[T]) Contains(item T) bool {
	_, ok := s.data[item]
	return ok
}

// Count returns the number of items in the set.
func (s *Set[T]) Count() int {
	return len(s.data)
}

// Keys returns all elements in the set as a slice.
func (s *Set[T]) Keys() []T {
	keys := make([]T, 0, len(s.data))
	for key := range s.data {
		keys = append(keys, key)
	}
	return keys
}

// Remove removes an item from the set.
func (s *Set[T]) Remove(item T) {
	delete(s.data, item)
}
