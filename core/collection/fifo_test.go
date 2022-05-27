package collection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFifo(t *testing.T) {
	elements := [][]byte{
		[]byte("hello"),
		[]byte("world"),
		[]byte("again"),
	}
	queue := NewQueue(8)
	for i := range elements {
		queue.Put(elements[i])
	}

	for _, element := range elements {
		body, ok := queue.Take()
		assert.True(t, ok)
		assert.Equal(t, string(element), string(body.([]byte)))
	}
}

func TestTakeTooMany(t *testing.T) {
	elements := [][]byte{
		[]byte("hello"),
		[]byte("world"),
		[]byte("again"),
	}
	queue := NewQueue(8)
	for i := range elements {
		queue.Put(elements[i])
	}

	for range elements {
		queue.Take()
	}

	assert.True(t, queue.Empty())
	_, ok := queue.Take()
	assert.False(t, ok)
}

func TestPutMore(t *testing.T) {
	elements := [][]byte{
		[]byte("hello"),
		[]byte("world"),
		[]byte("again"),
	}
	queue := NewQueue(2)
	for i := range elements {
		queue.Put(elements[i])
	}

	for _, element := range elements {
		body, ok := queue.Take()
		assert.True(t, ok)
		assert.Equal(t, string(element), string(body.([]byte)))
	}
}

func TestPutMoreWithHeaderNotZero(t *testing.T) {
	elements := [][]byte{
		[]byte("hello"),
		[]byte("world"),
		[]byte("again"),
	}
	queue := NewQueue(4)
	for i := range elements {
		queue.Put(elements[i])
	}

	// take 1
	body, ok := queue.Take()
	assert.True(t, ok)
	element, ok := body.([]byte)
	assert.True(t, ok)
	assert.Equal(t, element, []byte("hello"))

	// put more
	queue.Put([]byte("b4"))
	queue.Put([]byte("b5")) // will store in elements[0]
	queue.Put([]byte("b6")) // cause expansion

	results := [][]byte{
		[]byte("world"),
		[]byte("again"),
		[]byte("b4"),
		[]byte("b5"),
		[]byte("b6"),
	}

	for _, element := range results {
		body, ok := queue.Take()
		assert.True(t, ok)
		assert.Equal(t, string(element), string(body.([]byte)))
	}
}
