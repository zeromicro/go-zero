package prof

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

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

var pc = &profileCenter{
	slots: make(map[string]*profileSlot),
}

func init() {
	flushRepeatedly()
}

func flushRepeatedly() {
	threading.GoSafe(func() {
		for {
			time.Sleep(flushInterval)
			logx.Stat(generateReport())
		}
	})
}

func report(name string, duration time.Duration) {
	slot := loadOrStoreSlot(name, duration)

	atomic.AddInt64(&slot.lifecount, 1)
	atomic.AddInt64(&slot.lastcount, 1)
	atomic.AddInt64(&slot.lifecycle, int64(duration))
	atomic.AddInt64(&slot.lastcycle, int64(duration))
}

func loadOrStoreSlot(name string, duration time.Duration) *profileSlot {
	pc.lock.RLock()
	slot, ok := pc.slots[name]
	pc.lock.RUnlock()

	if ok {
		return slot
	}

	pc.lock.Lock()
	defer pc.lock.Unlock()

	// double-check
	if slot, ok = pc.slots[name]; ok {
		return slot
	}

	slot = &profileSlot{}
	pc.slots[name] = slot
	return slot
}

func generateReport() string {
	var builder strings.Builder
	builder.WriteString("Profiling report\n")
	builder.WriteString("QUEUE,LIFECOUNT,LIFECYCLE,LASTCOUNT,LASTCYCLE\n")

	calcFn := func(total, count int64) string {
		if count == 0 {
			return "-"
		}
		return (time.Duration(total) / time.Duration(count)).String()
	}

	pc.lock.Lock()
	for key, slot := range pc.slots {
		builder.WriteString(fmt.Sprintf("%s,%d,%s,%d,%s\n",
			key,
			slot.lifecount,
			calcFn(slot.lifecycle, slot.lifecount),
			slot.lastcount,
			calcFn(slot.lastcycle, slot.lastcount),
		))

		// reset last cycle stats
		atomic.StoreInt64(&slot.lastcount, 0)
		atomic.StoreInt64(&slot.lastcycle, 0)
	}
	pc.lock.Unlock()

	return builder.String()
}
