//go:build linux || darwin

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
