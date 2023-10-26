package clientinterceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// TimeoutCallOption is a call option that controls timeout.
type TimeoutCallOption struct {
	grpc.EmptyCallOption
	timeout time.Duration
}

// TimeoutInterceptor is an interceptor that controls timeout.
func TimeoutInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		t := getTimeoutFromCallOptions(opts, timeout)
		if t <= 0 {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		ctx, cancel := context.WithTimeout(ctx, t)
		defer cancel()

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// WithCallTimeout returns a call option that controls method call timeout.
func WithCallTimeout(timeout time.Duration) grpc.CallOption {
	return TimeoutCallOption{
		timeout: timeout,
	}
}

func getTimeoutFromCallOptions(opts []grpc.CallOption, defaultTimeout time.Duration) time.Duration {
	for _, opt := range opts {
		if o, ok := opt.(TimeoutCallOption); ok {
			return o.timeout
		}
	}

	return defaultTimeout
}
