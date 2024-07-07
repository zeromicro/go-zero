package clientinterceptors

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/metainfo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestUnaryCustomKeysInterceptor(t *testing.T) {
	testKey := metainfo.PrefixPass + "test"
	testCustomKeysData := map[string]string{testKey: testKey}

	t.Run("No existing metadata", func(t *testing.T) {
		ctx := context.Background()

		err := UnaryCustomKeysInterceptor(ctx, "/foo", nil, nil, new(grpc.ClientConn),
			func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
				opts ...grpc.CallOption) error {
				md, ok := metadata.FromOutgoingContext(ctx)
				assert.True(t, ok)
				assert.Empty(t, md[testKey])
				return nil
			})
		assert.Nil(t, err)
	})

	t.Run("With existing metadata", func(t *testing.T) {
		ctx := metainfo.CustomKeysMapPropagator.Extract(context.Background(),
			metainfo.GrpcHeaderCarrier(metadata.New(testCustomKeysData)))
		ctx = metadata.AppendToOutgoingContext(ctx, "foo", "bar")

		err := UnaryCustomKeysInterceptor(ctx, "/foo", nil, nil, new(grpc.ClientConn),
			func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
				opts ...grpc.CallOption) error {
				md, ok := metadata.FromOutgoingContext(ctx)
				assert.True(t, ok)
				assert.Equal(t, testKey, md.Get(testKey)[0])
				return nil
			})
		assert.Nil(t, err)
	})

	t.Run("Invoker returns error", func(t *testing.T) {
		expectedErr := errors.New("invoker error")
		ctx := context.Background()

		err := UnaryCustomKeysInterceptor(ctx, "/foo", nil, nil, new(grpc.ClientConn),
			func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
				opts ...grpc.CallOption) error {
				return expectedErr
			})
		assert.Equal(t, expectedErr, err)
	})
}
