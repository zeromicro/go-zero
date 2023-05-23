package cache

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/stat"
	"github.com/zeromicro/go-zero/core/stringx"
	"github.com/zeromicro/go-zero/core/threading"
)

const (
	timingWheelSlots = 300
	cleanWorkers     = 5
	taskKeyLen       = 8
)

var (
	// use atomic to avoid data race in unit tests
	timingWheel atomic.Value
	taskRunner  = threading.NewTaskRunner(cleanWorkers)
)

type delayTask struct {
	delay time.Duration
	task  func() error
	keys  []string
}

func init() {
	tw, err := collection.NewTimingWheel(time.Second, timingWheelSlots, clean)
	logx.Must(err)
	timingWheel.Store(tw)

	proc.AddShutdownListener(func() {
		if err := tw.Drain(clean); err != nil {
			logx.Errorf("failed to drain timing wheel: %v", err)
		}
	})
}

// AddCleanTask adds a clean task on given keys.
func AddCleanTask(task func() error, keys ...string) {
	tw := timingWheel.Load().(*collection.TimingWheel)
	if err := tw.SetTimer(stringx.Randn(taskKeyLen), delayTask{
		delay: time.Second,
		task:  task,
		keys:  keys,
	}, time.Second); err != nil {
		logx.Errorf("failed to set timer for keys: %q, error: %v", formatKeys(keys), err)
	}
}

func clean(key, value any) {
	taskRunner.Schedule(func() {
		dt := value.(delayTask)
		err := dt.task()
		if err == nil {
			return
		}

		next, ok := nextDelay(dt.delay)
		if ok {
			dt.delay = next
			tw := timingWheel.Load().(*collection.TimingWheel)
			if err = tw.SetTimer(key, dt, next); err != nil {
				logx.Errorf("failed to set timer for key: %s, error: %v", key, err)
			}
		} else {
			msg := fmt.Sprintf("retried but failed to clear cache with keys: %q, error: %v",
				formatKeys(dt.keys), err)
			logx.Error(msg)
			stat.Report(msg)
		}
	})
}

func nextDelay(delay time.Duration) (time.Duration, bool) {
	switch delay {
	case time.Second:
		return time.Second * 5, true
	case time.Second * 5:
		return time.Minute, true
	case time.Minute:
		return time.Minute * 5, true
	case time.Minute * 5:
		return time.Hour, true
	default:
		return 0, false
	}
}
