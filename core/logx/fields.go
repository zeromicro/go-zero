package logx

import "context"

var fieldsContextKey contextKey

type contextKey struct{}

// WithFields returns a new context with the given fields.
func WithFields(ctx context.Context, fields ...LogField) context.Context {
	if val := ctx.Value(fieldsContextKey); val != nil {
		if arr, ok := val.([]LogField); ok {
			return context.WithValue(ctx, fieldsContextKey, append(arr, fields...))
		}
	}

	return context.WithValue(ctx, fieldsContextKey, fields)
}
