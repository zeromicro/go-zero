package syncx

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManagedResource(t *testing.T) {
	var count int32
	resource := NewManagedResource(func() any {
		return atomic.AddInt32(&count, 1)
	}, func(a, b any) bool {
		return a == b
	})

	assert.Equal(t, resource.Take(), resource.Take())
	old := resource.Take()
	resource.MarkBroken(old)
	assert.NotEqual(t, old, resource.Take())
}
