package resolver

import "google.golang.org/grpc/resolver"

type discovResolver struct {
}

func (r discovResolver) ResolveNow(options resolver.ResolveNowOptions) {
}

func (r discovResolver) Close() {
}
