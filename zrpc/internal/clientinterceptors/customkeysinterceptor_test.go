package clientinterceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/zeromicro/go-zero/core/metainfo"
)

func TestUnaryCustomKeysInterceptor(t *testing.T) {
	testKey := metainfo.PrefixPass + "test"
	testCustomKeysData := map[string]string{testKey: testKey}
	ctx := metainfo.CustomKeysMapPropagator.Extract(context.Background(),
		metainfo.GrpcHeaderCarrier(metadata.New(testCustomKeysData)))

	err := UnaryCustomKeysInterceptor(ctx, "/foo", nil, nil, new(grpc.ClientConn),
		func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
			opts ...grpc.CallOption) error {
			md, ok := metadata.FromOutgoingContext(ctx)
			assert.True(t, ok)
			assert.Equal(t, testKey, md.Get(testKey)[0])
			return nil
		})
	assert.Nil(t, err)
}
