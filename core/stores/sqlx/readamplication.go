package sqlx

import (
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	concurrencyThreshold = 3
	logInterval          = 60 * 1000 // 1 minute
)

var logger = logx.NewLessLogger(logInterval)

type (
	concurrentReads struct {
		reads map[string]*queryReference
		lock  sync.Mutex
	}

	queryReference struct {
		concurrency    uint32
		maxConcurrency uint32
	}
)

func newConcurrentReads() *concurrentReads {
	return &concurrentReads{
		reads: make(map[string]*queryReference),
	}
}

func (r *concurrentReads) add(query string) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if ref, ok := r.reads[query]; ok {
		ref.concurrency++
		if ref.maxConcurrency < ref.concurrency {
			ref.maxConcurrency = ref.concurrency
		}
	} else {
		r.reads[query] = &queryReference{
			concurrency:    1,
			maxConcurrency: 1,
		}
	}
}

func (r *concurrentReads) remove(query string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	ref, ok := r.reads[query]
	if !ok {
		return
	}

	if ref.concurrency > 1 {
		ref.concurrency--
		return
	}

	// last reference to remove
	delete(r.reads, query)
	if ref.maxConcurrency >= concurrencyThreshold {
		logger.Errorf("sql query amplified, query: %q, maxConcurrency: %d",
			query, ref.maxConcurrency)
	}
}
