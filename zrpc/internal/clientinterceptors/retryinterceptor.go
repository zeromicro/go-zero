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
		max:                0, // disabled
		perCallTimeout:     0, // disabled
		includeRetryHeader: true,
		codes:              DefaultRetriableCodes,
		backoffFunc:        retrybackoff.BackoffLinearWithJitter(50*time.Millisecond /*jitter*/, 0.10),
	}
)

type (
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

// RetryDisable disables the retry behaviour on this call, or this interceptor.
//
// Its semantically the same to `RetryWithMax`
func RetryDisable() *RetryCallOption {
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

func RetryWithCodes(retryCodes ...codes.Code) *RetryCallOption {
	return &RetryCallOption{applyFunc: func(o *retryOptions) {
		o.codes = retryCodes
	}}
}

func RetryWithPerRetryTimeout(timeout time.Duration) *RetryCallOption {
	return &RetryCallOption{applyFunc: func(o *retryOptions) {
		o.perCallTimeout = timeout
	}}
}

func RetryInterceptor(opts ...*RetryCallOption) grpc.UnaryClientInterceptor {
	intOpts := reuseOrNewWithCallOptions(defaultOptions, opts)
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		logger := logx.WithContext(ctx)
		grpcOpts, retryOpts := filterCallOptions(opts)
		callOpts := reuseOrNewWithCallOptions(intOpts, retryOpts)
		// short circuit for simplicity, and avoiding allocations.
		if callOpts.max == 0 {
			return invoker(ctx, method, req, reply, cc, grpcOpts...)
		}

		var lastErr error
		for attempt := 0; attempt <= callOpts.max; attempt++ {
			if err := waitRetryBackoff(attempt, ctx, callOpts); err != nil {
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
func waitRetryBackoff(attempt int, ctx context.Context, retryOptions *retryOptions) error {
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
