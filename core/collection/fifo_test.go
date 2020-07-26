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
