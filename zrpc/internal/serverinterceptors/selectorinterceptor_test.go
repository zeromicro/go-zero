package serverinterceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/selector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestUnarySelectorInterceptor(t *testing.T) {
	t.Run("no select and colors", func(t *testing.T) {
		_, err := UnarySelectorInterceptor(context.Background(), nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
			assert.Equal(t, "", selector.SelectorFromContext(ctx))
			assert.Len(t, selector.ColorsFromContext(ctx).Colors(), 0)
			return nil, nil
		})

		assert.NoError(t, err)
	})

	t.Run("has select and colors", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{
			"selector": []string{selector.DefaultSelector},
			"colors":   []string{"v1", "v2"},
		})
		_, err := UnarySelectorInterceptor(ctx, nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
			assert.Equal(t, selector.DefaultSelector, selector.SelectorFromContext(ctx))
			assert.Equal(t, []string{"v1", "v2"}, selector.ColorsFromContext(ctx).Colors())
			return nil, nil
		})

		assert.NoError(t, err)
	})
}

func TestStreamSelectorInterceptor(t *testing.T) {
	t.Run("no select and colors", func(t *testing.T) {
		err := StreamSelectorInterceptor(nil, &mockedServerStream{ctx: context.Background()}, nil, func(srv interface{}, stream grpc.ServerStream) error {
			assert.Equal(t, "", selector.SelectorFromContext(stream.Context()))
			assert.Len(t, selector.ColorsFromContext(stream.Context()).Colors(), 0)
			return nil
		})

		assert.NoError(t, err)
	})

	t.Run("has select and colors", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{
			"selector": []string{selector.DefaultSelector},
			"colors":   []string{"v1", "v2"},
		})
		err := StreamSelectorInterceptor(nil, &mockedServerStream{ctx: ctx}, nil, func(srv interface{}, stream grpc.ServerStream) error {
			assert.Equal(t, selector.DefaultSelector, selector.SelectorFromContext(stream.Context()))
			assert.Equal(t, []string{"v1", "v2"}, selector.ColorsFromContext(stream.Context()).Colors())
			return nil
		})

		assert.NoError(t, err)
	})
}
