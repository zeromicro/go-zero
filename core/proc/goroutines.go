//go:build linux || darwin || freebsd

package proc

import (
	"fmt"
	"os"
	"path"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	goroutineProfile = "goroutine"
	debugLevel       = 2
)

type creator interface {
	Create(name string) (file *os.File, err error)
}

func dumpGoroutines(ctor creator) {
	command := path.Base(os.Args[0])
	pid := syscall.Getpid()
	dumpFile := path.Join(os.TempDir(), fmt.Sprintf("%s-%d-goroutines-%s.dump",
		command, pid, time.Now().Format(timeFormat)))

	logx.Infof("Got dump goroutine signal, printing goroutine profile to %s", dumpFile)

	if f, err := ctor.Create(dumpFile); err != nil {
		logx.Errorf("Failed to dump goroutine profile, error: %v", err)
	} else {
		defer f.Close()
		pprof.Lookup(goroutineProfile).WriteTo(f, debugLevel)
	}
}

type fileCreator struct{}

func (fc fileCreator) Create(name string) (file *os.File, err error) {
	return os.Create(name)
}
