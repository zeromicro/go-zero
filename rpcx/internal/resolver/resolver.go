package resolver

import (
	"zero/core/discov"

	"google.golang.org/grpc/resolver"
)

type discovResolver struct {
	scheme string
	etcd   discov.EtcdConf
	cc     resolver.ClientConn
}

func (r *discovResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (
	resolver.Resolver, error) {
	sub, err := discov.NewSubscriber(r.etcd.Hosts, r.etcd.Key)
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
		r.cc.UpdateState(resolver.State{
			Addresses: addrs,
		})
	})

	return r, nil
}

func (r *discovResolver) Close() {
}

func (r *discovResolver) ResolveNow(options resolver.ResolveNowOptions) {
}

func (r *discovResolver) Scheme() string {
	return r.scheme
}

func RegisterResolver(scheme string, etcd discov.EtcdConf) {
	resolver.Register(&discovResolver{
		scheme: scheme,
		etcd:   etcd,
	})
}
