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
			name:     "valid read-primary mode",
			mode:     readPrimaryMode,
			expected: true,
		},
		{
			name:     "valid read-replica mode",
			mode:     readReplicaMode,
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
	readPrimaryCtx := WithReadPrimary(ctx)

	val := readPrimaryCtx.Value(readWriteModeKey{})
	assert.Equal(t, readPrimaryMode, val)

	readReplicaCtx := WithReadReplica(ctx)
	val = readReplicaCtx.Value(readWriteModeKey{})
	assert.Equal(t, readReplicaMode, val)
}

func TestWithWriteMode(t *testing.T) {
	ctx := context.Background()
	writeCtx := WithWrite(ctx)

	val := writeCtx.Value(readWriteModeKey{})
	assert.Equal(t, writeMode, val)
}

func TestGetReadWriteMode(t *testing.T) {
	t.Run("valid read-primary mode", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), readWriteModeKey{}, readPrimaryMode)
		assert.Equal(t, readPrimaryMode, getReadWriteMode(ctx))
	})

	t.Run("valid read-replica mode", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), readWriteModeKey{}, readReplicaMode)
		assert.Equal(t, readReplicaMode, getReadWriteMode(ctx))
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

func TestUsePrimary(t *testing.T) {
	t.Run("context with read-replica mode", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), readWriteModeKey{}, readReplicaMode)
		assert.False(t, usePrimary(ctx))
	})

	t.Run("context with read-primary mode", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), readWriteModeKey{}, readPrimaryMode)
		assert.True(t, usePrimary(ctx))
	})

	t.Run("context with write mode", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), readWriteModeKey{}, writeMode)
		assert.True(t, usePrimary(ctx))
	})

	t.Run("context with invalid mode", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), readWriteModeKey{}, readWriteMode("invalid"))
		assert.True(t, usePrimary(ctx))
	})

	t.Run("context with no mode set", func(t *testing.T) {
		ctx := context.Background()
		assert.True(t, usePrimary(ctx))
	})
}

func TestWithModeTwice(t *testing.T) {
	ctx := context.Background()
	ctx = WithReadPrimary(ctx)
	writeCtx := WithWrite(ctx)

	val := writeCtx.Value(readWriteModeKey{})
	assert.Equal(t, writeMode, val)
}
