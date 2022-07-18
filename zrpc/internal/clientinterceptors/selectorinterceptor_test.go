package clientinterceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/selector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestUnarySelectorInterceptor(t *testing.T) {
	t.Run("use all defaults", func(t *testing.T) {
		interceptor := UnarySelectorInterceptor("defaultSelector", []string{"v1", "v2"})
		err := interceptor(context.Background(), "foo", nil, nil, nil, func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			assert.Equal(t, "defaultSelector", selector.SelectorFromContext(ctx))
			assert.Equal(t, []string{"v1", "v2"}, selector.ColorsFromContext(ctx).Colors())

			md, b := metadata.FromOutgoingContext(ctx)
			assert.True(t, b)
			assert.Equal(t, "defaultSelector", md.Get("selector")[0])
			assert.Equal(t, []string{"v1", "v2"}, md.Get("colors"))
			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("default selector", func(t *testing.T) {
		interceptor := UnarySelectorInterceptor("defaultSelector", []string{})
		err := interceptor(context.Background(), "foo", nil, nil, nil, func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			assert.Equal(t, "defaultSelector", selector.SelectorFromContext(ctx))
			assert.Equal(t, []string(nil), selector.ColorsFromContext(ctx).Colors())

			md, b := metadata.FromOutgoingContext(ctx)
			assert.True(t, b)
			assert.Equal(t, "defaultSelector", md.Get("selector")[0])
			assert.Equal(t, []string(nil), md.Get("colors"))
			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("default colors", func(t *testing.T) {
		interceptor := UnarySelectorInterceptor("", []string{"v1", ""})
		err := interceptor(context.Background(), "foo", nil, nil, nil, func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			assert.Equal(t, "", selector.SelectorFromContext(ctx))
			assert.Equal(t, []string{"v1"}, selector.ColorsFromContext(ctx).Colors())

			md, b := metadata.FromOutgoingContext(ctx)
			assert.True(t, b)
			assert.Equal(t, []string(nil), md.Get("selector"))
			assert.Equal(t, []string{"v1"}, md.Get("colors"))
			return nil
		})
		assert.NoError(t, err)
	})
}

func TestStreamSelectorInterceptor(t *testing.T) {
	t.Run("use all defaults", func(t *testing.T) {
		interceptor := StreamSelectorInterceptor("defaultSelector", []string{"v1", "v2"})
		_, err := interceptor(context.Background(), nil, nil, "foo", func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
			assert.Equal(t, "defaultSelector", selector.SelectorFromContext(ctx))
			assert.Equal(t, []string{"v1", "v2"}, selector.ColorsFromContext(ctx).Colors())

			md, b := metadata.FromOutgoingContext(ctx)
			assert.True(t, b)
			assert.Equal(t, "defaultSelector", md.Get("selector")[0])
			assert.Equal(t, []string{"v1", "v2"}, md.Get("colors"))
			return nil, nil
		})
		assert.NoError(t, err)
	})

	t.Run("default selector", func(t *testing.T) {
		interceptor := StreamSelectorInterceptor("defaultSelector", []string{})
		_, err := interceptor(context.Background(), nil, nil, "foo", func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
			assert.Equal(t, "defaultSelector", selector.SelectorFromContext(ctx))
			assert.Equal(t, []string(nil), selector.ColorsFromContext(ctx).Colors())

			md, b := metadata.FromOutgoingContext(ctx)
			assert.True(t, b)
			assert.Equal(t, "defaultSelector", md.Get("selector")[0])
			assert.Equal(t, []string(nil), md.Get("colors"))
			return nil, nil
		})
		assert.NoError(t, err)
	})

	t.Run("default colors", func(t *testing.T) {
		interceptor := StreamSelectorInterceptor("", []string{"v1", ""})
		_, err := interceptor(context.Background(), nil, nil, "foo", func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
			assert.Equal(t, "", selector.SelectorFromContext(ctx))
			assert.Equal(t, []string{"v1"}, selector.ColorsFromContext(ctx).Colors())

			md, b := metadata.FromOutgoingContext(ctx)
			assert.True(t, b)
			assert.Equal(t, []string(nil), md.Get("selector"))
			assert.Equal(t, []string{"v1"}, md.Get("colors"))
			return nil, nil
		})
		assert.NoError(t, err)
	})
}
