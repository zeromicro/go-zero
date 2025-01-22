package internal

import (
	"github.com/zeromicro/go-zero/core/discov"
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

func TestNopResolver_Close(t *testing.T) {
	var isChanged bool
	r := nopResolver{}
	r.Close()
	assert.False(t, isChanged)
	r = nopResolver{
		closeFunc: func() {
			isChanged = true
		},
	}
	r.Close()
	assert.True(t, isChanged)
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
