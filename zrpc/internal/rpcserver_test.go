package internal

import (
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/internal/mock"
	"google.golang.org/grpc"
)

func TestRpcServer(t *testing.T) {
	server := NewRpcServer("localhost:54321", WithRpcHealth(true))
	server.SetName("mock")
	var wg, wgDone sync.WaitGroup
	var grpcServer *grpc.Server
	var lock sync.Mutex
	wg.Add(2)
	wgDone.Add(2)
	go func() {
		err := server.Start(func(server *grpc.Server) {
			lock.Lock()
			mock.RegisterDepositServiceServer(server, new(mock.DepositServer))
			grpcServer = server
			lock.Unlock()
			wg.Done()
		})
		assert.Nil(t, err)
		wgDone.Done()
	}()

	go func() {
		listener, err := net.Listen("tcp", "localhost:54322")
		assert.Nil(t, err)
		serverWithListener := NewRpcServer(listener.Addr().String(), WithRpcHealth(true))
		err = serverWithListener.StartWithListener(listener, func(server *grpc.Server) {
			lock.Lock()
			mock.RegisterDepositServiceServer(server, new(mock.DepositServer))
			grpcServer = server
			lock.Unlock()
			wg.Done()
		})
		assert.Nil(t, err)
		wgDone.Done()
	}()

	wg.Wait()
	time.Sleep(100 * time.Millisecond)

	lock.Lock()
	grpcServer.GracefulStop()
	lock.Unlock()

	proc.Shutdown()
	wgDone.Wait()
}

func TestRpcServer_WithBadAddress(t *testing.T) {
	server := NewRpcServer("localhost:111111", WithRpcHealth(true))
	server.SetName("mock")
	err := server.Start(func(server *grpc.Server) {
		mock.RegisterDepositServiceServer(server, new(mock.DepositServer))
	})
	assert.NotNil(t, err)

	proc.WrapUp()
}
