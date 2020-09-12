package contextx

import (
	"context"

	"github.com/tal-tech/go-zero/core/mapping"
)

const contextTagKey = "ctx"

var unmarshaler = mapping.NewUnmarshaler(contextTagKey)

type contextValuer struct {
	context.Context
}

func (cv contextValuer) Value(key string) (interface{}, bool) {
	v := cv.Context.Value(key)
	return v, v != nil
}

func For(ctx context.Context, v interface{}) error {
	return unmarshaler.UnmarshalValuer(contextValuer{
		Context: ctx,
	}, v)
}
