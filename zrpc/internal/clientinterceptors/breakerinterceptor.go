package clientinterceptors

import (
	"context"
	"path"

	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/zrpc/internal/codes"
	"google.golang.org/grpc"
)

// BreakerInterceptor is an interceptor that acts as a circuit breaker.
//
// Deprecated: use balancer/breaker instead, which provides both service-level
// and instance-level circuit breaker capabilities with unified implementation.
func BreakerInterceptor(ctx context.Context, method string, req, reply any,
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	breakerName := path.Join(cc.Target(), method)
	return breaker.DoWithAcceptableCtx(ctx, breakerName, func() error {
		return invoker(ctx, method, req, reply, cc, opts...)
	}, codes.Acceptable)
}
