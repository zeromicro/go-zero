package breaker

import (
	"context"
	"sync"
)

var (
	lock     sync.RWMutex
	breakers = make(map[string]Breaker)
)

// Do calls Breaker.Do on the Breaker with given name.
func Do(name string, req func() error) error {
	return do(name, func(b Breaker) error {
		return b.Do(req)
	})
}

// DoCtx calls Breaker.DoCtx on the Breaker with given name.
func DoCtx(ctx context.Context, name string, req func() error) error {
	return do(name, func(b Breaker) error {
		return b.DoCtx(ctx, req)
	})
}

// DoWithAcceptable calls Breaker.DoWithAcceptable on the Breaker with given name.
func DoWithAcceptable(name string, req func() error, acceptable Acceptable) error {
	return do(name, func(b Breaker) error {
		return b.DoWithAcceptable(req, acceptable)
	})
}

// DoWithAcceptableCtx calls Breaker.DoWithAcceptableCtx on the Breaker with given name.
func DoWithAcceptableCtx(ctx context.Context, name string, req func() error,
	acceptable Acceptable) error {
	return do(name, func(b Breaker) error {
		return b.DoWithAcceptableCtx(ctx, req, acceptable)
	})
}

// DoWithFallback calls Breaker.DoWithFallback on the Breaker with given name.
func DoWithFallback(name string, req func() error, fallback Fallback) error {
	return do(name, func(b Breaker) error {
		return b.DoWithFallback(req, fallback)
	})
}

// DoWithFallbackCtx calls Breaker.DoWithFallbackCtx on the Breaker with given name.
func DoWithFallbackCtx(ctx context.Context, name string, req func() error, fallback Fallback) error {
	return do(name, func(b Breaker) error {
		return b.DoWithFallbackCtx(ctx, req, fallback)
	})
}

// DoWithFallbackAcceptable calls Breaker.DoWithFallbackAcceptable on the Breaker with given name.
func DoWithFallbackAcceptable(name string, req func() error, fallback Fallback,
	acceptable Acceptable) error {
	return do(name, func(b Breaker) error {
		return b.DoWithFallbackAcceptable(req, fallback, acceptable)
	})
}

// DoWithFallbackAcceptableCtx calls Breaker.DoWithFallbackAcceptableCtx on the Breaker with given name.
func DoWithFallbackAcceptableCtx(ctx context.Context, name string, req func() error,
	fallback Fallback, acceptable Acceptable) error {
	return do(name, func(b Breaker) error {
		return b.DoWithFallbackAcceptableCtx(ctx, req, fallback, acceptable)
	})
}

// GetBreaker returns the Breaker with the given name.
func GetBreaker(name string) Breaker {
	lock.RLock()
	b, ok := breakers[name]
	lock.RUnlock()
	if ok {
		return b
	}

	lock.Lock()
	b, ok = breakers[name]
	if !ok {
		b = NewBreaker(WithName(name))
		breakers[name] = b
	}
	lock.Unlock()

	return b
}

// NoBreakerFor disables the circuit breaker for the given name.
func NoBreakerFor(name string) {
	lock.Lock()
	breakers[name] = NopBreaker()
	lock.Unlock()
}

func do(name string, execute func(b Breaker) error) error {
	return execute(GetBreaker(name))
}
