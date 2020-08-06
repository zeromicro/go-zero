package resolver

import (
	"fmt"
	"strings"

	"zero/core/discov"

	"google.golang.org/grpc/resolver"
)

const (
	DiscovScheme = "discov"
	EndpointSep  = ","
)

var builder discovBuilder

type discovBuilder struct{}

func (b *discovBuilder) Scheme() string {
	return DiscovScheme
}

func (b *discovBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (
	resolver.Resolver, error) {
	if target.Scheme != DiscovScheme {
		return nil, fmt.Errorf("bad scheme: %s", target.Scheme)
	}

	hosts := strings.Split(target.Authority, EndpointSep)
	sub, err := discov.NewSubscriber(hosts, target.Endpoint)
	if err != nil {
		return nil, err
	}

	update := func() {
		var addrs []resolver.Address
		for _, val := range sub.Values() {
			addrs = append(addrs, resolver.Address{
				Addr: val,
			})
		}
		cc.UpdateState(resolver.State{
			Addresses: addrs,
		})
	}
	sub.AddListener(update)
	update()

	return &discovResolver{
		cc: cc,
	}, nil
}

type discovResolver struct {
	cc resolver.ClientConn
}

func (r *discovResolver) Close() {
}

func (r *discovResolver) ResolveNow(options resolver.ResolveNowOptions) {
}

func RegisterResolver() {
	resolver.Register(&builder)
}
