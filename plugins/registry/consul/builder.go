package consul

import (
	"context"

	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"google.golang.org/grpc/resolver"
)

// schemeName for the urls
// All target URLs like 'consul://.../...' will be resolved by this resolver
const schemeName = "consul"

// builder implements resolver.Builder and use for constructing all consul resolvers
type builder struct{}

func (b *builder) Build(url resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	tgt, err := parseURL(url.URL)
	if err != nil {
		return nil, errors.Wrap(err, "Wrong consul URL")
	}
	cli, err := api.NewClient(tgt.consulConfig())
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't connect to the Consul API")
	}
	cli.Health()
	ctx, cancel := context.WithCancel(context.Background())
	pipe := make(chan []*consulAddr)
	go watchConsulService(ctx, cli.Health(), tgt, pipe)
	go populateEndpoints(ctx, cc, pipe)

	return &resolvr{cancelFunc: cancel}, nil
}

// Scheme returns the scheme supported by this resolver.
// Scheme is defined at https://github.com/grpc/grpc/blob/master/doc/naming.md.
func (b *builder) Scheme() string {
	return schemeName
}
