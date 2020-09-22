package serverinterceptors

import (
	"context"
	"sync"

	"github.com/tal-tech/go-zero/core/load"
	"github.com/tal-tech/go-zero/core/stat"
	"google.golang.org/grpc"
)

const serviceType = "rpc"

var (
	sheddingStat *load.SheddingStat
	lock         sync.Mutex
)

func UnarySheddingInterceptor(shedder load.Shedder, metrics *stat.Metrics) grpc.UnaryServerInterceptor {
	ensureSheddingStat()

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (val interface{}, err error) {
		sheddingStat.IncrementTotal()
		var promise load.Promise
		promise, err = shedder.Allow()
		if err != nil {
			metrics.AddDrop()
			sheddingStat.IncrementDrop()
			return
		}

		defer func() {
			if err == context.DeadlineExceeded {
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
