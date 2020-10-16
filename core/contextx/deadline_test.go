package contextx

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShrinkDeadlineLess(t *testing.T) {
	deadline := time.Now().Add(time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	ctx, cancel = ShrinkDeadline(ctx, time.Minute)
	defer cancel()
	dl, ok := ctx.Deadline()
	assert.True(t, ok)
	assert.Equal(t, deadline, dl)
}

func TestShrinkDeadlineMore(t *testing.T) {
	deadline := time.Now().Add(time.Minute)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	ctx, cancel = ShrinkDeadline(ctx, time.Second)
	defer cancel()
	dl, ok := ctx.Deadline()
	assert.True(t, ok)
	assert.True(t, dl.Before(deadline))
}
