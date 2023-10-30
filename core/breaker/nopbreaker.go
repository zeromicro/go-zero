package breaker

const nopBreakerName = "nopBreaker"

type nopBreaker struct{}

func newNopBreaker() Breaker {
	return nopBreaker{}
}

func (b nopBreaker) Name() string {
	return nopBreakerName
}

func (b nopBreaker) Allow() (Promise, error) {
	return nopPromise{}, nil
}

func (b nopBreaker) Do(req func() error) error {
	return req()
}

func (b nopBreaker) DoWithAcceptable(req func() error, _ Acceptable) error {
	return req()
}

func (b nopBreaker) DoWithFallback(req func() error, _ func(err error) error) error {
	return req()
}

func (b nopBreaker) DoWithFallbackAcceptable(req func() error, _ func(err error) error,
	_ Acceptable) error {
	return req()
}

type nopPromise struct{}

func (p nopPromise) Accept() {
}

func (p nopPromise) Reject(_ string) {
}
