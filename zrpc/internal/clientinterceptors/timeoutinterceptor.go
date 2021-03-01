package clientinterceptors

import (
	"context"
	"time"

	"github.com/tal-tech/go-zero/core/contextx"
	"google.golang.org/grpc"
)

// TimeoutInterceptor is an interceptor that controls timeout.
func TimeoutInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if timeout <= 0 {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		ctx, cancel := contextx.ShrinkDeadline(ctx, timeout)
		defer cancel()
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
