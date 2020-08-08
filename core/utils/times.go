package utils

import (
	"fmt"
	"time"

	"github.com/tal-tech/go-zero/core/timex"
)

type ElapsedTimer struct {
	start time.Duration
}

func NewElapsedTimer() *ElapsedTimer {
	return &ElapsedTimer{
		start: timex.Now(),
	}
}

func (et *ElapsedTimer) Duration() time.Duration {
	return timex.Since(et.start)
}

func (et *ElapsedTimer) Elapsed() string {
	return timex.Since(et.start).String()
}

func (et *ElapsedTimer) ElapsedMs() string {
	return fmt.Sprintf("%.1fms", float32(timex.Since(et.start))/float32(time.Millisecond))
}

func CurrentMicros() int64 {
	return time.Now().UnixNano() / int64(time.Microsecond)
}

func CurrentMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
