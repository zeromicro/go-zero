package clientinterceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/md"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestUnaryMdInterceptor(t *testing.T) {
	t.Run("has defaultMd", func(t *testing.T) {
		interceptor := UnaryMdInterceptor(map[string][]string{"colors": {"v2"}})
		err := interceptor(context.Background(), "foo", nil, nil, nil, func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			assert.EqualValues(t, map[string][]string{"colors": {"v2"}}, md.FromContext(ctx))

			m, b := metadata.FromOutgoingContext(ctx)
			assert.True(t, b)
			assert.Equal(t, `{"colors":["v2"]}`, m.Get("metadata")[0])
			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("no defaultMd", func(t *testing.T) {
		interceptor := UnaryMdInterceptor(nil)
		err := interceptor(context.Background(), "foo", nil, nil, nil, func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			assert.Equal(t, md.Metadata{}, md.FromContext(ctx))

			m, b := metadata.FromOutgoingContext(ctx)
			assert.True(t, b)
			assert.Equal(t, "{}", m.Get("metadata")[0])
			return nil
		})
		assert.NoError(t, err)
	})

}

func TestStreamMdInterceptor(t *testing.T) {
	t.Run("has defaultMd", func(t *testing.T) {
		interceptor := StreamMdInterceptor(map[string][]string{"colors": {"v1"}})
		_, err := interceptor(context.Background(), nil, nil, "foo", func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
			assert.EqualValues(t, map[string][]string{"colors": {"v1"}}, md.FromContext(ctx))

			m, b := metadata.FromOutgoingContext(ctx)
			assert.True(t, b)
			assert.Equal(t, `{"colors":["v1"]}`, m.Get("metadata")[0])
			return nil, nil
		})
		assert.NoError(t, err)
	})

	t.Run("no defaultMd", func(t *testing.T) {
		interceptor := StreamMdInterceptor(map[string][]string{})
		_, err := interceptor(context.Background(), nil, nil, "foo", func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
			assert.EqualValues(t, map[string][]string{}, md.FromContext(ctx))

			m, b := metadata.FromOutgoingContext(ctx)
			assert.True(t, b)
			assert.Equal(t, "{}", m.Get("metadata")[0])
			return nil, nil
		})
		assert.NoError(t, err)
	})
}
