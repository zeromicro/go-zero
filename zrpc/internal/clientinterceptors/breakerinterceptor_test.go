package clientinterceptors

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/zrpc/internal/balancer/breaker"
	"google.golang.org/grpc"
)

func TestBreakerInterceptor(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{
			name: "nil",
			err:  nil,
		},
		{
			name: "with error",
			err:  errors.New("mock"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cc := new(grpc.ClientConn)
			err := BreakerInterceptor(context.Background(), "/foo", nil, nil, cc,
				func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
					opts ...grpc.CallOption) error {
					// verify breaker is enabled in context
					assert.True(t, breaker.HasBreaker(ctx))
					return test.err
				})
			assert.Equal(t, test.err, err)
		})
	}
}
