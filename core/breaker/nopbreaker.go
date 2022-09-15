package breaker

const noOpBreakerName = "nopBreaker"

// NewBreaker returns a Breaker object.
// opts can be used to customize the Breaker.
func NewNoOpBreaker() Breaker {
	return newNoOpBreaker()
}

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

func (b noOpBreaker) DoWithAcceptable(req func() error, acceptable Acceptable) error {
	return req()
}

func (b noOpBreaker) DoWithFallback(req func() error, fallback func(err error) error) error {
	return req()
}

func (b noOpBreaker) DoWithFallbackAcceptable(req func() error, fallback func(err error) error,
	acceptable Acceptable) error {
	return req()
}

type nopPromise struct{}

func (p nopPromise) Accept() {
}

func (p nopPromise) Reject(reason string) {
}
