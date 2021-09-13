//go:build windows
// +build windows

package proc

func StartProfile() Stopper {
	return noopStopper
}
