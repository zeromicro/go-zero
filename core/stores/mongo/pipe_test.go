package mongo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/breaker"
)

func TestRejectedPipe_All(t *testing.T) {
	assert.Equal(t, breaker.ErrServiceUnavailable, new(rejectedPipe).All(nil))
}

func TestRejectedPipe_AllowDiskUse(t *testing.T) {
	var p rejectedPipe
	assert.Equal(t, p, p.AllowDiskUse())
}

func TestRejectedPipe_Batch(t *testing.T) {
	var p rejectedPipe
	assert.Equal(t, p, p.Batch(1))
}

func TestRejectedPipe_Collation(t *testing.T) {
	var p rejectedPipe
	assert.Equal(t, p, p.Collation(nil))
}

func TestRejectedPipe_Explain(t *testing.T) {
	assert.Equal(t, breaker.ErrServiceUnavailable, new(rejectedPipe).Explain(nil))
}

func TestRejectedPipe_Iter(t *testing.T) {
	assert.EqualValues(t, rejectedIter{}, new(rejectedPipe).Iter())
}

func TestRejectedPipe_One(t *testing.T) {
	assert.Equal(t, breaker.ErrServiceUnavailable, new(rejectedPipe).One(nil))
}

func TestRejectedPipe_SetMaxTime(t *testing.T) {
	var p rejectedPipe
	assert.Equal(t, p, p.SetMaxTime(0))
}
