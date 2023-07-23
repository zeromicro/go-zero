//go:build windows

package fs

import "os"

func CloseOnExec(*os.File) {
}
