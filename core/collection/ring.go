package collection

import "sync"

type Ring struct {
	elements []interface{}
	index    int
	lock     sync.Mutex
}

func NewRing(n int) *Ring {
	if n < 1 {
		panic("n should be greater than 0")
	}

	return &Ring{
		elements: make([]interface{}, n),
	}
}

func (r *Ring) Add(v interface{}) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.elements[r.index%len(r.elements)] = v
	r.index++
}

func (r *Ring) Take() []interface{} {
	r.lock.Lock()
	defer r.lock.Unlock()

	var size int
	var start int
	if r.index > len(r.elements) {
		size = len(r.elements)
		start = r.index % len(r.elements)
	} else {
		size = r.index
	}

	elements := make([]interface{}, size)
	for i := 0; i < size; i++ {
		elements[i] = r.elements[(start+i)%len(r.elements)]
	}

	return elements
}
