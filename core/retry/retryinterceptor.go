package retry

import (
	"context"
	"strconv"
	"time"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/retry/backoff"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const AttemptMetadataKey = "x-retry-attempt"

var (
	// DefaultRetriableCodes default retry code
	DefaultRetriableCodes = []codes.Code{codes.ResourceExhausted, codes.Unavailable}
	// defaultRetryOptions default retry configuration
	defaultRetryOptions = &options{
		max:                0, // disabled
		perCallTimeout:     0, // disabled
		includeRetryHeader: true,
		codes:              DefaultRetriableCodes,
		backoffFunc:        backoff.LinearWithJitter(50*time.Millisecond /*jitter*/, 0.10),
	}
)

type (
	// options retry the configuration
	options struct {
		max                int
		perCallTimeout     time.Duration
		includeRetryHeader bool
		codes              []codes.Code
		backoffFunc        backoff.Func
	}

	// CallOption is a grpc.CallOption that is local to grpc retry.
	CallOption struct {
		grpc.EmptyCallOption // make sure we implement private after() and before() fields so we don't panic.
		apply                func(opt *options)
	}
)

func waitRetryBackoff(logger logx.Logger, attempt int, ctx context.Context, retryOptions *options) error {
	var waitTime time.Duration = 0
	if attempt > 0 {
		waitTime = retryOptions.backoffFunc(attempt)
	}
	if waitTime > 0 {
		timer := time.NewTimer(waitTime)
		logger.Infof("grpc retry attempt: %d, backoff for %v", attempt, waitTime)
		select {
		case <-ctx.Done():
			timer.Stop()
			return status.FromContextError(ctx.Err()).Err()
		case <-timer.C:
			// double check
			err := ctx.Err()
			if err != nil {
				return status.FromContextError(err).Err()
			}
		}
	}

	return nil
}

func isRetriable(err error, retryOptions *options) bool {
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

func reuseOrNewWithCallOptions(opt *options, retryCallOptions []*CallOption) *options {
	if len(retryCallOptions) == 0 {
		return opt
	}

	return parseRetryCallOptions(opt, retryCallOptions...)
}

func parseRetryCallOptions(opt *options, opts ...*CallOption) *options {
	for _, option := range opts {
		option.apply(opt)
	}

	return opt
}

func perCallContext(ctx context.Context, callOpts *options, attempt int) context.Context {
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
		return metadata.MD{}
	}

	return md.Copy()
}

func filterCallOptions(callOptions []grpc.CallOption) (grpcOptions []grpc.CallOption, retryOptions []*CallOption) {
	for _, opt := range callOptions {
		if co, ok := opt.(*CallOption); ok {
			retryOptions = append(retryOptions, co)
		} else {
			grpcOptions = append(grpcOptions, opt)
		}
	}

	return grpcOptions, retryOptions
}

func Do(ctx context.Context, call func(ctx context.Context, opts ...grpc.CallOption) error, opts ...grpc.CallOption) error {
	logger := logx.WithContext(ctx)
	grpcOpts, retryOpts := filterCallOptions(opts)
	callOpts := reuseOrNewWithCallOptions(defaultRetryOptions, retryOpts)

	if callOpts.max == 0 {
		return call(ctx, opts...)
	}

	var lastErr error
	for attempt := 0; attempt <= callOpts.max; attempt++ {
		if err := waitRetryBackoff(logger, attempt, ctx, callOpts); err != nil {
			return err
		}

		callCtx := perCallContext(ctx, callOpts, attempt)
		lastErr = call(callCtx, grpcOpts...)

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
