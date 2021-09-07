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

func TestUnaryTimeoutInterceptor(t *testing.T) {
	interceptor := UnaryTimeoutInterceptor(time.Millisecond * 10)
	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{
		FullMethod: "/",
	}, func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	})
	assert.Nil(t, err)
}

func TestUnaryTimeoutInterceptor_panic(t *testing.T) {
	interceptor := UnaryTimeoutInterceptor(time.Millisecond * 10)
	assert.Panics(t, func() {
		_, _ = interceptor(context.Background(), nil, &grpc.UnaryServerInfo{
			FullMethod: "/",
		}, func(ctx context.Context, req interface{}) (interface{}, error) {
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
	}, func(ctx context.Context, req interface{}) (interface{}, error) {
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
	}, func(ctx context.Context, req interface{}) (interface{}, error) {
		defer wg.Done()
		time.Sleep(time.Millisecond * 50)
		return nil, nil
	})
	wg.Wait()
	assert.EqualValues(t, status.Error(codes.DeadlineExceeded, context.DeadlineExceeded.Error()), err)
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
	}, func(ctx context.Context, req interface{}) (interface{}, error) {
		defer wg.Done()
		time.Sleep(time.Millisecond * 50)
		return nil, nil
	})

	wg.Wait()
	assert.EqualValues(t, status.Error(codes.Canceled, context.Canceled.Error()), err)
}
