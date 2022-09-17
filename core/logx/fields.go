package logx

import "context"

var fieldsContextKey contextKey

type contextKey struct{}

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
