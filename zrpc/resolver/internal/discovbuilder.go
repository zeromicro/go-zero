package internal

import (
	"strings"

	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc/resolver/internal/targets"
	"google.golang.org/grpc/resolver"
)

type discovBuilder struct{}

func (b *discovBuilder) Build(target resolver.Target, cc resolver.ClientConn, _ resolver.BuildOptions) (
	resolver.Resolver, error) {
	hosts := strings.FieldsFunc(targets.GetAuthority(target), func(r rune) bool {
		return r == EndpointSepChar
	})
	sub, err := discov.NewSubscriber(hosts, targets.GetEndpoints(target))
	if err != nil {
		return nil, err
	}

	update := func() {
		vals := subset(sub.Values(), subsetSize)
		addrs := make([]resolver.Address, 0, len(vals))
		for _, val := range vals {
			addrs = append(addrs, resolver.Address{
				Addr: val,
			})
		}
		if err := cc.UpdateState(resolver.State{
			Addresses: addrs,
		}); err != nil {
			logx.Error(err)
		}
	}
	sub.AddListener(update)
	update()

	return &discovResolver{
		cc:  cc,
		sub: sub,
	}, nil
}

func (b *discovBuilder) Scheme() string {
	return DiscovScheme
}

type discovResolver struct {
	cc  resolver.ClientConn
	sub *discov.Subscriber
}

func (r *discovResolver) Close() {
	r.sub.Close()
}

func (r *discovResolver) ResolveNow(_ resolver.ResolveNowOptions) {
}
