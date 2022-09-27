package logx

import (
	"context"
	"sync/atomic"
)

var (
	fieldsContextKey contextKey
	// we store *[]LogField as the value, because []LogField is not comparable,
	// we need to use CompareAndSwap to compare the stored value.
	globalFields atomic.Value
)

type contextKey struct{}

// AddGlobalFields adds global fields.
func AddGlobalFields(fields ...LogField) {
	for {
		old := globalFields.Load()
		if old == nil {
			val := append([]LogField(nil), fields...)
			if globalFields.CompareAndSwap(old, &val) {
				return
			}
		} else {
			oldFields := old.(*[]LogField)
			val := append(*oldFields, fields...)
			if globalFields.CompareAndSwap(old, &val) {
				return
			}
		}
	}
}

// ContextWithFields returns a new context with the given fields.
func ContextWithFields(ctx context.Context, fields ...LogField) context.Context {
	if val := ctx.Value(fieldsContextKey); val != nil {
		if arr, ok := val.([]LogField); ok {
			return context.WithValue(ctx, fieldsContextKey, append(arr, fields...))
		}
	}

	return context.WithValue(ctx, fieldsContextKey, fields)
}

// WithFields returns a new logger with the given fields.
// deprecated: use ContextWithFields instead.
func WithFields(ctx context.Context, fields ...LogField) context.Context {
	return ContextWithFields(ctx, fields...)
}
