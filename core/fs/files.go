//go:build linux || darwin || freebsd

package fs

import (
	"os"
	"syscall"
)

// CloseOnExec makes sure closing the file on process forking.
func CloseOnExec(file *os.File) {
	if file != nil {
		syscall.CloseOnExec(int(file.Fd()))
	}
}
