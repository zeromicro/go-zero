package breaker

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/stat"
)

func init() {
	stat.SetReporter(nil)
}

func TestCircuitBreaker_Allow(t *testing.T) {
	b := NewBreaker()
	assert.True(t, len(b.Name()) > 0)
	_, err := b.Allow()
	assert.Nil(t, err)
}

func TestLogReason(t *testing.T) {
	b := NewBreaker()
	assert.True(t, len(b.Name()) > 0)

	for i := 0; i < 1000; i++ {
		_ = b.Do(func() error {
			return errors.New(strconv.Itoa(i))
		})
	}
	errs := b.(*circuitBreaker).throttle.(loggedThrottle).errWin
	assert.Equal(t, numHistoryReasons, errs.count)
}

func BenchmarkGoogleBreaker(b *testing.B) {
	br := NewBreaker()
	for i := 0; i < b.N; i++ {
		_ = br.Do(func() error {
			return nil
		})
	}
}
