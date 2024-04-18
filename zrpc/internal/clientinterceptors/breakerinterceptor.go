package clientinterceptors

import (
	"context"
	"errors"
	"path"
	"time"

	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/timex"
	"github.com/zeromicro/go-zero/zrpc/internal/codes"
	"google.golang.org/grpc"
)

const minTimeout = time.Millisecond * 100

// BreakerInterceptor is an interceptor that acts as a circuit breaker.
func BreakerInterceptor(ctx context.Context, method string, req, reply any,
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	breakerName := path.Join(cc.Target(), method)
	start := timex.Now()
	return breaker.DoWithAcceptable(breakerName, func() error {
		return invoker(ctx, method, req, reply, cc, opts...)
	}, func(err error) bool {
		if codes.Acceptable(err) {
			return true
		}

		return errors.Is(err, context.DeadlineExceeded) && timex.Since(start) < minTimeout
	})
}
