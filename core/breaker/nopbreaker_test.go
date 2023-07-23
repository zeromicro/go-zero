package breaker

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNopBreaker(t *testing.T) {
	b := newNoOpBreaker()
	assert.Equal(t, noOpBreakerName, b.Name())
	p, err := b.Allow()
	assert.Nil(t, err)
	p.Accept()
	for i := 0; i < 1000; i++ {
		p, err := b.Allow()
		assert.Nil(t, err)
		p.Reject("any")
	}
	assert.Nil(t, b.Do(func() error {
		return nil
	}))
	assert.Nil(t, b.DoWithAcceptable(func() error {
		return nil
	}, defaultAcceptable))
	errDummy := errors.New("any")
	assert.Equal(t, errDummy, b.DoWithFallback(func() error {
		return errDummy
	}, func(err error) error {
		return nil
	}))
	assert.Equal(t, errDummy, b.DoWithFallbackAcceptable(func() error {
		return errDummy
	}, func(err error) error {
		return nil
	}, defaultAcceptable))
}
