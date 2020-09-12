package load

import (
	"sync/atomic"
	"time"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stat"
)

type (
	SheddingStat struct {
		name  string
		total int64
		pass  int64
		drop  int64
	}

	snapshot struct {
		Total int64
		Pass  int64
		Drop  int64
	}
)

func NewSheddingStat(name string) *SheddingStat {
	st := &SheddingStat{
		name: name,
	}
	go st.run()
	return st
}

func (s *SheddingStat) IncrementTotal() {
	atomic.AddInt64(&s.total, 1)
}

func (s *SheddingStat) IncrementPass() {
	atomic.AddInt64(&s.pass, 1)
}

func (s *SheddingStat) IncrementDrop() {
	atomic.AddInt64(&s.drop, 1)
}

func (s *SheddingStat) reset() snapshot {
	return snapshot{
		Total: atomic.SwapInt64(&s.total, 0),
		Pass:  atomic.SwapInt64(&s.pass, 0),
		Drop:  atomic.SwapInt64(&s.drop, 0),
	}
}

func (s *SheddingStat) run() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		c := stat.CpuUsage()
		st := s.reset()
		if st.Drop == 0 {
			logx.Statf("(%s) shedding_stat [1m], cpu: %d, total: %d, pass: %d, drop: %d",
				s.name, c, st.Total, st.Pass, st.Drop)
		} else {
			logx.Statf("(%s) shedding_stat_drop [1m], cpu: %d, total: %d, pass: %d, drop: %d",
				s.name, c, st.Total, st.Pass, st.Drop)
		}
	}
}
