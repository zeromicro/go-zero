package contextx

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShrinkDeadlineLess(t *testing.T) {
	deadline := time.Now().Add(time.Second)
	ctx, _ := context.WithDeadline(context.Background(), deadline)
	ctx, _ = ShrinkDeadline(ctx, time.Minute)
	dl, ok := ctx.Deadline()
	assert.True(t, ok)
	assert.Equal(t, deadline, dl)
}

func TestShrinkDeadlineMore(t *testing.T) {
	deadline := time.Now().Add(time.Minute)
	ctx, _ := context.WithDeadline(context.Background(), deadline)
	ctx, _ = ShrinkDeadline(ctx, time.Second)
	dl, ok := ctx.Deadline()
	assert.True(t, ok)
	assert.True(t, dl.Before(deadline))
}
