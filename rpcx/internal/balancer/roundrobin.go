package balancer

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
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

func (b *roundRobinPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	panic("implement me")
}
