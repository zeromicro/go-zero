package breaker

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNopBreaker(t *testing.T) {
	b := NopBreaker()
	assert.Equal(t, nopBreakerName, b.Name())
	_, err := b.Allow()
	assert.Nil(t, err)
	p, err := b.AllowCtx(context.Background())
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
	assert.Nil(t, b.DoCtx(context.Background(), func() error {
		return nil
	}))
	assert.Nil(t, b.DoWithAcceptable(func() error {
		return nil
	}, defaultAcceptable))
	assert.Nil(t, b.DoWithAcceptableCtx(context.Background(), func() error {
		return nil
	}, defaultAcceptable))
	errDummy := errors.New("any")
	assert.Equal(t, errDummy, b.DoWithFallback(func() error {
		return errDummy
	}, func(err error) error {
		return nil
	}))
	assert.Equal(t, errDummy, b.DoWithFallbackCtx(context.Background(), func() error {
		return errDummy
	}, func(err error) error {
		return nil
	}))
	assert.Equal(t, errDummy, b.DoWithFallbackAcceptable(func() error {
		return errDummy
	}, func(err error) error {
		return nil
	}, defaultAcceptable))
	assert.Equal(t, errDummy, b.DoWithFallbackAcceptableCtx(context.Background(), func() error {
		return errDummy
	}, func(err error) error {
		return nil
	}, defaultAcceptable))
}
