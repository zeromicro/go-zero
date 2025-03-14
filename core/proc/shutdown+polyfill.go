//go:build windows

package proc

import "time"

// ShutdownConf is empty on windows.
type ShutdownConf struct{}

// AddShutdownListener returns fn itself on windows, lets callers call fn on their own.
func AddShutdownListener(fn func()) func() {
	return fn
}

// AddWrapUpListener returns fn itself on windows, lets callers call fn on their own.
func AddWrapUpListener(fn func()) func() {
	return fn
}

// SetTimeToForceQuit does nothing on windows.
func SetTimeToForceQuit(duration time.Duration) {
}

// Setup does nothing on windows.
func Setup(conf ShutdownConf) {
}

// Shutdown does nothing on windows.
func Shutdown() {
}

// WrapUp does nothing on windows.
func WrapUp() {
}
