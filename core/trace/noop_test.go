package trace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoopSpan_Fork(t *testing.T) {
	ctx, span := emptyNoopSpan.Fork(context.Background(), "", "")
	assert.Equal(t, emptyNoopSpan, span)
	assert.Equal(t, context.Background(), ctx)
}

func TestNoopSpan_Follow(t *testing.T) {
	ctx, span := emptyNoopSpan.Follow(context.Background(), "", "")
	assert.Equal(t, emptyNoopSpan, span)
	assert.Equal(t, context.Background(), ctx)
}

func TestNoopSpan(t *testing.T) {
	emptyNoopSpan.Visit(func(key, val string) bool {
		assert.Fail(t, "should not go here")
		return true
	})

	ctx, span := emptyNoopSpan.Follow(context.Background(), "", "")
	assert.Equal(t, context.Background(), ctx)
	assert.Equal(t, "", span.TraceId())
	assert.Equal(t, "", span.SpanId())
}
