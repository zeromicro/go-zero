//go:build linux || darwin || freebsd

package proc

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
)

const (
	// defaultWrapUpTime is the default time to wait before calling wrap up listeners.
	defaultWrapUpTime = time.Second
	// defaultWaitTime is the default time to wait before force quitting.
	// why we use 5500 milliseconds is because most of our queues are blocking mode with 5 seconds
	defaultWaitTime = 5500 * time.Millisecond
)

var (
	wrapUpListeners   = new(listenerManager)
	shutdownListeners = new(listenerManager)
	wrapUpTime        = defaultWrapUpTime
	waitTime          = defaultWaitTime
	shutdownLock      sync.Mutex
)

// ShutdownConf defines the shutdown configuration for the process.
type ShutdownConf struct {
	// WrapUpTime is the time to wait before calling shutdown listeners.
	WrapUpTime time.Duration `json:",default=1s"`
	// WaitTime is the time to wait before force quitting.
	WaitTime time.Duration `json:",default=5.5s"`
}

// AddShutdownListener adds fn as a shutdown listener.
// The returned func can be used to wait for fn getting called.
func AddShutdownListener(fn func()) (waitForCalled func()) {
	return shutdownListeners.addListener(fn)
}

// AddWrapUpListener adds fn as a wrap up listener.
// The returned func can be used to wait for fn getting called.
func AddWrapUpListener(fn func()) (waitForCalled func()) {
	return wrapUpListeners.addListener(fn)
}

// SetTimeToForceQuit sets the waiting time before force quitting.
func SetTimeToForceQuit(duration time.Duration) {
	shutdownLock.Lock()
	defer shutdownLock.Unlock()
	waitTime = duration
}

func Setup(conf ShutdownConf) {
	shutdownLock.Lock()
	defer shutdownLock.Unlock()

	if conf.WrapUpTime > 0 {
		wrapUpTime = conf.WrapUpTime
	}
	if conf.WaitTime > 0 {
		waitTime = conf.WaitTime
	}
}

// Shutdown calls the registered shutdown listeners, only for test purpose.
func Shutdown() {
	shutdownListeners.notifyListeners()
}

// WrapUp wraps up the process, only for test purpose.
func WrapUp() {
	wrapUpListeners.notifyListeners()
}

func gracefulStop(signals chan os.Signal, sig syscall.Signal) {
	signal.Stop(signals)

	logx.Infof("Got signal %d, shutting down...", sig)
	go wrapUpListeners.notifyListeners()

	time.Sleep(wrapUpTime)
	go shutdownListeners.notifyListeners()

	shutdownLock.Lock()
	remainingTime := waitTime - wrapUpTime
	shutdownLock.Unlock()

	time.Sleep(remainingTime)
	logx.Infof("Still alive after %v, going to force kill the process...", waitTime)
	_ = syscall.Kill(syscall.Getpid(), sig)
}

type listenerManager struct {
	lock      sync.Mutex
	waitGroup sync.WaitGroup
	listeners []func()
}

func (lm *listenerManager) addListener(fn func()) (waitForCalled func()) {
	lm.waitGroup.Add(1)

	lm.lock.Lock()
	lm.listeners = append(lm.listeners, func() {
		defer lm.waitGroup.Done()
		fn()
	})
	lm.lock.Unlock()

	// we can return lm.waitGroup.Wait directly,
	// but we want to make the returned func more readable.
	// creating an extra closure would be negligible in practice.
	return func() {
		lm.waitGroup.Wait()
	}
}

func (lm *listenerManager) notifyListeners() {
	lm.lock.Lock()
	defer lm.lock.Unlock()

	group := threading.NewRoutineGroup()
	for _, listener := range lm.listeners {
		group.RunSafe(listener)
	}
	group.Wait()

	lm.listeners = nil
}
