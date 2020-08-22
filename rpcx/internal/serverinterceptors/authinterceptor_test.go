package serverinterceptors

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/stores/redis"
	"github.com/tal-tech/go-zero/rpcx/internal/auth"
	"google.golang.org/grpc/metadata"
)

func TestUnaryAuthorizeInterceptor(t *testing.T) {
	tests := []struct {
		name   string
		strict bool
	}{
		{
			name:   "strict=true",
			strict: true,
		},
		{
			name:   "strict=false",
			strict: false,
		},
	}

	r := miniredis.NewMiniRedis()
	assert.Nil(t, r.Start())
	defer r.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := redis.NewRedis(r.Addr(), redis.NodeType)
			authenticator, err := auth.NewAuthenticator(store, "apps", test.strict)
			assert.Nil(t, err)
			interceptor := UnaryAuthorizeInterceptor(authenticator)
			md := metadata.New(map[string]string{
				"app":   "name",
				"token": "key",
			})
			ctx := metadata.NewIncomingContext(context.Background(), md)
			_, err = interceptor(ctx, nil, nil,
				func(ctx context.Context, req interface{}) (interface{}, error) {
					return nil, nil
				})
			if test.strict {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
