package sqlx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValid(t *testing.T) {
	testCases := []struct {
		name     string
		mode     readWriteMode
		expected bool
	}{
		{
			name:     "valid read mode",
			mode:     readMode,
			expected: true,
		},
		{
			name:     "valid write mode",
			mode:     writeMode,
			expected: true,
		},
		{
			name:     "not specified mode (empty)",
			mode:     notSpecifiedMode,
			expected: false,
		},
		{
			name:     "invalid custom string",
			mode:     readWriteMode("delete"),
			expected: false,
		},
		{
			name:     "case sensitive check",
			mode:     readWriteMode("READ"),
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := tc.mode.isValid()
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestWithReadMode(t *testing.T) {
	ctx := context.Background()
	readCtx := WithReadMode(ctx)

	val := readCtx.Value(readWriteModeKey{})
	assert.Equal(t, readMode, val)
}

func TestWithWriteMode(t *testing.T) {
	ctx := context.Background()
	writeCtx := WithWriteMode(ctx)

	val := writeCtx.Value(readWriteModeKey{})
	assert.Equal(t, writeMode, val)
}

func TestGetReadWriteMode(t *testing.T) {
	t.Run("valid read mode", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), readWriteModeKey{}, readMode)
		assert.Equal(t, readMode, getReadWriteMode(ctx))
	})

	t.Run("valid write mode", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), readWriteModeKey{}, writeMode)
		assert.Equal(t, writeMode, getReadWriteMode(ctx))
	})

	t.Run("invalid mode value (wrong type)", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), readWriteModeKey{}, "not-a-mode")
		assert.Equal(t, notSpecifiedMode, getReadWriteMode(ctx))
	})

	t.Run("invalid mode value (wrong value)", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), readWriteModeKey{}, readWriteMode("delete"))
		assert.Equal(t, notSpecifiedMode, getReadWriteMode(ctx))
	})

	t.Run("no mode set", func(t *testing.T) {
		ctx := context.Background()
		assert.Equal(t, notSpecifiedMode, getReadWriteMode(ctx))
	})
}

func TestIsReadonly(t *testing.T) {
	t.Run("context with read mode", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), readWriteModeKey{}, readMode)
		assert.True(t, isReadonly(ctx))
	})

	t.Run("context with write mode", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), readWriteModeKey{}, writeMode)
		assert.False(t, isReadonly(ctx))
	})

	t.Run("context with invalid mode", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), readWriteModeKey{}, readWriteMode("invalid"))
		assert.False(t, isReadonly(ctx))
	})

	t.Run("context with wrong type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), readWriteModeKey{}, "read") // string, not readWriteMode
		assert.False(t, isReadonly(ctx))
	})

	t.Run("context with no mode set", func(t *testing.T) {
		ctx := context.Background()
		assert.False(t, isReadonly(ctx))
	})
}
