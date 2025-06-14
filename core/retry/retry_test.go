package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/retry"
)

var errRetry = errors.New("Testing")

func TestRetryNotifyRecoverNoRetries(t *testing.T) {
	config := retry.DefaultConfig()
	config.MaxRetries = 0
	config.Duration = 1

	var operationCalls, notifyCalls, recoveryCalls int

	b := config.NewBackOff()
	err := retry.NotifyRecover(func() error {
		operationCalls++

		return errRetry
	}, b, func(err error, d time.Duration, times int) {
		notifyCalls++
	}, func(ts int) {
		recoveryCalls++
	}, false)

	assert.Error(t, err)
	assert.Equal(t, errRetry, err)
	assert.Equal(t, 1, operationCalls)
	assert.Equal(t, 0, notifyCalls)
	assert.Equal(t, 0, recoveryCalls)
}

func TestRetryNotifyRecoverMaxRetries(t *testing.T) {
	config := retry.DefaultConfig()
	config.Policy = retry.PolicyConstant
	config.MaxRetries = 3
	config.Duration = 1

	var operationCalls, notifyCalls, recoveryCalls int

	b := config.NewBackOff()
	err := retry.NotifyRecover(func() error {
		operationCalls++

		return errRetry
	}, b, func(err error, d time.Duration, times int) {
		notifyCalls++
	}, func(ts int) {
		recoveryCalls++
	}, false)

	assert.Error(t, err)
	assert.Equal(t, errRetry, err)
	assert.Equal(t, 4, operationCalls)
	assert.Equal(t, 1, notifyCalls)
	assert.Equal(t, 0, recoveryCalls)
}

func TestRetryNotifyRecoverRecovery(t *testing.T) {
	config := retry.DefaultConfig()
	config.Policy = retry.PolicyConstant
	config.MaxRetries = 3
	config.Duration = 1

	var operationCalls, notifyCalls, recoveryCalls, times int

	b := config.NewBackOff()
	err := retry.NotifyRecover(func() error {
		operationCalls++

		if operationCalls >= 2 {
			return nil
		}

		return errRetry
	}, b, func(err error, d time.Duration, times int) {
		notifyCalls++
	}, func(ts int) {
		times = ts
		recoveryCalls++
	}, false)

	assert.NoError(t, err)
	assert.Equal(t, 2, operationCalls)
	assert.Equal(t, 1, notifyCalls)
	assert.Equal(t, 1, recoveryCalls)
	assert.Equal(t, 1, times)
}

func TestRetryNotifyRecoverCancel(t *testing.T) {
	config := retry.DefaultConfig()
	config.Policy = retry.PolicyConstant
	config.Duration = 1 * time.Minute

	var notifyCalls, recoveryCalls int

	ctx, cancel := context.WithCancel(context.Background())
	b := config.NewBackOffWithContext(ctx)
	errC := make(chan error, 1)
	startedC := make(chan struct{}, 100)

	go func() {
		errC <- retry.NotifyRecover(func() error {
			return errRetry
		}, b, func(err error, d time.Duration, times int) {
			notifyCalls++
			startedC <- struct{}{}
		}, func(ts int) {
			recoveryCalls++
		}, false)
	}()

	<-startedC
	cancel()

	err := <-errC
	assert.Error(t, err)
	assert.True(t, errors.Is(err, context.Canceled))
	assert.Equal(t, 1, notifyCalls)
	assert.Equal(t, 0, recoveryCalls)
}

func TestPermanent(t *testing.T) {
	config := retry.DefaultConfig()
	config.Policy = retry.PolicyConstant
	config.MaxRetries = 3
	config.Duration = 1

	var operationCalls, notifyCalls, recoveryCalls, times int

	b := config.NewBackOff()
	err := retry.NotifyRecover(func() error {
		operationCalls++

		if operationCalls >= 2 {
			return retry.Permanent(errRetry)
		}

		return errRetry
	}, b, func(err error, d time.Duration, times int) {
		notifyCalls++
	}, func(ts int) {
		times = ts
		recoveryCalls++
	}, false)

	assert.ErrorIs(t, err, errRetry)
	assert.Equal(t, 2, operationCalls)
	assert.Equal(t, 1, notifyCalls)
	assert.Equal(t, 0, recoveryCalls)
	assert.Equal(t, 0, times)
}
