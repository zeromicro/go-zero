package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/service"
	"github.com/tal-tech/go-zero/rest"
	"github.com/tal-tech/go-zero/rest/httpx"
)

var (
	port    = flag.Int("port", 3333, "the port to listen")
	timeout = flag.Int64("timeout", 0, "timeout of milliseconds")
)

type Request struct {
	User string `json:"user"`
}

func handleGet(w http.ResponseWriter, r *http.Request) {
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	var req Request
	err := httpx.Parse(r, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	httpx.OkJson(w, fmt.Sprintf("Content-Length: %d, UserLen: %d", r.ContentLength, len(req.User)))
}

func main() {
	flag.Parse()

	engine := rest.MustNewServer(rest.RestConf{
		ServiceConf: service.ServiceConf{
			Log: logx.LogConf{
				Mode: "console",
			},
		},
		Port:         *port,
		Timeout:      *timeout,
		MaxConns:     500,
		MaxBytes:     50,
		CpuThreshold: 500,
	})
	defer engine.Stop()

	engine.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/",
		Handler: handleGet,
	})
	engine.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/",
		Handler: handlePost,
	})
	engine.Start()
}
