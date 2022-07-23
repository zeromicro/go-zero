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
	t.Run("no selector", func(t *testing.T) {
		_, err := UnarySelectorInterceptor(context.Background(), nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
			assert.Equal(t, "", selector.FromContext(ctx))
			return nil, nil
		})

		assert.NoError(t, err)
	})

	t.Run("has selector", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{
			"selector": []string{selector.DefaultSelector},
		})
		_, err := UnarySelectorInterceptor(ctx, nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
			assert.Equal(t, selector.DefaultSelector, selector.FromContext(ctx))
			return nil, nil
		})

		assert.NoError(t, err)
	})
}

func TestStreamSelectorInterceptor(t *testing.T) {
	t.Run("no selector", func(t *testing.T) {
		err := StreamSelectorInterceptor(nil, &mockedServerStream{ctx: context.Background()}, nil, func(srv interface{}, stream grpc.ServerStream) error {
			assert.Equal(t, "", selector.FromContext(stream.Context()))
			return nil
		})

		assert.NoError(t, err)
	})

	t.Run("has selector", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{
			"selector": []string{selector.DefaultSelector},
		})
		err := StreamSelectorInterceptor(nil, &mockedServerStream{ctx: ctx}, nil, func(srv interface{}, stream grpc.ServerStream) error {
			assert.Equal(t, selector.DefaultSelector, selector.FromContext(stream.Context()))
			return nil
		})

		assert.NoError(t, err)
	})
}
