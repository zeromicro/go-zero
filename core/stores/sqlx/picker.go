package sqlx

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/core/timex"
)

const (
	decayTime       = int64(time.Second * 10) // default value from finagle
	forcePick       = int64(time.Second)
	initSuccess     = 1000
	throttleSuccess = initSuccess / 2
	penalty         = int64(math.MaxInt32)
	pickTimes       = 3
	logInterval     = time.Minute
)

var ErrNoFollowerAvailable = errors.New("no follower available")

type (
	picker interface {
		pick() (*pickResult, error)
	}
	pickResult struct {
		conn SqlConn
		done func(err error)
	}
	p2cPicker struct {
		conns []*subConn
		r     *rand.Rand
		stamp *syncx.AtomicDuration
		lock  sync.Mutex
	}
)

func newP2cPicker(followers []SqlConn) *p2cPicker {
	var conns []*subConn
	for i, follower := range followers {
		conns = append(conns, &subConn{
			success: initSuccess,
			db:      i,
			conn:    follower,
		})
	}

	return &p2cPicker{
		conns: conns,
		r:     rand.New(rand.NewSource(time.Now().UnixNano())),
		stamp: syncx.NewAtomicDuration(),
	}
}

func (p *p2cPicker) pick() (*pickResult, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	var chosen *subConn
	switch len(p.conns) {
	case 0:
		return nil, ErrNoFollowerAvailable
	case 1:
		chosen = p.choose(p.conns[0], nil)
	case 2:
		chosen = p.choose(p.conns[0], p.conns[1])
	default:
		var node1, node2 *subConn
		for i := 0; i < pickTimes; i++ {
			a := p.r.Intn(len(p.conns))
			b := p.r.Intn(len(p.conns) - 1)
			if b >= a {
				b++
			}
			node1 = p.conns[a]
			node2 = p.conns[b]
			if node1.healthy() && node2.healthy() {
				break
			}
		}

		chosen = p.choose(node1, node2)
	}

	atomic.AddInt64(&chosen.inflight, 1)
	atomic.AddInt64(&chosen.requests, 1)

	return &pickResult{
		conn: chosen.conn,
		done: p.buildDoneFunc(chosen),
	}, nil
}

func (p *p2cPicker) buildDoneFunc(c *subConn) func(err error) {
	start := int64(timex.Now())
	return func(err error) {
		atomic.AddInt64(&c.inflight, -1)
		now := timex.Now()
		last := atomic.SwapInt64(&c.last, int64(now))
		td := int64(now) - last
		if td < 0 {
			td = 0
		}
		w := math.Exp(float64(-td) / float64(decayTime))
		lag := int64(now) - start
		if lag < 0 {
			lag = 0
		}
		olag := atomic.LoadUint64(&c.lag)
		if olag == 0 {
			w = 0
		}
		atomic.StoreUint64(&c.lag, uint64(float64(olag)*w+float64(lag)*(1-w)))
		success := initSuccess

		if err != nil && !acceptable(err) {
			success = 0
		}

		osucc := atomic.LoadUint64(&c.success)
		atomic.StoreUint64(&c.success, uint64(float64(osucc)*w+float64(success)*(1-w)))

		stamp := p.stamp.Load()
		if now-stamp >= logInterval {
			if p.stamp.CompareAndSwap(stamp, now) {
				p.logStats()
			}
		}
	}
}

func (p *p2cPicker) choose(c1, c2 *subConn) *subConn {
	start := int64(timex.Now())
	if c2 == nil {
		atomic.StoreInt64(&c1.pick, start)
		return c1
	}

	if c1.load() > c2.load() {
		c1, c2 = c2, c1
	}

	pick := atomic.LoadInt64(&c2.pick)
	if start-pick > forcePick && atomic.CompareAndSwapInt64(&c2.pick, pick, start) {
		return c2
	}

	atomic.StoreInt64(&c1.pick, start)
	return c1
}

func (p *p2cPicker) logStats() {
	var stats []string

	p.lock.Lock()
	defer p.lock.Unlock()

	for _, conn := range p.conns {
		stats = append(stats, fmt.Sprintf("db: %d, load: %d, reqs: %d",
			conn.db, conn.load(), atomic.SwapInt64(&conn.requests, 0)))
	}

	logx.Statf("follower-db - p2c - %s", strings.Join(stats, "; "))
}

type subConn struct {
	lag      uint64
	inflight int64
	success  uint64
	requests int64
	last     int64
	pick     int64
	db       int
	conn     SqlConn
}

func (c *subConn) healthy() bool {
	return atomic.LoadUint64(&c.success) > throttleSuccess
}

func (c *subConn) load() int64 {
	// plus one to avoid multiply zero
	lag := int64(math.Sqrt(float64(atomic.LoadUint64(&c.lag) + 1)))
	load := lag * (atomic.LoadInt64(&c.inflight) + 1)
	if load == 0 {
		return penalty
	}

	return load
}

func acceptable(err error) bool {
	if errors.Is(err, ErrNotFound) {
		return true
	}

	return false
}
