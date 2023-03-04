package collection

import "sync"

// A Queue is a FIFO queue.
type Queue struct {
	lock     sync.Mutex
	elements []any
	size     int
	head     int
	tail     int
	count    int
}

// NewQueue returns a Queue object.
func NewQueue(size int) *Queue {
	return &Queue{
		elements: make([]any, size),
		size:     size,
	}
}

// Empty checks if q is empty.
func (q *Queue) Empty() bool {
	q.lock.Lock()
	empty := q.count == 0
	q.lock.Unlock()

	return empty
}

// Put puts element into q at the last position.
func (q *Queue) Put(element any) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.head == q.tail && q.count > 0 {
		nodes := make([]any, len(q.elements)+q.size)
		copy(nodes, q.elements[q.head:])
		copy(nodes[len(q.elements)-q.head:], q.elements[:q.head])
		q.head = 0
		q.tail = len(q.elements)
		q.elements = nodes
	}

	q.elements[q.tail] = element
	q.tail = (q.tail + 1) % len(q.elements)
	q.count++
}

// Take takes the first element out of q if not empty.
func (q *Queue) Take() (any, bool) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.count == 0 {
		return nil, false
	}

	element := q.elements[q.head]
	q.head = (q.head + 1) % len(q.elements)
	q.count--

	return element, true
}
