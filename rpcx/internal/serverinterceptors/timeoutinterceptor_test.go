package serverinterceptors

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
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
