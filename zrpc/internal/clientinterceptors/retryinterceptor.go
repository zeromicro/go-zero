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

const AttemptMetadataKey = "x-retry-attempt"

var (
	// DefaultRetriableCodes default retry code
	DefaultRetriableCodes = []codes.Code{codes.ResourceExhausted, codes.Unavailable}
	// defaultRetryOptions default retry configuration
	defaultRetryOptions = &retryOptions{
		max:                0, // disabled
		perCallTimeout:     0, // disabled
		includeRetryHeader: true,
		codes:              DefaultRetriableCodes,
		backoffFunc:        retrybackoff.BackoffLinearWithJitter(50*time.Millisecond /*jitter*/, 0.10),
	}
)

type (
	// retryOptions retry the configuration
	retryOptions struct {
		max                int
		perCallTimeout     time.Duration
		includeRetryHeader bool
		codes              []codes.Code
		backoffFunc        retrybackoff.BackoffFunc
	}
	// RetryCallOption is a grpc.CallOption that is local to grpc retry.
	RetryCallOption struct {
		grpc.EmptyCallOption // make sure we implement private after() and before() fields so we don't panic.
		applyFunc            func(opt *retryOptions)
	}
)

// RetryWithDisable disables the retry behaviour on this call, or this interceptor.
//
// Its semantically the same to `RetryWithMax`
func RetryWithDisable() *RetryCallOption {
	return RetryWithMax(0)
}

// RetryWithMax sets the maximum number of retries on this call, or this interceptor.
func RetryWithMax(maxRetries int) *RetryCallOption {
	return &RetryCallOption{applyFunc: func(options *retryOptions) {
		options.max = maxRetries
	}}
}

// RetryWithBackoff sets the `BackoffFunc` used to control time between retries.
func RetryWithBackoff(backoffFunc func(attempt int) time.Duration) *RetryCallOption {
	return &RetryCallOption{applyFunc: func(o *retryOptions) {
		o.backoffFunc = func(ctx context.Context, attempt int) time.Duration {
			return backoffFunc(attempt)
		}
	}}
}

// RetryWithCodes Allow code to be retried.
func RetryWithCodes(retryCodes ...codes.Code) *RetryCallOption {
	return &RetryCallOption{applyFunc: func(o *retryOptions) {
		o.codes = retryCodes
	}}
}

// RetryWithPerRetryTimeout timeout for each retry
func RetryWithPerRetryTimeout(timeout time.Duration) *RetryCallOption {
	return &RetryCallOption{applyFunc: func(o *retryOptions) {
		o.perCallTimeout = timeout
	}}
}

// AutoRetryInterceptor retry interceptor
func AutoRetryInterceptor(enable bool) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if !enable {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		logger := logx.WithContext(ctx)
		grpcOpts, retryOpts := filterCallOptions(opts)
		callOpts := reuseOrNewWithCallOptions(defaultRetryOptions, retryOpts)
		// short circuit for simplicity, and avoiding allocations.
		if callOpts.max == 0 {
			return invoker(ctx, method, req, reply, cc, grpcOpts...)
		}

		var lastErr error
		for attempt := 0; attempt <= callOpts.max; attempt++ {
			if err := waitRetryBackoff(logger, attempt, ctx, callOpts); err != nil {
				return err
			}

			callCtx := perCallContext(ctx, callOpts, attempt)
			lastErr = invoker(callCtx, method, req, reply, cc, grpcOpts...)

			if lastErr == nil {
				return nil
			}
			if attempt == 0 {
				logger.Errorf("grpc call failed, got err: %v", lastErr)
			} else {
				logger.Errorf("grpc retry attempt: %d, got err: %v", attempt, lastErr)
			}
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

func waitRetryBackoff(logger logx.Logger, attempt int, ctx context.Context, retryOptions *retryOptions) error {
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

func reuseOrNewWithCallOptions(opt *retryOptions, retryCallOptions []*RetryCallOption) *retryOptions {
	if len(retryCallOptions) == 0 {
		return opt
	}
	return parseRetryCallOptions(opt, retryCallOptions...)
}

func parseRetryCallOptions(opt *retryOptions, opts ...*RetryCallOption) *retryOptions {
	for _, option := range opts {
		option.applyFunc(opt)
	}
	return opt
}

func perCallContext(ctx context.Context, callOpts *retryOptions, attempt int) context.Context {
	if attempt > 0 {
		if callOpts.perCallTimeout != 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, callOpts.perCallTimeout)
			_ = cancel
		}
		if callOpts.includeRetryHeader {
			cloneMd := extractIncomingAndClone(ctx)
			cloneMd.Set(AttemptMetadataKey, strconv.Itoa(attempt))
			ctx = metadata.NewOutgoingContext(ctx, cloneMd)
		}
	}

	return ctx
}

func extractIncomingAndClone(ctx context.Context) metadata.MD {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return metadata.Pairs()
	}
	// clone
	cloneMd := metadata.Pairs()
	for k, v := range md {
		cloneMd[k] = make([]string, len(v))
		copy(cloneMd[k], v)
	}
	return cloneMd
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
