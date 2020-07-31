package resolver

import "google.golang.org/grpc/resolver"

type discovResolver struct {
	cc resolver.ClientConn
}

func (r discovResolver) ResolveNow(options resolver.ResolveNowOptions) {
}

func (r discovResolver) Close() {
}
