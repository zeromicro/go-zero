//go:build windows

package proc

func StartProfile() Stopper {
	return noopStopper
}
