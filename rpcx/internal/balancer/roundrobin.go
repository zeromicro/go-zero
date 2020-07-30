package balancer

import (
	"context"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

const (
	Name = "roundrobin"
)

func init() {
	balancer.Register(newBuilder())
}

type roundRobinPickerBuilder struct {
}

func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, new(roundRobinPickerBuilder))
}

func (b *roundRobinPickerBuilder) Build(readySCs map[resolver.Address]balancer.SubConn) balancer.Picker {
	panic("implement me")
}

type roundRobinPicker struct {
}

func (p *roundRobinPicker) Pick(ctx context.Context, info balancer.PickInfo) (
	conn balancer.SubConn, done func(balancer.DoneInfo), err error) {
	panic("implement me")
}
