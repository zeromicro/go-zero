package internal

import (
	"sync/atomic"
	"time"

	"github.com/tal-tech/go-zero/core/logx"
)

const statInterval = time.Minute

type CacheStat struct {
	name string
	// export the fields to let the unit tests working,
	// reside in internal package, doesn't matter.
	Total   uint64
	Hit     uint64
	Miss    uint64
	DbFails uint64
}

func NewCacheStat(name string) *CacheStat {
	ret := &CacheStat{
		name: name,
	}
	go ret.statLoop()

	return ret
}

func (cs *CacheStat) IncrementTotal() {
	atomic.AddUint64(&cs.Total, 1)
}

func (cs *CacheStat) IncrementHit() {
	atomic.AddUint64(&cs.Hit, 1)
}

func (cs *CacheStat) IncrementMiss() {
	atomic.AddUint64(&cs.Miss, 1)
}

func (cs *CacheStat) IncrementDbFails() {
	atomic.AddUint64(&cs.DbFails, 1)
}

func (cs *CacheStat) statLoop() {
	ticker := time.NewTicker(statInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			total := atomic.SwapUint64(&cs.Total, 0)
			if total == 0 {
				continue
			}

			hit := atomic.SwapUint64(&cs.Hit, 0)
			percent := 100 * float32(hit) / float32(total)
			miss := atomic.SwapUint64(&cs.Miss, 0)
			dbf := atomic.SwapUint64(&cs.DbFails, 0)
			logx.Statf("dbcache(%s) - qpm: %d, hit_ratio: %.1f%%, hit: %d, miss: %d, db_fails: %d",
				cs.name, total, percent, hit, miss, dbf)
		}
	}
}
