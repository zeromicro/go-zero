package main

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/tal-tech/go-zero/core/cmdline"
	"github.com/tal-tech/go-zero/core/logx"
)

const (
	goroutineProfile = "goroutine"
	debugLevel       = 2
	timeFormat       = "0102150405"
)

func init() {
	go func() {
		// https://golang.org/pkg/os/signal/#Notify
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGUSR1, syscall.SIGTERM)

		for {
			v := <-signals
			switch v {
			case syscall.SIGUSR1:
				dumpGoroutines()
			case syscall.SIGTERM:
				gracefulStop(signals)
			default:
				logx.Error("Got unregistered signal:", v)
			}
		}
	}()
}

func dumpGoroutines() {
	command := path.Base(os.Args[0])
	pid := syscall.Getpid()
	dumpFile := path.Join(os.TempDir(), fmt.Sprintf("%s-%d-goroutines-%s.dump",
		command, pid, time.Now().Format(timeFormat)))

	logx.Infof("Got dump goroutine signal, printing goroutine profile to %s", dumpFile)

	if f, err := os.Create(dumpFile); err != nil {
		logx.Errorf("Failed to dump goroutine profile, error: %v", err)
	} else {
		defer f.Close()
		pprof.Lookup(goroutineProfile).WriteTo(f, debugLevel)
	}
}

func gracefulStop(signals chan os.Signal) {
	signal.Stop(signals)

	logx.Info("Got signal SIGTERM, shutting down...")

	time.Sleep(time.Second * 5)
	logx.Infof("Still alive after %v, going to force kill the process...", time.Second*5)
	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
}

func main() {
	cmdline.EnterToContinue()
}
