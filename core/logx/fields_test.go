package logx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextWithFields(t *testing.T) {
	ctx := ContextWithFields(context.Background(), Field("a", 1), Field("b", 2))
	vals := ctx.Value(fieldsContextKey)
	assert.NotNil(t, vals)
	fields, ok := vals.([]LogField)
	assert.True(t, ok)
	assert.EqualValues(t, []LogField{Field("a", 1), Field("b", 2)}, fields)
}

func TestWithFieldsAppend(t *testing.T) {
	var dummyKey struct{}
	ctx := context.WithValue(context.Background(), dummyKey, "dummy")
	ctx = ContextWithFields(ctx, Field("a", 1), Field("b", 2))
	ctx = ContextWithFields(ctx, Field("c", 3), Field("d", 4))
	vals := ctx.Value(fieldsContextKey)
	assert.NotNil(t, vals)
	fields, ok := vals.([]LogField)
	assert.True(t, ok)
	assert.Equal(t, "dummy", ctx.Value(dummyKey))
	assert.EqualValues(t, []LogField{
		Field("a", 1),
		Field("b", 2),
		Field("c", 3),
		Field("d", 4),
	}, fields)
}
