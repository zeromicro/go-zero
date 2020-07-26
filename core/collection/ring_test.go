package collection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRingLess(t *testing.T) {
	ring := NewRing(5)
	for i := 0; i < 3; i++ {
		ring.Add(i)
	}
	elements := ring.Take()
	assert.ElementsMatch(t, []interface{}{0, 1, 2}, elements)
}

func TestRingMore(t *testing.T) {
	ring := NewRing(5)
	for i := 0; i < 11; i++ {
		ring.Add(i)
	}
	elements := ring.Take()
	assert.ElementsMatch(t, []interface{}{6, 7, 8, 9, 10}, elements)
}
