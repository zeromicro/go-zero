package mon

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	mopt "go.mongodb.org/mongo-driver/mongo/options"
)

func TestSetSlowThreshold(t *testing.T) {
	assert.Equal(t, defaultSlowThreshold, slowThreshold.Load())
	SetSlowThreshold(time.Second)
	assert.Equal(t, time.Second, slowThreshold.Load())
}

func Test_defaultTimeoutOption(t *testing.T) {
	opts := mopt.Client()
	defaultTimeoutOption()(opts)
	assert.Equal(t, defaultTimeout, *opts.Timeout)
}

func TestWithTimeout(t *testing.T) {
	opts := mopt.Client()
	WithTimeout(time.Second)(opts)
	assert.Equal(t, time.Second, *opts.Timeout)
}

func TestDisableLog(t *testing.T) {
	assert.True(t, logMon.True())
	assert.True(t, logSlowMon.True())
	defer func() {
		logMon.Set(true)
		logSlowMon.Set(true)
	}()

	DisableLog()
	assert.False(t, logMon.True())
	assert.False(t, logSlowMon.True())
}

func TestDisableInfoLog(t *testing.T) {
	assert.True(t, logMon.True())
	assert.True(t, logSlowMon.True())
	defer func() {
		logMon.Set(true)
		logSlowMon.Set(true)
	}()

	DisableInfoLog()
	assert.False(t, logMon.True())
	assert.True(t, logSlowMon.True())
}
