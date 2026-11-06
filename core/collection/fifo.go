package collection

import "sync"

const queueGrowThreshold = 256

// A Queue is a FIFO queue.
type Queue struct {
	lock     sync.Mutex
	elements []any
	head     int
	tail     int
	count    int
}

// NewQueue returns a Queue object.
func NewQueue(size int) *Queue {
	if size < 1 {
		panic("size must be greater than 0")
	}

	return &Queue{
		elements: make([]any, size),
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

	if q.count == len(q.elements) {
		nodes := make([]any, nextQueueCapacity(len(q.elements)))
		n := copy(nodes, q.elements[q.head:])
		copy(nodes[n:], q.elements[:q.head])
		q.head = 0
		q.tail = q.count
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
	q.elements[q.head] = nil
	q.head = (q.head + 1) % len(q.elements)
	q.count--

	return element, true
}

func nextQueueCapacity(capacity int) int {
	if capacity < queueGrowThreshold {
		return capacity << 1
	}

	// Use a growth curve similar to Go slices: double small queues, then
	// transition smoothly toward 1.25x growth for larger queues.
	return capacity + ((capacity + 3*queueGrowThreshold) >> 2)
}
