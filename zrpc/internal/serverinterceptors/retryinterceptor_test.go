package serverinterceptors

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/retry"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestRetryInterceptor(t *testing.T) {
	t.Run("retries exceeded", func(t *testing.T) {
		interceptor := RetryInterceptor(2)
		ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{retry.AttemptMetadataKey: "3"}))
		resp, err := interceptor(ctx, nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, nil
		})
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("reasonable retries", func(t *testing.T) {
		interceptor := RetryInterceptor(2)
		ctx := metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{retry.AttemptMetadataKey: "2"}))
		resp, err := interceptor(ctx, nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, nil
		})
		assert.NoError(t, err)
		assert.Nil(t, resp)
	})
	t.Run("no retries", func(t *testing.T) {
		interceptor := RetryInterceptor(0)
		resp, err := interceptor(context.Background(), nil, nil, func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, nil
		})
		assert.NoError(t, err)
		assert.Nil(t, resp)
	})

}
