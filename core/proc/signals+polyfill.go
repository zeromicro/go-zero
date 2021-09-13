//go:build windows
// +build windows

package proc

import "context"

func Done() <-chan struct{} {
	return context.Background().Done()
}
