package internal

import (
	"strings"
	"unicode"

	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/selector"
	"github.com/zeromicro/go-zero/zrpc/resolver/internal/targets"
	"google.golang.org/grpc/attributes"
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
		var addrs []resolver.Address
		for _, val := range subset(sub.Values(), subsetSize) {
			addr := val
			var attrs *attributes.Attributes

			valSplit := strings.SplitN(val, "@", 2)
			if len(valSplit) == 2 {
				addr = valSplit[0]
				colors := strings.FieldsFunc(valSplit[1], func(r rune) bool {
					return r == ',' || unicode.IsSpace(r)
				})
				attrs = attrs.WithValue("colors", selector.NewColors(colors...))
			}

			addrs = append(addrs, resolver.Address{
				Addr:               addr,
				BalancerAttributes: attrs,
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

	return &nopResolver{cc: cc}, nil
}

func (b *discovBuilder) Scheme() string {
	return DiscovScheme
}
