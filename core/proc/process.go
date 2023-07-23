package proc

import (
	"os"
	"path/filepath"
)

var (
	procName string
	pid      int
)

func init() {
	procName = filepath.Base(os.Args[0])
	pid = os.Getpid()
}

// Pid returns pid of current process.
func Pid() int {
	return pid
}

// ProcessName returns the processname, same as the command name.
func ProcessName() string {
	return procName
}
