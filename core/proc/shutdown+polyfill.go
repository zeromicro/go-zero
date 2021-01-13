// +build windows

package proc

import "time"

func AddShutdownListener(fn func()) func() {
	return func() {}
}

func AddWrapUpListener(fn func()) func() {
	return func() {}
}

func SetTimeoutToForceQuit(duration time.Duration) {
}
