package proc

import (
	"os"
	"strconv"
	"sync"
)

var (
	envs    = make(map[string]string)
	envLock sync.RWMutex
)

// Env returns the value of the given environment variable.
func Env(name string) string {
	envLock.RLock()
	val, ok := envs[name]
	envLock.RUnlock()

	if ok {
		return val
	}

	val = os.Getenv(name)
	envLock.Lock()
	envs[name] = val
	envLock.Unlock()

	return val
}

// EnvInt returns an int value of the given environment variable.
func EnvInt(name string) (int, bool) {
	val := Env(name)
	if len(val) == 0 {
		return 0, false
	}

	n, err := strconv.Atoi(val)
	if err != nil {
		return 0, false
	}

	return n, true
}
