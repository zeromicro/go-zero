package sqlx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	const (
		query = "select foo"
		times = 10
	)
	cr := newConcurrentReads()
	assert.NotPanics(t, func() {
		cr.remove(query)
	})

	for i := 0; i < times; i++ {
		cr.add(query)
	}

	ref := cr.reads[query]
	assert.Equal(t, uint32(times), ref.concurrency)

	for i := 0; i < times; i++ {
		cr.remove(query)
	}

	// just removed, not decremented
	assert.Equal(t, uint32(1), ref.concurrency)
	assert.Equal(t, uint32(times), ref.maxConcurrency)
}
