package prometheus

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/3Rivers/go-zero/core/logx"
	"github.com/3Rivers/go-zero/core/threading"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var once sync.Once

func StartAgent(c Config) {
	once.Do(func() {
		if len(c.Host) == 0 {
			return
		}

		threading.GoSafe(func() {
			http.Handle(c.Path, promhttp.Handler())
			addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
			logx.Infof("Starting prometheus agent at %s", addr)
			if err := http.ListenAndServe(addr, nil); err != nil {
				logx.Error(err)
			}
		})
	})
}
