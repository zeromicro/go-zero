package roundrobin

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

const Name = "zero_rr"

func init() {
	balancer.Register(newRoundRobinBuilder())
}

type roundRobinPickerBuilder struct {
}

func newRoundRobinBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, new(roundRobinPickerBuilder))
}

func (b *roundRobinPickerBuilder) Build(readySCs map[resolver.Address]balancer.SubConn) balancer.Picker {
	rand.Seed(time.Now().UnixNano())
	picker := &roundRobinPicker{
		index: rand.Int(),
	}

	for addr, conn := range readySCs {
		picker.conns = append(picker.conns, &subConn{
			addr: addr,
			conn: conn,
		})
	}

	return picker
}

type roundRobinPicker struct {
	conns []*subConn
	index int
	lock  sync.Mutex
}

func (p *roundRobinPicker) Pick(ctx context.Context, info balancer.PickInfo) (
	conn balancer.SubConn, done func(balancer.DoneInfo), err error) {
	fmt.Println(p.conns)
	p.lock.Lock()
	defer p.lock.Unlock()

	p.index = (p.index + 1) % len(p.conns)
	return p.conns[p.index].conn, func(info balancer.DoneInfo) {
	}, nil
}

type subConn struct {
	addr resolver.Address
	conn balancer.SubConn
}
