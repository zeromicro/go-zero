package internal

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/syncx"
)

func TestStartHttp(t *testing.T) {
	svr := httptest.NewUnstartedServer(http.NotFoundHandler())
	fields := strings.Split(svr.Listener.Addr().String(), ":")
	port, err := strconv.Atoi(fields[1])
	assert.Nil(t, err)
	err = StartHttp(fields[0], port, http.NotFoundHandler(), &mockProbe{}, func(svr *http.Server) {
		svr.IdleTimeout = 0
	})
	assert.NotNil(t, err)
	proc.WrapUp()
}

func TestStartHttps(t *testing.T) {
	svr := httptest.NewTLSServer(http.NotFoundHandler())
	fields := strings.Split(svr.Listener.Addr().String(), ":")
	port, err := strconv.Atoi(fields[1])
	assert.Nil(t, err)
	err = StartHttps(fields[0], port, "", "", http.NotFoundHandler(), &mockProbe{}, func(svr *http.Server) {
		svr.IdleTimeout = 0
	})
	assert.NotNil(t, err)
	proc.WrapUp()
}
func TestStartWithShutdownListener(t *testing.T) {
	probe := &mockProbe{}
	shutdownCalled := make(chan struct{})
	serverStarted := make(chan struct{})
	serverClosed := make(chan struct{})

	run := func(svr *http.Server) error {
		close(serverStarted)
		<-shutdownCalled
		return http.ErrServerClosed
	}

	go func() {
		err := start("localhost", 8888, http.NotFoundHandler(), probe, run)
		assert.Equal(t, http.ErrServerClosed, err)
		close(serverClosed)
	}()

	select {
	case <-serverStarted:
		assert.True(t, probe.IsReady(), "server should be marked as ready")
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for server to start")
	}

	proc.WrapUp()
	time.Sleep(time.Millisecond * 50)
	close(shutdownCalled)
}

type mockProbe struct {
	ready syncx.AtomicBool
}

func (m *mockProbe) MarkReady() {
	m.ready.Set(true)
}

func (m *mockProbe) MarkNotReady() {
	m.ready.Set(false)
}

func (m *mockProbe) IsReady() bool {
	return m.ready.True()
}

func (m *mockProbe) Name() string { return "" }
