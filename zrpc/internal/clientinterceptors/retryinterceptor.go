package clientinterceptors

import (
	"context"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/zrpc/internal/clientinterceptors/retrybackoff"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"strconv"
	"time"
)

const AttemptMetadataKey = "X-retry-attempt"

var (
	DefaultRetriableCodes = []codes.Code{codes.ResourceExhausted, codes.Unavailable}

	defaultOptions = &retryOptions{
		max:            0, // disabled
		perCallTimeout: 0, // disabled
		includeHeader:  true,
		codes:          DefaultRetriableCodes,
		backoffFunc:    retrybackoff.BackoffLinearWithJitter(50*time.Millisecond /*jitter*/, 0.10),
	}
)

type (
	retryOptions struct {
		max            uint
		perCallTimeout time.Duration
		includeHeader  bool
		codes          []codes.Code
		backoffFunc    retrybackoff.BackoffFunc
	}
	// RetryCallOption is a grpc.CallOption that is local to grpc retry.
	RetryCallOption struct {
		grpc.EmptyCallOption // make sure we implement private after() and before() fields so we don't panic.
		applyFunc            func(opt *retryOptions)
	}
)

// Disable disables the retry behaviour on this call, or this interceptor.
//
// Its semantically the same to `WithMax`
func Disable() *RetryCallOption {
	return WithMax(0)
}

// WithMax sets the maximum number of retries on this call, or this interceptor.
func WithMax(maxRetries uint) *RetryCallOption {
	return &RetryCallOption{applyFunc: func(options *retryOptions) {
		options.max = maxRetries
	}}
}

// WithBackoff sets the `BackoffFunc` used to control time between retries.
func WithBackoff(backoffFunc func(attempt uint) time.Duration) *RetryCallOption {
	return &RetryCallOption{applyFunc: func(o *retryOptions) {
		o.backoffFunc = func(ctx context.Context, attempt uint) time.Duration {
			return backoffFunc(attempt)
		}
	}}
}
func WithCodes(retryCodes ...codes.Code) *RetryCallOption {
	return &RetryCallOption{applyFunc: func(o *retryOptions) {
		o.codes = retryCodes
	}}
}
func WithPerRetryTimeout(timeout time.Duration) *RetryCallOption {
	return &RetryCallOption{applyFunc: func(o *retryOptions) {
		o.perCallTimeout = timeout
	}}
}

func RetryInterceptor(optFuncs ...*RetryCallOption) grpc.UnaryClientInterceptor {
	intOpts := reuseOrNewWithCallOptions(defaultOptions, optFuncs)
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		logger := logx.WithContext(ctx)
		grpcOpts, retryOpts := filterCallOptions(opts)
		callOpts := reuseOrNewWithCallOptions(intOpts, retryOpts)
		// short circuit for simplicity, and avoiding allocations.
		if callOpts.max == 0 {
			return invoker(ctx, method, req, reply, cc, grpcOpts...)
		}
		var lastErr error
		for attempt := uint(1); attempt <= callOpts.max; attempt++ {
			// wait a while
			if err := waitRetryBackoff(attempt, ctx, callOpts); err != nil {
				return err
			}
			// retry
			callCtx := perCallContext(ctx, callOpts, attempt)
			lastErr = invoker(callCtx, method, req, reply, cc, grpcOpts...)

			if lastErr == nil {
				return nil
			}
			logger.Errorf("grpc retry attempt: %d, got err: %v", attempt, lastErr)
			if isContextError(lastErr) {
				if ctx.Err() != nil {
					logger.Errorf("grpc retry attempt: %d, parent context error: %v", attempt, ctx.Err())
					return lastErr
				} else if callOpts.perCallTimeout != 0 {
					logger.Errorf("grpc retry attempt: %d, context error from retry call", attempt)
					continue
				}
			}
			if !isRetriable(lastErr, callOpts) {
				return lastErr
			}
		}
		return lastErr
	}
}
func waitRetryBackoff(attempt uint, ctx context.Context, retryOptions *retryOptions) error {
	logger := logx.WithContext(ctx)
	var waitTime time.Duration = 0
	if attempt > 0 {
		waitTime = retryOptions.backoffFunc(ctx, attempt)
	}
	if waitTime > 0 {
		timer := time.NewTimer(waitTime)
		logger.Infof("grpc retry attempt: %d, backoff for %v", attempt, waitTime)
		select {
		case <-ctx.Done():
			timer.Stop()
			return status.FromContextError(ctx.Err()).Err()
		case <-timer.C:
		}
	}
	return nil
}

func isRetriable(err error, retryOptions *retryOptions) bool {
	errCode := status.Code(err)
	if isContextError(err) {
		return false
	}
	for _, code := range retryOptions.codes {
		if code == errCode {
			return true
		}
	}
	return false
}

func isContextError(err error) bool {
	code := status.Code(err)
	return code == codes.DeadlineExceeded || code == codes.Canceled
}
func reuseOrNewWithCallOptions(opt *retryOptions, callOptions []*RetryCallOption) *retryOptions {
	if len(callOptions) == 0 {
		return opt
	}
	optCopy := &retryOptions{}
	*optCopy = *opt
	for _, f := range callOptions {
		f.applyFunc(optCopy)
	}
	return optCopy
}

func perCallContext(ctx context.Context, callOpts *retryOptions, attempt uint) context.Context {
	if callOpts.perCallTimeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, callOpts.perCallTimeout)
		_ = cancel
	}
	if attempt > 0 && callOpts.includeHeader {
		mdClone := extractIncomingAndClone(ctx)
		mdClone.Set(AttemptMetadataKey, strconv.Itoa(int(attempt)))
		ctx = metadata.NewOutgoingContext(ctx, mdClone)
	}
	return ctx
}
func extractIncomingAndClone(ctx context.Context) metadata.MD {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return metadata.Pairs()
	}
	newMd := metadata.Pairs()
	for k, v := range md {
		newMd[k] = make([]string, len(v))
		copy(newMd[k], v)
	}
	return newMd
}
func filterCallOptions(callOptions []grpc.CallOption) (grpcOptions []grpc.CallOption, retryOptions []*RetryCallOption) {
	for _, opt := range callOptions {
		if co, ok := opt.(*RetryCallOption); ok {
			retryOptions = append(retryOptions, co)
		} else {
			grpcOptions = append(grpcOptions, opt)
		}
	}
	return grpcOptions, retryOptions
}
