package main

import (
	"flag"
	"math"
	"net/http"
	"time"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/service"
	"github.com/tal-tech/go-zero/rest"
	"github.com/tal-tech/go-zero/rest/httpx"
)

var (
	port    = flag.Int("port", 3333, "the port to listen")
	timeout = flag.Int64("timeout", 1000, "timeout of milliseconds")
	cpu     = flag.Int64("cpu", 500, "cpu threshold")
)

type Request struct {
	User string `form:"user,optional"`
}

func handle(w http.ResponseWriter, r *http.Request) {
	var req Request
	err := httpx.Parse(r, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var result float64
	for i := 0; i < 30000; i++ {
		result += math.Sqrt(float64(i))
	}
	time.Sleep(time.Millisecond * 5)
	httpx.OkJson(w, result)
}

func main() {
	flag.Parse()

	logx.Disable()
	engine := rest.MustNewServer(rest.RestConf{
		ServiceConf: service.ServiceConf{
			Log: logx.LogConf{
				Mode: "console",
			},
		},
		Port:         *port,
		Timeout:      *timeout,
		CpuThreshold: *cpu,
	})
	defer engine.Stop()

	engine.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/",
		Handler: handle,
	})
	engine.Start()
}
