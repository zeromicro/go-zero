package resolver

import (
	"strings"

	"github.com/tal-tech/go-zero/core/discov"
	"google.golang.org/grpc/resolver"
)

type discovBuilder struct{}

func (d *discovBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (
	resolver.Resolver, error) {
	hosts := strings.FieldsFunc(target.Authority, func(r rune) bool {
		return r == EndpointSepChar
	})
	sub, err := discov.NewSubscriber(hosts, target.Endpoint)
	if err != nil {
		return nil, err
	}

	update := func() {
		var addrs []resolver.Address
		for _, val := range subset(sub.Values(), subsetSize) {
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

	return &nopResolver{cc: cc}, nil
}

func (d *discovBuilder) Scheme() string {
	return DiscovScheme
}
