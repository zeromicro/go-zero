package serverinterceptors

import (
	"context"
	"time"

	"github.com/tal-tech/go-zero/core/contextx"
	"google.golang.org/grpc"
)

// UnaryTimeoutInterceptor returns a func that sets timeout to incoming unary requests.
func UnaryTimeoutInterceptor(timeout time.Duration) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx, cancel := contextx.ShrinkDeadline(ctx, timeout)
		defer cancel()
		return handler(ctx, req)
	}
}
