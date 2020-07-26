package httpserver

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

func StartHttp(host string, port int, handler http.Handler) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	server := buildHttpServer(addr, handler)
	return StartServer(server)
}

func StartHttps(host string, port int, certFile, keyFile string, handler http.Handler) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	if server, err := buildHttpsServer(addr, handler, certFile, keyFile); err != nil {
		return err
	} else {
		return StartServer(server)
	}
}

func buildHttpServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{Addr: addr, Handler: handler}
}

func buildHttpsServer(addr string, handler http.Handler, certFile, keyFile string) (*http.Server, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	config := tls.Config{Certificates: []tls.Certificate{cert}}
	return &http.Server{
		Addr:      addr,
		Handler:   handler,
		TLSConfig: &config,
	}, nil
}
