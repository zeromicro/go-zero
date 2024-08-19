package serverinterceptors

import (
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/zrpc/internal/codes"
	"google.golang.org/grpc"
	gcodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// StreamBreakerInterceptor is an interceptor that acts as a circuit breaker.
func StreamBreakerInterceptor(svr any, stream grpc.ServerStream, info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	breakerName := info.FullMethod
	err := breaker.DoWithAcceptable(breakerName, func() error {
		return handler(svr, stream)
	}, serverSideAcceptable)

	return convertError(err)
}

// UnaryBreakerInterceptor is an interceptor that acts as a circuit breaker.
func UnaryBreakerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp any, err error) {
	breakerName := info.FullMethod
	err = breaker.DoWithAcceptableCtx(ctx, breakerName, func() error {
		var err error
		resp, err = handler(ctx, req)
		return err
	}, serverSideAcceptable)

	return resp, convertError(err)
}

func convertError(err error) error {
	if err == nil {
		return nil
	}

	// we don't convert context.DeadlineExceeded to status error,
	// because grpc will convert it and return to the client.
	if errors.Is(err, breaker.ErrServiceUnavailable) {
		return status.Error(gcodes.Unavailable, err.Error())
	}

	return err
}

func serverSideAcceptable(err error) bool {
	if errorx.In(err, context.DeadlineExceeded, breaker.ErrServiceUnavailable) {
		return false
	}
	return codes.Acceptable(err)
}
