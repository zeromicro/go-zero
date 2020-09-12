package breaker

import "sync"

var (
	lock     sync.RWMutex
	breakers = make(map[string]Breaker)
)

func Do(name string, req func() error) error {
	return do(name, func(b Breaker) error {
		return b.Do(req)
	})
}

func DoWithAcceptable(name string, req func() error, acceptable Acceptable) error {
	return do(name, func(b Breaker) error {
		return b.DoWithAcceptable(req, acceptable)
	})
}

func DoWithFallback(name string, req func() error, fallback func(err error) error) error {
	return do(name, func(b Breaker) error {
		return b.DoWithFallback(req, fallback)
	})
}

func DoWithFallbackAcceptable(name string, req func() error, fallback func(err error) error,
	acceptable Acceptable) error {
	return do(name, func(b Breaker) error {
		return b.DoWithFallbackAcceptable(req, fallback, acceptable)
	})
}

func GetBreaker(name string) Breaker {
	lock.RLock()
	b, ok := breakers[name]
	lock.RUnlock()
	if ok {
		return b
	}

	lock.Lock()
	defer lock.Unlock()

	b = NewBreaker()
	breakers[name] = b
	return b
}

func NoBreakFor(name string) {
	lock.Lock()
	breakers[name] = newNoOpBreaker()
	lock.Unlock()
}

func do(name string, execute func(b Breaker) error) error {
	lock.RLock()
	b, ok := breakers[name]
	lock.RUnlock()
	if ok {
		return execute(b)
	}

	lock.Lock()
	b, ok = breakers[name]
	if !ok {
		b = NewBreaker(WithName(name))
		breakers[name] = b
	}
	lock.Unlock()

	return execute(b)
}
