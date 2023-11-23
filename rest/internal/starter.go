package internal

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/internal/health"
)

const probeNamePrefix = "rest"

// StartOption defines the method to customize http.Server.
type StartOption func(svr *http.Server)

// StartHttp starts a http server.
func StartHttp(host string, port int, handler http.Handler, opts ...StartOption) error {
	return start(host, port, handler, func(svr *http.Server) error {
		return svr.ListenAndServe()
	}, opts...)
}

// StartHttps starts a https server.
func StartHttps(host string, port int, certFile, keyFile string, handler http.Handler,
	opts ...StartOption) error {
	return start(host, port, handler, func(svr *http.Server) error {
		// certFile and keyFile are set in buildHttpsServer
		return svr.ListenAndServeTLS(certFile, keyFile)
	}, opts...)
}

func start(host string, port int, handler http.Handler, run func(svr *http.Server) error,
	opts ...StartOption) (err error) {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: handler,
	}
	for _, opt := range opts {
		opt(server)
	}
	healthManager := health.NewHealthManager(fmt.Sprintf("%s-%s:%d", probeNamePrefix, host, port))

	waitForCalled := proc.AddShutdownListener(func() {
		healthManager.MarkNotReady()
		if e := server.Shutdown(context.Background()); e != nil {
			logx.Error(e)
		}
	})
	defer func() {
		if errors.Is(err, http.ErrServerClosed) {
			waitForCalled()
		}
	}()

	healthManager.MarkReady()
	health.AddProbe(healthManager)
	return run(server)
}
