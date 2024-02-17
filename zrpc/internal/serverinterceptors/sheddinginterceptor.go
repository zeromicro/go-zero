package serverinterceptors

import (
	"context"
	"errors"
	"sync"

	"github.com/zeromicro/go-zero/core/load"
	"github.com/zeromicro/go-zero/core/stat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const serviceType = "rpc"

var (
	sheddingStat *load.SheddingStat
	lock         sync.Mutex
)

// UnarySheddingInterceptor returns a func that does load shedding on processing unary requests.
func UnarySheddingInterceptor(shedder load.Shedder, metrics *stat.Metrics) grpc.UnaryServerInterceptor {
	ensureSheddingStat()

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (val any, err error) {
		sheddingStat.IncrementTotal()
		var promise load.Promise
		promise, err = shedder.Allow()
		if err != nil {
			metrics.AddDrop()
			sheddingStat.IncrementDrop()
			err = status.Error(codes.ResourceExhausted, err.Error())
			return
		}

		defer func() {
			if errors.Is(err, context.DeadlineExceeded) {
				promise.Fail()
			} else {
				sheddingStat.IncrementPass()
				promise.Pass()
			}
		}()

		return handler(ctx, req)
	}
}

func ensureSheddingStat() {
	lock.Lock()
	if sheddingStat == nil {
		sheddingStat = load.NewSheddingStat(serviceType)
	}
	lock.Unlock()
}
