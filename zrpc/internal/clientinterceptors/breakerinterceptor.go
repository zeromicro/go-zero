package clientinterceptors

import (
	"context"

	"github.com/zeromicro/go-zero/zrpc/internal/balancer/breaker"
	"google.golang.org/grpc"
)

// BreakerInterceptor is an interceptor that enables the circuit breaker.
func BreakerInterceptor(ctx context.Context, method string, req, reply any,
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	return invoker(breaker.WithBreaker(ctx), method, req, reply, cc, opts...)
}
