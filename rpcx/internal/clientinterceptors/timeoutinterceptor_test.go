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
				func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
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
		func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
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
