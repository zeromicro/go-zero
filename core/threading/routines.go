package threading

import (
	"bytes"
	"runtime"
	"strconv"

	"github.com/zeromicro/go-zero/core/rescue"
)

// GoSafe runs the given fn using another goroutine, recovers if fn panics.
func GoSafe(fn func()) {
	go RunSafe(fn)
}

// RoutineId is only for debug, never use it in production.
func RoutineId() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	// if error, just return 0
	n, _ := strconv.ParseUint(string(b), 10, 64)

	return n
}

// RunSafe runs the given fn, recovers if fn panics.
func RunSafe(fn func()) {
	defer rescue.Recover()

	fn()
}
