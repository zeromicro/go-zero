package prof

import (
	"bytes"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

type (
	profileSlot struct {
		lifecount int64
		lastcount int64
		lifecycle int64
		lastcycle int64
	}

	profileCenter struct {
		lock  sync.RWMutex
		slots map[string]*profileSlot
	}
)

const flushInterval = 5 * time.Minute

var (
	pc = &profileCenter{
		slots: make(map[string]*profileSlot),
	}
	once sync.Once
)

func report(name string, duration time.Duration) {
	updated := func() bool {
		pc.lock.RLock()
		defer pc.lock.RUnlock()

		slot, ok := pc.slots[name]
		if ok {
			atomic.AddInt64(&slot.lifecount, 1)
			atomic.AddInt64(&slot.lastcount, 1)
			atomic.AddInt64(&slot.lifecycle, int64(duration))
			atomic.AddInt64(&slot.lastcycle, int64(duration))
		}
		return ok
	}()

	if !updated {
		func() {
			pc.lock.Lock()
			defer pc.lock.Unlock()

			pc.slots[name] = &profileSlot{
				lifecount: 1,
				lastcount: 1,
				lifecycle: int64(duration),
				lastcycle: int64(duration),
			}
		}()
	}

	once.Do(flushRepeatly)
}

func flushRepeatly() {
	threading.GoSafe(func() {
		for {
			time.Sleep(flushInterval)
			logx.Stat(generateReport())
		}
	})
}

func generateReport() string {
	var buffer bytes.Buffer
	buffer.WriteString("Profiling report\n")
	var data [][]string
	calcFn := func(total, count int64) string {
		if count == 0 {
			return "-"
		}

		return (time.Duration(total) / time.Duration(count)).String()
	}

	func() {
		pc.lock.Lock()
		defer pc.lock.Unlock()

		for key, slot := range pc.slots {
			data = append(data, []string{
				key,
				strconv.FormatInt(slot.lifecount, 10),
				calcFn(slot.lifecycle, slot.lifecount),
				strconv.FormatInt(slot.lastcount, 10),
				calcFn(slot.lastcycle, slot.lastcount),
			})

			// reset the data for last cycle
			slot.lastcount = 0
			slot.lastcycle = 0
		}
	}()

	table := tablewriter.NewWriter(&buffer)
	table.SetHeader([]string{"QUEUE", "LIFECOUNT", "LIFECYCLE", "LASTCOUNT", "LASTCYCLE"})
	table.SetBorder(false)
	table.AppendBulk(data)
	table.Render()

	return buffer.String()
}
