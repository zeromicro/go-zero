package internal

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tal-tech/go-zero/core/proc"
)

func StartHttp(host string, port int, handler http.Handler) error {
	return start(host, port, handler, func(srv *http.Server) error {
		return srv.ListenAndServe()
	})
}

func StartHttps(host string, port int, certFile, keyFile string, handler http.Handler) error {
	return start(host, port, handler, func(srv *http.Server) error {
		// certFile and keyFile are set in buildHttpsServer
		return srv.ListenAndServeTLS(certFile, keyFile)
	})
}

func start(host string, port int, handler http.Handler, run func(srv *http.Server) error) error {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: handler,
	}
	waitForCalled := proc.AddWrapUpListener(func() {
		server.Shutdown(context.Background())
	})
	defer waitForCalled()

	return run(server)
}
