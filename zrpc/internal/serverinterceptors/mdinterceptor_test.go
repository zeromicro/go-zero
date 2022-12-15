package serverinterceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/md"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestUnaryMdInterceptor(t *testing.T) {
	t.Run("no md", func(t *testing.T) {
		_, err := UnaryMdInterceptor(context.Background(), nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
			assert.Len(t, md.FromContext(ctx), 0)
			return nil, nil
		})

		assert.NoError(t, err)
	})

	t.Run("has md", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{
			"metadata": []string{`{"colors":["v1"]}`},
		})
		_, err := UnaryMdInterceptor(ctx, nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
			assert.Equal(t, []string{"v1"}, md.ValuesFromContext(ctx, "colors"))
			return nil, nil
		})

		assert.NoError(t, err)
	})
}

func TestStreamMdInterceptor(t *testing.T) {
	t.Run("no md", func(t *testing.T) {
		err := StreamMdInterceptor(nil, &mockedServerStream{ctx: context.Background()}, nil, func(srv interface{}, stream grpc.ServerStream) error {
			assert.Len(t, md.FromContext(stream.Context()), 0)
			return nil
		})

		assert.NoError(t, err)
	})

	t.Run("has md", func(t *testing.T) {
		ctx := metadata.NewIncomingContext(context.Background(), metadata.MD{
			"metadata": []string{`{"colors":["v1"]}`},
		})
		err := StreamMdInterceptor(nil, &mockedServerStream{ctx: ctx}, nil, func(srv interface{}, stream grpc.ServerStream) error {
			assert.Equal(t, []string{"v1"}, md.ValuesFromContext(stream.Context(), "colors"))
			return nil
		})

		assert.NoError(t, err)
	})
}
