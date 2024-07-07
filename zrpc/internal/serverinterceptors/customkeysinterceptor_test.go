package serverinterceptors

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
	ctx := metadata.NewIncomingContext(context.Background(), metadata.New(testCustomKeysData))

	_, err := UnaryCustomKeysInterceptor()(ctx, nil, &grpc.UnaryServerInfo{
		FullMethod: "/",
	}, func(ctx context.Context, req any) (any, error) {
		md := metainfo.GetMapFromContext(ctx)
		assert.Equal(t, testCustomKeysData, md)
		// todo(cong): test logger fields.
		return nil, nil
	})

	assert.Nil(t, err)
}
