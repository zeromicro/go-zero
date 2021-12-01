package internal

import (
	"strings"

	"google.golang.org/grpc/resolver"
)

type directBuilder struct{}

func (d *directBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (
	resolver.Resolver, error) {
	var addrs []resolver.Address
	endpoints := strings.FieldsFunc(target.Endpoint, func(r rune) bool {
		return r == EndpointSepChar
	})

	for _, val := range subset(endpoints, subsetSize) {
		addrs = append(addrs, resolver.Address{
			Addr: val,
		})
	}
	if err := cc.UpdateState(resolver.State{
		Addresses: addrs,
	}); err != nil {
		return nil, err
	}

	return &nopResolver{cc: cc}, nil
}

func (d *directBuilder) Scheme() string {
	return DirectScheme
}
