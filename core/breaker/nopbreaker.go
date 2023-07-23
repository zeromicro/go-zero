package breaker

const noOpBreakerName = "nopBreaker"

type noOpBreaker struct{}

func newNoOpBreaker() Breaker {
	return noOpBreaker{}
}

func (b noOpBreaker) Name() string {
	return noOpBreakerName
}

func (b noOpBreaker) Allow() (Promise, error) {
	return nopPromise{}, nil
}

func (b noOpBreaker) Do(req func() error) error {
	return req()
}

func (b noOpBreaker) DoWithAcceptable(req func() error, _ Acceptable) error {
	return req()
}

func (b noOpBreaker) DoWithFallback(req func() error, _ func(err error) error) error {
	return req()
}

func (b noOpBreaker) DoWithFallbackAcceptable(req func() error, _ func(err error) error,
	_ Acceptable) error {
	return req()
}

type nopPromise struct{}

func (p nopPromise) Accept() {
}

func (p nopPromise) Reject(_ string) {
}
