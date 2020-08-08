package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/tal-tech/go-zero/core/fx"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/service"
	"github.com/tal-tech/go-zero/core/stat"
	"github.com/tal-tech/go-zero/rest"
)

const duration = time.Millisecond

func main() {
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			fmt.Printf("cpu: %d\n", stat.CpuUsage())
		}
	}()

	logx.Disable()
	engine := rest.MustNewServer(rest.RestConf{
		ServiceConf: service.ServiceConf{
			Log: logx.LogConf{
				Mode: "console",
			},
		},
		Host:         "0.0.0.0",
		Port:         3333,
		CpuThreshold: 800,
	})
	defer engine.Stop()
	engine.AddRoute(rest.Route{
		Method: http.MethodGet,
		Path:   "/",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			if err := fx.DoWithTimeout(func() error {
				job(duration)
				return nil
			}, time.Millisecond*100); err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
			}
		},
	})
	engine.Start()
}

func job(duration time.Duration) {
	done := make(chan int)

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				select {
				case <-done:
					return
				default:
				}
			}
		}()
	}

	time.Sleep(duration)
	close(done)
}
