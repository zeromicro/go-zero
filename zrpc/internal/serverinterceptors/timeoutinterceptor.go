package serverinterceptors

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	// ServerSpecifiedTimeoutConf defines specified timeout for gRPC method.
	ServerSpecifiedTimeoutConf struct {
		FullMethod string
		Timeout    time.Duration
	}

	specifiedTimeoutCache map[string]time.Duration
)

// UnaryTimeoutInterceptor returns a func that sets timeout to incoming unary requests.
func UnaryTimeoutInterceptor(timeout time.Duration, specifiedTimeouts ...ServerSpecifiedTimeoutConf) grpc.UnaryServerInterceptor {
	cache := cacheSpecifiedTimeout(specifiedTimeouts)
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (any, error) {
		t := getTimeoutByUnaryServerInfo(info, timeout, cache)
		ctx, cancel := context.WithTimeout(ctx, t)
		defer cancel()

		var resp any
		var err error
		var lock sync.Mutex
		done := make(chan struct{})
		// create channel with buffer size 1 to avoid goroutine leak
		panicChan := make(chan any, 1)
		go func() {
			defer func() {
				if p := recover(); p != nil {
					// attach call stack to avoid missing in different goroutine
					panicChan <- fmt.Sprintf("%+v\n\n%s", p, strings.TrimSpace(string(debug.Stack())))
				}
			}()

			lock.Lock()
			defer lock.Unlock()
			resp, err = handler(ctx, req)
			close(done)
		}()

		select {
		case p := <-panicChan:
			panic(p)
		case <-done:
			lock.Lock()
			defer lock.Unlock()
			return resp, err
		case <-ctx.Done():
			err := ctx.Err()
			if errors.Is(err, context.Canceled) {
				err = status.Error(codes.Canceled, err.Error())
			} else if errors.Is(err, context.DeadlineExceeded) {
				err = status.Error(codes.DeadlineExceeded, err.Error())
			}
			return nil, err
		}
	}
}

func cacheSpecifiedTimeout(specifiedTimeouts []ServerSpecifiedTimeoutConf) specifiedTimeoutCache {
	cache := make(specifiedTimeoutCache, len(specifiedTimeouts))
	for _, st := range specifiedTimeouts {
		if st.FullMethod != "" {
			cache[st.FullMethod] = st.Timeout
		}
	}

	return cache
}

func getTimeoutByUnaryServerInfo(info *grpc.UnaryServerInfo, defaultTimeout time.Duration, specifiedTimeout specifiedTimeoutCache) time.Duration {
	if ts, ok := info.Server.(TimeoutStrategy); ok {
		return ts.GetTimeoutByFullMethod(info.FullMethod, defaultTimeout)
	} else if v, ok := specifiedTimeout[info.FullMethod]; ok {
		return v
	}

	return defaultTimeout
}

type TimeoutStrategy interface {
	GetTimeoutByFullMethod(fullMethod string, defaultTimeout time.Duration) time.Duration
}
