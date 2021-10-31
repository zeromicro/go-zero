package internal

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/tal-tech/go-zero/core/proc"
)

var (
	cipherSuites = []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	}
)

// StartHttp starts a http server.
func StartHttp(host string, port int, handler http.Handler) error {
	return start(host, port, handler, nil, func(srv *http.Server) error {
		return srv.ListenAndServe()
	})
}

// StartHttps starts a https server.
func StartHttps(host string, port int, certFile, keyFile string, handler http.Handler) error {
	safeTlsConfig := &tls.Config{
		CipherSuites: cipherSuites,
	}
	return start(host, port, handler, safeTlsConfig, func(srv *http.Server) error {
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
