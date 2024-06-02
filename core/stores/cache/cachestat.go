package cache

import (
	"sync/atomic"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/timex"
)

const statInterval = time.Minute

// A Stat is used to stat the cache.
type Stat struct {
	name string
	// export the fields to let the unit tests working,
	// reside in internal package, doesn't matter.
	Total   uint64
	Hit     uint64
	Miss    uint64
	DbFails uint64
}

// NewStat returns a Stat.
func NewStat(name string) *Stat {
	ret := &Stat{
		name: name,
	}

	go func() {
		ticker := timex.NewTicker(statInterval)
		defer ticker.Stop()

		ret.statLoop(ticker)
	}()

	return ret
}

// IncrementTotal increments the total count.
func (s *Stat) IncrementTotal() {
	atomic.AddUint64(&s.Total, 1)
}

// IncrementHit increments the hit count.
func (s *Stat) IncrementHit() {
	atomic.AddUint64(&s.Hit, 1)
}

// IncrementMiss increments the miss count.
func (s *Stat) IncrementMiss() {
	atomic.AddUint64(&s.Miss, 1)
}

// IncrementDbFails increments the db fail count.
func (s *Stat) IncrementDbFails() {
	atomic.AddUint64(&s.DbFails, 1)
}

func (s *Stat) statLoop(ticker timex.Ticker) {
	for range ticker.Chan() {
		total := atomic.SwapUint64(&s.Total, 0)
		if total == 0 {
			continue
		}

		hit := atomic.SwapUint64(&s.Hit, 0)
		percent := 100 * float32(hit) / float32(total)
		miss := atomic.SwapUint64(&s.Miss, 0)
		dbf := atomic.SwapUint64(&s.DbFails, 0)
		logx.Statf("dbcache(%s) - qpm: %d, hit_ratio: %.1f%%, hit: %d, miss: %d, db_fails: %d",
			s.name, total, percent, hit, miss, dbf)
	}
}
