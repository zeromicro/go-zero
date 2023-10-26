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
	// MethodTimeoutConf defines specified timeout for gRPC method.
	MethodTimeoutConf struct {
		FullMethod string
		Timeout    time.Duration
	}

	methodTimeouts map[string]time.Duration
)

// UnaryTimeoutInterceptor returns a func that sets timeout to incoming unary requests.
func UnaryTimeoutInterceptor(timeout time.Duration,
	methodTimeouts ...MethodTimeoutConf) grpc.UnaryServerInterceptor {
	timeouts := buildMethodTimeouts(methodTimeouts)
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (any, error) {
		t := getTimeoutByUnaryServerInfo(info.FullMethod, timeouts, timeout)
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

func buildMethodTimeouts(timeouts []MethodTimeoutConf) methodTimeouts {
	mt := make(methodTimeouts, len(timeouts))
	for _, st := range timeouts {
		if st.FullMethod != "" {
			mt[st.FullMethod] = st.Timeout
		}
	}

	return mt
}

func getTimeoutByUnaryServerInfo(method string, timeouts methodTimeouts,
	defaultTimeout time.Duration) time.Duration {
	if v, ok := timeouts[method]; ok {
		return v
	}

	return defaultTimeout
}
