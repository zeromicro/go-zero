//go:build linux || darwin || freebsd

package proc

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	profileDuration = time.Minute
	timeFormat      = "0102150405"
)

var done = make(chan struct{})

func init() {
	go func() {
		// https://golang.org/pkg/os/signal/#Notify
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGTERM, syscall.SIGINT)

		for {
			v := <-signals
			switch v {
			case syscall.SIGUSR1:
				dumpGoroutines(fileCreator{})
			case syscall.SIGUSR2:
				profiler := StartProfile()
				time.AfterFunc(profileDuration, profiler.Stop)
			case syscall.SIGTERM:
				stopOnSignal()
				gracefulStop(signals, syscall.SIGTERM)
			case syscall.SIGINT:
				stopOnSignal()
				gracefulStop(signals, syscall.SIGINT)
			default:
				logx.Error("Got unregistered signal:", v)
			}
		}
	}()
}

// Done returns the channel that notifies the process quitting.
func Done() <-chan struct{} {
	return done
}

func stopOnSignal() {
	select {
	case <-done:
		// already closed
	default:
		close(done)
	}
}
