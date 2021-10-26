package retry

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

func TestRetryWithDisable(t *testing.T) {
	opt := &options{}
	assert.EqualValues(t, &options{}, parseRetryCallOptions(opt, WithDisable()))
}

func TestRetryWithMax(t *testing.T) {
	n := 5
	for i := 0; i < n; i++ {
		opt := &options{}
		assert.EqualValues(t, &options{max: i}, parseRetryCallOptions(opt, WithMax(i)))
	}
}

func TestRetryWithBackoff(t *testing.T) {
	opt := &options{}

	retryCallOptions := parseRetryCallOptions(opt, WithBackoff(func(attempt int) time.Duration {
		return time.Millisecond
	}))
	assert.EqualValues(t, time.Millisecond, retryCallOptions.backoffFunc(1))

}

func TestRetryWithCodes(t *testing.T) {
	opt := &options{}
	c := []codes.Code{codes.Unknown, codes.NotFound}
	options := parseRetryCallOptions(opt, WithCodes(c...))
	assert.EqualValues(t, c, options.codes)
}

func TestRetryWithPerRetryTimeout(t *testing.T) {
	opt := &options{}
	options := parseRetryCallOptions(opt, WithPerRetryTimeout(time.Millisecond))
	assert.EqualValues(t, time.Millisecond, options.perCallTimeout)
}

func Test_waitRetryBackoff(t *testing.T) {

	opt := &options{perCallTimeout: time.Second, backoffFunc: func(attempt int) time.Duration {
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
	assert.False(t, isRetriable(status.FromContextError(context.DeadlineExceeded).Err(), &options{codes: DefaultRetriableCodes}))
	assert.True(t, isRetriable(status.Error(codes.ResourceExhausted, ""), &options{codes: DefaultRetriableCodes}))
	assert.False(t, isRetriable(errors.New("error"), &options{}))
}

func Test_perCallContext(t *testing.T) {
	opt := &options{perCallTimeout: time.Second, includeRetryHeader: true}
	ctx := metadata.NewIncomingContext(context.Background(), map[string][]string{"1": {"1"}})
	callContext := perCallContext(ctx, opt, 1)
	md, ok := metadata.FromOutgoingContext(callContext)
	assert.True(t, ok)
	assert.EqualValues(t, metadata.MD{"1": {"1"}, AttemptMetadataKey: {"1"}}, md)

}

func Test_filterCallOptions(t *testing.T) {
	grpcEmptyCallOpt := &grpc.EmptyCallOption{}
	retryCallOpt := &CallOption{}
	options, retryCallOptions := filterCallOptions([]grpc.CallOption{
		grpcEmptyCallOpt,
		retryCallOpt,
	})
	assert.EqualValues(t, []grpc.CallOption{grpcEmptyCallOpt}, options)
	assert.EqualValues(t, []*CallOption{retryCallOpt}, retryCallOptions)

}
