//@File     pusher.go
//@Time     2024/5/10
//@Author   #Suyghur,

package prometheus

import (
	"time"

	prom "github.com/prometheus/client_golang/prometheus"
	prompush "github.com/prometheus/client_golang/prometheus/push"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/core/threading"
	"github.com/zeromicro/go-zero/core/timex"
)

var (
	pusher       *prompush.Pusher
	pusherTicker timex.Ticker
	pusherDone   *syncx.DoneChan
)

// StartPusher starts a pusher to push metrics to prometheus pushgateway.
// see https://github.com/prometheus/pushgateway for details.
func StartPusher(c MetricsPusherConfig) {
	if len(c.Url) == 0 {
		return
	}

	if len(c.JobName) == 0 {
		return
	}

	syncx.Once(func() {
		pusher = prompush.New(c.Url, c.JobName).Gatherer(prom.DefaultGatherer)
		pusherTicker = timex.NewTicker(time.Duration(c.Interval) * time.Second)
		pusherDone = syncx.NewDoneChan()

		threading.GoSafe(func() {
			for {
				select {
				case <-pusherTicker.Chan():
					_ = pusher.Push()
				case <-pusherDone.Done():
					pusherTicker.Stop()
					return
				}
			}
		})

		proc.AddShutdownListener(func() {
			pusherDone.Close()
		})
	})
}
