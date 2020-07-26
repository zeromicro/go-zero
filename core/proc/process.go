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

func Pid() int {
	return pid
}

func ProcessName() string {
	return procName
}
