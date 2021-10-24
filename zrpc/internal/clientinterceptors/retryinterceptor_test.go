package clientinterceptors

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
	"time"
)

func TestRetryInterceptor(t *testing.T) {

	cc := new(grpc.ClientConn)
	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := RetryInterceptor()(context.Background(), "/1", nil, nil, cc,
		func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
			var err error
			go func() {
				time.Sleep(time.Second)
				err = status.Error(codes.ResourceExhausted, "ResourceExhausted")
			}()
			select {
			case <-ctx.Done():
				err = ctx.Err()
				fmt.Println("超时", err)
			}

			return err
		}, WithMax(200), WithPerRetryTimeout(time.Microsecond))
	assert.Error(t, err)

}
