//go:build windows
// +build windows

package fs

import "os"

func CloseOnExec(*os.File) {
}
