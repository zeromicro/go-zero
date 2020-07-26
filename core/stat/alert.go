// +build linux

package stat

import (
	"flag"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"zero/core/executors"
	"zero/core/proc"
	"zero/core/sysx"
	"zero/core/timex"
	"zero/core/utils"
)

const (
	clusterNameKey = "CLUSTER_NAME"
	testEnv        = "test.v"
	timeFormat     = "2006-01-02 15:04:05"
)

var (
	reporter     = utils.Report
	lock         sync.RWMutex
	lessExecutor = executors.NewLessExecutor(time.Minute * 5)
	dropped      int32
	clusterName  = proc.Env(clusterNameKey)
)

func init() {
	if flag.Lookup(testEnv) != nil {
		SetReporter(nil)
	}
}

func Report(msg string) {
	lock.RLock()
	fn := reporter
	lock.RUnlock()

	if fn != nil {
		reported := lessExecutor.DoOrDiscard(func() {
			var builder strings.Builder
			fmt.Fprintf(&builder, "%s\n", timex.Time().Format(timeFormat))
			if len(clusterName) > 0 {
				fmt.Fprintf(&builder, "cluster: %s\n", clusterName)
			}
			fmt.Fprintf(&builder, "host: %s\n", sysx.Hostname())
			dp := atomic.SwapInt32(&dropped, 0)
			if dp > 0 {
				fmt.Fprintf(&builder, "dropped: %d\n", dp)
			}
			builder.WriteString(strings.TrimSpace(msg))
			fn(builder.String())
		})
		if !reported {
			atomic.AddInt32(&dropped, 1)
		}
	}
}

func SetReporter(fn func(string)) {
	lock.Lock()
	defer lock.Unlock()
	reporter = fn
}
