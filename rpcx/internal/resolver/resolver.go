package resolver

import (
	"zero/core/discov"

	"google.golang.org/grpc/resolver"
)

const discovScheme = "discov"

type discovBuilder struct {
	etcd discov.EtcdConf
}

func (b *discovBuilder) Scheme() string {
	return discovScheme
}

func (b *discovBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (
	resolver.Resolver, error) {
	sub, err := discov.NewSubscriber(b.etcd.Hosts, b.etcd.Key)
	if err != nil {
		return nil, err
	}

	sub.AddListener(func() {
		vals := sub.Values()
		var addrs []resolver.Address
		for _, val := range vals {
			addrs = append(addrs, resolver.Address{
				Addr: val,
			})
		}
		cc.UpdateState(resolver.State{
			Addresses: addrs,
		})
	})

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

func RegisterResolver(etcd discov.EtcdConf) {
	resolver.Register(&discovBuilder{
		etcd: etcd,
	})
}
