//go:build linux

package stat

import (
	"flag"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zeromicro/go-zero/core/executors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/sysx"
)

const (
	clusterNameKey = "CLUSTER_NAME"
	testEnv        = "test.v"
	timeFormat     = "2006-01-02 15:04:05"
)

var (
	reporter     = logx.Alert
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

// Report reports given message.
func Report(msg string) {
	lock.RLock()
	fn := reporter
	lock.RUnlock()

	if fn != nil {
		reported := lessExecutor.DoOrDiscard(func() {
			var builder strings.Builder
			builder.WriteString(fmt.Sprintln(time.Now().Format(timeFormat)))
			if len(clusterName) > 0 {
				builder.WriteString(fmt.Sprintf("cluster: %s\n", clusterName))
			}
			builder.WriteString(fmt.Sprintf("host: %s\n", sysx.Hostname()))
			dp := atomic.SwapInt32(&dropped, 0)
			if dp > 0 {
				builder.WriteString(fmt.Sprintf("dropped: %d\n", dp))
			}
			builder.WriteString(strings.TrimSpace(msg))
			fn(builder.String())
		})
		if !reported {
			atomic.AddInt32(&dropped, 1)
		}
	}
}

// SetReporter sets the given reporter.
func SetReporter(fn func(string)) {
	lock.Lock()
	defer lock.Unlock()
	reporter = fn
}
