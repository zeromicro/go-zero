package internal

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/client/v3/mock/mockserver"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)

func TestDiscovBuilder_Scheme(t *testing.T) {
	var b discovBuilder
	assert.Equal(t, DiscovScheme, b.Scheme())
}

func TestDiscovBuilder_Build(t *testing.T) {
	servers, err := mockserver.StartMockServers(2)
	assert.NoError(t, err)
	t.Cleanup(func() {
		servers.Stop()
	})

	var addrs []string
	for _, server := range servers.Servers {
		addrs = append(addrs, server.Address)
	}
	u, err := url.Parse(fmt.Sprintf("%s://%s", DiscovScheme, strings.Join(addrs, ",")))
	assert.NoError(t, err)

	var b discovBuilder
	_, err = b.Build(resolver.Target{
		URL: *u,
	}, mockClientConn{}, resolver.BuildOptions{})
	assert.Error(t, err)
}

type mockClientConn struct{}

func (m mockClientConn) UpdateState(_ resolver.State) error {
	return nil
}

func (m mockClientConn) ReportError(_ error) {
}

func (m mockClientConn) NewAddress(_ []resolver.Address) {
}

func (m mockClientConn) NewServiceConfig(_ string) {
}

func (m mockClientConn) ParseServiceConfig(_ string) *serviceconfig.ParseResult {
	return nil
}
