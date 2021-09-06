package internal

import (
	"errors"
	"github.com/tal-tech/go-zero/core/discov"
	"google.golang.org/grpc"
	"reflect"
	"testing"

	. "github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/netx"
)

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
	}

	for _, test := range tests {
		val := figureOutListenOn(test.input)
		assert.Equal(t, test.expect, val)
	}
}

func TestNewRpcPubServerWithEtcdAuth(t *testing.T) {
	pub := &discov.Publisher{}
	auth := false
	patches := ApplyFunc(discov.NewPublisherWithAuth, func(endpoints []string, user, pass, key, value string, opts ...discov.PublisherOption) *discov.Publisher {

		if user == "admin" && pass == "pass" {
			auth = true
		}
		return pub
	})
	defer patches.Reset()
	patches.ApplyMethod(reflect.TypeOf(pub), "KeepAliveWithAuth", func(*discov.Publisher) error {
		if auth {
			return nil
		} else {
			return errors.New("")
		}
	})

	patches.ApplyFunc(NewRpcServer, func(address string, opts ...ServerOption) Server {
		return &mockServer{}
	})

	s, err := NewRpcPubServerWithEtcdAuth([]string{}, "admin", "pass", "key", "listenon")

	assert.Nil(t, err)
	assert.Nil(t, s.Start(func(server *grpc.Server) {

	}))
}

type mockServer struct {
}

func (m *mockServer) AddOptions(options ...grpc.ServerOption) {

}

func (m *mockServer) AddStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) {

}

func (m *mockServer) AddUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) {

}

func (m *mockServer) SetName(string) {

}

func (m *mockServer) Start(register RegisterFn) error {
	return nil
}
