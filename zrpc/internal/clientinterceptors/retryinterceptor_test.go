package clientinterceptors

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"testing"
	"time"
)

func TestRetryInterceptor_WithMax(t *testing.T) {
	n := 4
	for i := 0; i < n; i++ {
		count := 0
		cc := new(grpc.ClientConn)
		err := AutoRetryInterceptor(true)(context.Background(), "/1", nil, nil, cc,
			func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
				count++
				return status.Error(codes.ResourceExhausted, "ResourceExhausted")
			}, RetryWithMax(i))
		assert.Error(t, err)
		assert.Equal(t, i+1, count)
	}

}

func TestRetryWithDisable(t *testing.T) {
	opt := &retryOptions{}
	assert.EqualValues(t, &retryOptions{}, parseRetryCallOptions(opt, RetryWithDisable()))
}

func TestRetryWithMax(t *testing.T) {
	n := 5
	for i := 0; i < n; i++ {
		opt := &retryOptions{}
		assert.EqualValues(t, &retryOptions{max: i}, parseRetryCallOptions(opt, RetryWithMax(i)))
	}
}

func TestRetryWithBackoff(t *testing.T) {
	opt := &retryOptions{}

	retryCallOptions := parseRetryCallOptions(opt, RetryWithBackoff(func(attempt int) time.Duration {
		return time.Millisecond
	}))
	assert.EqualValues(t, time.Millisecond, retryCallOptions.backoffFunc(context.Background(), 1))

}

func TestRetryWithCodes(t *testing.T) {
	opt := &retryOptions{}
	c := []codes.Code{codes.Unknown, codes.NotFound}
	options := parseRetryCallOptions(opt, RetryWithCodes(c...))
	assert.EqualValues(t, c, options.codes)
}

func TestRetryWithPerRetryTimeout(t *testing.T) {
	opt := &retryOptions{}
	options := parseRetryCallOptions(opt, RetryWithPerRetryTimeout(time.Millisecond))
	assert.EqualValues(t, time.Millisecond, options.perCallTimeout)
}

func Test_waitRetryBackoff(t *testing.T) {

	opt := &retryOptions{perCallTimeout: time.Second, backoffFunc: func(ctx context.Context, attempt int) time.Duration {
		return time.Second
	}}
	logger := logx.WithContext(context.Background())
	err := waitRetryBackoff(logger, 1, context.Background(), opt)
	assert.NoError(t, err)
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancelFunc()
	err = waitRetryBackoff(logger, 1, ctx, opt)
	assert.ErrorIs(t, err, status.FromContextError(context.DeadlineExceeded).Err())
}

func Test_isRetriable(t *testing.T) {
	assert.False(t, isRetriable(status.FromContextError(context.DeadlineExceeded).Err(), &retryOptions{codes: DefaultRetriableCodes}))
	assert.True(t, isRetriable(status.Error(codes.ResourceExhausted, ""), &retryOptions{codes: DefaultRetriableCodes}))
	assert.False(t, isRetriable(errors.New("error"), &retryOptions{}))
}

func Test_perCallContext(t *testing.T) {
	opt := &retryOptions{perCallTimeout: time.Second, includeRetryHeader: true}
	ctx := metadata.NewIncomingContext(context.Background(), map[string][]string{"1": {"1"}})
	callContext := perCallContext(ctx, opt, 1)
	md, ok := metadata.FromOutgoingContext(callContext)
	assert.True(t, ok)
	assert.EqualValues(t, metadata.MD{"1": {"1"}, AttemptMetadataKey: {"1"}}, md)

}

func Test_filterCallOptions(t *testing.T) {
	grpcEmptyCallOpt := &grpc.EmptyCallOption{}
	retryCallOpt := &RetryCallOption{}
	options, retryCallOptions := filterCallOptions([]grpc.CallOption{
		grpcEmptyCallOpt,
		retryCallOpt,
	})
	assert.EqualValues(t, []grpc.CallOption{grpcEmptyCallOpt}, options)
	assert.EqualValues(t, []*RetryCallOption{retryCallOpt}, retryCallOptions)

}
