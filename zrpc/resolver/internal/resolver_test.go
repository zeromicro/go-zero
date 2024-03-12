package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)

func TestNopResolver(t *testing.T) {
	assert.NotPanics(t, func() {
		RegisterResolver()
		// make sure ResolveNow & Close don't panic
		var r nopResolver
		r.ResolveNow(resolver.ResolveNowOptions{})
		r.Close()
	})
}

type mockedClientConn struct {
	state resolver.State
	err   error
}

func (m *mockedClientConn) UpdateState(state resolver.State) error {
	m.state = state
	return m.err
}

func (m *mockedClientConn) ReportError(_ error) {
}

func (m *mockedClientConn) NewAddress(_ []resolver.Address) {
}

func (m *mockedClientConn) NewServiceConfig(_ string) {
}

func (m *mockedClientConn) ParseServiceConfig(_ string) *serviceconfig.ParseResult {
	return nil
}
