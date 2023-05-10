package fx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	assert.NotNil(t, DoWithRetry(func() error {
		return errors.New("any")
	}))

	var times int
	assert.Nil(t, DoWithRetry(func() error {
		times++
		if times == defaultRetryTimes {
			return nil
		}
		return errors.New("any")
	}))

	times = 0
	assert.NotNil(t, DoWithRetry(func() error {
		times++
		if times == defaultRetryTimes+1 {
			return nil
		}
		return errors.New("any")
	}))

	total := 2 * defaultRetryTimes
	times = 0
	assert.Nil(t, DoWithRetry(func() error {
		times++
		if times == total {
			return nil
		}
		return errors.New("any")
	}, WithRetry(total)))
}

func TestErrHandler(t *testing.T) {
	total := 2 * defaultRetryTimes
	breakErr := errors.New("break here")
	times := 0

	err := DoWithRetry(func() error {
		times++
		if times > 1 {
			return breakErr
		}
		return errors.New("common err")
	},
		WithRetry(total),
		WithErrHandler(func(err error) bool {
			return errors.Is(err, breakErr)
		}),
	)

	assert.Equal(t, err.Error(), "common err\nbreak here")

	times = 0
	err = DoWithRetry(func() error {
		return errors.New("common err")
	}, WithRetry(2))

	assert.Equal(t, err.Error(), "common err\ncommon err")

	times = 0
	err = DoWithRetry(func() error {
		times++
		if times > 1 {
			return breakErr
		}
		return nil
	},
		WithRetry(total),
		WithErrHandler(func(err error) bool {
			return errors.Is(err, breakErr)
		}),
	)

	assert.Nil(t, err)
}
