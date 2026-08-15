package internal

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/core/netx"
	"google.golang.org/grpc"
)

func TestNewRpcPubServer(t *testing.T) {
	s, err := NewRpcPubServer(discov.EtcdConf{
		User: "user",
		Pass: "pass",
		ID:   10,
	}, "")
	assert.NoError(t, err)
	assert.NotPanics(t, func() {
		s.Start(nil)
	})
}

func TestFigureOutListenOn(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{
			input:  "192.168.0.5:1234",
			expect: "192.168.0.5:1234",
		},
		{
			input:  "0.0.0.0:8080",
			expect: netx.InternalIp() + ":8080",
		},
		{
			input:  ":8080",
			expect: netx.InternalIp() + ":8080",
		},
		{
			input:  "",
			expect: netx.InternalIp(),
		},
	}

	for _, test := range tests {
		val := figureOutListenOn(test.input)
		assert.Equal(t, test.expect, val)
	}
}

func TestKeepAliveServerStart(t *testing.T) {
	pub := &mockEtcdPublisher{}
	server := &mockRpcServer{}
	keepalive := keepAliveServer{
		publisher: pub,
		Server:    server,
	}

	err := keepalive.Start(nil)

	assert.NoError(t, err)
	assert.Equal(t, 1, pub.keepaliveCalls)
	assert.Equal(t, 1, server.startCalls)
}

func TestKeepAliveServerStartKeepAliveError(t *testing.T) {
	pub := &mockEtcdPublisher{
		keepaliveErr: errors.New("keepalive error"),
	}
	server := &mockRpcServer{}
	keepalive := keepAliveServer{
		publisher: pub,
		Server:    server,
	}

	err := keepalive.Start(nil)

	assert.EqualError(t, err, "keepalive error")
	assert.Equal(t, 0, server.startCalls)
}

func TestKeepAliveServerPauseResumeEtcdRegister(t *testing.T) {
	pub := &mockEtcdPublisher{}
	keepalive := keepAliveServer{
		publisher: pub,
		Server:    &mockRpcServer{},
	}

	keepalive.PauseEtcdRegister()
	keepalive.ResumeEtcdRegister()

	assert.Equal(t, 1, pub.pauseCalls)
	assert.Equal(t, 1, pub.resumeCalls)
}

type mockEtcdPublisher struct {
	keepaliveCalls int
	pauseCalls     int
	resumeCalls    int
	keepaliveErr   error
}

func (m *mockEtcdPublisher) KeepAlive() error {
	m.keepaliveCalls++
	return m.keepaliveErr
}

func (m *mockEtcdPublisher) Pause() {
	m.pauseCalls++
}

func (m *mockEtcdPublisher) Resume() {
	m.resumeCalls++
}

type mockRpcServer struct {
	startCalls int
}

func (m *mockRpcServer) AddOptions(_ ...grpc.ServerOption) {}

func (m *mockRpcServer) AddStreamInterceptors(_ ...grpc.StreamServerInterceptor) {}

func (m *mockRpcServer) AddUnaryInterceptors(_ ...grpc.UnaryServerInterceptor) {}

func (m *mockRpcServer) SetName(_ string) {}

func (m *mockRpcServer) Start(_ RegisterFn) error {
	m.startCalls++
	return nil
}
