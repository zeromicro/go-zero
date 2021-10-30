package serverinterceptors

import (
	"context"
	"strconv"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func RetryInterceptor(maxAttempt int) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		var md metadata.MD
		requestMd, ok := metadata.FromIncomingContext(ctx)
		if ok {
			md = requestMd.Copy()
			attemptMd := md.Get(retry.AttemptMetadataKey)
			if len(attemptMd) != 0 && attemptMd[0] != "" {
				if attempt, err := strconv.Atoi(attemptMd[0]); err == nil {
					if attempt > maxAttempt {
						logx.WithContext(ctx).Errorf("retries exceeded:%d, max retries:%d", attempt, maxAttempt)
						return nil, status.Error(codes.FailedPrecondition, "Retries exceeded")
					}
				}
			}
		}

		return handler(ctx, req)
	}
}
