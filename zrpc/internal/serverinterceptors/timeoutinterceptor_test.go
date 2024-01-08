package serverinterceptors

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	deadlineExceededErr = status.Error(codes.DeadlineExceeded, context.DeadlineExceeded.Error())
	canceledErr         = status.Error(codes.Canceled, context.Canceled.Error())
)

func TestUnaryTimeoutInterceptor(t *testing.T) {
	interceptor := UnaryTimeoutInterceptor(time.Millisecond * 10)
	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/",
	}, func(ctx context.Context, req any) (any, error) {
		return nil, nil
	})
	assert.Nil(t, err)
}

func TestUnaryTimeoutInterceptor_panic(t *testing.T) {
	interceptor := UnaryTimeoutInterceptor(time.Millisecond * 10)
	assert.Panics(t, func() {
		_, _ = interceptor(context.Background(), nil, &grpc.UnaryServerInfo{
			FullMethod: "/",
		}, func(ctx context.Context, req any) (any, error) {
			panic("any")
		})
	})
}

func TestUnaryTimeoutInterceptor_timeout(t *testing.T) {
	const timeout = time.Millisecond * 10
	interceptor := UnaryTimeoutInterceptor(timeout)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)
	_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{
		FullMethod: "/",
	}, func(ctx context.Context, req any) (any, error) {
		defer wg.Done()
		tm, ok := ctx.Deadline()
		assert.True(t, ok)
		assert.True(t, tm.Before(time.Now().Add(timeout+time.Millisecond)))
		return nil, nil
	})
	wg.Wait()
	assert.Nil(t, err)
}

func TestUnaryTimeoutInterceptor_timeoutExpire(t *testing.T) {
	const timeout = time.Millisecond * 10
	interceptor := UnaryTimeoutInterceptor(timeout)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)
	_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{
		FullMethod: "/",
	}, func(ctx context.Context, req any) (any, error) {
		defer wg.Done()
		time.Sleep(time.Millisecond * 50)
		return nil, nil
	})
	wg.Wait()
	assert.EqualValues(t, deadlineExceededErr, err)
}

func TestUnaryTimeoutInterceptor_cancel(t *testing.T) {
	const timeout = time.Minute * 10
	interceptor := UnaryTimeoutInterceptor(timeout)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{
		FullMethod: "/",
	}, func(ctx context.Context, req any) (any, error) {
		defer wg.Done()
		time.Sleep(time.Millisecond * 50)
		return nil, nil
	})

	wg.Wait()
	assert.EqualValues(t, canceledErr, err)
}

type tempServer struct {
	timeout time.Duration
}

func (s *tempServer) run(duration time.Duration) {
	time.Sleep(duration)
}

func TestUnaryTimeoutInterceptor_TimeoutStrategy(t *testing.T) {
	type args struct {
		interceptorTimeout time.Duration
		contextTimeout     time.Duration
		serverTimeout      time.Duration
		runTime            time.Duration

		fullMethod string
	}
	var tests = []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "do not timeout with interceptor timeout",
			args: args{
				interceptorTimeout: time.Second,
				contextTimeout:     time.Second * 5,
				serverTimeout:      time.Second * 3,
				runTime:            time.Millisecond * 50,
				fullMethod:         "/",
			},
			wantErr: nil,
		},
		{
			name: "timeout with interceptor timeout",
			args: args{
				interceptorTimeout: time.Second,
				contextTimeout:     time.Second * 5,
				serverTimeout:      time.Second * 3,
				runTime:            time.Second * 2,
				fullMethod:         "/",
			},
			wantErr: deadlineExceededErr,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			interceptor := UnaryTimeoutInterceptor(tt.args.interceptorTimeout)
			ctx, cancel := context.WithTimeout(context.Background(), tt.args.contextTimeout)
			defer cancel()

			svr := &tempServer{timeout: tt.args.serverTimeout}

			_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{
				Server:     svr,
				FullMethod: tt.args.fullMethod,
			}, func(ctx context.Context, req interface{}) (interface{}, error) {
				svr.run(tt.args.runTime)
				return nil, nil
			})
			t.Logf("error: %+v", err)

			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}

func TestUnaryTimeoutInterceptor_SpecifiedTimeout(t *testing.T) {
	type args struct {
		interceptorTimeout time.Duration
		contextTimeout     time.Duration
		method             string
		methodTimeout      time.Duration
		runTime            time.Duration
	}
	var tests = []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "do not timeout without set timeout for full method",
			args: args{
				interceptorTimeout: time.Second,
				contextTimeout:     time.Second * 5,
				method:             "/run",
				runTime:            time.Millisecond * 50,
			},
			wantErr: nil,
		},
		{
			name: "do not timeout with set timeout for full method",
			args: args{
				interceptorTimeout: time.Second,
				contextTimeout:     time.Second * 5,
				method:             "/run/do_not_timeout",
				methodTimeout:      time.Second * 3,
				runTime:            time.Second * 2,
			},
			wantErr: nil,
		},
		{
			name: "timeout with set timeout for full method",
			args: args{
				interceptorTimeout: time.Second,
				contextTimeout:     time.Second * 5,
				method:             "/run/timeout",
				methodTimeout:      time.Millisecond * 100,
				runTime:            time.Millisecond * 500,
			},
			wantErr: deadlineExceededErr,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var specifiedTimeouts []MethodTimeoutConf
			if tt.args.methodTimeout > 0 {
				specifiedTimeouts = []MethodTimeoutConf{
					{
						FullMethod: tt.args.method,
						Timeout:    tt.args.methodTimeout,
					},
				}
			}

			interceptor := UnaryTimeoutInterceptor(tt.args.interceptorTimeout, specifiedTimeouts...)
			ctx, cancel := context.WithTimeout(context.Background(), tt.args.contextTimeout)
			defer cancel()

			_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{
				FullMethod: tt.args.method,
			}, func(ctx context.Context, req interface{}) (interface{}, error) {
				time.Sleep(tt.args.runTime)
				return nil, nil
			})
			t.Logf("error: %+v", err)

			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}
