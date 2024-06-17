package clientinterceptors

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestTimeoutInterceptor(t *testing.T) {
	timeouts := []time.Duration{0, time.Millisecond * 10}
	for _, timeout := range timeouts {
		t.Run(strconv.FormatInt(int64(timeout), 10), func(t *testing.T) {
			interceptor := TimeoutInterceptor(timeout)
			cc := new(grpc.ClientConn)
			err := interceptor(context.Background(), "/foo", nil, nil, cc,
				func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
					opts ...grpc.CallOption) error {
					return nil
				},
			)
			assert.Nil(t, err)
		})
	}
}

func TestTimeoutInterceptor_timeout(t *testing.T) {
	const timeout = time.Millisecond * 10
	interceptor := TimeoutInterceptor(timeout)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	var wg sync.WaitGroup
	wg.Add(1)
	cc := new(grpc.ClientConn)
	err := interceptor(ctx, "/foo", nil, nil, cc,
		func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
			opts ...grpc.CallOption) error {
			defer wg.Done()
			tm, ok := ctx.Deadline()
			assert.True(t, ok)
			assert.True(t, tm.Before(time.Now().Add(timeout+time.Millisecond)))
			return nil
		})
	wg.Wait()
	assert.Nil(t, err)
}

func TestTimeoutInterceptor_panic(t *testing.T) {
	timeouts := []time.Duration{0, time.Millisecond * 10}
	for _, timeout := range timeouts {
		t.Run(strconv.FormatInt(int64(timeout), 10), func(t *testing.T) {
			interceptor := TimeoutInterceptor(timeout)
			cc := new(grpc.ClientConn)
			assert.Panics(t, func() {
				_ = interceptor(context.Background(), "/foo", nil, nil, cc,
					func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
						opts ...grpc.CallOption) error {
						panic("any")
					},
				)
			})
		})
	}
}

func TestTimeoutInterceptor_TimeoutCallOption(t *testing.T) {
	type args struct {
		interceptorTimeout time.Duration
		callOptionTimeout  time.Duration
		runTime            time.Duration
	}
	var tests = []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "do not timeout without call option timeout",
			args: args{
				interceptorTimeout: time.Second,
				runTime:            time.Millisecond * 50,
			},
			wantErr: nil,
		},
		{
			name: "timeout without call option timeout",
			args: args{
				interceptorTimeout: time.Second,
				runTime:            time.Second * 2,
			},
			wantErr: context.DeadlineExceeded,
		},
		{
			name: "do not timeout with call option timeout",
			args: args{
				interceptorTimeout: time.Second,
				callOptionTimeout:  time.Second * 3,
				runTime:            time.Second * 2,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			interceptor := TimeoutInterceptor(tt.args.interceptorTimeout)

			cc := new(grpc.ClientConn)
			var co []grpc.CallOption
			if tt.args.callOptionTimeout > 0 {
				co = append(co, WithCallTimeout(tt.args.callOptionTimeout))
			}

			err := interceptor(context.Background(), "/foo", nil, nil, cc,
				func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
					opts ...grpc.CallOption) error {
					timer := time.NewTimer(tt.args.runTime)
					defer timer.Stop()

					select {
					case <-timer.C:
						return nil
					case <-ctx.Done():
						return ctx.Err()
					}
				}, co...,
			)
			t.Logf("error: %+v", err)

			assert.EqualValues(t, tt.wantErr, err)
		})
	}
}
