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
		if timeout <= 0 {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		// create channel with buffer size 1 to avoid goroutine leak
		done := make(chan error, 1)
		panicChan := make(chan interface{}, 1)
		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()

			done <- invoker(ctx, method, req, reply, cc, opts...)
		}()

		select {
		case p := <-panicChan:
			panic(p)
		case err := <-done:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
