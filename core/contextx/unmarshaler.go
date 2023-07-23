package contextx

import (
	"context"

	"github.com/zeromicro/go-zero/core/mapping"
)

const contextTagKey = "ctx"

var unmarshaler = mapping.NewUnmarshaler(contextTagKey)

type contextValuer struct {
	context.Context
}

func (cv contextValuer) Value(key string) (any, bool) {
	v := cv.Context.Value(key)
	return v, v != nil
}

// For unmarshals ctx into v.
func For(ctx context.Context, v any) error {
	return unmarshaler.UnmarshalValuer(contextValuer{
		Context: ctx,
	}, v)
}
