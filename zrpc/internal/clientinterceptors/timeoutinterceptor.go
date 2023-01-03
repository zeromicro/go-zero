package clientinterceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// TimeoutInterceptor is an interceptor that controls timeout.
func TimeoutInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		timeout = getTimeoutByCallOptions(opts, timeout)

		if timeout <= 0 {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func getTimeoutByCallOptions(callOptions []grpc.CallOption, defaultTimeout time.Duration) time.Duration {
	for _, callOption := range callOptions {
		if o, ok := callOption.(TimeoutCallOption); ok {
			return o.timeout
		}
	}

	return defaultTimeout
}

type TimeoutCallOption struct {
	grpc.EmptyCallOption

	timeout time.Duration
}

func WithTimeoutCallOption(timeout time.Duration) grpc.CallOption {
	return TimeoutCallOption{
		timeout: timeout,
	}
}
