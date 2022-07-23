package clientinterceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/selector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestStreamSelectorInterceptor(t *testing.T) {
	t.Run("defaultSelector", func(t *testing.T) {
		interceptor := StreamSelectorInterceptor("defaultSelector")
		_, err := interceptor(context.Background(), nil, nil, "foo", func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
			assert.Equal(t, "defaultSelector", selector.FromContext(ctx))

			m, b := metadata.FromOutgoingContext(ctx)
			assert.True(t, b)
			assert.Equal(t, "defaultSelector", m.Get("selector")[0])
			return nil, nil
		})
		assert.NoError(t, err)
	})

	t.Run("mockSelector", func(t *testing.T) {
		interceptor := StreamSelectorInterceptor("defaultSelector")
		_, err := interceptor(selector.NewContext(context.Background(), "mockSelector"), nil, nil, "foo", func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
			assert.Equal(t, "mockSelector", selector.FromContext(ctx))

			m, b := metadata.FromOutgoingContext(ctx)
			assert.True(t, b)
			assert.Equal(t, "mockSelector", m.Get("selector")[0])
			return nil, nil
		})
		assert.NoError(t, err)
	})

}

func TestUnarySelectorInterceptor(t *testing.T) {
	t.Run("defaultSelector", func(t *testing.T) {
		interceptor := UnarySelectorInterceptor("defaultSelector")
		err := interceptor(context.Background(), "foo", nil, nil, nil, func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			assert.Equal(t, "defaultSelector", selector.FromContext(ctx))

			m, b := metadata.FromOutgoingContext(ctx)
			assert.True(t, b)
			assert.Equal(t, "defaultSelector", m.Get("selector")[0])
			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("mockSelector", func(t *testing.T) {
		interceptor := UnarySelectorInterceptor("defaultSelector")
		err := interceptor(selector.NewContext(context.Background(), "mockSelector"), "foo", nil, nil, nil, func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			assert.Equal(t, "mockSelector", selector.FromContext(ctx))

			m, b := metadata.FromOutgoingContext(ctx)
			assert.True(t, b)
			assert.Equal(t, "mockSelector", m.Get("selector")[0])
			return nil
		})
		assert.NoError(t, err)
	})

}
