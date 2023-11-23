package collection

import (
	"sync"
)

// A Ring can be used as fixed size ring.
type Ring struct {
	elements []any
	index    int
	lock     sync.RWMutex
}

// NewRing returns a Ring object with the given size n.
func NewRing(n int) *Ring {
	if n < 1 {
		panic("n should be greater than 0")
	}

	return &Ring{
		elements: make([]any, n),
	}
}

// Add adds v into r.
func (r *Ring) Add(v any) {
	r.lock.Lock()
	defer r.lock.Unlock()

	ringLength := len(r.elements)

	r.elements[r.index%ringLength] = v
	r.index++

	// prevent ring index overflow
	if r.index/ringLength >= 2 {
		r.index = r.index - ringLength
	}
}

// Take takes all items from r.
func (r *Ring) Take() []any {
	r.lock.RLock()
	defer r.lock.RUnlock()

	var size int
	var start int
	ringLength := len(r.elements)

	if r.index > ringLength {
		size = ringLength
		start = r.index % ringLength
	} else {
		size = r.index
	}

	elements := make([]any, size)
	for i := 0; i < size; i++ {
		elements[i] = r.elements[(start+i)%ringLength]
	}

	return elements
}
