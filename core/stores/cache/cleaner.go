package cache

import (
	"fmt"
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
	timingWheel *collection.TimingWheel
	taskRunner  = threading.NewTaskRunner(cleanWorkers)
)

type delayTask struct {
	delay time.Duration
	task  func() error
	keys  []string
}

func init() {
	var err error
	timingWheel, err = collection.NewTimingWheel(time.Second, timingWheelSlots, clean)
	logx.Must(err)

	proc.AddShutdownListener(func() {
		timingWheel.Drain(clean)
	})
}

// AddCleanTask adds a clean task on given keys.
func AddCleanTask(task func() error, keys ...string) {
	timingWheel.SetTimer(stringx.Randn(taskKeyLen), delayTask{
		delay: time.Second,
		task:  task,
		keys:  keys,
	}, time.Second)
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
			timingWheel.SetTimer(key, dt, next)
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
