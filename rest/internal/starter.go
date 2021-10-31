package internal

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/tal-tech/go-zero/core/proc"
)

// StartHttp starts a http server.
func StartHttp(host string, port int, handler http.Handler) error {
	return start(host, port, handler, nil, func(srv *http.Server) error {
		return srv.ListenAndServe()
	})
}

// StartHttps starts a https server.
func StartHttps(host string, port int, certFile, keyFile string, tlsConfig *tls.Config, handler http.Handler) error {
	return start(host, port, handler, tlsConfig, func(srv *http.Server) error {
		// certFile and keyFile are set in buildHttpsServer
		return srv.ListenAndServeTLS(certFile, keyFile)
	})
}

func start(host string, port int, handler http.Handler, tlsConfig *tls.Config, run func(srv *http.Server) error) (err error) {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: handler,
	}
	if tlsConfig != nil {
		server.TLSConfig = tlsConfig
	}
	waitForCalled := proc.AddWrapUpListener(func() {
		server.Shutdown(context.Background())
	})
	defer func() {
		if err == http.ErrServerClosed {
			waitForCalled()
		}
	}()

	return run(server)
}
