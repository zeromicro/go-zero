// +build windows

package proc

import "time"

func AddShutdownListener(fn func()) func() {
	return fn
}

func AddWrapUpListener(fn func()) func() {
	return fn
}

func SetTimeoutToForceQuit(duration time.Duration) {
}
